package application_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/storage/storagetest"
)

func TestHandlersDigital(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	s := e.digital(t, store, &pipeFalso{paginas: 1, dpi: 300})
	quem := digitalizador(t, e, e.unidade)
	h := router(s, quem)
	arq := router(s, arquivoPublico(uuid.Nil))
	u := e.criar(codPerm)
	base := "/acervo/unidades/" + u.ID.String() + "/digitais"

	if rec := do(h, http.MethodPost, base+"/upload", `{"filename":"nota.docx"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("extensão: %d %s", rec.Code, rec.Body)
	}
	rec := do(h, http.MethodPost, base+"/upload", `{"filename":"lote.pdf"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("upload: %d %s", rec.Code, rec.Body)
	}
	chave := campo(t, rec.Body.String(), "object_key")
	_ = store.Put(ctx, "app", chave, strings.NewReader(pdfOriginal), int64(len(pdfOriginal)), "application/pdf")
	rec = do(h, http.MethodPost, base, `{"object_key":"`+chave+`","tipo_documento":"TEXTO","equipamento":"Fujitsu","dpi_declarado":300,"responsavel_cpf":"12345678901"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("registro: %d %s", rec.Code, rec.Body)
	}
	id := campo(t, rec.Body.String(), "id")
	if _, err := s.ProcessarProximo(ctx); err != nil {
		t.Fatal(err)
	}
	d := "/acervo/digitais/" + id

	casos := []struct {
		h                  http.Handler
		method, path, body string
		want               int
	}{
		{h, http.MethodPost, "/acervo/unidades/x/digitais/upload", `{"filename":"a.pdf"}`, http.StatusBadRequest},
		{h, http.MethodPost, base + "/upload", `{`, http.StatusBadRequest},
		{h, http.MethodPost, "/acervo/unidades/" + uuid.NewString() + "/digitais/upload", `{"filename":"a.pdf"}`, http.StatusNotFound},
		{h, http.MethodPost, "/acervo/unidades/x/digitais", `{}`, http.StatusBadRequest},
		{h, http.MethodPost, base, `{`, http.StatusBadRequest},
		{h, http.MethodPost, base, `{"object_key":"k","tipo_documento":"FOTO","equipamento":"e","dpi_declarado":300}`, http.StatusUnprocessableEntity},
		{h, http.MethodPost, base, `{"object_key":"k","tipo_documento":"TEXTO","origem":"FOTO","equipamento":"e","dpi_declarado":300}`, http.StatusUnprocessableEntity},
		{h, http.MethodGet, "/acervo/digitais?unidade_id=" + u.ID.String(), "", http.StatusOK},
		{h, http.MethodGet, "/acervo/digitais?unidade_id=x", "", http.StatusBadRequest},
		{arq, http.MethodGet, "/acervo/digitais?falhas=true", "", http.StatusOK},
		{h, http.MethodGet, d, "", http.StatusOK},
		{h, http.MethodGet, "/acervo/digitais/x", "", http.StatusBadRequest},
		{h, http.MethodGet, "/acervo/digitais/" + uuid.NewString(), "", http.StatusNotFound},
		{h, http.MethodGet, "/acervo/digitais/x/arquivo", "", http.StatusBadRequest},
		{h, http.MethodGet, "/acervo/digitais/" + uuid.NewString() + "/arquivo", "", http.StatusNotFound},
		{arq, http.MethodPost, d + "/verificar", ``, http.StatusOK},
		{arq, http.MethodPost, "/acervo/digitais/x/verificar", ``, http.StatusBadRequest},
		{arq, http.MethodPost, d + "/reprocessar", ``, http.StatusConflict},
	}
	for _, c := range casos {
		if rec := do(c.h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}

	// O arquivo em fluxo: inteiro, por trecho, sufixo e fora do tamanho.
	ler := func(rng string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, d+"/arquivo", nil)
		if rng != "" {
			req.Header.Set("Range", rng)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}
	inteiro := ler("")
	tamanho := inteiro.Body.Len()
	if inteiro.Code != http.StatusOK || inteiro.Header().Get("Accept-Ranges") != "bytes" || !strings.HasPrefix(inteiro.Body.String(), "%PDF-1.7") ||
		inteiro.Header().Get("ETag") == "" || inteiro.Header().Get("Content-Type") != "application/pdf" {
		t.Fatalf("inteiro: %d %v", inteiro.Code, inteiro.Header())
	}
	parcial := ler("bytes=1-4")
	if parcial.Code != http.StatusPartialContent || parcial.Body.String() != "PDF-" || parcial.Header().Get("Content-Range") != "bytes 1-4/"+strconv.Itoa(tamanho) {
		t.Fatalf("trecho: %d %q %v", parcial.Code, parcial.Body, parcial.Header())
	}
	if r := ler("bytes=5-"); r.Code != http.StatusPartialContent || r.Body.Len() != tamanho-5 {
		t.Fatalf("até o fim: %d %d", r.Code, r.Body.Len())
	}
	if r := ler("bytes=0-99999"); r.Code != http.StatusPartialContent || r.Body.Len() != tamanho {
		t.Fatalf("fim além do tamanho: %d %d", r.Code, r.Body.Len())
	}
	if r := ler("bytes=-3"); r.Code != http.StatusPartialContent || r.Body.Len() != 3 {
		t.Fatalf("sufixo: %d %d", r.Code, r.Body.Len())
	}
	for _, ruim := range []string{"bytes=99999-", "bytes=4-1", "bytes=0-1,4-5", "linhas=0-1", "bytes=-0", "bytes=x-", "bytes=1-y", "bytes=-"} {
		if r := ler(ruim); r.Code != http.StatusRequestedRangeNotSatisfiable || r.Header().Get("Content-Range") != "bytes */"+strconv.Itoa(tamanho) {
			t.Errorf("%q: %d", ruim, r.Code)
		}
	}
	// Armazenamento fora no meio: erro antes de começar a resposta.
	store.FailGet = true
	if r := ler("bytes=0-1"); r.Code != http.StatusInternalServerError {
		t.Fatalf("armazenamento fora: %d", r.Code)
	}
	store.FailGet = false
	// Leitura interrompida depois do cabeçalho: só registra.
	store.FailRead = true
	if r := ler(""); r.Code != http.StatusOK {
		t.Fatalf("leitura interrompida: %d", r.Code)
	}
	store.FailRead = false
	down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), quem)
	if r := do(down, http.MethodGet, "/acervo/digitais", ""); r.Code != http.StatusInternalServerError {
		t.Fatalf("lista com o banco fora: %d", r.Code)
	}
}

func TestHandlersNatoDigitalETIFF(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.limparFila(t)
	store := storagetest.New()
	s := e.digital(t, store, &pipeFalso{paginas: 1, dpi: 0})
	quem := digitalizador(t, e, e.unidade)
	h := router(s, quem)
	u := e.criar(codPerm)
	base := "/acervo/unidades/" + u.ID.String() + "/digitais"
	rec := do(h, http.MethodPost, base+"/upload", `{"filename":"oficio.tiff"}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), ".tiff") {
		t.Fatalf("upload tiff: %d %s", rec.Code, rec.Body)
	}
	chave := campo(t, rec.Body.String(), "object_key")
	if err := store.Put(ctx, "app", chave, strings.NewReader("II*tiff"), 7, "image/tiff"); err != nil {
		t.Fatal(err)
	}
	rec = do(h, http.MethodPost, base, `{"object_key":"`+chave+`","tipo_documento":"TEXTO","origem":"NATO_DIGITAL","equipamento":"SEI","dpi_declarado":0}`)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"origem":"NATO_DIGITAL"`) {
		t.Fatalf("nato: %d %s", rec.Code, rec.Body)
	}
	if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatal(err)
	}
	id := campo(t, rec.Body.String(), "id")
	if rec := do(h, http.MethodGet, "/acervo/digitais/"+id, ""); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"situacao":"PRONTO"`) {
		t.Fatalf("pronto: %d %s", rec.Code, rec.Body)
	}
}
