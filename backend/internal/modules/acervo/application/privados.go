package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Arquivos privados de interesse público e social (arts. 55 a 58; ADR
// 036). Só o Arquivo Público (acervo:central, global) lê e escreve: o
// cadastro tem dados pessoais do titular (LGPD).

// Eventos dos arquivos privados (outbox e histórico).
const (
	EventPrivadoCadastrado = "acervo.privado.cadastrado"
	EventPrivadoAlterado   = "acervo.privado.alterado"
	EventPrivadoParecer    = "acervo.privado.parecer"
	EventPrivadoDeclarado  = "acervo.privado.declarado"
	EventPrivadoGuarda     = "acervo.privado.guarda"
	EventPrivadoAlienacao  = "acervo.privado.alienacao"
)

// registrarPrivado grava o passo no histórico, o evento e a auditoria (sem
// os dados pessoais no evento).
func (s *Service) registrarPrivado(ctx context.Context, tx pgx.Tx, identity auth.Identity, evento string, p domain.Privado, descricao string, antes, depois any) error {
	tipo := strings.TrimPrefix(evento, "acervo.privado.")
	if err := s.repo.InsertEventoPrivado(ctx, tx, p.ID, domain.Evento{Tipo: tipo, Descricao: descricao, Por: ator(identity)}); err != nil {
		return err
	}
	if err := s.outbox.Write(ctx, tx, evento, "acervo_privado", p.ID.String(), uuid.Nil,
		map[string]any{"id": p.ID, "numero": p.Numero, "situacao": p.Situacao, "guarda": p.Guarda}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_privado", p.ID.String(), antes, depois))
}

// CadastrarPrivado registra o arquivo privado. A entidade privada de
// serviço público municipal já é de interesse público e social (§ 3º).
func (s *Service) CadastrarPrivado(ctx context.Context, identity auth.Identity, dados domain.Privado) (domain.Privado, error) {
	if err := central(identity); err != nil {
		return domain.Privado{}, MapError(err)
	}
	p, err := domain.NovoPrivado(dados)
	if err != nil {
		return domain.Privado{}, MapError(err)
	}
	p.CriadoEm = s.agora().UTC()
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if p.Numero, err = s.repo.InsertPrivado(ctx, tx, p, ator(identity)); err != nil {
			return err
		}
		desc := "Cadastrado; aguarda o parecer da Comissão Especial (art. 56, § 1º)"
		if p.ServicoPublico {
			desc = "Cadastrado e de interesse público e social por força do art. 56, § 3º (entidade privada de serviço público)"
		}
		return s.registrarPrivado(ctx, tx, identity, EventPrivadoCadastrado, p, desc, nil, map[string]any{"numero": p.Numero, "situacao": p.Situacao})
	})
	if err != nil {
		return domain.Privado{}, MapError(err)
	}
	return s.Privado(ctx, identity, p.ID)
}

// etapaPrivado carrega o arquivo para alterar, aplica a mudança e registra.
func (s *Service) etapaPrivado(ctx context.Context, identity auth.Identity, id uuid.UUID, evento string,
	fn func(p *domain.Privado) (string, error)) (domain.Privado, error) {
	if err := central(identity); err != nil {
		return domain.Privado{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		p, err := s.repo.GetPrivado(ctx, tx, id, true)
		if err != nil {
			return err
		}
		antes := map[string]any{"situacao": p.Situacao, "guarda": p.Guarda}
		descricao, err := fn(&p)
		if err != nil {
			return err
		}
		if err := s.repo.UpdatePrivado(ctx, tx, p); err != nil {
			return err
		}
		return s.registrarPrivado(ctx, tx, identity, evento, p, descricao, antes, map[string]any{"situacao": p.Situacao, "guarda": p.Guarda})
	})
	if err != nil {
		return domain.Privado{}, MapError(err)
	}
	return s.Privado(ctx, identity, id)
}

