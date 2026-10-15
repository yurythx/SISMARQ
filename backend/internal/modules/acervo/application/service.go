// Package application tem os casos de uso do Acervo. Cada escrita grava, na
// mesma transação, a unidade, o histórico, o evento (outbox) e a auditoria
// (ADR 016).
package application

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/events"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/outbox"
)

// Eventos do Acervo (outbox).
const (
	EventCadastrada    = "acervo.unidade.cadastrada"
	EventAlterada      = "acervo.unidade.alterada"
	EventTransferida   = "acervo.unidade.transferida"
	EventLocalizada    = "acervo.unidade.localizada"
	EventBloqueada     = "acervo.unidade.bloqueada"
	EventDesbloqueada  = "acervo.unidade.desbloqueada"
	EventLocalSalvo    = "acervo.local.salvo"
	recalculoPorPagina = 200
)

// SituacaoSerieRevogada: a série saiu da TTDD em vigor — a unidade precisa
// ser reclassificada antes de seguir o ciclo (art. 52).
const SituacaoSerieRevogada = "SERIE_REVOGADA"

// Service implementa os casos de uso.
type Service struct {
	pool   *pgxpool.Pool
	repo   domain.Repository
	ttdd   domain.TTDD
	outbox *outbox.Writer
	logger *slog.Logger
	agora  func() time.Time
}

// NewService cria o serviço.
func NewService(pool *pgxpool.Pool, repo domain.Repository, ttdd domain.TTDD, ob *outbox.Writer, logger *slog.Logger) *Service {
	return &Service{pool: pool, repo: repo, ttdd: ttdd, outbox: ob, logger: logger, agora: time.Now}
}

// WithRelogio troca o relógio (testes dos prazos).
func (s *Service) WithRelogio(agora func() time.Time) *Service {
	s.agora = agora
	return s
}

