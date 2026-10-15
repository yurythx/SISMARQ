// Package domain é o núcleo do Acervo (SISMARQ, Fase 1 do
// docs/ROADMAP_SISMARQ.md): a unidade de arquivamento, a localização física,
// o bloqueio judicial e as travas da Lei Municipal nº 10.301/2019. Cada regra
// cita o artigo da lei. O domínio não conhece banco nem HTTP.
package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase é a idade do documento (art. 18) ou o seu fim.
type Fase string

const (
	FaseCorrente      Fase = "CORRENTE"      // art. 18, § 1º
	FaseIntermediaria Fase = "INTERMEDIARIA" // art. 18, § 2º
	FasePermanente    Fase = "PERMANENTE"    // art. 18, § 3º
	FaseEliminado     Fase = "ELIMINADO"     // art. 38 (Termo de Eliminação)
)

// Destinacao é a destinação final que a TTDD dá à série (art. 28).
type Destinacao string

const (
	DestinacaoEliminacao       Destinacao = "ELIMINACAO"
	DestinacaoGuardaPermanente Destinacao = "GUARDA_PERMANENTE"
	// DestinacaoIndefinida: a série não traz destinação. Nada sai do
	// arquivo intermediário até a Comissão Central defini-la (art. 52).
	DestinacaoIndefinida Destinacao = "INDEFINIDA"
)

// Tipo é a espécie da unidade de arquivamento (a coluna "Unidade de
// arquivamento" dos Anexos I a III).
type Tipo string

const (
	TipoCaixa     Tipo = "CAIXA"
	TipoPasta     Tipo = "PASTA"
	TipoLivro     Tipo = "LIVRO"
	TipoProcesso  Tipo = "PROCESSO"
	TipoDossie    Tipo = "DOSSIE"
	TipoDocumento Tipo = "DOCUMENTO"
	TipoMapa      Tipo = "MAPA"
	TipoOutro     Tipo = "OUTRO"
)

// NivelAcesso: o acesso é a regra e o sigilo, a exceção (art. 2º).
type NivelAcesso string

const (
	AcessoPublico  NivelAcesso = "PUBLICO"
	AcessoRestrito NivelAcesso = "RESTRITO"
	AcessoSigiloso NivelAcesso = "SIGILOSO"
)

var (
	ErrNaoEncontrada     = errors.New("acervo: unidade de arquivamento não encontrada")
	ErrLocalNaoExiste    = errors.New("acervo: local não encontrado")
	ErrForaDoEscopo      = errors.New("acervo: fora do escopo de quem pede")
	ErrSerieInvalida     = errors.New("acervo: código de classificação inexistente ou revogado na TTDD")
	ErrBloqueioNaoExiste = errors.New("acervo: bloqueio não encontrado")
	ErrBloqueioLiberado  = errors.New("acervo: bloqueio já liberado")

	// As travas da lei: nenhuma combinação de permissões as contorna.
	ErrPermanente           = errors.New("acervo: documento de guarda permanente é inalienável e não pode ser eliminado (arts. 20 e 42)")
	ErrBloqueioJudicial     = errors.New("acervo: há bloqueio judicial ativo (arts. 32, II, e 44)")
	ErrNaoPrevistaNaTTDD    = errors.New("acervo: a TTDD não destina a série à eliminação (art. 19)")
	ErrPrazoNaoCumprido     = errors.New("acervo: o prazo de guarda ainda não foi cumprido (arts. 19 e 30)")
	ErrAguardaEvento        = errors.New("acervo: o prazo depende da data do evento (ex.: fim da vigência), ainda não informada")
	ErrDestinacaoIndefinida = errors.New("acervo: a série não tem destinação definida na TTDD")
	ErrFaseInvalida         = errors.New("acervo: operação não permitida nesta fase")
	ErrEliminado            = errors.New("acervo: unidade eliminada não pode ser alterada")
)

// ErrValidacao descreve um campo inválido.
type ErrValidacao struct{ Campo, Motivo string }

func (e *ErrValidacao) Error() string { return "acervo: " + e.Campo + ": " + e.Motivo }

func invalido(campo, motivo string) error { return &ErrValidacao{Campo: campo, Motivo: motivo} }

