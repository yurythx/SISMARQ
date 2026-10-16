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

func TestDigitalPropagatesDatabaseErrors(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()
	id := uuid.New()
	p := pagination.New(1, 10, 10)
	f := domain.FiltroDigital{UnidadeID: &id, Situacao: domain.DigitalPronto, SoFalhas: true, Restrito: true, Unidades: []uuid.UUID{id}}
	calls := map[string]func(db database.DBTX) error{
		"InsertDigital":        func(db database.DBTX) error { _, err := r.InsertDigital(ctx, db, domain.Digital{}); return err },
		"GetDigital":           func(db database.DBTX) error { _, err := r.GetDigital(ctx, db, id); return err },
		"ListDigitais":         func(db database.DBTX) error { _, _, err := r.ListDigitais(ctx, db, f, p); return err },
		"ProximoDigital":       func(db database.DBTX) error { _, _, err := r.ProximoDigital(ctx, db, time.Now()); return err },
		"ConcluirDigital":      func(db database.DBTX) error { return r.ConcluirDigital(ctx, db, domain.Digital{}, "texto") },
		"ParaVerificar":        func(db database.DBTX) error { _, err := r.ParaVerificar(ctx, db, time.Now(), 10); return err },
		"RegistrarVerificacao": func(db database.DBTX) error { return r.RegistrarVerificacao(ctx, db, id, domain.Verificacao{}) },
		"ListVerificacoes":     func(db database.DBTX) error { _, err := r.ListVerificacoes(ctx, db, id, 10); return err },
		"EquipeCentral":        func(db database.DBTX) error { _, err := r.EquipeCentral(ctx, db); return err },
		"EquipeSEMAD":          func(db database.DBTX) error { _, err := r.EquipeSEMAD(ctx, db); return err },
	}
	for name, call := range calls {
		if err := call(dbtest.Fail{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com o banco fora: %v", name, err)
		}
	}
	for _, name := range []string{"ParaVerificar", "ListVerificacoes", "EquipeCentral", "EquipeSEMAD"} {
		if err := calls[name](dbtest.ScanFail{}); err == nil {
			t.Errorf("%s com linha ilegível deveria falhar", name)
		}
		if err := calls[name](dbtest.RowsErr{}); !errors.Is(err, dbtest.ErrInjected) {
			t.Errorf("%s com erro durante a leitura: %v", name, err)
		}
	}
	// Segunda gravação falha (o texto do OCR; o estado da verificação).
	if err := r.ConcluirDigital(ctx, &dbtest.Seq{Execs: []dbtest.ExecResult{dbtest.OK, {Err: dbtest.ErrInjected}}}, domain.Digital{}, "texto"); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("texto do OCR: %v", err)
	}
	if err := r.RegistrarVerificacao(ctx, &dbtest.Seq{Execs: []dbtest.ExecResult{dbtest.OK, {Err: dbtest.ErrInjected}}}, id, domain.Verificacao{}); !errors.Is(err, dbtest.ErrInjected) {
		t.Errorf("estado da verificação: %v", err)
	}
	for _, nova := range []func() dbtest.QueryResult{
		func() dbtest.QueryResult { return dbtest.QueryResult{Err: dbtest.ErrInjected} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.BadRows()} },
		func() dbtest.QueryResult { return dbtest.QueryResult{Rows: dbtest.ErrRows()} },
	} {
		if _, _, err := r.ListDigitais(ctx, &dbtest.Seq{Row: countRow{}, Queries: []dbtest.QueryResult{nova()}}, domain.FiltroDigital{}, p); err == nil {
			t.Error("página de representantes ilegível")
		}
	}
}
