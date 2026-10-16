package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
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
	"github.com/yurythx/sismarq/internal/platform/notificacoes"
	"github.com/yurythx/sismarq/internal/platform/storage"
)

// Fase 5 (ADR 037): digitalização com valor probatório. O original vai
// direto do navegador ao MinIO (URL pré-assinada); o worker gera o PDF/A
// com OCR e metadados, calcula o SHA-256, guarda com object lock e apaga o
// original; outro worker confere periodicamente se o hash ainda bate.

// Eventos da digitalização (outbox).
const (
	EventDigitalRecebido   = "acervo.digital.recebido"
	EventDigitalPronto     = "acervo.digital.pronto"
	EventDigitalReprovado  = "acervo.digital.reprovado"
	EventDigitalFalhou     = "acervo.digital.falhou"
	EventDigitalReaberto   = "acervo.digital.reprocessar"
	EventIntegridadeFalhou = "acervo.digital.integridade_falhou"
)

const (
	prefixoDigital  = "acervo/digitalizacao"
	verificarLote   = 100
	verificacoesMax = 20
)

// TiposDigital: o scanner entrega PDF (uma imagem por página).
var TiposDigital = []string{"application/pdf"}

// ConfigDigital parametriza a digitalização (config.AcervoDigitalConfig).
type ConfigDigital struct {
	Store       storage.Provider
	Bucket      string // onde o original chega (o bucket da aplicação)
	WORMBucket  string // onde o PDF/A fica, com object lock
	MaxBytes    int64
	URLExpiry   time.Duration
	DPITexto    int
	DPIPlanta   int
	TmpDir      string
	Pipeline    domain.Pipeline
	Travado     time.Duration // PROCESSANDO há mais que isto volta à fila
	Verificacao time.Duration // intervalo entre verificações de cada arquivo
}

// WithDigitalizacao liga a digitalização.
func (s *Service) WithDigitalizacao(c ConfigDigital) *Service {
	if c.Travado <= 0 {
		c.Travado = 3 * time.Hour
	}
	if c.Verificacao <= 0 {
		c.Verificacao = 24 * time.Hour
	}
	s.dig = c
	return s
}

// unidadeParaDigitalizar lê a unidade e exige a gestão dela.
func (s *Service) unidadeParaDigitalizar(ctx context.Context, db database.DBTX, identity auth.Identity, id uuid.UUID, lock bool) (domain.Unidade, domain.Calculo, error) {
	u, c, err := s.repo.GetUnidade(ctx, db, id, lock)
	if err != nil {
		return u, c, err
	}
	if !le(identity, u.UnidadeID) {
		return u, c, domain.ErrNaoEncontrada
	}
	if !gerePara(identity, u) {
		return u, c, domain.ErrForaDoEscopo
	}
	return u, c, domain.PodeDigitalizar(u)
}

// UploadDigital autoriza o envio direto do PDF digitalizado.
func (s *Service) UploadDigital(ctx context.Context, identity auth.Identity, unidade uuid.UUID, filename string) (modkit.UploadTicket, error) {
	if _, _, err := s.unidadeParaDigitalizar(ctx, s.pool, identity, unidade, false); err != nil {
		return modkit.UploadTicket{}, MapError(err)
	}
	return modkit.NewUpload(ctx, s.dig.Store, s.dig.Bucket, prefixoDigital+"/"+unidade.String(), filename, "application/pdf",
		TiposDigital, s.dig.URLExpiry)
}

// NovoEnvio é o registro do PDF enviado, com os metadados do Decreto.
type NovoEnvio struct {
	ObjectKey      string
	Descricao      string
	TipoDoc        string
	Equipamento    string
	DPI            int
	ResponsavelCPF string
}

