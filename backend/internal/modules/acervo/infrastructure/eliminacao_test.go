package infrastructure

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

// umaLinha entrega uma linha legível (valores zero) e acaba.
type umaLinha struct{ lida bool }

func (r *umaLinha) Close()                                       {}
func (r *umaLinha) Err() error                                   { return nil }
func (r *umaLinha) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *umaLinha) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *umaLinha) Scan(...any) error                            { return nil }
func (r *umaLinha) Values() ([]any, error)                       { return nil, nil }
func (r *umaLinha) RawValues() [][]byte                          { return nil }
func (r *umaLinha) Conn() *pgx.Conn                              { return nil }
func (r *umaLinha) Next() bool {
	if r.lida {
		return false
	}
	r.lida = true
	return true
}

func TestEliminacaoPropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	f := domain.Filtro{Restrito: true, Unidades: []uuid.UUID{id}}
	calls := map[string]func(db database.DBTX) error{
		"SalvarComissao":    func(db database.DBTX) error { return r.SalvarComissao(ctx, db, domain.Comissao{}, true) },
		"GetComissao":       func(db database.DBTX) error { _, err := r.GetComissao(ctx, db, id); return err },
		"ComissaoDaUnidade": func(db database.DBTX) error { _, err := r.ComissaoDaUnidade(ctx, db, id); return err },
		"ListComissoes":     func(db database.DBTX) error { _, err := r.ListComissoes(ctx, db); return err },
		"InsertListagem":    func(db database.DBTX) error { _, err := r.InsertListagem(ctx, db, domain.Listagem{}); return err },
		"UpdateListagem":    func(db database.DBTX) error { return r.UpdateListagem(ctx, db, domain.Listagem{}) },
		"GetListagem":       func(db database.DBTX) error { _, err := r.GetListagem(ctx, db, id, false); return err },
		"ListListagens":     func(db database.DBTX) error { _, _, err := r.ListListagens(ctx, db, f, p); return err },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	if err := calls["ListComissoes"](dbtest.ScanFail{}); err == nil {
		t.Error("comissão ilegível")
	}
	if err := calls["ListComissoes"](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("leitura das comissões: %v", err)
	}

	// Comissão: a troca dos membros falha (exclusão, depois inclusão).
	membro := domain.Comissao{Membros: []domain.Membro{{}}}
	for _, execs := range [][]dbtest.ExecResult{{dbtest.OK, {Err: dbtest.ErrInjected}}, {dbtest.OK, dbtest.OK, {Err: dbtest.ErrInjected}}} {
		if err := r.SalvarComissao(ctx, &dbtest.Seq{Execs: execs}, membro, false); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("membros da comissão: %v", err)
		}
	}

	// Cada chamada recebe linhas novas (as falsas se consomem na leitura).
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, err := r.GetComissao(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id); err == nil {
			t.Error("membros ilegíveis")
		}
		if _, err := r.ListComissoes(ctx, &dbtest.Seq{Queries: []dbtest.QueryResult{{Rows: &umaLinha{}}, nova()}}); err == nil {
			t.Error("membros da lista ilegíveis")
		}
		if _, err := r.GetListagem(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{nova()}}, id, true); err == nil {
			t.Error("itens da relação ilegíveis")
		}
		if _, err := r.GetListagem(ctx, &dbtest.Seq{Row: okRow{}, Queries: []dbtest.QueryResult{{Rows: dbtest.NoRows()}, nova()}}, id, false); err == nil {
			t.Error("pedidos da relação ilegíveis")
		}
		if _, _, err := r.ListListagens(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, f, p); err == nil {
			t.Error("página de relações ilegível")
		}
	}
}
