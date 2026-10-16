import { describe, expect, it } from "vitest";

import { isProtectedPath, PROTECTED_PREFIXES } from "./protected";

describe("rotas protegidas", () => {
  it("a caixa de notificações e o acervo exigem sessão", () => {
    expect(PROTECTED_PREFIXES).toContain("/notificacoes");
    expect(PROTECTED_PREFIXES).toContain("/acervo");
    expect(isProtectedPath("/notificacoes")).toBe(true);
    expect(isProtectedPath("/acervo/busca")).toBe(true);
  });

  it("o site público e o arquivo não pedem login", () => {
    expect(isProtectedPath("/")).toBe(false);
    expect(isProtectedPath("/arquivo")).toBe(false);
    expect(isProtectedPath("/arquivo/busca")).toBe(false);
    expect(isProtectedPath("/login")).toBe(false);
  });
});
