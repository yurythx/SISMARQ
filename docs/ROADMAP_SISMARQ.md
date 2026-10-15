# Roadmap do SISMARQ — Sistema Municipal de Arquivos de Rondonópolis

Para: a equipe que desenvolve o sistema e a Divisão (Superintendência) de
Arquivo Público, que vai validar cada fase.

Este roadmap leva a plataforma de hoje até o **sistema oficial do Arquivo
Público Municipal**. Ele segue três documentos, que são a referência do
projeto. Quando o código ou um ADR divergir deles, os documentos prevalecem
e a divergência é registrada na seção 6.

| Documento (em [`docs/SISMARQ/`](SISMARQ/)) | O que define |
|---|---|
| **Lei Municipal nº 10.301/2019** | Cria o Arquivo Público e o SISMARQ. Define as regras: as três idades, a classificação, a TTDD, as comissões, a eliminação com os Anexos I a IV, o recolhimento e as penalidades. |
| **Tabela de Temporalidade (TTDD)** | As TTDDs aprovadas pela CCPAD e publicadas no DIORONDON-e, de ago/2025 a set/2026: 13 tabelas, por secretaria. É o mesmo arquivo de `docs/ttdd.pdf`, já carregado no Atlas com 1.686 séries. |
| **Artigo científico (SISMARQ e RDC-Arq)** | A solução técnica. O SEI cuida da fase corrente e o SISMARQ das fases intermediária e permanente. O artigo descreve três módulos (pipeline, taxonomia e descarte), as exigências do Decreto 10.278/2020 e a busca no Typesense. |

Registrado em 2026-10-15; **Fases 1 e 2 entregues** no mesmo dia (ADRs 032 e 033 — falta a validação da contagem dos prazos pela Divisão de Arquivo e a ponte com o Trâmite). Os artigos citados são da Lei 10.301/2019, salvo
indicação em contrário.

---

## 1. Onde estamos

O SISMARQ é construído sobre a plataforma que já existe. O núcleo e boa
parte do que a lei exige já estão prontos.

| Já existe | Atende a |
|---|---|
| Kernel de plug-ins, IAM com escopo por entidade, unidade e departamento | Órgão central e setoriais (arts. 8 a 12): cada secretaria só gere o que é dela |
| Auditoria imutável com cadeia SHA-256 e cópia WORM | Rastreabilidade e prova de cada ato (art. 38, parágrafo único; art. 65) |
| **Atlas: a TTDD oficial** com vigência, histórico, revogação, organograma e atualização pela tela | Instrumentos de gestão (arts. 22 a 30) e reexame da TTDD (art. 52) |
| Trâmite com classificação pela TTDD e listagem de guarda dos processos encerrados (ADR 026) | Ponte entre a fase corrente e o SISMARQ |
| Signum: assinatura eletrônica com reautenticação | Atas e termos (Anexos II e IV). **Não** tem certificado ICP-Brasil (ver Fase 5) |
| Arquivos no MinIO, Typesense (ADR 031), notificações, LGPD, site público acessível (e-MAG) | Armazenamento, busca, avisos, acesso do cidadão (arts. 2 e 7, I) |

**O que falta** é o coração do SISMARQ:
- o **acervo**: o documento de arquivo, com fase, destinação e
  localização;
- as **travas legais** no domínio;
- o fluxo de **avaliação e eliminação** das comissões, com os Anexos I a IV;
- o **recolhimento** à guarda permanente;
- a **digitalização com valor probatório**;
- a **busca do acervo**.

---

## 2. Princípios que valem para todas as fases

1. **A lei manda no código.** Toda regra da Lei 10.301 vira regra do
   domínio, testada e com o artigo citado no código. A tela nunca é a única
   barreira.
2. **Travas que não se contornam.** Nenhuma combinação de permissões, nem a
   do administrador, consegue:
   - eliminar um documento permanente (arts. 20 e 42);
   - eliminar um documento com ação judicial pendente (art. 44);
   - eliminar sem autorização do Arquivo Público (art. 65);
   - eliminar antes do fim do prazo do edital (art. 37, § 2º).
