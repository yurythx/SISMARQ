# Homologação do Acervo com a Superintendência de Arquivo Público

Para: Superintendência de Arquivo Público e a equipe que desenvolve o
SISMARQ. Órgão piloto: **SEMAD**. Decisões da Fase 0:
[FASE0.md](FASE0.md) e [ADR 043](../adr/043-fase0-decisoes-institucionais.md).

O código das Fases 1 a 6 está no repositório (ADRs 032 a 044). Mapa de
domínio para PRs: [DOMINIO.md](DOMINIO.md) e [EARQ_RDC.md](EARQ_RDC.md)
([ADR 045](../adr/045-linha-de-base-arquivistica.md)). Esta página é o que
**não se fecha no git**: a conferência com a lei, a TTDD e os Anexos. Sem
ela o sistema oficial não está pronto (art. 62).

Na VM, o compose `migrate` precisa estar em **goose ≥ 157** (origem SCAN/NATO_DIGITAL
e remoção do catálogo/blog PaperMoon) **antes** de `make homolog-acervo` /
`bash deploy/homolog/aplicar-na-vm.sh`. Conferência automática de metadados:
`make homolog-conferir`. C1–C4 no papel e C6 com PDF real no MinIO continuam humanos.

## 1. Prazos (Fase 1) — C1

Amostrar 5 a 10 códigos reais da TTDD **em vigor** e conferir o painel.
A tabela pode ser republicada: vale a versão vigente na época do
documento; série revogada não apaga a unidade (ADR 019).

| Conferir | Como (prática adotada) |
|---|---|
| Prazo em anos | Data-limite final + anos da série = fim da fase (arts. 29 e 30; ADR 032). |
| “Enquanto vigorar” / “vigorando” | Data do **evento** comprovada: fim de contrato/convênio, fim de gestão/mandato, desligamento, ou o fato que a condição da série descreve. Quem cadastra informa; o sistema não inventa. |
| Contas do TCE | Séries financeiras: campo “contas aprovadas em” do Anexo I. Sem essa data a eliminação financeira não se confunde com o fim da fase por anos. |
| Série revogada | Unidade permanece com situação de série fora da TTDD vigente; não some. |

Registrar o resultado (data, servidor, códigos conferidos) na seção 5.
Enquanto isso, o sistema conta o que está no ADR 032.

## 2. Anexos I a IV (Fase 3) — C2 a C4

Eliminação piloto, ponta a ponta, na **SEMAD**. Imprimir Relação, Ata,
Edital e Termo e confrontar com os modelos da lei. Nome do órgão nos
PDFs: Superintendência de Arquivo Público. Veículo do edital: DIORONDON-e
(ou cópia de teste da mesma redação).

- **C2 (ordinário):** Relação → ata CSPAD → Autorizo do Arquivo → edital
  30 dias → amostragem → Termo.
- **C3 (art. 35):** unidade fora da TTDD ou destinação indefinida.
  Conferir CCPAD → Arquivo → SEMAD → edital → termo, e a cópia na caixa
  de `acervo:central` e `acervo:semad`.
- **C4 (art. 36):** só se o órgão **não** estiver na lista 0.5 (TTDD no
  Atlas). Comissão = CCPAD; `rito=ART_36`.

Ata/termo continuam `AUSENTE` no Signum e no PAdES até haver certificado
(ADR 040). Ainda **não** entra a aprovação da TTDD de atividades-fim
(art. 32, §§ 4º a 8º).

## 3. Acesso público (Fase 6) — C5

| Conferir | Esperado |
|---|---|
| Documento permanente digitalizado e ostensivo | Aparece em `/arquivo` e na busca sem login |
| Unidade restrita ou sigilosa | Não aparece na consulta pública; o SIC (`2.0.03.02`) é o caminho |
| Mudança de nível de acesso | O PDF some da consulta na hora |
| Facetas | Órgão, classe, espécie, fase, destinação, ano |

## 3.1. Digitalização (Fase 5) — C6

Dois caminhos na ficha da unidade ([DIGITALIZACAO.md](DIGITALIZACAO.md),
ADR 044). Popular metadados: `make homolog-digitais` (ou
`make homolog-acervo`, que já inclui). O PDF/A real só existe após o
envio pela tela.

| Conferir | Esperado |
|---|---|
| PDF do scanner, origem SCAN, 300 DPI texto | `PRONTO`, PDF/A-2b, OCR |
| PDF sem imagens, origem SCAN | `REPROVADO` (não é scan) |
| PDF sem imagens, origem NATO_DIGITAL | `PRONTO` (não exige DPI) |
| TIFF multipágina, origem SCAN, DPI ok | `PRONTO` |
| Permanente digitalizado | unidade continua permanente (art. 42) |

Conferência **automática** dos metadados do seed (não é a ata do art. 62):

```bash
make homolog-conferir
```

C2–C4 (Anexos impressos) e C6 com arquivo real no MinIO continuam
humanos. O seed `RD-HOM-SCAN` / `RD-HOM-NATO` é conferido também por
`TestSeedHomologAcervo` (origem e DPI 0 no nato-digital); isso não
substitui o envio pela tela. Como rodar a suíte:
[DIGITALIZACAO.md](DIGITALIZACAO.md) (seção Testes automatizados).
IAM: [IAM.md](IAM.md). CONARQ 01–07: importar em Atualizar TTDD
quando a tabela oficial estiver em mãos (o Atlas já avisa o art. 26);
roteiro: [atlas/ATUALIZAR.md](../atlas/ATUALIZAR.md).
RIPD: [RIPD_MODELO.md](../RIPD_MODELO.md) — o DPO preenche; não inventar.

## 4. O que a Fase 0 já fechou

Quadro em [FASE0.md](FASE0.md). IAM: mapear grupos LDAP para os perfis
com `acervo:central` (Arquivo, global), `acervo:semad` (SEMAD, global) e
`acervo:avaliacao` (CSPAD, no órgão). Portarias, certificado e SEI
continuam posteriores.

## 5. Registro

| Data | Quem | O quê | Resultado |
|---|---|---|---|
| | | | |

Quando a primeira linha estiver preenchida com “bateu com a TTDD / com os
Anexos”, atualize a coluna “Hoje” do roadmap no mesmo espírito das outras
fases.
