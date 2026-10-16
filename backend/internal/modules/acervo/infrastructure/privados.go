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

const privadoCols = `id, numero, titular, tipo_titular, documento, contato, descricao, datas_limite, mensuracao, local_guarda,
	servico_publico, situacao, parecer_comissao, parecer_texto, parecer_favoravel, parecer_em, decreto_numero, decreto_em,
	decreto_diario, acesso_autorizado, acesso_condicoes, guarda, guarda_em, guarda_termo, convenio, criado_em`

func scanPrivado(row pgx.Row) (domain.Privado, error) {
	var p domain.Privado
	err := row.Scan(&p.ID, &p.Numero, &p.Titular, &p.TipoTitular, &p.Documento, &p.Contato, &p.Descricao, &p.DatasLimite,
		&p.Mensuracao, &p.LocalGuarda, &p.ServicoPublico, &p.Situacao, &p.ParecerComissao, &p.ParecerTexto, &p.ParecerFavoravel,
		&p.ParecerEm, &p.DecretoNumero, &p.DecretoEm, &p.DecretoDiario, &p.AcessoAutorizado, &p.AcessoCondicoes, &p.Guarda,
		&p.GuardaEm, &p.GuardaTermo, &p.Convenio, &p.CriadoEm)
	return p, err
}

// InsertPrivado grava o cadastro e devolve o número (AP-AAAA-NNNNN).
func (r *Repository) InsertPrivado(ctx context.Context, db database.DBTX, p domain.Privado, por *uuid.UUID) (string, error) {
	var numero string
	if err := db.QueryRow(ctx, `INSERT INTO acervo_privados (id, titular, tipo_titular, documento, contato, descricao, datas_limite,
		mensuracao, local_guarda, servico_publico, situacao, acesso_autorizado, acesso_condicoes, guarda, convenio, criado_por, criado_em)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING numero`,
		p.ID, p.Titular, p.TipoTitular, p.Documento, p.Contato, p.Descricao, p.DatasLimite, p.Mensuracao, p.LocalGuarda,
		p.ServicoPublico, p.Situacao, p.AcessoAutorizado, p.AcessoCondicoes, p.Guarda, p.Convenio, por, p.CriadoEm).Scan(&numero); err != nil {
		return "", fmt.Errorf("acervo.InsertPrivado: %w", err)
	}
	return numero, nil
}

