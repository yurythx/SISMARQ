package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Fase 6 (ADR 038): consulta pública, índice da busca e exposições.

// ostensivoSQL: a unidade vai para a consulta pública.
const ostensivoSQL = `u.nivel_acesso = 'PUBLICO' AND u.fase <> 'ELIMINADO'`

// PublicarDigital libera (em != nil) ou retira o representante da consulta
// pública. O gatilho do banco recusa o que não é ostensivo.
func (r *Repository) PublicarDigital(ctx context.Context, db database.DBTX, id uuid.UUID, em *time.Time, por *uuid.UUID) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_digitais SET publicado_em = $2, publicado_por = $3 WHERE id = $1`, id, em, por); err != nil {
		return fmt.Errorf("acervo.PublicarDigital: %w", err)
	}
	return nil
}

// PendentesIndice lista as unidades marcadas para reindexar (mais antigas
// primeiro).
func (r *Repository) PendentesIndice(ctx context.Context, db database.DBTX, limite int) ([]domain.PendenteIndice, error) {
	rows, err := db.Query(ctx, `SELECT unidade_id, desde FROM acervo_indice_pendentes ORDER BY desde, unidade_id LIMIT $1`, limite)
	if err != nil {
		return nil, fmt.Errorf("acervo.PendentesIndice: %w", err)
	}
	defer rows.Close()
	out := []domain.PendenteIndice{}
	for rows.Next() {
		var p domain.PendenteIndice
		if err := rows.Scan(&p.UnidadeID, &p.Desde); err != nil {
			return nil, fmt.Errorf("acervo.PendentesIndice scan: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.PendentesIndice rows: %w", err)
	}
	return out, nil
}

// ConcluirPendentes tira a marca das unidades indexadas — só se não
// mudaram de novo depois da leitura (a marca nova fica para a próxima).
func (r *Repository) ConcluirPendentes(ctx context.Context, db database.DBTX, ps []domain.PendenteIndice) error {
	ids, desde := make([]uuid.UUID, len(ps)), make([]time.Time, len(ps))
	for i, p := range ps {
		ids[i], desde[i] = p.UnidadeID, p.Desde
	}
	if _, err := db.Exec(ctx, `DELETE FROM acervo_indice_pendentes p USING unnest($1::uuid[], $2::timestamptz[]) AS x(id, desde)
		WHERE p.unidade_id = x.id AND p.desde = x.desde`, ids, desde); err != nil {
		return fmt.Errorf("acervo.ConcluirPendentes: %w", err)
	}
	return nil
}

// RegistrosIndice carrega as unidades dadas, com o órgão e os
// representantes prontos.
func (r *Repository) RegistrosIndice(ctx context.Context, db database.DBTX, ids []uuid.UUID, comTexto bool) ([]domain.RegistroIndice, error) {
	return r.registros(ctx, db, `u.id = ANY($1) ORDER BY u.id`, comTexto, ids)
}

// PaginaIndice percorre o acervo por id (carga completa do índice).
func (r *Repository) PaginaIndice(ctx context.Context, db database.DBTX, depoisDe uuid.UUID, limite int) ([]domain.RegistroIndice, error) {
	return r.registros(ctx, db, `u.id > $1 ORDER BY u.id LIMIT $2`, true, depoisDe, limite)
}

func (r *Repository) registros(ctx context.Context, db database.DBTX, where string, comTexto bool, args ...any) ([]domain.RegistroIndice, error) {
	rows, err := db.Query(ctx, `SELECT `+unidadeCols+`, o.nome FROM acervo_unidades u JOIN unidades o ON o.id = u.unidade_id WHERE `+where, args...)
	if err != nil {
		return nil, fmt.Errorf("acervo.RegistrosIndice: %w", err)
	}
	out := []domain.RegistroIndice{}
	pos := map[uuid.UUID]int{}
	err = func() error {
		defer rows.Close()
		for rows.Next() {
			var reg domain.RegistroIndice
			if reg.Unidade, _, err = scanUnidade(rows, &reg.Orgao); err != nil {
				return fmt.Errorf("acervo.RegistrosIndice scan: %w", err)
			}
			pos[reg.Unidade.ID] = len(out)
			out = append(out, reg)
		}
		return rows.Err()
	}()
	if err != nil || len(out) == 0 {
		return out, err
	}
	ids := make([]uuid.UUID, len(out))
	for i, reg := range out {
		ids[i] = reg.Unidade.ID
	}
	texto, junta := `''`, ``
	if comTexto {
		texto, junta = `COALESCE(t.texto, '')`, ` LEFT JOIN acervo_digital_textos t ON t.digital_id = d.id`
	}
	drows, err := db.Query(ctx, `SELECT d.unidade_id, d.id, d.numero, d.descricao, d.paginas, d.pdfa_bytes, d.pdfa_sha256, d.conformidade,
		d.digitalizado_em, d.publicado_em, `+texto+` FROM acervo_digitais d`+junta+`
		WHERE d.unidade_id = ANY($1) AND d.situacao = 'PRONTO' ORDER BY d.digitalizado_em, d.numero`, ids)
	if err != nil {
		return nil, fmt.Errorf("acervo.RegistrosIndice digitais: %w", err)
	}
	defer drows.Close()
	for drows.Next() {
		var unidade uuid.UUID
		var d domain.DigitalIndice
		if err := drows.Scan(&unidade, &d.ID, &d.Numero, &d.Descricao, &d.Paginas, &d.Bytes, &d.SHA256, &d.Conformidade,
			&d.DigitalizadoEm, &d.PublicadoEm, &d.Texto); err != nil {
			return nil, fmt.Errorf("acervo.RegistrosIndice digitais scan: %w", err)
		}
		i := pos[unidade]
		out[i].Digitais = append(out[i].Digitais, d)
	}
	if err := drows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.RegistrosIndice digitais rows: %w", err)
	}
	return out, nil
}

// BuscarBanco é a busca no Postgres (o índice fora do ar): os mesmos
// filtros e o mesmo recorte de acesso, sem facetas.
func (r *Repository) BuscarBanco(ctx context.Context, db database.DBTX, c domain.Consulta, f domain.FiltroBusca) ([]domain.ItemBusca, int64, error) {
	conds, args := []string{"true"}, []any{}
	add := func(cond string, v any) {
		args = append(args, v)
		conds = append(conds, strings.ReplaceAll(cond, "?", "$"+strconv.Itoa(len(args))))
	}
	digital := `d.situacao = 'PRONTO'`
	if f.Publico {
		conds = append(conds, ostensivoSQL)
		digital += ` AND d.publicado_em IS NOT NULL`
	}
	if f.Restrito {
		add("u.unidade_id = ANY(?)", f.Unidades)
	}
	rank := `0::real`
	if c.Q != "" {
		add(`(u.search @@ nexus_search_tsquery('portuguese', ?) OR EXISTS (SELECT 1 FROM acervo_digitais d
			JOIN acervo_digital_textos t ON t.digital_id = d.id WHERE d.unidade_id = u.id AND `+digital+`
			AND t.search @@ nexus_search_tsquery('portuguese', ?)))`, c.Q)
		rank = `ts_rank(u.search, nexus_search_tsquery('portuguese', $` + strconv.Itoa(len(args)) + `))`
	}
	for _, k := range domain.FacetasInternas {
		v, ok := c.Filtros[k]
		if !ok {
			continue
		}
		switch k {
		case domain.FacetaOrgao:
			add("o.nome = ?", v)
		case domain.FacetaClasse:
			add("split_part(u.codigo_classificacao, '.', 1) = ?", v)
		case domain.FacetaEspecie:
			add("u.tipo = ?", v)
		case domain.FacetaFase:
			add("u.fase = ?", v)
		case domain.FacetaDestinacao:
			add("u.destinacao = ?", v)
		case domain.FacetaAno:
			ano, _ := strconv.Atoi(v)
			add("? BETWEEN EXTRACT(YEAR FROM u.data_inicial) AND EXTRACT(YEAR FROM u.data_final)", ano)
		case domain.FacetaNivel:
			add("u.nivel_acesso = ?", v)
		}
	}
	if c.SoDigital {
		conds = append(conds, `EXISTS (SELECT 1 FROM acervo_digitais d WHERE d.unidade_id = u.id AND `+digital+`)`)
	}
	from := ` FROM acervo_unidades u JOIN unidades o ON o.id = u.unidade_id WHERE ` + strings.Join(conds, " AND ")
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*)`+from, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.BuscarBanco count: %w", err)
	}
	p := pagination.New(c.Pagina, c.PorPagina, c.PorPagina)
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT u.id, u.numero, u.tipo, u.codigo_classificacao, u.assunto, o.nome,
		EXTRACT(YEAR FROM u.data_inicial)::int, EXTRACT(YEAR FROM u.data_final)::int, u.fase, u.destinacao, u.nivel_acesso,
		(SELECT COUNT(*) FROM acervo_digitais d WHERE d.unidade_id = u.id AND `+digital+`)::int`+from+
		` ORDER BY `+rank+` DESC, u.created_at DESC, u.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.BuscarBanco: %w", err)
	}
	defer rows.Close()
	out := []domain.ItemBusca{}
	for rows.Next() {
		var it domain.ItemBusca
		if err := rows.Scan(&it.ID, &it.Numero, &it.Tipo, &it.Codigo, &it.Assunto, &it.Orgao, &it.AnoInicial, &it.AnoFinal,
			&it.Fase, &it.Destinacao, &it.NivelAcesso, &it.Digitais); err != nil {
			return nil, 0, fmt.Errorf("acervo.BuscarBanco scan: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.BuscarBanco rows: %w", err)
	}
	return out, total, nil
}

