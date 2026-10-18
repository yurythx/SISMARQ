# Domínio arquivístico do SISMARQ

Linha de base para **implementar ou modificar** acervo, digitalização,
busca, armazenamento, TTDD e procedimentos. Complementa os ADRs (o ADR
ganha se houver conflito pontual; este texto descreve o conjunto).
Travado pelo [ADR 045](../adr/045-linha-de-base-arquivistica.md).

**Não inventar** decreto municipal, portaria de comissão, certificado
ICP-Brasil, cliente SEI, códigos CONARQ 01–07 nem “bateu” na homologação
C1–C4. Ver [FASE0.md](FASE0.md) e [HOMOLOGACAO.md](HOMOLOGACAO.md).

---

## 1. O que o sistema é

O SISMARQ é o Sistema Municipal de Arquivos da Prefeitura de Rondonópolis
(Lei Municipal nº 10.301/2019). Gestão nas **três idades**, TTDD, comissões,
eliminação, recolhimento e consulta pública.

| Peça | Papel | Plug-in |
|---|---|---|
| TTDD e procedimentos | Norma e roteiro | Atlas (`DependsOn` nenhum; porta de TTDD exposta) |
| Unidade de arquivamento | Documento de arquivo | Acervo (`DependsOn: atlas`) |
| Processo na corrente (fora do SEI) | Ingresso quando conclui | Trâmite → evento `tramite.processo.concluido` |
| Assinatura eletrônica | Hash + reautenticação | Signum (`pades` / atas = `AUSENTE`) |
| Pasta de trabalho | ACL, não preservação | Arquivos (MinIO comum) |

O artigo científico em `docs/SISMARQ/` reserva o **SEI** à fase corrente
eletrônica. O Trâmite **não** substitui o SEI ([ADR 042](../adr/042-ingresso-sei-adiado.md)).

Princípios (roadmap §2): a lei manda no código; nenhuma permissão contorna
permanente, bloqueio judicial ou edital de 30 dias; Anexos I–IV iguais aos
da lei; implantação gradual (piloto **SEMAD**).

---

## 2. Unidade de arquivamento

Tabela `acervo_unidades`. Número `UA-AAAA-NNNNNN` gerado pelo banco.
**DELETE é proibido** (art. 38: o registro permanece depois da eliminação).

| Campo | Valores / regra |
|---|---|
| `tipo` | `CAIXA`, `PASTA`, `LIVRO`, `PROCESSO`, `DOSSIE`, `DOCUMENTO`, `MAPA`, `OUTRO` |
| `codigo_classificacao` | Série da TTDD (Atlas), via porta `domain.TTDD` — plug-ins não se importam |
| `unidade_id` | Produtora = escopo IAM |
| `data_inicial` / `data_final` | Datas-limite |
| `fase` | `CORRENTE`, `INTERMEDIARIA`, `PERMANENTE`, `ELIMINADO` |
| `destinacao` | Da série: `ELIMINACAO`, `GUARDA_PERMANENTE`, `INDEFINIDA` (art. 28) |
| `nivel_acesso` | `PUBLICO` (padrão), `RESTRITO`, `SIGILOSO`; hipótese legal se ≠ público (art. 2º) |
| `data_evento` | Condição “enquanto vigorar” |
| `contas_aprovadas_em` | Séries financeiras / Anexo I / TCE |
| `local_id` / `posicao` | Hierarquia física |

Locais (`acervo_locais`): `PREDIO` → `DEPOSITO` → `SALA` → `CORREDOR` →
`ESTANTE` → `PRATELEIRA`. Só `acervo:central` (global) cadastra o depósito
do Arquivo Público.

### 2.1. Situação gravada

Recalculada a cada escrita, no worker `prazos` (diário) e no evento
`atlas.ttdd.carga.aplicada` (fila `sismarq.acervo.ttdd`). A carga
`make ttdd-aplicar` **não** emite o evento; o recálculo diário cobre.

`NO_CORRENTE` · `TRANSFERIR` · `NO_INTERMEDIARIO` · `AVALIAR_ELIMINACAO` ·
`RECOLHER` · `AGUARDA_EVENTO` · `AGUARDA_DESTINACAO` · `BLOQUEADA` ·
`PERMANENTE` · `RECOLHIDA` · `ELIMINADA` · `SERIE_REVOGADA`

