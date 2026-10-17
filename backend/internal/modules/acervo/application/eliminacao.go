package application

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/notificacoes"
)

// Eventos da Fase 3 (ADR 034).
const (
	EventComissaoSalva      = "acervo.comissao.salva"
	EventListagemCriada     = "acervo.listagem.criada"
	EventListagemAlterada   = "acervo.listagem.alterada"
	EventListagemAprovada   = "acervo.listagem.aprovada"
	EventListagemAutorizada = "acervo.listagem.autorizada"
	EventListagemCCPAD      = "acervo.listagem.ccpad"
	EventListagemSEMAD      = "acervo.listagem.semad"
	EventListagemDevolvida  = "acervo.listagem.devolvida"
	EventEditalPublicado    = "acervo.listagem.edital_publicado"
	EventPedidoRegistrado   = "acervo.listagem.pedido_registrado"
	EventPedidoDecidido     = "acervo.listagem.pedido_decidido"
	EventAmostragem         = "acervo.listagem.amostragem"
	// EventEliminacao: o Termo de Eliminação (Anexo II). Uma cópia vai à
	// Superintendência de Arquivo Público e à SEMAD (art. 38, parágrafo único).
	EventEliminacao        = "acervo.listagem.eliminada"
	EventListagemCancelada = "acervo.listagem.cancelada"
)

// Comissoes lista as comissões (CSPADs e CCPAD) com os membros.
func (s *Service) Comissoes(ctx context.Context) ([]domain.Comissao, error) {
	c, err := s.repo.ListComissoes(ctx, s.pool)
	return c, MapError(err)
}

// SalvarComissao cria (id nil) ou altera uma comissão. Só o Arquivo Público
// (acervo:central) mantém o cadastro das comissões designadas por portaria.
func (s *Service) SalvarComissao(ctx context.Context, identity auth.Identity, id *uuid.UUID, c domain.Comissao) (domain.Comissao, error) {
	if !auth.CanGlobal(identity, auth.PermAcervoCentral) {
		return domain.Comissao{}, MapError(domain.ErrForaDoEscopo)
	}
	nova := id == nil
	if nova {
		c.ID = uuid.New()
	} else {
		c.ID = *id
	}
	if err := c.Validar(); err != nil {
		return domain.Comissao{}, MapError(err)
	}
	avisos := c.Avisos
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		var antes any
		if !nova {
			atual, err := s.repo.GetComissao(ctx, tx, c.ID)
			if err != nil {
				return err
			}
			antes = atual
		}
		if err := s.repo.SalvarComissao(ctx, tx, c, nova); err != nil {
			return err
		}
		if err := s.outbox.Write(ctx, tx, EventComissaoSalva, "acervo_comissao", c.ID.String(), uuid.Nil,
			map[string]any{"id": c.ID, "tipo": c.Tipo, "nome": c.Nome}); err != nil {
			return err
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventComissaoSalva, "acervo_comissao", c.ID.String(), antes, c))
	})
	if err != nil {
		return domain.Comissao{}, MapError(err)
	}
	out, err := s.repo.GetComissao(ctx, s.pool, c.ID)
	out.Avisos = avisos
	return out, MapError(err)
}

func avalia(identity auth.Identity, unidade uuid.UUID) bool {
	return auth.Can(identity, auth.PermAcervoAvaliacao, auth.InUnidade(unidade))
}

// NovaListagem é o cabeçalho da Relação de Eliminação (Anexo I).
type NovaListagem struct {
	UnidadeID          uuid.UUID
	Rito               string
	ResponsavelSelecao *uuid.UUID
	ParecerJuridico    string
	ContasExercicios   string
	ContasAprovadasEm  *time.Time
	Unidades           []uuid.UUID
}

func limparCabecalho(l *domain.Listagem) error {
	l.ParecerJuridico, l.ContasExercicios = strings.TrimSpace(l.ParecerJuridico), strings.TrimSpace(l.ContasExercicios)
	switch {
	case len([]rune(l.ParecerJuridico)) > 4000:
		return &domain.ErrValidacao{Campo: "parecer_juridico", Motivo: "até 4000 caracteres"}
	case len([]rune(l.ContasExercicios)) > 200:
		return &domain.ErrValidacao{Campo: "contas_exercicios", Motivo: "até 200 caracteres"}
	}
	return nil
}

