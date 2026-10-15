# ==============================================================================
# INFRAESTRUTURA PROXMOX VE - PROJETO NEXUS & SERVIÇOS
# ==============================================================================
# Arquitetura Single-VM Enterprise:
# VM KVM Dedicada (ID 201) para Docker Engine com isolamento de hardware.
# Borda e Ingress seguros via Cloudflare Tunnel (Zero Trust - sem portas expostas).
# ==============================================================================

# 1. Máquina Virtual KVM Dedicada: Servidor Docker Principal do Projeto Nexus
module "vm_nexus_app" {
  source = "./modules/proxmox_vm"

  node_name         = var.proxmox_node_name
  vm_id             = 201
  name              = "nexus-prod-app"
  description       = "VM KVM Dedicada Docker Host - Projeto Nexus (Ingress Cloudflare Tunnel)"
  tags              = ["prod", "docker", "nexus", "kvm"]
  cores             = 4
  memory            = 4096
  disk_storage      = var.storage_disk_pool
  disk_size         = 60
  network_bridge    = var.network_bridge
  ipv4_address      = "192.168.1.46/24"
  ipv4_gateway      = var.default_gateway
  admin_user        = "root"
  admin_password    = var.vm_admin_password
  ssh_public_keys   = compact([
    var.admin_ssh_public_key,
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPkIpwYsgE/WHWJEd73U93jWVxfvSSjQlin3D9D8F/VW ansible-deploy@projeto-assistencia"
  ])
  clone_template_id = var.vm_template_id
  start_on_boot     = true
}

# 2. Geração Automática do Inventário do Ansible
resource "local_file" "ansible_inventory" {
  filename        = "${path.module}/../ansible/inventory.ini"
  file_permission = "0644"
  content         = <<-EOT
    # ==============================================================================
    # INVENTÁRIO GERADO AUTOMATICAMENTE VIA TERRAFORM - PROJETO NEXUS
    # ==============================================================================

    [all:vars]
    ansible_user=root
    ansible_python_interpreter=/usr/bin/python3
    ansible_ssh_common_args='-o StrictHostKeyChecking=no'

    # --- Hosts com Docker Engine (VM KVM Dedicada de Produção) ---
    [docker_hosts]
    nexus-prod-app ansible_host=192.168.1.46 proxmox_vmid=${module.vm_nexus_app.id}

    # --- Grupos de Aplicação ---
    [nexus_apps]
    nexus-prod-app
  EOT
}
