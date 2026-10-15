package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
)

// Eventos da Fase 2 (ADR 033).
const (
	EventEmprestada    = "acervo.unidade.emprestada"
	EventDevolvida     = "acervo.unidade.devolvida"
	EventGuiaEnviada   = "acervo.guia.enviada"
	EventGuiaRecebida  = "acervo.guia.recebida"
	EventGuiaCancelada = "acervo.guia.cancelada"
	EventImportacao    = "acervo.importacao.aplicada"
)

// Emprestar registra a saída da unidade para quem a pediu (art. 12, II).
// Eliminada não sai; a que já está fora, também não.
func (s *Service) Emprestar(ctx context.Context, identity auth.Identity, id, solicitante uuid.UUID, finalidade string, devolverAte time.Time) (Detalhe, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		u, err := s.carregar(ctx, tx, identity, id, auth.PermAcervoManage)
		if err != nil {
			return err
		}
		if err := u.Editavel(); err != nil {
			return err
		}
		e, err := domain.NovoEmprestimo(u.ID, solicitante, finalidade, devolverAte, s.hoje())
		if err != nil {
			return err
		}
		e.RetiradoPor, e.RetiradoEm = ator(identity), s.agora().UTC()
		if err := s.repo.InsertEmprestimo(ctx, tx, e); err != nil {
			return err
		}
		return s.registrar(ctx, tx, identity, EventEmprestada, u,
			domain.Movimentacao{Tipo: domain.MovEmprestimo, Descricao: "Emprestada até " + e.DevolverAte.Format("02/01/2006") + ": " + e.Finalidade}, nil, e)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, id)
}

// Devolver registra a volta da unidade.
func (s *Service) Devolver(ctx context.Context, identity auth.Identity, emprestimoID uuid.UUID, observacao string) (Detalhe, error) {
	observacao = strings.TrimSpace(observacao)
	if len([]rune(observacao)) > 500 {
		return Detalhe{}, MapError(&domain.ErrValidacao{Campo: "observacao", Motivo: "até 500 caracteres"})
	}
	var unidade uuid.UUID
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		e, err := s.repo.GetEmprestimo(ctx, tx, emprestimoID)
		if err != nil {
			return err
		}
		u, err := s.carregar(ctx, tx, identity, e.UnidadeID, auth.PermAcervoManage)
		if err != nil {
			return err
		}
		if !e.Aberto() {
			return domain.ErrEmprestimoDevolvido
		}
		antes := e
		agora := s.agora().UTC()
		e.DevolvidoEm, e.RecebidoPor, e.Observacao = &agora, ator(identity), observacao
		if err := s.repo.DevolverEmprestimo(ctx, tx, e); err != nil {
			return err
		}
		unidade = u.ID
		desc := "Devolvida"
		if observacao != "" {
			desc += ": " + observacao
		}
		return s.registrar(ctx, tx, identity, EventDevolvida, u, domain.Movimentacao{Tipo: domain.MovDevolucao, Descricao: desc}, antes, e)
	})
	if err != nil {
		return Detalhe{}, MapError(err)
	}
	return s.Get(ctx, identity, unidade)
}

// Emprestimos pagina os empréstimos em aberto no escopo (atrasados: só os
// que passaram do prazo).
func (s *Service) Emprestimos(ctx context.Context, identity auth.Identity, atrasados bool, p pagination.Params) ([]domain.Emprestimo, int64, error) {
	var em *time.Time
	if atrasados {
		h := s.hoje()
		em = &h
	}
	items, total, err := s.repo.ListEmprestimosAbertos(ctx, s.pool, s.filtro(identity, domain.Filtro{}), em, p)
	return items, total, MapError(err)
}

