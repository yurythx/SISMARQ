# Digitalização no SISMARQ

Como o documento digital entra no repositório. Base: Decreto Federal
10.278/2020, Lei Municipal 10.301/2019 arts. 20, 42 e 48, ADRs
[037](../adr/037-acervo-digitalizacao.md) e
[044](../adr/044-ingresso-scan-nato-digital.md). O cidadão vê o mesmo
fluxo na página inicial (`/#digitalizacao`).

## Dois caminhos (não misturar)

| Origem | O que é | Arquivo | DPI | Se falhar a medição |
|---|---|---|---|---|
| `SCAN` | Papel digitalizado | PDF ou TIFF | 300 texto / 600 planta | `REPROVADO` — digitalizar de novo |
| `NATO_DIGITAL` | Já nasceu eletrônico | PDF | 0 (não se aplica) | Segue para PDF/A e OCR |

Digitalizar **não** autoriza eliminar permanente (arts. 20 e 42).

Assinatura ICP-Brasil: o JSON declara `pades=AUSENTE` até existir
certificado municipal (ADR 040). Pasta de rede (SMB/FTP): não habilitada
(Fase 0).

## O que o servidor faz na tela

1. Abre a unidade de arquivamento (quem gere o órgão; depois do
   recolhimento, o Arquivo Público).
2. Em **Representantes digitais**, escolhe a origem.
3. Envia o arquivo (até o limite `ACERVO_DIGITAL_MAX_MB`, padrão 2 GB).
   O navegador manda direto ao MinIO.
4. Informa scanner **ou** sistema de origem, tipo (texto/planta), parte
   ou volume, CPF (guardado mascarado).
5. Aguarda o worker: fila → processamento → `PRONTO` (PDF/A) ou
   `REPROVADO` / `FALHOU`.

Quem tem `acervo:central` verifica a integridade e reprocessa falhas.

## O que o worker faz

1. Mede páginas e DPI (`pdfimages` no PDF, `tiffinfo` no TIFF).
2. Se `SCAN` e não há imagens ou o DPI é baixo: reprova e apaga o
   original.
3. `ocrmypdf` → PDF/A-2b, OCR em português, deskew.
4. Metadados do Decreto no XMP (`govbr:`).
5. SHA-256, object lock, apaga o original quando `PRONTO`.
6. Texto do OCR para a busca; verificação periódica do hash.

## Homologação (popular o banco)

Pré-requisito: `make homolog-acervo` (organograma SEMAD + TTDD + caixas
`HOM-UA-*`).

```bash
make homolog-digitais
```

Isso aplica `deploy/homolog/seed-digitais.sql`: marca as caixas e grava
dois representantes **`REPROVADO`** (scan e nato-digital) só para a lista
aparecer. **Não** há PDF/A no MinIO nessas linhas — o pipeline de verdade
é o envio pela tela (situação `ENVIADO` dispara o worker).

Para testar o pipeline: na ficha de `HOM-UA-01` ou `HOM-UA-03`, envie um
PDF do scanner (`SCAN`, 300 DPI) e um PDF de ofício (`NATO_DIGITAL`).

## Testes automatizados

A meta de 100% por plug-in é do CI (`scripts/coverage-gate.sh`, ADRs 009
e 010): `go test ./...` com `coverpkg=./internal/...,./pkg/...`. Este
recorte (origem SCAN / NATO_DIGITAL, TIFF, seed de homologação) está
coberto nos pacotes do acervo e na UI.

Banco de teste **vazio** (não reutilize `template1` com schema antigo):

```bash
# Postgres de CI local + migrations 000001–000156
export TEST_DATABASE_URL="postgres://nexus:...@localhost:5432/nexus_test?sslmode=disable"
cd backend
go test -p 1 -count=1 -coverpkg=./internal/modules/acervo/... \
  -coverprofile=coverage-acervo.out ./internal/modules/acervo/...
go tool cover -func=coverage-acervo.out | grep -E 'digital.go|pipeline.go'
# Só o plug-in acervo: scripts/coverage-gate.sh coverage-acervo.out
# deve marcar acervo 100%. digital.go e pipeline.go ficam em 100%.

cd ../frontend
npx vitest run src/app/\(protected\)/acervo/fase5.test.tsx \
  src/app/faq/__tests__/faqPage.test.tsx \
  src/app/__tests__/public-site.test.tsx
```

`TestSeedHomologAcervo` lê `deploy/homolog/seed-acervo.sql` e
`seed-digitais.sql` a partir da raiz do repositório (cinco níveis acima
de `infrastructure/`). Monte o repo inteiro, não só `backend/`.

O Postgres de teste compartilhado acumula unidades; listas globais não
podem exigir “o item deste teste é o primeiro da página”. O CI sobe um
banco vazio. Localmente: `make test-db-reset` (recria `nexus_test` a
partir de `template0`). Representante digital **não se apaga** (Decreto
10.278): o seed é idempotente, sem `DELETE`.

Pipeline com `ocrmypdf` e `tiffinfo` de verdade: `make pipeline-real`.

## O que não fazer

- Enviar PDF nato-digital como `SCAN` (será reprovado sem imagens).
- Inventar PAdES, pasta SMB ou validação veraPDF neste repositório.
- Tratar digitalização como autorização para descartar permanente.
