"use client";

import { ArrowLeft, ExternalLink } from "lucide-react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import type { FormEvent } from "react";

import { fmtDia, situacaoNorma } from "@/components/acervo/labels";
import { ConfirmButton } from "@/components/sismarq/ConfirmButton";
import { DataState } from "@/components/sismarq/DataState";
import { Markdown } from "@/components/sismarq/Markdown";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoNorma } from "@/lib/sismarq/types";

const txt = (f: FormData, k: string) => String(f.get(k) ?? "").trim();

/** A instrução normativa: a redação em vigor, o rascunho da próxima (a
 * alteração), a publicação, a revogação e o histórico. */
export default function NormaPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const q = useApiQuery<AcervoNorma>(`v1/acervo/normas/${id}`);
  const { run, pending } = useAction();
  const n = q.data;
  const vigente = n?.redacoes.find((r) => r.situacao === "PUBLICADA");
  const rascunho = n?.redacoes.find((r) => r.situacao === "RASCUNHO");
  const historico = n?.redacoes.filter((r) => r.situacao === "PUBLICADA") ?? [];

  async function salvar(caminho: string, corpo: unknown, ok: string, metodo: "post" | "put" = "post") {
    const url = `v1/acervo/normas/${id}${caminho}`;
    const r = await run(() => (metodo === "put" ? apiClient.put<AcervoNorma>(url, corpo) : apiClient.post<AcervoNorma>(url, corpo)), ok);
    if (r) await q.mutate(r.data, { revalidate: false });
  }

  async function descartar() {
    const r = await run(() => apiClient.delete<AcervoNorma | null>(`v1/acervo/normas/${id}/rascunho`), "Rascunho descartado.");
    if (!r) return;
    if (n?.situacao === "RASCUNHO") router.push("/acervo/normas");
    else await q.mutate();
  }

  const form = (fn: (f: FormData) => Promise<void>) => (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void fn(new FormData(e.currentTarget));
  };
  const s = n ? situacaoNorma(n.situacao) : null;

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/normas" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Instruções normativas
      </Link>
      <DataState loading={q.isLoading} error={q.error} empty={!n}>
        {n && s && (
          <>
            <PageHeader
              eyebrow="Arquivo Público · art. 54"
              title={`Instrução Normativa nº ${n.numero}`}
              description={n.ementa}
              actions={
                n.situacao !== "RASCUNHO" && (
                  <Link href={`/normas/${n.id}`} className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
                    <ExternalLink size={14} aria-hidden="true" /> Página pública
                  </Link>
                )
              }
            />
            <div className="flex flex-wrap items-center gap-2 text-sm">
              <Badge tone={s.tone}>{s.label}</Badge>
              {vigente && (
                <span className="text-muted">
                  Redação {vigente.versao} em vigor desde {fmtDia(vigente.publicada_em)}
                </span>
              )}
              {n.revogada_em && (
                <span className="text-muted">
                  Revogada em {fmtDia(n.revogada_em)} — {n.revogacao_ato}
                </span>
              )}
            </div>

            {n.situacao !== "REVOGADA" && (
              <section aria-label="Redação em rascunho" className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4">
                <h2 className="text-base font-semibold">
                  {rascunho ? `Redação ${rascunho.versao} em rascunho` : "Alterar a norma (nova redação)"}
                </h2>
                <p className="text-sm text-muted">
                  {vigente
                    ? "A redação em vigor não muda: a alteração é uma nova redação, que só vale depois de publicada. A anterior fica no histórico."
                    : "A norma ainda não foi publicada; o público não a vê."}
                </p>
                <form
                  key={rascunho ? `r${rascunho.versao}` : "nova"}
                  aria-label="Redação"
                  className="flex flex-col gap-3"
                  onSubmit={form((f) =>
                    salvar("/rascunho", { ementa: txt(f, "ementa"), texto: String(f.get("texto") ?? ""), motivo: txt(f, "motivo") }, "Rascunho salvo.", "put"),
                  )}
                >
                  <Input id="ementa" name="ementa" label="Ementa *" required maxLength={500} defaultValue={(rascunho ?? vigente)?.ementa} />
                  <Textarea id="texto" name="texto" label="Texto (Markdown) *" required rows={12} maxLength={100000} defaultValue={(rascunho ?? vigente)?.texto} />
                  <Input
                    id="motivo"
                    name="motivo"
                    label="Motivo da redação *"
                    required
                    maxLength={500}
                    defaultValue={rascunho?.motivo}
                    placeholder={vigente ? "Ex.: Alterada pela IN 2/2027 — novo prazo do art. 4º" : "Redação original"}
                  />
                  <div className="flex flex-wrap gap-2">
                    <Button type="submit" variant="secondary" disabled={pending}>
                      Salvar rascunho
                    </Button>
                    {rascunho && (
                      <ConfirmButton
                        variant="ghost"
                        onConfirm={descartar}
                        confirmLabel="Descartar"
                        title="Descartar o rascunho?"
                        description={n.situacao === "RASCUNHO" ? "A norma nunca foi publicada e será apagada." : "A redação em vigor continua valendo."}
                      >
                        Descartar rascunho
                      </ConfirmButton>
                    )}
                  </div>
                </form>
                {rascunho && (
                  <form
                    aria-label="Publicar"
                    className="flex flex-wrap items-end gap-3 border-t border-surface-border pt-3"
                    onSubmit={form((f) =>
                      salvar("/publicar", { publicada_em: txt(f, "publicada_em"), diario: txt(f, "diario") }, `Redação ${rascunho.versao} publicada.`),
                    )}
                  >
                    <Input id="publicada_em" name="publicada_em" type="date" label="Publicada em *" required />
                    <div className="min-w-64 flex-1">
                      <Input id="diario" name="diario" label="Diário Oficial (edição, página)" maxLength={200} />
                    </div>
                    <Button type="submit" disabled={pending}>
                      Publicar redação {rascunho.versao}
                    </Button>
                  </form>
                )}
              </section>
            )}

            {n.situacao === "VIGENTE" && (
              <form
                aria-label="Revogar"
                className="flex flex-wrap items-end gap-3 rounded-xl border border-surface-border bg-surface p-4"
                onSubmit={form((f) => salvar("/revogar", { revogada_em: txt(f, "revogada_em"), ato: txt(f, "ato") }, "Instrução normativa revogada."))}
              >
                <div className="min-w-64 flex-1">
                  <Input id="ato" name="ato" label="Ato de revogação *" required maxLength={500} placeholder="Ex.: Revogada pela IN 5/2027" />
                </div>
                <Input id="revogada_em" name="revogada_em" type="date" label="Revogada em *" required />
                <Button type="submit" variant="danger" disabled={pending}>
                  Revogar
                </Button>
              </form>
            )}

            {historico.length > 0 && (
              <section aria-labelledby="historico" className="flex flex-col gap-3">
                <h2 id="historico" className="text-lg font-semibold">
                  Redações publicadas
                </h2>
                {historico.map((r) => (
                  <details key={r.versao} open={r === vigente} className="rounded-lg border border-surface-border p-3">
                    <summary className="cursor-pointer text-sm font-medium">
                      Redação {r.versao} — de {fmtDia(r.publicada_em)}
                      {r.vigente_ate ? ` a ${fmtDia(r.vigente_ate)}` : " até hoje"} · {r.motivo}
                    </summary>
                    <div className="mt-3">
                      <Markdown source={r.texto} />
                    </div>
                  </details>
                ))}
              </section>
            )}
          </>
        )}
      </DataState>
    </div>
  );
}
