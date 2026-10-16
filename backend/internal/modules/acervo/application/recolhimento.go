package application

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/notificacoes"
)

// Fase 4 (ADR 035): recolhimento ao Arquivo Público. A origem relaciona o
// conjunto de guarda permanente, declara o preparo do art. 53 e envia; o
// Arquivo Público confere, registra a entrada e a origem recebe a cópia do
// registro (art. 5º, VII) — ou devolve com o que falta.

// Eventos do recolhimento (outbox).
const (
	EventRecolhimentoEnviado   = "acervo.recolhimento.enviado"
	EventRecolhimentoRecebido  = "acervo.recolhimento.recebido"
	EventRecolhimentoDevolvido = "acervo.recolhimento.devolvido"
	EventRecolhimentoCancelado = "acervo.recolhimento.cancelado"
	EventRecolhida             = "acervo.unidade.recolhida"
)

// WithNotificacoes liga a entrega em tempo real dos avisos à origem.
func (s *Service) WithNotificacoes(n *notificacoes.Service) *Service {
	s.avisos = n
	return s
}

// gerePara: quem gere a unidade. Recolhida, a guarda é do Arquivo Público
// (acervo:central, global); antes, da unidade produtora.
func gerePara(identity auth.Identity, u domain.Unidade) bool {
	if u.RecolhidaEm != nil {
		return auth.CanGlobal(identity, auth.PermAcervoCentral)
	}
	return gere(identity, u.UnidadeID)
}

// conferirRecolhimento diz se a unidade pode ir (ou continuar) no
// recolhimento: da origem, pelas condições do domínio, fora de empréstimo e
// de guia; na inclusão, fora de outro recolhimento enviado.
func (s *Service) conferirRecolhimento(ctx context.Context, db database.DBTX, r domain.Recolhimento, u domain.Unidade, incluir bool) error {
	if u.UnidadeID != r.OrigemID {
		return &domain.ErrValidacao{Campo: "unidades", Motivo: u.Numero + " não é da unidade de origem"}
	}
	serie, err := s.ttdd.Serie(ctx, u.CodigoClassificacao)
	if err != nil {
		return err
	}
	if err := domain.PodeRecolher(u, serie, r.Motivo, s.hoje()); err != nil {
		return prefixar(u, err)
	}
	checks := []struct {
		fn  func(context.Context, database.DBTX, uuid.UUID) (bool, error)
		err error
	}{{s.repo.EmprestimoAberto, domain.ErrEmprestada}, {s.repo.EmGuiaEnviada, domain.ErrEmOutraGuia}}
	if incluir {
		checks = append(checks, struct {
			fn  func(context.Context, database.DBTX, uuid.UUID) (bool, error)
			err error
		}{s.repo.EmRecolhimentoAberto, domain.ErrEmOutroRecolhimento})
	}
	for _, c := range checks {
		sim, err := c.fn(ctx, db, u.ID)
		if err != nil {
			return err
		}
		if sim {
			return prefixar(u, c.err)
		}
	}
	return nil
}

// NovoRecolhimento é o pedido da origem.
type NovoRecolhimento struct {
	OrigemID   uuid.UUID
	Motivo     string
	Fundamento string
	Preparo    domain.Preparo
	Historico  string
	Observacao string
	Unidades   []uuid.UUID
}

// CriarRecolhimento relaciona as unidades e envia ao Arquivo Público. As
// unidades só mudam de guarda na entrada.
func (s *Service) CriarRecolhimento(ctx context.Context, identity auth.Identity, n NovoRecolhimento) (domain.Recolhimento, error) {
	r, ids, err := domain.NovoRecolhimento(domain.Recolhimento{OrigemID: n.OrigemID, Motivo: n.Motivo, Fundamento: n.Fundamento,
		Preparo: n.Preparo, Historico: n.Historico, Observacao: n.Observacao}, n.Unidades)
	if err != nil {
		return domain.Recolhimento{}, MapError(err)
	}
	if !gere(identity, r.OrigemID) {
		return domain.Recolhimento{}, MapError(domain.ErrForaDoEscopo)
	}
	r.CriadoPor, r.CriadoEm = ator(identity), s.agora().UTC()
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		lote := make([]domain.Unidade, 0, len(ids))
		for _, id := range ids {
			u, _, err := s.repo.GetUnidade(ctx, tx, id, true)
			if err != nil {
				return err
			}
			if err := s.conferirRecolhimento(ctx, tx, r, u, true); err != nil {
				return err
			}
			lote = append(lote, u)
		}
		if r.Numero, err = s.repo.InsertRecolhimento(ctx, tx, r, ids); err != nil {
			return err
		}
		for _, u := range lote {
			if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovRecolhimento, Por: ator(identity),
				Descricao: "Relacionada no recolhimento " + r.Numero + " ao Arquivo Público"}); err != nil {
				return err
			}
		}
		return s.registrarRecolhimento(ctx, tx, EventRecolhimentoEnviado, r, nil,
			map[string]any{"numero": r.Numero, "origem_id": r.OrigemID, "motivo": r.Motivo, "unidades": ids})
	})
	if err != nil {
		return domain.Recolhimento{}, MapError(err)
	}
	return s.Recolhimento(ctx, identity, r.ID)
}