### 2.2. Contagem de prazos (premissa — C1 ainda valida)

[ADR 032](../adr/032-acervo-unidades-de-arquivamento.md). Arts. 29 e 30.

- Corrente: a partir da **data-limite final** da unidade.
- Intermediária: a partir do fim da corrente.
- “Enquanto vigorar” / “vigorando”: fim = **data do evento** comprovada
  (contrato, gestão, vínculo). Sem data → `AGUARDA_EVENTO`.
- Meses/dias: do início da fase; “após o evento”: do evento.
- Série sem destinação na TTDD (6 casos “X”) → `INDEFINIDA` /
  `AGUARDA_DESTINACAO`: nada sai do intermediário até a CCPAD definir.
- Série revogada: UA existente permanece; cadastro novo recusado.

A Superintendência confirma amostra em [HOMOLOGACAO.md](HOMOLOGACAO.md) §1
antes de tratar a conta como definitiva.

### 2.3. Travas (domínio + banco)

Ninguém contorna, inclusive administrador. Funções `PodeTransferir`,
`PodeEliminar`, `PodeEliminarArt35`, `PodeDigitalizar`, `PodePublicarDigital`.
Gatilho `acervo_travas` e correlatos.

| Regra | Lei 10.301 | Código |
|---|---|---|
| Permanente não muda de fase nem se elimina, mesmo digitalizado | arts. 20 e 42 | `ErrPermanente` |
| Eliminada imutável; UA não se apaga | art. 38 | DELETE recusado |
| Eliminação só se destinação `ELIMINACAO` (ou rito 35/36) | arts. 19, 35, 36 | `PodeEliminar*` |
| Bloqueio judicial (`acervo_bloqueios`, `liberado_em` nulo) | arts. 32 II, 44 | `ErrBloqueioJudicial` |
| Transferir só corrente com prazo cumprido | art. 29 | `PodeTransferir` |
| Edital 30 dias no DIORONDON-e antes do Termo | art. 37 § 2º | CHECK |
| Amostra sorteada vira permanente | art. 39 | domínio da listagem |
| Digital PRONTO só publica se UA `PUBLICO` | art. 2º | gatilho de publicação |

Escrita atômica: UA + `acervo_movimentacoes` + outbox `acervo.unidade.*` +
auditoria ([ADR 016](../adr/016-camada-de-aplicacao-e-transacoes.md)).

### 2.4. Permissões

Detalhe: [IAM.md](IAM.md). Resumo: `acervo:read` / `manage` / `bloqueio` no
escopo; `acervo:avaliacao` (CSPAD no órgão); `acervo:central` e
`acervo:semad` **globais**. Fora do escopo: leitura 404, ação 403.

---

## 3. Fases 2 a 4 no código (o que já existe)

| Fluxo | ADR | Estado | Não fazer no PR |
|---|---|---|---|
| Empréstimo (1 aberto/UA, ≤ 1 ano) | 033 | Pronto | — |
| Guia de transferência + aceite | 033 | Pronto | Transferir sem `PodeTransferir` |
| Importação planilha simular→aplicar | 033 | Pronto | — |
| Etiqueta / QR | 033 | Pronto | — |
| Ingresso Trâmite | 039 | Pronto | FK entre plug-ins |
| Eliminação Anexos I–IV, ritos 35/36 | 034, 041 | Pronto; atas `AUSENTE` | Inventar composição sem portaria |
| Recolhimento art. 53 | 035 | Pronto | Alterar proveniência recolhida |
| Privados arts. 55–58 | 036 | Pronto | — |
| Normas art. 54 | 036 | Pronto | Mutar redação publicada |
| SEI | 042 | Envelope recusa `ErrSEIAdiado` | Cliente HTTP de laboratório |
| Aprovação TTDD art. 32 §§ 4º–8º | — | **Não implementado** | Fingir rito CSPAD→SEMAD |

---

## 4. Digitalização e armazenamento

Operação: [DIGITALIZACAO.md](DIGITALIZACAO.md). ADRs 037 e 044.

### 4.1. Origem (não misturar)

