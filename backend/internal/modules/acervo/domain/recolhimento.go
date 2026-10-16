package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 4 (ADR 035): recolhimento ao Arquivo Público e guarda permanente
// (Lei 10.301/2019, arts. 5º, VII, 16, 17, 41, 42, 47, 50, 51 e 53).

var (
	ErrRecolhimentoNaoExiste = errors.New("acervo: recolhimento não encontrado")
	ErrRecolhimentoSituacao  = errors.New("acervo: o recolhimento não está mais aguardando o Arquivo Público")
	ErrEmOutroRecolhimento   = errors.New("acervo: a unidade já está num recolhimento aguardando o Arquivo Público")
	ErrRecolhida             = errors.New("acervo: a unidade já foi recolhida ao Arquivo Público")
	ErrNaoPermanente         = errors.New("acervo: a TTDD não destina a série à guarda permanente (art. 41)")
	ErrRecolherDoCorrente    = errors.New("acervo: no recolhimento regular a unidade passa antes pelo arquivo intermediário (arts. 29 e 30); " +
		"direto do corrente, só na extinção ou na desestatização (arts. 50 e 51)")
)

// Situação do recolhimento.
const (
	RecolhimentoEnviado   = "ENVIADO"
	RecolhimentoRecebido  = "RECEBIDO"
	RecolhimentoDevolvido = "DEVOLVIDO"
	RecolhimentoCancelado = "CANCELADO"
)

// Motivo do recolhimento.
const (
	// MotivoRegular: prazos da TTDD cumpridos, no cronograma do Regimento
	// Interno do Arquivo Público (arts. 47 e 53).
	MotivoRegular = "REGULAR"
	// MotivoExtincao: extinção do órgão ou da entidade (arts. 17 e 50).
	MotivoExtincao = "EXTINCAO"
	// MotivoDesestatizacao: empresa em desestatização (arts. 16 e 51).
	MotivoDesestatizacao = "DESESTATIZACAO"
)

// MovRecolhimento é a movimentação do recolhimento (relação e entrada).
const MovRecolhimento = "RECOLHIMENTO"

// SituacaoRecolhida: sob a guarda do Arquivo Público (art. 5º, VII).
const SituacaoRecolhida = "RECOLHIDA"

// MaxItensRecolhimento limita um recolhimento (uma remessa física).
const MaxItensRecolhimento = 2000

// Preparo é o que o art. 53 exige do conjunto antes da entrada no Arquivo
// Público. O instrumento descritivo é gerado pelo sistema a partir das
// unidades relacionadas (identificação, acesso e controle).
type Preparo struct {
	Classificados  bool `json:"classificados"`
	Avaliados      bool `json:"avaliados"`
	Organizados    bool `json:"organizados"`
	Higienizados   bool `json:"higienizados"`
	Acondicionados bool `json:"acondicionados"`
}

// Completo: todas as exigências do art. 53 declaradas pela origem.
func (p Preparo) Completo() bool {
	return p.Classificados && p.Avaliados && p.Organizados && p.Higienizados && p.Acondicionados
}

// Recolhimento leva o conjunto de guarda permanente da origem ao Arquivo
// Público. Na entrada, o Arquivo registra o número de entrada e a origem
// recebe a cópia do registro (art. 5º, VII).
type Recolhimento struct {
	ID       uuid.UUID `json:"id"`
	Numero   string    `json:"numero"`
	OrigemID uuid.UUID `json:"origem_id"`
	Motivo   string    `json:"motivo"`
	// Fundamento: o ato de extinção ou o edital de desestatização.
	Fundamento string  `json:"fundamento"`
	Preparo    Preparo `json:"preparo"`
	// Historico: história administrativa e âmbito do conjunto, para o
	// instrumento descritivo.
	Historico       string         `json:"historico"`
	Observacao      string         `json:"observacao"`
	Situacao        string         `json:"situacao"`
	CriadoPor       *uuid.UUID     `json:"criado_por"`
	CriadoEm        time.Time      `json:"criado_em"`
	RecebidoPor     *uuid.UUID     `json:"recebido_por"`
	RecebidoEm      *time.Time     `json:"recebido_em"`
	RegistroEntrada string         `json:"registro_entrada"`
	LocalID         *uuid.UUID     `json:"local_id"`
	DevolucaoMotivo string         `json:"devolucao_motivo"`
	Itens           []UnidadeLista `json:"itens"`
	TotalItens      int            `json:"total_itens"`
}

