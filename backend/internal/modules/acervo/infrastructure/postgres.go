// Package infrastructure persiste o Acervo no PostgreSQL (consultas
// parametrizadas via pgx).
package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Repository implementa domain.Repository.
type Repository struct{}

// NewRepository cria o repositório.
func NewRepository() *Repository { return &Repository{} }

const unidadeCols = `u.id, u.numero, u.tipo, u.codigo_classificacao, u.assunto, u.unidade_id, u.data_inicial, u.data_final,
	u.quantidade, u.especificacao, u.observacoes, u.fase, u.destinacao, u.data_evento, u.contas_aprovadas_em,
	u.nivel_acesso, u.hipotese_legal, u.local_id, u.posicao, u.created_by, u.created_at, u.updated_at,
	u.fim_corrente, u.fim_intermediaria, u.situacao, u.recolhida_em, u.recolhimento_id`

func scanUnidade(row pgx.Row, extra ...any) (domain.Unidade, domain.Calculo, error) {
	var u domain.Unidade
	var c domain.Calculo
	dest := []any{&u.ID, &u.Numero, &u.Tipo, &u.CodigoClassificacao, &u.Assunto, &u.UnidadeID, &u.DataInicial, &u.DataFinal,
		&u.Quantidade, &u.Especificacao, &u.Observacoes, &u.Fase, &u.Destinacao, &u.DataEvento, &u.ContasAprovadasEm,
		&u.NivelAcesso, &u.HipoteseLegal, &u.LocalID, &u.Posicao, &u.CreatedBy, &u.CreatedAt, &u.UpdatedAt,
		&c.FimCorrente, &c.FimIntermediaria, &c.Situacao, &u.RecolhidaEm, &u.RecolhimentoID}
	err := row.Scan(append(dest, extra...)...)
	return u, c, err
}

// InsertUnidade grava a unidade e devolve o número gerado (UA-AAAA-NNNNNN).
func (r *Repository) InsertUnidade(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) (string, error) {
	var numero string
	err := db.QueryRow(ctx, `INSERT INTO acervo_unidades (id, tipo, codigo_classificacao, assunto, unidade_id, data_inicial,
		data_final, quantidade, especificacao, observacoes, fase, destinacao, data_evento, contas_aprovadas_em, nivel_acesso,
		hipotese_legal, local_id, posicao, created_by, fim_corrente, fim_intermediaria, situacao)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
		RETURNING numero`,
		u.ID, u.Tipo, u.CodigoClassificacao, u.Assunto, u.UnidadeID, u.DataInicial, u.DataFinal, u.Quantidade,
		u.Especificacao, u.Observacoes, u.Fase, u.Destinacao, u.DataEvento, u.ContasAprovadasEm, u.NivelAcesso,
		u.HipoteseLegal, u.LocalID, u.Posicao, u.CreatedBy, c.FimCorrente, c.FimIntermediaria, c.Situacao).Scan(&numero)
	if err != nil {
		return "", fmt.Errorf("acervo.InsertUnidade: %w", err)
	}
	return numero, nil
}

// GetUnidade busca a unidade (FOR UPDATE na transação que vai alterá-la).
func (r *Repository) GetUnidade(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Unidade, domain.Calculo, error) {
	q := `SELECT ` + unidadeCols + ` FROM acervo_unidades u WHERE u.id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	u, c, err := scanUnidade(db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return u, c, domain.ErrNaoEncontrada
	}
	if err != nil {
		return u, c, fmt.Errorf("acervo.GetUnidade: %w", err)
	}
	return u, c, nil
}

// UpdateUnidade grava os campos descritivos, a classificação e o cálculo
// (a fase muda só por SetFase).
func (r *Repository) UpdateUnidade(ctx context.Context, db database.DBTX, u domain.Unidade, c domain.Calculo) error {
	_, err := db.Exec(ctx, `UPDATE acervo_unidades SET tipo = $2, codigo_classificacao = $3, assunto = $4, unidade_id = $5,
		data_inicial = $6, data_final = $7, quantidade = $8, especificacao = $9, observacoes = $10, destinacao = $11,
		data_evento = $12, contas_aprovadas_em = $13, nivel_acesso = $14, hipotese_legal = $15, local_id = $16, posicao = $17,
		fim_corrente = $18, fim_intermediaria = $19, situacao = $20
		WHERE id = $1`,
		u.ID, u.Tipo, u.CodigoClassificacao, u.Assunto, u.UnidadeID, u.DataInicial, u.DataFinal, u.Quantidade,
		u.Especificacao, u.Observacoes, u.Destinacao, u.DataEvento, u.ContasAprovadasEm, u.NivelAcesso, u.HipoteseLegal,
		u.LocalID, u.Posicao, c.FimCorrente, c.FimIntermediaria, c.Situacao)
	if err != nil {
		return fmt.Errorf("acervo.UpdateUnidade: %w", err)
	}
	return nil
}

// SetCalculo grava só o resultado do cálculo de prazos.
func (r *Repository) SetCalculo(ctx context.Context, db database.DBTX, id uuid.UUID, c domain.Calculo) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_unidades SET fim_corrente = $2, fim_intermediaria = $3, situacao = $4 WHERE id = $1`,
		id, c.FimCorrente, c.FimIntermediaria, c.Situacao); err != nil {
		return fmt.Errorf("acervo.SetCalculo: %w", err)
	}
	return nil
}