3. **Os modelos oficiais saem iguais aos da lei.** A Relação, o Termo, o
   Edital e a Ata são gerados com os campos e textos dos Anexos I a IV.
4. **SEI e SISMARQ não se sobrepõem** (artigo, Quadro 1):
   - o SEI fica com a fase corrente (art. 18, § 1º);
   - o SISMARQ fica com as fases intermediária e permanente (art. 18,
     §§ 2º e 3º) e com a esteira das comissões;
   - o Trâmite da plataforma não substitui o SEI.
5. **Implantação gradual** (art. 62): uma secretaria piloto, depois as
   demais, sempre com a Divisão de Arquivo validando.
6. **Mesmo padrão da plataforma:** cada módulo é um plug-in, com transação
   e auditoria (ADR 016), 100% de cobertura (ADR 009) e um ADR por decisão.

---

## 3. Fases

Cada fase traz o objetivo, a base legal, as entregas e o critério de pronto.
As fases 1 a 4 formam o **mínimo para o sistema oficial**. As fases 5 a 7
completam o projeto do artigo. As datas são definidas com a Prefeitura no
fim da Fase 0, porque dependem das decisões dela.

### Fase 0 — Decisões institucionais (antes de programar o acervo)

**Objetivo:** fechar o que não depende do código e trava as fases seguintes.

| Item | Base | O que decidir ou obter |
|---|---|---|
| 0.1 | Art. 68 e art. 33 ("deste Decreto") | Existe o **decreto regulamentador**? O sistema segue o decreto, se ele existir. |
| 0.2 | Art. 47, parágrafo único | O **Regimento Interno** do Arquivo Público, com o cronograma de recolhimento. |
| 0.3 | Art. 31, §§ 3º, 6º e 7º | **Portarias** da CCPAD e de cada CSPAD: membros, presidentes e secretários. Viram cadastro na Fase 3. |
| 0.4 | Art. 26 | **Atividades-meio.** As TTDDs publicadas não têm as classes do CONARQ (01 a 07). Decidir se o município adota o código de atividades-meio do CONARQ e carregá-lo no Atlas. |
| 0.5 | Art. 31, § 5º | **Cobertura da TTDD.** Quais secretarias ainda não têm TTDD publicada (o prazo legal venceu em 2021). Para elas vale o art. 36. |
| 0.6 | Artigo, item 4.1 | **Política de digitalização:** 300 DPI para texto e 600 DPI para plantas, scanners, quem digitaliza e certificado ICP-Brasil. |
| 0.7 | Artigo, Quadro 1 | **SEI:** a versão em uso, se tem API e quem autoriza a integração (Fase 7). |
| 0.8 | Seção 6 | Resolver as **inconsistências** listadas na seção 6. |
| 0.9 | LGPD | RIPD do acervo (modelo em [`RIPD_MODELO.md`](RIPD_MODELO.md)), com dados pessoais de servidores, pacientes e contribuintes. |
| 0.10 | Infraestrutura | Estimativa de volume do passivo (metros lineares, páginas), disco, MinIO e backup. |

**Pronto quando** cada item tiver resposta registrada, num ADR ou na seção
6, e a secretaria piloto estiver escolhida. A sugestão é a SEMAD
(Administração, Gestão de Pessoas e Inovação): é onde fica o Arquivo
Público (art. 5º) e a TTDD dela está na versão II.

### Fase 1 — Núcleo arquivístico: o acervo e suas travas

**Objetivo:** o plug-in **Acervo**, com o documento de arquivo, a unidade de
arquivamento e a máquina de estados da lei.

**Base:** arts. 18, 19, 20, 24 a 30, 42, 44 e 65. Artigo, itens 4.2 e 4.4.