| Origem | Arquivo | DPI | Decreto 10.278 |
|---|---|---|---|
| `SCAN` (omissão = SCAN) | PDF ou TIFF | 300 texto / 600 planta; sem imagem = `REPROVADO` | Sim (papel) |
| `NATO_DIGITAL` | PDF já eletrônico | 0 | Não; só preservação PDF/A |

Nato-digital **não** é valor de original do Decreto 10.278.

Estados: `ENVIADO` → `PROCESSANDO` → `PRONTO` \| `REPROVADO` \| `FALHOU`
(máx. 3 tentativas; original preservado se falhou).

### 4.2. Pipeline

Upload pré-assinado (`ACERVO_DIGITAL_MAX_MB`, teto 2 GB) → `pdfinfo` /
`pdfimages` / `tiffinfo` → `ocrmypdf --output-type pdfa-2 --deskew --clean
--skip-text -l por` → pikepdf grava XMP `govbr:` **com** `pdfaExtension` →
`pdftotext` em `acervo_digital_textos` → SHA-256 do PDF/A → bucket
`ACERVO_DIGITAL_BUCKET` (padrão `sismarq-acervo-digital`) com **object lock
Compliance** → apaga o original se `PRONTO`.

Hash no XMP = original; hash do PDF/A (muda ao gravar XMP) fica no banco.
Gatilho impede apagar o `PRONTO` ou reduzir retenção.

Retenção: permanente ou destinação indefinida = 100 anos; eliminação = até
1 ano após fim intermediário (mínimo 1 ano). Eliminar o papel **não** apaga
sozinho o PDF/A.

Fixity: worker `integridade` (`ACERVO_FIXITY_INTERVALO`); tabela
`acervo_digital_verificacoes`; alerta a `acervo:central`.

Visualizador: `GET .../arquivo` com HTTP 206, `ETag` = SHA-256; pdf.js em
fatias. Art. 42: digitalizar não muda fase.

### 4.3. Adiado de propósito

- PAdES ICP-Brasil: `pades=AUSENTE`; porta `AssinadorICP` **não composta**
  ([ADR 040](../adr/040-pades-icp-brasil-adiado.md)).
- Pasta SMB/FTP.
- veraPDF (hoje a conformidade PDF/A é a do ocrmypdf/Ghostscript).

Homologação C6: seed `RD-HOM-*` é metadado `REPROVADO` sem objeto no MinIO.
PDF/A real só com envio pela ficha.

---

## 5. Busca, acesso e guarda lógica

[ADR 038](../adr/038-acervo-acesso-busca-e-difusao.md). O Typesense **não**
é consultado do navegador ([ADR 031](../adr/031-sismarq-e-typesense.md)).

| Canal | Índice | Quem vê |
|---|---|---|
| `/arquivo` (anônimo) | `sismarq_acervo_publico_v1` | Só `PUBLICO`, não eliminado; ficha sem local, observações, cadastrador |
| `/acervo/busca` | `sismarq_acervo_v1` | Escopo; faceta `nivel_acesso` |
| Fallback | Postgres `tsvector` | Sem Typesense; OCR; sem facetas |

Facetas públicas: órgão, classe, espécie, fase, destinação, ano.
Worker `indice` + `acervo_indice_pendentes`. SIC do não-ostensivo: série
`2.0.03.02`. Exposições virtuais: só documento já liberado.

Critério C5: unidade restrita (ex. `HOM-UA-08`) **não** entra em `/arquivo`.

Buckets: `MINIO_BUCKET` (upload temporário) ≠ `ACERVO_DIGITAL_BUCKET`
(PDF/A WORM) ≠ bucket da auditoria. Não gravar representante permanente
no plug-in Arquivos.

---

## 6. TTDD (Atlas)

ADRs [017](../adr/017-ttdd-oficial.md), [019](../adr/019-vigencia-da-ttdd-e-versoes.md),
[022](../adr/022-atualizar-ttdd-pela-tela.md).

Hierarquia duplex `N.0.FF.SS.II`: órgão → função → subfunção → série.
Carga oficial: **9 órgãos, 41 funções, 134 subfunções, 1.686 séries**
(`deploy/ttdd/`). Cada fase: **anos XOR condição**. Destinação no artefato:
eliminação 981, guarda permanente 699, indefinida 6.

