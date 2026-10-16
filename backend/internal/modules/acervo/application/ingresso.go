package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/yurythx/sismarq/internal/domain/events"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// payloadTramite é o envelope de tramite.processo.concluido (ADR 039).
type payloadTramite struct {
	ID              string `json:"id"`
	Numero          string `json:"numero"`
	Assunto         string `json:"assunto"`
	CodigoTTDD      string `json:"codigo_ttdd"`
	UnidadeOrigemID string `json:"unidade_origem_id"`
	Sigilo          string `json:"sigilo"`
	CreatedAt       string `json:"created_at"`
	ConcluidoAt     string `json:"concluido_at"`
}

// HandleTramite recebe o processo encerrado e, se classificado pela TTDD,
// cria a unidade de arquivamento. Sem série, ignora (não é falha). Já
// ingressado, é idempotente. Reabrir o processo no Trâmite não apaga a UA.
func (s *Service) HandleTramite(ctx context.Context, ev events.Event) error {
	var p payloadTramite
	if len(ev.Payload) > 0 {
		if err := json.Unmarshal(ev.Payload, &p); err != nil {
			return err
		}
	}
	procID, err := uuid.Parse(p.ID)
	if err != nil {
		return err
	}
	origem, err := uuid.Parse(p.UnidadeOrigemID)
	if err != nil {
		return err
	}
	pe := domain.ProcessoEncerrado{
		ID: procID, Numero: p.Numero, Assunto: p.Assunto, CodigoTTDD: p.CodigoTTDD,
		UnidadeOrigemID: origem, Sigilo: p.Sigilo,
		CreatedAt: parseInstante(p.CreatedAt), ConcluidoAt: parseInstante(p.ConcluidoAt),
	}
	if pe.ConcluidoAt.IsZero() {
		pe.ConcluidoAt = s.hoje()
	}
	u, err := domain.UnidadeDoTramite(pe)
	if errors.Is(err, domain.ErrSemClassificacao) {
		return nil
	}
	if err != nil {
		return err
	}
	_, _, err = s.repo.GetUnidadePorTramite(ctx, s.pool, procID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, domain.ErrNaoEncontrada) {
		return err
	}
	desc := "Ingresso do processo " + p.Numero + " encerrado no Trâmite"
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		_, err := s.criarTx(ctx, tx, auth.Identity{}, u, domain.MovCadastro, desc)
		return err
	})
	if uniqueViolation(err) {
		return nil
	}
	return err
}

func parseInstante(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, time.DateOnly} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func uniqueViolation(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}
