package domain

import (
	"errors"
	"strings"
	"testing"
)

func redacao(v int, situacao, em string) Redacao {
	r := Redacao{Versao: v, Ementa: "Ementa " + string(rune('0'+v)), Texto: "Texto", Motivo: "M", Situacao: situacao}
	if em != "" {
		d := dia(em)
		r.PublicadaEm = &d
	}
	return r
}

func TestRedacaoENumero(t *testing.T) {
	r := Redacao{Ementa: " Recolhimento ", Texto: " Art. 1º ", Motivo: " Redação original "}
	if err := ValidarRedacao(&r); err != nil || r.Ementa != "Recolhimento" || r.Texto != "Art. 1º" {
		t.Fatalf("válida: %+v %v", r, err)
	}
	for nome, x := range map[string]Redacao{
		"sem ementa":   {Texto: "t", Motivo: "m"},
		"ementa longa": {Ementa: strings.Repeat("x", 501), Texto: "t", Motivo: "m"},
		"sem texto":    {Ementa: "e", Motivo: "m"},
		"texto longo":  {Ementa: "e", Texto: strings.Repeat("x", MaxTextoNorma+1), Motivo: "m"},
		"sem motivo":   {Ementa: "e", Texto: "t"},
	} {
		var v *ErrValidacao
		if err := ValidarRedacao(&x); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
	if n, err := ValidarNumero(" 3/2026 "); err != nil || n != "3/2026" {
		t.Fatalf("número: %q %v", n, err)
	}
	for _, n := range []string{"", "IN 3/2026", "3-2026", "3/26"} {
		if _, err := ValidarNumero(n); err == nil {
			t.Errorf("número %q aceito", n)
		}
	}
}

// A norma muda com o tempo: cada redação vigora até a seguinte e a última,
// até a revogação.
func TestMontarHistorico(t *testing.T) {
	revogada := dia("2028-01-10")
	n := Norma{Situacao: NormaRevogada, RevogadaEm: &revogada, Redacoes: []Redacao{
		redacao(4, RedacaoRascunho, ""), redacao(3, RedacaoPublicada, "2027-06-01"),
		redacao(2, RedacaoPublicada, "2026-12-01"), redacao(1, RedacaoPublicada, "2026-03-01"),
	}}
	n.Montar()
	if n.Versao != 3 || n.Ementa != "Ementa 3" || !n.PublicadaEm.Equal(dia("2027-06-01")) || !n.TemRascunho || n.TotalRedacoes != 3 {
		t.Fatalf("montada: %+v", n)
	}
	if !n.Redacoes[1].VigenteAte.Equal(revogada) || !n.Redacoes[2].VigenteAte.Equal(dia("2027-06-01")) || !n.Redacoes[3].VigenteAte.Equal(dia("2026-12-01")) {
		t.Fatalf("vigências: %+v", n.Redacoes)
	}
	if n.Vigente().Versao != 3 || n.Rascunho().Versao != 4 {
		t.Fatal("vigente e rascunho")
	}
	pub := n.Publicas()
	if len(pub.Redacoes) != 3 || pub.TemRascunho || len(n.Redacoes) != 4 {
		t.Fatalf("públicas: %+v", pub)
	}
	nova := Norma{Situacao: NormaRascunho, Redacoes: []Redacao{redacao(1, RedacaoRascunho, "")}}
	nova.Montar()
	if nova.Versao != 0 || nova.Ementa != "Ementa 1" || nova.Vigente() != nil || nova.PublicadaEm != nil {
		t.Fatalf("nunca publicada: %+v", nova)
	}
	if (Norma{}).Rascunho() != nil {
		t.Fatal("sem redações")
	}
}

func TestPublicarERevogar(t *testing.T) {
	hoje := dia("2026-10-15")
	vigente := Norma{Situacao: NormaVigente, Redacoes: []Redacao{redacao(2, RedacaoRascunho, ""), redacao(1, RedacaoPublicada, "2026-03-01")}}
	vigente.Montar()
	if err := PodePublicarNorma(vigente, dia("2026-10-01"), hoje, "DIORONDON-e nº 6.400"); err != nil {
		t.Fatal(err)
	}
	semRascunho := Norma{Situacao: NormaVigente, Redacoes: []Redacao{redacao(1, RedacaoPublicada, "2026-03-01")}}
	semRascunho.Montar()
	revogada := vigente
	revogada.Situacao = NormaRevogada
	var v *ErrValidacao
	casos := []struct {
		nome string
		err  error
	}{
		{"revogada", PodePublicarNorma(revogada, hoje, hoje, "")},
		{"sem rascunho", PodePublicarNorma(semRascunho, hoje, hoje, "")},
	}
	for _, c := range casos {
		if c.err == nil {
			t.Errorf("%s: aceita", c.nome)
		}
	}
	if !errors.Is(casos[0].err, ErrNormaSituacao) || !errors.Is(casos[1].err, ErrSemRascunho) {
		t.Fatal("erros de situação")
	}
	for nome, err := range map[string]error{
		"antes da vigente": PodePublicarNorma(vigente, dia("2026-01-01"), hoje, ""),
		"futura":           PodePublicarNorma(vigente, dia("2026-10-16"), hoje, ""),
		"diário longo":     PodePublicarNorma(vigente, hoje, hoje, strings.Repeat("x", 201)),
	} {
		if !errors.As(err, &v) {
			t.Errorf("publicar %s: %v", nome, err)
		}
	}
	if err := PodeRevogar(semRascunho, dia("2026-09-01"), hoje, "IN 9/2026"); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(PodeRevogar(revogada, hoje, hoje, "x"), ErrNormaSituacao) {
		t.Fatal("revogar revogada")
	}
	for nome, err := range map[string]error{
		"sem ato":          PodeRevogar(semRascunho, hoje, hoje, " "),
		"antes da vigente": PodeRevogar(semRascunho, dia("2026-01-01"), hoje, "x"),
		"futura":           PodeRevogar(semRascunho, dia("2026-10-16"), hoje, "x"),
	} {
		if !errors.As(err, &v) {
			t.Errorf("revogar %s: %v", nome, err)
		}
	}
}
