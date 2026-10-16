package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func dia(s string) time.Time {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		panic(err)
	}
	return t
}

func anos(n int) *int { return &n }

func valida() Unidade {
	return Unidade{Tipo: "caixa", CodigoClassificacao: " 2.0.01.00.01 ", UnidadeID: uuid.New(),
		DataInicial: dia("2018-01-01"), DataFinal: dia("2019-12-31")}
}

func TestNormalizarEValidar(t *testing.T) {
	u := valida()
	u.Normalizar()
	if u.Tipo != TipoCaixa || u.CodigoClassificacao != "2.0.01.00.01" || u.NivelAcesso != AcessoPublico || u.Quantidade != 1 {
		t.Fatalf("normalização: %+v", u)
	}
	if err := u.Validar(); err != nil {
		t.Fatal(err)
	}
	longo := func(n int) string { return strings.Repeat("x", n) }
	casos := map[string]func(*Unidade){
		"tipo":                   func(u *Unidade) { u.Tipo = "GAVETA" },
		"codigo":                 func(u *Unidade) { u.CodigoClassificacao = "" },
		"assunto":                func(u *Unidade) { u.Assunto = longo(501) },
		"unidade":                func(u *Unidade) { u.UnidadeID = uuid.Nil },
		"datas vazias":           func(u *Unidade) { u.DataInicial = time.Time{} },
		"datas invertidas":       func(u *Unidade) { u.DataFinal = dia("2017-01-01") },
		"quantidade":             func(u *Unidade) { u.Quantidade = -1 },
		"especificacao":          func(u *Unidade) { u.Especificacao = longo(121) },
		"observacoes":            func(u *Unidade) { u.Observacoes = longo(4001) },
		"nivel":                  func(u *Unidade) { u.NivelAcesso = "SECRETO" },
		"sigilo sem hipótese":    func(u *Unidade) { u.NivelAcesso = AcessoSigiloso },
		"hipotese longa":         func(u *Unidade) { u.NivelAcesso, u.HipoteseLegal = AcessoRestrito, longo(501) },
		"posicao":                func(u *Unidade) { u.Posicao = longo(121) },
		"evento antes do início": func(u *Unidade) { e := dia("2000-01-01"); u.DataEvento = &e },
	}
	for nome, mudar := range casos {
		v := u
		mudar(&v)
		var ev *ErrValidacao
		if err := v.Validar(); !errors.As(err, &ev) || ev.Error() == "" {
			t.Errorf("%s: aceito (%v)", nome, err)
		}
	}
}

func TestEditavelEBloqueio(t *testing.T) {
	if (Unidade{Fase: FaseEliminado}).Editavel() != ErrEliminado || (Unidade{Fase: FasePermanente}).Editavel() != nil {
		t.Fatal("eliminada é histórico; permanente aceita descrição e local")
	}
	id := uuid.New()
	if _, err := NovoBloqueio(id, " ", ""); err == nil {
		t.Fatal("processo obrigatório")
	}
	if _, err := NovoBloqueio(id, "1", strings.Repeat("x", 2001)); err == nil {
		t.Fatal("motivo longo")
	}
	b, err := NovoBloqueio(id, " 0001234-56.2026.8.11.0003 ", " ação de cobrança ")
	if err != nil || b.Processo != "0001234-56.2026.8.11.0003" || b.Motivo != "ação de cobrança" || !b.Ativo() {
		t.Fatalf("bloqueio: %+v %v", b, err)
	}
	agora := time.Now()
	b.LiberadoEm = &agora
	if b.Ativo() {
		t.Fatal("liberado não é ativo")
	}
}

func TestLocalValidar(t *testing.T) {
	l := Local{ID: uuid.New(), Tipo: " estante ", Nome: " Estante 4 "}
	if err := l.Validar(); err != nil || l.Tipo != LocalEstante || l.Nome != "Estante 4" {
		t.Fatalf("local: %+v %v", l, err)
	}
	for _, ruim := range []Local{
		{ID: l.ID, Tipo: "COFRE", Nome: "x"},
		{ID: l.ID, Tipo: LocalPredio, Nome: ""},
		{ID: l.ID, Tipo: LocalPredio, Nome: "x", ParentID: &l.ID},
	} {
		if err := ruim.Validar(); err == nil {
			t.Errorf("aceito: %+v", ruim)
		}
	}
}

func TestDestinacaoDa(t *testing.T) {
	if DestinacaoDa(Serie{}) != DestinacaoIndefinida || DestinacaoDa(Serie{Destinacao: DestinacaoEliminacao}) != DestinacaoEliminacao {
		t.Fatal("destinação")
	}
}

