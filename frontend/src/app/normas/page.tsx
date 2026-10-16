import type { Metadata } from "next";
import Link from "next/link";
import { Search } from "lucide-react";

import { fmtDia, situacaoNorma } from "@/components/acervo/labels";
import { PublicShell } from "@/components/layout/PublicShell";
import { Badge } from "@/components/ui/Badge";
import { publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoNorma } from "@/lib/sismarq/types";

const description =
  "Instruções normativas do Arquivo Público Municipal de Rondonópolis (Lei 10.301/2019, art. 54): o texto em vigor e o histórico de cada alteração.";

export const metadata: Metadata = {
  title: "Instruções normativas",
  description,
  alternates: { canonical: `${APP_URL}/normas` },
};

type Props = { searchParams: Promise<{ q?: string }> };

/** Normas do Arquivo Público, sem login (transparência ativa). */
export default async function NormasPublicasPage({ searchParams }: Props) {
  const { q = "" } = await searchParams;
  const filtro = new URLSearchParams({ page_size: "100" });
  if (q) filtro.set("q", q);
  const normas = await publicGet<AcervoNorma[]>(`acervo/publico/normas?${filtro}`, 60)
    .then((r) => r.data)
    .catch(() => null);

  return (
    <PublicShell>
      <div className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 py-12">
        <div>
          <p className="dateline">Arquivo Público Municipal</p>
          <h1 className="mt-2 text-3xl font-bold">Instruções normativas</h1>
          <p className="mt-1 text-muted">{description}</p>
        </div>
        <form action="/normas" method="get" role="search" className="flex gap-2">
          <label htmlFor="q" className="sr-only">
            Buscar nas normas
          </label>
          <input
            id="q"
            name="q"
            defaultValue={q}
            placeholder="Número, assunto ou trecho do texto"
            className="flex-1 rounded-md border border-surface-border bg-surface px-3 py-2 text-sm"
          />
          <button type="submit" className="inline-flex items-center gap-1 rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground">
            <Search size={16} aria-hidden="true" /> Buscar
          </button>
        </form>
        {normas === null ? (
          <p role="alert" className="text-sm text-muted">
            A consulta às normas está indisponível no momento. Tente de novo em instantes.
          </p>
        ) : normas.length === 0 ? (
          <p className="text-sm text-muted">{q ? "Nenhuma norma encontrada para a busca." : "Nenhuma instrução normativa publicada ainda."}</p>
        ) : (
          <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Instruções normativas">
            {normas.map((n) => {
              const s = situacaoNorma(n.situacao);
              return (
                <li key={n.id} className="flex flex-col gap-1 p-4">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <Link href={`/normas/${n.id}`} className="font-semibold hover:underline">
                      Instrução Normativa nº {n.numero}
                    </Link>
                    <Badge tone={s.tone}>{s.label}</Badge>
                  </div>
                  <p className="text-sm">{n.ementa}</p>
                  <p className="text-xs text-muted">
                    Redação em vigor desde {fmtDia(n.publicada_em)}
                    {n.total_redacoes > 1 && ` · alterada ${n.total_redacoes - 1} vez(es)`}
                    {n.revogada_em && ` · revogada em ${fmtDia(n.revogada_em)}`}
                  </p>
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </PublicShell>
  );
}
