import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

let pathname = "/acervo";
vi.mock("next/navigation", () => ({ usePathname: () => pathname }));
let perms: string[] = [];
vi.mock("@/lib/sismarq/SismarqProvider", () => ({
  useSismarq: () => ({ can: (p: string) => perms.includes(p), canGlobal: (p: string) => perms.includes(p), enabled: () => true }),
}));

import { ModuleTabs } from "./ModuleTabs";

const nomes = () => screen.getAllByRole("link").map((a) => a.textContent);

describe("ModuleTabs", () => {
  it("índice + as mesmas subpáginas do submenu, filtradas por permissão", () => {
    perms = [];
    const { unmount } = render(<ModuleTabs module="acervo" index={{ href: "/acervo", label: "Unidades" }} ariaLabel="Seções do Acervo" />);
    expect(nomes()[0]).toBe("Unidades");
    expect(nomes()).toContain("Eliminação");
    expect(nomes()).not.toContain("Importar");
    expect(nomes()).not.toContain("Exposições");
    unmount();

    perms = ["acervo:manage", "acervo:central"];
    render(<ModuleTabs module="acervo" index={{ href: "/acervo", label: "Unidades" }} ariaLabel="Seções do Acervo" />);
    expect(nomes()).toEqual(expect.arrayContaining(["Importar", "Arquivos privados", "Exposições"]));
  });

  it("marca a aba da página atual (o índice só na rota exata)", () => {
    perms = [];
    pathname = "/atlas/ttdd/1.0";
    render(<ModuleTabs module="atlas" index={{ href: "/atlas", label: "Procedimentos" }} ariaLabel="Seções do Atlas" />);
    expect(screen.getByRole("navigation", { name: "Seções do Atlas" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Tabela de Temporalidade" })).toHaveAttribute("aria-current", "page");
    expect(screen.getByRole("link", { name: "Procedimentos" })).not.toHaveAttribute("aria-current");
  });
});