// SetFase muda a fase (as travas do banco conferem a transição).
func (r *Repository) SetFase(ctx context.Context, db database.DBTX, id uuid.UUID, fase domain.Fase, c domain.Calculo) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_unidades SET fase = $2, fim_corrente = $3, fim_intermediaria = $4, situacao = $5 WHERE id = $1`,
		id, fase, c.FimCorrente, c.FimIntermediaria, c.Situacao); err != nil {
		return fmt.Errorf("acervo.SetFase: %w", err)
	}
	return nil
}

func filtroSQL(f domain.Filtro) (string, []any) {
	conds, args := []string{"true"}, []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		conds = append(conds, strings.ReplaceAll(cond, "?", "$"+strconv.Itoa(len(args))))
	}
	if f.Restrito {
		add("u.unidade_id = ANY(?)", f.Unidades)
	}
	if f.UnidadeID != nil {
		add("u.unidade_id = ?", *f.UnidadeID)
	}
	if f.Fase != "" {
		add("u.fase = ?", f.Fase)
	}
	if f.Situacao != "" {
		add("u.situacao = ?", f.Situacao)
	}
	if f.Codigo != "" {
		add("starts_with(u.codigo_classificacao, ?)", f.Codigo)
	}
	if q := strings.TrimSpace(f.Query); q != "" {
		add("u.search @@ nexus_search_tsquery('portuguese', ?)", q)
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

// caminhoLocal monta "Prédio > Depósito > Estante" de cada local.
const caminhoLocal = `WITH RECURSIVE arvore AS (
	SELECT id, parent_id, nome::text AS caminho FROM acervo_locais WHERE parent_id IS NULL
	UNION ALL
	SELECT l.id, l.parent_id, a.caminho || ' > ' || l.nome FROM acervo_locais l JOIN arvore a ON l.parent_id = a.id
)`

// ListUnidades pagina as unidades (mais recentes primeiro). Nunca nil.
func (r *Repository) ListUnidades(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.UnidadeLista, int64, error) {
	where, args := filtroSQL(f)
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_unidades u`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListUnidades count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, caminhoLocal+` SELECT `+unidadeCols+`, COALESCE(a.caminho, '') FROM acervo_unidades u
		LEFT JOIN arvore a ON a.id = u.local_id`+where+
		` ORDER BY u.created_at DESC, u.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListUnidades: %w", err)
	}
	defer rows.Close()
	out := []domain.UnidadeLista{}
	for rows.Next() {
		var l domain.UnidadeLista
		if l.Unidade, l.Calculo, err = scanUnidade(rows, &l.LocalCaminho); err != nil {
			return nil, 0, fmt.Errorf("acervo.ListUnidades scan: %w", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListUnidades rows: %w", err)
	}
	return out, total, nil
}

// Painel conta as unidades por situação.
func (r *Repository) Painel(ctx context.Context, db database.DBTX, f domain.Filtro) (map[string]int64, error) {
	where, args := filtroSQL(f)
	rows, err := db.Query(ctx, `SELECT u.situacao, COUNT(*) FROM acervo_unidades u`+where+` GROUP BY u.situacao`, args...)
	if err != nil {
		return nil, fmt.Errorf("acervo.Painel: %w", err)
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var s string
		var n int64
		if err := rows.Scan(&s, &n); err != nil {
			return nil, fmt.Errorf("acervo.Painel scan: %w", err)
		}
		out[s] = n
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.Painel rows: %w", err)
	}
	return out, nil
}

// ParaRecalculo pagina (por id) as unidades fora das fases finais.
func (r *Repository) ParaRecalculo(ctx context.Context, db database.DBTX, depoisDe uuid.UUID, limite int) ([]domain.Unidade, error) {
	rows, err := db.Query(ctx, `SELECT `+unidadeCols+` FROM acervo_unidades u
		WHERE u.fase IN ('CORRENTE', 'INTERMEDIARIA') AND u.id > $1 ORDER BY u.id LIMIT $2`, depoisDe, limite)
	if err != nil {
		return nil, fmt.Errorf("acervo.ParaRecalculo: %w", err)
	}
	defer rows.Close()
	out := []domain.Unidade{}
	for rows.Next() {
		u, _, err := scanUnidade(rows)
		if err != nil {
			return nil, fmt.Errorf("acervo.ParaRecalculo scan: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ParaRecalculo rows: %w", err)
	}
	return out, nil
}

const bloqueioCols = `id, unidade_id, processo, motivo, criado_por, criado_em, liberado_por, liberado_em, motivo_liberacao`

func scanBloqueio(row pgx.Row) (domain.Bloqueio, error) {
	var b domain.Bloqueio
	err := row.Scan(&b.ID, &b.UnidadeID, &b.Processo, &b.Motivo, &b.CriadoPor, &b.CriadoEm, &b.LiberadoPor, &b.LiberadoEm, &b.MotivoLiberacao)
	return b, err
}

// InsertBloqueio registra o bloqueio judicial.
func (r *Repository) InsertBloqueio(ctx context.Context, db database.DBTX, b domain.Bloqueio) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_bloqueios (id, unidade_id, processo, motivo, criado_por, criado_em)
		VALUES ($1, $2, $3, $4, $5, $6)`, b.ID, b.UnidadeID, b.Processo, b.Motivo, b.CriadoPor, b.CriadoEm); err != nil {
		return fmt.Errorf("acervo.InsertBloqueio: %w", err)
	}
	return nil
}

