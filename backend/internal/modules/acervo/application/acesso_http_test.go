package application_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/modules/acervo/transport"
)

func publico(s *application.Service) http.Handler {
	r := chi.NewRouter()
	transport.NewHandlers(s, logger, 100).RegisterPublicRoutes(r)
	return r
}

func TestHandlersAcesso(t *testing.T) {
	a := novoAcesso(t)
	ctx := context.Background()
	marca := "Feira" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	u, d := a.pronta(t, "Planta da "+marca)
	h, arq, pub := router(a.s, a.quem), router(a.s, arquivoPublico(uuid.Nil)), publico(a.s)
	dg := "/acervo/digitais/" + d.ID.String()

	if rec := do(h, http.MethodPost, dg+"/publicacao", ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"publicado_em":"`) {
		t.Fatalf("liberar: %d %s", rec.Code, rec.Body)
	}
	rec := do(arq, http.MethodPost, "/acervo/exposicoes", `{"titulo":"Feiras `+marca+`","itens":[{"digital_id":"`+d.ID.String()+`","legenda":"Planta"}]}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("exposição: %d %s", rec.Code, rec.Body)
	}
	xid, slug := campo(t, rec.Body.String(), "id"), campo(t, rec.Body.String(), "slug")
	x := "/acervo/exposicoes/" + xid

	casos := []struct {
		h                  http.Handler
		method, path, body string
		want               int
	}{
		{h, http.MethodGet, "/acervo/busca?q=" + marca + "&fase=CORRENTE&digital=true&page=1&page_size=10", "", http.StatusOK},
		{h, http.MethodGet, "/acervo/busca?ano=x", "", http.StatusUnprocessableEntity},
		{arq, http.MethodGet, "/acervo/digitais?publicados=true", "", http.StatusOK},

		{h, http.MethodPost, "/acervo/digitais/x/publicacao", "", http.StatusBadRequest},
		{h, http.MethodPost, dg + "/publicacao", "", http.StatusConflict},
		{h, http.MethodDelete, "/acervo/digitais/" + uuid.NewString() + "/publicacao", "", http.StatusNotFound},
		{arq, http.MethodGet, "/acervo/exposicoes", "", http.StatusOK},
		{h, http.MethodGet, "/acervo/exposicoes", "", http.StatusForbidden},
		{arq, http.MethodPost, "/acervo/exposicoes", `{`, http.StatusBadRequest},
		{arq, http.MethodPost, "/acervo/exposicoes", `{"titulo":"!!!"}`, http.StatusUnprocessableEntity},
		{arq, http.MethodGet, x, "", http.StatusOK},
		{arq, http.MethodGet, "/acervo/exposicoes/x", "", http.StatusBadRequest},
		{arq, http.MethodGet, "/acervo/exposicoes/" + uuid.NewString(), "", http.StatusNotFound},
		{arq, http.MethodPut, "/acervo/exposicoes/x", `{"titulo":"a"}`, http.StatusBadRequest},
		{arq, http.MethodPut, x, `{`, http.StatusBadRequest},
		{arq, http.MethodPut, "/acervo/exposicoes/" + uuid.NewString(), `{"titulo":"a"}`, http.StatusNotFound},
		{arq, http.MethodPut, x, `{"titulo":"Feiras e mercados","itens":[{"digital_id":"` + d.ID.String() + `"}]}`, http.StatusOK},
		{arq, http.MethodPost, x + "/publicacao", "", http.StatusOK},
		{arq, http.MethodPost, x + "/publicacao", "", http.StatusConflict},
		{arq, http.MethodPost, "/acervo/exposicoes/x/publicacao", "", http.StatusBadRequest},

		{pub, http.MethodGet, "/acervo/publico/arquivo", "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/busca?q=" + marca, "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/busca?nivel_acesso=RESTRITO", "", http.StatusUnprocessableEntity},
		{pub, http.MethodGet, "/acervo/publico/unidades/" + u.ID.String(), "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/unidades/x", "", http.StatusBadRequest},
		{pub, http.MethodGet, "/acervo/publico/unidades/" + uuid.NewString(), "", http.StatusNotFound},
		{pub, http.MethodGet, "/acervo/publico/digitais/x/arquivo", "", http.StatusBadRequest},
		{pub, http.MethodGet, "/acervo/publico/digitais/" + uuid.NewString() + "/arquivo", "", http.StatusNotFound},
		{pub, http.MethodGet, "/acervo/publico/exposicoes", "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/exposicoes/" + slug, "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/exposicoes/nao-existe", "", http.StatusNotFound},
		{arq, http.MethodDelete, x + "/publicacao", "", http.StatusOK},
		{h, http.MethodDelete, dg + "/publicacao", "", http.StatusOK},
	}
	for _, c := range casos {
		if rec := do(c.h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}

	// A busca: a página, as facetas (vazias no banco) e a origem.
	body := do(pub, http.MethodGet, "/acervo/publico/busca?q="+marca, "").Body.String()
	if !strings.Contains(body, `"origem":"banco"`) || !strings.Contains(body, `"total_items":1`) || !strings.Contains(body, `"facetas":{}`) ||
		strings.Contains(body, "nivel_acesso") {
		t.Fatalf("resposta da busca: %s", body)
	}

	// O documento liberado em fluxo, sem login: inteiro (conta), trecho e 416.
	if _, err := a.s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	antes := consultas(t, a.e, "DOCUMENTO")
	arquivo := "/acervo/publico/digitais/" + d.ID.String() + "/arquivo"
	rec = do(pub, http.MethodGet, arquivo, "")
	if rec.Code != http.StatusOK || rec.Header().Get("Cache-Control") != "public, max-age=3600" || !strings.HasPrefix(rec.Body.String(), "%PDF") {
		t.Fatalf("documento: %d %v", rec.Code, rec.Header())
	}
	trecho := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, arquivo, nil)
	req.Header.Set("Range", "bytes=0-3")
	pub.ServeHTTP(trecho, req)
	if trecho.Code != http.StatusPartialContent || trecho.Body.String() != "%PDF" {
		t.Fatalf("trecho: %d %q", trecho.Code, trecho.Body)
	}
	fora := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, arquivo, nil)
	req.Header.Set("Range", "bytes=999999-")
	pub.ServeHTTP(fora, req)
	if fora.Code != http.StatusRequestedRangeNotSatisfiable {
		t.Fatalf("fora do arquivo: %d", fora.Code)
	}
	if depois := consultas(t, a.e, "DOCUMENTO"); depois != antes+1 {
		t.Fatalf("só a abertura conta: %d -> %d", antes, depois)
	}

	// Banco fora: 500 sem vazar o erro.
	down := func() *application.Service { return a.e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}) }
	for _, c := range []struct {
		h    http.Handler
		path string
	}{
		{router(down(), a.quem), "/acervo/busca?q=x"}, {router(down(), arquivoPublico(uuid.Nil)), "/acervo/exposicoes"},
		{publico(down()), "/acervo/publico/arquivo"}, {publico(down()), "/acervo/publico/busca"}, {publico(down()), "/acervo/publico/exposicoes"},
	} {
		if rec := do(c.h, http.MethodGet, c.path, ""); rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), errBoom.Error()) {

			t.Errorf("%s com o banco fora: %d %s", c.path, rec.Code, rec.Body)
		}
	}
}
