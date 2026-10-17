package application

import (
	"context"
	"log/slog"
	"strings"

	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/audit"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database"
	"github.com/yurythx/sismarq/internal/platform/modkit"
)

// Fase 6 (ADR 038): acesso, busca e difusão.

// Eventos da Fase 6 (outbox).
const (
	EventDigitalPublicado    = "acervo.digital.publicado"
	EventDigitalDespublicado = "acervo.digital.despublicado"
	EventExposicaoSalva      = "acervo.exposicao.salva"
	EventExposicaoPublicada  = "acervo.exposicao.publicada"
	EventExposicaoRecolhida  = "acervo.exposicao.despublicada"
	// indiceLote: unidades por gravação no índice.
	indiceLote = 200
	// BuscaMaxPorPagina limita a página da busca.
	BuscaMaxPorPagina = 50
	destaques         = 6
	recentes          = 8
)

// IndiceBusca é a porta do índice da busca (o Typesense, em
// infrastructure.Indice).
type IndiceBusca interface {
	Garantir(ctx context.Context) (bool, error)
	Gravar(ctx context.Context, regs []domain.RegistroIndice, geracao int64) error
	Limpar(ctx context.Context, geracao int64) error
	Buscar(ctx context.Context, c domain.Consulta, f domain.FiltroBusca) (domain.ResultadoBusca, error)
}

// ConfigBusca liga o índice (nil: a busca fica no banco) e o canal do SIC.
type ConfigBusca struct {
	Indice IndiceBusca
	SICURL string
}

// WithBusca liga a busca e a consulta pública.
func (s *Service) WithBusca(c ConfigBusca) *Service {
	if !strings.HasPrefix(c.SICURL, "https://") && !strings.HasPrefix(c.SICURL, "http://") {
		c.SICURL = "" // o link vai para a página: só http(s)
	}
	s.busca = c
	return s
}

// PublicarDigital libera (publicar) ou retira o representante da consulta
// pública. Decide quem gere a unidade (recolhida, o Arquivo Público); só o
// ostensivo vai.
func (s *Service) PublicarDigital(ctx context.Context, identity auth.Identity, id uuid.UUID, publicar bool) (DetalheDigital, error) {
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		d, err := s.repo.GetDigital(ctx, tx, id)
		if err != nil {
			return err
		}
		u, _, err := s.repo.GetUnidade(ctx, tx, d.UnidadeID, true)
		if err != nil {
			return err
		}
		switch {
		case !le(identity, u.UnidadeID) && !auth.CanGlobal(identity, auth.PermAcervoCentral):
			return domain.ErrDigitalNaoExiste
		case !gerePara(identity, u):
			return domain.ErrForaDoEscopo
		case publicar == (d.PublicadoEm != nil):
			return domain.ErrPublicacao
		}
		var em *time.Time
		evento := EventDigitalDespublicado
		if publicar {
			if err := domain.PodePublicarDigital(u, d); err != nil {
				return err
			}
			agora := s.agora().UTC()
			em, evento = &agora, EventDigitalPublicado
		}
		if err := s.repo.PublicarDigital(ctx, tx, id, em, ator(identity)); err != nil {
			return err
		}
		payload := map[string]string{"id": id.String(), "numero": d.Numero, "unidade_arquivamento_id": u.ID.String()}
		if err := s.outbox.Write(ctx, tx, evento, "acervo_digital", id.String(), uuid.Nil, payload); err != nil {
			return err
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_digital", id.String(),
			map[string]any{"publicado_em": d.PublicadoEm}, map[string]any{"publicado_em": em}))
	})
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	return s.Digital(ctx, identity, id)
}

// descritores preenche o descritor da série (TTDD do Atlas) de cada
// registro; cache evita reler a mesma série.
func (s *Service) descritor(ctx context.Context, codigo string, cache map[string]string) string {
	d, ok := cache[codigo]
	if !ok {
		if serie, err := s.ttdd.Serie(ctx, codigo); err == nil {
			d = serie.Descritor
		}
		cache[codigo] = d
	}
	return d
}

func (s *Service) descritores(ctx context.Context, regs []domain.RegistroIndice, cache map[string]string) {
	for i := range regs {
		regs[i].Descritor = s.descritor(ctx, regs[i].Unidade.CodigoClassificacao, cache)
	}
}

