package application_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/storage"
	"github.com/yurythx/sismarq/internal/platform/storage/storagetest"
)

// acesso prepara o serviço com a digitalização (falsa) e a busca no banco.
type acesso struct {
	e     *env
	s     *application.Service
	store *wormMem
	pipe  *pipeFalso
	quem  auth.Identity
}

func novoAcesso(t *testing.T) *acesso {
	e := novoEnv(t)
	e.agora = dia("2026-10-16")
	e.limparFila(t)
	store := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	pipe := &pipeFalso{paginas: 2, dpi: 300}
	s := e.digital(t, store, pipe).WithBusca(application.ConfigBusca{SICURL: "https://sic.exemplo.gov.br"})
	return &acesso{e: e, s: s, store: store, pipe: pipe, quem: digitalizador(t, e, e.unidade)}
}

// pronta cadastra a unidade permanente com o assunto dado e digitaliza um
// documento (pronto).
func (a *acesso) pronta(t *testing.T, assunto string) (application.Detalhe, application.DetalheDigital) {
	t.Helper()
	ctx := context.Background()
	n := a.e.nova(codPerm)
	n.Assunto = assunto
	u, err := a.e.real().Criar(ctx, admin, n)
	if err != nil {
		t.Fatal(err)
	}
	d := enviar(t, a.s, a.store, a.quem, u.ID, "TEXTO")
	if ok, err := a.s.ProcessarProximo(ctx); !ok || err != nil {
		t.Fatalf("processar: %v %v", ok, err)
	}
	return u, d
}

func (a *acesso) nivel(t *testing.T, u application.Detalhe, nivel domain.NivelAcesso, hipotese string) {
	t.Helper()
	novo := u.Unidade
	novo.NivelAcesso, novo.HipoteseLegal = nivel, hipotese
	if _, err := a.s.Alterar(context.Background(), admin, u.ID, novo); err != nil {
		t.Fatal(err)
	}
}

func consultas(t *testing.T, e *env, tipo string) int64 {
	t.Helper()
	var n int64
	_ = e.pool.QueryRow(context.Background(), `SELECT COALESCE(SUM(total), 0) FROM acervo_consultas WHERE dia = $1 AND tipo = $2`, e.agora, tipo).Scan(&n)
	return n
}

