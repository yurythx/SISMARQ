package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 2 (ADR 033): empréstimo e devolução, guia de transferência e
// importação do acervo existente.

var (
	ErrEmprestimoAberto    = errors.New("acervo: a unidade já está emprestada")
	ErrEmprestimoNaoExiste = errors.New("acervo: empréstimo não encontrado")
	ErrEmprestimoDevolvido = errors.New("acervo: empréstimo já devolvido")
	ErrEmprestada          = errors.New("acervo: a unidade está emprestada; registre a devolução antes")
	ErrGuiaNaoExiste       = errors.New("acervo: guia de transferência não encontrada")
	ErrGuiaSituacao        = errors.New("acervo: a guia não está mais aguardando o recebimento")
	ErrEmOutraGuia         = errors.New("acervo: a unidade já está numa guia de transferência aguardando o recebimento")
	ErrUnidadeProdutora    = errors.New("acervo: unidade produtora não encontrada (sigla ou id) ou desativada")
)

// Movimentações da Fase 2.
const (
	MovEmprestimo = "EMPRESTIMO"
	MovDevolucao  = "DEVOLUCAO"
	MovGuia       = "GUIA"
	MovImportacao = "IMPORTACAO"
	MovListagem   = "LISTAGEM"   // Fase 3: passos da relação de eliminação
	MovAmostra    = "AMOSTRA"    // art. 39
	MovEliminacao = "ELIMINACAO" // Termo de Eliminação (Anexo II)
)

// Emprestimo é a saída controlada de uma unidade (art. 12, II).
type Emprestimo struct {
	ID            uuid.UUID  `json:"id"`
	UnidadeID     uuid.UUID  `json:"unidade_arquivamento_id"`
	SolicitanteID uuid.UUID  `json:"solicitante_id"`
	Solicitante   string     `json:"solicitante"`
	Finalidade    string     `json:"finalidade"`
	DevolverAte   time.Time  `json:"devolver_ate"`
	RetiradoPor   *uuid.UUID `json:"retirado_por"`
	RetiradoEm    time.Time  `json:"retirado_em"`
	DevolvidoEm   *time.Time `json:"devolvido_em"`
	RecebidoPor   *uuid.UUID `json:"recebido_por"`
	Observacao    string     `json:"observacao"`
	// Numero e Assunto da unidade, na listagem de empréstimos.
	Numero  string `json:"numero,omitempty"`
	Assunto string `json:"assunto,omitempty"`
}

// Aberto informa se a unidade ainda não voltou.
func (e Emprestimo) Aberto() bool { return e.DevolvidoEm == nil }

// Atrasado informa se o prazo de devolução passou.
func (e Emprestimo) Atrasado(hoje time.Time) bool { return e.Aberto() && hoje.After(e.DevolverAte) }

// NovoEmprestimo valida o pedido. A devolução vem depois de hoje e em até
// um ano.
func NovoEmprestimo(unidade, solicitante uuid.UUID, finalidade string, devolverAte, hoje time.Time) (Emprestimo, error) {
	finalidade = strings.TrimSpace(finalidade)
	switch {
	case solicitante == uuid.Nil:
		return Emprestimo{}, invalido("solicitante_id", "informe quem retira a unidade")
	case finalidade == "" || utf8.RuneCountInString(finalidade) > 500:
		return Emprestimo{}, invalido("finalidade", "informe a finalidade (até 500 caracteres)")
	case devolverAte.Before(hoje) || devolverAte.After(hoje.AddDate(1, 0, 0)):
		return Emprestimo{}, invalido("devolver_ate", "entre hoje e um ano")
	}
	return Emprestimo{ID: uuid.New(), UnidadeID: unidade, SolicitanteID: solicitante, Finalidade: finalidade, DevolverAte: devolverAte}, nil
}

// Situação da guia de transferência.
const (
	GuiaEnviada   = "ENVIADA"
	GuiaRecebida  = "RECEBIDA"
	GuiaCancelada = "CANCELADA"
)

// MaxItensGuia limita uma guia (uma remessa física).
const MaxItensGuia = 500

// Guia é a guia de transferência do arquivo corrente ao intermediário
// (art. 29): a relação das unidades remetidas e o aceite de quem recebe. É
// ela mesma documento de arquivo (TTDD 2.0.03.01.00).
type Guia struct {
	ID          uuid.UUID      `json:"id"`
	Numero      string         `json:"numero"`
	OrigemID    uuid.UUID      `json:"origem_id"`
	Situacao    string         `json:"situacao"`
	Observacao  string         `json:"observacao"`
	CriadaPor   *uuid.UUID     `json:"criada_por"`
	CriadaEm    time.Time      `json:"criada_em"`
	RecebidaPor *uuid.UUID     `json:"recebida_por"`
	RecebidaEm  *time.Time     `json:"recebida_em"`
	Itens       []UnidadeLista `json:"itens"`
	TotalItens  int            `json:"total_itens"`
}

// NovaGuia valida a remessa: ao menos uma unidade, sem repetição, até o
// limite.
func NovaGuia(origem uuid.UUID, unidades []uuid.UUID, observacao string) (Guia, []uuid.UUID, error) {
	observacao = strings.TrimSpace(observacao)
	vistos := map[uuid.UUID]bool{}
	var ids []uuid.UUID
	for _, id := range unidades {
		if id != uuid.Nil && !vistos[id] {
			vistos[id] = true
			ids = append(ids, id)
		}
	}
	switch {
	case origem == uuid.Nil:
		return Guia{}, nil, invalido("origem_id", "informe a unidade de origem")
	case len(ids) == 0 || len(ids) > MaxItensGuia:
		return Guia{}, nil, invalido("unidades", "de 1 a 500 unidades de arquivamento por guia")
	case utf8.RuneCountInString(observacao) > 1000:
		return Guia{}, nil, invalido("observacao", "até 1000 caracteres")
	}
	return Guia{ID: uuid.New(), OrigemID: origem, Situacao: GuiaEnviada, Observacao: observacao}, ids, nil
}
