// Package transport expõe o Acervo via HTTP.
package transport

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// Handlers agrupa os handlers do Acervo.
type Handlers struct {
	svc         *application.Service
	logger      *slog.Logger
	maxPageSize int
}

// NewHandlers cria os handlers.
func NewHandlers(svc *application.Service, logger *slog.Logger, maxPageSize int) *Handlers {
	return &Handlers{svc: svc, logger: logger, maxPageSize: maxPageSize}
}

// RegisterRoutes monta as rotas (todas autenticadas). A permissão da rota
// diz o que a pessoa faz em ALGUM lugar; o escopo de cada unidade é decidido
// no serviço (ADR 013).
func (h *Handlers) RegisterRoutes(r chi.Router) {
	ler := auth.RequirePermission(h.logger, auth.PermAcervoRead)
	gerir := auth.RequirePermission(h.logger, auth.PermAcervoManage)
	bloquear := auth.RequirePermission(h.logger, auth.PermAcervoBloqueio)
	central := auth.RequirePermission(h.logger, auth.PermAcervoCentral)
	r.Route("/acervo", func(r chi.Router) {
		r.With(ler).Get("/unidades", h.List)
		r.With(ler).Get("/painel", h.Painel)
		r.With(ler).Get("/unidades/{id}", h.Get)
		r.With(gerir).Post("/unidades", h.Create)
		r.With(gerir).Put("/unidades/{id}", h.Update)
		r.With(gerir).Post("/unidades/{id}/transferir", h.Transferir)
		r.With(gerir).Put("/unidades/{id}/local", h.Localizar)
		r.With(bloquear).Post("/unidades/{id}/bloqueios", h.Bloquear)
		r.With(bloquear).Post("/bloqueios/{id}/liberar", h.Liberar)
		r.With(ler).Get("/locais", h.ListLocais)
		r.With(central).Post("/locais", h.CreateLocal)
		r.With(central).Put("/locais/{id}", h.UpdateLocal)
		h.registerFase2(r, ler, gerir)
	})
}

func (h *Handlers) fail(w http.ResponseWriter, r *http.Request, err error) {
	httputil.WriteError(w, r, h.logger, application.MapError(err))
}

// UnidadeRequest é o corpo do cadastro e da alteração. Datas em AAAA-MM-DD.
type UnidadeRequest struct {
	Tipo                string     `json:"tipo" validate:"required,max=16"`
	CodigoClassificacao string     `json:"codigo_classificacao" validate:"required,max=32"`
	Assunto             string     `json:"assunto" validate:"max=500"`
	UnidadeID           uuid.UUID  `json:"unidade_id" validate:"required"`
	DataInicial         string     `json:"data_inicial" validate:"required"`
	DataFinal           string     `json:"data_final" validate:"required"`
	Quantidade          int        `json:"quantidade" validate:"min=0,max=100000"`
	Especificacao       string     `json:"especificacao" validate:"max=120"`
	Observacoes         string     `json:"observacoes" validate:"max=4000"`
	DataEvento          string     `json:"data_evento"`
	ContasAprovadasEm   string     `json:"contas_aprovadas_em"`
	NivelAcesso         string     `json:"nivel_acesso" validate:"max=16"`
	HipoteseLegal       string     `json:"hipotese_legal" validate:"max=500"`
	LocalID             *uuid.UUID `json:"local_id"`
	Posicao             string     `json:"posicao" validate:"max=120"`
}

func data(campo, v string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, strings.TrimSpace(v))
	if err != nil {
		return t, apperrors.Validation(campo + ": use AAAA-MM-DD")
	}
	return t, nil
}

func dataOpcional(campo, v string) (*time.Time, error) {
	if strings.TrimSpace(v) == "" {
		return nil, nil
	}
	t, err := data(campo, v)
	return &t, err
}

func (req UnidadeRequest) unidade() (domain.Unidade, error) {
	u := domain.Unidade{
		Tipo: domain.Tipo(req.Tipo), CodigoClassificacao: req.CodigoClassificacao, Assunto: req.Assunto, UnidadeID: req.UnidadeID,
		Quantidade: req.Quantidade, Especificacao: req.Especificacao, Observacoes: req.Observacoes,
		NivelAcesso: domain.NivelAcesso(req.NivelAcesso), HipoteseLegal: req.HipoteseLegal, LocalID: req.LocalID, Posicao: req.Posicao,
	}
	var err error
	if u.DataInicial, err = data("data_inicial", req.DataInicial); err != nil {
		return u, err
	}
	if u.DataFinal, err = data("data_final", req.DataFinal); err != nil {
		return u, err
	}
	if u.DataEvento, err = dataOpcional("data_evento", req.DataEvento); err != nil {
		return u, err
	}
	u.ContasAprovadasEm, err = dataOpcional("contas_aprovadas_em", req.ContasAprovadasEm)
	return u, err
}

