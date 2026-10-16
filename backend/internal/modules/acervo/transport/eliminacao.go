package transport

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	apperrors "github.com/yurythx/sismarq/internal/domain/errors"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// registerFase3 monta as rotas do rito da eliminação (ADR 034).
func (h *Handlers) registerFase3(r chi.Router, ler, central func(http.Handler) http.Handler) {
	avaliar := auth.RequirePermission(h.logger, auth.PermAcervoAvaliacao)
	r.With(ler).Get("/comissoes", h.Comissoes)
	r.With(central).Post("/comissoes", h.salvarComissao(false))
	r.With(central).Put("/comissoes/{id}", h.salvarComissao(true))
	r.With(ler).Get("/listagens", h.Listagens)
	r.With(ler).Get("/listagens/{id}", h.Listagem)
	r.With(avaliar).Post("/listagens", h.CriarListagem)
	r.With(avaliar).Put("/listagens/{id}", h.AlterarListagem)
	r.With(avaliar).Post("/listagens/{id}/aprovar", h.Aprovar)
	r.With(central).Post("/listagens/{id}/autorizar", h.listagemSimples(h.svc.Autorizar))
	r.With(central).Post("/listagens/{id}/devolver", h.DevolverListagem)
	r.With(avaliar).Post("/listagens/{id}/edital", h.PublicarEdital)
	r.With(avaliar).Post("/listagens/{id}/pedidos", h.RegistrarPedido)
	r.With(avaliar).Post("/listagens/{id}/pedidos/{pedido}/decidir", h.DecidirPedido)
	r.With(avaliar).Post("/listagens/{id}/amostragem", h.Amostrar)
	r.With(avaliar).Post("/listagens/{id}/termo", h.RegistrarTermo)
	r.With(avaliar).Post("/listagens/{id}/cancelar", h.listagemSimples(h.svc.CancelarListagem))
}

// Comissoes lista as comissões.
func (h *Handlers) Comissoes(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Comissoes(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, c)
}

// ComissaoRequest é o cadastro da comissão.
type ComissaoRequest struct {
	Tipo       string          `json:"tipo" validate:"required,max=16"`
	EntidadeID *uuid.UUID      `json:"entidade_id"`
	Nome       string          `json:"nome" validate:"required,max=200"`
	Portaria   string          `json:"portaria" validate:"max=200"`
	Ativa      *bool           `json:"ativa"`
	Membros    []domain.Membro `json:"membros" validate:"max=30"`
}

func (h *Handlers) salvarComissao(alterar bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id *uuid.UUID
		if alterar {
			v, err := httputil.UUIDParam(r, "id")
			if err != nil {
				h.fail(w, r, err)
				return
			}
			id = &v
		}
		var req ComissaoRequest
		if err := httputil.Bind(w, r, &req); err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		c, err := h.svc.SalvarComissao(r.Context(), identity, id, domain.Comissao{Tipo: strings.ToUpper(req.Tipo), EntidadeID: req.EntidadeID,
			Nome: req.Nome, Portaria: req.Portaria, Ativa: req.Ativa == nil || *req.Ativa, Membros: req.Membros})
		if err != nil {
			h.fail(w, r, err)
			return
		}
		if alterar {
			httputil.WriteOK(w, c)
			return
		}
		httputil.WriteCreated(w, c)
	}
}

