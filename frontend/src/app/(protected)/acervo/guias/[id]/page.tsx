"use client";

import { ArrowLeft, Printer } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";

import { TIPOS, fmtDia, situacaoGuia } from "@/components/acervo/labels";
import { ConfirmButton } from "@/components/sismarq/ConfirmButton";
import { DataState } from "@/components/sismarq/DataState";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoGuia } from "@/lib/sismarq/types";

/** A guia de transferência, para imprimir e assinar, e o aceite. */
export default function GuiaPage() {
  const { id } = useParams<{ id: string }>();
  const { can } = useSismarq();
  const guia = useApiQuery<AcervoGuia>(`v1/acervo/guias/${id}`);
  const { run, pending } = useAction();
  const g = guia.data;

  async function acao(qual: "receber" | "cancelar") {
    const r = await run(
      () => apiClient.post<AcervoGuia>(`v1/acervo/guias/${id}/${qual}`),
      qual === "receber" ? "Recebimento registrado: unidades no arquivo intermediário." : "Guia cancelada.",
    );
    if (r) await guia.mutate(r.data, { revalidate: false });
  }

  const s = g ? situacaoGuia(g.situacao) : null;
  return (
    <div className="mx-auto flex max-w-4xl flex-col gap-5 text-foreground">
      <div className="flex flex-wrap items-center justify-between gap-2 print:hidden">
        <Link href="/acervo/guias" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Guias de transferência
        </Link>
        <div className="flex gap-2">
          {g?.situacao === "ENVIADA" && can("acervo:manage") && (
            <>
              <ConfirmButton
                size="sm"
                variant="ghost"
                onConfirm={() => acao("cancelar")}
                confirmLabel="Cancelar a guia"
                title="Cancelar a guia?"
                description="As unidades ficam livres para outra remessa."
              >
                Cancelar guia
              </ConfirmButton>
              <Button size="sm" disabled={pending} onClick={() => void acao("receber")}>
                Registrar recebimento
              </Button>
            </>
          )}
          <Button size="sm" variant="secondary" onClick={() => window.print()}>
            <Printer size={14} aria-hidden="true" className="mr-1" /> Imprimir
          </Button>
        </div>
      </div>
      <DataState loading={guia.isLoading} error={guia.error} empty={!g}>
        {g && s && (
          <article aria-label={`Guia ${g.numero}`} className="flex flex-col gap-4">
            <header className="flex flex-col gap-1 border-b border-surface-border pb-3">
              <p className="text-xs font-semibold uppercase tracking-wide">
                Prefeitura Municipal de Rondonópolis · Sistema Municipal de Arquivos — SISMARQ
              </p>
              <h1 className="text-2xl font-bold">Guia de transferência {g.numero}</h1>
              <p className="text-sm">
                Transferência ao arquivo intermediário (Lei Municipal nº 10.301/2019, art. 29) · emitida em {fmtDia(g.criada_em)}{" "}
                <Badge tone={s.tone}>{s.label}</Badge>
              </p>
              {g.observacao && <p className="text-sm">Observação: {g.observacao}</p>}
            </header>
            <table className="w-full border-collapse text-sm">
              <caption className="sr-only">Unidades relacionadas na guia</caption>
              <thead>
                <tr className="border-b border-surface-border text-left text-xs uppercase">
                  <th scope="col" className="py-2 pr-2">
                    Unidade
                  </th>
                  <th scope="col" className="py-2 pr-2">
                    Código
                  </th>
                  <th scope="col" className="py-2 pr-2">
                    Assunto
                  </th>
                  <th scope="col" className="py-2 pr-2">
                    Datas-limite
                  </th>
                  <th scope="col" className="py-2">
                    Quantificação
                  </th>
                </tr>
              </thead>
              <tbody>
                {g.itens.map((u) => (
                  <tr key={u.id} className="border-b border-surface-border align-top">
                    <td className="py-2 pr-2 font-mono text-xs">{u.numero}</td>
                    <td className="py-2 pr-2 font-mono text-xs">{u.codigo_classificacao}</td>
                    <td className="py-2 pr-2">{u.assunto}</td>
                    <td className="py-2 pr-2">
                      {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}
                    </td>
                    <td className="py-2">
                      {u.quantidade} {u.especificacao || TIPOS[u.tipo].toLowerCase()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <p className="text-sm">Total: {g.total_itens} unidade(s) de arquivamento.</p>
            <div className="mt-8 grid gap-8 sm:grid-cols-2">
              <div className="border-t border-foreground pt-1 text-xs">Responsável pela remessa (unidade produtora)</div>
              <div className="border-t border-foreground pt-1 text-xs">
                Recebido pelo arquivo intermediário{g.recebida_em ? ` em ${fmtDia(g.recebida_em)}` : ""}
              </div>
            </div>
          </article>
        )}
      </DataState>
    </div>
  );
}