// conferirItem diz se a unidade pode ir (ou continuar) na relação: da
// unidade da relação, elegível pelas travas da lei (PodeEliminar — arts. 19,
// 20, 30, 42 e 44) e fora de empréstimo. Na inclusão, também fora de outra
// relação em andamento.
func (s *Service) conferirItem(ctx context.Context, tx pgx.Tx, l domain.Listagem, id uuid.UUID, incluir bool) (domain.Unidade, domain.Serie, error) {
	u, _, err := s.repo.GetUnidade(ctx, tx, id, true)
	if err != nil {
		return u, domain.Serie{}, err
	}
	if u.UnidadeID != l.UnidadeID {
		return u, domain.Serie{}, &domain.ErrValidacao{Campo: "unidades", Motivo: u.Numero + " não é da unidade da relação"}
	}
	serie, err := s.ttdd.Serie(ctx, u.CodigoClassificacao)
	existe := err == nil
	if err != nil && !(l.Rito == domain.RitoArt35 && errors.Is(err, domain.ErrSerieInvalida)) {
		return u, serie, err
	}
	bloqueada, err := s.repo.TemBloqueioAtivo(ctx, tx, u.ID)
	if err != nil {
		return u, serie, err
	}
	if l.Rito == domain.RitoArt35 {
		if err := domain.PodeEliminarArt35(u, serie, existe, bloqueada, s.hoje()); err != nil {
			return u, serie, prefixar(u, err)
		}
	} else if err := domain.PodeEliminar(u, serie, bloqueada, s.hoje()); err != nil {
		return u, serie, prefixar(u, err)
	}
	fora, err := s.repo.EmprestimoAberto(ctx, tx, u.ID)
	if err != nil {
		return u, serie, err
	}
	if fora {
		return u, serie, prefixar(u, domain.ErrEmprestada)
	}
	if incluir {
		em, err := s.repo.EmListagemAtiva(ctx, tx, u.ID)
		if err != nil {
			return u, serie, err
		}
		if em {
			return u, serie, prefixar(u, domain.ErrEmOutraListagem)
		}
	}
	return u, serie, nil
}

// prefixar diz qual unidade travou a etapa, sem perder o erro de origem.
type erroUnidade struct {
	numero string
	err    error
}

func (e *erroUnidade) Error() string {
	return e.numero + ": " + strings.TrimPrefix(e.err.Error(), "acervo: ")
}
func (e *erroUnidade) Unwrap() error { return e.err }

func prefixar(u domain.Unidade, err error) error { return &erroUnidade{numero: u.Numero, err: err} }

// registrarListagem grava o evento e a auditoria de uma etapa.
func (s *Service) registrarListagem(ctx context.Context, tx pgx.Tx, evento string, l domain.Listagem, antes, depois any) error {
	if err := s.outbox.Write(ctx, tx, evento, "acervo_listagem", l.ID.String(), uuid.Nil,
		map[string]any{"id": l.ID, "numero": l.Numero, "unidade_id": l.UnidadeID, "situacao": l.Situacao}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_listagem", l.ID.String(), antes, depois))
}

// movimentar registra o passo no histórico de cada unidade dada.
func (s *Service) movimentar(ctx context.Context, tx pgx.Tx, identity auth.Identity, itens []domain.ItemListagem, tipo, desc string) error {
	for _, i := range itens {
		if err := s.repo.InsertMovimentacao(ctx, tx, i.ID, domain.Movimentacao{Tipo: tipo, Por: ator(identity), Descricao: desc}); err != nil {
			return err
		}
	}
	return nil
}

// CriarListagem abre a Relação de Eliminação da unidade com a CSPAD do órgão
// (art. 31) e inclui as unidades elegíveis.
func (s *Service) CriarListagem(ctx context.Context, identity auth.Identity, n NovaListagem) (domain.Listagem, error) {
	if !avalia(identity, n.UnidadeID) {
		return domain.Listagem{}, MapError(domain.ErrForaDoEscopo)
	}
	rito, err := domain.NormalizarRito(n.Rito)
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	l := domain.Listagem{ID: uuid.New(), UnidadeID: n.UnidadeID, Situacao: domain.ListagemRascunho, Rito: rito,
		ResponsavelSelecao: n.ResponsavelSelecao, ParecerJuridico: n.ParecerJuridico, ContasExercicios: n.ContasExercicios,
		ContasAprovadasEm: n.ContasAprovadasEm, AtaSignum: domain.SignumAusente, TermoSignum: domain.SignumAusente,
		CriadaPor: ator(identity), CriadaEm: s.agora().UTC()}
	if err := limparCabecalho(&l); err != nil {
		return domain.Listagem{}, MapError(err)
	}
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		var c domain.Comissao
		var err error
		if rito == domain.RitoArt36 {
			c, err = s.repo.ComissaoCentral(ctx, tx)
		} else {
			c, err = s.repo.ComissaoDaUnidade(ctx, tx, n.UnidadeID)
		}
		if err != nil {
			return err
		}
		l.ComissaoID = c.ID
		if l.Numero, err = s.repo.InsertListagem(ctx, tx, l); err != nil {
			return err
		}
		if err := s.incluir(ctx, tx, identity, l, n.Unidades); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemCriada, l, nil, map[string]any{"numero": l.Numero, "unidades": n.Unidades})
	})
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	return s.Listagem(ctx, identity, l.ID)
}