// conferirRemessa diz se a unidade pode ir numa guia: da origem, no
// corrente, com o prazo corrente cumprido (art. 29), sem estar emprestada nem
// em outra guia.
func (s *Service) conferirRemessa(ctx context.Context, db database.DBTX, origem uuid.UUID, u domain.Unidade, guiaAtual uuid.UUID) error {
	if u.UnidadeID != origem {
		return &domain.ErrValidacao{Campo: "unidades", Motivo: u.Numero + " não é da unidade de origem"}
	}
	_, p, err := s.calcular(ctx, db, u)
	if err != nil {
		return err
	}
	if err := domain.PodeTransferir(u, p, s.hoje()); err != nil {
		return err
	}
	fora, err := s.repo.EmprestimoAberto(ctx, db, u.ID)
	if err != nil {
		return err
	}
	if fora {
		return domain.ErrEmprestada
	}
	if guiaAtual == uuid.Nil {
		em, err := s.repo.EmGuiaEnviada(ctx, db, u.ID)
		if err != nil {
			return err
		}
		if em {
			return domain.ErrEmOutraGuia
		}
	}
	return nil
}

// CriarGuia monta a guia de transferência das unidades da origem ao arquivo
// intermediário (art. 29). As unidades só mudam de fase no recebimento.
func (s *Service) CriarGuia(ctx context.Context, identity auth.Identity, origem uuid.UUID, unidades []uuid.UUID, observacao string) (domain.Guia, error) {
	g, ids, err := domain.NovaGuia(origem, unidades, observacao)
	if err != nil {
		return domain.Guia{}, MapError(err)
	}
	if !gere(identity, origem) {
		return domain.Guia{}, MapError(domain.ErrForaDoEscopo)
	}
	g.CriadaPor, g.CriadaEm = ator(identity), s.agora().UTC()
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		return s.criarGuiaTx(ctx, tx, identity, &g, ids)
	})
	if err != nil {
		return domain.Guia{}, MapError(err)
	}
	return s.Guia(ctx, identity, g.ID)
}

func (s *Service) criarGuiaTx(ctx context.Context, tx pgx.Tx, identity auth.Identity, g *domain.Guia, ids []uuid.UUID) error {
	var lote []domain.Unidade
	for _, id := range ids {
		u, _, err := s.repo.GetUnidade(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if err := s.conferirRemessa(ctx, tx, g.OrigemID, u, uuid.Nil); err != nil {
			return err
		}
		lote = append(lote, u)
	}
	var err error
	if g.Numero, err = s.repo.InsertGuia(ctx, tx, *g, ids); err != nil {
		return err
	}
	for _, u := range lote {
		if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovGuia, Por: ator(identity),
			Descricao: "Relacionada na guia de transferência " + g.Numero}); err != nil {
			return err
		}
	}
	if err := s.outbox.Write(ctx, tx, EventGuiaEnviada, "acervo_guia", g.ID.String(), uuid.Nil,
		map[string]any{"id": g.ID, "numero": g.Numero, "origem_id": g.OrigemID, "itens": len(ids)}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventGuiaEnviada, "acervo_guia", g.ID.String(), nil,
		map[string]any{"numero": g.Numero, "origem_id": g.OrigemID, "unidades": ids, "observacao": g.Observacao}))
}

// carregarGuia lê a guia enviada para recebê-la ou cancelá-la, no escopo de
// gestão da origem.
func (s *Service) carregarGuia(ctx context.Context, tx pgx.Tx, identity auth.Identity, id uuid.UUID) (domain.Guia, error) {
	g, err := s.repo.GetGuia(ctx, tx, id, true)
	if err != nil {
		return g, err
	}
	if !le(identity, g.OrigemID) {
		return g, domain.ErrGuiaNaoExiste
	}
	if !gere(identity, g.OrigemID) {
		return g, domain.ErrForaDoEscopo
	}
	if g.Situacao != domain.GuiaEnviada {
		return g, domain.ErrGuiaSituacao
	}
	return g, nil
}

