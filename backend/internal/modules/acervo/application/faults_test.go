package application_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/modules/acervo/transport"
	"github.com/yurythx/sismarq/internal/platform/auth"
)

// Matriz de falhas: cada chamada ao repositório falha (ou envenena a
// transação logo depois) e o caso de uso tem de propagar o erro.
func TestEveryRepositoryFailureIsPropagated(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2021-01-01")
	real := e.real()
	predio, err := real.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: domain.LocalPredio, Nome: "Falhas " + uuid.NewString()[:8], Ativo: true})
	if err != nil {
		t.Fatal(err)
	}
	comLocal := func() application.Detalhe {
		u := e.nova(codElim)
		u.LocalID = &predio.ID
		d, err := real.Criar(ctx, admin, u)
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	bloqueada := func() (application.Detalhe, uuid.UUID) {
		d := comLocal()
		d, err := real.Bloquear(ctx, admin, d.ID, "1", "")
		if err != nil {
			t.Fatal(err)
		}
		return d, d.Bloqueios[0].ID
	}
	ops := map[string]func() func(*application.Service) error{
		"Criar": func() func(*application.Service) error {
			return func(s *application.Service) error {
				u := e.nova(codElim)
				u.LocalID = &predio.ID
				_, err := s.Criar(ctx, admin, u)
				return err
			}
		},
		"Alterar": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error {
				u := e.nova(codElim)
				u.LocalID = &predio.ID
				_, err := s.Alterar(ctx, admin, d.ID, u)
				return err
			}
		},
		"Transferir": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Transferir(ctx, admin, d.ID, ""); return err }
		},
		"Localizar": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Localizar(ctx, admin, d.ID, &predio.ID, "p"); return err }
		},
		"Bloquear": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Bloquear(ctx, admin, d.ID, "1", ""); return err }
		},
		"Liberar": func() func(*application.Service) error {
			_, bid := bloqueada()
			return func(s *application.Service) error { _, err := s.Liberar(ctx, admin, bid, "ok"); return err }
		},
		"Get": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Get(ctx, admin, d.ID); return err }
		},
		"List": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.List(ctx, admin, domain.Filtro{}, pagination.New(1, 5, 100))
				return err
			}
		},
		"Painel": func() func(*application.Service) error {
			return func(s *application.Service) error { _, err := s.Painel(ctx, admin, domain.Filtro{}); return err }
		},
		"ListLocais": func() func(*application.Service) error {
			return func(s *application.Service) error { _, err := s.ListLocais(ctx); return err }
		},
		"CriarLocal": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.SalvarLocal(ctx, admin, nil, domain.Local{ParentID: &predio.ID, Tipo: domain.LocalSala, Nome: uuid.NewString(), Ativo: true})
				return err
			}
		},
		"AlterarLocal": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.SalvarLocal(ctx, admin, &predio.ID, predio)
				return err
			}
		},
	}
	for name, o := range ops {
		probe := &faultRepo{inner: infrastructure.NewRepository()}
		if err := o()(e.svc(probe)); err != nil {
			t.Fatalf("%s sem falha: %v", name, err)
		}
		for i := range probe.trace {
			for _, poison := range []bool{false, true} {
				r := &faultRepo{inner: infrastructure.NewRepository()}
				if poison {
					r.poisonAt = i + 1
				} else {
					r.failAt = i + 1
				}
				if err := o()(e.svc(r)); err == nil && !r.noTx {
					t.Errorf("%s: falha (veneno=%v) na chamada %d (%s) foi engolida", name, poison, i+1, probe.trace[i])
				}
			}
		}
	}
}

// tudo responde qualquer código (o recálculo percorre todo o banco de teste).
type tudo struct{ vigente bool }

func (t tudo) Serie(_ context.Context, codigo string) (domain.Serie, error) {
	if !t.vigente {
		return domain.Serie{}, domain.ErrSerieInvalida
	}
	return domain.Serie{Codigo: codigo, CorrenteAnos: anos(1), IntermAnos: anos(1), Destinacao: domain.DestinacaoEliminacao, Vigente: true}, nil
}

// O recálculo percorre o banco inteiro: falha em cada ponto, de forma
// determinística (série conhecida: lista, bloqueio, gravação; série
// inexistente: lista, gravação da situação).
func TestRecalculoPropagaFalhas(t *testing.T) {
	e := novoEnv(t)
	e.criar(codElim)
	ctx := context.Background()
	for _, c := range []struct {
		ttdd   domain.TTDD
		failAt int
	}{{tudo{true}, 1}, {tudo{true}, 2}, {tudo{true}, 3}, {tudo{false}, 2}} {
		s := application.NewService(e.pool, &faultRepo{inner: infrastructure.NewRepository(), failAt: c.failAt}, c.ttdd, nil, logger)
		if _, err := s.Recalcular(ctx); !errors.Is(err, errBoom) {
			t.Errorf("falha na chamada %d (vigente=%v) engolida: %v", c.failAt, c.ttdd, err)
		}
	}
}

