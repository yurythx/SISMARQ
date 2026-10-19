import type { ReactNode } from "react";

import { ModuleTabs } from "@/components/layout/ModuleTabs";

/** Todas as páginas do Acervo compartilham as abas do módulo. */
export default function AcervoLayout({ children }: { children: ReactNode }) {
  return (
    <>
      <ModuleTabs module="acervo" index={{ href: "/acervo", label: "Unidades" }} ariaLabel="Seções do Acervo" />
      {children}
    </>
  );
}
