// Rótulos do Acervo (Lei 10.301/2019; ADR 032) — um lugar só para a tela.
import type { AcervoDestinacao, AcervoFase, AcervoNivelAcesso, AcervoTipo, AcervoTipoLocal } from "@/lib/sismarq/types";

export const FASES: Record<AcervoFase, string> = {
  CORRENTE: "Corrente",
  INTERMEDIARIA: "Intermediária",
  PERMANENTE: "Permanente",
  ELIMINADO: "Eliminado",
};

export const DESTINACOES: Record<AcervoDestinacao, string> = {
  ELIMINACAO: "Eliminação",
  GUARDA_PERMANENTE: "Guarda permanente",
  INDEFINIDA: "Indefinida na TTDD",
};

export const TIPOS: Record<AcervoTipo, string> = {
  CAIXA: "Caixa",
  PASTA: "Pasta",
  LIVRO: "Livro",
  PROCESSO: "Processo",
  DOSSIE: "Dossiê",
  DOCUMENTO: "Documento",
  MAPA: "Mapa ou planta",
  OUTRO: "Outro",
};

export const NIVEIS: Record<AcervoNivelAcesso, string> = {
  PUBLICO: "Público",
  RESTRITO: "Restrito",
  SIGILOSO: "Sigiloso",
};

export const TIPOS_LOCAL: Record<AcervoTipoLocal, string> = {
  PREDIO: "Prédio",
  DEPOSITO: "Depósito",
  SALA: "Sala",
  CORREDOR: "Corredor",
  ESTANTE: "Estante",
  PRATELEIRA: "Prateleira",
};

type Tone = "neutral" | "success" | "danger" | "warning" | "info";

/** Situação do prazo de guarda: o que fazer com a unidade agora. */
export const SITUACOES: Record<string, { label: string; tone: Tone; ajuda: string }> = {
  NO_CORRENTE: { label: "No arquivo corrente", tone: "neutral", ajuda: "Prazo corrente correndo." },
  TRANSFERIR: { label: "Transferir ao intermediário", tone: "info", ajuda: "Prazo corrente cumprido (art. 29)." },
  NO_INTERMEDIARIO: { label: "No arquivo intermediário", tone: "neutral", ajuda: "Prazo intermediário correndo." },
  AVALIAR_ELIMINACAO: { label: "Avaliar eliminação", tone: "warning", ajuda: "Prazos cumpridos; destino eliminação: vai à CSPAD (arts. 19 e 33)." },
  RECOLHER: { label: "Recolher ao Arquivo Público", tone: "info", ajuda: "Prazos cumpridos; destino guarda permanente (art. 53)." },
  AGUARDA_EVENTO: { label: "Aguarda a data do evento", tone: "warning", ajuda: "Prazo por condição (ex.: fim da vigência): informe a data do evento." },
  AGUARDA_DESTINACAO: { label: "Aguarda destinação", tone: "warning", ajuda: "A série não tem destinação definida na TTDD." },
  BLOQUEADA: { label: "Bloqueio judicial", tone: "danger", ajuda: "Retida por ação judicial (art. 44)." },
  SERIE_REVOGADA: { label: "Série revogada", tone: "danger", ajuda: "Reclassifique pela TTDD em vigor (art. 52)." },
  PERMANENTE: { label: "Guarda permanente", tone: "success", ajuda: "Inalienável e imprescritível (arts. 20 e 42)." },
  ELIMINADA: { label: "Eliminada", tone: "neutral", ajuda: "Registro histórico (art. 38)." },
};

export function situacao(s: string) {
  return SITUACOES[s] ?? { label: s, tone: "neutral" as Tone, ajuda: "" };
}

/** "2024-12-31T00:00:00Z" ou "2024-12-31" → "31/12/2024" (sem fuso). */
export function fmtDia(iso?: string | null): string {
  if (!iso) return "—";
  const [a, m, d] = iso.slice(0, 10).split("-");
  return `${d}/${m}/${a}`;
}

/** Situação da guia de transferência. */
export const SITUACAO_GUIA: Record<string, { label: string; tone: Tone }> = {
  ENVIADA: { label: "Aguardando recebimento", tone: "info" },
  RECEBIDA: { label: "Recebida", tone: "success" },
  CANCELADA: { label: "Cancelada", tone: "neutral" },
};

export function situacaoGuia(s: string) {
  return SITUACAO_GUIA[s] ?? { label: s, tone: "neutral" as Tone };
}
