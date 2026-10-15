package app

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"

	acervoDomain "github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/atlas"
)

type acervoUnidade struct {
	ID          string `json:"id"`
	Numero      string `json:"numero"`
	Fase        string `json:"fase"`
	Destinacao  string `json:"destinacao"`
	Situacao    string `json:"situacao"`
	SerieResumo struct {
		Descritor string `json:"descritor"`
		Vigente   bool   `json:"vigente"`
	} `json:"serie"`
	Pode struct {
		Gerir    bool `json:"gerir"`
		Bloquear bool `json:"bloquear"`
	} `json:"pode"`
}

// Acervo (ADR 032) pela API real: classificação pela TTDD oficial do Atlas,
// escopo por unidade (ADR 013), locais só do Arquivo Público, bloqueio
// judicial e evento + auditoria na mesma transação.
func TestAcervoHTTP(t *testing.T) {
	h := newHarness(t)
	admin := h.admin()
	o := h.orgEscopo(t)
	agenteA := h.gestorEm(t, o.a, "acervo:read", "acervo:manage", "acervo:bloqueio")
	agenteB := h.gestorEm(t, o.b, "acervo:read", "acervo:manage")
	const serie = "2.0.02.00.07" // série vigente da TTDD oficial (Atlas)
	corpo := func(unidade, codigo string) string {
		return `{"tipo":"CAIXA","codigo_classificacao":"` + codigo + `","unidade_id":"` + unidade +
			`","data_inicial":"2015-01-01","data_final":"2015-12-31","quantidade":3,"especificacao":"caixas-arquivo"}`
	}

	u := data[acervoUnidade](t, h.expect(http.StatusCreated, http.MethodPost, "/api/v1/acervo/unidades", agenteA, corpo(o.a, serie)))
	if !strings.HasPrefix(u.Numero, "UA-") || u.Fase != "CORRENTE" || !u.SerieResumo.Vigente || u.SerieResumo.Descritor == "" ||
		u.Destinacao == "" || !u.Pode.Gerir || !u.Pode.Bloquear {
		t.Fatalf("cadastro: %+v", u)
	}
	if n := outboxCount(t, h, "acervo.unidade.cadastrada", u.ID); n != 1 {
		t.Fatalf("evento de cadastro: %d", n)
	}
	if n := auditCount(t, h, "acervo.unidade.cadastrada", u.ID); n != 1 {
		t.Fatalf("auditoria do cadastro: %d", n)
	}

	// Série inexistente; unidade de outra secretaria; leitura e escrita fora
	// do escopo.
	h.expect(http.StatusUnprocessableEntity, http.MethodPost, "/api/v1/acervo/unidades", agenteA, corpo(o.a, "99.9.99.99.99"))
	h.expect(http.StatusForbidden, http.MethodPost, "/api/v1/acervo/unidades", agenteA, corpo(o.b, serie))
	h.expect(http.StatusNotFound, http.MethodGet, "/api/v1/acervo/unidades/"+u.ID, agenteB, "")
	h.expect(http.StatusNotFound, http.MethodPut, "/api/v1/acervo/unidades/"+u.ID, agenteB, corpo(o.b, serie))
	if body := h.expect(http.StatusOK, http.MethodGet, "/api/v1/acervo/unidades?q="+u.Numero, agenteB, "").Body.String(); strings.Contains(body, u.ID) {
		t.Fatal("a outra secretaria não lista")
	}
	if body := h.expect(http.StatusOK, http.MethodGet, "/api/v1/acervo/unidades?q="+u.Numero, agenteA, "").Body.String(); !strings.Contains(body, u.ID) {
		t.Fatal("a própria secretaria lista")
	}
	h.expect(http.StatusOK, http.MethodGet, "/api/v1/acervo/painel?unidade_id="+o.a, agenteA, "")

	// Bloqueio judicial: só com acervo:bloqueio.
	h.expect(http.StatusForbidden, http.MethodPost, "/api/v1/acervo/unidades/"+u.ID+"/bloqueios", agenteB, `{"processo":"1"}`)
	b := data[acervoUnidade](t, h.expect(http.StatusCreated, http.MethodPost, "/api/v1/acervo/unidades/"+u.ID+"/bloqueios", agenteA,
		`{"processo":"0001234-56.2026.8.11.0003","motivo":"Ação de cobrança"}`))
	if b.Situacao != acervoDomain.SituacaoBloqueada {
		t.Fatalf("bloqueada: %+v", b)
	}

	// Locais: só o Arquivo Público (acervo:central).
	h.expect(http.StatusForbidden, http.MethodPost, "/api/v1/acervo/locais", agenteA, `{"tipo":"PREDIO","nome":"x"}`)
	h.expect(http.StatusCreated, http.MethodPost, "/api/v1/acervo/locais", admin, `{"tipo":"PREDIO","nome":"Arquivo Central `+uuid.NewString()[:6]+`"}`)
	h.expect(http.StatusOK, http.MethodGet, "/api/v1/acervo/locais", agenteA, "")
}

// A porta de TTDD do Acervo lê o Atlas: série inexistente vira
// ErrSerieInvalida; outras falhas sobem.
func TestTTDDPort(t *testing.T) {
	d := testDeps(t)
	var mod *atlas.Module
	for _, p := range d.Kernel.Plugins() {
		if m, ok := p.(*atlas.Module); ok {
			mod = m
		}
	}
	port := &ttddPort{svc: mod.Service()}
	s, err := port.Serie(context.Background(), "2.0.02.00.07")
	if err != nil || s.Codigo != "2.0.02.00.07" || !s.Vigente || s.Destinacao == "" {
		t.Fatalf("série: %+v %v", s, err)
	}
	if _, err := port.Serie(context.Background(), "99.9.99.99.99"); err != acervoDomain.ErrSerieInvalida {
		t.Fatalf("inexistente: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := port.Serie(ctx, "2.0.02.00.07"); err == nil || err == acervoDomain.ErrSerieInvalida {
		t.Fatalf("falha do banco: %v", err)
	}
}