// Unidade é a unidade de arquivamento: caixa, pasta, livro, processo,
// dossiê, documento avulso ou mapa — o que os Anexos I a III listam.
type Unidade struct {
	ID     uuid.UUID `json:"id"`
	Numero string    `json:"numero"`
	Tipo   Tipo      `json:"tipo"`
	// Código da série na TTDD (art. 25: órgão > classe > subclasse >
	// assunto > espécie; nas tabelas, órgão.função.subfunção.série).
	CodigoClassificacao string `json:"codigo_classificacao"`
	// Assunto: o descritor do código (coluna "Assunto" dos Anexos).
	Assunto string `json:"assunto"`
	// UnidadeID: a unidade produtora, dona da unidade de arquivamento
	// (escopo do ADR 013).
	UnidadeID uuid.UUID `json:"unidade_id"`
	// Datas-limite (colunas "ano inicial" e "ano final" dos Anexos).
	DataInicial time.Time `json:"data_inicial"`
	DataFinal   time.Time `json:"data_final"`
	// Quantificação e especificação (ex.: 2 + "caixas-arquivo").
	Quantidade    int        `json:"quantidade"`
	Especificacao string     `json:"especificacao"`
	Observacoes   string     `json:"observacoes"`
	Fase          Fase       `json:"fase"`
	Destinacao    Destinacao `json:"destinacao"`
	// DataEvento: o evento que encerra a fase corrente nas séries com prazo
	// por condição (ex.: "Enquanto estiver vigorando" = fim da vigência).
	DataEvento *time.Time `json:"data_evento"`
	// ContasAprovadasEm: aprovação das contas pelo Tribunal de Contas,
	// quando a eliminação exige (quadro do Anexo I).
	ContasAprovadasEm *time.Time  `json:"contas_aprovadas_em"`
	NivelAcesso       NivelAcesso `json:"nivel_acesso"`
	HipoteseLegal     string      `json:"hipotese_legal"`
	LocalID           *uuid.UUID  `json:"local_id"`
	Posicao           string      `json:"posicao"`
	CreatedBy         *uuid.UUID  `json:"created_by"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

var tipos = map[Tipo]bool{TipoCaixa: true, TipoPasta: true, TipoLivro: true, TipoProcesso: true, TipoDossie: true,
	TipoDocumento: true, TipoMapa: true, TipoOutro: true}

var niveis = map[NivelAcesso]bool{AcessoPublico: true, AcessoRestrito: true, AcessoSigiloso: true}

// Normalizar apara os textos e aplica os padrões.
func (u *Unidade) Normalizar() {
	u.Tipo = Tipo(strings.ToUpper(strings.TrimSpace(string(u.Tipo))))
	u.CodigoClassificacao = strings.TrimSpace(u.CodigoClassificacao)
	u.Assunto = strings.TrimSpace(u.Assunto)
	u.Especificacao = strings.TrimSpace(u.Especificacao)
	u.Observacoes = strings.TrimSpace(u.Observacoes)
	u.HipoteseLegal = strings.TrimSpace(u.HipoteseLegal)
	u.Posicao = strings.TrimSpace(u.Posicao)
	u.NivelAcesso = NivelAcesso(strings.ToUpper(strings.TrimSpace(string(u.NivelAcesso))))
	if u.NivelAcesso == "" {
		u.NivelAcesso = AcessoPublico
	}
	if u.Quantidade == 0 {
		u.Quantidade = 1
	}
}

// Validar confere os campos descritivos (não a fase, que só muda pelas
// operações do ciclo de vida).
func (u Unidade) Validar() error {
	switch {
	case !tipos[u.Tipo]:
		return invalido("tipo", "use CAIXA, PASTA, LIVRO, PROCESSO, DOSSIE, DOCUMENTO, MAPA ou OUTRO")
	case u.CodigoClassificacao == "" || utf8.RuneCountInString(u.CodigoClassificacao) > 32:
		return invalido("codigo_classificacao", "informe o código da série na TTDD")
	case utf8.RuneCountInString(u.Assunto) > 500:
		return invalido("assunto", "até 500 caracteres")
	case u.UnidadeID == uuid.Nil:
		return invalido("unidade_id", "informe a unidade produtora")
	case u.DataInicial.IsZero() || u.DataFinal.IsZero():
		return invalido("datas", "informe as datas-limite")
	case u.DataFinal.Before(u.DataInicial):
		return invalido("datas", "a data final não pode ser anterior à inicial")
	case u.Quantidade < 1 || u.Quantidade > 100000:
		return invalido("quantidade", "entre 1 e 100000")
	case utf8.RuneCountInString(u.Especificacao) > 120:
		return invalido("especificacao", "até 120 caracteres")
	case utf8.RuneCountInString(u.Observacoes) > 4000:
		return invalido("observacoes", "até 4000 caracteres")
	case !niveis[u.NivelAcesso]:
		return invalido("nivel_acesso", "use PUBLICO, RESTRITO ou SIGILOSO")
	case u.NivelAcesso != AcessoPublico && u.HipoteseLegal == "":
		return invalido("hipotese_legal", "obrigatória para acesso restrito ou sigiloso (art. 2º)")
	case utf8.RuneCountInString(u.HipoteseLegal) > 500:
		return invalido("hipotese_legal", "até 500 caracteres")
	case utf8.RuneCountInString(u.Posicao) > 120:
		return invalido("posicao", "até 120 caracteres")
	case u.DataEvento != nil && u.DataEvento.Before(u.DataInicial):
		return invalido("data_evento", "não pode ser anterior à data inicial")
	}
	return nil
}

// Editavel: unidade eliminada é registro histórico (art. 38); a permanente
// aceita só a descrição e a localização, nunca a reclassificação.
func (u Unidade) Editavel() error {
	if u.Fase == FaseEliminado {
		return ErrEliminado
	}
	return nil
}

// Bloqueio é a retenção por ação judicial (arts. 30, II, 32, II, e 44).
type Bloqueio struct {
	ID              uuid.UUID  `json:"id"`
	UnidadeID       uuid.UUID  `json:"unidade_arquivamento_id"`
	Processo        string     `json:"processo"`
	Motivo          string     `json:"motivo"`
	CriadoPor       *uuid.UUID `json:"criado_por"`
	CriadoEm        time.Time  `json:"criado_em"`
	LiberadoPor     *uuid.UUID `json:"liberado_por"`
	LiberadoEm      *time.Time `json:"liberado_em"`
	MotivoLiberacao string     `json:"motivo_liberacao"`
}

// Ativo informa se o bloqueio ainda retém a unidade.
func (b Bloqueio) Ativo() bool { return b.LiberadoEm == nil }

// NovoBloqueio valida o registro do bloqueio.
func NovoBloqueio(unidade uuid.UUID, processo, motivo string) (Bloqueio, error) {
	processo, motivo = strings.TrimSpace(processo), strings.TrimSpace(motivo)
	if processo == "" || utf8.RuneCountInString(processo) > 120 {
		return Bloqueio{}, invalido("processo", "informe o número do processo judicial (até 120 caracteres)")
	}
	if utf8.RuneCountInString(motivo) > 2000 {
		return Bloqueio{}, invalido("motivo", "até 2000 caracteres")
	}
	return Bloqueio{ID: uuid.New(), UnidadeID: unidade, Processo: processo, Motivo: motivo}, nil
}

// TipoLocal é o nível da localização física (artigo técnico, Quadro 1).
type TipoLocal string

const (
	LocalPredio     TipoLocal = "PREDIO"
	LocalDeposito   TipoLocal = "DEPOSITO"
	LocalSala       TipoLocal = "SALA"
	LocalCorredor   TipoLocal = "CORREDOR"
	LocalEstante    TipoLocal = "ESTANTE"
	LocalPrateleira TipoLocal = "PRATELEIRA"
)

var tiposLocal = map[TipoLocal]bool{LocalPredio: true, LocalDeposito: true, LocalSala: true, LocalCorredor: true,
	LocalEstante: true, LocalPrateleira: true}

// Local é um ponto da hierarquia física: prédio > depósito > … > prateleira.
type Local struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parent_id"`
	Tipo     TipoLocal  `json:"tipo"`
	Nome     string     `json:"nome"`
	Ativo    bool       `json:"ativo"`
	// Caminho: "Prédio Sede > Depósito 1 > Estante 4" (montado na leitura).
	Caminho string `json:"caminho"`
}

