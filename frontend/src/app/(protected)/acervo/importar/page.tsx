"use client";

import { ArrowLeft, Download, FileUp } from "lucide-react";
import Link from "next/link";
import { useState, type ChangeEvent } from "react";

import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Button, buttonClass } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { apiClient } from "@/lib/api/client";
import type { AcervoImportacao } from "@/lib/sismarq/types";

/** Importação do acervo existente por planilha (ADR 033): simula, mostra os
 * erros por linha e só grava a planilha conferida, sem nenhum erro — tudo ou
 * nada. */
export default function ImportarPage() {
  const { run, pending } = useAction();
  const [conteudo, setConteudo] = useState("");
  const [nome, setNome] = useState("");
  const [sim, setSim] = useState<AcervoImportacao | null>(null);
  const [feito, setFeito] = useState<AcervoImportacao | null>(null);

  async function escolher(e: ChangeEvent<HTMLInputElement>) {
    const arquivo = e.target.files?.[0];
    setSim(null);
    setFeito(null);
    if (!arquivo) return;
    setNome(arquivo.name);
    setConteudo(await arquivo.text());
  }

  async function simular() {
    const r = await run(() => apiClient.post<AcervoImportacao>("v1/acervo/importacao/simular", { conteudo }));
    if (r) setSim(r.data);
  }

  async function aplicar() {
    const r = await run(
      () => apiClient.post<AcervoImportacao>("v1/acervo/importacao/aplicar", { conteudo, hash: sim?.hash }),
      "Planilha importada.",
    );
    if (r) {
      setFeito(r.data);
      setSim(null);
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Importar o acervo existente"
        description="Cadastre de uma vez as caixas, livros e processos que já estão no arquivo. Cada linha passa pelas mesmas regras do cadastro; a planilha só entra se nenhuma linha tiver erro."
        actions={
          <a href={`/api/backend/${"v1/acervo/importacao/modelo"}`} download="modelo-importacao-acervo.csv" className={buttonClass("secondary", "md")}>
            <Download size={16} aria-hidden="true" /> Planilha-modelo
          </a>
        }
      />
      <Card>
        <CardHeader>
          <CardTitle as="h2" className="text-base">
            1. Escolha a planilha (CSV)
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-3 text-sm">
          <p className="text-muted">
            No Excel, &ldquo;Salvar como&rdquo; → CSV. Colunas obrigatórias: tipo, codigo_classificacao, unidade (sigla), data_inicial e
            data_final (dd/mm/aaaa). Até 5.000 linhas.
          </p>
          <label className="inline-flex w-fit cursor-pointer items-center gap-2 rounded-md border border-surface-border px-3 py-2">
            <FileUp size={16} aria-hidden="true" />
            <span>{nome || "Selecionar arquivo"}</span>
            <input type="file" accept=".csv,text/csv" className="sr-only" onChange={(e) => void escolher(e)} aria-label="Planilha CSV" />
          </label>
          <div>
            <Button onClick={() => void simular()} disabled={!conteudo || pending}>
              2. Conferir
            </Button>
          </div>
        </CardContent>
      </Card>

      {sim && (
        <Card>
          <CardHeader>
            <CardTitle as="h2" className="text-base">
              Resultado da conferência
            </CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col gap-3 text-sm">
            <p>
              {sim.linhas} linha(s) pronta(s) para entrar · {sim.erros.length} com erro.
            </p>
            {sim.erros.length > 0 ? (
              <>
                <p className="text-danger">Corrija as linhas abaixo na planilha e confira de novo. Nada foi gravado.</p>
                <ul className="flex max-h-80 flex-col gap-1 overflow-auto" aria-label="Erros por linha">
                  {sim.erros.map((e) => (
                    <li key={`${e.linha}-${e.motivo}`}>
                      <strong>Linha {e.linha}:</strong> {e.motivo}
                    </li>
                  ))}
                </ul>
              </>
            ) : (
              <div>
                <Button onClick={() => void aplicar()} disabled={pending}>
                  3. Importar {sim.linhas} unidade(s)
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {feito && (
        <p role="status" className="rounded-lg border border-success/40 bg-success/10 p-4 text-sm">
          {feito.linhas} unidade(s) de arquivamento cadastrada(s): {feito.numeros.slice(0, 5).join(", ")}
          {feito.numeros.length > 5 && "…"}. <Link href="/acervo" className="font-semibold underline">Ver o acervo</Link>
        </p>
      )}
    </div>
  );
}
