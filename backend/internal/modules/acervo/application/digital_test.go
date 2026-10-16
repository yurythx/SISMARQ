package application_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
	"github.com/yurythx/sismarq/internal/platform/storage"
	"github.com/yurythx/sismarq/internal/platform/storage/storagetest"
)

// wormMem é o armazenamento em memória com object lock (o MinIO real).
type wormMem struct {
	*storagetest.Memory
	falhaLock, falhaPut bool
	dias                map[string]int
}

func (w *wormMem) EnsureImmutableBucket(context.Context, string, int) error {
	if w.falhaLock {
		return storagetest.ErrInjected
	}
	return nil
}

func (w *wormMem) PutImmutable(ctx context.Context, bucket, obj string, r io.Reader, size int64, ct string, dias int) error {
	if w.falhaPut {
		return storagetest.ErrInjected
	}
	w.dias[obj] = dias
	return w.Memory.Put(ctx, bucket, obj, r, size, ct)
}

// semRange esconde o GetRange (armazenamento que só lê o arquivo inteiro).
type semRange struct{ storage.Provider }

// pipeFalso imita o ocrmypdf e o poppler.
type pipeFalso struct {
	paginas, dpi, paginasSaida int
	errMedir, errProc          error
	naoGrava, saidaPasta       bool // simula PDF/A ausente ou ilegível
	meta                       domain.Metadados
	texto                      string // o texto do OCR (padrão: um ofício)
}

func (p *pipeFalso) Medir(context.Context, string) (int, int, error) {
	return p.paginas, p.dpi, p.errMedir
}

func (p *pipeFalso) Processar(_ context.Context, _, saida string, m domain.Metadados) (domain.Resultado, error) {
	p.meta = m
	if p.errProc != nil {
		return domain.Resultado{}, p.errProc
	}
	switch {
	case p.saidaPasta:
		if err := os.Mkdir(saida, 0o700); err != nil {
			return domain.Resultado{}, err
		}
	case !p.naoGrava:
		if err := os.WriteFile(saida, []byte("%PDF-1.7 PDF/A de "+m.Identificador), 0o600); err != nil {
			return domain.Resultado{}, err
		}
	}
	pag := p.paginasSaida
	if pag == 0 {
		pag = p.paginas
	}
	texto := p.texto
	if texto == "" {
		texto = "Ofício nº 12/2015"
	}
	return domain.Resultado{Paginas: pag, Texto: texto, Conformidade: domain.ConformidadePDFA}, nil

}

const pdfOriginal = "%PDF-1.4 imagem digitalizada"

func (e *env) digital(t *testing.T, store storage.Provider, pipe domain.Pipeline) *application.Service {
	t.Helper()
	return e.real().WithDigitalizacao(application.ConfigDigital{Store: store, Bucket: "app", WORMBucket: "worm", MaxBytes: 1 << 20,
		URLExpiry: time.Minute, DPITexto: 300, DPIPlanta: 600, TmpDir: t.TempDir(), Pipeline: pipe})
}

// limparFila: os testes dividem o banco; nada de outro teste fica na fila.
func (e *env) limparFila(t *testing.T) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), `UPDATE acervo_digitais SET situacao = 'FALHOU' WHERE situacao IN ('ENVIADO', 'PROCESSANDO')`); err != nil {
		t.Fatal(err)
	}
}

