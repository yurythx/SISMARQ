// Package acervo é o plug-in Acervo do SISMARQ (ADR 032): as unidades de
// arquivamento da administração municipal nas três idades (Lei Municipal nº
// 10.301/2019, art. 18), com classificação pela TTDD do Atlas, prazos de
// guarda, localização física, bloqueio judicial e as travas da lei.
package acervo

import (
	"time"

	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/modules/acervo/transport"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/kernel"
	"github.com/yurythx/sismarq/internal/platform/messaging"
	"github.com/yurythx/sismarq/internal/platform/modkit"
)

// Key é a chave do módulo.
const Key = "acervo"

// IntervaloRecalculo: os prazos andam com o calendário — recálculo diário.
const IntervaloRecalculo = 24 * time.Hour

// QueueTTDD recebe as mudanças da TTDD (Atlas) para recalcular os prazos.
var QueueTTDD = messaging.QueueSpec{
	Name:        "sismarq.acervo.ttdd",
	DLQName:     "sismarq.acervo.ttdd.dlq",
	RoutingKeys: []string{"atlas.ttdd.#"},
}

// Module é o plug-in.
type Module struct {
	svc      *application.Service
	handlers *transport.Handlers
}

// New constrói o módulo. ttdd é a porta para a Tabela de Temporalidade
// (adaptador do Atlas, ligado em internal/app).
func New(deps modkit.Deps, ttdd domain.TTDD) *Module {
	svc := application.NewService(deps.Pool, infrastructure.NewRepository(), ttdd, deps.Outbox, deps.Logger).WithNotificacoes(deps.Notificacoes)
	return &Module{svc: svc, handlers: transport.NewHandlers(svc, deps.Logger, deps.Config.MaxPageSize)}
}

// Manifest implementa kernel.Plugin.
func (m *Module) Manifest() kernel.Manifest {
	return kernel.Manifest{
		Key:            Key,
		Name:           "Acervo — Gestão Arquivística",
		Description:    "Unidades de arquivamento nas três idades (Lei 10.301/2019): classificação pela TTDD, prazos de guarda, localização física, bloqueio judicial, eliminação e recolhimento ao Arquivo Público.",
		DefaultEnabled: true,
		DependsOn:      []string{"atlas"},
		Icon:           "archive",
		Route:          "/acervo",
		Permissions: []kernel.PermissionInfo{
			{Key: string(auth.PermAcervoRead), Description: "Consultar as unidades de arquivamento (no escopo concedido)"},
			{Key: string(auth.PermAcervoManage), Description: "Cadastrar, alterar, localizar e transferir unidades de arquivamento (agente setorial, art. 9º)"},
			{Key: string(auth.PermAcervoBloqueio), Description: "Registrar e liberar bloqueio judicial (arts. 32, II, e 44)"},
			{Key: string(auth.PermAcervoAvaliacao), Description: "Conduzir a eliminação como membro da CSPAD: relação, ata, edital, amostragem e termo (arts. 31 a 40)"},
			{Key: string(auth.PermAcervoCentral), Description: "Arquivo Público (órgão central): locais, comissões, autorização de eliminação e entrada dos recolhimentos (arts. 5º, 19 e 53)"},
		},
	}
}

// RegisterRoutes implementa kernel.RouteProvider.
func (m *Module) RegisterRoutes(r kernel.Routes) { m.handlers.RegisterRoutes(r.Authed) }

// Workers implementa kernel.WorkerProvider.
func (m *Module) Workers() []kernel.Worker {
	return []kernel.Worker{{Name: "prazos", Process: kernel.ProcessWorker, Run: m.svc.RunRecalculo(IntervaloRecalculo)}}
}

// Consumers implementa kernel.ConsumerProvider.
func (m *Module) Consumers() []kernel.Consumer {
	return []kernel.Consumer{{Queue: QueueTTDD, Handler: m.svc.HandleTTDD}}
}
