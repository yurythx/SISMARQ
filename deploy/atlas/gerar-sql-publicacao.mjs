// Gera deploy/homolog/seed-procedimentos.sql a partir dos JSON de rascunho,
// já como HOMOLOGADO (consulta pública). Idempotente: não altera os três
// procedimentos da migration 000132 se o código já existir.
import { createHash } from "node:crypto";
import { readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const dir = dirname(fileURLToPath(import.meta.url));
const a = JSON.parse(readFileSync(join(dir, "rascunhos-procedimentos.json"), "utf8"));
const b = JSON.parse(readFileSync(join(dir, "rascunhos-procedimentos-2.json"), "utf8"));
const procedimentos = [...a.procedimentos, ...b.procedimentos];

function uuid5(s) {
  const h = createHash("sha1").update("sismarq-atlas-procedimento:" + s).digest();
  h[6] = (h[6] & 0x0f) | 0x50;
  h[8] = (h[8] & 0x3f) | 0x80;
  const hex = h.subarray(0, 16).toString("hex");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20, 32)}`;
}

function q(s) {
  return "'" + String(s ?? "").replaceAll("'", "''") + "'";
}

const lines = [
  "-- GERADO por deploy/atlas/gerar-sql-publicacao.mjs — não edite à mão.",
  "-- Publica os 92 fluxos dos JSON como HOMOLOGADO (consulta /procedimentos).",
  "-- Não sobrescreve código já existente (ADM.LIC.001, ADM.DIR.002, ADM.MAT.003).",
  "BEGIN;",
];

for (const p of procedimentos) {
  const wf = uuid5(p.codigo_processual);
  const etapas = p.etapas ?? [];
  const etapaID = Object.fromEntries(etapas.map((e) => [e.ordem, uuid5(p.codigo_processual + ":e:" + e.ordem)]));
  lines.push(`
INSERT INTO atlas_workflows (id, codigo_processual, titulo, objetivo, publico_alvo, versao, ativo, situacao, nivel_acesso, hipotese_legal_restricao, codigo_ttdd)
SELECT ${q(wf)}, ${q(p.codigo_processual)}, ${q(p.titulo)}, ${q(p.objetivo)}, ${q(p.publico_alvo)}, 1, TRUE, 'HOMOLOGADO', ${q(p.nivel_acesso || "PUBLICO")}, NULLIF(${q(p.hipotese_legal_restricao || "")}, ''), ${q(p.codigo_ttdd)}
WHERE EXISTS (SELECT 1 FROM atlas_classificacao_ttdd c WHERE c.codigo = ${q(p.codigo_ttdd)} AND c.revogada_em IS NULL)
  AND NOT EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.codigo_processual = ${q(p.codigo_processual)});
`);
  for (const e of etapas) {
    const eid = etapaID[e.ordem];
    lines.push(`INSERT INTO atlas_etapas (id, workflow_id, ordem, unidade_administrativa, nome_setor, atribuicoes_setor, prazo_sla_em_dias, manter_aberto_apos_remessa)
SELECT ${q(eid)}, ${q(wf)}, ${Number(e.ordem)}, ${q(e.unidade_administrativa)}, ${q(e.nome_setor)}, ${q(e.atribuicoes_setor)}, ${Number(e.prazo_sla_em_dias || 0)}, ${e.manter_aberto_apos_remessa ? "TRUE" : "FALSE"}
WHERE EXISTS (SELECT 1 FROM atlas_workflows w WHERE w.id = ${q(wf)})
  AND NOT EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.workflow_id = ${q(wf)} AND x.ordem = ${Number(e.ordem)});
`);
    for (const d of e.documentos ?? []) {
      const did = uuid5(p.codigo_processual + ":d:" + e.ordem + ":" + d.nome_documento);
      lines.push(`INSERT INTO atlas_etapa_documentos (id, etapa_id, nome_documento, obrigatorio, formato, tipo_assinatura, exige_conferencia_copia, modelo_minuta_padrao_url)
SELECT ${q(did)}, ${q(eid)}, ${q(d.nome_documento)}, ${d.obrigatorio ? "TRUE" : "FALSE"}, ${q(d.formato)}, ${q(d.tipo_assinatura)}, ${d.exige_conferencia_copia ? "TRUE" : "FALSE"}, NULLIF(${q(d.modelo_minuta_padrao_url || "")}, '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas x WHERE x.id = ${q(eid)})
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_documentos y WHERE y.etapa_id = ${q(eid)} AND y.nome_documento = ${q(d.nome_documento)});
`);
    }
  }
  for (const e of etapas) {
    const origem = etapaID[e.ordem];
    for (const t of e.transicoes ?? []) {
      const dest = etapaID[t.destino_ordem];
      if (!dest) continue;
      const tid = uuid5(p.codigo_processual + ":t:" + e.ordem + ":" + t.destino_ordem + ":" + t.condicao_transicao);
      lines.push(`INSERT INTO atlas_etapa_transicoes (id, origem_etapa_id, destino_etapa_id, condicao_transicao, is_devolucao_diligencia, descricao_diligencia)
SELECT ${q(tid)}, ${q(origem)}, ${q(dest)}, ${q(t.condicao_transicao)}, ${t.is_devolucao_diligencia ? "TRUE" : "FALSE"}, NULLIF(${q(t.descricao_diligencia || "")}, '')
WHERE EXISTS (SELECT 1 FROM atlas_etapas o WHERE o.id = ${q(origem)})
  AND EXISTS (SELECT 1 FROM atlas_etapas d WHERE d.id = ${q(dest)})
  AND NOT EXISTS (SELECT 1 FROM atlas_etapa_transicoes z WHERE z.origem_etapa_id = ${q(origem)} AND z.destino_etapa_id = ${q(dest)} AND z.condicao_transicao = ${q(t.condicao_transicao)});
`);
    }
  }
}

const codes = procedimentos.map((p) => q(p.codigo_processual)).join(", ");
const values = procedimentos.map((p) => `(${q(p.codigo_processual)})`).join(", ");
lines.push(`
UPDATE atlas_workflows
   SET situacao = 'HOMOLOGADO', ativo = TRUE
 WHERE codigo_processual IN (${codes})
   AND situacao <> 'HOMOLOGADO';

DO $chk$
DECLARE
  faltam TEXT;
BEGIN
  SELECT string_agg(c.codigo, ', ' ORDER BY c.codigo) INTO faltam
    FROM (VALUES ${values}) AS c(codigo)
   WHERE NOT EXISTS (
     SELECT 1 FROM atlas_workflows w
      WHERE w.codigo_processual = c.codigo AND w.ativo AND w.situacao = 'HOMOLOGADO'
   );
  IF faltam IS NOT NULL THEN
    RAISE EXCEPTION 'procedimentos não publicados (série TTDD ausente ou inserção falhou): %', faltam;
  END IF;
END
$chk$;
COMMIT;
`);

const dest = join(dir, "..", "homolog", "seed-procedimentos.sql");
writeFileSync(dest, lines.join("\n"));
console.log("escreveu", dest, "procedimentos", procedimentos.length);