// enviar faz o upload (direto ao armazenamento) e registra o PDF.
func enviar(t *testing.T, s *application.Service, store storage.Provider, quem auth.Identity, unidade uuid.UUID, tipo string) application.DetalheDigital {
	t.Helper()
	ctx := context.Background()
	tk, err := s.UploadDigital(ctx, quem, unidade, "lote.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Put(ctx, "app", tk.ObjectKey, strings.NewReader(pdfOriginal), int64(len(pdfOriginal)), "application/pdf"); err != nil {
		t.Fatal(err)
	}
	d, err := s.RegistrarDigital(ctx, quem, unidade, application.NovoEnvio{ObjectKey: tk.ObjectKey, TipoDoc: tipo, Equipamento: "Fujitsu fi-7160 S/N 1",
		DPI: 300, ResponsavelCPF: "123.456.789-01"})
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func digitalizador(t *testing.T, e *env, unidade uuid.UUID) auth.Identity {
	id := agente(unidade)
	id.UserID, id.Name = dbtest.User(t, e.pool), "Ana Digitalizadora"
	return id
}

// Do upload ao PDF/A guardado com object lock, com os metadados do Decreto
// 10.278/2020 e o original apagado.
func TestDigitalizacaoAprovada(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	pipe := &pipeFalso{paginas: 2, dpi: 300}
	s := e.digital(t, store, pipe)
	u := e.criar(codPerm)
	quem := digitalizador(t, e, e.unidade)

	if _, err := s.UploadDigital(ctx, agente(e.outra), u.ID, "a.pdf"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("fora do escopo: %v", err)
	}
	if _, err := s.UploadDigital(ctx, agente(e.unidade, "acervo:read"), u.ID, "a.pdf"); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só leitura: %v", err)
	}
	tk, err := s.UploadDigital(ctx, quem, u.ID, "lote.pdf")
	if err != nil || !strings.HasPrefix(tk.ObjectKey, "acervo/digitalizacao/"+u.ID.String()+"/") {
		t.Fatalf("ticket: %+v %v", tk, err)
	}
	novo := application.NovoEnvio{ObjectKey: tk.ObjectKey, TipoDoc: "TEXTO", Equipamento: "Fujitsu", DPI: 300}
	if _, err := s.RegistrarDigital(ctx, agente(e.unidade), u.ID, novo); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("sem usuário local: %v", err)
	}
	ruim := novo
	ruim.ResponsavelCPF = "123"
	if _, err := s.RegistrarDigital(ctx, quem, u.ID, ruim); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("CPF inválido: %v", err)
	}
	ruim = novo
	ruim.Equipamento = ""
	if _, err := s.RegistrarDigital(ctx, quem, u.ID, ruim); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("sem equipamento: %v", err)
	}
	if _, err := s.RegistrarDigital(ctx, quem, u.ID, novo); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("antes do upload: %v", err)
	}
	if _, err := s.RegistrarDigital(ctx, agente(e.outra), u.ID, novo); codigoDe(err) == "" {
		t.Fatalf("fora do escopo: %v", err)
	}

	d := enviar(t, s, store, quem, u.ID, "texto")
	if d.Situacao != domain.DigitalEnviado || !strings.HasPrefix(d.Numero, "RD-") || d.ResponsavelCPF != "***.456.789-**" ||
		d.ResponsavelNome != "Ana Digitalizadora" || d.CodigoClass != codPerm || d.OriginalBytes != int64(len(pdfOriginal)) || d.Unidade != u.Numero {
		t.Fatalf("registro: %+v", d)
	}

	if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatalf("processar: %v %v", ok, err)
	}
	pronto, err := s.Digital(ctx, quem, d.ID)
	if err != nil || pronto.Situacao != domain.DigitalPronto || pronto.Paginas != 2 || pronto.DPIMedido != 300 ||
		pronto.Conformidade != domain.ConformidadePDFA || pronto.RetencaoAte == nil || len(pronto.PDFASHA256) != 64 || pronto.ProcessadoEm == nil {
		t.Fatalf("pronto: %+v %v", pronto, err)
	}
	if pipe.meta.Identificador != d.Numero || pipe.meta.HashOriginal != d.OriginalSHA256 || pipe.meta.ResponsavelCPF != "***.456.789-**" ||
		pipe.meta.Classificacao != codPerm || !strings.Contains(pipe.meta.Titulo, d.Numero) {
		t.Fatalf("metadados do Decreto: %+v", pipe.meta)
	}
	chave := "digitais/" + u.ID.String() + "/" + d.Numero + ".pdf"
	if !store.Has("worm", chave) || store.Has("app", tk.ObjectKey) || store.dias[chave] < 36000 {
		t.Fatalf("guarda: worm=%v original=%v dias=%d", store.Has("worm", chave), store.Has("app", tk.ObjectKey), store.dias[chave])
	}
	var texto string
	if err := e.pool.QueryRow(ctx, `SELECT texto FROM acervo_digital_textos WHERE digital_id = $1`, d.ID).Scan(&texto); err != nil || texto != "Ofício nº 12/2015" {
		t.Fatalf("texto do OCR: %q %v", texto, err)
	}
	// O pronto não muda no banco (hash, objeto e retenção).
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET pdfa_sha256 = repeat('0', 64) WHERE id = $1`, d.ID); err == nil {
		t.Fatal("o banco alterou o hash do pronto")
	}
	if _, err := e.pool.Exec(ctx, `DELETE FROM acervo_digitais WHERE id = $1`, d.ID); err == nil {
		t.Fatal("o banco apagou o representante")
	}
	if ok, err := s.ProcessarProximo(ctx); ok || err != nil {
		t.Fatalf("fila vazia: %v %v", ok, err)
	}
	if det, _ := s.Get(ctx, quem, u.ID); det.Movimentacoes[0].Tipo != domain.MovDigitalizacao || !strings.Contains(det.Movimentacoes[0].Descricao, "pronto") {
		t.Fatalf("histórico da unidade: %+v", det.Movimentacoes[0])
	}

	// Leitura por trecho (Range) e inteira; sem Range no armazenamento.
	arq, err := s.AbrirArquivo(ctx, quem, d.ID)
	if err != nil {
		t.Fatal(err)
	}
	rc, err := s.LerArquivo(ctx, arq, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	trecho, _ := io.ReadAll(rc)
	if string(trecho) != "DF-1." {
		t.Fatalf("trecho: %q", trecho)
	}
	rc, _ = s.LerArquivo(ctx, arq, 0, 0)
	if inteiro, _ := io.ReadAll(rc); !bytes.HasPrefix(inteiro, []byte("%PDF-1.7 PDF/A de RD-")) {
		t.Fatalf("inteiro: %q", inteiro)
	}
	semTrecho := e.digital(t, semRange{store}, pipe)
	if _, err := semTrecho.LerArquivo(ctx, arq, 0, 1); !errors.Is(err, application.ErrSemRange) {
		t.Fatalf("sem Range: %v", err)
	}
	if _, err := s.AbrirArquivo(ctx, agente(e.outra), d.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("arquivo fora do escopo: %v", err)
	}
	if _, err := s.Digital(ctx, agente(e.outra), uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}

	// Listas: a unidade, o escopo, o Arquivo Público.
	if items, total, err := s.Digitais(ctx, quem, domain.FiltroDigital{UnidadeID: &u.ID}, pagination.New(1, 10, 100)); err != nil || total != 1 || items[0].ID != d.ID {
		t.Fatalf("da unidade: %+v %v", items, err)
	}
	if _, total, _ := s.Digitais(ctx, agente(e.outra), domain.FiltroDigital{}, pagination.New(1, 10, 100)); total != 0 {
		t.Fatal("outra unidade não vê")
	}
	if _, total, err := s.Digitais(ctx, arquivoPublico(uuid.Nil), domain.FiltroDigital{}, pagination.New(1, 10, 100)); err != nil || total < 1 {
		t.Fatalf("Arquivo Público: %d %v", total, err)
	}
}

// Resolução abaixo do mínimo, PDF sem imagens e planta: reprovado, com o
// original apagado (é preciso digitalizar de novo).
func TestDigitalizacaoReprovada(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := storagetest.New()
	pipe := &pipeFalso{paginas: 1}
	s := e.digital(t, store, pipe)
	quem := digitalizador(t, e, e.unidade)
	u := e.criar(codElim)
	casos := []struct {
		tipo   string
		dpi    int
		motivo string
	}{
		{"TEXTO", 200, "200 DPI abaixo do mínimo de 300"},
		{"PLANTA", 400, "400 DPI abaixo do mínimo de 600 DPI para planta"},
		{"TEXTO", 0, "não traz imagens digitalizadas"},
	}
	for _, c := range casos {
		pipe.dpi = c.dpi
		d := enviar(t, s, store, quem, u.ID, c.tipo)
		if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
			t.Fatal(err)
		}
		got, _ := s.Digital(ctx, quem, d.ID)
		if got.Situacao != domain.DigitalReprovado || !strings.Contains(got.Motivo, c.motivo) || got.ProcessadoEm == nil {
			t.Errorf("%s %d: %+v", c.tipo, c.dpi, got.Digital)
		}
	}
	// Sem object lock no armazenamento: grava comum (e avisa).
	pipe.dpi = 600
	d := enviar(t, s, store, quem, u.ID, "PLANTA")
	if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatal(err)
	}
	got, _ := s.Digital(ctx, quem, d.ID)
	if got.Situacao != domain.DigitalPronto || !store.Has("worm", "digitais/"+u.ID.String()+"/"+d.Numero+".pdf") {
		t.Fatalf("sem WORM: %+v", got.Digital)
	}
	// Representante de eliminação: retido até um ano depois do prazo.
	if got.RetencaoAte == nil || got.RetencaoAte.Year() < 2027 {
		t.Fatalf("retenção: %v", got.RetencaoAte)
	}
	if _, err := s.UploadDigital(ctx, quem, uuid.New(), "a.pdf"); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("unidade inexistente: %v", err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_unidades SET fase = 'ELIMINADO' WHERE id = $1`, u.ID); err == nil {
		if _, err := s.UploadDigital(ctx, quem, u.ID, "a.pdf"); codigoDe(err) != apperrors.CodeConflict {
			t.Fatalf("eliminada: %v", err)
		}
	}
}

