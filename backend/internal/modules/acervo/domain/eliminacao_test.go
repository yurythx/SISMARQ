package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func comissao(n int) Comissao {
	areasSetorial := []string{AreaArquivologia, AreaJuridico, AreaAdmFinancas, AreaUnidadeProdutora, AreaOutra, AreaOutra}
	ent := uuid.New()
	c := Comissao{Tipo: ComissaoSetorial, EntidadeID: &ent, Nome: " CSPAD da SEMAD ", Portaria: "Portaria 12/2026"}
	for i := range n {
		papel := PapelMembro
		if i == 0 {
			papel = PapelPresidente
		}
		c.Membros = append(c.Membros, Membro{UserID: uuid.New(), Papel: papel, Area: areasSetorial[i]})
	}
	return c
}

func TestComissaoValidar(t *testing.T) {
	c := comissao(5)
	if err := c.Validar(); err != nil || c.Nome != "CSPAD da SEMAD" || len(c.Avisos) != 0 {
		t.Fatalf("válida: %+v %v", c, err)
	}
	par := comissao(4)
	if err := par.Validar(); err != nil || len(par.Avisos) != 1 {
		t.Fatalf("número par vira aviso: %+v %v", par.Avisos, err)
	}
	central := Comissao{Tipo: ComissaoCentral, Nome: "CCPAD", Membros: []Membro{
		{UserID: uuid.New(), Papel: PapelPresidente, Area: AreaArquivologia}, {UserID: uuid.New(), Papel: PapelMembro, Area: AreaJuridico},
		{UserID: uuid.New(), Papel: PapelMembro, Area: AreaAdmFinancas}, {UserID: uuid.New(), Papel: PapelSecretario, Area: AreaArquivoPublico},
		{UserID: uuid.New(), Papel: PapelMembro, Area: AreaOutra}}}
	if err := central.Validar(); err != nil {
		t.Fatalf("CCPAD: %v", err)
	}
	ent := uuid.New()
	casos := map[string]func(*Comissao){
		"tipo":                func(c *Comissao) { c.Tipo = "X" },
		"setorial sem órgão":  func(c *Comissao) { c.EntidadeID = nil },
		"central com órgão":   func(c *Comissao) { c.Tipo = ComissaoCentral; c.EntidadeID = &ent },
		"sem nome":            func(c *Comissao) { c.Nome = " " },
		"portaria longa":      func(c *Comissao) { c.Portaria = strings.Repeat("x", 201) },
		"membro repetido":     func(c *Comissao) { c.Membros = append(c.Membros, c.Membros[1]) },
		"papel inválido":      func(c *Comissao) { c.Membros[1].Papel = "VOGAL" },
		"área inválida":       func(c *Comissao) { c.Membros[1].Area = "TI" },
		"sem presidente":      func(c *Comissao) { c.Membros[0].Papel = PapelMembro },
		"dois presidentes":    func(c *Comissao) { c.Membros[1].Papel = PapelPresidente },
		"sem jurídico (§ 6º)": func(c *Comissao) { c.Membros[1].Area = AreaOutra },
	}
	for nome, mudar := range casos {
		v := comissao(5)
		mudar(&v)
		if err := v.Validar(); err == nil {
			t.Errorf("%s: aceita", nome)
		}
	}
}

func listagem(situacao string, itens int) Listagem {
	l := Listagem{Situacao: situacao, ParecerJuridico: "Parecer 1/2026"}
	for range itens {
		l.Itens = append(l.Itens, ItemListagem{UnidadeLista: UnidadeLista{Unidade: Unidade{ID: uuid.New()}}})
	}
	return l
}

