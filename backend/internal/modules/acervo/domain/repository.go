package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Filtro da listagem de unidades de arquivamento.
type Filtro struct {
	Query     string
	Fase      Fase
	Situacao  string
	Codigo    string // prefixo do código de classificação (ex.: "2.0.01")
	UnidadeID *uuid.UUID
	// Restrito: a listagem fica nas Unidades (cobertura de quem pede).
	Restrito bool
	Unidades []uuid.UUID
}

// Calculo é o que se grava do cálculo de prazos (para filtrar e contar).
// A destinação fica na própria Unidade.
type Calculo struct {
	FimCorrente      *time.Time `json:"fim_corrente"`
	FimIntermediaria *time.Time `json:"fim_intermediaria"`
	Situacao         string     `json:"situacao"`
}

// Repository é a porta de persistência do Acervo. Como em todos os
// plug-ins, recebe o executor: o pool na leitura, a transação na escrita.
type Repository interface {
	InsertUnidade(ctx context.Context, db database.DBTX, u Unidade, c Calculo) (string, error)
	GetUnidade(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Unidade, Calculo, error)
	UpdateUnidade(ctx context.Context, db database.DBTX, u Unidade, c Calculo) error
	SetCalculo(ctx context.Context, db database.DBTX, id uuid.UUID, c Calculo) error
	SetFase(ctx context.Context, db database.DBTX, id uuid.UUID, fase Fase, c Calculo) error
	ListUnidades(ctx context.Context, db database.DBTX, f Filtro, p pagination.Params) ([]UnidadeLista, int64, error)
	Painel(ctx context.Context, db database.DBTX, f Filtro) (map[string]int64, error)
	// ParaRecalculo pagina por id as unidades fora das fases finais.
	ParaRecalculo(ctx context.Context, db database.DBTX, depoisDe uuid.UUID, limite int) ([]Unidade, error)

	InsertBloqueio(ctx context.Context, db database.DBTX, b Bloqueio) error
	GetBloqueio(ctx context.Context, db database.DBTX, id uuid.UUID) (Bloqueio, error)
	LiberarBloqueio(ctx context.Context, db database.DBTX, b Bloqueio) error
	ListBloqueios(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]Bloqueio, error)
	TemBloqueioAtivo(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error)

	InsertMovimentacao(ctx context.Context, db database.DBTX, unidade uuid.UUID, m Movimentacao) error
	ListMovimentacoes(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]Movimentacao, error)

	// Fase 2 (ADR 033).
	InsertEmprestimo(ctx context.Context, db database.DBTX, e Emprestimo) error
	GetEmprestimo(ctx context.Context, db database.DBTX, id uuid.UUID) (Emprestimo, error)
	DevolverEmprestimo(ctx context.Context, db database.DBTX, e Emprestimo) error
	EmprestimoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error)
	ListEmprestimos(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]Emprestimo, error)
	ListEmprestimosAbertos(ctx context.Context, db database.DBTX, f Filtro, atrasadosEm *time.Time, p pagination.Params) ([]Emprestimo, int64, error)

	InsertGuia(ctx context.Context, db database.DBTX, g Guia, unidades []uuid.UUID) (string, error)
	GetGuia(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Guia, error)
	SetGuiaSituacao(ctx context.Context, db database.DBTX, g Guia) error
	EmGuiaEnviada(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error)
	ListGuias(ctx context.Context, db database.DBTX, f Filtro, p pagination.Params) ([]Guia, int64, error)

	// Fase 3 (ADR 034).
	SalvarComissao(ctx context.Context, db database.DBTX, c Comissao, nova bool) error
	GetComissao(ctx context.Context, db database.DBTX, id uuid.UUID) (Comissao, error)
	ListComissoes(ctx context.Context, db database.DBTX) ([]Comissao, error)
	ComissaoDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) (Comissao, error)

	InsertListagem(ctx context.Context, db database.DBTX, l Listagem) (string, error)
	UpdateListagem(ctx context.Context, db database.DBTX, l Listagem) error
	GetListagem(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Listagem, error)
	ListListagens(ctx context.Context, db database.DBTX, f Filtro, p pagination.Params) ([]Listagem, int64, error)
	AddItemListagem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID) error
	RemoveItemListagem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID) error
	MarcarItem(ctx context.Context, db database.DBTX, listagem, unidade uuid.UUID, amostra, retirada bool) error
	EmListagemAtiva(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error)
	InsertPedido(ctx context.Context, db database.DBTX, listagem uuid.UUID, p Pedido) error
	DecidirPedido(ctx context.Context, db database.DBTX, p Pedido) error

	// Fase 4 (ADR 035).
	InsertRecolhimento(ctx context.Context, db database.DBTX, r Recolhimento, unidades []uuid.UUID) (string, error)
	GetRecolhimento(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Recolhimento, error)
	ListRecolhimentos(ctx context.Context, db database.DBTX, f Filtro, p pagination.Params) ([]Recolhimento, int64, error)
	// ReceberRecolhimento grava a entrada e devolve o número do registro.
	ReceberRecolhimento(ctx context.Context, db database.DBTX, r Recolhimento) (string, error)
	SetRecolhimentoSituacao(ctx context.Context, db database.DBTX, r Recolhimento) error
	EmRecolhimentoAberto(ctx context.Context, db database.DBTX, unidade uuid.UUID) (bool, error)
	// MarcarRecolhida põe a unidade sob a guarda do Arquivo Público.
	MarcarRecolhida(ctx context.Context, db database.DBTX, u Unidade, c Calculo) error
	// GestoresDaUnidade: quem gere o acervo da unidade por lotação (para a
	// cópia do registro de entrada, art. 5º, VII).
	GestoresDaUnidade(ctx context.Context, db database.DBTX, unidade uuid.UUID) ([]uuid.UUID, error)

	// Fase 4, segunda parte (ADR 036): instruções normativas.
	InsertNorma(ctx context.Context, db database.DBTX, n Norma) error
	GetNorma(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Norma, error)
	ListNormas(ctx context.Context, db database.DBTX, q string, publicas bool, p pagination.Params) ([]Norma, int64, error)
	UpdateNorma(ctx context.Context, db database.DBTX, n Norma) error
	DeleteNorma(ctx context.Context, db database.DBTX, id uuid.UUID) error
	// SalvarRedacao grava a redação em rascunho (nova ou alterada).
	SalvarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, r Redacao, nova bool, por *uuid.UUID) error
	PublicarRedacao(ctx context.Context, db database.DBTX, norma uuid.UUID, r Redacao, por *uuid.UUID) error
	ApagarRascunho(ctx context.Context, db database.DBTX, norma uuid.UUID) error

	// Arquivos privados de interesse público e social.
	InsertPrivado(ctx context.Context, db database.DBTX, p Privado, por *uuid.UUID) (string, error)
	GetPrivado(ctx context.Context, db database.DBTX, id uuid.UUID, paraAtualizar bool) (Privado, error)
	ListPrivados(ctx context.Context, db database.DBTX, q string, p pagination.Params) ([]Privado, int64, error)
	UpdatePrivado(ctx context.Context, db database.DBTX, p Privado) error
	InsertEventoPrivado(ctx context.Context, db database.DBTX, privado uuid.UUID, e Evento) error
	InsertAlienacao(ctx context.Context, db database.DBTX, privado uuid.UUID, a Alienacao, por *uuid.UUID) error
	DecidirAlienacao(ctx context.Context, db database.DBTX, a Alienacao, por *uuid.UUID) error

	// Fase 5 (ADR 037): representantes digitais.
	InsertDigital(ctx context.Context, db database.DBTX, d Digital) (string, error)
	GetDigital(ctx context.Context, db database.DBTX, id uuid.UUID) (Digital, error)
	ListDigitais(ctx context.Context, db database.DBTX, f FiltroDigital, p pagination.Params) ([]Digital, int64, error)
	// ProximoDigital reserva o próximo da fila (ENVIADO, ou PROCESSANDO
	// abandonado antes de travadoAntes) para o worker.
	ProximoDigital(ctx context.Context, db database.DBTX, travadoAntes time.Time) (Digital, bool, error)
	// ConcluirDigital grava o resultado do pipeline (e o texto do OCR).
	ConcluirDigital(ctx context.Context, db database.DBTX, d Digital, texto string) error
	ParaVerificar(ctx context.Context, db database.DBTX, antesDe time.Time, limite int) ([]Digital, error)
	RegistrarVerificacao(ctx context.Context, db database.DBTX, id uuid.UUID, v Verificacao) error
	ListVerificacoes(ctx context.Context, db database.DBTX, id uuid.UUID, limite int) ([]Verificacao, error)
	// EquipeCentral: quem tem acervo:central global (alertas de integridade).
	EquipeCentral(ctx context.Context, db database.DBTX) ([]uuid.UUID, error)

	// ResolverUnidade acha a unidade produtora ativa pela sigla (ou id).
	ResolverUnidade(ctx context.Context, db database.DBTX, ref string) (uuid.UUID, error)

	InsertLocal(ctx context.Context, db database.DBTX, l Local) error
	UpdateLocal(ctx context.Context, db database.DBTX, l Local) error
	GetLocal(ctx context.Context, db database.DBTX, id uuid.UUID) (Local, error)
	ListLocais(ctx context.Context, db database.DBTX) ([]Local, error)
}

// UnidadeLista é a linha da listagem: a unidade e o resultado gravado do
// cálculo.
type UnidadeLista struct {
	Unidade
	Calculo
	LocalCaminho string `json:"local_caminho"`
}
