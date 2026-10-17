package application_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/yurythx/sismarq/internal/domain/events"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/platform/database"
)

func eventoTramite(t *testing.T, payload any) events.Event {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return events.Event{Type: "tramite.processo.concluido", Payload: raw}
}

func TestHandleTramiteIngressoIdempotente(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	s := e.real()
	repo := infrastructure.NewRepository()
	proc := uuid.New()
	payload := map[string]any{
		"id": proc.String(), "numero": "2025.000099", "assunto": "Alvará",
		"codigo_ttdd": codElim, "unidade_origem_id": e.unidade.String(), "sigilo": "publico",
		"created_at": "2024-02-01T00:00:00Z", "concluido_at": "2025-01-15T12:00:00Z",
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, payload)); err != nil {
		t.Fatal(err)
	}
	got, _, err := repo.GetUnidadePorTramite(ctx, e.pool, proc)
	if err != nil || got.Tipo != domain.TipoProcesso || got.CodigoClassificacao != codElim || got.Assunto != "Alvará" ||
		got.TramiteProcessoID == nil || *got.TramiteProcessoID != proc || got.UnidadeID != e.unidade {
		t.Fatalf("ingresso: %+v %v", got, err)
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, payload)); err != nil {
		t.Fatal(err)
	}

	semSerie := map[string]any{"id": uuid.NewString(), "unidade_origem_id": e.unidade.String(), "numero": "x"}
	if err := s.HandleTramite(ctx, eventoTramite(t, semSerie)); err != nil {
		t.Fatalf("sem série deveria ser ignorado: %v", err)
	}
	if err := s.HandleTramite(ctx, events.Event{Type: "tramite.processo.concluido", Payload: []byte(`{`)}); err == nil {
		t.Fatal("JSON inválido")
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, map[string]any{"id": "não-é-uuid"})); err == nil {
		t.Fatal("id inválido")
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, map[string]any{"id": uuid.NewString(), "unidade_origem_id": "x"})); err == nil {
		t.Fatal("origem inválida")
	}
	lixoData := map[string]any{
		"id": uuid.NewString(), "numero": "2025.000102", "codigo_ttdd": codElim,
		"unidade_origem_id": e.unidade.String(), "created_at": "ontem", "concluido_at": "hoje",
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, lixoData)); err != nil {
		t.Fatalf("instante ilegível usa a data de hoje: %v", err)
	}
	codigoLongo := map[string]any{
		"id": uuid.NewString(), "numero": "2025.000103", "codigo_ttdd": strings.Repeat("9", 40),
		"unidade_origem_id": e.unidade.String(), "created_at": "2024-01-01T00:00:00Z", "concluido_at": "2024-06-01T00:00:00Z",
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, codigoLongo)); err == nil {
		t.Fatal("código da TTDD longo")
	}

	sig := uuid.New()
	restrito := map[string]any{
		"id": sig.String(), "numero": "2025.000100", "codigo_ttdd": codElim,
		"unidade_origem_id": e.unidade.String(), "sigilo": "restrito",
		"created_at": "2024-01-01", "concluido_at": "2024-06-01T00:00:00Z",
	}
	if err := s.HandleTramite(ctx, eventoTramite(t, restrito)); err != nil {
		t.Fatal(err)
	}
	r, _, err := repo.GetUnidadePorTramite(ctx, e.pool, sig)
	if err != nil || r.NivelAcesso != domain.AcessoRestrito || r.HipoteseLegal == "" || r.Assunto != "Processo 2025.000100" {
		t.Fatalf("restrito: %+v %v", r, err)
	}

	corrida := e.svc(duplicaTramite{Repository: infrastructure.NewRepository()})
	mesmo := map[string]any{
		"id": uuid.NewString(), "numero": "2025.000104", "codigo_ttdd": codElim,
		"unidade_origem_id": e.unidade.String(), "created_at": "2024-01-01T00:00:00Z", "concluido_at": "2024-06-01T00:00:00Z",
	}
	if err := corrida.HandleTramite(ctx, eventoTramite(t, mesmo)); err != nil {
		t.Fatalf("unicidade no ingresso é idempotente: %v", err)
	}
}

// duplicaTramite simula dois eventos iguais: a leitura ainda não vê a
// unidade e o INSERT encontra a restrição 23505.
type duplicaTramite struct{ domain.Repository }

func (duplicaTramite) GetUnidadePorTramite(context.Context, database.DBTX, uuid.UUID) (domain.Unidade, domain.Calculo, error) {
	return domain.Unidade{}, domain.Calculo{}, domain.ErrNaoEncontrada
}

func (duplicaTramite) InsertUnidade(context.Context, database.DBTX, domain.Unidade, domain.Calculo) (string, error) {
	return "", &pgconn.PgError{Code: "23505"}
}
