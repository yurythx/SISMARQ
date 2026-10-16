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

const recolhimentoCols = `r.id, r.numero, r.origem_id, r.motivo, r.fundamento, r.classificados, r.avaliados, r.organizados,
	r.higienizados, r.acondicionados, r.historico, r.observacao, r.situacao, r.criado_por, r.criado_em, r.recebido_por,
	r.recebido_em, COALESCE(r.registro_entrada, ''), r.local_id, r.devolucao_motivo,
	(SELECT COUNT(*) FROM acervo_recolhimento_itens i WHERE i.recolhimento_id = r.id)`

func scanRecolhimento(row pgx.Row) (domain.Recolhimento, error) {
	var r domain.Recolhimento
	p := &r.Preparo
	err := row.Scan(&r.ID, &r.Numero, &r.OrigemID, &r.Motivo, &r.Fundamento, &p.Classificados, &p.Avaliados, &p.Organizados,
		&p.Higienizados, &p.Acondicionados, &r.Historico, &r.Observacao, &r.Situacao, &r.CriadoPor, &r.CriadoEm, &r.RecebidoPor,
		&r.RecebidoEm, &r.RegistroEntrada, &r.LocalID, &r.DevolucaoMotivo, &r.TotalItens)
	return r, err
}

// InsertRecolhimento grava o recolhimento e os itens; devolve o número
// (RC-AAAA-NNNNN).
func (r *Repository) InsertRecolhimento(ctx context.Context, db database.DBTX, rc domain.Recolhimento, unidades []uuid.UUID) (string, error) {
	var numero string
	p := rc.Preparo
	if err := db.QueryRow(ctx, `INSERT INTO acervo_recolhimentos (id, origem_id, motivo, fundamento, classificados, avaliados,
		organizados, higienizados, acondicionados, historico, observacao, situacao, criado_por, criado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING numero`,
		rc.ID, rc.OrigemID, rc.Motivo, rc.Fundamento, p.Classificados, p.Avaliados, p.Organizados, p.Higienizados, p.Acondicionados,
		rc.Historico, rc.Observacao, rc.Situacao, rc.CriadoPor, rc.CriadoEm).Scan(&numero); err != nil {
		return "", fmt.Errorf("acervo.InsertRecolhimento: %w", err)
	}
	if _, err := db.Exec(ctx, `INSERT INTO acervo_recolhimento_itens (recolhimento_id, unidade_id) SELECT $1, unnest($2::uuid[])`,
		rc.ID, unidades); err != nil {
		return "", fmt.Errorf("acervo.InsertRecolhimento itens: %w", err)
	}
	return numero, nil
}

// GetRecolhimento busca o recolhimento com as unidades (o instrumento
// descritivo).
func (r *Repository) GetRecolhimento(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Recolhimento, error) {
	q := `SELECT ` + recolhimentoCols + ` FROM acervo_recolhimentos r WHERE r.id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	rc, err := scanRecolhimento(db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return rc, domain.ErrRecolhimentoNaoExiste
	}
	if err != nil {
		return rc, fmt.Errorf("acervo.GetRecolhimento: %w", err)
	}
	rows, err := db.Query(ctx, caminhoLocal+` SELECT `+unidadeCols+`, COALESCE(a.caminho, '') FROM acervo_recolhimento_itens i
		JOIN acervo_unidades u ON u.id = i.unidade_id LEFT JOIN arvore a ON a.id = u.local_id
		WHERE i.recolhimento_id = $1 ORDER BY u.codigo_classificacao, u.data_inicial, u.numero`, id)
	if err != nil {
		return rc, fmt.Errorf("acervo.GetRecolhimento itens: %w", err)
	}
	defer rows.Close()
	rc.Itens = []domain.UnidadeLista{}
	for rows.Next() {
		var l domain.UnidadeLista
		if l.Unidade, l.Calculo, err = scanUnidade(rows, &l.LocalCaminho); err != nil {
			return rc, fmt.Errorf("acervo.GetRecolhimento itens scan: %w", err)
		}
		rc.Itens = append(rc.Itens, l)
	}
	if err := rows.Err(); err != nil {
		return rc, fmt.Errorf("acervo.GetRecolhimento itens rows: %w", err)
	}
	return rc, nil
}

// ListRecolhimentos pagina os recolhimentos da origem no escopo (mais
// recentes primeiro).
func (r *Repository) ListRecolhimentos(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Recolhimento, int64, error) {
	conds, args := "true", []any{}
	if f.Restrito {
		args = append(args, f.Unidades)
		conds += " AND r.origem_id = ANY($1)"
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_recolhimentos r WHERE `+conds, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListRecolhimentos count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT `+recolhimentoCols+` FROM acervo_recolhimentos r WHERE `+conds+
		` ORDER BY r.criado_em DESC, r.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListRecolhimentos: %w", err)
	}
	defer rows.Close()
	out := []domain.Recolhimento{}
	for rows.Next() {
		rc, err := scanRecolhimento(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("acervo.ListRecolhimentos scan: %w", err)
		}
		out = append(out, rc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListRecolhimentos rows: %w", err)
	}
	return out, total, nil
}

