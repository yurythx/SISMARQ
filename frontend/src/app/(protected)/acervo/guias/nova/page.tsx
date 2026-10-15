"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import { TIPOS, fmtDia } from "@/components/acervo/labels";
import { UnidadeDonaSelect } from "@/components/iam/UnidadeDonaSelect";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoGuia, AcervoUnidade } from "@/lib/sismarq/types";

/** Nova guia: as unidades da origem com o prazo corrente cumprido. */
export default function NovaGuiaPage() {
  const router = useRouter();
  const { run, pending } = useAction();
  const [origem, setOrigem] = useState("");
  const [marcadas, setMarcadas] = useState<Set<string>>(new Set());
  const aptas = useApiPage<AcervoUnidade>(
    origem ? withQuery("v1/acervo/unidades", { unidade_id: origem, situacao: "TRANSFERIR", page_size: 200 }) : null,
  );
  const itens = aptas.data?.items ?? [];

  function alternar(id: string) {
    const n = new Set(marcadas);
    if (n.has(id)) n.delete(id);
    else n.add(id);
    setMarcadas(n);
  }

  async function enviar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const r = await run(
      () =>
        apiClient.post<AcervoGuia>("v1/acervo/guias", {
          origem_id: origem,
          unidades: [...marcadas],
          observacao: String(f.get("observacao") ?? "").trim(),
        }),
      "Guia de transferência gerada.",
    );
    if (r) router.push(`/acervo/guias/${r.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/guias" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Guias de transferência
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Nova guia de transferência"
        description="Só entram as unidades com o prazo da fase corrente cumprido (art. 29), fora de empréstimo e de outra guia. Elas mudam de fase quando o arquivo intermediário registra o recebimento."
      />
      <form
        onSubmit={enviar}
        onChange={(e) => {
          const alvo = e.target as unknown as HTMLSelectElement;
          if (alvo.name === "origem_id") {
            setOrigem(alvo.value);
            setMarcadas(new Set());
          }
        }}
        className="flex flex-col gap-5"
        aria-label="Nova guia de transferência"
      >
        <UnidadeDonaSelect id="origem_id" name="origem_id" permission="acervo:manage" label="Unidade de origem" />
        {origem && (
          <DataState
            loading={aptas.isLoading}
            error={aptas.error}
            empty={itens.length === 0}
            emptyTitle="Nenhuma unidade pronta para transferir"
            emptyDescription="Aparecem aqui as unidades com o prazo da fase corrente cumprido."
          >
            <fieldset className="flex flex-col gap-2">
              <legend className="mb-2 text-sm font-semibold">Unidades ({marcadas.size} marcada(s))</legend>
              {itens.map((u) => (
                <label key={u.id} className="flex items-start gap-2 rounded-lg border border-surface-border p-3 text-sm">
                  <input type="checkbox" checked={marcadas.has(u.id)} onChange={() => alternar(u.id)} className="mt-1" />
                  <span>
                    <span className="font-mono text-xs text-muted">{u.numero}</span> {u.assunto} — {TIPOS[u.tipo]},{" "}
                    <span className="font-mono">{u.codigo_classificacao}</span>, {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}
                  </span>
                </label>
              ))}
            </fieldset>
          </DataState>
        )}
        <Input id="observacao" name="observacao" label="Observação" maxLength={1000} />
        <div>
          <Button type="submit" disabled={pending || marcadas.size === 0}>
            Gerar guia
          </Button>
        </div>
      </form>
    </div>
  );
}
