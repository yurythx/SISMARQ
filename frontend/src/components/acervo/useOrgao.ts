"use client";

import { useApiQuery } from "@/lib/api/swr";
import type { OrgTree } from "@/lib/sismarq/types";

/** Nome do órgão (entidade) e da unidade, para o cabeçalho dos Anexos da
 * Lei 10.301/2019 ("ÓRGÃO/ENTIDADE" e "UNIDADE/SETOR"). */
export function useOrgao(unidadeId?: string | null): { orgao: string; unidade: string; entidades: OrgTree[] } {
  const tree = useApiQuery<OrgTree[]>("v1/iam/org-tree");
  const entidades = tree.data ?? [];
  for (const e of entidades) {
    const u = e.unidades.find((x) => x.id === unidadeId);
    if (u) return { orgao: e.nome, unidade: u.sigla ? `${u.sigla} — ${u.nome}` : u.nome, entidades };
  }
  return { orgao: "", unidade: "", entidades };
}
