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

// registerNormas monta a gestão das instruções normativas (ADR 036).
func (h *Handlers) registerNormas(r chi.Router, central func(http.Handler) http.Handler) {
	r.With(central).Get("/normas", h.Normas)
	r.With(central).Post("/normas", h.CriarNorma)
	r.With(central).Get("/normas/{id}", h.Norma)
	r.With(central).Put("/normas/{id}/rascunho", h.SalvarRascunho)
	r.With(central).Delete("/normas/{id}/rascunho", h.DescartarRascunho)
	r.With(central).Post("/normas/{id}/publicar", h.PublicarNorma)
	r.With(central).Post("/normas/{id}/revogar", h.RevogarNorma)
}

// RegisterPublicRoutes: consulta anônima das instruções normativas
// publicadas, com o histórico das redações (art. 54). O módulo aplica o
// limite por IP.
func (h *Handlers) RegisterPublicRoutes(r chi.Router) {
	r.Get("/acervo/publico/normas", h.NormasPublicas)
	r.Get("/acervo/publico/normas/{id}", h.NormaPublica)
}

// NormaRequest é a criação da norma com a primeira redação.
type NormaRequest struct {
	Numero string `json:"numero" validate:"required,max=9"`
	RedacaoRequest
}

// RedacaoRequest é uma redação (a primeira ou a alteração).
type RedacaoRequest struct {
	Ementa string `json:"ementa" validate:"required,max=500"`
	Texto  string `json:"texto" validate:"required,max=100000"`
	Motivo string `json:"motivo" validate:"max=500"`
}

func (r RedacaoRequest) redacao() domain.Redacao {
	return domain.Redacao{Ementa: r.Ementa, Texto: r.Texto, Motivo: r.Motivo}
}

// Normas pagina as normas para a gestão (?q=).
func (h *Handlers) Normas(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Normas(r.Context(), identity, httputil.Query(r, "q", 100), p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// CriarNorma abre a norma com a primeira redação em rascunho.
func (h *Handlers) CriarNorma(w http.ResponseWriter, r *http.Request) {
	var req NormaRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.CriarNorma(r.Context(), identity, req.Numero, req.redacao())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, n)
}

// Norma devolve a norma com todas as redações.
func (h *Handlers) Norma(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.Norma(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, n)
}

// SalvarRascunho grava a redação em rascunho (nova ou alterada).
func (h *Handlers) SalvarRascunho(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req RedacaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.SalvarRascunho(r.Context(), identity, id, req.redacao())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, n)
}

// DescartarRascunho apaga o rascunho (e a norma nunca publicada).
func (h *Handlers) DescartarRascunho(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.DescartarRascunho(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if n.ID == uuid.Nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	httputil.WriteOK(w, n)
}

// PublicacaoRequest é a publicação da redação.
type PublicacaoRequest struct {
	PublicadaEm string `json:"publicada_em" validate:"required"`
	Diario      string `json:"diario" validate:"max=200"`
}

// PublicarNorma publica o rascunho.
func (h *Handlers) PublicarNorma(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req PublicacaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	em, err := data("publicada_em", req.PublicadaEm)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.PublicarNorma(r.Context(), identity, id, application.Publicacao{Em: em, Diario: req.Diario})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, n)
}

// RevogacaoRequest é a revogação da norma.
type RevogacaoRequest struct {
	RevogadaEm string `json:"revogada_em" validate:"required"`
	Ato        string `json:"ato" validate:"required,max=500"`
}

// RevogarNorma encerra a vigência.
func (h *Handlers) RevogarNorma(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req RevogacaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	em, err := data("revogada_em", req.RevogadaEm)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	n, err := h.svc.RevogarNorma(r.Context(), identity, id, em, req.Ato)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, n)
}

// NormasPublicas pagina as normas publicadas (?q=).
func (h *Handlers) NormasPublicas(w http.ResponseWriter, r *http.Request) {
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.NormasPublicas(r.Context(), httputil.Query(r, "q", 100), p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WritePage(w, items, p, total)
}

// NormaPublica devolve a norma publicada com o histórico das redações.
func (h *Handlers) NormaPublica(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	n, err := h.svc.NormaPublica(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WriteOK(w, n)
}
