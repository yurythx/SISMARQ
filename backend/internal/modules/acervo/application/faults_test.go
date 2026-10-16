package application_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yurythx/sismarq/internal/domain/pagination"
	"github.com/yurythx/sismarq/internal/modules/acervo/application"
	"github.com/yurythx/sismarq/internal/modules/acervo/domain"
	"github.com/yurythx/sismarq/internal/modules/acervo/infrastructure"
	"github.com/yurythx/sismarq/internal/modules/acervo/transport"
	"github.com/yurythx/sismarq/internal/platform/auth"
	"github.com/yurythx/sismarq/internal/platform/database/dbtest"
	"github.com/yurythx/sismarq/internal/platform/storage/storagetest"
)

// Matriz de falhas: cada chamada ao repositório falha (ou envenena a
// transação logo depois) e o caso de uso tem de propagar o erro.
func TestEveryRepositoryFailureIsPropagated(t *testing.T) {
	e := novoEnv(t)
	ctx := context.Background()
	e.agora = dia("2021-01-01")
	real := e.real()
	predio, err := real.SalvarLocal(ctx, admin, nil, domain.Local{Tipo: domain.LocalPredio, Nome: "Falhas " + uuid.NewString()[:8], Ativo: true})
	if err != nil {
		t.Fatal(err)
	}
	comLocal := func() application.Detalhe {
		u := e.nova(codElim)
		u.LocalID = &predio.ID
		d, err := real.Criar(ctx, admin, u)
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	solicitante := dbtest.User(t, e.pool)
	planilha := "tipo;codigo_classificacao;unidade;data_inicial;data_final;local\nCAIXA;" + codElim + ";" + e.unidade.String() +
		";01/01/2018;31/12/2018;" + predio.Caminho + "\n"
	_, membros := e.cspad(t, real)
	cspad := avaliador(e.unidade)
	orig, base := e.agora, dia("2026-09-01")
	// rito monta uma relação até a etapa dada: 0 rascunho, 1 aprovada,
	// 2 autorizada, 3 em edital, 4 com pedido, 5 amostrada.
	rito := func(etapa int) domain.Listagem {
		e.agora = base
		defer func() { e.agora = orig }()
		u, u2 := comLocal(), comLocal()
		l, err := real.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, ParecerJuridico: "P", Unidades: []uuid.UUID{u.ID, u2.ID}})
		passos := []func() (domain.Listagem, error){
			func() (domain.Listagem, error) {
				return real.Aprovar(ctx, cspad, l.ID, application.Ata{Numero: "1", Em: base, Presentes: membros[:3]})
			},
			func() (domain.Listagem, error) { return real.Autorizar(ctx, admin, l.ID) },
			func() (domain.Listagem, error) {
				return real.PublicarEdital(ctx, cspad, l.ID, application.Edital{Numero: "1", PublicadoEm: base})
			},
			func() (domain.Listagem, error) {
				return real.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{UnidadeID: &u.ID, Requerente: "R", Pedido: "P", RecebidoEm: base})
			},
			func() (domain.Listagem, error) { return real.Amostrar(ctx, cspad, l.ID, "c", 1) },
		}
		for i := 0; err == nil && i < etapa; i++ {
			l, err = passos[i]()
		}
		if err != nil {
			t.Fatal(err)
		}
		return l
	}
	// naData roda a operação com o relógio na data dada.
	naData := func(d time.Time, op func(*application.Service) error) func(*application.Service) error {
		return func(s *application.Service) error {
			e.agora = d
			defer func() { e.agora = orig }()
			return op(s)
		}
	}
	// recolhimento: extinção (qualquer fase) de uma unidade de guarda
	// permanente da origem, com um gestor lotado para a cópia do registro.
	e.gestorLotado(t, e.unidade)
	recolhimento := func() domain.Recolhimento {
		u := e.nova(codPerm)
		u.LocalID = &predio.ID
		d, err := real.Criar(ctx, admin, u)
		if err != nil {
			t.Fatal(err)
		}
		r, err := real.CriarRecolhimento(ctx, admin, application.NovoRecolhimento{OrigemID: e.unidade, Motivo: domain.MotivoExtincao,
			Fundamento: "Lei", Preparo: preparo, Unidades: []uuid.UUID{d.ID}})
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	// norma: com a primeira redação em rascunho ou já publicada.
	norma := func(publicada bool) domain.Norma {
		n, err := real.CriarNorma(ctx, admin, numeroNum(), domain.Redacao{Ementa: "e", Texto: "t"})
		if err == nil && publicada {
			n, err = real.PublicarNorma(ctx, admin, n.ID, application.Publicacao{Em: e.agora})
		}
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	// privado: até a etapa dada (0 cadastrado, 1 com parecer, 2 declarado,
	// 3 com alienação pendente).
	privado := func(etapa int) domain.Privado {
		p, err := real.CadastrarPrivado(ctx, admin, fundoFamilia())
		passos := []func() (domain.Privado, error){
			func() (domain.Privado, error) {
				return real.RegistrarParecer(ctx, admin, p.ID, application.Parecer{Comissao: "c", Texto: "t", Favoravel: true, Em: e.agora})
			},
			func() (domain.Privado, error) { return real.Declarar(ctx, admin, p.ID, "D", "", e.agora) },
			func() (domain.Privado, error) {
				return real.NotificarAlienacao(ctx, admin, p.ID, domain.Alienacao{NotificadoEm: e.agora})
			},
		}
		for i := 0; err == nil && i < etapa; i++ {
			p, err = passos[i]()
		}
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	// digital: o representante na situação pedida (processado com o
	// pipeline falso).
	dstore := &wormMem{Memory: storagetest.New(), dias: map[string]int{}}
	dcfg := func(s *application.Service) *application.Service {
		return s.WithDigitalizacao(application.ConfigDigital{Store: dstore, Bucket: "app", WORMBucket: "worm", MaxBytes: 1 << 20,
			URLExpiry: time.Minute, DPITexto: 300, DPIPlanta: 600, TmpDir: t.TempDir(), Pipeline: &pipeFalso{paginas: 1, dpi: 300}})
	}
	dreal := dcfg(e.real())
	dquem := digitalizador(t, e, e.unidade)
	digital := func(pronto bool) application.DetalheDigital {
		u := comLocal()
		e.limparFila(t)
		d := enviar(t, dreal, dstore, dquem, u.ID, "TEXTO")
		if pronto {
			if _, err := dreal.ProcessarProximo(ctx); err != nil {
				t.Fatal(err)
			}
		}
		return d
	}
	bloqueada := func() (application.Detalhe, uuid.UUID) {
		d := comLocal()
		d, err := real.Bloquear(ctx, admin, d.ID, "1", "")
		if err != nil {
			t.Fatal(err)
		}
		return d, d.Bloqueios[0].ID
	}
	ops := map[string]func() func(*application.Service) error{
		"Criar": func() func(*application.Service) error {
			return func(s *application.Service) error {
				u := e.nova(codElim)
				u.LocalID = &predio.ID
				_, err := s.Criar(ctx, admin, u)
				return err
			}
		},
		"Alterar": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error {
				u := e.nova(codElim)
				u.LocalID = &predio.ID
				_, err := s.Alterar(ctx, admin, d.ID, u)
				return err
			}
		},
		"Transferir": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Transferir(ctx, admin, d.ID, ""); return err }
		},
		"Localizar": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error {
				_, err := s.Localizar(ctx, admin, d.ID, &predio.ID, "p")
				return err
			}
		},
		"Bloquear": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Bloquear(ctx, admin, d.ID, "1", ""); return err }
		},
		"Liberar": func() func(*application.Service) error {
			_, bid := bloqueada()
			return func(s *application.Service) error { _, err := s.Liberar(ctx, admin, bid, "ok"); return err }
		},
		"Get": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error { _, err := s.Get(ctx, admin, d.ID); return err }
		},
		"List": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.List(ctx, admin, domain.Filtro{}, pagination.New(1, 5, 100))
				return err
			}
		},
		"Painel": func() func(*application.Service) error {
			return func(s *application.Service) error { _, err := s.Painel(ctx, admin, domain.Filtro{}); return err }
		},
		"Emprestar": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error {
				_, err := s.Emprestar(ctx, admin, d.ID, solicitante, "consulta", dia("2021-01-10"))
				return err
			}
		},
		"Devolver": func() func(*application.Service) error {
			d := comLocal()
			d, err := real.Emprestar(ctx, admin, d.ID, solicitante, "consulta", dia("2021-01-10"))
			if err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error {
				_, err := s.Devolver(ctx, admin, d.Emprestimos[0].ID, "ok")
				return err
			}
		},
		"Emprestimos": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Emprestimos(ctx, admin, true, pagination.New(1, 5, 100))
				return err
			}
		},
		"CriarGuia": func() func(*application.Service) error {
			d := comLocal()
			return func(s *application.Service) error {
				_, err := s.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{d.ID}, "")
				return err
			}
		},
		"ReceberGuia": func() func(*application.Service) error {
			d := comLocal()
			g, err := real.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{d.ID}, "")
			if err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error { _, err := s.ReceberGuia(ctx, admin, g.ID); return err }
		},
		"CancelarGuia": func() func(*application.Service) error {
			d := comLocal()
			g, err := real.CriarGuia(ctx, admin, e.unidade, []uuid.UUID{d.ID}, "")
			if err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error { _, err := s.CancelarGuia(ctx, admin, g.ID); return err }
		},
		"Guias": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Guias(ctx, admin, pagination.New(1, 5, 100))
				return err
			}
		},
		"CriarRecolhimento": func() func(*application.Service) error {
			u := e.nova(codPerm)
			u.LocalID = &predio.ID
			d, err := real.Criar(ctx, admin, u)
			if err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error {
				_, err := s.CriarRecolhimento(ctx, admin, application.NovoRecolhimento{OrigemID: e.unidade, Motivo: domain.MotivoExtincao,
					Fundamento: "Lei", Preparo: preparo, Unidades: []uuid.UUID{d.ID}})
				return err
			}
		},
		"ReceberRecolhimento": func() func(*application.Service) error {
			r := recolhimento()
			return func(s *application.Service) error {
				_, err := s.ReceberRecolhimento(ctx, admin, r.ID, application.Entrada{LocalID: &predio.ID})
				return err
			}
		},
		"DevolverRecolhimento": func() func(*application.Service) error {
			r := recolhimento()
			return func(s *application.Service) error {
				_, err := s.DevolverRecolhimento(ctx, admin, r.ID, "falta")
				return err
			}
		},
		"CancelarRecolhimento": func() func(*application.Service) error {
			r := recolhimento()
			return func(s *application.Service) error { _, err := s.CancelarRecolhimento(ctx, admin, r.ID); return err }
		},
		"Recolhimentos": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Recolhimentos(ctx, admin, pagination.New(1, 5, 100))
				return err
			}
		},
		"CriarNorma": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.CriarNorma(ctx, admin, numeroNum(), domain.Redacao{Ementa: "e", Texto: "t"})
				return err
			}
		},
		"SalvarRascunho": func() func(*application.Service) error {
			n := norma(false)
			return func(s *application.Service) error {
				_, err := s.SalvarRascunho(ctx, admin, n.ID, domain.Redacao{Ementa: "e", Texto: "t2", Motivo: "m"})
				return err
			}
		},
		"PublicarNorma": func() func(*application.Service) error {
			n := norma(false)
			return func(s *application.Service) error {
				_, err := s.PublicarNorma(ctx, admin, n.ID, application.Publicacao{Em: e.agora})
				return err
			}
		},
		"DescartarRascunho": func() func(*application.Service) error {
			n := norma(false)
			return func(s *application.Service) error { _, err := s.DescartarRascunho(ctx, admin, n.ID); return err }
		},
		"DescartarAlteracao": func() func(*application.Service) error {
			n := norma(true)
			if _, err := real.SalvarRascunho(ctx, admin, n.ID, domain.Redacao{Ementa: "e", Texto: "t2", Motivo: "m"}); err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error { _, err := s.DescartarRascunho(ctx, admin, n.ID); return err }
		},
		"RevogarNorma": func() func(*application.Service) error {
			n := norma(true)
			if _, err := real.SalvarRascunho(ctx, admin, n.ID, domain.Redacao{Ementa: "e", Texto: "t2", Motivo: "m"}); err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error {
				_, err := s.RevogarNorma(ctx, admin, n.ID, e.agora, "IN")
				return err
			}
		},
		"Normas": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Normas(ctx, admin, "", pagination.New(1, 5, 100))
				return err
			}
		},
		"NormaPublica": func() func(*application.Service) error {
			n := norma(true)
			return func(s *application.Service) error { _, err := s.NormaPublica(ctx, n.ID); return err }
		},
		"CadastrarPrivado": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.CadastrarPrivado(ctx, admin, fundoFamilia())
				return err
			}
		},
		"AlterarPrivado": func() func(*application.Service) error {
			p := privado(0)
			return func(s *application.Service) error {
				_, err := s.AlterarPrivado(ctx, admin, p.ID, fundoFamilia())
				return err
			}
		},
		"RegistrarParecer": func() func(*application.Service) error {
			p := privado(0)
			return func(s *application.Service) error {
				_, err := s.RegistrarParecer(ctx, admin, p.ID, application.Parecer{Comissao: "c", Texto: "t", Favoravel: true, Em: e.agora})
				return err
			}
		},
		"MudarGuarda": func() func(*application.Service) error {
			p := privado(2)
			return func(s *application.Service) error {
				_, err := s.MudarGuarda(ctx, admin, p.ID, "DOADO", "t", e.agora)
				return err
			}
		},
		"NotificarAlienacao": func() func(*application.Service) error {
			p := privado(2)
			return func(s *application.Service) error {
				_, err := s.NotificarAlienacao(ctx, admin, p.ID, domain.Alienacao{NotificadoEm: e.agora})
				return err
			}
		},
		"DecidirAlienacao": func() func(*application.Service) error {
			p := privado(3)
			return func(s *application.Service) error {
				_, err := s.DecidirAlienacao(ctx, admin, p.ID, p.Alienacoes[0].ID, "INTERESSE", "f")
				return err
			}
		},
		"Privados": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Privados(ctx, admin, "", pagination.New(1, 5, 100))
				return err
			}
		},
		"RegistrarDigital": func() func(*application.Service) error {
			u := comLocal()
			tk, err := dreal.UploadDigital(ctx, dquem, u.ID, "a.pdf")
			if err != nil {
				t.Fatal(err)
			}
			_ = dstore.Put(ctx, "app", tk.ObjectKey, strings.NewReader(pdfOriginal), int64(len(pdfOriginal)), "application/pdf")
			return func(s *application.Service) error {
				_, err := dcfg(s).RegistrarDigital(ctx, dquem, u.ID, application.NovoEnvio{ObjectKey: tk.ObjectKey, TipoDoc: "TEXTO", Equipamento: "e", DPI: 300})
				return err
			}
		},
		"UploadDigital": func() func(*application.Service) error {
			u := comLocal()
			return func(s *application.Service) error {
				_, err := dcfg(s).UploadDigital(ctx, dquem, u.ID, "a.pdf")
				return err
			}
		},
		"Reprocessar": func() func(*application.Service) error {
			d := digital(false)
			if _, err := e.pool.Exec(ctx, `UPDATE acervo_digitais SET situacao = 'FALHOU' WHERE id = $1`, d.ID); err != nil {
				t.Fatal(err)
			}
			return func(s *application.Service) error { _, err := dcfg(s).Reprocessar(ctx, admin, d.ID); return err }
		},
		"VerificarAgora": func() func(*application.Service) error {
			d := digital(true)
			return func(s *application.Service) error { _, err := dcfg(s).VerificarAgora(ctx, admin, d.ID); return err }
		},
		"AbrirArquivo": func() func(*application.Service) error {
			d := digital(true)
			return func(s *application.Service) error { _, err := dcfg(s).AbrirArquivo(ctx, dquem, d.ID); return err }
		},
		"Digitais": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := dcfg(s).Digitais(ctx, dquem, nil, false, pagination.New(1, 5, 100))
				return err
			}
		},
		"Importar": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.Importar(ctx, admin, planilha, "", false)
				return err
			}
		},
		"ImportarAplicar": func() func(*application.Service) error {
			sim, err := real.Importar(ctx, admin, planilha, "", false)
			if err != nil || len(sim.Erros) > 0 {
				t.Fatalf("simulação: %+v %v", sim, err)
			}
			return func(s *application.Service) error {
				_, err := s.Importar(ctx, admin, planilha, sim.Hash, true)
				return err
			}
		},
		"Comissoes": func() func(*application.Service) error {
			return func(s *application.Service) error { _, err := s.Comissoes(ctx); return err }
		},
		"SalvarComissao": func() func(*application.Service) error {
			return func(s *application.Service) error {
				c := domain.Comissao{Tipo: domain.ComissaoCentral, Nome: "CCPAD", Ativa: false, Membros: []domain.Membro{
					{UserID: membros[0], Papel: domain.PapelPresidente, Area: domain.AreaArquivologia},
					{UserID: membros[1], Papel: domain.PapelMembro, Area: domain.AreaJuridico},
					{UserID: membros[2], Papel: domain.PapelMembro, Area: domain.AreaAdmFinancas},
					{UserID: membros[3], Papel: domain.PapelMembro, Area: domain.AreaArquivoPublico},
					{UserID: membros[4], Papel: domain.PapelMembro, Area: domain.AreaOutra}}}
				salva, err := s.SalvarComissao(ctx, admin, nil, c)
				if err != nil {
					return err
				}
				_, err = s.SalvarComissao(ctx, admin, &salva.ID, c)
				return err
			}
		},
		"CriarListagem": func() func(*application.Service) error {
			u := comLocal()
			return naData(base, func(s *application.Service) error {
				_, err := s.CriarListagem(ctx, cspad, application.NovaListagem{UnidadeID: e.unidade, Unidades: []uuid.UUID{u.ID}})
				return err
			})
		},
		"AlterarListagem": func() func(*application.Service) error {
			l, u := rito(0), comLocal()
			return naData(base, func(s *application.Service) error {
				_, err := s.AlterarListagem(ctx, cspad, l.ID, application.NovaListagem{ParecerJuridico: "P2", Unidades: []uuid.UUID{u.ID}},
					[]uuid.UUID{l.Itens[0].ID})
				return err
			})
		},
		"Aprovar": func() func(*application.Service) error {
			l := rito(0)
			return naData(base, func(s *application.Service) error {
				_, err := s.Aprovar(ctx, cspad, l.ID, application.Ata{Numero: "1", Em: base, Presentes: membros[:3]})
				return err
			})
		},
		"Autorizar": func() func(*application.Service) error {
			l := rito(1)
			return naData(base, func(s *application.Service) error { _, err := s.Autorizar(ctx, admin, l.ID); return err })
		},
		"DevolverListagem": func() func(*application.Service) error {
			l := rito(1)
			return naData(base, func(s *application.Service) error {
				_, err := s.DevolverListagem(ctx, admin, l.ID, "motivo")
				return err
			})
		},
		"PublicarEdital": func() func(*application.Service) error {
			l := rito(2)
			return naData(base, func(s *application.Service) error {
				_, err := s.PublicarEdital(ctx, cspad, l.ID, application.Edital{Numero: "1", PublicadoEm: base})
				return err
			})
		},
		"RegistrarPedido": func() func(*application.Service) error {
			l := rito(3)
			return naData(base, func(s *application.Service) error {
				_, err := s.RegistrarPedido(ctx, cspad, l.ID, domain.Pedido{Requerente: "R", Pedido: "P", RecebidoEm: base})
				return err
			})
		},
		"DecidirPedido": func() func(*application.Service) error {
			l := rito(4)
			return naData(base, func(s *application.Service) error {
				_, err := s.DecidirPedido(ctx, cspad, l.ID, l.Pedidos[0].ID, true, true, "ok")
				return err
			})
		},
		"Amostrar": func() func(*application.Service) error {
			l := rito(3)
			return naData(base, func(s *application.Service) error { _, err := s.Amostrar(ctx, cspad, l.ID, "c", 1); return err })
		},
		"RegistrarTermo": func() func(*application.Service) error {
			l := rito(3)
			e.agora = base
			if _, err := real.Amostrar(ctx, cspad, l.ID, "c", 1); err != nil {
				t.Fatal(err)
			}
			e.agora = orig
			return naData(base.AddDate(0, 0, 31), func(s *application.Service) error {
				_, err := s.RegistrarTermo(ctx, cspad, l.ID, application.Termo{EliminadoEm: base.AddDate(0, 0, 31), Acompanhamento: "X"})
				return err
			})
		},
		"CancelarListagem": func() func(*application.Service) error {
			l := rito(0)
			return naData(base, func(s *application.Service) error { _, err := s.CancelarListagem(ctx, cspad, l.ID); return err })
		},
		"Listagens": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, _, err := s.Listagens(ctx, cspad, pagination.New(1, 5, 100))
				return err
			}
		},
		"ListLocais": func() func(*application.Service) error {
			return func(s *application.Service) error { _, err := s.ListLocais(ctx); return err }
		},
		"CriarLocal": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.SalvarLocal(ctx, admin, nil, domain.Local{ParentID: &predio.ID, Tipo: domain.LocalSala, Nome: uuid.NewString(), Ativo: true})
				return err
			}
		},
		"AlterarLocal": func() func(*application.Service) error {
			return func(s *application.Service) error {
				_, err := s.SalvarLocal(ctx, admin, &predio.ID, predio)
				return err
			}
		},
	}
	for name, o := range ops {
		probe := &faultRepo{inner: infrastructure.NewRepository()}
		if err := o()(e.svc(probe)); err != nil {
			t.Fatalf("%s sem falha: %v", name, err)
		}
		for i := range probe.trace {
			for _, poison := range []bool{false, true} {
				r := &faultRepo{inner: infrastructure.NewRepository()}
				if poison {
					r.poisonAt = i + 1
				} else {
					r.failAt = i + 1
				}
				if err := o()(e.svc(r)); err == nil && !r.noTx {
					t.Errorf("%s: falha (veneno=%v) na chamada %d (%s) foi engolida", name, poison, i+1, probe.trace[i])
				}
			}
		}
	}
}

