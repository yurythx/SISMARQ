package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 6 (ADR 038): acesso, busca e difusão. O ostensivo é a regra (art.
// 2º): a consulta pública mostra a descrição das unidades de acesso público
// e, quando liberado, o representante digital. O restrito e o sigiloso só
// aparecem para quem tem a unidade no escopo.

var (
	ErrNaoOstensivo       = errors.New("acervo: só o documento ostensivo (nível público) vai para a consulta pública (art. 2º)")
	ErrPublicacao         = errors.New("acervo: o representante já está nessa situação na consulta pública")
	ErrExposicaoNaoExiste = errors.New("acervo: exposição não encontrada")
	ErrExposicaoVazia     = errors.New("acervo: a exposição precisa de ao menos um documento liberado à consulta pública")
	ErrExposicaoSituacao  = errors.New("acervo: a exposição já está nessa situação")
	ErrSlugEmUso          = errors.New("acervo: já existe uma exposição com esse endereço")
	ErrConsulta           = errors.New("acervo: consulta inválida")
)

// SerieSIC é a série da TTDD da SEMAD dos pedidos de acesso à informação:
// o pedido do que não é ostensivo segue pelo Serviço de Informação ao
// Cidadão.
const SerieSIC = "2.0.03.02"

// PodePublicarDigital: só o representante pronto de unidade ostensiva e não
// eliminada vai para a consulta pública. O banco confere de novo.
func PodePublicarDigital(u Unidade, d Digital) error {
	switch {
	case d.Situacao != DigitalPronto:
		return ErrDigitalSituacao
	case u.Fase == FaseEliminado:
		return ErrEliminado
	case u.NivelAcesso != AcessoPublico:
		return ErrNaoOstensivo
	}
	return nil
}

// Facetas da busca do acervo (artigo, item 4.4). O nível de acesso só
// existe na busca interna.
const (
	FacetaOrgao      = "orgao"
	FacetaClasse     = "classe"
	FacetaEspecie    = "especie"
	FacetaFase       = "fase"
	FacetaDestinacao = "destinacao"
	FacetaAno        = "ano"
	FacetaNivel      = "nivel_acesso"
)

// FacetasPublicas e FacetasInternas, na ordem em que a tela mostra.
var (
	FacetasPublicas = []string{FacetaOrgao, FacetaClasse, FacetaEspecie, FacetaFase, FacetaDestinacao, FacetaAno}
	FacetasInternas = append(append([]string{}, FacetasPublicas...), FacetaNivel)
)

// Consulta é uma busca no acervo.
type Consulta struct {
	Q string
	// Filtros: faceta -> valor escolhido.
	Filtros   map[string]string
	SoDigital bool
	Pagina    int
	PorPagina int
}

// FiltroBusca é o que a busca pode mostrar a quem pede: a coleção pública
// ou as unidades do escopo (Restrito) na interna.
type FiltroBusca struct {
	Publico  bool
	Restrito bool
	Unidades []uuid.UUID
}

// NormalizarConsulta confere os filtros (só facetas conhecidas; o ano é um
// número) e limita a página.
func NormalizarConsulta(c Consulta, publico bool, maxPorPagina int) (Consulta, error) {
	c.Q = strings.TrimSpace(c.Q)
	if utf8.RuneCountInString(c.Q) > 200 {
		return c, &ErrValidacao{Campo: "q", Motivo: "no máximo 200 caracteres"}
	}
	validas := FacetasInternas
	if publico {
		validas = FacetasPublicas
	}
	limpos := map[string]string{}
	for k, v := range c.Filtros {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if !contem(validas, k) {
			return c, &ErrValidacao{Campo: k, Motivo: "filtro desconhecido"}
		}
		if k == FacetaAno {
			if _, err := strconv.Atoi(v); err != nil {
				return c, &ErrValidacao{Campo: k, Motivo: "ano inválido"}
			}
		}
		limpos[k] = v
	}
	c.Filtros = limpos
	c.Pagina = max(c.Pagina, 1)
	if c.PorPagina <= 0 || c.PorPagina > maxPorPagina {
		c.PorPagina = min(20, maxPorPagina)
	}
	return c, nil
}

func contem(lista []string, v string) bool {
	for _, x := range lista {
		if x == v {
			return true
		}
	}
	return false
}

// Classe é o primeiro nível do código de classificação ("2.0.03.02" -> "2").
func Classe(codigo string) string {
	c, _, _ := strings.Cut(codigo, ".")
	return c
}

// anosMax limita os anos indexados de uma unidade com datas-limite muito
// largas (as pontas ficam sempre).
const anosMax = 200

// Anos são os anos cobertos pelas datas-limite (a faceta "ano" acha a
// unidade em qualquer ano do intervalo).
func Anos(inicio, fim time.Time) []int {
	a, b := inicio.Year(), fim.Year()
	if b-a+1 > anosMax {
		return []int{a, b}
	}
	out := make([]int, 0, b-a+1)
	for y := a; y <= b; y++ {
		out = append(out, y)
	}
	return out
}

// TextoMaxIndice limita o texto do OCR que vai para o índice por unidade.
const TextoMaxIndice = 200_000

// RegistroIndice é a unidade com o que a busca precisa: a série, o órgão e
// os representantes prontos (com o texto do OCR quando pedido).
type RegistroIndice struct {
	Unidade   Unidade
	Orgao     string
	Descritor string
	Digitais  []DigitalIndice
}

// DigitalIndice é o representante pronto na busca e na ficha pública.
type DigitalIndice struct {
	ID             uuid.UUID  `json:"id"`
	Numero         string     `json:"numero"`
	Descricao      string     `json:"descricao"`
	Paginas        int        `json:"paginas"`
	Bytes          int64      `json:"pdfa_bytes"`
	SHA256         string     `json:"pdfa_sha256"`
	Conformidade   string     `json:"conformidade"`
	DigitalizadoEm time.Time  `json:"digitalizado_em"`
	PublicadoEm    *time.Time `json:"publicado_em"`
	Texto          string     `json:"-"`
}

// Ostensivo: a unidade vai para a consulta pública.
func (r RegistroIndice) Ostensivo() bool {
	return r.Unidade.NivelAcesso == AcessoPublico && r.Unidade.Fase != FaseEliminado
}

// Liberados são os representantes da consulta pública.
func (r RegistroIndice) Liberados() []DigitalIndice {
	out := []DigitalIndice{}
	for _, d := range r.Digitais {
		if d.PublicadoEm != nil {
			out = append(out, d)
		}
	}
	return out
}

// TextoIndice junta o texto do OCR dos representantes dados, no limite.
func TextoIndice(ds []DigitalIndice) string {
	var b strings.Builder
	for _, d := range ds {
		if d.Texto == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(d.Texto)
		if b.Len() >= TextoMaxIndice {
			break
		}
	}
	s := b.String()
	if len(s) > TextoMaxIndice {
		s = strings.ToValidUTF8(s[:TextoMaxIndice], "")
	}
	return s
}

// PendenteIndice é uma unidade marcada para reindexar.
type PendenteIndice struct {
	UnidadeID uuid.UUID
	Desde     time.Time
}

// ItemBusca é uma unidade encontrada.
type ItemBusca struct {
	ID          uuid.UUID `json:"id"`
	Numero      string    `json:"numero"`
	Tipo        string    `json:"tipo"`
	Codigo      string    `json:"codigo_classificacao"`
	Descritor   string    `json:"descritor"`
	Assunto     string    `json:"assunto"`
	Orgao       string    `json:"orgao"`
	AnoInicial  int       `json:"ano_inicial"`
	AnoFinal    int       `json:"ano_final"`
	Fase        string    `json:"fase"`
	Destinacao  string    `json:"destinacao"`
	NivelAcesso string    `json:"nivel_acesso,omitempty"`
	Digitais    int       `json:"digitais"`
	// Trecho: onde o termo apareceu, com os termos entre <mark>.
	Trecho string `json:"trecho"`
}

// Faceta é um valor de faceta e quantas unidades o têm.
type Faceta struct {
	Valor string `json:"valor"`
	Total int    `json:"total"`
}

// Origens do resultado: o índice (Typesense) ou o banco, quando o índice
// está fora (sem facetas).
const (
	OrigemIndice = "indice"
	OrigemBanco  = "banco"
)

// ResultadoBusca é a resposta da busca.
type ResultadoBusca struct {
	Itens   []ItemBusca         `json:"items"`
	Total   int64               `json:"total"`
	Facetas map[string][]Faceta `json:"facetas"`
	Origem  string              `json:"origem"`
	Ms      int64               `json:"ms"`
}

// Tipos de consulta pública contados por dia (indicador da SEMAD).
const (
	ConsultaBusca     = "BUSCA"
	ConsultaFicha     = "FICHA"
	ConsultaDocumento = "DOCUMENTO"
	ConsultaExposicao = "EXPOSICAO"
)

// FichaPublica é a descrição ostensiva da unidade: sem localização,
// observações ou quem cadastrou.
type FichaPublica struct {
	ID            uuid.UUID       `json:"id"`
	Numero        string          `json:"numero"`
	Tipo          string          `json:"tipo"`
	Codigo        string          `json:"codigo_classificacao"`
	Descritor     string          `json:"descritor"`
	Assunto       string          `json:"assunto"`
	Orgao         string          `json:"orgao"`
	DataInicial   time.Time       `json:"data_inicial"`
	DataFinal     time.Time       `json:"data_final"`
	Quantidade    int             `json:"quantidade"`
	Especificacao string          `json:"especificacao"`
	Fase          string          `json:"fase"`
	Destinacao    string          `json:"destinacao"`
	Recolhida     bool            `json:"recolhida"`
	Digitais      []DigitalIndice `json:"digitais"`
}

// NovaFicha monta a ficha pública do registro.
func NovaFicha(r RegistroIndice) FichaPublica {
	u := r.Unidade
	return FichaPublica{ID: u.ID, Numero: u.Numero, Tipo: string(u.Tipo), Codigo: u.CodigoClassificacao, Descritor: r.Descritor,
		Assunto: u.Assunto, Orgao: r.Orgao, DataInicial: u.DataInicial, DataFinal: u.DataFinal, Quantidade: u.Quantidade,
		Especificacao: u.Especificacao, Fase: string(u.Fase), Destinacao: string(u.Destinacao), Recolhida: u.RecolhidaEm != nil,
		Digitais: r.Liberados()}
}

// DigitalRecente é um representante recém-liberado (página do Arquivo).
type DigitalRecente struct {
	DigitalIndice
	UnidadeID     uuid.UUID `json:"unidade_arquivamento_id"`
	UnidadeNumero string    `json:"unidade_numero"`
	Assunto       string    `json:"assunto"`
}

// Exposicao é uma exposição virtual do Arquivo Público (difusão, art. 5º,
// X e XI).
type Exposicao struct {
	ID           uuid.UUID       `json:"id"`
	Slug         string          `json:"slug"`
	Titulo       string          `json:"titulo"`
	Resumo       string          `json:"resumo"`
	Texto        string          `json:"texto"`
	PublicadaEm  *time.Time      `json:"publicada_em"`
	CriadoPor    *uuid.UUID      `json:"-"`
	CriadoEm     time.Time       `json:"criado_em"`
	AtualizadoEm time.Time       `json:"atualizado_em"`
	Itens        []ItemExposicao `json:"itens,omitempty"`

}

// ItemExposicao é um documento da exposição, com a legenda.
type ItemExposicao struct {
	DigitalID     uuid.UUID `json:"digital_id"`
	Ordem         int       `json:"ordem"`
	Legenda       string    `json:"legenda"`
	Numero        string    `json:"numero"`
	Descricao     string    `json:"descricao"`
	Paginas       int       `json:"paginas"`
	UnidadeID     uuid.UUID `json:"unidade_arquivamento_id"`
	UnidadeNumero string    `json:"unidade_numero"`
	Publicado     bool      `json:"publicado"`
}

// ItensMax limita os documentos de uma exposição.
const ItensMax = 100

// NovaExposicao valida o que se grava da exposição.
func NovaExposicao(x Exposicao) (Exposicao, error) {
	x.Titulo, x.Resumo, x.Texto = strings.TrimSpace(x.Titulo), strings.TrimSpace(x.Resumo), strings.TrimSpace(x.Texto)
	switch {
	case x.Titulo == "" || utf8.RuneCountInString(x.Titulo) > 200:
		return x, &ErrValidacao{Campo: "titulo", Motivo: "obrigatório, até 200 caracteres"}
	case utf8.RuneCountInString(x.Resumo) > 500:
		return x, &ErrValidacao{Campo: "resumo", Motivo: "no máximo 500 caracteres"}
	case utf8.RuneCountInString(x.Texto) > 20000:
		return x, &ErrValidacao{Campo: "texto", Motivo: "no máximo 20.000 caracteres"}
	case x.Slug == "":
		return x, &ErrValidacao{Campo: "titulo", Motivo: "precisa de letras ou números"}
	case len(x.Itens) > ItensMax:
		return x, &ErrValidacao{Campo: "itens", Motivo: "no máximo " + strconv.Itoa(ItensMax) + " documentos"}
	}
	vistos := map[uuid.UUID]bool{}
	for i := range x.Itens {
		it := &x.Itens[i]
		it.Legenda, it.Ordem = strings.TrimSpace(it.Legenda), i
		switch {
		case vistos[it.DigitalID]:
			return x, &ErrValidacao{Campo: "itens", Motivo: "documento repetido"}
		case utf8.RuneCountInString(it.Legenda) > 1000:
			return x, &ErrValidacao{Campo: "itens", Motivo: "legenda de até 1.000 caracteres"}
		}
		vistos[it.DigitalID] = true
	}
	return x, nil
}

// PodePublicarExposicao: a exposição publicada mostra ao menos um
// documento liberado.
func PodePublicarExposicao(x Exposicao) error {
	if x.PublicadaEm != nil {
		return ErrExposicaoSituacao
	}
	for _, it := range x.Itens {
		if it.Publicado {
			return nil
		}
	}
	return ErrExposicaoVazia
}

// SoPublicados deixa na exposição só os documentos liberados (a visita
// pública).
func (x Exposicao) SoPublicados() Exposicao {
	itens := []ItemExposicao{}
	for _, it := range x.Itens {
		if it.Publicado {
			itens = append(itens, it)
		}
	}
	x.Itens = itens
	return x
}
