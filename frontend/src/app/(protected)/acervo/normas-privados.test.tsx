import { fireEvent, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { situacaoNorma, situacaoPrivado } from "@/components/acervo/labels";
import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation, router } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);

import NormaPage from "./normas/[id]/page";
import NormasPage from "./normas/page";
import PrivadoPage from "./privados/[id]/page";
import PrivadosPage from "./privados/page";

const R1 = {
  versao: 1,
  ementa: "Procedimentos de recolhimento",
  texto: "**Art. 1º** Recolher em março.",
  motivo: "Redação original",
  situacao: "PUBLICADA",
  publicada_em: "2026-03-01T00:00:00Z",
  diario: "DIORONDON-e",
  vigente_ate: "2027-06-01T00:00:00Z",
  criada_em: "2026-02-01T00:00:00Z",
};
const R2 = { ...R1, versao: 2, texto: "Art. 1º Recolher em junho.", motivo: "Alterada pela IN 2/2027", publicada_em: "2027-06-01T00:00:00Z", vigente_ate: null };
const RASC = { ...R1, versao: 3, texto: "Art. 1º Recolher em julho.", motivo: "IN 4/2027", situacao: "RASCUNHO", publicada_em: null, vigente_ate: null };

const NORMA = {
  id: "n1",
  numero: "1/2026",
  situacao: "VIGENTE",
  ementa: "Procedimentos de recolhimento",
  versao: 2,
  publicada_em: R2.publicada_em,
  revogada_em: null as string | null,
  revogacao_ato: "",
  tem_rascunho: true,
  criada_em: "2026-02-01T00:00:00Z",
  redacoes: [RASC, R2, R1],
  total_redacoes: 2,
};

const PRIVADO = {
  id: "p1",
  numero: "AP-2026-00001",
  titular: "Família Silva",
  tipo_titular: "PF",
  documento: "123",
  contato: "",
  descricao: "Fotografias da fundação",
  datas_limite: "1953-1970",
  mensuracao: "3 caixas",
  local_guarda: "",
  servico_publico: false,
  situacao: "CADASTRADO",
  parecer_comissao: "",
  parecer_texto: "",
  parecer_favoravel: null as boolean | null,
  parecer_em: null as string | null,
  decreto_numero: "",
  decreto_em: null as string | null,
  decreto_diario: "",
  acesso_autorizado: false,
  acesso_condicoes: "",
  guarda: "TITULAR",
  guarda_em: null as string | null,
  guarda_termo: "",
  convenio: "",
  criado_em: "2026-10-01T00:00:00Z",
  eventos: [{ tipo: "cadastrado", descricao: "Cadastrado; aguarda o parecer", em: "2026-10-01T10:00:00Z" }],
  alienacoes: [] as unknown[],
};

const DECLARADO = {
  ...PRIVADO,
  situacao: "DECLARADO",
  parecer_favoravel: true,
  parecer_em: "2026-09-01T00:00:00Z",
  parecer_texto: "Conjunto único",
  decreto_numero: "Decreto 10/2026",
  decreto_em: "2026-10-01T00:00:00Z",
  decreto_diario: "DIORONDON-e",
  acesso_autorizado: true,
  acesso_condicoes: "Com agendamento",
  guarda: "DEPOSITADO",
  guarda_em: "2026-10-05T00:00:00Z",
  guarda_termo: "Termo 1",
  alienacoes: [
    { id: "a1", notificado_em: "2026-10-10T00:00:00Z", adquirente: "Colecionador", condicoes: "R$ 10 mil", decisao: "PENDENTE", fundamento: "", decidido_em: null, prazo_final: "2026-12-09T00:00:00Z", vencida: false },
    { id: "a0", notificado_em: "2026-01-10T00:00:00Z", adquirente: "", condicoes: "", decisao: "PENDENTE", fundamento: "", decidido_em: null, prazo_final: "2026-03-11T00:00:00Z", vencida: true },
    { id: "az", notificado_em: "2025-01-10T00:00:00Z", adquirente: "", condicoes: "", decisao: "SEM_INTERESSE", fundamento: "Sem dotação", decidido_em: "2025-02-01T00:00:00Z", prazo_final: "2025-03-11T00:00:00Z", vencida: false },
  ],
};

