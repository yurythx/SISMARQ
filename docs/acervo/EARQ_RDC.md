# e-ARQ Brasil e RDC-Arq — mapa para implementações

Como o governo federal organiza SIGAD e repositório, e o que o SISMARQ
já cobre. **Não é avaliação CONARQ nem selo ISO 16363.** Serve à Fase 8
(art. 13 da Lei 10.301) e a PRs de preservação. Linha de base:
[DOMINIO.md](DOMINIO.md), [ADR 045](../adr/045-linha-de-base-arquivistica.md).

Fontes: [Resolução CONARQ nº 50/2022](https://www.gov.br/conarq/pt-br/legislacao-arquivistica/resolucoes-do-conarq/resolucao-no-50-de-06-de-maio-de-2022)
(e-ARQ Brasil v2); [Resolução CONARQ nº 51/2023](https://www.gov.br/conarq/pt-br/assuntos/noticias/diretrizes-para-a-implementacao-de-repositorios-arquivisticos-digitais-confiaveis-versao-2/DIRETRIZESRDCArq12DEZ2023.pdf)
(RDC-Arq); Decreto 10.278/2020; IN SGD/ME nº 1/2019 (Executivo federal —
referência, não obrigação municipal).

---

## 1. Dois sistemas, dois papéis

| | SIGAD (e-ARQ) | RDC-Arq |
|---|---|---|
| Função | Gestão: produzir, classificar, tramitar, avaliar, eliminar ou recolher | Preservação e acesso pelo tempo necessário; autenticidade e relação orgânica |
| Idades | Corrente e intermediária (o permanente recolhe para o RDC) | Permanente obrigatório; nas outras idades só se associado a um SIGAD |
| SISMARQ hoje | Núcleo municipal (Acervo + Atlas + Trâmite) | Guarda reforçada (PDF/A, hash, object lock, fixity) — **não certificado** |

O Arquivo Nacional recomenda SIGAD conforme e-ARQ e integração com RDC-Arq
para a cadeia de custódia. No município o art. 13 já pede o informatizado
alinhado ao e-ARQ. Evoluir o MinIO atual a um RDC exige ADR próprio
(pacote AIP, metadados de preservação, política de migração de formato).

---

## 2. Inventário temático (trabalho, não auditoria)

Legenda: **Coberto** no código e nos testes de domínio; **Parcial** com
lacuna explícita; **Fora** até decisão/ato; **Humano** (homologação ou DPO).

| Tema e-ARQ / federal | Situação | Onde | Próximo passo honesto |
|---|---|---|---|
| Captura e registro do documento de arquivo | Coberto | UA, ingresso Trâmite | SEI: envelope só (ADR 042) |
| Classificação (plano / TTDD) | Parcial | Atlas 1.686 séries fim; contrato [TTDD_API.md](TTDD_API.md) | Carga CONARQ 01–07 oficial; Educação |
| Tramitação na corrente | Parcial | Trâmite; não é SEI | Não fingir PEN |
| Avaliação e destinação | Coberto | Prazos + destinação da série | Validar conta (C1) |
| Eliminação com instrumento | Coberto | Anexos I–IV, 30 dias | C2–C4 no papel |
| Recolhimento permanente | Coberto | Art. 53, registro de entrada | Pacote AIP no permanente digital |
| Controle de acesso e sigilo | Coberto | Nível + escopo IAM; busca pela API | RIPD do DPO |
| Captura da digitalização 10.278 | Parcial | SCAN DPI, XMP, PDF/A-2b | PAdES; veraPDF; pasta scanner |
| Nato-digital | Coberto | `NATO_DIGITAL` sem DPI | Não rotular como 10.278 |
| Integridade (hash, WORM, fixity) | Coberto | Object lock + worker | Política de backup medida |
| Metadados e-ARQ parte II (6 entidades) | Parcial | UA + XMP govbr: + auditoria | Mapear agente/evento de preservação |
| Pesquisa e acesso | Coberto | API Typesense/PG; ostensivo | — |
| Auditoria imutável | Coberto | Cadeia SHA-256 + export WORM | — |
| Descrição multinível NOBRADE/ISAD(G) | Parcial | Ficha + instrumento de recolhimento | Níveis fundo/série/dossiê se a Superintendência pedir |
| Assinatura ICP-Brasil | Fora | `pades=AUSENTE` | Certificado municipal (0.6) |
| RDC certificado / ISO 16363 | Fora | — | Autoavaliação depois do AIP |
| e-ARQ checklist completo (Fase 8) | Humano | Este arquivo é o mapa | Planilha obrigatório/recomendado com a Divisão |
| Capacitação e ato art. 64 | Humano | — | Fora do código |

---

## 3. O que um PR de “RDC” pode e não pode fazer

**Pode (com ADR novo):** validação veraPDF no worker após o ocrmypdf;
envelope AIP interno (identificador, hashes, metadados, trilha) para
permanente já `PRONTO`; gravar evento de preservação na fixity. O contrato
da API da TTDD (art. 45) está em [TTDD_API.md](TTDD_API.md).

**Não pode:** afirmar repositório confiável certificado; assinar PAdES com
chave local; implementar cliente SEI; carregar classes 01–07 inventadas;
marcar a coluna do roadmap como “conforme e-ARQ”.

---

## 4. Pacote AIP (alvo, não implementação)

O recolhimento digital federal não é “um PDF no MinIO”. É um pacote com:

- componentes digitais (PDF/A);
- metadados de documento, classe, agente, evento de gestão e de preservação
  (e-ARQ v2 parte II);
- hashes conferidos na ingestão;
- trilha de custódia.

Quando o SEI existir, o ADR 042 liga o adaptador nesse espírito. Antes,
o AIP interno do permanente digitalizado pode nascer **sem** URL do SEI.

---

## 5. Como atualizar este mapa

Ao fechar um item da tabela, mude a coluna Situação no mesmo commit da
entrega e cite o ADR. Avaliação humana e-ARQ (Fase 8) preenche uma
planilha à parte; não substitua este arquivo por um “100% conforme”.
