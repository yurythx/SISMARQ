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
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoListagem, AcervoUnidade } from "@/lib/sismarq/types";

/** Nova relação de eliminação: as unidades da unidade com os prazos
 * cumpridos e destinação eliminação (arts. 19 e 33). */
export default function NovaRelacaoPage() {
  const router = useRouter();
  const { run, pending } = useAction();
  const [unidade, setUnidade] = useState("");
  const [marcadas, setMarcadas] = useState<Set<string>>(new Set());
  const aptas = useApiPage<AcervoUnidade>(
    unidade ? withQuery("v1/acervo/unidades", { unidade_id: unidade, situacao: "AVALIAR_ELIMINACAO", page_size: 200 }) : null,
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
        apiClient.post<AcervoListagem>("v1/acervo/listagens", {
          unidade_id: unidade,
          unidades: [...marcadas],
          parecer_juridico: String(f.get("parecer_juridico") ?? "").trim(),
          contas_exercicios: String(f.get("contas_exercicios") ?? "").trim(),
          contas_aprovadas_em: String(f.get("contas_aprovadas_em") ?? ""),
        }),
      "Relação de eliminação aberta.",
    );
    if (r) router.push(`/acervo/eliminacao/${r.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/eliminacao" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Eliminação de documentos
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Nova relação de eliminação"
        description="Só entram unidades com os prazos de guarda cumpridos, destinação eliminação na TTDD, sem bloqueio judicial, fora de empréstimo e de outra relação (arts. 19, 20, 30, 42 e 44). O órgão precisa ter a CSPAD ativa (art. 31)."
      />
      <form
        onSubmit={enviar}
        onChange={(e) => {
          const alvo = e.target as unknown as HTMLSelectElement;
          if (alvo.name === "unidade_id") {
            setUnidade(alvo.value);
            setMarcadas(new Set());
          }
        }}
        className="flex flex-col gap-5"
        aria-label="Nova relação de eliminação"
      >
        <UnidadeDonaSelect id="unidade_id" permission="acervo:avaliacao" label="Unidade/setor" />
        {unidade && (
          <DataState
            loading={aptas.isLoading}
            error={aptas.error}
            empty={itens.length === 0}
            emptyTitle="Nenhuma unidade para avaliar"
            emptyDescription="Aparecem aqui as unidades com os prazos cumpridos e destinação eliminação."
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
        <Textarea
          id="parecer_juridico"
          name="parecer_juridico"
          label="Parecer jurídico (art. 32, § 4º)"
          maxLength={4000}
          rows={3}
          placeholder="Número e conclusão do parecer — exigido antes da aprovação pela CSPAD"
        />
        <div className="grid gap-4 sm:grid-cols-2">
          <Input id="contas_exercicios" name="contas_exercicios" label="Contas do(s) exercício(s) de" maxLength={200} />
          <Input id="contas_aprovadas_em" name="contas_aprovadas_em" type="date" label="Contas aprovadas pelo Tribunal de Contas em" />
        </div>
        <p className="text-xs text-muted">
          O quadro das contas só é preenchido quando os documentos dependem da aprovação das contas pelo Tribunal de Contas (Anexo I).
        </p>
        <div>
          <Button type="submit" disabled={pending || marcadas.size === 0}>
            Abrir relação
          </Button>
        </div>
      </form>
    </div>
  );
}
