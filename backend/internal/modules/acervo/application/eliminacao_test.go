package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

// cspad cria a Comissão Setorial do órgão da unidade (art. 31, § 6º) e
// devolve os membros (o primeiro é o presidente).
func (e *env) cspad(t *testing.T, s *application.Service) (domain.Comissao, []uuid.UUID) {
	t.Helper()
	var ent uuid.UUID
	if err := e.pool.QueryRow(context.Background(), `SELECT entidade_id FROM unidades WHERE id = $1`, e.unidade).Scan(&ent); err != nil {
		t.Fatal(err)
	}
	areas := []string{domain.AreaArquivologia, domain.AreaJuridico, domain.AreaAdmFinancas, domain.AreaUnidadeProdutora, domain.AreaOutra}
	var ids []uuid.UUID
	c := domain.Comissao{Tipo: domain.ComissaoSetorial, EntidadeID: &ent, Nome: "CSPAD de teste", Portaria: "Portaria 1/2026", Ativa: true}
	for i, a := range areas {
		id := dbtest.User(t, e.pool)
		ids = append(ids, id)
		papel := domain.PapelMembro
		if i == 0 {
			papel = domain.PapelPresidente
		}
		c.Membros = append(c.Membros, domain.Membro{UserID: id, Papel: papel, Area: a})
	}
	salva, err := s.SalvarComissao(context.Background(), admin, nil, c)
	if err != nil {
		t.Fatal(err)
	}
	return salva, ids
}

func avaliador(unidade uuid.UUID) auth.Identity {
	return agente(unidade, "acervo:read", "acervo:avaliacao")
}

