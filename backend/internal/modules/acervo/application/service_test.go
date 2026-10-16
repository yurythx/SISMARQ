package application_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/events"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
	"github.com/yurythx/sismarq/internal/platform/outbox"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func anos(n int) *int { return &n }

func dia(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)
	return t
}

// ttdd é a TTDD de teste (a real vem do Atlas, ligada em internal/app).
type ttdd struct {
	mu     sync.Mutex
	series map[string]domain.Serie
	err    error
}

func (f *ttdd) Serie(_ context.Context, codigo string) (domain.Serie, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return domain.Serie{}, f.err
	}
	s, ok := f.series[codigo]
	if !ok {
		return s, domain.ErrSerieInvalida
	}
	return s, nil
}

func (f *ttdd) set(codigo string, s domain.Serie) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s.Codigo = codigo
	f.series[codigo] = s
}

const (
	codElim  = "9.9.01.00.01"
	codPerm  = "9.9.01.00.02"
	codVigor = "9.9.01.00.03"
)

func novaTTDD() *ttdd {
	t := &ttdd{series: map[string]domain.Serie{}}
	t.set(codElim, domain.Serie{Descritor: "Notas de empenho", CorrenteAnos: anos(1), IntermAnos: anos(4), Destinacao: domain.DestinacaoEliminacao, Vigente: true})
	t.set(codPerm, domain.Serie{Descritor: "Leis e decretos", CorrenteAnos: anos(2), IntermAnos: anos(3), Destinacao: domain.DestinacaoGuardaPermanente, Vigente: true})
	t.set(codVigor, domain.Serie{Descritor: "Contratos", CorrenteCondicao: "Enquanto estiver vigorando", IntermAnos: anos(5), Destinacao: domain.DestinacaoEliminacao, Vigente: true})
	return t
}

type env struct {
	t       *testing.T
	pool    *pgxpool.Pool
	ttdd    *ttdd
	agora   time.Time
	unidade uuid.UUID
	outra   uuid.UUID
}

func novoEnv(t *testing.T) *env {
	pool := dbtest.Pool(t)
	return &env{t: t, pool: pool, ttdd: novaTTDD(), agora: dia("2019-06-01"),
		unidade: dbtest.Unidade(t, pool), outra: dbtest.Unidade(t, pool)}
}

func (e *env) svc(repo domain.Repository) *application.Service {
	return application.NewService(e.pool, repo, e.ttdd, outbox.NewWriter("test"), logger).WithRelogio(func() time.Time { return e.agora })
}

func (e *env) real() *application.Service { return e.svc(infrastructure.NewRepository()) }

// admin: concessão global de tudo.
var admin = auth.Identity{Subject: "admin", Permissions: []string{"*"}, Scopes: []auth.Scope{{Perfil: "admin", Permissions: []string{"*"}}}}

// agente: gestão, consulta e bloqueio só na unidade dada.
func agente(unidade uuid.UUID, perms ...string) auth.Identity {
	if len(perms) == 0 {
		perms = []string{"acervo:read", "acervo:manage", "acervo:bloqueio"}
	}
	return auth.Identity{Subject: "agente", Permissions: perms, Scopes: []auth.Scope{{Perfil: "agente", UnidadeID: &unidade, Permissions: perms}}}
}

func (e *env) nova(codigo string) domain.Unidade {
	return domain.Unidade{Tipo: domain.TipoCaixa, CodigoClassificacao: codigo, UnidadeID: e.unidade,
		DataInicial: dia("2018-01-01"), DataFinal: dia("2018-12-31"), Quantidade: 2, Especificacao: "caixas-arquivo"}
}

func (e *env) criar(codigo string) application.Detalhe {
	e.t.Helper()
	d, err := e.real().Criar(context.Background(), admin, e.nova(codigo))
	if err != nil {
		e.t.Fatal(err)
	}
	return d
}

func codigoDe(err error) apperrors.Code {
	if e, ok := apperrors.As(err); ok {
		return e.Code
	}
	return ""
}

