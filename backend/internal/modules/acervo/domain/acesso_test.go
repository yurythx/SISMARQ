package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPodePublicarDigital(t *testing.T) {
	pub := Unidade{NivelAcesso: AcessoPublico, Fase: FasePermanente}
	pronto := Digital{Situacao: DigitalPronto}
	if err := PodePublicarDigital(pub, pronto); err != nil {
		t.Fatal(err)
	}
	casos := []struct {
		u    Unidade
		d    Digital
		want error
	}{
		{pub, Digital{Situacao: DigitalEnviado}, ErrDigitalSituacao},
		{Unidade{NivelAcesso: AcessoPublico, Fase: FaseEliminado}, pronto, ErrEliminado},
		{Unidade{NivelAcesso: AcessoRestrito, Fase: FasePermanente}, pronto, ErrNaoOstensivo},
		{Unidade{NivelAcesso: AcessoSigiloso, Fase: FaseCorrente}, pronto, ErrNaoOstensivo},
	}
	for _, c := range casos {
		if err := PodePublicarDigital(c.u, c.d); !errors.Is(err, c.want) {
			t.Errorf("%+v: %v, esperado %v", c.u, err, c.want)
		}
	}
}

func TestNormalizarConsulta(t *testing.T) {
	c, err := NormalizarConsulta(Consulta{Q: "  ofício ", Filtros: map[string]string{FacetaAno: "2015", FacetaFase: " ", FacetaOrgao: "SEMAD"}}, true, 50)
	if err != nil || c.Q != "ofício" || c.Pagina != 1 || c.PorPagina != 20 || len(c.Filtros) != 2 || c.Filtros[FacetaOrgao] != "SEMAD" {
		t.Fatalf("normalizada: %+v %v", c, err)
	}
	if c, _ := NormalizarConsulta(Consulta{Pagina: 3, PorPagina: 10}, false, 5); c.PorPagina != 5 || c.Pagina != 3 {
		t.Fatalf("teto da página: %+v", c)
	}
	if c, _ := NormalizarConsulta(Consulta{PorPagina: 7}, false, 50); c.PorPagina != 7 {
		t.Fatalf("página válida: %+v", c)
	}
	var v *ErrValidacao
	for _, ruim := range []Consulta{
		{Q: strings.Repeat("a", 201)},
		{Filtros: map[string]string{"local": "x"}},
		{Filtros: map[string]string{FacetaAno: "dois mil"}},
		{Filtros: map[string]string{FacetaNivel: "RESTRITO"}}, // a consulta pública não filtra por nível
	} {
		if _, err := NormalizarConsulta(ruim, true, 50); !errors.As(err, &v) {
			t.Errorf("%+v: %v", ruim, err)
		}
	}
	if _, err := NormalizarConsulta(Consulta{Filtros: map[string]string{FacetaNivel: "RESTRITO"}}, false, 50); err != nil {
		t.Fatalf("a interna filtra por nível: %v", err)
	}
}

func TestClasseEAnos(t *testing.T) {
	if Classe("2.0.03.02") != "2" || Classe("9") != "9" {
		t.Fatal("classe")
	}
	d := func(s string) time.Time { x, _ := time.Parse("2006-01-02", s); return x }
	if a := Anos(d("2015-03-01"), d("2017-01-01")); len(a) != 3 || a[0] != 2015 || a[2] != 2017 {
		t.Fatalf("anos: %v", a)
	}
	if a := Anos(d("1700-01-01"), d("2020-01-01")); len(a) != 2 || a[1] != 2020 {
		t.Fatalf("intervalo largo fica nas pontas: %v", a)
	}
}

