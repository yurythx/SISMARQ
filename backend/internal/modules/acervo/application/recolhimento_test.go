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
	"github.com/yurythx/sismarq/internal/platform/notificacoes"
)

var preparo = domain.Preparo{Classificados: true, Avaliados: true, Organizados: true, Higienizados: true, Acondicionados: true}

// arquivoPublico: a equipe do órgão central (concessão global).
func arquivoPublico(user uuid.UUID) auth.Identity {
	perms := []string{"acervo:read", "acervo:manage", "acervo:central"}
	return auth.Identity{Subject: "arquivo", UserID: user, Permissions: perms, Scopes: []auth.Scope{{Perfil: "arquivo", Permissions: perms}}}
}

// gestorLotado cria um usuário com acervo:manage lotado na unidade (quem
// recebe a cópia do registro de entrada).
func (e *env) gestorLotado(t *testing.T, unidade uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	user := dbtest.User(t, e.pool)
	var perfil uuid.UUID
	slug := "acervo-" + uuid.NewString()[:8]
	if err := e.pool.QueryRow(ctx, `INSERT INTO perfis (slug, nome, permissoes) VALUES ($1, $1, ARRAY['acervo:read', 'acervo:manage'])
		RETURNING id`, slug).Scan(&perfil); err != nil {
		t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `INSERT INTO user_scopes (user_id, perfil_id, unidade_id) VALUES ($1, $2, $3)`, user, perfil, unidade); err != nil {
		t.Fatal(err)
	}
	return user
}