// DigitaisRecentes lista os últimos liberados à consulta pública.
func (r *Repository) DigitaisRecentes(ctx context.Context, db database.DBTX, limite int) ([]domain.DigitalRecente, error) {
	rows, err := db.Query(ctx, `SELECT d.id, d.numero, d.descricao, d.paginas, d.pdfa_bytes, d.pdfa_sha256, d.conformidade,
		d.digitalizado_em, d.publicado_em, u.id, u.numero, u.assunto
		FROM acervo_digitais d JOIN acervo_unidades u ON u.id = d.unidade_id
		WHERE d.publicado_em IS NOT NULL AND `+ostensivoSQL+` ORDER BY d.publicado_em DESC, d.id LIMIT $1`, limite)
	if err != nil {
		return nil, fmt.Errorf("acervo.DigitaisRecentes: %w", err)
	}
	defer rows.Close()
	out := []domain.DigitalRecente{}
	for rows.Next() {
		var d domain.DigitalRecente
		if err := rows.Scan(&d.ID, &d.Numero, &d.Descricao, &d.Paginas, &d.Bytes, &d.SHA256, &d.Conformidade, &d.DigitalizadoEm,
			&d.PublicadoEm, &d.UnidadeID, &d.UnidadeNumero, &d.Assunto); err != nil {
			return nil, fmt.Errorf("acervo.DigitaisRecentes scan: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.DigitaisRecentes rows: %w", err)
	}
	return out, nil
}

// ContarPublico conta as unidades ostensivas e os representantes liberados.
func (r *Repository) ContarPublico(ctx context.Context, db database.DBTX) (unidades, digitais int64, err error) {
	if err = db.QueryRow(ctx, `SELECT (SELECT COUNT(*) FROM acervo_unidades u WHERE `+ostensivoSQL+`),
		(SELECT COUNT(*) FROM acervo_digitais d JOIN acervo_unidades u ON u.id = d.unidade_id
		  WHERE d.publicado_em IS NOT NULL AND `+ostensivoSQL+`)`).Scan(&unidades, &digitais); err != nil {
		return 0, 0, fmt.Errorf("acervo.ContarPublico: %w", err)
	}
	return unidades, digitais, nil
}

// RegistrarConsulta soma uma consulta pública do dia.
func (r *Repository) RegistrarConsulta(ctx context.Context, db database.DBTX, dia time.Time, tipo string) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_consultas (dia, tipo, total) VALUES ($1, $2, 1)
		ON CONFLICT (dia, tipo) DO UPDATE SET total = acervo_consultas.total + 1`, dia, tipo); err != nil {
		return fmt.Errorf("acervo.RegistrarConsulta: %w", err)
	}
	return nil
}

const exposicaoCols = `x.id, x.slug, x.titulo, x.resumo, x.texto, x.publicada_em, x.criado_por, x.criado_em, x.atualizado_em`

func scanExposicao(row pgx.Row) (domain.Exposicao, error) {
	var x domain.Exposicao
	err := row.Scan(&x.ID, &x.Slug, &x.Titulo, &x.Resumo, &x.Texto, &x.PublicadaEm, &x.CriadoPor, &x.CriadoEm, &x.AtualizadoEm)
	return x, err
}

// InsertExposicao grava a exposição e os itens.
func (r *Repository) InsertExposicao(ctx context.Context, db database.DBTX, x domain.Exposicao) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_exposicoes (id, slug, titulo, resumo, texto, criado_por, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`, x.ID, x.Slug, x.Titulo, x.Resumo, x.Texto, x.CriadoPor, x.CriadoEm); err != nil {
		if database.IsUniqueViolation(err) {
			return domain.ErrSlugEmUso
		}
		return fmt.Errorf("acervo.InsertExposicao: %w", err)
	}
	return gravarItens(ctx, db, x)
}