func (s *Service) registrarRecolhimento(ctx context.Context, tx pgx.Tx, evento string, r domain.Recolhimento, antes, depois any) error {
	if err := s.outbox.Write(ctx, tx, evento, "acervo_recolhimento", r.ID.String(), uuid.Nil,
		map[string]any{"id": r.ID, "numero": r.Numero, "origem_id": r.OrigemID, "situacao": r.Situacao}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_recolhimento", r.ID.String(), antes, depois))
}

// carregarRecolhimento lê o recolhimento enviado para decidir. Central:
// recebe ou devolve (Arquivo Público); senão, a origem cancela.
func (s *Service) carregarRecolhimento(ctx context.Context, tx pgx.Tx, identity auth.Identity, id uuid.UUID, central bool) (domain.Recolhimento, error) {
	r, err := s.repo.GetRecolhimento(ctx, tx, id, true)
	if err != nil {
		return r, err
	}
	if !podeLerRecolhimento(identity, r) {
		return r, domain.ErrRecolhimentoNaoExiste
	}
	if central && !auth.CanGlobal(identity, auth.PermAcervoCentral) || !central && !gere(identity, r.OrigemID) {
		return r, domain.ErrForaDoEscopo
	}
	return r, domain.PodeDecidir(r)
}

func podeLerRecolhimento(identity auth.Identity, r domain.Recolhimento) bool {
	return le(identity, r.OrigemID) || auth.CanGlobal(identity, auth.PermAcervoCentral)
}

// Entrada é o registro do Arquivo Público ao receber.
type Entrada struct {
	LocalID *uuid.UUID
}

// ReceberRecolhimento registra a entrada no Arquivo Público: cada unidade é
// conferida de novo e passa à guarda permanente do Arquivo, com o local
// dado; a origem recebe a cópia do registro (art. 5º, VII).
func (s *Service) ReceberRecolhimento(ctx context.Context, identity auth.Identity, id uuid.UUID, e Entrada) (domain.Recolhimento, error) {
	var enviadas []notificacoes.Enviada
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		r, err := s.carregarRecolhimento(ctx, tx, identity, id, true)
		if err != nil {
			return err
		}
		if err := s.conferirLocal(ctx, tx, e.LocalID); err != nil {
			return err
		}
		agora := s.agora().UTC()
		r.Situacao, r.RecebidoPor, r.RecebidoEm, r.LocalID = domain.RecolhimentoRecebido, ator(identity), &agora, e.LocalID
		if r.RegistroEntrada, err = s.repo.ReceberRecolhimento(ctx, tx, r); err != nil {
			return err
		}
		if err := s.registrarRecolhimento(ctx, tx, EventRecolhimentoRecebido, r, map[string]any{"situacao": domain.RecolhimentoEnviado},
			map[string]any{"situacao": r.Situacao, "registro_entrada": r.RegistroEntrada, "local_id": r.LocalID}); err != nil {
			return err
		}
		for _, item := range r.Itens {
			if err := s.recolherUnidade(ctx, tx, identity, r, item.ID); err != nil {
				return err
			}
		}
		enviadas, err = s.avisarOrigem(ctx, tx, identity, r, fmt.Sprintf("Recolhimento %s recebido: registro de entrada %s", r.Numero, r.RegistroEntrada),
			fmt.Sprintf("O Arquivo Público registrou a entrada de %d unidade(s) de arquivamento (Lei 10.301/2019, art. 5º, VII).", len(r.Itens)))
		return err
	})
	if err != nil {
		return domain.Recolhimento{}, MapError(err)
	}
	s.entregar(ctx, enviadas)
	return s.Recolhimento(ctx, identity, id)
}

func (s *Service) recolherUnidade(ctx context.Context, tx pgx.Tx, identity auth.Identity, r domain.Recolhimento, id uuid.UUID) error {
	u, _, err := s.repo.GetUnidade(ctx, tx, id, true)
	if err != nil {
		return err
	}
	if err := s.conferirRecolhimento(ctx, tx, r, u, false); err != nil {
		return err
	}
	antes := map[string]any{"fase": u.Fase, "local_id": u.LocalID}
	de := u.Fase
	u.Fase, u.Destinacao, u.RecolhidaEm, u.RecolhimentoID = domain.FasePermanente, domain.DestinacaoGuardaPermanente, r.RecebidoEm, &r.ID
	if r.LocalID != nil {
		u.LocalID = r.LocalID
	}
	_, p, err := s.calcular(ctx, tx, u)
	if err != nil {
		return err
	}
	if err := s.repo.MarcarRecolhida(ctx, tx, u, calculo(p)); err != nil {
		return err
	}
	para := u.Fase
	return s.registrar(ctx, tx, identity, EventRecolhida, u, domain.Movimentacao{Tipo: domain.MovRecolhimento, FaseDe: &de, FasePara: &para,
		Descricao: "Recolhida ao Arquivo Público pelo recolhimento " + r.Numero + ", registro de entrada " + r.RegistroEntrada + " (arts. 5º, VII, e 53)"},
		antes, map[string]any{"fase": para, "registro_entrada": r.RegistroEntrada, "local_id": u.LocalID})
}

