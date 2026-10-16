import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft } from "lucide-react";

import { fmtDia, situacaoNorma } from "@/components/acervo/labels";
import { PublicShell } from "@/components/layout/PublicShell";
import { Markdown } from "@/components/sismarq/Markdown";
import { Badge } from "@/components/ui/Badge";
import { PublicApiError, publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoNorma, AcervoRedacao } from "@/lib/sismarq/types";

type Props = { params: Promise<{ id: string }> };

async function carregar(id: string): Promise<AcervoNorma | null> {
  try {
    return (await publicGet<AcervoNorma>(`acervo/publico/normas/${encodeURIComponent(id)}`, 60)).data;
  } catch (e) {
    if (e instanceof PublicApiError && (e.status === 404 || e.status === 400)) return null;
    throw e;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const n = await carregar((await params).id);
  if (!n) return { title: "Norma não encontrada" };
  return {
    title: `Instrução Normativa nº ${n.numero}`,
    description: n.ementa,
    alternates: { canonical: `${APP_URL}/normas/${n.id}` },
  };
}

function periodo(r: AcervoRedacao) {
  return `de ${fmtDia(r.publicada_em)}${r.vigente_ate ? ` a ${fmtDia(r.vigente_ate)}` : " até hoje"}`;
}

/** A norma publicada: o texto em vigor e, abaixo, cada redação anterior com
 * o período em que vigorou e o que a alterou (texto compilado). */
export default async function NormaPublicaPage({ params }: Props) {
  const n = await carregar((await params).id);
  if (!n) notFound();
  const [atual, ...anteriores] = n.redacoes;
  const s = situacaoNorma(n.situacao);

  return (
    <PublicShell>
      <article className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 py-12">
        <Link href="/normas" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Instruções normativas
        </Link>
        <header className="flex flex-col gap-2">
          <p className="dateline">Arquivo Público Municipal · Lei 10.301/2019, art. 54</p>
          <h1 className="text-3xl font-bold">Instrução Normativa nº {n.numero}</h1>
          <p className="text-lg">{n.ementa}</p>
          <div className="flex flex-wrap items-center gap-2 text-sm text-muted">
            <Badge tone={s.tone}>{s.label}</Badge>
            {atual && (
              <span>
                Redação {atual.versao}, publicada em {fmtDia(atual.publicada_em)}
                {atual.diario && ` (${atual.diario})`}
              </span>
            )}
          </div>
          {n.revogada_em && (
            <p role="note" className="rounded-lg border border-danger p-3 text-sm">
              Revogada em {fmtDia(n.revogada_em)} — {n.revogacao_ato}. O texto fica disponível para consulta histórica.
            </p>
          )}
        </header>

        {atual && (
          <section aria-label={n.revogada_em ? "Último texto em vigor" : "Texto em vigor"}>
            <Markdown source={atual.texto} />
          </section>
        )}

        {anteriores.length > 0 && (
          <section aria-labelledby="historico" className="flex flex-col gap-3 border-t border-surface-border pt-6">
            <h2 id="historico" className="text-xl font-semibold">
              Redações anteriores
            </h2>
            <p className="text-sm text-muted">Cada alteração gera uma nova redação; as anteriores ficam aqui, com o período em que vigoraram.</p>
            {anteriores.map((r) => (
              <details key={r.versao} className="rounded-lg border border-surface-border p-3">
                <summary className="cursor-pointer text-sm font-medium">
                  Redação {r.versao} — vigente {periodo(r)}
                  <span className="block text-xs font-normal text-muted">
                    {r.motivo}
                    {r.ementa !== atual?.ementa && ` · ementa: ${r.ementa}`}
                  </span>
                </summary>
                <div className="mt-3">
                  <Markdown source={r.texto} />
                </div>
              </details>
            ))}
          </section>
        )}
        {atual && atual.motivo && anteriores.length > 0 && (
          <p className="text-xs text-muted">Motivo da redação em vigor: {atual.motivo}</p>
        )}
      </article>
    </PublicShell>
  );
}
