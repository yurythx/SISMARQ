package infrastructure

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/typesense"
)

// BackendIndice é o que o índice do acervo usa do Typesense (o
// *typesense.Client).
type BackendIndice interface {
	EnsureCollection(ctx context.Context, s typesense.Schema) (bool, error)
	Import(ctx context.Context, collection string, docs []map[string]any) error
	DeleteByFilter(ctx context.Context, collection, filter string) (int, error)
	Search(ctx context.Context, collection string, p typesense.SearchParams) (typesense.SearchResult, error)
}

// Indice é a busca do acervo no Typesense (ADR 038), em duas coleções:
//   - a INTERNA tem todas as unidades, com o nível de acesso, o órgão e o
//     texto do OCR de todos os representantes; a API filtra pelo escopo de
//     quem pede, sempre no servidor;
//   - a PÚBLICA só tem o ostensivo, e nela o texto do OCR é só o dos
//     representantes liberados. Um erro de filtro na consulta pública não
//     tem o que vazar: o restrito nunca entra nessa coleção.
//
// O navegador nunca fala com o Typesense: a chave fica no servidor.
type Indice struct {
	ts               BackendIndice
	interna, publica string
}

// NovoIndice cria o índice com os nomes das coleções (mudar o schema pede
// nomes novos: a coleção nova nasce vazia e recebe a carga completa).
func NovoIndice(ts BackendIndice, interna, publica string) *Indice {
	return &Indice{ts: ts, interna: interna, publica: publica}
}

// Schema é o schema de uma das coleções.
func (ix *Indice) Schema(publico bool) typesense.Schema {
	nome, noIndex := ix.interna, false
	if publico {
		nome = ix.publica
	}
	fields := []typesense.Field{
		{Name: "numero", Type: "string"},
		{Name: "codigo", Type: "string"},
		{Name: "descritor", Type: "string", Optional: true},
		{Name: "assunto", Type: "string", Optional: true},
		{Name: "texto", Type: "string", Optional: true},
		{Name: domain.FacetaOrgao, Type: "string", Facet: true},
		{Name: domain.FacetaClasse, Type: "string", Facet: true},
		{Name: domain.FacetaEspecie, Type: "string", Facet: true},
		{Name: domain.FacetaFase, Type: "string", Facet: true},
		{Name: domain.FacetaDestinacao, Type: "string", Facet: true},
		{Name: domain.FacetaAno, Type: "int32[]", Facet: true},
		{Name: "ano_inicial", Type: "int32", Sort: true},
		{Name: "ano_final", Type: "int32", Index: &noIndex, Optional: true},
		{Name: "digitais", Type: "int32"},
		{Name: "updated_at", Type: "int64", Sort: true},
		{Name: "geracao", Type: "int64"},
	}
	if !publico {
		fields = append(fields,
			typesense.Field{Name: "detalhes", Type: "string", Optional: true},
			typesense.Field{Name: "orgao_id", Type: "string"},
			typesense.Field{Name: domain.FacetaNivel, Type: "string", Facet: true})
	}
	return typesense.Schema{Name: nome, Fields: fields}
}

// Garantir cria as coleções que faltam. Devolve true quando criou alguma
// (quem chama faz a carga completa).
func (ix *Indice) Garantir(ctx context.Context) (bool, error) {
	criou := false
	for _, publico := range []bool{false, true} {
		c, err := ix.ts.EnsureCollection(ctx, ix.Schema(publico))
		if err != nil {
			return false, fmt.Errorf("acervo: garantir a coleção %s: %w", ix.Schema(publico).Name, err)
		}
		criou = criou || c
	}
	return criou, nil
}

func documento(r domain.RegistroIndice, publico bool, geracao int64) map[string]any {
	u, ds := r.Unidade, r.Digitais
	if publico {
		ds = r.Liberados()
	}
	doc := map[string]any{
		"id": u.ID.String(), "numero": u.Numero, "codigo": u.CodigoClassificacao, "descritor": r.Descritor, "assunto": u.Assunto,
		"texto": domain.TextoIndice(ds), domain.FacetaOrgao: r.Orgao, domain.FacetaClasse: domain.Classe(u.CodigoClassificacao),
		domain.FacetaEspecie: string(u.Tipo), domain.FacetaFase: string(u.Fase), domain.FacetaDestinacao: string(u.Destinacao),
		domain.FacetaAno: domain.Anos(u.DataInicial, u.DataFinal), "ano_inicial": u.DataInicial.Year(), "ano_final": u.DataFinal.Year(),
		"digitais": len(ds), "updated_at": u.UpdatedAt.Unix(), "geracao": geracao,
	}
	if !publico {
		doc["detalhes"] = strings.TrimSpace(u.Especificacao + " " + u.Observacoes)
		doc["orgao_id"] = u.UnidadeID.String()
		doc[domain.FacetaNivel] = string(u.NivelAcesso)
	}
	return doc
}

