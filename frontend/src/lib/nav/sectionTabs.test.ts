import { describe, expect, it } from "vitest";

import { hasPermission } from "@/lib/sismarq/permissions";

import { CONFIG_TABS, MODULE_CHILDREN, isNavHrefActive, visibleTabs } from "./sectionTabs";

describe("abas de Configurações", () => {
  it("cada aba só aparece com a permissão correspondente", () => {
    const tabs = visibleTabs(CONFIG_TABS, (p) => hasPermission(["users:read", "iam:*"], p), () => true);
    expect(tabs.map((t) => t.href)).toEqual([
      "/configuracao/organizacao",
      "/configuracao/perfis",
      "/configuracao/mapeamento-ad",
      "/configuracao/usuarios",
    ]);
  });

  it("Egress exige a permissão E o plug-in ativo", () => {
    const can = (p: string) => hasPermission(["*"], p);
    expect(visibleTabs(CONFIG_TABS, can, () => true).some((t) => t.href.endsWith("egress"))).toBe(true);
    expect(visibleTabs(CONFIG_TABS, can, (m) => m !== "egress").some((t) => t.href.endsWith("egress"))).toBe(false);
  });

  it("administração delegada (ADR 013): perfis e mapeamentos só com concessão global", () => {
    const can = (p: string) => hasPermission(["users:read", "iam:*"], p);
    const tabs = visibleTabs(CONFIG_TABS, can, () => true, () => false);
    expect(tabs.map((t) => t.href)).toEqual(["/configuracao/organizacao", "/configuracao/usuarios"]);
  });

  it("sem permissões: nenhuma aba", () => {
    expect(visibleTabs(CONFIG_TABS, () => false, () => true)).toEqual([]);
  });
});

describe("submenus dos módulos", () => {
  it("Acervo: privados e importar só com a permissão (global no Arquivo)", () => {
    const todos = visibleTabs(MODULE_CHILDREN.acervo!, (p) => hasPermission(["*"], p), () => true);
    expect(todos.some((t) => t.href === "/acervo/privados")).toBe(true);
    const leitura = visibleTabs(MODULE_CHILDREN.acervo!, (p) => hasPermission(["acervo:read"], p), () => true);
    expect(leitura.some((t) => t.href === "/acervo/eliminacao")).toBe(true);
    expect(leitura.some((t) => t.href === "/acervo/privados")).toBe(false);
  });

  it("Atlas: cobertura só para quem gere", () => {
    const leitura = visibleTabs(MODULE_CHILDREN.atlas!, (p) => hasPermission(["atlas:read"], p), () => true);
    expect(leitura.map((t) => t.href)).toEqual(["/atlas/ttdd", "/atlas/organograma", "/atlas/modelos"]);
  });

  it("isNavHrefActive: prefixo mais longo", () => {
    const irmaos = [{ href: "/acervo" }, { href: "/acervo/digitais" }];
    expect(isNavHrefActive("/acervo/digitais/d1", "/acervo", irmaos)).toBe(false);
    expect(isNavHrefActive("/acervo/digitais/d1", "/acervo/digitais", irmaos)).toBe(true);
  });
});

