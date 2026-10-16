package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Arquivos privados de interesse público e social (arts. 55 a 58; ADR 036).
// O Arquivo Público cadastra o arquivo de uma pessoa física ou jurídica, a
// Comissão Especial dá o parecer e o decreto o declara de interesse público
// e social. A declaração não transfere a guarda (§ 4º): o titular continua
// responsável, salvo doação ou depósito revogável ao Arquivo (§ 5º). Na
// alienação, o município tem preferência por 60 dias (art. 58).

var (
	ErrPrivadoNaoExiste   = errors.New("acervo: arquivo privado não encontrado")
	ErrPrivadoSituacao    = errors.New("acervo: o arquivo privado não está na etapa que a operação exige")
	ErrAlienacaoPendente  = errors.New("acervo: já há uma alienação aguardando a manifestação do município")
	ErrAlienacaoNaoExiste = errors.New("acervo: alienação não encontrada")
	ErrAlienacaoDecidida  = errors.New("acervo: a manifestação do município já foi registrada")
	ErrPreferenciaVencida = errors.New("acervo: passaram os 60 dias do direito de preferência (art. 58)")
)

// Situação do arquivo privado.
const (
	PrivadoCadastrado = "CADASTRADO"    // aguardando o parecer da Comissão Especial
	PrivadoParecer    = "COM_PARECER"   // parecer favorável: aguardando o decreto
	PrivadoNaoDeclar  = "NAO_DECLARADO" // parecer contrário
	PrivadoDeclarado  = "DECLARADO"     // decreto ou § 3º
)

// Guarda do acervo declarado (§§ 4º e 5º).
const (
	GuardaTitular    = "TITULAR"
	GuardaDoado      = "DOADO"
	GuardaDepositado = "DEPOSITADO" // a título revogável
)

// Titular: pessoa física ou jurídica (art. 55).
const (
	TitularPF = "PF"
	TitularPJ = "PJ"
)

// Manifestação do município na alienação (art. 58).
const (
	AlienacaoPendente     = "PENDENTE"
	AlienacaoInteresse    = "INTERESSE"
	AlienacaoSemInteresse = "SEM_INTERESSE"
)

// PrazoPreferencia: dias para o município manifestar interesse (art. 58).
const PrazoPreferencia = 60

// Privado é o arquivo privado cadastrado.
type Privado struct {
	ID          uuid.UUID `json:"id"`
	Numero      string    `json:"numero"`
	Titular     string    `json:"titular"`
	TipoTitular string    `json:"tipo_titular"`
	// Documento (CPF/CNPJ) e contato: dado pessoal — só o Arquivo Público
	// lê o cadastro (LGPD).
	Documento string `json:"documento"`
	Contato   string `json:"contato"`
	// Descricao: os conjuntos e a relevância para a história, a cultura e o
	// desenvolvimento científico e tecnológico do município (art. 56).
	Descricao   string `json:"descricao"`
	DatasLimite string `json:"datas_limite"`
	Mensuracao  string `json:"mensuracao"`
	LocalGuarda string `json:"local_guarda"`
	// ServicoPublico: entidade privada encarregada de serviço público
	// municipal — de interesse público e social por força do § 3º.
	ServicoPublico bool   `json:"servico_publico"`
	Situacao       string `json:"situacao"`

	ParecerComissao  string     `json:"parecer_comissao"`
	ParecerTexto     string     `json:"parecer_texto"`
	ParecerFavoravel *bool      `json:"parecer_favoravel"`
	ParecerEm        *time.Time `json:"parecer_em"`

	DecretoNumero string     `json:"decreto_numero"`
	DecretoEm     *time.Time `json:"decreto_em"`
	DecretoDiario string     `json:"decreto_diario"`

	// Acesso do público, com a autorização do titular (§ 2º).
	AcessoAutorizado bool   `json:"acesso_autorizado"`
	AcessoCondicoes  string `json:"acesso_condicoes"`

	Guarda      string     `json:"guarda"`
	GuardaEm    *time.Time `json:"guarda_em"`
	GuardaTermo string     `json:"guarda_termo"`
	// Convenio de assistência técnica (art. 57).
	Convenio string `json:"convenio"`

	CriadoEm   time.Time   `json:"criado_em"`
	Eventos    []Evento    `json:"eventos"`
	Alienacoes []Alienacao `json:"alienacoes"`
}

