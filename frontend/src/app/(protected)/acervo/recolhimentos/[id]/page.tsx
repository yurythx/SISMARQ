"use client";

import { ArrowLeft, Printer } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import type { FormEvent } from "react";

import { MOTIVOS_RECOLHIMENTO, PREPARO, TIPOS, ano, fmtDia, situacaoRecolhimento } from "@/components/acervo/labels";
import { useOrgao } from "@/components/acervo/useOrgao";
import { ConfirmButton } from "@/components/sismarq/ConfirmButton";
import { DataState } from "@/components/sismarq/DataState";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoLocal, AcervoRecolhimento } from "@/lib/sismarq/types";

/** O recolhimento: instrumento descritivo (art. 53) e, depois da entrada,
 * o registro de entrada no Arquivo Público (art. 5º, VII) — para imprimir e
 * assinar. */
export default function RecolhimentoPage() {
  const { id } = useParams<{ id: string }>();
  const { can, canGlobal } = useSismarq();
  const q = useApiQuery<AcervoRecolhimento>(`v1/acervo/recolhimentos/${id}`);
  const central = canGlobal("acervo:central");
  const locais = useApiQuery<AcervoLocal[]>(central ? "v1/acervo/locais" : null);
  const { run, pending } = useAction();
  const r = q.data;
  const { orgao, unidade } = useOrgao(r?.origem_id);

  async function enviar(caminho: string, corpo: unknown, ok: string) {
    const out = await run(() => apiClient.post<AcervoRecolhimento>(`v1/acervo/recolhimentos/${id}/${caminho}`, corpo), ok);
    if (out) await q.mutate(out.data, { revalidate: false });
  }
  const form = (fn: (f: FormData) => Promise<void>) => (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void fn(new FormData(e.currentTarget));
  };

  const s = r ? situacaoRecolhimento(r.situacao) : null;
  const volumes = r?.itens.reduce((t, u) => t + u.quantidade, 0) ?? 0;
  const anos = r?.itens.flatMap((u) => [ano(u.data_inicial), ano(u.data_final)]).sort() ?? [];
  return (
    <div className="mx-auto flex max-w-5xl flex-col gap-5 text-foreground">
      <div className="flex flex-wrap items-center justify-between gap-2 print:hidden">
        <Link href="/acervo/recolhimentos" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Recolhimentos
        </Link>
        <div className="flex gap-2">
          {r?.situacao === "ENVIADO" && can("acervo:manage") && (
            <ConfirmButton
              size="sm"
              variant="ghost"
              onConfirm={() => enviar("cancelar", undefined, "Recolhimento cancelado.")}
              confirmLabel="Cancelar o recolhimento"
              title="Cancelar o recolhimento?"
              description="As unidades continuam sob a guarda do órgão."
            >
              Cancelar recolhimento
            </ConfirmButton>
          )}
          <Button size="sm" variant="secondary" onClick={() => window.print()}>
            <Printer size={14} aria-hidden="true" className="mr-1" /> Imprimir
          </Button>
        </div>
      </div>

      <DataState loading={q.isLoading} error={q.error} empty={!r}>
        {r && s && (
          <>
            {r.situacao === "ENVIADO" && central && (
              <section aria-label="Decisão do Arquivo Público" className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4 print:hidden">
                <h2 className="text-base font-semibold">Entrada no Arquivo Público</h2>
                <form
                  aria-label="Registrar entrada"
                  className="flex flex-wrap items-end gap-3"
                  onSubmit={form((f) =>
                    enviar("receber", { local_id: String(f.get("local_id") ?? "") || null }, "Entrada registrada; a origem recebeu a cópia do registro."),
                  )}
                >
                  <div className="min-w-64">
                    <Select
                      id="local_id"
                      name="local_id"
                      label="Guardar em"
                      placeholder="Sem localização"
                      options={(locais.data ?? []).filter((l) => l.ativo).map((l) => ({ value: l.id, label: l.caminho }))}
                    />
                  </div>
                  <Button type="submit" disabled={pending}>
                    Registrar entrada
                  </Button>
                </form>
                <form
                  aria-label="Devolver à origem"
                  className="flex flex-col gap-2"
                  onSubmit={form((f) => enviar("devolver", { motivo: String(f.get("motivo") ?? "").trim() }, "Recolhimento devolvido à origem."))}
                >
                  <Textarea id="motivo" name="motivo" label="O que falta (devolução)" required maxLength={2000} rows={2} />
                  <div>
                    <Button type="submit" variant="secondary" disabled={pending}>
                      Devolver à origem
                    </Button>
                  </div>
                </form>
              </section>
            )}

            <article aria-label={`Recolhimento ${r.numero}`} className="flex flex-col gap-4 text-sm leading-relaxed">
              <header className="flex flex-col gap-1 border-b border-foreground pb-3 text-center">
                <p className="text-xs font-semibold uppercase">Prefeitura Municipal de Rondonópolis · Sistema Municipal de Arquivos — SISMARQ</p>
                <h1 className="text-lg font-bold uppercase">
                  {r.registro_entrada ? `Registro de entrada ${r.registro_entrada}` : `Recolhimento ${r.numero}`}
                </h1>
                <p>
                  Recolhimento {r.numero} ao Arquivo Público Municipal (Lei Municipal nº 10.301/2019, arts. 5º, VII, e 53){" "}
                  <Badge tone={s.tone}>{s.label}</Badge>
                </p>
              </header>
              <dl className="grid gap-1 sm:grid-cols-2">
                <div>
                  <dt className="inline font-semibold">ÓRGÃO/ENTIDADE: </dt>
                  <dd className="inline">{orgao}</dd>
                </div>
                <div>
                  <dt className="inline font-semibold">UNIDADE DE ORIGEM: </dt>
                  <dd className="inline">{unidade}</dd>
                </div>
                <div>
                  <dt className="inline font-semibold">MOTIVO: </dt>
                  <dd className="inline">{MOTIVOS_RECOLHIMENTO[r.motivo] ?? r.motivo}</dd>
                </div>
                {r.fundamento && (
                  <div>
                    <dt className="inline font-semibold">FUNDAMENTO: </dt>
                    <dd className="inline">{r.fundamento}</dd>
                  </div>
                )}
                <div>
                  <dt className="inline font-semibold">DATAS-LIMITE: </dt>
                  <dd className="inline">{anos.length ? `${anos[0]} a ${anos[anos.length - 1]}` : "—"}</dd>
                </div>
                <div>
                  <dt className="inline font-semibold">MENSURAÇÃO: </dt>
                  <dd className="inline">
                    {r.total_itens} unidade(s) de arquivamento ({volumes} volume(s))
                  </dd>
                </div>
              </dl>
              <p>
                <strong>Preparo declarado pela origem (art. 53):</strong>{" "}
                {PREPARO.map((p) => `${p.label}: ${r.preparo[p.campo] ? "sim" : "não"}`).join(" · ")}.
              </p>
              {r.historico && (
                <section aria-label="História administrativa">
                  <h2 className="font-semibold uppercase">História administrativa e âmbito</h2>
                  <p className="whitespace-pre-line">{r.historico}</p>
                </section>
              )}
              {r.situacao === "DEVOLVIDO" && (
                <p className="rounded-lg border border-warning p-2 print:hidden">
                  <strong>Devolvido pelo Arquivo Público:</strong> {r.devolucao_motivo}
                </p>
              )}
              <table className="w-full border-collapse border border-foreground text-xs">
                <caption className="sr-only">Instrumento descritivo: unidades recolhidas</caption>
                <thead>
                  <tr>
                    <th scope="col" className="border border-foreground p-1">Unidade</th>
                    <th scope="col" className="border border-foreground p-1">Código</th>
                    <th scope="col" className="border border-foreground p-1">Assunto</th>
                    <th scope="col" className="border border-foreground p-1">Datas-limite</th>
                    <th scope="col" className="border border-foreground p-1">Quantificação</th>
                    <th scope="col" className="border border-foreground p-1">Acesso</th>
                  </tr>
                </thead>
                <tbody>
                  {r.itens.map((u) => (
                    <tr key={u.id}>
                      <td className="border border-foreground p-1 font-mono">
                        <Link href={`/acervo/${u.id}`} className="hover:underline">
                          {u.numero}
                        </Link>
                      </td>
                      <td className="border border-foreground p-1 font-mono">{u.codigo_classificacao}</td>
                      <td className="border border-foreground p-1">{u.assunto}</td>
                      <td className="border border-foreground p-1">
                        {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}
                      </td>
                      <td className="border border-foreground p-1">
                        {u.quantidade} {u.especificacao || TIPOS[u.tipo].toLowerCase()}
                      </td>
                      <td className="border border-foreground p-1">{u.nivel_acesso === "PUBLICO" ? "Público" : `${u.nivel_acesso.toLowerCase()} — ${u.hipotese_legal}`}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {r.observacao && <p>Observação: {r.observacao}</p>}
              <div className="mt-8 grid gap-8 sm:grid-cols-2">
                <div className="border-t border-foreground pt-1 text-center text-xs uppercase">
                  Responsável pelo recolhimento (origem) · {fmtDia(r.criado_em)}
                </div>
                <div className="border-t border-foreground pt-1 text-center text-xs uppercase">
                  Arquivo Público Municipal{r.recebido_em ? ` · entrada em ${fmtDia(r.recebido_em)}` : ""}
                </div>
              </div>
            </article>
          </>
        )}
      </DataState>
    </div>
  );
}
