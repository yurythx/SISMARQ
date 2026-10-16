package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
	"github.com/yurythx/sismarq/internal/platform/typesense"
)

func TestAcessoPropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	agora := time.Now()
	x := domain.Exposicao{ID: id, Slug: "x", Titulo: "x", Itens: []domain.ItemExposicao{{DigitalID: id}}}
	c := domain.Consulta{Q: "ofício", Pagina: 1, PorPagina: 10}
	calls := map[string]func(db database.DBTX) error{
		"PublicarDigital":   func(db database.DBTX) error { return r.PublicarDigital(ctx, db, id, &agora, nil) },
		"PendentesIndice":   func(db database.DBTX) error { _, err := r.PendentesIndice(ctx, db, 10); return err },
		"ConcluirPendentes": func(db database.DBTX) error { return r.ConcluirPendentes(ctx, db, []domain.PendenteIndice{{UnidadeID: id}}) },
		"RegistrosIndice":   func(db database.DBTX) error { _, err := r.RegistrosIndice(ctx, db, []uuid.UUID{id}, false); return err },
		"PaginaIndice":      func(db database.DBTX) error { _, err := r.PaginaIndice(ctx, db, id, 10); return err },
		"BuscarBanco":       func(db database.DBTX) error { _, _, err := r.BuscarBanco(ctx, db, c, domain.FiltroBusca{}); return err },
		"DigitaisRecentes":  func(db database.DBTX) error { _, err := r.DigitaisRecentes(ctx, db, 5); return err },
		"ContarPublico":     func(db database.DBTX) error { _, _, err := r.ContarPublico(ctx, db); return err },
		"RegistrarConsulta": func(db database.DBTX) error { return r.RegistrarConsulta(ctx, db, agora, domain.ConsultaBusca) },
		"InsertExposicao":   func(db database.DBTX) error { return r.InsertExposicao(ctx, db, x) },
		"UpdateExposicao":   func(db database.DBTX) error { return r.UpdateExposicao(ctx, db, x) },
		"GetExposicao":      func(db database.DBTX) error { _, err := r.GetExposicao(ctx, db, id, ""); return err },
		"ListExposicoes":    func(db database.DBTX) error { _, _, err := r.ListExposicoes(ctx, db, true, p); return err },
		"PublicarExposicao": func(db database.DBTX) error { return r.PublicarExposicao(ctx, db, id, nil) },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	for _, name := range []string{"PendentesIndice", "RegistrosIndice", "PaginaIndice", "DigitaisRecentes"} {
		if err := calls[name](dbtest.ScanFail{}); err == nil {
			t.Errorf("%s com linha ilegível deveria falhar", name)
		}
		if err := calls[name](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com erro durante a leitura: %v", name, err)
		}
	}
	// Os representantes das unidades: a segunda consulta falha.
	for _, nova := range []dbtest.QueryResult{{Err: dbtest.ErrInjected}, {Rows: dbtest.BadRows()}, {Rows: dbtest.ErrRows()}} {
		seq := &dbtest.Seq{Queries: []dbtest.QueryResult{{Rows: &umaLinha{}}, nova}}
		if _, err := r.RegistrosIndice(ctx, seq, []uuid.UUID{id}, true); err == nil {
			t.Error("representantes ilegíveis")
		}
	}
	if regs, err := r.RegistrosIndice(ctx, &dbtest.Seq{Queries: []dbtest.QueryResult{{Rows: &umaLinha{}}, {Rows: &umaLinha{}}}}, nil, true); err != nil ||
		len(regs) != 1 || len(regs[0].Digitais) != 1 {
		t.Fatalf("registro com um representante: %+v %v", regs, err)
	}
	// A página da busca e a das exposições: a consulta, a linha e a leitura.
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, _, err := r.BuscarBanco(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, c, domain.FiltroBusca{}); err == nil {
			t.Error("página da busca ilegível")
		}
		if _, _, err := r.ListExposicoes(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, false, p); err == nil {
			t.Error("página das exposições ilegível")
		}
		if _, err := r.GetExposicao(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, ""); err == nil {
			t.Error("itens da exposição ilegíveis")
		}
	}

	// Itens das exposições: a gravação que falha depois do texto.
	falhaNa := func(n int) *dbtest.Seq {
		execs := make([]dbtest.ExecResult, n)
		for i := range execs {
			execs[i] = dbtest.OK
		}
		return &dbtest.Seq{Execs: append(execs, dbtest.ExecResult{Err: dbtest.ErrInjected})}
	}
	if err := r.InsertExposicao(ctx, falhaNa(1), x); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("itens na criação: %v", err)
	}
	for _, n := range []int{1, 2} {
		if err := r.UpdateExposicao(ctx, falhaNa(n), x); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("alteração, falha na gravação %d: %v", n+1, err)
		}
	}
}

