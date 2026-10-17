import { fireEvent, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fmtDia, situacao } from "@/components/acervo/labels";
import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import UnidadePage from "./[id]/page";
import BuscaAcervoPage from "./busca/page";
import ExposicaoEditorPage from "./exposicoes/[id]/page";
import ExposicoesPage from "./exposicoes/page";
import LocaisPage from "./locais/page";
import NovaUnidadePage from "./novo/page";
import AcervoPage from "./page";

const UNIDADE = {
  id: "u1",
  numero: "UA-2026-000001",
  tipo: "CAIXA",
  codigo_classificacao: "2.0.02.00.07",
  assunto: "Pregões de 2015",
  unidade_id: "un1",
  data_inicial: "2015-01-01T00:00:00Z",
  data_final: "2015-12-31T00:00:00Z",
  quantidade: 3,
  especificacao: "caixas-arquivo",
  observacoes: "",
  fase: "CORRENTE",
  destinacao: "GUARDA_PERMANENTE",
  data_evento: null,
  contas_aprovadas_em: null,
  nivel_acesso: "PUBLICO",
  hipotese_legal: "",
  local_id: null,
  posicao: "",
  created_at: "2026-10-15T00:00:00Z",
  updated_at: "2026-10-15T00:00:00Z",
  fim_corrente: "2016-12-31T00:00:00Z",
  fim_intermediaria: "2017-12-31T00:00:00Z",
  situacao: "TRANSFERIR",
  local_caminho: "",
  recolhida_em: null,
  recolhimento_id: null,
  tramite_processo_id: null,
};

const FICHA = {
  ...UNIDADE,
  serie: { codigo: "2.0.02.00.07", descritor: "Pregão Eletrônico", vigente: true, destinacao: "GUARDA_PERMANENTE" },
  prazos: { fim_corrente: UNIDADE.fim_corrente, fim_intermediaria: UNIDADE.fim_intermediaria, situacao: "TRANSFERIR", regra_corrente: "1 ano(s)", regra_intermediaria: "1 ano(s)" },
  bloqueios: [
    { id: "b1", unidade_arquivamento_id: "u1", processo: "0001234-56", motivo: "Ação", criado_em: "2026-10-01T00:00:00Z", liberado_em: null, motivo_liberacao: "" },
  ],
  movimentacoes: [{ id: "m1", tipo: "CADASTRO", fase_de: null, fase_para: "CORRENTE", descricao: "Cadastro na fase corrente", em: "2026-10-15T00:00:00Z" }],
  emprestimos: [] as unknown[],
  pode: { gerir: true, bloquear: true },
};

const LOCAL = { id: "l1", parent_id: null, tipo: "PREDIO", nome: "Arquivo Central", ativo: true, caminho: "Arquivo Central" };

beforeEach(() => resetNavigation());

describe("rótulos do Acervo", () => {
  it("formata datas sem fuso e conhece as situações", () => {
    expect(fmtDia("2015-12-31T00:00:00Z")).toBe("31/12/2015");
    expect(fmtDia(null)).toBe("—");
    expect(situacao("BLOQUEADA").tone).toBe("danger");
    expect(situacao("DESCONHECIDA").label).toBe("DESCONHECIDA");
  });
});

describe("AcervoPage", () => {
  it("mostra o painel e a lista, e filtra pela situação", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read", "acervo:manage"] }),
      "GET v1/acervo/painel": { data: { TRANSFERIR: 4, NO_CORRENTE: 2, AVALIAR_ELIMINACAO: 3, AGUARDA_EVENTO: 1, BLOQUEADA: 1 } },
      "GET v1/acervo/unidades": page([UNIDADE]),
    });
    renderApp(<AcervoPage />);
    expect(await screen.findByRole("link", { name: /UA-2026-000001/ })).toHaveAttribute("href", "/acervo/u1");
    expect(screen.getByRole("link", { name: /Nova unidade/ })).toHaveAttribute("href", "/acervo/novo");
    expect(await screen.findByRole("button", { name: /^3\s*Avaliar eliminação/ })).toBeInTheDocument();
    const card = await screen.findByRole("button", { name: /^4\s*Transferir ao intermediário/ });
    await userEvent.click(card);
    await waitFor(() => expect(backend.calls.some((c) => c.path === "v1/acervo/unidades" && c.query.get("situacao") === "TRANSFERIR")).toBe(true));
    expect(card).toHaveAttribute("aria-pressed", "true");
    fireEvent.change(screen.getByLabelText("Buscar"), { target: { value: "pregão" } });
    await waitFor(() => expect(backend.calls.some((c) => c.query.get("q") === "pregão")).toBe(true));
  });
});

