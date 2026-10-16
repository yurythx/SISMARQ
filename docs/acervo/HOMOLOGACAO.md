# Homologação do Acervo com a Divisão de Arquivo Público

Para: Superintendência / Divisão de Arquivo Público e a equipe que
desenvolve o SISMARQ.

O código das Fases 1 a 6 está no repositório (ADRs 032 a 042). Esta
página é o que **não se fecha no git**: a conferência com a lei, a TTDD
e os Anexos. Sem ela o sistema oficial não está pronto (art. 62).

## 1. Prazos (Fase 1)

Amostrar séries da TTDD em vigor e conferir o painel do Acervo.

| Conferir | Como |
|---|---|
| Prazo em anos | A data-limite final + anos da série = fim da fase (arts. 29 e 30) |
| “Enquanto vigorar” / “vigorando” | Qual **evento** cadastrar (fim do contrato, da vigência, da gestão) |
| Contas do TCE | Séries financeiras: o campo “contas aprovadas em” do Anexo I |
| Série revogada | Unidade fica com situação de série fora da TTDD; não some |

Registrar o resultado (data, servidor, códigos conferidos) neste arquivo
ou numa ata. Enquanto isso, o sistema conta o que está no ADR 032.

## 2. Anexos I a IV (Fase 3)

Uma eliminação piloto, ponta a ponta, na secretaria combinada (sugestão:
SEMAD). Imprimir Relação, Ata, Edital e Termo e confrontar com os
modelos da lei.

Incluir um piloto **art. 35** (série fora da TTDD) e, se a 0.5 listar
órgão sem TTDD, um **art. 36**. Conferir a cópia do Termo na caixa de quem
tem `acervo:central` e `acervo:semad`. A ata/termo continuam sem Signum
(`ata_signum=AUSENTE`) até a 0.6.

Ainda **não** entra: aprovação da TTDD de atividades-fim pelo fluxo da
CSPAD/CCPAD (art. 32, §§ 4º a 8º).

## 3. Acesso público (Fase 6)

| Conferir | Esperado |
|---|---|
| Documento permanente digitalizado e ostensivo | Aparece em `/arquivo` e na busca sem login |
| Unidade restrita ou sigilosa | Não aparece na consulta pública; o SIC (`2.0.03.02`) é o caminho |
| Mudança de nível de acesso | O PDF some da consulta na hora |
| Facetas | Órgão, classe, espécie, fase, destinação, ano |

## 4. O que a Divisão decide (Fase 0)

Quadro em [FASE0.md](FASE0.md). Nenhum decreto é inventado no git. Em
especial: portarias das comissões, lista 0.5 (art. 36), quem recebe
`acervo:semad`, política 0.6 (PAdES/Signum) e 0.7 (SEI).

## 5. Registro

| Data | Quem | O quê | Resultado |
|---|---|---|---|
| | | | |

Quando a primeira linha estiver preenchida com “bateu com a TTDD / com os
Anexos”, atualize a coluna “Hoje” do roadmap no mesmo espírito das outras
fases.
