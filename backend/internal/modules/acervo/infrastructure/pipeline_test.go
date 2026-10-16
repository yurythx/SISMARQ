package infrastructure

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
)

// fakes cria programas falsos (sh) que imitam o ocrmypdf, o poppler e o
// python do pipeline. O pipeline real roda na imagem do worker (ADR 037).
func fakes(t *testing.T, scripts map[string]string) Ferramentas {
	t.Helper()
	dir := t.TempDir()
	bin := func(nome string) string {
		p := filepath.Join(dir, nome)
		corpo, ok := scripts[nome]
		if !ok {
			corpo = "exit 0"
		}
		if err := os.WriteFile(p, []byte("#!/bin/sh\n"+corpo+"\n"), 0o755); err != nil { // #nosec G306 -- programa de teste
			t.Fatal(err)
		}
		return p
	}
	return Ferramentas{OCRmyPDF: bin("ocrmypdf"), PDFInfo: bin("pdfinfo"), PDFImages: bin("pdfimages"), PDFToText: bin("pdftotext"),
		Python: bin("python3"), Timeout: 10 * time.Second}
}

const pdfimagesSaida = `page   num  type   width height color comp bpc  enc interp  object ID x-ppi y-ppi size ratio
--------------------------------------------------------------------------------------------
   1     0 image    2480  3508  gray    1   8  jpeg   no         8  0   300   301  412K 4.8%
   2     1 image    2480  3508  gray    1   8  jpeg   no        13  0   298   300  401K 4.7%
`

func TestPipelineMedir(t *testing.T) {
	f := fakes(t, map[string]string{
		"pdfinfo":   `echo "Title: x"; echo "Pages:          2"`,
		"pdfimages": "cat <<'EOF'\n" + pdfimagesSaida + "EOF",
	})
	p := NovoPipeline(f)
	paginas, dpi, err := p.Medir(context.Background(), "original.pdf")
	if err != nil || paginas != 2 || dpi != 298 {
		t.Fatalf("medir: %d %d %v", paginas, dpi, err)
	}
	// Sem imagens (PDF nato-digital): DPI 0.
	f2 := fakes(t, map[string]string{"pdfinfo": `echo "Pages: 1"`, "pdfimages": `echo "page num type"`})
	if _, dpi, err := NovoPipeline(f2).Medir(context.Background(), "x.pdf"); err != nil || dpi != 0 {
		t.Fatalf("sem imagens: %d %v", dpi, err)
	}
	// Falhas: pdfinfo quebra, não diz as páginas, pdfimages quebra.
	for nome, s := range map[string]map[string]string{
		"pdfinfo falha":   {"pdfinfo": `echo "arquivo corrompido" >&2; exit 1`},
		"sem páginas":     {"pdfinfo": `echo "Title: x"`},
		"pdfimages falha": {"pdfinfo": `echo "Pages: 1"`, "pdfimages": `exit 3`},
	} {
		if _, _, err := NovoPipeline(fakes(t, s)).Medir(context.Background(), "x.pdf"); err == nil {
			t.Errorf("%s: sem erro", nome)
		}
	}
}

func TestPipelineProcessar(t *testing.T) {
	dir := t.TempDir()
	log := filepath.Join(dir, "args.log")
	f := fakes(t, map[string]string{
		// ocrmypdf: o último argumento é a saída.
		"ocrmypdf":  `echo "$@" >> ` + log + `; for a; do out="$a"; done; echo PDFA > "$out"`,
		"python3":   `echo "$4" >> ` + log + `; cp "$3" "$4"`,
		"pdftotext": `echo "  Ofício nº 12/2015  "`,
		"pdfinfo":   `echo "Pages: 3"`,
	})
	saida := filepath.Join(dir, "pdfa.pdf")
	res, err := NovoPipeline(f).Processar(context.Background(), filepath.Join(dir, "original.pdf"), saida, domain.Metadados{
		Identificador: "RD-2026-000001", DigitalizadoEm: time.Date(2026, 10, 16, 14, 30, 0, 0, time.UTC), DPI: 300, Titulo: "RD — Ofícios"})
	if err != nil || res.Paginas != 3 || res.Texto != "Ofício nº 12/2015" || res.Conformidade != domain.ConformidadePDFA {
		t.Fatalf("processar: %+v %v", res, err)
	}
	args, _ := os.ReadFile(log) // #nosec G304 -- arquivo do teste
	for _, quer := range []string{"--output-type pdfa-2", "--deskew", "--clean", "--skip-text", "-l por", saida} {
		if !strings.Contains(string(args), quer) {
			t.Errorf("argumentos sem %q: %s", quer, args)
		}
	}
	if b, _ := os.ReadFile(saida); string(b) != "PDFA\n" { // #nosec G304 -- arquivo do teste
		t.Fatalf("saída: %q", b)
	}
	for nome, s := range map[string]map[string]string{
		"ocrmypdf falha":  {"ocrmypdf": `echo "` + strings.Repeat("x", 600) + `" >&2; exit 2`},
		"xmp falha":       {"python3": `exit 1`},
		"pdftotext falha": {"python3": `exit 0`, "pdftotext": `exit 1`},
		"pdfinfo falha":   {"pdfinfo": `exit 1`},
	} {
		if _, err := NovoPipeline(fakes(t, s)).Processar(context.Background(), "o.pdf", filepath.Join(t.TempDir(), "s.pdf"), domain.Metadados{}); err == nil {
			t.Errorf("%s: sem erro", nome)
		} else if nome == "ocrmypdf falha" && len(err.Error()) > 700 {
			t.Errorf("mensagem não truncada: %d", len(err.Error()))
		}
	}
	// Padrões: nomes no PATH, OCR em português e 2 h por programa.
	if p := NovoPipeline(Ferramentas{}); p.f.OCRmyPDF != "ocrmypdf" || p.f.Idioma != "por" || p.f.Timeout != 2*time.Hour || p.f.Python != "python3" {
		t.Fatalf("padrões: %+v", p.f)
	}
}
