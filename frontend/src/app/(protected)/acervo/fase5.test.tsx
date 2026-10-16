import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fmtBytes, situacaoDigital } from "@/components/acervo/labels";
import { identityRoutes, mockBackend, page, renderApp } from "@/test/backend";
import { resetNavigation } from "@/test/navigation";

vi.mock("next/navigation", async () => (await import("@/test/navigation")).navigationModule);
vi.mock("pdfjs-dist/web/pdf_viewer.css", () => ({}));

// pdf.js simulado: o jsdom não tem canvas.
const pdf = vi.hoisted(() => ({
  getDocument: vi.fn(),
  setDocument: vi.fn(),
  handlers: {} as Record<string, () => void>,
  destroy: vi.fn(async () => {}),
}));
vi.mock("pdfjs-dist", () => ({ GlobalWorkerOptions: {}, getDocument: pdf.getDocument }));
vi.mock("pdfjs-dist/web/pdf_viewer.mjs", () => ({
  EventBus: class {
    on(ev: string, fn: () => void) {
      pdf.handlers[ev] = fn;
    }
  },
  PDFViewer: class {
    setDocument(doc: unknown) {
      pdf.setDocument(doc);
    }
  },
}));
const upload = vi.hoisted(() => ({ putToTicket: vi.fn(async () => {}) }));
vi.mock("@/lib/sismarq/upload", async (orig) => ({ ...(await orig<typeof import("@/lib/sismarq/upload")>()), putToTicket: upload.putToTicket }));

import { DigitaisDaUnidade } from "@/components/acervo/DigitaisDaUnidade";
import { VisualizadorPDF } from "@/components/acervo/VisualizadorPDF";
import DigitalPage from "./digitais/[id]/page";
import DigitaisPage from "./digitais/page";

const SHA = "a".repeat(64);
const DIGITAL = {
  id: "d1",
  numero: "RD-2026-000001",
  unidade_arquivamento_id: "u1",
  descricao: "Volume 1",
  tipo_documento: "TEXTO",
  equipamento: "Fujitsu fi-7160",
  dpi_declarado: 300,
  responsavel_nome: "Ana",
  responsavel_cpf: "***.456.789-**",
  codigo_classificacao: "2.0.01.00.02",
  digitalizado_em: "2026-10-16T14:30:00Z",
  situacao: "PRONTO",
  motivo: "",
  tentativas: 1,
  original_sha256: "b".repeat(64),
  original_bytes: 2048,
  pdfa_sha256: SHA,
  pdfa_bytes: 3 * 1024 * 1024,
  paginas: 12,
  dpi_medido: 300,
  conformidade: "PDF/A-2b",
  retencao_ate: "2126-10-16T00:00:00Z",
  processado_em: "2026-10-16T14:31:00Z",
  verificado_em: "2026-10-17T00:00:00Z",
  integro: true,
  criado_em: "2026-10-16T14:30:00Z",
  publicado_em: null,
  pades: "AUSENTE",
};
const DETALHE = {
  ...DIGITAL,
  unidade_numero: "UA-2026-000001",
  verificacoes: [
    { em: "2026-10-17T00:00:00Z", integro: true, sha256_lido: SHA, erro: "" },
    { em: "2026-10-18T00:00:00Z", integro: false, sha256_lido: "c".repeat(64), erro: "" },
    { em: "2026-10-19T00:00:00Z", integro: false, sha256_lido: "", erro: "objeto não encontrado" },
  ],
};

beforeEach(() => {
  resetNavigation({ id: "d1" });
  pdf.getDocument.mockReset();
  pdf.setDocument.mockReset();
  pdf.handlers = {};
  upload.putToTicket.mockClear();
});

describe("rótulos", () => {
  it("situação e tamanho", () => {
    expect(situacaoDigital("X").label).toBe("X");
    expect(fmtBytes(500)).toBe("500 B");
    expect(fmtBytes(1536)).toBe("1,5 KB");
    expect(fmtBytes(3 * 1024 * 1024 * 1024)).toBe("3 GB");
  });
});

