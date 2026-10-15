package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// ListLocais devolve a hierarquia física (prédio > depósito > … >
// prateleira) com o caminho de cada local.
func (s *Service) ListLocais(ctx context.Context) ([]domain.Local, error) {
	l, err := s.repo.ListLocais(ctx, s.pool)
	return l, MapError(err)
}

// SalvarLocal cria (id nil) ou altera um local. Só o Arquivo Público
// (acervo:central, concessão global) mantém a hierarquia física.
func (s *Service) SalvarLocal(ctx context.Context, identity auth.Identity, id *uuid.UUID, l domain.Local) (domain.Local, error) {
	if !auth.CanGlobal(identity, auth.PermAcervoCentral) {
		return domain.Local{}, MapError(domain.ErrForaDoEscopo)
	}
	if id != nil {
		l.ID = *id
	} else {
		l.ID = uuid.New()
	}
	if err := l.Validar(); err != nil {
		return domain.Local{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		var antes *domain.Local
		if id != nil {
			atual, err := s.repo.GetLocal(ctx, tx, l.ID)
			if err != nil {
				return err
			}
			antes = &atual
		}
		if err := s.conferirPai(ctx, tx, l); err != nil {
			return err
		}
		if antes == nil {
			if err := s.repo.InsertLocal(ctx, tx, l); err != nil {
				return err
			}
		} else if err := s.repo.UpdateLocal(ctx, tx, l); err != nil {
			return err
		}
		if err := s.outbox.Write(ctx, tx, EventLocalSalvo, "acervo_local", l.ID.String(), uuid.Nil,
			map[string]string{"id": l.ID.String(), "nome": l.Nome}); err != nil {
			return err
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventLocalSalvo, "acervo_local", l.ID.String(), antes, l))
	})
	if err != nil {
		return domain.Local{}, MapError(err)
	}
	out, err := s.repo.GetLocal(ctx, s.pool, l.ID)
	return out, MapError(err)
}

// conferirPai exige um pai existente e impede ciclos (um local não pode ficar
// dentro de um descendente seu).
func (s *Service) conferirPai(ctx context.Context, db database.DBTX, l domain.Local) error {
	for pai := l.ParentID; pai != nil; {
		if *pai == l.ID {
			return &domain.ErrValidacao{Campo: "parent_id", Motivo: "um local não pode ficar dentro de um local que ele contém"}
		}
		p, err := s.repo.GetLocal(ctx, db, *pai)
		if err != nil {
			return err
		}
		pai = p.ParentID
	}
	return nil
}
