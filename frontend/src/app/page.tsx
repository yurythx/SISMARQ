import type { Metadata } from "next";
import Link from "next/link";
import { Archive, ArrowRight, CalendarClock, FileSearch, HelpCircle, Landmark, MessageSquare, Newspaper, ScanLine } from "lucide-react";

import { PublicShell } from "@/components/layout/PublicShell";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/Card";
import { getServerBranding } from "@/lib/branding/server";
import { listPublishedServices } from "@/lib/catalog/catalog";
import { listPublicPosts } from "@/lib/blog/publicBlog";
import { ServiceIcon, getCategoryIconStyles } from "@/lib/catalog/icons";
import { APP_URL } from "@/lib/env";
import { DIGITALIZACAO_FLUXO, IDADES, LEI, RITO_ELIMINACAO } from "@/lib/sismarq/conteudo";

export const dynamic = "force-dynamic";

const description =
  "Arquivo Público Municipal de Rondonópolis e Sistema Municipal de Arquivos (SISMARQ), criados pela Lei Municipal nº 10.301/2019: acesso aos documentos públicos, Tabela de Temporalidade e gestão documental.";

export const metadata: Metadata = {
  title: "Início — Arquivo Público Municipal e SISMARQ",
  description,
  alternates: { canonical: `${APP_URL}` },
  openGraph: {
    title: "Início — Arquivo Público Municipal e SISMARQ",
    description,
    type: "website",
    url: `${APP_URL}`,
  },
};