// sha256Objeto lê o objeto inteiro e devolve o hash.
func (s *Service) sha256Objeto(ctx context.Context, bucket, key string) (string, error) {
	rc, err := s.dig.Store.Get(ctx, bucket, key)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	h := sha256.New()
	if _, err := io.Copy(h, rc); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// RegistrarDigital confirma o upload, guarda o hash do original e põe na
// fila do pipeline.
func (s *Service) RegistrarDigital(ctx context.Context, identity auth.Identity, unidade uuid.UUID, in NovoEnvio) (DetalheDigital, error) {
	if identity.UserID == uuid.Nil {
		return DetalheDigital{}, MapError(&domain.ErrValidacao{Campo: "responsavel", Motivo: "usuário sem cadastro local não digitaliza"})
	}
	cpf, err := domain.MascararCPF(in.ResponsavelCPF)
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	d, err := domain.NovoDigital(domain.Digital{UnidadeID: unidade, Descricao: in.Descricao, TipoDoc: in.TipoDoc, Equipamento: in.Equipamento,
		DPIDeclarado: in.DPI, ResponsavelID: identity.UserID, ResponsavelNome: identity.Name, ResponsavelCPF: cpf})
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	if _, _, err := s.unidadeParaDigitalizar(ctx, s.pool, identity, unidade, false); err != nil {
		return DetalheDigital{}, MapError(err)
	}
	info, err := modkit.ConfirmUpload(ctx, s.dig.Store, s.dig.Bucket, in.ObjectKey, prefixoDigital+"/"+unidade.String(), s.dig.MaxBytes, TiposDigital)
	if err != nil {
		return DetalheDigital{}, err
	}
	if d.OriginalSHA256, err = s.sha256Objeto(ctx, s.dig.Bucket, in.ObjectKey); err != nil {
		return DetalheDigital{}, MapError(fmt.Errorf("acervo: leitura do original: %w", err))
	}
	agora := s.agora().UTC()
	d.OriginalKey, d.OriginalBytes, d.DigitalizadoEm, d.CriadoEm = in.ObjectKey, info.Size, agora, agora
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		u, _, err := s.unidadeParaDigitalizar(ctx, tx, identity, unidade, true)
		if err != nil {
			return err
		}
		d.CodigoClass = u.CodigoClassificacao
		if d.Numero, err = s.repo.InsertDigital(ctx, tx, d); err != nil {
			return err
		}
		return s.registrarDigital(ctx, tx, identity, EventDigitalRecebido, u, d,
			"Representante digital "+d.Numero+" recebido para o pipeline (Decreto 10.278/2020)")
	})
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	return s.Digital(ctx, identity, d.ID)
}

// registrarDigital grava a movimentação da unidade, o evento e a auditoria.
func (s *Service) registrarDigital(ctx context.Context, tx pgx.Tx, identity auth.Identity, evento string, u domain.Unidade, d domain.Digital, descricao string) error {
	if err := s.repo.InsertMovimentacao(ctx, tx, u.ID, domain.Movimentacao{Tipo: domain.MovDigitalizacao, Por: ator(identity), Descricao: descricao}); err != nil {
		return err
	}
	if err := s.outbox.Write(ctx, tx, evento, "acervo_digital", d.ID.String(), uuid.Nil,
		map[string]any{"id": d.ID, "numero": d.Numero, "unidade_arquivamento_id": u.ID, "situacao": d.Situacao}); err != nil {
		return err
	}
	return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, evento, "acervo_digital", d.ID.String(), nil,
		map[string]any{"numero": d.Numero, "situacao": d.Situacao, "pdfa_sha256": d.PDFASHA256, "original_sha256": d.OriginalSHA256}))
}