// tsFalso imita o Typesense: guarda o que recebeu e falha onde pedido.
type tsFalso struct {
	criou                                 bool
	falhaEnsure, falhaImport, falhaDelete int // na chamada n (1 = primeira)
	falhaSearch                           bool
	ensures, imports, deletes             int
	importados                            map[string][]map[string]any
	filtrosDelete                         []string
	params                                typesense.SearchParams
	colecao                               string
	resposta                              typesense.SearchResult
}

func (f *tsFalso) EnsureCollection(_ context.Context, _ typesense.Schema) (bool, error) {
	f.ensures++
	if f.ensures == f.falhaEnsure {
		return false, dbtest.ErrInjected
	}
	return f.criou, nil
}

func (f *tsFalso) Import(_ context.Context, c string, docs []map[string]any) error {
	f.imports++
	if f.imports == f.falhaImport {
		return dbtest.ErrInjected
	}
	if f.importados == nil {
		f.importados = map[string][]map[string]any{}
	}
	f.importados[c] = append(f.importados[c], docs...)
	return nil
}

func (f *tsFalso) DeleteByFilter(_ context.Context, c, filtro string) (int, error) {
	f.deletes++
	if f.deletes == f.falhaDelete {
		return 0, dbtest.ErrInjected
	}
	f.filtrosDelete = append(f.filtrosDelete, c+" "+filtro)
	return 1, nil
}

func (f *tsFalso) Search(_ context.Context, c string, p typesense.SearchParams) (typesense.SearchResult, error) {
	f.colecao, f.params = c, p
	if f.falhaSearch {
		return typesense.SearchResult{}, dbtest.ErrInjected
	}
	return f.resposta, nil
}

func registros() []domain.RegistroIndice {
	agora := time.Now()
	d := func(s string) time.Time { x, _ := time.Parse("2006-01-02", s); return x }
	return []domain.RegistroIndice{
		{Unidade: domain.Unidade{ID: uuid.New(), Numero: "UA-1", Tipo: domain.TipoCaixa, CodigoClassificacao: "2.0.03.02", Assunto: "Ofícios",
			UnidadeID: uuid.New(), DataInicial: d("2015-01-01"), DataFinal: d("2016-12-31"), Fase: domain.FasePermanente,
			Destinacao: domain.DestinacaoGuardaPermanente, NivelAcesso: domain.AcessoPublico, Observacoes: "sala 3", UpdatedAt: agora},
			Orgao: "SEMAD", Descritor: "Pedidos de acesso", Digitais: []domain.DigitalIndice{{Texto: "liberado", PublicadoEm: &agora}, {Texto: "interno"}}},
		{Unidade: domain.Unidade{ID: uuid.New(), Numero: "UA-2", Tipo: domain.TipoPasta, CodigoClassificacao: "1.0", NivelAcesso: domain.AcessoRestrito,
			Fase: domain.FaseCorrente, UnidadeID: uuid.New(), DataInicial: d("2020-01-01"), DataFinal: d("2020-01-01")},
			Orgao: "SMS", Digitais: []domain.DigitalIndice{{Texto: "prontuário"}}},
	}
}

func TestIndiceGravaNasDuasColecoes(t *testing.T) {
	ctx := context.Background()
	ts := &tsFalso{criou: true}
	ix := NovoIndice(ts, "interna", "publica")
	if criou, err := ix.Garantir(ctx); err != nil || !criou || ts.ensures != 2 {
		t.Fatalf("garantir: %v %v", criou, err)
	}
	if s := ix.Schema(true); s.Name != "publica" || strings.Contains(fmt.Sprint(s.Fields), "orgao_id") || strings.Contains(fmt.Sprint(s.Fields), "nivel_acesso") {
		t.Fatalf("a coleção pública não tem nível nem órgão: %+v", s)
	}
	regs := registros()
	if err := ix.Gravar(ctx, regs, 7); err != nil {
		t.Fatal(err)
	}
	internos, publicos := ts.importados["interna"], ts.importados["publica"]
	if len(internos) != 2 || len(publicos) != 1 {
		t.Fatalf("documentos: %d internos, %d públicos", len(internos), len(publicos))
	}
	pub, int0 := publicos[0], internos[0]
	if pub["texto"] != "liberado" || pub["digitais"] != 1 || pub["nivel_acesso"] != nil || pub["orgao_id"] != nil || pub["detalhes"] != nil ||
		pub["classe"] != "2" || fmt.Sprint(pub["ano"]) != "[2015 2016]" || pub["geracao"] != int64(7) {
		t.Fatalf("documento público: %+v", pub)
	}
	if int0["texto"] != "liberado\ninterno" || int0["digitais"] != 2 || int0["nivel_acesso"] != "PUBLICO" || int0["detalhes"] != "sala 3" {
		t.Fatalf("documento interno: %+v", int0)
	}
	// A restrita sai da coleção pública.
	if len(ts.filtrosDelete) != 1 || !strings.Contains(ts.filtrosDelete[0], "publica id:[`"+regs[1].Unidade.ID.String()+"`]") {
		t.Fatalf("saída da pública: %v", ts.filtrosDelete)
	}
	if err := ix.Limpar(ctx, 9); err != nil || len(ts.filtrosDelete) != 3 || !strings.HasSuffix(ts.filtrosDelete[2], "publica geracao:<9") {
		t.Fatalf("limpar: %v %v", ts.filtrosDelete, err)
	}
	// Só públicas: nada a tirar.
	ts2 := &tsFalso{}
	if err := NovoIndice(ts2, "i", "p").Gravar(ctx, regs[:1], 1); err != nil || ts2.deletes != 0 {
		t.Fatalf("sem saída: %d %v", ts2.deletes, err)
	}
}

