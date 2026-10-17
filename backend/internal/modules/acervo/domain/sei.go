package domain

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 7 (ADR 042): ingresso do SEI. A decisão 0.7 registrou o envelope
// e adiou a integração (sem versão, API, URL nem autorização — ADR 043).
// A função recusa o ingresso: não há cliente HTTP nem evento.

// ErrSEIAdiado: o município ainda não autorizou a integração; inventar um
// cliente seria uma API de mentira.
var ErrSEIAdiado = errors.New("acervo: ingresso do SEI adiado até autorização, API e ambiente de homologação (decisão 0.7)")

// ProcessoSEI é o que um pacote AIP do SEI precisaria para o ingresso
// (metadados mínimos; documentos e hashes ficam para o ADR da 0.7).
type ProcessoSEI struct {
	Protocolo       string
	Assunto         string
	CodigoTTDD      string
	UnidadeOrigemID uuid.UUID
	NumeroSEI       string
}

// UnidadeDoSEI valida o envelope e recusa o ingresso. Quando a 0.7 existir,
// um ADR posterior substitui o retorno por UnidadeDoTramite-equivalente.
func UnidadeDoSEI(p ProcessoSEI) (Unidade, error) {
	switch {
	case p.UnidadeOrigemID == uuid.Nil:
		return Unidade{}, invalido("unidade_id", "informe a unidade produtora")
	case utf8.RuneCountInString(p.NumeroSEI) > 80:
		return Unidade{}, invalido("numero_sei", "até 80 caracteres")
	case p.NumeroSEI == "" && p.Protocolo == "":
		return Unidade{}, invalido("numero_sei", "informe o número no SEI ou o protocolo")
	}
	return Unidade{}, ErrSEIAdiado
}