// ProcessarProximo roda o pipeline no próximo da fila. Devolve false com a
// fila vazia.
func (s *Service) ProcessarProximo(ctx context.Context) (bool, error) {
	var d domain.Digital
	var achou bool
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		var err error
		d, achou, err = s.repo.ProximoDigital(ctx, tx, s.agora().Add(-s.dig.Travado))
		return err
	})
	if err != nil || !achou {
		return false, err
	}
	res, texto, perr := s.pipeline(ctx, &d)
	agora := s.agora().UTC()
	d.ProcessadoEm = &agora
	evento, descricao := EventDigitalPronto, ""
	switch {
	case perr == nil && res != "":
		d.Situacao, d.Motivo, evento = domain.DigitalReprovado, res, EventDigitalReprovado
		descricao = "Representante digital " + d.Numero + " reprovado: " + res
	case perr == nil:
		d.Situacao, d.Motivo = domain.DigitalPronto, ""
		descricao = fmt.Sprintf("Representante digital %s pronto: %d página(s), %s, SHA-256 %s", d.Numero, d.Paginas, d.Conformidade, d.PDFASHA256)
	case d.Tentativas >= domain.MaxTentativas:
		d.Situacao, d.Motivo, evento = domain.DigitalFalhou, truncar(perr.Error(), 2000), EventDigitalFalhou
		descricao = "Representante digital " + d.Numero + " falhou no pipeline; o Arquivo Público pode reprocessar"
	default:
		// Tenta de novo na próxima volta do worker.
		d.Situacao, d.Motivo, d.ProcessadoEm = domain.DigitalEnviado, truncar(perr.Error(), 2000), nil
		s.logger.WarnContext(ctx, "acervo: pipeline falhou; nova tentativa", slog.String("digital", d.Numero), slog.Any("error", perr))
	}
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.repo.ConcluirDigital(ctx, tx, d, texto); err != nil {
			return err
		}
		if descricao == "" {
			return nil
		}
		u, _, err := s.repo.GetUnidade(ctx, tx, d.UnidadeID, false)
		if err != nil {
			return err
		}
		return s.registrarDigital(ctx, tx, auth.Identity{}, evento, u, d, descricao)
	})
	if err != nil {
		return true, err
	}
	if d.Situacao != domain.DigitalEnviado && d.Situacao != domain.DigitalFalhou {
		// Pronto (o PDF/A é a cópia oficial) ou reprovado (digitalizar de
		// novo): o original não fica.
		if derr := s.dig.Store.Delete(ctx, s.dig.Bucket, d.OriginalKey); derr != nil {
			s.logger.WarnContext(ctx, "acervo: original não apagado", slog.String("digital", d.Numero), slog.Any("error", derr))
		}
	}
	return true, nil
}

// criarArquivo é a costura dos testes para a falha de disco.
var criarArquivo = func(nome string) (*os.File, error) {
	return os.Create(nome) // #nosec G304 -- caminho no diretório temporário do próprio worker
}

func truncar(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n])
	}
	return s
}

