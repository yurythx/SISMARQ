"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import { situacaoPrivado } from "@/components/acervo/labels";
import { CamposPrivado, corpoPrivado } from "@/components/acervo/privado";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoPrivado } from "@/lib/sismarq/types";

/** Arquivos privados de interesse público e social (arts. 55 a 58). Só o
 * Arquivo Público: o cadastro tem dados pessoais (LGPD). */
export default function PrivadosPage() {
  const router = useRouter();
  const { canGlobal } = useSismarq();
  const central = canGlobal("acervo:central");
  const [page, setPage] = useState(1);
  const [q, setQ] = useState("");
  const [novo, setNovo] = useState(false);
  const { run, pending } = useAction();
  const lista = useApiPage<AcervoPrivado>(central ? withQuery("v1/acervo/privados", { q: q || undefined, page, page_size: 20 }) : null);
  const itens = lista.data?.items ?? [];

  async function cadastrar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const r = await run(() => apiClient.post<AcervoPrivado>("v1/acervo/privados", corpoPrivado(new FormData(e.currentTarget))), "Arquivo privado cadastrado.");
    if (r) router.push(`/acervo/privados/${r.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Arquivo Público"
        title="Arquivos privados de interesse público e social"
        description="Arquivos de pessoas físicas ou jurídicas relevantes para a história, a cultura e o desenvolvimento do município: parecer da Comissão Especial, decreto, doação ou depósito e o direito de preferência na venda (arts. 55 a 58)."
        actions={central && !novo && <Button onClick={() => setNovo(true)}>Cadastrar arquivo</Button>}
      />
      {!central ? (
        <p className="text-sm text-muted">O cadastro dos arquivos privados é do Arquivo Público (acervo:central).</p>
      ) : (
        <>
          {novo && (
            <form onSubmit={cadastrar} aria-label="Cadastrar arquivo privado" className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4">
              <CamposPrivado />
              <label className="flex items-center gap-2 text-sm">
                <input type="checkbox" name="servico_publico" /> Entidade privada encarregada de serviço público municipal (art. 56, § 3º)
              </label>
              <div className="flex gap-2">
                <Button type="submit" disabled={pending}>
                  Cadastrar
                </Button>
                <Button type="button" variant="ghost" onClick={() => setNovo(false)}>
                  Cancelar
                </Button>
              </div>
            </form>
          )}
          <Input id="busca" label="Buscar" value={q} onChange={(e) => setQ(e.target.value)} placeholder="Titular, número ou descrição" />
          <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle="Nenhum arquivo privado cadastrado">
            <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Arquivos privados">
              {itens.map((p) => {
                const s = situacaoPrivado(p.situacao);
                return (
                  <li key={p.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                    <div>
                      <Link href={`/acervo/privados/${p.id}`} className="font-medium hover:underline">
                        {p.titular}
                      </Link>
                      <p className="text-xs text-muted">
                        <span className="font-mono">{p.numero}</span> · {p.descricao.slice(0, 120)}
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
