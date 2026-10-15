import { fireEvent, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import EtiquetaPage from "./[id]/etiqueta/page";
import EmprestimosPage from "./emprestimos/page";
import GuiaPage from "./guias/[id]/page";
import GuiasPage from "./guias/page";
import ImportarPage from "./importar/page";

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
  fase: "CORRENTE",
  destinacao: "GUARDA_PERMANENTE",
  local_caminho: "Arquivo Central",
  posicao: "P2",
  situacao: "TRANSFERIR",
};

const GUIA = {
  id: "g1",
  numero: "GT-2026-00001",
  origem_id: "un1",
  situacao: "ENVIADA",
  observacao: "Remessa anual",
  criada_em: "2026-10-15T00:00:00Z",
  recebida_em: null,
  itens: [UNIDADE],
  total_itens: 1,
};

beforeEach(() => resetNavigation({ id: "g1" }));

describe("Guias de transferência", () => {
  it("lista as guias", async () => {
    mockBackend({ ...identityRoutes({ permissions: ["acervo:read", "acervo:manage"] }), "GET v1/acervo/guias": page([GUIA]) });
    renderApp(<GuiasPage />);
    expect(await screen.findByRole("link", { name: "GT-2026-00001" })).toHaveAttribute("href", "/acervo/guias/g1");
    expect(screen.getByText("Aguardando recebimento")).toBeInTheDocument();
  });

  it("mostra a guia para imprimir e registra o recebimento", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read", "acervo:manage"] }),
      "GET v1/acervo/guias/g1": { data: GUIA },
      "POST v1/acervo/guias/g1/receber": { data: { ...GUIA, situacao: "RECEBIDA", recebida_em: "2026-10-16T00:00:00Z" } },
    });
    renderApp(<GuiaPage />);
    expect(await screen.findByRole("heading", { name: "Guia de transferência GT-2026-00001" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "Pregões de 2015" })).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Registrar recebimento" }));
    await waitFor(() => expect(backend.calls.some((c) => c.method === "POST" && c.path === "v1/acervo/guias/g1/receber")).toBe(true));
    expect(await screen.findByText("Recebida")).toBeInTheDocument();
  });
});

describe("ImportarPage", () => {
  it("confere a planilha, mostra os erros e só aplica sem erro", async () => {
    let simulacoes = 0;
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:manage"] }),
      "POST v1/acervo/importacao/simular": () => {
        simulacoes++;
        return simulacoes === 1
          ? { data: { hash: "h", linhas: 1, aplicada: false, erros: [{ linha: 3, motivo: "unidade: não encontrada" }], numeros: [] } }
          : { data: { hash: "h2", linhas: 2, aplicada: false, erros: [], numeros: [] } };
      },
      "POST v1/acervo/importacao/aplicar": { data: { hash: "h2", linhas: 2, aplicada: true, erros: [], numeros: ["UA-1", "UA-2"] } },
    });
    renderApp(<ImportarPage />);
    const arquivo = new File(["tipo;codigo_classificacao\n"], "acervo.csv", { type: "text/csv" });
    await userEvent.upload(screen.getByLabelText("Planilha CSV"), arquivo);
    await userEvent.click(await screen.findByRole("button", { name: "2. Conferir" }));
    expect(await screen.findByText(/Linha 3:/)).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /3\. Importar/ })).not.toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "2. Conferir" }));
    await userEvent.click(await screen.findByRole("button", { name: "3. Importar 2 unidade(s)" }));
    expect(await screen.findByText((_, el) => el?.tagName === "P" && /UA-1, UA-2/.test(el.textContent ?? ""))).toBeInTheDocument();
    expect(backend.calls.find((c) => c.path === "v1/acervo/importacao/aplicar")?.body).toMatchObject({ hash: "h2" });
  });
});

describe("EmprestimosPage", () => {
  it("lista os empréstimos e filtra os atrasados", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"] }),
      "GET v1/acervo/emprestimos": page([
        { id: "e1", unidade_arquivamento_id: "u1", solicitante: "Maria", finalidade: "consulta", devolver_ate: "2020-01-01T00:00:00Z", numero: "UA-1", assunto: "Pregões" },
      ]),
    });
    renderApp(<EmprestimosPage />);
    expect(await screen.findByText("Atrasado")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("switch", { name: "Só os atrasados" }));
    await waitFor(() => expect(backend.calls.some((c) => c.query.get("atrasados") === "true")).toBe(true));
  });
});

describe("EtiquetaPage", () => {
  it("mostra a etiqueta com o QR", async () => {
    resetNavigation({ id: "u1" });
    mockBackend({ ...identityRoutes(), "GET v1/acervo/unidades/u1": { data: { ...UNIDADE, emprestimos: [], bloqueios: [], movimentacoes: [] } } });
    renderApp(<EtiquetaPage />);
    const etiqueta = await screen.findByRole("region", { name: "Etiqueta UA-2026-000001" });
    expect(etiqueta).toHaveTextContent("Pregões de 2015");
    expect(etiqueta.querySelector("svg")).not.toBeNull();
  });
});
