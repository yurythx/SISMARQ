# Fase 0 — Registro das decisões institucionais

Quadro fechado em 16/10/2026 com as respostas da Prefeitura. **Não há**
decreto, portaria, certificado nem URL do SEI inventados. O que ainda não
existe opera com o fallback da lei (municipal + federal) e do código.
Detalhe normativo: [ADR 043](../adr/043-fase0-decisoes-institucionais.md).
Domínio e espelho federal (para implementar depois): [DOMINIO.md](DOMINIO.md),
[EARQ_RDC.md](EARQ_RDC.md), [ADR 045](../adr/045-linha-de-base-arquivistica.md).

## A. Quem e onde (IAM / Keycloak + LDAP)

Logins não entram no git: vêm do Keycloak (federação LDAP). O IAM concede
as permissões nos **perfis**; o mapeamento grupo AD → perfil + escopo é a
tela Configurações → Mapeamento AD.

| Papel | Permissão | Escopo | Quem |
|---|---|---|---|
| A1. Secretaria piloto | — | SEMAD | Primeira Relação de Eliminação: unidade produtora da SEMAD. |
| A2. Quem recebe o acervo / autoriza eliminação ordinária, locais, comissões, recolhimento | `acervo:central` | **global** | Superintendência de Arquivo Público (o próprio Arquivo). |
| A3. Arts. 35 e 36 + cópia do Termo (art. 38) | `acervo:semad` | **global** | Perfil da SEMAD (Administração), distinto do Arquivo. Sem esse grupo, `termo_copia_semad` fica 0. |
| A4. CSPAD do órgão piloto (Relação, ata, edital, amostragem, termo) | `acervo:avaliacao` | **órgão** (SEMAD) | Membros nomeados na portaria, quando existir. Permissão no perfil; composição em `/acervo/comissoes` (usuários). Sem portaria não se inventa membro. |

## B. Itens 0.1 a 0.10

| Item | Tema | Situação | Resposta |
|---|---|---|---|
| 0.1 | Decreto regulamentador (arts. 68 e 33) | **REGISTRADO — não existe** | Opera-se a **Lei Municipal 10.301/2019** e a legislação **federal** (Lei 8.159/1991; Decreto 10.278/2020; CONARQ). Não há decreto municipal que altere os Anexos I–IV. O art. 33 diz “deste Decreto” porque o modelo foi escrito como se um decreto republicasse os anexos; esse decreto não foi publicado. Os Anexos gerados são os da lei. |
| 0.2 | Regimento Interno do Arquivo (art. 47) | **REGISTRADO — em elaboração** | Ainda não existe. A tela de normas (ADR 036) já recebe o regimento e o cronograma de recolhimento quando forem publicados. |
| 0.3 | Portarias das comissões (art. 31) | **REGISTRADO — posterior** | CCPAD e CSPAD do piloto: cadastrar em `/acervo/comissoes` só com a portaria real (presidente, secretário, áreas do art. 31). Permissão `acervo:avaliacao` no perfil; membros = usuários. |
| 0.4 | Atividades-meio CONARQ (art. 26) | **REGISTRADO — adota** | O município **adota** as classes **01 a 07** do Código de Classificação e da TTDD de atividades-meio da Administração Pública (CONARQ), porque o **art. 26** da Lei 10.301/2019 manda usar esse código para o que não é atividade-fim. As TTDDs já no Atlas são as de **atividades-fim** publicadas no DIORONDON-e. A carga 01–07 entra por “Atualizar TTDD”, sem códigos inventados. O Atlas explica isso na consulta. |
| 0.5 | Órgãos sem TTDD de atividades-fim (art. 36) | **REGISTRADO — regra, sem lista inventada** | `rito=ART_36` só se o órgão **não** tiver TTDD de atividades-fim no Atlas. Com TTDD publicada: `ORDINARIO` ou `ART_35`. Órgãos **com** TTDD carregada hoje (`deploy/ttdd/ttdd.json`): Administração (SEMAD, 2.0); Fazenda (3.0); Transparência e Controle Interno (4.0); Habitação e Urbanismo (6.0); Infraestrutura (7.0); Meio Ambiente, Agricultura e Pecuária (8.0); Desenvolvimento Econômico (9.0); Promoção e Assistência Social (10.0); Saúde (12.0). Qualquer outro órgão do organograma, até publicar, usa art. 36. |
| 0.6 | Política de digitalização e ICP-Brasil | **REGISTRADO — preparado, sem certificado** | Quem digitaliza: **Arquivo Municipal**. Pasta SMB/FTP: ainda a decidir (não ligada). Certificado A1/A3/HSM: **não há**. Padrão a seguir, quando houver certificado: Decreto Federal 10.278/2020 (PAdES). Até lá `pades=AUSENTE`, `ata_signum=AUSENTE`, `termo_signum=AUSENTE`. Signum eletrônico atual **não** substitui ICP-Brasil. Porta `AssinadorICP` no domínio, **não composta** (ADR 040). |
| 0.7 | SEI | **REGISTRADO — envelope pronto, sem integração** | Sem versão, API, URL de homologação nem mapeamento órgão SEI → unidade. O Trâmite já ingressa sozinho (ADR 039). `UnidadeDoSEI` valida o envelope e recusa (ADR 042). |
| 0.8 | Inconsistências (seção 6 do roadmap) | **REGISTRADO** | (1) Edital: **DIORONDON-e** (Diário do Município), não o “Diário do Estado” do Anexo III. (2) Nome nos PDFs: **Superintendência de Arquivo Público**. (3) Ver 0.1 — “deste Decreto” não cria decreto. (5) “Enquanto vigorar” / “vigorando”: data do **evento** informada por quem cadastra, conforme o tipo da série (contrato, gestão, vínculo, fato da condição). TTDD pode mudar no tempo: série revogada não apaga a unidade (ADR 019). (6)(7) já no código. |
| 0.9 | RIPD do acervo | **REGISTRADO — DPO não começou** | Preencher [`../RIPD_MODELO.md`](../RIPD_MODELO.md) quando o encarregado iniciar. |
| 0.10 | Volume e infra | **REGISTRADO — qualitativo** | Sem metros lineares, páginas nem backup medido. Expectativa: **grande volume** (cidade com cerca de 70 anos e ~300 mil habitantes). Na VM: `docker system df` e o volume MinIO; backup já existe (`scripts/backup.sh`, Postgres + buckets). **Não** inventar TB nem política até a Superintendência medir o acervo físico. |

## C. Homologação

Continua reunião com a Superintendência (tela, não chat): [HOMOLOGACAO.md](HOMOLOGACAO.md). Sem a primeira linha da ata de conferência o sistema não é o oficial (art. 62). IAM na operação: [IAM.md](IAM.md).
