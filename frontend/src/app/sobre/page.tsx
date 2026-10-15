import { Archive, ArrowRight, Landmark, Scale, ShieldCheck } from "lucide-react";
import type { Metadata } from "next";
import Link from "next/link";
import { connection } from "next/server";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card";
import { PublicShell } from "@/components/layout/PublicShell";
import { Table, TableBody, TableCell, TableHead, TableHeaderCell, TableRow } from "@/components/ui/Table";
import { getServerBranding } from "@/lib/branding/server";
import {
  BASE_LEGAL,
  COMPETENCIAS,
  FINALIDADES,
  IDADES,
  LEI,
  MISSAO,
  SEI_X_SISMARQ,
  TTDD_PUBLICACOES,
} from "@/lib/sismarq/conteudo";

const description =
  "O Arquivo Público Municipal de Rondonópolis e o Sistema Municipal de Arquivos (SISMARQ), criados pela Lei Municipal nº 10.301/2019: missão, competências, estrutura e base legal.";

export const metadata: Metadata = {
  title: "Sobre o Arquivo Público e o SISMARQ",
  description,
  openGraph: { title: "Sobre o Arquivo Público e o SISMARQ", description, type: "website" },
};

// Como o sistema informatizado atende às exigências legais e às práticas do
// Governo Federal (art. 13 da lei: e-ARQ Brasil).
const SISTEMA = [
  {
    titulo: "Gestão arquivística (e-ARQ Brasil)",
    texto:
      "O art. 13 da lei prevê um sistema informatizado de gestão arquivística que atenda ao e-ARQ Brasil do CONARQ. A Tabela de Temporalidade em vigor está carregada no sistema, com histórico de versões.",
  },
  {
    titulo: "Prova e integridade",
    texto:
      "Todo ato fica registrado numa trilha de auditoria imutável, encadeada por hash. A digitalização segue o Decreto Federal nº 10.278/2020: PDF/A, metadados técnicos e hash SHA-256.",
  },
  {
    titulo: "Acesso e privacidade",
    texto:
      "Acesso pleno como regra e sigilo como exceção (art. 2º), com proteção dos dados pessoais conforme a LGPD.",
  },
  {
    titulo: "Acessibilidade",
    texto:
      "Portal no padrão de acessibilidade do Governo Federal (e-MAG), com atalhos de teclado, alto contraste e tradução para Libras (VLibras).",
  },
];

