import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { PublicShell } from "@/components/layout/PublicShell";
import { PublicApiError } from "@/lib/api/publicServer";
import { listPublishedServices } from "@/lib/catalog/catalog";
import { APP_URL } from "@/lib/env";
import { ServicosCatalogClient } from "./ServicosCatalogClient";

export const dynamic = "force-dynamic";

const description =
  "Serviços do Arquivo Público Municipal de Rondonópolis: acesso aos documentos, Tabela de Temporalidade, avaliação e eliminação, recolhimento e arquivos privados (Lei Municipal nº 10.301/2019).";

export const metadata: Metadata = {
  title: "Carta de Serviços do Arquivo Público",
  description,
  alternates: { canonical: `${APP_URL}/servicos` },
  openGraph: { title: "Carta de Serviços do Arquivo Público", description, type: "website", url: `${APP_URL}/servicos` },
};

export default async function ServicosPage() {
  let services: Awaited<ReturnType<typeof listPublishedServices>> = [];
  try {
    services = await listPublishedServices();
  } catch (err) {
    if (err instanceof PublicApiError && err.status === 404) notFound();
    throw err;
  }

  return (
    <PublicShell>
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-10 px-6 py-12">
        <div className="flex flex-col gap-3">
          <p className="dateline">Carta de serviços</p>
          <h1 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Serviços do Arquivo Público Municipal
          </h1>
          <p className="max-w-3xl text-base text-muted">{description}</p>
        </div>

        <ServicosCatalogClient initialServices={services} />
      </div>
    </PublicShell>
  );
}
