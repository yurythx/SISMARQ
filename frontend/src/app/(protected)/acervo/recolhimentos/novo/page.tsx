"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";

import { MOTIVOS_RECOLHIMENTO, PREPARO, TIPOS, fmtDia, recolhiveis, situacao } from "@/components/acervo/labels";
import { UnidadeDonaSelect } from "@/components/iam/UnidadeDonaSelect";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoPreparo, AcervoRecolhimento, AcervoUnidade } from "@/lib/sismarq/types";

/** Novo recolhimento ao Arquivo Público. */
export default function NovoRecolhimentoPage() {
  const router = useRouter();
  const { run, pending } = useAction();
  const [origem, setOrigem] = useState("");
  const [motivo, setMotivo] = useState("REGULAR");
  const [marcadas, setMarcadas] = useState<Set<string>>(new Set());
  const [preparo, setPreparo] = useState<AcervoPreparo>({
    classificados: false,
    avaliados: false,
    organizados: false,
    higienizados: false,
    acondicionados: false,
  });
  const lista = useApiPage<AcervoUnidade>(origem ? withQuery("v1/acervo/unidades", { unidade_id: origem, page_size: 200 }) : null);
  const itens = recolhiveis(lista.data?.items ?? [], motivo);
  const completo = PREPARO.every((p) => preparo[p.campo]);

  function alternar(id: string) {
    const n = new Set(marcadas);
    if (n.has(id)) n.delete(id);
    else n.add(id);
    setMarcadas(n);
  }

  async function enviar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const txt = (k: string) => String(f.get(k) ?? "").trim();
    const r = await run(
      () =>
        apiClient.post<AcervoRecolhimento>("v1/acervo/recolhimentos", {
          origem_id: origem,
          motivo,
          fundamento: motivo === "REGULAR" ? "" : txt("fundamento"),
          preparo,
          historico: txt("historico"),
          observacao: txt("observacao"),
          unidades: [...marcadas],
        }),
      "Recolhimento enviado ao Arquivo Público.",
    );
    if (r) router.push(`/acervo/recolhimentos/${r.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/recolhimentos" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Recolhimentos
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Novo recolhimento"
        description="Relacione as unidades de guarda permanente e declare o preparo do art. 53. As unidades só passam à guarda do Arquivo Público quando ele registra a entrada; as despesas de preparo e transporte são do órgão (art. 53, § 2º)."
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
        aria-label="Novo recolhimento"
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <UnidadeDonaSelect id="origem_id" name="origem_id" permission="acervo:manage" label="Unidade de origem" />
          <Select
            id="motivo"
            label="Motivo *"
            value={motivo}
            onChange={(e) => {
              setMotivo(e.target.value);
              setMarcadas(new Set());
            }}
            options={Object.entries(MOTIVOS_RECOLHIMENTO).map(([value, label]) => ({ value, label }))}
          />
        </div>
        {motivo !== "REGULAR" && (
          <Input
            id="fundamento"
            name="fundamento"
            label="Ato de extinção ou edital de desestatização *"
            required
            maxLength={500}
            placeholder="Ex.: Lei nº …, de …, que extingue a autarquia"
          />
        )}
        {origem && (
          <DataState
            loading={lista.isLoading}
            error={lista.error}
            empty={itens.length === 0}
            emptyTitle="Nenhuma unidade para recolher"
            emptyDescription={
              motivo === "REGULAR"
                ? "Aparecem aqui as unidades com os prazos cumpridos e destinação guarda permanente."
                : "Aparecem aqui as unidades de guarda permanente ainda no órgão."
            }
          >
            <fieldset className="flex flex-col gap-2">
              <legend className="mb-2 text-sm font-semibold">Unidades ({marcadas.size} marcada(s))</legend>
              {itens.map((u) => (
                <label key={u.id} className="flex items-start gap-2 rounded-lg border border-surface-border p-3 text-sm">
                  <input type="checkbox" checked={marcadas.has(u.id)} onChange={() => alternar(u.id)} className="mt-1" />
                  <span>
                    <span className="font-mono text-xs text-muted">{u.numero}</span> {u.assunto} — {TIPOS[u.tipo]},{" "}
                    <span className="font-mono">{u.codigo_classificacao}</span>, {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}{" "}
                    <span className="text-xs text-muted">({situacao(u.situacao).label})</span>
                  </span>
                </label>
              ))}
            </fieldset>
          </DataState>
        )}
        <fieldset className="flex flex-col gap-2">
          <legend className="mb-1 text-sm font-semibold">Preparo do conjunto (art. 53) *</legend>
          <div className="flex flex-wrap gap-4">
            {PREPARO.map((p) => (
              <label key={p.campo} className="flex items-center gap-2 text-sm">
                <input type="checkbox" checked={preparo[p.campo]} onChange={(e) => setPreparo({ ...preparo, [p.campo]: e.target.checked })} />
                {p.label}
              </label>
            ))}
          </div>
          <p className="text-xs text-muted">O instrumento descritivo é gerado pelo sistema a partir das unidades relacionadas.</p>
        </fieldset>
        <Textarea
          id="historico"
          name="historico"
          label="História administrativa e âmbito do conjunto"
          maxLength={8000}
          rows={3}
          placeholder="Quem produziu, em que período e para quê — compõe o instrumento descritivo"
        />
        <Input id="observacao" name="observacao" label="Observação" maxLength={1000} />
        <div>
          <Button type="submit" disabled={pending || marcadas.size === 0 || !completo}>
            Enviar ao Arquivo Público
          </Button>
        </div>
      </form>
    </div>
  );
}
