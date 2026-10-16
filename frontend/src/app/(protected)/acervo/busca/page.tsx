"use client";

import { ArrowLeft, Search } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { FACETAS_INTERNAS, ItemResultado, NOME_FACETA, rotuloFaceta } from "@/components/acervo/busca";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Button } from "@/components/ui/Button";
import { useApiQuery, withQuery } from "@/lib/api/swr";
import type { AcervoBusca } from "@/lib/sismarq/types";

/** Busca no acervo e no texto do OCR, com facetas, no escopo de quem pede:
 * o restrito e o sigiloso só aparecem para quem tem a unidade (ADR 038). */
export default function BuscaAcervoPage() {
  const [texto, setTexto] = useState("");
  const [q, setQ] = useState("");
  const [filtros, setFiltros] = useState<Record<string, string>>({});
  const [digital, setDigital] = useState(false);
  const [page, setPage] = useState(1);
  const res = useApiQuery<AcervoBusca>(
    withQuery("v1/acervo/busca", { q: q || undefined, ...filtros, digital: digital || undefined, page, page_size: 20 }),
  );
  const r = res.data;

  function alternar(k: string, v: string) {
    setFiltros((f) => {
      const n = { ...f };
      if (n[k] === v) delete n[k];
      else n[k] = v;
      return n;
    });
    setPage(1);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Busca no acervo"
        description="Nas descrições das unidades e no texto dos documentos digitalizados (OCR). Você vê o que está no seu escopo; o Arquivo Público vê todo o acervo."
      />
      <form
        role="search"
        className="flex flex-wrap gap-2"
        onSubmit={(e) => {
          e.preventDefault();
          setQ(texto.trim());
          setPage(1);
        }}
      >
        <label htmlFor="busca-acervo" className="sr-only">
          Buscar
        </label>
        <input
          id="busca-acervo"
          value={texto}
          onChange={(e) => setTexto(e.target.value)}
          placeholder="Assunto, número, código ou palavra do documento"
          className="min-w-0 flex-1 rounded-md border border-surface-border bg-surface px-3 py-2 text-sm"
        />
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={digital}
            onChange={(e) => {
              setDigital(e.target.checked);
              setPage(1);
            }}
          />
          Só com documento digital
        </label>
        <Button type="submit">
          <Search size={16} aria-hidden="true" /> Buscar
        </Button>
      </form>
      <div className="grid gap-6 md:grid-cols-[16rem_1fr]">
        <aside aria-label="Filtros" className="flex flex-col gap-4">
          {r?.origem === "banco" && <p className="text-xs text-muted">Índice da busca indisponível: busca simplificada, sem facetas.</p>}
          {FACETAS_INTERNAS.map((k) => {
            const vs = r?.facetas[k] ?? [];
            if (vs.length === 0) return null;
            return (
              <fieldset key={k}>
                <legend className="mb-1 text-sm font-semibold">{NOME_FACETA[k]}</legend>
                <ul className="flex flex-col gap-0.5 text-sm">
                  {vs.map((v) => (
                    <li key={v.valor}>
                      <button
                        type="button"
                        aria-pressed={filtros[k] === v.valor}
                        onClick={() => alternar(k, v.valor)}
                        className="flex w-full justify-between gap-2 text-left hover:underline aria-pressed:font-semibold"
                      >
                        <span>{rotuloFaceta(k, v.valor)}</span>
                        <span className="text-muted">{v.total}</span>
                      </button>
                    </li>
                  ))}
                </ul>
              </fieldset>
            );
          })}
        </aside>
        <DataState loading={res.isLoading} error={res.error} empty={!r || r.items.length === 0} emptyTitle="Nada encontrado no seu escopo">
          {r && (
            <section aria-label="Resultados" className="flex flex-col gap-3">
              <p className="text-sm text-muted" aria-live="polite">
                {r.meta.total_items.toLocaleString("pt-BR")} unidade(s){r.origem === "indice" && ` · ${r.ms} ms`}
              </p>
              <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Unidades encontradas">
                {r.items.map((it) => (
                  <ItemResultado key={it.id} item={it} href={`/acervo/${it.id}`} />
                ))}
              </ul>
              <Pagination meta={r.meta} onPage={setPage} />
            </section>
          )}
        </DataState>
      </div>
    </div>
  );
}
