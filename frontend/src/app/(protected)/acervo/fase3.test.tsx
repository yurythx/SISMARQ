import { fireEvent, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import ComissoesPage from "./comissoes/page";
import AnexoPage from "./eliminacao/[id]/anexo/[n]/page";
import RelacaoPage from "./eliminacao/[id]/page";
import NovaRelacaoPage from "./eliminacao/nova/page";
import EliminacaoPage from "./eliminacao/page";
import { porExtenso, somarDias, situacaoListagem } from "@/components/acervo/labels";

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

const P1 = "11111111-1111-4111-8111-111111111111";
const P2 = "22222222-2222-4222-8222-222222222222";

const COMISSAO = {
  id: "c1",
  tipo: "SETORIAL",
  entidade_id: "e1",
  nome: "CSPAD da SEMAD",
  portaria: "Portaria nº 12/2026",
  ativa: true,
  membros: [
    { user_id: P1, nome: "Ana Presidente", papel: "PRESIDENTE", area: "ARQUIVOLOGIA_HISTORIA" },
    { user_id: P2, nome: "Bruno Jurídico", papel: "MEMBRO", area: "JURIDICO" },
  ],
  avisos: ["Recomenda-se número ímpar de membros."],
};

const ITEM = {
  id: "u1",
  numero: "UA-2026-000001",
  tipo: "CAIXA",
  codigo_classificacao: "2.0.02.00.07",
  assunto: "Notas de empenho",
  unidade_id: "un1",
  data_inicial: "2015-01-01T00:00:00Z",
  data_final: "2015-12-31T00:00:00Z",
  quantidade: 3,
  especificacao: "",
  observacoes: "Sem valor secundário",
  situacao: "AVALIAR_ELIMINACAO",
  amostra: false,
  retirada: false,
};

const BASE = {
  id: "l1",
  numero: "LE-2026-00001",
  comissao_id: "c1",
  unidade_id: "un1",
  rito: "ORDINARIO",
  situacao: "RASCUNHO",
  ata_signum: "AUSENTE",
  termo_signum: "AUSENTE",
  ccpad_autorizada_em: null as string | null,
  semad_autorizada_em: null as string | null,
  termo_copia_em: null as string | null,
  termo_copia_arquivo: 0,
  termo_copia_semad: 0,
  responsavel_selecao: null,
  parecer_juridico: "Parecer 7/2026",
  contas_exercicios: "2015",
  contas_aprovadas_em: "2016-05-01T00:00:00Z",
  ata_numero: "",
  ata_em: null,
  ata_local: "",
  ata_manifestacoes: "",
  ata_presentes: [] as string[],
  autorizada_em: null as string | null,
  devolucao_motivo: "",
  edital_numero: "",
  edital_publicado_em: null as string | null,
  edital_diario: "",
  amostra_criterios: "",
  amostra_percentual: null as number | null,
  amostra_em: null as string | null,
  termo_eliminado_em: null as string | null,
  termo_acompanhamento: "",
  termo_destinacao_papel: "",
  criada_em: "2026-09-01T00:00:00Z",
  itens: [ITEM, { ...ITEM, id: "u2", numero: "UA-2026-000002", quantidade: 1, especificacao: "pasta" }],
  pedidos: [] as unknown[],
  total_itens: 2,
};

const COMPLETA = {
  ...BASE,
  situacao: "ELIMINADA",
  ata_numero: "1/2026",
  ata_em: "2026-09-01T10:30:00Z",
  ata_local: "Sala da SEMAD",
  ata_manifestacoes: "Nenhuma objeção.",
  ata_presentes: [P1, P2],
  autorizada_em: "2026-09-02T00:00:00Z",
  edital_numero: "3/2026",
  edital_publicado_em: "2026-09-03T00:00:00Z",
  edital_diario: "DIORONDON-e nº 6.300, p. 12",
  amostra_criterios: "Sorteio",
  amostra_percentual: 10,
  amostra_em: "2026-09-10T00:00:00Z",
  termo_eliminado_em: "2026-10-05T00:00:00Z",
  termo_acompanhamento: "Servidor X",
  termo_destinacao_papel: "Doação a cooperativa",
  itens: [{ ...ITEM, amostra: true }, { ...ITEM, id: "u2", numero: "UA-2026-000002" }, { ...ITEM, id: "u3", numero: "UA-2026-000003", retirada: true }],
};

const rotas = (l: unknown, extra = {}) => ({
  ...identityRoutes(),
  "GET v1/iam/org-tree": { data: TREE },
  "GET v1/acervo/comissoes": { data: [COMISSAO] },
  "GET v1/acervo/listagens/l1": { data: l },
  ...extra,
});

beforeEach(() => resetNavigation({ id: "l1", n: "1" }));

describe("rótulos da eliminação", () => {
  it("formata datas e etapas", () => {
    expect(somarDias("2026-09-01T00:00:00Z", 30)).toBe("2026-10-01");
    expect(porExtenso("2026-09-01")).toEqual({ dia: "1", mes: "setembro", ano: "2026" });
    expect(porExtenso(null).mes).toBe("_______");
    expect(porExtenso("2026-13-01").mes).toBe("");
    expect(situacaoListagem("X").label).toBe("X");
  });
});

describe("Relações de eliminação", () => {
  it("lista as relações", async () => {
    mockBackend({ ...identityRoutes(), "GET v1/acervo/listagens": page([{ ...COMPLETA }, { ...BASE, id: "l2", numero: "LE-2026-00002" }]) });
    renderApp(<EliminacaoPage />);
    expect(await screen.findByRole("link", { name: "LE-2026-00001" })).toHaveAttribute("href", "/acervo/eliminacao/l1");
    expect(screen.getByText(/eliminada em 05\/10\/2026/)).toBeInTheDocument();
    expect(screen.getByText("Em seleção")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Nova relação/ })).toBeInTheDocument();
  });

  it("abre a relação com as unidades marcadas", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/unidades": page([ITEM]),
      "POST v1/acervo/listagens": { data: BASE },
    });
    renderApp(<NovaRelacaoPage />);
    await screen.findByRole("option", { name: "DC — Departamento de Compras" });
    fireEvent.change(screen.getByLabelText(/Unidade\/setor/), { target: { value: "un1" } });
    const caixa = await screen.findByRole("checkbox");
    await userEvent.click(caixa);
    await userEvent.click(caixa);
    await userEvent.click(caixa);
    fireEvent.change(screen.getByLabelText(/Parecer jurídico/), { target: { value: "Parecer 1" } });
    fireEvent.change(screen.getByLabelText(/Contas do/), { target: { value: "2015" } });
    fireEvent.submit(screen.getByRole("form", { name: "Nova relação de eliminação" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/eliminacao/l1"));
    const post = backend.calls.find((c) => c.method === "POST");
    expect(post?.body).toMatchObject({ unidade_id: "un1", rito: "ORDINARIO", unidades: ["u1"], parecer_juridico: "Parecer 1", contas_exercicios: "2015" });
  });
});

