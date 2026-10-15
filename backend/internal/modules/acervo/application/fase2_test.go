package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

// Empréstimo e devolução (art. 12, II): um aberto por vez, atraso no painel.
func TestEmprestimoEDevolucao(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	d := e.criar(codElim)
	solicitante := dbtest.User(t, e.pool)

	if _, err := s.Emprestar(ctx, admin, d.ID, solicitante, " ", dia("2019-06-10")); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("sem finalidade: %v", err)
	}
	if _, err := s.Emprestar(ctx, agente(e.outra), d.ID, solicitante, "x", dia("2019-06-10")); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("fora do escopo: %v", err)
	}
	got, err := s.Emprestar(ctx, agente(e.unidade), d.ID, solicitante, "Instrução do processo 123", dia("2019-06-10"))
	if err != nil || len(got.Emprestimos) != 1 || !got.Emprestimos[0].Aberto() || got.Emprestimos[0].Solicitante == "" ||
		got.Movimentacoes[0].Tipo != domain.MovEmprestimo {
		t.Fatalf("empréstimo: %+v %v", got.Emprestimos, err)
	}
	if _, err := s.Emprestar(ctx, admin, d.ID, solicitante, "de novo", dia("2019-06-10")); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já emprestada: %v", err)
	}
	// Atraso: o relógio passa do prazo.
	e.agora = dia("2019-06-20")
	painel, err := s.Painel(ctx, agente(e.unidade), domain.Filtro{})
	if err != nil || painel[application.PainelEmprestadas] != 1 || painel[application.PainelAtrasadas] != 1 {
		t.Fatalf("painel: %v %v", painel, err)
	}
	atrasados, total, err := s.Emprestimos(ctx, agente(e.unidade), true, pagination.New(1, 10, 100))
	if err != nil || total != 1 || atrasados[0].Numero != d.Numero {
		t.Fatalf("atrasados: %+v %v", atrasados, err)
	}
	if _, total, _ := s.Emprestimos(ctx, agente(e.outra), false, pagination.New(1, 10, 100)); total != 0 {
		t.Fatal("outra unidade não vê o empréstimo")
	}
	eid := got.Emprestimos[0].ID
	if _, err := s.Devolver(ctx, admin, eid, strings.Repeat("x", 501)); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("observação longa: %v", err)
	}
	if _, err := s.Devolver(ctx, admin, uuid.New(), ""); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("empréstimo inexistente: %v", err)
	}
	got, err = s.Devolver(ctx, agente(e.unidade), eid, "Conferida")
	if err != nil || got.Emprestimos[0].Aberto() || got.Movimentacoes[0].Descricao != "Devolvida: Conferida" {
		t.Fatalf("devolução: %+v %v", got.Emprestimos, err)
	}
	if _, err := s.Devolver(ctx, admin, eid, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("devolver duas vezes: %v", err)
	}
	if got, err := s.Devolver(ctx, admin, eid, ""); err == nil {
		t.Fatalf("sem observação também conflita: %+v", got)
	}
}