// AlterarPrivado corrige o cadastro (titular, descrição, acesso, convênio).
// A etapa, a guarda e o § 3º não mudam por aqui.
func (s *Service) AlterarPrivado(ctx context.Context, identity auth.Identity, id uuid.UUID, dados domain.Privado) (domain.Privado, error) {
	return s.etapaPrivado(ctx, identity, id, EventPrivadoAlterado, func(p *domain.Privado) (string, error) {
		dados.TipoTitular, dados.ServicoPublico = p.TipoTitular, p.ServicoPublico
		if err := dados.ValidarCadastro(); err != nil {
			return "", err
		}
		p.Titular, p.Documento, p.Contato, p.Descricao = dados.Titular, dados.Documento, dados.Contato, dados.Descricao
		p.DatasLimite, p.Mensuracao, p.LocalGuarda, p.Convenio = dados.DatasLimite, dados.Mensuracao, dados.LocalGuarda, dados.Convenio
		p.AcessoAutorizado, p.AcessoCondicoes = dados.AcessoAutorizado, dados.AcessoCondicoes
		return "Cadastro alterado", nil
	})
}

// Parecer é o da Comissão Especial (art. 56, § 1º).
type Parecer struct {
	Comissao  string
	Texto     string
	Favoravel bool
	Em        time.Time
}

// RegistrarParecer grava o parecer: favorável, segue ao decreto;
// contrário, encerra a declaração.
func (s *Service) RegistrarParecer(ctx context.Context, identity auth.Identity, id uuid.UUID, pa Parecer) (domain.Privado, error) {
	return s.etapaPrivado(ctx, identity, id, EventPrivadoParecer, func(p *domain.Privado) (string, error) {
		if err := domain.PodeParecer(*p, pa.Comissao, pa.Texto, pa.Em, s.hoje()); err != nil {
			return "", err
		}
		em, fav := pa.Em, pa.Favoravel
		p.ParecerComissao, p.ParecerTexto, p.ParecerFavoravel, p.ParecerEm = strings.TrimSpace(pa.Comissao), strings.TrimSpace(pa.Texto), &fav, &em
		if fav {
			p.Situacao = domain.PrivadoParecer
			return "Parecer favorável da Comissão Especial; aguarda o decreto", nil
		}
		p.Situacao = domain.PrivadoNaoDeclar
		return "Parecer contrário da Comissão Especial", nil
	})
}

// Declarar registra o decreto de interesse público e social (art. 56).
func (s *Service) Declarar(ctx context.Context, identity auth.Identity, id uuid.UUID, numero, diario string, em time.Time) (domain.Privado, error) {
	return s.etapaPrivado(ctx, identity, id, EventPrivadoDeclarado, func(p *domain.Privado) (string, error) {
		if err := domain.PodeDeclarar(*p, numero, diario, em, s.hoje()); err != nil {
			return "", err
		}
		p.Situacao, p.DecretoNumero, p.DecretoEm, p.DecretoDiario = domain.PrivadoDeclarado, strings.TrimSpace(numero), &em, strings.TrimSpace(diario)
		return "Declarado de interesse público e social pelo Decreto " + p.DecretoNumero, nil
	})
}

var descGuarda = map[string]string{
	domain.GuardaDoado:      "Doado ao Arquivo Público (art. 56, § 5º)",
	domain.GuardaDepositado: "Depositado no Arquivo Público, a título revogável (art. 56, § 5º)",
	domain.GuardaTitular:    "Depósito revogado: o acervo voltou ao titular (art. 56, § 4º)",
}

// MudarGuarda registra a doação, o depósito revogável ou a revogação do
// depósito.
func (s *Service) MudarGuarda(ctx context.Context, identity auth.Identity, id uuid.UUID, guarda, termo string, em time.Time) (domain.Privado, error) {
	guarda = strings.ToUpper(strings.TrimSpace(guarda))
	return s.etapaPrivado(ctx, identity, id, EventPrivadoGuarda, func(p *domain.Privado) (string, error) {
		if err := domain.PodeMudarGuarda(*p, guarda, termo, em, s.hoje()); err != nil {
			return "", err
		}
		p.Guarda, p.GuardaEm, p.GuardaTermo = guarda, &em, strings.TrimSpace(termo)
		return descGuarda[guarda] + ": " + p.GuardaTermo, nil
	})
}

