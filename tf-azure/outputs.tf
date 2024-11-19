resource "local_file" "kubeconfig" {
  depends_on   = [azurerm_kubernetes_cluster.cluster]
  filename     = "kubeconfig"
  content      = azurerm_kubernetes_cluster.cluster.kube_config_raw
}

output "workload_identity_client_id" {
  value = azurerm_user_assigned_identity.dapr-identity.client_id
}

output "kubernetes_resource" {
  value = templatefile("${path.module}/sa.yaml", {
    workload_identity = azurerm_user_assigned_identity.dapr-identity.client_id
  })
}