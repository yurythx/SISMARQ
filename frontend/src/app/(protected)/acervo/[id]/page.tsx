"use client";

import { ArrowLeft, ArrowRightLeft, Gavel, HandHelping, MapPin, Pencil, QrCode } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useState, type FormEvent } from "react";

import { DigitaisDaUnidade } from "@/components/acervo/DigitaisDaUnidade";
import { DESTINACOES, FASES, NIVEIS, TIPOS, fmtDia, situacao } from "@/components/acervo/labels";
import { UserPicker, type PickedUser } from "@/components/iam/UserPicker";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { fmtDateTime, useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button, buttonClass } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Dialog } from "@/components/ui/Dialog";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoDetalhe, AcervoLocal } from "@/lib/sismarq/types";

function Campo({ nome, children }: { nome: string; children: React.ReactNode }) {
  return (
    <div>
      <dt className="text-xs font-medium uppercase tracking-wide text-muted">{nome}</dt>
      <dd className="mt-0.5 text-sm text-foreground">{children || "—"}</dd>
    </div>
  );
}

type Acao = "transferir" | "localizar" | "bloquear" | "emprestar" | { liberar: string } | { devolver: string } | null;

/** Ficha da unidade de arquivamento (ADR 032). */
export default function UnidadePage() {
  const { id } = useParams<{ id: string }>();
  const ficha = useApiQuery<AcervoDetalhe>(`v1/acervo/unidades/${id}`);
  const locais = useApiQuery<AcervoLocal[]>("v1/acervo/locais");
  const { run, pending } = useAction();
  const [acao, setAcao] = useState<Acao>(null);
  const [solicitante, setSolicitante] = useState<PickedUser[]>([]);
  const u = ficha.data;

  async function enviar(fn: () => Promise<{ data: AcervoDetalhe }>, ok: string) {
    const r = await run(fn, ok);
    if (r) {
      await ficha.mutate(r.data, { revalidate: false });
      setAcao(null);
    }
  }

  function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const txt = (k: string) => String(f.get(k) ?? "").trim();
    if (acao === "transferir") {
      void enviar(() => apiClient.post<AcervoDetalhe>(`v1/acervo/unidades/${id}/transferir`, { observacao: txt("observacao") }), "Transferida ao arquivo intermediário.");
    } else if (acao === "localizar") {
      void enviar(
        () => apiClient.put<AcervoDetalhe>(`v1/acervo/unidades/${id}/local`, { local_id: txt("local_id") || null, posicao: txt("posicao") }),
        "Localização registrada.",
      );
    } else if (acao === "bloquear") {
      void enviar(
        () => apiClient.post<AcervoDetalhe>(`v1/acervo/unidades/${id}/bloqueios`, { processo: txt("processo"), motivo: txt("motivo") }),
        "Bloqueio judicial registrado.",
      );
    } else if (acao === "emprestar") {
      void enviar(
        () =>
          apiClient.post<AcervoDetalhe>(`v1/acervo/unidades/${id}/emprestimos`, {
            solicitante_id: solicitante[0]?.id ?? "",
            finalidade: txt("finalidade"),
            devolver_ate: txt("devolver_ate"),
          }),
        "Empréstimo registrado.",
      );
    } else if (acao && typeof acao === "object" && "devolver" in acao) {
      void enviar(() => apiClient.post<AcervoDetalhe>(`v1/acervo/emprestimos/${acao.devolver}/devolver`, { observacao: txt("observacao") }), "Devolução registrada.");
    } else if (acao && typeof acao === "object") {
      void enviar(() => apiClient.post<AcervoDetalhe>(`v1/acervo/bloqueios/${acao.liberar}/liberar`, { motivo: txt("motivo") }), "Bloqueio liberado.");
    }
  }

  const titulos: Record<string, string> = {
    transferir: "Transferir ao arquivo intermediário",
    localizar: "Localização física",
    bloquear: "Registrar bloqueio judicial",
    emprestar: "Emprestar a unidade",
  };
  const tituloAcao =
    acao === null ? "" : typeof acao === "object" ? ("devolver" in acao ? "Registrar devolução" : "Liberar bloqueio judicial") : (titulos[acao] ?? "");
  const emprestimoAberto = u?.emprestimos.find((e) => !e.devolvido_em);
  const info = u ? situacao(u.situacao) : null;
  const final = u?.fase === "ELIMINADO";

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <DataState loading={ficha.isLoading} error={ficha.error} empty={!u}>
        {u && info && (
          <>
            <PageHeader
              eyebrow={`Acervo · ${TIPOS[u.tipo]}`}
              title={u.numero}
              description={u.assunto}
              actions={
                <>
                  {u.pode.gerir && !final && (
                    <>
                      <Link href={`/acervo/${id}/editar`} className={buttonClass("secondary", "md")}>
                        <Pencil size={16} aria-hidden="true" /> Alterar
                      </Link>
                      <Button variant="secondary" onClick={() => setAcao("localizar")}>
                        <MapPin size={16} aria-hidden="true" /> Localizar
                      </Button>
                    </>
                  )}
                  <Link href={`/acervo/${id}/etiqueta`} className={buttonClass("ghost", "md")}>
                    <QrCode size={16} aria-hidden="true" /> Etiqueta
                  </Link>
                  {u.pode.gerir && !final && !emprestimoAberto && (
                    <Button
                      variant="secondary"
                      onClick={() => {
                        setSolicitante([]);
                        setAcao("emprestar");
                      }}
                    >
                      <HandHelping size={16} aria-hidden="true" /> Emprestar
                    </Button>
                  )}
                  {u.pode.gerir && u.fase === "CORRENTE" && (
                    <Button onClick={() => setAcao("transferir")} disabled={u.situacao !== "TRANSFERIR"} title={situacao("TRANSFERIR").ajuda}>
                      <ArrowRightLeft size={16} aria-hidden="true" /> Transferir
                    </Button>
                  )}
                  {u.pode.bloquear && !final && (
                    <Button variant="danger" onClick={() => setAcao("bloquear")}>
                      <Gavel size={16} aria-hidden="true" /> Bloqueio judicial
                    </Button>
                  )}
                </>
              }
            />

            <div className="flex flex-wrap items-center gap-2">
              <Badge tone={info.tone}>{info.label}</Badge>
              <Badge>{FASES[u.fase]}</Badge>
              <Badge tone={u.destinacao === "GUARDA_PERMANENTE" ? "success" : "neutral"}>Destinação: {DESTINACOES[u.destinacao]}</Badge>
              {u.nivel_acesso !== "PUBLICO" && <Badge tone="warning">Acesso {NIVEIS[u.nivel_acesso].toLowerCase()}</Badge>}
            </div>
            {info.ajuda && <p className="-mt-2 text-sm text-muted">{info.ajuda}</p>}
            {u.recolhida_em && u.recolhimento_id && (
              <p className="-mt-2 text-sm">
                Sob a guarda do Arquivo Público desde {fmtDia(u.recolhida_em)} —{" "}
                <Link href={`/acervo/recolhimentos/${u.recolhimento_id}`} className="underline">
                  ver o recolhimento e o registro de entrada
                </Link>
                .
              </p>
            )}

            <div className="grid gap-4 lg:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle as="h2" className="text-base">Classificação e prazos (TTDD)</CardTitle>
                </CardHeader>
                <CardContent>
                  <dl className="grid gap-3 sm:grid-cols-2">
                    <Campo nome="Série">
                      <Link href={`/atlas/ttdd/${encodeURIComponent(u.codigo_classificacao)}`} className="hover:underline">
                        <span className="font-mono">{u.codigo_classificacao}</span> {u.serie.descritor}
                      </Link>
                      {!u.serie.vigente && <span className="block text-xs text-danger">Série fora da TTDD em vigor.</span>}
                    </Campo>
                    <Campo nome="Datas-limite">
                      {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)}
                    </Campo>
                    <Campo nome={`Fim da fase corrente (${u.prazos.regra_corrente || "—"})`}>{fmtDia(u.prazos.fim_corrente)}</Campo>
                    <Campo nome={`Fim da fase intermediária (${u.prazos.regra_intermediaria || "—"})`}>{fmtDia(u.prazos.fim_intermediaria)}</Campo>
                    <Campo nome="Data do evento">{fmtDia(u.data_evento)}</Campo>
                    <Campo nome="Contas aprovadas pelo TCE">{fmtDia(u.contas_aprovadas_em)}</Campo>
                  </dl>
                  <p className="mt-3 text-xs text-muted">
                    Prazos contados da data-limite final; a fase intermediária, do fim da corrente (Lei 10.301/2019, arts. 29 e 30).
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle as="h2" className="text-base">Unidade e localização</CardTitle>
                </CardHeader>
                <CardContent>
                  <dl className="grid gap-3 sm:grid-cols-2">
                    <Campo nome="Quantificação">
                      {u.quantidade} {u.especificacao}
                    </Campo>
                    <Campo nome="Local">{u.local_caminho ? `${u.local_caminho}${u.posicao ? `, ${u.posicao}` : ""}` : "Não localizada"}</Campo>
                    {u.hipotese_legal && <Campo nome="Hipótese legal">{u.hipotese_legal}</Campo>}
                    {u.observacoes && <Campo nome="Observações">{u.observacoes}</Campo>}
                  </dl>
                </CardContent>
              </Card>
            </div>

            <DigitaisDaUnidade unidadeId={u.id} podeGerir={u.pode.gerir && !final} />

            <Card>
              <CardHeader>
                <CardTitle as="h2" className="text-base">Empréstimos (art. 12, II)</CardTitle>
              </CardHeader>
              <CardContent>
                {u.emprestimos.length === 0 ? (
                  <p className="text-sm text-muted">Nenhum empréstimo registrado.</p>
                ) : (
                  <ul className="flex flex-col divide-y divide-surface-border text-sm" aria-label="Empréstimos">
                    {u.emprestimos.map((e) => (
                      <li key={e.id} className="flex flex-wrap items-start justify-between gap-2 py-2">
                        <div>
                          <p>
                            {e.solicitante || "—"} — {e.finalidade}
                          </p>
                          <p className="text-xs text-muted">
                            Retirada em {fmtDateTime(e.retirado_em)} · devolver até {fmtDia(e.devolver_ate)}
                            {e.devolvido_em && ` · devolvida em ${fmtDateTime(e.devolvido_em)}`}
                          </p>
                        </div>
                        {!e.devolvido_em && u.pode.gerir && (
                          <Button size="sm" variant="secondary" onClick={() => setAcao({ devolver: e.id })}>
                            Registrar devolução
                          </Button>
                        )}
                      </li>
                    ))}
                  </ul>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle as="h2" className="text-base">Bloqueios judiciais (arts. 32, II, e 44)</CardTitle>
              </CardHeader>
              <CardContent>
                {u.bloqueios.length === 0 ? (
                  <p className="text-sm text-muted">Nenhum bloqueio registrado.</p>
                ) : (
                  <ul className="flex flex-col divide-y divide-surface-border text-sm">
                    {u.bloqueios.map((b) => (
                      <li key={b.id} className="flex flex-wrap items-start justify-between gap-2 py-2">
                        <div>
                          <p>
                            Processo <span className="font-mono">{b.processo}</span>
                            {b.motivo && ` — ${b.motivo}`}
                          </p>
                          <p className="text-xs text-muted">
                            Desde {fmtDateTime(b.criado_em)}
                            {b.liberado_em && ` · liberado em ${fmtDateTime(b.liberado_em)}: ${b.motivo_liberacao}`}
                          </p>
                        </div>
                        {!b.liberado_em && u.pode.bloquear && (
                          <Button size="sm" variant="secondary" onClick={() => setAcao({ liberar: b.id })}>
                            Liberar
                          </Button>
                        )}
                      </li>
                    ))}
                  </ul>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle as="h2" className="text-base">Histórico</CardTitle>
              </CardHeader>
              <CardContent>
                <ol className="flex flex-col gap-2 text-sm" aria-label="Histórico da unidade">
                  {u.movimentacoes.map((m) => (
                    <li key={m.id}>
                      <span className="text-xs text-muted">{fmtDateTime(m.em)}</span> — {m.descricao}
                    </li>
                  ))}
                </ol>
              </CardContent>
            </Card>
          </>
        )}
      </DataState>

      <Dialog open={acao !== null} onClose={() => setAcao(null)} title={tituloAcao} size="lg">
        <form onSubmit={submit} className="flex flex-col gap-4" aria-label={tituloAcao}>
          {acao === "transferir" && (
            <>
              <p className="text-sm text-muted">O prazo da fase corrente foi cumprido. A unidade passa ao arquivo intermediário do órgão (art. 29).</p>
              <Input id="observacao" name="observacao" label="Observação" maxLength={500} />
            </>
          )}
          {acao === "localizar" && (
            <>
              <Select
                id="local_id"
                name="local_id"
                label="Local"
                placeholder="Sem localização"
                defaultValue={u?.local_id ?? ""}
                options={(locais.data ?? []).filter((l) => l.ativo).map((l) => ({ value: l.id, label: l.caminho }))}
              />
              <Input id="posicao" name="posicao" label="Posição" maxLength={120} placeholder="Ex.: prateleira 2, posição 7" defaultValue={u?.posicao} />
            </>
          )}
          {acao === "bloquear" && (
            <>
              <p className="text-sm text-muted">Enquanto o bloqueio estiver ativo, a unidade não pode ser eliminada (art. 44).</p>
              <Input id="processo" name="processo" label="Número do processo judicial *" required maxLength={120} />
              <Textarea id="motivo" name="motivo" label="Motivo" maxLength={2000} rows={3} />
            </>
          )}
          {acao === "emprestar" && (
            <>
              <UserPicker idPrefix="solicitante" label="Quem retira *" value={solicitante} onChange={(v) => setSolicitante(v.slice(-1))} />
              <Input id="finalidade" name="finalidade" label="Finalidade *" required maxLength={500} />
              <Input id="devolver_ate" name="devolver_ate" type="date" label="Devolver até *" required />
            </>
          )}
          {acao !== null && typeof acao === "object" && "devolver" in acao && (
            <Input id="observacao" name="observacao" label="Observação (estado da unidade, conferência)" maxLength={500} />
          )}
          {acao !== null && typeof acao === "object" && "liberar" in acao && (
            <Textarea id="motivo" name="motivo" label="Motivo da liberação *" required maxLength={2000} rows={3} />
          )}
          <div className="flex justify-end gap-2">
            <Button type="button" variant="ghost" onClick={() => setAcao(null)}>
              Cancelar
            </Button>
            <Button type="submit" disabled={pending}>
              Confirmar
            </Button>
          </div>
        </form>
      </Dialog>
    </div>
  );
}