export default async function AboutPage() {
  await connection();
  const brand = await getServerBranding();

  return (
    <PublicShell>
      <div className="mx-auto flex w-full max-w-4xl flex-1 flex-col gap-12 px-6 py-12">
        <section className="flex flex-col gap-4">
          <div className="flex items-center gap-2 text-xs font-bold font-mono text-accent uppercase tracking-wider">
            <Landmark size={16} className="text-success" aria-hidden="true" />
            {LEI.numero}, de {LEI.data}
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-foreground sm:text-4xl">
            Sobre o Arquivo Público Municipal e o SISMARQ
          </h1>
          <p className="max-w-3xl text-base leading-relaxed text-muted">{LEI.ementa}</p>
          <p className="max-w-3xl text-base leading-relaxed text-muted">
            O <strong>{brand.appName || "SISMARQ"}</strong> é o sistema do Arquivo Público Municipal de Rondonópolis
            para gerir os documentos da administração municipal em todo o ciclo de vida, da produção à eliminação ou à
            guarda permanente, conforme a lei municipal.
          </p>
        </section>

        <section className="flex flex-col gap-3 rounded-xl border border-surface-border bg-surface p-6 shadow-xs">
          <h2 className="flex items-center gap-2 text-xl font-semibold text-foreground">
            <Archive size={20} className="text-primary" aria-hidden="true" /> Missão
          </h2>
          <p className="text-sm leading-relaxed text-muted">{MISSAO}</p>
          <p className="text-xs text-muted">
            Base: Lei 10.301/2019, art. 5º. Pela lei, o Arquivo Público é subordinado ao Gabinete da Secretaria
            Municipal de Administração, sob a responsabilidade da Divisão de Arquivo Público.
          </p>
        </section>

        <section className="flex flex-col gap-4">
          <h2 className="text-xl font-semibold text-foreground">Competências do Arquivo Público</h2>
          <ol className="flex list-decimal flex-col gap-2 pl-6 text-sm leading-relaxed text-muted">
            {COMPETENCIAS.map((c) => (
              <li key={c}>{c}</li>
            ))}
          </ol>
          <p className="text-xs text-muted">Base: Lei 10.301/2019, art. 5º, I a XI.</p>
        </section>

        <section className="flex flex-col gap-4">
          <h2 className="text-xl font-semibold text-foreground">O Sistema Municipal de Arquivos</h2>
          <p className="text-sm leading-relaxed text-muted">
            O SISMARQ organiza as atividades de gestão de documentos de todos os órgãos e entidades da administração
            municipal. O <strong>Arquivo Público</strong> é o órgão central. As unidades das secretarias que coordenam o
            protocolo e a gestão de documentos são os <strong>órgãos setoriais</strong>, e os servidores que chefiam essas
            atividades são os <strong>agentes setoriais</strong>. Cada órgão tem uma{" "}
            <strong>Comissão Setorial Permanente de Avaliação de Documentos (CSPAD)</strong>, e a{" "}
            <strong>Comissão Central (CCPAD)</strong> funciona junto ao Arquivo Público.
          </p>
          <div className="grid gap-4 sm:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle as="h3" className="text-base">Finalidades</CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="flex list-disc flex-col gap-1.5 pl-5 text-sm text-muted">
                  {FINALIDADES.map((f) => (
                    <li key={f}>{f}</li>
                  ))}
                </ul>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle as="h3" className="text-base">As três idades dos documentos</CardTitle>
              </CardHeader>
              <CardContent className="flex flex-col gap-3">
                {IDADES.map((i) => (
                  <div key={i.nome} className="flex flex-col gap-0.5">
                    <span className="text-sm font-semibold text-foreground">{i.nome}</span>
                    <span className="text-sm text-muted">{i.descricao}</span>
                    <span className="text-xs text-muted">{i.onde}</span>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
          <p className="text-xs text-muted">Base: Lei 10.301/2019, arts. 6º a 12, 18 e 31.</p>
        </section>

        <section className="flex flex-col gap-4">
          <h2 className="text-xl font-semibold text-foreground">SISMARQ e SEI: cada um no seu papel</h2>
          <div className="grid gap-4 sm:grid-cols-2">
            {[SEI_X_SISMARQ.sei, SEI_X_SISMARQ.sismarq].map((lado) => (
              <Card key={lado.titulo}>
                <CardHeader>
                  <CardTitle as="h3" className="text-base">{lado.titulo}</CardTitle>
                  <p className="text-xs font-semibold text-primary">{lado.escopo}</p>
                </CardHeader>
                <CardContent>
                  <ul className="flex list-disc flex-col gap-1.5 pl-5 text-sm text-muted">
                    {lado.itens.map((it) => (
                      <li key={it}>{it}</li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            ))}
          </div>
          <p className="text-xs text-muted">
            Base: artigo técnico do SISMARQ, Quadro 1. Os processos concluídos no SEI são recolhidos ao SISMARQ.
          </p>
        </section>

        <section className="flex flex-col gap-4">
          <h2 className="text-xl font-semibold text-foreground">O sistema informatizado</h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {SISTEMA.map((s) => (
              <Card key={s.titulo}>
                <CardHeader className="flex flex-row items-center gap-3 space-y-0">
                  <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    <ShieldCheck size={18} aria-hidden="true" />
                  </span>
                  <CardTitle as="h3" className="text-base">{s.titulo}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted">{s.texto}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section className="flex flex-col gap-4">
          <h2 className="flex items-center gap-2 text-xl font-semibold text-foreground">
            <Scale size={20} className="text-primary" aria-hidden="true" /> Base legal
          </h2>
          <Table caption="Normas que fundamentam o Arquivo Público Municipal e o SISMARQ." className="min-w-[500px]">
            <TableHead>
              <TableRow>
                <TableHeaderCell>Norma</TableHeaderCell>
                <TableHeaderCell>Assunto</TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {BASE_LEGAL.map((b) => (
                <TableRow key={b.norma}>
                  <TableCell className="align-top font-medium text-foreground">{b.norma}</TableCell>
                  <TableCell className="text-muted">{b.tema}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="text-sm text-muted">
            <p>
              As Tabelas de Temporalidade e Destinação de Documentos em vigor foram aprovadas pela Comissão Central
              Permanente de Avaliação de Documentos e publicadas no Diário Oficial do Município:
            </p>
            <ul className="mt-2 flex list-disc flex-col gap-1 pl-5">
              {TTDD_PUBLICACOES.map((p) => (
                <li key={p}>{p}</li>
              ))}
            </ul>
          </div>
          <Link href="/temporalidade" className="inline-flex items-center gap-1 text-sm font-semibold text-primary hover:underline">
            Consultar a Tabela de Temporalidade <ArrowRight size={14} aria-hidden="true" />
          </Link>
        </section>
      </div>
    </PublicShell>
  );
}