describe("UnidadePage", () => {
  it("mostra a ficha e transfere ao intermediário", async () => {
    resetNavigation({ id: "u1" });
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read", "acervo:manage", "acervo:bloqueio"] }),
      "GET v1/acervo/unidades/u1": { data: FICHA },
      "GET v1/acervo/locais": { data: [LOCAL] },
      "POST v1/acervo/unidades/u1/transferir": { data: { ...FICHA, fase: "INTERMEDIARIA", situacao: "NO_INTERMEDIARIO" } },
      "POST v1/acervo/bloqueios/b1/liberar": { data: { ...FICHA, bloqueios: [] } },
    });
    renderApp(<UnidadePage />);
    expect(await screen.findByRole("heading", { name: "UA-2026-000001", level: 1 })).toBeInTheDocument();
    expect(screen.getByText("31/12/2016")).toBeInTheDocument();
    expect(screen.getByText(/0001234-56/)).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: /Transferir/ }));
    await userEvent.click(screen.getByRole("button", { name: "Confirmar" }));
    await waitFor(() => expect(backend.calls.some((c) => c.method === "POST" && c.path === "v1/acervo/unidades/u1/transferir")).toBe(true));
    expect(await screen.findByText("No arquivo intermediário")).toBeInTheDocument();
  });
});

describe("NovaUnidadePage", () => {
  it("cadastra e abre a ficha", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:manage"] }),
      "GET v1/iam/org-tree": { data: [] },
      "GET v1/atlas/ttdd/*": { data: { codigo: "2.0.02.00.07", descritor: "Pregão Eletrônico", revogada_em: null } },
      "POST v1/acervo/unidades": { data: FICHA },
    });
    renderApp(<NovaUnidadePage />);
    fireEvent.change(await screen.findByLabelText(/Código de classificação/), { target: { value: "2.0.02.00.07" } });
    expect(await screen.findByText("Pregão Eletrônico")).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("Data inicial *"), { target: { value: "2015-01-01" } });
    fireEvent.change(screen.getByLabelText("Data final *"), { target: { value: "2015-12-31" } });
    fireEvent.submit(screen.getByRole("form", { name: "Unidade de arquivamento" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/u1"));
    const post = backend.calls.find((c) => c.method === "POST");
    expect(post?.body).toMatchObject({ codigo_classificacao: "2.0.02.00.07", data_inicial: "2015-01-01", nivel_acesso: "PUBLICO" });
  });
});

describe("LocaisPage", () => {
  it("lista os locais e cria um novo (Arquivo Público)", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["*"] }),
      "GET v1/acervo/locais": { data: [LOCAL] },
      "POST v1/acervo/locais": { status: 201, data: LOCAL },
    });
    renderApp(<LocaisPage />);
    const lista = await screen.findByRole("list", { name: "Locais físicos" });
    expect(lista).toHaveTextContent("Arquivo Central");
    await userEvent.click(screen.getByRole("button", { name: /Novo local/ }));
    fireEvent.change(screen.getByLabelText("Nome *"), { target: { value: "Depósito 1" } });
    fireEvent.submit(screen.getByRole("form", { name: "Local físico" }));
    await waitFor(() => expect(backend.calls.some((c) => c.method === "POST" && c.path === "v1/acervo/locais")).toBe(true));
  });
});

const BUSCA = {
  items: [
    {
      id: "u1",
      numero: "UA-2026-000001",
      tipo: "CAIXA",
      codigo_classificacao: "2.0.02.00.07",
      descritor: "Pregão Eletrônico",
      assunto: "Pregões de 2015",
      orgao: "SEMAD",
      ano_inicial: 2015,
      ano_final: 2015,
      fase: "CORRENTE",
      destinacao: "GUARDA_PERMANENTE",
      nivel_acesso: "PUBLICO",
      digitais: 1,
      trecho: "pregão <mark>eletrônico</mark>",
    },
  ],
  meta: { page: 1, page_size: 20, total_items: 1, total_pages: 1 },
  facetas: { orgao: [{ valor: "SEMAD", total: 1 }], fase: [{ valor: "CORRENTE", total: 1 }] },
  origem: "indice" as const,
  ms: 12,
};

const EXPO = {
  id: "e1",
  slug: "memoria-municipal",
  titulo: "Memória municipal",
  resumo: "Série fotográfica",
  texto: "Apresentação",
  publicada_em: null as string | null,
  criado_em: "2026-10-01T00:00:00Z",
  atualizado_em: "2026-10-02T00:00:00Z",
  itens: [] as { digital_id: string; ordem: number; legenda: string; numero: string; descricao: string; paginas: number; unidade_arquivamento_id: string; unidade_numero: string; publicado: boolean }[],
};