// O cidadão acha, sem login, o permanente digitalizado e liberado; o
// restrito some da consulta pública e só aparece para quem tem o escopo.
func TestConsultaPublica(t *testing.T) {
	a := novoAcesso(t)
	e, s := a.e, a.s
	ctx := context.Background()
	marca := "Mercado" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	ocr := "Tombamento" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	a.pipe.texto = "Ofício do " + ocr
	u, d := a.pronta(t, "Planta do "+marca)
	a.pipe.texto = ""
	var orgao string
	if err := e.pool.QueryRow(ctx, `SELECT nome FROM unidades WHERE id = $1`, e.unidade).Scan(&orgao); err != nil {
		t.Fatal(err)
	}

	// Liberar: quem gere a unidade; o resto não.
	if _, err := s.PublicarDigital(ctx, agente(e.outra), d.ID, true); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("fora do escopo: %v", err)
	}
	if _, err := s.PublicarDigital(ctx, agente(e.unidade, "acervo:read"), d.ID, true); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só leitura: %v", err)
	}
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, false); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("retirar o que não está: %v", err)
	}
	if _, err := s.PublicarDigital(ctx, a.quem, uuid.New(), true); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	pub, err := s.PublicarDigital(ctx, a.quem, d.ID, true)
	if err != nil || pub.PublicadoEm == nil {
		t.Fatalf("liberar: %+v %v", pub, err)
	}
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("liberar de novo: %v", err)
	}

	// Busca sem login (o índice fora: o banco), pela descrição e pelo OCR.
	res, err := s.BuscarPublico(ctx, domain.Consulta{Q: marca})
	if err != nil || res.Total != 1 || res.Origem != domain.OrigemBanco || res.Itens[0].ID != u.ID || res.Itens[0].Digitais != 1 ||
		res.Itens[0].NivelAcesso != "" || res.Itens[0].Descritor != "Leis e decretos" || res.Itens[0].Orgao != orgao {
		t.Fatalf("busca pública: %+v %v", res, err)
	}
	filtros := map[string]string{domain.FacetaOrgao: orgao, domain.FacetaClasse: "9", domain.FacetaEspecie: "CAIXA", domain.FacetaFase: "CORRENTE",
		domain.FacetaDestinacao: "GUARDA_PERMANENTE", domain.FacetaAno: "2018"}
	if res, err := s.BuscarPublico(ctx, domain.Consulta{Q: marca, Filtros: filtros, SoDigital: true}); err != nil || res.Total != 1 {
		t.Fatalf("com as facetas: %+v %v", res, err)
	}
	if res, err := s.BuscarPublico(ctx, domain.Consulta{Q: ocr}); err != nil || res.Total != 1 || res.Itens[0].ID != u.ID {
		t.Fatalf("pelo texto do OCR: %+v %v", res, err)
	}
	if res, _ := s.BuscarPublico(ctx, domain.Consulta{Q: marca, Filtros: map[string]string{domain.FacetaAno: "1990"}}); res.Total != 0 {
		t.Fatalf("fora do ano: %+v", res)
	}
	if _, err := s.BuscarPublico(ctx, domain.Consulta{Filtros: map[string]string{domain.FacetaNivel: "RESTRITO"}}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("a pública não filtra por nível: %v", err)
	}
	if res, err := s.BuscarPublico(ctx, domain.Consulta{Filtros: map[string]string{domain.FacetaOrgao: orgao}}); err != nil || res.Total < 1 {
		t.Fatalf("sem termo: %+v %v", res, err)
	}

	// A ficha e o documento, sem login.
	f, err := s.FichaPublica(ctx, u.ID)
	if err != nil || f.Numero != u.Numero || len(f.Digitais) != 1 || f.Digitais[0].ID != d.ID || f.Descritor != "Leis e decretos" || f.Orgao != orgao {
		t.Fatalf("ficha: %+v %v", f, err)
	}
	doc, err := s.DocumentoPublico(ctx, d.ID)
	if err != nil {
		t.Fatal(err)
	}
	rc, _ := s.LerArquivo(ctx, doc, 0, 4)
	if b, _ := io.ReadAll(rc); string(b) != "%PDF" {
		t.Fatalf("documento: %q", b)
	}
	painel, err := s.ArquivoPublico(ctx)
	if err != nil || painel.SICURL != "https://sic.exemplo.gov.br" || painel.SerieSIC != domain.SerieSIC || painel.Unidades < 1 || painel.Digitais < 1 {
		t.Fatalf("página do Arquivo: %+v %v", painel, err)
	}
	if consultas(t, e, domain.ConsultaBusca) < 1 || consultas(t, e, domain.ConsultaFicha) < 1 {
		t.Fatal("as consultas públicas não foram contadas")
	}

	// Um digital ainda não liberado não aparece na ficha nem abre.
	outro := enviar(t, s, a.store, a.quem, u.ID, "TEXTO")
	if _, err := s.ProcessarProximo(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.DocumentoPublico(ctx, outro.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("não liberado: %v", err)
	}
	if f, _ := s.FichaPublica(ctx, u.ID); len(f.Digitais) != 1 {
		t.Fatalf("a ficha mostra só o liberado: %+v", f.Digitais)
	}

	// Restrita: sai da consulta pública na hora (gatilho do banco).
	a.nivel(t, u, domain.AcessoRestrito, "Informação pessoal (Lei 12.527/2011, art. 31)")
	for _, q := range []string{marca, ocr} {
		if res, _ := s.BuscarPublico(ctx, domain.Consulta{Q: q}); res.Total != 0 {
			t.Fatalf("restrita na busca pública (%s): %+v", q, res)
		}
	}
	if _, err := s.FichaPublica(ctx, u.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("ficha da restrita: %v", err)
	}
	if _, err := s.DocumentoPublico(ctx, d.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("documento da restrita: %v", err)
	}
	if det, _ := s.Digital(ctx, a.quem, d.ID); det.PublicadoEm != nil {
		t.Fatal("o banco não retirou o representante da consulta pública")
	}
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("liberar a restrita: %v", err)
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET publicado_em = now() WHERE id = $1`, d.ID); err == nil {
		t.Fatal("o banco liberou o representante da unidade restrita")
	}
	if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET publicado_em = now() WHERE id = $1`, outro.ID); err == nil {
		t.Fatal("o banco liberou o representante da unidade restrita")
	}

	// Na busca interna: quem tem o escopo vê (com o nível); os demais não.
	interna := domain.Consulta{Q: marca, Filtros: map[string]string{domain.FacetaNivel: "RESTRITO"}}
	if res, err := s.Buscar(ctx, a.quem, interna); err != nil || res.Total != 1 || res.Itens[0].NivelAcesso != "RESTRITO" || res.Itens[0].Digitais != 2 {
		t.Fatalf("interna no escopo: %+v %v", res, err)
	}
	if res, _ := s.Buscar(ctx, agente(e.outra), interna); res.Total != 0 {
		t.Fatalf("interna fora do escopo: %+v", res)
	}
	if res, err := s.Buscar(ctx, auth.Identity{Subject: "sem-acervo"}, interna); err != nil || res.Total != 0 || res.Itens == nil {
		t.Fatalf("sem escopo nenhum: %+v %v", res, err)
	}
	if res, _ := s.Buscar(ctx, arquivoPublico(uuid.Nil), interna); res.Total != 1 {
		t.Fatalf("Arquivo Público: %+v", res)
	}
	if _, err := s.Buscar(ctx, a.quem, domain.Consulta{Filtros: map[string]string{"x": "y"}}); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("filtro desconhecido: %v", err)
	}

	// Volta a ser pública: liberar e retirar.
	u2, _ := s.Get(ctx, admin, u.ID)
	a.nivel(t, u2, domain.AcessoPublico, "")
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	if ret, err := s.PublicarDigital(ctx, a.quem, d.ID, false); err != nil || ret.PublicadoEm != nil {
		t.Fatalf("retirar: %+v %v", ret, err)
	}
	if _, err := s.DocumentoPublico(ctx, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("documento inexistente: %v", err)
	}
	// Defesa em profundidade: mesmo que o gatilho falhasse (aqui, desligado
	// na sessão), a unidade que não é ostensiva não entrega o documento.
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	conn, err := e.pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(ctx, `SET session_replication_role = replica;
		UPDATE acervo_unidades SET nivel_acesso = 'SIGILOSO', hipotese_legal = 'Teste' WHERE id = '`+u.ID.String()+`';
		RESET session_replication_role`)
	conn.Release()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.DocumentoPublico(ctx, d.ID); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("documento de unidade sigilosa: %v", err)
	}

	if _, err := s.FichaPublica(ctx, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("ficha inexistente: %v", err)
	}
}