var motivos = map[string]bool{MotivoRegular: true, MotivoExtincao: true, MotivoDesestatizacao: true}

// NovoRecolhimento valida o pedido: origem, motivo (os especiais com o
// fundamento), o preparo do art. 53 e de 1 a 2.000 unidades sem repetição.
func NovoRecolhimento(r Recolhimento, unidades []uuid.UUID) (Recolhimento, []uuid.UUID, error) {
	r.Motivo = strings.ToUpper(strings.TrimSpace(r.Motivo))
	r.Fundamento, r.Historico, r.Observacao = strings.TrimSpace(r.Fundamento), strings.TrimSpace(r.Historico), strings.TrimSpace(r.Observacao)
	vistos := map[uuid.UUID]bool{}
	var ids []uuid.UUID
	for _, id := range unidades {
		if id != uuid.Nil && !vistos[id] {
			vistos[id] = true
			ids = append(ids, id)
		}
	}
	switch {
	case r.OrigemID == uuid.Nil:
		return r, nil, invalido("origem_id", "informe a unidade de origem")
	case !motivos[r.Motivo]:
		return r, nil, invalido("motivo", "use REGULAR, EXTINCAO ou DESESTATIZACAO")
	case r.Motivo != MotivoRegular && r.Fundamento == "":
		return r, nil, invalido("fundamento", "informe o ato de extinção ou o edital de desestatização (arts. 50 e 51)")
	case utf8.RuneCountInString(r.Fundamento) > 500:
		return r, nil, invalido("fundamento", "até 500 caracteres")
	case !r.Preparo.Completo():
		return r, nil, invalido("preparo", "os documentos devem estar classificados, avaliados, organizados, higienizados e acondicionados (art. 53)")
	case utf8.RuneCountInString(r.Historico) > 8000:
		return r, nil, invalido("historico", "até 8000 caracteres")
	case utf8.RuneCountInString(r.Observacao) > 1000:
		return r, nil, invalido("observacao", "até 1000 caracteres")
	case len(ids) == 0 || len(ids) > MaxItensRecolhimento:
		return r, nil, invalido("unidades", "de 1 a 2000 unidades de arquivamento por recolhimento")
	}
	r.ID, r.Situacao = uuid.New(), RecolhimentoEnviado
	return r, ids, nil
}

// PodeRecolher reúne as condições da unidade:
//   - nem recolhida nem eliminada;
//   - destinação guarda permanente (art. 41, I). A amostra da eliminação já
//     é permanente (art. 39) e vale a destinação gravada na unidade;
//   - no recolhimento regular, já no intermediário com o prazo cumprido, ou
//     permanente guardada no órgão (art. 42); nos especiais (extinção e
//     desestatização), em qualquer fase (arts. 50 e 51).
func PodeRecolher(u Unidade, s Serie, motivo string, hoje time.Time) error {
	switch {
	case u.RecolhidaEm != nil:
		return ErrRecolhida
	case u.Fase == FaseEliminado:
		return ErrEliminado
	}
	dest := DestinacaoDa(s)
	if u.Fase == FasePermanente {
		dest = u.Destinacao
	}
	switch {
	case dest == DestinacaoIndefinida:
		return ErrDestinacaoIndefinida
	case dest != DestinacaoGuardaPermanente:
		return ErrNaoPermanente
	case motivo != MotivoRegular || u.Fase == FasePermanente:
		return nil
	case u.Fase == FaseCorrente:
		return ErrRecolherDoCorrente
	}
	p := CalcularPrazos(u, s, false, hoje)
	switch {
	case p.FimIntermediaria == nil:
		return ErrAguardaEvento
	case !cumprido(p.FimIntermediaria, hoje):
		return ErrPrazoNaoCumprido
	}
	return nil
}

// PodeDecidir: só o recolhimento enviado é recebido, devolvido ou
// cancelado.
func PodeDecidir(r Recolhimento) error {
	if r.Situacao != RecolhimentoEnviado {
		return ErrRecolhimentoSituacao
	}
	return nil
}