// Evento é o histórico do arquivo privado.
type Evento struct {
	Tipo      string     `json:"tipo"`
	Descricao string     `json:"descricao"`
	Por       *uuid.UUID `json:"por"`
	Em        time.Time  `json:"em"`
}

// Alienacao é a notificação de venda ao município (art. 58).
type Alienacao struct {
	ID           uuid.UUID  `json:"id"`
	NotificadoEm time.Time  `json:"notificado_em"`
	Adquirente   string     `json:"adquirente"`
	Condicoes    string     `json:"condicoes"`
	Decisao      string     `json:"decisao"`
	Fundamento   string     `json:"fundamento"`
	DecididoEm   *time.Time `json:"decidido_em"`
	// PrazoFinal: último dia para o município manifestar interesse.
	PrazoFinal time.Time `json:"prazo_final"`
	// Vencida: pendente depois do prazo — a preferência caducou.
	Vencida bool `json:"vencida"`
}

// Completar calcula o prazo e se venceu.
func (a *Alienacao) Completar(hoje time.Time) {
	a.PrazoFinal = a.NotificadoEm.AddDate(0, 0, PrazoPreferencia)
	a.Vencida = a.Decisao == AlienacaoPendente && hoje.After(a.PrazoFinal)
}

func limite(campo, v string, max int, obrigatorio bool) error {
	switch {
	case obrigatorio && v == "":
		return invalido(campo, "obrigatório, até "+strconv.Itoa(max)+" caracteres")
	case utf8.RuneCountInString(v) > max:
		return invalido(campo, "até "+strconv.Itoa(max)+" caracteres")
	}
	return nil
}

// ValidarCadastro confere os dados do arquivo e do titular.
func (p *Privado) ValidarCadastro() error {
	for _, c := range []*string{&p.Titular, &p.Documento, &p.Contato, &p.Descricao, &p.DatasLimite, &p.Mensuracao, &p.LocalGuarda, &p.Convenio, &p.AcessoCondicoes} {
		*c = strings.TrimSpace(*c)
	}
	p.TipoTitular = strings.ToUpper(strings.TrimSpace(p.TipoTitular))
	if p.TipoTitular != TitularPF && p.TipoTitular != TitularPJ {
		return invalido("tipo_titular", "use PF (pessoa física) ou PJ (pessoa jurídica)")
	}
	if p.ServicoPublico && p.TipoTitular != TitularPJ {
		return invalido("servico_publico", "só entidade privada (pessoa jurídica) presta serviço público (§ 3º)")
	}
	for _, c := range []struct {
		campo, v string
		max      int
		obrig    bool
	}{
		{"titular", p.Titular, 200, true}, {"documento", p.Documento, 20, false}, {"contato", p.Contato, 500, false},
		{"descricao", p.Descricao, 8000, true}, {"datas_limite", p.DatasLimite, 100, false}, {"mensuracao", p.Mensuracao, 200, false},
		{"local_guarda", p.LocalGuarda, 500, false}, {"convenio", p.Convenio, 1000, false}, {"acesso_condicoes", p.AcessoCondicoes, 2000, false},
	} {
		if err := limite(c.campo, c.v, c.max, c.obrig); err != nil {
			return err
		}
	}
	return nil
}

// NovoPrivado prepara o cadastro. A entidade de serviço público já nasce
// declarada (§ 3º); as demais aguardam o parecer da Comissão Especial.
func NovoPrivado(p Privado) (Privado, error) {
	if err := p.ValidarCadastro(); err != nil {
		return p, err
	}
	p.ID, p.Situacao, p.Guarda = uuid.New(), PrivadoCadastrado, GuardaTitular
	if p.ServicoPublico {
		p.Situacao = PrivadoDeclarado
	}
	return p, nil
}

