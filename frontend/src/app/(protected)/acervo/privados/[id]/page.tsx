"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useRef, useState, type FormEvent, type ReactNode } from "react";

import { DECISOES_ALIENACAO, GUARDAS, fmtDia, situacaoPrivado } from "@/components/acervo/labels";
import { CamposPrivado, corpoPrivado } from "@/components/acervo/privado";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { fmtDateTime, useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoAlienacao, AcervoPrivado } from "@/lib/sismarq/types";

const txt = (f: FormData, k: string) => String(f.get(k) ?? "").trim();

/** Para onde a guarda pode ir (§§ 4º e 5º): do titular, doação ou
 * depósito; do depósito, de volta ao titular ou doação; a doação é final. */
const DESTINOS: Record<string, string[]> = {
  TITULAR: ["DEPOSITADO", "DOADO"],
  DEPOSITADO: ["TITULAR", "DOADO"],
  DOADO: [],
};

function Etapa({ titulo, children }: { titulo: string; children: ReactNode }) {
  return (
    <section aria-label={titulo} className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4">
      <h2 className="text-base font-semibold">{titulo}</h2>
      {children}
    </section>
  );
}

function Campo({ nome, children }: { nome: string; children: ReactNode }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-muted">{nome}</dt>
      <dd className="mt-0.5 text-sm">{children || "—"}</dd>
    </div>
  );
}

/** O arquivo privado: cadastro, parecer, decreto, guarda e o direito de
 * preferência (arts. 55 a 58). */