// ReceberRecolhimento grava a entrada no Arquivo Público e devolve o número
// do registro (RE-AAAA-NNNNN).
func (r *Repository) ReceberRecolhimento(ctx context.Context, db database.DBTX, rc domain.Recolhimento) (string, error) {
	var registro string
	if err := db.QueryRow(ctx, `UPDATE acervo_recolhimentos SET situacao = 'RECEBIDO', recebido_por = $2, recebido_em = $3,
		local_id = $4, registro_entrada = 'RE-' || to_char($3::timestamptz, 'YYYY') || '-' ||
		lpad(nextval('acervo_registro_entrada_seq')::text, 5, '0')
		WHERE id = $1 RETURNING registro_entrada`, rc.ID, rc.RecebidoPor, rc.RecebidoEm, rc.LocalID).Scan(&registro); err != nil {
		return "", fmt.Errorf("acervo.ReceberRecolhimento: %w", err)
	}
	return registro, nil
}

// SetRecolhimentoSituacao grava a devolução ou o cancelamento.
func (r *Repository) SetRecolhimentoSituacao(ctx context.Context, db database.DBTX, rc domain.Recolhimento) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_recolhimentos SET situacao = $2, devolucao_motivo = $3 WHERE id = $1`,
		rc.ID, rc.Situacao, rc.DevolucaoMotivo); err != nil {
		return fmt.Errorf("acervo.SetRecolhimentoSituacao: %w", err)
	}
	return nil
}

// EmRecolhimentoAberto informa se a unidade está num recolhimento enviado.
func (r *Repository) EmRecolhimentoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM acervo_recolhimento_itens i
		JOIN acervo_recolhimentos r ON r.id = i.recolhimento_id WHERE i.unidade_id = $1 AND r.situacao = 'ENVIADO')`,
		unidade).Scan(&ok); err != nil {
		return false, fmt.Errorf("acervo.EmRecolhimentoAberto: %w", err)
	}
	return ok, nil
}

// MarcarRecolhida grava a entrada da unidade no Arquivo Público: fase
// permanente, destinação, guarda e local.
func (r *Repository) MarcarRecolhida(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_unidades SET fase = $2, destinacao = $3, recolhida_em = $4, recolhimento_id = $5,
		local_id = $6, fim_corrente = $7, fim_intermediaria = $8, situacao = $9 WHERE id = $1`,
		u.ID, u.Fase, u.Destinacao, u.RecolhidaEm, u.RecolhimentoID, u.LocalID, c.FimCorrente, c.FimIntermediaria, c.Situacao); err != nil {
		return fmt.Errorf("acervo.MarcarRecolhida: %w", err)
	}
	return nil
}

// GestoresDaUnidade: os usuários ativos com acervo:manage por lotação na
// unidade, numa unidade acima dela ou na entidade inteira. A concessão
// global (administração da plataforma) fica de fora: o aviso é da origem.
func (r *Repository) GestoresDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.Query(ctx, `WITH RECURSIVE acima AS (
			SELECT id, parent_id FROM unidades WHERE id = $1
			UNION ALL
			SELECT u.id, u.parent_id FROM unidades u JOIN acima a ON u.id = a.parent_id
		)
		SELECT DISTINCT s.user_id FROM user_scopes s
		JOIN perfis p ON p.id = s.perfil_id AND p.ativo
		JOIN users us ON us.id = s.user_id AND us.active
		WHERE p.permissoes && ARRAY['acervo:manage', 'acervo:*', '*']::text[] AND s.departamento_id IS NULL
		  AND (s.unidade_id IN (SELECT id FROM acima)
		       OR (s.unidade_id IS NULL AND s.entidade_id = (SELECT entidade_id FROM unidades WHERE id = $1)))
		ORDER BY s.user_id`, unidade)
	if err != nil {
		return nil, fmt.Errorf("acervo.GestoresDaUnidade: %w", err)
	}
	defer rows.Close()
	out := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("acervo.GestoresDaUnidade scan: %w", err)
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.GestoresDaUnidade rows: %w", err)
	}
	return out, nil
}
