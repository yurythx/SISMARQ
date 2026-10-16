package application_test

// Código gerado por scripts/genfault.py a partir da interface Repository do
// domínio. Não edite à mão: regenere se a interface mudar.

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

type innerRepo = domain.Repository

var errBoom = errors.New("falha simulada no repositório")

// faultRepo envolve o repositório real: na chamada failAt devolve erro; depois
// da chamada poisonAt "envenena" a transação (a instrução SQL seguinte — do
// repositório, do outbox ou da auditoria — falha). noTx marca veneno em
// leitura fora de transação, onde não há instrução seguinte na mesma conexão.
type faultRepo struct {
	inner                   innerRepo
	calls, failAt, poisonAt int
	trace                   []string
	noTx                    bool
}

func (fr *faultRepo) hook(name string) error {
	fr.calls++
	fr.trace = append(fr.trace, name)
	if fr.calls == fr.failAt {
		return errBoom
	}
	return nil
}

func (fr *faultRepo) post(ctx context.Context, db database.DBTX) {
	if fr.calls != fr.poisonAt {
		return
	}
	if _, inTx := db.(pgx.Tx); !inTx {
		fr.noTx = true
		return
	}
	_, _ = db.Exec(ctx, `SELECT 1/0`)
}

func (fr *faultRepo) InsertUnidade(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) (string, error) {
	if err := fr.hook("InsertUnidade"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertUnidade(ctx, db, u, c)
}

func (fr *faultRepo) GetUnidade(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Unidade, domain.Calculo, error) {
	if err := fr.hook("GetUnidade"); err != nil {
		var z0 domain.Unidade
		var z1 domain.Calculo
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetUnidade(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) UpdateUnidade(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) error {
	if err := fr.hook("UpdateUnidade"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.UpdateUnidade(ctx, db, u, c)
}

func (fr *faultRepo) SetCalculo(ctx context.Context, db database.DBTX, id uuid.UUID, c domain.Calculo) error {
	if err := fr.hook("SetCalculo"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SetCalculo(ctx, db, id, c)
}

func (fr *faultRepo) SetFase(ctx context.Context, db database.DBTX, id uuid.UUID, fase domain.Fase, c domain.Calculo) error {
	if err := fr.hook("SetFase"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SetFase(ctx, db, id, fase, c)
}

func (fr *faultRepo) ListUnidades(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.UnidadeLista, int64, error) {
	if err := fr.hook("ListUnidades"); err != nil {
		var z0 []domain.UnidadeLista
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListUnidades(ctx, db, f, p)
}

func (fr *faultRepo) Painel(ctx context.Context, db database.DBTX, f domain.Filtro) (map[string]int64, error) {
	if err := fr.hook("Painel"); err != nil {
		var z0 map[string]int64
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.Painel(ctx, db, f)
}

func (fr *faultRepo) ParaRecalculo(ctx context.Context, db database.DBTX, depoisDe uuid.UUID, limite int) ([]domain.Unidade, error) {
	if err := fr.hook("ParaRecalculo"); err != nil {
		var z0 []domain.Unidade
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ParaRecalculo(ctx, db, depoisDe, limite)
}

func (fr *faultRepo) InsertBloqueio(ctx context.Context, db database.DBTX, b domain.Bloqueio) error {
	if err := fr.hook("InsertBloqueio"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertBloqueio(ctx, db, b)
}

func (fr *faultRepo) GetBloqueio(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Bloqueio, error) {
	if err := fr.hook("GetBloqueio"); err != nil {
		var z0 domain.Bloqueio
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetBloqueio(ctx, db, id)
}

func (fr *faultRepo) LiberarBloqueio(ctx context.Context, db database.DBTX, b domain.Bloqueio) error {
	if err := fr.hook("LiberarBloqueio"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.LiberarBloqueio(ctx, db, b)
}

func (fr *faultRepo) ListBloqueios(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Bloqueio, error) {
	if err := fr.hook("ListBloqueios"); err != nil {
		var z0 []domain.Bloqueio
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListBloqueios(ctx, db, unidade)
}

func (fr *faultRepo) TemBloqueioAtivo(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	if err := fr.hook("TemBloqueioAtivo"); err != nil {
		var z0 bool
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.TemBloqueioAtivo(ctx, db, unidade)
}

func (fr *faultRepo) InsertMovimentacao(ctx context.Context, db database.DBTX, unidade uuid.UUID, m domain.Movimentacao) error {
	if err := fr.hook("InsertMovimentacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertMovimentacao(ctx, db, unidade, m)
}

func (fr *faultRepo) ListMovimentacoes(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Movimentacao, error) {
	if err := fr.hook("ListMovimentacoes"); err != nil {
		var z0 []domain.Movimentacao
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListMovimentacoes(ctx, db, unidade)
}

func (fr *faultRepo) InsertEmprestimo(ctx context.Context, db database.DBTX, e domain.Emprestimo) error {
	if err := fr.hook("InsertEmprestimo"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertEmprestimo(ctx, db, e)
}

func (fr *faultRepo) GetEmprestimo(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Emprestimo, error) {
	if err := fr.hook("GetEmprestimo"); err != nil {
		var z0 domain.Emprestimo
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetEmprestimo(ctx, db, id)
}

func (fr *faultRepo) DevolverEmprestimo(ctx context.Context, db database.DBTX, e domain.Emprestimo) error {
	if err := fr.hook("DevolverEmprestimo"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.DevolverEmprestimo(ctx, db, e)
}

func (fr *faultRepo) EmprestimoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	if err := fr.hook("EmprestimoAberto"); err != nil {
		var z0 bool
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.EmprestimoAberto(ctx, db, unidade)
}

func (fr *faultRepo) ListEmprestimos(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Emprestimo, error) {
	if err := fr.hook("ListEmprestimos"); err != nil {
		var z0 []domain.Emprestimo
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListEmprestimos(ctx, db, unidade)
}

func (fr *faultRepo) ListEmprestimosAbertos(ctx context.Context, db database.DBTX, f domain.Filtro, atrasadosEm *time.Time, p pagination.Params) ([]domain.Emprestimo, int64, error) {
	if err := fr.hook("ListEmprestimosAbertos"); err != nil {
		var z0 []domain.Emprestimo
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListEmprestimosAbertos(ctx, db, f, atrasadosEm, p)
}

func (fr *faultRepo) InsertGuia(ctx context.Context, db database.DBTX, g domain.Guia, unidades []uuid.UUID) (string, error) {
	if err := fr.hook("InsertGuia"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertGuia(ctx, db, g, unidades)
}

func (fr *faultRepo) GetGuia(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Guia, error) {
	if err := fr.hook("GetGuia"); err != nil {
		var z0 domain.Guia
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetGuia(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) SetGuiaSituacao(ctx context.Context, db database.DBTX, g domain.Guia) error {
	if err := fr.hook("SetGuiaSituacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SetGuiaSituacao(ctx, db, g)
}

func (fr *faultRepo) EmGuiaEnviada(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	if err := fr.hook("EmGuiaEnviada"); err != nil {
		var z0 bool
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.EmGuiaEnviada(ctx, db, unidade)
}

func (fr *faultRepo) ListGuias(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Guia, int64, error) {
	if err := fr.hook("ListGuias"); err != nil {
		var z0 []domain.Guia
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListGuias(ctx, db, f, p)
}

func (fr *faultRepo) SalvarComissao(ctx context.Context, db database.DBTX, c domain.Comissao, nova bool) error {
	if err := fr.hook("SalvarComissao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SalvarComissao(ctx, db, c, nova)
}

func (fr *faultRepo) GetComissao(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Comissao, error) {
	if err := fr.hook("GetComissao"); err != nil {
		var z0 domain.Comissao
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetComissao(ctx, db, id)
}

func (fr *faultRepo) ListComissoes(ctx context.Context, db database.DBTX) ([]domain.Comissao, error) {
	if err := fr.hook("ListComissoes"); err != nil {
		var z0 []domain.Comissao
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListComissoes(ctx, db)
}

func (fr *faultRepo) ComissaoDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) (domain.Comissao, error) {
	if err := fr.hook("ComissaoDaUnidade"); err != nil {
		var z0 domain.Comissao
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ComissaoDaUnidade(ctx, db, unidade)
}

func (fr *faultRepo) InsertListagem(ctx context.Context, db database.DBTX, l domain.Listagem) (string, error) {
	if err := fr.hook("InsertListagem"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertListagem(ctx, db, l)
}

func (fr *faultRepo) UpdateListagem(ctx context.Context, db database.DBTX, l domain.Listagem) error {
	if err := fr.hook("UpdateListagem"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.UpdateListagem(ctx, db, l)
}

func (fr *faultRepo) GetListagem(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Listagem, error) {
	if err := fr.hook("GetListagem"); err != nil {
		var z0 domain.Listagem
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetListagem(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) ListListagens(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Listagem, int64, error) {
	if err := fr.hook("ListListagens"); err != nil {
		var z0 []domain.Listagem
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListListagens(ctx, db, f, p)
}

func (fr *faultRepo) AddItemListagem(ctx context.Context, db database.DBTX, listagem uuid.UUID, unidade uuid.UUID) error {
	if err := fr.hook("AddItemListagem"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.AddItemListagem(ctx, db, listagem, unidade)
}

func (fr *faultRepo) RemoveItemListagem(ctx context.Context, db database.DBTX, listagem uuid.UUID, unidade uuid.UUID) error {
	if err := fr.hook("RemoveItemListagem"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.RemoveItemListagem(ctx, db, listagem, unidade)
}

func (fr *faultRepo) MarcarItem(ctx context.Context, db database.DBTX, listagem uuid.UUID, unidade uuid.UUID, amostra bool, retirada bool) error {
	if err := fr.hook("MarcarItem"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.MarcarItem(ctx, db, listagem, unidade, amostra, retirada)
}

func (fr *faultRepo) EmListagemAtiva(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	if err := fr.hook("EmListagemAtiva"); err != nil {
		var z0 bool
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.EmListagemAtiva(ctx, db, unidade)
}

func (fr *faultRepo) InsertPedido(ctx context.Context, db database.DBTX, listagem uuid.UUID, p domain.Pedido) error {
	if err := fr.hook("InsertPedido"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertPedido(ctx, db, listagem, p)
}

func (fr *faultRepo) DecidirPedido(ctx context.Context, db database.DBTX, p domain.Pedido) error {
	if err := fr.hook("DecidirPedido"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.DecidirPedido(ctx, db, p)
}

func (fr *faultRepo) InsertRecolhimento(ctx context.Context, db database.DBTX, r domain.Recolhimento, unidades []uuid.UUID) (string, error) {
	if err := fr.hook("InsertRecolhimento"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertRecolhimento(ctx, db, r, unidades)
}

func (fr *faultRepo) GetRecolhimento(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Recolhimento, error) {
	if err := fr.hook("GetRecolhimento"); err != nil {
		var z0 domain.Recolhimento
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetRecolhimento(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) ListRecolhimentos(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Recolhimento, int64, error) {
	if err := fr.hook("ListRecolhimentos"); err != nil {
		var z0 []domain.Recolhimento
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListRecolhimentos(ctx, db, f, p)
}

func (fr *faultRepo) ReceberRecolhimento(ctx context.Context, db database.DBTX, r domain.Recolhimento) (string, error) {
	if err := fr.hook("ReceberRecolhimento"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ReceberRecolhimento(ctx, db, r)
}

func (fr *faultRepo) SetRecolhimentoSituacao(ctx context.Context, db database.DBTX, r domain.Recolhimento) error {
	if err := fr.hook("SetRecolhimentoSituacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SetRecolhimentoSituacao(ctx, db, r)
}

func (fr *faultRepo) EmRecolhimentoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	if err := fr.hook("EmRecolhimentoAberto"); err != nil {
		var z0 bool
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.EmRecolhimentoAberto(ctx, db, unidade)
}

func (fr *faultRepo) MarcarRecolhida(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) error {
	if err := fr.hook("MarcarRecolhida"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.MarcarRecolhida(ctx, db, u, c)
}

func (fr *faultRepo) GestoresDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]uuid.UUID, error) {
	if err := fr.hook("GestoresDaUnidade"); err != nil {
		var z0 []uuid.UUID
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GestoresDaUnidade(ctx, db, unidade)
}

func (fr *faultRepo) InsertNorma(ctx context.Context, db database.DBTX, n domain.Norma) error {
	if err := fr.hook("InsertNorma"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertNorma(ctx, db, n)
}

func (fr *faultRepo) GetNorma(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Norma, error) {
	if err := fr.hook("GetNorma"); err != nil {
		var z0 domain.Norma
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetNorma(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) ListNormas(ctx context.Context, db database.DBTX, q string, publicas bool, p pagination.Params) ([]domain.Norma, int64, error) {
	if err := fr.hook("ListNormas"); err != nil {
		var z0 []domain.Norma
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListNormas(ctx, db, q, publicas, p)
}

func (fr *faultRepo) UpdateNorma(ctx context.Context, db database.DBTX, n domain.Norma) error {
	if err := fr.hook("UpdateNorma"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.UpdateNorma(ctx, db, n)
}

func (fr *faultRepo) DeleteNorma(ctx context.Context, db database.DBTX, id uuid.UUID) error {
	if err := fr.hook("DeleteNorma"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.DeleteNorma(ctx, db, id)
}

func (fr *faultRepo) SalvarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, r domain.Redacao, nova bool, por *uuid.UUID) error {
	if err := fr.hook("SalvarRedacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.SalvarRedacao(ctx, db, norma, r, nova, por)
}

func (fr *faultRepo) PublicarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, r domain.Redacao, por *uuid.UUID) error {
	if err := fr.hook("PublicarRedacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.PublicarRedacao(ctx, db, norma, r, por)
}

func (fr *faultRepo) ApagarRascunho(ctx context.Context, db database.DBTX, norma uuid.UUID) error {
	if err := fr.hook("ApagarRascunho"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.ApagarRascunho(ctx, db, norma)
}

func (fr *faultRepo) InsertPrivado(ctx context.Context, db database.DBTX, p domain.Privado, por *uuid.UUID) (string, error) {
	if err := fr.hook("InsertPrivado"); err != nil {
		var z0 string
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertPrivado(ctx, db, p, por)
}

func (fr *faultRepo) GetPrivado(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Privado, error) {
	if err := fr.hook("GetPrivado"); err != nil {
		var z0 domain.Privado
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetPrivado(ctx, db, id, paraAtualizar)
}

func (fr *faultRepo) ListPrivados(ctx context.Context, db database.DBTX, q string, p pagination.Params) ([]domain.Privado, int64, error) {
	if err := fr.hook("ListPrivados"); err != nil {
		var z0 []domain.Privado
		var z1 int64
		return z0, z1, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListPrivados(ctx, db, q, p)
}

func (fr *faultRepo) UpdatePrivado(ctx context.Context, db database.DBTX, p domain.Privado) error {
	if err := fr.hook("UpdatePrivado"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.UpdatePrivado(ctx, db, p)
}

func (fr *faultRepo) InsertEventoPrivado(ctx context.Context, db database.DBTX, privado uuid.UUID, e domain.Evento) error {
	if err := fr.hook("InsertEventoPrivado"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertEventoPrivado(ctx, db, privado, e)
}

func (fr *faultRepo) InsertAlienacao(ctx context.Context, db database.DBTX, privado uuid.UUID, a domain.Alienacao, por *uuid.UUID) error {
	if err := fr.hook("InsertAlienacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertAlienacao(ctx, db, privado, a, por)
}

func (fr *faultRepo) DecidirAlienacao(ctx context.Context, db database.DBTX, a domain.Alienacao, por *uuid.UUID) error {
	if err := fr.hook("DecidirAlienacao"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.DecidirAlienacao(ctx, db, a, por)
}

func (fr *faultRepo) ResolverUnidade(ctx context.Context, db database.DBTX, ref string) (uuid.UUID, error) {
	if err := fr.hook("ResolverUnidade"); err != nil {
		var z0 uuid.UUID
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ResolverUnidade(ctx, db, ref)
}

func (fr *faultRepo) InsertLocal(ctx context.Context, db database.DBTX, l domain.Local) error {
	if err := fr.hook("InsertLocal"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.InsertLocal(ctx, db, l)
}

func (fr *faultRepo) UpdateLocal(ctx context.Context, db database.DBTX, l domain.Local) error {
	if err := fr.hook("UpdateLocal"); err != nil {
		return err
	}
	defer fr.post(ctx, db)
	return fr.inner.UpdateLocal(ctx, db, l)
}

func (fr *faultRepo) GetLocal(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Local, error) {
	if err := fr.hook("GetLocal"); err != nil {
		var z0 domain.Local
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.GetLocal(ctx, db, id)
}

func (fr *faultRepo) ListLocais(ctx context.Context, db database.DBTX) ([]domain.Local, error) {
	if err := fr.hook("ListLocais"); err != nil {
		var z0 []domain.Local
		return z0, err
	}
	defer fr.post(ctx, db)
	return fr.inner.ListLocais(ctx, db)
}
