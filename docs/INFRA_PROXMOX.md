# SISMARQ no Proxmox: Terraform + Ansible

Hypervisor: **`192.168.1.13:8006`**. Esse IP é a API do Proxmox, **não** a URL do SISMARQ. A aplicação sobe numa **VM KVM nova** (VMID **220**, `sismarq-prod-app`, `192.168.1.52`, **4 GB**), clonada do template **9000** (`ubuntu-2404-cloud`). Não reutiliza a 201 (`nexus-prod-app`) nem a 210 (`projeto-assistencia`).

SSO: **`https://sso.rondonopolis.mt.gov.br/realms/rondonopolis`** (client `sismarq`). O Keycloak de teste (`docker-compose.keycloak.yml`) **não** sobe neste fluxo.

```
Internet / rede municipal
        │
sso.rondonopolis.mt.gov.br     (Keycloak já existente)
        │
Proxmox 192.168.1.13
        └── VM 220 sismarq-prod-app  (192.168.1.52)
             ├── Docker: frontend :3010, API :8010, postgres, redis,
             │           rabbitmq, minio, typesense, worker
             └── .env com KEYCLOAK_ISSUER_URL = issuer público HTTPS
```

## 1. No Proxmox (`pve1`, 192.168.1.13)

Token de API:

```bash
pveum user add terraform-prov@pve
pveum acl modify / -user terraform-prov@pve -role Administrator
pveum user token add terraform-prov@pve tf-token
```

Template Cloud-Init: a VM **9000** (`ubuntu-2404-cloud`) **já existe**. Não rode `create-cloudinit-template.sh` (e não precisa clonar o repo no Proxmox só por isso). Confira IPs das VMs atuais para não colidir:

```bash
qm config 201 | grep ipconfig
qm config 210 | grep ipconfig
```

Usuário/token `terraform-prov@pve!tf-token` **já existem**. Não rode `pveum user add` de novo. O secret do token só aparece na criação: se você não tiver salvo, recrie:

```bash
pveum user token delete terraform-prov@pve tf-token
pveum user token add terraform-prov@pve tf-token
```

Guarde o secret no `terraform.tfvars` da máquina que roda o Terraform (não no Proxmox).

O padrão do Terraform é `192.168.1.52`. Se esse IP já aparecer no `ipconfig` acima, mude `app_ipv4` no `terraform.tfvars`.

## 2. Na máquina que opera o IaC (com Terraform, Ansible e SSH)

```bash
cp deploy/infra/terraform/terraform.tfvars.example deploy/infra/terraform/terraform.tfvars
# Preencha token, senha do console, chave SSH pública e, se preciso, app_ipv4 / node_name / storage.

cp deploy/infra/ansible/secrets.yml.example deploy/infra/ansible/secrets.yml
# Cole o Client Secret do client sismarq (o mesmo da tela de Configurações).

make infra-init
make infra-plan
make infra-apply
make infra-deploy
```

O Ansible instala Docker, clona o repositório, roda `scripts/deploy.sh` e grava o SSO oficial no `.env`. Sem `secrets.yml` a stack sobe, mas o botão Keycloak não autentica.

No client `sismarq` do Keycloak, Valid redirect URI:

`http://192.168.1.52:3010/api/auth/callback/keycloak`

(ajuste o IP se `app_ipv4` for outro; com HTTPS local, use `https://…`.)

## 3. Depois do deploy

| URL | O quê |
|---|---|
| `http://192.168.1.52:3010` | SISMARQ (navegador) |
| `http://192.168.1.52:8010/health` | API |
| Tela Configurações → Keycloak | Deve mostrar issuer `https://sso.rondonopolis.mt.gov.br/realms/rondonopolis` (já gravado no `.env`; pode salvar de novo pela tela) |

Login local (`make prod-seed-admin` na VM) funciona mesmo se o SSO falhar.
Mapeamento AD: [IAM.md](acervo/IAM.md) — só com o CN real do grupo; não
invente membros sem portaria.

Atualizar código na VM **sem** `git pull` do GitHub `main`:
`deploy/homolog/sincronizar-vm.sh` (ou `.ps1` no Windows). A sync recarrega
a TTDD e só **insere** procedimentos com código ainda inexistente; para
mudar fluxo já publicado ou aplicar nova TTDD pela gestão, ver
[atlas/ATUALIZAR.md](atlas/ATUALIZAR.md).

Volume (item 0.10 da Fase 0): ainda **sem** metros lineares. Na VM, conferir disco
sem inventar número:

```bash
docker system df
docker volume ls
df -h
```

Backup da stack: `scripts/backup.sh` (dump Postgres + espelho MinIO). Copiar
`BACKUP_DIR` para fora da máquina — o backup no mesmo disco não cobre perda da VM.

Não ligue `docker-compose.keycloak.yml` nesta VM.

## 4. Cloudflare Tunnel (opcional)

Só se houver hostname público. Sem túnel, o acesso é pela LAN (`192.168.1.52`). Token em `CLOUDFLARE_TUNNEL_TOKEN` e `COMPOSE_FILE=docker-compose.yml:docker-compose.cloudflare.yml`.

## 5. State antigo (`vm_nexus_app`)

Se já existir state Terraform com o módulo `vm_nexus_app`, o `apply` vai **substituir** o recurso pelo `vm_sismarq_app`. Confira o plan antes. Para destruir só a VM antiga: `terraform destroy -target=module.vm_nexus_app` no state antigo.
