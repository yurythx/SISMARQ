package domain

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Serie é o que o Acervo precisa da TTDD (Atlas). A porta é do Acervo; o
// adaptador é ligado em internal/app — os plugins não se importam.
type Serie struct {
	Codigo           string
	Descritor        string
	CorrenteAnos     *int
	CorrenteCondicao string
	IntermAnos       *int
	IntermCondicao   string
	// Destinacao vazia = a tabela não define (DestinacaoIndefinida).
	Destinacao  Destinacao
	Observacoes string
	Vigente     bool
}

// TTDD consulta a Tabela de Temporalidade em vigor.
type TTDD interface {
	// Serie devolve a série pelo código; ErrSerieInvalida se não existir.
	Serie(ctx context.Context, codigo string) (Serie, error)
}

// DestinacaoDa traduz a destinação da série.
func DestinacaoDa(s Serie) Destinacao {
	if s.Destinacao == "" {
		return DestinacaoIndefinida
	}
	return s.Destinacao
}

// Situação do prazo de guarda.
const (
	SituacaoNoCorrente        = "NO_CORRENTE"        // prazo corrente correndo
	SituacaoTransferir        = "TRANSFERIR"         // corrente cumprido: vai ao intermediário (art. 29)
	SituacaoNoIntermediario   = "NO_INTERMEDIARIO"   // prazo intermediário correndo
	SituacaoAvaliarEliminacao = "AVALIAR_ELIMINACAO" // prazos cumpridos, destino eliminação: vai à CSPAD
	SituacaoRecolher          = "RECOLHER"           // prazos cumpridos, destino permanente: recolhimento (art. 53)
	SituacaoAguardaEvento     = "AGUARDA_EVENTO"     // prazo por condição sem a data do evento
	SituacaoAguardaDestinacao = "AGUARDA_DESTINACAO" // série sem destinação definida
	SituacaoBloqueada         = "BLOQUEADA"          // bloqueio judicial ativo (art. 44)
	SituacaoPermanente        = "PERMANENTE"
	SituacaoEliminada         = "ELIMINADA"
)

// Prazos é o resultado do cálculo para uma unidade.
type Prazos struct {
	// FimCorrente/FimIntermediaria: nil quando dependem de um evento.
	FimCorrente      *time.Time `json:"fim_corrente"`
	FimIntermediaria *time.Time `json:"fim_intermediaria"`
	Situacao         string     `json:"situacao"`
	// Regra: o texto da TTDD que gerou o prazo, para conferência humana.
	RegraCorrente      string `json:"regra_corrente"`
	RegraIntermediaria string `json:"regra_intermediaria"`
}

var (
	reMeses       = regexp.MustCompile(`^(\d{1,3})\s+m[eê]s(es)?$`)
	reDias        = regexp.MustCompile(`^(\d{1,4})\s+dias?$`)
	reDiasEvento  = regexp.MustCompile(`^(\d{1,4})\s+dias?\s+ap[oó]s\s+a\s+data\s+do\s+evento$`)
	reMesesEvento = regexp.MustCompile(`^(\d{1,3})\s+m[eê]s(es)?\s+ap[oó]s\s+a\s+data\s+do\s+evento$`)
)

// fimDaFase aplica o prazo de uma fase a partir de inicio. Condições:
//   - "Enquanto estiver vigorando" (727 séries da TTDD): termina no evento;
//   - "NN meses" / "NN dias": a partir do início;
//   - "NN dias após a data do evento": a partir do evento.
//
// Sem anos e sem condição, a fase não tem permanência (prazo zero).
// Condição desconhecida: depende de um evento informado à mão.
func fimDaFase(inicio time.Time, anos *int, condicao string, evento *time.Time) (*time.Time, string) {
	c := strings.ToLower(strings.TrimSpace(condicao))
	regra := condicao
	if anos != nil {
		regra = strconv.Itoa(*anos) + " ano(s)"
		if c != "" {
			regra += " / " + condicao
		}
	}
	switch {
	case strings.Contains(c, "vigor"):
		if evento == nil {
			return nil, regra
		}
		fim := *evento
		if anos != nil {
			fim = fim.AddDate(*anos, 0, 0)
		}
		return &fim, regra
	case reDiasEvento.MatchString(c) || reMesesEvento.MatchString(c):
		if evento == nil {
			return nil, regra
		}
		return ptr(somar(*evento, c)), regra
	case reMeses.MatchString(c) || reDias.MatchString(c):
		return ptr(somar(inicio, c)), regra
	case c != "":
		// Condição que o sistema não interpreta: a data do evento decide.
		if evento == nil {
			return nil, regra
		}
		return ptr(*evento), regra
	case anos != nil:
		return ptr(inicio.AddDate(*anos, 0, 0)), regra
	}
	return ptr(inicio), "sem permanência"
}