// IndexarPendentes reindexa as unidades marcadas pelo banco (toda mudança
// numa unidade, num representante ou no texto do OCR). Devolve quantas.
func (s *Service) IndexarPendentes(ctx context.Context) (int, error) {
	total, cache := 0, map[string]string{}
	for s.busca.Indice != nil {
		ps, err := s.repo.PendentesIndice(ctx, s.pool, indiceLote)
		if err != nil || len(ps) == 0 {
			return total, err
		}
		ids := make([]uuid.UUID, len(ps))
		for i, p := range ps {
			ids[i] = p.UnidadeID
		}
		regs, err := s.repo.RegistrosIndice(ctx, s.pool, ids, true)
		if err != nil {
			return total, err
		}
		s.descritores(ctx, regs, cache)
		if err := s.busca.Indice.Gravar(ctx, regs, time.Now().UnixNano()); err != nil {
			return total, err
		}
		if err := s.repo.ConcluirPendentes(ctx, s.pool, ps); err != nil {
			return total, err
		}
		total += len(ps)
		if len(ps) < indiceLote {
			break
		}
	}
	return total, nil
}

// Recarregar faz a carga completa: grava o acervo inteiro numa geração nova
// e apaga o que sobrou das anteriores.
func (s *Service) Recarregar(ctx context.Context) (int, error) {
	if s.busca.Indice == nil {
		return 0, nil
	}
	geracao, depois, n, cache := time.Now().UnixNano(), uuid.Nil, 0, map[string]string{}
	for {
		regs, err := s.repo.PaginaIndice(ctx, s.pool, depois, indiceLote)
		if err != nil {
			return n, err
		}
		if len(regs) > 0 {
			s.descritores(ctx, regs, cache)
			if err := s.busca.Indice.Gravar(ctx, regs, geracao); err != nil {
				return n, err
			}
			n += len(regs)
			depois = regs[len(regs)-1].Unidade.ID
		}
		if len(regs) < indiceLote {
			break
		}
	}
	return n, s.busca.Indice.Limpar(ctx, geracao)
}

// RunIndice mantém o índice do acervo: garante as coleções (carga completa
// quando nascem), reindexa o que o banco marcou a cada intervalo e refaz a
// carga completa a cada recarga. Sem índice configurado, só espera.
func (s *Service) RunIndice(intervalo, recarga time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if s.busca.Indice == nil || intervalo <= 0 {
			<-ctx.Done()
			return nil
		}
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()
		pronto, ultima := false, time.Time{}
		for {
			if !pronto {
				criou, err := s.busca.Indice.Garantir(ctx)
				if err != nil {
					s.logger.ErrorContext(ctx, "acervo: índice da busca indisponível", slog.Any("error", err))
				} else {
					pronto = true
					if !criou {
						ultima = time.Now()
					}
				}

			}
			if pronto && time.Since(ultima) >= recarga {
				if n, err := s.Recarregar(ctx); err != nil {
					s.logger.ErrorContext(ctx, "acervo: carga completa do índice falhou", slog.Any("error", err))
				} else {
					ultima = time.Now()
					s.logger.InfoContext(ctx, "acervo: índice recarregado", slog.Int("unidades", n))
				}
			}
			if pronto {
				if _, err := s.IndexarPendentes(ctx); err != nil {
					s.logger.ErrorContext(ctx, "acervo: indexação falhou", slog.Any("error", err))
				}
			}
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}
		}
	}
}

// Buscar é a busca interna: o escopo de quem pede (o Arquivo Público vê
// todo o acervo). O restrito e o sigiloso só aparecem para quem tem a
// unidade no escopo.
func (s *Service) Buscar(ctx context.Context, identity auth.Identity, c domain.Consulta) (domain.ResultadoBusca, error) {
	f := domain.FiltroBusca{}
	if !auth.CanGlobal(identity, auth.PermAcervoCentral) {
		fl := s.filtro(identity, domain.Filtro{})
		f.Restrito, f.Unidades = fl.Restrito, fl.Unidades
	}
	return s.buscar(ctx, c, f)
}

// BuscarPublico é a consulta sem login: só o ostensivo.
func (s *Service) BuscarPublico(ctx context.Context, c domain.Consulta) (domain.ResultadoBusca, error) {
	r, err := s.buscar(ctx, c, domain.FiltroBusca{Publico: true})
	if err == nil {
		s.contar(ctx, domain.ConsultaBusca)
	}
	return r, err
}

func (s *Service) buscar(ctx context.Context, c domain.Consulta, f domain.FiltroBusca) (domain.ResultadoBusca, error) {
	c, err := domain.NormalizarConsulta(c, f.Publico, BuscaMaxPorPagina)
	if err != nil {
		return domain.ResultadoBusca{}, MapError(err)
	}
	if f.Restrito && len(f.Unidades) == 0 {
		return domain.ResultadoBusca{Itens: []domain.ItemBusca{}, Facetas: map[string][]domain.Faceta{}, Origem: domain.OrigemBanco}, nil
	}
	if s.busca.Indice != nil {
		r, err := s.busca.Indice.Buscar(ctx, c, f)
		if err == nil {
			return r, nil
		}
		s.logger.WarnContext(ctx, "acervo: índice da busca fora; a consulta vai ao banco", slog.Any("error", err))
	}
	inicio := time.Now()
	itens, total, err := s.repo.BuscarBanco(ctx, s.pool, c, f)
	if err != nil {
		return domain.ResultadoBusca{}, MapError(err)
	}
	cache := map[string]string{}
	for i := range itens {
		itens[i].Descritor = s.descritor(ctx, itens[i].Codigo, cache)
		if f.Publico {
			itens[i].NivelAcesso = ""
		}
	}
	return domain.ResultadoBusca{Itens: itens, Total: total, Facetas: map[string][]domain.Faceta{}, Origem: domain.OrigemBanco,
		Ms: time.Since(inicio).Milliseconds()}, nil
}

