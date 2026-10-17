import { describe, expect, it } from "vitest";

import { buildNav } from "./Sidebar";
import { isNavHrefActive } from "@/lib/nav/sectionTabs";
import type { ModuleStatus } from "@/lib/sismarq/types";
import { hasPermission } from "@/lib/sismarq/permissions";

function mod(key: string, extra: Partial<ModuleStatus> = {}): ModuleStatus {
  return {
    key, name: key, description: "", core: false, default_enabled: true, depends_on: [], permissions: [],
    public: false, icon: "box", route: `/${key}`, enabled: true, configured: true, dependents: [], blocked_by: [], ...extra,
  };
}

describe("buildNav", () => {
  const modules = [
    mod("iam", { core: true }), mod("audit", { core: true }), mod("blog"), mod("wiki", { enabled: false }),
    mod("egress"), mod("search"), mod("tramite"), mod("acervo"),
  ];

  it("lista só plugins ativos, fora do núcleo/ocultos", () => {
    const nav = buildNav(modules, () => false);
    expect(nav.modules.map((m) => m.href)).toEqual(["/blog", "/tramite", "/acervo"]);
    expect(nav.admin).toEqual([]);
  });

  it("submenus do módulo (Trâmite, Acervo) e esconde o que exige permissão", () => {
    const nav = buildNav(modules, () => false);
    const tramite = nav.modules.find((m) => m.href === "/tramite");
    expect(tramite?.children?.map((c) => c.href)).toEqual(["/tramite/guarda"]);
    const acervo = nav.modules.find((m) => m.href === "/acervo");
    expect(acervo?.children?.some((c) => c.href === "/acervo/eliminacao")).toBe(true);
    expect(acervo?.children?.some((c) => c.href === "/acervo/privados")).toBe(false);
    expect(acervo?.children?.some((c) => c.href === "/acervo/importar")).toBe(false);
  });

  it("Acervo: arquivos privados e importar só com a permissão", () => {
    const nav = buildNav(modules, (p) => hasPermission(["acervo:manage", "acervo:central"], p));
    const hrefs = nav.modules.find((m) => m.href === "/acervo")?.children?.map((c) => c.href) ?? [];
    expect(hrefs).toContain("/acervo/privados");
    expect(hrefs).toContain("/acervo/importar");
  });

  it("administração depende das permissões efetivas (com curinga)", () => {
    const nav = buildNav(modules, (p) => hasPermission(["audit:*"], p));
    expect(nav.admin.map((a) => a.href)).toEqual(["/auditoria"]);
    const all = buildNav(modules, (p) => hasPermission(["*"], p));
    expect(all.admin.map((a) => a.href)).toEqual(["/auditoria", "/monitoramento", "/configuracao"]);
    expect(all.admin.find((a) => a.href === "/configuracao")?.children?.some((c) => c.href === "/configuracao/modulos")).toBe(true);
  });

  it("módulos de gestão (catálogo, contato) exigem a permissão para entrar no menu", () => {
    const mods = [mod("catalog", { route: "/gestao/servicos" }), mod("contact", { route: "/gestao/contato" })];
    expect(buildNav(mods, () => false).modules).toEqual([]);
    const nav = buildNav(mods, (p) => hasPermission(["contact:read"], p));
    expect(nav.modules.map((m) => m.href)).toEqual(["/gestao/contato"]);
  });
});

describe("isNavHrefActive", () => {
  const acervo = [
    { href: "/acervo" },
    { href: "/acervo/eliminacao" },
    { href: "/acervo/guias" },
  ];
  it("a subpágina mais específica ganha do índice do módulo", () => {
    expect(isNavHrefActive("/acervo/eliminacao/x", "/acervo", acervo)).toBe(false);
    expect(isNavHrefActive("/acervo/eliminacao/x", "/acervo/eliminacao", acervo)).toBe(true);
    expect(isNavHrefActive("/acervo/abc", "/acervo", acervo)).toBe(true);
  });
});