// indiceFalso faz o papel do Typesense.
type indiceFalso struct {
	mu                                    sync.Mutex
	criou                                 bool
	falhaGarantir                         int // falha nas n primeiras
	falhaGravar, falhaLimpar, falhaBuscar bool
	garantias, gravacoes, limpezas        int
	gravados                              map[uuid.UUID]domain.RegistroIndice
	resultado                             domain.ResultadoBusca
	filtro                                domain.FiltroBusca
}

func (f *indiceFalso) Garantir(context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.garantias++
	if f.garantias <= f.falhaGarantir {
		return false, errors.New("typesense fora")
	}
	return f.criou, nil
}

func (f *indiceFalso) Gravar(_ context.Context, regs []domain.RegistroIndice, _ int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.gravacoes++
	if f.falhaGravar {
		return errors.New("typesense fora")
	}
	if f.gravados == nil {
		f.gravados = map[uuid.UUID]domain.RegistroIndice{}
	}
	for _, r := range regs {
		f.gravados[r.Unidade.ID] = r
	}
	return nil
}

func (f *indiceFalso) Limpar(context.Context, int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.limpezas++
	if f.falhaLimpar {
		return errors.New("typesense fora")
	}
	return nil
}

func (f *indiceFalso) Buscar(_ context.Context, _ domain.Consulta, filtro domain.FiltroBusca) (domain.ResultadoBusca, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.filtro = filtro
	if f.falhaBuscar {
		return domain.ResultadoBusca{}, errors.New("typesense fora")
	}
	return f.resultado, nil
}

func (f *indiceFalso) contagens() (int, int, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.garantias, f.gravacoes, f.limpezas
}

