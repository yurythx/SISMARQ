package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// InsertNorma grava a norma (as redações vão por SalvarRedacao). Número
// repetido vira ErrNumeroNorma.
func (r *Repository) InsertNorma(ctx context.Context, db database.DBTX, n domain.Norma) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_normas (id, numero, situacao, criada_em) VALUES ($1, $2, $3, $4)`,
		n.ID, n.Numero, n.Situacao, n.CriadaEm); err != nil {
		if unico(err) {
			return domain.ErrNumeroNorma
		}
		return fmt.Errorf("acervo.InsertNorma: %w", err)
	}
	return nil
}

// GetNorma busca a norma com todas as redações (a mais nova primeiro).
func (r *Repository) GetNorma(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Norma, error) {
	q := `SELECT id, numero, situacao, revogada_em, revogacao_ato, criada_em FROM acervo_normas WHERE id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	var n domain.Norma
	err := db.QueryRow(ctx, q, id).Scan(&n.ID, &n.Numero, &n.Situacao, &n.RevogadaEm, &n.RevogacaoAto, &n.CriadaEm)
	if errors.Is(err, pgx.ErrNoRows) {
		return n, domain.ErrNormaNaoExiste
	}
	if err != nil {
		return n, fmt.Errorf("acervo.GetNorma: %w", err)
	}
	rows, err := db.Query(ctx, `SELECT versao, ementa, texto, motivo, situacao, publicada_em, diario, criada_em
		FROM acervo_norma_redacoes WHERE norma_id = $1 ORDER BY versao DESC`, id)
	if err != nil {
		return n, fmt.Errorf("acervo.GetNorma redações: %w", err)
	}
	defer rows.Close()
	n.Redacoes = []domain.Redacao{}
	for rows.Next() {
		var x domain.Redacao
		if err := rows.Scan(&x.Versao, &x.Ementa, &x.Texto, &x.Motivo, &x.Situacao, &x.PublicadaEm, &x.Diario, &x.CriadaEm); err != nil {
			return n, fmt.Errorf("acervo.GetNorma redações scan: %w", err)
		}
		n.Redacoes = append(n.Redacoes, x)
	}
	if err := rows.Err(); err != nil {
		return n, fmt.Errorf("acervo.GetNorma redações rows: %w", err)
	}
	n.Montar()
	return n, nil
}

// ListNormas pagina as normas pelo número (mais novas primeiro), com a
// ementa e a publicação da redação em vigor. publicas: só as já publicadas.
func (r *Repository) ListNormas(ctx context.Context, db database.DBTX, q string, publicas bool, p pagination.Params) ([]domain.Norma, int64, error) {
	conds, args := "true", []any{}
	if publicas {
		conds += " AND n.situacao <> 'RASCUNHO'"
	}
	if q != "" {
		args = append(args, "%"+q+"%")
		conds += ` AND (n.numero ILIKE $1 OR EXISTS (SELECT 1 FROM acervo_norma_redacoes x WHERE x.norma_id = n.id
			AND x.situacao = 'PUBLICADA' AND (x.ementa ILIKE $1 OR x.texto ILIKE $1)))`
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_normas n WHERE `+conds, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListNormas count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT n.id, n.numero, n.situacao, n.revogada_em, n.revogacao_ato, n.criada_em,
		COALESCE(v.ementa, rs.ementa, ''), COALESCE(v.versao, 0), v.publicada_em, rs.versao IS NOT NULL,
		(SELECT COUNT(*) FROM acervo_norma_redacoes c WHERE c.norma_id = n.id AND c.situacao = 'PUBLICADA')
		FROM acervo_normas n
		LEFT JOIN LATERAL (SELECT ementa, versao, publicada_em FROM acervo_norma_redacoes x
			WHERE x.norma_id = n.id AND x.situacao = 'PUBLICADA' ORDER BY x.versao DESC LIMIT 1) v ON true
		LEFT JOIN acervo_norma_redacoes rs ON rs.norma_id = n.id AND rs.situacao = 'RASCUNHO'
		WHERE `+conds+` ORDER BY split_part(n.numero, '/', 2)::int DESC, split_part(n.numero, '/', 1)::int DESC
		LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListNormas: %w", err)
	}
	defer rows.Close()
	out := []domain.Norma{}
	for rows.Next() {
		var n domain.Norma
		if err := rows.Scan(&n.ID, &n.Numero, &n.Situacao, &n.RevogadaEm, &n.RevogacaoAto, &n.CriadaEm, &n.Ementa, &n.Versao,
			&n.PublicadaEm, &n.TemRascunho, &n.TotalRedacoes); err != nil {
			return nil, 0, fmt.Errorf("acervo.ListNormas scan: %w", err)
		}
		if publicas {
			n.TemRascunho = false
		}
		n.Redacoes = []domain.Redacao{}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListNormas rows: %w", err)
	}
	return out, total, nil
}

// UpdateNorma grava a situação e a revogação.
func (r *Repository) UpdateNorma(ctx context.Context, db database.DBTX, n domain.Norma) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_normas SET situacao = $2, revogada_em = $3, revogacao_ato = $4 WHERE id = $1`,
		n.ID, n.Situacao, n.RevogadaEm, n.RevogacaoAto); err != nil {
		return fmt.Errorf("acervo.UpdateNorma: %w", err)
	}
	return nil
}

// DeleteNorma apaga a norma nunca publicada (o rascunho descartado).
func (r *Repository) DeleteNorma(ctx context.Context, db database.DBTX, id uuid.UUID) error {
	if _, err := db.Exec(ctx, `DELETE FROM acervo_normas WHERE id = $1 AND situacao = 'RASCUNHO'`, id); err != nil {
		return fmt.Errorf("acervo.DeleteNorma: %w", err)
	}
	return nil
}

// SalvarRedacao grava a redação em rascunho: nova (versão seguinte) ou a
// alteração do rascunho existente.
func (r *Repository) SalvarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, x domain.Redacao, nova bool, por *uuid.UUID) error {
	q := `UPDATE acervo_norma_redacoes SET ementa = $3, texto = $4, motivo = $5 WHERE norma_id = $1 AND versao = $2 AND situacao = 'RASCUNHO'`
	if nova {
		q = `INSERT INTO acervo_norma_redacoes (norma_id, versao, ementa, texto, motivo, criada_por) VALUES ($1, $2, $3, $4, $5, $6)`
	}
	args := []any{norma, x.Versao, x.Ementa, x.Texto, x.Motivo}
	if nova {
		args = append(args, por)
	}
	if _, err := db.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("acervo.SalvarRedacao: %w", err)
	}
	return nil
}

// PublicarRedacao publica o rascunho.
func (r *Repository) PublicarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, x domain.Redacao, por *uuid.UUID) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_norma_redacoes SET situacao = 'PUBLICADA', publicada_em = $3, diario = $4, publicada_por = $5
		WHERE norma_id = $1 AND versao = $2 AND situacao = 'RASCUNHO'`, norma, x.Versao, x.PublicadaEm, x.Diario, por); err != nil {
		return fmt.Errorf("acervo.PublicarRedacao: %w", err)
	}
	return nil
}

// ApagarRascunho descarta a redação em rascunho.
func (r *Repository) ApagarRascunho(ctx context.Context, db database.DBTX, norma uuid.UUID) error {
	if _, err := db.Exec(ctx, `DELETE FROM acervo_norma_redacoes WHERE norma_id = $1 AND situacao = 'RASCUNHO'`, norma); err != nil {
		return fmt.Errorf("acervo.ApagarRascunho: %w", err)
	}
	return nil
}