func TestCadastroConsultaEAlteracao(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()

	d := e.criar(codElim)
	if !strings.HasPrefix(d.Numero, "UA-") || d.Fase != domain.FaseCorrente || d.Destinacao != domain.DestinacaoEliminacao ||
		d.Assunto != "Notas de empenho" || d.Situacao != domain.SituacaoNoCorrente || !d.Pode.Gerir || !d.Pode.Bloquear ||
		d.SerieResumo.Descritor != "Notas de empenho" || len(d.Movimentacoes) != 1 || d.Movimentacoes[0].Tipo != domain.MovCadastro {
		t.Fatalf("cadastro: %+v", d)
	}
	if !d.Prazos.FimCorrente.Equal(dia("2019-12-31")) || !d.Prazos.FimIntermediaria.Equal(dia("2023-12-31")) {
		t.Fatalf("prazos: %+v", d.Prazos)
	}

	// Validação, série inexistente e escopo.
	ruim := e.nova(codElim)
	ruim.Tipo = "GAVETA"
	if _, err := s.Criar(ctx, admin, ruim); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("validação: %v", err)
	}
	if _, err := s.Criar(ctx, admin, e.nova("0.0.00.00.00")); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("série inexistente: %v", err)
	}
	e.ttdd.set("9.9.01.00.09", domain.Serie{Destinacao: domain.DestinacaoEliminacao, Vigente: false})
	if _, err := s.Criar(ctx, admin, e.nova("9.9.01.00.09")); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("série revogada: %v", err)
	}
	if _, err := s.Criar(ctx, agente(e.outra), e.nova(codElim)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("fora do escopo: %v", err)
	}
	if _, err := s.Get(ctx, agente(e.outra), d.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("leitura fora do escopo é 404: %v", err)
	}
	if _, err := s.Get(ctx, admin, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	// Consulta sem gestão: lê, mas não altera.
	leitor := agente(e.unidade, "acervo:read")
	if got, err := s.Get(ctx, leitor, d.ID); err != nil || got.Pode.Gerir {
		t.Fatalf("leitor: %+v %v", got.Pode, err)
	}
	if _, err := s.Alterar(ctx, leitor, d.ID, e.nova(codElim)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("leitor não altera: %v", err)
	}

	// Alteração: reclassifica (destinação vem da nova série); mover para fora
	// do escopo, não.
	novo := e.nova(codPerm)
	novo.Assunto = "Leis municipais de 2018"
	novo.NivelAcesso, novo.HipoteseLegal = domain.AcessoRestrito, "Dados pessoais (LGPD)"
	alt, err := s.Alterar(ctx, agente(e.unidade), d.ID, novo)
	if err != nil || alt.Destinacao != domain.DestinacaoGuardaPermanente || alt.Assunto != "Leis municipais de 2018" || alt.Numero != d.Numero {
		t.Fatalf("alteração: %+v %v", alt, err)
	}
	mover := e.nova(codPerm)
	mover.UnidadeID = e.outra
	if _, err := s.Alterar(ctx, agente(e.unidade), d.ID, mover); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("mover para fora do escopo: %v", err)
	}
	if _, err := s.Alterar(ctx, admin, d.ID, e.nova("0.0.00.00.00")); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("reclassificar para série inexistente: %v", err)
	}
	e.ttdd.set(codPerm+"9", domain.Serie{Vigente: false})
	if _, err := s.Alterar(ctx, admin, d.ID, e.nova(codPerm+"9")); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("reclassificar para série revogada: %v", err)
	}
	ruim = e.nova(codPerm)
	ruim.Quantidade = -5
	if _, err := s.Alterar(ctx, admin, d.ID, ruim); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("alteração inválida: %v", err)
	}

	// Listagem, filtros e painel, no escopo.
	f := domain.Filtro{Query: d.Numero, Fase: domain.FaseCorrente, Codigo: "9.9.01", UnidadeID: &e.unidade, Situacao: domain.SituacaoNoCorrente}
	items, total, err := s.List(ctx, agente(e.unidade), f, pagination.New(1, 10, 100))
	if err != nil || total != 1 || len(items) != 1 || items[0].ID != d.ID {
		t.Fatalf("listagem: %d %+v %v", total, items, err)
	}
	if items, _, _ := s.List(ctx, agente(e.outra), domain.Filtro{Query: d.Numero}, pagination.New(1, 10, 100)); len(items) != 0 {
		t.Fatal("outra unidade não vê")
	}
	painel, err := s.Painel(ctx, agente(e.unidade), domain.Filtro{})
	if err != nil || painel[domain.SituacaoNoCorrente] != 1 {
		t.Fatalf("painel: %v %v", painel, err)
	}
}

