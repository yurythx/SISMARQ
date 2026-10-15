package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNovoEmprestimo(t *testing.T) {
	hoje := dia("2026-10-15")
	u, s := uuid.New(), uuid.New()
	e, err := NovoEmprestimo(u, s, " Instrução do processo 123 ", dia("2026-10-30"), hoje)
	if err != nil || e.Finalidade != "Instrução do processo 123" || !e.Aberto() || e.Atrasado(hoje) || !e.Atrasado(dia("2026-10-31")) {
		t.Fatalf("empréstimo: %+v %v", e, err)
	}
	agora := time.Now()
	e.DevolvidoEm = &agora
	if e.Aberto() || e.Atrasado(dia("2027-01-01")) {
		t.Fatal("devolvido não está aberto nem atrasado")
	}
	for nome, c := range map[string]struct {
		s   uuid.UUID
		fin string
		ate time.Time
	}{
		"sem solicitante":   {uuid.Nil, "x", dia("2026-10-20")},
		"sem finalidade":    {s, " ", dia("2026-10-20")},
		"finalidade longa":  {s, strings.Repeat("x", 501), dia("2026-10-20")},
		"devolução passada": {s, "x", dia("2026-10-14")},
		"mais de um ano":    {s, "x", dia("2027-10-16")},
	} {
		var v *ErrValidacao
		if _, err := NovoEmprestimo(u, c.s, c.fin, c.ate, hoje); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
}

func TestNovaGuia(t *testing.T) {
	o, a, b := uuid.New(), uuid.New(), uuid.New()
	g, ids, err := NovaGuia(o, []uuid.UUID{a, b, a, uuid.Nil}, " remessa de 2026 ")
	if err != nil || len(ids) != 2 || g.Situacao != GuiaEnviada || g.Observacao != "remessa de 2026" || g.OrigemID != o {
		t.Fatalf("guia: %+v %v %v", g, ids, err)
	}
	muitas := make([]uuid.UUID, MaxItensGuia+1)
	for i := range muitas {
		muitas[i] = uuid.New()
	}
	for nome, c := range map[string]struct {
		origem uuid.UUID
		ids    []uuid.UUID
		obs    string
	}{
		"sem origem":       {uuid.Nil, []uuid.UUID{a}, ""},
		"vazia":            {o, []uuid.UUID{uuid.Nil}, ""},
		"acima do limite":  {o, muitas, ""},
		"observação longa": {o, []uuid.UUID{a}, strings.Repeat("x", 1001)},
	} {
		if _, _, err := NovaGuia(c.origem, c.ids, c.obs); err == nil {
			t.Errorf("%s: aceita", nome)
		}
	}
}

func TestLerPlanilha(t *testing.T) {
	csv := "\uFEFFtipo;codigo_classificacao;assunto;unidade;data_inicial;data_final;quantidade;data_evento;local;nivel_acesso\n" +
		"Dossiê;2.0.01.00.01;Organogramas;SEMAD;01/01/2015;2015-12-31;2;10/03/2016;Sede > Estante 1;Público\n" +
		";;;;;\n" +
		"CAIXA;2.0.01.00.01;;SEMAD;32/01/2015;31/12/2015;;;;\n" +
		"CAIXA;2.0.01.00.01;;SEMAD;01/01/2015;x;;;;\n" +
		"CAIXA;2.0.01.00.01;;SEMAD;01/01/2015;31/12/2015;dois;;;\n" +
		"CAIXA;2.0.01.00.01;;SEMAD;01/01/2015;31/12/2015;;ontem;;\n" +
		"CAIXA;2.0.01.00.01;;;01/01/2015;31/12/2015;;;;\n"
	linhas, erros, err := LerPlanilha(csv)
	if err != nil || len(linhas) != 1 || len(erros) != 5 {
		t.Fatalf("leitura: %+v %+v %v", linhas, erros, err)
	}
	l := linhas[0]
	if l.Linha != 2 || l.Unidade.Tipo != "DOSSIE" || l.UnidadeRef != "SEMAD" || l.LocalRef != "Sede > Estante 1" || l.Unidade.Quantidade != 2 ||
		!l.Unidade.DataFinal.Equal(dia("2015-12-31")) || l.Unidade.DataEvento == nil || l.Unidade.NivelAcesso != "PUBLICO" {
		t.Fatalf("linha: %+v", l)
	}
	if erros[0].Linha != 4 || !strings.Contains(erros[0].Motivo, "data_inicial") {
		t.Fatalf("erro aponta a linha: %+v", erros[0])
	}
	// Separador "," também vale.
	if linhas, _, err := LerPlanilha("tipo,codigo_classificacao,unidade,data_inicial,data_final\nCAIXA,2.0,X,01/01/2015,31/12/2015\n"); err != nil || len(linhas) != 1 {
		t.Fatalf("vírgula: %v", err)
	}
	for nome, ruim := range map[string]string{
		"vazia":           "",
		"sem coluna":      "tipo;unidade\nCAIXA;X\n",
		"só cabeçalho":    "tipo;codigo_classificacao;unidade;data_inicial;data_final\n",
		"aspas quebradas": "tipo;codigo_classificacao;unidade;data_inicial;data_final\n\"CAIXA;x\n",
	} {
		if _, _, err := LerPlanilha(ruim); !errors.Is(err, ErrPlanilha) {
			t.Errorf("%s: %v", nome, err)
		}
	}
	grande := "tipo;codigo_classificacao;unidade;data_inicial;data_final\n" + strings.Repeat("CAIXA;2.0;X;01/01/2015;31/12/2015\n", MaxLinhasImportacao+1)
	if _, _, err := LerPlanilha(grande); !errors.Is(err, ErrPlanilha) {
		t.Fatalf("acima do limite: %v", err)
	}
}