// pipeline baixa o original, confere a resolução, gera e guarda o PDF/A.
// Devolve o motivo da reprovação ("" se aprovado) e o texto do OCR.
func (s *Service) pipeline(ctx context.Context, d *domain.Digital) (string, string, error) {
	dir, err := os.MkdirTemp(s.dig.TmpDir, "digital-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(dir)
	original := filepath.Join(dir, "original.pdf")
	if err := s.baixar(ctx, d.OriginalKey, original); err != nil {
		return "", "", err
	}
	paginas, dpi, err := s.dig.Pipeline.Medir(ctx, original)
	if err != nil {
		return "", "", err
	}
	d.DPIMedido = dpi
	minimo := domain.DPIMinimo(d.TipoDoc, s.dig.DPITexto, s.dig.DPIPlanta)
	switch {
	case dpi == 0:
		return "o PDF não traz imagens digitalizadas: o Decreto 10.278 trata da digitalização do papel", "", nil
	case dpi < minimo:
		return fmt.Sprintf("resolução de %d DPI abaixo do mínimo de %d DPI para %s; digitalize de novo", dpi, minimo,
			strings.ToLower(d.TipoDoc)), "", nil
	}
	u, c, err := s.repo.GetUnidade(ctx, s.pool, d.UnidadeID, false)
	if err != nil {
		return "", "", err
	}
	saida := filepath.Join(dir, "pdfa.pdf")
	res, err := s.dig.Pipeline.Processar(ctx, original, saida, domain.Metadados{Identificador: d.Numero, ResponsavelNome: d.ResponsavelNome,
		ResponsavelCPF: d.ResponsavelCPF, DigitalizadoEm: d.DigitalizadoEm, Equipamento: d.Equipamento, DPI: dpi,
		HashOriginal: d.OriginalSHA256, Classificacao: d.CodigoClass, Titulo: d.Numero + " — " + u.Assunto})
	if err != nil {
		return "", "", err
	}
	if res.Paginas != paginas {
		return "", "", fmt.Errorf("o PDF/A tem %d página(s) e o original %d", res.Paginas, paginas)
	}
	ate := domain.Retencao(u, c, s.hoje())
	d.PDFAKey = "digitais/" + d.UnidadeID.String() + "/" + d.Numero + ".pdf"
	if d.PDFASHA256, d.PDFABytes, err = s.guardar(ctx, saida, d.PDFAKey, ate); err != nil {
		return "", "", err
	}
	d.Paginas, d.Conformidade, d.RetencaoAte = res.Paginas, res.Conformidade, &ate
	return "", res.Texto, nil
}

func (s *Service) baixar(ctx context.Context, key, destino string) error {
	rc, err := s.dig.Store.Get(ctx, s.dig.Bucket, key)
	if err != nil {
		return err
	}
	defer rc.Close()
	f, err := criarArquivo(destino)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, rc); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// guardar calcula o hash e grava o PDF/A com object lock até ate (sem WORM
// no provider, grava comum e avisa).
func (s *Service) guardar(ctx context.Context, arquivo, key string, ate time.Time) (string, int64, error) {
	f, err := os.Open(arquivo) // #nosec G304 -- arquivo gerado pelo pipeline no diretório temporário
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	tamanho, err := io.Copy(h, f)
	if err != nil {
		return "", 0, err
	}
	// SectionReader relê do início sem Seek (o mesmo descritor).
	corpo := io.NewSectionReader(f, 0, tamanho)
	dias := max(1, int(time.Until(ate).Hours()/24)+1)
	if w, ok := s.dig.Store.(storage.WORMWriter); ok {
		// Idempotente: cria o bucket dedicado com object lock na primeira vez
		// (a retenção de cada objeto vem abaixo, por objeto).
		if err := w.EnsureImmutableBucket(ctx, s.dig.WORMBucket, 1); err != nil {
			return "", 0, err
		}
		err = w.PutImmutable(ctx, s.dig.WORMBucket, key, corpo, tamanho, "application/pdf", dias)
	} else {
		s.logger.WarnContext(ctx, "acervo: armazenamento sem object lock; o PDF/A foi gravado sem retenção imutável", slog.String("objeto", key))
		err = s.dig.Store.Put(ctx, s.dig.WORMBucket, key, corpo, tamanho, "application/pdf")
	}
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), tamanho, nil
}