// Guia de transferência (art. 29): remessa da origem, aceite que muda a fase,
// cancelamento, e as recusas.
func TestGuiaDeTransferencia(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	a, b := e.criar(codElim), e.criar(codElim)
	novo := e.criar(codPerm) // prazo corrente ainda correndo em 2020
	e.agora = dia("2020-06-01")

	if _, err := s.CriarGuia(ctx, admin, e.unidade, nil, ""); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("guia vazia: %v", err)
	}
	if _, err := s.CriarGuia(ctx, agente(e.outra), e.unidade, []uuid.UUID{a.ID}, ""); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("origem fora do escopo: %v", err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.outra, []uuid.UUID{a.ID}, ""); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("unidade de outra origem: %v", err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{novo.ID}, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("prazo corrente correndo: %v", err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{uuid.New()}, ""); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("unidade inexistente: %v", err)
	}
	g, err := s.CriarGuia(ctx, agente(e.unidade), e.unidade, []uuid.UUID{a.ID, b.ID}, "Remessa anual")
	if err != nil || !strings.HasPrefix(g.Numero, "GT-") || g.Situacao != domain.GuiaEnviada || g.TotalItens != 2 || len(g.Itens) != 2 {
		t.Fatalf("guia: %+v %v", g, err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{a.ID}, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já em outra guia: %v", err)
	}
	if _, err := s.Guia(ctx, agente(e.outra), g.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("guia fora do escopo: %v", err)
	}
	if guias, total, err := s.Guias(ctx, agente(e.unidade), pagination.New(1, 10, 100)); err != nil || total < 1 || guias[0].ID != g.ID {
		t.Fatalf("guias: %+v %v", guias, err)
	}
	// Emprestada no meio do caminho: o aceite recusa.
	solicitante := dbtest.User(t, e.pool)
	emp, err := s.Emprestar(ctx, admin, b.ID, solicitante, "consulta", dia("2020-06-05"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ReceberGuia(ctx, admin, g.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("aceite com unidade emprestada: %v", err)
	}
	if _, err := s.Devolver(ctx, admin, emp.Emprestimos[0].ID, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ReceberGuia(ctx, agente(e.outra), g.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("aceite fora do escopo: %v", err)
	}
	if _, err := s.ReceberGuia(ctx, agente(e.unidade, "acervo:read"), g.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("aceite sem gestão: %v", err)
	}
	rec, err := s.ReceberGuia(ctx, agente(e.unidade), g.ID)
	if err != nil || rec.Situacao != domain.GuiaRecebida || rec.RecebidaEm == nil || rec.Itens[0].Fase != domain.FaseIntermediaria {
		t.Fatalf("aceite: %+v %v", rec, err)
	}
	if _, err := s.CancelarGuia(ctx, admin, g.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("cancelar recebida: %v", err)
	}
	if _, err := s.ReceberGuia(ctx, admin, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("guia inexistente: %v", err)
	}
	// Cancelamento libera a unidade para outra guia.
	c := e.criar(codElim)
	g2, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{c.ID}, "")
	if err != nil {
		t.Fatal(err)
	}
	if can, err := s.CancelarGuia(ctx, agente(e.unidade), g2.ID); err != nil || can.Situacao != domain.GuiaCancelada {
		t.Fatalf("cancelamento: %+v %v", can, err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{c.ID}, ""); err != nil {
		t.Fatalf("depois do cancelamento: %v", err)
	}
	// Emprestada não entra na guia.
	d := e.criar(codElim)
	if _, err := s.Emprestar(ctx, admin, d.ID, solicitante, "x", dia("2020-06-05")); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{d.ID}, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("emprestada na guia: %v", err)
	}
}

// Importação do acervo existente: simula, aponta os erros por linha, recusa
// aplicar com erro ou com outro arquivo, e grava tudo ou nada.
func TestImportacao(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	sigla := "IMP" + strings.ToUpper(uuid.NewString()[:5])
	if _, err := e.pool.Exec(ctx, `UPDATE unidades SET sigla = $2 WHERE id = $1`, e.unidade, sigla); err != nil {
		t.Fatal(err)
	}
	predio, err := s.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: domain.LocalPredio, Nome: "Imp " + sigla, Ativo: true})
	if err != nil {
		t.Fatal(err)
	}
	cab := "tipo;codigo_classificacao;unidade;data_inicial;data_final;local\n"
	boa := "CAIXA;" + codElim + ";" + sigla + ";01/01/2015;31/12/2015;" + predio.Caminho + "\n"
	ruins := "CAIXA;0.0.00.00.00;" + sigla + ";01/01/2015;31/12/2015;\n" + // série inexistente
		"CAIXA;" + codElim + ";NAOEXISTE;01/01/2015;31/12/2015;\n" + // unidade produtora
		"CAIXA;" + codElim + ";" + sigla + ";01/01/2015;31/12/2015;Lugar nenhum\n" + // local
		"GAVETA;" + codElim + ";" + sigla + ";01/01/2015;31/12/2015;\n" + // validação
		"CAIXA;" + codElim + ";" + sigla + ";xx;31/12/2015;\n" // leitura

	sim, err := s.Importar(ctx, admin, cab+boa+ruins, "", false)
	if err != nil || sim.Aplicada || sim.Linhas != 1 || len(sim.Erros) != 5 || len(sim.Hash) != 64 || len(sim.Numeros) != 0 {
		t.Fatalf("simulação: %+v %v", sim, err)
	}
	if _, err := s.Importar(ctx, admin, cab+boa+ruins, sim.Hash, true); err != nil {
		t.Fatalf("aplicar com erro não grava, mas responde o relatório: %v", err)
	}
	if items, _, _ := s.List(ctx, admin, domain.Filtro{UnidadeID: &e.unidade}, pagination.New(1, 50, 100)); len(items) != 0 {
		t.Fatalf("nada pode ter sido gravado: %d", len(items))
	}
	if _, err := s.Importar(ctx, admin, cab+boa, sim.Hash, true); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("outro arquivo: %v", err)
	}
	if _, err := s.Importar(ctx, admin, "nada", "", false); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("planilha ilegível: %v", err)
	}
	// Fora do escopo: a linha vira erro.
	if out, err := s.Importar(ctx, agente(e.outra), cab+boa, "", false); err != nil || len(out.Erros) != 1 {
		t.Fatalf("fora do escopo: %+v %v", out, err)
	}
	ok, err := s.Importar(ctx, agente(e.unidade), cab+boa+boa, "", false)
	if err != nil || len(ok.Erros) != 0 {
		t.Fatalf("simulação limpa: %+v %v", ok, err)
	}
	ap, err := s.Importar(ctx, agente(e.unidade), cab+boa+boa, ok.Hash, true)
	if err != nil || !ap.Aplicada || ap.Linhas != 2 || len(ap.Numeros) != 2 {
		t.Fatalf("aplicação: %+v %v", ap, err)
	}
	items, _, _ := s.List(ctx, admin, domain.Filtro{UnidadeID: &e.unidade}, pagination.New(1, 50, 100))
	if len(items) != 2 || items[0].LocalCaminho != predio.Caminho {
		t.Fatalf("importadas: %+v", items)
	}
	if got, err := s.Get(ctx, admin, items[0].ID); err != nil || got.Movimentacoes[0].Tipo != domain.MovImportacao {
		t.Fatalf("histórico da importação: %+v %v", got.Movimentacoes, err)
	}
}