func (s *Service) incluir(ctx context.Context, tx pgx.Tx, identity auth.Identity, l domain.Listagem, ids []uuid.UUID) error {
	for _, id := range ids {
		u, _, err := s.conferirItem(ctx, tx, l, id, true)
		if err != nil {
			return err
		}
		if err := s.repo.AddItemListagem(ctx, tx, l.ID, u.ID); err != nil {
			return err
		}
		if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovListagem, Por: ator(identity),
			Descricao: "Incluída na Relação de Eliminação " + l.Numero}); err != nil {
			return err
		}
	}
	return nil
}

// carregarListagem lê a relação para alterá-la, no escopo da CSPAD (ou do
// Arquivo Público, com central).
func (s *Service) carregarListagem(ctx context.Context, tx pgx.Tx, identity auth.Identity, id uuid.UUID, central bool) (domain.Listagem, error) {
	l, err := s.repo.GetListagem(ctx, tx, id, true)
	if err != nil {
		return l, err
	}
	if !podeLerListagem(identity, l) {
		return l, domain.ErrListagemNaoExiste
	}
	if central && !auth.CanGlobal(identity, auth.PermAcervoCentral) || !central && !avalia(identity, l.UnidadeID) {
		return l, domain.ErrForaDoEscopo
	}
	return l, nil
}

func podeLerListagem(identity auth.Identity, l domain.Listagem) bool {
	return le(identity, l.UnidadeID) || avalia(identity, l.UnidadeID) || auth.CanGlobal(identity, auth.PermAcervoCentral) ||
		auth.CanGlobal(identity, auth.PermAcervoSEMAD)
}

// Etapa aplica uma mudança à relação na transação e devolve a relação lida
// de novo.
func (s *Service) etapa(ctx context.Context, identity auth.Identity, id uuid.UUID, central bool,
	fn func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error) (domain.Listagem, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		l, err := s.carregarListagem(ctx, tx, identity, id, central)
		if err != nil {
			return err
		}
		return fn(ctx, tx, &l)
	})
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	return s.Listagem(ctx, identity, id)
}

// AlterarListagem muda o cabeçalho e inclui ou retira unidades (rascunho).
func (s *Service) AlterarListagem(ctx context.Context, identity auth.Identity, id uuid.UUID, n NovaListagem, retirar []uuid.UUID) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if l.Situacao != domain.ListagemRascunho {
			return domain.ErrListagemSituacao
		}
		antes := *l
		l.ResponsavelSelecao, l.ParecerJuridico, l.ContasExercicios, l.ContasAprovadasEm = n.ResponsavelSelecao, n.ParecerJuridico,
			n.ContasExercicios, n.ContasAprovadasEm
		if err := limparCabecalho(l); err != nil {
			return err
		}
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		for _, u := range retirar {
			if err := s.repo.RemoveItemListagem(ctx, tx, l.ID, u); err != nil {
				return err
			}
		}
		if err := s.incluir(ctx, tx, identity, *l, n.Unidades); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemAlterada, *l, cabecalho(antes),
			map[string]any{"cabecalho": cabecalho(*l), "incluidas": n.Unidades, "retiradas": retirar})
	})
}