// avisarOrigem grava o aviso para quem gere o acervo da origem e para quem
// enviou o recolhimento (menos quem decidiu).
func (s *Service) avisarOrigem(ctx context.Context, tx pgx.Tx, identity auth.Identity, r domain.Recolhimento, titulo, mensagem string) ([]notificacoes.Enviada, error) {
	usuarios, err := s.repo.GestoresDaUnidade(ctx, tx, r.OrigemID)
	if err != nil {
		return nil, err
	}
	if r.CriadoPor != nil {
		usuarios = append(usuarios, *r.CriadoPor)
	}
	vistos := map[uuid.UUID]bool{identity.UserID: true}
	destino := usuarios[:0:0]
	for _, u := range usuarios {
		if !vistos[u] {
			vistos[u] = true
			destino = append(destino, u)
		}
	}
	return notificacoes.Registrar(ctx, tx, destino, notificacoes.Nova{Modulo: "acervo", Titulo: titulo, Mensagem: mensagem,
		Link: "/acervo/recolhimentos/" + r.ID.String(), Chave: "acervo:recolhimento:" + r.ID.String() + ":" + r.Situacao})
}

// entregar publica os avisos gravados (depois do commit).
func (s *Service) entregar(ctx context.Context, enviadas []notificacoes.Enviada) {
	if s.avisos != nil {
		s.avisos.Entregar(ctx, enviadas)
	}
}

// DevolverRecolhimento: o Arquivo Público devolve à origem com o que falta
// (ex.: preparo do art. 53 incompleto). A origem envia de novo.
func (s *Service) DevolverRecolhimento(ctx context.Context, identity auth.Identity, id uuid.UUID, motivo string) (domain.Recolhimento, error) {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" || utf8.RuneCountInString(motivo) > 2000 {
		return domain.Recolhimento{}, MapError(&domain.ErrValidacao{Campo: "motivo", Motivo: "informe o que falta (até 2000 caracteres)"})
	}
	var enviadas []notificacoes.Enviada
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		r, err := s.carregarRecolhimento(ctx, tx, identity, id, true)
		if err != nil {
			return err
		}
		r.Situacao, r.DevolucaoMotivo = domain.RecolhimentoDevolvido, motivo
		if err := s.repo.SetRecolhimentoSituacao(ctx, tx, r); err != nil {
			return err
		}
		if err := s.registrarRecolhimento(ctx, tx, EventRecolhimentoDevolvido, r, map[string]any{"situacao": domain.RecolhimentoEnviado},
			map[string]any{"situacao": r.Situacao, "motivo": motivo}); err != nil {
			return err
		}
		enviadas, err = s.avisarOrigem(ctx, tx, identity, r, "Recolhimento "+r.Numero+" devolvido pelo Arquivo Público", motivo)
		return err
	})
	if err != nil {
		return domain.Recolhimento{}, MapError(err)
	}
	s.entregar(ctx, enviadas)
	return s.Recolhimento(ctx, identity, id)
}

// CancelarRecolhimento: a origem desiste antes da entrada.
func (s *Service) CancelarRecolhimento(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Recolhimento, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		r, err := s.carregarRecolhimento(ctx, tx, identity, id, false)
		if err != nil {
			return err
		}
		r.Situacao = domain.RecolhimentoCancelado
		if err := s.repo.SetRecolhimentoSituacao(ctx, tx, r); err != nil {
			return err
		}
		return s.registrarRecolhimento(ctx, tx, EventRecolhimentoCancelado, r, map[string]any{"situacao": domain.RecolhimentoEnviado},
			map[string]any{"situacao": r.Situacao})
	})
	if err != nil {
		return domain.Recolhimento{}, MapError(err)
	}
	return s.Recolhimento(ctx, identity, id)
}

// Recolhimento devolve o recolhimento com as unidades (o instrumento
// descritivo). Fora do escopo de leitura, 404.
func (s *Service) Recolhimento(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Recolhimento, error) {
	r, err := s.repo.GetRecolhimento(ctx, s.pool, id, false)
	if err == nil && !podeLerRecolhimento(identity, r) {
		err = domain.ErrRecolhimentoNaoExiste
	}
	return r, MapError(err)
}

// Recolhimentos pagina os recolhimentos no escopo; o Arquivo Público vê
// todos.
func (s *Service) Recolhimentos(ctx context.Context, identity auth.Identity, p pagination.Params) ([]domain.Recolhimento, int64, error) {
	f := s.filtro(identity, domain.Filtro{})
	if auth.CanGlobal(identity, auth.PermAcervoCentral) {
		f.Restrito = false
	}
	items, total, err := s.repo.ListRecolhimentos(ctx, s.pool, f, p)
	return items, total, MapError(err)
}