describe("VisualizadorPDF", () => {
  it("abre por trechos (Range), sem baixar o arquivo, e mede a primeira página", async () => {
    pdf.getDocument.mockReturnValue({ promise: Promise.resolve({ numPages: 600 }), destroy: pdf.destroy });
    const { unmount } = render(<VisualizadorPDF url="/api/backend/v1/acervo/digitais/d1/arquivo" titulo="RD-1" />);
    expect(screen.getByText("Abrindo o documento…")).toBeInTheDocument();
    expect(await screen.findByText(/600 página\(s\)/)).toBeInTheDocument();
    expect(pdf.getDocument).toHaveBeenCalledWith(
      expect.objectContaining({ url: "/api/backend/v1/acervo/digitais/d1/arquivo", rangeChunkSize: 65536, disableAutoFetch: true, disableStream: true }),
    );
    expect(pdf.setDocument).toHaveBeenCalledWith({ numPages: 600 });
    pdf.handlers.pagerendered!();
    pdf.handlers.pagerendered!();
    expect(await screen.findByText(/primeira página em \d+ ms/)).toBeInTheDocument();
    unmount();
    expect(pdf.destroy).toHaveBeenCalled();
  });

  it("mostra o erro de abertura", async () => {
    pdf.getDocument.mockReturnValue({ promise: Promise.reject(new Error("PDF inválido")), destroy: pdf.destroy });
    render(<VisualizadorPDF url="/x" titulo="RD-1" />);
    expect(await screen.findByText("Não foi possível abrir: PDF inválido")).toBeInTheDocument();
  });

  it("erro sem mensagem e desmontagem no meio da abertura", async () => {
    pdf.getDocument.mockReturnValue({ promise: Promise.reject("x"), destroy: pdf.destroy });
    render(<VisualizadorPDF url="/x" titulo="RD-1" />);
    expect(await screen.findByText(/não foi possível abrir o documento/)).toBeInTheDocument();
    let abrir: (v: unknown) => void = () => {};
    pdf.getDocument.mockReturnValue({ promise: new Promise((r) => (abrir = r)), destroy: pdf.destroy });
    const { unmount } = render(<VisualizadorPDF url="/y" titulo="RD-2" />);
    await waitFor(() => expect(pdf.getDocument).toHaveBeenCalledTimes(2));
    unmount();
    abrir({ numPages: 1 });
    await waitFor(() => expect(pdf.setDocument).not.toHaveBeenCalled());
  });
});

describe("DigitaisDaUnidade", () => {
  it("lista e envia um PDF: ticket, upload direto e registro com os metadados", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/digitais": page([
        DIGITAL,
        { ...DIGITAL, id: "d2", numero: "RD-2026-000002", situacao: "REPROVADO", motivo: "200 DPI abaixo do mínimo", descricao: "", integro: null },
        { ...DIGITAL, id: "d3", numero: "RD-2026-000003", situacao: "ENVIADO", motivo: "", integro: false },
      ]),
      "POST v1/acervo/unidades/u1/digitais/upload": { data: { object_key: "acervo/digitalizacao/u1/x.pdf", upload_url: "https://minio/x", method: "PUT", headers: {}, expires_at: "" } },
      "POST v1/acervo/unidades/u1/digitais": { data: DETALHE },
    });
    renderApp(<DigitaisDaUnidade unidadeId="u1" podeGerir />);
    expect(await screen.findByRole("link", { name: "RD-2026-000001" })).toHaveAttribute("href", "/acervo/digitais/d1");
    expect(screen.getByText(/12 página\(s\) · 300 DPI · 3 MB · PDF\/A-2b/)).toBeInTheDocument();
    expect(screen.getByText("200 DPI abaixo do mínimo")).toBeInTheDocument();
    expect(screen.getByText("2 KB enviados")).toBeInTheDocument();
    expect(screen.getByText("Integridade falhou")).toBeInTheDocument();

    const form = screen.getByRole("form", { name: "Enviar digitalização" });
    const enviar = within(form).getByRole("button", { name: "Enviar para processamento" });
    expect(enviar).toBeDisabled();
    const arquivo = new File(["%PDF"], "lote.pdf", { type: "application/pdf" });
    await userEvent.upload(within(form).getByLabelText("PDF do scanner *"), arquivo);
    fireEvent.change(within(form).getByLabelText("Scanner (modelo e nº de série) *"), { target: { value: "Fujitsu fi-7160" } });
    fireEvent.change(within(form).getByLabelText("Tipo *"), { target: { value: "PLANTA" } });
    fireEvent.change(within(form).getByLabelText("Resolução (DPI) *"), { target: { value: "600" } });
    fireEvent.change(within(form).getByLabelText("CPF de quem digitalizou"), { target: { value: "123.456.789-01" } });
    fireEvent.submit(form);
    await waitFor(() => expect(backend.calls.some((c) => c.method === "POST" && c.path === "v1/acervo/unidades/u1/digitais")).toBe(true));
    expect(upload.putToTicket).toHaveBeenCalledWith(expect.objectContaining({ object_key: "acervo/digitalizacao/u1/x.pdf" }), arquivo);
    expect(backend.calls.find((c) => c.path === "v1/acervo/unidades/u1/digitais/upload")?.body).toEqual({ filename: "lote.pdf" });
    expect(backend.calls.find((c) => c.path === "v1/acervo/unidades/u1/digitais")?.body).toEqual({
      object_key: "acervo/digitalizacao/u1/x.pdf",
      descricao: "",
      tipo_documento: "PLANTA",
      equipamento: "Fujitsu fi-7160",
      dpi_declarado: 600,
      responsavel_cpf: "123.456.789-01",
    });
    await waitFor(() => expect(enviar).toBeDisabled());
  });

  it("sem gestão só consulta; sem digitalizações", async () => {
    mockBackend({ ...identityRoutes(), "GET v1/acervo/digitais": page([]) });
    renderApp(<DigitaisDaUnidade unidadeId="u1" podeGerir={false} />);
    expect(await screen.findByText("Nenhuma digitalização.")).toBeInTheDocument();
    expect(screen.queryByRole("form", { name: "Enviar digitalização" })).not.toBeInTheDocument();
  });

  it("envio sem arquivo não faz nada", async () => {
    const backend = mockBackend({ ...identityRoutes(), "GET v1/acervo/digitais": page([]) });
    renderApp(<DigitaisDaUnidade unidadeId="u1" podeGerir />);
    fireEvent.submit(await screen.findByRole("form", { name: "Enviar digitalização" }));
    expect(backend.calls.some((c) => c.method === "POST")).toBe(false);
  });
});

