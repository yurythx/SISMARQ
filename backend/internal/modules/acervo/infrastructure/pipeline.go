package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
)

// Ferramentas são os programas do pipeline (artigo técnico, item 4.3),
// instalados na imagem do worker (Dockerfile.worker).
type Ferramentas struct {
	OCRmyPDF  string
	PDFInfo   string
	PDFImages string
	PDFToText string
	Python    string
	// Timeout de cada programa (o OCR de centenas de páginas demora).
	Timeout time.Duration
	// Idioma do OCR (tesseract).
	Idioma string
}

func (f Ferramentas) comPadroes() Ferramentas {
	padrao := func(v *string, d string) {
		if *v == "" {
			*v = d
		}
	}
	padrao(&f.OCRmyPDF, "ocrmypdf")
	padrao(&f.PDFInfo, "pdfinfo")
	padrao(&f.PDFImages, "pdfimages")
	padrao(&f.PDFToText, "pdftotext")
	padrao(&f.Python, "python3")
	padrao(&f.Idioma, "por")
	if f.Timeout <= 0 {
		f.Timeout = 2 * time.Hour
	}
	return f
}

// PipelineExec é o domain.Pipeline com os programas do sistema.
type PipelineExec struct{ f Ferramentas }

// NovoPipeline cria o pipeline (os programas ausentes da config usam o
// nome padrão, procurado no PATH).
func NovoPipeline(f Ferramentas) *PipelineExec { return &PipelineExec{f: f.comPadroes()} }

func (p *PipelineExec) rodar(ctx context.Context, bin string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, p.f.Timeout)
	defer cancel()
	var out, errb bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, args...) // #nosec G204 -- programas da configuração, argumentos montados aqui
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if len(msg) > 500 {
			msg = msg[len(msg)-500:]
		}
		return nil, fmt.Errorf("%s: %w: %s", bin, err, msg)
	}
	return out.Bytes(), nil
}

var rePaginas = regexp.MustCompile(`(?m)^Pages:\s+(\d+)`)

func (p *PipelineExec) paginas(ctx context.Context, arquivo string) (int, error) {
	out, err := p.rodar(ctx, p.f.PDFInfo, arquivo)
	if err != nil {
		return 0, err
	}
	m := rePaginas.FindSubmatch(out)
	if m == nil {
		return 0, fmt.Errorf("pdfinfo: número de páginas não encontrado")
	}
	return strconv.Atoi(string(m[1]))
}

// Medir conta as páginas e acha a menor resolução das imagens do original
// (pdfimages -list: colunas x-ppi e y-ppi).
func (p *PipelineExec) Medir(ctx context.Context, original string) (int, int, error) {
	paginas, err := p.paginas(ctx, original)
	if err != nil {
		return 0, 0, err
	}
	out, err := p.rodar(ctx, p.f.PDFImages, "-list", original)
	if err != nil {
		return 0, 0, err
	}
	menor := 0
	for _, linha := range strings.Split(string(out), "\n") {
		c := strings.Fields(linha)
		if len(c) < 14 {
			continue
		}
		x, errX := strconv.Atoi(c[12])
		y, errY := strconv.Atoi(c[13])
		if errX != nil || errY != nil {
			continue // cabeçalho
		}
		if d := min(x, y); menor == 0 || d < menor {
			menor = d
		}
	}
	return paginas, menor, nil
}

// Processar: OCR em português com alinhamento e limpeza, saída PDF/A-2b
// (ocrmypdf); metadados do Decreto no XMP (pikepdf); texto (pdftotext).
func (p *PipelineExec) Processar(ctx context.Context, original, saida string, m domain.Metadados) (domain.Resultado, error) {
	intermediario := saida + ".ocr.pdf"
	if _, err := p.rodar(ctx, p.f.OCRmyPDF, "--output-type", "pdfa-2", "--deskew", "--clean", "--skip-text",
		"-l", p.f.Idioma, "--title", m.Titulo, original, intermediario); err != nil {
		return domain.Resultado{}, err
	}
	meta, _ := json.Marshal(map[string]string{
		"IdentificadorUnico":          m.Identificador,
		"ResponsavelDigitalizacao":    m.ResponsavelNome,
		"ResponsavelDigitalizacaoCPF": m.ResponsavelCPF,
		"DataHoraDigitalizacaoUTC":    m.DigitalizadoEm.UTC().Format(time.RFC3339),
		"EquipamentoScanner":          m.Equipamento,
		"ResolucaoDPI":                strconv.Itoa(m.DPI),
		"HashSHA256Original":          m.HashOriginal,
		"CodigoClassificacao":         m.Classificacao,
	})
	if _, err := p.rodar(ctx, p.f.Python, "-c", scriptXMP, intermediario, saida, string(meta)); err != nil {
		return domain.Resultado{}, err
	}
	texto, err := p.rodar(ctx, p.f.PDFToText, "-enc", "UTF-8", saida, "-")
	if err != nil {
		return domain.Resultado{}, err
	}
	paginas, err := p.paginas(ctx, saida)
	if err != nil {
		return domain.Resultado{}, err
	}
	return domain.Resultado{Paginas: paginas, Texto: strings.TrimSpace(string(texto)), Conformidade: domain.ConformidadePDFA}, nil
}

