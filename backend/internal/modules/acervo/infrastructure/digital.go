package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

const digitalCols = `d.id, d.numero, d.unidade_id, d.descricao, d.tipo_documento, d.equipamento, d.dpi_declarado, d.responsavel_id,
	d.responsavel_nome, d.responsavel_cpf, d.codigo_classificacao, d.digitalizado_em, d.situacao, d.motivo, d.tentativas,
	d.original_key, d.original_sha256, d.original_bytes, d.pdfa_key, d.pdfa_sha256, d.pdfa_bytes, d.paginas, d.dpi_medido,
	d.conformidade, d.retencao_ate, d.processado_em, d.verificado_em, d.integro, d.criado_em, d.publicado_em`

func scanDigital(row pgx.Row) (domain.Digital, error) {
	var d domain.Digital
	err := row.Scan(&d.ID, &d.Numero, &d.UnidadeID, &d.Descricao, &d.TipoDoc, &d.Equipamento, &d.DPIDeclarado, &d.ResponsavelID,
		&d.ResponsavelNome, &d.ResponsavelCPF, &d.CodigoClass, &d.DigitalizadoEm, &d.Situacao, &d.Motivo, &d.Tentativas,
		&d.OriginalKey, &d.OriginalSHA256, &d.OriginalBytes, &d.PDFAKey, &d.PDFASHA256, &d.PDFABytes, &d.Paginas, &d.DPIMedido,
		&d.Conformidade, &d.RetencaoAte, &d.ProcessadoEm, &d.VerificadoEm, &d.Integro, &d.CriadoEm, &d.PublicadoEm)
	return d, err
}

// InsertDigital grava o envio e devolve o identificador (RD-AAAA-NNNNNN).
func (r *Repository) InsertDigital(ctx context.Context, db database.DBTX, d domain.Digital) (string, error) {
	var numero string
	if err := db.QueryRow(ctx, `INSERT INTO acervo_digitais (id, unidade_id, descricao, tipo_documento, equipamento, dpi_declarado,
		responsavel_id, responsavel_nome, responsavel_cpf, codigo_classificacao, digitalizado_em, situacao, original_key,
		original_sha256, original_bytes, criado_por, criado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $7, $16) RETURNING numero`,
		d.ID, d.UnidadeID, d.Descricao, d.TipoDoc, d.Equipamento, d.DPIDeclarado, d.ResponsavelID, d.ResponsavelNome,
		d.ResponsavelCPF, d.CodigoClass, d.DigitalizadoEm, d.Situacao, d.OriginalKey, d.OriginalSHA256, d.OriginalBytes,
		d.CriadoEm).Scan(&numero); err != nil {
		return "", fmt.Errorf("acervo.InsertDigital: %w", err)
	}
	return numero, nil
}