// Prazos por anos: corrente a partir da data-limite final; intermediária a
// partir do fim da corrente (arts. 29 e 30).
func TestPrazosPorAnos(t *testing.T) {
	u := valida()
	u.Fase = FaseCorrente
	s := Serie{CorrenteAnos: anos(2), IntermAnos: anos(5), Destinacao: DestinacaoEliminacao, Vigente: true}
	p := CalcularPrazos(u, s, false, dia("2020-06-01"))
	if !p.FimCorrente.Equal(dia("2021-12-31")) || !p.FimIntermediaria.Equal(dia("2026-12-31")) || p.Situacao != SituacaoNoCorrente {
		t.Fatalf("prazos: %+v", p)
	}
	if p.RegraCorrente != "2 ano(s)" || p.RegraIntermediaria != "5 ano(s)" {
		t.Fatalf("regras: %q %q", p.RegraCorrente, p.RegraIntermediaria)
	}
	if CalcularPrazos(u, s, false, dia("2022-01-01")).Situacao != SituacaoTransferir {
		t.Fatal("corrente cumprido: transferir")
	}
	u.Fase = FaseIntermediaria
	for hoje, esperado := range map[string]string{"2025-01-01": SituacaoNoIntermediario, "2027-01-01": SituacaoAvaliarEliminacao} {
		if got := CalcularPrazos(u, s, false, dia(hoje)).Situacao; got != esperado {
			t.Errorf("%s: %s, esperado %s", hoje, got, esperado)
		}
	}
	s.Destinacao = DestinacaoGuardaPermanente
	if CalcularPrazos(u, s, false, dia("2027-01-01")).Situacao != SituacaoRecolher {
		t.Fatal("guarda permanente: recolher")
	}
	s.Destinacao = ""
	if CalcularPrazos(u, s, false, dia("2027-01-01")).Situacao != SituacaoAguardaDestinacao {
		t.Fatal("sem destinação")
	}
	if CalcularPrazos(u, s, true, dia("2027-01-01")).Situacao != SituacaoBloqueada {
		t.Fatal("bloqueada")
	}
	for fase, sit := range map[Fase]string{FasePermanente: SituacaoPermanente, FaseEliminado: SituacaoEliminada} {
		u.Fase = fase
		if CalcularPrazos(u, s, true, dia("2027-01-01")).Situacao != sit {
			t.Errorf("%s", fase)
		}
	}
}

// "Enquanto estiver vigorando": a fase corrente termina no evento (fim da
// vigência); sem a data, o prazo aguarda.
func TestPrazosPorCondicao(t *testing.T) {
	u := valida()
	u.Fase = FaseCorrente
	s := Serie{CorrenteCondicao: "Enquanto estiver vigorando", IntermAnos: anos(5), Destinacao: DestinacaoEliminacao}
	p := CalcularPrazos(u, s, false, dia("2030-01-01"))
	if p.FimCorrente != nil || p.FimIntermediaria != nil || p.Situacao != SituacaoAguardaEvento || p.RegraIntermediaria != "5 ano(s)" {
		t.Fatalf("sem evento: %+v", p)
	}
	ev := dia("2022-03-10")
	u.DataEvento = &ev
	p = CalcularPrazos(u, s, false, dia("2030-01-01"))
	if !p.FimCorrente.Equal(ev) || !p.FimIntermediaria.Equal(dia("2027-03-10")) || p.Situacao != SituacaoTransferir {
		t.Fatalf("com evento: %+v", p)
	}
	// Anos somados ao evento; condição em meses e dias; dias após o evento;
	// condição desconhecida (vale o evento); sem prazo (zero).
	casos := []struct {
		anos     *int
		condicao string
		evento   *time.Time
		fim      string
	}{
		{anos(1), "Enquanto estiver vigorando", &ev, "2023-03-10"},
		{nil, "03 meses", nil, "2020-03-31"},
		{nil, "90 dias", nil, "2020-03-30"},
		{nil, "30 dias após a data do evento", &ev, "2022-04-09"},
		{nil, "06 meses após a data do evento", &ev, "2022-09-10"},
		{anos(2), "até a homologação", &ev, "2022-03-10"},
		{nil, "", nil, "2019-12-31"},
	}
	for _, c := range casos {
		fim, regra := fimDaFase(dia("2019-12-31"), c.anos, c.condicao, c.evento)
		if fim == nil || !fim.Equal(dia(c.fim)) || regra == "" {
			t.Errorf("%q: %v %q, esperado %s", c.condicao, fim, regra, c.fim)
		}
	}
	for _, cond := range []string{"30 dias após a data do evento", "até a homologação"} {
		if fim, _ := fimDaFase(dia("2019-12-31"), nil, cond, nil); fim != nil {
			t.Errorf("%q sem evento deveria aguardar", cond)
		}
	}
	// Intermediária que depende de evento: aguarda.
	u.DataEvento = nil
	s = Serie{CorrenteAnos: anos(1), IntermCondicao: "até a homologação", Destinacao: DestinacaoEliminacao}
	u.Fase = FaseIntermediaria
	if CalcularPrazos(u, s, false, dia("2030-01-01")).Situacao != SituacaoAguardaEvento {
		t.Fatal("intermediária sem evento")
	}
}

