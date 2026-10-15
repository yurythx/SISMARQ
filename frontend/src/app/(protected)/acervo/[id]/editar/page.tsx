"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";

import { UnidadeForm, type UnidadeBody } from "@/components/acervo/UnidadeForm";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoDetalhe } from "@/lib/sismarq/types";

/** Alteração da unidade de arquivamento. */
export default function EditarUnidadePage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { run, pending } = useAction();
  const ficha = useApiQuery<AcervoDetalhe>(`v1/acervo/unidades/${id}`);

  async function salvar(body: UnidadeBody) {
    const d = await run(() => apiClient.put<AcervoDetalhe>(`v1/acervo/unidades/${id}`, body), "Unidade de arquivamento alterada.");
    if (d) router.push(`/acervo/${id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href={`/acervo/${id}`} className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Ficha da unidade
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title={ficha.data ? `Alterar ${ficha.data.numero}` : "Alterar unidade"}
        description="Unidade permanente não muda de série (arts. 20 e 42); unidade eliminada não muda (art. 38)."
      />
      <DataState loading={ficha.isLoading} error={ficha.error} empty={!ficha.data}>
        {ficha.data && <UnidadeForm inicial={ficha.data} pending={pending} onSubmit={salvar} enviar="Salvar" />}
      </DataState>
    </div>
  );
}