// ReceberGuia é o aceite do arquivo intermediário: cada unidade passa à fase
// intermediária (art. 29), conferida de novo.
func (s *Service) ReceberGuia(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Guia, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		g, err := s.carregarGuia(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		return s.receberGuiaTx(ctx, tx, identity, g)
	})
	if err != nil {
		return domain.Guia{}, MapError(err)
	}
	return s.Guia(ctx, identity, id)
}

func (s *Service) receberGuiaTx(ctx context.Context, tx pgx.Tx, identity auth.Identity, g domain.Guia) error {
	for _, item := range g.Itens {
		u, _, err := s.repo.GetUnidade(ctx, tx, item.ID, true)
		if err != nil {
			return err
		}
		if err := s.conferirRemessa(ctx, tx, g.OrigemID, u, g.ID); err != nil {
			return err
		}
		de := u.Fase
		u.Fase = domain.FaseIntermediaria
		_, p, err := s.calcular(ctx, tx, u)
		if err != nil {
			return err
		}
		if err := s.repo.SetFase(ctx, tx, u.ID, u.Fase, calculo(p)); err != nil {
			return err
		}
		para := u.Fase
		if err := s.registrar(ctx, tx, identity, EventTransferida, u, domain.Movimentacao{Tipo: domain.MovTransferencia, FaseDe: &de,
			FasePara: &para, Descricao: "Transferida ao arquivo intermediário pela guia " + g.Numero + " (art. 29)"},
			map[string]any{"fase": de}, map[string]any{"fase": para, "guia": g.Numero}); err != nil {
			return err
		}
	}
	antes := g.Situacao
	agora := s.agora().UTC()
	g.Situacao, g.RecebidaPor, g.RecebidaEm = domain.GuiaRecebida, ator(identity), &agora
	if err := s.repo.SetGuiaSituacao(ctx, tx, g); err != nil {
		return err
	}
	if err := s.outbox.Write(ctx, tx, EventGuiaRecebida, "acervo_guia", g.ID.String(), uuid.Nil,
		map[string]any{"id": g.ID, "numero": g.Numero, "itens": len(g.Itens)}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventGuiaRecebida, "acervo_guia", g.ID.String(),
		map[string]any{"situacao": antes}, map[string]any{"situacao": g.Situacao}))
}

// CancelarGuia desfaz a remessa ainda não recebida.
func (s *Service) CancelarGuia(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Guia, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		g, err := s.carregarGuia(ctx, tx, identity, id)
		if err != nil {
			return err
		}
		g.Situacao = domain.GuiaCancelada
		if err := s.repo.SetGuiaSituacao(ctx, tx, g); err != nil {
			return err
		}
		if err := s.outbox.Write(ctx, tx, EventGuiaCancelada, "acervo_guia", g.ID.String(), uuid.Nil,
			map[string]any{"id": g.ID, "numero": g.Numero}); err != nil {
			return err
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventGuiaCancelada, "acervo_guia", g.ID.String(),
			map[string]any{"situacao": domain.GuiaEnviada}, map[string]any{"situacao": g.Situacao}))
	})
	if err != nil {
		return domain.Guia{}, MapError(err)
	}
	return s.Guia(ctx, identity, id)
}

// Guia devolve a guia com as unidades. Fora do escopo de leitura, 404.
func (s *Service) Guia(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Guia, error) {
	g, err := s.repo.GetGuia(ctx, s.pool, id, false)
	if err == nil && !le(identity, g.OrigemID) {
		err = domain.ErrGuiaNaoExiste
	}
	return g, MapError(err)
}

// Guias pagina as guias no escopo.
func (s *Service) Guias(ctx context.Context, identity auth.Identity, p pagination.Params) ([]domain.Guia, int64, error) {
	items, total, err := s.repo.ListGuias(ctx, s.pool, s.filtro(identity, domain.Filtro{}), p)
	return items, total, MapError(err)
}