func TestPodeTransferir(t *testing.T) {
	u := valida()
	u.Fase = FaseCorrente
	hoje := dia("2030-01-01")
	fim := dia("2025-01-01")
	futuro := dia("2031-01-01")
	casos := []struct {
		fase Fase
		fim  *time.Time
		err  error
	}{
		{FaseIntermediaria, &fim, ErrFaseInvalida},
		{FaseCorrente, nil, ErrAguardaEvento},
		{FaseCorrente, &futuro, ErrPrazoNaoCumprido},
		{FaseCorrente, &fim, nil},
	}
	for _, c := range casos {
		u.Fase = c.fase
		if err := PodeTransferir(u, Prazos{FimCorrente: c.fim}, hoje); !errors.Is(err, c.err) && err != c.err {
			t.Errorf("%s %v: %v, esperado %v", c.fase, c.fim, err, c.err)
		}
	}
}

// As travas da lei, uma a uma.
func TestPodeEliminar(t *testing.T) {
	u := valida()
	u.Fase = FaseIntermediaria
	elim := Serie{CorrenteAnos: anos(1), IntermAnos: anos(1), Destinacao: DestinacaoEliminacao}
	hoje := dia("2030-01-01")
	if err := PodeEliminar(u, elim, false, hoje); err != nil {
		t.Fatalf("elegível: %v", err)
	}
	casos := map[string]struct {
		mudar     func(*Unidade, *Serie)
		bloqueada bool
		err       error
	}{
		"permanente (arts. 20 e 42)":  {func(u *Unidade, _ *Serie) { u.Fase = FasePermanente }, false, ErrPermanente},
		"já eliminada":                {func(u *Unidade, _ *Serie) { u.Fase = FaseEliminado }, false, ErrEliminado},
		"bloqueio judicial (art. 44)": {func(*Unidade, *Serie) {}, true, ErrBloqueioJudicial},
		"sem destinação":              {func(_ *Unidade, s *Serie) { s.Destinacao = "" }, false, ErrDestinacaoIndefinida},
		"guarda permanente (art. 19)": {func(_ *Unidade, s *Serie) { s.Destinacao = DestinacaoGuardaPermanente }, false, ErrNaoPrevistaNaTTDD},
		"prazo correndo (art. 30)":    {func(_ *Unidade, s *Serie) { s.IntermAnos = anos(50) }, false, ErrPrazoNaoCumprido},
		"aguarda evento":              {func(_ *Unidade, s *Serie) { s.CorrenteCondicao = "Enquanto estiver vigorando" }, false, ErrAguardaEvento},
	}
	for nome, c := range casos {
		v, s := u, elim
		c.mudar(&v, &s)
		if err := PodeEliminar(v, s, c.bloqueada, hoje); !errors.Is(err, c.err) {
			t.Errorf("%s: %v, esperado %v", nome, err, c.err)
		}
	}
	if err := PodeEliminarArt35(u, Serie{}, false, false, hoje); err != nil {
		t.Fatalf("fora da TTDD: %v", err)
	}
	if err := PodeEliminarArt35(u, elim, true, false, hoje); !errors.Is(err, ErrRitoOrdinario) {
		t.Fatalf("ordinário no art. 35: %v", err)
	}
	perm := elim
	perm.Destinacao = DestinacaoGuardaPermanente
	if err := PodeEliminarArt35(u, perm, true, false, hoje); !errors.Is(err, ErrNaoPrevistaNaTTDD) {
		t.Fatalf("permanente no art. 35: %v", err)
	}
	permU := u
	permU.Fase = FasePermanente
	if err := PodeEliminarArt35(permU, Serie{}, false, false, hoje); !errors.Is(err, ErrPermanente) {
		t.Fatalf("fase permanente no art. 35: %v", err)
	}
	elimU := u
	elimU.Fase = FaseEliminado
	if err := PodeEliminarArt35(elimU, Serie{}, false, false, hoje); !errors.Is(err, ErrEliminado) {
		t.Fatalf("já eliminada no art. 35: %v", err)
	}
	if err := PodeEliminarArt35(u, Serie{}, false, true, hoje); !errors.Is(err, ErrBloqueioJudicial) {
		t.Fatalf("bloqueio no art. 35: %v", err)
	}
	longa := elim
	longa.IntermAnos = anos(50)
	if err := PodeEliminarArt35(u, longa, true, false, hoje); !errors.Is(err, ErrPrazoNaoCumprido) {
		t.Fatalf("prazo no art. 35: %v", err)
	}
}
