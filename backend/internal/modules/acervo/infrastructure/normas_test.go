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

func TestNormasEPrivadosPropagateDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	calls := map[string]func(db database.DBTX) error{
		"InsertNorma":         func(db database.DBTX) error { return r.InsertNorma(ctx, db, domain.Norma{}) },
		"GetNorma":            func(db database.DBTX) error { _, err := r.GetNorma(ctx, db, id, true); return err },
		"ListNormas":          func(db database.DBTX) error { _, _, err := r.ListNormas(ctx, db, "x", true, p); return err },
		"UpdateNorma":         func(db database.DBTX) error { return r.UpdateNorma(ctx, db, domain.Norma{}) },
		"DeleteNorma":         func(db database.DBTX) error { return r.DeleteNorma(ctx, db, id) },
		"SalvarRedacao":       func(db database.DBTX) error { return r.SalvarRedacao(ctx, db, id, domain.Redacao{}, false, nil) },
		"PublicarRedacao":     func(db database.DBTX) error { return r.PublicarRedacao(ctx, db, id, domain.Redacao{}, nil) },
		"ApagarRascunho":      func(db database.DBTX) error { return r.ApagarRascunho(ctx, db, id) },
		"InsertPrivado":       func(db database.DBTX) error { _, err := r.InsertPrivado(ctx, db, domain.Privado{}, nil); return err },
		"GetPrivado":          func(db database.DBTX) error { _, err := r.GetPrivado(ctx, db, id, true); return err },
		"ListPrivados":        func(db database.DBTX) error { _, _, err := r.ListPrivados(ctx, db, "x", p); return err },
		"UpdatePrivado":       func(db database.DBTX) error { return r.UpdatePrivado(ctx, db, domain.Privado{}) },
		"InsertEventoPrivado": func(db database.DBTX) error { return r.InsertEventoPrivado(ctx, db, id, domain.Evento{}) },
		"InsertAlienacao":     func(db database.DBTX) error { return r.InsertAlienacao(ctx, db, id, domain.Alienacao{}, nil) },
		"DecidirAlienacao":    func(db database.DBTX) error { return r.DecidirAlienacao(ctx, db, domain.Alienacao{}, nil) },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	// Cada chamada recebe linhas novas (as falsas se consomem na leitura).
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, err := r.GetNorma(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, false); err == nil {
			t.Error("redações ilegíveis")
		}
		if _, _, err := r.ListNormas(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, "", false, p); err == nil {
			t.Error("página de normas ilegível")
		}
		if _, err := r.GetPrivado(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, false); err == nil {
			t.Error("histórico ilegível")
		}
		if _, err := r.GetPrivado(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{{Rows: dbtest.NoRows()}, nova()}}, id, false); err == nil {
			t.Error("alienações ilegíveis")
		}
		if _, _, err := r.ListPrivados(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, "", p); err == nil {
			t.Error("página de privados ilegível")
		}
	}
}