func TestIndiceDoAcervo(t *testing.T) {
	a := novoAcesso(t)
	ctx := context.Background()
	marca := "Indice" + strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	u, d := a.pronta(t, marca)
	fi := &indiceFalso{resultado: domain.ResultadoBusca{Origem: domain.OrigemIndice, Total: 9}}
	s := a.s.WithBusca(application.ConfigBusca{Indice: fi})

	// Toda mudança marca a unidade; a indexação leva o texto do OCR e a série.
	if n, err := s.IndexarPendentes(ctx); err != nil || n < 1 {
		t.Fatalf("pendentes: %d %v", n, err)
	}
	reg, ok := fi.gravados[u.ID]
	if !ok || reg.Descritor != "Leis e decretos" || len(reg.Digitais) != 1 || reg.Digitais[0].Texto != "Ofício nº 12/2015" || reg.Orgao == "" {
		t.Fatalf("registro indexado: %+v", reg)
	}
	if n, err := s.IndexarPendentes(ctx); err != nil || n != 0 {
		t.Fatalf("fila vazia: %d %v", n, err)
	}
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	delete(fi.gravados, u.ID)
	if _, err := s.IndexarPendentes(ctx); err != nil || fi.gravados[u.ID].Digitais[0].PublicadoEm == nil {
		t.Fatalf("liberação reindexada: %+v %v", fi.gravados[u.ID], err)
	}
	if n, err := s.Recarregar(ctx); err != nil || n < 1 || fi.limpezas != 1 {
		t.Fatalf("carga completa: %d %v", n, err)
	}
	fi.falhaGravar = true
	if _, err := s.Recarregar(ctx); err == nil {
		t.Fatal("falha do índice na carga completa")
	}
	a.e.criar(codElim) // marca uma pendente
	if _, err := s.IndexarPendentes(ctx); err == nil {
		t.Fatal("falha do índice na indexação")
	}
	fi.falhaGravar = false

	// A busca usa o índice com o escopo; com ele fora, volta ao banco.
	if res, err := s.Buscar(ctx, a.quem, domain.Consulta{Q: "x"}); err != nil || res.Origem != domain.OrigemIndice || res.Total != 9 ||
		!fi.filtro.Restrito || len(fi.filtro.Unidades) == 0 {
		t.Fatalf("interna pelo índice: %+v %+v %v", res, fi.filtro, err)
	}
	if _, _ = s.Buscar(ctx, arquivoPublico(uuid.Nil), domain.Consulta{}); fi.filtro.Restrito {
		t.Fatal("o Arquivo Público vê todo o acervo")
	}
	if _, _ = s.BuscarPublico(ctx, domain.Consulta{}); !fi.filtro.Publico {
		t.Fatal("a consulta pública vai à coleção pública")
	}
	fi.falhaBuscar = true
	if res, err := s.Buscar(ctx, a.quem, domain.Consulta{Q: marca}); err != nil || res.Origem != domain.OrigemBanco || res.Total != 1 {

		t.Fatalf("volta ao banco: %+v %v", res, err)
	}

	// Sem índice: nada a indexar.
	sem := a.e.real()
	if n, err := sem.IndexarPendentes(ctx); n != 0 || err != nil {
		t.Fatalf("sem índice: %d %v", n, err)
	}
	if n, err := sem.Recarregar(ctx); n != 0 || err != nil {
		t.Fatalf("sem índice: %d %v", n, err)
	}
}

func esperar(t *testing.T, cond func() bool) {
	t.Helper()
	for limite := time.Now().Add(10 * time.Second); !cond(); time.Sleep(5 * time.Millisecond) {
		if time.Now().After(limite) {
			t.Fatal("tempo esgotado")
		}
	}
}

func TestRunIndice(t *testing.T) {
	e := novoEnv(t)
	rodar := func(s *application.Service, recarga time.Duration, ate func() bool) {
		ctx, cancel := context.WithCancel(context.Background())
		fim := make(chan error, 1)
		go func() { fim <- s.RunIndice(2*time.Millisecond, recarga)(ctx) }()
		esperar(t, ate)
		cancel()
		if err := <-fim; err != nil {
			t.Fatal(err)
		}
	}
	// Sem índice, só espera.
	rodar(e.real(), time.Hour, func() bool { return true })
	ctx0, cancel0 := context.WithCancel(context.Background())
	cancel0()
	if err := e.real().WithBusca(application.ConfigBusca{Indice: &indiceFalso{}}).RunIndice(0, time.Hour)(ctx0); err != nil {
		t.Fatal(err)
	}

	// O Typesense volta; a coleção nasce vazia: carga completa.
	fi := &indiceFalso{criou: true, falhaGarantir: 1}
	rodar(e.real().WithBusca(application.ConfigBusca{Indice: fi}), time.Hour, func() bool {
		g, _, l := fi.contagens()
		return g >= 2 && l >= 1
	})
	// Coleção existente, índice com falha: a carga e a indexação registram
	// o erro e seguem.
	ruim := &indiceFalso{falhaGravar: true}
	e.criar(codElim)
	rodar(e.real().WithBusca(application.ConfigBusca{Indice: ruim}), 0, func() bool {
		_, gr, _ := ruim.contagens()
		return gr >= 3
	})
	// A coleção existente não recarrega antes da hora.
	quieto := &indiceFalso{}
	rodar(e.real().WithBusca(application.ConfigBusca{Indice: quieto}), time.Hour, func() bool {
		g, _, _ := quieto.contagens()
		return g >= 1
	})
	if _, _, l := quieto.contagens(); l != 0 {
		t.Fatal("recarregou a coleção existente")
	}
}