// Listagens pagina as relações de eliminação.
func (h *Handlers) Listagens(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Listagens(r.Context(), identity, p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// Listagem devolve a relação.
func (h *Handlers) Listagem(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.Listagem(r.Context(), identity, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// ListagemRequest é o cabeçalho da relação (Anexo I) e as unidades.
type ListagemRequest struct {
	UnidadeID          uuid.UUID   `json:"unidade_id" validate:"required"`
	ResponsavelSelecao *uuid.UUID  `json:"responsavel_selecao"`
	ParecerJuridico    string      `json:"parecer_juridico" validate:"max=4000"`
	ContasExercicios   string      `json:"contas_exercicios" validate:"max=200"`
	ContasAprovadasEm  string      `json:"contas_aprovadas_em"`
	Unidades           []uuid.UUID `json:"unidades" validate:"max=2000"`
	Retirar            []uuid.UUID `json:"retirar" validate:"max=2000"`
}

func (req ListagemRequest) nova() (application.NovaListagem, error) {
	contas, err := dataOpcional("contas_aprovadas_em", req.ContasAprovadasEm)
	return application.NovaListagem{UnidadeID: req.UnidadeID, ResponsavelSelecao: req.ResponsavelSelecao, ParecerJuridico: req.ParecerJuridico,
		ContasExercicios: req.ContasExercicios, ContasAprovadasEm: contas, Unidades: req.Unidades}, err
}

// CriarListagem abre a relação.
func (h *Handlers) CriarListagem(w http.ResponseWriter, r *http.Request) {
	var req ListagemRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	n, err := req.nova()
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.CriarListagem(r.Context(), identity, n)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, l)
}

// AlterarListagem muda o cabeçalho e as unidades (rascunho).
func (h *Handlers) AlterarListagem(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req ListagemRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	n, err := req.nova()
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.AlterarListagem(r.Context(), identity, id, n, req.Retirar)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// AtaRequest é a ata da CSPAD (Anexo IV). ata_em em AAAA-MM-DDTHH:MM.
type AtaRequest struct {
	Numero        string      `json:"ata_numero" validate:"required,max=40"`
	Em            string      `json:"ata_em" validate:"required"`
	Local         string      `json:"ata_local" validate:"max=200"`
	Manifestacoes string      `json:"ata_manifestacoes" validate:"max=8000"`
	Presentes     []uuid.UUID `json:"presentes" validate:"max=30"`
}

func dataHora(campo, v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, apperrors.Validation(campo + ": use AAAA-MM-DDTHH:MM")
}

// Aprovar registra a ata.
func (h *Handlers) Aprovar(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req AtaRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	em, err := dataHora("ata_em", req.Em)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.Aprovar(r.Context(), identity, id, application.Ata{Numero: req.Numero, Em: em, Local: req.Local,
		Manifestacoes: req.Manifestacoes, Presentes: req.Presentes})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// listagemSimples é uma etapa sem corpo (autorizar, cancelar).
func (h *Handlers) listagemSimples(fn func(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Listagem, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := httputil.UUIDParam(r, "id")
		if err != nil {
			h.fail(w, r, err)
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		l, err := fn(r.Context(), identity, id)
		if err != nil {
			h.fail(w, r, err)
			return
		}
		httputil.WriteOK(w, l)
	}
}

// MotivoRequest traz o motivo da devolução.
type MotivoRequest struct {
	Motivo string `json:"motivo" validate:"required,max=2000"`
}

// DevolverListagem devolve a relação à CSPAD.
func (h *Handlers) DevolverListagem(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req MotivoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.DevolverListagem(r.Context(), identity, id, req.Motivo)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// EditalRequest é a publicação do Edital de Ciência (Anexo III).
type EditalRequest struct {
	Numero      string `json:"edital_numero" validate:"required,max=40"`
	PublicadoEm string `json:"edital_publicado_em" validate:"required"`
	Diario      string `json:"edital_diario" validate:"max=200"`
}

// PublicarEdital registra a publicação.
func (h *Handlers) PublicarEdital(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req EditalRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	pub, err := data("edital_publicado_em", req.PublicadoEm)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.PublicarEdital(r.Context(), identity, id, application.Edital{Numero: req.Numero, PublicadoEm: pub, Diario: req.Diario})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// PedidoRequest é o pedido de desentranhamento ou de cópia.
type PedidoRequest struct {
	UnidadeID  *uuid.UUID `json:"unidade_arquivamento_id"`
	Requerente string     `json:"requerente" validate:"required,max=200"`
	Pedido     string     `json:"pedido" validate:"required,max=4000"`
	RecebidoEm string     `json:"recebido_em" validate:"required"`
}

// RegistrarPedido recebe o pedido.
func (h *Handlers) RegistrarPedido(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req PedidoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	em, err := data("recebido_em", req.RecebidoEm)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.RegistrarPedido(r.Context(), identity, id, domain.Pedido{UnidadeID: req.UnidadeID, Requerente: req.Requerente,
		Pedido: req.Pedido, RecebidoEm: em})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, l)
}

// DecisaoRequest é a decisão sobre o pedido.
type DecisaoRequest struct {
	Deferido bool   `json:"deferido"`
	Retirar  bool   `json:"retirar"`
	Decisao  string `json:"decisao" validate:"required,max=4000"`
}

// DecidirPedido defere ou indefere.
func (h *Handlers) DecidirPedido(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	pedido, err := httputil.UUIDParam(r, "pedido")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req DecisaoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.DecidirPedido(r.Context(), identity, id, pedido, req.Deferido, req.Retirar, req.Decisao)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// AmostragemRequest traz os critérios e o percentual (art. 39).
type AmostragemRequest struct {
	Criterios  string  `json:"criterios" validate:"required,max=2000"`
	Percentual float64 `json:"percentual" validate:"required"`
}

// Amostrar faz a amostragem.
func (h *Handlers) Amostrar(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req AmostragemRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.Amostrar(r.Context(), identity, id, req.Criterios, req.Percentual)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}

// TermoRequest é o Termo de Eliminação (Anexo II).
type TermoRequest struct {
	EliminadoEm     string `json:"termo_eliminado_em" validate:"required"`
	Acompanhamento  string `json:"termo_acompanhamento" validate:"required,max=200"`
	DestinacaoPapel string `json:"termo_destinacao_papel" validate:"max=500"`
}

// RegistrarTermo registra a eliminação.
func (h *Handlers) RegistrarTermo(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req TermoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	em, err := data("termo_eliminado_em", req.EliminadoEm)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	l, err := h.svc.RegistrarTermo(r.Context(), identity, id, application.Termo{EliminadoEm: em, Acompanhamento: req.Acompanhamento,
		DestinacaoPapel: req.DestinacaoPapel})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, l)
}
