# Guia de Automação de Infraestrutura: Proxmox VE + Terraform + Ansible + Cloudflare Tunnel

Este guia documenta o padrão corporativo de **Infraestrutura como Código (IaC)**, **Gerenciamento de Configuração** e **Borda Zero Trust** do SISMARQ no Proxmox VE 8/9.

---

## 1. Visão Geral da Arquitetura (Single VM + Cloudflare Zero Trust)

A arquitetura elimina proxies reversos locais expostos (como containers LXC de borda ou portas abertas no roteador). A entrega de tráfego ocorre via **Cloudflare Tunnel (conexão outbound criptografada)** diretamente para a VM de aplicação:

```
      INTERNET / USUÁRIOS
               │
      [ Cloudflare Edge ] (WAF, Proteção DDoS, SSL/TLS, Políticas Zero Trust)
               │
        (Túnel Seguro Outbound / Sem Portas Abertas no Firewall)
               │
   Proxmox VE Host (pve1)
               └── VM 201: nexus-prod-app (KVM Dedicada - IP 192.168.1.46)
                    ├── Cloud-Init + QEMU Guest Agent
                    ├── Docker Engine CE + Compose Plugin v2
                    │    ├── cloudflared (conector de túnel)
                    │    ├── frontend (Next.js 16 - :3010)
                    │    ├── backend-api (Go Microkernel - :8010)
                    │    └── postgres, rabbitmq, redis, minio
```

### Por que esta arquitetura é superior?

1. **Zero Portas Abertas (Zero Trust Ingress):** Não há NAT, DMZ nem Port Forwarding (80/443) no roteador corporativo. Todo tráfego vem da Cloudflare através de túnel seguro.
2. **Eliminação do LXC 200 (Edge Proxy):** Sem proxy intermediário consumindo memória e CPU desnecessariamente. A borda da Cloudflare já faz terminação TLS, cache estático e proteção contra ataques.
3. **Isolamento de Kernel & Estabilidade de Storage:** A VM KVM 201 opera com kernel próprio e sistema de arquivos virtualizado `ext4`, garantindo que o driver `overlay2` do Docker funcione com desempenho nativo, sem conflitos com o pool ZFS do Proxmox host.

---

## 2. Estrutura de Arquivos

```
deploy/infra/
├── terraform/
│   ├── versions.tf                  # Provider bpg/proxmox
│   ├── provider.tf                  # Autenticação via API Token
│   ├── variables.tf                 # Declaração de variáveis e template ID
│   ├── terraform.tfvars             # Valores locais (protegido por .gitignore)
│   ├── terraform.tfvars.example     # Exemplo versionado no Git
│   ├── main.tf                      # Declaração da VM 201 + gerador de inventário
│   ├── outputs.tf                   # IPs, IDs e inventário gerado
│   ├── create-cloudinit-template.sh # Script para criar template Cloud-Init no Proxmox
│   └── modules/
│       ├── proxmox_vm/              # Módulo para criar VMs KVM com Cloud-Init
│       └── proxmox_lxc/             # Módulo legado (mantido para serviços auxiliares)
│
└── ansible/
    ├── ansible.cfg                  # Otimizações de SSH e pipelining
    ├── inventory.ini                # Gerado dinamicamente pelo Terraform
    ├── site.yml                     # Playbook orquestrador mestre
    └── roles/
        ├── common/                  # Hardening, timezone, UFW, QEMU guest agent
        ├── docker/                  # Instalação oficial do Docker Engine + Compose
        └── nexus_deploy/            # Deploy automatizado da stack via deploy.sh
```

---

## 3. Pré-Requisitos no Proxmox (`pve1`)

### A. Permissão do Token de API
O usuário de automação e o token de API:
* **Usuário:** `terraform-prov@pve`
* **Token ID:** `tf-token`
* **Full Token ID:** `terraform-prov@pve!tf-token`

Garanta a permissão executando no shell do Proxmox:
```bash
pveum acl modify / -user terraform-prov@pve -role Administrator
```

### B. Template Cloud-Init (VM 9000)
Para que o Terraform possa clonar a VM KVM com Cloud-Init instantaneamente, você precisa de um template base no Proxmox (ID 9000).

Execute o script auxiliar diretamente no terminal do Proxmox:
```bash
# No shell do Proxmox Host:
bash deploy/infra/terraform/create-cloudinit-template.sh 9000 debian-12-cloudinit zfs
```
Em menos de 1 minuto o template `debian-12-cloudinit` (VM 9000) estará pronto para ser clonado pelo Terraform.

---

## 4. Como Executar via Makefile ou Linha de Comando

Você pode executar diretamente da raiz do projeto usando os atalhos do `Makefile`:

### Passo 1: Inicializar e Planejar
```bash
make infra-init
make infra-plan
```

### Passo 2: Provisionar a Infraestrutura
```bash
make infra-apply
```
O Terraform irá:
1. Clonar o template e criar a VM KVM 201 (`nexus-prod-app`) com 4 vCPUs, 8GB RAM, 60GB disco e IP `192.168.1.46`.
2. Injetar sua chave SSH pública via Cloud-Init.
3. Gerar automaticamente o arquivo `deploy/infra/ansible/inventory.ini`.

### Passo 3: Configurar o Docker e Fazer Deploy do SISMARQ
```bash
make infra-deploy
```
O Ansible irá:
1. Configurar hardening básico, firewall UFW e iniciar o `qemu-guest-agent`.
2. Instalar o Docker Engine oficial e o Compose plugin na VM 201.
3. Executar o script `scripts/deploy.sh`, que gera `.env` (segredos criptográficos aleatórios fortes), chaves RSA locais e sobe todos os containers (`backend-api`, `frontend`, `postgres`, `rabbitmq`, `redis`, `minio`, `migrate`).

---

## 5. Conectando ao Cloudflare Tunnel

Após o deploy da VM, configure o roteamento no **Cloudflare Zero Trust** (em *Access > Tunnels*):

* **Frontend:** `nexus.seudominio.com` ➔ `http://192.168.1.46:3010`
* **API / Backend:** `api.seudominio.com` ➔ `http://192.168.1.46:8010`

Caso queira que o conector `cloudflared` suba junto com a stack dentro do Docker:
1. Adicione `CLOUDFLARE_TUNNEL_TOKEN=seu_token` no arquivo `.env`.
2. Acrescente `docker-compose.cloudflare.yml` à variável `COMPOSE_FILE` no `.env`.

---

## 6. Manutenção e Comandos Úteis

### Recriar a Máquina do Zero (Disposable Infra)
Se precisar resetar completamente o host de aplicação:
```bash
cd deploy/infra/terraform
terraform destroy -target=module.vm_nexus_app
terraform apply
make infra-deploy
```

### Testar Conexão Ansible
```bash
cd deploy/infra/ansible
ansible all -m ping
```
