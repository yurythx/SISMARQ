"use client";

import Link from "next/link";
import { useState, type FormEvent } from "react";

import { fmtBytes, situacaoDigital } from "@/components/acervo/labels";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { apiClient } from "@/lib/api/client";
import { useApiPage, withQuery } from "@/lib/api/swr";
import type { AcervoDigital, AcervoDetalheDigital, UploadTicket } from "@/lib/sismarq/types";
import { putToTicket } from "@/lib/sismarq/upload";

const EM_ANDAMENTO = ["ENVIADO", "PROCESSANDO"];

/** Os representantes digitais da unidade e o envio de um novo PDF do
 * scanner (Decreto 10.278/2020): o arquivo vai direto ao armazenamento e o
 * worker gera o PDF/A com OCR, metadados e hash. */
export function DigitaisDaUnidade({ unidadeId, podeGerir }: { unidadeId: string; podeGerir: boolean }) {
  const { run, pending } = useAction();
  const [arquivo, setArquivo] = useState<File | null>(null);
  const lista = useApiPage<AcervoDigital>(withQuery("v1/acervo/digitais", { unidade_id: unidadeId, page_size: 50 }), {
    // Enquanto há algo na fila, acompanha o pipeline.
    refreshInterval: (dados) => ((dados?.items ?? []).some((d) => EM_ANDAMENTO.includes(d.situacao)) ? 5000 : 0),
  });
  const itens = lista.data?.items ?? [];

  async function enviar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!arquivo) return;
    const form = e.currentTarget;
    const f = new FormData(form);
    const txt = (k: string) => String(f.get(k) ?? "").trim();
    const base = `v1/acervo/unidades/${unidadeId}/digitais`;
    const r = await run(async () => {
      const { data: ticket } = await apiClient.post<UploadTicket>(`${base}/upload`, { filename: arquivo.name });
      await putToTicket(ticket, arquivo);
      return apiClient.post<AcervoDetalheDigital>(base, {
        object_key: ticket.object_key,
        descricao: txt("descricao"),
        tipo_documento: txt("tipo_documento"),
        equipamento: txt("equipamento"),
        dpi_declarado: Number(f.get("dpi_declarado")),
        responsavel_cpf: txt("responsavel_cpf"),
      });
    }, "Enviado: o PDF/A fica pronto depois do processamento.");
    if (r) {
      form.reset();
      setArquivo(null);
      await lista.mutate();
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle as="h2" className="text-base">
          Representantes digitais (Decreto 10.278/2020)
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        {itens.length === 0 ? (
          <p className="text-sm text-muted">Nenhuma digitalização.</p>
        ) : (
          <ul className="flex flex-col divide-y divide-surface-border text-sm" aria-label="Representantes digitais">
            {itens.map((d) => {
              const s = situacaoDigital(d.situacao);
              return (
                <li key={d.id} className="flex flex-wrap items-center justify-between gap-2 py-2">
                  <div>
                    <Link href={`/acervo/digitais/${d.id}`} className="font-mono hover:underline">
                      {d.numero}
                    </Link>
                    {d.descricao && <span> — {d.descricao}</span>}
                    <p className="text-xs text-muted">
                      {d.situacao === "PRONTO"
                        ? `${d.paginas} página(s) · ${d.dpi_medido} DPI · ${fmtBytes(d.pdfa_bytes)} · ${d.conformidade}`
                        : d.motivo || `${fmtBytes(d.original_bytes)} enviados`}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    {d.integro === false && <Badge tone="danger">Integridade falhou</Badge>}
                    <Badge tone={s.tone}>{s.label}</Badge>
                  </div>
                </li>
              );
            })}
          </ul>
        )}
        {podeGerir && (
          <form onSubmit={enviar} aria-label="Enviar digitalização" className="flex flex-col gap-3 border-t border-surface-border pt-3">
            <div className="flex flex-col gap-1">
              <label htmlFor="arquivo-digital" className="text-sm font-medium">
                PDF do scanner *
              </label>
              <input
                id="arquivo-digital"
                type="file"
                accept="application/pdf"
                required
                onChange={(e) => setArquivo(e.target.files?.[0] ?? null)}
                className="text-sm"
              />
            </div>
            <div className="grid gap-3 sm:grid-cols-3">
              <Select
                id="tipo_documento"
                name="tipo_documento"
                label="Tipo *"
                defaultValue="TEXTO"
                options={[
                  { value: "TEXTO", label: "Texto (mínimo 300 DPI)" },
                  { value: "PLANTA", label: "Planta ou mapa (mínimo 600 DPI)" },
                ]}
              />
              <Input id="equipamento" name="equipamento" label="Scanner (modelo e nº de série) *" required maxLength={200} />
              <Input id="dpi_declarado" name="dpi_declarado" type="number" label="Resolução (DPI) *" required min={72} max={4800} defaultValue={300} />
              <Input id="descricao" name="descricao" label="Parte ou volume" maxLength={500} placeholder="Ex.: volume 1 de 3" />
              <Input id="responsavel_cpf" name="responsavel_cpf" label="CPF de quem digitalizou" maxLength={14} placeholder="Guardado mascarado" />
            </div>
            <div>
              <Button type="submit" disabled={pending || !arquivo}>
                Enviar para processamento
              </Button>
            </div>
          </form>
        )}
      </CardContent>
    </Card>
  );
}
