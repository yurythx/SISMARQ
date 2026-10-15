"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { fmtDia } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { Toggle } from "@/components/ui/Toggle";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoEmprestimo } from "@/lib/sismarq/types";

/** Empréstimos em aberto (art. 12, II), do prazo mais antigo ao mais novo. */
export default function EmprestimosPage() {
  const [atrasados, setAtrasados] = useState(false);
  const [page, setPage] = useState(1);
  const lista = useApiPage<AcervoEmprestimo>(withQuery("v1/acervo/emprestimos", { atrasados: atrasados || undefined, page, page_size: 20 }));
  const itens = lista.data?.items ?? [];
  const hoje = new Date().toISOString().slice(0, 10);
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Empréstimos em aberto"
        description="Unidades de arquivamento fora do arquivo, com quem retirou e até quando (Lei 10.301/2019, art. 12, II)."
        actions={
          <Toggle
            checked={atrasados}
            onChange={(v) => {
              setAtrasados(v);
              setPage(1);
            }}
            label="Só os atrasados"
          />
        }
      />
      <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle="Nenhum empréstimo em aberto">
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Empréstimos">
          {itens.map((e) => (
            <li key={e.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
              <div>
                <Link href={`/acervo/${e.unidade_arquivamento_id}`} className="font-medium hover:underline">
                  <span className="font-mono text-xs text-muted">{e.numero}</span> {e.assunto}
                </Link>
                <p className="text-xs text-muted">
                  {e.solicitante || "—"} · {e.finalidade} · devolver até {fmtDia(e.devolver_ate)}
                </p>
              </div>
              {e.devolver_ate.slice(0, 10) < hoje && <Badge tone="danger">Atrasado</Badge>}
            </li>
          ))}
        </ul>
        <Pagination meta={lista.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
