package domain

import "testing"

// Amostra C1 da homologação (docs/acervo/HOMOLOGACAO.md): códigos reais da
// TTDD da SEMAD e o que o painel deve mostrar em 16/10/2026. As caixas em
// deploy/homolog/seed-acervo.sql usam os mesmos números.
func TestAmostraC1PrazosTTDD(t *testing.T) {
	hoje := dia("2026-10-16")
	um := anos(1)
	evento := dia("2024-06-01")
	casos := []struct {
		nome, codigo, situacao string
		u                      Unidade
		s                      Serie
		bloqueada              bool
	}{
		{
			nome: "HOM-UA-01 organogramas a avaliar", codigo: "2.0.01.00.01", situacao: SituacaoAvaliarEliminacao,
			u: Unidade{Fase: FaseIntermediaria, DataFinal: dia("2015-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoEliminacao, Vigente: true},
		},
		{
			nome: "HOM-UA-02 pregão a recolher", codigo: "2.0.02.00.07", situacao: SituacaoRecolher,
			u: Unidade{Fase: FaseIntermediaria, DataFinal: dia("2015-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoGuardaPermanente, Vigente: true},
		},
		{
			nome: "HOM-UA-03 permanente", codigo: "2.0.01.03.01", situacao: SituacaoPermanente,
			u: Unidade{Fase: FasePermanente, DataFinal: dia("2010-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoGuardaPermanente, Vigente: true},
		},
		{
			nome: "HOM-UA-04 contrato sem evento", codigo: "2.0.03.00.12", situacao: SituacaoAguardaEvento,
			u: Unidade{Fase: FaseCorrente, DataFinal: dia("2023-12-31")},
			s: Serie{CorrenteCondicao: "Enquanto estiver vigorando", IntermAnos: um, Destinacao: DestinacaoGuardaPermanente, Vigente: true},
		},
		{
			nome: "HOM-UA-05 contrato encerrado", codigo: "2.0.03.00.12", situacao: SituacaoTransferir,
			u: Unidade{Fase: FaseCorrente, DataFinal: dia("2023-12-31"), DataEvento: &evento},
			s: Serie{CorrenteCondicao: "Enquanto estiver vigorando", IntermAnos: um, Destinacao: DestinacaoGuardaPermanente, Vigente: true},
		},
		{
			nome: "HOM-UA-06 ainda no corrente", codigo: "2.0.01.00.01", situacao: SituacaoNoCorrente,
			u: Unidade{Fase: FaseCorrente, DataFinal: dia("2026-06-30")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoEliminacao, Vigente: true},
		},
		{
			nome: "HOM-UA-07 bloqueio judicial", codigo: "2.0.01.00.01", situacao: SituacaoBloqueada,
			u: Unidade{Fase: FaseIntermediaria, DataFinal: dia("2015-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoEliminacao, Vigente: true},
			bloqueada: true,
		},
		{
			nome: "HOM-UA-08 sindicância a transferir", codigo: "2.0.05.02.00", situacao: SituacaoTransferir,
			u: Unidade{Fase: FaseCorrente, DataFinal: dia("2022-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoGuardaPermanente, Vigente: true},
		},
		{
			nome: "HOM-UA-09 SIC a avaliar", codigo: "2.0.03.02.00", situacao: SituacaoAvaliarEliminacao,
			u: Unidade{Fase: FaseIntermediaria, DataFinal: dia("2018-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoEliminacao, Vigente: true},
		},
		{
			nome: "HOM-UA-10 contas TCE a avaliar", codigo: "2.0.01.00.05", situacao: SituacaoAvaliarEliminacao,
			u: Unidade{Fase: FaseIntermediaria, DataFinal: dia("2014-12-31")},
			s: Serie{CorrenteAnos: um, IntermAnos: um, Destinacao: DestinacaoEliminacao, Vigente: true},
		},
	}
	for _, c := range casos {
		p := CalcularPrazos(c.u, c.s, c.bloqueada, hoje)
		if p.Situacao != c.situacao {
			t.Errorf("%s (%s): situacao=%s queria %s (prazos=%+v)", c.nome, c.codigo, p.Situacao, c.situacao, p)
		}
		if c.nome == "HOM-UA-03 permanente" {
			if err := PodeEliminar(c.u, c.s, c.bloqueada, hoje); err != ErrPermanente {
				t.Errorf("permanente não se elimina: %v", err)
			}
		}
		if c.bloqueada {
			if err := PodeEliminar(c.u, c.s, true, hoje); err != ErrBloqueioJudicial {
				t.Errorf("bloqueio não se elimina: %v", err)
			}
		}
	}
}
