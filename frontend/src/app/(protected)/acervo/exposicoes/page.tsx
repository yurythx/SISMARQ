"use client";

import { ArrowLeft, Plus } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { buttonClass } from "@/components/ui/Button";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoExposicao } from "@/lib/sismarq/types";

/** Exposições virtuais do Arquivo Público (difusão, art. 5º, X e XI). */
export default function ExposicoesPage() {
  const [page, setPage] = useState(1);
  const lista = useApiPage<AcervoExposicao>(withQuery("v1/acervo/exposicoes", { page, page_size: 20 }));
  const itens = lista.data?.items ?? [];
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Arquivo Público"
        title="Exposições virtuais"
        description="Documentos liberados à consulta pública reunidos em exposições, com texto de apresentação e legendas. Só os documentos liberados aparecem na visita."
        actions={
          <Link href="/acervo/exposicoes/nova" className={buttonClass("primary", "md")}>
            <Plus size={16} aria-hidden="true" /> Nova exposição
          </Link>
        }
      />
      <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle="Nenhuma exposição ainda">
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Exposições">
          {itens.map((x) => (
            <li key={x.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
              <div>
                <Link href={`/acervo/exposicoes/${x.id}`} className="font-medium hover:underline">
                  {x.titulo}
                </Link>
                <p className="text-xs text-muted">
                  /arquivo/exposicoes/{x.slug} · alterada em {fmtDia(x.atualizado_em)}
                </p>
              </div>
              {x.publicada_em ? <Badge tone="success">Publicada em {fmtDia(x.publicada_em)}</Badge> : <Badge>Rascunho</Badge>}
            </li>
          ))}
        </ul>
        <Pagination meta={lista.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
