"use client";

import { ArrowLeft, Printer } from "lucide-react";
import Link from "next/link";
import { useParams } from "next/navigation";
import type { ReactNode } from "react";

import { TIPOS, ano, fmtDia, porExtenso } from "@/components/acervo/labels";
import { useOrgao } from "@/components/acervo/useOrgao";
import { DataState } from "@/components/sismarq/DataState";
import { Button } from "@/components/ui/Button";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoComissao, AcervoItemListagem, AcervoListagem } from "@/lib/sismarq/types";

const TITULOS: Record<string, string> = {
  "1": "Relação de eliminação de documentos",
  "2": "Termo de eliminação de documentos",
  "3": "Edital de ciência de eliminação de documentos",
  "4": "Ata de reunião para aprovação de eliminação de documentos",
};
const ROMANOS: Record<string, string> = { "1": "I", "2": "II", "3": "III", "4": "IV" };

function Assinatura({ children }: { children: ReactNode }) {
  return <div className="mt-10 border-t border-foreground pt-1 text-center text-xs uppercase">{children}</div>;
}

/** Tabela comum aos Anexos I, II e III. */
function Conjuntos({ itens }: { itens: AcervoItemListagem[] }) {
  return (
    <table className="w-full border-collapse border border-foreground text-xs">
      <caption className="sr-only">Conjuntos documentais</caption>
      <thead>
        <tr>
          <th scope="col" rowSpan={2} className="border border-foreground p-1">
            Código referente à classificação
          </th>
          <th scope="col" rowSpan={2} className="border border-foreground p-1">
            Assunto (descritor do código)
          </th>
          <th scope="colgroup" colSpan={2} className="border border-foreground p-1">
            Datas-limite
          </th>
          <th scope="colgroup" colSpan={2} className="border border-foreground p-1">
            Unidade de arquivamento
          </th>
          <th scope="col" rowSpan={2} className="border border-foreground p-1">
            Observações e/ou justificativas
          </th>
        </tr>
        <tr>
          <th scope="col" className="border border-foreground p-1">Ano inicial</th>
          <th scope="col" className="border border-foreground p-1">Ano final</th>
          <th scope="col" className="border border-foreground p-1">Quantificação</th>
          <th scope="col" className="border border-foreground p-1">Especificação</th>
        </tr>
      </thead>
      <tbody>
        {itens.map((u) => (
          <tr key={u.id}>
            <td className="border border-foreground p-1 font-mono">{u.codigo_classificacao}</td>
            <td className="border border-foreground p-1">{u.assunto}</td>
            <td className="border border-foreground p-1 text-center">{ano(u.data_inicial)}</td>
            <td className="border border-foreground p-1 text-center">{ano(u.data_final)}</td>
            <td className="border border-foreground p-1 text-center">{u.quantidade}</td>
            <td className="border border-foreground p-1">{u.especificacao || TIPOS[u.tipo].toLowerCase()}</td>
            <td className="border border-foreground p-1">
              {u.numero}
              {u.observacoes && ` — ${u.observacoes}`}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

/** "N unidade(s) de arquivamento (Q caixa(s)…)" — a mensuração dos Anexos. */
function mensuracao(itens: AcervoItemListagem[]): string {
  const q = itens.reduce((t, u) => t + u.quantidade, 0);
  return `${itens.length} unidade(s) de arquivamento (${q} ${q === 1 ? "volume" : "volumes"})`;
}

/** Anexos I a IV da Lei 10.301/2019, preenchidos com a relação, para
 * imprimir e assinar. */
export default function AnexoPage() {
  const { id, n } = useParams<{ id: string; n: string }>();
  const q = useApiQuery<AcervoListagem>(`v1/acervo/listagens/${id}`);
  const coms = useApiQuery<AcervoComissao[]>("v1/acervo/comissoes");
  const l = q.data;
  const { orgao, unidade } = useOrgao(l?.unidade_id);
  const comissao = coms.data?.find((c) => c.id === l?.comissao_id);
  const titulo = TITULOS[n];
  const relacionados = l?.itens.filter((u) => !u.retirada) ?? [];
  const eliminados = relacionados.filter((u) => !u.amostra);
  const anos = relacionados.flatMap((u) => [ano(u.data_inicial), ano(u.data_final)]).sort();

  return (
    <div className="mx-auto flex max-w-5xl flex-col gap-5 text-foreground">
      <div className="flex flex-wrap items-center justify-between gap-2 print:hidden">
        <Link href={`/acervo/eliminacao/${id}`} className="inline-flex items-center gap-1 text-sm text-muted hover:text-foreground">
          <ArrowLeft size={14} aria-hidden="true" /> Relação {l?.numero}
        </Link>
        <Button size="sm" variant="secondary" onClick={() => window.print()}>
          <Printer size={14} aria-hidden="true" className="mr-1" /> Imprimir
        </Button>
      </div>
      <DataState loading={q.isLoading} error={q.error} empty={!l || !titulo} emptyTitle="Documento não encontrado">
        {l && titulo && (
          <article aria-label={`Anexo ${ROMANOS[n]}`} className="flex flex-col gap-4 text-sm leading-relaxed">
            <header className="flex flex-col gap-1 border-b border-foreground pb-3 text-center">
              <p className="text-xs font-semibold">ANEXO {ROMANOS[n]}</p>
              <p className="text-xs font-semibold uppercase">Prefeitura Municipal de Rondonópolis · Sistema Municipal de Arquivos — SISMARQ</p>
              <h1 className="text-lg font-bold uppercase">{titulo}</h1>
            </header>
            {n !== "4" && (
              <dl className="grid gap-1 sm:grid-cols-2">
                <div>
                  <dt className="inline font-semibold">ÓRGÃO/ENTIDADE: </dt>
                  <dd className="inline">{orgao}</dd>
                </div>
                <div>
                  <dt className="inline font-semibold">UNIDADE/SETOR: </dt>
                  <dd className="inline">{unidade}</dd>
                </div>
              </dl>
            )}

            {n === "1" && (
              <>
                <p>
                  <strong>LISTAGEM Nº:</strong> {l.numero}
                </p>
                <Conjuntos itens={relacionados} />
                <dl className="grid gap-1 border border-foreground p-2 sm:grid-cols-2">
                  <div>
                    <dt className="inline font-semibold">DATAS-LIMITE GERAIS: </dt>
                    <dd className="inline">{anos.length ? `${anos[0]} a ${anos[anos.length - 1]}` : "—"}</dd>
                  </div>
                  <div>
                    <dt className="inline font-semibold">MENSURAÇÃO TOTAL: </dt>
                    <dd className="inline">{mensuracao(relacionados)}</dd>
                  </div>
                  <div>
                    <dt className="inline font-semibold">Conta(s) do(s) exercício(s) de: </dt>
                    <dd className="inline">{l.contas_exercicios || "—"}</dd>
                  </div>
                  <div>
                    <dt className="inline font-semibold">Conta(s) aprovada(s) pelo Tribunal de Contas em: </dt>
                    <dd className="inline">{fmtDia(l.contas_aprovadas_em)}</dd>
                  </div>
                </dl>
                <div className="grid gap-6 sm:grid-cols-3">
                  <Assinatura>Responsável pela seleção</Assinatura>
                  <Assinatura>Presidente da Comissão Permanente de Avaliação de Documentos</Assinatura>
                  <Assinatura>Gestor do órgão/entidade</Assinatura>
                </div>
                <p className="mt-6">
                  LOCAL/DATA: Rondonópolis, {l.autorizada_em ? fmtDia(l.autorizada_em) : "____/____/______"}
                </p>
                <Assinatura>Autorizo — Titular da Divisão de Arquivo Público</Assinatura>
              </>
            )}

            {n === "2" && (
              <>
                <p className="text-justify">
                  Aos {porExtenso(l.termo_eliminado_em).dia} dias do mês de {porExtenso(l.termo_eliminado_em).mes} do ano de{" "}
                  {porExtenso(l.termo_eliminado_em).ano}, o(a) {unidade}, de acordo com o que consta da Listagem de Eliminação de Documentos nº{" "}
                  {l.numero}, aprovada pelo(a) titular da Divisão de Arquivo Público do Município de Rondonópolis, Estado de Mato Grosso, por intermédio
                  da autorização de {fmtDia(l.autorizada_em)} e respectivo Edital de Ciência de Eliminação de Documentos nº {l.edital_numero}, publicado no
                  Diário Oficial do Município de Rondonópolis em {fmtDia(l.edital_publicado_em)}
                  {l.edital_diario && ` (${l.edital_diario})`}, procedeu à eliminação de {mensuracao(eliminados)} dos documentos relativos aos conjuntos
                  documentais abaixo, integrantes do acervo do(a) {orgao}.
                </p>
                <Conjuntos itens={eliminados} />
                <p>
                  <strong>SERVIDOR RESPONSÁVEL PELO ACOMPANHAMENTO DA ELIMINAÇÃO:</strong> {l.termo_acompanhamento}
                </p>
                {l.termo_destinacao_papel && (
                  <p>
                    <strong>Destinação do papel (art. 40):</strong> {l.termo_destinacao_papel}
                  </p>
                )}
                <p>LOCAL/DATA: Rondonópolis, {fmtDia(l.termo_eliminado_em)}</p>
                <Assinatura>Presidente da Comissão Permanente de Avaliação de Documentos</Assinatura>
              </>
            )}

            {n === "3" && (
              <>
                <p className="text-justify">
                  <strong>Nº {l.edital_numero}</strong> — O(A) Presidente da Comissão Permanente de Avaliação, designado(a) pela{" "}
                  {comissao?.portaria || "Portaria nº ____"}, de acordo com a Listagem de Eliminação de Documentos nº {l.numero}, autorizada pela
                  Divisão de Arquivo Público do Município de Rondonópolis, Estado de Mato Grosso — Órgão Central do SISMARQ, faz saber a quem possa
                  interessar que a partir do 30º dia subsequente à data de publicação deste Edital no Diário Oficial do Município (Lei Municipal nº
                  10.301/2019, art. 37), se não houver oposição, o(a) {unidade} eliminará {mensuracao(relacionados)} dos documentos relativos aos
                  conjuntos documentais relacionados abaixo, do(a) {orgao}.
                </p>
                <Conjuntos itens={relacionados} />
                <p className="text-justify">
                  Os interessados, no prazo citado, poderão requerer, às suas expensas, o desentranhamento de documentos ou cópias de peças do processo,
                  mediante petição, desde que tenham a respectiva qualificação e demonstração de legitimidade do pedido, dirigida à Comissão Permanente
                  de Avaliação do(a) {orgao}.
                </p>
                <p>LOCAL/DATA: Rondonópolis, {fmtDia(l.edital_publicado_em)}</p>
                <Assinatura>Presidente da Comissão Permanente de Avaliação de Documentos</Assinatura>
              </>
            )}

            {n === "4" && (
              <>
                <p className="text-center font-semibold">
                  ATA Nº {l.ata_numero} — REUNIÃO PARA APROVAÇÃO DE ELIMINAÇÃO DE DOCUMENTOS
                </p>
                <p className="text-justify">
                  Ao(s) {porExtenso(l.ata_em).dia} dia(s) do mês de {porExtenso(l.ata_em).mes} de {porExtenso(l.ata_em).ano}, às{" "}
                  {l.ata_em ? l.ata_em.slice(11, 16).replace(":", "h") : "____"}, nas dependências do(a) {l.ata_local || "________"}, reuniram-se os
                  membros da Comissão Setorial Permanente de Avaliação de Documentos — CSPAD do(a) {orgao}, para aprovação da Listagem de Eliminação de
                  Documentos nº {l.numero}.
                </p>
                <section aria-label="Manifestações dos presentes" className="border border-foreground p-2">
                  <h2 className="font-semibold uppercase">Manifestações dos presentes</h2>
                  <p className="whitespace-pre-line">{l.ata_manifestacoes || "Sem manifestações."}</p>
                </section>
                <p className="text-justify">
                  Após as manifestações e esclarecimentos, a referida eliminação foi aprovada pelos membros da CSPAD presentes na reunião, conforme a
                  lista de presença abaixo, com o parecer jurídico juntado (art. 32, § 4º).
                </p>
                <section aria-label="Lista de presença">
                  <h2 className="font-semibold uppercase">Lista de presença</h2>
                  <div className="grid gap-4 sm:grid-cols-2">
                    {l.ata_presentes.map((p) => (
                      <Assinatura key={p}>{comissao?.membros.find((m) => m.user_id === p)?.nome || p}</Assinatura>
                    ))}
                  </div>
                </section>
                <Assinatura>Secretário(a) da reunião</Assinatura>
              </>
            )}
          </article>
        )}
      </DataState>
    </div>
  );
}
