"use client";

import { ArrowLeft, Plus } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia, situacaoRecolhimento } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { buttonClass } from "@/components/ui/Button";
import { useApiPage, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoRecolhimento } from "@/lib/sismarq/types";

/** Recolhimentos ao Arquivo Público (Lei 10.301/2019, arts. 5º, VII, e 53). */
export default function RecolhimentosPage() {
  const { can } = useSismarq();
  const [page, setPage] = useState(1);
  const lista = useApiPage<AcervoRecolhimento>(withQuery("v1/acervo/recolhimentos", { page, page_size: 20 }));
  const itens = lista.data?.items ?? [];
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Recolhimentos ao Arquivo Público"
        description="O conjunto de guarda permanente vai ao Arquivo Público classificado, avaliado, organizado, higienizado, acondicionado e com instrumento descritivo (art. 53). Na entrada, o Arquivo registra o número de entrada e a origem recebe a cópia do registro (art. 5º, VII)."
        actions={
          can("acervo:manage") && (
            <Link href="/acervo/recolhimentos/novo" className={buttonClass("primary", "md")}>
              <Plus size={16} aria-hidden="true" /> Novo recolhimento
            </Link>
          )
        }
      />
      <DataState
        loading={lista.isLoading}
        error={lista.error}
        empty={itens.length === 0}
        emptyTitle="Nenhum recolhimento"
        emptyDescription="As unidades com os prazos cumpridos e destinação guarda permanente aparecem no painel do Acervo como “Recolher ao Arquivo Público”."
      >
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Recolhimentos">
          {itens.map((r) => {
            const s = situacaoRecolhimento(r.situacao);
            return (
              <li key={r.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                <div>
                  <Link href={`/acervo/recolhimentos/${r.id}`} className="font-mono font-medium hover:underline">
                    {r.numero}
                  </Link>
                  <p className="text-xs text-muted">
                    {r.total_itens} unidade(s) · enviado em {fmtDia(r.criado_em)}
                    {r.registro_entrada && ` · registro de entrada ${r.registro_entrada}`}
                  </p>
                </div>
                <Badge tone={s.tone}>{s.label}</Badge>
              </li>
            );
          })}
        </ul>
        <Pagination meta={lista.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