// Falha no pipeline: tenta de novo e, na terceira, fica FALHOU com o
// original guardado; o Arquivo Público reprocessa.
func TestDigitalizacaoFalhaEReprocessa(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	pipe := &pipeFalso{paginas: 2, dpi: 300, errMedir: errors.New("ocrmypdf: exit 2")}
	s := e.digital(t, store, pipe)
	quem := digitalizador(t, e, e.unidade)
	u := e.criar(codPerm)
	d := enviar(t, s, store, quem, u.ID, "TEXTO")
	for i := 1; i <= domain.MaxTentativas; i++ {
		if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
			t.Fatal(err)
		}
		got, _ := s.Digital(ctx, quem, d.ID)
		quer := domain.DigitalEnviado
		if i == domain.MaxTentativas {
			quer = domain.DigitalFalhou
		}
		if got.Situacao != quer || got.Tentativas != i || !strings.Contains(got.Motivo, "ocrmypdf") {
			t.Fatalf("tentativa %d: %+v", i, got.Digital)
		}
	}
	if !store.Has("app", d.OriginalKey) {
		t.Fatal("o original fica para reprocessar")
	}
	if _, err := s.Reprocessar(ctx, quem, d.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público: %v", err)
	}
	arquivo := arquivoPublico(uuid.Nil)
	got, err := s.Reprocessar(ctx, arquivo, d.ID)
	if err != nil || got.Situacao != domain.DigitalEnviado || got.Tentativas != 0 {
		t.Fatalf("reprocessar: %+v %v", got, err)
	}
	if _, err := s.Reprocessar(ctx, arquivo, d.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("reprocessar o que não falhou: %v", err)
	}
	if _, err := s.Reprocessar(ctx, arquivo, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	if _, err := s.AbrirArquivo(ctx, quem, d.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("arquivo antes de pronto: %v", err)
	}

	// Cada etapa do pipeline que falha devolve à fila.
	criar := *application.CriarArquivo
	defer func() { *application.CriarArquivo = criar }()
	falhas := []func(){
		func() { pipe.errMedir, pipe.errProc = errors.New(strings.Repeat("x", 3000)), nil }, // mensagem longa: truncada
		func() { pipe.errMedir, pipe.errProc = nil, errors.New("pikepdf") },
		func() { pipe.errProc, pipe.paginasSaida = nil, 5 },
		func() { pipe.paginasSaida, store.falhaLock = 0, true },
		func() { store.falhaLock, store.falhaPut = false, true },
		func() { store.falhaPut, pipe.naoGrava = false, true },
		func() { pipe.naoGrava, pipe.saidaPasta = false, true },
		func() { pipe.saidaPasta, store.FailRead = false, true },
		func() {
			store.FailRead = false
			*application.CriarArquivo = func(string) (*os.File, error) { return nil, errors.New("disco cheio") }
		},
		func() { *application.CriarArquivo, store.FailGet = criar, true },
	}
	for i, prepara := range falhas {
		prepara()
		if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
			t.Fatal(err)
		}
		if got, _ := s.Digital(ctx, quem, d.ID); got.Situacao != domain.DigitalEnviado || got.Motivo == "" || len([]rune(got.Motivo)) > 2000 {
			t.Fatalf("falha %d: %+v", i, got.Digital)
		}
		if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET situacao = 'ENVIADO', tentativas = 0 WHERE id = $1`, d.ID); err != nil {
			t.Fatal(err)
		}
	}
	store.FailGet = false
	ruim := e.real().WithDigitalizacao(application.ConfigDigital{Store: store, Bucket: "app", WORMBucket: "worm", MaxBytes: 1 << 20,
		DPITexto: 300, DPIPlanta: 600, TmpDir: "/nao/existe", Pipeline: pipe})
	if ok, err := ruim.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Digital(ctx, quem, d.ID); got.Situacao != domain.DigitalEnviado || got.Motivo == "" {
		t.Fatalf("sem diretório temporário: %+v", got.Digital)
	}
	if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Digital(ctx, quem, d.ID); got.Situacao != domain.DigitalPronto {
		t.Fatalf("por fim, pronto: %+v", got.Digital)
	}
	// O original apagado com o armazenamento fora: só avisa.
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET situacao = 'FALHOU' WHERE situacao IN ('ENVIADO', 'PROCESSANDO')`); err != nil {
		t.Fatal(err)
	}
	outro := enviar(t, s, store, quem, u.ID, "TEXTO")
	store.FailDelete = true
	if ok, err := s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatal(err)
	}
	store.FailDelete = false
	if got, _ := s.Digital(ctx, quem, outro.ID); got.Situacao != domain.DigitalPronto {
		t.Fatalf("pronto mesmo sem apagar o original: %+v", got.Digital)
	}
	// O registro falha se o original não pode ser lido.
	tk, _ := s.UploadDigital(ctx, quem, u.ID, "b.pdf")
	_ = store.Put(ctx, "app", tk.ObjectKey, strings.NewReader(pdfOriginal), int64(len(pdfOriginal)), "application/pdf")
	novo := application.NovoEnvio{ObjectKey: tk.ObjectKey, TipoDoc: "TEXTO", Equipamento: "e", DPI: 300}
	for _, liga := range []*bool{&store.FailRead, &store.FailGet} {
		*liga = true
		if _, err := s.RegistrarDigital(ctx, quem, u.ID, novo); err == nil {
			t.Fatal("original ilegível")
		}
		*liga = false
	}
}

