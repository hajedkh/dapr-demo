output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

resource "local_file" "kubeconfig" {
  depends_on   = [azurerm_kubernetes_cluster.cluster]
  filename     = "kubeconfig"
  content      = azurerm_kubernetes_cluster.cluster.kube_config_raw
}

output "acr_login_server" {
  description = "The URL of the Azure Container Registry"
  value       = azurerm_container_registry.acr.login_server
}

output "user_assigned_id" {
  value       = azurerm_container_registry.acr.login_server
}
output "workload_identity_client_id" {
  value = azurerm_user_assigned_identity.dapr-identity.client_id
}

output "servicebus_connection_string" {
  sensitive = true
  value = azurerm_servicebus_namespace.sb-tvt.default_primary_connection_string
  description = "The connection string for the Azure Service Bus Namespace."
}

output "key_vault_uri" {
  value = azurerm_key_vault.kv.vault_uri
  description = "The URI for the Azure Key Vault"
}