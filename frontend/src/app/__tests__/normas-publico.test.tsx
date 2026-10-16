import { render, screen } from "@testing-library/react";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("server-only", () => ({}));
vi.mock("@/lib/env", async (orig) => ({
  ...(await orig<typeof import("@/lib/env")>()),
  APP_URL: "https://portal.gov.br",
}));
vi.mock("@/components/layout/PublicShell", () => ({
  PublicShell: ({ children }: { children: ReactNode }) => <main>{children}</main>,
}));

class NotFound extends Error {}
vi.mock("next/navigation", () => ({
  notFound: () => {
    throw new NotFound("notFound");
  },
}));

const api = vi.hoisted(() => ({ routes: {} as Record<string, unknown>, paths: [] as string[] }));
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
    api.paths.push(path);
    const k = Object.keys(api.routes)
      .filter((x) => path.startsWith(x))
      .sort((a, b) => b.length - a.length)[0];
    const v = k ? api.routes[k] : new PublicApiError(404, "NOT_FOUND", "x");
    if (v instanceof Error) throw v;
    return { data: v };
  };
  return { PublicApiError, publicGet: async (path: string) => lookup(path) };
});

import { PublicApiError } from "@/lib/api/publicServer";

import NormaPage, { generateMetadata } from "../normas/[id]/page";
import NormasPage from "../normas/page";

const R1 = {
  versao: 1,
  ementa: "Recolhimento ao Arquivo Público",
  texto: "**Art. 1º** Recolher em março.",
  motivo: "Redação original",
  situacao: "PUBLICADA",
  publicada_em: "2026-03-01T00:00:00Z",
  diario: "",
  vigente_ate: "2027-06-01T00:00:00Z",
  criada_em: "2026-02-01T00:00:00Z",
};
const R2 = { ...R1, versao: 2, ementa: "Procedimentos de recolhimento", texto: "Art. 1º Recolher em junho.", motivo: "Alterada pela IN 2/2027", publicada_em: "2027-06-01T00:00:00Z", diario: "DIORONDON-e nº 7.000", vigente_ate: null };
const NORMA = {
  id: "n1",
  numero: "1/2026",
  situacao: "VIGENTE",
  ementa: "Procedimentos de recolhimento",
  versao: 2,
  publicada_em: R2.publicada_em,
  revogada_em: null as string | null,
  revogacao_ato: "",
  tem_rascunho: false,
  criada_em: "2026-02-01T00:00:00Z",
  redacoes: [R2, R1],
  total_redacoes: 2,
};

const params = (id: string) => ({ params: Promise.resolve({ id }) });

beforeEach(() => {
  api.routes = {};
  api.paths = [];
});

describe("Normas — página pública", () => {
  it("lista as normas, com a busca", async () => {
    api.routes["acervo/publico/normas"] = [NORMA, { ...NORMA, id: "n2", numero: "2/2026", situacao: "REVOGADA", total_redacoes: 1, revogada_em: "2027-01-01T00:00:00Z" }];
    render(await NormasPage({ searchParams: Promise.resolve({ q: "recolh" }) }));
    expect(screen.getByRole("link", { name: "Instrução Normativa nº 1/2026" })).toHaveAttribute("href", "/normas/n1");
    expect(screen.getByText(/alterada 1 vez\(es\)/)).toBeInTheDocument();
    expect(screen.getByText(/revogada em 01\/01\/2027/)).toBeInTheDocument();
    expect(api.paths[0]).toContain("q=recolh");
  });

  it("vazia, sem busca e com busca; indisponível", async () => {
    api.routes["acervo/publico/normas"] = [];
    const { unmount } = render(await NormasPage({ searchParams: Promise.resolve({}) }));
    expect(screen.getByText("Nenhuma instrução normativa publicada ainda.")).toBeInTheDocument();
    unmount();
    const b = render(await NormasPage({ searchParams: Promise.resolve({ q: "x" }) }));
    expect(screen.getByText("Nenhuma norma encontrada para a busca.")).toBeInTheDocument();
    b.unmount();
    api.routes["acervo/publico/normas"] = new PublicApiError(503, "DOWN", "fora");
    render(await NormasPage({ searchParams: Promise.resolve({}) }));
    expect(screen.getByRole("alert")).toHaveTextContent("indisponível");
  });

  it("mostra o texto em vigor e as redações anteriores com a vigência", async () => {
    api.routes["acervo/publico/normas/n1"] = NORMA;
    render(await NormaPage(params("n1")));
    expect(screen.getByRole("heading", { level: 1, name: "Instrução Normativa nº 1/2026" })).toBeInTheDocument();
    expect(screen.getByText(/Redação 2, publicada em 01\/06\/2027 \(DIORONDON-e nº 7.000\)/)).toBeInTheDocument();
    expect(screen.getByRole("region", { name: "Texto em vigor" })).toHaveTextContent("Recolher em junho.");
    expect(screen.getByText(/Redação 1 — vigente de 01\/03\/2026 a 01\/06\/2027/)).toBeInTheDocument();
    expect(screen.getByText(/ementa: Recolhimento ao Arquivo Público/)).toBeInTheDocument();
    expect(screen.getByText("Motivo da redação em vigor: Alterada pela IN 2/2027")).toBeInTheDocument();
    expect(await generateMetadata(params("n1"))).toMatchObject({ title: "Instrução Normativa nº 1/2026", description: NORMA.ementa });
  });

  it("revogada continua consultável; inexistente é 404", async () => {
    api.routes["acervo/publico/normas/n1"] = {
      ...NORMA,
      situacao: "REVOGADA",
      revogada_em: "2027-08-01T00:00:00Z",
      revogacao_ato: "IN 9/2027",
      redacoes: [{ ...R1, vigente_ate: "2027-08-01T00:00:00Z" }],
      total_redacoes: 1,
    };
    render(await NormaPage(params("n1")));
    expect(screen.getByRole("note")).toHaveTextContent("Revogada em 01/08/2027 — IN 9/2027");
    expect(screen.getByRole("region", { name: "Último texto em vigor" })).toBeInTheDocument();
    expect(screen.queryByRole("heading", { name: "Redações anteriores" })).not.toBeInTheDocument();

    await expect(NormaPage(params("x"))).rejects.toBeInstanceOf(NotFound);
    expect(await generateMetadata(params("x"))).toEqual({ title: "Norma não encontrada" });
    api.routes["acervo/publico/normas/n9"] = new PublicApiError(503, "DOWN", "fora");
    await expect(NormaPage(params("n9"))).rejects.toBeInstanceOf(PublicApiError);
  });
});
