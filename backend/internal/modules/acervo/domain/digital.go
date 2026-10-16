package domain

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Fase 5 (ADR 037): o representante digital com valor probatório (Decreto
// Federal 10.278/2020; Lei 10.301/2019, art. 48; artigo técnico, itens 4.1 e
// 4.3). O papel digitalizado vira PDF/A com OCR, metadados no XMP, hash
// SHA-256 e guarda imutável, verificada periodicamente. Digitalizar nunca
// libera a eliminação de documento permanente (arts. 20 e 42): a trava é a
// mesma de antes, PodeEliminar.

var (
	ErrDigitalNaoExiste = errors.New("acervo: representante digital não encontrado")
	ErrDigitalSituacao  = errors.New("acervo: o representante digital não está na situação que a operação exige")
)

// Situação do representante digital.
const (
	DigitalEnviado     = "ENVIADO"     // original recebido; aguarda o worker
	DigitalProcessando = "PROCESSANDO" // no pipeline
	DigitalPronto      = "PRONTO"      // PDF/A guardado, com hash e metadados
	DigitalReprovado   = "REPROVADO"   // resolução abaixo do mínimo: digitalizar de novo
	DigitalFalhou      = "FALHOU"      // o pipeline falhou nas tentativas; o Arquivo reprocessa
)

// Tipo do documento, que define a resolução mínima (artigo, item 4.1).
const (
	DocTexto  = "TEXTO"
	DocPlanta = "PLANTA" // plantas, mapas e imagens de alta densidade
)

// MaxTentativas antes de o representante ficar FALHOU.
const MaxTentativas = 3

// ConformidadePDFA é o padrão gerado (ISO 19005-2, nível B).
const ConformidadePDFA = "PDF/A-2b"

var reCPF = regexp.MustCompile(`^\d{11}$`)

// MascararCPF guarda só o meio do CPF (***.456.789-**): o Decreto pede o
// responsável identificado; a LGPD, o mínimo necessário.
func MascararCPF(cpf string) (string, error) {
	d := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, cpf)
	if d == "" {
		return "", nil
	}
	if !reCPF.MatchString(d) {
		return "", invalido("responsavel_cpf", "CPF com 11 dígitos")
	}
	return "***." + d[3:6] + "." + d[6:9] + "-**", nil
}

// Digital é o representante digital de (parte de) uma unidade de
// arquivamento.
type Digital struct {
	ID        uuid.UUID `json:"id"`
	Numero    string    `json:"numero"` // identificador único (RD-AAAA-NNNNNN)
	UnidadeID uuid.UUID `json:"unidade_arquivamento_id"`
	Descricao string    `json:"descricao"`
	TipoDoc   string    `json:"tipo_documento"`
	// Metadados do Decreto 10.278/2020 (art. 9º e Anexo II).
	Equipamento     string     `json:"equipamento"`
	DPIDeclarado    int        `json:"dpi_declarado"`
	ResponsavelID   uuid.UUID  `json:"responsavel_id"`
	ResponsavelNome string     `json:"responsavel_nome"`
	ResponsavelCPF  string     `json:"responsavel_cpf"` // mascarado
	CodigoClass     string     `json:"codigo_classificacao"`
	DigitalizadoEm  time.Time  `json:"digitalizado_em"` // UTC
	Situacao        string     `json:"situacao"`
	Motivo          string     `json:"motivo"`
	Tentativas      int        `json:"tentativas"`
	OriginalKey     string     `json:"-"`
	OriginalSHA256  string     `json:"original_sha256"`
	OriginalBytes   int64      `json:"original_bytes"`
	PDFAKey         string     `json:"-"`
	PDFASHA256      string     `json:"pdfa_sha256"`
	PDFABytes       int64      `json:"pdfa_bytes"`
	Paginas         int        `json:"paginas"`
	DPIMedido       int        `json:"dpi_medido"`
	Conformidade    string     `json:"conformidade"`
	RetencaoAte     *time.Time `json:"retencao_ate"`
	ProcessadoEm    *time.Time `json:"processado_em"`
	VerificadoEm    *time.Time `json:"verificado_em"`
	Integro         *bool      `json:"integro"`
	CriadoEm        time.Time  `json:"criado_em"`
	// PublicadoEm: liberado à consulta pública (Fase 6, ADR 038).
	PublicadoEm *time.Time `json:"publicado_em"`
	// Pades: estado da assinatura ICP-Brasil. Sem política municipal
	// (decisão 0.6) o valor é sempre AUSENTE — o JSON declara a ausência
	// para a tela nunca assumir certificado (ADR 040).
	Pades string `json:"pades"`
}

// PadesAusente: o PDF/A não tem assinatura ICP-Brasil. Não é falha do
// pipeline; a política ainda não existe.
const PadesAusente = "AUSENTE"

