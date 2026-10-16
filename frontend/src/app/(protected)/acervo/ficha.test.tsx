import { fireEvent, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fail, identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import EditarUnidadePage from "./[id]/editar/page";
import UnidadePage from "./[id]/page";
import NovaGuiaPage from "./guias/nova/page";

const SOLICITANTE = "33333333-3333-4333-8333-333333333333";

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
  observacoes: "Caixa úmida",
  fase: "CORRENTE",
  destinacao: "ELIMINACAO",
  data_evento: null,
  contas_aprovadas_em: null,
  nivel_acesso: "RESTRITO",
  hipotese_legal: "LAI, art. 31",
  local_id: "l1",
  posicao: "P2",
  created_at: "2026-10-15T00:00:00Z",
  updated_at: "2026-10-15T00:00:00Z",
  fim_corrente: "2016-12-31T00:00:00Z",
  fim_intermediaria: "2017-12-31T00:00:00Z",
  situacao: "NO_CORRENTE",
  local_caminho: "Arquivo Central",
  recolhida_em: null,
  recolhimento_id: null,
  tramite_processo_id: null,
};

const FICHA = {
  ...UNIDADE,
  serie: { codigo: "2.0.02.00.07", descritor: "Pregão Eletrônico", vigente: false, destinacao: "ELIMINACAO" },
  prazos: { fim_corrente: UNIDADE.fim_corrente, fim_intermediaria: UNIDADE.fim_intermediaria, situacao: "NO_CORRENTE", regra_corrente: "", regra_intermediaria: "" },
  bloqueios: [
    { id: "b1", unidade_arquivamento_id: "u1", processo: "0001234-56", motivo: "", criado_em: "2026-10-01T00:00:00Z", liberado_em: null, motivo_liberacao: "" },
    { id: "b2", unidade_arquivamento_id: "u1", processo: "0009999-00", motivo: "Ação popular", criado_em: "2026-09-01T00:00:00Z", liberado_em: "2026-09-20T00:00:00Z", motivo_liberacao: "Trânsito em julgado" },
  ],
  movimentacoes: [{ id: "m1", tipo: "CADASTRO", fase_de: null, fase_para: "CORRENTE", descricao: "Cadastro na fase corrente", em: "2026-10-15T00:00:00Z" }],
  emprestimos: [
    { id: "e1", unidade_arquivamento_id: "u1", solicitante: "", finalidade: "Auditoria", devolver_ate: "2026-11-01T00:00:00Z", retirado_em: "2026-10-01T00:00:00Z", devolvido_em: null, observacao: "" },
    { id: "e2", unidade_arquivamento_id: "u1", solicitante: "Rita", finalidade: "Consulta", devolver_ate: "2026-09-01T00:00:00Z", retirado_em: "2026-08-01T00:00:00Z", devolvido_em: "2026-08-20T00:00:00Z", observacao: "" },
  ],
  pode: { gerir: true, bloquear: true },
};

const LOCAIS = [
  { id: "l1", parent_id: null, tipo: "PREDIO", nome: "Arquivo Central", ativo: true, caminho: "Arquivo Central" },
  { id: "l2", parent_id: null, tipo: "PREDIO", nome: "Antigo", ativo: false, caminho: "Antigo" },
];

const posts = (backend: ReturnType<typeof mockBackend>, path: string) => backend.calls.filter((c) => c.path === path && c.method !== "GET");

