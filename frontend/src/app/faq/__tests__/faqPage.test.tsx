import type { ReactNode } from "react";
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

vi.mock("server-only", () => ({}));
vi.mock("@/components/layout/PublicShell", () => ({
  PublicShell: ({ children }: { children: ReactNode }) => <main>{children}</main>,
}));

import FaqPage from "../page";

describe("FaqPage", () => {
  it("renderiza o cabeçalho, categorias e perguntas frequentes", () => {
    render(<FaqPage />);

    expect(screen.getByRole("heading", { name: "Perguntas Frequentes (FAQ)", level: 1 })).toBeInTheDocument();
    expect(screen.getByText("O Arquivo Público e o SISMARQ")).toBeInTheDocument();
    expect(screen.getByText("Avaliação e eliminação")).toBeInTheDocument();
    expect(screen.getByText("O que é o SISMARQ?")).toBeInTheDocument();
    // Toda resposta cita a base legal (docs/SISMARQ).
    expect(screen.getAllByText(/^Base: /).length).toBeGreaterThan(10);
    expect(screen.getByRole("link", { name: /consultar a temporalidade/i })).toHaveAttribute("href", "/temporalidade");
    expect(screen.getByText("Como entra um documento digital no SISMARQ?")).toBeInTheDocument();
    expect(screen.getByText(/PDF nato-digital/)).toBeInTheDocument();
    expect(screen.getByText("Depois de digitalizar, o papel permanente pode ir embora?")).toBeInTheDocument();
    expect(screen.getByText(/veda a eliminação de documento de guarda permanente/)).toBeInTheDocument();
    expect(screen.getByText("Ainda tem dúvidas?")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /fale conosco/i })).toBeInTheDocument();
  });
});