// GetDigital busca o representante.
func (r *Repository) GetDigital(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Digital, error) {
	d, err := scanDigital(db.QueryRow(ctx, `SELECT `+digitalCols+` FROM acervo_digitais d WHERE d.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return d, domain.ErrDigitalNaoExiste
	}
	if err != nil {
		return d, fmt.Errorf("acervo.GetDigital: %w", err)
	}
	return d, nil
}

// ListDigitais pagina os representantes (mais recentes primeiro).
func (r *Repository) ListDigitais(ctx context.Context, db database.DBTX, f domain.FiltroDigital, p pagination.Params) ([]domain.Digital, int64, error) {
	conds, args := "true", []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		conds += " AND " + fmt.Sprintf(cond, len(args))
	}
	if f.UnidadeID != nil {
		add("d.unidade_id = $%d", *f.UnidadeID)
	}
	if f.Situacao != "" {
		add("d.situacao = $%d", f.Situacao)
	}
	if f.SoFalhas {
		conds += " AND (d.integro = false OR d.situacao IN ('FALHOU', 'REPROVADO'))"
	}
	if f.Publicados {
		conds += " AND d.publicado_em IS NOT NULL"
	}

	if f.Restrito {
		add("u.unidade_id = ANY($%d)", f.Unidades)
	}
	from := ` FROM acervo_digitais d JOIN acervo_unidades u ON u.id = d.unidade_id WHERE ` + conds
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*)`+from, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListDigitais count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT `+digitalCols+from+` ORDER BY d.criado_em DESC, d.id LIMIT $`+strconv.Itoa(len(args)-1)+
		` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListDigitais: %w", err)
	}
	out, err := lerDigitais(rows)
	return out, total, err
}

func lerDigitais(rows pgx.Rows) ([]domain.Digital, error) {
	defer rows.Close()
	out := []domain.Digital{}
	for rows.Next() {
		d, err := scanDigital(rows)
		if err != nil {
			return nil, fmt.Errorf("acervo.digitais scan: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.digitais rows: %w", err)
	}
	return out, nil
}

// ProximoDigital reserva o próximo da fila: marca PROCESSANDO e conta a
// tentativa. SKIP LOCKED deixa vários workers dividirem a fila.
func (r *Repository) ProximoDigital(ctx context.Context, db database.DBTX, travadoAntes time.Time) (domain.Digital, bool, error) {
	d, err := scanDigital(db.QueryRow(ctx, `UPDATE acervo_digitais d SET situacao = 'PROCESSANDO', tentativas = tentativas + 1,
		processando_desde = now()
		WHERE d.id = (SELECT x.id FROM acervo_digitais x WHERE x.situacao = 'ENVIADO'
			OR (x.situacao = 'PROCESSANDO' AND x.processando_desde < $1)
			ORDER BY x.criado_em LIMIT 1 FOR UPDATE SKIP LOCKED)
		RETURNING `+digitalCols, travadoAntes))
	if errors.Is(err, pgx.ErrNoRows) {
		return d, false, nil
	}
	if err != nil {
		return d, false, fmt.Errorf("acervo.ProximoDigital: %w", err)
	}
	return d, true, nil
}

// ConcluirDigital grava o resultado do pipeline e o texto do OCR.
func (r *Repository) ConcluirDigital(ctx context.Context, db database.DBTX, d domain.Digital, texto string) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_digitais SET situacao = $2, motivo = $3, pdfa_key = $4, pdfa_sha256 = $5, pdfa_bytes = $6,
		paginas = $7, dpi_medido = $8, conformidade = $9, retencao_ate = $10, processado_em = $11, processando_desde = NULL,
		tentativas = $12 WHERE id = $1`,
		d.ID, d.Situacao, d.Motivo, d.PDFAKey, d.PDFASHA256, d.PDFABytes, d.Paginas, d.DPIMedido, d.Conformidade, d.RetencaoAte,
		d.ProcessadoEm, d.Tentativas); err != nil {
		return fmt.Errorf("acervo.ConcluirDigital: %w", err)
	}
	if texto == "" {
		return nil
	}
	if _, err := db.Exec(ctx, `INSERT INTO acervo_digital_textos (digital_id, texto) VALUES ($1, $2)
		ON CONFLICT (digital_id) DO UPDATE SET texto = EXCLUDED.texto`, d.ID, texto); err != nil {
		return fmt.Errorf("acervo.ConcluirDigital texto: %w", err)
	}
	return nil
}

// ParaVerificar: os prontos verificados antes de antesDe (ou nunca), os
// mais antigos primeiro.
func (r *Repository) ParaVerificar(ctx context.Context, db database.DBTX, antesDe time.Time, limite int) ([]domain.Digital, error) {
	rows, err := db.Query(ctx, `SELECT `+digitalCols+` FROM acervo_digitais d WHERE d.situacao = 'PRONTO'
		AND (d.verificado_em IS NULL OR d.verificado_em < $1) ORDER BY d.verificado_em NULLS FIRST, d.id LIMIT $2`, antesDe, limite)
	if err != nil {
		return nil, fmt.Errorf("acervo.ParaVerificar: %w", err)
	}
	return lerDigitais(rows)
}

// RegistrarVerificacao grava a verificação e o estado atual.
func (r *Repository) RegistrarVerificacao(ctx context.Context, db database.DBTX, id uuid.UUID, v domain.Verificacao) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_digital_verificacoes (digital_id, em, integro, sha256_lido, erro) VALUES ($1, $2, $3, $4, $5)`,
		id, v.Em, v.Integro, v.SHA256Lido, v.Erro); err != nil {
		return fmt.Errorf("acervo.RegistrarVerificacao: %w", err)
	}
	if _, err := db.Exec(ctx, `UPDATE acervo_digitais SET verificado_em = $2, integro = $3 WHERE id = $1`, id, v.Em, v.Integro); err != nil {
		return fmt.Errorf("acervo.RegistrarVerificacao estado: %w", err)
	}
	return nil
}

// ListVerificacoes devolve as últimas verificações.
func (r *Repository) ListVerificacoes(ctx context.Context, db database.DBTX, id uuid.UUID, limite int) ([]domain.Verificacao, error) {
	rows, err := db.Query(ctx, `SELECT em, integro, sha256_lido, erro FROM acervo_digital_verificacoes WHERE digital_id = $1
		ORDER BY em DESC LIMIT $2`, id, limite)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListVerificacoes: %w", err)
	}
	defer rows.Close()
	out := []domain.Verificacao{}
	for rows.Next() {
		var v domain.Verificacao
		if err := rows.Scan(&v.Em, &v.Integro, &v.SHA256Lido, &v.Erro); err != nil {
			return nil, fmt.Errorf("acervo.ListVerificacoes scan: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListVerificacoes rows: %w", err)
	}
	return out, nil
}

// EquipeCentral: usuários ativos com acervo:central em concessão global.
func (r *Repository) EquipeCentral(ctx context.Context, db database.DBTX) ([]uuid.UUID, error) {
	return r.equipeGlobal(ctx, db, "acervo:central", "EquipeCentral")
}

// EquipeSEMAD: usuários ativos com acervo:semad em concessão global.
func (r *Repository) EquipeSEMAD(ctx context.Context, db database.DBTX) ([]uuid.UUID, error) {
	return r.equipeGlobal(ctx, db, "acervo:semad", "EquipeSEMAD")
}

func (r *Repository) equipeGlobal(ctx context.Context, db database.DBTX, perm, nome string) ([]uuid.UUID, error) {
	rows, err := db.Query(ctx, `SELECT DISTINCT s.user_id FROM user_scopes s
		JOIN perfis p ON p.id = s.perfil_id AND p.ativo JOIN users us ON us.id = s.user_id AND us.active
		WHERE p.permissoes && ARRAY[$1, 'acervo:*']::text[]
		  AND s.entidade_id IS NULL AND s.unidade_id IS NULL AND s.departamento_id IS NULL ORDER BY s.user_id`, perm)
	if err != nil {
		return nil, fmt.Errorf("acervo.%s: %w", nome, err)
	}
	defer rows.Close()
	out := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("acervo.%s scan: %w", nome, err)
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.%s rows: %w", nome, err)
	}
	return out, nil
}