const chamou = (backend: ReturnType<typeof mockBackend>, method: string, path: string) =>
  waitFor(() => expect(backend.calls.find((c) => c.method === method && c.path === path)).toBeTruthy());

describe("Ficha da relação", () => {
  it("rascunho: parecer, retirada, ata e cancelamento", async () => {
    const backend = mockBackend(
      rotas({ ...BASE, devolucao_motivo: "Conferir datas" }, {
        "PUT v1/acervo/listagens/l1": { data: BASE },
        "POST v1/acervo/listagens/l1/aprovar": { data: { ...BASE, situacao: "APROVADA", ata_em: "2026-09-01T10:00:00Z" } },
        "POST v1/acervo/listagens/l1/cancelar": { data: { ...BASE, situacao: "CANCELADA" } },
      }),
    );
    renderApp(<RelacaoPage />);
    expect(await screen.findByRole("heading", { name: "Relação de eliminação LE-2026-00001" })).toBeInTheDocument();
    expect(await screen.findByText("Secretaria Municipal de Administração · DC — Departamento de Compras")).toBeInTheDocument();
    expect(screen.getByText(/Devolvida pelo Arquivo Público: Conferir datas/)).toBeInTheDocument();

    fireEvent.submit(screen.getByRole("form", { name: "Parecer jurídico e contas" }));
    await chamou(backend, "PUT", "v1/acervo/listagens/l1");
    expect(backend.calls.find((c) => c.method === "PUT")?.body).toMatchObject({ parecer_juridico: "Parecer 7/2026", contas_aprovadas_em: "2016-05-01" });

    await userEvent.click(screen.getByRole("button", { name: "Retirar UA-2026-000001" }));
    await waitFor(() => expect(backend.calls.filter((c) => c.method === "PUT")[1]?.body).toMatchObject({ retirar: ["u1"] }));

    await screen.findByText(/Ana Presidente — Presidente/);
    fireEvent.change(screen.getByLabelText("Ata nº *"), { target: { value: "1/2026" } });
    fireEvent.change(screen.getByLabelText("Data e hora *"), { target: { value: "2026-09-01T10:00" } });
    fireEvent.submit(screen.getByRole("form", { name: "Ata da reunião" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/aprovar");
    expect(backend.calls.find((c) => c.path.endsWith("/aprovar"))?.body).toMatchObject({ ata_numero: "1/2026", presentes: [P1, P2] });
    expect(await screen.findByText("Aprovada pela CSPAD")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Anexo IV/ })).toHaveAttribute("href", "/acervo/eliminacao/l1/anexo/4");

    await userEvent.click(screen.getByRole("button", { name: "Cancelar relação" }));
    await userEvent.click(await screen.findByRole("button", { name: "Cancelar a relação" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/cancelar");
  });

  it("aprovada: o Arquivo Público autoriza ou devolve", async () => {
    const backend = mockBackend(
      rotas({ ...BASE, situacao: "APROVADA" }, {
        "POST v1/acervo/listagens/l1/devolver": { data: { ...BASE, situacao: "APROVADA" } },
        "POST v1/acervo/listagens/l1/autorizar": { data: { ...BASE, situacao: "AUTORIZADA" } },
      }),
    );
    renderApp(<RelacaoPage />);
    fireEvent.change(await screen.findByLabelText("Motivo da devolução"), { target: { value: "Rever" } });
    fireEvent.submit(screen.getByRole("form", { name: "Devolver à CSPAD" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/devolver");
    await userEvent.click(screen.getByRole("button", { name: "Autorizar eliminação" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/autorizar");
    expect(await screen.findByText(/Publique o edital/)).toBeInTheDocument();
  });

  it("art. 35: CCPAD e SEMAD", async () => {
    const backend = mockBackend(
      rotas({ ...BASE, situacao: "APROVADA", rito: "ART_35" }, {
        "POST v1/acervo/listagens/l1/autorizar-ccpad": { data: { ...BASE, situacao: "APROVADA", rito: "ART_35", ccpad_autorizada_em: "2026-09-02T00:00:00Z" } },
      }),
    );
    renderApp(<RelacaoPage />);
    await userEvent.click(await screen.findByRole("button", { name: "Autorizar pela CCPAD" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/autorizar-ccpad");
    const sem = mockBackend(
      rotas({ ...BASE, situacao: "AGUARDANDO_SEMAD", rito: "ART_35", ccpad_autorizada_em: "2026-09-02T00:00:00Z" }, {
        "POST v1/acervo/listagens/l1/autorizar-semad": { data: { ...BASE, situacao: "AUTORIZADA", rito: "ART_35" } },
      }),
    );
    renderApp(<RelacaoPage />);
    await userEvent.click(await screen.findByRole("button", { name: "Aprovar pela SEMAD" }));
    await chamou(sem, "POST", "v1/acervo/listagens/l1/autorizar-semad");
  });

  it("autorizada: registra o edital", async () => {
    const backend = mockBackend(
      rotas({ ...BASE, situacao: "AUTORIZADA" }, {
        "POST v1/acervo/listagens/l1/edital": { data: { ...BASE, situacao: "EM_EDITAL", edital_numero: "3/2026", edital_publicado_em: "2026-09-03T00:00:00Z" } },
      }),
    );
    renderApp(<RelacaoPage />);
    fireEvent.change(await screen.findByLabelText("Edital nº *"), { target: { value: "3/2026" } });
    fireEvent.change(screen.getByLabelText("Publicado em *"), { target: { value: "2026-09-03" } });
    fireEvent.submit(screen.getByRole("form", { name: "Edital de ciência" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/edital");
    expect(await screen.findByText("03/10/2026")).toBeInTheDocument();
    expect(screen.getByText("04/10/2026")).toBeInTheDocument();
  });

  it("em edital: pedidos, decisões e amostragem", async () => {
    const pedidos = [
      { id: "p1", unidade_arquivamento_id: "u1", requerente: "Maria", pedido: "Cópia", recebido_em: "2026-09-05T00:00:00Z", situacao: "PENDENTE", decisao: "", decidido_em: null },
      { id: "p2", unidade_arquivamento_id: null, requerente: "João", pedido: "Vista", recebido_em: "2026-09-06T00:00:00Z", situacao: "PENDENTE", decisao: "", decidido_em: null },
      { id: "p3", unidade_arquivamento_id: null, requerente: "Rita", pedido: "x", recebido_em: "2026-09-06T00:00:00Z", situacao: "DEFERIDO", decisao: "Ok", decidido_em: null },
      { id: "p4", unidade_arquivamento_id: null, requerente: "Caio", pedido: "y", recebido_em: "2026-09-06T00:00:00Z", situacao: "INDEFERIDO", decisao: "Não", decidido_em: null },
    ];
    const l = { ...BASE, situacao: "EM_EDITAL", edital_numero: "3/2026", edital_publicado_em: "2026-09-03T00:00:00Z", pedidos };
    const backend = mockBackend(
      rotas(l, {
        "POST v1/acervo/listagens/l1/pedidos": { data: l },
        "POST v1/acervo/listagens/l1/pedidos/p1/decidir": { data: l },
        "POST v1/acervo/listagens/l1/pedidos/p2/decidir": { data: l },
        "POST v1/acervo/listagens/l1/amostragem": { data: { ...l, amostra_em: "2026-09-10T00:00:00Z", amostra_percentual: 10, amostra_criterios: "Sorteio" } },
      }),
    );
    renderApp(<RelacaoPage />);
    expect(await screen.findByText("Deferido")).toBeInTheDocument();
    expect(screen.getByText("Indeferido")).toBeInTheDocument();
    expect(screen.getByText(/· UA-2026-000001/)).toBeInTheDocument();

    const maria = screen.getByRole("form", { name: "Decidir pedido de Maria" });
    fireEvent.change(within(maria).getByLabelText("Fundamentação *"), { target: { value: "Legítimo" } });
    await userEvent.click(within(maria).getByRole("button", { name: "Deferir" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/pedidos/p1/decidir");
    expect(backend.calls.find((c) => c.path.endsWith("p1/decidir"))?.body).toEqual({ deferido: true, retirar: true, decisao: "Legítimo" });

    const joao = screen.getByRole("form", { name: "Decidir pedido de João" });
    fireEvent.change(within(joao).getByLabelText("Fundamentação *"), { target: { value: "Sem legitimidade" } });
    await userEvent.click(within(joao).getByRole("button", { name: "Indeferir" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/pedidos/p2/decidir");
    expect(backend.calls.find((c) => c.path.endsWith("p2/decidir"))?.body).toEqual({ deferido: false, retirar: false, decisao: "Sem legitimidade" });

    fireEvent.change(screen.getByLabelText("Requerente *"), { target: { value: "Pedro" } });
    fireEvent.change(screen.getByLabelText("Recebido em *"), { target: { value: "2026-09-07" } });
    fireEvent.change(screen.getByLabelText("Pedido *"), { target: { value: "Desentranhamento" } });
    fireEvent.submit(screen.getByRole("form", { name: "Registrar pedido" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/pedidos");
    expect(backend.calls.find((c) => c.path.endsWith("/pedidos"))?.body).toMatchObject({ requerente: "Pedro", unidade_arquivamento_id: null });

    fireEvent.change(screen.getByLabelText("Critérios da CSPAD *"), { target: { value: "Sorteio" } });
    fireEvent.submit(screen.getByRole("form", { name: "Amostragem" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/amostragem");
    expect(backend.calls.find((c) => c.path.endsWith("/amostragem"))?.body).toEqual({ criterios: "Sorteio", percentual: 10 });
    expect(await screen.findByRole("form", { name: "Termo de eliminação" })).toBeInTheDocument();
  });

  it("em edital com amostra: registra o termo", async () => {
    const l = { ...COMPLETA, situacao: "EM_EDITAL", termo_eliminado_em: null, pedidos: [] };
    const backend = mockBackend(rotas(l, { "POST v1/acervo/listagens/l1/termo": { data: COMPLETA } }));
    renderApp(<RelacaoPage />);
    expect(await screen.findByText("Nenhum pedido registrado.")).toBeInTheDocument();
    expect(screen.getByText("Amostra — guarda permanente")).toBeInTheDocument();
    expect(screen.getByText("Retirada")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Cancelar relação" })).not.toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("Eliminado em *"), { target: { value: "2026-10-05" } });
    fireEvent.change(screen.getByLabelText("Servidor que acompanhou *"), { target: { value: "Servidor X" } });
    fireEvent.submit(screen.getByRole("form", { name: "Termo de eliminação" }));
    await chamou(backend, "POST", "v1/acervo/listagens/l1/termo");
    expect((await screen.findAllByText("Eliminada")).length).toBeGreaterThan(1);
    expect(screen.getByRole("link", { name: /Anexo II —/ })).toBeInTheDocument();
  });

  it("só leitura: sem ações", async () => {
    mockBackend({
      ...rotas({ ...BASE, situacao: "APROVADA" }),
      ...identityRoutes({ permissions: ["acervo:read"], scopes: [{ perfil: "consulta", origem: "manual", permissions: ["acervo:read"] }] }),
      "GET v1/acervo/comissoes": { data: [] },
      "GET v1/iam/org-tree": { data: [] },
    });
    renderApp(<RelacaoPage />);
    expect(await screen.findByRole("heading", { name: /LE-2026-00001/ })).toBeInTheDocument();
    expect(screen.getAllByText("A eliminar")).toHaveLength(2);
    expect(screen.queryByRole("button", { name: "Autorizar eliminação" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Cancelar relação" })).not.toBeInTheDocument();
  });
});

describe("Anexos I a IV", () => {
  it.each([
    ["1", "Relação de eliminação de documentos", /DATAS-LIMITE GERAIS/],
    ["2", "Termo de eliminação de documentos", /procedeu à eliminação de 1 unidade\(s\)/],
    ["3", "Edital de ciência de eliminação de documentos", /a partir do 30º dia subsequente/],
    ["4", "Ata de reunião para aprovação de eliminação de documentos", /às 10h30, nas dependências do\(a\) Sala da SEMAD/],
  ])("Anexo %s", async (n, titulo, texto) => {
    resetNavigation({ id: "l1", n });
    mockBackend(rotas(COMPLETA));
    renderApp(<AnexoPage />);
    expect(await screen.findByRole("heading", { name: titulo })).toBeInTheDocument();
    expect(screen.getByText(texto)).toBeInTheDocument();
    if (n !== "4") expect(await screen.findByText("Secretaria Municipal de Administração", { selector: "dd" })).toBeInTheDocument();
    if (n === "1") expect(screen.getByText(/Titular da Superintendência de Arquivo Público/)).toBeInTheDocument();
    if (n === "4") expect(await screen.findByText("Bruno Jurídico")).toBeInTheDocument();
  });

  it("Anexos sem os dados da etapa e documento inexistente", async () => {
    resetNavigation({ id: "l1", n: "4" });
    mockBackend({ ...rotas({ ...BASE, itens: [], ata_presentes: [P1] }), "GET v1/acervo/comissoes": { data: [] } });
    const { unmount } = renderApp(<AnexoPage />);
    expect(await screen.findByText("Sem manifestações.")).toBeInTheDocument();
    expect(screen.getByText(P1)).toBeInTheDocument();
    unmount();

    resetNavigation({ id: "l1", n: "1" });
    renderApp(<AnexoPage />);
    expect(await screen.findByText(/____\/____\/______/)).toBeInTheDocument();
    window.print = vi.fn();
    await userEvent.click(screen.getByRole("button", { name: /Imprimir/ }));
    expect(window.print).toHaveBeenCalled();
  });

  it("Anexo III sem portaria e documento inexistente", async () => {
    resetNavigation({ id: "l1", n: "3" });
    mockBackend({ ...rotas({ ...COMPLETA, edital_diario: "" }), "GET v1/acervo/comissoes": { data: [] } });
    const { unmount } = renderApp(<AnexoPage />);
    expect(await screen.findByText(/designado\(a\) pela Portaria nº ____/)).toBeInTheDocument();
    unmount();
    resetNavigation({ id: "l1", n: "9" });
    renderApp(<AnexoPage />);
    expect(await screen.findByText("Documento não encontrado")).toBeInTheDocument();
  });
});

describe("Comissões", () => {
  it("lista e cadastra a comissão", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/iam/org-tree": { data: TREE },
      "GET v1/acervo/comissoes": { data: [COMISSAO, { ...COMISSAO, id: "c2", tipo: "CENTRAL", entidade_id: null, nome: "CCPAD", portaria: "", ativa: false, avisos: [] }] },
      "POST v1/acervo/comissoes": { data: COMISSAO },
      "PUT v1/acervo/comissoes/c1": { data: COMISSAO },
    });
    renderApp(<ComissoesPage />);
    expect(await screen.findByRole("heading", { name: "CSPAD da SEMAD" })).toBeInTheDocument();
    expect(screen.getByText("Recomenda-se número ímpar de membros.")).toBeInTheDocument();
    expect(screen.getByText(/Órgão Central do SISMARQ/)).toBeInTheDocument();
    expect(await screen.findByText(/Secretaria Municipal de Administração · Portaria/)).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "Nova comissão" }));
    fireEvent.change(screen.getByLabelText("Nome *"), { target: { value: "CSPAD nova" } });
    fireEvent.change(screen.getByLabelText("Órgão/entidade *"), { target: { value: "e1" } });
    const busca = screen.getByLabelText("Membros");
    await userEvent.type(busca, `${P1}{Enter}`);
    fireEvent.change(screen.getByLabelText("Papel"), { target: { value: "PRESIDENTE" } });
    fireEvent.change(screen.getByLabelText("Área"), { target: { value: "JURIDICO" } });
    fireEvent.submit(screen.getByRole("form", { name: "Comissão" }));
    await chamou(backend, "POST", "v1/acervo/comissoes");
    expect(backend.calls.find((c) => c.method === "POST")?.body).toEqual({
      tipo: "SETORIAL",
      entidade_id: "e1",
      nome: "CSPAD nova",
      portaria: "",
      ativa: true,
      membros: [{ user_id: P1, papel: "PRESIDENTE", area: "JURIDICO" }],
    });

    await userEvent.click(screen.getAllByRole("button", { name: "Editar" })[1]!);
    expect(screen.getByLabelText("Nome *")).toHaveValue("CCPAD");
    expect(screen.queryByLabelText("Órgão/entidade *")).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Cancelar" }));
    await userEvent.click(screen.getAllByRole("button", { name: "Editar" })[0]!);
    await userEvent.click(screen.getByRole("button", { name: "Remover Bruno Jurídico" }));
    fireEvent.change(screen.getByLabelText("Tipo *"), { target: { value: "CENTRAL" } });
    fireEvent.submit(screen.getByRole("form", { name: "Comissão" }));
    await chamou(backend, "PUT", "v1/acervo/comissoes/c1");
    expect(backend.calls.find((c) => c.method === "PUT")?.body).toMatchObject({ tipo: "CENTRAL", entidade_id: null, membros: [{ user_id: P1 }] });
  });

  it("sem acervo:central só consulta", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"], scopes: [] }),
      "GET v1/iam/org-tree": { data: [] },
      "GET v1/acervo/comissoes": { data: [{ ...COMISSAO, membros: [{ ...COMISSAO.membros[0]!, nome: "", papel: "X", area: "Y" }] }] },
    });
    renderApp(<ComissoesPage />);
    expect(await screen.findByText(`${P1} — X, Y`)).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Nova comissão" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Editar" })).not.toBeInTheDocument();
  });
});