func somar(t time.Time, c string) time.Time {
	n, _ := strconv.Atoi(strings.Fields(c)[0])
	if strings.Contains(c, "m") && !strings.Contains(c, "dia") {
		return t.AddDate(0, n, 0)
	}
	return t.AddDate(0, 0, n)
}

func ptr(t time.Time) *time.Time { return &t }

// CalcularPrazos aplica a TTDD à unidade (arts. 29 e 30). A contagem da fase
// corrente parte da data-limite final da unidade; a intermediária, do fim da
// corrente. (Critério registrado no ADR 032, a validar pela Divisão de
// Arquivo — ROADMAP_SISMARQ.md, seção 6, item 5.)
func CalcularPrazos(u Unidade, s Serie, bloqueada bool, hoje time.Time) Prazos {
	p := Prazos{}
	p.FimCorrente, p.RegraCorrente = fimDaFase(u.DataFinal, s.CorrenteAnos, s.CorrenteCondicao, u.DataEvento)
	if p.FimCorrente != nil {
		p.FimIntermediaria, p.RegraIntermediaria = fimDaFase(*p.FimCorrente, s.IntermAnos, s.IntermCondicao, u.DataEvento)
	} else {
		_, p.RegraIntermediaria = fimDaFase(u.DataFinal, s.IntermAnos, s.IntermCondicao, nil)
	}
	p.Situacao = situacao(u, p, DestinacaoDa(s), bloqueada, hoje)
	return p
}

func cumprido(fim *time.Time, hoje time.Time) bool { return fim != nil && !hoje.Before(*fim) }

func situacao(u Unidade, p Prazos, dest Destinacao, bloqueada bool, hoje time.Time) string {
	if u.RecolhidaEm != nil {
		return SituacaoRecolhida
	}
	switch u.Fase {
	case FasePermanente:
		return SituacaoPermanente
	case FaseEliminado:
		return SituacaoEliminada
	}
	if bloqueada {
		return SituacaoBloqueada
	}
	if p.FimCorrente == nil {
		return SituacaoAguardaEvento
	}
	if u.Fase == FaseCorrente {
		if cumprido(p.FimCorrente, hoje) {
			return SituacaoTransferir
		}
		return SituacaoNoCorrente
	}
	switch {
	case p.FimIntermediaria == nil:
		return SituacaoAguardaEvento
	case !cumprido(p.FimIntermediaria, hoje):
		return SituacaoNoIntermediario
	case dest == DestinacaoEliminacao:
		return SituacaoAvaliarEliminacao
	case dest == DestinacaoGuardaPermanente:
		return SituacaoRecolher
	}
	return SituacaoAguardaDestinacao
}

// PodeTransferir: do corrente ao intermediário, só com o prazo corrente
// cumprido (art. 29, I).
func PodeTransferir(u Unidade, p Prazos, hoje time.Time) error {
	switch {
	case u.Fase != FaseCorrente:
		return ErrFaseInvalida
	case p.FimCorrente == nil:
		return ErrAguardaEvento
	case !cumprido(p.FimCorrente, hoje):
		return ErrPrazoNaoCumprido
	}
	return nil
}

// PodeEliminar reúne as travas da lei. É a única porta para a eliminação
// (Fase 3) e vale para qualquer pessoa, inclusive o administrador:
//   - permanente nunca (arts. 20 e 42);
//   - bloqueio judicial ativo, nunca (art. 44);
//   - só o que a TTDD destina à eliminação (art. 19);
//   - só com os prazos cumpridos (arts. 19 e 30).
func PodeEliminar(u Unidade, s Serie, bloqueada bool, hoje time.Time) error {
	switch {
	case u.Fase == FasePermanente:
		return ErrPermanente
	case u.Fase == FaseEliminado:
		return ErrEliminado
	case bloqueada:
		return ErrBloqueioJudicial
	case DestinacaoDa(s) == DestinacaoIndefinida:
		return ErrDestinacaoIndefinida
	case DestinacaoDa(s) != DestinacaoEliminacao:
		return ErrNaoPrevistaNaTTDD
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
