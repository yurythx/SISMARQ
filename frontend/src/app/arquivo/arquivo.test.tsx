import { render, screen } from "@testing-library/react";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("server-only", () => ({}));
vi.mock("@/lib/env", async (orig) => ({ ...(await orig<typeof import("@/lib/env")>()), APP_URL: "https://arquivo.gov.br" }));
vi.mock("@/components/layout/PublicShell", () => ({ PublicShell: ({ children }: { children: ReactNode }) => <main>{children}</main> }));
vi.mock("@/components/acervo/VisualizadorPDF", () => ({ VisualizadorPDF: ({ url }: { url: string }) => <p>PDF {url}</p> }));

class NotFound extends Error {}
vi.mock("next/navigation", () => ({
  notFound: () => {
    throw new NotFound("notFound");
  },
}));

const api = vi.hoisted(() => ({ routes: {} as Record<string, unknown> }));
vi.mock("@/lib/api/publicServer", async () => {
  class PublicApiError extends Error {
    constructor(
      readonly status: number,
      readonly code: string,
      message: string,
    ) {
      super(message);
    }
  }
  const lookup = (path: string) => {
    const key = Object.keys(api.routes).find((k) => path.startsWith(k));
    const v = key ? api.routes[key] : new PublicApiError(404, "NOT_FOUND", "x");
    if (v instanceof Error) throw v;
    return v;
  };
  return {
    PublicApiError,
    publicGet: async (path: string) => ({ data: lookup(path) }),
  };
});

import { PublicApiError } from "@/lib/api/publicServer";

import ArquivoPublicoPage from "./page";
import BuscaPublicaPage from "./busca/page";
import ExposicaoPage, { generateMetadata as metaExpo } from "./exposicoes/[slug]/page";
import FichaPublicaPage, { generateMetadata as metaFicha } from "./unidades/[id]/page";

const ITEM = {
  id: "u1",
  numero: "UA-2026-000001",
  tipo: "CAIXA",
  codigo_classificacao: "2.0.02.00.07",
  descritor: "Pregão Eletrônico",
  assunto: "Pregões de 2015",
  orgao: "SEMAD",
  ano_inicial: 2015,
  ano_final: 2015,
  fase: "INTERMEDIARIA",
  destinacao: "GUARDA_PERMANENTE",
  digitais: 2,
  trecho: "pregão <mark>eletrônico</mark>",
};

const BUSCA = {
  items: [ITEM],
  meta: { page: 2, page_size: 20, total_items: 40, total_pages: 2 },
  facetas: {
    orgao: [{ valor: "SEMAD", total: 12 }],
    classe: [{ valor: "2.0", total: 3 }],
    especie: [{ valor: "CAIXA", total: 4 }],
    fase: [{ valor: "INTERMEDIARIA", total: 8 }],
    destinacao: [{ valor: "GUARDA_PERMANENTE", total: 8 }],
    ano: [{ valor: "2015", total: 8 }],
  },
  origem: "indice",
  ms: 9,
};

const DIGITAL = {
  id: "d1",
  numero: "DG-1",
  descricao: "Ofício",
  paginas: 3,
  pdfa_bytes: 1200,
  pdfa_sha256: "abc",
  conformidade: "PDF/A-2b",
  digitalizado_em: "2026-01-01T00:00:00Z",
  publicado_em: "2026-01-02T00:00:00Z",
};

beforeEach(() => {
  api.routes = {};
});

