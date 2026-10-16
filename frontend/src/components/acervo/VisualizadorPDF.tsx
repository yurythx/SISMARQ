"use client";

import "pdfjs-dist/web/pdf_viewer.css";

import { useEffect, useRef, useState } from "react";

/** Trecho pedido por vez (64 KB): o pdf.js busca só o que a página visível
 * precisa (HTTP 206), sem baixar o arquivo inteiro. */
const TRECHO = 64 * 1024;

type Estado = { paginas?: number; ms?: number; erro?: string };

/** Visualizador em fluxo do PDF/A (ADR 037): pdf.js com Range, renderizando
 * só as páginas na tela. Abre processos de centenas de páginas sem baixar o
 * arquivo; a meta do artigo técnico é a primeira página em menos de 300 ms. */
export function VisualizadorPDF({ url, titulo }: { url: string; titulo: string }) {
  const container = useRef<HTMLDivElement>(null);
  const [estado, setEstado] = useState<Estado>({});

  useEffect(() => {
    let cancelado = false;
    let destruir: (() => Promise<void>) | undefined;
    const inicio = performance.now();
    (async () => {
      const pdfjs = await import("pdfjs-dist");
      const { EventBus, PDFViewer } = await import("pdfjs-dist/web/pdf_viewer.mjs");
      pdfjs.GlobalWorkerOptions.workerSrc = new URL("pdfjs-dist/build/pdf.worker.min.mjs", import.meta.url).toString();
      const tarefa = pdfjs.getDocument({ url, rangeChunkSize: TRECHO, disableAutoFetch: true, disableStream: true });
      destruir = () => tarefa.destroy();
      const doc = await tarefa.promise;
      if (cancelado || !container.current) return;
      const eventBus = new EventBus();
      let primeira = true;
      eventBus.on("pagerendered", () => {
        if (primeira && !cancelado) {
          primeira = false;
          setEstado((e) => ({ ...e, ms: Math.round(performance.now() - inicio) }));
        }
      });
      const viewer = new PDFViewer({ container: container.current, eventBus });
      viewer.setDocument(doc);
      setEstado({ paginas: doc.numPages });
    })().catch((e: unknown) => {
      if (!cancelado) setEstado({ erro: e instanceof Error ? e.message : "não foi possível abrir o documento" });
    });
    return () => {
      cancelado = true;
      void destruir?.();
    };
  }, [url]);

  return (
    <section aria-label={`Documento ${titulo}`} className="flex flex-col gap-2">
      <p className="text-xs text-muted" aria-live="polite">
        {estado.erro
          ? `Não foi possível abrir: ${estado.erro}`
          : estado.paginas
            ? `${estado.paginas} página(s)${estado.ms !== undefined ? ` · primeira página em ${estado.ms} ms` : ""} · carregado sob demanda`
            : "Abrindo o documento…"}
      </p>
      <div className="relative h-[75vh] overflow-hidden rounded-lg border border-surface-border bg-surface-2">
        <div ref={container} data-testid="pdf-container" className="absolute inset-0 overflow-auto">
          <div className="pdfViewer" />
        </div>
      </div>
    </section>
  );
}
