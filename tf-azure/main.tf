data "azurerm_client_config" "client" {}

resource "azurerm_resource_group" "rg" {
  location = var.resource_group_location
  name     = var.resource_group_name_prefix
}

// AKS Cluster
resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "dapr-demo-cluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "dapr-demo-cluster"
  oidc_issuer_enabled = true
  workload_identity_enabled = true

  default_node_pool {
    name       = "default"
    node_count = "1"
    vm_size    = "standard_d2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
  }
}

// Configuration for Managed Identity
resource "azurerm_user_assigned_identity" "dapr-identity" {
  name                = "dapr-identity"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
}

resource "azurerm_federated_identity_credential" "dapr-federated" {
  name                = "dapr-federated"
  resource_group_name = azurerm_resource_group.rg.name
  audience            = ["api://AzureADTokenExchange"]
  issuer              = azurerm_kubernetes_cluster.cluster.oidc_issuer_url
  parent_id           = azurerm_user_assigned_identity.dapr-identity.id
  subject             = "system:serviceaccount:default:workload-identity-sa"
}


// Pub/Sub Service Bus Queues
resource "azurerm_servicebus_namespace" "sb-tvt" {
  name                = "dapr-service-bus"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku                 = "Basic"
}
resource "azurerm_servicebus_queue" "sb_queue" {
  name                = "orders"
  namespace_id        = azurerm_servicebus_namespace.sb-tvt.id
}
resource "azurerm_role_assignment" "servicebus_receiver" {
  principal_id   = azurerm_user_assigned_identity.dapr-identity.principal_id
  role_definition_name = "Azure Service Bus Data Receiver"
  scope          = azurerm_servicebus_namespace.sb-tvt.id
}

// Key Vault
resource "azurerm_key_vault" "kv" {
  name                        = "dapr-key-vault-jalbrecht"
  location                    = azurerm_resource_group.rg.location
  resource_group_name         = azurerm_resource_group.rg.name
  sku_name                    = "standard"
  tenant_id                   = data.azurerm_client_config.client.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.client.tenant_id
    object_id = data.azurerm_client_config.client.object_id

    key_permissions = [
      "Get", "List", "Update", "Create", "Delete"
    ]

    secret_permissions = [
      "Get", "List", "Set"
    ]

    storage_permissions = [
      "Get", "List", "Delete"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.client.tenant_id
    object_id = azurerm_user_assigned_identity.dapr-identity.principal_id

    key_permissions = [
      "Get", "List", "Update", "Create", "Delete"
    ]

    secret_permissions = [
      "Get", "List", "Set"
    ]

    storage_permissions = [
      "Get", "List", "Delete"
    ]
  }
}

resource "azurerm_key_vault_secret" "pubSubName" {
  name         = "pubSubName"
  value        = "azure-pubsub"
  key_vault_id = azurerm_key_vault.kv.id
}

resource "azurerm_key_vault_secret" "payment-key" {
  name         = "payment-key"
  value        = "ec37aa25501f5aea74d5eb3d19b08333"
  key_vault_id = azurerm_key_vault.kv.id
}

resource "azurerm_key_vault_secret" "servicebus-connection-string" {
  name         = "servicebus-connection-string"
  value        = azurerm_servicebus_namespace.sb-tvt.default_primary_connection_string
  key_vault_id = azurerm_key_vault.kv.id
}