func (e *env) avisos(t *testing.T, user uuid.UUID) int {
	t.Helper()
	var n int
	if err := e.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM notificacoes WHERE user_id = $1 AND modulo = 'acervo'`, user).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// intermediaria cria a unidade e a leva ao intermediário.
func (e *env) intermediaria(t *testing.T, s *application.Service, codigo string) application.Detalhe {
	t.Helper()
	agora := e.agora
	e.agora = dia("2019-06-01")
	d := e.criar(codigo)
	e.agora = dia("2021-01-01")
	if _, err := s.Transferir(context.Background(), admin, d.ID, ""); err != nil {
		t.Fatal(err)
	}
	e.agora = agora
	return d
}

// O recolhimento regular, do pedido à entrada no Arquivo Público, com a
// cópia do registro para a origem (arts. 5º, VII, e 53).
func TestRecolhimentoRegular(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2024-06-01")
	s := e.real()
	a, b := e.intermediaria(t, s, codPerm), e.intermediaria(t, s, codPerm)
	c := e.intermediaria(t, s, codElim)
	e.agora = dia("2024-06-01")
	gestor := e.gestorLotado(t, e.unidade)
	equipe := dbtest.User(t, e.pool)
	arquivo := arquivoPublico(equipe)
	origem := agente(e.unidade)
	novo := func(ids ...uuid.UUID) application.NovoRecolhimento {
		return application.NovoRecolhimento{OrigemID: e.unidade, Motivo: domain.MotivoRegular, Preparo: preparo, Historico: "Leis desde 1990", Unidades: ids}
	}
	if d, _ := s.Get(ctx, admin, a.ID); d.Prazos.Situacao != domain.SituacaoRecolher {
		t.Fatalf("situação antes: %s", d.Prazos.Situacao)
	}

	// Recusas do pedido.
	if _, err := s.CriarRecolhimento(ctx, agente(e.outra), novo(a.ID)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("fora do escopo: %v", err)
	}
	sem := novo(a.ID)
	sem.Preparo.Acondicionados = false
	if _, err := s.CriarRecolhimento(ctx, origem, sem); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("preparo incompleto (art. 53): %v", err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, novo(c.ID)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("destino eliminação: %v", err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, novo(uuid.New())); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("unidade inexistente: %v", err)
	}
	alheia := e.nova(codPerm)
	alheia.UnidadeID = e.outra
	de, err := s.Criar(ctx, admin, alheia)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, novo(de.ID)); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("unidade de outra origem: %v", err)
	}

	// Envio: as unidades só mudam de guarda na entrada.
	r, err := s.CriarRecolhimento(ctx, origem, novo(a.ID, b.ID, a.ID))
	if err != nil || r.Situacao != domain.RecolhimentoEnviado || !strings.HasPrefix(r.Numero, "RC-") || len(r.Itens) != 2 || r.TotalItens != 2 {
		t.Fatalf("envio: %+v %v", r, err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, novo(a.ID)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já em outro recolhimento: %v", err)
	}
	if u, _ := s.Get(ctx, origem, a.ID); u.Prazos.Situacao != domain.SituacaoRecolher || u.Movimentacoes[0].Tipo != domain.MovRecolhimento || !u.Pode.Gerir {
		t.Fatalf("relacionada: %+v", u.Unidade)
	}

	// Entrada: só o Arquivo Público; local conferido.
	if _, err := s.ReceberRecolhimento(ctx, origem, r.ID, application.Entrada{}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("a origem não registra a entrada: %v", err)
	}
	nada := uuid.New()
	if _, err := s.ReceberRecolhimento(ctx, arquivo, r.ID, application.Entrada{LocalID: &nada}); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("local inexistente: %v", err)
	}
	local, err := s.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: domain.LocalDeposito, Nome: "Depósito do Arquivo Público " + uuid.NewString()[:8], Ativo: true})
	if err != nil {
		t.Fatal(err)
	}
	if r, err = s.ReceberRecolhimento(ctx, arquivo, r.ID, application.Entrada{LocalID: &local.ID}); err != nil ||
		r.Situacao != domain.RecolhimentoRecebido || !strings.HasPrefix(r.RegistroEntrada, "RE-") || r.RecebidoEm == nil || *r.LocalID != local.ID {
		t.Fatalf("entrada: %+v %v", r, err)
	}
	if _, err := s.ReceberRecolhimento(ctx, arquivo, r.ID, application.Entrada{}); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("entrada repetida: %v", err)
	}
	u, err := s.Get(ctx, origem, a.ID)
	if err != nil || u.Fase != domain.FasePermanente || u.Situacao != domain.SituacaoRecolhida || u.RecolhidaEm == nil ||
		*u.RecolhimentoID != r.ID || *u.LocalID != local.ID || u.Pode.Gerir || !strings.Contains(u.Movimentacoes[0].Descricao, r.RegistroEntrada) {
		t.Fatalf("recolhida: %+v %v", u, err)
	}
	// Cópia do registro à origem (art. 5º, VII): o gestor lotado recebe; quem
	// registrou, não.
	if e.avisos(t, gestor) != 1 || e.avisos(t, equipe) != 0 {
		t.Fatalf("avisos: gestor %d, equipe %d", e.avisos(t, gestor), e.avisos(t, equipe))
	}

	// Sob a guarda do Arquivo Público: a origem consulta, o Arquivo gere.
	if _, err := s.Localizar(ctx, origem, a.ID, nil, "P1"); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("a origem não gere a recolhida: %v", err)
	}
	if got, err := s.Localizar(ctx, arquivo, a.ID, &local.ID, "Estante 3"); err != nil || got.Posicao != "Estante 3" || !got.Pode.Gerir {
		t.Fatalf("o Arquivo gere: %+v %v", got, err)
	}
	mudar := e.nova(codPerm)
	mudar.UnidadeID = e.outra
	if _, err := s.Alterar(ctx, admin, a.ID, mudar); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("proveniência fixa: %v", err)
	}
	if _, err := s.Transferir(ctx, arquivo, a.ID, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("recolhida não se transfere: %v", err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET unidade_id = $2 WHERE id = $1`, a.ID, e.outra); err == nil {
		t.Fatal("o banco mudou a proveniência da recolhida")
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET recolhida_em = NULL, recolhimento_id = NULL WHERE id = $1`, a.ID); err == nil {
		t.Fatal("o banco devolveu a recolhida ao órgão")
	}
	if _, err := s.CriarRecolhimento(ctx, origem, novo(a.ID)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já recolhida: %v", err)
	}

	// Leitura no escopo; o Arquivo Público vê todos.
	if _, err := s.Recolhimento(ctx, agente(e.outra), r.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("fora do escopo: %v", err)
	}
	if _, err := s.Recolhimento(ctx, arquivo, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	if _, total, _ := s.Recolhimentos(ctx, agente(e.outra), pagination.New(1, 10, 100)); total != 0 {
		t.Fatal("outra origem não vê")
	}
	if _, total, err := s.Recolhimentos(ctx, arquivo, pagination.New(1, 10, 100)); err != nil || total < 1 {
		t.Fatalf("o Arquivo vê: %v %v", total, err)
	}
	if got, err := s.Recolhimento(ctx, arquivo, r.ID); err != nil || got.ID != r.ID {
		t.Fatalf("o Arquivo lê o recolhimento: %+v %v", got, err)
	}
	if _, total, _ := s.Recolhimentos(ctx, origem, pagination.New(1, 10, 100)); total != 1 {
		t.Fatalf("a origem vê os seus: %d", total)
	}
}

// Extinção, devolução, cancelamento e a amostra permanente (arts. 39, 42 e
// 50).
func TestRecolhimentoEspecialDevolucaoECancelamento(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2019-06-01")
	s := e.real()
	corrente := e.criar(codPerm)
	amostra := e.criar(codElim)
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'PERMANENTE', destinacao = 'GUARDA_PERMANENTE' WHERE id = $1`, amostra.ID); err != nil {
		t.Fatal(err)
	}
	gestor := e.gestorLotado(t, e.unidade)
	// Só o órgão central, sem gestão global: não age pela origem.
	arquivo := auth.Identity{Subject: "arquivo", Permissions: []string{"acervo:read", "acervo:central"},
		Scopes: []auth.Scope{{Perfil: "arquivo", Permissions: []string{"acervo:read", "acervo:central"}}}}
	origem := agente(e.unidade)
	regular := application.NovoRecolhimento{OrigemID: e.unidade, Motivo: domain.MotivoRegular, Preparo: preparo, Unidades: []uuid.UUID{corrente.ID}}
	if _, err := s.CriarRecolhimento(ctx, origem, regular); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("regular direto do corrente: %v", err)
	}
	extincao := regular
	extincao.Motivo, extincao.Fundamento = domain.MotivoExtincao, "Lei 11.000/2026 extingue a autarquia"
	extincao.Unidades = []uuid.UUID{corrente.ID, amostra.ID}
	r, err := s.CriarRecolhimento(ctx, origem, extincao)
	if err != nil || r.Motivo != domain.MotivoExtincao || r.TotalItens != 2 {
		t.Fatalf("extinção: %+v %v", r, err)
	}

	// Devolução com o que falta.
	if _, err := s.DevolverRecolhimento(ctx, arquivo, r.ID, " "); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("devolução sem motivo: %v", err)
	}
	if _, err := s.DevolverRecolhimento(ctx, origem, r.ID, "x"); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("a origem não devolve: %v", err)
	}
	if r, err = s.DevolverRecolhimento(ctx, arquivo, r.ID, "Faltam as etiquetas das caixas"); err != nil ||
		r.Situacao != domain.RecolhimentoDevolvido || r.DevolucaoMotivo == "" {
		t.Fatalf("devolução: %+v %v", r, err)
	}
	if e.avisos(t, gestor) != 1 {
		t.Fatal("a origem é avisada da devolução")
	}
	if _, err := s.CancelarRecolhimento(ctx, origem, r.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("cancelar devolvido: %v", err)
	}

	// Reenvio e cancelamento pela origem.
	r, err = s.CriarRecolhimento(ctx, origem, extincao)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.CancelarRecolhimento(ctx, agente(e.outra), r.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("cancelar fora do escopo: %v", err)
	}
	if _, err := s.CancelarRecolhimento(ctx, arquivo, r.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("o Arquivo não cancela pela origem: %v", err)
	}
	if r, err = s.CancelarRecolhimento(ctx, origem, r.ID); err != nil || r.Situacao != domain.RecolhimentoCancelado {
		t.Fatalf("cancelamento: %+v %v", r, err)
	}

	// Emprestada não vai; TTDD fora do ar falha.
	if _, err := s.Emprestar(ctx, admin, corrente.ID, dbtest.User(t, e.pool), "consulta", e.agora); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, extincao); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("emprestada: %v", err)
	}
	e.ttdd.err = errors.New("atlas fora")
	if _, err := s.CriarRecolhimento(ctx, origem, extincao); err == nil {
		t.Fatal("TTDD fora deveria falhar")
	}
	e.ttdd.err = nil

	// Em guia de transferência também não.
	e.agora = dia("2021-01-01")
	outra := e.criar(codPerm)
	if _, err := s.CriarGuia(ctx, origem, e.unidade, []uuid.UUID{outra.ID}, ""); err != nil {
		t.Fatal(err)
	}
	extincao.Unidades = []uuid.UUID{outra.ID}
	if _, err := s.CriarRecolhimento(ctx, origem, extincao); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("em guia: %v", err)
	}
	if _, err := s.CriarRecolhimento(ctx, origem, application.NovoRecolhimento{}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("pedido vazio: %v", err)
	}
}

