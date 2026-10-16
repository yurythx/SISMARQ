import { fireEvent, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { recolhiveis, situacaoRecolhimento } from "@/components/acervo/labels";
import type { AcervoUnidade } from "@/lib/sismarq/types";
import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import UnidadePage from "./[id]/page";
import RecolhimentoPage from "./recolhimentos/[id]/page";
import NovoRecolhimentoPage from "./recolhimentos/novo/page";
import RecolhimentosPage from "./recolhimentos/page";

const TREE = [
  {
    id: "e1",
    nome: "Secretaria Municipal de Administração",
    sigla: "SEMAD",
    slug: "semad",
    documento: "",
    ativo: true,
    unidades: [{ id: "un1", entidade_id: "e1", nome: "Departamento de Compras", sigla: "DC", slug: "dc", ad_group: "", email: "", departamentos: [] }],
  },
];

const UNIDADE = {
  id: "u1",
  numero: "UA-2026-000001",
  tipo: "CAIXA",
  codigo_classificacao: "2.0.01.00.02",
  assunto: "Leis e decretos",
  unidade_id: "un1",
  data_inicial: "2015-01-01T00:00:00Z",
  data_final: "2015-12-31T00:00:00Z",
  quantidade: 2,
  especificacao: "",
  observacoes: "",
  fase: "INTERMEDIARIA",
  destinacao: "GUARDA_PERMANENTE",
  data_evento: null,
  contas_aprovadas_em: null,
  nivel_acesso: "PUBLICO",
  hipotese_legal: "",
  local_id: null,
  posicao: "",
  created_at: "2026-10-15T00:00:00Z",
  updated_at: "2026-10-15T00:00:00Z",
  fim_corrente: null,
  fim_intermediaria: null,
  situacao: "RECOLHER",
  local_caminho: "",
  recolhida_em: null,
  recolhimento_id: null,
};

const PREPARO = { classificados: true, avaliados: true, organizados: true, higienizados: true, acondicionados: true };

const RECOLHIMENTO = {
  id: "r1",
  numero: "RC-2026-00001",
  origem_id: "un1",
  motivo: "EXTINCAO",
  fundamento: "Lei 11.000/2026",
  preparo: PREPARO,
  historico: "Autarquia criada em 1990.",
  observacao: "Caixas numeradas",
  situacao: "ENVIADO",
  criado_em: "2026-10-15T00:00:00Z",
  recebido_em: null as string | null,
  registro_entrada: "",
  local_id: null,
  devolucao_motivo: "",
  itens: [UNIDADE, { ...UNIDADE, id: "u2", numero: "UA-2026-000002", nivel_acesso: "RESTRITO", hipotese_legal: "LAI, art. 31", especificacao: "pastas" }],
  total_itens: 2,
};

const u = (extra: Partial<AcervoUnidade>) => ({ ...UNIDADE, ...extra }) as AcervoUnidade;

beforeEach(() => resetNavigation({ id: "r1" }));

describe("rótulos do recolhimento", () => {
  it("filtra o que pode ir em cada motivo", () => {
    const lista = [
      u({ id: "a", situacao: "RECOLHER" }),
      u({ id: "b", situacao: "PERMANENTE", fase: "PERMANENTE" }),
      u({ id: "c", situacao: "RECOLHIDA", fase: "PERMANENTE", recolhida_em: "2026-01-01" }),
      u({ id: "d", situacao: "NO_CORRENTE", fase: "CORRENTE" }),
      u({ id: "e", situacao: "AVALIAR_ELIMINACAO", destinacao: "ELIMINACAO" }),
      u({ id: "f", situacao: "ELIMINADA", fase: "ELIMINADO" }),
    ];
    expect(recolhiveis(lista, "REGULAR").map((x) => x.id)).toEqual(["a", "b"]);
    expect(recolhiveis(lista, "EXTINCAO").map((x) => x.id)).toEqual(["a", "b", "d"]);
    expect(situacaoRecolhimento("X").label).toBe("X");
  });
});

describe("Recolhimentos", () => {
  it("lista os recolhimentos", async () => {
    mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/recolhimentos": page([{ ...RECOLHIMENTO, situacao: "RECEBIDO", registro_entrada: "RE-2026-00001" }]),
    });
    renderApp(<RecolhimentosPage />);
    expect(await screen.findByRole("link", { name: "RC-2026-00001" })).toHaveAttribute("href", "/acervo/recolhimentos/r1");
    expect(screen.getByText(/registro de entrada RE-2026-00001/)).toBeInTheDocument();
    expect(screen.getByText("Recebido")).toBeInTheDocument();
  });

  it("envia o recolhimento: motivo, unidades e preparo do art. 53", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/unidades": page([UNIDADE, u({ id: "d", numero: "UA-2026-000009", situacao: "NO_CORRENTE", fase: "CORRENTE" })]),
      "POST v1/acervo/recolhimentos": { data: RECOLHIMENTO },
    });
    renderApp(<NovoRecolhimentoPage />);
    await screen.findByRole("option", { name: "DC — Departamento de Compras" });
    fireEvent.change(screen.getByLabelText(/Unidade de origem/), { target: { value: "un1" } });
    expect(await screen.findByText(/UA-2026-000001/)).toBeInTheDocument();
    expect(screen.queryByText(/UA-2026-000009/)).not.toBeInTheDocument();

    // Na extinção entra também a unidade ainda no corrente (arts. 50 e 51).
    fireEvent.change(screen.getByLabelText("Motivo *"), { target: { value: "EXTINCAO" } });
    expect(await screen.findByText(/UA-2026-000009/)).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText(/Ato de extinção/), { target: { value: "Lei 11.000/2026" } });
    const caixas = screen.getAllByRole("checkbox", { name: /UA-2026/ });
    await userEvent.click(caixas[0]!);
    await userEvent.click(caixas[1]!);
    await userEvent.click(caixas[1]!);
    const enviar = screen.getByRole("button", { name: "Enviar ao Arquivo Público" });
    expect(enviar).toBeDisabled();
    for (const nome of ["Classificados", "Avaliados", "Organizados", "Higienizados", "Acondicionados"]) {
      await userEvent.click(screen.getByRole("checkbox", { name: nome }));
    }
    expect(enviar).toBeEnabled();
    fireEvent.change(screen.getByLabelText(/História administrativa/), { target: { value: " Autarquia de 1990 " } });
    fireEvent.submit(screen.getByRole("form", { name: "Novo recolhimento" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/recolhimentos/r1"));
    expect(backend.calls.find((c) => c.method === "POST")?.body).toEqual({
      origem_id: "un1",
      motivo: "EXTINCAO",
      fundamento: "Lei 11.000/2026",
      preparo: PREPARO,
      historico: "Autarquia de 1990",
      observacao: "",
      unidades: ["u1"],
    });
  });

  it("sem unidades para recolher, nos dois motivos", async () => {
    mockBackend({ ...identityRoutes(), "GET v1/iam/org-tree": { data: TREE }, "GET v1/acervo/unidades": page([]) });
    renderApp(<NovoRecolhimentoPage />);
    await screen.findByRole("option", { name: "DC — Departamento de Compras" });
    fireEvent.change(screen.getByLabelText(/Unidade de origem/), { target: { value: "un1" } });
    expect(await screen.findByText(/prazos cumpridos e destinação guarda permanente/)).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("Motivo *"), { target: { value: "DESESTATIZACAO" } });
    expect(await screen.findByText(/guarda permanente ainda no órgão/)).toBeInTheDocument();
  });
});

