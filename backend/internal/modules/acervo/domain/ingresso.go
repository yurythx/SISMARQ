package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 2, ponte com o Trâmite (ADR 039): o processo encerrado vira unidade
// de arquivamento do tipo PROCESSO, classificada pela série que o Trâmite
// já gravou (ADR 026). Sem série, o processo não entra — a guarda documental
// do Trâmite continua orientando, e o arquivo classifica depois. Reabrir o
// processo não apaga a unidade (o ingresso é fato arquivístico).

// ErrSemClassificacao: o evento chegou sem código da TTDD.
var ErrSemClassificacao = errors.New("acervo: processo encerrado sem série da TTDD não entra no acervo")

// ProcessoEncerrado é o que o evento do Trâmite precisa para o ingresso.
// Assunto pode vir vazio (processo restrito ou sigiloso: o conteúdo não
// sai no barramento).
type ProcessoEncerrado struct {
	ID              uuid.UUID
	Numero          string
	Assunto         string
	CodigoTTDD      string
	UnidadeOrigemID uuid.UUID
	CreatedAt       time.Time
	ConcluidoAt     time.Time
	Sigilo          string
}

// UnidadeDoTramite monta a unidade a partir do processo encerrado. A
// destinação e a fase saem da TTDD na aplicação, como no cadastro pela tela.
func UnidadeDoTramite(p ProcessoEncerrado) (Unidade, error) {
	if p.ID == uuid.Nil {
		return Unidade{}, invalido("tramite_processo_id", "informe o processo")
	}
	if p.CodigoTTDD == "" {
		return Unidade{}, ErrSemClassificacao
	}
	if p.UnidadeOrigemID == uuid.Nil {
		return Unidade{}, invalido("unidade_id", "informe a unidade produtora")
	}
	inicio, fim := p.CreatedAt, p.ConcluidoAt
	if inicio.IsZero() {
		inicio = fim
	}
	if fim.IsZero() {
		fim = inicio
	}
	if fim.Before(inicio) {
		inicio, fim = fim, inicio
	}
	assunto := p.Assunto
	if assunto == "" {
		assunto = "Processo " + p.Numero
	}
	if utf8.RuneCountInString(assunto) > 500 {
		assunto = string([]rune(assunto)[:500])
	}
	nivel, hipotese := AcessoPublico, ""
	switch p.Sigilo {
	case "restrito":
		nivel, hipotese = AcessoRestrito, "Processo administrativo restrito (Trâmite)"
	case "sigiloso":
		nivel, hipotese = AcessoSigiloso, "Processo administrativo sigiloso (Trâmite)"
	}
	obs := "Ingresso automático do processo " + p.Numero + " (Trâmite, ADR 039)."
	id := p.ID
	u := Unidade{
		Tipo: TipoProcesso, CodigoClassificacao: p.CodigoTTDD, Assunto: assunto,
		UnidadeID: p.UnidadeOrigemID, DataInicial: inicio, DataFinal: fim, Quantidade: 1,
		Especificacao: "processo", Observacoes: obs, NivelAcesso: nivel, HipoteseLegal: hipotese,
		TramiteProcessoID: &id,
	}
	u.Normalizar()
	if err := u.Validar(); err != nil {
		return u, err
	}
	return u, nil
}
