"use client";

import { ArrowLeft, Printer } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";

import { DESTINACOES, TIPOS, fmtDia } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { Button } from "@/components/ui/Button";
import { QRCodeSVG } from "@/components/ui/QRCodeSVG";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoDetalhe } from "@/lib/sismarq/types";

/** Etiqueta da unidade de arquivamento para a lombada da caixa ou a capa do
 * livro: o QR abre a ficha no sistema. */
export default function EtiquetaPage() {
  const { id } = useParams<{ id: string }>();
  const ficha = useApiQuery<AcervoDetalhe>(`v1/acervo/unidades/${id}`);
  const u = ficha.data;
  const url = typeof window === "undefined" ? `/acervo/${id}` : `${window.location.origin}/acervo/${id}`;

  return (
    <div className="mx-auto flex max-w-xl flex-col gap-5 text-foreground">
      <div className="flex items-center justify-between gap-2 print:hidden">
        <Link href={`/acervo/${id}`} className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Ficha da unidade
        </Link>
        <Button size="sm" onClick={() => window.print()}>
          <Printer size={14} aria-hidden="true" className="mr-1" /> Imprimir
        </Button>
      </div>
      <DataState loading={ficha.isLoading} error={ficha.error} empty={!u}>
        {u && (
          <section
            aria-label={`Etiqueta ${u.numero}`}
            className="flex gap-4 rounded-lg border-2 border-foreground p-4 print:break-inside-avoid"
          >
            <QRCodeSVG value={url} size={120} />
            <div className="flex min-w-0 flex-col gap-1">
              <p className="text-xs font-semibold uppercase tracking-wide">Arquivo Público Municipal · SISMARQ</p>
              <p className="font-mono text-xl font-bold">{u.numero}</p>
              <p className="text-sm">
                {TIPOS[u.tipo]} · <span className="font-mono">{u.codigo_classificacao}</span>
              </p>
              <p className="text-sm font-medium">{u.assunto}</p>
              <p className="text-xs">
                {fmtDia(u.data_inicial)} a {fmtDia(u.data_final)} · {u.quantidade} {u.especificacao}
              </p>
              <p className="text-xs">Destinação: {DESTINACOES[u.destinacao]}</p>
              {u.local_caminho && (
                <p className="text-xs">
                  {u.local_caminho}
                  {u.posicao && `, ${u.posicao}`}
                </p>
              )}
            </div>
          </section>
        )}
      </DataState>
    </div>
  );
}
