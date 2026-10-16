package application_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/events"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
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
}