// Validar confere o local.
func (l *Local) Validar() error {
	l.Nome = strings.TrimSpace(l.Nome)
	l.Tipo = TipoLocal(strings.ToUpper(strings.TrimSpace(string(l.Tipo))))
	if !tiposLocal[l.Tipo] {
		return invalido("tipo", "use PREDIO, DEPOSITO, SALA, CORREDOR, ESTANTE ou PRATELEIRA")
	}
	if l.Nome == "" || utf8.RuneCountInString(l.Nome) > 120 {
		return invalido("nome", "obrigatório, até 120 caracteres")
	}
	if l.ParentID != nil && *l.ParentID == l.ID {
		return invalido("parent_id", "um local não pode estar dentro de si mesmo")
	}
	return nil
}

// Movimentacao é o histórico da unidade (cadastro, transferência,
// localização, bloqueio…), além da trilha de auditoria.
type Movimentacao struct {
	ID        uuid.UUID  `json:"id"`
	Tipo      string     `json:"tipo"`
	FaseDe    *Fase      `json:"fase_de"`
	FasePara  *Fase      `json:"fase_para"`
	Descricao string     `json:"descricao"`
	Por       *uuid.UUID `json:"por"`
	Em        time.Time  `json:"em"`
}

// Tipos de movimentação.
const (
	MovCadastro      = "CADASTRO"
	MovAlteracao     = "ALTERACAO"
	MovTransferencia = "TRANSFERENCIA"
	MovLocalizacao   = "LOCALIZACAO"
	MovBloqueio      = "BLOQUEIO"
	MovDesbloqueio   = "DESBLOQUEIO"
)