// router monta as rotas com a identidade dada já no contexto.
func router(s *application.Service, identity auth.Identity) http.Handler {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req.WithContext(auth.WithIdentity(req.Context(), identity)))
		})
	})
	transport.NewHandlers(s, logger, 100).RegisterRoutes(r)
	return r
}

func do(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandlers(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2021-01-01")
	h := router(e.real(), admin)
	corpo := func(extra string) string {
		return `{"tipo":"CAIXA","codigo_classificacao":"` + codElim + `","unidade_id":"` + e.unidade.String() +
			`","data_inicial":"2018-01-01","data_final":"2018-12-31","quantidade":1` + extra + `}`
	}
	rec := do(h, http.MethodPost, "/acervo/unidades", corpo(`,"data_evento":"2019-01-01","contas_aprovadas_em":"2019-05-01"`))
	if rec.Code != http.StatusCreated {
		t.Fatalf("criar: %d %s", rec.Code, rec.Body)
	}
	id := strings.Split(strings.Split(rec.Body.String(), `"id":"`)[1], `"`)[0]
	predio := do(h, http.MethodPost, "/acervo/locais", `{"tipo":"PREDIO","nome":"H `+uuid.NewString()[:8]+`"}`)
	if predio.Code != http.StatusCreated {
		t.Fatalf("local: %d %s", predio.Code, predio.Body)
	}
	localID := strings.Split(strings.Split(predio.Body.String(), `"id":"`)[1], `"`)[0]
	casos := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/acervo/unidades?q=UA&fase=corrente&unidade_id=" + e.unidade.String(), "", http.StatusOK},
		{http.MethodGet, "/acervo/unidades?unidade_id=x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/painel", "", http.StatusOK},
		{http.MethodGet, "/acervo/painel?unidade_id=x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/unidades/" + id, "", http.StatusOK},
		{http.MethodGet, "/acervo/unidades/x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/unidades/" + uuid.NewString(), "", http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "2018-01-01", "01/01/2018", 1), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "2018-12-31", "x", 1), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", corpo(`,"data_evento":"x"`), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", corpo(`,"contas_aprovadas_em":"x"`), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "CAIXA", "GAVETA", 1), http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/unidades/x", corpo(""), http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + id, `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + uuid.NewString(), corpo(""), http.StatusNotFound},
		{http.MethodPut, "/acervo/unidades/" + id, corpo(`,"observacoes":"revisada"`), http.StatusOK},
		{http.MethodPut, "/acervo/unidades/x/local", `{}`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + id + "/local", `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + uuid.NewString() + "/local", `{}`, http.StatusNotFound},
		{http.MethodPut, "/acervo/unidades/" + id + "/local", `{"local_id":"` + localID + `","posicao":"P1"}`, http.StatusOK},
		{http.MethodPost, "/acervo/unidades/x/transferir", `{}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/transferir", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + uuid.NewString() + "/transferir", `{}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades/" + id + "/transferir", `{"observacao":"ok"}`, http.StatusOK},
		{http.MethodPost, "/acervo/unidades/x/bloqueios", `{"processo":"1"}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/bloqueios", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + uuid.NewString() + "/bloqueios", `{"processo":"1"}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades/" + id + "/bloqueios", `{"processo":"0001"}`, http.StatusCreated},
		{http.MethodPost, "/acervo/bloqueios/x/liberar", `{"motivo":"ok"}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/bloqueios/" + uuid.NewString() + "/liberar", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/bloqueios/" + uuid.NewString() + "/liberar", `{"motivo":"ok"}`, http.StatusNotFound},
		{http.MethodGet, "/acervo/locais", "", http.StatusOK},
		{http.MethodPost, "/acervo/locais", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/locais", `{"tipo":"COFRE","nome":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/locais/x", `{"tipo":"PREDIO","nome":"x"}`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/locais/" + localID, `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/locais/" + uuid.NewString(), `{"tipo":"PREDIO","nome":"x"}`, http.StatusNotFound},
		{http.MethodPut, "/acervo/locais/" + localID, `{"tipo":"PREDIO","nome":"Renomeado ` + uuid.NewString()[:6] + `","ativo":true}`, http.StatusOK},
	}
	for _, c := range casos {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	// Liberar o bloqueio criado acima.
	ficha := do(h, http.MethodGet, "/acervo/unidades/"+id, "").Body.String()
	bid := strings.Split(strings.Split(ficha, `"bloqueios":[{"id":"`)[1], `"`)[0]
	if rec := do(h, http.MethodPost, "/acervo/bloqueios/"+bid+"/liberar", `{"motivo":"Encerrada"}`); rec.Code != http.StatusOK {
		t.Fatalf("liberar: %d %s", rec.Code, rec.Body)
	}
	// Banco fora: a listagem e o painel respondem 500 sem vazar o erro.
	for _, p := range []string{"/acervo/unidades", "/acervo/painel", "/acervo/locais"} {
		down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
		if rec := do(down, http.MethodGet, p, ""); rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), errBoom.Error()) {
			t.Errorf("%s com o banco fora: %d %s", p, rec.Code, rec.Body)
		}
	}
}