**Entregas:**
- **Unidade de arquivamento**, cobrindo caixa, pasta, livro, processo e
  documento avulso. Cada uma tem:
  - código de classificação, ligado à série da TTDD do Atlas (art. 25:
    órgão > classe > subclasse > assunto > espécie, que na TTDD aparece
    como órgão.função.subfunção.série);
  - órgão produtor e unidade;
  - datas-limite (ano inicial e final);
  - quantificação e especificação (as colunas dos Anexos I a III);
  - observações.
- **Localização física** em hierarquia: prédio, depósito, corredor,
  estante, prateleira e caixa (artigo, Quadro 1; art. 12, II).
- **Fases** CORRENTE → INTERMEDIÁRIA → PERMANENTE, mais os estados do
  descarte (EM_LISTAGEM, EM_EDITAL, ELIMINADO). As transições valem só no
  domínio:
  - **permanente** não é eliminável, nem depois de digitalizado (arts. 20,
    41 e 42);
  - **bloqueio judicial**: registro do processo judicial, de quem bloqueou
    e de quando. Enquanto houver bloqueio, a eliminação é recusada (arts.
    30, II, 32, II e 44);
  - eliminação só de série que a TTDD prevê, ou pelo rito do art. 35.
- **Cálculo dos prazos** a partir da TTDD (arts. 29 e 30):
  - os prazos em anos contam das datas-limite;
  - os prazos por condição precisam de uma **data do evento**. Na TTDD,
    "vigorando", "enquanto vigorar" e "até a aprovação das contas"
    aparecem em centenas de séries;
  - o campo "contas aprovadas pelo Tribunal de Contas em", do Anexo I, é
    obrigatório para as séries financeiras.
- **Perfis:**
  - Arquivo Público, órgão central (art. 11);
  - agente setorial (art. 9º);
  - membro de CSPAD e membro da CCPAD (Fase 3);
  - consulta.
- **Painel de vencimentos** por secretaria: o que já pode ser transferido,
  recolhido ou avaliado.

**Pronto quando:**
- cada trava tiver teste que prova que nem o administrador a contorna;
- o cálculo de prazo bater com a TTDD numa amostra validada pela Divisão de
  Arquivo;
- a ADR do modelo do acervo estiver aceita.

### Fase 2 — Fases corrente e intermediária: cadastro, transferência e empréstimo

**Objetivo:** pôr o acervo físico existente dentro do sistema e controlar o
dia a dia.

**Base:** arts. 12, II, 29, 49 e 52.

**Entregas:**
- **Importação do passivo** por planilha. O padrão é o mesmo da importação
  do Atlas: simular, conferir e aplicar, com relatório de erros. Permite
  cadastrar caixas já existentes sem digitar uma a uma.
- **Transferência** da unidade produtora (arquivo corrente) para o arquivo
  intermediário do órgão (art. 29, I e II), com guia de transferência e
  aceite de quem recebe.
- **Empréstimo e consulta** (art. 12, II): quem retirou, para quê, prazo de
  devolução e cobrança por notificação.
- **Ponte com o Trâmite:** o processo encerrado entra no acervo já
  classificado. A guarda documental do ADR 026 passa a ser alimentada daqui.
- **Etiquetas** com código de barras ou QR para caixas e unidades de
  arquivamento, apontando para a ficha no sistema.

**Pronto quando** a secretaria piloto tiver seu arquivo intermediário
cadastrado e localizável, e o primeiro ciclo de empréstimo e devolução
tiver sido feito pelo sistema.

**Situação (ADR 033):** importação, guia de transferência com aceite,
empréstimo e etiqueta prontos. A ponte com o Trâmite depende da decisão 0.7
e do item 7 da seção 6.

### Fase 3 — Avaliação e eliminação: a esteira das comissões

**Objetivo:** levar a eliminação do começo ao fim pelo sistema, com os
modelos oficiais. É o que tira o arquivo do "cemitério de papel" (artigo,
item 2.3) e protege o gestor da responsabilidade penal (art. 66).

**Base:** arts. 31 a 40 e 65. Anexos I, II, III e IV.

**Entregas, na ordem do rito:**
1. **Comissões:** CCPAD e uma CSPAD por órgão (art. 31).
   - Composição conferida: formação em Arquivologia ou História, jurídico,
     administração e finanças, área produtora; número ímpar recomendado
     (§§ 6º e 7º).
   - Portaria e vigência registradas.