const posted = (backend: ReturnType<typeof mockBackend>, suffix: string) => backend.calls.find((c) => c.method !== "GET" && c.path.endsWith(suffix));

beforeEach(() => resetNavigation({ id: "n1" }));

describe("rótulos", () => {
  it("normas e privados", () => {
    expect(situacaoNorma("X").label).toBe("X");
    expect(situacaoPrivado("X").label).toBe("X");
    expect(situacaoNorma("VIGENTE").tone).toBe("success");
  });
});

describe("Instruções normativas — gestão", () => {
  it("lista, busca e cria em rascunho", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/normas": page([NORMA, { ...NORMA, id: "n2", numero: "2/2026", situacao: "RASCUNHO", versao: 0, publicada_em: null }]),
      "POST v1/acervo/normas": { data: { ...NORMA, id: "n3" } },
    });
    renderApp(<NormasPage />);
    expect(await screen.findByRole("link", { name: "IN nº 1/2026" })).toHaveAttribute("href", "/acervo/normas/n1");
    expect(screen.getByText(/Nunca publicada/)).toBeInTheDocument();
    expect(screen.getAllByText(/há redação em rascunho/)).toHaveLength(2);
    fireEvent.change(screen.getByLabelText("Buscar"), { target: { value: "recolh" } });
    await waitFor(() => expect(backend.calls.some((c) => c.query.get("q") === "recolh")).toBe(true));
    await userEvent.click(screen.getByRole("button", { name: "Nova instrução normativa" }));
    fireEvent.change(screen.getByLabelText("Número *"), { target: { value: "3/2026" } });
    fireEvent.change(screen.getByLabelText("Ementa *"), { target: { value: " Eliminação " } });
    fireEvent.change(screen.getByLabelText("Texto (Markdown) *"), { target: { value: "Art. 1º" } });
    fireEvent.submit(screen.getByRole("form", { name: "Nova instrução normativa" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/normas/n3"));
    expect(posted(backend, "/normas")?.body).toEqual({ numero: "3/2026", ementa: "Eliminação", texto: "Art. 1º" });
    await userEvent.click(screen.getByRole("button", { name: "Cancelar" }));
    expect(screen.queryByRole("form", { name: "Nova instrução normativa" })).not.toBeInTheDocument();
  });

  it("sem acervo:central só aponta a página pública", async () => {
    mockBackend({ ...identityRoutes({ permissions: ["acervo:read"], scopes: [] }) });
    renderApp(<NormasPage />);
    expect(await screen.findByText(/A gestão das instruções normativas é do Arquivo Público/)).toBeInTheDocument();
  });

  it("rascunho da alteração: salva, publica e mostra o histórico", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/normas/n1": { data: NORMA },
      "PUT v1/acervo/normas/n1/rascunho": { data: NORMA },
      "POST v1/acervo/normas/n1/publicar": {
        data: { ...NORMA, versao: 3, tem_rascunho: false, redacoes: [{ ...RASC, situacao: "PUBLICADA", publicada_em: "2027-07-01T00:00:00Z" }, { ...R2, vigente_ate: "2027-07-01T00:00:00Z" }, R1] },
      },
    });
    renderApp(<NormaPage />);
    expect(await screen.findByRole("heading", { name: "Instrução Normativa nº 1/2026" })).toBeInTheDocument();
    expect(screen.getByText(/Redação 2 em vigor desde 01\/06\/2027/)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Página pública/ })).toHaveAttribute("href", "/normas/n1");
    expect(screen.getByRole("heading", { name: "Redação 3 em rascunho" })).toBeInTheDocument();
    expect(screen.getByText(/Redação 1 — de 01\/03\/2026 a 01\/06\/2027/)).toBeInTheDocument();
    expect(screen.getByLabelText("Texto (Markdown) *")).toHaveValue("Art. 1º Recolher em julho.");

    fireEvent.submit(screen.getByRole("form", { name: "Redação" }));
    await waitFor(() =>
      expect(posted(backend, "/rascunho")?.body).toEqual({ ementa: "Procedimentos de recolhimento", texto: "Art. 1º Recolher em julho.", motivo: "IN 4/2027" }),
    );
    fireEvent.change(screen.getByLabelText("Publicada em *"), { target: { value: "2027-07-01" } });
    fireEvent.submit(screen.getByRole("form", { name: "Publicar" }));
    await waitFor(() => expect(posted(backend, "/publicar")?.body).toEqual({ publicada_em: "2027-07-01", diario: "" }));
    expect(await screen.findByRole("heading", { name: "Alterar a norma (nova redação)" })).toBeInTheDocument();
    expect(screen.getByText(/Redação 3 em vigor desde 01\/07\/2027/)).toBeInTheDocument();
  });

  it("descarta o rascunho e revoga", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/normas/n1": { data: NORMA },
      "DELETE v1/acervo/normas/n1/rascunho": { data: { ...NORMA, tem_rascunho: false, redacoes: [R2, R1] } },
      "POST v1/acervo/normas/n1/revogar": { data: { ...NORMA, situacao: "REVOGADA", revogada_em: "2027-08-01T00:00:00Z", revogacao_ato: "IN 9/2027", redacoes: [R2, R1] } },
    });
    renderApp(<NormaPage />);
    await userEvent.click(await screen.findByRole("button", { name: "Descartar rascunho" }));
    expect(await screen.findByText("A redação em vigor continua valendo.")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Descartar" }));
    await waitFor(() => expect(backend.calls.some((c) => c.method === "DELETE")).toBe(true));

    fireEvent.change(screen.getByLabelText("Ato de revogação *"), { target: { value: "IN 9/2027" } });
    fireEvent.change(screen.getByLabelText("Revogada em *"), { target: { value: "2027-08-01" } });
    fireEvent.submit(screen.getByRole("form", { name: "Revogar" }));
    await waitFor(() => expect(posted(backend, "/revogar")?.body).toEqual({ revogada_em: "2027-08-01", ato: "IN 9/2027" }));
    expect(await screen.findByText(/Revogada em 01\/08\/2027 — IN 9\/2027/)).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Redação" })).not.toBeInTheDocument();
  });

  it("norma nunca publicada: descartar volta à lista", async () => {
    const nova = { ...NORMA, situacao: "RASCUNHO", versao: 0, publicada_em: null, redacoes: [{ ...RASC, versao: 1, motivo: "Redação original" }] };
    mockBackend({ ...identityRoutes(), "GET v1/acervo/normas/n1": { data: nova }, "DELETE v1/acervo/normas/n1/rascunho": { status: 204 } });
    renderApp(<NormaPage />);
    expect(await screen.findByText("A norma ainda não foi publicada; o público não a vê.")).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /Página pública/ })).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Descartar rascunho" }));
    await userEvent.click(await screen.findByRole("button", { name: "Descartar" }));
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/normas"));
  });
});