func (s *Service) hoje() time.Time {
	t := s.agora().UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// MapError traduz os erros de domínio para a API.
func MapError(err error) error {
	var v *domain.ErrValidacao
	switch {
	case err == nil:
		return nil
	case errors.As(err, &v):
		return apperrors.Validation(v.Campo + ": " + v.Motivo)
	case errors.Is(err, domain.ErrNaoEncontrada):
		return apperrors.NotFound("unidade de arquivamento não encontrada")
	case errors.Is(err, domain.ErrForaDoEscopo):
		return apperrors.Forbidden("a unidade de arquivamento está fora do seu escopo")
	case errors.Is(err, domain.ErrLocalNaoExiste):
		return apperrors.NotFound("local não encontrado")
	case errors.Is(err, domain.ErrBloqueioNaoExiste):
		return apperrors.NotFound("bloqueio não encontrado")
	case errors.Is(err, domain.ErrSerieInvalida):
		return apperrors.Validation("codigo_classificacao: série inexistente ou revogada na Tabela de Temporalidade em vigor")
	case errors.Is(err, domain.ErrBloqueioLiberado), errors.Is(err, domain.ErrPermanente), errors.Is(err, domain.ErrBloqueioJudicial),
		errors.Is(err, domain.ErrNaoPrevistaNaTTDD), errors.Is(err, domain.ErrPrazoNaoCumprido), errors.Is(err, domain.ErrAguardaEvento),
		errors.Is(err, domain.ErrDestinacaoIndefinida), errors.Is(err, domain.ErrFaseInvalida), errors.Is(err, domain.ErrEliminado):
		return apperrors.Conflict(strings.TrimPrefix(err.Error(), "acervo: "))
	}
	return err
}

// Detalhe é a ficha da unidade: dados, prazos calculados na hora, série,
// bloqueios, histórico e o que quem pede pode fazer.
type Detalhe struct {
	domain.UnidadeLista
	Serie         domain.Serie          `json:"-"`
	SerieResumo   SerieResumo           `json:"serie"`
	Prazos        domain.Prazos         `json:"prazos"`
	Bloqueios     []domain.Bloqueio     `json:"bloqueios"`
	Movimentacoes []domain.Movimentacao `json:"movimentacoes"`
	Pode          Permissoes            `json:"pode"`
}

// SerieResumo é a série da TTDD aplicada à unidade.
type SerieResumo struct {
	Codigo     string `json:"codigo"`
	Descritor  string `json:"descritor"`
	Vigente    bool   `json:"vigente"`
	Destinacao string `json:"destinacao"`
}

// Permissoes diz o que a pessoa pode fazer com a unidade.
type Permissoes struct {
	Gerir    bool `json:"gerir"`
	Bloquear bool `json:"bloquear"`
}

func ator(identity auth.Identity) *uuid.UUID {
	if identity.UserID == uuid.Nil {
		return nil
	}
	id := identity.UserID
	return &id
}

func gere(identity auth.Identity, unidade uuid.UUID) bool {
	return auth.Can(identity, auth.PermAcervoManage, auth.InUnidade(unidade))
}

func le(identity auth.Identity, unidade uuid.UUID) bool {
	return auth.Can(identity, auth.PermAcervoRead, auth.InUnidade(unidade)) || gere(identity, unidade)
}

// calcular aplica a TTDD em vigor à unidade.
func (s *Service) calcular(ctx context.Context, db database.DBTX, u domain.Unidade) (domain.Serie, domain.Prazos, error) {
	serie, err := s.ttdd.Serie(ctx, u.CodigoClassificacao)
	if err != nil {
		return serie, domain.Prazos{}, err
	}
	bloqueada, err := s.repo.TemBloqueioAtivo(ctx, db, u.ID)
	if err != nil {
		return serie, domain.Prazos{}, err
	}
	p := domain.CalcularPrazos(u, serie, bloqueada, s.hoje())
	if !serie.Vigente && (u.Fase == domain.FaseCorrente || u.Fase == domain.FaseIntermediaria) {
		p.Situacao = SituacaoSerieRevogada
	}
	return serie, p, nil
}

func calculo(p domain.Prazos) domain.Calculo {
	return domain.Calculo{FimCorrente: p.FimCorrente, FimIntermediaria: p.FimIntermediaria, Situacao: p.Situacao}
}

// registrar grava o histórico, o evento e a auditoria da operação.
func (s *Service) registrar(ctx context.Context, tx pgx.Tx, identity auth.Identity, evento string, u domain.Unidade,
	m domain.Movimentacao, antes, depois any) error {
	m.Por = ator(identity)
	if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, m); err != nil {
		return err
	}
	payload := map[string]string{"id": u.ID.String(), "numero": u.Numero, "codigo": u.CodigoClassificacao,
		"unidade_id": u.UnidadeID.String(), "fase": string(u.Fase)}
	if err := s.outbox.Write(ctx, tx, evento, "acervo_unidade", u.ID.String(), uuid.Nil, payload); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_unidade", u.ID.String(), antes, depois))
}

// Criar cadastra a unidade de arquivamento na fase corrente. A série precisa
// existir e estar em vigor na TTDD; a destinação vem dela (art. 28).
func (s *Service) Criar(ctx context.Context, identity auth.Identity, u domain.Unidade) (Detalhe, error) {
	u.Normalizar()
	if err := u.Validar(); err != nil {
		return Detalhe{}, MapError(err)
	}
	if !gere(identity, u.UnidadeID) {
		return Detalhe{}, MapError(domain.ErrForaDoEscopo)
	}
	u.ID, u.Fase, u.CreatedBy = uuid.New(), domain.FaseCorrente, ator(identity)
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.conferirLocal(ctx, tx, u.LocalID); err != nil {
			return err
		}
		serie, p, err := s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if !serie.Vigente {
			return domain.ErrSerieInvalida
		}
		u.Destinacao = domain.DestinacaoDa(serie)
		if u.Assunto == "" {
			u.Assunto = serie.Descritor
		}
		if u.Numero, err = s.repo.InsertUnidade(ctx, tx, u, calculo(p)); err != nil {
			return err
		}
		fase := u.Fase
		return s.registrar(ctx, tx, identity, EventCadastrada, u,
			domain.Movimentacao{Tipo: domain.MovCadastro, FasePara: &fase, Descricao: "Cadastro na fase corrente"}, nil, u)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, u.ID)
}

