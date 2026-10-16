package application_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/modules/acervo/transport"
)

func TestHandlersNormasEPrivados(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2026-10-15")
	s := e.real()
	h := router(s, admin)
	pub := chi.NewRouter()
	transport.NewHandlers(s, logger, 100).RegisterPublicRoutes(pub)

	rec := do(h, http.MethodPost, "/acervo/normas", `{"numero":"`+numeroNum()+`","ementa":"Recolhimento","texto":"Art. 1º"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("norma: %d %s", rec.Code, rec.Body)
	}
	id := campo(t, rec.Body.String(), "id")
	descartavel := campo(t, do(h, http.MethodPost, "/acervo/normas", `{"numero":"`+numeroNum()+`","ementa":"x","texto":"y"}`).Body.String(), "id")
	n := "/acervo/normas/" + id

	rec = do(h, http.MethodPost, "/acervo/privados", `{"titular":"Família Silva","tipo_titular":"PF","descricao":"Fotografias"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("privado: %d %s", rec.Code, rec.Body)
	}
	pid := campo(t, rec.Body.String(), "id")
	p := "/acervo/privados/" + pid

	casos := []struct {
		h                  http.Handler
		method, path, body string
		want               int
	}{
		{h, http.MethodGet, "/acervo/normas?q=Recolhimento", "", http.StatusOK},
		{h, http.MethodPost, "/acervo/normas", `{`, http.StatusBadRequest},
		{h, http.MethodPost, "/acervo/normas", `{"numero":"x","ementa":"e","texto":"t"}`, http.StatusUnprocessableEntity},
		{h, http.MethodGet, n, "", http.StatusOK},
		{h, http.MethodGet, "/acervo/normas/x", "", http.StatusBadRequest},
		{h, http.MethodGet, "/acervo/normas/" + uuid.NewString(), "", http.StatusNotFound},
		{h, http.MethodPut, n + "/rascunho", `{`, http.StatusBadRequest},
		{h, http.MethodPut, "/acervo/normas/x/rascunho", `{}`, http.StatusBadRequest},
		{h, http.MethodPost, n + "/revogar", `{"revogada_em":"2026-10-15","ato":"IN"}`, http.StatusConflict},
		{h, http.MethodPut, n + "/rascunho", `{"ementa":"Recolhimento","texto":"Art. 1º Novo","motivo":"Redação original"}`, http.StatusOK},
		{h, http.MethodPost, n + "/publicar", `{`, http.StatusBadRequest},
		{h, http.MethodPost, "/acervo/normas/x/publicar", `{}`, http.StatusBadRequest},
		{h, http.MethodPost, n + "/publicar", `{"publicada_em":"ontem"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, n + "/publicar", `{"publicada_em":"2026-10-15","diario":"DIORONDON-e"}`, http.StatusOK},
		{h, http.MethodPost, n + "/publicar", `{"publicada_em":"2026-10-15"}`, http.StatusConflict},
		{h, http.MethodDelete, "/acervo/normas/x/rascunho", ``, http.StatusBadRequest},
		{h, http.MethodDelete, n + "/rascunho", ``, http.StatusConflict},
		{h, http.MethodPut, n + "/rascunho", `{"ementa":"Recolhimento","texto":"Art. 1º Alterado","motivo":"IN 2/2026"}`, http.StatusOK},
		{h, http.MethodDelete, n + "/rascunho", ``, http.StatusOK},
		{h, http.MethodDelete, "/acervo/normas/" + descartavel + "/rascunho", ``, http.StatusNoContent},
		{h, http.MethodPost, n + "/revogar", `{`, http.StatusBadRequest},
		{h, http.MethodPost, "/acervo/normas/x/revogar", `{}`, http.StatusBadRequest},
		{h, http.MethodPost, n + "/revogar", `{"revogada_em":"x","ato":"IN 3"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, n + "/revogar", `{"revogada_em":"2026-10-15","ato":"IN 3/2026"}`, http.StatusOK},
		{h, http.MethodPut, n + "/rascunho", `{"ementa":"e","texto":"t","motivo":"m"}`, http.StatusConflict},
		{pub, http.MethodGet, "/acervo/publico/normas?q=Alterado", "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/" + "normas/" + id, "", http.StatusOK},
		{pub, http.MethodGet, "/acervo/publico/normas/x", "", http.StatusBadRequest},
		{pub, http.MethodGet, "/acervo/publico/normas/" + uuid.NewString(), "", http.StatusNotFound},

		{h, http.MethodGet, "/acervo/privados?q=Silva", "", http.StatusOK},
		{h, http.MethodPost, "/acervo/privados", `{`, http.StatusBadRequest},
		{h, http.MethodPost, "/acervo/privados", `{"titular":"X","tipo_titular":"XX","descricao":"d"}`, http.StatusUnprocessableEntity},
		{h, http.MethodGet, p, "", http.StatusOK},
		{h, http.MethodGet, "/acervo/privados/x", "", http.StatusBadRequest},
		{h, http.MethodPut, p, `{`, http.StatusBadRequest},
		{h, http.MethodPut, p, `{"titular":"Família Silva","tipo_titular":"PF","descricao":"Fotografias e cartas"}`, http.StatusOK},
		{h, http.MethodPost, p + "/parecer", `{"parecer_comissao":"C","parecer_texto":"T","parecer_favoravel":true,"parecer_em":"x"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, p + "/parecer", `{"parecer_comissao":"C","parecer_texto":"T","parecer_favoravel":true,"parecer_em":"2026-09-01"}`, http.StatusOK},
		{h, http.MethodPost, p + "/declarar", `{"decreto_numero":"1","decreto_em":"x"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, p + "/declarar", `{"decreto_numero":"Decreto 1/2026","decreto_em":"2026-10-01"}`, http.StatusOK},
		{h, http.MethodPost, p + "/guarda", `{"guarda":"DEPOSITADO","termo":"T","em":"x"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, p + "/guarda", `{"guarda":"DEPOSITADO","termo":"Termo 1","em":"2026-10-10"}`, http.StatusOK},
		{h, http.MethodPost, p + "/alienacoes", `{"notificado_em":"x"}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, p + "/alienacoes", `{"notificado_em":"2026-10-14","adquirente":"X"}`, http.StatusOK},
		{h, http.MethodPost, p + "/alienacoes/x/decidir", `{"decisao":"INTERESSE","fundamento":"f"}`, http.StatusBadRequest},
		{h, http.MethodPost, p + "/alienacoes/" + uuid.NewString() + "/decidir", `{"decisao":"INTERESSE","fundamento":"f"}`, http.StatusNotFound},
	}
	for _, c := range casos {
		if rec := do(c.h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	for _, path := range []string{"/acervo/normas", "/acervo/privados"} {
		down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
		if rec := do(down, http.MethodGet, path, ""); rec.Code != http.StatusInternalServerError {
			t.Errorf("%s com o banco fora: %d", path, rec.Code)
		}
	}
	downPub := chi.NewRouter()
	transport.NewHandlers(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), logger, 100).RegisterPublicRoutes(downPub)
	if rec := do(downPub, http.MethodGet, "/acervo/publico/normas", ""); rec.Code != http.StatusInternalServerError {
		t.Errorf("lista pública com o banco fora: %d", rec.Code)
	}
}