func TestIndicePropagaFalhas(t *testing.T) {
	ctx := context.Background()
	for _, n := range []int{1, 2} {
		if _, err := NovoIndice(&tsFalso{falhaEnsure: n}, "i", "p").Garantir(ctx); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("garantir %d: %v", n, err)
		}
		if err := NovoIndice(&tsFalso{falhaImport: n}, "i", "p").Gravar(ctx, registros(), 1); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("import %d: %v", n, err)
		}
		if err := NovoIndice(&tsFalso{falhaDelete: n}, "i", "p").Limpar(ctx, 1); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("limpar %d: %v", n, err)
		}
	}
	if err := NovoIndice(&tsFalso{falhaDelete: 1}, "i", "p").Gravar(ctx, registros(), 1); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("tirar da pública: %v", err)
	}
	if _, err := NovoIndice(&tsFalso{falhaSearch: true}, "i", "p").Buscar(ctx, domain.Consulta{}, domain.FiltroBusca{}); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("buscar: %v", err)
	}
}

func TestIndiceBusca(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	ts := &tsFalso{resposta: typesense.SearchResult{Found: 3, SearchTime: 4,
		Hits: []typesense.Hit{
			{Document: map[string]any{"id": id.String(), "numero": "UA-1", "especie": "CAIXA", "ano_inicial": float64(2015), "digitais": float64(1),
				"nivel_acesso": "RESTRITO"}, Highlights: []typesense.Highlight{{Field: "assunto"}, {Field: "texto", Snippet: "o <mark>ofício</mark>"}}},
			{Document: map[string]any{"id": "x"}},
		},
		FacetCounts: []typesense.FacetCount{{FieldName: "fase", Counts: []struct {
			Value string `json:"value"`
			Count int    `json:"count"`
		}{{Value: "PERMANENTE", Count: 3}}}}}}
	ix := NovoIndice(ts, "interna", "publica")
	unidade := uuid.New()
	res, err := ix.Buscar(ctx, domain.Consulta{Q: "ofício", Pagina: 2, PorPagina: 10, SoDigital: true,
		Filtros: map[string]string{domain.FacetaAno: "2015", domain.FacetaOrgao: "SEMAD`x", domain.FacetaNivel: "RESTRITO"}},
		domain.FiltroBusca{Restrito: true, Unidades: []uuid.UUID{unidade}})
	if err != nil || ts.colecao != "interna" || res.Total != 3 || res.Ms != 4 || res.Origem != domain.OrigemIndice || len(res.Itens) != 2 {
		t.Fatalf("busca interna: %+v %v", res, err)
	}
	it := res.Itens[0]
	if it.ID != id || it.Tipo != "CAIXA" || it.AnoInicial != 2015 || it.Digitais != 1 || it.Trecho != "o <mark>ofício</mark>" || it.NivelAcesso != "RESTRITO" ||
		res.Itens[1].ID != uuid.Nil || res.Facetas["fase"][0] != (domain.Faceta{Valor: "PERMANENTE", Total: 3}) {
		t.Fatalf("itens e facetas: %+v %+v", res.Itens, res.Facetas)
	}
	f := ts.params.FilterBy
	for _, parte := range []string{"orgao_id:[`" + unidade.String() + "`]", "ano:=2015", "orgao:=`SEMAD\\`x`", "nivel_acesso:=`RESTRITO`", "digitais:>0"} {
		if !strings.Contains(f, parte) {
			t.Fatalf("filtro %q ausente em %q", parte, f)
		}
	}
	if ts.params.Page != 2 || ts.params.SortBy != "_text_match:desc,updated_at:desc" || !strings.Contains(ts.params.FacetBy, "nivel_acesso") ||
		!strings.Contains(ts.params.QueryBy, "detalhes") {
		t.Fatalf("parâmetros: %+v", ts.params)
	}
	// Pública: outra coleção, sem escopo nem nível, ordem pela data.
	if _, err := ix.Buscar(ctx, domain.Consulta{Filtros: map[string]string{domain.FacetaNivel: "RESTRITO"}},
		domain.FiltroBusca{Publico: true, Restrito: true, Unidades: []uuid.UUID{unidade}}); err != nil {
		t.Fatal(err)
	}
	if ts.colecao != "publica" || ts.params.FilterBy != "" || ts.params.Q != "*" || ts.params.SortBy != "updated_at:desc" ||
		strings.Contains(ts.params.FacetBy, "nivel") || strings.Contains(ts.params.QueryBy, "detalhes") {
		t.Fatalf("busca pública: %s %+v", ts.colecao, ts.params)
	}
}

