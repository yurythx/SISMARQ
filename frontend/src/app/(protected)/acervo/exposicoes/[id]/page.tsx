"use client";

import { ArrowDown, ArrowLeft, ArrowUp, ExternalLink, Trash2 } from "lucide-react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useState } from "react";

import { fmtDia } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { apiClient } from "@/lib/api/client";
import { useApiPage, useApiQuery } from "@/lib/api/swr";
import type { AcervoDigital, AcervoExposicao } from "@/lib/sismarq/types";

type Item = { digital_id: string; legenda: string; rotulo: string; publicado: boolean };

const campo = "w-full rounded-md border border-surface-border bg-surface px-3 py-2 text-sm";

/** Criar (id "nova") ou editar a exposição: texto, documentos e legendas. */
export default function ExposicaoEditorPage() {
  const { id } = useParams<{ id: string }>();
  const nova = id === "nova";
  const q = useApiQuery<AcervoExposicao>(nova ? null : `v1/acervo/exposicoes/${id}`);
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/exposicoes" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Exposições
      </Link>
      {nova ? (
        <Editor />
      ) : (
        <DataState loading={q.isLoading} error={q.error} empty={!q.data}>
          {q.data && <Editor key={q.data.atualizado_em} inicial={q.data} onSalva={(x) => void q.mutate(x, { revalidate: false })} />}
        </DataState>
      )}
    </div>
  );
}

