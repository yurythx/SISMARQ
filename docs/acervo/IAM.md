# IAM do Acervo (Keycloak / LDAP)

Logins **não** entram no git. A tela **Configurações → Mapeamento AD**
liga o **CN do grupo** que já existe no AD ao perfil + escopo. Sem o CN
real, não se cadastra mapeamento de exemplo (LGPD; ADR 043).

## Redirect do client `sismarq`

Realm: `rondonopolis`. Issuer:
`https://sso.rondonopolis.mt.gov.br/realms/rondonopolis`.

Valid redirect URI (homologação na VM 220):

`http://192.168.1.52:3010/api/auth/callback/keycloak`

Web origins: `http://192.168.1.52:3010`

Se o Client Secret tiver aparecido em log, `.env` compartilhado ou
chat: **gere outro no Keycloak**, grave só em `deploy/infra/ansible/secrets.yml`
(gitignore) e no `.env` da VM. Não commitar.

## Mapeamentos a criar (quando o AD informar o CN)

| Perfil no SISMARQ | Permissão | Escopo | Quem (papel, não lista nominal) |
|---|---|---|---|
| Arquivo Público | `acervo:central` | global | Superintendência de Arquivo Público |
| SEMAD (arts. 35/36 e cópia do Termo) | `acervo:semad` | global | Administração, distinto do Arquivo |
| CSPAD do piloto | `acervo:avaliacao` | órgão SEMAD | Só com portaria (art. 31) |

Composição da comissão: `/acervo/comissoes` com usuários nomeados na
portaria. Sem portaria, a permissão no perfil pode existir; membro
inventado, não.

## Login de emergência

`make prod-seed-admin` na VM cria/redefine o admin **local**. Não
substitui o SSO. Não ligue `docker-compose.keycloak.yml` nesta VM.
