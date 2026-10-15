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