// scriptXMP grava os metadados do Decreto 10.278/2020 no XMP do PDF/A, no
// namespace do artigo técnico (item 4.1), e declara o schema de extensão
// (pdfaExtension) que o PDF/A exige para propriedades fora do padrão — sem
// ele o arquivo deixa de ser PDF/A válido. Salva linearizado (o
// visualizador abre a primeira página sem baixar o arquivo inteiro).
const scriptXMP = `
import json, sys
import xml.etree.ElementTree as ET
import pikepdf

src, dst, meta = sys.argv[1], sys.argv[2], json.loads(sys.argv[3])
X = "adobe:ns:meta/"
RDF = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
GOV = "http://www.governodigital.gov.br/decreto10278/"
EXT = "http://www.aiim.org/pdfa/ns/extension/"
SCH = "http://www.aiim.org/pdfa/ns/schema#"
PRO = "http://www.aiim.org/pdfa/ns/property#"
for p, u in (("x", X), ("rdf", RDF), ("govbr", GOV), ("pdfaExtension", EXT), ("pdfaSchema", SCH), ("pdfaProperty", PRO)):
    ET.register_namespace(p, u)

def q(ns, n):
    return "{%s}%s" % (ns, n)

with pikepdf.open(src) as pdf:
    atual = pdf.Root.get("/Metadata")
    if atual is not None:
        raiz = ET.fromstring(atual.read_bytes())
        for prefixo, uri in (("pdf", "http://ns.adobe.com/pdf/1.3/"), ("xmp", "http://ns.adobe.com/xap/1.0/"),
                             ("dc", "http://purl.org/dc/elements/1.1/"), ("pdfaid", "http://www.aiim.org/pdfa/ns/id/")):
            ET.register_namespace(prefixo, uri)
    else:
        raiz = ET.Element(q(X, "xmpmeta"))
    rdf = raiz.find(q(RDF, "RDF"))
    if rdf is None:
        rdf = ET.SubElement(raiz, q(RDF, "RDF"))

    desc = ET.SubElement(rdf, q(RDF, "Description"), {q(RDF, "about"): ""})
    for k, v in meta.items():
        ET.SubElement(desc, q(GOV, k)).text = v

    bag = rdf.find(".//" + q(EXT, "schemas") + "/" + q(RDF, "Bag"))
    if bag is None:
        ext = ET.SubElement(rdf, q(RDF, "Description"), {q(RDF, "about"): ""})
        bag = ET.SubElement(ET.SubElement(ext, q(EXT, "schemas")), q(RDF, "Bag"))
    li = ET.SubElement(bag, q(RDF, "li"), {q(RDF, "parseType"): "Resource"})
    ET.SubElement(li, q(SCH, "schema")).text = "Metadados do Decreto Federal 10.278/2020 (SISMARQ)"
    ET.SubElement(li, q(SCH, "namespaceURI")).text = GOV
    ET.SubElement(li, q(SCH, "prefix")).text = "govbr"
    seq = ET.SubElement(ET.SubElement(li, q(SCH, "property")), q(RDF, "Seq"))
    for k in meta:
        p = ET.SubElement(seq, q(RDF, "li"), {q(RDF, "parseType"): "Resource"})
        ET.SubElement(p, q(PRO, "name")).text = k
        ET.SubElement(p, q(PRO, "valueType")).text = "Text"
        ET.SubElement(p, q(PRO, "category")).text = "external"
        ET.SubElement(p, q(PRO, "description")).text = k

    corpo = ET.tostring(raiz, encoding="unicode")
    xmp = '<?xpacket begin="\ufeff" id="W5M0MpCehiHzreSzNTczkc9d"?>' + corpo + '<?xpacket end="w"?>'
    stream = pikepdf.Stream(pdf, xmp.encode("utf-8"))
    stream["/Type"] = pikepdf.Name("/Metadata")
    stream["/Subtype"] = pikepdf.Name("/XML")
    pdf.Root["/Metadata"] = stream
    pdf.save(dst, linearize=True)
`
