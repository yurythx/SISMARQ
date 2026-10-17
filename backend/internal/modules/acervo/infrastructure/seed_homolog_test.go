package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
)

func TestSeedHomologAcervo(t *testing.T) {
	pool := dbtest.Pool(t)
	ctx := context.Background()
	_, this, _, _ := runtime.Caller(0)
	sqlPath := filepath.Join(filepath.Dir(this), "..", "..", "..", "..", "..", "deploy", "homolog", "seed-acervo.sql")
	raw, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO entidades (nome, sigla, slug)
		SELECT 'Administração', 'SEMAD', 'homolog-semad-' || gen_random_uuid()::text
		WHERE NOT EXISTS (SELECT 1 FROM entidades WHERE sigla = 'SEMAD')`); err != nil {
		t.Fatal(err)
	}
	var ent string
	if err := pool.QueryRow(ctx, `SELECT id FROM entidades WHERE sigla = 'SEMAD' LIMIT 1`).Scan(&ent); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO unidades (entidade_id, nome, sigla, slug)
		SELECT $1::uuid, 'Sede', 'SEMAD', 'sede-homolog-' || gen_random_uuid()::text
		WHERE NOT EXISTS (SELECT 1 FROM unidades WHERE sigla = 'SEMAD')`, ent); err != nil {
		t.Fatal(err)
	}

	if _, err := pool.Exec(ctx, string(raw)); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(raw)); err != nil {
		t.Fatal(err)
	}

	rows, err := pool.Query(ctx, `SELECT numero, situacao FROM acervo_unidades WHERE numero LIKE 'HOM-UA-%' ORDER BY numero`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	quer := map[string]string{
		"HOM-UA-01": "AVALIAR_ELIMINACAO",
		"HOM-UA-02": "RECOLHER",
		"HOM-UA-03": "PERMANENTE",
		"HOM-UA-04": "AGUARDA_EVENTO",
		"HOM-UA-05": "TRANSFERIR",
		"HOM-UA-06": "NO_CORRENTE",
		"HOM-UA-07": "BLOQUEADA",
		"HOM-UA-08": "TRANSFERIR",
		"HOM-UA-09": "AVALIAR_ELIMINACAO",
		"HOM-UA-10": "AVALIAR_ELIMINACAO",
	}
	got := map[string]string{}
	for rows.Next() {
		var n, s string
		if err := rows.Scan(&n, &s); err != nil {
			t.Fatal(err)
		}
		got[n] = s
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != len(quer) {
		t.Fatalf("caixas: %d queria %d (%v)", len(got), len(quer), got)
	}
	for n, s := range quer {
		if got[n] != s && got[n] != "SERIE_REVOGADA" {
			t.Errorf("%s situacao=%s queria %s", n, got[n], s)
		}
	}

	var bloqueios int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_bloqueios b
		JOIN acervo_unidades u ON u.id = b.unidade_id WHERE u.numero = 'HOM-UA-07' AND b.liberado_em IS NULL`).Scan(&bloqueios); err != nil {
		t.Fatal(err)
	}
	if bloqueios != 1 {
		t.Fatalf("bloqueio HOM-UA-07: %d", bloqueios)
	}

	digSQL := filepath.Join(filepath.Dir(this), "..", "..", "..", "..", "..", "deploy", "homolog", "seed-digitais.sql")
	rawDig, err := os.ReadFile(digSQL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(rawDig)); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(rawDig)); err != nil {
		t.Fatal(err)
	}
	var origemScan, origemNato string
	var dpiNato int
	if err := pool.QueryRow(ctx, `SELECT origem FROM acervo_digitais WHERE numero = 'RD-HOM-SCAN'`).Scan(&origemScan); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT origem, dpi_declarado FROM acervo_digitais WHERE numero = 'RD-HOM-NATO'`).Scan(&origemNato, &dpiNato); err != nil {
		t.Fatal(err)
	}
	if origemScan != "SCAN" || origemNato != "NATO_DIGITAL" || dpiNato != 0 {
		t.Fatalf("digitais: SCAN=%s NATO=%s dpi=%d", origemScan, origemNato, dpiNato)
	}
}
