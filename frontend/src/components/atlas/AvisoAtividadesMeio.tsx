/** Lei 10.301/2019, art. 26: o município adota o código CONARQ de
 * atividades-meio (classes 01 a 07). As tabelas já carregadas são as de
 * atividades-fim publicadas no DIORONDON-e (ADR 043). */
export function AvisoAtividadesMeio() {
  return (
    <aside className="rounded-lg border border-surface-border bg-surface px-4 py-3 text-sm text-muted" role="note">
      <p>
        <strong className="text-foreground">Atividades-meio (Lei 10.301/2019, art. 26).</strong> O
        município adota o Código de Classificação e a Tabela de Temporalidade de Documentos de Arquivo
        para as atividades-meio da Administração Pública (CONARQ), classes 01 a 07, porque o art. 26
        determina essa adoção para o que não é atividade-fim do órgão. As tabelas já nesta consulta são
        as TTDDs de atividades-fim publicadas no DIORONDON-e. A carga das classes 01–07 entra quando a
        Superintendência de Arquivo Público importar a tabela CONARQ em Atualizar TTDD — sem inventar
        códigos municipais.
      </p>
    </aside>
  );
}