export default async function HomePage() {
  const [services, posts, brand] = await Promise.all([
    listPublishedServices().catch(() => []),
    listPublicPosts("", "", 1, 6).catch(() => []),
    getServerBranding().catch(() => null),
  ]);

  return (
    <PublicShell>
      <div className="flex flex-col gap-20 pb-20">
        {/* Apresentação */}
        <section className="border-b border-surface-border bg-gradient-to-b from-surface-1 to-surface-2/40 px-6 py-20 sm:py-24">
          <div className="mx-auto flex max-w-5xl flex-col items-center text-center">
            <span className="inline-flex items-center gap-2 rounded-full border border-surface-border bg-surface px-4 py-1.5 text-xs font-semibold text-primary shadow-xs">
              <Landmark size={14} aria-hidden="true" /> {brand?.orgName || "Prefeitura Municipal de Rondonópolis"}
            </span>
            <h1 className="mt-6 text-4xl font-extrabold tracking-tight text-foreground sm:text-5xl">
              Arquivo Público Municipal <span className="text-accent">SISMARQ</span>
            </h1>
            <p className="mt-6 max-w-3xl text-lg leading-relaxed text-muted">
              O Sistema Municipal de Arquivos, criado pela <strong className="text-foreground">{LEI.numero}</strong>,
              organiza a gestão dos documentos da administração municipal e assegura ao cidadão o acesso pleno aos
              documentos públicos, ressalvados o sigilo legal e a proteção da intimidade.
            </p>
            <div className="mt-10 flex w-full flex-col items-stretch justify-center gap-3 sm:w-auto sm:flex-row sm:items-center">
              <Link
                href="/temporalidade"
                className="inline-flex items-center justify-center gap-2 rounded-xl bg-accent px-6 py-3.5 text-base font-semibold text-accent-foreground shadow-md hover:bg-accent/90"
              >
                <CalendarClock size={18} aria-hidden="true" /> Consultar a Tabela de Temporalidade
              </Link>
              <Link
                href="/servicos"
                className="inline-flex items-center justify-center gap-2 rounded-xl border border-surface-border bg-surface px-6 py-3.5 text-base font-semibold text-foreground shadow-xs hover:bg-surface-2"
              >
                <FileSearch size={18} aria-hidden="true" /> Serviços do Arquivo
              </Link>
            </div>
            <p className="mt-8 text-xs text-muted">
              {LEI.numero}, de {LEI.data}. Publicada no {LEI.publicacao}.
            </p>
          </div>
        </section>

        {/* Ciclo de vida */}
        <section className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6" aria-labelledby="ciclo">
          <div className="flex flex-col gap-2">
            <p className="dateline">Lei 10.301/2019, art. 18</p>
            <h2 id="ciclo" className="text-3xl font-bold tracking-tight text-foreground">
              O ciclo de vida dos documentos
            </h2>
            <p className="max-w-3xl text-base text-muted">
              Todo documento público municipal passa por três idades. A Tabela de Temporalidade define quanto tempo
              ele fica em cada uma e qual o seu destino: eliminação ou guarda permanente.
            </p>
          </div>
          <ol className="grid gap-6 md:grid-cols-3">
            {IDADES.map((idade, i) => (
              <li key={idade.nome}>
                <Card className="h-full border-surface-border">
                  <CardHeader className="flex flex-col gap-2">
                    <span className="text-xs font-bold uppercase tracking-wider text-primary">{i + 1}ª idade</span>
                    <CardTitle as="h3" className="text-lg">
                      {idade.nome}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="flex flex-col gap-2 text-sm text-muted">
                    <p>{idade.descricao}</p>
                    <p className="text-xs font-medium text-foreground/80">{idade.onde}</p>
                    <p className="text-xs">{idade.fonte}</p>
                  </CardContent>
                </Card>
              </li>
            ))}
          </ol>
        </section>

        {/* Digitalização */}
        <section className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6" aria-labelledby="digitalizacao">
          <div className="flex flex-col gap-2">
            <p className="dateline">Decreto 10.278/2020 · Lei 10.301/2019, art. 48</p>
            <h2 id="digitalizacao" className="text-3xl font-bold tracking-tight text-foreground">
              Como o documento digital entra no SISMARQ
            </h2>
            <p className="max-w-3xl text-base text-muted">
              Há dois caminhos, e eles não se misturam. O papel passa pelo scanner (PDF ou TIFF) com a resolução do
              Decreto federal. O PDF que já nasceu eletrônico (nato-digital) segue para PDF/A e OCR, sem fingir que foi
              escaneado. Em ambos os casos o hash e a guarda imutável ficam no repositório; digitalizar não libera a
              eliminação do permanente.
            </p>
            <p className="inline-flex items-center gap-2 text-sm text-primary">
              <ScanLine size={16} aria-hidden="true" /> Servidor autenticado envia o arquivo na ficha da unidade de
              arquivamento.
            </p>
          </div>
          <ol className="grid gap-4 md:grid-cols-2">
            {DIGITALIZACAO_FLUXO.map((etapa) => (
              <li key={etapa.titulo} className="flex flex-col gap-2 rounded-xl border border-surface-border bg-surface p-5 shadow-xs">
                <h3 className="text-base font-bold text-foreground">{etapa.titulo}</h3>
                <p className="text-sm text-muted">{etapa.descricao}</p>
                <p className="mt-auto text-xs text-muted">{etapa.fonte}</p>
              </li>
            ))}
          </ol>
        </section>

        {/* Rito de eliminação */}
        <section className="border-y border-surface-border bg-surface-1/70 px-6 py-16" aria-labelledby="rito">
          <div className="mx-auto flex max-w-6xl flex-col gap-8">
            <div className="flex flex-col gap-2">
              <p className="dateline">Lei 10.301/2019, arts. 19 e 33 a 40</p>
              <h2 id="rito" className="text-3xl font-bold tracking-tight text-foreground">
                Nenhum documento é eliminado sem rito
              </h2>
              <p className="max-w-3xl text-base text-muted">
                A eliminação depende da Tabela de Temporalidade, da avaliação da comissão, da autorização do Arquivo
                Público e de publicidade. Documentos permanentes nunca são eliminados, nem depois de digitalizados.
              </p>
            </div>
            <ol className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {RITO_ELIMINACAO.map((etapa, i) => (
                <li key={etapa.titulo} className="flex flex-col gap-2 rounded-xl border border-surface-border bg-surface p-5 shadow-xs">
                  <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-sm font-bold text-primary">
                    {i + 1}
                  </span>
                  <h3 className="text-base font-bold text-foreground">{etapa.titulo}</h3>
                  <p className="text-sm text-muted">{etapa.descricao}</p>
                  <p className="mt-auto text-xs text-muted">{etapa.fonte}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* Serviços */}
        <section className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6" aria-labelledby="servicos">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div className="flex flex-col gap-2">
              <p className="dateline">Carta de serviços</p>
              <h2 id="servicos" className="text-3xl font-bold tracking-tight text-foreground">
                Serviços do Arquivo Público
              </h2>
            </div>
            <Link href="/servicos" className="inline-flex items-center gap-1.5 text-sm font-semibold text-accent hover:underline">
              Ver todos <ArrowRight size={15} aria-hidden="true" />
            </Link>
          </div>
          {services.length === 0 ? (
            <p className="rounded-xl border border-dashed border-surface-border p-8 text-center text-sm text-muted">
              Nenhum serviço publicado no momento.
            </p>
          ) : (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {services.slice(0, 8).map((svc) => {
                const iconStyle = getCategoryIconStyles(svc.category);
                return (
                  <Link key={svc.slug} href={`/servicos/${svc.slug}`} className="group block h-full">
                    <Card className="flex h-full flex-col border-surface-border transition-all hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg">
                      <CardHeader className="flex flex-col gap-3">
                        <span className={`inline-flex h-11 w-11 items-center justify-center rounded-xl ${iconStyle.badge}`}>
                          <ServiceIcon name={svc.icon} size={22} aria-hidden="true" />
                        </span>
                        <CardTitle as="h3" className="text-base font-bold text-foreground group-hover:text-primary">
                          {svc.title}
                        </CardTitle>
                      </CardHeader>
                      <CardContent>
                        <CardDescription className="line-clamp-3 text-xs leading-relaxed text-muted">{svc.summary}</CardDescription>
                      </CardContent>
                    </Card>
                  </Link>
                );
              })}
            </div>
          )}
        </section>

        {/* Notícias */}
        <section className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6" aria-labelledby="noticias">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div className="flex flex-col gap-2">
              <span className="inline-flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-primary">
                <Newspaper size={16} aria-hidden="true" /> Notícias e orientações
              </span>
              <h2 id="noticias" className="text-3xl font-bold tracking-tight text-foreground">
                Gestão documental em Rondonópolis
              </h2>
            </div>
            <Link href="/blog" className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:underline">
              Ver todas <ArrowRight size={15} aria-hidden="true" />
            </Link>
          </div>
          {posts.length === 0 ? (
            <p className="rounded-xl border border-dashed border-surface-border p-8 text-center text-sm text-muted">
              Nenhuma publicação no momento.
            </p>
          ) : (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {posts.map((post) => (
                <Link key={post.id} href={`/blog/${post.slug}`} className="group block h-full">
                  <Card className="flex h-full flex-col border-surface-border transition-all hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg">
                    <CardHeader>
                      <CardTitle as="h3" className="text-base font-bold leading-snug text-foreground group-hover:text-primary">
                        {post.title}
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-1 flex-col justify-between gap-4">
                      <CardDescription className="line-clamp-3 text-xs leading-relaxed text-muted">{post.summary}</CardDescription>
                      <span className="inline-flex items-center gap-1 text-xs font-semibold text-accent">
                        Ler <ArrowRight size={13} aria-hidden="true" />
                      </span>
                    </CardContent>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </section>

        {/* Atalhos */}
        <section className="mx-auto grid w-full max-w-6xl gap-6 px-6 sm:grid-cols-3">
          {[
            { href: "/sobre", icon: Archive, titulo: "Sobre o Arquivo Público", texto: "Missão, competências, estrutura do SISMARQ e base legal." },
            { href: "/faq", icon: HelpCircle, titulo: "Perguntas frequentes", texto: "Acesso, Tabela de Temporalidade, eliminação e recolhimento, com o artigo da lei." },
            { href: "/contato", icon: MessageSquare, titulo: "Fale com o Arquivo", texto: "Envie sua dúvida ou solicitação ao Arquivo Público Municipal." },
          ].map(({ href, icon: Icon, titulo, texto }) => (
            <Link key={href} href={href} className="group block h-full">
              <div className="flex h-full flex-col gap-3 rounded-2xl border border-surface-border bg-surface p-6 transition-all hover:border-primary/40 hover:shadow-md">
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                  <Icon size={20} aria-hidden="true" />
                </span>
                <h3 className="text-base font-bold text-foreground group-hover:text-primary">{titulo}</h3>
                <p className="text-xs leading-relaxed text-muted">{texto}</p>
                <span className="mt-auto inline-flex items-center gap-1 text-xs font-semibold text-primary">
                  Acessar <ArrowRight size={12} aria-hidden="true" />
                </span>
              </div>
            </Link>
          ))}
        </section>
      </div>
    </PublicShell>
  );
}
