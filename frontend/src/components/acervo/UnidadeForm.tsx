"use client";

import { useState, type FormEvent } from "react";

import { NIVEIS, TIPOS, ajudaEvento } from "@/components/acervo/labels";
import { UnidadeDonaSelect } from "@/components/iam/UnidadeDonaSelect";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import { useApiQuery } from "@/lib/api/swr";
import type { AcervoUnidade, ClassificacaoTTDD } from "@/lib/sismarq/types";

/** Corpo enviado à API (datas em AAAA-MM-DD). */
export interface UnidadeBody {
  tipo: string;
  codigo_classificacao: string;
  assunto: string;
  unidade_id: string;
  data_inicial: string;
  data_final: string;
  quantidade: number;
  especificacao: string;
  observacoes: string;
  data_evento: string;
  contas_aprovadas_em: string;
  nivel_acesso: string;
  hipotese_legal: string;
  local_id: string | null;
  posicao: string;
}

const dia = (v?: string | null) => (v ? v.slice(0, 10) : "");

/** A série da TTDD que o código digitado alcança (prazos e destinação). */
function SerieDaTTDD({ codigo }: { codigo: string }) {
  const valido = /^\d+(\.\d+){2,4}$/.test(codigo);
  const serie = useApiQuery<ClassificacaoTTDD>(valido ? `v1/atlas/ttdd/${encodeURIComponent(codigo)}` : null);
  if (!valido) return <p className="text-xs text-muted">Informe o código completo da série (ex.: 2.0.01.00.01).</p>;
  if (serie.error) return <p className="text-xs text-danger">Série não encontrada na Tabela de Temporalidade.</p>;
  const s = serie.data;
  if (!s) return <p className="text-xs text-muted">Consultando a TTDD…</p>;
  return (
    <p className="text-xs text-muted" aria-live="polite">
      <strong className="text-foreground">{s.descritor}</strong>
      {s.revogada_em ? " — série revogada: escolha uma série em vigor." : ""}
    </p>
  );
}