// Gravar grava as unidades nas duas coleções (com a geração dada) e tira
// da pública as que deixaram de ser ostensivas.
func (ix *Indice) Gravar(ctx context.Context, regs []domain.RegistroIndice, geracao int64) error {
	internos, publicos, sair := make([]map[string]any, 0, len(regs)), []map[string]any{}, []string{}
	for _, r := range regs {
		internos = append(internos, documento(r, false, geracao))
		if r.Ostensivo() {
			publicos = append(publicos, documento(r, true, geracao))
		} else {
			sair = append(sair, r.Unidade.ID.String())
		}
	}
	if err := ix.ts.Import(ctx, ix.interna, internos); err != nil {
		return fmt.Errorf("acervo: gravar no índice interno: %w", err)
	}
	if err := ix.ts.Import(ctx, ix.publica, publicos); err != nil {
		return fmt.Errorf("acervo: gravar no índice público: %w", err)
	}
	if len(sair) > 0 {
		if _, err := ix.ts.DeleteByFilter(ctx, ix.publica, typesense.In("id", sair)); err != nil {
			return fmt.Errorf("acervo: tirar do índice público: %w", err)
		}
	}
	return nil
}

// Limpar apaga, nas duas coleções, o que ficou de gerações anteriores (o
// que saiu do acervo desde a última carga completa).
func (ix *Indice) Limpar(ctx context.Context, geracao int64) error {
	for _, c := range []string{ix.interna, ix.publica} {
		if _, err := ix.ts.DeleteByFilter(ctx, c, "geracao:<"+strconv.FormatInt(geracao, 10)); err != nil {
			return fmt.Errorf("acervo: limpar o índice %s: %w", c, err)
		}
	}
	return nil
}

// Buscar consulta a coleção pública ou a interna (com o escopo de quem
// pede no filtro).
func (ix *Indice) Buscar(ctx context.Context, c domain.Consulta, f domain.FiltroBusca) (domain.ResultadoBusca, error) {
	colecao, facetas, queryBy, pesos := ix.interna, domain.FacetasInternas, "numero,codigo,descritor,assunto,texto,detalhes", "4,4,3,3,1,1"
	if f.Publico {
		colecao, facetas, queryBy, pesos = ix.publica, domain.FacetasPublicas, "numero,codigo,descritor,assunto,texto", "4,4,3,3,1"
	}
	filtros := []string{}
	if f.Restrito && !f.Publico {
		ids := make([]string, len(f.Unidades))
		for i, id := range f.Unidades {
			ids[i] = id.String()
		}
		filtros = append(filtros, typesense.In("orgao_id", ids))
	}
	for _, k := range facetas {
		v, ok := c.Filtros[k]
		if !ok {
			continue
		}
		if k == domain.FacetaAno {
			ano, _ := strconv.Atoi(v)
			filtros = append(filtros, k+":="+strconv.Itoa(ano))
			continue
		}
		filtros = append(filtros, k+":="+typesense.Quote(v))
	}
	if c.SoDigital {
		filtros = append(filtros, "digitais:>0")
	}
	q, ordem := c.Q, "_text_match:desc,updated_at:desc"
	if q == "" {
		q, ordem = "*", "updated_at:desc"
	}
	res, err := ix.ts.Search(ctx, colecao, typesense.SearchParams{
		Q: q, QueryBy: queryBy, Weights: pesos, FilterBy: strings.Join(filtros, " && "), PerPage: c.PorPagina, Page: c.Pagina,
		FacetBy: strings.Join(facetas, ","), MaxFacetValues: 30, SortBy: ordem, ExcludeFields: "texto,detalhes",
		HighlightFields: "assunto,texto",
	})
	if err != nil {
		return domain.ResultadoBusca{}, err
	}
	out := domain.ResultadoBusca{Itens: make([]domain.ItemBusca, 0, len(res.Hits)), Total: int64(res.Found),
		Facetas: map[string][]domain.Faceta{}, Origem: domain.OrigemIndice, Ms: int64(res.SearchTime)}
	for _, h := range res.Hits {
		out.Itens = append(out.Itens, item(h))
	}
	for _, fc := range res.FacetCounts {
		vs := make([]domain.Faceta, 0, len(fc.Counts))
		for _, v := range fc.Counts {
			vs = append(vs, domain.Faceta{Valor: v.Value, Total: v.Count})
		}
		out.Facetas[fc.FieldName] = vs
	}
	return out, nil
}

func item(h typesense.Hit) domain.ItemBusca {
	d := h.Document
	id, _ := uuid.Parse(texto(d, "id"))
	it := domain.ItemBusca{ID: id, Numero: texto(d, "numero"), Tipo: texto(d, domain.FacetaEspecie), Codigo: texto(d, "codigo"),
		Descritor: texto(d, "descritor"), Assunto: texto(d, "assunto"), Orgao: texto(d, domain.FacetaOrgao),
		AnoInicial: numero(d, "ano_inicial"), AnoFinal: numero(d, "ano_final"), Fase: texto(d, domain.FacetaFase),
		Destinacao: texto(d, domain.FacetaDestinacao), NivelAcesso: texto(d, domain.FacetaNivel), Digitais: numero(d, "digitais")}
	for _, hl := range h.Highlights {
		if hl.Snippet != "" {
			it.Trecho = hl.Snippet
			break
		}
	}
	return it
}

func texto(d map[string]any, k string) string {
	s, _ := d[k].(string)
	return s
}

func numero(d map[string]any, k string) int {
	n, _ := d[k].(float64)
	return int(n)
}
