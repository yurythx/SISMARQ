import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft, Download, FileText } from "lucide-react";

import { rotuloFaceta } from "@/components/acervo/busca";
import { fmtBytes, fmtDia } from "@/components/acervo/labels";
import { VisualizadorPDF } from "@/components/acervo/VisualizadorPDF";
import { PublicShell } from "@/components/layout/PublicShell";
import { PublicApiError, publicGet } from "@/lib/api/publicServer";
import { APP_URL } from "@/lib/env";
import type { AcervoFichaPublica } from "@/lib/sismarq/types";

type Props = { params: Promise<{ id: string }>; searchParams: Promise<{ doc?: string }> };

async function carregar(id: string): Promise<AcervoFichaPublica | null> {
  try {
    return (await publicGet<AcervoFichaPublica>(`acervo/publico/unidades/${encodeURIComponent(id)}`, 60)).data;
  } catch (e) {
    if (e instanceof PublicApiError && (e.status === 404 || e.status === 400)) return null;
    throw e;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const f = await carregar((await params).id);
  if (!f) return { title: "Unidade não encontrada" };
  return {
    title: f.assunto || f.descritor || f.numero,
    description: `${f.numero} · ${f.codigo_classificacao} ${f.descritor} · ${f.orgao}`,
    alternates: { canonical: `${APP_URL}/arquivo/unidades/${f.id}` },
  };
}

function Campo({ nome, children }: { nome: string; children: React.ReactNode }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-muted">{nome}</dt>
      <dd className="mt-0.5 text-sm">{children || "—"}</dd>
    </div>
  );
}

/** A ficha ostensiva da unidade e o documento liberado em fluxo (sem login). */
export default async function FichaPublicaPage({ params, searchParams }: Props) {
  const [{ id }, { doc }] = await Promise.all([params, searchParams]);
  const f = await carregar(id);
  if (!f) notFound();
  const atual = f.digitais.find((d) => d.id === doc) ?? f.digitais[0];
  const arquivo = atual ? `/api/public/v1/acervo/publico/digitais/${atual.id}/arquivo` : "";

  return (
    <PublicShell>
      <article className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-6 py-12">
        <Link href="/arquivo/busca" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Consulta ao acervo
        </Link>
        <header className="flex flex-col gap-1">
          <p className="dateline">
            {f.orgao} · {f.recolhida ? "sob a guarda do Arquivo Público" : "no órgão produtor"}
          </p>
          <h1 className="text-3xl font-bold">{f.assunto || f.descritor || f.numero}</h1>
          <p className="font-mono text-sm text-muted">{f.numero}</p>
        </header>

        <dl className="grid gap-3 rounded-xl border border-surface-border bg-surface p-4 sm:grid-cols-3">
          <Campo nome="Classificação (TTDD)">
            {f.codigo_classificacao} {f.descritor && `— ${f.descritor}`}
          </Campo>
          <Campo nome="Órgão produtor">{f.orgao}</Campo>
          <Campo nome="Datas-limite">
            {fmtDia(f.data_inicial)} a {fmtDia(f.data_final)}
          </Campo>
          <Campo nome="Espécie">
            {rotuloFaceta("especie", f.tipo)} · {f.quantidade} {f.especificacao}
          </Campo>
          <Campo nome="Fase">{rotuloFaceta("fase", f.fase)}</Campo>
          <Campo nome="Destinação">{rotuloFaceta("destinacao", f.destinacao)}</Campo>
        </dl>

        {f.digitais.length === 0 ? (
          <p className="text-sm text-muted">
            Esta unidade ainda não tem documento digitalizado liberado para consulta. A consulta ao original pode ser agendada com o Arquivo
            Público.
          </p>
        ) : (
          <section aria-labelledby="documentos" className="flex flex-col gap-3">
            <h2 id="documentos" className="text-xl font-semibold">
              Documentos digitalizados
            </h2>
            {f.digitais.length > 1 && (
              <ul className="flex flex-wrap gap-2">
                {f.digitais.map((d) => (
                  <li key={d.id}>
                    <Link
                      href={`/arquivo/unidades/${f.id}?doc=${d.id}`}
                      aria-current={d.id === atual?.id ? "true" : undefined}
                      className="inline-flex items-center gap-1 rounded-md border border-surface-border px-3 py-1 text-sm aria-[current=true]:border-primary aria-[current=true]:font-semibold"
                    >
                      <FileText size={14} aria-hidden="true" /> {d.descricao || d.numero}
                    </Link>
                  </li>
                ))}
              </ul>
            )}
            {atual && (
              <>
                <p className="flex flex-wrap items-center gap-3 text-xs text-muted">
                  <span>
                    {atual.numero} · {atual.conformidade} · {atual.paginas} página(s) · {fmtBytes(atual.pdfa_bytes)}
                  </span>
                  <a href={arquivo} download={`${atual.numero}.pdf`} className="inline-flex items-center gap-1 text-primary hover:underline">
                    <Download size={14} aria-hidden="true" /> Baixar
                  </a>
                </p>
                <VisualizadorPDF url={arquivo} titulo={atual.numero} />
                <p className="break-all text-xs text-muted">
                  Representante digital com valor probatório (Decreto Federal 10.278/2020). SHA-256 para conferência:{" "}
                  <span className="font-mono">{atual.pdfa_sha256}</span>
                </p>
              </>
            )}
          </section>
        )}
      </article>
    </PublicShell>
  );
}
