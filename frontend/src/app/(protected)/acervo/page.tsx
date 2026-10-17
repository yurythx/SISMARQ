"use client";

import { FileUp, Plus } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { FASES, SITUACOES, TIPOS, fmtDia, situacao } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { Pagination } from "@/components/sismarq/Pagination";
import { Badge } from "@/components/ui/Badge";
import { buttonClass } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { useApiPage, useApiQuery, withQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoUnidade } from "@/lib/sismarq/types";

// Ordem do painel: o que pede ação primeiro.
const ORDEM_PAINEL = [
  "TRANSFERIR",
  "AVALIAR_ELIMINACAO",
  "RECOLHER",
  "AGUARDA_EVENTO",
  "BLOQUEADA",
  "SERIE_REVOGADA",
  "AGUARDA_DESTINACAO",
  "NO_CORRENTE",
  "NO_INTERMEDIARIO",
  "PERMANENTE",
  "RECOLHIDA",
  "ELIMINADA",
];

/** Acervo (ADR 032): as unidades de arquivamento no escopo de quem consulta,
 * com o painel do que transferir, avaliar ou recolher (Lei 10.301/2019). */
export default function AcervoPage() {
  const { can } = useSismarq();
  const [q, setQ] = useState("");
  const [fase, setFase] = useState("");
  const [sit, setSit] = useState("");
  const [page, setPage] = useState(1);
  const painel = useApiQuery<Record<string, number>>("v1/acervo/painel");
  const lista = useApiPage<AcervoUnidade>(
    withQuery("v1/acervo/unidades", { q: q || undefined, fase: fase || undefined, situacao: sit || undefined, page, page_size: 20 }),
  );
  const itens = lista.data?.items ?? [];

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        eyebrow="SISMARQ"
        title="Acervo"
        description="Unidades de arquivamento nas três idades — corrente, intermediária e permanente (Lei 10.301/2019, art. 18) — com os prazos da Tabela de Temporalidade."
        actions={
          can("acervo:manage") ? (
            <>
              <Link href="/acervo/importar" className={buttonClass("secondary", "md")}>
                <FileUp size={16} aria-hidden="true" /> Importar
              </Link>
              <Link href="/acervo/novo" className={buttonClass("primary", "md")}>
                <Plus size={16} aria-hidden="true" /> Nova unidade
              </Link>
            </>
          ) : undefined
        }
      />

      <section aria-labelledby="painel-titulo">
        <h2 id="painel-titulo" className="sr-only">
          Painel de situações
        </h2>
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {(painel.data?.EMPRESTADAS ?? 0) > 0 && (
            <li>
              <Link
                href="/acervo/emprestimos"
                className="flex w-full flex-col items-start gap-1 rounded-xl border border-surface-border bg-surface p-4 hover:bg-surface-2"
              >
                <span className="text-2xl font-bold text-foreground">{painel.data?.EMPRESTADAS}</span>
                <Badge tone={(painel.data?.EMPRESTIMOS_ATRASADOS ?? 0) > 0 ? "danger" : "info"}>Emprestadas</Badge>
                <span className="text-xs text-muted">
                  {painel.data?.EMPRESTIMOS_ATRASADOS ?? 0} com a devolução atrasada (art. 12, II).
                </span>
              </Link>
            </li>
          )}
          {ORDEM_PAINEL.filter((s) => (painel.data?.[s] ?? 0) > 0).map((s) => {
            const info = situacao(s);
            const ativo = sit === s;
            return (
              <li key={s}>
                <button
                  type="button"
                  aria-pressed={ativo}
                  onClick={() => {
                    setSit(ativo ? "" : s);
                    setPage(1);
                  }}
                  className={`flex w-full flex-col items-start gap-1 rounded-xl border p-4 text-left transition-colors ${
                    ativo ? "border-primary bg-primary/5" : "border-surface-border bg-surface hover:bg-surface-2"
                  }`}
                >
                  <span className="text-2xl font-bold text-foreground">{painel.data?.[s]}</span>
                  <Badge tone={info.tone}>{info.label}</Badge>
                  <span className="text-xs text-muted">{info.ajuda}</span>
                </button>
              </li>
            );
          })}
        </ul>
      </section>

      <form
        role="search"
        className="grid gap-3 sm:grid-cols-[2fr_1fr_1fr]"
        onSubmit={(e) => e.preventDefault()}
        aria-label="Filtrar o acervo"
      >
        <Input
          id="acervo-q"
          label="Buscar"
          placeholder="Número, código de classificação ou assunto"
          value={q}
          onChange={(e) => {
            setQ(e.target.value);
            setPage(1);
          }}
        />
        <Select
          id="acervo-fase"
          label="Fase"
          value={fase}
          placeholder="Todas"
          options={Object.entries(FASES).map(([value, label]) => ({ value, label }))}
          onChange={(e) => {
            setFase(e.target.value);
            setPage(1);
          }}
        />
        <Select
          id="acervo-situacao"
          label="Situação"
          value={sit}
          placeholder="Todas"
          options={Object.entries(SITUACOES).map(([value, s]) => ({ value, label: s.label }))}
          onChange={(e) => {
            setSit(e.target.value);
            setPage(1);
          }}
        />
      </form>

      <DataState
        loading={lista.isLoading}
        error={lista.error}
        empty={itens.length === 0}
        emptyTitle="Nenhuma unidade de arquivamento"
        emptyDescription="Cadastre as caixas, livros e processos do arquivo da sua secretaria."
      >
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Unidades de arquivamento">
          {itens.map((u) => {
            const info = situacao(u.situacao);
            return (
              <li key={u.id} className="flex flex-col gap-1 p-4 sm:flex-row sm:items-center sm:justify-between">
                <div className="min-w-0">
                  <Link href={`/acervo/${u.id}`} className="font-medium text-foreground hover:text-primary hover:underline">
                    <span className="font-mono text-xs text-muted">{u.numero}</span> {u.assunto || u.codigo_classificacao}
                  </Link>
                  <p className="text-xs text-muted">
                    {TIPOS[u.tipo]} · <span className="font-mono">{u.codigo_classificacao}</span> · {fmtDia(u.data_inicial)} a{" "}
                    {fmtDia(u.data_final)} · {FASES[u.fase]}
                    {u.local_caminho && ` · ${u.local_caminho}${u.posicao ? `, ${u.posicao}` : ""}`}
                  </p>
                </div>
                <Badge tone={info.tone} title={info.ajuda}>
                  {info.label}
                </Badge>
              </li>
            );
          })}
        </ul>
        <Pagination meta={lista.data?.meta} onPage={setPage} />
      </DataState>
    </div>
  );
}
