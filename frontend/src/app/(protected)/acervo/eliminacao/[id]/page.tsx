"use client";

import { ArrowLeft, FileText } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useRef, type FormEvent, type ReactNode } from "react";

import { AREAS, PAPEIS, TIPOS, fmtDia, situacaoListagem, somarDias } from "@/components/acervo/labels";
import { useOrgao } from "@/components/acervo/useOrgao";
import { ConfirmButton } from "@/components/sismarq/ConfirmButton";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button, buttonClass } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoComissao, AcervoListagem, AcervoPedido } from "@/lib/sismarq/types";

const txt = (f: FormData, k: string) => String(f.get(k) ?? "").trim();

function Etapa({ titulo, children }: { titulo: string; children: ReactNode }) {
  return (
    <section aria-label={titulo} className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-4">
      <h2 className="text-base font-semibold">{titulo}</h2>
      {children}
    </section>
  );
}

/** A relação de eliminação e as etapas do rito (arts. 31 a 40). */
export default function RelacaoPage() {
  const { id } = useParams<{ id: string }>();
  const { can, canGlobal } = useSismarq();
  const q = useApiQuery<AcervoListagem>(`v1/acervo/listagens/${id}`);
  const coms = useApiQuery<AcervoComissao[]>("v1/acervo/comissoes");
  const { run, pending } = useAction();
  const l = q.data;
  const { orgao, unidade } = useOrgao(l?.unidade_id);
  const avalia = can("acervo:avaliacao");
  const central = canGlobal("acervo:central");
  const semad = canGlobal("acervo:semad");
  const comissao = coms.data?.find((c) => c.id === l?.comissao_id);

  async function enviar(caminho: string, corpo: unknown, ok: string, metodo: "post" | "put" = "post") {
    const url = `v1/acervo/listagens/${id}${caminho}`;
    const r = await run(() => (metodo === "put" ? apiClient.put<AcervoListagem>(url, corpo) : apiClient.post<AcervoListagem>(url, corpo)), ok);
    if (r) await q.mutate(r.data, { revalidate: false });
  }
  const form = (fn: (f: FormData) => Promise<void>) => (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    void fn(new FormData(e.currentTarget));
  };
  // PUT troca o cabeçalho inteiro: reenvia o que já está gravado.
  const cabecalho = (l: AcervoListagem) => ({
    unidade_id: l.unidade_id,
    responsavel_selecao: l.responsavel_selecao,
    parecer_juridico: l.parecer_juridico,
    contas_exercicios: l.contas_exercicios,
    contas_aprovadas_em: l.contas_aprovadas_em?.slice(0, 10) ?? "",
  });

  const s = l ? situacaoListagem(l.situacao) : null;
  const cancelavel = avalia && !!l && (["RASCUNHO", "APROVADA", "AUTORIZADA", "AGUARDANDO_SEMAD"].includes(l.situacao) || (l.situacao === "EM_EDITAL" && !l.amostra_em));
  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/eliminacao" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Eliminação de documentos
      </Link>
      <DataState loading={q.isLoading} error={q.error} empty={!l}>
        {l && s && (
          <>
            <PageHeader
              eyebrow={[orgao, unidade].filter(Boolean).join(" · ") || "Acervo"}
              title={`Relação de eliminação ${l.numero}`}
              description={s.ajuda}
              actions={
                cancelavel && (
                  <ConfirmButton
                    variant="ghost"
                    onConfirm={() => enviar("/cancelar", undefined, "Relação cancelada.")}
                    confirmLabel="Cancelar a relação"
                    title="Cancelar a relação?"
                    description="As unidades ficam livres para outra relação."
                  >
                    Cancelar relação
                  </ConfirmButton>
                )
              }
            />
            <div className="flex flex-wrap items-center gap-2">
              <Badge tone={s.tone}>{s.label}</Badge>
              {l.rito === "ART_35" && <Badge tone="warning">Art. 35 — fora da TTDD</Badge>}
              {l.rito === "ART_36" && <Badge tone="warning">Art. 36 — órgão sem TTDD</Badge>}
              {comissao && <span className="text-sm text-muted">{comissao.nome}</span>}
              {(l.ata_signum === "AUSENTE" || l.termo_signum === "AUSENTE") && (
                <span className="text-sm text-muted">Ata e termo sem Signum (política 0.6 pendente).</span>
              )}
              {l.devolucao_motivo && l.situacao === "RASCUNHO" && (
                <span className="text-sm text-warning">Devolvida pelo Arquivo Público: {l.devolucao_motivo}</span>
              )}
            </div>

            <nav aria-label="Documentos" className="flex flex-wrap gap-2">
              {[
                { n: 1, nome: "Anexo I — Relação", ok: true },
                { n: 4, nome: "Anexo IV — Ata da CSPAD", ok: !!l.ata_em },
                { n: 3, nome: "Anexo III — Edital de Ciência", ok: !!l.edital_publicado_em },
                { n: 2, nome: "Anexo II — Termo de Eliminação", ok: !!l.termo_eliminado_em },
              ]
                .filter((d) => d.ok)
                .map((d) => (
                  <Link key={d.n} href={`/acervo/eliminacao/${l.id}/anexo/${d.n}`} className={buttonClass("secondary", "sm")}>
                    <FileText size={14} aria-hidden="true" /> {d.nome}
                  </Link>
                ))}
            </nav>

            <table className="w-full border-collapse text-sm">
              <caption className="sr-only">Unidades da relação</caption>
              <thead>
                <tr className="border-b border-surface-border text-left text-xs uppercase">
                  <th scope="col" className="py-2 pr-2">Unidade</th>
                  <th scope="col" className="py-2 pr-2">Código</th>
                  <th scope="col" className="py-2 pr-2">Assunto</th>
                  <th scope="col" className="py-2 pr-2">Datas-limite</th>
                  <th scope="col" className="py-2">Situação</th>
                </tr>
              </thead>
              <tbody>
                {l.itens.map((u) => (
                  <tr key={u.id} className="border-b border-surface-border align-top">
                    <td className="py-2 pr-2 font-mono text-xs">
                      <Link href={`/acervo/${u.id}`} className="hover:underline">
                        {u.numero}
                      </Link>
                    </td>
                    <td className="py-2 pr-2 font-mono text-xs">{u.codigo_classificacao}</td>
                    <td className="py-2 pr-2">
                      {u.assunto} <span className="text-xs text-muted">({u.quantidade} {u.especificacao || TIPOS[u.tipo].toLowerCase()})</span>
                    </td>
                    <td className="py-2 pr-2">
                      {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}
                    </td>
                    <td className="py-2">
                      {u.amostra ? (
                        <Badge tone="success">Amostra — guarda permanente</Badge>
                      ) : u.retirada ? (
                        <Badge tone="neutral">Retirada</Badge>
                      ) : l.situacao === "RASCUNHO" && avalia ? (
                        <Button
                          size="sm"
                          variant="ghost"
                          disabled={pending}
                          aria-label={`Retirar ${u.numero}`}
                          onClick={() => void enviar("", { ...cabecalho(l), retirar: [u.id] }, "Unidade retirada da relação.", "put")}
                        >
                          Retirar
                        </Button>
                      ) : l.situacao === "ELIMINADA" ? (
                        <Badge tone="neutral">Eliminada</Badge>
                      ) : (
                        <Badge tone="warning">A eliminar</Badge>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {l.situacao === "RASCUNHO" && avalia && (
              <>
                <Etapa titulo="Parecer jurídico e contas">
                  <form
                    aria-label="Parecer jurídico e contas"
                    className="flex flex-col gap-3"
                    onSubmit={form((f) =>
                      enviar(
                        "",
                        {
                          ...cabecalho(l),
                          parecer_juridico: txt(f, "parecer_juridico"),
                          contas_exercicios: txt(f, "contas_exercicios"),
                          contas_aprovadas_em: txt(f, "contas_aprovadas_em"),
                        },
                        "Relação atualizada.",
                        "put",
                      ),
                    )}
                  >
                    <Textarea id="parecer_juridico" name="parecer_juridico" label="Parecer jurídico (art. 32, § 4º)" maxLength={4000} rows={3} defaultValue={l.parecer_juridico} />
                    <div className="grid gap-3 sm:grid-cols-2">
                      <Input id="contas_exercicios" name="contas_exercicios" label="Contas do(s) exercício(s) de" maxLength={200} defaultValue={l.contas_exercicios} />
                      <Input id="contas_aprovadas_em" name="contas_aprovadas_em" type="date" label="Contas aprovadas em" defaultValue={l.contas_aprovadas_em?.slice(0, 10) ?? ""} />
                    </div>
                    <div>
                      <Button type="submit" variant="secondary" disabled={pending}>
                        Salvar
                      </Button>
                    </div>
                  </form>
                </Etapa>
                <Etapa titulo="Ata da reunião da CSPAD (Anexo IV)">
                  <form
                    aria-label="Ata da reunião"
                    className="flex flex-col gap-3"
                    onSubmit={form((f) =>
                      enviar(
                        "/aprovar",
                        {
                          ata_numero: txt(f, "ata_numero"),
                          ata_em: txt(f, "ata_em"),
                          ata_local: txt(f, "ata_local"),
                          ata_manifestacoes: txt(f, "ata_manifestacoes"),
                          presentes: f.getAll("presentes").map(String),
                        },
                        "Relação aprovada pela CSPAD.",
                      ),
                    )}
                  >
                    <div className="grid gap-3 sm:grid-cols-3">
                      <Input id="ata_numero" name="ata_numero" label="Ata nº *" required maxLength={40} placeholder="1/2026" />
                      <Input id="ata_em" name="ata_em" type="datetime-local" label="Data e hora *" required />
                      <Input id="ata_local" name="ata_local" label="Local" maxLength={200} />
                    </div>
                    <Textarea id="ata_manifestacoes" name="ata_manifestacoes" label="Manifestações dos presentes" maxLength={8000} rows={3} />
                    <fieldset className="flex flex-col gap-1">
                      <legend className="mb-1 text-sm font-medium">Presentes (quórum: presidente e maioria dos membros)</legend>
                      {(comissao?.membros ?? []).map((m) => (
                        <label key={m.user_id} className="flex items-center gap-2 text-sm">
                          <input type="checkbox" name="presentes" value={m.user_id} defaultChecked />
                          {m.nome || m.user_id} — {PAPEIS[m.papel] ?? m.papel}, {AREAS[m.area] ?? m.area}
                        </label>
                      ))}
                    </fieldset>
                    <div>
                      <Button type="submit" disabled={pending}>
                        Aprovar relação
                      </Button>
                    </div>
                  </form>
                </Etapa>
              </>
            )}

            {l.situacao === "APROVADA" && l.rito === "ART_35" && !l.ccpad_autorizada_em && (
              <Etapa titulo="Autorização da CCPAD (art. 35)">
                <p className="text-sm text-muted">Um membro da Comissão Central autoriza o documento que não consta da TTDD.</p>
                <div>
                  <Button disabled={pending} onClick={() => void enviar("/autorizar-ccpad", undefined, "CCPAD autorizou.")}>
                    Autorizar pela CCPAD
                  </Button>
                </div>
              </Etapa>
            )}

            {l.situacao === "APROVADA" && central && (l.rito !== "ART_35" || !!l.ccpad_autorizada_em) && (
              <Etapa titulo="Autorização do Arquivo Público (arts. 19 e 65)">
                <div>
                  <Button disabled={pending} onClick={() => void enviar("/autorizar", undefined, "Eliminação autorizada.")}>
                    Autorizar eliminação
                  </Button>
                </div>
                <form
                  aria-label="Devolver à CSPAD"
                  className="flex flex-col gap-2"
                  onSubmit={form((f) => enviar("/devolver", { motivo: txt(f, "motivo") }, "Relação devolvida à CSPAD."))}
                >
                  <Textarea id="motivo" name="motivo" label="Motivo da devolução" required maxLength={2000} rows={2} />
                  <div>
                    <Button type="submit" variant="secondary" disabled={pending}>
                      Devolver à CSPAD
                    </Button>
                  </div>
                </form>
              </Etapa>
            )}

            {l.situacao === "AGUARDANDO_SEMAD" && (semad || central) && (
              <Etapa titulo="Aprovação da SEMAD (arts. 35 e 36)">
                {semad && (
                  <div>
                    <Button disabled={pending} onClick={() => void enviar("/autorizar-semad", undefined, "SEMAD aprovou a eliminação.")}>
                      Aprovar pela SEMAD
                    </Button>
                  </div>
                )}
                {central && (
                  <form
                    aria-label="Devolver à comissão"
                    className="flex flex-col gap-2"
                    onSubmit={form((f) => enviar("/devolver", { motivo: txt(f, "motivo") }, "Relação devolvida."))}
                  >
                    <Textarea id="motivo_semad" name="motivo" label="Motivo da devolução" required maxLength={2000} rows={2} />
                    <div>
                      <Button type="submit" variant="secondary" disabled={pending}>
                        Devolver
                      </Button>
                    </div>
                  </form>
                )}
              </Etapa>
            )}

            {l.situacao === "AUTORIZADA" && avalia && (
              <Etapa titulo="Edital de Ciência de Eliminação (Anexo III)">
                <p className="text-sm text-muted">Publique o edital no Diário Oficial do Município; o prazo para os interessados é de 30 dias (art. 37).</p>
                <form
                  aria-label="Edital de ciência"
                  className="grid gap-3 sm:grid-cols-3"
                  onSubmit={form((f) =>
                    enviar(
                      "/edital",
                      { edital_numero: txt(f, "edital_numero"), edital_publicado_em: txt(f, "edital_publicado_em"), edital_diario: txt(f, "edital_diario") },
                      "Edital registrado: prazo de 30 dias aberto.",
                    ),
                  )}
                >
                  <Input id="edital_numero" name="edital_numero" label="Edital nº *" required maxLength={40} />
                  <Input id="edital_publicado_em" name="edital_publicado_em" type="date" label="Publicado em *" required />
                  <Input id="edital_diario" name="edital_diario" label="Diário Oficial (edição, página)" maxLength={200} placeholder="DIORONDON-e nº …, p. …" />
                  <div>
                    <Button type="submit" disabled={pending}>
                      Registrar publicação
                    </Button>
                  </div>
                </form>
              </Etapa>
            )}

            {l.situacao === "EM_EDITAL" && l.edital_publicado_em && (
              <>
                <p className="text-sm">
                  Edital nº {l.edital_numero} publicado em {fmtDia(l.edital_publicado_em)} · prazo para pedidos até{" "}
                  <strong>{fmtDia(somarDias(l.edital_publicado_em, 30))}</strong> · eliminação a partir de{" "}
                  <strong>{fmtDia(somarDias(l.edital_publicado_em, 31))}</strong>.
                </p>
                <Etapa titulo="Pedidos dos interessados (art. 37)">
                  {l.pedidos.length === 0 && <p className="text-sm text-muted">Nenhum pedido registrado.</p>}
                  <ul className="flex flex-col gap-3">
                    {l.pedidos.map((p) => (
                      <PedidoItem key={p.id} p={p} l={l} avalia={avalia} pending={pending} enviar={enviar} />
                    ))}
                  </ul>
                  {avalia && (
                    <form
                      aria-label="Registrar pedido"
                      className="grid gap-3 border-t border-surface-border pt-3 sm:grid-cols-2"
                      onSubmit={form((f) =>
                        enviar(
                          "/pedidos",
                          {
                            requerente: txt(f, "requerente"),
                            pedido: txt(f, "pedido"),
                            recebido_em: txt(f, "recebido_em"),
                            unidade_arquivamento_id: txt(f, "unidade_arquivamento_id") || null,
                          },
                          "Pedido registrado.",
                        ),
                      )}
                    >
                      <Input id="requerente" name="requerente" label="Requerente *" required maxLength={200} />
                      <Input id="recebido_em" name="recebido_em" type="date" label="Recebido em *" required />
                      <Select
                        id="unidade_arquivamento_id"
                        name="unidade_arquivamento_id"
                        label="Unidade"
                        placeholder="Toda a relação"
                        options={l.itens.filter((u) => !u.retirada).map((u) => ({ value: u.id, label: `${u.numero} — ${u.assunto}` }))}
                      />
                      <Textarea id="pedido" name="pedido" label="Pedido *" required maxLength={4000} rows={2} />
                      <div>
                        <Button type="submit" variant="secondary" disabled={pending}>
                          Registrar pedido
                        </Button>
                      </div>
                    </form>
                  )}
                </Etapa>
                {avalia && !l.amostra_em && (
                  <Etapa titulo="Amostragem (art. 39)">
                    <form
                      aria-label="Amostragem"
                      className="flex flex-col gap-3"
                      onSubmit={form((f) =>
                        enviar("/amostragem", { criterios: txt(f, "criterios"), percentual: Number(f.get("percentual")) }, "Amostra sorteada e levada à guarda permanente."),
                      )}
                    >
                      <Textarea id="criterios" name="criterios" label="Critérios da CSPAD *" required maxLength={2000} rows={2} />
                      <Input id="percentual" name="percentual" type="number" min={0.1} max={100} step={0.1} defaultValue={10} label="Percentual sorteado *" required />
                      <p className="text-xs text-muted">A amostra (ao menos uma unidade) passa à guarda permanente e não pode mais ser eliminada (arts. 20 e 42).</p>
                      <div>
                        <Button type="submit" disabled={pending}>
                          Sortear amostra
                        </Button>
                      </div>
                    </form>
                  </Etapa>
                )}
                {avalia && l.amostra_em && (
                  <Etapa titulo="Termo de Eliminação (Anexo II)">
                    <p className="text-sm text-muted">
                      Amostra de {l.amostra_percentual}% sorteada em {fmtDia(l.amostra_em)} ({l.amostra_criterios}). O termo só pode ser registrado depois do prazo
                      do edital e com os pedidos decididos.
                    </p>
                    <form
                      aria-label="Termo de eliminação"
                      className="grid gap-3 sm:grid-cols-2"
                      onSubmit={form((f) =>
                        enviar(
                          "/termo",
                          {
                            termo_eliminado_em: txt(f, "termo_eliminado_em"),
                            termo_acompanhamento: txt(f, "termo_acompanhamento"),
                            termo_destinacao_papel: txt(f, "termo_destinacao_papel"),
                          },
                          "Termo de Eliminação registrado.",
                        ),
                      )}
                    >
                      <Input id="termo_eliminado_em" name="termo_eliminado_em" type="date" label="Eliminado em *" required />
                      <Input id="termo_acompanhamento" name="termo_acompanhamento" label="Servidor que acompanhou *" required maxLength={200} />
                      <Input
                        id="termo_destinacao_papel"
                        name="termo_destinacao_papel"
                        label="Destinação do papel (art. 40)"
                        maxLength={500}
                        placeholder="Fragmentação; doação a cooperativa de catadores"
                      />
                      <div>
                        <Button type="submit" disabled={pending}>
                          Registrar termo
                        </Button>
                      </div>
                    </form>
                  </Etapa>
                )}
              </>
            )}
          </>
        )}
      </DataState>
    </div>
  );
}

function PedidoItem({
  p,
  l,
  avalia,
  pending,
  enviar,
}: {
  p: AcervoPedido;
  l: AcervoListagem;
  avalia: boolean;
  pending: boolean;
  enviar: (caminho: string, corpo: unknown, ok: string) => Promise<void>;
}) {
  const unidade = l.itens.find((u) => u.id === p.unidade_arquivamento_id);
  const tone = p.situacao === "DEFERIDO" ? "success" : p.situacao === "INDEFERIDO" ? "neutral" : "warning";
  const deferir = useRef(false);
  return (
    <li className="flex flex-col gap-2 rounded-lg border border-surface-border p-3 text-sm">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span>
          <strong>{p.requerente}</strong> · recebido em {fmtDia(p.recebido_em)}
          {unidade && ` · ${unidade.numero}`}
        </span>
        <Badge tone={tone}>{p.situacao === "PENDENTE" ? "Pendente" : p.situacao === "DEFERIDO" ? "Deferido" : "Indeferido"}</Badge>
      </div>
      <p>{p.pedido}</p>
      {p.decisao && <p className="text-muted">Decisão: {p.decisao}</p>}
      {p.situacao === "PENDENTE" && avalia && (
        <form
          aria-label={`Decidir pedido de ${p.requerente}`}
          className="flex flex-col gap-2"
          onSubmit={(e) => {
            e.preventDefault();
            const f = new FormData(e.currentTarget);
            const deferido = deferir.current;
            void enviar(
              `/pedidos/${p.id}/decidir`,
              { deferido, retirar: deferido && f.get("retirar") === "on", decisao: txt(f, "decisao") },
              deferido ? "Pedido deferido." : "Pedido indeferido.",
            );
          }}
        >
          <Textarea id={`decisao-${p.id}`} name="decisao" label="Fundamentação *" required maxLength={4000} rows={2} />
          {p.unidade_arquivamento_id && (
            <label className="flex items-center gap-2">
              <input type="checkbox" name="retirar" defaultChecked /> Deferido, retirar a unidade da relação
            </label>
          )}
          <div className="flex gap-2">
            <Button type="submit" size="sm" disabled={pending} onClick={() => (deferir.current = true)}>
              Deferir
            </Button>
            <Button type="submit" size="sm" variant="secondary" disabled={pending} onClick={() => (deferir.current = false)}>
              Indeferir
            </Button>
          </div>
        </form>
      )}
    </li>
  );
}