// NotificarAlienacao registra a notificação do titular que quer vender: o
// município tem 60 dias para manifestar interesse (art. 58).
func (s *Service) NotificarAlienacao(ctx context.Context, identity auth.Identity, id uuid.UUID, a domain.Alienacao) (domain.Privado, error) {
	a.ID, a.Decisao = uuid.New(), domain.AlienacaoPendente
	a.Adquirente, a.Condicoes = strings.TrimSpace(a.Adquirente), strings.TrimSpace(a.Condicoes)
	return s.alienacao(ctx, identity, id, func(tx pgx.Tx, p *domain.Privado) (string, error) {
		if err := domain.PodeNotificarAlienacao(*p, a, s.hoje()); err != nil {
			return "", err
		}
		a.Completar(s.hoje())
		if err := s.repo.InsertAlienacao(ctx, tx, p.ID, a, ator(identity)); err != nil {
			return "", err
		}
		return fmt.Sprintf("Notificação de alienação recebida em %s; o município manifesta interesse até %s (art. 58)",
			a.NotificadoEm.Format("02/01/2006"), a.PrazoFinal.Format("02/01/2006")), nil
	})
}

// DecidirAlienacao registra a manifestação do município dentro do prazo.
func (s *Service) DecidirAlienacao(ctx context.Context, identity auth.Identity, id, alienacaoID uuid.UUID, decisao, fundamento string) (domain.Privado, error) {
	decisao = strings.ToUpper(strings.TrimSpace(decisao))
	return s.alienacao(ctx, identity, id, func(tx pgx.Tx, p *domain.Privado) (string, error) {
		var a *domain.Alienacao
		for i := range p.Alienacoes {
			if p.Alienacoes[i].ID == alienacaoID {
				a = &p.Alienacoes[i]
			}
		}
		if a == nil {
			return "", domain.ErrAlienacaoNaoExiste
		}
		if err := domain.PodeDecidirAlienacao(*a, decisao, fundamento, s.hoje()); err != nil {
			return "", err
		}
		agora := s.agora().UTC()
		a.Decisao, a.Fundamento, a.DecididoEm = decisao, strings.TrimSpace(fundamento), &agora
		if err := s.repo.DecidirAlienacao(ctx, tx, *a, ator(identity)); err != nil {
			return "", err
		}
		if decisao == domain.AlienacaoInteresse {
			return "O município manifestou interesse na aquisição (art. 58)", nil
		}
		return "O município declarou não ter interesse na aquisição (art. 58)", nil
	})
}

// alienacao: as operações do art. 58 gravam a alienação e o histórico, sem
// mudar o cadastro.
func (s *Service) alienacao(ctx context.Context, identity auth.Identity, id uuid.UUID, fn func(tx pgx.Tx, p *domain.Privado) (string, error)) (domain.Privado, error) {
	if err := central(identity); err != nil {
		return domain.Privado{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		p, err := s.repo.GetPrivado(ctx, tx, id, true)
		if err != nil {
			return err
		}
		descricao, err := fn(tx, &p)
		if err != nil {
			return err
		}
		return s.registrarPrivado(ctx, tx, identity, EventPrivadoAlienacao, p, descricao, nil, map[string]any{"descricao": descricao})
	})
	if err != nil {
		return domain.Privado{}, MapError(err)
	}
	return s.Privado(ctx, identity, id)
}

// Privado devolve o arquivo com o histórico e as alienações (prazo do art.
// 58 calculado na hora).
func (s *Service) Privado(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Privado, error) {
	if err := central(identity); err != nil {
		return domain.Privado{}, MapError(err)
	}
	p, err := s.repo.GetPrivado(ctx, s.pool, id, false)
	for i := range p.Alienacoes {
		p.Alienacoes[i].Completar(s.hoje())
	}
	return p, MapError(err)
}

// Privados pagina os arquivos privados.
func (s *Service) Privados(ctx context.Context, identity auth.Identity, q string, p pagination.Params) ([]domain.Privado, int64, error) {
	if err := central(identity); err != nil {
		return nil, 0, MapError(err)
	}
	items, total, err := s.repo.ListPrivados(ctx, s.pool, buscaNorma(q), p)
	return items, total, MapError(err)
}
