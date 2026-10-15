package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

const emprestimoCols = `e.id, e.unidade_id, e.solicitante_id, COALESCE(us.display_name, ''), e.finalidade, e.devolver_ate,
	e.retirado_por, e.retirado_em, e.devolvido_em, e.recebido_por, e.observacao, u.numero, u.assunto`

const emprestimoFrom = ` FROM acervo_emprestimos e JOIN acervo_unidades u ON u.id = e.unidade_id
	LEFT JOIN users us ON us.id = e.solicitante_id`

func scanEmprestimo(row pgx.Row) (domain.Emprestimo, error) {
	var e domain.Emprestimo
	err := row.Scan(&e.ID, &e.UnidadeID, &e.SolicitanteID, &e.Solicitante, &e.Finalidade, &e.DevolverAte,
		&e.RetiradoPor, &e.RetiradoEm, &e.DevolvidoEm, &e.RecebidoPor, &e.Observacao, &e.Numero, &e.Assunto)
	return e, err
}

func unico(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}

// InsertEmprestimo registra a saída (o índice único garante um só aberto).
func (r *Repository) InsertEmprestimo(ctx context.Context, db database.DBTX, e domain.Emprestimo) error {
	_, err := db.Exec(ctx, `INSERT INTO acervo_emprestimos (id, unidade_id, solicitante_id, finalidade, devolver_ate, retirado_por, retirado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, e.ID, e.UnidadeID, e.SolicitanteID, e.Finalidade, e.DevolverAte, e.RetiradoPor, e.RetiradoEm)
	if unico(err) {
		return domain.ErrEmprestimoAberto
	}
	if err != nil {
		return fmt.Errorf("acervo.InsertEmprestimo: %w", err)
	}
	return nil
}

// GetEmprestimo busca o empréstimo (FOR UPDATE: só é lido para devolver).
func (r *Repository) GetEmprestimo(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Emprestimo, error) {
	e, err := scanEmprestimo(db.QueryRow(ctx, `SELECT `+emprestimoCols+emprestimoFrom+` WHERE e.id = $1 FOR UPDATE OF e`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return e, domain.ErrEmprestimoNaoExiste
	}
	if err != nil {
		return e, fmt.Errorf("acervo.GetEmprestimo: %w", err)
	}
	return e, nil
}

// DevolverEmprestimo grava a devolução.
func (r *Repository) DevolverEmprestimo(ctx context.Context, db database.DBTX, e domain.Emprestimo) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_emprestimos SET devolvido_em = $2, recebido_por = $3, observacao = $4 WHERE id = $1`,
		e.ID, e.DevolvidoEm, e.RecebidoPor, e.Observacao); err != nil {
		return fmt.Errorf("acervo.DevolverEmprestimo: %w", err)
	}
	return nil
}