// RunDigitalizacao processa a fila até esvaziar e espera o intervalo.
func (s *Service) RunDigitalizacao(intervalo time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()
		for {
			for {
				mais, err := s.ProcessarProximo(ctx)
				if err != nil {
					s.logger.ErrorContext(ctx, "acervo: digitalização falhou", slog.Any("error", err))
				}
				if !mais || err != nil || ctx.Err() != nil {
					break
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

// Reprocessar devolve à fila o que falhou no pipeline (só o Arquivo
// Público). O reprovado por resolução precisa de nova digitalização.
func (s *Service) Reprocessar(ctx context.Context, identity auth.Identity, id uuid.UUID) (DetalheDigital, error) {
	if err := central(identity); err != nil {
		return DetalheDigital{}, MapError(err)
	}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		d, err := s.repo.GetDigital(ctx, tx, id)
		if err != nil {
			return err
		}
		if d.Situacao != domain.DigitalFalhou {
			return domain.ErrDigitalSituacao
		}
		d.Situacao, d.Tentativas, d.Motivo = domain.DigitalEnviado, 0, ""
		if err := s.repo.ConcluirDigital(ctx, tx, d, ""); err != nil {
			return err
		}
		u, _, err := s.repo.GetUnidade(ctx, tx, d.UnidadeID, false)
		if err != nil {
			return err
		}
		return s.registrarDigital(ctx, tx, identity, EventDigitalReaberto, u, d, "Representante digital "+d.Numero+" devolvido ao pipeline")
	})
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	return s.Digital(ctx, identity, id)
}

// verificar recalcula o hash do PDF/A guardado e registra o resultado; no
// que não bate (ou sumiu), avisa a equipe do Arquivo Público.
func (s *Service) verificar(ctx context.Context, d domain.Digital) (domain.Verificacao, error) {
	v := domain.Verificacao{Em: s.agora().UTC()}
	lido, err := s.sha256DoWORM(ctx, d.PDFAKey)
	switch {
	case err != nil:
		v.Erro = truncar(err.Error(), 1000)
	default:
		v.SHA256Lido, v.Integro = lido, lido == d.PDFASHA256
	}
	var enviadas []notificacoes.Enviada
	err = database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.repo.RegistrarVerificacao(ctx, tx, d.ID, v); err != nil {
			return err
		}
		if v.Integro {
			return nil
		}
		if err := s.outbox.Write(ctx, tx, EventIntegridadeFalhou, "acervo_digital", d.ID.String(), uuid.Nil,
			map[string]any{"id": d.ID, "numero": d.Numero, "esperado": d.PDFASHA256, "lido": v.SHA256Lido, "erro": v.Erro}); err != nil {
			return err
		}
		equipe, err := s.repo.EquipeCentral(ctx, tx)
		if err != nil {
			return err
		}
		if enviadas, err = notificacoes.Registrar(ctx, tx, equipe, notificacoes.Nova{Modulo: "acervo",
			Titulo:   "Integridade: o representante digital " + d.Numero + " não confere",
			Mensagem: "O SHA-256 do PDF/A guardado não bate com o registrado (ou o arquivo sumiu). Confira a cópia e a trilha de auditoria.",
			Link:     "/acervo/digitais/" + d.ID.String(), Chave: "acervo:integridade:" + d.ID.String() + ":" + v.Em.Format("20060102")}); err != nil {
			return err
		}
		return audit.NewWriter(tx).Record(ctx, audit.Meta(ctx, EventIntegridadeFalhou, "acervo_digital", d.ID.String(),
			map[string]any{"pdfa_sha256": d.PDFASHA256}, map[string]any{"sha256_lido": v.SHA256Lido, "erro": v.Erro}))
	})
	if err != nil {
		return v, err
	}
	s.entregar(ctx, enviadas)
	if !v.Integro {
		s.logger.ErrorContext(ctx, "acervo: falha de integridade do representante digital", slog.String("digital", d.Numero),
			slog.String("esperado", d.PDFASHA256), slog.String("lido", v.SHA256Lido), slog.String("erro", v.Erro))
	}
	return v, nil
}

func (s *Service) sha256DoWORM(ctx context.Context, key string) (string, error) {
	rc, err := s.dig.Store.Get(ctx, s.dig.WORMBucket, key)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	h := sha256.New()
	if _, err := io.Copy(h, rc); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerificarIntegridade confere os arquivos não verificados no intervalo.
func (s *Service) VerificarIntegridade(ctx context.Context) (int, error) {
	lote, err := s.repo.ParaVerificar(ctx, s.pool, s.agora().Add(-s.dig.Verificacao), verificarLote)
	if err != nil {
		return 0, err
	}
	for i, d := range lote {
		if _, err := s.verificar(ctx, d); err != nil {
			return i, err
		}
	}
	return len(lote), nil
}

// RunIntegridade é o worker da verificação periódica (fixity).
func (s *Service) RunIntegridade(intervalo time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()
		for {
			if n, err := s.VerificarIntegridade(ctx); err != nil {
				s.logger.ErrorContext(ctx, "acervo: verificação de integridade falhou", slog.Any("error", err))
			} else if n > 0 {
				s.logger.InfoContext(ctx, "acervo: integridade verificada", slog.Int("arquivos", n))
			}
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
			}
		}
	}
}