func TestPodeAprovar(t *testing.T) {
	c := comissao(5)
	todos := []uuid.UUID{c.Membros[0].UserID, c.Membros[1].UserID, c.Membros[2].UserID}
	em := time.Now()
	if err := PodeAprovar(listagem(ListagemRascunho, 2), c, "1/2026", em, todos); err != nil {
		t.Fatalf("aprovação: %v", err)
	}
	semParecer := listagem(ListagemRascunho, 2)
	semParecer.ParecerJuridico = " "
	vazia := listagem(ListagemRascunho, 1)
	vazia.Itens[0].Amostra = true
	casos := []struct {
		nome      string
		l         Listagem
		numero    string
		em        time.Time
		presentes []uuid.UUID
		err       error
	}{
		{"etapa", listagem(ListagemAprovada, 2), "1", em, todos, ErrListagemSituacao},
		{"vazia", vazia, "1", em, todos, ErrListagemVazia},
		{"sem parecer (art. 32, § 4º)", semParecer, "1", em, todos, ErrSemParecer},
		{"sem presidente", listagem(ListagemRascunho, 1), "1", em, []uuid.UUID{c.Membros[1].UserID, c.Membros[2].UserID, c.Membros[3].UserID}, ErrQuorum},
		{"minoria", listagem(ListagemRascunho, 1), "1", em, []uuid.UUID{c.Membros[0].UserID, c.Membros[1].UserID}, ErrQuorum},
	}
	for _, k := range casos {
		if err := PodeAprovar(k.l, c, k.numero, k.em, k.presentes); !errors.Is(err, k.err) {
			t.Errorf("%s: %v", k.nome, err)
		}
	}
	for nome, f := range map[string]func() error{
		"sem número":   func() error { return PodeAprovar(listagem(ListagemRascunho, 1), c, " ", em, todos) },
		"sem data":     func() error { return PodeAprovar(listagem(ListagemRascunho, 1), c, "1", time.Time{}, todos) },
		"não é membro": func() error { return PodeAprovar(listagem(ListagemRascunho, 1), c, "1", em, append(todos, uuid.New())) },
	} {
		var v *ErrValidacao
		if err := f(); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
}

func TestEditalPedidosAmostraTermo(t *testing.T) {
	hoje := dia("2026-10-15")
	pub := dia("2026-09-01")
	if err := PodePublicar(listagem(ListagemAutorizada, 1), "5/2026", pub, hoje); err != nil {
		t.Fatal(err)
	}
	for nome, err := range map[string]error{
		"etapa":  PodePublicar(listagem(ListagemAprovada, 1), "5", pub, hoje),
		"número": PodePublicar(listagem(ListagemAutorizada, 1), "", pub, hoje),
		"futura": PodePublicar(listagem(ListagemAutorizada, 1), "5", dia("2026-10-16"), hoje),
	} {
		if err == nil {
			t.Errorf("edital %s: aceito", nome)
		}
	}

	l := listagem(ListagemEmEdital, 10)
	l.EditalPublicadoEm = &pub
	if !l.PrazoFinal().Equal(dia("2026-10-01")) || (Listagem{}).PrazoFinal() != nil || !l.Ativa() {
		t.Fatalf("prazo: %v", l.PrazoFinal())
	}
	if err := PodeReceberPedido(l, dia("2026-09-20")); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"2026-08-31", "2026-10-02"} {
		if err := PodeReceberPedido(l, dia(d)); !errors.Is(err, ErrPedidoForaDoPrazo) {
			t.Errorf("pedido em %s: %v", d, err)
		}
	}
	if err := PodeReceberPedido(listagem(ListagemAutorizada, 1), hoje); !errors.Is(err, ErrListagemSituacao) {
		t.Fatalf("pedido fora do edital: %v", err)
	}

	if TamanhoAmostra(10, 5) != 1 || TamanhoAmostra(10, 25) != 3 || TamanhoAmostra(3, 100) != 3 || TamanhoAmostra(3, 0.1) != 1 {
		t.Fatal("tamanho da amostra")
	}
	if err := PodeAmostrar(l, "Critérios", 10); err != nil {
		t.Fatal(err)
	}
	feita := l
	agora := time.Now()
	feita.AmostraEm = &agora
	vazia := listagem(ListagemEmEdital, 0)
	for nome, err := range map[string]error{
		"etapa":        PodeAmostrar(listagem(ListagemAutorizada, 1), "c", 10),
		"já feita":     PodeAmostrar(feita, "c", 10),
		"sem critério": PodeAmostrar(l, " ", 10),
		"zero":         PodeAmostrar(l, "c", 0),
		"acima de 100": PodeAmostrar(l, "c", 101),
		"vazia":        PodeAmostrar(vazia, "c", 10),
	} {
		if err == nil {
			t.Errorf("amostragem %s: aceita", nome)
		}
	}

	if err := PodeRegistrarTermo(feita, dia("2026-10-10"), hoje); err != nil {
		t.Fatalf("termo: %v", err)
	}
	pendente := feita
	pendente.Pedidos = []Pedido{{Situacao: PedidoDeferido}, {Situacao: PedidoPendente}}
	for nome, c := range map[string]struct {
		l      Listagem
		em, hj string
		err    error
	}{
		"no prazo (hoje)":     {feita, "2026-09-30", "2026-09-30", ErrEditalNoPrazo},
		"data antes do prazo": {feita, "2026-09-30", "2026-10-15", ErrEditalNoPrazo},
		"sem amostra":         {l, "2026-10-10", "2026-10-15", ErrSemAmostra},
		"pedido pendente":     {pendente, "2026-10-10", "2026-10-15", ErrPedidoPendente},
		"etapa":               {listagem(ListagemAutorizada, 1), "2026-10-10", "2026-10-15", ErrListagemSituacao},
	} {
		if c.l.EditalPublicadoEm == nil {
			c.l.EditalPublicadoEm = &pub
		}
		if err := PodeRegistrarTermo(c.l, dia(c.em), dia(c.hj)); !errors.Is(err, c.err) {
			t.Errorf("termo %s: %v", nome, err)
		}
	}
	var v *ErrValidacao
	if err := PodeRegistrarTermo(feita, dia("2026-10-20"), hoje); !errors.As(err, &v) {
		t.Fatalf("termo com data futura: %v", err)
	}

	if PodeCancelar(l) != nil || PodeCancelar(feita) == nil || PodeCancelar(listagem(ListagemEliminada, 1)) == nil {
		t.Fatal("cancelamento")
	}
	l.Itens[0].Amostra, l.Itens[1].Retirada = true, true
	if len(l.AEliminar()) != 8 {
		t.Fatalf("a eliminar: %d", len(l.AEliminar()))
	}
}
