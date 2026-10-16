package domain

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var preparoOk = Preparo{Classificados: true, Avaliados: true, Organizados: true, Higienizados: true, Acondicionados: true}

func TestNovoRecolhimento(t *testing.T) {
	origem, a := uuid.New(), uuid.New()
	r, ids, err := NovoRecolhimento(Recolhimento{OrigemID: origem, Motivo: " regular ", Preparo: preparoOk, Historico: " Secretaria criada em 1990 "},
		[]uuid.UUID{a, a, uuid.Nil})
	if err != nil || len(ids) != 1 || r.Motivo != MotivoRegular || r.Situacao != RecolhimentoEnviado || r.ID == uuid.Nil ||
		r.Historico != "Secretaria criada em 1990" {
		t.Fatalf("válido: %+v %v %v", r, ids, err)
	}
	if _, _, err := NovoRecolhimento(Recolhimento{OrigemID: origem, Motivo: MotivoExtincao, Fundamento: "Lei 11.000/2026", Preparo: preparoOk},
		[]uuid.UUID{a}); err != nil {
		t.Fatalf("extinção: %v", err)
	}
	muitos := make([]uuid.UUID, MaxItensRecolhimento+1)
	for i := range muitos {
		muitos[i] = uuid.New()
	}
	incompleto := preparoOk
	incompleto.Higienizados = false
	casos := map[string]struct {
		r   Recolhimento
		ids []uuid.UUID
	}{
		"sem origem":            {Recolhimento{Motivo: MotivoRegular, Preparo: preparoOk}, []uuid.UUID{a}},
		"motivo":                {Recolhimento{OrigemID: origem, Motivo: "DOACAO", Preparo: preparoOk}, []uuid.UUID{a}},
		"especial sem ato":      {Recolhimento{OrigemID: origem, Motivo: MotivoDesestatizacao, Preparo: preparoOk}, []uuid.UUID{a}},
		"fundamento longo":      {Recolhimento{OrigemID: origem, Motivo: MotivoExtincao, Fundamento: strings.Repeat("x", 501), Preparo: preparoOk}, []uuid.UUID{a}},
		"preparo (art. 53)":     {Recolhimento{OrigemID: origem, Motivo: MotivoRegular, Preparo: incompleto}, []uuid.UUID{a}},
		"histórico longo":       {Recolhimento{OrigemID: origem, Motivo: MotivoRegular, Preparo: preparoOk, Historico: strings.Repeat("x", 8001)}, []uuid.UUID{a}},
		"observação longa":      {Recolhimento{OrigemID: origem, Motivo: MotivoRegular, Preparo: preparoOk, Observacao: strings.Repeat("x", 1001)}, []uuid.UUID{a}},
		"sem unidades":          {Recolhimento{OrigemID: origem, Motivo: MotivoRegular, Preparo: preparoOk}, nil},
		"unidades além do teto": {Recolhimento{OrigemID: origem, Motivo: MotivoRegular, Preparo: preparoOk}, muitos},
	}
	for nome, c := range casos {
		var v *ErrValidacao
		if _, _, err := NovoRecolhimento(c.r, c.ids); !errors.As(err, &v) {
			t.Errorf("%s: %v", nome, err)
		}
	}
}

func TestPodeRecolher(t *testing.T) {
	hoje := dia("2026-10-15")
	perm := Serie{CorrenteAnos: anos(2), IntermAnos: anos(3), Destinacao: DestinacaoGuardaPermanente, Vigente: true}
	elim := perm
	elim.Destinacao = DestinacaoEliminacao
	sem := perm
	sem.Destinacao = ""
	evento := perm
	evento.CorrenteAnos, evento.CorrenteCondicao = nil, "Enquanto estiver vigorando"
	velha := Unidade{Fase: FaseIntermediaria, DataFinal: dia("2018-12-31")}
	nova := Unidade{Fase: FaseIntermediaria, DataFinal: dia("2025-12-31")}
	recolhida := Unidade{Fase: FasePermanente, Destinacao: DestinacaoGuardaPermanente, RecolhidaEm: &hoje}
	amostra := Unidade{Fase: FasePermanente, Destinacao: DestinacaoGuardaPermanente}
	casos := []struct {
		nome   string
		u      Unidade
		s      Serie
		motivo string
		err    error
	}{
		{"intermediária com prazo cumprido", velha, perm, MotivoRegular, nil},
		{"amostra permanente no órgão (art. 42)", amostra, elim, MotivoRegular, nil},
		{"extinção, ainda no corrente", Unidade{Fase: FaseCorrente, DataFinal: dia("2025-12-31")}, perm, MotivoExtincao, nil},
		{"desestatização, prazo correndo", nova, perm, MotivoDesestatizacao, nil},
		{"já recolhida", recolhida, perm, MotivoRegular, ErrRecolhida},
		{"eliminada", Unidade{Fase: FaseEliminado}, perm, MotivoRegular, ErrEliminado},
		{"destino eliminação", velha, elim, MotivoRegular, ErrNaoPermanente},
		{"destino eliminação na extinção", velha, elim, MotivoExtincao, ErrNaoPermanente},
		{"destino indefinido", velha, sem, MotivoRegular, ErrDestinacaoIndefinida},
		{"regular a partir do corrente", Unidade{Fase: FaseCorrente, DataFinal: dia("2010-12-31")}, perm, MotivoRegular, ErrRecolherDoCorrente},
		{"prazo intermediário correndo", nova, perm, MotivoRegular, ErrPrazoNaoCumprido},
		{"aguarda o evento", velha, evento, MotivoRegular, ErrAguardaEvento},
	}
	for _, c := range casos {
		if err := PodeRecolher(c.u, c.s, c.motivo, hoje); !errors.Is(err, c.err) {
			t.Errorf("%s: %v", c.nome, err)
		}
	}
	if PodeDecidir(Recolhimento{Situacao: RecolhimentoEnviado}) != nil || !errors.Is(PodeDecidir(Recolhimento{Situacao: RecolhimentoRecebido}), ErrRecolhimentoSituacao) {
		t.Fatal("decisão")
	}
	// Recolhida: a situação é a guarda do Arquivo Público.
	if p := CalcularPrazos(recolhida, perm, true, hoje); p.Situacao != SituacaoRecolhida {
		t.Fatalf("situação recolhida: %s", p.Situacao)
	}
}