// EmprestimoAberto informa se a unidade está fora.
func (r *Repository) EmprestimoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM acervo_emprestimos WHERE unidade_id = $1 AND devolvido_em IS NULL)`,
		unidade).Scan(&ok); err != nil {
		return false, fmt.Errorf("acervo.EmprestimoAberto: %w", err)
	}
	return ok, nil
}

func listarEmprestimos(ctx context.Context, db database.DBTX, q string, args ...any) ([]domain.Emprestimo, error) {
	rows, err := db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListEmprestimos: %w", err)
	}
	defer rows.Close()
	out := []domain.Emprestimo{}
	for rows.Next() {
		e, err := scanEmprestimo(rows)
		if err != nil {
			return nil, fmt.Errorf("acervo.ListEmprestimos scan: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListEmprestimos rows: %w", err)
	}
	return out, nil
}

// ListEmprestimos devolve o histórico de empréstimos da unidade.
func (r *Repository) ListEmprestimos(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Emprestimo, error) {
	return listarEmprestimos(ctx, db, `SELECT `+emprestimoCols+emprestimoFrom+` WHERE e.unidade_id = $1 ORDER BY e.retirado_em DESC`, unidade)
}

// ListEmprestimosAbertos pagina os empréstimos em aberto no escopo (com
// atrasadosEm, só os vencidos nessa data), do prazo mais antigo ao mais novo.
func (r *Repository) ListEmprestimosAbertos(ctx context.Context, db database.DBTX, f domain.Filtro, atrasadosEm *time.Time, p pagination.Params) ([]domain.Emprestimo, int64, error) {
	where, args := filtroSQL(f)
	where += " AND e.devolvido_em IS NULL"
	if atrasadosEm != nil {
		args = append(args, *atrasadosEm)
		where += " AND e.devolver_ate < $" + strconv.Itoa(len(args))
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*)`+emprestimoFrom+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListEmprestimosAbertos count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	out, err := listarEmprestimos(ctx, db, `SELECT `+emprestimoCols+emprestimoFrom+where+
		` ORDER BY e.devolver_ate, e.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	return out, total, err
}

const guiaCols = `g.id, g.numero, g.origem_id, g.situacao, g.observacao, g.criada_por, g.criada_em, g.recebida_por, g.recebida_em,
	(SELECT COUNT(*) FROM acervo_guia_itens i WHERE i.guia_id = g.id)`

func scanGuia(row pgx.Row) (domain.Guia, error) {
	var g domain.Guia
	err := row.Scan(&g.ID, &g.Numero, &g.OrigemID, &g.Situacao, &g.Observacao, &g.CriadaPor, &g.CriadaEm, &g.RecebidaPor, &g.RecebidaEm, &g.TotalItens)
	return g, err
}

// InsertGuia grava a guia e os itens; devolve o número (GT-AAAA-NNNNN).
func (r *Repository) InsertGuia(ctx context.Context, db database.DBTX, g domain.Guia, unidades []uuid.UUID) (string, error) {
	var numero string
	if err := db.QueryRow(ctx, `INSERT INTO acervo_guias (id, origem_id, situacao, observacao, criada_por, criada_em)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING numero`, g.ID, g.OrigemID, g.Situacao, g.Observacao, g.CriadaPor, g.CriadaEm).Scan(&numero); err != nil {
		return "", fmt.Errorf("acervo.InsertGuia: %w", err)
	}
	if _, err := db.Exec(ctx, `INSERT INTO acervo_guia_itens (guia_id, unidade_id) SELECT $1, unnest($2::uuid[])`, g.ID, unidades); err != nil {
		return "", fmt.Errorf("acervo.InsertGuia itens: %w", err)
	}
	return numero, nil
}

// GetGuia busca a guia com os itens.
func (r *Repository) GetGuia(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Guia, error) {
	q := `SELECT ` + guiaCols + ` FROM acervo_guias g WHERE g.id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	g, err := scanGuia(db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return g, domain.ErrGuiaNaoExiste
	}
	if err != nil {
		return g, fmt.Errorf("acervo.GetGuia: %w", err)
	}
	rows, err := db.Query(ctx, caminhoLocal+` SELECT `+unidadeCols+`, COALESCE(a.caminho, '') FROM acervo_guia_itens i
		JOIN acervo_unidades u ON u.id = i.unidade_id LEFT JOIN arvore a ON a.id = u.local_id
		WHERE i.guia_id = $1 ORDER BY u.numero`, id)
	if err != nil {
		return g, fmt.Errorf("acervo.GetGuia itens: %w", err)
	}
	defer rows.Close()
	g.Itens = []domain.UnidadeLista{}
	for rows.Next() {
		var l domain.UnidadeLista
		if l.Unidade, l.Calculo, err = scanUnidade(rows, &l.LocalCaminho); err != nil {
			return g, fmt.Errorf("acervo.GetGuia itens scan: %w", err)
		}
		g.Itens = append(g.Itens, l)
	}
	if err := rows.Err(); err != nil {
		return g, fmt.Errorf("acervo.GetGuia itens rows: %w", err)
	}
	return g, nil
}

// SetGuiaSituacao grava o recebimento ou o cancelamento.
func (r *Repository) SetGuiaSituacao(ctx context.Context, db database.DBTX, g domain.Guia) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_guias SET situacao = $2, recebida_por = $3, recebida_em = $4 WHERE id = $1`,
		g.ID, g.Situacao, g.RecebidaPor, g.RecebidaEm); err != nil {
		return fmt.Errorf("acervo.SetGuiaSituacao: %w", err)
	}
	return nil
}

// EmGuiaEnviada informa se a unidade está numa guia aguardando o recebimento.
func (r *Repository) EmGuiaEnviada(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM acervo_guia_itens i JOIN acervo_guias g ON g.id = i.guia_id
		WHERE i.unidade_id = $1 AND g.situacao = 'ENVIADA')`, unidade).Scan(&ok); err != nil {
		return false, fmt.Errorf("acervo.EmGuiaEnviada: %w", err)
	}
	return ok, nil
}

// ListGuias pagina as guias da origem no escopo (mais recentes primeiro).
func (r *Repository) ListGuias(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Guia, int64, error) {
	conds, args := "true", []any{}
	if f.Restrito {
		args = append(args, f.Unidades)
		conds += " AND g.origem_id = ANY($1)"
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_guias g WHERE `+conds, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListGuias count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT `+guiaCols+` FROM acervo_guias g WHERE `+conds+
		` ORDER BY g.criada_em DESC, g.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListGuias: %w", err)
	}
	defer rows.Close()
	out := []domain.Guia{}
	for rows.Next() {
		g, err := scanGuia(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("acervo.ListGuias scan: %w", err)
		}
		out = append(out, g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListGuias rows: %w", err)
	}
	return out, total, nil
}

// ResolverUnidade acha a unidade produtora ativa (e de entidade ativa) pelo
// id ou pela sigla; sigla repetida em duas entidades é ambígua.
func (r *Repository) ResolverUnidade(ctx context.Context, db database.DBTX, ref string) (uuid.UUID, error) {
	rows, err := db.Query(ctx, `SELECT un.id FROM unidades un JOIN entidades e ON e.id = un.entidade_id
		WHERE un.ativo AND e.ativo AND (un.id::text = $1 OR (un.sigla <> '' AND upper(un.sigla) = upper($1))) LIMIT 2`, ref)
	if err != nil {
		return uuid.Nil, fmt.Errorf("acervo.ResolverUnidade: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return uuid.Nil, fmt.Errorf("acervo.ResolverUnidade scan: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return uuid.Nil, fmt.Errorf("acervo.ResolverUnidade rows: %w", err)
	}
	if len(ids) != 1 {
		return uuid.Nil, domain.ErrUnidadeProdutora
	}
	return ids[0], nil
}