describe("Ficha do recolhimento", () => {
  it("o Arquivo Público registra a entrada com o local", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/recolhimentos/r1": { data: RECOLHIMENTO },
      "GET v1/acervo/locais": {
        data: [
          { id: "l1", parent_id: null, tipo: "DEPOSITO", nome: "Depósito", ativo: true, caminho: "Arquivo Público > Depósito" },
          { id: "l2", parent_id: null, tipo: "SALA", nome: "Velha", ativo: false, caminho: "Velha" },
        ],
      },
      "POST v1/acervo/recolhimentos/r1/receber": {
        data: { ...RECOLHIMENTO, situacao: "RECEBIDO", registro_entrada: "RE-2026-00001", recebido_em: "2026-10-16T00:00:00Z" },
      },
    });
    renderApp(<RecolhimentoPage />);
    expect(await screen.findByRole("heading", { name: "Recolhimento RC-2026-00001" })).toBeInTheDocument();
    expect(await screen.findByText("Secretaria Municipal de Administração")).toBeInTheDocument();
    expect(screen.getByText(/Classificados: sim/)).toBeInTheDocument();
    expect(screen.getByText("Autarquia criada em 1990.")).toBeInTheDocument();
    expect(screen.getByText(/4 volume\(s\)/)).toBeInTheDocument();
    expect(screen.getByText("restrito — LAI, art. 31")).toBeInTheDocument();
    const entrada = screen.getByRole("form", { name: "Registrar entrada" });
    await waitFor(() => expect(within(entrada).getByRole("option", { name: "Arquivo Público > Depósito" })).toBeInTheDocument());
    expect(within(entrada).queryByRole("option", { name: "Velha" })).not.toBeInTheDocument();
    fireEvent.change(within(entrada).getByLabelText("Guardar em"), { target: { value: "l1" } });
    fireEvent.submit(entrada);
    await waitFor(() => expect(backend.calls.find((c) => c.path.endsWith("/receber"))?.body).toEqual({ local_id: "l1" }));
    expect(await screen.findByRole("heading", { name: "Registro de entrada RE-2026-00001" })).toBeInTheDocument();
    expect(screen.getByText(/entrada em 16\/10\/2026/)).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Registrar entrada" })).not.toBeInTheDocument();
  });

  it("o Arquivo Público devolve; sem local a entrada vai sem localização", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/recolhimentos/r1": { data: { ...RECOLHIMENTO, historico: "", observacao: "", fundamento: "" } },
      "GET v1/acervo/locais": { data: [] },
      "POST v1/acervo/recolhimentos/r1/devolver": { data: { ...RECOLHIMENTO, situacao: "DEVOLVIDO", devolucao_motivo: "Faltam etiquetas" } },
    });
    renderApp(<RecolhimentoPage />);
    fireEvent.change(await screen.findByLabelText("O que falta (devolução)"), { target: { value: " Faltam etiquetas " } });
    fireEvent.submit(screen.getByRole("form", { name: "Devolver à origem" }));
    await waitFor(() => expect(backend.calls.find((c) => c.path.endsWith("/devolver"))?.body).toEqual({ motivo: "Faltam etiquetas" }));
    expect(await screen.findByText(/Devolvido pelo Arquivo Público:/)).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Cancelar recolhimento" })).not.toBeInTheDocument();
  });

  it("a origem cancela e imprime; sem acervo:central não decide", async () => {
    const backend = mockBackend({
      ...identityRoutes({ permissions: ["acervo:read", "acervo:manage"], scopes: [{ perfil: "agente", origem: "manual", unidade_id: "un1", permissions: ["acervo:read", "acervo:manage"] }] }),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/recolhimentos/r1": { data: { ...RECOLHIMENTO, motivo: "OUTRO" } },
      "POST v1/acervo/recolhimentos/r1/cancelar": { data: { ...RECOLHIMENTO, situacao: "CANCELADO" } },
    });
    renderApp(<RecolhimentoPage />);
    expect(await screen.findByText("OUTRO")).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Registrar entrada" })).not.toBeInTheDocument();
    window.print = vi.fn();
    await userEvent.click(screen.getByRole("button", { name: /Imprimir/ }));
    expect(window.print).toHaveBeenCalled();
    await userEvent.click(screen.getByRole("button", { name: "Cancelar recolhimento" }));
    await userEvent.click(await screen.findByRole("button", { name: "Cancelar o recolhimento" }));
    await waitFor(() => expect(backend.calls.some((c) => c.path.endsWith("/cancelar"))).toBe(true));
    expect(await screen.findByText("Cancelado")).toBeInTheDocument();
  });
});

describe("Ficha da unidade recolhida", () => {
  it("mostra a guarda do Arquivo Público e o caminho até o registro", async () => {
    resetNavigation({ id: "u1" });
    mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/unidades/u1": {
        data: {
          ...UNIDADE,
          fase: "PERMANENTE",
          situacao: "RECOLHIDA",
          recolhida_em: "2026-10-16T00:00:00Z",
          recolhimento_id: "r1",
          serie: { codigo: UNIDADE.codigo_classificacao, descritor: "Leis", vigente: true, destinacao: "GUARDA_PERMANENTE" },
          prazos: { fim_corrente: null, fim_intermediaria: null, situacao: "RECOLHIDA", regra_corrente: "", regra_intermediaria: "" },
          bloqueios: [],
          movimentacoes: [],
          emprestimos: [],
          pode: { gerir: false, bloquear: false },
        },
      },
      "GET v1/acervo/locais": { data: [] },
    });
    renderApp(<UnidadePage />);
    expect(await screen.findByText("No Arquivo Público")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /ver o recolhimento e o registro de entrada/ })).toHaveAttribute("href", "/acervo/recolhimentos/r1");
  });
});
