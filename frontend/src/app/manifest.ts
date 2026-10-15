import type { MetadataRoute } from "next";

import { getServerBranding } from "@/lib/branding/server";

export default async function manifest(): Promise<MetadataRoute.Manifest> {
  const brand = await getServerBranding().catch(() => null);
  const name = brand?.appName || "SISMARQ";
  const shortName = "SISMARQ";
  const description =
    brand?.appDescription ||
    "Sistema Municipal de Arquivos — Arquivo Público Municipal de Rondonópolis (Lei nº 10.301/2019).";

  return {
    name,
    short_name: shortName,
    description,
    start_url: "/",
    display: "standalone",
    background_color: "#0f172a",
    theme_color: "#1b365d",
    orientation: "portrait-primary",
    icons: [
      {
        src: "/icon.svg",
        sizes: "any",
        type: "image/svg+xml",
        purpose: "any",
      },
      {
        src: "/apple-icon.png",
        sizes: "180x180",
        type: "image/png",
        purpose: "maskable",
      },
    ],
  };
}
