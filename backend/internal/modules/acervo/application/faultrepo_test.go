package application_test

// Código gerado por scripts/genfault.py a partir da interface Repository do
// domínio. Não edite à mão: regenere se a interface mudar.

import (
	"context"
	"errors"

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