Revogada (`revogada_em`): sai da árvore, exportação e assistente; abre pelo
código; **não** aceita UA nova nem procedimento novo (422). Histórico
`atlas_ttdd_historico`: `ALTERADA` / `REVOGADA` / `RESTABELECIDA`.

| Caminho | Evento ao Acervo |
|---|---|
| Tela `/atlas/ttdd/atualizar` (simular → aplicar + hash) | `atlas.ttdd.carga.aplicada` |
| `make ttdd-impacto` / `ttdd-aplicar` | Não; worker diário |

**Atividades-meio CONARQ 01–07:** o município **adota** (art. 26). O Atlas
só tem atividades-fim do DIORONDON-e. Carga **somente** pela tela, arquivo
oficial, **sem códigos inventados**. Educação não tem órgão na TTDD.

Divergências do PDF (prefixo 5.0, duplicatas `-2`, seção 13.0 truncada,
correntes em branco, “X” da Saúde etc.) ficam em `deploy/ttdd/ttdd.json`
para a CCPAD — o sistema não “corrige” a norma.

Rito de **aprovação** da TTDD de atividades-fim (art. 32 §§ 4º–8º) **não**
está no produto; só o relatório de impacto da carga.

---

## 7. Procedimentos (Atlas)

[ADR 028](../adr/028-validacao-dos-procedimentos.md). Procedimento **não** é
unidade de arquivamento: é o fluxo (etapas, peças, transições, SLA)
classificado numa série vigente.

`RASCUNHO` → `EM_VALIDACAO` → `HOMOLOGADO`. CHECK: ativo só se homologado.
Importação padrão = rascunho. Homologado não volta; revisão = nova versão
na mesma transação. Entrevistas `atlas_validacoes` imutáveis. Caderno:
`/atlas/secretarias/{prefixo}/caderno`. Organograma = funcional da TTDD,
não o AD ([ADR 029](../adr/029-organograma-da-ttdd.md)).

Assistente ([ADR 023](../adr/023-assistente-fluxos-e-ttdd.md)): fluxos
**homologados** + TTDD (relevância ≥ 0,65). Recusa tarefas (“escreva”).
O ADR 021 (só TTDD) foi ampliado pelo 023.

Débitos operacionais: [PENDENCIAS.md](../atlas/PENDENCIAS.md) (siglas,
modelos vazios na data do registro, compras SEMAD vs SEFIN, Educação).

Rotas: públicas `/procedimentos`, `/temporalidade`; autenticadas `/atlas/*`.

---

## 8. Código de referência

Domínio Acervo: `backend/internal/modules/acervo/domain/` (`acervo.go`,
`prazos.go`, `eliminacao.go`, `digital.go`, `acesso.go`, `sei.go`, …).

Migrações 000146–000156 (núcleo UA até origem SCAN/NATO_DIGITAL).

Workers: `prazos`, `digitalizacao`, `integridade`, `indice`.

Frontend autenticado `/acervo/*`; público `/arquivo/*`. Contrato HTTP da TTDD
para outros sistemas: [TTDD_API.md](TTDD_API.md).

Testes: travas no domínio; `TestPipelineReal` / `make pipeline-real`;
`make homolog-conferir` (metadados, **não** ata do art. 62).

---

## 9. Ordem para o próximo PR ou reunião

1. Homologação **C1** (prazos) e **C2–C4** (Anexos) — Superintendência + SEMAD.
2. **C6** com PDF/TIFF reais na ficha.
3. CONARQ 01–07 com tabela oficial.
4. Checklist e-ARQ / rumo a RDC: [EARQ_RDC.md](EARQ_RDC.md).
5. Protocolo / outros sistemas: contrato [TTDD_API.md](TTDD_API.md) (art. 45).
6. Biblioteca de modelos + entrevistas (PENDENCIAS do Atlas).
7. Rito art. 32 no sistema (ADR novo).
8. veraPDF e pacote AIP (ADR novo; SIGAD já opera).
9. PAdES quando houver certificado; SEI quando houver API.
10. RIPD (DPO), backup medido, capacitação, ato art. 64 — Fase 8.

Invariantes de implementação: plug-ins isolados; lei no domínio e no
banco; cobertura do módulo no CI; não afirmar ICP/SEI/e-ARQ certificado.
