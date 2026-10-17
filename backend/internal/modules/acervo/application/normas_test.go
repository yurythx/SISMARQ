package application_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
)

// numeroNum sorteia um número N/AAAA (o número é único no banco, que os
// testes compartilham).
func numeroNum() string { return fmt.Sprintf("%d/%d", rand.IntN(9999)+1, 2100+rand.IntN(800)) }

// A instrução normativa nasce, é publicada, muda (nova redação) e é
// revogada — e o público consulta cada redação no seu período (art. 54).
func TestNormaAoLongoDoTempo(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	arquivo := arquivoPublico(uuid.Nil)
	numero := numeroNum()

	// Só o Arquivo Público escreve e lê a gestão.
	if _, err := s.CriarNorma(ctx, agente(e.unidade), numero, domain.Redacao{Ementa: "e", Texto: "t"}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("fora do Arquivo: %v", err)
	}
	if _, err := s.CriarNorma(ctx, arquivo, "IN 1", domain.Redacao{Ementa: "e", Texto: "t"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("número inválido: %v", err)
	}
	if _, err := s.CriarNorma(ctx, arquivo, numero, domain.Redacao{Texto: "t"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("sem ementa: %v", err)
	}
	e.agora = dia("2026-03-01")
	n, err := s.CriarNorma(ctx, arquivo, numero, domain.Redacao{Ementa: "Procedimentos de recolhimento", Texto: "Art. 1º Recolher em março."})
	if err != nil || n.Situacao != domain.NormaRascunho || !n.TemRascunho || n.Redacoes[0].Motivo != application.MotivoOriginal {
		t.Fatalf("criação: %+v %v", n, err)
	}
	if _, err := s.CriarNorma(ctx, arquivo, numero, domain.Redacao{Ementa: "e", Texto: "t"}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("número repetido: %v", err)
	}
	// Nunca publicada: o público não vê.
	if _, err := s.NormaPublica(ctx, n.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("rascunho público: %v", err)
	}
	// Ajuste do rascunho (mesma versão).
	if n, err = s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "Procedimentos de recolhimento", Texto: "Art. 1º Recolher em março e abril.",
		Motivo: application.MotivoOriginal}); err != nil || len(n.Redacoes) != 1 {
		t.Fatalf("rascunho: %+v %v", n, err)
	}
	if _, err := s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "e"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("rascunho inválido: %v", err)
	}
	if _, err := s.PublicarNorma(ctx, arquivo, n.ID, application.Publicacao{Em: dia("2026-03-02")}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("publicação futura: %v", err)
	}
	if n, err = s.PublicarNorma(ctx, arquivo, n.ID, application.Publicacao{Em: dia("2026-03-01"), Diario: "DIORONDON-e nº 6.100"}); err != nil ||
		n.Situacao != domain.NormaVigente || n.Versao != 1 || n.TemRascunho {
		t.Fatalf("publicação: %+v %v", n, err)
	}
	if _, err := s.PublicarNorma(ctx, arquivo, n.ID, application.Publicacao{Em: dia("2026-03-01")}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("publicar sem rascunho: %v", err)
	}
	if _, err := s.DescartarRascunho(ctx, arquivo, n.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("descartar sem rascunho: %v", err)
	}

	// A norma é alterada: nova redação, publicada depois; a anterior fica no
	// histórico, vigente até a nova.
	e.agora = dia("2027-06-01")
	if n, err = s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "Procedimentos de recolhimento", Texto: "Art. 1º Recolher em junho.",
		Motivo: "Alterada pela IN 2/2027"}); err != nil || n.Versao != 1 || !n.TemRascunho || n.Redacoes[0].Versao != 2 {
		t.Fatalf("nova redação: %+v %v", n, err)
	}
	// O público segue vendo só a redação em vigor.
	if pub, err := s.NormaPublica(ctx, n.ID); err != nil || len(pub.Redacoes) != 1 || pub.TemRascunho {
		t.Fatalf("público durante a alteração: %+v %v", pub, err)
	}
	if n, err = s.PublicarNorma(ctx, arquivo, n.ID, application.Publicacao{Em: dia("2027-06-01")}); err != nil || n.Versao != 2 {
		t.Fatalf("publicação da alteração: %+v %v", n, err)
	}
	pub, err := s.NormaPublica(ctx, n.ID)
	if err != nil || len(pub.Redacoes) != 2 || !pub.Redacoes[1].VigenteAte.Equal(dia("2027-06-01")) || pub.Redacoes[0].VigenteAte != nil {
		t.Fatalf("histórico público: %+v %v", pub, err)
	}
	// A redação publicada não muda no banco.
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_norma_redacoes SET texto = 'x' WHERE norma_id = $1 AND versao = 1`, n.ID); err == nil {
		t.Fatal("o banco alterou a redação publicada")
	}
	if _, err := e.pool.Exec(ctx, `DELETE FROM acervo_norma_redacoes WHERE norma_id = $1 AND versao = 1`, n.ID); err == nil {
		t.Fatal("o banco apagou a redação publicada")
	}

	// Rascunho descartado; revogação descarta o pendente.
	if _, err := s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "e", Texto: "t", Motivo: "m"}); err != nil {
		t.Fatal(err)
	}
	if n, err = s.DescartarRascunho(ctx, arquivo, n.ID); err != nil || n.TemRascunho || n.ID == uuid.Nil {
		t.Fatalf("descarte: %+v %v", n, err)
	}
	if _, err := s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "e", Texto: "t", Motivo: "m"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.RevogarNorma(ctx, arquivo, n.ID, e.agora, " "); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("revogação sem ato: %v", err)
	}
	if n, err = s.RevogarNorma(ctx, arquivo, n.ID, e.agora, "Revogada pela IN 5/2027"); err != nil || n.Situacao != domain.NormaRevogada || n.TemRascunho {
		t.Fatalf("revogação: %+v %v", n, err)
	}
	if _, err := s.SalvarRascunho(ctx, arquivo, n.ID, domain.Redacao{Ementa: "e", Texto: "t", Motivo: "m"}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("alterar revogada: %v", err)
	}
	if pub, err := s.NormaPublica(ctx, n.ID); err != nil || pub.Situacao != domain.NormaRevogada || !pub.Redacoes[0].VigenteAte.Equal(e.agora) {
		t.Fatalf("revogada consultável: %+v %v", pub, err)
	}

	// Rascunho de norma nunca publicada: descartar apaga a norma.
	outra, err := s.CriarNorma(ctx, arquivo, numeroNum(), domain.Redacao{Ementa: "Provisória", Texto: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if gone, err := s.DescartarRascunho(ctx, arquivo, outra.ID); err != nil || gone.ID != uuid.Nil {
		t.Fatalf("descartar a nunca publicada: %+v %v", gone, err)
	}
	if _, err := s.Norma(ctx, arquivo, outra.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("apagada: %v", err)
	}

	// Listas: a gestão vê rascunhos; o público, só as publicadas.
	rasc, err := s.CriarNorma(ctx, arquivo, numeroNum(), domain.Redacao{Ementa: "Ainda em estudo", Texto: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if _, total, err := s.Normas(ctx, arquivo, "", pagination.New(1, 50, 100)); err != nil || total < 2 {
		t.Fatalf("gestão: %d %v", total, err)
	}
	if items, _, err := s.Normas(ctx, arquivo, rasc.Numero, pagination.New(1, 50, 100)); err != nil || len(items) != 1 || !items[0].TemRascunho {
		t.Fatalf("busca pelo número: %+v %v", items, err)
	}
	items, _, err := s.NormasPublicas(ctx, "Recolher em junho", pagination.New(1, 50, 100))
	if err != nil || len(items) == 0 {
		t.Fatalf("busca no texto: %+v %v", items, err)
	}
	achou := false
	for _, it := range items {
		if it.ID == n.ID {
			achou = true
			break
		}
	}
	if !achou {
		t.Fatalf("busca no texto não trouxe a norma %s: %+v", n.ID, items)
	}
	if items, _, _ := s.NormasPublicas(ctx, "Ainda em estudo", pagination.New(1, 50, 100)); len(items) != 0 {
		t.Fatal("rascunho na lista pública")
	}
	if _, _, err := s.Normas(ctx, agente(e.unidade), "", pagination.New(1, 50, 100)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("lista fora do Arquivo: %v", err)
	}
	if _, err := s.Norma(ctx, agente(e.unidade), n.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("ficha fora do Arquivo: %v", err)
	}
	if _, _, err := s.NormasPublicas(ctx, strings.Repeat("x", 150), pagination.New(1, 50, 100)); err != nil {
		t.Fatalf("busca longa: %v", err)
	}
	if _, err := s.NormaPublica(ctx, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	if _, err := s.SalvarRascunho(ctx, agente(e.unidade), n.ID, domain.Redacao{Ementa: "e", Texto: "t", Motivo: "m"}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("alterar fora do Arquivo: %v", err)
	}
}
