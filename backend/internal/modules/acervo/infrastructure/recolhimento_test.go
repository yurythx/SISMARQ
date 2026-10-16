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

func TestRecolhimentoPropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	f := domain.Filtro{Restrito: true, Unidades: []uuid.UUID{id}}
	calls := map[string]func(db database.DBTX) error{
		"InsertRecolhimento": func(db database.DBTX) error {
			_, err := r.InsertRecolhimento(ctx, db, domain.Recolhimento{}, nil)
			return err
		},
		"GetRecolhimento":   func(db database.DBTX) error { _, err := r.GetRecolhimento(ctx, db, id, true); return err },
		"ListRecolhimentos": func(db database.DBTX) error { _, _, err := r.ListRecolhimentos(ctx, db, f, p); return err },
		"ReceberRecolhimento": func(db database.DBTX) error {
			_, err := r.ReceberRecolhimento(ctx, db, domain.Recolhimento{})
			return err
		},
		"SetRecolhimentoSituacao": func(db database.DBTX) error { return r.SetRecolhimentoSituacao(ctx, db, domain.Recolhimento{}) },
		"EmRecolhimentoAberto":    func(db database.DBTX) error { _, err := r.EmRecolhimentoAberto(ctx, db, id); return err },
		"MarcarRecolhida":         func(db database.DBTX) error { return r.MarcarRecolhida(ctx, db, domain.Unidade{}, domain.Calculo{}) },
		"GestoresDaUnidade":       func(db database.DBTX) error { _, err := r.GestoresDaUnidade(ctx, db, id); return err },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	if err := calls["GestoresDaUnidade"](dbtest.ScanFail{}); err == nil {
		t.Error("gestor ilegível")
	}
	if err := calls["GestoresDaUnidade"](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("leitura dos gestores: %v", err)
	}
	// Itens do recolhimento: a gravação falha depois do cabeçalho.
	if _, err := r.InsertRecolhimento(ctx, &dbtest.Seq{Row: okRow{}}, domain.Recolhimento{}, nil); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("itens do recolhimento: %v", err)
	}
	// Cada chamada recebe linhas novas (as falsas se consomem na leitura).
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, err := r.GetRecolhimento(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, false); err == nil {
			t.Error("unidades do recolhimento ilegíveis")
		}
		if _, _, err := r.ListRecolhimentos(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, f, p); err == nil {
			t.Error("página de recolhimentos ilegível")
		}
	}
}