/** Formulário da unidade de arquivamento (cadastro e alteração). */
export function UnidadeForm({
  inicial,
  pending,
  onSubmit,
  enviar,
}: {
  inicial?: AcervoUnidade;
  pending: boolean;
  onSubmit: (body: UnidadeBody) => void;
  enviar: string;
}) {
  const [codigo, setCodigo] = useState(inicial?.codigo_classificacao ?? "");
  const [nivel, setNivel] = useState<string>(inicial?.nivel_acesso ?? "PUBLICO");
  const serieOk = /^\d+(\.\d+){2,4}$/.test(codigo);
  const serie = useApiQuery<ClassificacaoTTDD>(serieOk ? `v1/atlas/ttdd/${encodeURIComponent(codigo)}` : null);

  function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const txt = (k: string) => String(f.get(k) ?? "").trim();
    onSubmit({
      tipo: txt("tipo"),
      codigo_classificacao: txt("codigo_classificacao"),
      assunto: txt("assunto"),
      unidade_id: txt("unidade_id"),
      data_inicial: txt("data_inicial"),
      data_final: txt("data_final"),
      quantidade: Number(txt("quantidade") || "1"),
      especificacao: txt("especificacao"),
      observacoes: txt("observacoes"),
      data_evento: txt("data_evento"),
      contas_aprovadas_em: txt("contas_aprovadas_em"),
      nivel_acesso: nivel,
      hipotese_legal: nivel === "PUBLICO" ? "" : txt("hipotese_legal"),
      local_id: inicial?.local_id ?? null,
      posicao: inicial?.posicao ?? "",
    });
  }

  return (
    <form onSubmit={submit} className="flex flex-col gap-6" aria-label="Unidade de arquivamento">
      <fieldset className="grid gap-4 sm:grid-cols-2">
        <legend className="mb-2 text-base font-semibold text-foreground">Identificação e classificação</legend>
        <Select
          id="tipo"
          name="tipo"
          label="Tipo *"
          required
          defaultValue={inicial?.tipo ?? "CAIXA"}
          options={Object.entries(TIPOS).map(([value, label]) => ({ value, label }))}
        />
        <UnidadeDonaSelect id="unidade_id" permission="acervo:manage" label="Unidade produtora" defaultValue={inicial?.unidade_id} />
        <div className="flex flex-col gap-1 sm:col-span-2">
          <Input
            id="codigo_classificacao"
            name="codigo_classificacao"
            label="Código de classificação (série da TTDD) *"
            required
            value={codigo}
            onChange={(e) => setCodigo(e.target.value.trim())}
            aria-describedby="codigo-ajuda"
          />
          <div id="codigo-ajuda">
            <SerieDaTTDD codigo={codigo} />
          </div>
        </div>
        <div className="sm:col-span-2">
          <Input
            id="assunto"
            name="assunto"
            label="Assunto (descritor)"
            maxLength={500}
            defaultValue={inicial?.assunto}
            placeholder="Em branco: o descritor da série"
          />
        </div>
      </fieldset>

      <fieldset className="grid gap-4 sm:grid-cols-2">
        <legend className="mb-2 text-base font-semibold text-foreground">Datas-limite e quantificação (Anexos I a III)</legend>
        <Input id="data_inicial" name="data_inicial" type="date" label="Data inicial *" required defaultValue={dia(inicial?.data_inicial)} />
        <Input id="data_final" name="data_final" type="date" label="Data final *" required defaultValue={dia(inicial?.data_final)} />
        <Input id="quantidade" name="quantidade" type="number" min={1} max={100000} label="Quantidade" defaultValue={inicial?.quantidade ?? 1} />
        <Input
          id="especificacao"
          name="especificacao"
          label="Especificação"
          maxLength={120}
          placeholder="Ex.: caixas-arquivo, volumes, metros lineares"
          defaultValue={inicial?.especificacao}
        />
        <div className="flex flex-col gap-1">
          <Input id="data_evento" name="data_evento" type="date" label="Data do evento" defaultValue={dia(inicial?.data_evento)} aria-describedby="evento-ajuda" />
          <p id="evento-ajuda" className="text-xs text-muted">
            {ajudaEvento(serie.data?.fase_corrente_condicao || serie.data?.fase_interm_condicao, serie.data?.descritor)}
          </p>
        </div>
        <div className="flex flex-col gap-1">
          <Input
            id="contas_aprovadas_em"
            name="contas_aprovadas_em"
            type="date"
            label="Contas aprovadas pelo Tribunal de Contas em"
            defaultValue={dia(inicial?.contas_aprovadas_em)}
            aria-describedby="contas-ajuda"
          />
          <p id="contas-ajuda" className="text-xs text-muted">
            Séries financeiras: data em que o Tribunal de Contas aprovou as contas do exercício (quadro do Anexo I). Sem essa data a
            eliminação não se confunde com o fim da fase por anos. A TTDD pode ser republicada; a unidade já cadastrada não some.
          </p>
        </div>
      </fieldset>

      <fieldset className="grid gap-4 sm:grid-cols-2">
        <legend className="mb-2 text-base font-semibold text-foreground">Acesso (art. 2º)</legend>
        <Select
          id="nivel_acesso"
          label="Nível de acesso"
          value={nivel}
          onChange={(e) => setNivel(e.target.value)}
          options={Object.entries(NIVEIS).map(([value, label]) => ({ value, label }))}
        />
        {nivel !== "PUBLICO" && (
          <Input
            id="hipotese_legal"
            name="hipotese_legal"
            label="Hipótese legal da restrição *"
            required
            maxLength={500}
            defaultValue={inicial?.hipotese_legal}
          />
        )}
      </fieldset>

      <Textarea id="observacoes" name="observacoes" label="Observações e justificativas" maxLength={4000} rows={3} defaultValue={inicial?.observacoes} />

      <div>
        <Button type="submit" disabled={pending}>
          {enviar}
        </Button>
      </div>
    </form>
  );
}