// tudo responde qualquer código (o recálculo percorre todo o banco de teste).
type tudo struct{ vigente bool }

func (t tudo) Serie(_ context.Context, codigo string) (domain.Serie, error) {
	if !t.vigente {
		return domain.Serie{}, domain.ErrSerieInvalida
	}
	return domain.Serie{Codigo: codigo, CorrenteAnos: anos(1), IntermAnos: anos(1), Destinacao: domain.DestinacaoEliminacao, Vigente: true}, nil
}

// O recálculo percorre o banco inteiro: falha em cada ponto, de forma
// determinística (série conhecida: lista, bloqueio, gravação; série
// inexistente: lista, gravação da situação).
func TestRecalculoPropagaFalhas(t *testing.T) {
	e := novoEnv(t)
	e.criar(codElim)
	ctx := context.Background()
	for _, c := range []struct {
		ttdd   domain.TTDD
		failAt int
	}{{tudo{true}, 1}, {tudo{true}, 2}, {tudo{true}, 3}, {tudo{false}, 2}} {
		s := application.NewService(e.pool, &faultRepo{inner: infrastructure.NewRepository(), failAt: c.failAt}, c.ttdd, nil, logger)
		if _, err := s.Recalcular(ctx); !errors.Is(err, errBoom) {
			t.Errorf("falha na chamada %d (vigente=%v) engolida: %v", c.failAt, c.ttdd, err)
		}
	}
}

