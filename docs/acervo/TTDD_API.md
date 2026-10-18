# Contrato da TTDD para outros sistemas (art. 45)

A Lei 10.301/2019, art. 45, manda que a Tabela de Temporalidade esteja
disponível aos sistemas informatizados de protocolo. Este texto é o
contrato **já em produção** no SISMARQ — leitura anônima, sem inventar
cliente SEI ([ADR 042](../adr/042-ingresso-sei-adiado.md)).

Quem for ligar o SEI, o protocolo de outra secretaria ou um SIGAD usa
**estas URLs**. Não consultar o Typesense. Não chamar `/atlas/admin` nem
`POST /atlas/chat` sem autenticação.

Índice do domínio: [README.md](README.md). Campos da série: [DOMINIO.md](DOMINIO.md) §6.

---

## 1. Como chegar

| Caminho | Uso |
|---|---|
| `GET /api/v1/atlas/ttdd…` na API (rede interna, VM `8010`) | Integração servidor a servidor |
| `GET /api/public/v1/atlas/ttdd…` no frontend (BFF allowlist) | Site e clientes que só falam com a borda HTTPS |

O BFF **não** manda token. O IP do visitante só segue em `X-Forwarded-For`
atrás de proxy confiável. Há rate limit por IP no backend.

Envelope JSON da plataforma: `{ "data": …, "error": null }` (ou `error`
RFC 7807 negociado). Listagens paginadas trazem `meta` (`page`,
`page_size`, `total`). Padrão: `page=1`, `page_size=20`, teto 100
(`MAX_PAGE_SIZE`; absoluto 200).

---

## 2. Recursos

Todas as listagens e a exportação **omitam séries revogadas**. A série
revogada continua em `GET …/{codigo}` e em `…/historico` (documentos já
produzidos).

### `GET /atlas/ttdd`

Lista séries vigentes.

| Query | Significado |
|---|---|
| `q` | Busca full-text (até 200 caracteres) |
| `codigo` | Prefixo hierárquico válido: `2.0`, `2.0.01`, `2.0.01.00` ou série `2.0.01.00.00` |
| `page`, `page_size` | Paginação |

`codigo` inválido → 400.

### `GET /atlas/ttdd/estrutura`

Árvore órgão → função → subfunção com **contagem de séries vigentes**.
Sem paginação.

### `GET /atlas/ttdd/exportar`

CSV da mesma filtragem da lista, todas as páginas (teto interno 20 000).

- `Content-Type: text/csv; charset=utf-8`
- BOM UTF-8, separador `;`, arquivo `ttdd.csv`
- Colunas: Código, Série documental, Órgão, Função, Subfunção, Fase
  corrente, Fase intermediária, Destinação final, Observações,
  Recomendação da subfunção, Fonte

As fases no CSV são o texto da TTDD (anos ou condição), não os inteiros
do JSON.

### `GET /atlas/ttdd/{codigo}`

Uma série, **inclusive revogada**. 400 se o código não for duplex; 404 se
não existir.

JSON (campos estáveis):

```json
{
  "codigo": "2.0.02.00.07",
  "descritor": "…",
  "fase_corrente_anos": 5,
  "fase_corrente_condicao": "",
  "fase_interm_anos": null,
  "fase_interm_condicao": "Enquanto estiver vigorando",
  "destinacao_final": "ELIMINACAO",
  "observacoes": "",
  "subfuncao": { "codigo": "2.0.02.00", "nome": "…", "funcao": { "codigo": "2.0.02", "nome": "…", "orgao": { "prefixo": "2.0", "nome": "…", "edicao_diario": "…", "data_publicacao": "…", "versao": "…" } } },
  "revogada_em": null,
  "revogada_edicao": "",
  "created_at": "…"
}
```

- Cada fase: **anos XOR condição** (um dos dois preenchido; intermediária
  pode não ter fase).
- `destinacao_final`: `ELIMINACAO` \| `GUARDA_PERMANENTE` \| `null`
  (indefinida na TTDD — as seis séries “X”).
- `revogada_em` não nulo = fora da TTDD vigente; o código permanece.

### `GET /atlas/ttdd/{codigo}/historico`

Mudanças da série, da mais recente à mais antiga. Eventos: `ALTERADA`,
`REVOGADA`, `RESTABELECIDA`. Cada item traz `anterior` (prazos de antes),
`edicao_diario` e `registrado_em`.

### `GET /atlas/ttdd/{codigo}/modelos`

Modelos de documento ligados à série (biblioteca do Atlas). Pode estar
vazio.

---

## 3. O que não é este contrato

| Recurso | Por quê |
|---|---|
| `POST /atlas/chat` | Assistente; exige `atlas:read` |
| `/atlas/admin/*` | Gestão; `atlas:manage` global |
| Carga da TTDD | Tela autenticada ou `make ttdd-aplicar` |
| Acervo (`/acervo`, `/arquivo`) | Unidades de arquivamento, não a norma |
| Pacote SEI / AIP | ADR 042; sem URL nem token |

Organograma funcional da TTDD: `GET /atlas/organograma` (só séries
publicadas; rascunhos de procedimento não entram). O BFF público libera
o prefixo `v1/atlas/organograma`.

Procedimentos homologados: `GET /atlas/workflows` (mesmo BFF). Não
substituem a TTDD.

---

## 4. Como integrar (checklist)

1. Classificar o processo com `codigo` de **série vigente**
   (`revogada_em == null`).
2. Se a carga da TTDD revogar a série, o processo antigo **não** muda de
   código; o sistema de protocolo consulta o histórico.
3. Prazos: preferir os campos `fase_*_anos` / `fase_*_condicao` ao texto
   do CSV, para não reinterpretar “enquanto vigorar”.
4. Não cachear a tabela inteira sem `ETag`/revalidação; republicação no
   Diário muda destinação. Recarregar após aviso operacional ou no
   máximo no dia.
5. Erros: 400 código inválido, 404 série inexistente, 429 rate limit.

Quando a Prefeitura entregar versão e API do SEI, o adaptador do Acervo
usa este contrato para validar a série do envelope — não uma cópia
paralela da tabela.