// A verificação periódica acha o arquivo adulterado e o que sumiu, e avisa
// a equipe do Arquivo Público.
func TestIntegridade(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	s := e.digital(t, store, &pipeFalso{paginas: 1, dpi: 300})
	quem := digitalizador(t, e, e.unidade)
	u := e.criar(codPerm)
	equipe := dbtest.User(t, e.pool)
	var perfil uuid.UUID
	if err := e.pool.QueryRow(ctx, `INSERT INTO perfis (slug, nome, permissoes) VALUES ($1, $1, ARRAY['acervo:central']) RETURNING id`,
		"central-"+uuid.NewString()[:8]).Scan(&perfil); err != nil {
		t.Fatal(err)
	}
	if _, err := e.pool.Exec(ctx, `INSERT INTO user_scopes (user_id, perfil_id) VALUES ($1, $2)`, equipe, perfil); err != nil {
		t.Fatal(err)
	}
	a, b := enviar(t, s, store, quem, u.ID, "TEXTO"), enviar(t, s, store, quem, u.ID, "TEXTO")
	for range 2 {
		if _, err := s.ProcessarProximo(ctx); err != nil {
			t.Fatal(err)
		}
	}
	arquivo := arquivoPublico(uuid.Nil)
	if got, err := s.VerificarAgora(ctx, arquivo, a.ID); err != nil || got.Integro == nil || !*got.Integro || len(got.Verificacoes) != 1 {
		t.Fatalf("íntegro: %+v %v", got, err)
	}
	// Adulterado: o hash não bate; sumiu: erro de leitura.
	store.Corromper("worm", "digitais/"+u.ID.String()+"/"+a.Numero+".pdf", []byte("adulterado"))
	if err := store.Delete(ctx, "worm", "digitais/"+u.ID.String()+"/"+b.Numero+".pdf"); err != nil {
		t.Fatal(err)
	}
	e.agora = e.agora.Add(25 * time.Hour)
	n, err := s.VerificarIntegridade(ctx)
	if err != nil || n < 2 {
		t.Fatalf("verificação periódica: %d %v", n, err)
	}
	ga, _ := s.Digital(ctx, arquivo, a.ID)
	gb, _ := s.Digital(ctx, arquivo, b.ID)
	if *ga.Integro || ga.Verificacoes[0].SHA256Lido == a.PDFASHA256 || *gb.Integro || gb.Verificacoes[0].Erro == "" {
		t.Fatalf("falhas: %+v / %+v", ga.Verificacoes, gb.Verificacoes)
	}
	var avisos int
	if err := e.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notificacoes WHERE user_id = $1 AND modulo = 'acervo' AND link IN ($2, $3)`, equipe,
		"/acervo/digitais/"+a.ID.String(), "/acervo/digitais/"+b.ID.String()).Scan(&avisos); err != nil || avisos != 2 {
		t.Fatalf("avisos à equipe: %d %v", avisos, err)
	}
	if items, _, err := s.Digitais(ctx, arquivo, domain.FiltroDigital{UnidadeID: &u.ID, SoFalhas: true}, pagination.New(1, 10, 100)); err != nil || len(items) != 2 {
		t.Fatalf("painel de falhas: %d %v", len(items), err)
	}
	// Na hora: só o Arquivo Público, só o pronto.
	if _, err := s.VerificarAgora(ctx, quem, a.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("verificar fora do Arquivo: %v", err)
	}
	c := enviar(t, s, store, quem, u.ID, "TEXTO")
	if _, err := s.VerificarAgora(ctx, arquivo, c.ID); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("verificar o não pronto: %v", err)
	}
	if _, err := s.VerificarAgora(ctx, arquivo, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}

	// Os workers rodam, registram e param com o contexto.
	cctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	if err := s.RunDigitalizacao(10 * time.Millisecond)(cctx); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Digital(ctx, quem, c.ID); got.Situacao != domain.DigitalPronto {
		t.Fatalf("o worker processou a fila: %+v", got.Digital)
	}
	e.agora = e.agora.Add(25 * time.Hour)
	cctx2, cancel2 := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel2()
	if err := s.RunIntegridade(10 * time.Millisecond)(cctx2); err != nil {
		t.Fatal(err)
	}
}

// O banco fora em cada passo do worker e da verificação: o erro sobe (a
// falha do pipeline em si vira nova tentativa, testada acima).
func TestDigitalBancoFora(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2026-10-16")
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	quem := digitalizador(t, e, e.unidade)
	u := e.criar(codPerm)
	falho := func(n int, pipe domain.Pipeline) *application.Service {
		return e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: n}).WithDigitalizacao(application.ConfigDigital{Store: store,
			Bucket: "app", WORMBucket: "worm", MaxBytes: 1 << 20, DPITexto: 300, DPIPlanta: 600, TmpDir: t.TempDir(), Pipeline: pipe})
	}
	real := e.digital(t, store, &pipeFalso{paginas: 1, dpi: 300})
	// Reprovado (sem leitura da unidade no pipeline): reserva, conclusão,
	// unidade, histórico.
	for n := 1; n <= 4; n++ {
		e.limparFila(t)
		enviar(t, real, store, quem, u.ID, "TEXTO")
		if _, err := falho(n, &pipeFalso{paginas: 1, dpi: 100}).ProcessarProximo(ctx); err == nil {
			t.Errorf("falha na chamada %d engolida", n)
		}
	}
	// Aprovado: a leitura da unidade no pipeline falha e volta à fila.
	e.limparFila(t)
	b := enviar(t, real, store, quem, u.ID, "TEXTO")
	if ok, err := falho(2, &pipeFalso{paginas: 1, dpi: 300}).ProcessarProximo(ctx); !ok || err != nil {
		t.Fatalf("unidade ilegível no pipeline: %v %v", ok, err)
	}
	if got, _ := real.Digital(ctx, quem, b.ID); got.Situacao != domain.DigitalEnviado || got.Motivo == "" {
		t.Fatalf("de volta à fila: %+v", got.Digital)
	}
	e.limparFila(t)
	a := enviar(t, real, store, quem, u.ID, "TEXTO")
	if _, err := real.ProcessarProximo(ctx); err != nil {
		t.Fatal(err)
	}
	store.Corromper("worm", "digitais/"+u.ID.String()+"/"+a.Numero+".pdf", []byte("x"))
	e.agora = e.agora.Add(25 * time.Hour)
	// Verificação: lista, registro, equipe central; o evento e o aviso
	// falham com a transação envenenada.
	for n := 1; n <= 3; n++ {
		if _, err := falho(n, nil).VerificarIntegridade(ctx); err == nil {
			t.Errorf("verificação: falha na chamada %d engolida", n)
		}
	}
	for n := 2; n <= 3; n++ {
		envenenado := e.svc(&faultRepo{inner: infrastructure.NewRepository(), poisonAt: n}).WithDigitalizacao(application.ConfigDigital{Store: store,
			Bucket: "app", WORMBucket: "worm", MaxBytes: 1 << 20, TmpDir: t.TempDir()})
		if _, err := envenenado.VerificarIntegridade(ctx); err == nil {
			t.Errorf("verificação envenenada em %d engolida", n)
		}
	}
	// Leitura do PDF/A interrompida: registrada como erro.
	store.FailRead = true
	if v, err := real.VerificarAgora(ctx, arquivoPublico(uuid.Nil), a.ID); err != nil || v.Verificacoes[0].Erro == "" {
		t.Fatalf("leitura interrompida: %+v %v", v.Verificacoes, err)
	}
	store.FailRead = false
	cctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	if err := falho(1, nil).RunDigitalizacao(10 * time.Millisecond)(cctx); err != nil {
		t.Fatal(err)
	}
	cctx2, cancel2 := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel2()
	if err := falho(1, nil).RunIntegridade(10 * time.Millisecond)(cctx2); err != nil {
		t.Fatal(err)
	}
}