// VerificarAgora confere um arquivo na hora (Arquivo Público).
func (s *Service) VerificarAgora(ctx context.Context, identity auth.Identity, id uuid.UUID) (DetalheDigital, error) {
	if err := central(identity); err != nil {
		return DetalheDigital{}, MapError(err)
	}
	d, err := s.repo.GetDigital(ctx, s.pool, id)
	if err == nil && d.Situacao != domain.DigitalPronto {
		err = domain.ErrDigitalSituacao
	}
	if err == nil {
		_, err = s.verificar(ctx, d)
	}
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	return s.Digital(ctx, identity, id)
}

// DetalheDigital é o representante com as últimas verificações.
type DetalheDigital struct {
	domain.Digital
	Verificacoes []domain.Verificacao `json:"verificacoes"`
	Unidade      string               `json:"unidade_numero"`
}

// acessoDigital lê o representante e confere a leitura da unidade.
func (s *Service) acessoDigital(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Digital, domain.Unidade, error) {
	d, err := s.repo.GetDigital(ctx, s.pool, id)
	if err != nil {
		return d, domain.Unidade{}, err
	}
	u, _, err := s.repo.GetUnidade(ctx, s.pool, d.UnidadeID, false)
	if err != nil {
		return d, u, err
	}
	if !le(identity, u.UnidadeID) && !auth.CanGlobal(identity, auth.PermAcervoCentral) {
		return d, u, domain.ErrDigitalNaoExiste
	}
	return d.ComEstadoPAdES(), u, nil
}

// Digital devolve o representante com as verificações.
func (s *Service) Digital(ctx context.Context, identity auth.Identity, id uuid.UUID) (DetalheDigital, error) {
	d, u, err := s.acessoDigital(ctx, identity, id)
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	vs, err := s.repo.ListVerificacoes(ctx, s.pool, id, verificacoesMax)
	if err != nil {
		return DetalheDigital{}, MapError(err)
	}
	return DetalheDigital{Digital: d, Verificacoes: vs, Unidade: u.Numero}, nil
}

// Digitais pagina os representantes no escopo; o Arquivo Público vê todos
// (SoFalhas: o painel de integridade e do pipeline; Publicados: os
// liberados à consulta pública, para as exposições).
func (s *Service) Digitais(ctx context.Context, identity auth.Identity, filtro domain.FiltroDigital, p pagination.Params) ([]domain.Digital, int64, error) {
	f := s.filtro(identity, domain.Filtro{})
	filtro.Restrito, filtro.Unidades = f.Restrito, f.Unidades
	if auth.CanGlobal(identity, auth.PermAcervoCentral) {

		filtro.Restrito = false
	}
	items, total, err := s.repo.ListDigitais(ctx, s.pool, filtro, p)
	for i := range items {
		items[i] = items[i].ComEstadoPAdES()
	}
	return items, total, MapError(err)
}

// ErrSemRange: o armazenamento não lê trechos (o visualizador cai para o
// arquivo inteiro).
var ErrSemRange = errors.New("acervo: armazenamento sem leitura por trecho")

// AbrirArquivo confere a leitura e devolve o representante pronto, para o
// visualizador.
func (s *Service) AbrirArquivo(ctx context.Context, identity auth.Identity, id uuid.UUID) (domain.Digital, error) {
	d, _, err := s.acessoDigital(ctx, identity, id)
	if err == nil && d.Situacao != domain.DigitalPronto {
		err = domain.ErrDigitalSituacao
	}
	return d, MapError(err)
}

// LerArquivo devolve o trecho [offset, offset+length) do PDF/A (length 0 =
// o arquivo inteiro).
func (s *Service) LerArquivo(ctx context.Context, d domain.Digital, offset, length int64) (io.ReadCloser, error) {
	if length == 0 {
		return s.dig.Store.Get(ctx, s.dig.WORMBucket, d.PDFAKey)
	}
	rr, ok := s.dig.Store.(storage.RangeReader)
	if !ok {
		return nil, ErrSemRange
	}
	return rr.GetRange(ctx, s.dig.WORMBucket, d.PDFAKey, offset, length)
}