// GetBloqueio busca um bloqueio (FOR UPDATE: só é lido para ser liberado).
func (r *Repository) GetBloqueio(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Bloqueio, error) {
	b, err := scanBloqueio(db.QueryRow(ctx, `SELECT `+bloqueioCols+` FROM acervo_bloqueios WHERE id = $1 FOR UPDATE`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return b, domain.ErrBloqueioNaoExiste
	}
	if err != nil {
		return b, fmt.Errorf("acervo.GetBloqueio: %w", err)
	}
	return b, nil
}

// LiberarBloqueio grava a liberação.
func (r *Repository) LiberarBloqueio(ctx context.Context, db database.DBTX, b domain.Bloqueio) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_bloqueios SET liberado_por = $2, liberado_em = $3, motivo_liberacao = $4 WHERE id = $1`,
		b.ID, b.LiberadoPor, b.LiberadoEm, b.MotivoLiberacao); err != nil {
		return fmt.Errorf("acervo.LiberarBloqueio: %w", err)
	}
	return nil
}

// ListBloqueios lista os bloqueios da unidade (ativos primeiro).
func (r *Repository) ListBloqueios(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Bloqueio, error) {
	rows, err := db.Query(ctx, `SELECT `+bloqueioCols+` FROM acervo_bloqueios WHERE unidade_id = $1
		ORDER BY liberado_em IS NOT NULL, criado_em DESC`, unidade)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListBloqueios: %w", err)
	}
	defer rows.Close()
	out := []domain.Bloqueio{}
	for rows.Next() {
		b, err := scanBloqueio(rows)
		if err != nil {
			return nil, fmt.Errorf("acervo.ListBloqueios scan: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListBloqueios rows: %w", err)
	}
	return out, nil
}

// TemBloqueioAtivo informa se a unidade está retida.
func (r *Repository) TemBloqueioAtivo(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM acervo_bloqueios WHERE unidade_id = $1 AND liberado_em IS NULL)`,
		unidade).Scan(&ok); err != nil {
		return false, fmt.Errorf("acervo.TemBloqueioAtivo: %w", err)
	}
	return ok, nil
}

