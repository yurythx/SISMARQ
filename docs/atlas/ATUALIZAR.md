# Como atualizar procedimentos e a TTDD

Para: gestão do Atlas (`atlas:manage`) e quem opera o servidor de homologação.

Procedimento (fluxo) e TTDD (prazos da tabela) são cargas diferentes. Um
não substitui o outro. Entrevistas e padronização das peças:
[PLANO_DE_VALIDACAO.md](PLANO_DE_VALIDACAO.md). Regras: ADRs
[019](../adr/019-vigencia-da-ttdd-e-versoes.md),
[022](../adr/022-atualizar-ttdd-pela-tela.md),
[025](../adr/025-atlas-em-escala.md) e
[028](../adr/028-validacao-dos-procedimentos.md).

Na homologação (VM 220, `http://192.168.1.52:3010`) os 95 fluxos (3 da
migration + 92 da carga) já estão **homologados** na consulta pública
`/procedimentos`. Continuam hipóteses até a entrevista; correção entra
como **nova versão**.

## 1. Procedimentos já publicados

O homologado **não volta a rascunho**. Enquanto a nova versão não for
homologada, o site, a busca e o assistente continuam na versão em vigor.

### Um fluxo (prazo, setor, peça, classificação)

1. Atlas → abra o procedimento → **Nova versão** (já vem como rascunho).
2. Ajuste etapas e documentos. Se houver entrevista, **Registrar
   entrevista** (a primeira muda a situação para em validação).
3. **Homologar**. Exige série da TTDD **vigente**. Publica, desativa as
   versões anteriores do mesmo código e avisa quem segue o fluxo e as
   unidades das etapas.

### Vários de uma vez

1. Atlas → **Importar** → baixe os procedimentos em vigor.
2. Edite o JSON e importe de novo **como rascunho** (já vem marcado).
   Simule, confira e aplique. Rascunho não publica e não tira a versão
   em vigor.
3. Homologue cada um (ou o lote) depois da conferência.

Procedimento novo revelado na entrevista: Atlas → **Novo procedimento**
(também vem como rascunho) ou inclua no arquivo de importação.

### O que a sync da VM não faz

`deploy/homolog/seed-procedimentos.sql` (gerado por
`deploy/atlas/gerar-sql-publicacao.mjs` a partir dos JSON em
`deploy/atlas/`) só **insere código que ainda não existe**. Não
sobrescreve `ADM.LIC.001`, `ADM.DIR.002`, `ADM.MAT.003` nem os 92 já
gravados. `powershell -File deploy/homolog/sincronizar-vm.ps1` (ou o
`.sh`) **não atualiza o conteúdo** de um fluxo já publicado.

Use a tela (nova versão / importação) para mudar o que já está no banco.
Use o gerador + seed só para **códigos novos** que ainda não existem.

Para regenerar os JSON de hipótese (antes de um código existir no banco):
`node deploy/atlas/gerar-rascunhos.mjs` e
`node deploy/atlas/gerar-rascunhos-2.mjs`; depois
`node deploy/atlas/gerar-sql-publicacao.mjs` e a sync.

## 2. TTDD

Série que sai da tabela **não é apagada**: fica **revogada** (data e
edição do Diário Oficial), some da consulta, da árvore, da exportação e
do assistente, mas abre pelo código com o histórico dos prazos. Série
revogada não aceita procedimento novo nem nova versão (422).

### Pela tela (gestão)

Atlas → **Tabela de Temporalidade** → **Atualizar TTDD**
(`/atlas/ttdd/atualizar`).

1. Envie o **CSV da própria exportação** (revisado no Excel, mesmo
   cabeçalho) **ou** o `deploy/ttdd/ttdd.json`.
2. **Simule** — totais por situação, antes/depois das alteradas,
   revogadas, restabelecidas, novas e a lista de **procedimentos
   afetados**. Nada é gravado.
3. Confira o arquivo (o aplicar exige o **hash** da simulação) e
   **aplique**.

Atividades-meio CONARQ 01–07: só por esta tela, com o arquivo oficial,
sem inventar código municipal.

### Nova publicação no Diário Oficial (PDF)

1. Troque `docs/ttdd.pdf`.
2. `cd scripts/ttdd && npm ci && npm run tudo` — regenera `deploy/ttdd/`
   (detalhe em [scripts/ttdd/README.md](../../scripts/ttdd/README.md)).
3. Revise `divergencias` em `deploy/ttdd/ttdd.json` (levam à CCPAD; o
   sistema não “corrige” a norma).
4. Aplique pela tela (envie o JSON) **ou** no servidor:

```bash
make ttdd-impacto   # relatório, não grava
make ttdd-aplicar   # o mesmo relatório e grava
```

Na VM 220, `sincronizar-vm` já recarrega `deploy/ttdd/ttdd.sql` +
`carga.sql`. Se as 1.686 séries não mudaram, o relatório fica
`INALTERADA`.

### Depois da carga

Abra os procedimentos listados como afetados. Se prazos ou destinação
mudam o fluxo, publique **nova versão** (seção 1). Homologar continua
exigindo série vigente.

## 3. Onde está cada coisa

| O quê | Onde |
|---|---|
| Consulta pública | `/procedimentos`, `/temporalidade` |
| Gestão dos fluxos | `/atlas`, `/atlas/procedimentos/{id}` |
| Importação em lote | `/atlas/importar` |
| Atualizar a TTDD | `/atlas/ttdd/atualizar` |
| JSON de hipótese | `deploy/atlas/rascunhos-procedimentos.json` e `-2.json` |
| Seed de homologação (só código novo) | `deploy/homolog/seed-procedimentos.sql` |
| TTDD extraída do PDF | `deploy/ttdd/` |
| Sync da VM sem `git pull` | `deploy/homolog/sincronizar-vm.ps1` |