// GetPrivado busca o arquivo com o histórico e as alienações.
func (r *Repository) GetPrivado(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (domain.Privado, error) {
	q := `SELECT ` + privadoCols + ` FROM acervo_privados WHERE id = $1`
	if paraAtualizar {
		q += ` FOR UPDATE`
	}
	p, err := scanPrivado(db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return p, domain.ErrPrivadoNaoExiste
	}
	if err != nil {
		return p, fmt.Errorf("acervo.GetPrivado: %w", err)
	}
	rows, err := db.Query(ctx, `SELECT tipo, descricao, por, em FROM acervo_privado_eventos WHERE privado_id = $1 ORDER BY em DESC, id`, id)
	if err != nil {
		return p, fmt.Errorf("acervo.GetPrivado eventos: %w", err)
	}
	p.Eventos = []domain.Evento{}
	for rows.Next() {
		var e domain.Evento
		if err := rows.Scan(&e.Tipo, &e.Descricao, &e.Por, &e.Em); err != nil {
			rows.Close()
			return p, fmt.Errorf("acervo.GetPrivado eventos scan: %w", err)
		}
		p.Eventos = append(p.Eventos, e)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("acervo.GetPrivado eventos rows: %w", err)
	}
	arows, err := db.Query(ctx, `SELECT id, notificado_em, adquirente, condicoes, decisao, fundamento, decidido_em
		FROM acervo_privado_alienacoes WHERE privado_id = $1 ORDER BY notificado_em DESC, id`, id)
	if err != nil {
		return p, fmt.Errorf("acervo.GetPrivado alienações: %w", err)
	}
	defer arows.Close()
	p.Alienacoes = []domain.Alienacao{}
	for arows.Next() {
		var a domain.Alienacao
		if err := arows.Scan(&a.ID, &a.NotificadoEm, &a.Adquirente, &a.Condicoes, &a.Decisao, &a.Fundamento, &a.DecididoEm); err != nil {
			return p, fmt.Errorf("acervo.GetPrivado alienações scan: %w", err)
		}
		p.Alienacoes = append(p.Alienacoes, a)
	}
	if err := arows.Err(); err != nil {
		return p, fmt.Errorf("acervo.GetPrivado alienações rows: %w", err)
	}
	return p, nil
}

// ListPrivados pagina os arquivos privados (mais recentes primeiro), com
// busca pelo titular, pelo número ou pela descrição.
func (r *Repository) ListPrivados(ctx context.Context, db database.DBTX, q string, p pagination.Params) ([]domain.Privado, int64, error) {
	conds, args := "true", []any{}
	if q != "" {
		args = append(args, "%"+q+"%")
		conds = "(titular ILIKE $1 OR numero ILIKE $1 OR descricao ILIKE $1)"
	}
	var total int64
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM acervo_privados WHERE `+conds, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListPrivados count: %w", err)
	}
	args = append(args, p.Limit(), p.Offset())
	rows, err := db.Query(ctx, `SELECT `+privadoCols+` FROM acervo_privados WHERE `+conds+
		` ORDER BY criado_em DESC, id LIMIT $`+strconv.Itoa(len(args)-1)+` OFFSET $`+strconv.Itoa(len(args)), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("acervo.ListPrivados: %w", err)
	}
	defer rows.Close()
	out := []domain.Privado{}
	for rows.Next() {
		x, err := scanPrivado(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("acervo.ListPrivados scan: %w", err)
		}
		x.Eventos, x.Alienacoes = []domain.Evento{}, []domain.Alienacao{}
		out = append(out, x)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("acervo.ListPrivados rows: %w", err)
	}
	return out, total, nil
}

// UpdatePrivado grava o cadastro e as etapas.
func (r *Repository) UpdatePrivado(ctx context.Context, db database.DBTX, p domain.Privado) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_privados SET titular = $2, documento = $3, contato = $4, descricao = $5,
		datas_limite = $6, mensuracao = $7, local_guarda = $8, situacao = $9, parecer_comissao = $10, parecer_texto = $11,
		parecer_favoravel = $12, parecer_em = $13, decreto_numero = $14, decreto_em = $15, decreto_diario = $16,
		acesso_autorizado = $17, acesso_condicoes = $18, guarda = $19, guarda_em = $20, guarda_termo = $21, convenio = $22
		WHERE id = $1`,
		p.ID, p.Titular, p.Documento, p.Contato, p.Descricao, p.DatasLimite, p.Mensuracao, p.LocalGuarda, p.Situacao,
		p.ParecerComissao, p.ParecerTexto, p.ParecerFavoravel, p.ParecerEm, p.DecretoNumero, p.DecretoEm, p.DecretoDiario,
		p.AcessoAutorizado, p.AcessoCondicoes, p.Guarda, p.GuardaEm, p.GuardaTermo, p.Convenio); err != nil {
		return fmt.Errorf("acervo.UpdatePrivado: %w", err)
	}
	return nil
}

// InsertEventoPrivado grava um passo no histórico.
func (r *Repository) InsertEventoPrivado(ctx context.Context, db database.DBTX, privado uuid.UUID, e domain.Evento) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_privado_eventos (privado_id, tipo, descricao, por) VALUES ($1, $2, $3, $4)`,
		privado, e.Tipo, e.Descricao, e.Por); err != nil {
		return fmt.Errorf("acervo.InsertEventoPrivado: %w", err)
	}
	return nil
}

// InsertAlienacao grava a notificação de alienação (art. 58).
func (r *Repository) InsertAlienacao(ctx context.Context, db database.DBTX, privado uuid.UUID, a domain.Alienacao, por *uuid.UUID) error {
	if _, err := db.Exec(ctx, `INSERT INTO acervo_privado_alienacoes (id, privado_id, notificado_em, adquirente, condicoes, criada_por)
		VALUES ($1, $2, $3, $4, $5, $6)`, a.ID, privado, a.NotificadoEm, a.Adquirente, a.Condicoes, por); err != nil {
		return fmt.Errorf("acervo.InsertAlienacao: %w", err)
	}
	return nil
}

// DecidirAlienacao grava a manifestação do município.
func (r *Repository) DecidirAlienacao(ctx context.Context, db database.DBTX, a domain.Alienacao, por *uuid.UUID) error {
	if _, err := db.Exec(ctx, `UPDATE acervo_privado_alienacoes SET decisao = $2, fundamento = $3, decidido_em = $4, decidido_por = $5
		WHERE id = $1`, a.ID, a.Decisao, a.Fundamento, a.DecididoEm, por); err != nil {
		return fmt.Errorf("acervo.DecidirAlienacao: %w", err)
	}
	return nil
}