// UpdateExposicao troca o texto e os itens (o endereço fica).
func (r *Repository) UpdateExposicao(ctx context.Context, db database.DBTX, x domain.Exposicao) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_exposicoes SET titulo = $2, resumo = $3, texto = $4, atualizado_em = $5 WHERE id = $1`,
		x.ID, x.Titulo, x.Resumo, x.Texto, x.AtualizadoEm); err != nil {
		return fmt.Errorf("acervo.UpdateExposicao: %w", err)
	}
	if _, err := db.Exec(ctx, `DELETE FROM acervo_exposicao_itens WHERE exposicao_id = $1`, x.ID); err != nil {
		return fmt.Errorf("acervo.UpdateExposicao itens: %w", err)
	}
	return gravarItens(ctx, db, x)
}

func gravarItens(ctx context.Context, db database.DBTX, x domain.Exposicao) error {
	ids, ordens, legendas := make([]uuid.UUID, len(x.Itens)), make([]int32, len(x.Itens)), make([]string, len(x.Itens))
	for i, it := range x.Itens {
		ids[i], ordens[i], legendas[i] = it.DigitalID, int32(it.Ordem), it.Legenda // #nosec G115 -- até ItensMax
	}
	if _, err := db.Exec(ctx, `INSERT INTO acervo_exposicao_itens (exposicao_id, digital_id, ordem, legenda)
		SELECT $1, x.d, x.o, x.l FROM unnest($2::uuid[], $3::int[], $4::text[]) AS x(d, o, l)`, x.ID, ids, ordens, legendas); err != nil {
		return fmt.Errorf("acervo.exposicao itens: %w", err)
	}
	return nil
}

// GetExposicao busca pelo id (ou pelo slug, com id nulo), com os itens.
func (r *Repository) GetExposicao(ctx context.Context, db database.DBTX, id uuid.UUID, slug string) (domain.Exposicao, error) {
	x, err := scanExposicao(db.QueryRow(ctx, `SELECT `+exposicaoCols+` FROM acervo_exposicoes x
		WHERE ($1 <> '00000000-0000-0000-0000-000000000000'::uuid AND x.id = $1) OR ($1 = '00000000-0000-0000-0000-000000000000'::uuid AND x.slug = $2)`,
		id, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return x, domain.ErrExposicaoNaoExiste
	}
	if err != nil {
		return x, fmt.Errorf("acervo.GetExposicao: %w", err)
	}
	rows, err := db.Query(ctx, `SELECT i.digital_id, i.ordem, i.legenda, d.numero, d.descricao, d.paginas, u.id, u.numero,
		(d.publicado_em IS NOT NULL AND `+ostensivoSQL+`)
		FROM acervo_exposicao_itens i JOIN acervo_digitais d ON d.id = i.digital_id JOIN acervo_unidades u ON u.id = d.unidade_id
		WHERE i.exposicao_id = $1 ORDER BY i.ordem`, x.ID)
	if err != nil {
		return x, fmt.Errorf("acervo.GetExposicao itens: %w", err)
	}
	defer rows.Close()
	x.Itens = []domain.ItemExposicao{}
	for rows.Next() {
		var it domain.ItemExposicao
		if err := rows.Scan(&it.DigitalID, &it.Ordem, &it.Legenda, &it.Numero, &it.Descricao, &it.Paginas, &it.UnidadeID,
			&it.UnidadeNumero, &it.Publicado); err != nil {
			return x, fmt.Errorf("acervo.GetExposicao itens scan: %w", err)
		}
		x.Itens = append(x.Itens, it)
	}
	if err := rows.Err(); err != nil {
		return x, fmt.Errorf("acervo.GetExposicao itens rows: %w", err)
	}
	return x, nil
}

// ListExposicoes pagina as exposições (as publicadas, na consulta pública).
func (r *Repository) ListExposicoes(ctx context.Context, db database.DBTX, soPublicadas bool, p pagination.Params) ([]domain.Exposicao, int64, error) {
	where, ordem := ``, `x.atualizado_em DESC`
	if soPublicadas {
		where, ordem = ` WHERE x.publicada_em IS NOT NULL`, `x.publicada_em DESC`
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_exposicoes x`+where).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListExposicoes count: %w", err)
	}
	rows, err := db.Query(ctx, `SELECT `+exposicaoCols+` FROM acervo_exposicoes x`+where+` ORDER BY `+ordem+`, x.id LIMIT $1 OFFSET $2`,
		p.Limit(), p.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListExposicoes: %w", err)
	}
	defer rows.Close()
	out := []domain.Exposicao{}
	for rows.Next() {
		x, err := scanExposicao(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("acervo.ListExposicoes scan: %w", err)
		}
		out = append(out, x)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListExposicoes rows: %w", err)
	}
	return out, total, nil
}

// PublicarExposicao publica (em != nil) ou recolhe a exposição.
func (r *Repository) PublicarExposicao(ctx context.Context, db database.DBTX, id uuid.UUID, em *time.Time) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_exposicoes SET publicada_em = $2 WHERE id = $1`, id, em); err != nil {
		return fmt.Errorf("acervo.PublicarExposicao: %w", err)
	}
	return nil
}
