import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft, FileText } from "lucide-react";

import { fmtDia } from "@/components/acervo/labels";
import { PublicShell } from "@/components/layout/PublicShell";
import { Markdown } from "@/components/sismarq/Markdown";
import { PublicApiError, publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoExposicao } from "@/lib/sismarq/types";

type Props = { params: Promise<{ slug: string }> };

async function carregar(slug: string): Promise<AcervoExposicao | null> {
  try {
    return (await publicGet<AcervoExposicao>(`acervo/publico/exposicoes/${encodeURIComponent(slug)}`, 60)).data;
  } catch (e) {
    if (e instanceof PublicApiError && e.status === 404) return null;
    throw e;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const x = await carregar((await params).slug);
  if (!x) return { title: "Exposição não encontrada" };
  return { title: x.titulo, description: x.resumo, alternates: { canonical: `${APP_URL}/arquivo/exposicoes/${x.slug}` } };
}

/** A exposição virtual do Arquivo Público (difusão, art. 5º, X e XI). */
export default async function ExposicaoPage({ params }: Props) {
  const x = await carregar((await params).slug);
  if (!x) notFound();
  const itens = x.itens ?? [];

  return (
    <PublicShell>
      <article className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 py-12">
        <Link href="/arquivo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Arquivo Público
        </Link>
        <header className="flex flex-col gap-2">
          <p className="dateline">Exposição virtual · publicada em {fmtDia(x.publicada_em)}</p>
          <h1 className="text-3xl font-bold">{x.titulo}</h1>
          {x.resumo && <p className="text-lg text-muted">{x.resumo}</p>}
        </header>
        {x.texto && <Markdown source={x.texto} />}
        <section aria-labelledby="documentos" className="flex flex-col gap-3">
          <h2 id="documentos" className="text-xl font-semibold">
            Documentos da exposição
          </h2>
          <ol className="flex flex-col gap-3">
            {itens.map((it) => (
              <li key={it.digital_id} className="rounded-xl border border-surface-border bg-surface p-4">
                <Link href={`/arquivo/unidades/${it.unidade_arquivamento_id}?doc=${it.digital_id}`} className="inline-flex items-center gap-2 font-medium hover:underline">
                  <FileText size={16} aria-hidden="true" /> {it.descricao || it.numero}
                </Link>
                {it.legenda && <p className="mt-1 text-sm">{it.legenda}</p>}
                <p className="mt-1 text-xs text-muted">
                  {it.numero} · {it.paginas} página(s) · unidade {it.unidade_numero}
                </p>
              </li>
            ))}
          </ol>
        </section>
      </article>
    </PublicShell>
  );
}
