import Link from "next/link";
import { FileText } from "lucide-react";

import { DESTINACOES, FASES, NIVEIS, TIPOS } from "@/components/acervo/labels";
import { Badge } from "@/components/ui/Badge";
import type { AcervoDestinacao, AcervoFase, AcervoItemBusca, AcervoNivelAcesso, AcervoTipo } from "@/lib/sismarq/types";

/** Facetas da busca do acervo (artigo técnico, item 4.4), na ordem da tela.
 * O nível de acesso só existe na busca interna. */
export const FACETAS_PUBLICAS = ["orgao", "classe", "especie", "fase", "destinacao", "ano"] as const;
export const FACETAS_INTERNAS = [...FACETAS_PUBLICAS, "nivel_acesso"] as const;

export const NOME_FACETA: Record<string, string> = {
  orgao: "Órgão",
  classe: "Classe",
  especie: "Espécie",
  fase: "Fase",
  destinacao: "Destinação",
  ano: "Ano",
  nivel_acesso: "Nível de acesso",
};

/** O rótulo de um valor de faceta (o código vira o nome). */
export function rotuloFaceta(faceta: string, valor: string): string {
  switch (faceta) {
    case "fase":
      return FASES[valor as AcervoFase] ?? valor;
    case "destinacao":
      return DESTINACOES[valor as AcervoDestinacao] ?? valor;
    case "especie":
      return TIPOS[valor as AcervoTipo] ?? valor;
    case "nivel_acesso":
      return NIVEIS[valor as AcervoNivelAcesso] ?? valor;
    case "classe":
      return `Classe ${valor}`;
    default:
      return valor;
  }
}

/** O trecho do índice com os termos destacados. Só o <mark> vira marcação:
 * todo o resto é texto (o React escapa), nunca HTML. */
export function Trecho({ texto }: { texto: string }) {
  if (!texto) return null;
  const partes = texto.split(/(<mark>|<\/mark>)/);
  let dentro = false;
  const nos: React.ReactNode[] = [];
  partes.forEach((p, i) => {
    if (p === "<mark>") dentro = true;
    else if (p === "</mark>") dentro = false;
    else if (p) nos.push(dentro ? <mark key={i}>{p}</mark> : <span key={i}>{p}</span>);
  });
  return <p className="text-sm text-muted">…{nos}…</p>;
}

function anos(it: AcervoItemBusca): string {
  return it.ano_inicial === it.ano_final ? String(it.ano_inicial) : `${it.ano_inicial}–${it.ano_final}`;
}

/** Uma unidade encontrada: número, assunto, série, órgão, anos e fase. */
export function ItemResultado({ item, href }: { item: AcervoItemBusca; href: string }) {
  return (
    <li className="flex flex-col gap-1 p-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <Link href={href} className="font-semibold hover:underline">
          {item.assunto || item.descritor || item.numero}
        </Link>
        <div className="flex flex-wrap gap-1">
          {item.nivel_acesso && item.nivel_acesso !== "PUBLICO" && (
            <Badge tone={item.nivel_acesso === "SIGILOSO" ? "danger" : "warning"}>{rotuloFaceta("nivel_acesso", item.nivel_acesso)}</Badge>
          )}
          <Badge tone="neutral">{rotuloFaceta("fase", item.fase)}</Badge>
          {item.digitais > 0 && (
            <Badge tone="success">
              <FileText size={12} aria-hidden="true" /> {item.digitais} digital(is)
            </Badge>
          )}
        </div>
      </div>
      <p className="text-xs text-muted">
        <span className="font-mono">{item.numero}</span> · {item.codigo_classificacao}
        {item.descritor && ` — ${item.descritor}`} · {item.orgao} · {anos(item)}
      </p>
      <Trecho texto={item.trecho} />
    </li>
  );
}