func (h *Handlers) bindUnidade(w http.ResponseWriter, r *http.Request) (domain.Unidade, bool) {
	var req UnidadeRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return domain.Unidade{}, false
	}
	u, err := req.unidade()
	if err != nil {
		h.fail(w, r, err)
		return u, false
	}
	return u, true
}

// Create cadastra uma unidade de arquivamento.
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	u, ok := h.bindUnidade(w, r)
	if !ok {
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Criar(r.Context(), identity, u)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, d)
}

// Update altera uma unidade de arquivamento.
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	u, ok := h.bindUnidade(w, r)
	if !ok {
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Alterar(r.Context(), identity, id, u)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

// Get devolve a ficha da unidade.
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Get(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

func (h *Handlers) filtro(r *http.Request) (domain.Filtro, error) {
	unidade, err := httputil.OptionalUUIDQuery(r, "unidade_id")
	return domain.Filtro{
		Query: httputil.Query(r, "q", 200), Fase: domain.Fase(strings.ToUpper(httputil.Query(r, "fase", 16))),
		Situacao: strings.ToUpper(httputil.Query(r, "situacao", 32)), Codigo: httputil.Query(r, "codigo", 32), UnidadeID: unidade,
	}, err
}

// List pagina as unidades no escopo de quem pede.
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	f, err := h.filtro(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.List(r.Context(), identity, f, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Painel conta as unidades por situação.
func (h *Handlers) Painel(w http.ResponseWriter, r *http.Request) {
	f, err := h.filtro(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	out, err := h.svc.Painel(r.Context(), identity, f)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, out)
}

// TransferirRequest é a observação opcional da transferência.
type TransferirRequest struct {
	Observacao string `json:"observacao" validate:"max=500"`
}

// Transferir leva a unidade ao arquivo intermediário.
func (h *Handlers) Transferir(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req TransferirRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Transferir(r.Context(), identity, id, req.Observacao)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

// LocalizarRequest põe a unidade num local (local_id nulo a retira).
type LocalizarRequest struct {
	LocalID *uuid.UUID `json:"local_id"`
	Posicao string     `json:"posicao" validate:"max=120"`
}

// Localizar altera a localização física.
func (h *Handlers) Localizar(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req LocalizarRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Localizar(r.Context(), identity, id, req.LocalID, req.Posicao)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

// BloquearRequest registra a ação judicial.
type BloquearRequest struct {
	Processo string `json:"processo" validate:"required,max=120"`
	Motivo   string `json:"motivo" validate:"max=2000"`
}

// Bloquear registra o bloqueio judicial.
func (h *Handlers) Bloquear(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req BloquearRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Bloquear(r.Context(), identity, id, req.Processo, req.Motivo)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, d)
}

// LiberarRequest traz o motivo da liberação.
type LiberarRequest struct {
	Motivo string `json:"motivo" validate:"required,max=2000"`
}

// Liberar encerra o bloqueio.
func (h *Handlers) Liberar(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req LiberarRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Liberar(r.Context(), identity, id, req.Motivo)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

// LocalRequest é o corpo de um local.
type LocalRequest struct {
	ParentID *uuid.UUID `json:"parent_id"`
	Tipo     string     `json:"tipo" validate:"required,max=16"`
	Nome     string     `json:"nome" validate:"required,max=120"`
	Ativo    *bool      `json:"ativo"`
}

func (req LocalRequest) local() domain.Local {
	ativo := req.Ativo == nil || *req.Ativo
	return domain.Local{ParentID: req.ParentID, Tipo: domain.TipoLocal(req.Tipo), Nome: req.Nome, Ativo: ativo}
}

// ListLocais devolve a hierarquia física.
func (h *Handlers) ListLocais(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.ListLocais(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// CreateLocal cria um local.
func (h *Handlers) CreateLocal(w http.ResponseWriter, r *http.Request) {
	var req LocalRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.SalvarLocal(r.Context(), identity, nil, req.local())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, l)
}

// UpdateLocal altera um local.
func (h *Handlers) UpdateLocal(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req LocalRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.SalvarLocal(r.Context(), identity, &id, req.local())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}
