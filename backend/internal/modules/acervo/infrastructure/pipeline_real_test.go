package infrastructure

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
)

// TestPipelineReal roda o pipeline com os programas de verdade — na imagem
// do worker (Dockerfile.worker), com ACERVO_PIPELINE_REAL=1. Gera uma
// página "escaneada" a 300 DPI, passa pelo ocrmypdf e confere o PDF/A, o
// OCR em português e os metadados do Decreto 10.278 no XMP.
func TestPipelineReal(t *testing.T) {
	if os.Getenv("ACERVO_PIPELINE_REAL") != "1" {
		t.Skip("só na imagem do worker (ACERVO_PIPELINE_REAL=1)")
	}
	ctx := context.Background()
	dir := t.TempDir()
	original := filepath.Join(dir, "escaneado.pdf")
	gerar := `
import sys
from PIL import Image, ImageDraw, ImageFont
img = Image.new("L", (2480, 3508), 255)
d = ImageDraw.Draw(img)
try:
    f = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 90)
except OSError:
    f = ImageFont.load_default()
d.text((200, 300), "PREFEITURA DE RONDONOPOLIS", fill=0, font=f)
d.text((200, 500), "Oficio numero 12 de 2015", fill=0, font=f)
img = img.rotate(1.5, fillcolor=255)
img.save(sys.argv[1], "PDF", resolution=300)
`
	if out, err := exec.Command("python3", "-c", gerar, original).CombinedOutput(); err != nil { // #nosec G204 -- teste
		t.Fatalf("gerar o escaneado: %v %s", err, out)
	}
	p := NovoPipeline(Ferramentas{Timeout: 5 * time.Minute})
	paginas, dpi, err := p.Medir(ctx, original)
	if err != nil || paginas != 1 || dpi != 300 {
		t.Fatalf("medir: %d páginas, %d DPI, %v", paginas, dpi, err)
	}
	saida := filepath.Join(dir, "pdfa.pdf")
	inicio := time.Now()
	res, err := p.Processar(ctx, original, saida, domain.Metadados{Identificador: "RD-2026-000001", ResponsavelNome: "Ana",
		ResponsavelCPF: "***.456.789-**", DigitalizadoEm: time.Date(2026, 10, 16, 14, 30, 0, 0, time.UTC), Equipamento: "Fujitsu fi-7160",
		DPI: dpi, HashOriginal: strings.Repeat("a", 64), Classificacao: "2.0.01.00.02", Titulo: "RD-2026-000001 — Ofícios"})
	if err != nil {
		t.Fatalf("processar: %v", err)
	}
	t.Logf("pipeline: %s, %d página(s), texto: %q", time.Since(inicio), res.Paginas, res.Texto)
	if res.Paginas != 1 || !strings.Contains(strings.ToUpper(res.Texto), "RONDONOPOLIS") {
		t.Fatalf("resultado: %+v", res)
	}
	// O XMP tem a identificação PDF/A-2b, os metadados do Decreto e o
	// schema de extensão.
	ler := `
import sys, pikepdf
with pikepdf.open(sys.argv[1]) as pdf:
    print(pdf.Root.Metadata.read_bytes().decode("utf-8"))
    print("LINEARIZADO" if pdf.is_linearized else "NAO_LINEARIZADO")
`
	out, err := exec.Command("python3", "-c", ler, saida).CombinedOutput() // #nosec G204 -- teste
	if err != nil {
		t.Fatalf("ler o XMP: %v %s", err, out)
	}
	xmp := string(out)
	for _, quer := range []string{`pdfaid:part="2"`, `pdfaid:conformance="B"`, "govbr:IdentificadorUnico>RD-2026-000001<",
		"govbr:ResponsavelDigitalizacaoCPF>***.456.789-**<", "govbr:DataHoraDigitalizacaoUTC>2026-10-16T14:30:00Z<",
		"govbr:ResolucaoDPI>300<", "govbr:CodigoClassificacao>2.0.01.00.02<", "pdfaSchema:namespaceURI>http://www.governodigital.gov.br/decreto10278/<",
		"LINEARIZADO"} {
		if !strings.Contains(xmp, quer) {
			t.Errorf("XMP sem %q", quer)
		}
	}
	if t.Failed() {
		t.Logf("XMP:\n%s", xmp)
	}
}
