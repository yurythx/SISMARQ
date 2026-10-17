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
  AGUARDANDO_SEMAD: { label: "Aguardando a SEMAD", tone: "warning", ajuda: "Arts. 35 e 36: a Secretaria Municipal de Administração ainda não aprovou." },
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

/** Situação da instrução normativa. */
export const SITUACAO_NORMA: Record<string, { label: string; tone: Tone }> = {
  RASCUNHO: { label: "Em elaboração", tone: "neutral" },
  VIGENTE: { label: "Vigente", tone: "success" },
  REVOGADA: { label: "Revogada", tone: "danger" },
};

export function situacaoNorma(s: string) {
  return SITUACAO_NORMA[s] ?? { label: s, tone: "neutral" as Tone };
}

/** Etapa do arquivo privado (arts. 55 e 56). */
export const SITUACAO_PRIVADO: Record<string, { label: string; tone: Tone }> = {
  CADASTRADO: { label: "Aguardando parecer", tone: "info" },
  COM_PARECER: { label: "Parecer favorável — aguardando decreto", tone: "warning" },
  NAO_DECLARADO: { label: "Parecer contrário", tone: "neutral" },
  DECLARADO: { label: "De interesse público e social", tone: "success" },
};

export function situacaoPrivado(s: string) {
  return SITUACAO_PRIVADO[s] ?? { label: s, tone: "neutral" as Tone };
}

/** Guarda do arquivo declarado (art. 56, §§ 4º e 5º). */
export const GUARDAS: Record<string, string> = {
  TITULAR: "Com o titular",
  DEPOSITADO: "Depositado no Arquivo Público (revogável)",
  DOADO: "Doado ao Arquivo Público",
};

/** Manifestação do município na alienação (art. 58). */
export const DECISOES_ALIENACAO: Record<string, string> = {
  PENDENTE: "Aguardando a manifestação",
  INTERESSE: "O município tem interesse",
  SEM_INTERESSE: "O município não tem interesse",
};

/** Situação do representante digital (Fase 5). */
export const SITUACAO_DIGITAL: Record<string, { label: string; tone: Tone }> = {
  ENVIADO: { label: "Na fila do processamento", tone: "neutral" },
  PROCESSANDO: { label: "Processando (OCR e PDF/A)", tone: "info" },
  PRONTO: { label: "Pronto — PDF/A guardado", tone: "success" },
  REPROVADO: { label: "Reprovado — digitalize de novo", tone: "warning" },
  FALHOU: { label: "Falhou no processamento", tone: "danger" },
};

export function situacaoDigital(s: string) {
  return SITUACAO_DIGITAL[s] ?? { label: s, tone: "neutral" as Tone };
}

/** Origem do arquivo enviado (ADR 044). */
export const ORIGEM_DIGITAL: Record<string, string> = {
  SCAN: "Digitalização do papel",
  NATO_DIGITAL: "Nato-digital",
};

export function origemDigital(s?: string | null): string {
  return ORIGEM_DIGITAL[s ?? ""] ?? (s || "Digitalização do papel");
}

/** Nome atual do órgão central nos Anexos gerados (TTDD 2025; ADR 043). A lei
 * ainda fala em Divisão de Arquivo Público. */
export const ORGAO_ARQUIVO = "Superintendência de Arquivo Público";
export const ORGAO_ARQUIVO_EXTENSO =
  "Superintendência de Arquivo Público do Município de Rondonópolis, Estado de Mato Grosso";

/** Ajuda da data do evento conforme a condição da TTDD (ADR 043, item 0.8.5). */
export function ajudaEvento(condicao?: string | null, descritor?: string | null): string {
  const c = (condicao ?? "").toLowerCase();
  const d = (descritor ?? "").toLowerCase();
  if (c.includes("vigor")) {
    if (/contrat|convên|conven|termo de/.test(d)) {
      return "Informe a data em que o contrato, convênio ou termo deixou de vigorar (encerramento, rescisão ou término do prazo). Use o documento que comprova o evento.";
    }
    if (/gestão|gestao|mandato|exercício|exercicio/.test(d)) {
      return "Informe a data do término da gestão, do mandato ou do exercício a que a série se refere.";
    }
    if (/funcional|servidor|lotação|lotacao|cargo/.test(d)) {
      return "Informe a data do desligamento, da vacância ou do fim do vínculo.";
    }
    return "Informe a data em que a relação jurídica deixou de vigorar (fim de contrato, gestão, licença ou outro fato previsto na série). Não estime: use o documento que comprova o evento.";
  }
  if (/ap[oó]s a data do evento/.test(c)) {
    return "Informe a data do fato que a TTDD chama de evento; os dias ou meses da condição contam a partir dela.";
  }
  if (/prescri/.test(c)) {
    return "Informe a data em que se considera prescrita a pretensão (marco legal ou administrativo da série).";
  }
  return "Preencha só se a série da TTDD tiver prazo por condição (vigência, evento, prescrição). A data-limite (ano inicial e final) continua obrigatória. Se a TTDD for republicada, a unidade já classificada não some.";
}

/** Tamanho legível (KB, MB, GB). */
export function fmtBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  const u = ["KB", "MB", "GB"];
  let v = n / 1024;
  let i = 0;
  while (v >= 1024 && i < u.length - 1) {
    v /= 1024;
    i++;
  }
  return `${v.toLocaleString("pt-BR", { maximumFractionDigits: 1 })} ${u[i]}`;
}