// O rito inteiro (arts. 31 a 40; Anexos I a IV) e as recusas de cada etapa.
func TestRitoDaEliminacao(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-09-01")
	s := e.real()
	a, b, c, d := e.criar(codElim), e.criar(codElim), e.criar(codElim), e.criar(codElim)
	perm := e.criar(codPerm)
	cspad := avaliador(e.unidade)

	// Sem comissão, não há relação (art. 31).
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Unidades: []uuid.UUID{a.ID}}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("sem comissão: %v", err)
	}
	if _, err := s.SalvarComissao(ctx, cspad, nil, domain.Comissao{}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público cadastra comissões: %v", err)
	}
	com, membros := e.cspad(t, s)
	if lista, err := s.Comissoes(ctx); err != nil || len(lista) == 0 {
		t.Fatalf("comissões: %v", err)
	}

	// Relação: escopo, unidade não elegível (guarda permanente).
	if _, err := s.CriarListagem(ctx, avaliador(e.outra), application.NovaListagem{UnidadeID: e.unidade}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("fora do escopo: %v", err)
	}
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Unidades: []uuid.UUID{perm.ID}}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("guarda permanente não entra (art. 19): %v", err)
	}
	l, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, ContasExercicios: "2018", Unidades: []uuid.UUID{a.ID, b.ID, c.ID}})
	if err != nil || l.Situacao != domain.ListagemRascunho || len(l.Itens) != 3 || l.ComissaoID != com.ID || l.Numero == "" ||
		l.Rito != domain.RitoOrdinario || l.AtaSignum != domain.SignumAusente {
		t.Fatalf("relação: %+v %v", l, err)
	}
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Unidades: []uuid.UUID{a.ID}}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já em outra relação: %v", err)
	}

	// Ata: sem parecer jurídico, sem quórum; depois aprovada.
	ata := application.Ata{Numero: "1/2026", Em: e.agora, Local: "Sala da SEMAD", Presentes: membros[:3]}
	if _, err := s.Aprovar(ctx, cspad, l.ID, ata); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("sem parecer (art. 32, § 4º): %v", err)
	}
	l, err = s.AlterarListagem(ctx, cspad, l.ID, application.NovaListagem{ParecerJuridico: "Parecer 7/2026 da assessoria jurídica",
		ContasExercicios: "2018", Unidades: []uuid.UUID{d.ID}}, []uuid.UUID{c.ID})
	if err != nil || len(l.Itens) != 3 {
		t.Fatalf("alteração: %+v %v", l.Itens, err)
	}
	if _, err := s.Aprovar(ctx, cspad, l.ID, application.Ata{Numero: "1/2026", Em: e.agora, Presentes: membros[1:3]}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("sem quórum: %v", err)
	}
	if l, err = s.Aprovar(ctx, cspad, l.ID, ata); err != nil || l.Situacao != domain.ListagemAprovada {
		t.Fatalf("aprovação: %+v %v", l, err)
	}
	if _, err := s.AlterarListagem(ctx, cspad, l.ID, application.NovaListagem{}, nil); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("aprovada não se altera: %v", err)
	}

	// O Arquivo Público devolve; nova ata; autoriza.
	if _, err := s.Autorizar(ctx, cspad, l.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público autoriza: %v", err)
	}
	if _, err := s.DevolverListagem(ctx, admin, l.ID, " "); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("devolução sem motivo: %v", err)
	}
	if l, err = s.DevolverListagem(ctx, admin, l.ID, "Conferir as datas-limite"); err != nil || l.Situacao != domain.ListagemRascunho || l.AtaEm != nil {
		t.Fatalf("devolução: %+v %v", l, err)
	}
	if _, err := s.DevolverListagem(ctx, admin, l.ID, "de novo"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("devolver rascunho: %v", err)
	}
	if _, err := s.Autorizar(ctx, admin, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("autorizar rascunho: %v", err)
	}
	if _, err = s.Aprovar(ctx, cspad, l.ID, ata); err != nil {
		t.Fatal(err)
	}
	if l, err = s.Autorizar(ctx, admin, l.ID); err != nil || l.Situacao != domain.ListagemAutorizada || l.AutorizadaEm == nil {
		t.Fatalf("autorização: %+v %v", l, err)
	}

	// Edital, pedido de desentranhamento, amostragem.
	if _, err := s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{Requerente: "João", Pedido: "Cópia", RecebidoEm: e.agora}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("pedido antes do edital: %v", err)
	}
	if l, err = s.PublicarEdital(ctx, cspad, l.ID, application.Edital{Numero: "3/2026", PublicadoEm: e.agora, Diario: "DIORONDON-e nº 6.300, p. 12"}); err != nil ||
		l.Situacao != domain.ListagemEmEdital {
		t.Fatalf("edital: %+v %v", l, err)
	}
	if _, err := s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{Requerente: " ", Pedido: "x", RecebidoEm: e.agora}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("pedido sem requerente: %v", err)
	}
	if _, err := s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{Requerente: "x", Pedido: " ", RecebidoEm: e.agora}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("pedido vazio: %v", err)
	}
	fora := perm.ID
	if _, err := s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{UnidadeID: &fora, Requerente: "x", Pedido: "y", RecebidoEm: e.agora}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("pedido de unidade fora da relação: %v", err)
	}
	if l, err = s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{UnidadeID: &d.ID, Requerente: "Maria (CPF ***.456.789-**)",
		Pedido: "Desentranhamento do contrato original", RecebidoEm: e.agora.AddDate(0, 0, 5)}); err != nil || len(l.Pedidos) != 1 {
		t.Fatalf("pedido: %+v %v", l.Pedidos, err)
	}
	pedido := l.Pedidos[0].ID
	e.agora = e.agora.AddDate(0, 0, 31)
	if _, err := s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: e.agora, Acompanhamento: "Servidor X"}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("termo com pedido pendente: %v", err)
	}
	if _, err := s.DecidirPedido(ctx, cspad, l.ID, pedido, true, true, " "); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("decisão sem fundamento: %v", err)
	}
	if _, err := s.DecidirPedido(ctx, cspad, l.ID, uuid.New(), true, false, "ok"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("pedido inexistente: %v", err)
	}
	if l, err = s.DecidirPedido(ctx, cspad, l.ID, pedido, true, true, "Deferido: a unidade sai da relação"); err != nil || len(l.AEliminar()) != 2 {
		t.Fatalf("decisão: %+v %v", l, err)
	}
	if _, err := s.DecidirPedido(ctx, cspad, l.ID, pedido, false, false, "de novo"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("decidir duas vezes: %v", err)
	}
	if _, err := s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: e.agora, Acompanhamento: "Servidor X"}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("termo sem amostragem (art. 39): %v", err)
	}
	if l, err = s.Amostrar(ctx, cspad, l.ID, "Uma caixa por exercício, sorteada", 10); err != nil || l.AmostraEm == nil {
		t.Fatalf("amostragem: %+v %v", l, err)
	}
	var amostra uuid.UUID
	for _, i := range l.Itens {
		if i.Amostra {
			amostra = i.ID
		}
	}
	if amostra == uuid.Nil || len(l.AEliminar()) != 1 {
		t.Fatalf("uma unidade na amostra: %+v", l.Itens)
	}
	if u, _ := s.Get(ctx, admin, amostra); u.Fase != domain.FasePermanente || u.Destinacao != domain.DestinacaoGuardaPermanente {
		t.Fatalf("amostra vai à guarda permanente: %+v", u.Unidade)
	}
	if _, err := s.CancelarListagem(ctx, cspad, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("cancelar depois da amostragem: %v", err)
	}

	// Termo de Eliminação (Anexo II).
	if _, err := s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: e.agora, Acompanhamento: " "}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("termo sem acompanhamento: %v", err)
	}
	l, err = s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: e.agora, Acompanhamento: "Servidor X",
		DestinacaoPapel: "Doado à cooperativa de catadores (Decreto 5.940/2006)"})
	if err != nil || l.Situacao != domain.ListagemEliminada || l.TermoEliminadoEm == nil || l.TermoCopiaEm == nil ||
		l.TermoSignum != domain.SignumAusente {
		t.Fatalf("termo: %+v %v", l, err)
	}
	for _, id := range []uuid.UUID{a.ID, b.ID} {
		if id == amostra {
			continue
		}
		if u, _ := s.Get(ctx, admin, id); u.Fase != domain.FaseEliminado || u.Movimentacoes[0].Tipo != domain.MovEliminacao {
			t.Fatalf("eliminada: %+v", u.Unidade)
		}
	}
	if u, _ := s.Get(ctx, admin, d.ID); u.Fase != domain.FaseCorrente {
		t.Fatalf("retirada por pedido não é eliminada: %s", u.Fase)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`, amostra); err == nil {
		t.Fatal("o banco eliminou a amostra")
	}

	// Leitura e listagem no escopo.
	if _, err := s.Listagem(ctx, avaliador(e.outra), l.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("relação fora do escopo: %v", err)
	}
	if items, total, err := s.Listagens(ctx, cspad, pagination.New(1, 10, 100)); err != nil || total < 1 || items[0].ID != l.ID {
		t.Fatalf("relações: %+v %v", items, err)
	}
	if _, total, _ := s.Listagens(ctx, avaliador(e.outra), pagination.New(1, 10, 100)); total != 0 {
		t.Fatal("outra unidade não vê")
	}
}

// Bloqueio judicial depois da ata trava a autorização; cancelamento libera
// as unidades.
func TestRelacaoBloqueioECancelamento(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-09-01")
	s := e.real()
	_, membros := e.cspad(t, s)
	cspad := avaliador(e.unidade)
	u := e.criar(codElim)
	l, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, ParecerJuridico: "Parecer", Unidades: []uuid.UUID{u.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Aprovar(ctx, cspad, l.ID, application.Ata{Numero: "2/2026", Em: e.agora, Presentes: membros[:3]}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Bloquear(ctx, admin, u.ID, "0009999-00.2026.8.11.0003", "Ação popular"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Autorizar(ctx, admin, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("bloqueio judicial depois da ata (art. 44): %v", err)
	}
	if l, err = s.CancelarListagem(ctx, cspad, l.ID); err != nil || l.Situacao != domain.ListagemCancelada {
		t.Fatalf("cancelamento: %+v %v", l, err)
	}
	if _, err := s.CancelarListagem(ctx, cspad, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("cancelar duas vezes: %v", err)
	}
	if _, err := s.Listagem(ctx, admin, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("relação inexistente: %v", err)
	}
	// Comissão: alteração, inexistente e segunda ativa no mesmo órgão.
	coms, _ := s.Comissoes(ctx)
	var minha domain.Comissao
	for _, c := range coms {
		if c.Nome == "CSPAD de teste" && c.Ativa && len(c.Membros) == 5 && c.EntidadeID != nil {
			minha = c
		}
	}
	minha.Portaria = "Portaria 2/2026"
	if got, err := s.SalvarComissao(ctx, admin, &minha.ID, minha); err != nil || got.Portaria != "Portaria 2/2026" {
		t.Fatalf("alteração da comissão: %+v %v", got, err)
	}
	nada := uuid.New()
	if _, err := s.SalvarComissao(ctx, admin, &nada, minha); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("comissão inexistente: %v", err)
	}
	if _, err := s.SalvarComissao(ctx, admin, nil, minha); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("segunda comissão ativa no órgão: %v", err)
	}
	if _, err := s.SalvarComissao(ctx, admin, nil, domain.Comissao{Tipo: "X"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("comissão inválida: %v", err)
	}
}

// Recusas que o rito não percorre: campos longos, itens inelegíveis, leitura
// fora do escopo e etapas fora de ordem.
func TestRelacaoRecusas(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-09-01")
	s := e.real()
	_, membros := e.cspad(t, s)
	cspad := avaliador(e.unidade)
	longo := func(n int) string { return strings.Repeat("x", n+1) }
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, ParecerJuridico: longo(4000)}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("parecer longo: %v", err)
	}
	alheia := e.nova(codElim)
	alheia.UnidadeID = e.outra
	de, err := s.Criar(ctx, admin, alheia)
	if err != nil {
		t.Fatal(err)
	}
	emprestada := e.criar(codElim)
	if _, err := s.Emprestar(ctx, admin, emprestada.ID, dbtest.User(t, e.pool), "consulta", e.agora); err != nil {
		t.Fatal(err)
	}
	nova := func(id uuid.UUID) error {
		_, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Unidades: []uuid.UUID{id}})
		return err
	}
	if err := nova(uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("unidade inexistente: %v", err)
	}
	if err := nova(de.ID); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("unidade de outra unidade: %v", err)
	}
	if err := nova(emprestada.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("unidade emprestada: %v", err)
	}
	u := e.criar(codElim)
	e.ttdd.err = errors.New("atlas fora")
	if err := nova(u.ID); err == nil {
		t.Fatal("TTDD fora deveria falhar")
	}
	e.ttdd.err = nil

	l, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, ParecerJuridico: "Parecer", Unidades: []uuid.UUID{u.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AlterarListagem(ctx, cspad, l.ID, application.NovaListagem{ContasExercicios: longo(200)}, nil); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("exercícios longos: %v", err)
	}
	ata := application.Ata{Numero: "1/2026", Em: e.agora, Presentes: membros[:3]}
	if _, err := s.Aprovar(ctx, avaliador(e.outra), l.ID, ata); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("aprovar fora do escopo: %v", err)
	}
	edital := application.Edital{Numero: "1/2026", PublicadoEm: e.agora}
	if _, err := s.PublicarEdital(ctx, cspad, l.ID, edital); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("edital de rascunho: %v", err)
	}
	if _, err := s.Aprovar(ctx, cspad, l.ID, ata); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Autorizar(ctx, admin, l.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.DecidirPedido(ctx, cspad, l.ID, uuid.New(), false, false, "x"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("decidir antes do edital: %v", err)
	}
	if _, err := s.Amostrar(ctx, cspad, l.ID, "c", 10); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("amostrar antes do edital: %v", err)
	}
	edital.Diario = longo(200)
	if _, err := s.PublicarEdital(ctx, cspad, l.ID, edital); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("diário longo: %v", err)
	}
	if _, err := s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: e.agora, Acompanhamento: "X", DestinacaoPapel: longo(500)}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("destinação longa: %v", err)
	}
}

func (e *env) ccpad(t *testing.T, s *application.Service) (domain.Comissao, []uuid.UUID) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), `UPDATE acervo_comissoes SET ativa = false WHERE tipo = 'CENTRAL' AND ativa`); err != nil {
		t.Fatal(err)
	}
	areas := []string{domain.AreaArquivologia, domain.AreaJuridico, domain.AreaAdmFinancas, domain.AreaArquivoPublico, domain.AreaOutra}
	var ids []uuid.UUID
	c := domain.Comissao{Tipo: domain.ComissaoCentral, Nome: "CCPAD de teste", Portaria: "Portaria 2/2026", Ativa: true}
	for i, a := range areas {
		id := dbtest.User(t, e.pool)
		ids = append(ids, id)
		papel := domain.PapelMembro
		if i == 0 {
			papel = domain.PapelPresidente
		}
		c.Membros = append(c.Membros, domain.Membro{UserID: id, Papel: papel, Area: a})
	}
	salva, err := s.SalvarComissao(context.Background(), admin, nil, c)
	if err != nil {
		t.Fatal(err)
	}
	return salva, ids
}

func membroCCPAD(user, unidade uuid.UUID) auth.Identity {
	id := agente(unidade, "acervo:read", "acervo:avaliacao")
	id.UserID = user
	return id
}

func semad() auth.Identity {
	return auth.Identity{Subject: "semad", Permissions: []string{"acervo:read", "acervo:semad"},
		Scopes: []auth.Scope{{Perfil: "semad", Permissions: []string{"acervo:read", "acervo:semad"}}}}
}

func TestRitosEspeciais(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-09-01")
	s := e.real()
	_, membros := e.cspad(t, s)
	_, ccpadIDs := e.ccpad(t, s)
	cspad, cc, sem := avaliador(e.unidade), membroCCPAD(ccpadIDs[0], e.unidade), semad()
	indef, elim := e.criar(codIndef), e.criar(codElim)

	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Rito: "X"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("rito inválido: %v", err)
	}
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Rito: domain.RitoArt35, ParecerJuridico: "P",
		Unidades: []uuid.UUID{elim.ID}}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("art. 35 não aceita o que cabe no rito ordinário: %v", err)
	}
	if _, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Rito: domain.RitoArt35, ParecerJuridico: "P",
		Unidades: []uuid.UUID{e.criar(codPerm).ID}}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("guarda permanente não entra no art. 35: %v", err)
	}

	l, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Rito: domain.RitoArt35, ParecerJuridico: "P",
		Unidades: []uuid.UUID{indef.ID}})
	if err != nil || l.Rito != domain.RitoArt35 {
		t.Fatalf("art. 35: %+v %v", l, err)
	}
	ata := application.Ata{Numero: "1", Em: e.agora, Presentes: membros[:3]}
	if l, err = s.Aprovar(ctx, cspad, l.ID, ata); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Autorizar(ctx, admin, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("art. 35 sem a CCPAD: %v", err)
	}
	if _, err := s.AutorizarCCPAD(ctx, cspad, l.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só membro da CCPAD: %v", err)
	}
	if l, err = s.AutorizarCCPAD(ctx, cc, l.ID); err != nil || l.CcpadAutorizadaEm == nil {
		t.Fatalf("CCPAD: %+v %v", l, err)
	}
	if _, err := s.AutorizarCCPAD(ctx, cc, l.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("CCPAD duas vezes: %v", err)
	}
	if l, err = s.Autorizar(ctx, admin, l.ID); err != nil || l.Situacao != domain.ListagemAguardandoSEMAD {
		t.Fatalf("aguarda SEMAD: %+v %v", l, err)
	}
	if _, err := s.AutorizarSEMAD(ctx, cspad, l.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só a SEMAD: %v", err)
	}
	if _, err := s.PublicarEdital(ctx, cspad, l.ID, application.Edital{Numero: "1", PublicadoEm: e.agora}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("edital sem a SEMAD: %v", err)
	}
	if l, err = s.AutorizarSEMAD(ctx, sem, l.ID); err != nil || l.Situacao != domain.ListagemAutorizada || l.SemadAutorizadaEm == nil {
		t.Fatalf("SEMAD: %+v %v", l, err)
	}

	l36, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Rito: domain.RitoArt36, ParecerJuridico: "P",
		Unidades: []uuid.UUID{e.criar(codElim).ID}})
	if err != nil || l36.Rito != domain.RitoArt36 {
		t.Fatalf("art. 36: %+v %v", l36, err)
	}
	if _, err := s.AutorizarCCPAD(ctx, cc, l36.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("CCPAD só no art. 35: %v", err)
	}
	ata36 := application.Ata{Numero: "2", Em: e.agora, Presentes: ccpadIDs[:3]}
	if l36, err = s.Aprovar(ctx, cspad, l36.ID, ata36); err != nil {
		t.Fatal(err)
	}
	if l36, err = s.Autorizar(ctx, admin, l36.ID); err != nil || l36.Situacao != domain.ListagemAguardandoSEMAD {
		t.Fatalf("art. 36 aguarda SEMAD: %+v %v", l36, err)
	}
	if l36, err = s.DevolverListagem(ctx, admin, l36.ID, "Revisar a classificação"); err != nil || l36.Situacao != domain.ListagemRascunho {
		t.Fatalf("devolver da SEMAD: %+v %v", l36, err)
	}
}
