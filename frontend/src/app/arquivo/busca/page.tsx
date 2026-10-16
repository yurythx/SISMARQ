import type { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft, Search, X } from "lucide-react";

import { FACETAS_PUBLICAS, ItemResultado, NOME_FACETA, rotuloFaceta } from "@/components/acervo/busca";
import { PublicShell } from "@/components/layout/PublicShell";
import { PublicApiError, publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoBusca } from "@/lib/sismarq/types";

export const metadata: Metadata = {
  title: "Consulta ao acervo",
  description: "Busca no acervo de acesso público do Arquivo Público Municipal de Rondonópolis, com filtros por órgão, classe, espécie, fase, destinação e ano.",
  alternates: { canonical: `${APP_URL}/arquivo/busca` },
};

type Params = Record<string, string | undefined>;
type Props = { searchParams: Promise<Params> };

const CAMPOS = ["q", ...FACETAS_PUBLICAS, "digital"] as const;

/** O endereço da busca com um parâmetro trocado (vazio tira); a página
 * volta para a primeira. */
function com(atual: Params, chave: string, valor?: string): string {
  const u = new URLSearchParams();
  for (const k of CAMPOS) {
    const v = k === chave ? valor : atual[k];
    if (v) u.set(k, v);
  }
  if (chave === "page" && valor) u.set("page", valor);
  const s = u.toString();
  return s ? `/arquivo/busca?${s}` : "/arquivo/busca";
}

/** Consulta pública sem login: só o ostensivo (ADR 038). */
export default async function BuscaPublicaPage({ searchParams }: Props) {
  const sp = await searchParams;
  const pagina = Math.max(1, Number(sp.page) || 1);
  const consulta = new URLSearchParams({ page: String(pagina), page_size: "20" });
  for (const k of CAMPOS) if (sp[k]) consulta.set(k, sp[k]!);
  let res: AcervoBusca | null = null;
  let invalida = false;
  try {
    res = (await publicGet<AcervoBusca>(`acervo/publico/busca?${consulta}`, 30)).data;
  } catch (e) {
    invalida = e instanceof PublicApiError && e.status === 422;
  }
  const ativos = FACETAS_PUBLICAS.filter((k) => sp[k]);

  return (
    <PublicShell>
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-6 py-12">
        <Link href="/arquivo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Arquivo Público
        </Link>
        <h1 className="text-3xl font-bold">Consulta ao acervo</h1>
        <form action="/arquivo/busca" method="get" role="search" className="flex flex-wrap gap-2">
          <label htmlFor="q" className="sr-only">
            Buscar no acervo
          </label>
          <input
            id="q"
            name="q"
            defaultValue={sp.q ?? ""}
            placeholder="Assunto, número, código ou palavra do documento"
            className="min-w-0 flex-1 rounded-md border border-surface-border bg-surface px-3 py-2 text-sm"
          />
          {FACETAS_PUBLICAS.map((k) => (sp[k] ? <input key={k} type="hidden" name={k} value={sp[k]} /> : null))}
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" name="digital" value="true" defaultChecked={sp.digital === "true"} /> Só com documento digital
          </label>
          <button type="submit" className="inline-flex items-center gap-1 rounded-md bg-primary px-4 py-2 text-sm text-primary-foreground">
            <Search size={16} aria-hidden="true" /> Buscar
          </button>
        </form>

        {ativos.length > 0 && (
          <ul className="flex flex-wrap gap-2" aria-label="Filtros aplicados">
            {ativos.map((k) => (
              <li key={k}>
                <Link href={com(sp, k)} className="inline-flex items-center gap-1 rounded-full border border-surface-border px-3 py-1 text-xs hover:bg-surface">
                  {NOME_FACETA[k]}: {rotuloFaceta(k, sp[k]!)} <X size={12} aria-label="tirar o filtro" />
                </Link>
              </li>
            ))}
          </ul>
        )}

        {invalida ? (
          <p role="alert" className="text-sm">
            Consulta inválida. Revise os filtros.
          </p>
        ) : res === null ? (
          <p role="alert" className="text-sm text-muted">
            A consulta ao acervo está indisponível no momento. Tente de novo em instantes.
          </p>
        ) : (
          <div className="grid gap-6 md:grid-cols-[16rem_1fr]">
            <aside aria-label="Filtros" className="flex flex-col gap-4">
              {res.origem === "banco" && (
                <p className="text-xs text-muted">Busca simplificada no momento: os filtros por faceta voltam em instantes.</p>
              )}
              {FACETAS_PUBLICAS.map((k) => {
                const vs = res.facetas[k] ?? [];
                if (vs.length === 0) return null;
                return (
                  <section key={k} aria-labelledby={`f-${k}`}>
                    <h2 id={`f-${k}`} className="mb-1 text-sm font-semibold">
                      {NOME_FACETA[k]}
                    </h2>
                    <ul className="flex flex-col gap-0.5 text-sm">
                      {vs.map((v) => (
                        <li key={v.valor}>
                          <Link href={com(sp, k, sp[k] === v.valor ? undefined : v.valor)} aria-current={sp[k] === v.valor ? "true" : undefined} className="flex justify-between gap-2 hover:underline aria-[current=true]:font-semibold">
                            <span>{rotuloFaceta(k, v.valor)}</span>
                            <span className="text-muted">{v.total}</span>
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </section>
                );
              })}
            </aside>
            <section aria-labelledby="resultados" className="flex flex-col gap-3">
              <h2 id="resultados" className="text-sm text-muted" aria-live="polite">
                {res.meta.total_items.toLocaleString("pt-BR")} unidade(s) encontrada(s)
              </h2>
              {res.items.length === 0 ? (
                <p className="text-sm text-muted">
                  Nada encontrado. Documentos de acesso restrito não aparecem na consulta pública: o pedido de acesso é feito pelo{" "}
                  <Link href="/arquivo#acesso" className="text-primary hover:underline">
                    Serviço de Informação ao Cidadão
                  </Link>
                  .
                </p>
              ) : (
                <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Unidades encontradas">
                  {res.items.map((it) => (
                    <ItemResultado key={it.id} item={it} href={`/arquivo/unidades/${it.id}`} />
                  ))}
                </ul>
              )}
              {res.meta.total_pages > 1 && (
                <nav aria-label="Páginas" className="flex items-center justify-between text-sm">
                  {pagina > 1 ? <Link href={com(sp, "page", String(pagina - 1))}>← Anterior</Link> : <span />}
                  <span className="text-muted">
                    Página {pagina} de {res.meta.total_pages}
                  </span>
                  {pagina < res.meta.total_pages ? <Link href={com(sp, "page", String(pagina + 1))}>Próxima →</Link> : <span />}
                </nav>
              )}
            </section>
          </div>
        )}
      </div>
    </PublicShell>
  );
}