func cabecalho(l domain.Listagem) map[string]any {
	return map[string]any{"parecer_juridico": l.ParecerJuridico, "contas_exercicios": l.ContasExercicios,
		"contas_aprovadas_em": l.ContasAprovadasEm, "responsavel_selecao": l.ResponsavelSelecao}
}

// Ata é a reunião da CSPAD que aprova a relação (Anexo IV).
type Ata struct {
	Numero        string
	Em            time.Time
	Local         string
	Manifestacoes string
	Presentes     []uuid.UUID
}

// Aprovar registra a ata da CSPAD (Anexo IV) e confere de novo cada unidade.
func (s *Service) Aprovar(ctx context.Context, identity auth.Identity, id uuid.UUID, a Ata) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		c, err := s.repo.GetComissao(ctx, tx, l.ComissaoID)
		if err != nil {
			return err
		}
		if err := domain.PodeAprovar(*l, c, a.Numero, a.Em, a.Presentes); err != nil {
			return err
		}
		if err := s.reconferir(ctx, tx, *l); err != nil {
			return err
		}
		em := a.Em.UTC()
		l.AtaNumero, l.AtaEm, l.AtaLocal, l.AtaManifestacoes, l.AtaPresentes = strings.TrimSpace(a.Numero), &em,
			strings.TrimSpace(a.Local), strings.TrimSpace(a.Manifestacoes), a.Presentes
		l.Situacao, l.DevolucaoMotivo = domain.ListagemAprovada, ""
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		quem := "CSPAD"
		if c.Tipo == domain.ComissaoCentral {
			quem = "CCPAD"
		}
		if err := s.movimentar(ctx, tx, identity, l.AEliminar(), domain.MovListagem,
			"Relação "+l.Numero+" aprovada pela "+quem+" (ata "+l.AtaNumero+")"); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemAprovada, *l, map[string]any{"situacao": domain.ListagemRascunho},
			map[string]any{"situacao": l.Situacao, "ata": l.AtaNumero, "presentes": l.AtaPresentes})
	})
}

// reconferir aplica de novo as travas a cada unidade que vai à eliminação.
func (s *Service) reconferir(ctx context.Context, tx pgx.Tx, l domain.Listagem) error {
	for _, i := range l.AEliminar() {
		if _, _, err := s.conferirItem(ctx, tx, l, i.ID, false); err != nil {
			return err
		}
	}
	return nil
}

// Autorizar é o "Autorizo" do titular da Superintendência de Arquivo Público (Anexo
// I; arts. 19 e 65). Só com acervo:central.
func (s *Service) Autorizar(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, true, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if l.Situacao != domain.ListagemAprovada {
			return domain.ErrListagemSituacao
		}
		if l.Rito == domain.RitoArt35 && l.CcpadAutorizadaEm == nil {
			return domain.ErrSemCCPAD
		}
		if err := s.reconferir(ctx, tx, *l); err != nil {
			return err
		}
		agora := s.agora().UTC()
		l.AutorizadaPor, l.AutorizadaEm = ator(identity), &agora
		if l.RitoEspecial() {
			l.Situacao = domain.ListagemAguardandoSEMAD
		} else {
			l.Situacao = domain.ListagemAutorizada
		}
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		if err := s.movimentar(ctx, tx, identity, l.AEliminar(), domain.MovListagem,
			"Eliminação autorizada pelo Arquivo Público (relação "+l.Numero+")"); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemAutorizada, *l, map[string]any{"situacao": domain.ListagemAprovada},
			map[string]any{"situacao": l.Situacao})
	})
}

// AutorizarCCPAD registra a autorização da CCPAD no rito do art. 35.
func (s *Service) AutorizarCCPAD(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		l, err := s.repo.GetListagem(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if !podeLerListagem(identity, l) {
			return domain.ErrListagemNaoExiste
		}
		c, err := s.repo.ComissaoCentral(ctx, tx)
		if err != nil {
			return err
		}
		if err := domain.PodeAutorizarCCPAD(l, c, identity.UserID); err != nil {
			return err
		}
		if err := s.reconferir(ctx, tx, l); err != nil {
			return err
		}
		agora := s.agora().UTC()
		l.CcpadAutorizadaPor, l.CcpadAutorizadaEm = ator(identity), &agora
		if err := s.repo.UpdateListagem(ctx, tx, l); err != nil {
			return err
		}
		if err := s.movimentar(ctx, tx, identity, l.AEliminar(), domain.MovListagem,
			"Eliminação autorizada pela CCPAD (relação "+l.Numero+", art. 35)"); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemCCPAD, l, map[string]any{"situacao": domain.ListagemAprovada},
			map[string]any{"ccpad_autorizada_em": agora})
	})
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	return s.Listagem(ctx, identity, id)
}

