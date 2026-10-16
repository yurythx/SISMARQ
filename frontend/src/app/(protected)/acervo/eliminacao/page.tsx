"use client";

import { ArrowLeft, Plus, Users } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia, situacaoListagem } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { buttonClass } from "@/components/ui/Button";
import { useApiPage, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoListagem } from "@/lib/sismarq/types";

/** Relações de eliminação (Lei 10.301/2019, arts. 31 a 40): da seleção pela
 * CSPAD ao Termo de Eliminação. */
export default function EliminacaoPage() {
  const { can } = useSismarq();
  const [page, setPage] = useState(1);
  const lista = useApiPage<AcervoListagem>(withQuery("v1/acervo/listagens", { page, page_size: 20 }));
  const itens = lista.data?.items ?? [];
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Eliminação de documentos"
        description="Relação aprovada pela Comissão Setorial (CSPAD), autorizada pelo Arquivo Público, Edital de Ciência com 30 dias no Diário Oficial do Município, amostragem e Termo de Eliminação (Lei 10.301/2019, arts. 31 a 40; Anexos I a IV)."
        actions={
          <>
            <Link href="/acervo/comissoes" className={buttonClass("secondary", "md")}>
              <Users size={16} aria-hidden="true" /> Comissões
            </Link>
            {can("acervo:avaliacao") && (
              <Link href="/acervo/eliminacao/nova" className={buttonClass("primary", "md")}>
                <Plus size={16} aria-hidden="true" /> Nova relação
              </Link>
            )}
          </>
        }
      />
      <DataState
        loading={lista.isLoading}
        error={lista.error}
        empty={itens.length === 0}
        emptyTitle="Nenhuma relação de eliminação"
        emptyDescription="As unidades com os prazos cumpridos e destinação eliminação aparecem no painel do Acervo como “Avaliar eliminação”."
      >
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Relações de eliminação">
          {itens.map((l) => {
            const s = situacaoListagem(l.situacao);
            return (
              <li key={l.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                <div>
                  <Link href={`/acervo/eliminacao/${l.id}`} className="font-mono font-medium hover:underline">
                    {l.numero}
                  </Link>
                  <p className="text-xs text-muted">
                    {l.total_itens} unidade(s) · aberta em {fmtDia(l.criada_em)}
                    {l.edital_publicado_em && ` · edital em ${fmtDia(l.edital_publicado_em)}`}
                    {l.termo_eliminado_em && ` · eliminada em ${fmtDia(l.termo_eliminado_em)}`}
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