describe("Representante digital", () => {
  it("metadados, verificações, visualizador e verificação pelo Arquivo Público", async () => {
    pdf.getDocument.mockReturnValue({ promise: Promise.resolve({ numPages: 12 }), destroy: pdf.destroy });
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/digitais/d1": { data: DETALHE },
      "POST v1/acervo/digitais/d1/verificar": { data: { ...DETALHE, integro: false } },
    });
    renderApp(<DigitalPage />);
    expect(await screen.findByRole("heading", { name: "RD-2026-000001" })).toBeInTheDocument();
    expect(screen.getByText("Íntegro na última verificação")).toBeInTheDocument();
    expect(screen.getByText(/300 DPI medidos \(declarado 300 DPI, texto\)/)).toBeInTheDocument();
    expect(screen.getByText("Ana · ***.456.789-**")).toBeInTheDocument();
    expect(screen.getByText(SHA)).toBeInTheDocument();
    expect(screen.getByText(/hash diferente \(cccccccccccccccc…\)/)).toBeInTheDocument();
    expect(screen.getByText(/falhou: objeto não encontrado/)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Baixar PDF\/A/ })).toHaveAttribute("href", "/api/backend/v1/acervo/digitais/d1/arquivo");
    expect(screen.getByRole("link", { name: /Unidade UA-2026-000001/ })).toHaveAttribute("href", "/acervo/u1");
    expect(await screen.findByText(/12 página\(s\) · carregado sob demanda/)).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Verificar integridade" }));
    await waitFor(() => expect(backend.calls.some((c) => c.path.endsWith("/verificar"))).toBe(true));
    expect(await screen.findByText("Integridade falhou — o hash não confere")).toBeInTheDocument();
  });

  it("falhou: o Arquivo Público reprocessa; sem verificações nem visualizador", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/digitais/d1": {
        data: { ...DETALHE, situacao: "FALHOU", motivo: "ocrmypdf: exit 2", descricao: "", integro: null, verificacoes: [], dpi_medido: 0, conformidade: "" },
      },
      "POST v1/acervo/digitais/d1/reprocessar": { data: { ...DETALHE, situacao: "ENVIADO", verificacoes: [] } },
    });
    renderApp(<DigitalPage />);
    expect(await screen.findByText("ocrmypdf: exit 2")).toBeInTheDocument();
    expect(screen.getByText("Ainda não verificado.")).toBeInTheDocument();
    expect(screen.getByText(/Digitalização da unidade UA-2026-000001/)).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /Baixar/ })).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Reprocessar" }));
    await waitFor(() => expect(backend.calls.some((c) => c.path.endsWith("/reprocessar"))).toBe(true));
    expect(await screen.findByText("Na fila do processamento")).toBeInTheDocument();
  });

  it("sem acervo:central não há ações", async () => {
    mockBackend({
      ...identityRoutes({ permissions: ["acervo:read"], scopes: [] }),
      "GET v1/acervo/digitais/d1": { data: { ...DETALHE, situacao: "FALHOU", verificacoes: [] } },
    });
    renderApp(<DigitalPage />);
    expect(await screen.findByRole("heading", { name: "RD-2026-000001" })).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Reprocessar" })).not.toBeInTheDocument();
  });
});

describe("Digitalização — lista", () => {
  it("lista e filtra as falhas", async () => {
    const backend = mockBackend({
      ...identityRoutes(),
      "GET v1/acervo/digitais": page([{ ...DIGITAL, integro: false, verificado_em: null }]),
    });
    renderApp(<DigitaisPage />);
    expect(await screen.findByRole("link", { name: "RD-2026-000001" })).toBeInTheDocument();
    expect(screen.getByText("Integridade falhou")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("checkbox", { name: /Só os que falharam/ }));
    await waitFor(() => expect(backend.calls.some((c) => c.query.get("falhas") === "true")).toBe(true));
  });

  it("vazia, com e sem o filtro", async () => {
    mockBackend({ ...identityRoutes(), "GET v1/acervo/digitais": page([]) });
    renderApp(<DigitaisPage />);
    expect(await screen.findByText("Nenhuma digitalização")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("checkbox", { name: /Só os que falharam/ }));
    expect(await screen.findByText("Nenhuma falha")).toBeInTheDocument();
  });
});