// AutorizarSEMAD registra a aprovação da Secretaria Municipal de Administração
// (arts. 35 e 36). Só com acervo:semad global.
func (s *Service) AutorizarSEMAD(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error) {
	if !auth.CanGlobal(identity, auth.PermAcervoSEMAD) {
		return domain.Listagem{}, MapError(domain.ErrForaDoEscopo)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		l, err := s.repo.GetListagem(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if err := domain.PodeAutorizarSEMAD(l); err != nil {
			return err
		}
		if err := s.reconferir(ctx, tx, l); err != nil {
			return err
		}
		agora := s.agora().UTC()
		l.Situacao, l.SemadAutorizadaPor, l.SemadAutorizadaEm = domain.ListagemAutorizada, ator(identity), &agora
		if err := s.repo.UpdateListagem(ctx, tx, l); err != nil {
			return err
		}
		if err := s.movimentar(ctx, tx, identity, l.AEliminar(), domain.MovListagem,
			"Eliminação aprovada pela SEMAD (relação "+l.Numero+")"); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemSEMAD, l, map[string]any{"situacao": domain.ListagemAguardandoSEMAD},
			map[string]any{"situacao": l.Situacao})
	})
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	return s.Listagem(ctx, identity, id)
}

// DevolverListagem: o Arquivo Público devolve a relação à CSPAD, com o
// motivo; ela volta ao rascunho e precisa de nova ata.
func (s *Service) DevolverListagem(ctx context.Context, identity auth.Identity, id uuid.UUID, motivo string) (domain.Listagem, error) {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" || len([]rune(motivo)) > 2000 {
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "motivo", Motivo: "informe o motivo da devolução (até 2000 caracteres)"})
	}
	return s.etapa(ctx, identity, id, true, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if l.Situacao != domain.ListagemAprovada && l.Situacao != domain.ListagemAguardandoSEMAD {
			return domain.ErrListagemSituacao
		}
		antes := l.Situacao
		l.Situacao, l.DevolucaoMotivo = domain.ListagemRascunho, motivo
		l.AtaNumero, l.AtaEm, l.AtaLocal, l.AtaManifestacoes, l.AtaPresentes = "", nil, "", "", []uuid.UUID{}
		l.AutorizadaPor, l.AutorizadaEm, l.CcpadAutorizadaPor, l.CcpadAutorizadaEm = nil, nil, nil, nil
		l.SemadAutorizadaPor, l.SemadAutorizadaEm = nil, nil
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemDevolvida, *l, map[string]any{"situacao": antes},
			map[string]any{"situacao": l.Situacao, "motivo": motivo})
	})
}

// Edital é a publicação do Edital de Ciência (Anexo III).
type Edital struct {
	Numero      string
	PublicadoEm time.Time
	Diario      string
}

// PublicarEdital registra a publicação no Diário Oficial do Município (art.
// 37) e abre o prazo de 30 dias.
func (s *Service) PublicarEdital(ctx context.Context, identity auth.Identity, id uuid.UUID, e Edital) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if err := domain.PodePublicar(*l, e.Numero, e.PublicadoEm, s.hoje()); err != nil {
			return err
		}
		if len([]rune(e.Diario)) > 200 {
			return &domain.ErrValidacao{Campo: "edital_diario", Motivo: "até 200 caracteres"}
		}
		pub := e.PublicadoEm
		l.Situacao, l.EditalNumero, l.EditalPublicadoEm, l.EditalDiario = domain.ListagemEmEdital, strings.TrimSpace(e.Numero), &pub,
			strings.TrimSpace(e.Diario)
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		if err := s.movimentar(ctx, tx, identity, l.AEliminar(), domain.MovListagem,
			"Edital de Ciência "+l.EditalNumero+" publicado em "+pub.Format("02/01/2006")+" (prazo de 30 dias)"); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventEditalPublicado, *l, map[string]any{"situacao": domain.ListagemAutorizada},
			map[string]any{"situacao": l.Situacao, "edital": l.EditalNumero, "publicado_em": pub, "diario": l.EditalDiario})
	})
}