// hubFalso guarda os tópicos publicados (entrega em tempo real).
type hubFalso struct{ topicos []string }

func (h *hubFalso) Publish(_ context.Context, topic, _ string, _ any) error {
	h.topicos = append(h.topicos, topic)
	return nil
}

// A cópia do registro chega na hora, pelo tópico pessoal de quem gere a
// origem (ADR 027).
func TestRecolhimentoAvisoEmTempoReal(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2019-06-01")
	hub := &hubFalso{}
	s := e.real().WithNotificacoes(notificacoes.NewService(e.pool, hub, logger))
	gestor := e.gestorLotado(t, e.unidade)
	d := e.criar(codPerm)
	remetente := agente(e.unidade)
	remetente.UserID = dbtest.User(t, e.pool)
	r, err := s.CriarRecolhimento(ctx, remetente, application.NovoRecolhimento{OrigemID: e.unidade, Motivo: domain.MotivoExtincao,
		Fundamento: "Lei", Preparo: preparo, Unidades: []uuid.UUID{d.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ReceberRecolhimento(ctx, admin, r.ID, application.Entrada{}); err != nil {
		t.Fatal(err)
	}
	// O gestor lotado e quem enviou recebem.
	if len(hub.topicos) != 2 || !strings.Contains(strings.Join(hub.topicos, " "), gestor.String()) ||
		!strings.Contains(strings.Join(hub.topicos, " "), remetente.UserID.String()) {
		t.Fatalf("entrega: %v", hub.topicos)
	}
}
