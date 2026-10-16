package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// registerFase4 monta as rotas do recolhimento ao Arquivo Público (ADR 035).
func (h *Handlers) registerFase4(r chi.Router, ler, gerir, central func(http.Handler) http.Handler) {
	r.With(ler).Get("/recolhimentos", h.Recolhimentos)
	r.With(ler).Get("/recolhimentos/{id}", h.Recolhimento)
	r.With(gerir).Post("/recolhimentos", h.CriarRecolhimento)
	r.With(gerir).Post("/recolhimentos/{id}/cancelar", h.CancelarRecolhimento)
	r.With(central).Post("/recolhimentos/{id}/receber", h.ReceberRecolhimento)
	r.With(central).Post("/recolhimentos/{id}/devolver", h.DevolverRecolhimento)
}

// Recolhimentos pagina os recolhimentos no escopo.
func (h *Handlers) Recolhimentos(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Recolhimentos(r.Context(), identity, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Recolhimento devolve o recolhimento com as unidades.
func (h *Handlers) Recolhimento(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	rc, err := h.svc.Recolhimento(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, rc)
}

// RecolhimentoRequest é o pedido da origem.
type RecolhimentoRequest struct {
	OrigemID   uuid.UUID      `json:"origem_id" validate:"required"`
	Motivo     string         `json:"motivo" validate:"required,max=20"`
	Fundamento string         `json:"fundamento" validate:"max=500"`
	Preparo    domain.Preparo `json:"preparo"`
	Historico  string         `json:"historico" validate:"max=8000"`
	Observacao string         `json:"observacao" validate:"max=1000"`
	Unidades   []uuid.UUID    `json:"unidades" validate:"required,min=1,max=2000"`
}

// CriarRecolhimento relaciona as unidades e envia ao Arquivo Público.
func (h *Handlers) CriarRecolhimento(w http.ResponseWriter, r *http.Request) {
	var req RecolhimentoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	rc, err := h.svc.CriarRecolhimento(r.Context(), identity, application.NovoRecolhimento{OrigemID: req.OrigemID, Motivo: req.Motivo,
		Fundamento: req.Fundamento, Preparo: req.Preparo, Historico: req.Historico, Observacao: req.Observacao, Unidades: req.Unidades})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, rc)
}

// EntradaRequest é o registro de entrada no Arquivo Público.
type EntradaRequest struct {
	LocalID *uuid.UUID `json:"local_id"`
}

// ReceberRecolhimento registra a entrada no Arquivo Público.
func (h *Handlers) ReceberRecolhimento(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req EntradaRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	rc, err := h.svc.ReceberRecolhimento(r.Context(), identity, id, application.Entrada{LocalID: req.LocalID})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, rc)
}

// DevolverRecolhimento devolve o recolhimento à origem.
func (h *Handlers) DevolverRecolhimento(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req MotivoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	rc, err := h.svc.DevolverRecolhimento(r.Context(), identity, id, req.Motivo)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, rc)
}

// CancelarRecolhimento: a origem desiste antes da entrada.
func (h *Handlers) CancelarRecolhimento(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	rc, err := h.svc.CancelarRecolhimento(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, rc)
}