// Transferência só com o prazo corrente cumprido (art. 29); permanente não
// é reclassificada; eliminada não muda.
func TestTransferenciaEFases(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	d := e.criar(codPerm)
	if _, err := s.Transferir(ctx, admin, d.ID, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("antes do prazo: %v", err)
	}
	e.agora = dia("2021-01-01")
	if _, err := s.Transferir(ctx, admin, d.ID, strings.Repeat("x", 1001)); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("observação longa: %v", err)
	}
	got, err := s.Transferir(ctx, agente(e.unidade), d.ID, " fim do exercício ")
	if err != nil || got.Fase != domain.FaseIntermediaria || got.Situacao != domain.SituacaoNoIntermediario ||
		got.Movimentacoes[0].Tipo != domain.MovTransferencia || !strings.Contains(got.Movimentacoes[0].Descricao, "pela guia GT-") ||
		got.Movimentacoes[1].Tipo != domain.MovGuia {
		t.Fatalf("transferência: %+v %v", got, err)
	}
	if _, err := s.Transferir(ctx, admin, d.ID, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("já no intermediário: %v", err)
	}
	// Sem evento: aguarda.
	v := e.criar(codVigor)
	if v.Situacao != domain.SituacaoAguardaEvento {
		t.Fatalf("vigorando sem evento: %s", v.Situacao)
	}
	if _, err := s.Transferir(ctx, admin, v.ID, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("aguarda evento: %v", err)
	}
	// Recolher (simulado no banco): permanente não se reclassifica, e as
	// travas do banco seguram o resto.
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'PERMANENTE' WHERE id = $1`, d.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Alterar(ctx, admin, d.ID, e.nova(codElim)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("permanente não se reclassifica: %v", err)
	}
	desc := e.nova(codPerm)
	desc.Observacoes = "Higienizado"
	if got, err := s.Alterar(ctx, admin, d.ID, desc); err != nil || got.Destinacao != domain.DestinacaoGuardaPermanente || got.Situacao != domain.SituacaoPermanente {
		t.Fatalf("permanente aceita descrição: %+v %v", got, err)
	}
	for _, sql := range []string{
		`UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`,
		`UPDATE acervo_unidades SET fase = 'INTERMEDIARIA' WHERE id = $1`,
		`DELETE FROM acervo_unidades WHERE id = $1`,
	} {
		if _, err := e.pool.Exec(ctx, sql, d.ID); err == nil {
			t.Errorf("o banco aceitou: %s", sql)
		}
	}
	// Eliminada (simulado): não muda mais.
	el := e.criar(codElim)
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`, el.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Alterar(ctx, admin, el.ID, e.nova(codElim)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("eliminada: %v", err)
	}
	if _, err := s.Localizar(ctx, admin, el.ID, nil, ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("eliminada não se localiza: %v", err)
	}
	if _, err := s.Bloquear(ctx, admin, el.ID, "1", ""); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("eliminada não se bloqueia: %v", err)
	}
	if _, err := s.Emprestar(ctx, admin, el.ID, dbtest.User(t, e.pool), "x", e.agora.AddDate(0, 0, 5)); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("eliminada não se empresta: %v", err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'CORRENTE' WHERE id = $1`, el.ID); err == nil {
		t.Fatal("eliminada não volta")
	}
	// Destinação de guarda permanente não vai a ELIMINADO, nem com bloqueio.
	g := e.criar(codPerm)
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`, g.ID); err == nil {
		t.Fatal("guarda permanente eliminada pelo banco")
	}
	b := e.criar(codElim)
	if _, err := s.Bloquear(ctx, admin, b.ID, "0001", "ação"); err != nil {
		t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`, b.ID); err == nil {
		t.Fatal("eliminada com bloqueio judicial")
	}
}

func TestBloqueioJudicial(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	d := e.criar(codElim)
	if _, err := s.Bloquear(ctx, agente(e.unidade, "acervo:read", "acervo:manage"), d.ID, "1", ""); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("sem acervo:bloqueio: %v", err)
	}
	if _, err := s.Bloquear(ctx, admin, d.ID, " ", ""); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("processo vazio: %v", err)
	}
	comUsuario := agente(e.unidade)
	comUsuario.UserID = dbtest.User(t, e.pool)
	got, err := s.Bloquear(ctx, comUsuario, d.ID, "0001234-56.2026.8.11.0003", "Ação de cobrança")
	if err != nil || len(got.Bloqueios) != 1 || !got.Bloqueios[0].Ativo() || got.Situacao != domain.SituacaoBloqueada ||
		got.Bloqueios[0].CriadoPor == nil || *got.Bloqueios[0].CriadoPor != comUsuario.UserID || got.Movimentacoes[0].Por == nil {
		t.Fatalf("bloqueio: %+v %v", got, err)
	}
	bid := got.Bloqueios[0].ID
	if _, err := s.Liberar(ctx, admin, bid, " "); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("liberar sem motivo: %v", err)
	}
	if _, err := s.Liberar(ctx, admin, uuid.New(), "x"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("bloqueio inexistente: %v", err)
	}
	if _, err := s.Liberar(ctx, agente(e.outra), bid, "x"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("liberar fora do escopo: %v", err)
	}
	got, err = s.Liberar(ctx, agente(e.unidade), bid, "Trânsito em julgado")
	if err != nil || got.Bloqueios[0].Ativo() || got.Situacao != domain.SituacaoNoCorrente || got.Movimentacoes[0].Tipo != domain.MovDesbloqueio {
		t.Fatalf("liberação: %+v %v", got, err)
	}
	if _, err := s.Liberar(ctx, admin, bid, "de novo"); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("liberar duas vezes: %v", err)
	}
}

func TestLocaisELocalizacao(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	sfx := uuid.NewString()[:8]
	if _, err := s.SalvarLocal(ctx, agente(e.unidade), nil, domain.Local{Tipo: domain.LocalPredio, Nome: "x"}); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público mantém os locais: %v", err)
	}
	if _, err := s.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: "COFRE", Nome: "x"}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("tipo inválido: %v", err)
	}
	predio, err := s.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: domain.LocalPredio, Nome: "Sede " + sfx, Ativo: true})
	if err != nil {
		t.Fatal(err)
	}
	estante, err := s.SalvarLocal(ctx, admin, nil, domain.Local{ParentID: &predio.ID, Tipo: domain.LocalEstante, Nome: "Estante 4", Ativo: true})
	if err != nil || estante.Caminho != "Sede "+sfx+" > Estante 4" {
		t.Fatalf("estante: %+v %v", estante, err)
	}
	// Ciclo, pai inexistente, local inexistente.
	if _, err := s.SalvarLocal(ctx, admin, &predio.ID, domain.Local{ParentID: &estante.ID, Tipo: domain.LocalPredio, Nome: "Sede " + sfx, Ativo: true}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("ciclo: %v", err)
	}
	nada := uuid.New()
	if _, err := s.SalvarLocal(ctx, admin, nil, domain.Local{ParentID: &nada, Tipo: domain.LocalSala, Nome: "x", Ativo: true}); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("pai inexistente: %v", err)
	}
	if _, err := s.SalvarLocal(ctx, admin, &nada, domain.Local{Tipo: domain.LocalSala, Nome: "x", Ativo: true}); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("local inexistente: %v", err)
	}
	locais, err := s.ListLocais(ctx)
	if err != nil || len(locais) < 2 {
		t.Fatalf("locais: %d %v", len(locais), err)
	}

	d := e.criar(codElim)
	got, err := s.Localizar(ctx, agente(e.unidade), d.ID, &estante.ID, " Prateleira 2, posição 7 ")
	if err != nil || got.LocalCaminho != estante.Caminho || got.Posicao != "Prateleira 2, posição 7" || got.Movimentacoes[0].Tipo != domain.MovLocalizacao {
		t.Fatalf("localização: %+v %v", got, err)
	}
	if items, _, _ := s.List(ctx, admin, domain.Filtro{Query: d.Numero}, pagination.New(1, 5, 100)); len(items) != 1 || items[0].LocalCaminho != estante.Caminho {
		t.Fatalf("caminho na listagem: %+v", items)
	}
	if _, err := s.Localizar(ctx, admin, d.ID, &estante.ID, strings.Repeat("x", 121)); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("posição longa: %v", err)
	}
	if _, err := s.Localizar(ctx, admin, d.ID, &nada, ""); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("local inexistente: %v", err)
	}
	// Local desativado não recebe unidade.
	estante.Ativo = false
	if _, err := s.SalvarLocal(ctx, admin, &estante.ID, estante); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Localizar(ctx, admin, d.ID, &estante.ID, ""); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("local desativado: %v", err)
	}
	nova := e.nova(codElim)
	nova.LocalID = &estante.ID
	if _, err := s.Criar(ctx, admin, nova); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("cadastro em local desativado: %v", err)
	}
}

// Recálculo: o calendário anda, a TTDD muda e a série pode ser revogada.
func TestRecalculo(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	d := e.criar(codElim)
	// Série some da TTDD: fica para reclassificar, e a ficha ainda abre.
	e.ttdd.set(codElim, domain.Serie{Descritor: "Notas de empenho", CorrenteAnos: anos(1), IntermAnos: anos(4), Destinacao: domain.DestinacaoEliminacao, Vigente: false})
	if n, err := s.Recalcular(ctx); err != nil || n < 1 {
		t.Fatalf("recálculo: %d %v", n, err)
	}
	got, err := s.Get(ctx, admin, d.ID)
	if err != nil || got.Situacao != application.SituacaoSerieRevogada || got.SerieResumo.Vigente {
		t.Fatalf("revogada: %+v %v", got, err)
	}
	delete(e.ttdd.series, codElim)
	if err := s.HandleTTDD(ctx, events.Event{Type: "atlas.ttdd.carga.aplicada"}); err != nil {
		t.Fatal(err)
	}
	if got, err := s.Get(ctx, admin, d.ID); err != nil || got.SerieResumo.Codigo != codElim || got.Situacao != application.SituacaoSerieRevogada {
		t.Fatalf("série inexistente: %+v %v", got, err)
	}
	// Erro da TTDD: o recálculo para e informa; a ficha também.
	e.ttdd.err = errors.New("atlas fora")
	if _, err := s.Recalcular(ctx); err == nil {
		t.Fatal("recálculo deveria falhar")
	}
	if _, err := s.Get(ctx, admin, d.ID); err == nil {
		t.Fatal("ficha deveria falhar")
	}
	// O worker roda, registra a falha e para com o contexto.
	cctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	if err := s.RunRecalculo(10 * time.Millisecond)(cctx); err != nil {
		t.Fatal(err)
	}
	// Sem erro: para logo depois de registrar o recálculo (sem prazo de relógio).
	e.ttdd.err = nil
	cctx2, cancel2 := context.WithCancel(ctx)
	defer cancel2()
	ok := application.NewService(e.pool, infrastructure.NewRepository(), e.ttdd, outbox.NewWriter("test"), slog.New(cancelaAoLogar{cancel2}))
	if err := ok.RunRecalculo(time.Hour)(cctx2); err != nil {
		t.Fatal(err)
	}
}

func TestMapError(t *testing.T) {
	casos := map[error]apperrors.Code{
		&domain.ErrValidacao{Campo: "x", Motivo: "y"}: apperrors.CodeValidation,
		domain.ErrNaoEncontrada:                       apperrors.CodeNotFound,
		domain.ErrForaDoEscopo:                        apperrors.CodeForbidden,
		domain.ErrLocalNaoExiste:                      apperrors.CodeNotFound,
		domain.ErrBloqueioNaoExiste:                   apperrors.CodeNotFound,
		domain.ErrSerieInvalida:                       apperrors.CodeValidation,
		domain.ErrEmprestimoNaoExiste:                 apperrors.CodeNotFound,
		domain.ErrGuiaNaoExiste:                       apperrors.CodeNotFound,
		domain.ErrUnidadeProdutora:                    apperrors.CodeValidation,
	}
	for _, trava := range []error{domain.ErrBloqueioLiberado, domain.ErrPermanente, domain.ErrBloqueioJudicial, domain.ErrNaoPrevistaNaTTDD,
		domain.ErrPrazoNaoCumprido, domain.ErrAguardaEvento, domain.ErrDestinacaoIndefinida, domain.ErrFaseInvalida, domain.ErrEliminado,
		domain.ErrEmprestimoAberto, domain.ErrEmprestimoDevolvido, domain.ErrEmprestada, domain.ErrGuiaSituacao, domain.ErrEmOutraGuia} {
		casos[trava] = apperrors.CodeConflict
	}
	for err, code := range casos {
		if got := codigoDe(application.MapError(err)); got != code {
			t.Errorf("%v: %s, esperado %s", err, got, code)
		}
	}
	if application.MapError(nil) != nil || !errors.Is(application.MapError(errBoom), errBoom) {
		t.Fatal("nil e erro desconhecido passam")
	}
}

// cancelaAoLogar encerra o worker no primeiro registro de log.
type cancelaAoLogar struct{ cancel context.CancelFunc }

func (h cancelaAoLogar) Enabled(context.Context, slog.Level) bool { return true }
func (h cancelaAoLogar) Handle(context.Context, slog.Record) error {
	h.cancel()
	return nil
}
func (h cancelaAoLogar) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h cancelaAoLogar) WithGroup(string) slog.Handler      { return h }
