package domain

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 7 (ADR 042): ingresso do SEI. A decisão 0.7 (versão, API, quem
// autoriza) ainda não existe. O tipo descreve o envelope que um processo
// concluído precisaria para virar unidade de arquivamento; a função
// recusa o ingresso até essa decisão — não há cliente HTTP nem evento.

// ErrSEIAdiado: o município ainda não fechou a 0.7; inventar um cliente
// seria uma API de mentira.
var ErrSEIAdiado = errors.New("acervo: ingresso do SEI adiado até a decisão 0.7 (versão, API e autorização)")

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