// InsertMovimentacao registra um passo do histórico.
func (r *Repository) InsertMovimentacao(ctx context.Context, db database.DBTX, unidade uuid.UUID, m domain.Movimentacao) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_movimentacoes (unidade_id, tipo, fase_de, fase_para, descricao, por)
		VALUES ($1, $2, $3, $4, $5, $6)`, unidade, m.Tipo, m.FaseDe, m.FasePara, m.Descricao, m.Por); err != nil {
		return fmt.Errorf("acervo.InsertMovimentacao: %w", err)
	}
	return nil
}

// ListMovimentacoes devolve o histórico (mais recente primeiro).
func (r *Repository) ListMovimentacoes(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]domain.Movimentacao, error) {
	rows, err := db.Query(ctx, `SELECT id, tipo, fase_de, fase_para, descricao, por, em FROM acervo_movimentacoes
		WHERE unidade_id = $1 ORDER BY em DESC, id`, unidade)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListMovimentacoes: %w", err)
	}
	defer rows.Close()
	out := []domain.Movimentacao{}
	for rows.Next() {
		var m domain.Movimentacao
		if err := rows.Scan(&m.ID, &m.Tipo, &m.FaseDe, &m.FasePara, &m.Descricao, &m.Por, &m.Em); err != nil {
			return nil, fmt.Errorf("acervo.ListMovimentacoes scan: %w", err)
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListMovimentacoes rows: %w", err)
	}
	return out, nil
}

// InsertLocal grava um local.
func (r *Repository) InsertLocal(ctx context.Context, db database.DBTX, l domain.Local) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_locais (id, parent_id, tipo, nome, ativo) VALUES ($1, $2, $3, $4, $5)`,
		l.ID, l.ParentID, l.Tipo, l.Nome, l.Ativo); err != nil {
		return fmt.Errorf("acervo.InsertLocal: %w", err)
	}
	return nil
}

// UpdateLocal altera um local.
func (r *Repository) UpdateLocal(ctx context.Context, db database.DBTX, l domain.Local) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_locais SET parent_id = $2, tipo = $3, nome = $4, ativo = $5 WHERE id = $1`,
		l.ID, l.ParentID, l.Tipo, l.Nome, l.Ativo); err != nil {
		return fmt.Errorf("acervo.UpdateLocal: %w", err)
	}
	return nil
}

// GetLocal busca um local com o caminho.
func (r *Repository) GetLocal(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Local, error) {
	var l domain.Local
	err := db.QueryRow(ctx, caminhoLocal+` SELECT l.id, l.parent_id, l.tipo, l.nome, l.ativo, COALESCE(a.caminho, l.nome)
		FROM acervo_locais l LEFT JOIN arvore a ON a.id = l.id WHERE l.id = $1`, id).
		Scan(&l.ID, &l.ParentID, &l.Tipo, &l.Nome, &l.Ativo, &l.Caminho)
	if errors.Is(err, pgx.ErrNoRows) {
		return l, domain.ErrLocalNaoExiste
	}
	if err != nil {
		return l, fmt.Errorf("acervo.GetLocal: %w", err)
	}
	return l, nil
}

// ListLocais devolve todos os locais com o caminho, em ordem de caminho.
func (r *Repository) ListLocais(ctx context.Context, db database.DBTX) ([]domain.Local, error) {
	rows, err := db.Query(ctx, caminhoLocal+` SELECT l.id, l.parent_id, l.tipo, l.nome, l.ativo, COALESCE(a.caminho, l.nome)
		FROM acervo_locais l LEFT JOIN arvore a ON a.id = l.id ORDER BY 6`)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListLocais: %w", err)
	}
	defer rows.Close()
	out := []domain.Local{}
	for rows.Next() {
		var l domain.Local
		if err := rows.Scan(&l.ID, &l.ParentID, &l.Tipo, &l.Nome, &l.Ativo, &l.Caminho); err != nil {
			return nil, fmt.Errorf("acervo.ListLocais scan: %w", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListLocais rows: %w", err)
	}
	return out, nil
}