describe("UnidadePage — ações", () => {
  beforeEach(() => resetNavigation({ id: "u1" }));

  it("localiza, bloqueia, libera, devolve e cancela", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/unidades/u1": { data: FICHA },
      "GET v1/acervo/locais": { data: LOCAIS },
      "PUT v1/acervo/unidades/u1/local": { data: FICHA },
      "POST v1/acervo/unidades/u1/bloqueios": { data: FICHA },
      "POST v1/acervo/bloqueios/b1/liberar": { data: FICHA },
      "POST v1/acervo/emprestimos/e1/devolver": { data: { ...FICHA, emprestimos: [] } },
    });
    renderApp(<UnidadePage />);
    expect(await screen.findByText("Série fora da TTDD em vigor.")).toBeInTheDocument();
    expect(screen.getByText("Acesso restrito")).toBeInTheDocument();
    expect(screen.getByText("Arquivo Central, P2")).toBeInTheDocument();
    expect(screen.getByText("LAI, art. 31")).toBeInTheDocument();
    expect(screen.getByText("Caixa úmida")).toBeInTheDocument();
    expect(screen.getByText(/liberado em .*Trânsito em julgado/)).toBeInTheDocument();
    expect(screen.getByText(/devolvida em/)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Transferir/ })).toBeDisabled();
    expect(screen.queryByRole("button", { name: /Emprestar/ })).not.toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: /Localizar/ }));
    const local = await screen.findByRole("form", { name: "Localização física" });
    expect(within(local).queryByRole("option", { name: "Antigo" })).not.toBeInTheDocument();
    fireEvent.change(within(local).getByLabelText("Posição"), { target: { value: "P9" } });
    fireEvent.submit(local);
    await waitFor(() => expect(posts(backend, "v1/acervo/unidades/u1/local")[0]?.body).toEqual({ local_id: "l1", posicao: "P9" }));

    await userEvent.click(await screen.findByRole("button", { name: /Bloqueio judicial/ }));
    const bloq = await screen.findByRole("form", { name: "Registrar bloqueio judicial" });
    fireEvent.change(within(bloq).getByLabelText("Número do processo judicial *"), { target: { value: "123" } });
    fireEvent.submit(bloq);
    await waitFor(() => expect(posts(backend, "v1/acervo/unidades/u1/bloqueios")[0]?.body).toEqual({ processo: "123", motivo: "" }));

    await userEvent.click(await screen.findByRole("button", { name: "Liberar" }));
    const lib = await screen.findByRole("form", { name: "Liberar bloqueio judicial" });
    fireEvent.change(within(lib).getByLabelText("Motivo da liberação *"), { target: { value: "Decisão" } });
    fireEvent.submit(lib);
    await waitFor(() => expect(posts(backend, "v1/acervo/bloqueios/b1/liberar")).toHaveLength(1));

    await userEvent.click(await screen.findByRole("button", { name: "Registrar devolução" }));
    const dev = await screen.findByRole("form", { name: "Registrar devolução" });
    fireEvent.submit(dev);
    await waitFor(() => expect(posts(backend, "v1/acervo/emprestimos/e1/devolver")[0]?.body).toEqual({ observacao: "" }));
    expect(await screen.findByText("Nenhum empréstimo registrado.")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: /Bloqueio judicial/ }));
    await userEvent.click(await screen.findByRole("button", { name: "Cancelar" }));
    await waitFor(() => expect(screen.queryByRole("form", { name: "Registrar bloqueio judicial" })).not.toBeInTheDocument());
  });

  it("empresta a unidade a quem retira", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/unidades/u1": { data: { ...FICHA, emprestimos: [], bloqueios: [], local_caminho: "", local_id: null } },
      "GET v1/acervo/locais": { data: [] },
      "POST v1/acervo/unidades/u1/emprestimos": { data: FICHA },
    });
    renderApp(<UnidadePage />);
    expect(await screen.findByText("Não localizada")).toBeInTheDocument();
    expect(screen.getByText("Nenhum bloqueio registrado.")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: /Emprestar/ }));
    const form = await screen.findByRole("form", { name: "Emprestar a unidade" });
    await userEvent.type(within(form).getByLabelText("Quem retira *"), `${SOLICITANTE}{Enter}`);
    fireEvent.change(within(form).getByLabelText("Finalidade *"), { target: { value: "Instrução" } });
    fireEvent.change(within(form).getByLabelText("Devolver até *"), { target: { value: "2026-11-01" } });
    fireEvent.submit(form);
    await waitFor(() =>
      expect(posts(backend, "v1/acervo/unidades/u1/emprestimos")[0]?.body).toEqual({ solicitante_id: SOLICITANTE, finalidade: "Instrução", devolver_ate: "2026-11-01" }),
    );
  });

  it("unidade eliminada e sem permissão: só consulta", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"] }),
      "GET v1/acervo/unidades/u1": { data: { ...FICHA, fase: "ELIMINADO", situacao: "ELIMINADA", nivel_acesso: "PUBLICO", pode: { gerir: false, bloquear: false } } },
      "GET v1/acervo/locais": { data: [] },
    });
    renderApp(<UnidadePage />);
    expect(await screen.findByText("Eliminado")).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /Alterar/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Bloqueio judicial/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Liberar" })).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Etiqueta/ })).toHaveAttribute("href", "/acervo/u1/etiqueta");
  });
});

