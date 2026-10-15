package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

// okRow responde a qualquer Scan sem erro (valores zero).
type okRow struct{}

func (okRow) Scan(...any) error { return nil }

func TestFase2PropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	hoje := time.Now()
	f := domain.Filtro{Restrito: true, Unidades: []uuid.UUID{id}}
	calls := map[string]func(db database.DBTX) error{
		"InsertEmprestimo":   func(db database.DBTX) error { return r.InsertEmprestimo(ctx, db, domain.Emprestimo{}) },
		"GetEmprestimo":      func(db database.DBTX) error { _, err := r.GetEmprestimo(ctx, db, id); return err },
		"DevolverEmprestimo": func(db database.DBTX) error { return r.DevolverEmprestimo(ctx, db, domain.Emprestimo{}) },
		"EmprestimoAberto":   func(db database.DBTX) error { _, err := r.EmprestimoAberto(ctx, db, id); return err },
		"ListEmprestimos":    func(db database.DBTX) error { _, err := r.ListEmprestimos(ctx, db, id); return err },
		"ListEmprestimosAbertos": func(db database.DBTX) error {
			_, _, err := r.ListEmprestimosAbertos(ctx, db, f, &hoje, p)
			return err
		},
		"InsertGuia":      func(db database.DBTX) error { _, err := r.InsertGuia(ctx, db, domain.Guia{}, nil); return err },
		"GetGuia":         func(db database.DBTX) error { _, err := r.GetGuia(ctx, db, id, true); return err },
		"SetGuiaSituacao": func(db database.DBTX) error { return r.SetGuiaSituacao(ctx, db, domain.Guia{}) },
		"EmGuiaEnviada":   func(db database.DBTX) error { _, err := r.EmGuiaEnviada(ctx, db, id); return err },
		"ListGuias":       func(db database.DBTX) error { _, _, err := r.ListGuias(ctx, db, f, p); return err },
		"ResolverUnidade": func(db database.DBTX) error { _, err := r.ResolverUnidade(ctx, db, "X"); return err },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	for _, name := range []string{"ListEmprestimos", "ResolverUnidade"} {
		if err := calls[name](dbtest.ScanFail{}); err == nil {
			t.Errorf("%s com linha ilegível deveria falhar", name)
		}
		if err := calls[name](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com erro durante a leitura: %v", name, err)
		}
	}
	// Primeira consulta ok, segunda falha (itens, página).
	if _, err := r.InsertGuia(ctx, &dbtest.Seq{Row: okRow{}}, domain.Guia{}, nil); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("itens da guia: %v", err)
	}
	// Cada chamada recebe linhas novas (as falsas se consomem na leitura).
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, err := r.GetGuia(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, false); err == nil {
			t.Error("itens da guia ilegíveis")
		}
		if _, _, err := r.ListGuias(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, f, p); err == nil {
			t.Error("página de guias ilegível")
		}
		if _, _, err := r.ListEmprestimosAbertos(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, f, nil, p); err == nil {
			t.Error("página de empréstimos ilegível")
		}
	}
}