describe("Arquivo Público — página inicial", () => {
  it("mostra números, exposições, recentes e o SIC", async () => {
    api.routes["acervo/publico/arquivo"] = {
      sic_url: "https://sic.exemplo.gov.br",
      serie_sic: "SIC",
      unidades: 1200,
      digitais: 80,
      exposicoes: [{ id: "e1", slug: "memoria", titulo: "Memória", resumo: "Série", texto: "", publicada_em: "2026-10-01T00:00:00Z", criado_em: "", atualizado_em: "" }],
      recentes: [{ ...DIGITAL, unidade_arquivamento_id: "u1", unidade_numero: "UA-1", assunto: "Pregões" }],
    };
    render(await ArquivoPublicoPage());
    expect(screen.getByRole("heading", { name: "Arquivo Público Municipal" })).toBeInTheDocument();
    expect(screen.getByText("1.200")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Memória" })).toHaveAttribute("href", "/arquivo/exposicoes/memoria");
    expect(screen.getByRole("link", { name: /Pregões|Ofício/ })).toHaveAttribute("href", "/arquivo/unidades/u1?doc=d1");
    expect(screen.getByRole("link", { name: /pedido de acesso pelo SIC/ })).toHaveAttribute("href", "https://sic.exemplo.gov.br");
  });

  it("sem painel e sem URL do SIC", async () => {
    api.routes["acervo/publico/arquivo"] = new PublicApiError(503, "E", "x");
    const { unmount } = render(await ArquivoPublicoPage());
    expect(screen.getByRole("alert")).toHaveTextContent(/indisponível/);
    unmount();
    api.routes["acervo/publico/arquivo"] = {
      sic_url: "",
      serie_sic: "",
      unidades: 0,
      digitais: 0,
      exposicoes: [],
      recentes: [],
    };
    render(await ArquivoPublicoPage());
    expect(screen.getByText(/Procure o SIC/)).toBeInTheDocument();
  });
});

describe("Consulta pública ao acervo", () => {
  it("lista, pagina, aplica e tira faceta", async () => {
    api.routes["acervo/publico/busca"] = BUSCA;
    render(await BuscaPublicaPage({ searchParams: Promise.resolve({ q: "pregao", orgao: "SEMAD", digital: "true", page: "2" }) }));
    expect(screen.getByRole("link", { name: "Pregões de 2015" })).toHaveAttribute("href", "/arquivo/unidades/u1");
    expect(screen.getByRole("navigation", { name: "Páginas" })).toHaveTextContent("Página 2 de 2");
    expect(screen.getByRole("link", { name: /Anterior/ })).toHaveAttribute("href", expect.stringContaining("page=1"));
    expect(screen.getByRole("list", { name: "Filtros aplicados" })).toHaveTextContent(/Órgão/);
    expect(screen.getByRole("link", { current: true })).toHaveAttribute("href", "/arquivo/busca?q=pregao&digital=true");
  });

  it("vazio, índice no banco, consulta inválida e indisponível", async () => {
    api.routes["acervo/publico/busca"] = { ...BUSCA, items: [], origem: "banco", facetas: {}, meta: { page: 1, page_size: 20, total_items: 0, total_pages: 1 } };
    const a = render(await BuscaPublicaPage({ searchParams: Promise.resolve({}) }));
    expect(screen.getByText(/Nada encontrado/)).toBeInTheDocument();
    expect(screen.getByText(/Busca simplificada/)).toBeInTheDocument();
    a.unmount();

    api.routes["acervo/publico/busca"] = new PublicApiError(422, "INVALID", "x");
    const b = render(await BuscaPublicaPage({ searchParams: Promise.resolve({ fase: "NOPE" }) }));
    expect(screen.getByRole("alert")).toHaveTextContent(/Consulta inválida/);
    b.unmount();

    api.routes["acervo/publico/busca"] = new PublicApiError(503, "E", "x");
    render(await BuscaPublicaPage({ searchParams: Promise.resolve({}) }));
    expect(screen.getByRole("alert")).toHaveTextContent(/indisponível/);
  });
});

describe("Exposição virtual pública", () => {
  it("mostra texto, legendas e metadados; 404 e erro", async () => {
    const expo = {
      id: "e1",
      slug: "memoria",
      titulo: "Memória municipal",
      resumo: "Série fotográfica",
      texto: "## Abertura",
      publicada_em: "2026-10-01T00:00:00Z",
      criado_em: "",
      atualizado_em: "",
      itens: [
        { digital_id: "d1", ordem: 1, legenda: "Vista da praça", numero: "DG-1", descricao: "Foto", paginas: 1, unidade_arquivamento_id: "u1", unidade_numero: "UA-1", publicado: true },
        { digital_id: "d2", ordem: 2, legenda: "", numero: "DG-2", descricao: "", paginas: 2, unidade_arquivamento_id: "u2", unidade_numero: "UA-2", publicado: true },
      ],
    };
    api.routes["acervo/publico/exposicoes/memoria"] = expo;
    expect(await metaExpo({ params: Promise.resolve({ slug: "memoria" }) })).toMatchObject({ title: "Memória municipal" });
    render(await ExposicaoPage({ params: Promise.resolve({ slug: "memoria" }) }));
    expect(screen.getByRole("heading", { name: "Memória municipal" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Foto/ })).toHaveAttribute("href", "/arquivo/unidades/u1?doc=d1");
    expect(screen.getByText("Vista da praça")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /DG-2/ })).toBeInTheDocument();

    api.routes = {};
    expect(await metaExpo({ params: Promise.resolve({ slug: "sumiu" }) })).toEqual({ title: "Exposição não encontrada" });
    await expect(ExposicaoPage({ params: Promise.resolve({ slug: "sumiu" }) })).rejects.toBeInstanceOf(NotFound);
    api.routes["acervo/publico/exposicoes/x"] = new PublicApiError(500, "E", "x");
    await expect(ExposicaoPage({ params: Promise.resolve({ slug: "x" }) })).rejects.toMatchObject({ status: 500 });
  });
});

