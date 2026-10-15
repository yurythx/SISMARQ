output "nexus_app_vm_id" {
  description = "ID da VM KVM de Produção do Nexus"
  value       = module.vm_nexus_app.id
}

output "nexus_app_vm_ip" {
  description = "IP da VM KVM de Produção do Nexus"
  value       = module.vm_nexus_app.ipv4_address
}

output "ansible_inventory_path" {
  description = "Caminho do inventário gerado para o Ansible"
  value       = local_file.ansible_inventory.filename
}
