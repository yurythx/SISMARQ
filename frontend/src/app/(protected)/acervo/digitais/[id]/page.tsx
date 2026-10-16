"use client";

import { ArrowLeft, Download } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";

import { fmtBytes, fmtDia, situacaoDigital } from "@/components/acervo/labels";
import { VisualizadorPDF } from "@/components/acervo/VisualizadorPDF";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { fmtDateTime, useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button, buttonClass } from "@/components/ui/Button";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoDetalheDigital } from "@/lib/sismarq/types";

function Campo({ nome, children }: { nome: string; children: React.ReactNode }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-muted">{nome}</dt>
      <dd className="mt-0.5 break-all text-sm">{children || "—"}</dd>
    </div>
  );
}

/** O representante digital: metadados do Decreto 10.278/2020, hashes,
 * integridade e o documento em fluxo. */
export default function DigitalPage() {
  const { id } = useParams<{ id: string }>();
  const { canGlobal } = useSismarq();
  const central = canGlobal("acervo:central");
  const q = useApiQuery<AcervoDetalheDigital>(`v1/acervo/digitais/${id}`);
  const { run, pending } = useAction();
  const d = q.data;

  async function acao(qual: "verificar" | "reprocessar") {
    const r = await run(
      () => apiClient.post<AcervoDetalheDigital>(`v1/acervo/digitais/${id}/${qual}`),
      qual === "verificar" ? "Integridade verificada." : "Devolvido ao processamento.",
    );
    if (r) await q.mutate(r.data, { revalidate: false });
  }

  const s = d ? situacaoDigital(d.situacao) : null;
  const arquivo = `/api/backend/v1/acervo/digitais/${id}/arquivo`;
  return (
    <div className="flex flex-col gap-6">
      <DataState loading={q.isLoading} error={q.error} empty={!d}>
        {d && s && (
          <>
            <Link href={`/acervo/${d.unidade_arquivamento_id}`} className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
              <ArrowLeft size={14} aria-hidden="true" /> Unidade {d.unidade_numero}
            </Link>
            <PageHeader
              eyebrow="Representante digital · Decreto 10.278/2020"
              title={d.numero}
              description={d.descricao || `Digitalização da unidade ${d.unidade_numero}`}
              actions={
                <>
                  {d.situacao === "PRONTO" && (
                    <a href={arquivo} download={`${d.numero}.pdf`} className={buttonClass("secondary", "md")}>
                      <Download size={16} aria-hidden="true" /> Baixar PDF/A
                    </a>
                  )}
                  {central && d.situacao === "PRONTO" && (
                    <Button variant="secondary" disabled={pending} onClick={() => void acao("verificar")}>
                      Verificar integridade
                    </Button>
                  )}
                  {central && d.situacao === "FALHOU" && (
                    <Button disabled={pending} onClick={() => void acao("reprocessar")}>
                      Reprocessar
                    </Button>
                  )}
                </>
              }
            />
            <div className="flex flex-wrap items-center gap-2">
              <Badge tone={s.tone}>{s.label}</Badge>
              {d.integro === true && <Badge tone="success">Íntegro na última verificação</Badge>}
              {d.integro === false && <Badge tone="danger">Integridade falhou — o hash não confere</Badge>}
            </div>
            {d.motivo && <p className="text-sm">{d.motivo}</p>}

            <dl className="grid gap-3 rounded-xl border border-surface-border bg-surface p-4 sm:grid-cols-3">
              <Campo nome="Identificador único">{d.numero}</Campo>
              <Campo nome="Classificação (TTDD)">{d.codigo_classificacao}</Campo>
              <Campo nome="Digitalizado em (UTC)">{fmtDateTime(d.digitalizado_em)}</Campo>
              <Campo nome="Responsável">{[d.responsavel_nome, d.responsavel_cpf].filter(Boolean).join(" · ")}</Campo>
              <Campo nome="Equipamento">{d.equipamento}</Campo>
              <Campo nome="Resolução">
                {d.dpi_medido ? `${d.dpi_medido} DPI medidos` : "—"} (declarado {d.dpi_declarado} DPI, {d.tipo_documento.toLowerCase()})
              </Campo>
              <Campo nome="Formato">{d.conformidade && `${d.conformidade} · ${d.paginas} página(s) · ${fmtBytes(d.pdfa_bytes)}`}</Campo>
              <Campo nome="Guarda imutável até">{fmtDia(d.retencao_ate)}</Campo>
              <Campo nome="Tentativas">{String(d.tentativas)}</Campo>
              <Campo nome="SHA-256 do PDF/A">
                <span className="font-mono text-xs">{d.pdfa_sha256}</span>
              </Campo>
              <Campo nome="SHA-256 do original">
                <span className="font-mono text-xs">{d.original_sha256}</span>
              </Campo>
            </dl>

            {d.situacao === "PRONTO" && <VisualizadorPDF url={arquivo} titulo={d.numero} />}

            <section aria-labelledby="verificacoes" className="flex flex-col gap-2">
              <h2 id="verificacoes" className="text-base font-semibold">
                Verificações de integridade
              </h2>
              {d.verificacoes.length === 0 ? (
                <p className="text-sm text-muted">Ainda não verificado.</p>
              ) : (
                <ol className="flex flex-col gap-1 text-sm">
                  {d.verificacoes.map((v) => (
                    <li key={v.em}>
                      <span className="text-xs text-muted">{fmtDateTime(v.em)}</span> —{" "}
                      {v.integro ? "íntegro" : v.erro ? `falhou: ${v.erro}` : `hash diferente (${v.sha256_lido.slice(0, 16)}…)`}
                    </li>
                  ))}
                </ol>
              )}
            </section>
          </>
        )}
      </DataState>
    </div>
  );
}