export default function PrivadoPage() {
  const { id } = useParams<{ id: string }>();
  const q = useApiQuery<AcervoPrivado>(`v1/acervo/privados/${id}`);
  const { run, pending } = useAction();
  const [editando, setEditando] = useState(false);
  const p = q.data;

  async function enviar(caminho: string, corpo: unknown, ok: string, metodo: "post" | "put" = "post") {
    const url = `v1/acervo/privados/${id}${caminho}`;
    const r = await run(() => (metodo === "put" ? apiClient.put<AcervoPrivado>(url, corpo) : apiClient.post<AcervoPrivado>(url, corpo)), ok);
    if (r) {
      await q.mutate(r.data, { revalidate: false });
      setEditando(false);
    }
  }
  const form = (fn: (f: FormData) => Promise<void>) => (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void fn(new FormData(e.currentTarget));
  };
  const s = p ? situacaoPrivado(p.situacao) : null;
  const pendente = p?.alienacoes.find((a) => a.decisao === "PENDENTE" && !a.vencida);

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/privados" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Arquivos privados
      </Link>
      <DataState loading={q.isLoading} error={q.error} empty={!p}>
        {p && s && (
          <>
            <PageHeader
              eyebrow={`Arquivo privado · ${p.numero}`}
              title={p.titular}
              description={p.descricao}
              actions={!editando && <Button variant="secondary" onClick={() => setEditando(true)}>Alterar cadastro</Button>}
            />
            <div className="flex flex-wrap items-center gap-2">
              <Badge tone={s.tone}>{s.label}</Badge>
              {p.situacao === "DECLARADO" && <Badge>{GUARDAS[p.guarda] ?? p.guarda}</Badge>}
              {p.servico_publico && <Badge tone="info">Serviço público (§ 3º)</Badge>}
            </div>

            {editando ? (
              <form
                aria-label="Alterar cadastro"
                className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4"
                onSubmit={form((f) => enviar("", corpoPrivado(f), "Cadastro alterado.", "put"))}
              >
                <CamposPrivado p={p} />
                <div className="flex gap-2">
                  <Button type="submit" disabled={pending}>
                    Salvar
                  </Button>
                  <Button type="button" variant="ghost" onClick={() => setEditando(false)}>
                    Cancelar
                  </Button>
                </div>
              </form>
            ) : (
              <dl className="grid gap-3 rounded-xl border border-surface-border bg-surface p-4 sm:grid-cols-3">
                <Campo nome="Titular">
                  {p.titular} ({p.tipo_titular === "PJ" ? "pessoa jurídica" : "pessoa física"}){p.documento && ` · ${p.documento}`}
                </Campo>
                <Campo nome="Contato">{p.contato}</Campo>
                <Campo nome="Datas-limite">{p.datas_limite}</Campo>
                <Campo nome="Mensuração">{p.mensuracao}</Campo>
                <Campo nome="Onde está">{p.local_guarda}</Campo>
                <Campo nome="Acesso do público (§ 2º)">
                  {p.acesso_autorizado ? `Autorizado pelo titular${p.acesso_condicoes ? ` — ${p.acesso_condicoes}` : ""}` : "Não autorizado"}
                </Campo>
                <Campo nome="Convênio (art. 57)">{p.convenio}</Campo>
                {p.parecer_em && (
                  <Campo nome="Parecer da Comissão Especial">
                    {p.parecer_favoravel ? "Favorável" : "Contrário"} em {fmtDia(p.parecer_em)} — {p.parecer_texto}
                  </Campo>
                )}
                {p.decreto_em && (
                  <Campo nome="Decreto">
                    {p.decreto_numero}, de {fmtDia(p.decreto_em)}
                    {p.decreto_diario && ` (${p.decreto_diario})`}
                  </Campo>
                )}
                {p.guarda_em && <Campo nome="Guarda desde">{`${fmtDia(p.guarda_em)} — ${p.guarda_termo}`}</Campo>}
              </dl>
            )}

            {p.situacao === "CADASTRADO" && (
              <Etapa titulo="Parecer da Comissão Especial (art. 56, § 1º)">
                <form
                  aria-label="Parecer"
                  className="flex flex-col gap-3"
                  onSubmit={form((f) =>
                    enviar(
                      "/parecer",
                      {
                        parecer_comissao: txt(f, "parecer_comissao"),
                        parecer_texto: txt(f, "parecer_texto"),
                        parecer_favoravel: f.get("parecer_favoravel") === "on",
                        parecer_em: txt(f, "parecer_em"),
                      },
                      "Parecer registrado.",
                    ),
                  )}
                >
                  <Input id="parecer_comissao" name="parecer_comissao" label="Comissão Especial (composição) *" required maxLength={2000} />
                  <Textarea id="parecer_texto" name="parecer_texto" label="Parecer com a avaliação técnica *" required rows={3} maxLength={8000} />
                  <div className="flex flex-wrap items-end gap-4">
                    <Input id="parecer_em" name="parecer_em" type="date" label="Data *" required />
                    <label className="flex items-center gap-2 text-sm">
                      <input type="checkbox" name="parecer_favoravel" /> Favorável à declaração
                    </label>
                    <Button type="submit" disabled={pending}>
                      Registrar parecer
                    </Button>
                  </div>
                </form>
              </Etapa>
            )}

            {p.situacao === "COM_PARECER" && (
              <Etapa titulo="Decreto de declaração (art. 56)">
                <form
                  aria-label="Decreto"
                  className="flex flex-wrap items-end gap-3"
                  onSubmit={form((f) =>
                    enviar(
                      "/declarar",
                      { decreto_numero: txt(f, "decreto_numero"), decreto_em: txt(f, "decreto_em"), decreto_diario: txt(f, "decreto_diario") },
                      "Arquivo declarado de interesse público e social.",
                    ),
                  )}
                >
                  <Input id="decreto_numero" name="decreto_numero" label="Decreto nº *" required maxLength={40} />
                  <Input id="decreto_em" name="decreto_em" type="date" label="Data *" required />
                  <Input id="decreto_diario" name="decreto_diario" label="Diário Oficial" maxLength={200} />
                  <Button type="submit" disabled={pending}>
                    Registrar decreto
                  </Button>
                </form>
              </Etapa>
            )}

            {p.situacao === "DECLARADO" && (DESTINOS[p.guarda] ?? []).length > 0 && (
              <Etapa titulo="Guarda (art. 56, §§ 4º e 5º)">
                <p className="text-sm text-muted">A declaração não transfere a guarda: o titular segue responsável, salvo doação ou depósito no Arquivo Público.</p>
                <form
                  aria-label="Guarda"
                  className="flex flex-wrap items-end gap-3"
                  onSubmit={form((f) => enviar("/guarda", { guarda: txt(f, "guarda"), termo: txt(f, "termo"), em: txt(f, "em") }, "Guarda registrada."))}
                >
                  <Select
                    id="guarda"
                    name="guarda"
                    label="Nova guarda *"
                    options={(DESTINOS[p.guarda] ?? []).map((g) => ({ value: g, label: GUARDAS[g] ?? g }))}
                  />
                  <div className="min-w-64 flex-1">
                    <Input id="termo" name="termo" label="Termo *" required maxLength={1000} placeholder="Termo de doação ou de depósito nº …" />
                  </div>
                  <Input id="em" name="em" type="date" label="Data *" required />
                  <Button type="submit" disabled={pending}>
                    Registrar
                  </Button>
                </form>
              </Etapa>
            )}

            {p.situacao === "DECLARADO" && (
              <Etapa titulo="Direito de preferência (art. 58)">
                <p className="text-sm text-muted">Na venda, o titular notifica o município, que tem 60 dias para manifestar interesse.</p>
                {p.alienacoes.length > 0 && (
                  <ul className="flex flex-col gap-3">
                    {p.alienacoes.map((a) => (
                      <AlienacaoItem key={a.id} a={a} pending={pending} enviar={enviar} />
                    ))}
                  </ul>
                )}
                {!pendente && p.guarda !== "DOADO" && (
                  <form
                    aria-label="Notificação de alienação"
                    className="grid gap-3 border-t border-surface-border pt-3 sm:grid-cols-3"
                    onSubmit={form((f) =>
                      enviar(
                        "/alienacoes",
                        { notificado_em: txt(f, "notificado_em"), adquirente: txt(f, "adquirente"), condicoes: txt(f, "condicoes") },
                        "Notificação registrada: prazo de 60 dias aberto.",
                      ),
                    )}
                  >
                    <Input id="notificado_em" name="notificado_em" type="date" label="Notificado em *" required />
                    <Input id="adquirente" name="adquirente" label="Pretendente" maxLength={200} />
                    <Input id="condicoes" name="condicoes" label="Condições (preço, prazo)" maxLength={2000} />
                    <div>
                      <Button type="submit" variant="secondary" disabled={pending}>
                        Registrar notificação
                      </Button>
                    </div>
                  </form>
                )}
              </Etapa>
            )}

            <section aria-labelledby="historico" className="flex flex-col gap-2">
              <h2 id="historico" className="text-base font-semibold">
                Histórico
              </h2>
              <ol className="flex flex-col gap-1 text-sm">
                {p.eventos.map((ev, i) => (
                  <li key={i}>
                    <span className="text-xs text-muted">{fmtDateTime(ev.em)}</span> — {ev.descricao}
                  </li>
                ))}
              </ol>
            </section>
          </>
        )}
      </DataState>
    </div>
  );
}

