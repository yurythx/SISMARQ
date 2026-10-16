package domain

import (
	"errors"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 3 (ADR 034): o rito da eliminação da Lei 10.301/2019, arts. 31 a 40 e
// Anexos I a IV.

var (
	ErrComissaoNaoExiste = errors.New("acervo: comissão não encontrada")
	ErrSemComissao       = errors.New("acervo: o órgão não tem Comissão Setorial Permanente de Avaliação de Documentos ativa (art. 31)")
	ErrComposicao        = errors.New("acervo: composição da comissão em desacordo com o art. 31")
	ErrListagemNaoExiste = errors.New("acervo: relação de eliminação não encontrada")
	ErrListagemSituacao  = errors.New("acervo: a relação de eliminação não está na etapa que esta ação exige")
	ErrListagemVazia     = errors.New("acervo: a relação não tem unidades a eliminar")
	ErrSemParecer        = errors.New("acervo: falta o parecer da assessoria jurídica (art. 32, § 4º)")
	ErrQuorum            = errors.New("acervo: a aprovação exige o presidente e a maioria dos membros da CSPAD presentes (Anexo IV)")
	ErrEmOutraListagem   = errors.New("acervo: a unidade já está em outra relação de eliminação em andamento")
	ErrEditalNoPrazo     = errors.New("acervo: a eliminação só pode ocorrer a partir do 30º dia da publicação do edital (art. 37, § 2º)")
	ErrPedidoPendente    = errors.New("acervo: há pedido de desentranhamento sem decisão (art. 37, § 2º)")
	ErrSemAmostra        = errors.New("acervo: falta a amostragem (art. 39)")
	ErrAmostraFeita      = errors.New("acervo: a amostragem já foi feita")
	ErrPedidoForaDoPrazo = errors.New("acervo: pedido fora do prazo de 30 dias do edital")
	ErrPedidoNaoExiste   = errors.New("acervo: pedido não encontrado")
	ErrPedidoDecidido    = errors.New("acervo: o pedido já foi decidido")
)

// Tipos de comissão (art. 31).
const (
	ComissaoSetorial = "SETORIAL" // CSPAD de cada órgão
	ComissaoCentral  = "CENTRAL"  // CCPAD, junto ao Arquivo Público
)

// Papéis e áreas dos membros (art. 31, §§ 6º e 7º).
const (
	PapelPresidente = "PRESIDENTE"
	PapelSecretario = "SECRETARIO"
	PapelMembro     = "MEMBRO"

	AreaArquivologia     = "ARQUIVOLOGIA_HISTORIA"
	AreaJuridico         = "JURIDICO"
	AreaAdmFinancas      = "ADMINISTRACAO_FINANCAS"
	AreaUnidadeProdutora = "UNIDADE_PRODUTORA"
	AreaArquivoPublico   = "ARQUIVO_PUBLICO"
	AreaOutra            = "OUTRA"
)

var papeis = map[string]bool{PapelPresidente: true, PapelSecretario: true, PapelMembro: true}
var areas = map[string]bool{AreaArquivologia: true, AreaJuridico: true, AreaAdmFinancas: true, AreaUnidadeProdutora: true,
	AreaArquivoPublico: true, AreaOutra: true}

// Membro de comissão.
type Membro struct {
	UserID uuid.UUID `json:"user_id"`
	Nome   string    `json:"nome"`
	Papel  string    `json:"papel"`
	Area   string    `json:"area"`
}

// Comissao é a CSPAD de um órgão ou a CCPAD (art. 31).
type Comissao struct {
	ID         uuid.UUID  `json:"id"`
	Tipo       string     `json:"tipo"`
	EntidadeID *uuid.UUID `json:"entidade_id"`
	Nome       string     `json:"nome"`
	Portaria   string     `json:"portaria"`
	Ativa      bool       `json:"ativa"`
	Membros    []Membro   `json:"membros"`
	// Avisos: recomendações da lei não cumpridas (ex.: número par).
	Avisos []string `json:"avisos"`
}

// exigidas são as áreas que a lei pede em cada tipo de comissão.
var exigidas = map[string][]string{
	ComissaoSetorial: {AreaArquivologia, AreaJuridico, AreaAdmFinancas, AreaUnidadeProdutora}, // art. 31, § 6º
	ComissaoCentral:  {AreaArquivologia, AreaJuridico, AreaAdmFinancas, AreaArquivoPublico},   // art. 31, § 7º
}

var nomesArea = map[string]string{AreaArquivologia: "Arquivologia ou História", AreaJuridico: "assessoria jurídica",
	AreaAdmFinancas: "administração e finanças", AreaUnidadeProdutora: "unidade produtora", AreaArquivoPublico: "Arquivo Público"}

// Validar confere a comissão: tipo, órgão, nome, um presidente, membros sem
// repetição e as áreas exigidas pelo art. 31. O número ímpar é preferência
// da lei (§§ 6º e 7º): vira aviso, não erro.
func (c *Comissao) Validar() error {
	c.Nome, c.Portaria = strings.TrimSpace(c.Nome), strings.TrimSpace(c.Portaria)
	switch {
	case c.Tipo != ComissaoSetorial && c.Tipo != ComissaoCentral:
		return invalido("tipo", "use SETORIAL ou CENTRAL")
	case (c.Tipo == ComissaoCentral) != (c.EntidadeID == nil):
		return invalido("entidade_id", "a CSPAD é de um órgão; a CCPAD é do Arquivo Público, sem órgão")
	case c.Nome == "" || utf8.RuneCountInString(c.Nome) > 200:
		return invalido("nome", "obrigatório, até 200 caracteres")
	case utf8.RuneCountInString(c.Portaria) > 200:
		return invalido("portaria", "até 200 caracteres")
	}
	vistos, presidentes, temArea := map[uuid.UUID]bool{}, 0, map[string]bool{}
	for _, m := range c.Membros {
		if m.UserID == uuid.Nil || vistos[m.UserID] || !papeis[m.Papel] || !areas[m.Area] {
			return invalido("membros", "cada membro aparece uma vez, com papel e área válidos")
		}
		vistos[m.UserID] = true
		temArea[m.Area] = true
		if m.Papel == PapelPresidente {
			presidentes++
		}
	}
	if presidentes != 1 {
		return invalido("membros", "a comissão tem exatamente um presidente")
	}
	var faltam []string
	for _, a := range exigidas[c.Tipo] {
		if !temArea[a] {
			faltam = append(faltam, nomesArea[a])
		}
	}
	if len(faltam) > 0 {
		return invalido("membros", "falta servidor de: "+strings.Join(faltam, ", ")+" (art. 31)")
	}
	c.Avisos = []string{}
	if len(c.Membros)%2 == 0 {
		c.Avisos = append(c.Avisos, "A lei recomenda número ímpar de membros (art. 31, §§ 6º e 7º).")
	}
	return nil
}

// Situações da Relação de Eliminação.
const (
	ListagemRascunho   = "RASCUNHO"
	ListagemAprovada   = "APROVADA"   // ata da CSPAD (Anexo IV)
	ListagemAutorizada = "AUTORIZADA" // "Autorizo" do Arquivo Público (Anexo I)
	ListagemEmEdital   = "EM_EDITAL"  // edital publicado, prazo de 30 dias (Anexo III)
	ListagemEliminada  = "ELIMINADA"  // termo registrado (Anexo II)
	ListagemCancelada  = "CANCELADA"
)

// PrazoEdital: a eliminação só a partir do 30º dia da publicação (art. 37,
// § 2º; Anexo III).
const PrazoEdital = 30

// ItemListagem é uma unidade da relação.
type ItemListagem struct {
	UnidadeLista
	Amostra  bool `json:"amostra"`
	Retirada bool `json:"retirada"`
}

// Pedido de desentranhamento ou de cópia durante o edital (art. 37, § 2º).
type Pedido struct {
	ID          uuid.UUID  `json:"id"`
	UnidadeID   *uuid.UUID `json:"unidade_arquivamento_id"`
	Requerente  string     `json:"requerente"`
	Pedido      string     `json:"pedido"`
	RecebidoEm  time.Time  `json:"recebido_em"`
	Situacao    string     `json:"situacao"`
	Decisao     string     `json:"decisao"`
	DecididoPor *uuid.UUID `json:"decidido_por"`
	DecididoEm  *time.Time `json:"decidido_em"`
}

// Situação dos pedidos.
const (
	PedidoPendente   = "PENDENTE"
	PedidoDeferido   = "DEFERIDO"
	PedidoIndeferido = "INDEFERIDO"
)

// Listagem é a Relação de Eliminação de Documentos (Anexo I) e o processo
// que ela percorre até o Termo (Anexo II).
type Listagem struct {
	ID                 uuid.UUID  `json:"id"`
	Numero             string     `json:"numero"`
	ComissaoID         uuid.UUID  `json:"comissao_id"`
	UnidadeID          uuid.UUID  `json:"unidade_id"`
	Situacao           string     `json:"situacao"`
	ResponsavelSelecao *uuid.UUID `json:"responsavel_selecao"`
	ParecerJuridico    string     `json:"parecer_juridico"`
	ContasExercicios   string     `json:"contas_exercicios"`
	ContasAprovadasEm  *time.Time `json:"contas_aprovadas_em"`

	AtaNumero        string      `json:"ata_numero"`
	AtaEm            *time.Time  `json:"ata_em"`
	AtaLocal         string      `json:"ata_local"`
	AtaManifestacoes string      `json:"ata_manifestacoes"`
	AtaPresentes     []uuid.UUID `json:"ata_presentes"`

	AutorizadaPor   *uuid.UUID `json:"autorizada_por"`
	AutorizadaEm    *time.Time `json:"autorizada_em"`
	DevolucaoMotivo string     `json:"devolucao_motivo"`

	EditalNumero      string     `json:"edital_numero"`
	EditalPublicadoEm *time.Time `json:"edital_publicado_em"`
	EditalDiario      string     `json:"edital_diario"`

	AmostraCriterios  string     `json:"amostra_criterios"`
	AmostraPercentual *float64   `json:"amostra_percentual"`
	AmostraEm         *time.Time `json:"amostra_em"`

	TermoEliminadoEm     *time.Time `json:"termo_eliminado_em"`
	TermoAcompanhamento  string     `json:"termo_acompanhamento"`
	TermoDestinacaoPapel string     `json:"termo_destinacao_papel"`
	TermoRegistradoPor   *uuid.UUID `json:"termo_registrado_por"`

	CriadaPor *uuid.UUID `json:"criada_por"`
	CriadaEm  time.Time  `json:"criada_em"`

	Itens   []ItemListagem `json:"itens"`
	Pedidos []Pedido       `json:"pedidos"`
	// Resumo, para a listagem de relações.
	TotalItens int `json:"total_itens"`
}

// Ativa informa se a relação ainda está em andamento.
func (l Listagem) Ativa() bool {
	return l.Situacao != ListagemEliminada && l.Situacao != ListagemCancelada
}

// AEliminar são os itens que vão de fato à eliminação (sem amostra e sem os
// retirados por pedido deferido).
func (l Listagem) AEliminar() []ItemListagem {
	var out []ItemListagem
	for _, i := range l.Itens {
		if !i.Amostra && !i.Retirada {
			out = append(out, i)
		}
	}
	return out
}

// PrazoFinal é a primeira data em que a eliminação pode ocorrer.
func (l Listagem) PrazoFinal() *time.Time {
	if l.EditalPublicadoEm == nil {
		return nil
	}
	t := l.EditalPublicadoEm.AddDate(0, 0, PrazoEdital)
	return &t
}

func exigir(l Listagem, situacoes ...string) error {
	for _, s := range situacoes {
		if l.Situacao == s {
			return nil
		}
	}
	return ErrListagemSituacao
}

// PodeAprovar confere a ata da CSPAD (Anexo IV): relação com itens, parecer
// jurídico (art. 32, § 4º), número e data da ata, e quórum — o presidente e a
// maioria dos membros presentes.
func PodeAprovar(l Listagem, c Comissao, numero string, em time.Time, presentes []uuid.UUID) error {
	if err := exigir(l, ListagemRascunho); err != nil {
		return err
	}
	switch {
	case len(l.AEliminar()) == 0:
		return ErrListagemVazia
	case strings.TrimSpace(l.ParecerJuridico) == "":
		return ErrSemParecer
	case strings.TrimSpace(numero) == "" || utf8.RuneCountInString(numero) > 40:
		return invalido("ata_numero", "informe o número da ata (até 40 caracteres)")
	case em.IsZero():
		return invalido("ata_em", "informe a data e a hora da reunião")
	}
	membros, presidente := map[uuid.UUID]string{}, false
	for _, m := range c.Membros {
		membros[m.UserID] = m.Papel
	}
	contados := map[uuid.UUID]bool{}
	for _, p := range presentes {
		papel, ok := membros[p]
		if !ok {
			return invalido("presentes", "só membros da comissão constam da lista de presença")
		}
		contados[p] = true
		presidente = presidente || papel == PapelPresidente
	}
	if !presidente || len(contados)*2 <= len(c.Membros) {
		return ErrQuorum
	}
	return nil
}

// PodePublicar confere o edital: autorizada, com número e data de publicação
// que não é futura.
func PodePublicar(l Listagem, numero string, publicadoEm, hoje time.Time) error {
	if err := exigir(l, ListagemAutorizada); err != nil {
		return err
	}
	switch {
	case strings.TrimSpace(numero) == "" || utf8.RuneCountInString(numero) > 40:
		return invalido("edital_numero", "informe o número do edital (até 40 caracteres)")
	case publicadoEm.IsZero() || publicadoEm.After(hoje):
		return invalido("edital_publicado_em", "informe a data de publicação (não pode ser futura)")
	}
	return nil
}

// PodeReceberPedido: só durante o edital e dentro dos 30 dias.
func PodeReceberPedido(l Listagem, recebidoEm time.Time) error {
	if err := exigir(l, ListagemEmEdital); err != nil {
		return err
	}
	if recebidoEm.Before(*l.EditalPublicadoEm) || recebidoEm.After(*l.PrazoFinal()) {
		return ErrPedidoForaDoPrazo
	}
	return nil
}

// TamanhoAmostra é quantos itens vão à guarda permanente (ao menos um).
func TamanhoAmostra(itens int, percentual float64) int {
	n := int(math.Ceil(float64(itens) * percentual / 100))
	return max(1, min(n, itens))
}

// PodeAmostrar confere a amostragem (art. 39): durante o edital, uma vez,
// com critérios (art. 32, III) e percentual entre 0 e 100.
func PodeAmostrar(l Listagem, criterios string, percentual float64) error {
	if err := exigir(l, ListagemEmEdital); err != nil {
		return err
	}
	switch {
	case l.AmostraEm != nil:
		return ErrAmostraFeita
	case strings.TrimSpace(criterios) == "" || utf8.RuneCountInString(criterios) > 2000:
		return invalido("amostra_criterios", "descreva os critérios qualitativos e quantitativos (art. 39)")
	case percentual <= 0 || percentual > 100:
		return invalido("amostra_percentual", "entre 0 e 100")
	case len(l.AEliminar()) == 0:
		return ErrListagemVazia
	}
	return nil
}

// PodeRegistrarTermo confere o Termo de Eliminação (Anexo II): edital
// vencido, sem pedido pendente, amostragem feita, data da eliminação entre o
// fim do prazo e hoje.
func PodeRegistrarTermo(l Listagem, eliminadoEm, hoje time.Time) error {
	if err := exigir(l, ListagemEmEdital); err != nil {
		return err
	}
	prazo := *l.PrazoFinal()
	switch {
	case hoje.Before(prazo) || eliminadoEm.Before(prazo):
		return ErrEditalNoPrazo
	case eliminadoEm.After(hoje):
		return invalido("termo_eliminado_em", "a data da eliminação não pode ser futura")
	case l.AmostraEm == nil:
		return ErrSemAmostra
	}
	for _, p := range l.Pedidos {
		if p.Situacao == PedidoPendente {
			return ErrPedidoPendente
		}
	}
	return nil
}

// PodeCancelar: antes do termo e antes da amostragem (a amostra já foi para
// a guarda permanente e não volta).
func PodeCancelar(l Listagem) error {
	if !l.Ativa() || l.AmostraEm != nil {
		return ErrListagemSituacao
	}
	return nil
}
