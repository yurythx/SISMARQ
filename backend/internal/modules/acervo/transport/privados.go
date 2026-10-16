package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/pkg/httputil"
)

// domainID é o id do arquivo privado lido da rota.
type domainID = uuid.UUID

// registerPrivados monta as rotas dos arquivos privados (ADR 036), só do
// Arquivo Público.
func (h *Handlers) registerPrivados(r chi.Router, central func(http.Handler) http.Handler) {
	r.With(central).Get("/privados", h.Privados)
	r.With(central).Post("/privados", h.CadastrarPrivado)
	r.With(central).Get("/privados/{id}", h.Privado)
	r.With(central).Put("/privados/{id}", h.AlterarPrivado)
	r.With(central).Post("/privados/{id}/parecer", h.RegistrarParecer)
	r.With(central).Post("/privados/{id}/declarar", h.Declarar)
	r.With(central).Post("/privados/{id}/guarda", h.MudarGuarda)
	r.With(central).Post("/privados/{id}/alienacoes", h.NotificarAlienacao)
	r.With(central).Post("/privados/{id}/alienacoes/{alienacao}/decidir", h.DecidirAlienacao)
}

// PrivadoRequest é o cadastro do arquivo privado.
type PrivadoRequest struct {
	Titular          string `json:"titular" validate:"required,max=200"`
	TipoTitular      string `json:"tipo_titular" validate:"required,max=2"`
	Documento        string `json:"documento" validate:"max=20"`
	Contato          string `json:"contato" validate:"max=500"`
	Descricao        string `json:"descricao" validate:"required,max=8000"`
	DatasLimite      string `json:"datas_limite" validate:"max=100"`
	Mensuracao       string `json:"mensuracao" validate:"max=200"`
	LocalGuarda      string `json:"local_guarda" validate:"max=500"`
	ServicoPublico   bool   `json:"servico_publico"`
	AcessoAutorizado bool   `json:"acesso_autorizado"`
	AcessoCondicoes  string `json:"acesso_condicoes" validate:"max=2000"`
	Convenio         string `json:"convenio" validate:"max=1000"`
}

func (r PrivadoRequest) privado() domain.Privado {
	return domain.Privado{Titular: r.Titular, TipoTitular: r.TipoTitular, Documento: r.Documento, Contato: r.Contato,
		Descricao: r.Descricao, DatasLimite: r.DatasLimite, Mensuracao: r.Mensuracao, LocalGuarda: r.LocalGuarda,
		ServicoPublico: r.ServicoPublico, AcessoAutorizado: r.AcessoAutorizado, AcessoCondicoes: r.AcessoCondicoes, Convenio: r.Convenio}
}

// Privados pagina os arquivos privados (?q=).
func (h *Handlers) Privados(w http.ResponseWriter, r *http.Request) {
	identity, _ := auth.IdentityFromContext(r.Context())
	p := httputil.Page(r, h.maxPageSize)
	items, total, err := h.svc.Privados(r.Context(), identity, httputil.Query(r, "q", 100), p)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WritePage(w, items, p, total)
}

// CadastrarPrivado registra o arquivo privado.
func (h *Handlers) CadastrarPrivado(w http.ResponseWriter, r *http.Request) {
	var req PrivadoRequest
	if err := httputil.Bind(w, r, &req); err != nil {
		h.fail(w, r, err)
		return
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	p, err := h.svc.CadastrarPrivado(r.Context(), identity, req.privado())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteCreated(w, p)
}

// Privado devolve o arquivo com o histórico e as alienações.
func (h *Handlers) Privado(w http.ResponseWriter, r *http.Request) {
	h.privadoAcao(w, r, func(r *http.Request, id domainID) (domain.Privado, error) {
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.Privado(r.Context(), identity, id)
	})
}

// AlterarPrivado corrige o cadastro.
func (h *Handlers) AlterarPrivado(w http.ResponseWriter, r *http.Request) {
	var req PrivadoRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.AlterarPrivado(r.Context(), identity, id, req.privado())
	})
}

// ParecerRequest é o parecer da Comissão Especial.
type ParecerRequest struct {
	Comissao  string `json:"parecer_comissao" validate:"required,max=2000"`
	Texto     string `json:"parecer_texto" validate:"required,max=8000"`
	Favoravel bool   `json:"parecer_favoravel"`
	Em        string `json:"parecer_em" validate:"required"`
}