func (s *Service) conferirLocal(ctx context.Context, db database.DBTX, id *uuid.UUID) error {
	if id == nil {
		return nil
	}
	l, err := s.repo.GetLocal(ctx, db, *id)
	if err != nil {
		return err
	}
	if !l.Ativo {
		return &domain.ErrValidacao{Campo: "local_id", Motivo: "local desativado"}
	}
	return nil
}

// carregar lê a unidade para alterá-la e confere o escopo de gestão.
func (s *Service) carregar(ctx context.Context, tx pgx.Tx, identity auth.Identity, id uuid.UUID, perm auth.Permission) (domain.Unidade, error) {
	u, _, err := s.repo.GetUnidade(ctx, tx, id, true)
	if err != nil {
		return u, err
	}
	if !le(identity, u.UnidadeID) {
		return u, domain.ErrNaoEncontrada
	}
	if !auth.Can(identity, perm, auth.InUnidade(u.UnidadeID)) {
		return u, domain.ErrForaDoEscopo
	}
	return u, nil
}

// Alterar muda a descrição, a classificação, as datas, o evento, o acesso ou
// a unidade dona. Eliminada não muda; permanente não é reclassificada.
func (s *Service) Alterar(ctx context.Context, identity auth.Identity, id uuid.UUID, novo domain.Unidade) (Detalhe, error) {
	novo.Normalizar()
	if err := novo.Validar(); err != nil {
		return Detalhe{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		antes, err := s.carregar(ctx, tx, identity, id, auth.PermAcervoManage)
		if err != nil {
			return err
		}
		if err := antes.Editavel(); err != nil {
			return err
		}
		if antes.Fase == domain.FasePermanente && novo.CodigoClassificacao != antes.CodigoClassificacao {
			return domain.ErrPermanente
		}
		if novo.UnidadeID != antes.UnidadeID && !gere(identity, novo.UnidadeID) {
			return domain.ErrForaDoEscopo
		}
		if err := s.conferirLocal(ctx, tx, novo.LocalID); err != nil {
			return err
		}
		u := novo
		u.ID, u.Numero, u.Fase, u.CreatedBy, u.CreatedAt = antes.ID, antes.Numero, antes.Fase, antes.CreatedBy, antes.CreatedAt
		serie, p, err := s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if !serie.Vigente && u.CodigoClassificacao != antes.CodigoClassificacao {
			return domain.ErrSerieInvalida
		}
		u.Destinacao = antes.Destinacao
		if u.Fase != domain.FasePermanente {
			u.Destinacao = domain.DestinacaoDa(serie)
		}
		if err := s.repo.UpdateUnidade(ctx, tx, u, calculo(p)); err != nil {
			return err
		}
		return s.registrar(ctx, tx, identity, EventAlterada, u, domain.Movimentacao{Tipo: domain.MovAlteracao, Descricao: "Dados alterados"}, antes, u)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, id)
}

// Transferir leva a unidade do arquivo corrente ao intermediário do órgão,
// só com o prazo corrente cumprido (art. 29).
func (s *Service) Transferir(ctx context.Context, identity auth.Identity, id uuid.UUID, observacao string) (Detalhe, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		u, err := s.carregar(ctx, tx, identity, id, auth.PermAcervoManage)
		if err != nil {
			return err
		}
		_, p, err := s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if err := domain.PodeTransferir(u, p, s.hoje()); err != nil {
			return err
		}
		de := u.Fase
		u.Fase = domain.FaseIntermediaria
		_, p, err = s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if err := s.repo.SetFase(ctx, tx, u.ID, u.Fase, calculo(p)); err != nil {
			return err
		}
		para := u.Fase
		desc := "Transferência ao arquivo intermediário (art. 29)"
		if o := strings.TrimSpace(observacao); o != "" {
			desc += ": " + o
		}
		return s.registrar(ctx, tx, identity, EventTransferida, u,
			domain.Movimentacao{Tipo: domain.MovTransferencia, FaseDe: &de, FasePara: &para, Descricao: desc},
			map[string]any{"fase": de}, map[string]any{"fase": para})
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, id)
}

