package application

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Instruções normativas do Arquivo Público (art. 54; ADR 036). Só o Arquivo
// Público (acervo:central, global) escreve; qualquer pessoa consulta o que
// foi publicado, com o histórico das redações.

// Eventos das normas (outbox).
const (
	EventNormaCriada    = "acervo.norma.criada"
	EventNormaRascunho  = "acervo.norma.rascunho"
	EventNormaPublicada = "acervo.norma.publicada"
	EventNormaDescarte  = "acervo.norma.rascunho_descartado"
	EventNormaRevogada  = "acervo.norma.revogada"
)

// MotivoOriginal é o motivo padrão da primeira redação.
const MotivoOriginal = "Redação original"

func central(identity auth.Identity) error {
	if !auth.CanGlobal(identity, auth.PermAcervoCentral) {
		return domain.ErrForaDoEscopo
	}
	return nil
}

func (s *Service) registrarNorma(ctx context.Context, tx pgx.Tx, evento string, n domain.Norma, antes, depois any) error {
	if err := s.outbox.Write(ctx, tx, evento, "acervo_norma", n.ID.String(), uuid.Nil,
		map[string]any{"id": n.ID, "numero": n.Numero, "situacao": n.Situacao}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_norma", n.ID.String(), antes, depois))
}

// CriarNorma abre a instrução normativa com a primeira redação em rascunho.
func (s *Service) CriarNorma(ctx context.Context, identity auth.Identity, numero string, r domain.Redacao) (domain.Norma, error) {
	if err := central(identity); err != nil {
		return domain.Norma{}, MapError(err)
	}
	numero, err := domain.ValidarNumero(numero)
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	if strings.TrimSpace(r.Motivo) == "" {
		r.Motivo = MotivoOriginal
	}
	if err := domain.ValidarRedacao(&r); err != nil {
		return domain.Norma{}, MapError(err)
	}
	n := domain.Norma{ID: uuid.New(), Numero: numero, Situacao: domain.NormaRascunho, CriadaEm: s.agora().UTC()}
	r.Versao = 1
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.repo.InsertNorma(ctx, tx, n); err != nil {
			return err
		}
		if err := s.repo.SalvarRedacao(ctx, tx, n.ID, r, true, ator(identity)); err != nil {
			return err
		}
		return s.registrarNorma(ctx, tx, EventNormaCriada, n, nil, map[string]any{"numero": n.Numero, "ementa": r.Ementa})
	})
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	return s.Norma(ctx, identity, n.ID)
}

// carregarNorma lê a norma para alterá-la (só o Arquivo Público).
func (s *Service) carregarNorma(ctx context.Context, tx pgx.Tx, identity auth.Identity, id uuid.UUID) (domain.Norma, error) {
	if err := central(identity); err != nil {
		return domain.Norma{}, err
	}
	return s.repo.GetNorma(ctx, tx, id, true)
}

// SalvarRascunho grava a redação em rascunho: altera o rascunho existente
// ou abre a nova redação da norma em vigor (a alteração futura da norma).
// A redação publicada nunca muda.
func (s *Service) SalvarRascunho(ctx context.Context, identity auth.Identity, id uuid.UUID, r domain.Redacao) (domain.Norma, error) {
	if err := domain.ValidarRedacao(&r); err != nil {
		return domain.Norma{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		n, err := s.carregarNorma(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		if n.Situacao == domain.NormaRevogada {
			return domain.ErrNormaSituacao
		}
		nova := n.Rascunho() == nil
		r.Versao = n.Versao + 1
		if !nova {
			r.Versao = n.Rascunho().Versao
		}
		if err := s.repo.SalvarRedacao(ctx, tx, n.ID, r, nova, ator(identity)); err != nil {
			return err
		}
		return s.registrarNorma(ctx, tx, EventNormaRascunho, n, nil, map[string]any{"versao": r.Versao, "motivo": r.Motivo})
	})
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	return s.Norma(ctx, identity, id)
}

// Publicacao é o registro da publicação da redação.
type Publicacao struct {
	Em     time.Time
	Diario string
}

// PublicarNorma publica o rascunho: ele passa a ser a redação em vigor e a
// anterior fica no histórico, vigente até esta data.
func (s *Service) PublicarNorma(ctx context.Context, identity auth.Identity, id uuid.UUID, pub Publicacao) (domain.Norma, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		n, err := s.carregarNorma(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		if err := domain.PodePublicarNorma(n, pub.Em, s.hoje(), pub.Diario); err != nil {
			return err
		}
		r := *n.Rascunho()
		em := pub.Em
		r.PublicadaEm, r.Diario = &em, strings.TrimSpace(pub.Diario)
		if err := s.repo.PublicarRedacao(ctx, tx, n.ID, r, ator(identity)); err != nil {
			return err
		}
		antes := n.Situacao
		n.Situacao = domain.NormaVigente
		if err := s.repo.UpdateNorma(ctx, tx, n); err != nil {
			return err
		}
		return s.registrarNorma(ctx, tx, EventNormaPublicada, n, map[string]any{"situacao": antes, "versao": n.Versao},
			map[string]any{"situacao": n.Situacao, "versao": r.Versao, "publicada_em": em, "diario": r.Diario})
	})
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	return s.Norma(ctx, identity, id)
}

