package transport

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// Fase 6 (ADR 038): busca, consulta pública e exposições.

// registerAcesso monta as rotas autenticadas da Fase 6.
func (h *Handlers) registerAcesso(r chi.Router, ler, central func(http.Handler) http.Handler) {
	r.With(ler).Get("/busca", h.Buscar)
	r.With(ler).Post("/digitais/{id}/publicacao", h.publicacaoDigital(true))
	r.With(ler).Delete("/digitais/{id}/publicacao", h.publicacaoDigital(false))
	r.With(central).Get("/exposicoes", h.Exposicoes)
	r.With(central).Post("/exposicoes", h.CriarExposicao)
	r.With(central).Get("/exposicoes/{id}", h.Exposicao)
	r.With(central).Put("/exposicoes/{id}", h.AlterarExposicao)
	r.With(central).Post("/exposicoes/{id}/publicacao", h.publicacaoExposicao(true))
	r.With(central).Delete("/exposicoes/{id}/publicacao", h.publicacaoExposicao(false))
}

// registerAcessoPublico monta a consulta sem login: só o ostensivo.
func (h *Handlers) registerAcessoPublico(r chi.Router) {
	r.Get("/acervo/publico/arquivo", h.ArquivoPublico)
	r.Get("/acervo/publico/busca", h.BuscarPublico)
	r.Get("/acervo/publico/unidades/{id}", h.FichaPublica)
	r.Get("/acervo/publico/digitais/{id}/arquivo", h.ArquivoPublicoDigital)
	r.Get("/acervo/publico/exposicoes", h.ExposicoesPublicas)
	r.Get("/acervo/publico/exposicoes/{slug}", h.ExposicaoPublica)
}

// BuscaResponse é o resultado da busca: a página, as facetas e de onde
// veio (o índice ou, com ele fora, o banco, sem facetas).
type BuscaResponse struct {
	Items   []domain.ItemBusca         `json:"items"`
	Meta    pagination.Meta            `json:"meta"`
	Facetas map[string][]domain.Faceta `json:"facetas"`
	Origem  string                     `json:"origem"`
	Ms      int64                      `json:"ms"`
}

func consulta(r *http.Request) domain.Consulta {
	p := httputil.Page(r, application.BuscaMaxPorPagina)
	c := domain.Consulta{Q: httputil.Query(r, "q", 400), Filtros: map[string]string{}, SoDigital: httputil.Query(r, "digital", 5) == "true",
		Pagina: p.Page, PorPagina: p.PageSize}
	for _, k := range domain.FacetasInternas {
		if v := httputil.Query(r, k, 200); v != "" {
			c.Filtros[k] = v
		}
	}
	return c
}

func responderBusca(w http.ResponseWriter, c domain.Consulta, res domain.ResultadoBusca) {
	httputil.WriteOK(w, BuscaResponse{Items: res.Itens, Meta: pagination.NewMeta(pagination.Params{Page: c.Pagina, PageSize: c.PorPagina}, res.Total),
		Facetas: res.Facetas, Origem: res.Origem, Ms: res.Ms})
}

// Buscar é a busca interna (?q=, facetas, ?digital=true), no escopo de
// quem pede.
func (h *Handlers) Buscar(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	c := consulta(r)
	res, err := h.svc.Buscar(r.Context(), identity, c)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	responderBusca(w, c, res)
}

// BuscarPublico é a busca sem login, só no ostensivo.
func (h *Handlers) BuscarPublico(w http.ResponseWriter, r *http.Request) {
	c := consulta(r)
	res, err := h.svc.BuscarPublico(r.Context(), c)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=30")
	responderBusca(w, c, res)
}

func (h *Handlers) publicacaoDigital(publicar bool) http.HandlerFunc {
	return h.digitalAcao(func(ctx context.Context, identity auth.Identity, id uuid.UUID) (application.DetalheDigital, error) {
		return h.svc.PublicarDigital(ctx, identity, id, publicar)
	})
}

// ArquivoPublico é a página do Arquivo Público.
func (h *Handlers) ArquivoPublico(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.ArquivoPublico(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WriteOK(w, p)
}

// FichaPublica é a descrição ostensiva da unidade.
func (h *Handlers) FichaPublica(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	f, err := h.svc.FichaPublica(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WriteOK(w, f)
}

// ArquivoPublicoDigital entrega o PDF/A liberado, em fluxo (Range).
func (h *Handlers) ArquivoPublicoDigital(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	d, err := h.svc.DocumentoPublico(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if h.enviarPDF(w, r, d, "public, max-age=3600") {
		h.svc.ContarDocumento(r.Context())
	}
}

// ExposicaoRequest é a exposição com os documentos, na ordem.
type ExposicaoRequest struct {
	Titulo string             `json:"titulo" validate:"required,max=200"`
	Resumo string             `json:"resumo" validate:"max=500"`
	Texto  string             `json:"texto" validate:"max=20000"`
	Itens  []ItemExposicaoReq `json:"itens" validate:"max=100,dive"`
}

// ItemExposicaoReq é um documento da exposição.
type ItemExposicaoReq struct {
	DigitalID uuid.UUID `json:"digital_id" validate:"required"`
	Legenda   string    `json:"legenda" validate:"max=1000"`
}

func (req ExposicaoRequest) nova() application.NovaExposicao {
	n := application.NovaExposicao{Titulo: req.Titulo, Resumo: req.Resumo, Texto: req.Texto, Itens: make([]application.NovoItem, len(req.Itens))}
	for i, it := range req.Itens {
		n.Itens[i] = application.NovoItem{DigitalID: it.DigitalID, Legenda: it.Legenda}
	}
	return n
}

// Exposicoes pagina as exposições (Arquivo Público).
func (h *Handlers) Exposicoes(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Exposicoes(r.Context(), identity, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// CriarExposicao cria a exposição (rascunho).
func (h *Handlers) CriarExposicao(w http.ResponseWriter, r *http.Request) {
	var req ExposicaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	x, err := h.svc.SalvarExposicao(r.Context(), identity, nil, req.nova())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, x)
}

// Exposicao devolve a exposição com todos os documentos.
func (h *Handlers) Exposicao(w http.ResponseWriter, r *http.Request) {
	h.exposicaoAcao(h.svc.Exposicao)(w, r)
}

// AlterarExposicao troca o texto e os documentos.
func (h *Handlers) AlterarExposicao(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req ExposicaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	x, err := h.svc.SalvarExposicao(r.Context(), identity, &id, req.nova())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, x)
}

func (h *Handlers) publicacaoExposicao(publicar bool) http.HandlerFunc {
	return h.exposicaoAcao(func(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Exposicao, error) {
		return h.svc.PublicarExposicao(ctx, identity, id, publicar)
	})
}

func (h *Handlers) exposicaoAcao(fn func(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Exposicao, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputil.UUIDParam(r, "id")
		if err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		x, err := fn(r.Context(), identity, id)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		httputil.WriteOK(w, x)
	}
}

// ExposicoesPublicas pagina as exposições publicadas.
func (h *Handlers) ExposicoesPublicas(w http.ResponseWriter, r *http.Request) {
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.ExposicoesPublicas(r.Context(), p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WritePage(w, items, p, total)
}

// ExposicaoPublica é a visita à exposição publicada.
func (h *Handlers) ExposicaoPublica(w http.ResponseWriter, r *http.Request) {
	x, err := h.svc.ExposicaoPublica(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=60")
	httputil.WriteOK(w, x)
}