// router monta as rotas com a identidade dada já no contexto.
func router(s *application.Service, identity auth.Identity) http.Handler {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req.WithContext(auth.WithIdentity(req.Context(), identity)))
		})
	})
	transport.NewHandlers(s, logger, 100).RegisterRoutes(r)
	return r
}

func do(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandlers(t *testing.T) {
	e := novoEnv(t)
	e.agora = dia("2021-01-01")
	h := router(e.real(), admin)
	corpo := func(extra string) string {
		return `{"tipo":"CAIXA","codigo_classificacao":"` + codElim + `","unidade_id":"` + e.unidade.String() +
			`","data_inicial":"2018-01-01","data_final":"2018-12-31","quantidade":1` + extra + `}`
	}
	rec := do(h, http.MethodPost, "/acervo/unidades", corpo(`,"data_evento":"2019-01-01","contas_aprovadas_em":"2019-05-01"`))
	if rec.Code != http.StatusCreated {
		t.Fatalf("criar: %d %s", rec.Code, rec.Body)
	}
	id := strings.Split(strings.Split(rec.Body.String(), `"id":"`)[1], `"`)[0]
	predio := do(h, http.MethodPost, "/acervo/locais", `{"tipo":"PREDIO","nome":"H `+uuid.NewString()[:8]+`"}`)
	if predio.Code != http.StatusCreated {
		t.Fatalf("local: %d %s", predio.Code, predio.Body)
	}
	localID := strings.Split(strings.Split(predio.Body.String(), `"id":"`)[1], `"`)[0]
	casos := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodGet, "/acervo/unidades?q=UA&fase=corrente&unidade_id=" + e.unidade.String(), "", http.StatusOK},
		{http.MethodGet, "/acervo/unidades?unidade_id=x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/painel", "", http.StatusOK},
		{http.MethodGet, "/acervo/painel?unidade_id=x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/unidades/" + id, "", http.StatusOK},
		{http.MethodGet, "/acervo/unidades/x", "", http.StatusBadRequest},
		{http.MethodGet, "/acervo/unidades/" + uuid.NewString(), "", http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "2018-01-01", "01/01/2018", 1), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "2018-12-31", "x", 1), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", corpo(`,"data_evento":"x"`), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", corpo(`,"contas_aprovadas_em":"x"`), http.StatusUnprocessableEntity},
		{http.MethodPost, "/acervo/unidades", strings.Replace(corpo(""), "CAIXA", "GAVETA", 1), http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/unidades/x", corpo(""), http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + id, `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + uuid.NewString(), corpo(""), http.StatusNotFound},
		{http.MethodPut, "/acervo/unidades/" + id, corpo(`,"observacoes":"revisada"`), http.StatusOK},
		{http.MethodPut, "/acervo/unidades/x/local", `{}`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + id + "/local", `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/unidades/" + uuid.NewString() + "/local", `{}`, http.StatusNotFound},
		{http.MethodPut, "/acervo/unidades/" + id + "/local", `{"local_id":"` + localID + `","posicao":"P1"}`, http.StatusOK},
		{http.MethodPost, "/acervo/unidades/x/transferir", `{}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/transferir", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + uuid.NewString() + "/transferir", `{}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades/" + id + "/transferir", `{"observacao":"ok"}`, http.StatusOK},
		{http.MethodPost, "/acervo/unidades/x/bloqueios", `{"processo":"1"}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + id + "/bloqueios", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/unidades/" + uuid.NewString() + "/bloqueios", `{"processo":"1"}`, http.StatusNotFound},
		{http.MethodPost, "/acervo/unidades/" + id + "/bloqueios", `{"processo":"0001"}`, http.StatusCreated},
		{http.MethodPost, "/acervo/bloqueios/x/liberar", `{"motivo":"ok"}`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/bloqueios/" + uuid.NewString() + "/liberar", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/bloqueios/" + uuid.NewString() + "/liberar", `{"motivo":"ok"}`, http.StatusNotFound},
		{http.MethodGet, "/acervo/locais", "", http.StatusOK},
		{http.MethodPost, "/acervo/locais", `{`, http.StatusBadRequest},
		{http.MethodPost, "/acervo/locais", `{"tipo":"COFRE","nome":"x"}`, http.StatusUnprocessableEntity},
		{http.MethodPut, "/acervo/locais/x", `{"tipo":"PREDIO","nome":"x"}`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/locais/" + localID, `{`, http.StatusBadRequest},
		{http.MethodPut, "/acervo/locais/" + uuid.NewString(), `{"tipo":"PREDIO","nome":"x"}`, http.StatusNotFound},
		{http.MethodPut, "/acervo/locais/" + localID, `{"tipo":"PREDIO","nome":"Renomeado ` + uuid.NewString()[:6] + `","ativo":true}`, http.StatusOK},
	}
	for _, c := range casos {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s %s: %d, esperado %d (%s)", c.method, c.path, rec.Code, c.want, rec.Body)
		}
	}
	// Liberar o bloqueio criado acima.
	ficha := do(h, http.MethodGet, "/acervo/unidades/"+id, "").Body.String()
	bid := strings.Split(strings.Split(ficha, `"bloqueios":[{"id":"`)[1], `"`)[0]
	if rec := do(h, http.MethodPost, "/acervo/bloqueios/"+bid+"/liberar", `{"motivo":"Encerrada"}`); rec.Code != http.StatusOK {
		t.Fatalf("liberar: %d %s", rec.Code, rec.Body)
	}
	// Banco fora: a listagem e o painel respondem 500 sem vazar o erro.
	for _, p := range []string{"/acervo/unidades", "/acervo/painel", "/acervo/locais"} {
		down := router(e.svc(&faultRepo{inner: infrastructure.NewRepository(), failAt: 1}), admin)
		if rec := do(down, http.MethodGet, p, ""); rec.Code != http.StatusInternalServerError || strings.Contains(rec.Body.String(), errBoom.Error()) {
			t.Errorf("%s com o banco fora: %d %s", p, rec.Code, rec.Body)
		}
	}
}
