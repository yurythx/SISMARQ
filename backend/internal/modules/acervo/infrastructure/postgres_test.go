package infrastructure

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

// countRow responde à contagem da listagem.
type countRow struct{}

func (countRow) Scan(dest ...any) error {
	*(dest[0].(*int64)) = 1
	return nil
}

func TestRepositoryPropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	q := domain.Filtro{Query: "caixa", Fase: domain.FaseCorrente, Situacao: domain.SituacaoNoCorrente, Codigo: "2.0", UnidadeID: &id,
		Restrito: true, Unidades: []uuid.UUID{id}}
	calls := map[string]func(db database.DBTX) error{
		"InsertUnidade": func(db database.DBTX) error {
			_, err := r.InsertUnidade(ctx, db, domain.Unidade{}, domain.Calculo{})
			return err
		},
		"GetUnidade":    func(db database.DBTX) error { _, _, err := r.GetUnidade(ctx, db, id, false); return err },
		"UpdateUnidade": func(db database.DBTX) error { return r.UpdateUnidade(ctx, db, domain.Unidade{}, domain.Calculo{}) },
		"SetCalculo":    func(db database.DBTX) error { return r.SetCalculo(ctx, db, id, domain.Calculo{}) },
		"SetFase": func(db database.DBTX) error {
			return r.SetFase(ctx, db, id, domain.FaseIntermediaria, domain.Calculo{})
		},
		"ListUnidades":     func(db database.DBTX) error { _, _, err := r.ListUnidades(ctx, db, q, p); return err },
		"Painel":           func(db database.DBTX) error { _, err := r.Painel(ctx, db, q); return err },
		"ParaRecalculo":    func(db database.DBTX) error { _, err := r.ParaRecalculo(ctx, db, uuid.Nil, 10); return err },
		"InsertBloqueio":   func(db database.DBTX) error { return r.InsertBloqueio(ctx, db, domain.Bloqueio{}) },
		"GetBloqueio":      func(db database.DBTX) error { _, err := r.GetBloqueio(ctx, db, id); return err },
		"LiberarBloqueio":  func(db database.DBTX) error { return r.LiberarBloqueio(ctx, db, domain.Bloqueio{}) },
		"ListBloqueios":    func(db database.DBTX) error { _, err := r.ListBloqueios(ctx, db, id); return err },
		"TemBloqueioAtivo": func(db database.DBTX) error { _, err := r.TemBloqueioAtivo(ctx, db, id); return err },
		"InsertMovimentacao": func(db database.DBTX) error {
			return r.InsertMovimentacao(ctx, db, id, domain.Movimentacao{})
		},
		"ListMovimentacoes": func(db database.DBTX) error { _, err := r.ListMovimentacoes(ctx, db, id); return err },
		"InsertLocal":       func(db database.DBTX) error { return r.InsertLocal(ctx, db, domain.Local{}) },
		"UpdateLocal":       func(db database.DBTX) error { return r.UpdateLocal(ctx, db, domain.Local{}) },
		"GetLocal":          func(db database.DBTX) error { _, err := r.GetLocal(ctx, db, id); return err },
		"ListLocais":        func(db database.DBTX) error { _, err := r.ListLocais(ctx, db); return err },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	for _, name := range []string{"Painel", "ParaRecalculo", "ListBloqueios", "ListMovimentacoes", "ListLocais"} {
		if err := calls[name](dbtest.ScanFail{}); err == nil {
			t.Errorf("%s com linha ilegível deveria falhar", name)
		}
		if err := calls[name](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com erro durante a leitura: %v", name, err)
		}
	}
	page := &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{{Err: dbtest.ErrInjected}}}
	if _, _, err := r.ListUnidades(ctx, page, q, p); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("contagem ok, página falha: %v", err)
	}
	for _, rows := range []dbtest.QueryResult{{Rows: dbtest.BadRows()}, {Rows: dbtest.ErrRows()}} {
		if _, _, err := r.ListUnidades(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{rows}}, q, p); err == nil {
			t.Error("página ilegível ou interrompida")
		}
	}
}
