"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia, situacaoDigital } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoDigital } from "@/lib/sismarq/types";

/** Representantes digitais no escopo; o filtro de falhas é o painel do
 * Arquivo Público (integridade e processamento). */
export default function DigitaisPage() {
  const [page, setPage] = useState(1);
  const [falhas, setFalhas] = useState(false);
  const lista = useApiPage<AcervoDigital>(withQuery("v1/acervo/digitais", { falhas: falhas || undefined, page, page_size: 20 }));
  const itens = lista.data?.items ?? [];
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Digitalização"
        description="Representantes digitais com valor probatório (Decreto Federal 10.278/2020): PDF/A com OCR, metadados no XMP, SHA-256 e guarda imutável, conferida periodicamente. Digitalizar não libera a eliminação de documento permanente (arts. 20 e 42)."
      />
      <label className="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          checked={falhas}
          onChange={(e) => {
            setFalhas(e.target.checked);
            setPage(1);
          }}
        />
        Só os que falharam (integridade ou processamento)
      </label>
      <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle={falhas ? "Nenhuma falha" : "Nenhuma digitalização"}>
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Representantes digitais">
          {itens.map((d) => {
            const s = situacaoDigital(d.situacao);
            return (
              <li key={d.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                <div>
                  <Link href={`/acervo/digitais/${d.id}`} className="font-mono font-medium hover:underline">
                    {d.numero}
                  </Link>
                  <p className="text-xs text-muted">
                    {d.codigo_classificacao} · enviado em {fmtDia(d.criado_em)}
                    {d.verificado_em && ` · verificado em ${fmtDia(d.verificado_em)}`}
                  </p>
                </div>
                <div className="flex gap-2">
                  {d.integro === false && <Badge tone="danger">Integridade falhou</Badge>}
                  <Badge tone={s.tone}>{s.label}</Badge>
                </div>
              </li>
            );
          })}
        </ul>
        <Pagination meta={lista.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
