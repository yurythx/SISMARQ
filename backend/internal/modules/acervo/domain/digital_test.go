package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestMascararCPF(t *testing.T) {
	for entrada, esperado := range map[string]string{"123.456.789-01": "***.456.789-**", "12345678901": "***.456.789-**", "": ""} {
		if got, err := MascararCPF(entrada); err != nil || got != esperado {
			t.Errorf("%q: %q %v", entrada, got, err)
		}
	}
	var v *ErrValidacao
	if _, err := MascararCPF("123"); !errors.As(err, &v) {
		t.Fatalf("CPF curto: %v", err)
	}
}

func TestNovoDigital(t *testing.T) {
	d, err := NovoDigital(Digital{TipoDoc: " texto ", Equipamento: " Fujitsu fi-7160 ", DPIDeclarado: 300})
	if err != nil || d.TipoDoc != DocTexto || d.Situacao != DigitalEnviado || d.Equipamento != "Fujitsu fi-7160" || d.ID.String() == "" {
		t.Fatalf("válido: %+v %v", d, err)
	}
	for nome, x := range map[string]Digital{
		"tipo":              {TipoDoc: "FOTO", Equipamento: "e", DPIDeclarado: 300},
		"descrição longa":   {TipoDoc: DocTexto, Equipamento: "e", DPIDeclarado: 300, Descricao: strings.Repeat("x", 501)},
		"sem equipamento":   {TipoDoc: DocTexto, DPIDeclarado: 300},
		"DPI fora da faixa": {TipoDoc: DocPlanta, Equipamento: "e", DPIDeclarado: 50},
	} {
		var v *ErrValidacao
		if _, err := NovoDigital(x); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
	if DPIMinimo(DocTexto, 300, 600) != 300 || DPIMinimo(DocPlanta, 300, 600) != 600 {
		t.Fatal("DPI mínimo")
	}
	if PodeDigitalizar(Unidade{Fase: FaseIntermediaria}) != nil || !errors.Is(PodeDigitalizar(Unidade{Fase: FaseEliminado}), ErrEliminado) {
		t.Fatal("pode digitalizar")
	}
}

// O permanente fica para sempre; o que vai à eliminação, até um ano depois
// do prazo intermediário — e digitalizar não muda a trava do art. 42.
func TestRetencao(t *testing.T) {
	hoje := dia("2026-10-16")
	fim := dia("2030-12-31")
	velho := dia("2020-12-31")
	perm := Unidade{Fase: FasePermanente, Destinacao: DestinacaoGuardaPermanente}
	elim := Unidade{Fase: FaseIntermediaria, Destinacao: DestinacaoEliminacao}
	casos := []struct {
		nome string
		u    Unidade
		c    Calculo
		quer string
	}{
		{"permanente", perm, Calculo{FimIntermediaria: &fim}, "2126-10-16"},
		{"sem destinação definida", Unidade{Fase: FaseCorrente, Destinacao: DestinacaoIndefinida}, Calculo{FimIntermediaria: &fim}, "2126-10-16"},
		{"aguarda evento", elim, Calculo{}, "2126-10-16"},
		{"eliminação: prazo + 1 ano", elim, Calculo{FimIntermediaria: &fim}, "2031-12-31"},
		{"eliminação com prazo vencido: mínimo de 1 ano", elim, Calculo{FimIntermediaria: &velho}, "2027-10-16"},
	}
	for _, c := range casos {
		if got := Retencao(c.u, c.c, hoje).Format("2006-01-02"); got != c.quer {
			t.Errorf("%s: %s, esperado %s", c.nome, got, c.quer)
		}
	}
	if !errors.Is(PodeEliminar(perm, Serie{Destinacao: DestinacaoEliminacao}, false, hoje), ErrPermanente) {
		t.Fatal("digitalizado ou não, o permanente não se elimina (art. 42)")
	}
}
