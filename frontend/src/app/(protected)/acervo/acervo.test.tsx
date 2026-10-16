import { fireEvent, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fmtDia, situacao } from "@/components/acervo/labels";
import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import UnidadePage from "./[id]/page";
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
      "GET v1/acervo/painel": { data: { TRANSFERIR: 4, NO_CORRENTE: 2 } },
      "GET v1/acervo/unidades": page([UNIDADE]),
    });
    renderApp(<AcervoPage />);
    expect(await screen.findByRole("link", { name: /UA-2026-000001/ })).toHaveAttribute("href", "/acervo/u1");
    expect(screen.getByRole("link", { name: /Nova unidade/ })).toHaveAttribute("href", "/acervo/novo");
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