// Exposições virtuais: o Arquivo Público monta e publica; a visita mostra só
// o liberado.
func TestExposicoes(t *testing.T) {
	a := novoAcesso(t)
	s := a.s
	ctx := context.Background()
	arq := arquivoPublico(uuid.Nil)
	_, d := a.pronta(t, "Exposição")
	_, oculto := a.pronta(t, "Exposição")
	if _, err := s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	titulo := "Rondonópolis em 1953 " + uuid.NewString()[:8]
	nova := application.NovaExposicao{Titulo: titulo, Resumo: "A fundação", Texto: "## Contexto",
		Itens: []application.NovoItem{{DigitalID: oculto.ID, Legenda: "Interno"}, {DigitalID: d.ID, Legenda: "Planta"}}}

	if _, err := s.SalvarExposicao(ctx, a.quem, nil, nova); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("só o Arquivo Público: %v", err)
	}
	ruim := nova
	ruim.Titulo = ""
	if _, err := s.SalvarExposicao(ctx, arq, nil, ruim); codigoDe(err) != apperrors.CodeValidation {
		t.Fatalf("sem título: %v", err)
	}
	ruim = nova
	ruim.Itens = []application.NovoItem{{DigitalID: uuid.New()}}
	if _, err := s.SalvarExposicao(ctx, arq, nil, ruim); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("documento inexistente: %v", err)
	}
	x, err := s.SalvarExposicao(ctx, arq, nil, nova)
	if err != nil || !strings.HasPrefix(x.Slug, "rondonopolis-em-1953-") || len(x.Itens) != 2 || x.Itens[0].Publicado || !x.Itens[1].Publicado ||
		x.Itens[1].Legenda != "Planta" || x.PublicadaEm != nil {
		t.Fatalf("criar: %+v %v", x, err)
	}
	if _, err := s.SalvarExposicao(ctx, arq, nil, nova); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("mesmo endereço: %v", err)
	}
	if _, err := s.ExposicaoPublica(ctx, x.Slug); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("rascunho na visita: %v", err)
	}

	// Alterar: o endereço não muda.
	nova.Titulo, nova.Itens = "Outro título", nova.Itens[:1]
	if _, err := s.SalvarExposicao(ctx, arq, &x.ID, nova); err != nil {
		t.Fatal(err)
	}
	alt, _ := s.Exposicao(ctx, arq, x.ID)
	if alt.Slug != x.Slug || alt.Titulo != "Outro título" || len(alt.Itens) != 1 {
		t.Fatalf("alterar: %+v", alt)
	}
	if _, err := s.PublicarExposicao(ctx, arq, x.ID, true); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("sem documento liberado: %v", err)
	}
	nova.Itens = []application.NovoItem{{DigitalID: oculto.ID}, {DigitalID: d.ID}}
	if _, err := s.SalvarExposicao(ctx, arq, &x.ID, nova); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PublicarExposicao(ctx, arq, x.ID, false); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("recolher a não publicada: %v", err)
	}
	if p, err := s.PublicarExposicao(ctx, arq, x.ID, true); err != nil || p.PublicadaEm == nil {
		t.Fatalf("publicar: %+v %v", p, err)
	}
	if _, err := s.PublicarExposicao(ctx, arq, x.ID, true); codigoDe(err) != apperrors.CodeConflict {
		t.Fatalf("publicar de novo: %v", err)
	}

	// A visita: só o documento liberado.
	v, err := s.ExposicaoPublica(ctx, x.Slug)
	if err != nil || len(v.Itens) != 1 || v.Itens[0].DigitalID != d.ID {
		t.Fatalf("visita: %+v %v", v, err)
	}
	if consultas(t, a.e, domain.ConsultaExposicao) < 1 {
		t.Fatal("visita não contada")
	}
	if items, total, err := s.ExposicoesPublicas(ctx, pagination.New(1, 50, 50)); err != nil || total < 1 || !temExposicao(items, x.ID) {
		t.Fatalf("publicadas: %d %v", total, err)
	}
	if painel, _ := s.ArquivoPublico(ctx); len(painel.Exposicoes) == 0 {
		t.Fatal("exposição em destaque")
	}
	if items, total, err := s.Exposicoes(ctx, arq, pagination.New(1, 50, 50)); err != nil || total < 1 || !temExposicao(items, x.ID) {
		t.Fatalf("gestão: %d %v", total, err)
	}
	if _, _, err := s.Exposicoes(ctx, a.quem, pagination.New(1, 50, 50)); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("lista só para o Arquivo: %v", err)
	}
	if _, err := s.Exposicao(ctx, a.quem, x.ID); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("ver só para o Arquivo: %v", err)
	}
	if _, err := s.PublicarExposicao(ctx, a.quem, x.ID, true); codigoDe(err) != apperrors.CodeForbidden {
		t.Fatalf("publicar só para o Arquivo: %v", err)
	}
	if _, err := s.Exposicao(ctx, arq, uuid.New()); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("inexistente: %v", err)
	}
	if _, err := s.SalvarExposicao(ctx, arq, ptr(uuid.New()), nova); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("alterar inexistente: %v", err)
	}
	if _, err := s.PublicarExposicao(ctx, arq, uuid.New(), true); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("publicar inexistente: %v", err)
	}

	// Recolher: some da visita.
	if r, err := s.PublicarExposicao(ctx, arq, x.ID, false); err != nil || r.PublicadaEm != nil {
		t.Fatalf("recolher: %+v %v", r, err)
	}
	if _, err := s.ExposicaoPublica(ctx, x.Slug); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("recolhida na visita: %v", err)
	}
	if _, err := s.ExposicaoPublica(ctx, "nao-existe-"+uuid.NewString()[:8]); codigoDe(err) != apperrors.CodeNotFound {
		t.Fatalf("endereço inexistente: %v", err)
	}
}