// Localizar põe a unidade num local físico (ou a retira, com local nil).
func (s *Service) Localizar(ctx context.Context, identity auth.Identity, id uuid.UUID, local *uuid.UUID, posicao string) (Detalhe, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		u, err := s.carregar(ctx, tx, identity, id, auth.PermAcervoManage)
		if err != nil {
			return err
		}
		if err := u.Editavel(); err != nil {
			return err
		}
		if err := s.conferirLocal(ctx, tx, local); err != nil {
			return err
		}
		antes := map[string]any{"local_id": u.LocalID, "posicao": u.Posicao}
		u.LocalID, u.Posicao = local, strings.TrimSpace(posicao)
		if err := u.Validar(); err != nil {
			return err
		}
		_, p, err := s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateUnidade(ctx, tx, u, calculo(p)); err != nil {
			return err
		}
		return s.registrar(ctx, tx, identity, EventLocalizada, u,
			domain.Movimentacao{Tipo: domain.MovLocalizacao, Descricao: "Localização física alterada"},
			antes, map[string]any{"local_id": u.LocalID, "posicao": u.Posicao})
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, id)
}

// Bloquear retém a unidade por ação judicial (arts. 32, II, e 44): enquanto
// ativo, nada a elimina.
func (s *Service) Bloquear(ctx context.Context, identity auth.Identity, id uuid.UUID, processo, motivo string) (Detalhe, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		u, err := s.carregar(ctx, tx, identity, id, auth.PermAcervoBloqueio)
		if err != nil {
			return err
		}
		if err := u.Editavel(); err != nil {
			return err
		}
		b, err := domain.NovoBloqueio(u.ID, processo, motivo)
		if err != nil {
			return err
		}
		b.CriadoPor, b.CriadoEm = ator(identity), s.agora().UTC()
		if err := s.repo.InsertBloqueio(ctx, tx, b); err != nil {
			return err
		}
		if err := s.recalcularUma(ctx, tx, u); err != nil {
			return err
		}
		return s.registrar(ctx, tx, identity, EventBloqueada, u,
			domain.Movimentacao{Tipo: domain.MovBloqueio, Descricao: "Bloqueio judicial: processo " + b.Processo}, nil, b)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, id)
}

// Liberar encerra um bloqueio judicial, com o motivo.
func (s *Service) Liberar(ctx context.Context, identity auth.Identity, bloqueioID uuid.UUID, motivo string) (Detalhe, error) {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" || len([]rune(motivo)) > 2000 {
		return Detalhe{}, MapError(&domain.ErrValidacao{Campo: "motivo", Motivo: "informe o motivo da liberação (até 2000 caracteres)"})
	}
	var unidade uuid.UUID
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		b, err := s.repo.GetBloqueio(ctx, tx, bloqueioID)
		if err != nil {
			return err
		}
		u, err := s.carregar(ctx, tx, identity, b.UnidadeID, auth.PermAcervoBloqueio)
		if err != nil {
			return err
		}
		if !b.Ativo() {
			return domain.ErrBloqueioLiberado
		}
		agora := s.agora().UTC()
		antes := b
		b.LiberadoPor, b.LiberadoEm, b.MotivoLiberacao = ator(identity), &agora, motivo
		if err := s.repo.LiberarBloqueio(ctx, tx, b); err != nil {
			return err
		}
		if err := s.recalcularUma(ctx, tx, u); err != nil {
			return err
		}
		unidade = u.ID
		return s.registrar(ctx, tx, identity, EventDesbloqueada, u,
			domain.Movimentacao{Tipo: domain.MovDesbloqueio, Descricao: "Bloqueio liberado (processo " + b.Processo + "): " + motivo}, antes, b)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, unidade)
}

func (s *Service) recalcularUma(ctx context.Context, db database.DBTX, u domain.Unidade) error {
	_, p, err := s.calcular(ctx, db, u)
	if err != nil {
		return err
	}
	return s.repo.SetCalculo(ctx, db, u.ID, calculo(p))
}

