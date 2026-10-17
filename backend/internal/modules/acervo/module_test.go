package acervo

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/config"
	"github.com/yurythx/sismarq/internal/platform/kernel"
	"github.com/yurythx/sismarq/internal/platform/modkit"
	"github.com/yurythx/sismarq/internal/platform/outbox"
	"github.com/yurythx/sismarq/internal/platform/storage/storagetest"
)

type ttddVazio struct{}

func (ttddVazio) Serie(context.Context, string) (domain.Serie, error) {
	return domain.Serie{}, domain.ErrSerieInvalida
}

type allowAll struct{}

func (allowAll) Allow(context.Context, string) (bool, error) { return true, nil }

func TestManifestEConsumers(t *testing.T) {
	m := &Module{}
	man := m.Manifest()
	if man.Key != Key || man.DependsOn[0] != "atlas" || !man.DefaultEnabled {
		t.Fatalf("manifest: %+v", man)
	}
	cs := m.Consumers()
	if len(cs) != 2 || cs[0].Queue.Name != QueueTTDD.Name || cs[1].Queue.Name != QueueTramite.Name {
		t.Fatalf("consumers: %+v", cs)
	}
}

func TestNewRotasEWorkers(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		MaxPageSize: 50,
		MinIO:       config.MinIOConfig{Bucket: "app"},
		AcervoDigital: config.AcervoDigitalConfig{
			Bucket: "worm", Intervalo: time.Second, Fixity: 2 * time.Hour, MaxFileBytes: 1 << 20,
		},
		AcervoBusca: config.AcervoBuscaConfig{Intervalo: time.Second, Recarga: time.Hour, SICURL: "https://sic.exemplo.gov.br"},
	}
	deps := modkit.Deps{
		Config: cfg, Outbox: outbox.NewWriter("test"), Storage: storagetest.New(),
		Logger: log, PublicLimiter: allowAll{},
	}
	m := New(deps, ttddVazio{})
	if m.Manifest().Key != Key || len(m.Manifest().Permissions) < 6 {
		t.Fatalf("manifest composto: %+v", m.Manifest())
	}
	ws := m.Workers()
	if len(ws) != 4 || ws[0].Name != "prazos" || ws[1].Name != "digitalizacao" || ws[2].Name != "integridade" || ws[3].Name != "indice" {
		t.Fatalf("workers: %+v", ws)
	}
	pub, authed := chi.NewRouter(), chi.NewRouter()
	m.RegisterRoutes(kernel.Routes{Public: pub, Authed: authed})
	if pub.Match(chi.NewRouteContext(), http.MethodGet, "/acervo/publico/normas") == false {
		t.Fatal("rota pública das normas")
	}

	cfg.Typesense = config.TypesenseConfig{URL: "http://127.0.0.1:9", APIKey: "k", Timeout: time.Second}
	comIndice := New(deps, ttddVazio{})
	if len(comIndice.Workers()) != 4 {
		t.Fatal("com TYPESENSE_URL o worker do índice permanece")
	}

	cfg.AcervoDigital, cfg.AcervoBusca = config.AcervoDigitalConfig{}, config.AcervoBuscaConfig{}
	zero := New(deps, ttddVazio{})
	if n := len(zero.Workers()); n != 1 || zero.Workers()[0].Name != "prazos" {
		t.Fatalf("intervalo zero só o recálculo: %d", n)
	}
	cfg.AcervoBusca = config.AcervoBuscaConfig{Intervalo: time.Second}
	semRecarga := New(deps, ttddVazio{})
	if n := len(semRecarga.Workers()); n != 2 || semRecarga.Workers()[1].Name != "indice" {
		t.Fatalf("busca sem recarga: %d", n)
	}
}
