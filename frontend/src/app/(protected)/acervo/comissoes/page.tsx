"use client";

import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useState, type FormEvent } from "react";

import { AREAS, PAPEIS } from "@/components/acervo/labels";
import { useOrgao } from "@/components/acervo/useOrgao";
import { UserPicker, type PickedUser } from "@/components/iam/UserPicker";
import { DataState } from "@/components/sismarq/DataState";
import { PageHeader } from "@/components/sismarq/PageHeader";
import { useAction } from "@/components/sismarq/useAction";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { apiClient } from "@/lib/api/client";
import { useApiQuery } from "@/lib/api/swr";
import { useSismarq } from "@/lib/sismarq/SismarqProvider";
import type { AcervoComissao, AcervoMembro } from "@/lib/sismarq/types";

const opcoes = (r: Record<string, string>) => Object.entries(r).map(([value, label]) => ({ value, label }));

/** Comissões de avaliação (Lei 10.301/2019, art. 31): a CSPAD de cada órgão
 * e a CCPAD. Só o Arquivo Público (acervo:central) as cadastra. */
export default function ComissoesPage() {
  const { canGlobal } = useSismarq();
  const central = canGlobal("acervo:central");
  const lista = useApiQuery<AcervoComissao[]>("v1/acervo/comissoes");
  const { entidades } = useOrgao();
  const { run, pending } = useAction();
  const [editando, setEditando] = useState<AcervoComissao | null>(null);
  const [aberto, setAberto] = useState(false);
  const [tipo, setTipo] = useState("SETORIAL");
  const [membros, setMembros] = useState<AcervoMembro[]>([]);
  const itens = lista.data ?? [];
  const nomeOrgao = (id: string | null) => entidades.find((e) => e.id === id)?.nome ?? "";

  function abrir(c: AcervoComissao | null) {
    setEditando(c);
    setTipo(c?.tipo ?? "SETORIAL");
    setMembros(c?.membros ?? []);
    setAberto(true);
  }

  function escolher(v: PickedUser[]) {
    setMembros(v.map((p) => membros.find((m) => m.user_id === p.id) ?? { user_id: p.id, nome: p.name, papel: "MEMBRO", area: "OUTRA" }));
  }

  function mudar(i: number, campo: "papel" | "area", valor: string) {
    setMembros(membros.map((m, j) => (j === i ? { ...m, [campo]: valor } : m)));
  }

  async function salvar(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const corpo = {
      tipo,
      entidade_id: tipo === "SETORIAL" ? String(f.get("entidade_id") ?? "") || null : null,
      nome: String(f.get("nome") ?? "").trim(),
      portaria: String(f.get("portaria") ?? "").trim(),
      ativa: f.get("ativa") === "on",
      membros: membros.map(({ user_id, papel, area }) => ({ user_id, papel, area })),
    };
    const r = await run(
      () => (editando ? apiClient.put<AcervoComissao>(`v1/acervo/comissoes/${editando.id}`, corpo) : apiClient.post<AcervoComissao>("v1/acervo/comissoes", corpo)),
      "Comissão salva.",
    );
    if (r) {
      setAberto(false);
      await lista.mutate();
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <Link href="/acervo/eliminacao" className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
        <ArrowLeft size={14} aria-hidden="true" /> Eliminação de documentos
      </Link>
      <PageHeader
        eyebrow="Acervo"
        title="Comissões de avaliação"
        description="A Comissão Setorial Permanente de Avaliação de Documentos (CSPAD) de cada órgão aprova as relações de eliminação; a Comissão Central (CCPAD) coordena o sistema. Composição mínima: arquivologia ou história, área jurídica, administração ou finanças e a unidade produtora — na CCPAD, também o Arquivo Público (art. 31, §§ 6º e 7º). A permissão acervo:avaliacao fica no perfil (IAM); os membros são usuários nomeados na portaria — sem portaria não se inventa composição."
        actions={
          central &&
          !aberto && (
            <Button onClick={() => abrir(null)}>Nova comissão</Button>
          )
        }
      />

      {aberto && (
        <form onSubmit={salvar} aria-label="Comissão" className="flex flex-col gap-4 rounded-xl border border-surface-border bg-surface p-4">
          <div className="grid gap-3 sm:grid-cols-2">
            <Select
              id="tipo"
              name="tipo"
              label="Tipo *"
              value={tipo}
              onChange={(e) => setTipo(e.target.value)}
              options={[
                { value: "SETORIAL", label: "Setorial (CSPAD) — de um órgão" },
                { value: "CENTRAL", label: "Central (CCPAD)" },
              ]}
            />
            {tipo === "SETORIAL" && (
              <Select
                id="entidade_id"
                name="entidade_id"
                label="Órgão/entidade *"
                required
                placeholder="Selecione…"
                defaultValue={editando?.entidade_id ?? ""}
                options={entidades.map((e) => ({ value: e.id, label: e.nome }))}
              />
            )}
            <Input id="nome" name="nome" label="Nome *" required maxLength={200} defaultValue={editando?.nome} />
            <Input id="portaria" name="portaria" label="Portaria de designação" maxLength={200} defaultValue={editando?.portaria} placeholder="Portaria nº …, de …" />
          </div>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" name="ativa" defaultChecked={editando?.ativa ?? true} /> Ativa (uma comissão ativa por órgão)
          </label>
          <UserPicker idPrefix="membros" label="Membros" value={membros.map((m) => ({ id: m.user_id, name: m.nome }))} onChange={escolher} />
          {membros.length > 0 && (
            <ul className="flex flex-col gap-2" aria-label="Composição">
              {membros.map((m, i) => (
                <li key={m.user_id} className="grid items-end gap-2 rounded-lg border border-surface-border p-2 sm:grid-cols-3">
                  <span className="text-sm">{m.nome || m.user_id}</span>
                  <Select id={`papel-${i}`} label="Papel" value={m.papel} onChange={(e) => mudar(i, "papel", e.target.value)} options={opcoes(PAPEIS)} />
                  <Select id={`area-${i}`} label="Área" value={m.area} onChange={(e) => mudar(i, "area", e.target.value)} options={opcoes(AREAS)} />
                </li>
              ))}
            </ul>
          )}
          <div className="flex gap-2">
            <Button type="submit" disabled={pending}>
              Salvar comissão
            </Button>
            <Button type="button" variant="ghost" onClick={() => setAberto(false)}>
              Cancelar
            </Button>
          </div>
        </form>
      )}

      <DataState loading={lista.isLoading} error={lista.error} empty={itens.length === 0} emptyTitle="Nenhuma comissão cadastrada">
        <ul className="grid gap-4 lg:grid-cols-2" aria-label="Comissões">
          {itens.map((c) => (
            <li key={c.id} className="flex flex-col gap-2 rounded-xl border border-surface-border bg-surface p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <h2 className="font-semibold">{c.nome}</h2>
                <div className="flex gap-2">
                  <Badge tone="info">{c.tipo === "CENTRAL" ? "CCPAD" : "CSPAD"}</Badge>
                  <Badge tone={c.ativa ? "success" : "neutral"}>{c.ativa ? "Ativa" : "Inativa"}</Badge>
                </div>
              </div>
              <p className="text-xs text-muted">
                {c.tipo === "SETORIAL" ? nomeOrgao(c.entidade_id) : "Órgão Central do SISMARQ"}
                {c.portaria && ` · ${c.portaria}`}
              </p>
              <ul className="text-sm">
                {c.membros.map((m) => (
                  <li key={m.user_id}>
                    {m.nome || m.user_id} — {PAPEIS[m.papel] ?? m.papel}, {AREAS[m.area] ?? m.area}
                  </li>
                ))}
              </ul>
              {c.avisos.map((a) => (
                <p key={a} className="text-xs text-warning">
                  {a}
                </p>
              ))}
              {central && (
                <div>
                  <Button size="sm" variant="secondary" onClick={() => abrir(c)}>
                    Editar
                  </Button>
                </div>
              )}
            </li>
          ))}
        </ul>
      </DataState>
    </div>
  );
}