// contar soma a consulta pública do dia (indicador da SEMAD); a falha não
// derruba a consulta.
func (s *Service) contar(ctx context.Context, tipo string) {
	if err := s.repo.RegistrarConsulta(ctx, s.pool, s.hoje(), tipo); err != nil {
		s.logger.WarnContext(ctx, "acervo: consulta pública não contada", slog.String("tipo", tipo), slog.Any("error", err))
	}
}

// ContarDocumento soma a abertura de um documento na consulta pública.
func (s *Service) ContarDocumento(ctx context.Context) { s.contar(ctx, domain.ConsultaDocumento) }

// FichaPublica é a descrição ostensiva da unidade, com os representantes
// liberados. O que não é ostensivo responde como inexistente.
func (s *Service) FichaPublica(ctx context.Context, id uuid.UUID) (domain.FichaPublica, error) {
	regs, err := s.repo.RegistrosIndice(ctx, s.pool, []uuid.UUID{id}, false)
	if err == nil && (len(regs) == 0 || !regs[0].Ostensivo()) {
		err = domain.ErrNaoEncontrada
	}
	if err != nil {
		return domain.FichaPublica{}, MapError(err)
	}
	regs[0].Descritor = s.descritor(ctx, regs[0].Unidade.CodigoClassificacao, map[string]string{})
	s.contar(ctx, domain.ConsultaFicha)
	return domain.NovaFicha(regs[0]), nil
}

// DocumentoPublico confere que o representante está liberado (e a unidade
// segue ostensiva) para o visualizador da consulta pública.
func (s *Service) DocumentoPublico(ctx context.Context, id uuid.UUID) (domain.Digital, error) {
	d, err := s.repo.GetDigital(ctx, s.pool, id)
	if err == nil && (d.PublicadoEm == nil || d.Situacao != domain.DigitalPronto) {
		err = domain.ErrDigitalNaoExiste
	}
	if err == nil {
		var u domain.Unidade
		if u, _, err = s.repo.GetUnidade(ctx, s.pool, d.UnidadeID, false); err == nil && (u.NivelAcesso != domain.AcessoPublico || u.Fase == domain.FaseEliminado) {
			err = domain.ErrDigitalNaoExiste
		}
	}
	return d, MapError(err)
}

// PainelPublico é a página do Arquivo Público: números, exposições em
// destaque, os últimos documentos liberados e o canal do SIC.
type PainelPublico struct {
	SICURL     string                  `json:"sic_url"`
	SerieSIC   string                  `json:"serie_sic"`
	Unidades   int64                   `json:"unidades"`
	Digitais   int64                   `json:"digitais"`
	Exposicoes []domain.Exposicao      `json:"exposicoes"`
	Recentes   []domain.DigitalRecente `json:"recentes"`
}

// ArquivoPublico monta a página do Arquivo Público.
func (s *Service) ArquivoPublico(ctx context.Context) (PainelPublico, error) {
	p := PainelPublico{SICURL: s.busca.SICURL, SerieSIC: domain.SerieSIC}
	var err error
	if p.Unidades, p.Digitais, err = s.repo.ContarPublico(ctx, s.pool); err != nil {
		return p, MapError(err)
	}
	if p.Exposicoes, _, err = s.repo.ListExposicoes(ctx, s.pool, true, pagination.New(1, destaques, destaques)); err != nil {
		return p, MapError(err)
	}
	if p.Recentes, err = s.repo.DigitaisRecentes(ctx, s.pool, recentes); err != nil {
		return p, MapError(err)
	}
	return p, nil
}

// NovaExposicao é o que se grava da exposição.
type NovaExposicao struct {
	Titulo string
	Resumo string
	Texto  string
	Itens  []NovoItem
}

// NovoItem é um documento da exposição.
type NovoItem struct {
	DigitalID uuid.UUID
	Legenda   string
}

