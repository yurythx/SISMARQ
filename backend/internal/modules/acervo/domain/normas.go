package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Instruções normativas do Arquivo Público (art. 54; ADR 036). A norma muda
// com o tempo: cada alteração é uma nova redação (versão). A redação
// publicada nunca é editada nem apagada — fica no histórico, com o período
// em que vigorou e o ato que a alterou, como o texto compilado de uma lei.

var (
	ErrNormaNaoExiste = errors.New("acervo: instrução normativa não encontrada")
	ErrNormaSituacao  = errors.New("acervo: a instrução normativa não está na situação que a operação exige")
	ErrSemRascunho    = errors.New("acervo: não há redação em rascunho para publicar")
	ErrNumeroNorma    = errors.New("acervo: já existe instrução normativa com este número")
)

// Situação da norma.
const (
	NormaRascunho = "RASCUNHO" // nunca publicada
	NormaVigente  = "VIGENTE"
	NormaRevogada = "REVOGADA"
)

// Situação da redação.
const (
	RedacaoRascunho  = "RASCUNHO"
	RedacaoPublicada = "PUBLICADA"
)

// MaxTextoNorma limita o texto de uma redação (Markdown).
const MaxTextoNorma = 100000

var reNumeroNorma = regexp.MustCompile(`^\d{1,4}/\d{4}$`)

// Redacao é uma versão do texto da norma.
type Redacao struct {
	Versao int    `json:"versao"`
	Ementa string `json:"ementa"`
	Texto  string `json:"texto"`
	// Motivo: "Redação original" ou o ato que alterou (ex.: "Alterada
	// pela IN 3/2027: novo prazo do art. 4º").
	Motivo      string     `json:"motivo"`
	Situacao    string     `json:"situacao"`
	PublicadaEm *time.Time `json:"publicada_em"`
	Diario      string     `json:"diario"`
	// VigenteAte: a publicação da redação seguinte (ou a revogação).
	VigenteAte *time.Time `json:"vigente_ate"`
	CriadaEm   time.Time  `json:"criada_em"`
}

// Norma é a instrução normativa com as suas redações (a mais nova
// primeiro).
type Norma struct {
	ID            uuid.UUID  `json:"id"`
	Numero        string     `json:"numero"`
	Situacao      string     `json:"situacao"`
	Ementa        string     `json:"ementa"` // da redação vigente (ou do rascunho, se nunca publicada)
	Versao        int        `json:"versao"` // a redação vigente (0 = nunca publicada)
	PublicadaEm   *time.Time `json:"publicada_em"`
	RevogadaEm    *time.Time `json:"revogada_em"`
	RevogacaoAto  string     `json:"revogacao_ato"`
	TemRascunho   bool       `json:"tem_rascunho"`
	CriadaEm      time.Time  `json:"criada_em"`
	Redacoes      []Redacao  `json:"redacoes"`
	TotalRedacoes int        `json:"total_redacoes"`
}

// Vigente devolve a redação em vigor (nil se nunca publicada).
func (n Norma) Vigente() *Redacao {
	for i := range n.Redacoes {
		if n.Redacoes[i].Situacao == RedacaoPublicada {
			return &n.Redacoes[i]
		}
	}
	return nil
}

// Rascunho devolve a redação em rascunho (nil se não houver).
func (n Norma) Rascunho() *Redacao {
	for i := range n.Redacoes {
		if n.Redacoes[i].Situacao == RedacaoRascunho {
			return &n.Redacoes[i]
		}
	}
	return nil
}

// Publicas: só as redações publicadas (o que o público consulta).
func (n Norma) Publicas() Norma {
	out := n
	out.Redacoes = []Redacao{}
	for _, r := range n.Redacoes {
		if r.Situacao == RedacaoPublicada {
			out.Redacoes = append(out.Redacoes, r)
		}
	}
	out.TemRascunho = false
	return out
}

// ValidarRedacao confere ementa, texto e motivo de uma redação.
func ValidarRedacao(r *Redacao) error {
	r.Ementa, r.Texto, r.Motivo = strings.TrimSpace(r.Ementa), strings.TrimSpace(r.Texto), strings.TrimSpace(r.Motivo)
	switch {
	case r.Ementa == "" || utf8.RuneCountInString(r.Ementa) > 500:
		return invalido("ementa", "informe a ementa (até 500 caracteres)")
	case r.Texto == "" || utf8.RuneCountInString(r.Texto) > MaxTextoNorma:
		return invalido("texto", "informe o texto (até 100.000 caracteres)")
	case r.Motivo == "" || utf8.RuneCountInString(r.Motivo) > 500:
		return invalido("motivo", "informe o motivo da redação: original ou o ato que altera (até 500 caracteres)")
	}
	return nil
}

// ValidarNumero: "N/AAAA" (ex.: 1/2026).
func ValidarNumero(numero string) (string, error) {
	numero = strings.TrimSpace(numero)
	if !reNumeroNorma.MatchString(numero) {
		return numero, invalido("numero", "use número/ano, ex.: 1/2026")
	}
	return numero, nil
}

// PodePublicarNorma: há rascunho; a data da publicação não é futura nem anterior
// à da redação em vigor; a norma não foi revogada.
func PodePublicarNorma(n Norma, em, hoje time.Time, diario string) error {
	if n.Situacao == NormaRevogada {
		return ErrNormaSituacao
	}
	if n.Rascunho() == nil {
		return ErrSemRascunho
	}
	if v := n.Vigente(); v != nil && em.Before(*v.PublicadaEm) {
		return invalido("publicada_em", "não pode ser anterior à publicação da redação em vigor")
	}
	switch {
	case em.IsZero() || em.After(hoje):
		return invalido("publicada_em", "informe a data da publicação (não futura)")
	case utf8.RuneCountInString(strings.TrimSpace(diario)) > 200:
		return invalido("diario", "até 200 caracteres")
	}
	return nil
}

// PodeRevogar: só a vigente, com o ato e a data (não anterior à redação em
// vigor nem futura).
func PodeRevogar(n Norma, em, hoje time.Time, ato string) error {
	if n.Situacao != NormaVigente {
		return ErrNormaSituacao
	}
	ato = strings.TrimSpace(ato)
	switch {
	case ato == "" || utf8.RuneCountInString(ato) > 500:
		return invalido("ato", "informe o ato de revogação (até 500 caracteres)")
	case em.IsZero() || em.After(hoje) || em.Before(*n.Vigente().PublicadaEm):
		return invalido("revogada_em", "entre a publicação da redação em vigor e hoje")
	}
	return nil
}

// Montar preenche, a partir das redações (a mais nova primeiro), a ementa,
// a versão e a publicação em vigor, se há rascunho e até quando cada redação
// vigorou: até a publicação da seguinte ou, na última, até a revogação.
func (n *Norma) Montar() {
	n.TemRascunho, n.TotalRedacoes = false, 0
	var seguinte *time.Time = n.RevogadaEm
	for i := range n.Redacoes {
		r := &n.Redacoes[i]
		if r.Situacao == RedacaoRascunho {
			n.TemRascunho = true
			if n.Versao == 0 {
				n.Ementa = r.Ementa
			}
			continue
		}
		n.TotalRedacoes++
		r.VigenteAte = seguinte
		seguinte = r.PublicadaEm
		if n.Versao == 0 {
			n.Ementa, n.Versao, n.PublicadaEm = r.Ementa, r.Versao, r.PublicadaEm
		}
	}
}