// RegistrarPedido recebe o pedido de desentranhamento ou de cópia durante o
// prazo do edital (art. 37, § 2º).
func (s *Service) RegistrarPedido(ctx context.Context, identity auth.Identity, id uuid.UUID, p domain.Pedido) (domain.Listagem, error) {
	p.Requerente, p.Pedido = strings.TrimSpace(p.Requerente), strings.TrimSpace(p.Pedido)
	switch {
	case p.Requerente == "" || len([]rune(p.Requerente)) > 200:
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "requerente", Motivo: "informe quem pede (até 200 caracteres)"})
	case p.Pedido == "" || len([]rune(p.Pedido)) > 4000:
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "pedido", Motivo: "descreva o pedido (até 4000 caracteres)"})
	}
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if err := domain.PodeReceberPedido(*l, p.RecebidoEm); err != nil {
			return err
		}
		if p.UnidadeID != nil && !temItem(*l, *p.UnidadeID) {
			return &domain.ErrValidacao{Campo: "unidade_arquivamento_id", Motivo: "a unidade não está na relação"}
		}
		p.ID, p.Situacao = uuid.New(), domain.PedidoPendente
		if err := s.repo.InsertPedido(ctx, tx, l.ID, p); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventPedidoRegistrado, *l, nil, p)
	})
}

func temItem(l domain.Listagem, unidade uuid.UUID) bool {
	for _, i := range l.Itens {
		if i.ID == unidade {
			return true
		}
	}
	return false
}

// DecidirPedido defere ou indefere. Deferido com retirar, a unidade sai da
// relação (não será eliminada).
func (s *Service) DecidirPedido(ctx context.Context, identity auth.Identity, id, pedidoID uuid.UUID, deferido, retirar bool, decisao string) (domain.Listagem, error) {
	decisao = strings.TrimSpace(decisao)
	if decisao == "" || len([]rune(decisao)) > 4000 {
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "decisao", Motivo: "fundamente a decisão (até 4000 caracteres)"})
	}
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if l.Situacao != domain.ListagemEmEdital {
			return domain.ErrListagemSituacao
		}
		var p *domain.Pedido
		for i := range l.Pedidos {
			if l.Pedidos[i].ID == pedidoID {
				p = &l.Pedidos[i]
			}
		}
		if p == nil {
			return domain.ErrPedidoNaoExiste
		}
		if p.Situacao != domain.PedidoPendente {
			return domain.ErrPedidoDecidido
		}
		antes := *p
		agora := s.agora().UTC()
		p.Situacao, p.Decisao, p.DecididoPor, p.DecididoEm = domain.PedidoIndeferido, decisao, ator(identity), &agora
		if deferido {
			p.Situacao = domain.PedidoDeferido
		}
		if err := s.repo.DecidirPedido(ctx, tx, *p); err != nil {
			return err
		}
		if deferido && retirar && p.UnidadeID != nil {
			if err := s.repo.MarcarItem(ctx, tx, l.ID, *p.UnidadeID, false, true); err != nil {
				return err
			}
			if err := s.repo.InsertMovimentacao(ctx, tx, *p.UnidadeID, domain.Movimentacao{Tipo: domain.MovListagem, Por: ator(identity),
				Descricao: "Retirada da relação " + l.Numero + " por pedido deferido"}); err != nil {
				return err
			}
		}
		return s.registrarListagem(ctx, tx, EventPedidoDecidido, *l, antes, *p)
	})
}

// sortear escolhe n itens ao acaso (crypto/rand: a amostra não é escolhida
// por ninguém).
func sortear(itens []domain.ItemListagem, n int) []domain.ItemListagem {
	itens = append([]domain.ItemListagem(nil), itens...)
	for i := len(itens) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		k := int(j.Int64())
		itens[i], itens[k] = itens[k], itens[i]
	}
	return itens[:n]
}