2. **Listagem de Eliminação (Anexo I)**, gerada a partir das unidades com
   prazo cumprido.
   - Conferência obrigatória de ação judicial (art. 32, II e art. 44).
   - Parecer jurídico anexado (art. 32, § 4º).
   - Quadro de contas aprovadas pelo TCE, quando a série exigir.
   - Assinaturas: responsável pela seleção, presidente da CSPAD e gestor
     do órgão. Depois, o **"Autorizo"** do titular da Divisão de Arquivo
     Público.
3. **Ata da CSPAD (Anexo IV)**: manifestações, encaminhamentos e lista de
   presença, assinada no Signum (artigo, item 4.6).
4. **Amostragem (art. 39):**
   - a comissão define critérios qualitativos e quantitativos (art. 32,
     III);
   - o sistema sorteia e registra a amostra, que passa a ser **permanente**
     e não pode mais ser eliminada;
   - o percentual de referência do artigo é de 5% a 10%, mas o critério é
     da comissão.
5. **Edital de Ciência (Anexo III)**:
   - gerado e publicado no Diário Oficial, com edição, data e página
     registradas;
   - **trava de 30 dias** a partir da publicação (art. 37, § 2º);
   - pedidos de desentranhamento ou de cópia recebidos e decididos no
     período.
6. **Termo de Eliminação (Anexo II)**:
   - só depois do prazo e sem oposição pendente;
   - registra a fragmentação (art. 40) e a doação do papel (parágrafo
     único);
   - envia cópia à Divisão de Arquivo e à SEMAD (art. 38, parágrafo único).
7. **Ritos especiais:**
   - eliminação de documento que não está na TTDD, com autorização da
     CCPAD, do gerente da Divisão e da SEMAD (art. 35);
   - órgão sem TTDD oficializada (art. 36).
8. **Aprovação da TTDD de atividades-fim** pelo sistema:
   - a CSPAD envia, a CCPAD analisa a consistência técnica (não os prazos)
     e devolve com parecer, e a SEMAD aprova para publicação (art. 32, §§
     4º a 8º);
   - ao ser publicada, a nova versão entra no Atlas pela atualização com
     relatório de impacto que já existe (ADR 022).

**Pronto quando:**
- uma eliminação real da secretaria piloto tiver sido feita de ponta a
  ponta pelo sistema;
- os quatro documentos gerados tiverem sido conferidos pela Divisão de
  Arquivo contra os Anexos;
- a auditoria mostrar a cadeia completa da listagem ao termo.

### Fase 4 — Recolhimento e guarda permanente

**Objetivo:** levar o permanente ao Arquivo Público com o que a lei exige e
manter a guarda.

**Base:** arts. 5º, VII, 16, 17, 41, 42, 47, 50, 51, 53 e 54.

**Entregas:**
- **Recolhimento** com checklist do art. 53: classificado, avaliado,
  organizado, higienizado, acondicionado e com instrumento descritivo.
  - Segue o cronograma do Regimento Interno (art. 47).
  - Os custos de preparo e transporte ficam com o órgão (art. 53, § 2º).
- **Registro de entrada** no Arquivo Público, com **cópia enviada à
  unidade de origem** (art. 5º, VII).
- **Instrumento descritivo** gerado a partir do acervo recolhido.
- **Recolhimentos especiais:** extinção de órgão (arts. 17 e 50) e
  desestatização (arts. 16 e 51).
- **Arquivos privados de interesse público e social** (arts. 55 a 58):
  - cadastro;
  - parecer da Comissão Especial e decreto de declaração;
  - doação ou depósito revogável;
  - **direito de preferência**, com prazo de 60 dias na alienação.
- **Normas do Arquivo:** instruções normativas (art. 54) publicadas e
  consultáveis no sistema.

**Pronto quando** um recolhimento da secretaria piloto entrar pelo sistema
e o órgão de origem receber o registro de entrada.