// PodeParecer: só o cadastrado (ainda sem parecer).
func PodeParecer(p Privado, comissao, texto string, em, hoje time.Time) error {
	if p.Situacao != PrivadoCadastrado {
		return ErrPrivadoSituacao
	}
	if err := limite("parecer_comissao", strings.TrimSpace(comissao), 2000, true); err != nil {
		return err
	}
	if err := limite("parecer_texto", strings.TrimSpace(texto), 8000, true); err != nil {
		return err
	}
	if em.IsZero() || em.After(hoje) {
		return invalido("parecer_em", "informe a data do parecer (não futura)")
	}
	return nil
}

// PodeDeclarar: o decreto vem depois do parecer favorável (art. 56, § 1º).
func PodeDeclarar(p Privado, numero, diario string, em, hoje time.Time) error {
	if p.Situacao != PrivadoParecer {
		return ErrPrivadoSituacao
	}
	if err := limite("decreto_numero", strings.TrimSpace(numero), 40, true); err != nil {
		return err
	}
	if err := limite("decreto_diario", strings.TrimSpace(diario), 200, false); err != nil {
		return err
	}
	if em.IsZero() || em.After(hoje) || em.Before(*p.ParecerEm) {
		return invalido("decreto_em", "entre o parecer e hoje")
	}
	return nil
}

// PodeMudarGuarda: só o declarado (§§ 4º e 5º). Doação é definitiva; o
// depósito é revogável e volta ao titular (do titular, só se doa ou
// deposita: as outras combinações já caem nas regras acima).
func PodeMudarGuarda(p Privado, guarda, termo string, em, hoje time.Time) error {
	if p.Situacao != PrivadoDeclarado {
		return ErrPrivadoSituacao
	}
	switch {
	case p.Guarda == GuardaDoado:
		return invalido("guarda", "o acervo doado ao Arquivo Público não volta ao titular")
	case guarda != GuardaDoado && guarda != GuardaDepositado && guarda != GuardaTitular:
		return invalido("guarda", "use DOADO, DEPOSITADO ou TITULAR")
	case guarda == p.Guarda:
		return invalido("guarda", "o acervo já está nesta guarda")
	}
	if err := limite("termo", strings.TrimSpace(termo), 1000, true); err != nil {
		return err
	}
	if em.IsZero() || em.After(hoje) {
		return invalido("em", "informe a data (não futura)")
	}
	return nil
}

// PodeNotificarAlienacao: só o declarado, sem outra alienação pendente
// dentro do prazo.
func PodeNotificarAlienacao(p Privado, a Alienacao, hoje time.Time) error {
	if p.Situacao != PrivadoDeclarado || p.Guarda == GuardaDoado {
		return ErrPrivadoSituacao
	}
	for _, x := range p.Alienacoes {
		x.Completar(hoje)
		if x.Decisao == AlienacaoPendente && !x.Vencida {
			return ErrAlienacaoPendente
		}
	}
	if err := limite("adquirente", strings.TrimSpace(a.Adquirente), 200, false); err != nil {
		return err
	}
	if err := limite("condicoes", strings.TrimSpace(a.Condicoes), 2000, false); err != nil {
		return err
	}
	if a.NotificadoEm.IsZero() || a.NotificadoEm.After(hoje) {
		return invalido("notificado_em", "informe a data da notificação (não futura)")
	}
	return nil
}

// PodeDecidirAlienacao: a manifestação vem dentro dos 60 dias; depois
// deles, a preferência caducou.
func PodeDecidirAlienacao(a Alienacao, decisao, fundamento string, hoje time.Time) error {
	a.Completar(hoje)
	switch {
	case a.Decisao != AlienacaoPendente:
		return ErrAlienacaoDecidida
	case a.Vencida:
		return ErrPreferenciaVencida
	case decisao != AlienacaoInteresse && decisao != AlienacaoSemInteresse:
		return invalido("decisao", "use INTERESSE ou SEM_INTERESSE")
	}
	return limite("fundamento", strings.TrimSpace(fundamento), 2000, true)
}
