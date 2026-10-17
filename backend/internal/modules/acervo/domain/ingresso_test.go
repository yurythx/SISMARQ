package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

func TestUnidadeDoSEI(t *testing.T) {
	origem := uuid.New()
	_, err := UnidadeDoSEI(ProcessoSEI{NumeroSEI: "SEI-1", UnidadeOrigemID: origem, CodigoTTDD: "9.9.01.00.01"})
	if !errors.Is(err, ErrSEIAdiado) {
		t.Fatalf("adiado: %v", err)
	}
	if _, err := UnidadeDoSEI(ProcessoSEI{}); err == nil {
		t.Fatal("sem unidade")
	}
	if _, err := UnidadeDoSEI(ProcessoSEI{UnidadeOrigemID: origem, NumeroSEI: strings.Repeat("x", 81)}); err == nil {
		t.Fatal("número longo")
	}
	if _, err := UnidadeDoSEI(ProcessoSEI{UnidadeOrigemID: origem}); err == nil {
		t.Fatal("sem número SEI nem protocolo")
	}
	_, err = UnidadeDoSEI(ProcessoSEI{UnidadeOrigemID: origem, Protocolo: "P-1"})
	if !errors.Is(err, ErrSEIAdiado) {
		t.Fatalf("protocolo sem SEI ainda é adiado: %v", err)
	}
}

func TestUnidadeDoTramite(t *testing.T) {
	id, origem := uuid.New(), uuid.New()
	ini := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	fim := time.Date(2025, 3, 2, 9, 0, 0, 0, time.UTC)
	u, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "2025.000001", Assunto: "Licença", CodigoTTDD: "9.9.01.00.01",
		UnidadeOrigemID: origem, CreatedAt: ini, ConcluidoAt: fim, Sigilo: "publico",
	})
	if err != nil || u.Tipo != TipoProcesso || u.UnidadeID != origem || u.TramiteProcessoID == nil || *u.TramiteProcessoID != id ||
		u.NivelAcesso != AcessoPublico || u.Assunto != "Licença" || u.Quantidade != 1 {
		t.Fatalf("público: %+v %v", u, err)
	}
	if !u.DataInicial.Equal(ini) || !u.DataFinal.Equal(fim) {
		t.Fatalf("datas: %v %v", u.DataInicial, u.DataFinal)
	}

	semSerie, err := UnidadeDoTramite(ProcessoEncerrado{ID: id, UnidadeOrigemID: origem, Numero: "n"})
	if !errors.Is(err, ErrSemClassificacao) || semSerie.Tipo != "" {
		t.Fatalf("sem série: %+v %v", semSerie, err)
	}

	sig, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "2025.000002", CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem,
		CreatedAt: ini, ConcluidoAt: fim, Sigilo: "sigiloso",
	})
	if err != nil || sig.NivelAcesso != AcessoSigiloso || sig.HipoteseLegal == "" || sig.Assunto != "Processo 2025.000002" {
		t.Fatalf("sigiloso: %+v %v", sig, err)
	}

	inv := []ProcessoEncerrado{
		{CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem},
		{ID: id, CodigoTTDD: "9.9.01.00.01"},
	}
	for _, p := range inv {
		if _, err := UnidadeDoTramite(p); err == nil {
			t.Errorf("deveria recusar %+v", p)
		}
	}

	soFim, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "n", CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem, ConcluidoAt: fim,
	})
	if err != nil || !soFim.DataInicial.Equal(fim) || !soFim.DataFinal.Equal(fim) {
		t.Fatalf("só conclusão: %+v %v", soFim, err)
	}
	soIni, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "n", CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem, CreatedAt: ini,
	})
	if err != nil || !soIni.DataInicial.Equal(ini) || !soIni.DataFinal.Equal(ini) {
		t.Fatalf("só abertura: %+v %v", soIni, err)
	}
	invertido, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "n", CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem, CreatedAt: fim, ConcluidoAt: ini,
	})
	if err != nil || invertido.DataFinal.Before(invertido.DataInicial) {
		t.Fatalf("datas invertidas: %+v %v", invertido, err)
	}
	longo := strings.Repeat("ã", 520)
	trunc, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "n", Assunto: longo, CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem, CreatedAt: ini, ConcluidoAt: fim,
	})
	if err != nil || utf8.RuneCountInString(trunc.Assunto) != 500 {
		t.Fatalf("assunto truncado: %d %v", utf8.RuneCountInString(trunc.Assunto), err)
	}
	if _, err := UnidadeDoTramite(ProcessoEncerrado{
		ID: id, Numero: "n", CodigoTTDD: "9.9.01.00.01", UnidadeOrigemID: origem,
	}); err == nil {
		t.Fatal("sem datas")
	}
}

func TestAfirmaAssinaturaICP(t *testing.T) {
	vazio := Digital{}.ComEstadoPAdES()
	if vazio.Pades != PadesAusente || AfirmaAssinaturaICP(vazio) {
		t.Fatalf("ausente: %+v", vazio)
	}
	if AfirmaAssinaturaICP(Digital{Pades: "ICP-BRASIL"}) != true {
		t.Fatal("ICP gravada deveria afirmar")
	}
}