// Exposicoes pagina as exposições (Arquivo Público).
func (s *Service) Exposicoes(ctx context.Context, identity auth.Identity, p pagination.Params) ([]domain.Exposicao, int64, error) {
	if err := central(identity); err != nil {
		return nil, 0, MapError(err)
	}
	items, total, err := s.repo.ListExposicoes(ctx, s.pool, false, p)
	return items, total, MapError(err)
}

// Exposicao devolve a exposição com todos os itens (Arquivo Público).
func (s *Service) Exposicao(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Exposicao, error) {
	if err := central(identity); err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	x, err := s.repo.GetExposicao(ctx, s.pool, id, "")
	return x, MapError(err)
}

// SalvarExposicao cria (id nil) ou altera a exposição. O endereço sai do
// título na criação e não muda depois (os links já divulgados seguem
// valendo).
func (s *Service) SalvarExposicao(ctx context.Context, identity auth.Identity, id *uuid.UUID, n NovaExposicao) (domain.Exposicao, error) {
	if err := central(identity); err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	x := domain.Exposicao{Titulo: n.Titulo, Resumo: n.Resumo, Texto: n.Texto, Itens: make([]domain.ItemExposicao, len(n.Itens))}
	for i, it := range n.Itens {
		x.Itens[i] = domain.ItemExposicao{DigitalID: it.DigitalID, Legenda: it.Legenda}
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		agora := s.agora().UTC()
		var antes any
		if id == nil {
			x.ID, x.Slug, x.CriadoPor, x.CriadoEm = uuid.New(), modkit.Slugify(n.Titulo), ator(identity), agora
		} else {
			atual, err := s.repo.GetExposicao(ctx, tx, *id, "")
			if err != nil {
				return err
			}
			x.ID, x.Slug, x.CriadoPor, x.CriadoEm, antes = atual.ID, atual.Slug, atual.CriadoPor, atual.CriadoEm, atual
		}
		x.AtualizadoEm = agora
		var err error
		if x, err = domain.NovaExposicao(x); err != nil {
			return err
		}
		for _, it := range x.Itens {
			if _, err := s.repo.GetDigital(ctx, tx, it.DigitalID); err != nil {
				return err
			}
		}
		if id == nil {
			err = s.repo.InsertExposicao(ctx, tx, x)
		} else {
			err = s.repo.UpdateExposicao(ctx, tx, x)
		}
		if err != nil {
			return err
		}
		return s.registrarExposicao(ctx, tx, EventExposicaoSalva, x, antes)
	})
	if err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	return s.Exposicao(ctx, identity, x.ID)
}

func (s *Service) registrarExposicao(ctx context.Context, tx pgx.Tx, evento string, x domain.Exposicao, antes any) error {
	if err := s.outbox.Write(ctx, tx, evento, "acervo_exposicao", x.ID.String(), uuid.Nil,
		map[string]string{"id": x.ID.String(), "slug": x.Slug, "titulo": x.Titulo}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_exposicao", x.ID.String(), antes, x))
}

// PublicarExposicao publica (com ao menos um documento liberado) ou recolhe
// a exposição.
func (s *Service) PublicarExposicao(ctx context.Context, identity auth.Identity, id uuid.UUID, publicar bool) (domain.Exposicao, error) {
	if err := central(identity); err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		x, err := s.repo.GetExposicao(ctx, tx, id, "")
		if err != nil {
			return err
		}
		antes, evento := x, EventExposicaoRecolhida
		var em *time.Time
		switch {
		case publicar:
			if err := domain.PodePublicarExposicao(x); err != nil {
				return err
			}
			agora := s.agora().UTC()
			em, evento = &agora, EventExposicaoPublicada
		case x.PublicadaEm == nil:
			return domain.ErrExposicaoSituacao
		}
		if err := s.repo.PublicarExposicao(ctx, tx, id, em); err != nil {
			return err
		}
		x.PublicadaEm = em
		return s.registrarExposicao(ctx, tx, evento, x, antes)
	})
	if err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	return s.Exposicao(ctx, identity, id)
}

// ExposicoesPublicas pagina as exposições publicadas (sem os itens).
func (s *Service) ExposicoesPublicas(ctx context.Context, p pagination.Params) ([]domain.Exposicao, int64, error) {
	items, total, err := s.repo.ListExposicoes(ctx, s.pool, true, p)
	return items, total, MapError(err)
}

// ExposicaoPublica é a visita: só a publicada, só os documentos liberados.
func (s *Service) ExposicaoPublica(ctx context.Context, slug string) (domain.Exposicao, error) {
	x, err := s.repo.GetExposicao(ctx, s.pool, uuid.Nil, slug)
	if err == nil && x.PublicadaEm == nil {
		err = domain.ErrExposicaoNaoExiste
	}
	if err != nil {
		return domain.Exposicao{}, MapError(err)
	}
	s.contar(ctx, domain.ConsultaExposicao)
	return x.SoPublicados(), nil
}
