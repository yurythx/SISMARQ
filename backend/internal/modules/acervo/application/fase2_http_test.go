package application_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

func campo(t *testing.T, body, nome string) string {
	t.Helper()
	var env struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatal(err)
	}
	v, _ := env.Data[nome].(string)
	return v
}

func TestHandlersFase2(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2021-01-01")
	h := router(e.real(), admin)
	criar := func() string {
		rec := do(h, http.MethodPost, "/acervo/unidades", `{"tipo":"CAIXA","codigo_classificacao":"`+codElim+`","unidade_id":"`+e.unidade.String()+
			`","data_inicial":"2018-01-01","data_final":"2018-12-31"}`)
		if rec.Code != http.StatusCreated {
			t.Fatalf("criar: %d %s", rec.Code, rec.Body)
		}
		return campo(t, rec.Body.String(), "id")
	}
	id := criar()
	sol := dbtest.User(t, e.pool).String()
	emp := do(h, http.MethodPost, "/acervo/unidades/"+id+"/emprestimos", `{"solicitante_id":"`+sol+`","finalidade":"consulta","devolver_ate":"2021-01-10"}`)
	if emp.Code != http.StatusCreated {
		t.Fatalf("emprestar: %d %s", emp.Code, emp.Body)
	}
	eid := strings.Split(strings.Split(emp.Body.String(), `"emprestimos":[{"id":"`)[1], `"`)[0]
	g := do(h, http.MethodPost, "/acervo/guias", `{"origem_id":"`+e.unidade.String()+`","unidades":["`+criar()+`"]}`)
	if g.Code != http.StatusCreated {
		t.Fatalf("guia: %d %s", g.Code, g.Body)
	}
	gid := campo(t, g.Body.String(), "id")
	g2 := campo(t, do(h, http.MethodPost, "/acervo/guias", `{"origem_id":"`+e.unidade.String()+`","unidades":["`+criar()+`"]}`).Body.String(), "id")
	cab := "tipo;codigo_classificacao;unidade;data_inicial;data_final\\n"
	linha := "CAIXA;" + codElim + ";" + e.unidade.String() + ";01/01/2018;31/12/2018\\n"
	sim := do(h, http.MethodPost, "/acervo/importacao/simular", `{"conteudo":"`+cab+linha+`"}`)
	hash := campo(t, sim.Body.String(), "hash")

	casos := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodPost, "/acervo/unidades/x/emprestimos", `{}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/emprestimos", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/emprestimos", `{"solicitante_id":"` + sol + `","finalidade":"x","devolver_ate":"10/01/2021"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades/" + id + "/emprestimos", `{"solicitante_id":"` + sol + `","finalidade":"x","devolver_ate":"2021-01-10"}`, http.StatusConflict},
		{http.MethodPost, "/acervo/emprestimos/x/devolver", `{}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/emprestimos/" + eid + "/devolver", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/emprestimos/" + uuid.NewString() + "/devolver", `{}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/emprestimos/" + eid + "/devolver", `{"observacao":"ok"}`, http.StatusOK},
		{http.MethodGet, "/acervo/emprestimos?atrasados=true", "", http.StatusOK},
		{http.MethodGet, "/acervo/guias", "", http.StatusOK},
		{http.MethodGet, "/acervo/guias/x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/guias/" + uuid.NewString(), "", http.StatusNotFound},
		{http.MethodGet, "/acervo/guias/" + gid, "", http.StatusOK},
		{http.MethodPost, "/acervo/guias", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/guias", `{"origem_id":"` + e.unidade.String() + `","unidades":["` + uuid.NewString() + `"]}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/guias/x/receber", "", http.StatusBadRequest},
		{http.MethodPost, "/acervo/guias/" + uuid.NewString() + "/receber", "", http.StatusNotFound},
		{http.MethodPost, "/acervo/guias/" + gid + "/receber", "", http.StatusOK},
		{http.MethodPost, "/acervo/guias/" + g2 + "/cancelar", "", http.StatusOK},
		{http.MethodGet, "/acervo/importacao/modelo", "", http.StatusOK},
		{http.MethodPost, "/acervo/importacao/simular", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/importacao/simular", `{"conteudo":""}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/importacao/simular", `{"conteudo":"nada"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/importacao/aplicar", `{"conteudo":"` + cab + linha + `","hash":"` + hash + `"}`, http.StatusOK},
	}
	for _, c := range casos {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	modelo := do(h, http.MethodGet, "/acervo/importacao/modelo", "")
	if !strings.HasPrefix(modelo.Header().Get("Content-Type"), "text/csv") || !strings.Contains(modelo.Body.String(), "codigo_classificacao") {
		t.Fatalf("modelo: %s %s", modelo.Header(), modelo.Body)
	}
	for _, p := range []string{"/acervo/emprestimos", "/acervo/guias"} {
		down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
		if rec := do(down, http.MethodGet, p, ""); rec.Code != http.StatusInternalServerError {
			t.Errorf("%s com o banco fora: %d", p, rec.Code)
		}
	}
}