// Get devolve a ficha. Fora do escopo de leitura, 404 (ADR 013).
func (s *Service) Get(ctx context.Context, identity auth.Identity, id uuid.UUID) (Detalhe, error) {
	u, c, err := s.repo.GetUnidade(ctx, s.pool, id, false)
	if err == nil && !le(identity, u.UnidadeID) {
		err = domain.ErrNaoEncontrada
	}
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	d := Detalhe{UnidadeLista: domain.UnidadeLista{Unidade: u, Calculo: c}}
	if u.LocalID != nil {
		l, err := s.repo.GetLocal(ctx, s.pool, *u.LocalID)
		if err != nil {
			return Detalhe{}, MapError(err)
		}
		d.LocalCaminho = l.Caminho
	}
	serie, p, err := s.calcular(ctx, s.pool, u)
	switch {
	case errors.Is(err, domain.ErrSerieInvalida):
		d.Prazos = domain.Prazos{FimCorrente: c.FimCorrente, FimIntermediaria: c.FimIntermediaria, Situacao: c.Situacao}
		d.SerieResumo = SerieResumo{Codigo: u.CodigoClassificacao}
	case err != nil:
		return Detalhe{}, MapError(err)
	default:
		d.Serie, d.Prazos = serie, p
		d.SerieResumo = SerieResumo{Codigo: serie.Codigo, Descritor: serie.Descritor, Vigente: serie.Vigente, Destinacao: string(domain.DestinacaoDa(serie))}
	}
	if d.Bloqueios, err = s.repo.ListBloqueios(ctx, s.pool, u.ID); err != nil {
		return Detalhe{}, MapError(err)
	}
	if d.Movimentacoes, err = s.repo.ListMovimentacoes(ctx, s.pool, u.ID); err != nil {
		return Detalhe{}, MapError(err)
	}
	d.Pode = Permissoes{Gerir: gere(identity, u.UnidadeID), Bloquear: auth.Can(identity, auth.PermAcervoBloqueio, auth.InUnidade(u.UnidadeID))}
	return d, nil
}

func (s *Service) filtro(identity auth.Identity, f domain.Filtro) domain.Filtro {
	leitura := auth.CoverageOf(identity, auth.PermAcervoRead)
	gestao := auth.CoverageOf(identity, auth.PermAcervoManage)
	f.Restrito = !leitura.All && !gestao.All
	f.Unidades = append(append([]uuid.UUID{}, leitura.Unidades...), gestao.Unidades...)
	return f
}

// List pagina as unidades que a pessoa pode ver.
func (s *Service) List(ctx context.Context, identity auth.Identity, f domain.Filtro, p pagination.Params) ([]domain.UnidadeLista, int64, error) {
	items, total, err := s.repo.ListUnidades(ctx, s.pool, s.filtro(identity, f), p)
	return items, total, MapError(err)
}

// Painel conta as unidades por situação (o que transferir, avaliar,
// recolher…), no escopo de quem pede.
func (s *Service) Painel(ctx context.Context, identity auth.Identity, f domain.Filtro) (map[string]int64, error) {
	out, err := s.repo.Painel(ctx, s.pool, s.filtro(identity, f))
	return out, MapError(err)
}

// Recalcular refaz o cálculo de todas as unidades das fases corrente e
// intermediária (os prazos andam com o calendário e com a TTDD). Devolve
// quantas foram recalculadas.
func (s *Service) Recalcular(ctx context.Context) (int, error) {
	n, depois := 0, uuid.Nil
	for {
		lote, err := s.repo.ParaRecalculo(ctx, s.pool, depois, recalculoPorPagina)
		if err != nil {
			return n, err
		}
		for _, u := range lote {
			err := s.recalcularUma(ctx, s.pool, u)
			if errors.Is(err, domain.ErrSerieInvalida) {
				err = s.repo.SetCalculo(ctx, s.pool, u.ID, domain.Calculo{Situacao: SituacaoSerieRevogada})
			}
			if err != nil {
				return n, err
			}
			n++
			depois = u.ID
		}
		if len(lote) < recalculoPorPagina {
			return n, nil
		}
	}
}

// RunRecalculo recalcula no início e a cada intervalo, até ctx acabar.
func (s *Service) RunRecalculo(intervalo time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()
		for {
			if n, err := s.Recalcular(ctx); err != nil {
				s.logger.Error("acervo: recálculo dos prazos falhou", slog.Any("error", err))
			} else {
				s.logger.Info("acervo: prazos recalculados", slog.Int("unidades", n))
			}
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}
		}
	}
}

// HandleTTDD recalcula quando a TTDD muda (evento do Atlas).
func (s *Service) HandleTTDD(ctx context.Context, _ events.Event) error {
	_, err := s.Recalcular(ctx)
	return err
}
