"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { UnidadeForm, type UnidadeBody } from "@/components/acervo/UnidadeForm";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { apiClient } from "@/lib/api/client";
import type { AcervoDetalhe } from "@/lib/sismarq/types";

/** Cadastro de unidade de arquivamento na fase corrente (ADR 032). */
export default function NovaUnidadePage() {
  const router = useRouter();
  const { run, pending } = useAction();

  async function criar(body: UnidadeBody) {
    const d = await run(() => apiClient.post<AcervoDetalhe>("v1/acervo/unidades", body), "Unidade de arquivamento cadastrada.");
    if (d) router.push(`/acervo/${d.data.id}`);
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Nova unidade de arquivamento"
        description="A unidade entra na fase corrente. A série da Tabela de Temporalidade define os prazos de guarda e a destinação (Lei 10.301/2019, arts. 28 a 30)."
      />
      <UnidadeForm pending={pending} onSubmit={criar} enviar="Cadastrar" />
    </div>
  );
}