describe("Ficha pública da unidade", () => {
  it("ficha com dois digitais, escolhe o doc e cobre ausência", async () => {
    const ficha = {
      id: "u1",
      numero: "UA-2026-000001",
      tipo: "CAIXA",
      codigo_classificacao: "2.0.02.00.07",
      descritor: "Pregão Eletrônico",
      assunto: "Pregões de 2015",
      orgao: "SEMAD",
      data_inicial: "2015-01-01T00:00:00Z",
      data_final: "2015-12-31T00:00:00Z",
      quantidade: 3,
      especificacao: "caixas",
      fase: "CORRENTE",
      destinacao: "GUARDA_PERMANENTE",
      recolhida: true,
      digitais: [DIGITAL, { ...DIGITAL, id: "d2", numero: "DG-2", descricao: "" }],
    };
    api.routes["acervo/publico/unidades/u1"] = ficha;
    expect(await metaFicha({ params: Promise.resolve({ id: "u1" }) })).toMatchObject({ title: "Pregões de 2015" });
    render(await FichaPublicaPage({ params: Promise.resolve({ id: "u1" }), searchParams: Promise.resolve({ doc: "d2" }) }));
    expect(screen.getByText(/sob a guarda do Arquivo Público/)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /DG-2/ })).toHaveAttribute("aria-current", "true");
    expect(screen.getByText(/PDF \/api\/public/)).toBeInTheDocument();
  });

  it("sem digitais, não recolhida, 404 e 500", async () => {
    api.routes["acervo/publico/unidades/u2"] = {
      id: "u2",
      numero: "UA-2",
      tipo: "PROCESSO",
      codigo_classificacao: "1",
      descritor: "",
      assunto: "",
      orgao: "SMS",
      data_inicial: "2010-01-01T00:00:00Z",
      data_final: "2010-12-31T00:00:00Z",
      quantidade: 1,
      especificacao: "processo",
      fase: "CORRENTE",
      destinacao: "ELIMINACAO",
      recolhida: false,
      digitais: [],
    };
    expect(await metaFicha({ params: Promise.resolve({ id: "u2" }) })).toMatchObject({ title: "UA-2" });
    render(await FichaPublicaPage({ params: Promise.resolve({ id: "u2" }), searchParams: Promise.resolve({}) }));
    expect(screen.getByText(/no órgão produtor/)).toBeInTheDocument();
    expect(screen.getByText(/ainda não tem documento digitalizado/)).toBeInTheDocument();

    api.routes = {};
    expect(await metaFicha({ params: Promise.resolve({ id: "x" }) })).toEqual({ title: "Unidade não encontrada" });
    await expect(FichaPublicaPage({ params: Promise.resolve({ id: "x" }), searchParams: Promise.resolve({}) })).rejects.toBeInstanceOf(NotFound);
    api.routes["acervo/publico/unidades/z"] = new PublicApiError(500, "E", "x");
    await expect(FichaPublicaPage({ params: Promise.resolve({ id: "z" }), searchParams: Promise.resolve({}) })).rejects.toMatchObject({ status: 500 });
  });
});