### Fase 5 — Digitalização com valor probatório e preservação digital

**Objetivo:** o pipeline do artigo, que transforma papel em representante
digital com valor de prova e o mantém íntegro.

**Base:** art. 48. Decreto Federal 10.278/2020. Artigo, itens 4.1 e 4.3 e
Quadro 2 (módulo SISMARQ_PIPELINE).

**Entregas:**
- **Entrada:**
  - pasta monitorada (SMB/FTP) e envio pela tela;
  - cada lote ligado à unidade de arquivamento da Fase 1.
- **Processamento no worker:**
  - `ocrmypdf` com alinhamento (deskew), limpeza e OCR em português;
  - saída **PDF/A** (ISO 19005);
  - **SHA-256** do arquivo;
  - extração do texto para a busca;
  - conferência da **resolução**: 300 DPI para texto e 600 DPI para plantas.
- **Metadados do Decreto 10.278** gravados no XMP do PDF e no banco:
  - identificador único;
  - responsável pela digitalização, com o CPF mascarado nos logs;
  - data e hora em UTC;
  - equipamento;
  - DPI;
  - hash;
  - código de classificação.
- **Assinatura ICP-Brasil (PAdES)** do representante digital. O Signum hoje
  é só eletrônico: integrar certificado ICP-Brasil, conforme a decisão 0.6.
- **Guarda e integridade:**
  - armazenamento no MinIO com object lock, a mesma base da cópia WORM da
    auditoria;
  - **verificação periódica de integridade (fixity)**, com alerta quando
    o hash não bate.
- **Visualizador em fluxo:**
  - HTTP 206 (Range) com pdf.js, que renderiza só as páginas visíveis;
  - sem limite de tamanho;
  - meta do artigo: abrir em menos de 300 ms.
- **Art. 42 no domínio:** digitalizar um permanente nunca libera a
  eliminação do papel.

**Pronto quando:**
- um lote real for digitalizado com os metadados conferidos;
- a verificação de integridade detectar um arquivo adulterado num teste;
- um processo de 600 páginas abrir no visualizador sem baixar o arquivo
  inteiro.

### Fase 6 — Acesso, busca e difusão

**Objetivo:** cumprir o direito de acesso com as restrições da lei.

**Base:** arts. 2º, 5º, X e XI, 7º, I, 13 e 56, § 2º.

**Entregas:**
- **Busca do acervo no Typesense:**
  - coleção própria com facetas por órgão, código, classe, espécie, fase,
    destinação e ano (artigo, item 4.4);
  - texto do OCR;
  - meta do artigo: menos de 50 ms.
  - **Decisão de segurança antes (ADR):** o artigo propõe consultar o
    Typesense direto do navegador. A plataforma só indexa conteúdo público
    (ADR 031). O acervo tem sigilo e dados pessoais, então é preciso
    escolher entre uma coleção pública e outra interna, ou chaves de busca
    com filtro embutido por perfil.
- **Consulta pública** sem login, com o que é ostensivo: descrição e,
  quando liberado, o documento digital. Pedido de acesso ao restrito pelo
  SIC (série 2.0.03.02 da TTDD da SEMAD).
- **Nível de acesso por unidade e documento:** público, restrito ou
  sigiloso, com hipótese legal. O ostensivo é a regra (art. 2º).
- **Difusão** (art. 5º, X e XI): página do Arquivo Público, acervo em
  destaque e exposições virtuais. Prioridade baixa.

**Pronto quando** o cidadão encontrar, sem login, um documento permanente
digitalizado, e um documento restrito não aparecer para quem não tem
acesso.

### Fase 7 — Integração com o SEI e os demais sistemas

**Objetivo:** o recolhimento automático do que nasce digital no SEI e a
TTDD aplicada em todos os sistemas.

