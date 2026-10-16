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

// SalvarComissao grava a comissão e troca a lista de membros.
func (r *Repository) SalvarComissao(ctx context.Context, db database.DBTX, c domain.Comissao, nova bool) error {
	q := `UPDATE acervo_comissoes SET tipo = $2, entidade_id = $3, nome = $4, portaria = $5, ativa = $6 WHERE id = $1`
	if nova {
		q = `INSERT INTO acervo_comissoes (id, tipo, entidade_id, nome, portaria, ativa) VALUES ($1, $2, $3, $4, $5, $6)`
	}
	if _, err := db.Exec(ctx, q, c.ID, c.Tipo, c.EntidadeID, c.Nome, c.Portaria, c.Ativa); err != nil {
		if unico(err) {
			return &domain.ErrValidacao{Campo: "ativa", Motivo: "já existe comissão ativa para este órgão"}
		}
		return fmt.Errorf("acervo.SalvarComissao: %w", err)
	}
	if _, err := db.Exec(ctx, `DELETE FROM acervo_comissao_membros WHERE comissao_id = $1`, c.ID); err != nil {
		return fmt.Errorf("acervo.SalvarComissao membros: %w", err)
	}
	for _, m := range c.Membros {
		if _, err := db.Exec(ctx, `INSERT INTO acervo_comissao_membros (comissao_id, user_id, papel, area) VALUES ($1, $2, $3, $4)`,
			c.ID, m.UserID, m.Papel, m.Area); err != nil {
			return fmt.Errorf("acervo.SalvarComissao membro: %w", err)
		}
	}
	return nil
}

const comissaoCols = `c.id, c.tipo, c.entidade_id, c.nome, c.portaria, c.ativa`

func (r *Repository) membros(ctx context.Context, db database.DBTX, c *domain.Comissao) error {
	rows, err := db.Query(ctx, `SELECT m.user_id, COALESCE(u.display_name, ''), m.papel, m.area FROM acervo_comissao_membros m
		LEFT JOIN users u ON u.id = m.user_id WHERE m.comissao_id = $1 ORDER BY m.papel = 'PRESIDENTE' DESC, u.display_name`, c.ID)
	if err != nil {
		return fmt.Errorf("acervo.membros: %w", err)
	}
	defer rows.Close()
	c.Membros = []domain.Membro{}
	for rows.Next() {
		var m domain.Membro
		if err := rows.Scan(&m.UserID, &m.Nome, &m.Papel, &m.Area); err != nil {
			return fmt.Errorf("acervo.membros scan: %w", err)
		}
		c.Membros = append(c.Membros, m)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("acervo.membros rows: %w", err)
	}
	return nil
}

func (r *Repository) umaComissao(ctx context.Context, db database.DBTX, q string, args ...any) (domain.Comissao, error) {
	var c domain.Comissao
	err := db.QueryRow(ctx, q, args...).Scan(&c.ID, &c.Tipo, &c.EntidadeID, &c.Nome, &c.Portaria, &c.Ativa)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, domain.ErrComissaoNaoExiste
	}
	if err != nil {
		return c, fmt.Errorf("acervo.GetComissao: %w", err)
	}
	return c, r.membros(ctx, db, &c)
}

// GetComissao busca a comissão com os membros.
func (r *Repository) GetComissao(ctx context.Context, db database.DBTX, id uuid.UUID) (domain.Comissao, error) {
	return r.umaComissao(ctx, db, `SELECT `+comissaoCols+` FROM acervo_comissoes c WHERE c.id = $1`, id)
}

// ComissaoDaUnidade acha a CSPAD ativa do órgão (entidade) da unidade.
func (r *Repository) ComissaoDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) (domain.Comissao, error) {
	c, err := r.umaComissao(ctx, db, `SELECT `+comissaoCols+` FROM acervo_comissoes c JOIN unidades un ON un.entidade_id = c.entidade_id
		WHERE un.id = $1 AND c.ativa AND c.tipo = 'SETORIAL'`, unidade)
	if errors.Is(err, domain.ErrComissaoNaoExiste) {
		return c, domain.ErrSemComissao
	}
	return c, err
}

// ComissaoCentral acha a CCPAD ativa.
func (r *Repository) ComissaoCentral(ctx context.Context, db database.DBTX) (domain.Comissao, error) {
	c, err := r.umaComissao(ctx, db, `SELECT `+comissaoCols+` FROM acervo_comissoes c WHERE c.ativa AND c.tipo = 'CENTRAL'`)
	if errors.Is(err, domain.ErrComissaoNaoExiste) {
		return c, domain.ErrSemComissao
	}
	return c, err
}

