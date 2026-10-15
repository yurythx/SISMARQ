"use client";

import { ArrowLeft, Plus } from "lucide-react";
import Link from "next/link";
import { useState, type FormEvent } from "react";

import { TIPOS_LOCAL } from "@/components/acervo/labels";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Dialog } from "@/components/ui/Dialog";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Toggle } from "@/components/ui/Toggle";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoLocal } from "@/lib/sismarq/types";

/** Locais físicos do acervo: prédio > depósito > sala > corredor > estante >
 * prateleira. Só o Arquivo Público (acervo:central) os mantém. */
export default function LocaisPage() {
  const { canGlobal } = useSismarq();
  const central = canGlobal("acervo:central");
  const locais = useApiQuery<AcervoLocal[]>("v1/acervo/locais");
  const { run, pending } = useAction();
  const [editando, setEditando] = useState<AcervoLocal | "novo" | null>(null);
  const [ativo, setAtivo] = useState(true);
  const lista = locais.data ?? [];
  const atual = editando && editando !== "novo" ? editando : null;

  function abrir(l: AcervoLocal | "novo") {
    setAtivo(l === "novo" ? true : l.ativo);
    setEditando(l);
  }

  async function salvar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const body = {
      tipo: String(f.get("tipo")),
      nome: String(f.get("nome")).trim(),
      parent_id: String(f.get("parent_id") || "") || null,
      ativo,
    };
    const ok = await run(
      () => (atual ? apiClient.put(`v1/acervo/locais/${atual.id}`, body) : apiClient.post("v1/acervo/locais", body)),
      atual ? "Local alterado." : "Local criado.",
    );
    if (ok) {
      setEditando(null);
      await locais.mutate();
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Acervo
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Locais físicos"
        description="Onde estão as caixas, os livros e os processos: prédio, depósito, sala, corredor, estante e prateleira."
        actions={
          central && (
            <Button onClick={() => abrir("novo")}>
              <Plus size={16} aria-hidden="true" /> Novo local
            </Button>
          )
        }
      />
      <DataState loading={locais.isLoading} error={locais.error} empty={lista.length === 0} emptyTitle="Nenhum local cadastrado">
        <ul className="flex flex-col divide-y divide-surface-border rounded-xl border border-surface-border bg-surface" aria-label="Locais físicos">
          {lista.map((l) => (
            <li key={l.id} className="flex flex-wrap items-center justify-between gap-2 px-4 py-3">
              <div className="min-w-0">
                <p className="text-sm text-foreground">{l.caminho}</p>
                <p className="text-xs text-muted">{TIPOS_LOCAL[l.tipo]}</p>
              </div>
              <div className="flex items-center gap-2">
                {!l.ativo && <Badge tone="warning">Desativado</Badge>}
                {central && (
                  <Button size="sm" variant="secondary" onClick={() => abrir(l)}>
                    Alterar
                  </Button>
                )}
              </div>
            </li>
          ))}
        </ul>
      </DataState>

      <Dialog open={editando !== null} onClose={() => setEditando(null)} title={atual ? `Alterar ${atual.nome}` : "Novo local"} size="lg">
        <form onSubmit={salvar} className="flex flex-col gap-4" aria-label="Local físico">
          <Select
            id="tipo"
            name="tipo"
            label="Tipo *"
            required
            defaultValue={atual?.tipo ?? "PREDIO"}
            options={Object.entries(TIPOS_LOCAL).map(([value, label]) => ({ value, label }))}
          />
          <Input id="nome" name="nome" label="Nome *" required maxLength={120} defaultValue={atual?.nome} />
          <Select
            id="parent_id"
            name="parent_id"
            label="Dentro de"
            placeholder="Nenhum (nível mais alto)"
            defaultValue={atual?.parent_id ?? ""}
            options={lista.filter((l) => l.id !== atual?.id).map((l) => ({ value: l.id, label: l.caminho }))}
          />
          <Toggle checked={ativo} onChange={setAtivo} label="Ativo (recebe unidades de arquivamento)" />
          <div className="flex justify-end gap-2">
            <Button type="button" variant="ghost" onClick={() => setEditando(null)}>
              Cancelar
            </Button>
            <Button type="submit" disabled={pending}>
              Salvar
            </Button>
          </div>
        </form>
      </Dialog>
    </div>
  );
}