describe("EditarUnidadePage", () => {
  beforeEach(() => resetNavigation({ id: "u1" }));

  it("salva e volta à ficha", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/unidades/u1": { data: FICHA },
      "GET v1/acervo/locais": { data: LOCAIS },
      "GET v1/iam/org-tree": { data: [] },
      "GET v1/atlas/ttdd/*": { data: { codigo: "2.0.02.00.07", descritor: "Pregão Eletrônico", revogada_em: null } },
      "PUT v1/acervo/unidades/u1": { data: FICHA },
    });
    renderApp(<EditarUnidadePage />);
    expect(await screen.findByRole("heading", { name: "Alterar UA-2026-000001" })).toBeInTheDocument();
    fireEvent.submit(await screen.findByRole("form", { name: "Unidade de arquivamento" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/u1"));
    expect(backend.calls.find((c) => c.method === "PUT")?.body).toMatchObject({ codigo_classificacao: "2.0.02.00.07" });
  });

  it("falha ao salvar fica na tela", async () => {
    mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/unidades/u1": { data: FICHA },
      "GET v1/acervo/locais": { data: [] },
      "GET v1/iam/org-tree": { data: [] },
      "GET v1/atlas/ttdd/*": { data: { codigo: "2.0.02.00.07", descritor: "Pregão Eletrônico", revogada_em: null } },
      "PUT v1/acervo/unidades/u1": fail(409, "Unidade eliminada não muda", "CONFLICT"),
    });
    renderApp(<EditarUnidadePage />);
    fireEvent.submit(await screen.findByRole("form", { name: "Unidade de arquivamento" }));
    expect(await screen.findByText("Operação não concluída")).toBeInTheDocument();
    expect(router.push).not.toHaveBeenCalled();
  });
});

describe("NovaGuiaPage", () => {
  beforeEach(() => resetNavigation());

  it("gera a guia com as unidades marcadas", async () => {
    const tree = [
      { id: "e1", nome: "SEMAD", sigla: "SEMAD", slug: "semad", documento: "", ativo: true, unidades: [{ id: "un1", entidade_id: "e1", nome: "Compras", sigla: "", slug: "c", ad_group: "", email: "", departamentos: [] }] },
    ];
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: tree },
      "GET v1/acervo/unidades": page([{ ...UNIDADE, situacao: "TRANSFERIR" }]),
      "POST v1/acervo/guias": { data: { id: "g1" } },
    });
    renderApp(<NovaGuiaPage />);
    await screen.findByRole("option", { name: "Compras" });
    fireEvent.change(screen.getByLabelText(/Unidade de origem/), { target: { value: "un1" } });
    const caixa = await screen.findByRole("checkbox");
    await userEvent.click(caixa);
    await userEvent.click(caixa);
    await userEvent.click(caixa);
    fireEvent.change(screen.getByLabelText("Observação"), { target: { value: " Remessa " } });
    fireEvent.submit(screen.getByRole("form", { name: "Nova guia de transferência" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/guias/g1"));
    expect(backend.calls.find((c) => c.method === "POST")?.body).toEqual({ origem_id: "un1", unidades: ["u1"], observacao: "Remessa" });
  });

  it("sem unidades prontas para transferir", async () => {
    const tree = [{ id: "e1", nome: "SEMAD", sigla: "", slug: "s", documento: "", ativo: true, unidades: [{ id: "un1", entidade_id: "e1", nome: "Compras", sigla: "", slug: "c", ad_group: "", email: "", departamentos: [] }] }];
    mockBackend({ ...identityRoutes(), "GET v1/iam/org-tree": { data: tree }, "GET v1/acervo/unidades": page([]) });
    renderApp(<NovaGuiaPage />);
    await screen.findByRole("option", { name: "Compras" });
    fireEvent.change(screen.getByLabelText(/Unidade de origem/), { target: { value: "un1" } });
    expect(await screen.findByText("Nenhuma unidade pronta para transferir")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Gerar guia" })).toBeDisabled();
  });
});