describe("Arquivos privados", () => {
  beforeEach(() => resetNavigation({ id: "p1" }));

  it("lista e cadastra", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/privados": page([PRIVADO]),
      "POST v1/acervo/privados": { data: PRIVADO },
    });
    renderApp(<PrivadosPage />);
    expect(await screen.findByRole("link", { name: "Família Silva" })).toHaveAttribute("href", "/acervo/privados/p1");
    expect(screen.getByText("Aguardando parecer")).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("Buscar"), { target: { value: "Silva" } });
    await waitFor(() => expect(backend.calls.some((c) => c.query.get("q") === "Silva")).toBe(true));
    await userEvent.click(screen.getByRole("button", { name: "Cadastrar arquivo" }));
    const form = screen.getByRole("form", { name: "Cadastrar arquivo privado" });
    fireEvent.change(within(form).getByLabelText("Titular *"), { target: { value: "Concessionária" } });
    fireEvent.change(within(form).getByLabelText("Tipo *"), { target: { value: "PJ" } });
    fireEvent.change(within(form).getByLabelText(/Conjuntos e relevância/), { target: { value: "Plantas da rede" } });
    await userEvent.click(within(form).getByRole("checkbox", { name: /serviço público municipal/ }));
    fireEvent.submit(form);
    await waitFor(() => expect(router.push).toHaveBeenCalledWith("/acervo/privados/p1"));
    expect(posted(backend, "/privados")?.body).toMatchObject({ titular: "Concessionária", tipo_titular: "PJ", servico_publico: true, acesso_autorizado: false });
    await userEvent.click(screen.getByRole("button", { name: "Cancelar" }));
    expect(screen.getByRole("button", { name: "Cadastrar arquivo" })).toBeInTheDocument();
  });

  it("sem acervo:central não vê o cadastro", async () => {
    mockBackend({ ...identityRoutes({ permissions: ["acervo:read"], scopes: [] }) });
    renderApp(<PrivadosPage />);
    expect(await screen.findByText(/O cadastro dos arquivos privados é do Arquivo Público/)).toBeInTheDocument();
  });

  it("parecer, alteração do cadastro e decreto", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/privados/p1": { data: PRIVADO },
      "PUT v1/acervo/privados/p1": { data: { ...PRIVADO, contato: "(66) 9999" } },
      "POST v1/acervo/privados/p1/parecer": { data: { ...PRIVADO, situacao: "COM_PARECER", parecer_favoravel: true, parecer_em: "2026-09-01T00:00:00Z", parecer_texto: "Único" } },
      "POST v1/acervo/privados/p1/declarar": { data: { ...DECLARADO, guarda: "TITULAR", alienacoes: [] } },
    });
    renderApp(<PrivadoPage />);
    expect(await screen.findByRole("heading", { name: "Família Silva" })).toBeInTheDocument();
    expect(screen.getByText("Não autorizado")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "Alterar cadastro" }));
    const edit = screen.getByRole("form", { name: "Alterar cadastro" });
    fireEvent.change(within(edit).getByLabelText("Contato do titular"), { target: { value: "(66) 9999" } });
    fireEvent.submit(edit);
    await waitFor(() => expect(posted(backend, "/privados/p1")?.body).toMatchObject({ contato: "(66) 9999", tipo_titular: "PF" }));
    await waitFor(() => expect(screen.queryByRole("form", { name: "Alterar cadastro" })).not.toBeInTheDocument());
    await userEvent.click(screen.getByRole("button", { name: "Alterar cadastro" }));
    await userEvent.click(screen.getByRole("button", { name: "Cancelar" }));

    fireEvent.change(screen.getByLabelText("Comissão Especial (composição) *"), { target: { value: "A, B e C" } });
    fireEvent.change(screen.getByLabelText("Parecer com a avaliação técnica *"), { target: { value: "Único" } });
    fireEvent.change(screen.getByLabelText("Data *"), { target: { value: "2026-09-01" } });
    await userEvent.click(screen.getByRole("checkbox", { name: "Favorável à declaração" }));
    fireEvent.submit(screen.getByRole("form", { name: "Parecer" }));
    await waitFor(() =>
      expect(posted(backend, "/parecer")?.body).toEqual({ parecer_comissao: "A, B e C", parecer_texto: "Único", parecer_favoravel: true, parecer_em: "2026-09-01" }),
    );
    expect(await screen.findByText(/Favorável em 01\/09\/2026/)).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Decreto nº *"), { target: { value: "Decreto 10/2026" } });
    fireEvent.change(screen.getByLabelText("Data *"), { target: { value: "2026-10-01" } });
    fireEvent.submit(screen.getByRole("form", { name: "Decreto" }));
    await waitFor(() => expect(posted(backend, "/declarar")?.body).toEqual({ decreto_numero: "Decreto 10/2026", decreto_em: "2026-10-01", decreto_diario: "" }));
    expect(await screen.findByText("De interesse público e social")).toBeInTheDocument();
    expect(screen.getByRole("form", { name: "Notificação de alienação" })).toBeInTheDocument();
  });

  it("guarda e direito de preferência", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/privados/p1": { data: { ...DECLARADO, servico_publico: true, tipo_titular: "PJ", documento: "", contato: "x" } },
      "POST v1/acervo/privados/p1/guarda": { data: { ...DECLARADO, guarda: "DOADO", alienacoes: [] } },
      "POST v1/acervo/privados/p1/alienacoes/a1/decidir": { data: DECLARADO },
    });
    renderApp(<PrivadoPage />);
    expect(await screen.findByText("Serviço público (§ 3º)")).toBeInTheDocument();
    expect(screen.getByText(/Autorizado pelo titular — Com agendamento/)).toBeInTheDocument();
    expect(screen.getByText(/Decreto 10\/2026, de 01\/10\/2026 \(DIORONDON-e\)/)).toBeInTheDocument();
    expect(screen.getByText("Prazo vencido — a preferência caducou")).toBeInTheDocument();
    expect(screen.getByText("O município não tem interesse")).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Notificação de alienação" })).not.toBeInTheDocument();

    const decisao = screen.getByRole("form", { name: "Manifestação do município" });
    fireEvent.change(within(decisao).getByLabelText("Fundamento *"), { target: { value: "Sem dotação" } });
    await userEvent.click(within(decisao).getByRole("button", { name: "Não tem interesse" }));
    await waitFor(() => expect(posted(backend, "/decidir")?.body).toEqual({ decisao: "SEM_INTERESSE", fundamento: "Sem dotação" }));
    fireEvent.change(within(decisao).getByLabelText("Fundamento *"), { target: { value: "Acervo da fundação" } });
    await userEvent.click(within(decisao).getByRole("button", { name: "Tem interesse" }));
    await waitFor(() => expect(backend.calls.filter((c) => c.path.endsWith("/decidir"))[1]?.body).toEqual({ decisao: "INTERESSE", fundamento: "Acervo da fundação" }));

    const guarda = screen.getByRole("form", { name: "Guarda" });
    expect(within(guarda).getByRole("option", { name: "Com o titular" })).toBeInTheDocument();
    fireEvent.change(within(guarda).getByLabelText("Nova guarda *"), { target: { value: "DOADO" } });
    fireEvent.change(within(guarda).getByLabelText("Termo *"), { target: { value: "Termo de doação" } });
    fireEvent.change(within(guarda).getByLabelText("Data *"), { target: { value: "2026-10-15" } });
    fireEvent.submit(guarda);
    await waitFor(() => expect(posted(backend, "/guarda")?.body).toEqual({ guarda: "DOADO", termo: "Termo de doação", em: "2026-10-15" }));
    expect(await screen.findByText("Doado ao Arquivo Público")).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Guarda" })).not.toBeInTheDocument();
  });

  it("notifica a alienação; parecer contrário encerra", async () => {
    resetNavigation({ id: "p1" });
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/privados/p1": { data: { ...DECLARADO, guarda: "TITULAR", guarda_em: null, alienacoes: [] } },
      "POST v1/acervo/privados/p1/alienacoes": { data: DECLARADO },
    });
    renderApp(<PrivadoPage />);
    const f = await screen.findByRole("form", { name: "Notificação de alienação" });
    fireEvent.change(within(f).getByLabelText("Notificado em *"), { target: { value: "2026-10-10" } });
    fireEvent.change(within(f).getByLabelText("Pretendente"), { target: { value: "Colecionador" } });
    fireEvent.submit(f);
    await waitFor(() => expect(posted(backend, "/alienacoes")?.body).toEqual({ notificado_em: "2026-10-10", adquirente: "Colecionador", condicoes: "" }));
  });

  it("parecer contrário", async () => {
    mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/privados/p1": { data: { ...PRIVADO, situacao: "NAO_DECLARADO", parecer_favoravel: false, parecer_em: "2026-09-01T00:00:00Z", parecer_texto: "Sem relevância" } },
    });
    renderApp(<PrivadoPage />);
    expect(await screen.findByText(/Contrário em 01\/09\/2026 — Sem relevância/)).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Parecer" })).not.toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Decreto" })).not.toBeInTheDocument();
  });
});