// ComEstadoPAdES preenche o campo declarado. Nada no disco muda.
func (d Digital) ComEstadoPAdES() Digital {
	if d.Pades == "" {
		d.Pades = PadesAusente
	}
	return d
}

// AfirmaAssinaturaICP é verdadeiro só com assinatura ICP-Brasil gravada.
// Enquanto a política 0.6 não existir, permanece falso.
func AfirmaAssinaturaICP(d Digital) bool {
	return d.Pades != "" && d.Pades != PadesAusente
}

// DPIMinimo devolve a resolução exigida para o tipo.
func DPIMinimo(tipo string, texto, planta int) int {
	if tipo == DocPlanta {
		return planta
	}
	return texto
}

// NovoDigital valida o registro do envio.
func NovoDigital(d Digital) (Digital, error) {
	d.Descricao, d.Equipamento = strings.TrimSpace(d.Descricao), strings.TrimSpace(d.Equipamento)
	d.TipoDoc = strings.ToUpper(strings.TrimSpace(d.TipoDoc))
	switch {
	case d.TipoDoc != DocTexto && d.TipoDoc != DocPlanta:
		return d, invalido("tipo_documento", "use TEXTO ou PLANTA")
	case utf8.RuneCountInString(d.Descricao) > 500:
		return d, invalido("descricao", "até 500 caracteres")
	case d.Equipamento == "" || utf8.RuneCountInString(d.Equipamento) > 200:
		return d, invalido("equipamento", "informe o scanner (modelo e número de série), até 200 caracteres")
	case d.DPIDeclarado < 72 || d.DPIDeclarado > 4800:
		return d, invalido("dpi_declarado", "entre 72 e 4800")
	}
	d.ID, d.Situacao = uuid.New(), DigitalEnviado
	return d, nil
}

// PodeDigitalizar: a unidade eliminada não tem mais papel a digitalizar.
func PodeDigitalizar(u Unidade) error {
	if u.Fase == FaseEliminado {
		return ErrEliminado
	}
	return nil
}

// Retencao: por quanto tempo o object lock segura o PDF/A. O permanente (e
// a série sem destinação definida) fica para sempre — 100 anos renováveis;
// o que vai à eliminação, até um ano depois do fim do prazo intermediário
// (o representante segue a TTDD do papel), ao menos um ano.
func Retencao(u Unidade, c Calculo, hoje time.Time) time.Time {
	if u.Fase == FasePermanente || u.Destinacao != DestinacaoEliminacao || c.FimIntermediaria == nil {
		return hoje.AddDate(100, 0, 0)
	}
	fim := c.FimIntermediaria.AddDate(1, 0, 0)
	if minimo := hoje.AddDate(1, 0, 0); fim.Before(minimo) {
		return minimo
	}
	return fim
}

// Metadados são os campos gravados no XMP do PDF/A (Decreto 10.278, art.
// 9º; artigo, item 4.1). O hash do XMP é o do original recebido: o do
// próprio PDF/A muda ao gravar o XMP e fica no banco.
type Metadados struct {
	Identificador   string
	ResponsavelNome string
	ResponsavelCPF  string
	DigitalizadoEm  time.Time
	Equipamento     string
	DPI             int
	HashOriginal    string
	Classificacao   string
	Titulo          string
}

// Resultado é o que o pipeline devolve do PDF/A gerado.
type Resultado struct {
	Paginas      int
	DPIMinimo    int // a menor resolução das imagens (0 = sem imagens)
	Texto        string
	Conformidade string
}

// Pipeline transforma o original em PDF/A com OCR e metadados (ocrmypdf,
// poppler, pikepdf). A porta é do domínio; o adaptador, da infraestrutura.
type Pipeline interface {
	// Medir devolve as páginas e a menor resolução do original.
	Medir(ctx context.Context, original string) (paginas, dpi int, err error)
	// Processar gera o PDF/A em saida, com o OCR e os metadados.
	Processar(ctx context.Context, original, saida string, m Metadados) (Resultado, error)
}

// MovDigitalizacao é a movimentação do representante digital na unidade.
const MovDigitalizacao = "DIGITALIZACAO"

// Verificacao é uma conferência de integridade (fixity) do PDF/A guardado.
type Verificacao struct {
	Em         time.Time `json:"em"`
	Integro    bool      `json:"integro"`
	SHA256Lido string    `json:"sha256_lido"`
	Erro       string    `json:"erro"`
}

// FiltroDigital filtra a lista de representantes.
type FiltroDigital struct {
	UnidadeID *uuid.UUID
	Situacao  string
	// SoFalhas: os que falharam na integridade ou no pipeline (painel do
	// Arquivo Público).
	SoFalhas bool
	// Publicados: só os liberados à consulta pública (Fase 6).
	Publicados bool
	Restrito   bool

	Unidades []uuid.UUID
}