// ListComissoes devolve as comissões com os membros (ativas primeiro).
func (r *Repository) ListComissoes(ctx context.Context, db database.DBTX) ([]domain.Comissao, error) {
	rows, err := db.Query(ctx, `SELECT `+comissaoCols+` FROM acervo_comissoes c ORDER BY c.ativa DESC, c.tipo, c.nome`)
	if err != nil {
		return nil, fmt.Errorf("acervo.ListComissoes: %w", err)
	}
	out := []domain.Comissao{}
	for rows.Next() {
		var c domain.Comissao
		if err := rows.Scan(&c.ID, &c.Tipo, &c.EntidadeID, &c.Nome, &c.Portaria, &c.Ativa); err != nil {
			rows.Close()
			return nil, fmt.Errorf("acervo.ListComissoes scan: %w", err)
		}
		out = append(out, c)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("acervo.ListComissoes rows: %w", err)
	}
	for i := range out {
		if err := r.membros(ctx, db, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

const listagemCols = `l.id, l.numero, l.comissao_id, l.unidade_id, l.situacao, l.rito, l.responsavel_selecao, l.parecer_juridico,
	l.contas_exercicios, l.contas_aprovadas_em, l.ata_numero, l.ata_em, l.ata_local, l.ata_manifestacoes, l.ata_presentes,
	l.autorizada_por, l.autorizada_em, l.devolucao_motivo, l.ccpad_autorizada_por, l.ccpad_autorizada_em,
	l.semad_autorizada_por, l.semad_autorizada_em, l.ata_signum, l.termo_signum, l.edital_numero, l.edital_publicado_em,
	l.edital_diario, l.amostra_criterios, l.amostra_percentual::float8, l.amostra_em, l.termo_eliminado_em,
	l.termo_acompanhamento, l.termo_destinacao_papel, l.termo_registrado_por, l.termo_copia_em, l.termo_copia_arquivo,
	l.termo_copia_semad, l.criada_por, l.criada_em,
	(SELECT COUNT(*) FROM acervo_listagem_itens i WHERE i.listagem_id = l.id)`

func scanListagem(row pgx.Row) (domain.Listagem, error) {
	var l domain.Listagem
	err := row.Scan(&l.ID, &l.Numero, &l.ComissaoID, &l.UnidadeID, &l.Situacao, &l.Rito, &l.ResponsavelSelecao, &l.ParecerJuridico,
		&l.ContasExercicios, &l.ContasAprovadasEm, &l.AtaNumero, &l.AtaEm, &l.AtaLocal, &l.AtaManifestacoes, &l.AtaPresentes,
		&l.AutorizadaPor, &l.AutorizadaEm, &l.DevolucaoMotivo, &l.CcpadAutorizadaPor, &l.CcpadAutorizadaEm,
		&l.SemadAutorizadaPor, &l.SemadAutorizadaEm, &l.AtaSignum, &l.TermoSignum, &l.EditalNumero, &l.EditalPublicadoEm,
		&l.EditalDiario, &l.AmostraCriterios, &l.AmostraPercentual, &l.AmostraEm, &l.TermoEliminadoEm,
		&l.TermoAcompanhamento, &l.TermoDestinacaoPapel, &l.TermoRegistradoPor, &l.TermoCopiaEm, &l.TermoCopiaArquivo,
		&l.TermoCopiaSEMAD, &l.CriadaPor, &l.CriadaEm, &l.TotalItens)
	if l.AtaSignum == "" {
		l.AtaSignum = domain.SignumAusente
	}
	if l.TermoSignum == "" {
		l.TermoSignum = domain.SignumAusente
	}
	return l, err
}

// InsertListagem grava a relação; devolve o número (LE-AAAA-NNNNN).
func (r *Repository) InsertListagem(ctx context.Context, db database.DBTX, l domain.Listagem) (string, error) {
	var numero string
	if err := db.QueryRow(ctx, `INSERT INTO acervo_listagens (id, comissao_id, unidade_id, rito, responsavel_selecao, parecer_juridico,
		contas_exercicios, contas_aprovadas_em, criada_por, criada_em) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING numero`,
		l.ID, l.ComissaoID, l.UnidadeID, l.Rito, l.ResponsavelSelecao, l.ParecerJuridico, l.ContasExercicios, l.ContasAprovadasEm,
		l.CriadaPor, l.CriadaEm).Scan(&numero); err != nil {
		return "", fmt.Errorf("acervo.InsertListagem: %w", err)
	}
	return numero, nil
}

// UpdateListagem grava o cabeçalho e a etapa (as CHECKs do banco conferem a
// ordem do rito).
func (r *Repository) UpdateListagem(ctx context.Context, db database.DBTX, l domain.Listagem) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_listagens SET situacao = $2, rito = $3, responsavel_selecao = $4, parecer_juridico = $5,
		contas_exercicios = $6, contas_aprovadas_em = $7, ata_numero = $8, ata_em = $9, ata_local = $10, ata_manifestacoes = $11,
		ata_presentes = $12, autorizada_por = $13, autorizada_em = $14, devolucao_motivo = $15, ccpad_autorizada_por = $16,
		ccpad_autorizada_em = $17, semad_autorizada_por = $18, semad_autorizada_em = $19, ata_signum = $20, termo_signum = $21,
		edital_numero = $22, edital_publicado_em = $23, edital_diario = $24, amostra_criterios = $25, amostra_percentual = $26,
		amostra_em = $27, termo_eliminado_em = $28, termo_acompanhamento = $29, termo_destinacao_papel = $30,
		termo_registrado_por = $31, termo_copia_em = $32, termo_copia_arquivo = $33, termo_copia_semad = $34
		WHERE id = $1`,
		l.ID, l.Situacao, l.Rito, l.ResponsavelSelecao, l.ParecerJuridico, l.ContasExercicios, l.ContasAprovadasEm, l.AtaNumero, l.AtaEm,
		l.AtaLocal, l.AtaManifestacoes, l.AtaPresentes, l.AutorizadaPor, l.AutorizadaEm, l.DevolucaoMotivo, l.CcpadAutorizadaPor,
		l.CcpadAutorizadaEm, l.SemadAutorizadaPor, l.SemadAutorizadaEm, l.AtaSignum, l.TermoSignum, l.EditalNumero,
		l.EditalPublicadoEm, l.EditalDiario, l.AmostraCriterios, l.AmostraPercentual, l.AmostraEm, l.TermoEliminadoEm,
		l.TermoAcompanhamento, l.TermoDestinacaoPapel, l.TermoRegistradoPor, l.TermoCopiaEm, l.TermoCopiaArquivo, l.TermoCopiaSEMAD); err != nil {
		return fmt.Errorf("acervo.UpdateListagem: %w", err)
	}
	return nil
}

// GetListagem busca a relação com os itens e os pedidos.
func (r *Repository) GetListagem(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Listagem, error) {
	q := `SELECT ` + listagemCols + ` FROM acervo_listagens l WHERE l.id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	l, err := scanListagem(db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return l, domain.ErrListagemNaoExiste
	}
	if err != nil {
		return l, fmt.Errorf("acervo.GetListagem: %w", err)
	}
	rows, err := db.Query(ctx, caminhoLocal+` SELECT `+unidadeCols+`, COALESCE(a.caminho, ''), i.amostra, i.retirada
		FROM acervo_listagem_itens i JOIN acervo_unidades u ON u.id = i.unidade_id LEFT JOIN arvore a ON a.id = u.local_id
		WHERE i.listagem_id = $1 ORDER BY u.codigo_classificacao, u.data_inicial, u.numero`, id)
	if err != nil {
		return l, fmt.Errorf("acervo.GetListagem itens: %w", err)
	}
	l.Itens = []domain.ItemListagem{}
	for rows.Next() {
		var it domain.ItemListagem
		if it.Unidade, it.Calculo, err = scanUnidade(rows, &it.LocalCaminho, &it.Amostra, &it.Retirada); err != nil {
			rows.Close()
			return l, fmt.Errorf("acervo.GetListagem itens scan: %w", err)
		}
		l.Itens = append(l.Itens, it)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return l, fmt.Errorf("acervo.GetListagem itens rows: %w", err)
	}
	prows, err := db.Query(ctx, `SELECT id, unidade_id, requerente, pedido, recebido_em, situacao, decisao, decidido_por, decidido_em
		FROM acervo_listagem_pedidos WHERE listagem_id = $1 ORDER BY recebido_em, id`, id)
	if err != nil {
		return l, fmt.Errorf("acervo.GetListagem pedidos: %w", err)
	}
	defer prows.Close()
	l.Pedidos = []domain.Pedido{}
	for prows.Next() {
		var p domain.Pedido
		if err := prows.Scan(&p.ID, &p.UnidadeID, &p.Requerente, &p.Pedido, &p.RecebidoEm, &p.Situacao, &p.Decisao, &p.DecididoPor, &p.DecididoEm); err != nil {
			return l, fmt.Errorf("acervo.GetListagem pedidos scan: %w", err)
		}
		l.Pedidos = append(l.Pedidos, p)
	}
	if err := prows.Err(); err != nil {
		return l, fmt.Errorf("acervo.GetListagem pedidos rows: %w", err)
	}
	return l, nil
}

// ListListagens pagina as relações no escopo (mais recentes primeiro).
func (r *Repository) ListListagens(ctx context.Context, db database.DBTX, f domain.Filtro, p pagination.Params) ([]domain.Listagem, int64, error) {
	conds, args := "true", []any{}
	if f.Restrito {
		args = append(args, f.Unidades)
		conds += " AND l.unidade_id = ANY($1)"
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_listagens l WHERE `+conds, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListListagens count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT `+listagemCols+` FROM acervo_listagens l WHERE `+conds+
		` ORDER BY l.criada_em DESC, l.id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListListagens: %w", err)
	}
	defer rows.Close()
	out := []domain.Listagem{}
	for rows.Next() {
		l, err := scanListagem(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("acervo.ListListagens scan: %w", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListListagens rows: %w", err)
	}
	return out, total, nil
}

// AddItemListagem inclui a unidade na relação.
func (r *Repository) AddItemListagem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_listagem_itens (listagem_id, unidade_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		listagem, unidade); err != nil {
		return fmt.Errorf("acervo.AddItemListagem: %w", err)
	}
	return nil
}

// RemoveItemListagem tira a unidade da relação (rascunho).
func (r *Repository) RemoveItemListagem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID) error {
	if _, err := db.Exec(ctx, `DELETE FROM acervo_listagem_itens WHERE listagem_id = $1 AND unidade_id = $2`, listagem, unidade); err != nil {
		return fmt.Errorf("acervo.RemoveItemListagem: %w", err)
	}
	return nil
}

// MarcarItem grava se o item é amostra ou foi retirado por pedido deferido.
func (r *Repository) MarcarItem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID, amostra, retirada bool) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_listagem_itens SET amostra = $3, retirada = $4 WHERE listagem_id = $1 AND unidade_id = $2`,
		listagem, unidade, amostra, retirada); err != nil {
		return fmt.Errorf("acervo.MarcarItem: %w", err)
	}
	return nil
}

// EmListagemAtiva informa se a unidade está numa relação em andamento.
func (r *Repository) EmListagemAtiva(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error) {
	var ok bool
	if err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM acervo_listagem_itens i JOIN acervo_listagens l ON l.id = i.listagem_id
		WHERE i.unidade_id = $1 AND NOT i.retirada AND l.situacao NOT IN ('ELIMINADA', 'CANCELADA'))`, unidade).Scan(&ok); err != nil {
		return false, fmt.Errorf("acervo.EmListagemAtiva: %w", err)
	}
	return ok, nil
}

// InsertPedido registra o pedido de desentranhamento ou de cópia.
func (r *Repository) InsertPedido(ctx context.Context, db database.DBTX, listagem uuid.UUID, p domain.Pedido) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_listagem_pedidos (id, listagem_id, unidade_id, requerente, pedido, recebido_em)
		VALUES ($1, $2, $3, $4, $5, $6)`, p.ID, listagem, p.UnidadeID, p.Requerente, p.Pedido, p.RecebidoEm); err != nil {
		return fmt.Errorf("acervo.InsertPedido: %w", err)
	}
	return nil
}

// DecidirPedido grava a decisão.
func (r *Repository) DecidirPedido(ctx context.Context, db database.DBTX, p domain.Pedido) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_listagem_pedidos SET situacao = $2, decisao = $3, decidido_por = $4, decidido_em = $5 WHERE id = $1`,
		p.ID, p.Situacao, p.Decisao, p.DecididoPor, p.DecididoEm); err != nil {
		return fmt.Errorf("acervo.DecidirPedido: %w", err)
	}
	return nil
}
