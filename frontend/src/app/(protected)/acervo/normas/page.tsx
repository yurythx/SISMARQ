"use client";

import { ArrowLeft, ExternalLink } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import { fmtDia, situacaoNorma } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoNorma } from "@/lib/sismarq/types";

/** Instruções normativas do Arquivo Público (art. 54): a gestão das
 * redações. O público consulta em /normas. */
export default function NormasPage() {
  const router = useRouter();
  const { canGlobal } = useSismarq();
  const central = canGlobal("acervo:central");
  const [page, setPage] = useState(1);
  const [q, setQ] = useState("");
  const [nova, setNova] = useState(false);
  const { run, pending } = useAction();
  const lista = useApiPage<AcervoNorma>(central ? withQuery("v1/acervo/normas", { q: q || undefined, page, page_size: 20 }) : null);
  const itens = lista.data?.items ?? [];

  async function criar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const r = await run(
      () =>
        apiClient.post<AcervoNorma>("v1/acervo/normas", {
          numero: String(f.get("numero") ?? "").trim(),
          ementa: String(f.get("ementa") ?? "").trim(),
          texto: String(f.get("texto") ?? ""),
        }),
      "Instrução normativa criada em rascunho.",
    );
    if (r) router.push(`/acervo/normas/${r.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Arquivo Público"
        title="Instruções normativas"
        description="Normas do Arquivo Público para os órgãos do SISMARQ (art. 54). Cada alteração é uma nova redação: a publicada nunca muda e fica no histórico público, com o período em que vigorou."
        actions={
          <>
            <Link href="/normas" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
              <ExternalLink size={14} aria-hidden="true" /> Ver como o público
            </Link>
            {central && !nova && <Button onClick={() => setNova(true)}>Nova instrução normativa</Button>}
          </>
        }
      />
      {!central ? (
        <p className="text-sm text-muted">A gestão das instruções normativas é do Arquivo Público (acervo:central). O texto publicado está em /normas.</p>
      ) : (
        <>
          {nova && (
            <form onSubmit={criar} aria-label="Nova instrução normativa" className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4">
              <div className="grid gap-3 sm:grid-cols-[10rem_1fr]">
                <Input id="numero" name="numero" label="Número *" required placeholder="1/2026" pattern="\d{1,4}/\d{4}" />
                <Input id="ementa" name="ementa" label="Ementa *" required maxLength={500} />
              </div>
              <Textarea id="texto" name="texto" label="Texto (Markdown) *" required rows={10} maxLength={100000} />
              <div className="flex gap-2">
                <Button type="submit" disabled={pending}>
                  Criar em rascunho
                </Button>
                <Button type="button" variant="ghost" onClick={() => setNova(false)}>
                  Cancelar
                </Button>
              </div>
            </form>
          )}
          <Input id="busca" label="Buscar" value={q} onChange={(e) => setQ(e.target.value)} placeholder="Número, ementa ou texto" />
          <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle="Nenhuma instrução normativa">
            <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Normas">
              {itens.map((n) => {
                const s = situacaoNorma(n.situacao);
                return (
                  <li key={n.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                    <div>
                      <Link href={`/acervo/normas/${n.id}`} className="font-medium hover:underline">
                        IN nº {n.numero}
                      </Link>
                      <p className="text-sm">{n.ementa}</p>
                      <p className="text-xs text-muted">
                        {n.versao > 0 ? `Redação ${n.versao} desde ${fmtDia(n.publicada_em)}` : "Nunca publicada"}
                        {n.tem_rascunho && " · há redação em rascunho"}
                      </p>
                    </div>
                    <Badge tone={s.tone}>{s.label}</Badge>
                  </li>
                );
              })}
            </ul>
            <Pagination meta={lista.data?.meta} onPage={setPage} />
          </DataState>
        </>
      )}
    </div>
  );
}
