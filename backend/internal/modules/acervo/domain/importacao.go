package domain

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Importação do acervo existente por planilha (ADR 033). A planilha é o CSV
// que o Excel exporta (separador ";" ou ","), com cabeçalho. Colunas:
// obrigatórias tipo, codigo_classificacao, unidade (sigla ou id), data_inicial
// e data_final; opcionais assunto, quantidade, especificacao, observacoes,
// data_evento, nivel_acesso, hipotese_legal, local (caminho, como aparece na
// tela de locais) e posicao. Datas em dd/mm/aaaa ou aaaa-mm-dd.

// MaxLinhasImportacao limita uma planilha.
const MaxLinhasImportacao = 5000

// ColunasImportacao é o cabeçalho do modelo, na ordem.
var ColunasImportacao = []string{"tipo", "codigo_classificacao", "assunto", "unidade", "data_inicial", "data_final",
	"quantidade", "especificacao", "observacoes", "data_evento", "nivel_acesso", "hipotese_legal", "local", "posicao"}

var obrigatorias = []string{"tipo", "codigo_classificacao", "unidade", "data_inicial", "data_final"}

// ErrPlanilha: a planilha não pode ser lida (cabeçalho, formato, tamanho).
var ErrPlanilha = errors.New("acervo: planilha inválida")

// LinhaImportada é uma linha lida: a unidade (sem a unidade produtora, que
// vem da referência) e as referências a resolver no banco.
type LinhaImportada struct {
	Linha      int
	Unidade    Unidade
	UnidadeRef string
	LocalRef   string
}

// ErroLinha aponta o problema de uma linha da planilha.
type ErroLinha struct {
	Linha  int    `json:"linha"`
	Motivo string `json:"motivo"`
}

func dataPlanilha(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	for _, layout := range []string{"02/01/2006", time.DateOnly} {
		if t, err := time.Parse(layout, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("data %q inválida (use dd/mm/aaaa)", v)
}

var semAcento = strings.NewReplacer("Ê", "E", "ê", "E", "Ú", "U", "ú", "U", "Í", "I", "í", "I", "Á", "A", "á", "A")

// LerPlanilha lê o CSV. Erro de leitura geral devolve ErrPlanilha; erros de
// cada linha vêm na lista (a importação só é aplicada sem nenhum).
func LerPlanilha(conteudo string) ([]LinhaImportada, []ErroLinha, error) {
	conteudo = strings.TrimPrefix(conteudo, "\uFEFF") // BOM do Excel
	primeira, _, _ := strings.Cut(conteudo, "\n")
	r := csv.NewReader(strings.NewReader(conteudo))
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	if strings.Count(primeira, ";") > strings.Count(primeira, ",") {
		r.Comma = ';'
	}
	cab, err := r.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("%w: sem cabeçalho", ErrPlanilha)
	}
	col := map[string]int{}
	for i, c := range cab {
		col[strings.ToLower(strings.TrimSpace(c))] = i
	}
	for _, c := range obrigatorias {
		if _, ok := col[c]; !ok {
			return nil, nil, fmt.Errorf("%w: falta a coluna %q", ErrPlanilha, c)
		}
	}
	var linhas []LinhaImportada
	var erros []ErroLinha
	for n := 2; ; n++ {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("%w: linha %d ilegível", ErrPlanilha, n)
		}
		if len(linhas)+len(erros) >= MaxLinhasImportacao {
			return nil, nil, fmt.Errorf("%w: mais de %d linhas", ErrPlanilha, MaxLinhasImportacao)
		}
		v := func(c string) string {
			if i, ok := col[c]; ok && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		if strings.Join(rec, "") == "" {
			continue // linha em branco
		}
		l, err := linha(n, v)
		if err != nil {
			erros = append(erros, ErroLinha{Linha: n, Motivo: err.Error()})
			continue
		}
		linhas = append(linhas, l)
	}
	if len(linhas)+len(erros) == 0 {
		return nil, nil, fmt.Errorf("%w: nenhuma linha de dados", ErrPlanilha)
	}
	return linhas, erros, nil
}

func linha(n int, v func(string) string) (LinhaImportada, error) {
	u := Unidade{
		Tipo: Tipo(strings.ToUpper(semAcento.Replace(v("tipo")))), CodigoClassificacao: v("codigo_classificacao"), Assunto: v("assunto"),
		Especificacao: v("especificacao"), Observacoes: v("observacoes"), NivelAcesso: NivelAcesso(strings.ToUpper(semAcento.Replace(v("nivel_acesso")))),
		HipoteseLegal: v("hipotese_legal"), Posicao: v("posicao"),
	}
	var err error
	if u.DataInicial, err = dataPlanilha(v("data_inicial")); err != nil {
		return LinhaImportada{}, fmt.Errorf("data_inicial: %w", err)
	}
	if u.DataFinal, err = dataPlanilha(v("data_final")); err != nil {
		return LinhaImportada{}, fmt.Errorf("data_final: %w", err)
	}
	if e := v("data_evento"); e != "" {
		t, err := dataPlanilha(e)
		if err != nil {
			return LinhaImportada{}, fmt.Errorf("data_evento: %w", err)
		}
		u.DataEvento = &t
	}
	if q := v("quantidade"); q != "" {
		if u.Quantidade, err = strconv.Atoi(q); err != nil {
			return LinhaImportada{}, fmt.Errorf("quantidade %q não é um número", q)
		}
	}
	ref := v("unidade")
	if ref == "" {
		return LinhaImportada{}, errors.New("unidade: informe a sigla ou o id da unidade produtora")
	}
	return LinhaImportada{Linha: n, Unidade: u, UnidadeRef: ref, LocalRef: v("local")}, nil
}
