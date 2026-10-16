package application_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
)

func TestHandlersFase3(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2026-09-01")
	s := e.real()
	_, membros := e.cspad(t, s)
	h := router(s, admin)
	criar := func() string {
		rec := do(h, http.MethodPost, "/acervo/unidades", `{"tipo":"CAIXA","codigo_classificacao":"`+codElim+`","unidade_id":"`+e.unidade.String()+
			`","data_inicial":"2018-01-01","data_final":"2018-12-31"}`)
		return campo(t, rec.Body.String(), "id")
	}
	a, b := criar(), criar()
	rec := do(h, http.MethodPost, "/acervo/listagens", `{"unidade_id":"`+e.unidade.String()+`","parecer_juridico":"P","contas_aprovadas_em":"2019-05-01","unidades":["`+a+`","`+b+`"]}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("relação: %d %s", rec.Code, rec.Body)
	}
	id := campo(t, rec.Body.String(), "id")
	presentes := `["` + membros[0].String() + `","` + membros[1].String() + `","` + membros[2].String() + `"]`
	base := "/acervo/listagens/" + id
	membrosJSON := `[{"user_id":"` + membros[0].String() + `","papel":"PRESIDENTE","area":"ARQUIVOLOGIA_HISTORIA"},` +
		`{"user_id":"` + membros[1].String() + `","papel":"MEMBRO","area":"JURIDICO"},` +
		`{"user_id":"` + membros[2].String() + `","papel":"MEMBRO","area":"ADMINISTRACAO_FINANCAS"},` +
		`{"user_id":"` + membros[3].String() + `","papel":"MEMBRO","area":"ARQUIVO_PUBLICO"},` +
		`{"user_id":"` + membros[4].String() + `","papel":"MEMBRO","area":"OUTRA"}]`
	central := do(h, http.MethodPost, "/acervo/comissoes", `{"tipo":"central","nome":"CCPAD","ativa":false,"membros":`+membrosJSON+`}`)
	if central.Code != http.StatusCreated {
		t.Fatalf("CCPAD: %d %s", central.Code, central.Body)
	}
	cid := campo(t, central.Body.String(), "id")

	casos := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/acervo/comissoes", "", http.StatusOK},
		{http.MethodPost, "/acervo/comissoes", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/comissoes", `{"tipo":"X","nome":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/comissoes/x", `{"tipo":"CENTRAL","nome":"x"}`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/comissoes/" + cid, `{"tipo":"CENTRAL","nome":"CCPAD","ativa":false,"membros":` + membrosJSON + `}`, http.StatusOK},
		{http.MethodGet, "/acervo/listagens", "", http.StatusOK},
		{http.MethodGet, "/acervo/listagens/x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/listagens/" + uuid.NewString(), "", http.StatusNotFound},
		{http.MethodGet, base, "", http.StatusOK},
		{http.MethodPost, "/acervo/listagens", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/listagens", `{"unidade_id":"` + e.unidade.String() + `","contas_aprovadas_em":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/listagens", `{"unidade_id":"` + e.outra.String() + `"}`, http.StatusConflict},
		{http.MethodPut, "/acervo/listagens/x", `{}`, http.StatusBadRequest},
		{http.MethodPut, base, `{`, http.StatusBadRequest},
		{http.MethodPut, base, `{"unidade_id":"` + e.unidade.String() + `","contas_aprovadas_em":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/listagens/" + uuid.NewString(), `{"unidade_id":"` + e.unidade.String() + `"}`, http.StatusNotFound},
		{http.MethodPut, base, `{"unidade_id":"` + e.unidade.String() + `","parecer_juridico":"Parecer 1"}`, http.StatusOK},
		{http.MethodPost, "/acervo/listagens/x/aprovar", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + "/aprovar", `{`, http.StatusBadRequest},
		{http.MethodPost, base + "/aprovar", `{"ata_numero":"1","ata_em":"ontem"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/aprovar", `{"ata_numero":"1","ata_em":"2026-09-01T10:00"}`, http.StatusNotFound},
		{http.MethodPost, base + "/aprovar", `{"ata_numero":"1/2026","ata_em":"2026-09-01T10:00","presentes":` + presentes + `}`, http.StatusOK},
		{http.MethodPost, base + "/devolver", `{"motivo":"Conferir as datas-limite"}`, http.StatusOK},
		{http.MethodPost, base + "/aprovar", `{"ata_numero":"1/2026","ata_em":"2026-09-01T10:00","presentes":` + presentes + `}`, http.StatusOK},
		{http.MethodPost, base + "/devolver", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/listagens/x/devolver", `{"motivo":"x"}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/devolver", `{"motivo":"x"}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/listagens/x/autorizar", ``, http.StatusBadRequest},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/autorizar", ``, http.StatusNotFound},
		{http.MethodPost, base + "/autorizar", ``, http.StatusOK},
		{http.MethodPost, "/acervo/listagens/x/edital", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + "/edital", `{`, http.StatusBadRequest},
		{http.MethodPost, base + "/edital", `{"edital_numero":"1","edital_publicado_em":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/edital", `{"edital_numero":"1","edital_publicado_em":"2026-09-01"}`, http.StatusNotFound},
		{http.MethodPost, base + "/edital", `{"edital_numero":"3/2026","edital_publicado_em":"2026-09-01","edital_diario":"DIORONDON-e"}`, http.StatusOK},
		{http.MethodPost, "/acervo/listagens/x/pedidos", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + "/pedidos", `{`, http.StatusBadRequest},
		{http.MethodPost, base + "/pedidos", `{"requerente":"R","pedido":"P","recebido_em":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/pedidos", `{"requerente":"R","pedido":"P","recebido_em":"2026-09-02"}`, http.StatusNotFound},
		{http.MethodPost, base + "/pedidos", `{"unidade_arquivamento_id":"` + a + `","requerente":"R","pedido":"P","recebido_em":"2026-09-02"}`, http.StatusCreated},
		{http.MethodPost, "/acervo/listagens/x/pedidos/" + uuid.NewString() + "/decidir", `{"decisao":"x"}`, http.StatusBadRequest},
		{http.MethodPost, base + "/pedidos/x/decidir", `{"decisao":"x"}`, http.StatusBadRequest},
		{http.MethodPost, base + "/pedidos/" + uuid.NewString() + "/decidir", `{`, http.StatusBadRequest},
		{http.MethodPost, base + "/pedidos/" + uuid.NewString() + "/decidir", `{"decisao":"x"}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/listagens/x/amostragem", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + "/amostragem", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/amostragem", `{"criterios":"c","percentual":10}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/listagens/x/termo", `{}`, http.StatusBadRequest},
		{http.MethodPost, base + "/termo", `{`, http.StatusBadRequest},
		{http.MethodPost, base + "/termo", `{"termo_eliminado_em":"x","termo_acompanhamento":"X"}`, http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/listagens/" + uuid.NewString() + "/termo", `{"termo_eliminado_em":"2026-10-02","termo_acompanhamento":"X"}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/listagens/x/cancelar", ``, http.StatusBadRequest},
	}
	for _, c := range casos {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	// Decide o pedido, amostra e registra o termo depois do prazo.
	ficha := do(h, http.MethodGet, base, "").Body.String()
	pid := strings.Split(strings.Split(ficha, `"pedidos":[{"id":"`)[1], `"`)[0]
	if rec := do(h, http.MethodPost, base+"/pedidos/"+pid+"/decidir", `{"deferido":false,"decisao":"Indeferido"}`); rec.Code != http.StatusOK {
		t.Fatalf("decisão: %d %s", rec.Code, rec.Body)
	}
	if rec := do(h, http.MethodPost, base+"/amostragem", `{"criterios":"sorteio","percentual":10}`); rec.Code != http.StatusOK {
		t.Fatalf("amostragem: %d %s", rec.Code, rec.Body)
	}
	e.agora = dia("2026-10-02")
	if rec := do(h, http.MethodPost, base+"/termo", `{"termo_eliminado_em":"2026-10-02","termo_acompanhamento":"Servidor X"}`); rec.Code != http.StatusOK {
		t.Fatalf("termo: %d %s", rec.Code, rec.Body)
	}
	outra := campo(t, do(h, http.MethodPost, "/acervo/listagens", `{"unidade_id":"`+e.unidade.String()+`","unidades":["`+criar()+`"]}`).Body.String(), "id")
	if rec := do(h, http.MethodPost, "/acervo/listagens/"+outra+"/cancelar", ``); rec.Code != http.StatusOK {
		t.Fatalf("cancelar: %d %s", rec.Code, rec.Body)
	}
	for _, p := range []string{"/acervo/comissoes", "/acervo/listagens"} {
		down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
		if rec := do(down, http.MethodGet, p, ""); rec.Code != http.StatusInternalServerError {
			t.Errorf("%s com o banco fora: %d", p, rec.Code)
		}
	}
}