// DescartarRascunho apaga a redação em rascunho. Se a norma nunca foi
// publicada, ela some junto (nada dela chegou ao público).
func (s *Service) DescartarRascunho(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Norma, error) {
	var apagada bool
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		n, err := s.carregarNorma(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		if n.Rascunho() == nil {
			return domain.ErrSemRascunho
		}
		if err := s.repo.ApagarRascunho(ctx, tx, n.ID); err != nil {
			return err
		}
		if apagada = n.Situacao == domain.NormaRascunho; apagada {
			if err := s.repo.DeleteNorma(ctx, tx, n.ID); err != nil {
				return err
			}
		}
		return s.registrarNorma(ctx, tx, EventNormaDescarte, n, map[string]any{"versao": n.Rascunho().Versao}, map[string]any{"norma_apagada": apagada})
	})
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	if apagada {
		return domain.Norma{}, nil
	}
	return s.Norma(ctx, identity, id)
}

// RevogarNorma encerra a vigência com o ato e a data. O texto continua
// consultável; um rascunho pendente é descartado.
func (s *Service) RevogarNorma(ctx context.Context, identity auth.Identity, id uuid.UUID, em time.Time, ato string) (domain.Norma, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		n, err := s.carregarNorma(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		if err := domain.PodeRevogar(n, em, s.hoje(), ato); err != nil {
			return err
		}
		if n.TemRascunho {
			if err := s.repo.ApagarRascunho(ctx, tx, n.ID); err != nil {
				return err
			}
		}
		n.Situacao, n.RevogadaEm, n.RevogacaoAto = domain.NormaRevogada, &em, strings.TrimSpace(ato)
		if err := s.repo.UpdateNorma(ctx, tx, n); err != nil {
			return err
		}
		return s.registrarNorma(ctx, tx, EventNormaRevogada, n, map[string]any{"situacao": domain.NormaVigente},
			map[string]any{"situacao": n.Situacao, "revogada_em": em, "ato": n.RevogacaoAto})
	})
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	return s.Norma(ctx, identity, id)
}

// Norma devolve a norma com todas as redações, inclusive o rascunho (só o
// Arquivo Público).
func (s *Service) Norma(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Norma, error) {
	if err := central(identity); err != nil {
		return domain.Norma{}, MapError(err)
	}
	n, err := s.repo.GetNorma(ctx, s.pool, id, false)
	return n, MapError(err)
}

// Normas pagina as normas para a gestão (inclui as nunca publicadas).
func (s *Service) Normas(ctx context.Context, identity auth.Identity, q string, p pagination.Params) ([]domain.Norma, int64, error) {
	if err := central(identity); err != nil {
		return nil, 0, MapError(err)
	}
	items, total, err := s.repo.ListNormas(ctx, s.pool, buscaNorma(q), false, p)
	return items, total, MapError(err)
}

// NormaPublica: o texto em vigor e as redações anteriores, para qualquer
// pessoa. A norma nunca publicada não existe para o público.
func (s *Service) NormaPublica(ctx context.Context, id uuid.UUID) (domain.Norma, error) {
	n, err := s.repo.GetNorma(ctx, s.pool, id, false)
	if err == nil && n.Situacao == domain.NormaRascunho {
		err = domain.ErrNormaNaoExiste
	}
	if err != nil {
		return domain.Norma{}, MapError(err)
	}
	return n.Publicas(), nil
}

// NormasPublicas pagina as normas publicadas (vigentes e revogadas).
func (s *Service) NormasPublicas(ctx context.Context, q string, p pagination.Params) ([]domain.Norma, int64, error) {
	items, total, err := s.repo.ListNormas(ctx, s.pool, buscaNorma(q), true, p)
	return items, total, MapError(err)
}

// buscaNorma limita o termo de busca.
func buscaNorma(q string) string {
	q = strings.TrimSpace(q)
	if utf8.RuneCountInString(q) > 100 {
		q = string([]rune(q)[:100])
	}
	return q
}