// TestIndiceReal grava e consulta num Typesense de verdade (CI):
// TEST_TYPESENSE_URL e TEST_TYPESENSE_API_KEY.
func TestIndiceReal(t *testing.T) {
	base := os.Getenv("TEST_TYPESENSE_URL")
	if base == "" {
		t.Skip("TEST_TYPESENSE_URL não definido")
	}
	ctx := context.Background()
	c := typesense.New(base, os.Getenv("TEST_TYPESENSE_API_KEY"), 5*time.Second)
	sufixo := time.Now().UnixNano()
	ix := NovoIndice(c, fmt.Sprintf("acervo_i_%d", sufixo), fmt.Sprintf("acervo_p_%d", sufixo))
	t.Cleanup(func() {
		for _, publico := range []bool{false, true} {
			_, _ = c.DeleteByFilter(context.Background(), ix.Schema(publico).Name, "geracao:>0")
		}
	})
	if criou, err := ix.Garantir(ctx); err != nil || !criou {
		t.Fatalf("garantir: %v %v", criou, err)
	}
	regs := registros()
	if err := ix.Gravar(ctx, regs, 1); err != nil {
		t.Fatal(err)
	}
	// Sem acento, no texto do OCR liberado e com faceta.
	res, err := ix.Buscar(ctx, domain.Consulta{Q: "liberado", Pagina: 1, PorPagina: 10}, domain.FiltroBusca{Publico: true})
	if err != nil || res.Total != 1 || res.Itens[0].Numero != "UA-1" || !strings.Contains(res.Itens[0].Trecho, "<mark>") || len(res.Facetas["ano"]) != 2 {
		t.Fatalf("pública: %+v %v", res, err)
	}
	// O texto interno e a unidade restrita não estão na coleção pública.
	for _, q := range []string{"interno", "prontuario"} {
		if res, err := ix.Buscar(ctx, domain.Consulta{Q: q, Pagina: 1, PorPagina: 10}, domain.FiltroBusca{Publico: true}); err != nil || res.Total != 0 {
			t.Fatalf("%q vazou na pública: %+v %v", q, res, err)
		}
	}
	// Na interna, o escopo filtra: quem não tem o órgão da restrita não a vê.
	fora := domain.FiltroBusca{Restrito: true, Unidades: []uuid.UUID{regs[0].Unidade.UnidadeID}}
	if res, err := ix.Buscar(ctx, domain.Consulta{Q: "prontuario", Pagina: 1, PorPagina: 10}, fora); err != nil || res.Total != 0 {
		t.Fatalf("restrita fora do escopo: %+v %v", res, err)
	}
	dentro := domain.FiltroBusca{Restrito: true, Unidades: []uuid.UUID{regs[1].Unidade.UnidadeID}}
	if res, err := ix.Buscar(ctx, domain.Consulta{Q: "prontuario", Pagina: 1, PorPagina: 10, Filtros: map[string]string{domain.FacetaNivel: "RESTRITO"}}, dentro); err != nil ||
		res.Total != 1 || res.Itens[0].NivelAcesso != "RESTRITO" {
		t.Fatalf("restrita no escopo: %+v %v", res, err)
	}
	// A unidade que fica restrita sai da pública; a carga nova limpa a antiga.
	regs[0].Unidade.NivelAcesso = domain.AcessoRestrito
	if err := ix.Gravar(ctx, regs[:1], 2); err != nil {
		t.Fatal(err)
	}
	if res, _ := ix.Buscar(ctx, domain.Consulta{Q: "*", Pagina: 1, PorPagina: 10}, domain.FiltroBusca{Publico: true}); res.Total != 0 {
		t.Fatalf("restringida segue na pública: %+v", res)
	}
	if err := ix.Limpar(ctx, 2); err != nil {
		t.Fatal(err)
	}
	if res, _ := ix.Buscar(ctx, domain.Consulta{Q: "*", Pagina: 1, PorPagina: 10}, domain.FiltroBusca{}); res.Total != 1 {
		t.Fatalf("limpeza da geração anterior: %+v", res)
	}
}