describe("BuscaAcervoPage", () => {
  it("lista resultados, aplica faceta e restringe a digitais", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"] }),
      "GET v1/acervo/busca": { data: BUSCA },
    });
    renderApp(<BuscaAcervoPage />);
    expect(await screen.findByRole("link", { name: "Pregões de 2015" })).toHaveAttribute("href", "/acervo/u1");
    expect(screen.getByText("eletrônico").tagName).toBe("MARK");
    await userEvent.click(screen.getByRole("button", { name: /SEMAD/ }));
    await userEvent.click(screen.getByLabelText("Só com documento digital"));
    fireEvent.submit(screen.getByRole("search"));
    await waitFor(() => {
      const q = backend.calls.filter((c) => c.path === "v1/acervo/busca").at(-1)?.query;
      expect(q?.get("orgao")).toBe("SEMAD");
      expect(q?.get("digital")).toBe("true");
    });
  });

  it("avisa quando a busca cai no banco, sem facetas", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"] }),
      "GET v1/acervo/busca": { data: { ...BUSCA, origem: "banco", facetas: {} } },
    });
    renderApp(<BuscaAcervoPage />);
    expect(await screen.findByText(/Índice da busca indisponível/)).toBeInTheDocument();
  });
});

describe("Exposições virtuais", () => {
  it("lista rascunho e publicada", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:central"] }),
      "GET v1/acervo/exposicoes": page([
        { ...EXPO, titulo: "Rascunho" },
        { ...EXPO, id: "e2", slug: "pub", titulo: "Publicada", publicada_em: "2026-10-10T00:00:00Z" },
      ]),
    });
    renderApp(<ExposicoesPage />);
    expect(await screen.findByRole("link", { name: "Rascunho" })).toHaveAttribute("href", "/acervo/exposicoes/e1");
    expect(screen.getByText(/Publicada em/)).toBeInTheDocument();
  });

  it("cria exposição, inclui documento, publica e recolhe", async () => {
    const criada = { ...EXPO, itens: [{ digital_id: "d1", ordem: 1, legenda: "Legenda 1", numero: "DG-1", descricao: "Ofício", paginas: 2, unidade_arquivamento_id: "u1", unidade_numero: "UA-1", publicado: true }] };
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:central"] }),
      "GET v1/acervo/digitais": page([{ id: "d1", numero: "DG-1", descricao: "Ofício", codigo_classificacao: "2.0.02.00.07" }]),
      "POST v1/acervo/exposicoes": { status: 201, data: criada },
      "GET v1/acervo/exposicoes/e1": { data: criada },
      "PUT v1/acervo/exposicoes/e1": { data: { ...criada, titulo: "Memória municipal" } },
      "POST v1/acervo/exposicoes/e1/publicacao": { data: { ...criada, publicada_em: "2026-10-17T00:00:00Z" } },
      "DELETE v1/acervo/exposicoes/e1/publicacao": { data: { ...criada, publicada_em: null } },
    });
    resetNavigation({ id: "nova" });
    renderApp(<ExposicaoEditorPage />);
    fireEvent.change(await screen.findByLabelText("Título"), { target: { value: "Memória municipal" } });
    fireEvent.change(screen.getByLabelText("Adicionar documento liberado"), { target: { value: "d1" } });
    await userEvent.click(screen.getByRole("button", { name: "Adicionar" }));
    fireEvent.submit(screen.getByRole("button", { name: "Salvar" }).closest("form")!);
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/exposicoes/e1"));
    expect(backend.to("POST v1/acervo/exposicoes")[0]?.body).toMatchObject({ titulo: "Memória municipal", itens: [{ digital_id: "d1", legenda: "" }] });

    resetNavigation({ id: "e1" });
    const { unmount } = renderApp(<ExposicaoEditorPage />);
    await userEvent.click(await screen.findByRole("button", { name: "Publicar" }));
    expect(await screen.findByRole("link", { name: /Ver a visita/ })).toHaveAttribute("href", "/arquivo/exposicoes/memoria-municipal");
    await userEvent.click(screen.getByRole("button", { name: "Recolher" }));
    await waitFor(() => expect(backend.to("DELETE v1/acervo/exposicoes/e1/publicacao")).toHaveLength(1));
    unmount();
  });

  it("reordena e tira documento no editor", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:central"] }),
      "GET v1/acervo/digitais": page([]),
      "GET v1/acervo/exposicoes/e1": {
        data: {
          ...EXPO,
          itens: [
            { digital_id: "d1", ordem: 1, legenda: "A", numero: "DG-1", descricao: "Um", paginas: 1, unidade_arquivamento_id: "u1", unidade_numero: "UA-1", publicado: true },
            { digital_id: "d2", ordem: 2, legenda: "B", numero: "DG-2", descricao: "Dois", paginas: 1, unidade_arquivamento_id: "u1", unidade_numero: "UA-1", publicado: false },
          ],
        },
      },
    });
    resetNavigation({ id: "e1" });
    renderApp(<ExposicaoEditorPage />);
    expect(await screen.findByText(/Não liberado/)).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Descer Um" }));
    expect(screen.getByText("1. Dois")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Tirar Um" }));
    expect(screen.queryByText(/Um/)).not.toBeInTheDocument();
  });
});