function Editor({ inicial, onSalva }: { inicial?: AcervoExposicao; onSalva?: (x: AcervoExposicao) => void }) {
  const router = useRouter();
  const { run, pending } = useAction();
  const [titulo, setTitulo] = useState(inicial?.titulo ?? "");
  const [resumo, setResumo] = useState(inicial?.resumo ?? "");
  const [texto, setTexto] = useState(inicial?.texto ?? "");
  const [itens, setItens] = useState<Item[]>(
    (inicial?.itens ?? []).map((it) => ({ digital_id: it.digital_id, legenda: it.legenda, rotulo: it.descricao || it.numero, publicado: it.publicado })),
  );
  const [escolha, setEscolha] = useState("");
  const liberados = useApiPage<AcervoDigital>("v1/acervo/digitais?publicados=true&page_size=100");
  const disponiveis = (liberados.data?.items ?? []).filter((d) => !itens.some((it) => it.digital_id === d.id));

  function mover(i: number, delta: number) {
    setItens((xs) => {
      const n = [...xs];
      const [x] = n.splice(i, 1);
      n.splice(i + delta, 0, x!);
      return n;
    });
  }

  function adicionar() {
    const d = disponiveis.find((x) => x.id === escolha);
    if (!d) return;
    setItens((xs) => [...xs, { digital_id: d.id, legenda: "", rotulo: d.descricao || d.numero, publicado: true }]);
    setEscolha("");
  }

  async function salvar(e: React.FormEvent) {
    e.preventDefault();
    const corpo = { titulo, resumo, texto, itens: itens.map(({ digital_id, legenda }) => ({ digital_id, legenda })) };
    const r = await run(
      () => (inicial ? apiClient.put<AcervoExposicao>(`v1/acervo/exposicoes/${inicial.id}`, corpo) : apiClient.post<AcervoExposicao>("v1/acervo/exposicoes", corpo)),
      "Exposição salva.",
    );
    if (!r?.data) return;
    if (inicial) onSalva?.(r.data);
    else router.push(`/acervo/exposicoes/${r.data.id}`);
  }

  async function publicacao(publicar: boolean) {
    const caminho = `v1/acervo/exposicoes/${inicial!.id}/publicacao`;
    const r = await run(
      () => (publicar ? apiClient.post<AcervoExposicao>(caminho) : apiClient.delete<AcervoExposicao>(caminho)),
      publicar ? "Exposição publicada." : "Exposição recolhida.",
    );
    if (r?.data) onSalva?.(r.data);
  }

  return (
    <>
      <PageHeader
        eyebrow="Exposição virtual"
        title={inicial ? inicial.titulo : "Nova exposição"}
        description={inicial ? `Endereço: /arquivo/exposicoes/${inicial.slug}` : "O endereço sai do título e não muda depois."}
        actions={
          inicial && (
            <>
              {inicial.publicada_em ? (
                <>
                  <Link href={`/arquivo/exposicoes/${inicial.slug}`} className="inline-flex items-center gap-1 text-sm text-primary hover:underline">
                    <ExternalLink size={14} aria-hidden="true" /> Ver a visita
                  </Link>
                  <Button variant="secondary" disabled={pending} onClick={() => void publicacao(false)}>
                    Recolher
                  </Button>
                </>
              ) : (
                <Button disabled={pending} onClick={() => void publicacao(true)}>
                  Publicar
                </Button>
              )}
            </>
          )
        }
      />
      {inicial && (
        <div>{inicial.publicada_em ? <Badge tone="success">Publicada em {fmtDia(inicial.publicada_em)}</Badge> : <Badge>Rascunho</Badge>}</div>
      )}
      <form onSubmit={(e) => void salvar(e)} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1 text-sm">
          Título
          <input required maxLength={200} value={titulo} onChange={(e) => setTitulo(e.target.value)} className={campo} />
        </label>
        <label className="flex flex-col gap-1 text-sm">
          Resumo
          <input maxLength={500} value={resumo} onChange={(e) => setResumo(e.target.value)} className={campo} />
        </label>
        <label className="flex flex-col gap-1 text-sm">
          Texto de apresentação (Markdown)
          <textarea rows={8} maxLength={20000} value={texto} onChange={(e) => setTexto(e.target.value)} className={campo} />
        </label>

        <fieldset className="flex flex-col gap-2">
          <legend className="mb-1 text-sm font-semibold">Documentos ({itens.length})</legend>
          {itens.length === 0 && <p className="text-sm text-muted">Nenhum documento ainda. Só os liberados à consulta pública aparecem na visita.</p>}
          <ol className="flex flex-col gap-2">
            {itens.map((it, i) => (
              <li key={it.digital_id} className="flex flex-col gap-2 rounded-lg border border-surface-border p-3">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="text-sm font-medium">
                    {i + 1}. {it.rotulo}
                  </span>
                  <div className="flex items-center gap-1">
                    {!it.publicado && <Badge tone="warning">Não liberado — não aparece na visita</Badge>}
                    <Button type="button" variant="ghost" size="sm" aria-label={`Subir ${it.rotulo}`} disabled={i === 0} onClick={() => mover(i, -1)}>
                      <ArrowUp size={14} aria-hidden="true" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      aria-label={`Descer ${it.rotulo}`}
                      disabled={i === itens.length - 1}
                      onClick={() => mover(i, 1)}
                    >
                      <ArrowDown size={14} aria-hidden="true" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      aria-label={`Tirar ${it.rotulo}`}
                      onClick={() => setItens((xs) => xs.filter((x) => x.digital_id !== it.digital_id))}
                    >
                      <Trash2 size={14} aria-hidden="true" />
                    </Button>
                  </div>
                </div>
                <label className="flex flex-col gap-1 text-xs">
                  Legenda
                  <input
                    maxLength={1000}
                    value={it.legenda}
                    onChange={(e) => setItens((xs) => xs.map((x) => (x.digital_id === it.digital_id ? { ...x, legenda: e.target.value } : x)))}
                    className={campo}
                  />
                </label>
              </li>
            ))}
          </ol>
          <div className="flex flex-wrap items-end gap-2">
            <label className="flex flex-1 flex-col gap-1 text-sm">
              Adicionar documento liberado
              <select value={escolha} onChange={(e) => setEscolha(e.target.value)} className={campo}>
                <option value="">Escolha…</option>
                {disponiveis.map((d) => (
                  <option key={d.id} value={d.id}>
                    {d.numero} — {d.descricao || d.codigo_classificacao}
                  </option>
                ))}
              </select>
            </label>
            <Button type="button" variant="secondary" disabled={!escolha} onClick={adicionar}>
              Adicionar
            </Button>
          </div>
        </fieldset>
        <div>
          <Button type="submit" disabled={pending}>
            Salvar
          </Button>
        </div>
      </form>
    </>
  );
}
