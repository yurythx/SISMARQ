import type { ReactNode } from "react";

import { AssistenteProvider } from "@/components/atlas/AssistenteGaveta";
import { ModuleTabs } from "@/components/layout/ModuleTabs";

/** Todas as páginas do Atlas compartilham as abas do módulo e o
 * assistente em gaveta lateral. */
export default function AtlasLayout({ children }: { children: ReactNode }) {
  return (
    <AssistenteProvider>
      <ModuleTabs module="atlas" index={{ href: "/atlas", label: "Procedimentos" }} ariaLabel="Seções do Atlas" />
      {children}
    </AssistenteProvider>
  );
}