// Amostrar faz a amostragem do art. 39: com os critérios da comissão (art.
// 32, III), sorteia o percentual dado (ao menos uma unidade) e a leva à guarda
// permanente — dali ela não sai mais (arts. 20 e 42).
func (s *Service) Amostrar(ctx context.Context, identity auth.Identity, id uuid.UUID, criterios string, percentual float64) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if err := domain.PodeAmostrar(*l, criterios, percentual); err != nil {
			return err
		}
		elegiveis := l.AEliminar()
		amostra := sortear(elegiveis, domain.TamanhoAmostra(len(elegiveis), percentual))
		for _, i := range amostra {
			u, _, err := s.repo.GetUnidade(ctx, tx, i.ID, true)
			if err != nil {
				return err
			}
			if err := s.repo.MarcarItem(ctx, tx, l.ID, u.ID, true, false); err != nil {
				return err
			}
			u.Fase, u.Destinacao = domain.FasePermanente, domain.DestinacaoGuardaPermanente
			_, p, err := s.calcular(ctx, tx, u)
			if err != nil {
				return err
			}
			if err := s.repo.UpdateUnidade(ctx, tx, u, calculo(p)); err != nil {
				return err
			}
			if err := s.repo.SetFase(ctx, tx, u.ID, u.Fase, calculo(p)); err != nil {
				return err
			}
			para := u.Fase
			if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovAmostra, FasePara: &para, Por: ator(identity),
				Descricao: "Selecionada na amostragem da relação " + l.Numero + ": guarda permanente (art. 39)"}); err != nil {
				return err
			}
		}
		agora := s.agora().UTC()
		l.AmostraCriterios, l.AmostraPercentual, l.AmostraEm = strings.TrimSpace(criterios), &percentual, &agora
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		ids := make([]uuid.UUID, len(amostra))
		for i, a := range amostra {
			ids[i] = a.ID
		}
		return s.registrarListagem(ctx, tx, EventAmostragem, *l, nil,
			map[string]any{"criterios": l.AmostraCriterios, "percentual": percentual, "amostra": ids})
	})
}

// Termo é o Termo de Eliminação (Anexo II).
type Termo struct {
	EliminadoEm     time.Time
	Acompanhamento  string
	DestinacaoPapel string
}

// RegistrarTermo registra a eliminação (Anexo II; arts. 38 e 40): depois do
// prazo do edital, sem pedido pendente e com a amostragem feita. Cada unidade
// é conferida de novo pelas travas da lei e passa a ELIMINADO.
func (s *Service) RegistrarTermo(ctx context.Context, identity auth.Identity, id uuid.UUID, t Termo) (domain.Listagem, error) {
	t.Acompanhamento, t.DestinacaoPapel = strings.TrimSpace(t.Acompanhamento), strings.TrimSpace(t.DestinacaoPapel)
	switch {
	case t.Acompanhamento == "" || len([]rune(t.Acompanhamento)) > 200:
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "termo_acompanhamento", Motivo: "informe o servidor responsável pelo acompanhamento (Anexo II)"})
	case len([]rune(t.DestinacaoPapel)) > 500:
		return domain.Listagem{}, MapError(&domain.ErrValidacao{Campo: "termo_destinacao_papel", Motivo: "até 500 caracteres"})
	}
	var enviadas []notificacoes.Enviada
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		l, err := s.carregarListagem(ctx, tx, identity, id, false)
		if err != nil {
			return err
		}
		if err := domain.PodeRegistrarTermo(l, t.EliminadoEm, s.hoje()); err != nil {
			return err
		}
		for _, i := range l.AEliminar() {
			u, _, err := s.conferirItem(ctx, tx, l, i.ID, false)
			if err != nil {
				return err
			}
			de := u.Fase
			u.Fase = domain.FaseEliminado
			_, p, err := s.calcular(ctx, tx, u)
			if err != nil {
				return err
			}
			if err := s.repo.SetFase(ctx, tx, u.ID, u.Fase, calculo(p)); err != nil {
				return err
			}
			para := u.Fase
			if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovEliminacao, FaseDe: &de, FasePara: &para,
				Por: ator(identity), Descricao: "Eliminada em " + t.EliminadoEm.Format("02/01/2006") + " (Termo da relação " + l.Numero + ")"}); err != nil {
				return err
			}
		}
		em := t.EliminadoEm
		l.Situacao, l.TermoEliminadoEm, l.TermoAcompanhamento, l.TermoDestinacaoPapel, l.TermoRegistradoPor = domain.ListagemEliminada,
			&em, t.Acompanhamento, t.DestinacaoPapel, ator(identity)
		if enviadas, err = s.avisarCopiaTermo(ctx, tx, identity, &l); err != nil {
			return err
		}
		if err := s.repo.UpdateListagem(ctx, tx, l); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventEliminacao, l, map[string]any{"situacao": domain.ListagemEmEdital},
			map[string]any{"situacao": l.Situacao, "eliminado_em": em, "acompanhamento": t.Acompanhamento,
				"eliminadas": len(l.AEliminar()), "copia_arquivo": l.TermoCopiaArquivo, "copia_semad": l.TermoCopiaSEMAD})
	})
	if err != nil {
		return domain.Listagem{}, MapError(err)
	}
	s.entregar(ctx, enviadas)
	return s.Listagem(ctx, identity, id)
}

