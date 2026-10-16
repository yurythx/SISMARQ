package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
)

func fundoFamilia() domain.Privado {
	return domain.Privado{Titular: "Família Silva", TipoTitular: "PF", Documento: "123.456.789-00",
		Descricao: "Fotografias e cartas da fundação de Rondonópolis (1953-1970)", Mensuracao: "3 caixas"}
}

// Do cadastro ao decreto, a guarda e o direito de preferência (arts. 55 a
// 58).
func TestArquivoPrivado(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-15")
	s := e.real()
	arquivo := arquivoPublico(uuid.Nil)

	if _, err := s.CadastrarPrivado(ctx, agente(e.unidade), fundoFamilia()); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público: %v", err)
	}
	if _, err := s.CadastrarPrivado(ctx, arquivo, domain.Privado{TipoTitular: "PF"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("cadastro inválido: %v", err)
	}
	p, err := s.CadastrarPrivado(ctx, arquivo, fundoFamilia())
	if err != nil || p.Situacao != domain.PrivadoCadastrado || p.Numero == "" || len(p.Eventos) != 1 {
		t.Fatalf("cadastro: %+v %v", p, err)
	}

	// Decreto só depois do parecer favorável.
	if _, err := s.Declarar(ctx, arquivo, p.ID, "1", "", e.agora); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("decreto sem parecer: %v", err)
	}
	pa := application.Parecer{Comissao: "Comissão Especial: historiadora A, arquivista B", Texto: "Conjunto único da fundação",
		Favoravel: true, Em: dia("2026-09-01")}
	if p, err = s.RegistrarParecer(ctx, arquivo, p.ID, pa); err != nil || p.Situacao != domain.PrivadoParecer || !*p.ParecerFavoravel {
		t.Fatalf("parecer: %+v %v", p, err)
	}
	if _, err := s.RegistrarParecer(ctx, arquivo, p.ID, pa); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("segundo parecer: %v", err)
	}
	if _, err := s.NotificarAlienacao(ctx, arquivo, p.ID, domain.Alienacao{NotificadoEm: e.agora}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("alienação antes do decreto: %v", err)
	}
	if p, err = s.Declarar(ctx, arquivo, p.ID, "Decreto 10.000/2026", "DIORONDON-e", dia("2026-10-01")); err != nil || p.Situacao != domain.PrivadoDeclarado {
		t.Fatalf("decreto: %+v %v", p, err)
	}

	// Cadastro corrigido (o § 3º e o tipo não mudam por aqui).
	dados := fundoFamilia()
	dados.AcessoAutorizado, dados.AcessoCondicoes, dados.TipoTitular = true, "Consulta com agendamento", "PJ"
	if p, err = s.AlterarPrivado(ctx, arquivo, p.ID, dados); err != nil || !p.AcessoAutorizado || p.TipoTitular != domain.TitularPF {
		t.Fatalf("alteração: %+v %v", p, err)
	}
	dados.Titular = ""
	if _, err := s.AlterarPrivado(ctx, arquivo, p.ID, dados); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("alteração inválida: %v", err)
	}

	// Depósito revogável e revogação (§ 5º).
	if p, err = s.MudarGuarda(ctx, arquivo, p.ID, " depositado ", "Termo de depósito 1/2026", e.agora); err != nil || p.Guarda != domain.GuardaDepositado {
		t.Fatalf("depósito: %+v %v", p, err)
	}
	if p, err = s.MudarGuarda(ctx, arquivo, p.ID, domain.GuardaTitular, "A pedido da família", e.agora); err != nil || p.Guarda != domain.GuardaTitular {
		t.Fatalf("revogação do depósito: %+v %v", p, err)
	}
	if _, err := s.MudarGuarda(ctx, arquivo, p.ID, domain.GuardaTitular, "x", e.agora); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("guarda repetida: %v", err)
	}

	// Alienação: 60 dias para o município (art. 58).
	if p, err = s.NotificarAlienacao(ctx, arquivo, p.ID, domain.Alienacao{NotificadoEm: dia("2026-10-01"), Adquirente: "Colecionador X"}); err != nil ||
		len(p.Alienacoes) != 1 || !p.Alienacoes[0].PrazoFinal.Equal(dia("2026-11-30")) || p.Alienacoes[0].Vencida {
		t.Fatalf("notificação: %+v %v", p.Alienacoes, err)
	}
	if _, err := s.NotificarAlienacao(ctx, arquivo, p.ID, domain.Alienacao{NotificadoEm: e.agora}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("duas pendentes: %v", err)
	}
	al := p.Alienacoes[0].ID
	if _, err := s.DecidirAlienacao(ctx, arquivo, p.ID, uuid.New(), "INTERESSE", "x"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("alienação inexistente: %v", err)
	}
	if _, err := s.DecidirAlienacao(ctx, arquivo, p.ID, al, "talvez", "x"); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("decisão inválida: %v", err)
	}
	if p, err = s.DecidirAlienacao(ctx, arquivo, p.ID, al, " interesse ", "Acervo da fundação da cidade"); err != nil ||
		p.Alienacoes[0].Decisao != domain.AlienacaoInteresse || p.Alienacoes[0].DecididoEm == nil {
		t.Fatalf("manifestação: %+v %v", p.Alienacoes, err)
	}
	if _, err := s.DecidirAlienacao(ctx, arquivo, p.ID, al, "SEM_INTERESSE", "x"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("decidir duas vezes: %v", err)
	}
	// Notificação nova, sem resposta: depois de 60 dias, a preferência caduca.
	if p, err = s.NotificarAlienacao(ctx, arquivo, p.ID, domain.Alienacao{NotificadoEm: e.agora}); err != nil {
		t.Fatal(err)
	}
	e.agora = dia("2026-12-20")
	if _, err := s.DecidirAlienacao(ctx, arquivo, p.ID, p.Alienacoes[0].ID, "INTERESSE", "tarde"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("depois do prazo: %v", err)
	}
	if p, _ = s.Privado(ctx, arquivo, p.ID); !p.Alienacoes[0].Vencida {
		t.Fatalf("vencida: %+v", p.Alienacoes[0])
	}
	// Vencida a anterior, outra notificação; o município declina.
	if p, err = s.NotificarAlienacao(ctx, arquivo, p.ID, domain.Alienacao{NotificadoEm: e.agora}); err != nil {
		t.Fatal(err)
	}
	if p, err = s.DecidirAlienacao(ctx, arquivo, p.ID, p.Alienacoes[0].ID, "SEM_INTERESSE", "Sem dotação"); err != nil ||
		p.Alienacoes[0].Decisao != domain.AlienacaoSemInteresse {
		t.Fatalf("sem interesse: %+v %v", p.Alienacoes, err)
	}

	// Doação é definitiva.
	if p, err = s.MudarGuarda(ctx, arquivo, p.ID, domain.GuardaDoado, "Termo de doação 2/2026", e.agora); err != nil || p.Guarda != domain.GuardaDoado {
		t.Fatalf("doação: %+v %v", p, err)
	}
	if len(p.Eventos) < 8 {
		t.Fatalf("histórico: %d eventos", len(p.Eventos))
	}

	// Parecer contrário encerra; entidade de serviço público nasce declarada.
	outro, _ := s.CadastrarPrivado(ctx, arquivo, fundoFamilia())
	pa.Favoravel = false
	if outro, err = s.RegistrarParecer(ctx, arquivo, outro.ID, pa); err != nil || outro.Situacao != domain.PrivadoNaoDeclar {
		t.Fatalf("parecer contrário: %+v %v", outro, err)
	}
	conc := fundoFamilia()
	conc.Titular, conc.TipoTitular, conc.ServicoPublico = "Concessionária de transporte", "PJ", true
	if c, err := s.CadastrarPrivado(ctx, arquivo, conc); err != nil || c.Situacao != domain.PrivadoDeclarado {
		t.Fatalf("§ 3º: %+v %v", c, err)
	}

	// Leitura só do Arquivo Público.
	if _, err := s.Privado(ctx, agente(e.unidade), p.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("ficha fora do Arquivo: %v", err)
	}
	if _, err := s.Privado(ctx, arquivo, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	if _, err := s.RegistrarParecer(ctx, agente(e.unidade), p.ID, pa); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("parecer fora do Arquivo: %v", err)
	}
	if _, err := s.NotificarAlienacao(ctx, agente(e.unidade), p.ID, domain.Alienacao{}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("alienação fora do Arquivo: %v", err)
	}
	if _, _, err := s.Privados(ctx, agente(e.unidade), "", pagination.New(1, 10, 100)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("lista fora do Arquivo: %v", err)
	}
	if items, total, err := s.Privados(ctx, arquivo, "Concessionária de transporte", pagination.New(1, 10, 100)); err != nil || total < 1 || items[0].Titular != conc.Titular {
		t.Fatalf("busca: %+v %v", items, err)
	}
}