// ResultadoImportacao é o relatório da planilha: o hash do arquivo
// conferido, quantas linhas entrariam (ou entraram) e os erros por linha.
type ResultadoImportacao struct {
	Hash     string             `json:"hash"`
	Linhas   int                `json:"linhas"`
	Aplicada bool               `json:"aplicada"`
	Erros    []domain.ErroLinha `json:"erros"`
	Numeros  []string           `json:"numeros"`
}

var errSimulacao = errors.New("acervo: simulação")

// Importar lê a planilha do acervo existente e confere cada linha com as
// mesmas regras do cadastro. Simula (aplicar=false) ou grava — tudo ou nada:
// só a planilha conferida (mesmo hash) e sem nenhum erro é aplicada.
func (s *Service) Importar(ctx context.Context, identity auth.Identity, conteudo, hashConferido string, aplicar bool) (ResultadoImportacao, error) {
	soma := sha256.Sum256([]byte(conteudo))
	out := ResultadoImportacao{Hash: hex.EncodeToString(soma[:]), Erros: []domain.ErroLinha{}, Numeros: []string{}}
	if aplicar && hashConferido != out.Hash {
		return out, MapError(&domain.ErrValidacao{Campo: "hash", Motivo: "o arquivo enviado não é o que foi simulado: simule de novo antes de aplicar"})
	}
	linhas, erros, err := domain.LerPlanilha(conteudo)
	if err != nil {
		return out, MapError(&domain.ErrValidacao{Campo: "planilha", Motivo: strings.TrimPrefix(err.Error(), "acervo: planilha inválida: ")})
	}
	out.Erros = append(out.Erros, erros...)
	locais, err := s.repo.ListLocais(ctx, s.pool)
	if err != nil {
		return out, MapError(err)
	}
	porCaminho := map[string]uuid.UUID{}
	for _, l := range locais {
		if l.Ativo {
			porCaminho[strings.ToLower(l.Caminho)] = l.ID
		}
	}
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		for _, l := range linhas {
			numero, err := s.importarLinha(ctx, tx, identity, l, porCaminho)
			var v *domain.ErrValidacao
			switch {
			case errors.As(err, &v), errors.Is(err, domain.ErrSerieInvalida), errors.Is(err, domain.ErrUnidadeProdutora), errors.Is(err, domain.ErrForaDoEscopo):
				out.Erros = append(out.Erros, domain.ErroLinha{Linha: l.Linha, Motivo: strings.TrimPrefix(err.Error(), "acervo: ")})
				continue
			case err != nil:
				return err
			}
			out.Numeros = append(out.Numeros, numero)
		}
		out.Linhas = len(out.Numeros)
		if !aplicar || len(out.Erros) > 0 {
			return errSimulacao // nada fica gravado
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventImportacao, "acervo_importacao", out.Hash, nil,
			map[string]any{"linhas": out.Linhas}))
	})
	if errors.Is(err, errSimulacao) {
		out.Numeros = []string{}
		return out, nil
	}
	if err != nil {
		return out, MapError(err)
	}
	out.Aplicada = true
	return out, nil
}

// importarLinha resolve as referências e cadastra a unidade na transação.
func (s *Service) importarLinha(ctx context.Context, tx pgx.Tx, identity auth.Identity, l domain.LinhaImportada, locais map[string]uuid.UUID) (string, error) {
	u := l.Unidade
	var err error
	if u.UnidadeID, err = s.repo.ResolverUnidade(ctx, tx, l.UnidadeRef); err != nil {
		return "", err
	}
	if l.LocalRef != "" {
		id, ok := locais[strings.ToLower(strings.TrimSpace(l.LocalRef))]
		if !ok {
			return "", &domain.ErrValidacao{Campo: "local", Motivo: "local " + l.LocalRef + " não encontrado ou desativado"}
		}
		u.LocalID = &id
	}
	criada, err := s.criarTx(ctx, tx, identity, u, domain.MovImportacao, "Importada da planilha do acervo existente")
	return criada.Numero, err
}