// RegistrarParecer grava o parecer.
func (h *Handlers) RegistrarParecer(w http.ResponseWriter, r *http.Request) {
	var req ParecerRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		em, err := data("parecer_em", req.Em)
		if err != nil {
			return domain.Privado{}, err
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.RegistrarParecer(r.Context(), identity, id, application.Parecer{Comissao: req.Comissao, Texto: req.Texto,
			Favoravel: req.Favoravel, Em: em})
	})
}

// DecretoRequest é o decreto de declaração.
type DecretoRequest struct {
	Numero string `json:"decreto_numero" validate:"required,max=40"`
	Em     string `json:"decreto_em" validate:"required"`
	Diario string `json:"decreto_diario" validate:"max=200"`
}

// Declarar registra o decreto.
func (h *Handlers) Declarar(w http.ResponseWriter, r *http.Request) {
	var req DecretoRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		em, err := data("decreto_em", req.Em)
		if err != nil {
			return domain.Privado{}, err
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.Declarar(r.Context(), identity, id, req.Numero, req.Diario, em)
	})
}

// GuardaRequest é a doação, o depósito ou a revogação do depósito.
type GuardaRequest struct {
	Guarda string `json:"guarda" validate:"required,max=20"`
	Termo  string `json:"termo" validate:"required,max=1000"`
	Em     string `json:"em" validate:"required"`
}

// MudarGuarda registra a mudança de guarda.
func (h *Handlers) MudarGuarda(w http.ResponseWriter, r *http.Request) {
	var req GuardaRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		em, err := data("em", req.Em)
		if err != nil {
			return domain.Privado{}, err
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.MudarGuarda(r.Context(), identity, id, req.Guarda, req.Termo, em)
	})
}

// AlienacaoRequest é a notificação de alienação (art. 58).
type AlienacaoRequest struct {
	NotificadoEm string `json:"notificado_em" validate:"required"`
	Adquirente   string `json:"adquirente" validate:"max=200"`
	Condicoes    string `json:"condicoes" validate:"max=2000"`
}

// NotificarAlienacao registra a notificação.
func (h *Handlers) NotificarAlienacao(w http.ResponseWriter, r *http.Request) {
	var req AlienacaoRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		em, err := data("notificado_em", req.NotificadoEm)
		if err != nil {
			return domain.Privado{}, err
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.NotificarAlienacao(r.Context(), identity, id, domain.Alienacao{NotificadoEm: em, Adquirente: req.Adquirente, Condicoes: req.Condicoes})
	})
}

// PreferenciaRequest é a manifestação do município.
type PreferenciaRequest struct {
	Decisao    string `json:"decisao" validate:"required,max=20"`
	Fundamento string `json:"fundamento" validate:"required,max=2000"`
}

// DecidirAlienacao registra a manifestação do município.
func (h *Handlers) DecidirAlienacao(w http.ResponseWriter, r *http.Request) {
	alienacao, err := httputil.UUIDParam(r, "alienacao")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var req PreferenciaRequest
	h.privadoComCorpo(w, r, &req, func(r *http.Request, id domainID) (domain.Privado, error) {
		identity, _ := auth.IdentityFromContext(r.Context())
		return h.svc.DecidirAlienacao(r.Context(), identity, id, alienacao, req.Decisao, req.Fundamento)
	})
}

// privadoAcao lê o {id} e responde o arquivo devolvido pela ação.
func (h *Handlers) privadoAcao(w http.ResponseWriter, r *http.Request, fn func(*http.Request, domainID) (domain.Privado, error)) {
	id, err := httputil.UUIDParam(r, "id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	p, err := fn(r, id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	httputil.WriteOK(w, p)
}

// privadoComCorpo lê o {id} e o corpo antes da ação.
func (h *Handlers) privadoComCorpo(w http.ResponseWriter, r *http.Request, req any, fn func(*http.Request, domainID) (domain.Privado, error)) {
	h.privadoAcao(w, r, func(r *http.Request, id domainID) (domain.Privado, error) {
		if err := httputil.Bind(w, r, req); err != nil {
			return domain.Privado{}, err
		}
		return fn(r, id)
	})
}
