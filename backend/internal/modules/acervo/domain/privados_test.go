package domain

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func privado() Privado {
	return Privado{Titular: " Família Silva ", TipoTitular: "pf", Descricao: "Fotografias da fundação da cidade (1953-1970)"}
}

func TestNovoPrivado(t *testing.T) {
	p, err := NovoPrivado(privado())
	if err != nil || p.Situacao != PrivadoCadastrado || p.Guarda != GuardaTitular || p.TipoTitular != TitularPF || p.Titular != "Família Silva" {
		t.Fatalf("cadastro: %+v %v", p, err)
	}
	concessionaria := privado()
	concessionaria.TipoTitular, concessionaria.ServicoPublico = "PJ", true
	if p, err := NovoPrivado(concessionaria); err != nil || p.Situacao != PrivadoDeclarado {
		t.Fatalf("§ 3º: %+v %v", p, err)
	}
	for nome, mudar := range map[string]func(*Privado){
		"tipo":                   func(p *Privado) { p.TipoTitular = "X" },
		"serviço público por PF": func(p *Privado) { p.ServicoPublico = true },
		"sem titular":            func(p *Privado) { p.Titular = " " },
		"sem descrição":          func(p *Privado) { p.Descricao = "" },
		"documento longo":        func(p *Privado) { p.Documento = strings.Repeat("1", 21) },
		"condições longas":       func(p *Privado) { p.AcessoCondicoes = strings.Repeat("x", 2001) },
	} {
		x := privado()
		mudar(&x)
		var v *ErrValidacao
		if _, err := NovoPrivado(x); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
}

func TestEtapasDoPrivado(t *testing.T) {
	hoje := dia("2026-10-15")
	p, _ := NovoPrivado(privado())
	if err := PodeParecer(p, "Comissão: A, B e C", "Relevante", dia("2026-09-01"), hoje); err != nil {
		t.Fatal(err)
	}
	var v *ErrValidacao
	for nome, err := range map[string]error{
		"sem comissão": PodeParecer(p, " ", "x", hoje, hoje),
		"sem texto":    PodeParecer(p, "c", "", hoje, hoje),
		"data futura":  PodeParecer(p, "c", "x", dia("2026-10-16"), hoje),
	} {
		if !errors.As(err, &v) {
			t.Errorf("parecer %s: %v", nome, err)
		}
	}
	if !errors.Is(PodeDeclarar(p, "1", "", hoje, hoje), ErrPrivadoSituacao) {
		t.Fatal("decreto sem parecer")
	}
	parecer := dia("2026-09-01")
	p.Situacao, p.ParecerEm = PrivadoParecer, &parecer
	if !errors.Is(PodeParecer(p, "c", "x", hoje, hoje), ErrPrivadoSituacao) {
		t.Fatal("segundo parecer")
	}
	if err := PodeDeclarar(p, "Decreto 10.000/2026", "DIORONDON-e", hoje, hoje); err != nil {
		t.Fatal(err)
	}
	for nome, err := range map[string]error{
		"sem número":       PodeDeclarar(p, "", "", hoje, hoje),
		"diário longo":     PodeDeclarar(p, "1", strings.Repeat("x", 201), hoje, hoje),
		"antes do parecer": PodeDeclarar(p, "1", "", dia("2026-08-01"), hoje),
	} {
		if !errors.As(err, &v) {
			t.Errorf("decreto %s: %v", nome, err)
		}
	}

	// Guarda (§§ 4º e 5º).
	if !errors.Is(PodeMudarGuarda(p, GuardaDoado, "t", hoje, hoje), ErrPrivadoSituacao) {
		t.Fatal("guarda antes do decreto")
	}
	p.Situacao = PrivadoDeclarado
	if err := PodeMudarGuarda(p, GuardaDepositado, "Termo de depósito 1/2026", hoje, hoje); err != nil {
		t.Fatal(err)
	}
	deposito := p
	deposito.Guarda = GuardaDepositado
	if err := PodeMudarGuarda(deposito, GuardaTitular, "Revogação a pedido do titular", hoje, hoje); err != nil {
		t.Fatal(err)
	}
	doado := p
	doado.Guarda = GuardaDoado
	for nome, err := range map[string]error{
		"doado não volta": PodeMudarGuarda(doado, GuardaTitular, "t", hoje, hoje),
		"guarda inválida": PodeMudarGuarda(p, "VENDIDO", "t", hoje, hoje),
		"mesma guarda":    PodeMudarGuarda(p, GuardaTitular, "t", hoje, hoje),
		"sem termo":       PodeMudarGuarda(p, GuardaDoado, " ", hoje, hoje),
		"data futura":     PodeMudarGuarda(p, GuardaDoado, "t", dia("2026-10-16"), hoje),
	} {
		if !errors.As(err, &v) {
			t.Errorf("guarda %s: %v", nome, err)
		}
	}
	if err := PodeMudarGuarda(Privado{Situacao: PrivadoDeclarado, Guarda: GuardaDoado}, GuardaTitular, "t", hoje, hoje); !errors.As(err, &v) {
		t.Fatal("doado")
	}
}

// Art. 58: preferência de 60 dias na alienação.
func TestPreferencia(t *testing.T) {
	hoje := dia("2026-10-15")
	p := Privado{Situacao: PrivadoDeclarado, Guarda: GuardaTitular}
	a := Alienacao{ID: uuid.New(), NotificadoEm: dia("2026-09-01"), Decisao: AlienacaoPendente}
	a.Completar(hoje)
	if !a.PrazoFinal.Equal(dia("2026-10-31")) || a.Vencida {
		t.Fatalf("prazo: %+v", a)
	}
	if err := PodeNotificarAlienacao(p, a, hoje); err != nil {
		t.Fatal(err)
	}
	p.Alienacoes = []Alienacao{a}
	if !errors.Is(PodeNotificarAlienacao(p, a, hoje), ErrAlienacaoPendente) {
		t.Fatal("duas pendentes")
	}
	// Vencida a primeira, outra notificação pode chegar.
	if err := PodeNotificarAlienacao(p, Alienacao{NotificadoEm: dia("2026-11-10")}, dia("2026-11-10")); err != nil {
		t.Fatalf("depois do prazo: %v", err)
	}
	var v *ErrValidacao
	for nome, err := range map[string]error{
		"adquirente longo": PodeNotificarAlienacao(Privado{Situacao: PrivadoDeclarado}, Alienacao{NotificadoEm: hoje, Adquirente: strings.Repeat("x", 201)}, hoje),
		"condições longas": PodeNotificarAlienacao(Privado{Situacao: PrivadoDeclarado}, Alienacao{NotificadoEm: hoje, Condicoes: strings.Repeat("x", 2001)}, hoje),
		"data futura":      PodeNotificarAlienacao(Privado{Situacao: PrivadoDeclarado}, Alienacao{NotificadoEm: dia("2026-10-16")}, hoje),
	} {
		if !errors.As(err, &v) {
			t.Errorf("notificação %s: %v", nome, err)
		}
	}
	for nome, x := range map[string]Privado{"não declarado": {Situacao: PrivadoCadastrado}, "doado": {Situacao: PrivadoDeclarado, Guarda: GuardaDoado}} {
		if !errors.Is(PodeNotificarAlienacao(x, a, hoje), ErrPrivadoSituacao) {
			t.Errorf("notificação %s", nome)
		}
	}

	if err := PodeDecidirAlienacao(a, AlienacaoInteresse, "Interesse histórico", hoje); err != nil {
		t.Fatal(err)
	}
	decidida := a
	decidida.Decisao = AlienacaoSemInteresse
	casos := map[string]struct {
		err, want error
	}{
		"decidida": {PodeDecidirAlienacao(decidida, AlienacaoInteresse, "x", hoje), ErrAlienacaoDecidida},
		"vencida":  {PodeDecidirAlienacao(a, AlienacaoInteresse, "x", dia("2026-11-01")), ErrPreferenciaVencida},
	}
	for nome, c := range casos {
		if !errors.Is(c.err, c.want) {
			t.Errorf("decisão %s: %v", nome, c.err)
		}
	}
	for nome, err := range map[string]error{
		"decisão inválida": PodeDecidirAlienacao(a, "TALVEZ", "x", hoje),
		"sem fundamento":   PodeDecidirAlienacao(a, AlienacaoSemInteresse, " ", hoje),
	} {
		if !errors.As(err, &v) {
			t.Errorf("decisão %s: %v", nome, err)
		}
	}
}
