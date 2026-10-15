package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// maxPlanilha limita o corpo da importação (até 5.000 linhas).
const maxPlanilha = 10 << 20

// registerFase2 monta as rotas da Fase 2 (ADR 033).
func (h *Handlers) registerFase2(r chi.Router, ler, gerir func(http.Handler) http.Handler) {
	r.With(gerir).Post("/unidades/{id}/emprestimos", h.Emprestar)
	r.With(gerir).Post("/emprestimos/{id}/devolver", h.Devolver)
	r.With(ler).Get("/emprestimos", h.Emprestimos)
	r.With(ler).Get("/guias", h.Guias)
	r.With(ler).Get("/guias/{id}", h.Guia)
	r.With(gerir).Post("/guias", h.CriarGuia)
	r.With(gerir).Post("/guias/{id}/receber", h.guiaAcao(true))
	r.With(gerir).Post("/guias/{id}/cancelar", h.guiaAcao(false))
	r.With(gerir).Get("/importacao/modelo", h.ModeloImportacao)
	r.With(gerir).Post("/importacao/simular", h.importacao(false))
	r.With(gerir).Post("/importacao/aplicar", h.importacao(true))
}

// EmprestarRequest é o pedido de empréstimo.
type EmprestarRequest struct {
	SolicitanteID uuid.UUID `json:"solicitante_id" validate:"required"`
	Finalidade    string    `json:"finalidade" validate:"required,max=500"`
	DevolverAte   string    `json:"devolver_ate" validate:"required"`
}

// Emprestar registra a saída da unidade.
func (h *Handlers) Emprestar(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req EmprestarRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	ate, err := data("devolver_ate", req.DevolverAte)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Emprestar(r.Context(), identity, id, req.SolicitanteID, req.Finalidade, ate)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, d)
}

// DevolverRequest traz a observação da devolução.
type DevolverRequest struct {
	Observacao string `json:"observacao" validate:"max=500"`
}

// Devolver registra a volta da unidade.
func (h *Handlers) Devolver(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req DevolverRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.Devolver(r.Context(), identity, id, req.Observacao)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, d)
}

// Emprestimos pagina os empréstimos em aberto (?atrasados=true).
func (h *Handlers) Emprestimos(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Emprestimos(r.Context(), identity, httputil.Query(r, "atrasados", 5) == "true", p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Guias pagina as guias de transferência.
func (h *Handlers) Guias(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Guias(r.Context(), identity, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Guia devolve a guia com as unidades.
func (h *Handlers) Guia(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	g, err := h.svc.Guia(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, g)
}

// GuiaRequest é a remessa ao arquivo intermediário.
type GuiaRequest struct {
	OrigemID   uuid.UUID   `json:"origem_id" validate:"required"`
	Unidades   []uuid.UUID `json:"unidades" validate:"required,min=1,max=500"`
	Observacao string      `json:"observacao" validate:"max=1000"`
}

// CriarGuia gera a guia de transferência.
func (h *Handlers) CriarGuia(w http.ResponseWriter, r *http.Request) {
	var req GuiaRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	g, err := h.svc.CriarGuia(r.Context(), identity, req.OrigemID, req.Unidades, req.Observacao)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, g)
}

// guiaAcao recebe (aceite) ou cancela a guia.
func (h *Handlers) guiaAcao(receber bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputil.UUIDParam(r, "id")
		if err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		acao := h.svc.CancelarGuia
		if receber {
			acao = h.svc.ReceberGuia
		}
		g, err := acao(r.Context(), identity, id)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		httputil.WriteOK(w, g)
	}
}

// ModeloImportacao devolve a planilha-modelo (CSV com ";", como o Excel).
func (h *Handlers) ModeloImportacao(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="modelo-importacao-acervo.csv"`)
	_, _ = w.Write([]byte("\uFEFF" + strings.Join(domain.ColunasImportacao, ";") + "\n" +
		"CAIXA;2.0.01.00.01;Organogramas de 2015;SEMAD;01/01/2015;31/12/2015;2;caixas-arquivo;;;PUBLICO;;;\n"))
}

// ImportacaoRequest traz a planilha (texto do CSV) e, para aplicar, o hash
// devolvido pela simulação.
type ImportacaoRequest struct {
	Conteudo string `json:"conteudo" validate:"required"`
	Hash     string `json:"hash" validate:"omitempty,len=64"`
}

func (h *Handlers) importacao(aplicar bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxPlanilha)
		var req ImportacaoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.fail(w, r, apperrors.BadRequest("corpo inválido ou acima de 10 MB"))
			return
		}
		if err := httputil.Validate(req); err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		out, err := h.svc.Importar(r.Context(), identity, req.Conteudo, req.Hash, aplicar)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		httputil.WriteOK(w, out)
	}
}
