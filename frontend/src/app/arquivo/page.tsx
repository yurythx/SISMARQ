import type { Metadata } from "next";
import Link from "next/link";
import { BookOpen, FileText, Landmark, Search } from "lucide-react";

import { fmtDia } from "@/components/acervo/labels";
import { PublicShell } from "@/components/layout/PublicShell";
import { publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoPainelPublico } from "@/lib/sismarq/types";

const description =
  "Consulta ao acervo do Arquivo Público Municipal de Rondonópolis (Lei 10.301/2019): a descrição dos documentos de acesso público e os documentos digitalizados liberados, sem login.";

export const metadata: Metadata = {
  title: "Arquivo Público Municipal",
  description,
  alternates: { canonical: `${APP_URL}/arquivo` },
};

/** A página do Arquivo Público: a busca no acervo ostensivo, as exposições
 * virtuais, os últimos documentos liberados e o caminho do SIC para o que
 * não é ostensivo (arts. 2º e 5º, X e XI). */
export default async function ArquivoPublicoPage() {
  const p = await publicGet<AcervoPainelPublico>("acervo/publico/arquivo", 60)
    .then((r) => r.data)
    .catch(() => null);

  return (
    <PublicShell>
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-10 px-6 py-12">
        <header className="flex flex-col gap-3">
          <p className="dateline">Sistema Municipal de Arquivos · Lei 10.301/2019</p>
          <h1 className="text-3xl font-bold">Arquivo Público Municipal</h1>
          <p className="max-w-3xl text-muted">{description}</p>
          <form action="/arquivo/busca" method="get" role="search" className="flex max-w-2xl gap-2">
            <label htmlFor="q" className="sr-only">
              Buscar no acervo
            </label>
            <input
              id="q"
              name="q"
              placeholder="Assunto, número, código ou palavra do documento"
              className="flex-1 rounded-md border border-surface-border bg-surface px-3 py-2 text-sm"
            />
            <button type="submit" className="inline-flex items-center gap-1 rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground">
              <Search size={16} aria-hidden="true" /> Buscar
            </button>
          </form>
        </header>

        {p === null ? (
          <p role="alert" className="text-sm text-muted">
            A consulta ao acervo está indisponível no momento. Tente de novo em instantes.
          </p>
        ) : (
          <>
            <dl className="grid gap-4 sm:grid-cols-2" aria-label="O acervo em números">
              <div className="rounded-xl border border-surface-border bg-surface p-4">
                <dt className="text-sm text-muted">Unidades de arquivamento de acesso público</dt>
                <dd className="text-2xl font-bold">{p.unidades.toLocaleString("pt-BR")}</dd>
              </div>
              <div className="rounded-xl border border-surface-border bg-surface p-4">
                <dt className="text-sm text-muted">Documentos digitalizados para consulta</dt>
                <dd className="text-2xl font-bold">{p.digitais.toLocaleString("pt-BR")}</dd>
              </div>
            </dl>

            {p.exposicoes.length > 0 && (
              <section aria-labelledby="exposicoes" className="flex flex-col gap-3">
                <h2 id="exposicoes" className="text-xl font-semibold">
                  Exposições virtuais
                </h2>
                <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {p.exposicoes.map((x) => (
                    <li key={x.id} className="flex flex-col gap-1 rounded-xl border border-surface-border bg-surface p-4">
                      <Link href={`/arquivo/exposicoes/${x.slug}`} className="font-semibold hover:underline">
                        {x.titulo}
                      </Link>
                      {x.resumo && <p className="text-sm text-muted">{x.resumo}</p>}
                    </li>
                  ))}
                </ul>
              </section>
            )}

            {p.recentes.length > 0 && (
              <section aria-labelledby="recentes" className="flex flex-col gap-3">
                <h2 id="recentes" className="text-xl font-semibold">
                  Documentos digitalizados recentemente
                </h2>
                <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface">
                  {p.recentes.map((d) => (
                    <li key={d.id} className="flex flex-wrap items-center justify-between gap-2 p-4">
                      <Link href={`/arquivo/unidades/${d.unidade_arquivamento_id}?doc=${d.id}`} className="inline-flex items-center gap-2 hover:underline">
                        <FileText size={16} aria-hidden="true" /> {d.descricao || d.assunto || d.numero}
                      </Link>
                      <span className="text-xs text-muted">
                        {d.paginas} página(s) · liberado em {fmtDia(d.publicado_em)}
                      </span>
                    </li>
                  ))}
                </ul>
              </section>
            )}

            <section aria-labelledby="acesso" className="flex flex-col gap-2 rounded-xl border border-surface-border bg-surface p-5">
              <h2 id="acesso" className="inline-flex items-center gap-2 text-lg font-semibold">
                <Landmark size={18} aria-hidden="true" /> Não encontrou o que procura?
              </h2>
              <p className="text-sm text-muted">
                A consulta mostra o que é de acesso público. Documentos com informação pessoal ou sigilosa têm acesso restrito por lei
                (Lei 12.527/2011) e não aparecem aqui. O pedido de acesso é feito pelo Serviço de Informação ao Cidadão (SIC) da Prefeitura,
                que responde no prazo da Lei de Acesso à Informação.
              </p>
              {p.sic_url ? (
                <a href={p.sic_url} rel="noopener noreferrer" target="_blank" className="text-sm font-medium text-primary hover:underline">
                  Fazer um pedido de acesso pelo SIC
                </a>
              ) : (
                <p className="text-sm">Procure o SIC na sede da Prefeitura ou pelos canais de atendimento.</p>
              )}
            </section>

            <Link href="/normas" className="inline-flex items-center gap-2 text-sm text-primary hover:underline">
              <BookOpen size={16} aria-hidden="true" /> Instruções normativas do Arquivo Público
            </Link>
          </>
        )}
      </div>
    </PublicShell>
  );
}