**Base:** arts. 5º, VII, 45 e 53. Artigo, item 3 ("barramento REST com
pacote AIP").

**Entregas:**
- **Recolhimento do SEI:** processos concluídos chegam ao SISMARQ como
  pacote (documentos, metadados, hashes e trilha), conferidos na entrada.
  Depende da decisão 0.7.
- **A TTDD disponível para os sistemas de protocolo** (art. 45): a API
  pública do Atlas já serve a TTDD. Falta documentar o contrato e
  combinar com quem mantém o SEI.
- **Egress** (webhooks) para avisar os sistemas da Prefeitura sobre
  eliminações e recolhimentos.

**Pronto quando** um processo concluído no SEI chegar ao acervo sem
redigitação, com o hash conferido.

### Fase 8 — Implantação oficial

**Objetivo:** transformar o sistema em ferramenta oficial da Prefeitura.

**Base:** arts. 5º, VIII, 11, VIII, 13, 54, 61, 62 e 64.

**Entregas:**
- **Conformidade com o e-ARQ Brasil** (art. 13): checklist de requisitos
  avaliado com a Divisão de Arquivo, com as lacunas registradas.
- **Ato normativo** do Secretário de Administração (art. 64) e
  **instrução normativa** do Arquivo Público (art. 54) adotando o sistema
  e os procedimentos.
- **Capacitação** dos agentes setoriais e das comissões (arts. 5º, VIII e
  11, VIII), com manuais e vídeos curtos.
- **Implantação por secretaria** (art. 62), na ordem: piloto, depois as
  secretarias com TTDD publicada, depois as demais, conforme o art. 36.
- **Indicadores para a SEMAD** (art. 38, parágrafo único):
  - metros lineares eliminados, recolhidos e digitalizados;
  - prazos vencidos por secretaria;
  - consultas atendidas.
- **Operação:** backup e restauração testados, monitoramento, plano de
  continuidade e VPAT atualizado.

**Pronto quando** o ato oficial estiver publicado e todas as secretarias
com TTDD estiverem operando no sistema.

---

## 4. Ordem e dependências

```
Fase 0 ──► Fase 1 ──► Fase 2 ──► Fase 3 ──► Fase 4 ──► Fase 8
              │                     ▲           ▲
              └──► Fase 5 ──────────┴──► Fase 6 ┘
                                 Fase 7 (depende da 0.7; em paralelo a partir da Fase 4)
```

- A **Fase 1** é a base de tudo: o acervo, as travas e os prazos.
- **Mínimo para o sistema oficial:** Fases 0 a 4 e a Fase 8, mesmo sem
  digitalização.
- A **Fase 5** pode começar junto com a Fase 2, se houver equipe. Ela só
  precisa da unidade de arquivamento da Fase 1.
- A **Fase 6** precisa da Fase 5 para buscar no texto do OCR. A busca por
  metadados sai antes.

## 5. Matriz de rastreabilidade

| Exigência | Base | Fase | Hoje |
|---|---|---|---|
| Acesso pleno, com sigilo como exceção | Arts. 2º e 7º, I | 6 | Parcial (site público) |
| Órgão central e órgãos setoriais | Arts. 8º a 12 | 1 | **Pronto** (escopo do Acervo e `acervo:central` — ADR 032) |
| Sistema informatizado que atenda ao e-ARQ Brasil | Art. 13 | 8 | Falta |
| Três idades (corrente, intermediária, permanente) | Art. 18 | 1 | **Pronto** (ADR 032) |
| Eliminação só prevista na TTDD e autorizada | Arts. 19 e 65 | 1 e 3 | Parcial (travas prontas — ADR 032; rito na Fase 3) |
| Permanente inalienável | Arts. 20 e 42 | 1 | **Pronto** (domínio e banco — ADR 032) |
| Classificação duplex com código de 5 níveis | Arts. 24 e 25 | 1 | **Pronto** (unidade classificada pela TTDD — ADR 032) |
| Classes das atividades-meio do CONARQ | Art. 26 | 0 e 1 | Falta |
| TTDD com prazos e destinação | Arts. 27 a 30 | — | **Pronto** (Atlas) |
| Comissões CCPAD e CSPAD | Arts. 31 e 32 | 3 | Falta |
| Aprovação da TTDD de atividades-fim | Art. 32, §§ 4º a 8º | 3 | Parcial (carga com impacto) |
| Relação, Ata, Edital e Termo | Arts. 33, 34, 37 e 38 | 3 | Falta |
| Trava de 30 dias do edital | Art. 37, § 2º | 3 | Falta |
| Amostragem | Art. 39 | 3 | Falta |
| Fragmentação e doação | Art. 40 | 3 | Falta |
| Bloqueio por ação judicial | Arts. 32, II e 44 | 1 | **Pronto** (ADR 032) |
| TTDD nos sistemas informatizados | Art. 45 | 7 | Parcial (API pública) |
| Documentos digitais e híbridos | Art. 48 | 5 | Falta |
| Recolhimento com registro de entrada | Arts. 5º, VII e 53 | 4 | Falta |
| Arquivos privados de interesse público | Arts. 55 a 58 | 4 | Falta |
| Capacitação | Arts. 5º, VIII e 11, VIII | 8 | Falta |
| Digitalização com PDF/A, metadados e hash | Decreto 10.278; artigo 4.1 | 5 | Falta |
| Assinatura ICP-Brasil (PAdES) | Artigo, Quadro 3 | 5 | Falta (Signum é eletrônico) |
| Auditoria imutável | Artigo, Quadro 2 | — | **Pronto** |
| Busca facetada (Typesense) | Artigo, item 4.4 | 6 | Parcial (infraestrutura, ADR 031) |
| Visualizador em fluxo (HTTP 206) | Artigo, item 4.5 | 5 | Falta |

## 6. Inconsistências e pontos a confirmar nos documentos

Resolver na Fase 0. Até lá, o sistema adota o que está entre parênteses.

1. **Onde se publica o edital.** O art. 37 manda publicar no Diário Oficial
   **do Município**. O modelo do Anexo III fala em Diário Oficial **do
   Estado**. (Adotado: DIORONDON-e, conforme o artigo da lei e o artigo
   científico.)
2. **Nome do órgão.** A lei fala em "Divisão de Arquivo Público" da
   Secretaria Municipal de Administração. A TTDD de 2025 fala em
   "Superintendência de Arquivo Público" da Secretaria Municipal de
   Administração, Gestão de Pessoas e Inovação. (Adotado: o nome atual em
   cadastro configurável, que aparece nos modelos gerados.)
3. **Lei ou decreto.** O art. 33 chama o Anexo IV de "parte integrante
   deste Decreto", e o art. 68 prevê um decreto regulamentador. É preciso
   confirmar se o decreto existe e se altera algum modelo.
4. **Atividades-meio.** As 13 TTDDs publicadas cobrem atividades-fim por
   secretaria, e nenhuma traz as classes 01 a 07 do art. 26.
5. **Prazos por condição.** Centenas de séries usam "vigorando" ou
   "enquanto vigorar". O sistema precisa de uma data de evento informada
   por quem cadastra, e a Divisão de Arquivo deve dizer qual evento vale
   para cada tipo de série.
6. **Typesense no navegador.** O artigo propõe que o navegador consulte o
   Typesense diretamente. Isso conflita com o sigilo do acervo, e a
   solução fica para a decisão da Fase 6. O artigo também usa `tenant_id`,
   mas a plataforma atende um município por implantação (ADR 015), então o
   campo não entra.
7. **Trâmite e SEI.** O artigo põe a fase corrente no SEI. O Trâmite da
   plataforma deve ficar como ponte para o que não passa pelo SEI, sem
   competir com ele. A Prefeitura precisa confirmar.

## 7. Como manter este roadmap

- Ao fechar uma entrega, atualize a coluna "Hoje" da seção 5 no mesmo
  commit e cite o ADR.
- Toda decisão das seções 3 (Fase 0) e 6 vira ADR em [`adr/`](adr/README.md).
- Mudou a lei, um decreto ou a TTDD? Primeiro coloque o documento novo em
  `docs/SISMARQ/`, depois revise este roadmap.