func TestRegistroIndice(t *testing.T) {
	agora := time.Now()
	r := RegistroIndice{Unidade: Unidade{ID: uuid.New(), Numero: "UA-1", NivelAcesso: AcessoPublico, Fase: FasePermanente, Tipo: TipoCaixa},
		Orgao: "SEMAD", Descritor: "Ofícios", Digitais: []DigitalIndice{
			{Numero: "RD-1", Texto: "liberado", PublicadoEm: &agora},
			{Numero: "RD-2", Texto: "interno"},
			{Numero: "RD-3"},
		}}
	if !r.Ostensivo() || len(r.Liberados()) != 1 || r.Liberados()[0].Numero != "RD-1" {
		t.Fatalf("liberados: %+v", r.Liberados())
	}
	if got := TextoIndice(r.Digitais); got != "liberado\ninterno" {
		t.Fatalf("texto: %q", got)
	}
	f := NovaFicha(r)
	if f.Numero != "UA-1" || f.Orgao != "SEMAD" || f.Descritor != "Ofícios" || len(f.Digitais) != 1 || f.Recolhida || f.Tipo != "CAIXA" {
		t.Fatalf("ficha: %+v", f)
	}
	r.Unidade.Fase = FaseEliminado
	if r.Ostensivo() {
		t.Fatal("eliminada não é consulta pública")
	}
	r.Unidade.Fase, r.Unidade.NivelAcesso = FasePermanente, AcessoRestrito
	if r.Ostensivo() {
		t.Fatal("restrita não é consulta pública")
	}
	// O texto para no limite, sem cortar um caractere ao meio.
	longo := []DigitalIndice{{Texto: "a" + strings.Repeat("é", TextoMaxIndice/2+10)}, {Texto: "depois"}}
	if got := TextoIndice(longo); len(got) > TextoMaxIndice || strings.Contains(got, "depois") || !strings.HasSuffix(got, "é") {
		t.Fatalf("limite: %d", len(got))
	}
}

func TestNovaExposicao(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	x, err := NovaExposicao(Exposicao{Titulo: " Memória ", Slug: "memoria", Itens: []ItemExposicao{{DigitalID: a, Legenda: " Planta "}, {DigitalID: b}}})
	if err != nil || x.Titulo != "Memória" || x.Itens[0].Legenda != "Planta" || x.Itens[1].Ordem != 1 {
		t.Fatalf("exposição: %+v %v", x, err)
	}
	var v *ErrValidacao
	muitos := make([]ItemExposicao, ItensMax+1)
	for i := range muitos {
		muitos[i].DigitalID = uuid.New()
	}
	for _, ruim := range []Exposicao{
		{Slug: "x"},
		{Titulo: strings.Repeat("a", 201), Slug: "x"},
		{Titulo: "t", Slug: "x", Resumo: strings.Repeat("a", 501)},
		{Titulo: "t", Slug: "x", Texto: strings.Repeat("a", 20001)},
		{Titulo: "!!!"},
		{Titulo: "t", Slug: "x", Itens: muitos},
		{Titulo: "t", Slug: "x", Itens: []ItemExposicao{{DigitalID: a}, {DigitalID: a}}},
		{Titulo: "t", Slug: "x", Itens: []ItemExposicao{{DigitalID: a, Legenda: strings.Repeat("a", 1001)}}},
	} {
		if _, err := NovaExposicao(ruim); !errors.As(err, &v) {
			t.Errorf("%.40v: %v", ruim, err)
		}
	}
}

func TestPublicarExposicao(t *testing.T) {
	agora := time.Now()
	x := Exposicao{Itens: []ItemExposicao{{Numero: "RD-1"}, {Numero: "RD-2", Publicado: true}}}
	if err := PodePublicarExposicao(x); err != nil {
		t.Fatal(err)
	}
	if v := x.SoPublicados(); len(v.Itens) != 1 || v.Itens[0].Numero != "RD-2" {
		t.Fatalf("visita: %+v", v.Itens)
	}
	if err := PodePublicarExposicao(Exposicao{Itens: []ItemExposicao{{}}}); !errors.Is(err, ErrExposicaoVazia) {
		t.Fatalf("vazia: %v", err)
	}
	x.PublicadaEm = &agora
	if err := PodePublicarExposicao(x); !errors.Is(err, ErrExposicaoSituacao) {
		t.Fatalf("já publicada: %v", err)
	}
}