// avisarCopiaTermo entrega a cópia do Termo à Superintendência de Arquivo e à SEMAD
// (art. 38, parágrafo único) na caixa de quem tem acervo:central e
// acervo:semad globais. Sem destinatário SEMAD, o ato fica registrado com
// zero cópias — a permissão ainda não foi concedida (Fase 0).
func (s *Service) avisarCopiaTermo(ctx context.Context, tx pgx.Tx, identity auth.Identity, l *domain.Listagem) ([]notificacoes.Enviada, error) {
	arquivo, err := s.repo.EquipeCentral(ctx, tx)
	if err != nil {
		return nil, err
	}
	semad, err := s.repo.EquipeSEMAD(ctx, tx)
	if err != nil {
		return nil, err
	}
	agora := s.agora().UTC()
	l.TermoCopiaEm, l.TermoCopiaArquivo, l.TermoCopiaSEMAD = &agora, len(arquivo), len(semad)
	titulo := "Cópia do Termo de Eliminação " + l.Numero
	msg := "Art. 38, parágrafo único: cópia do Termo da relação " + l.Numero + "."
	link := "/acervo/eliminacao/" + l.ID.String()
	vistos := map[uuid.UUID]bool{identity.UserID: true}
	var dest []uuid.UUID
	for _, u := range append(append([]uuid.UUID{}, arquivo...), semad...) {
		if !vistos[u] {
			vistos[u] = true
			dest = append(dest, u)
		}
	}
	return notificacoes.Registrar(ctx, tx, dest, notificacoes.Nova{Modulo: "acervo", Titulo: titulo, Mensagem: msg,
		Link: link, Chave: "acervo:termo:" + l.ID.String()})
}

// CancelarListagem encerra a relação sem eliminar nada (antes da amostragem).
func (s *Service) CancelarListagem(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error) {
	return s.etapa(ctx, identity, id, false, func(ctx context.Context, tx pgx.Tx, l *domain.Listagem) error {
		if err := domain.PodeCancelar(*l); err != nil {
			return err
		}
		antes := l.Situacao
		l.Situacao = domain.ListagemCancelada
		if err := s.repo.UpdateListagem(ctx, tx, *l); err != nil {
			return err
		}
		return s.registrarListagem(ctx, tx, EventListagemCancelada, *l, map[string]any{"situacao": antes}, map[string]any{"situacao": l.Situacao})
	})
}

// Listagem devolve a relação. Fora do escopo, 404.
func (s *Service) Listagem(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error) {
	l, err := s.repo.GetListagem(ctx, s.pool, id, false)
	if err == nil && !podeLerListagem(identity, l) {
		err = domain.ErrListagemNaoExiste
	}
	return l, MapError(err)
}

// Listagens pagina as relações no escopo (o Arquivo Público vê todas).
func (s *Service) Listagens(ctx context.Context, identity auth.Identity, p pagination.Params) ([]domain.Listagem, int64, error) {
	f := s.filtro(identity, domain.Filtro{})
	aval := auth.CoverageOf(identity, auth.PermAcervoAvaliacao)
	f.Unidades = append(f.Unidades, aval.Unidades...)
	f.Restrito = f.Restrito && !aval.All && !auth.CanGlobal(identity, auth.PermAcervoCentral) &&
		!auth.CanGlobal(identity, auth.PermAcervoSEMAD)
	items, total, err := s.repo.ListListagens(ctx, s.pool, f, p)
	return items, total, MapError(err)
}