function AlienacaoItem({
  a,
  pending,
  enviar,
}: {
  a: AcervoAlienacao;
  pending: boolean;
  enviar: (caminho: string, corpo: unknown, ok: string) => Promise<void>;
}) {
  const interesse = useRef(true);
  const tone = a.decisao === "INTERESSE" ? "success" : a.decisao === "SEM_INTERESSE" || a.vencida ? "neutral" : "warning";
  return (
    <li className="flex flex-col gap-2 rounded-lg border border-surface-border p-3 text-sm">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span>
          Notificado em {fmtDia(a.notificado_em)}
          {a.adquirente && ` · pretendente: ${a.adquirente}`} · prazo até <strong>{fmtDia(a.prazo_final)}</strong>
        </span>
        <Badge tone={tone}>{a.vencida ? "Prazo vencido — a preferência caducou" : (DECISOES_ALIENACAO[a.decisao] ?? a.decisao)}</Badge>
      </div>
      {a.condicoes && <p className="text-muted">Condições: {a.condicoes}</p>}
      {a.fundamento && <p className="text-muted">Fundamento: {a.fundamento}</p>}
      {a.decisao === "PENDENTE" && !a.vencida && (
        <form
          aria-label="Manifestação do município"
          className="flex flex-col gap-2"
          onSubmit={(e) => {
            e.preventDefault();
            const f = new FormData(e.currentTarget);
            void enviar(
              `/alienacoes/${a.id}/decidir`,
              { decisao: interesse.current ? "INTERESSE" : "SEM_INTERESSE", fundamento: txt(f, "fundamento") },
              "Manifestação registrada.",
            );
          }}
        >
          <Textarea id={`fundamento-${a.id}`} name="fundamento" label="Fundamento *" required rows={2} maxLength={2000} />
          <div className="flex gap-2">
            <Button type="submit" size="sm" disabled={pending} onClick={() => (interesse.current = true)}>
              Tem interesse
            </Button>
            <Button type="submit" size="sm" variant="secondary" disabled={pending} onClick={() => (interesse.current = false)}>
              Não tem interesse
            </Button>
          </div>
        </form>
      )}
    </li>
  );
}
