// Rótulos do Acervo (Lei 10.301/2019; ADR 032) — um lugar só para a tela.
import type { AcervoDestinacao, AcervoFase, AcervoNivelAcesso, AcervoPreparo, AcervoTipo, AcervoTipoLocal, AcervoUnidade } from "@/lib/sismarq/types";

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
  PERMANENTE: { label: "Permanente no órgão", tone: "info", ajuda: "Inalienável (arts. 20 e 42); recolha ao Arquivo Público (art. 53)." },
  RECOLHIDA: { label: "No Arquivo Público", tone: "success", ajuda: "Recolhida: sob a guarda do Arquivo Público (art. 5º, VII)." },
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

/** Etapas da relação de eliminação (arts. 31 a 40). */
export const SITUACAO_LISTAGEM: Record<string, { label: string; tone: Tone; ajuda: string }> = {
  RASCUNHO: { label: "Em seleção", tone: "neutral", ajuda: "A CSPAD monta a relação e junta o parecer jurídico (art. 32, § 4º)." },
  APROVADA: { label: "Aprovada pela CSPAD", tone: "info", ajuda: "Aguarda a autorização do Arquivo Público (arts. 19 e 65)." },
  AUTORIZADA: { label: "Autorizada", tone: "info", ajuda: "Publique o Edital de Ciência no Diário Oficial do Município (art. 37)." },
  EM_EDITAL: { label: "Em edital", tone: "warning", ajuda: "Prazo de 30 dias para pedidos; faça a amostragem (art. 39)." },
  ELIMINADA: { label: "Eliminada", tone: "success", ajuda: "Termo de Eliminação registrado (Anexo II)." },
  CANCELADA: { label: "Cancelada", tone: "neutral", ajuda: "As unidades ficaram livres." },
};

export function situacaoListagem(s: string) {
  return SITUACAO_LISTAGEM[s] ?? { label: s, tone: "neutral" as Tone, ajuda: "" };
}

export const PAPEIS: Record<string, string> = {
  PRESIDENTE: "Presidente",
  SECRETARIO: "Secretário(a)",
  MEMBRO: "Membro",
};

/** Áreas exigidas na composição (art. 31, §§ 6º e 7º). */
export const AREAS: Record<string, string> = {
  ARQUIVOLOGIA_HISTORIA: "Arquivologia ou História",
  JURIDICO: "Área jurídica",
  ADMINISTRACAO_FINANCAS: "Administração ou finanças",
  UNIDADE_PRODUTORA: "Unidade produtora",
  ARQUIVO_PUBLICO: "Arquivo Público",
  OUTRA: "Outra área",
};

/** Ano (AAAA) de uma data ISO. */
export function ano(iso?: string | null): string {
  return iso ? iso.slice(0, 4) : "—";
}

/** Data ISO + n dias, em "AAAA-MM-DD" (sem fuso). */
export function somarDias(iso: string, n: number): string {
  const d = new Date(`${iso.slice(0, 10)}T00:00:00Z`);
  d.setUTCDate(d.getUTCDate() + n);
  return d.toISOString().slice(0, 10);
}

const MESES = ["janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"];

/** "2026-09-01…" → { dia: "1", mes: "setembro", ano: "2026" } (sem fuso). */
export function porExtenso(iso?: string | null): { dia: string; mes: string; ano: string } {
  if (!iso) return { dia: "___", mes: "_______", ano: "____" };
  return { dia: String(Number(iso.slice(8, 10))), mes: MESES[Number(iso.slice(5, 7)) - 1] ?? "", ano: iso.slice(0, 4) };
}

/** Situação do recolhimento ao Arquivo Público. */
export const SITUACAO_RECOLHIMENTO: Record<string, { label: string; tone: Tone }> = {
  ENVIADO: { label: "Aguardando o Arquivo Público", tone: "info" },
  RECEBIDO: { label: "Recebido", tone: "success" },
  DEVOLVIDO: { label: "Devolvido à origem", tone: "warning" },
  CANCELADO: { label: "Cancelado", tone: "neutral" },
};

export function situacaoRecolhimento(s: string) {
  return SITUACAO_RECOLHIMENTO[s] ?? { label: s, tone: "neutral" as Tone };
}

/** Motivo do recolhimento. */
export const MOTIVOS_RECOLHIMENTO: Record<string, string> = {
  REGULAR: "Regular — prazos da TTDD cumpridos (art. 53)",
  EXTINCAO: "Extinção do órgão ou entidade (arts. 17 e 50)",
  DESESTATIZACAO: "Desestatização (arts. 16 e 51)",
};

/** O preparo do art. 53, na ordem da lei. */
export const PREPARO: { campo: keyof AcervoPreparo; label: string }[] = [
  { campo: "classificados", label: "Classificados" },
  { campo: "avaliados", label: "Avaliados" },
  { campo: "organizados", label: "Organizados" },
  { campo: "higienizados", label: "Higienizados" },
  { campo: "acondicionados", label: "Acondicionados" },
];

/** Quais unidades da origem podem ir: no regular, as com os prazos cumpridos
 * (ou permanentes guardadas no órgão); na extinção e na desestatização, toda
 * a guarda permanente ainda no órgão (arts. 50 e 51). A API confere de novo. */
export function recolhiveis(itens: AcervoUnidade[], motivo: string): AcervoUnidade[] {
  return itens.filter((u) =>
    motivo === "REGULAR"
      ? u.situacao === "RECOLHER" || (u.situacao === "PERMANENTE" && !u.recolhida_em)
      : u.destinacao === "GUARDA_PERMANENTE" && !u.recolhida_em && u.fase !== "ELIMINADO",
  );
}
