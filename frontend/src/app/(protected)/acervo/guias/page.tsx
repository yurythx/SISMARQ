"use client";

import { ArrowLeft, Plus } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia, situacaoGuia } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { buttonClass } from "@/components/ui/Button";
import { useApiPage, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoGuia } from "@/lib/sismarq/types";

/** Guias de transferência ao arquivo intermediário (art. 29; a guia é
 * documento de arquivo — TTDD 2.0.03.01.00). */
export default function GuiasPage() {
  const { can } = useSismarq();
  const [page, setPage] = useState(1);
  const guias = useApiPage<AcervoGuia>(withQuery("v1/acervo/guias", { page, page_size: 20 }));
  const itens = guias.data?.items ?? [];
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Guias de transferência"
        description="A remessa das unidades com o prazo corrente cumprido ao arquivo intermediário do órgão, com o aceite de quem recebe (Lei 10.301/2019, art. 29)."
        actions={
          can("acervo:manage") && (
            <Link href="/acervo/guias/nova" className={buttonClass("primary", "md")}>
              <Plus size={16} aria-hidden="true" /> Nova guia
            </Link>
          )
        }
      />
      <DataState loading={guias.isLoading} error={guias.error} empty={itens.length === 0} emptyTitle="Nenhuma guia de transferência">
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Guias">
          {itens.map((g) => {
            const s = situacaoGuia(g.situacao);
            return (
              <li key={g.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                <div>
                  <Link href={`/acervo/guias/${g.id}`} className="font-mono font-medium hover:underline">
                    {g.numero}
                  </Link>
                  <p className="text-xs text-muted">
                    {g.total_itens} unidade(s) · emitida em {fmtDia(g.criada_em)}
                    {g.recebida_em && ` · recebida em ${fmtDia(g.recebida_em)}`}
                  </p>
                </div>
                <Badge tone={s.tone}>{s.label}</Badge>
              </li>
            );
          })}
        </ul>
        <Pagination meta={guias.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
