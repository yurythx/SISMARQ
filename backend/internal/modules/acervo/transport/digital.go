package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"

	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// registerDigital monta as rotas da digitalização (ADR 037).
func (h *Handlers) registerDigital(r chi.Router, ler, gerir, central func(http.Handler) http.Handler) {
	r.With(gerir).Post("/unidades/{id}/digitais/upload", h.UploadDigital)
	r.With(gerir).Post("/unidades/{id}/digitais", h.RegistrarDigital)
	r.With(ler).Get("/digitais", h.Digitais)
	r.With(ler).Get("/digitais/{id}", h.Digital)
	r.With(ler).Get("/digitais/{id}/arquivo", h.ArquivoDigital)
	r.With(central).Post("/digitais/{id}/verificar", h.digitalAcao(h.svc.VerificarAgora))
	r.With(central).Post("/digitais/{id}/reprocessar", h.digitalAcao(h.svc.Reprocessar))
}

// UploadDigitalRequest pede a URL de envio do PDF digitalizado.
type UploadDigitalRequest struct {
	Filename string `json:"filename" validate:"required,max=255"`
}

// UploadDigital emite a URL de upload direto ao MinIO.
func (h *Handlers) UploadDigital(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req UploadDigitalRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	t, err := h.svc.UploadDigital(r.Context(), identity, id, req.Filename)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, t)
}

// DigitalRequest registra o PDF enviado, com os metadados do Decreto.
type DigitalRequest struct {
	ObjectKey      string `json:"object_key" validate:"required,max=300"`
	Descricao      string `json:"descricao" validate:"max=500"`
	TipoDoc        string `json:"tipo_documento" validate:"required,max=10"`
	Equipamento    string `json:"equipamento" validate:"required,max=200"`
	DPI            int    `json:"dpi_declarado" validate:"required"`
	ResponsavelCPF string `json:"responsavel_cpf" validate:"max=20"`
}

// RegistrarDigital põe o PDF enviado na fila do pipeline.
func (h *Handlers) RegistrarDigital(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req DigitalRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.RegistrarDigital(r.Context(), identity, id, application.NovoEnvio{ObjectKey: req.ObjectKey, Descricao: req.Descricao,
		TipoDoc: req.TipoDoc, Equipamento: req.Equipamento, DPI: req.DPI, ResponsavelCPF: req.ResponsavelCPF})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, d)
}

// Digitais pagina os representantes (?unidade_id=, ?falhas=true,
// ?publicados=true).

func (h *Handlers) Digitais(w http.ResponseWriter, r *http.Request) {
	unidade, err := httputil.OptionalUUIDQuery(r, "unidade_id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Digitais(r.Context(), identity, domain.FiltroDigital{UnidadeID: unidade,
		SoFalhas: httputil.Query(r, "falhas", 5) == "true", Publicados: httputil.Query(r, "publicados", 5) == "true"}, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Digital devolve o representante com as verificações.
func (h *Handlers) Digital(w http.ResponseWriter, r *http.Request) {
	h.digitalAcao(h.svc.Digital)(w, r)
}

// digitalAcao lê o {id} e responde o representante devolvido pela ação
// (ver, verificar, reprocessar).
func (h *Handlers) digitalAcao(fn func(ctx context.Context, identity auth.Identity, id uuid.UUID) (application.DetalheDigital, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputil.UUIDParam(r, "id")
		if err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		d, err := fn(r.Context(), identity, id)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		httputil.WriteOK(w, d)
	}
}

// intervalo lê um único "Range: bytes=…" contra o tamanho do arquivo.
// ok=false: sem Range (arquivo inteiro); erro: intervalo fora (416).
func intervalo(h string, tamanho int64) (inicio, fim int64, ok bool, err error) {
	h = strings.TrimSpace(h)
	if h == "" {
		return 0, tamanho - 1, false, nil
	}
	spec, achou := strings.CutPrefix(h, "bytes=")
	if !achou || strings.Contains(spec, ",") {
		return 0, 0, false, errRange
	}
	a, b, _ := strings.Cut(spec, "-")
	switch {
	case a == "" && b != "":
		// Sufixo: os últimos N bytes.
		n, perr := strconv.ParseInt(b, 10, 64)
		if perr != nil || n <= 0 {
			return 0, 0, false, errRange
		}
		return max(0, tamanho-n), tamanho - 1, true, nil
	case a != "":
		inicio, perr := strconv.ParseInt(a, 10, 64)
		if perr != nil || inicio < 0 || inicio >= tamanho {
			return 0, 0, false, errRange
		}
		fim = tamanho - 1
		if b != "" {
			f, perr := strconv.ParseInt(b, 10, 64)
			if perr != nil || f < inicio {
				return 0, 0, false, errRange
			}
			fim = min(f, tamanho-1)
		}
		return inicio, fim, true, nil
	}
	return 0, 0, false, errRange
}

var errRange = errors.New("range inválido")

// ArquivoDigital entrega o PDF/A em fluxo: com Range, só o trecho pedido
// (206) — o visualizador (pdf.js) busca as páginas visíveis de um processo
// de centenas de páginas sem baixá-lo inteiro.
func (h *Handlers) ArquivoDigital(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	d, err := h.svc.AbrirArquivo(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.enviarPDF(w, r, d, "private, max-age=300")
}

// enviarPDF responde o PDF/A inteiro ou o trecho do Range (206). Devolve
// true no pedido sem Range (a abertura do documento ou o download; os
// trechos que o visualizador pede depois não contam).
func (h *Handlers) enviarPDF(w http.ResponseWriter, r *http.Request, d domain.Digital, cache string) bool {
	tamanho := d.PDFABytes
	inicio, fim, parcial, err := intervalo(r.Header.Get("Range"), tamanho)
	if err != nil {
		w.Header().Set("Content-Range", "bytes */"+strconv.FormatInt(tamanho, 10))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return false
	}
	n := fim - inicio + 1
	comprimento := int64(0)
	if parcial {
		comprimento = n
	}
	rc, err := h.svc.LerArquivo(r.Context(), d, inicio, comprimento)
	if err != nil {
		h.fail(w, r, err)
		return false
	}
	defer rc.Close()
	hd := w.Header()
	hd.Set("Content-Type", "application/pdf")
	hd.Set("Accept-Ranges", "bytes")
	hd.Set("ETag", `"`+d.PDFASHA256+`"`)
	hd.Set("Cache-Control", cache)
	hd.Set("Content-Disposition", `inline; filename="`+d.Numero+`.pdf"`)
	hd.Set("Content-Length", strconv.FormatInt(n, 10))
	status := http.StatusOK
	if parcial {
		hd.Set("Content-Range", "bytes "+strconv.FormatInt(inicio, 10)+"-"+strconv.FormatInt(fim, 10)+"/"+strconv.FormatInt(tamanho, 10))
		status = http.StatusPartialContent
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, io.LimitReader(rc, n)); err != nil {
		h.logger.WarnContext(r.Context(), "acervo: envio do PDF/A interrompido", "digital", d.Numero, "error", err)
	}
	return !parcial
}