func temExposicao(items []domain.Exposicao, id uuid.UUID) bool {
	for _, x := range items {
		if x.ID == id {
			return true
		}
	}
	return false
}

func ptr[T any](v T) *T { return &v }

// A consulta pública com o banco fora: a falha na busca, na ficha e na
// visita aparece; a falha só na contagem não derruba a consulta.
func TestConsultaPublicaBancoFora(t *testing.T) {
	a := novoAcesso(t)
	ctx := context.Background()
	u, d := a.pronta(t, "Banco fora")
	if _, err := a.s.PublicarDigital(ctx, a.quem, d.ID, true); err != nil {
		t.Fatal(err)
	}
	arq := arquivoPublico(uuid.Nil)
	x, err := a.s.SalvarExposicao(ctx, arq, nil, application.NovaExposicao{Titulo: "Banco fora " + uuid.NewString()[:8], Itens: []application.NovoItem{{DigitalID: d.ID}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.s.PublicarExposicao(ctx, arq, x.ID, true); err != nil {
		t.Fatal(err)
	}
	com := func(failAt int) *application.Service {
		return a.e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: failAt}).WithBusca(application.ConfigBusca{})
	}
	if _, err := com(1).BuscarPublico(ctx, domain.Consulta{Q: "x"}); err == nil {
		t.Error("busca com o banco fora")
	}
	if _, err := com(1).FichaPublica(ctx, u.ID); err == nil {
		t.Error("ficha com o banco fora")
	}
	if _, err := com(1).ExposicaoPublica(ctx, x.Slug); err == nil {
		t.Error("visita com o banco fora")
	}
	if _, _, err := com(1).ExposicoesPublicas(ctx, pagination.New(1, 5, 5)); err == nil {
		t.Error("exposições com o banco fora")
	}
	// A contagem é a última chamada: falha registrada, consulta servida.
	if res, err := com(2).BuscarPublico(ctx, domain.Consulta{Q: "Banco fora"}); err != nil || res.Total < 1 {
		t.Errorf("busca com a contagem fora: %+v %v", res, err)
	}
	if _, err := com(2).FichaPublica(ctx, u.ID); err != nil {
		t.Errorf("ficha com a contagem fora: %v", err)
	}
	if _, err := com(2).ExposicaoPublica(ctx, x.Slug); err != nil {
		t.Errorf("visita com a contagem fora: %v", err)
	}
	com(1).ContarDocumento(ctx)
}

var _ storage.Provider = (*wormMem)(nil)
