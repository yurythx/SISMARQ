package application_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
)

func TestHandlersFase4(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2019-06-01")
	s := e.real()
	h := router(s, admin)
	criar := func() string {
		rec := do(h, http.MethodPost, "/acervo/unidades", `{"tipo":"CAIXA","codigo_classificacao":"`+codPerm+`","unidade_id":"`+e.unidade.String()+
			`","data_inicial":"2018-01-01","data_final":"2018-12-31"}`)
		return campo(t, rec.Body.String(), "id")
	}
	preparoJSON := `"preparo":{"classificados":true,"avaliados":true,"organizados":true,"higienizados":true,"acondicionados":true}`
	pedido := func(u string) string {
		return `{"origem_id":"` + e.unidade.String() + `","motivo":"EXTINCAO","fundamento":"Lei 1/2026",` + preparoJSON + `,"unidades":["` + u + `"]}`
	}
	rec := do(h, http.MethodPost, "/acervo/recolhimentos", pedido(criar()))
	if rec.Code != http.StatusCreated {
		t.Fatalf("recolhimento: %d %s", rec.Code, rec.Body)
	}
	id := campo(t, rec.Body.String(), "id")
	outro := campo(t, do(h, http.MethodPost, "/acervo/recolhimentos", pedido(criar())).Body.String(), "id")
	terceiro := campo(t, do(h, http.MethodPost, "/acervo/recolhimentos", pedido(criar())).Body.String(), "id")
	base := "/acervo/recolhimentos/"

	casos := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/acervo/recolhimentos", "", http.StatusOK},
		{http.MethodGet, base + id, "", http.StatusOK},
		{http.MethodGet, base + "x", "", http.StatusBadRequest},
		{http.MethodGet, base + uuid.NewString(), "", http.StatusNotFound},
		{http.MethodPost, "/acervo/recolhimentos", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/recolhimentos", `{"origem_id":"` + e.unidade.String() + `","motivo":"REGULAR","unidades":["` + uuid.NewString() + `"]}`,
			http.StatusUnprocessableEntity},
		{http.MethodPost, base + "x/receber", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + id + "/receber", `{`, http.StatusBadRequest},
		{http.MethodPost, base + uuid.NewString() + "/receber", `{}`, http.StatusNotFound},
		{http.MethodPost, base + id + "/receber", `{}`, http.StatusOK},
		{http.MethodPost, base + id + "/receber", `{}`, http.StatusConflict},
		{http.MethodPost, base + "x/devolver", `{"motivo":"x"}`, http.StatusBadRequest},
		{http.MethodPost, base + outro + "/devolver", `{`, http.StatusBadRequest},
		{http.MethodPost, base + uuid.NewString() + "/devolver", `{"motivo":"x"}`, http.StatusNotFound},
		{http.MethodPost, base + outro + "/devolver", `{"motivo":"Faltam etiquetas"}`, http.StatusOK},
		{http.MethodPost, base + "x/cancelar", ``, http.StatusBadRequest},
		{http.MethodPost, base + uuid.NewString() + "/cancelar", ``, http.StatusNotFound},
		{http.MethodPost, base + terceiro + "/cancelar", ``, http.StatusOK},
	}
	for _, c := range casos {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
	if rec := do(down, http.MethodGet, "/acervo/recolhimentos", ""); rec.Code != http.StatusInternalServerError {
		t.Errorf("com o banco fora: %d", rec.Code)
	}
}
