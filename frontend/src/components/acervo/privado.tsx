"use client";

import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { Textarea } from "@/components/ui/Textarea";
import type { AcervoPrivado } from "@/lib/sismarq/types";

/** Corpo do cadastro a partir do formulário (cadastro e alteração). */
export function corpoPrivado(f: FormData) {
  const txt = (k: string) => String(f.get(k) ?? "").trim();
  return {
    titular: txt("titular"),
    tipo_titular: txt("tipo_titular"),
    documento: txt("documento"),
    contato: txt("contato"),
    descricao: txt("descricao"),
    datas_limite: txt("datas_limite"),
    mensuracao: txt("mensuracao"),
    local_guarda: txt("local_guarda"),
    servico_publico: f.get("servico_publico") === "on",
    acesso_autorizado: f.get("acesso_autorizado") === "on",
    acesso_condicoes: txt("acesso_condicoes"),
    convenio: txt("convenio"),
  };
}

/** Campos do cadastro, com os valores atuais na alteração. */
export function CamposPrivado({ p }: { p?: AcervoPrivado }) {
  return (
    <>
      <div className="grid gap-3 sm:grid-cols-[1fr_12rem_12rem]">
        <Input id="titular" name="titular" label="Titular *" required maxLength={200} defaultValue={p?.titular} />
        {p ? (
          <input type="hidden" name="tipo_titular" value={p.tipo_titular} />
        ) : (
          <Select
            id="tipo_titular"
            name="tipo_titular"
            label="Tipo *"
            defaultValue="PF"
            options={[
              { value: "PF", label: "Pessoa física" },
              { value: "PJ", label: "Pessoa jurídica" },
            ]}
          />
        )}
        <Input id="documento" name="documento" label="CPF/CNPJ" maxLength={20} defaultValue={p?.documento} />
      </div>
      <Textarea
        id="descricao"
        name="descricao"
        label="Conjuntos e relevância (art. 56) *"
        required
        rows={3}
        maxLength={8000}
        defaultValue={p?.descricao}
        placeholder="O que há no arquivo e por que é relevante para a história, a cultura ou o desenvolvimento do município"
      />
      <div className="grid gap-3 sm:grid-cols-3">
        <Input id="datas_limite" name="datas_limite" label="Datas-limite" maxLength={100} defaultValue={p?.datas_limite} placeholder="1953-1970" />
        <Input id="mensuracao" name="mensuracao" label="Mensuração" maxLength={200} defaultValue={p?.mensuracao} />
        <Input id="contato" name="contato" label="Contato do titular" maxLength={500} defaultValue={p?.contato} />
      </div>
      <Input id="local_guarda" name="local_guarda" label="Onde o acervo está guardado" maxLength={500} defaultValue={p?.local_guarda} />
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" name="acesso_autorizado" defaultChecked={p?.acesso_autorizado} /> O titular autoriza o acesso do público (§ 2º)
      </label>
      <Input id="acesso_condicoes" name="acesso_condicoes" label="Condições de acesso" maxLength={2000} defaultValue={p?.acesso_condicoes} />
      <Input id="convenio" name="convenio" label="Convênio de assistência técnica (art. 57)" maxLength={1000} defaultValue={p?.convenio} />
    </>
  );
}
