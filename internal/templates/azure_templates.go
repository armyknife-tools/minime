// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package templates

// generateAzureTemplates generates templates for Azure resources
func generateAzureTemplates() []Template {
	var templates []Template

	// Azure Virtual Machine
	templates = append(templates, Template{
		Provider:    "azure",
		Resource:    "virtual_machine",
		DisplayName: "Virtual Machine",
		Description: "Azure Virtual Machine with all common configuration options",
		Category:    "Compute",
		Tags:        "compute,vm,virtual machine",
		Content: `# Azure Virtual Machine Template
# This template includes all common properties for an Azure VM

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.example.id
  }
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "HTTP"
    priority                   = 1002
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_network_interface_security_group_association" "example" {
  network_interface_id      = azurerm_network_interface.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  
  # Use either admin_password or admin_ssh_key
  admin_password                  = "P@ssw0rd1234!"  # Not recommended for production
  disable_password_authentication = false
  
  # Uncomment to use SSH key instead of password
  # admin_ssh_key {
  #   username   = "adminuser"
  #   public_key = file("~/.ssh/id_rsa.pub")
  # }
  
  # Zone redundancy
  zone = "1"
  
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
    disk_size_gb         = 64
    
    # Disk encryption
    disk_encryption_set_id = azurerm_disk_encryption_set.example.id
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  
  # Custom data for cloud-init
  custom_data = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    echo "Hello from OpenTofu on Azure" > /var/www/html/index.html
    EOF
  )
  
  # Boot diagnostics
  boot_diagnostics {
    storage_account_uri = azurerm_storage_account.example.primary_blob_endpoint
  }
  
  # Identity (system or user assigned)
  identity {
    type = "SystemAssigned"
  }
  
  # Additional data disks
  additional_capabilities {
    ultra_ssd_enabled = true
  }
  
  # Automatic updates
  provision_vm_agent = true
  
  # Tags
  tags = {
    environment = "dev"
    application = "example"
  }
  
  # Proximity placement group for low latency
  proximity_placement_group_id = azurerm_proximity_placement_group.example.id
}

# Storage account for boot diagnostics
resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Disk encryption set
resource "azurerm_key_vault" "example" {
  name                        = "examplekeyvault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  enabled_for_disk_encryption = true
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true
  sku_name                    = "standard"
}

# Proximity placement group
resource "azurerm_proximity_placement_group" "example" {
  name                = "example-ppg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  tags = {
    environment = "dev"
  }
}
`,
	})

	// Azure Storage Account
	templates = append(templates, Template{
		Provider:    "azure",
		Resource:    "storage_account",
		DisplayName: "Storage Account",
		Description: "Azure Storage Account with all common configuration options",
		Category:    "Storage",
		Tags:        "storage,blob,file,queue,table",
		Content: `# Azure Storage Account Template
# This template includes all common properties for an Azure Storage Account

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  
  # Account configuration
  account_tier             = "Standard"  # Standard or Premium
  account_replication_type = "GRS"       # LRS, GRS, RAGRS, ZRS, GZRS, RAGZRS
  account_kind             = "StorageV2" # Storage, StorageV2, BlobStorage, BlockBlobStorage, FileStorage
  
  # Access configuration
  allow_nested_items_to_be_public = false
  shared_access_key_enabled       = true
  public_network_access_enabled   = true
  default_to_oauth_authentication = false
  
  # Blob configuration
  blob_properties {
    # Data protection
    versioning_enabled            = true
    change_feed_enabled           = true
    last_access_time_enabled      = true
    container_delete_retention_policy {
      days = 7
    }
    delete_retention_policy {
      days = 30
    }
    
    # CORS configuration
    cors_rule {
      allowed_headers    = ["*"]
      allowed_methods    = ["GET", "POST", "PUT"]
      allowed_origins    = ["https://example.com"]
      exposed_headers    = ["*"]
      max_age_in_seconds = 3600
    }
  }
  
  # Network configuration
  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["203.0.113.0/24"]
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
    bypass                     = ["Metrics", "AzureServices"]
  }
  
  # Identity (system or user assigned)
  identity {
    type = "SystemAssigned"
  }
  
  # Encryption configuration
  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }
  
  # Hierarchical namespace (for Data Lake Storage Gen2)
  is_hns_enabled = false
  
  # Lifecycle management
  lifecycle_rules {
    name    = "archiveold"
    enabled = true
    
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
    }
    
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 30
        tier_to_archive_after_days_since_modification_greater_than = 90
        delete_after_days_since_modification_greater_than          = 365
      }
      
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
      
      version {
        delete_after_days_since_creation = 90
      }
    }
  }
  
  # Tags
  tags = {
    environment = "dev"
    application = "example"
  }
}

# Storage container
resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private" # private, blob, container
}

# Storage queue
resource "azurerm_storage_queue" "example" {
  name                 = "messages"
  storage_account_name = azurerm_storage_account.example.name
}

# Storage table
resource "azurerm_storage_table" "example" {
  name                 = "records"
  storage_account_name = azurerm_storage_account.example.name
}

# Storage file share
resource "azurerm_storage_share" "example" {
  name                 = "sharename"
  storage_account_name = azurerm_storage_account.example.name
  quota                = 50
}

# Private endpoint for secure access
resource "azurerm_private_endpoint" "example" {
  name                = "example-endpoint"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  private_service_connection {
    name                           = "example-privateserviceconnection"
    private_connection_resource_id = azurerm_storage_account.example.id
    subresource_names              = ["blob"]
    is_manual_connection           = false
  }
}
`,
	})

	// Azure App Service
	templates = append(templates, Template{
		Provider:    "azure",
		Resource:    "app_service",
		DisplayName: "App Service",
		Description: "Azure App Service with all common configuration options",
		Category:    "Web",
		Tags:        "web,app,service,webapp",
		Content: `# Azure App Service Template
# This template includes all common properties for an Azure App Service

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_service_plan" "example" {
  name                = "example-appserviceplan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  # Pricing tier
  os_type  = "Linux"
  sku_name = "P1v2"  # Options: F1, B1, B2, B3, S1, S2, S3, P1v2, P2v2, P3v2, etc.
  
  # Zone redundancy
  zone_balancing_enabled = true
  
  tags = {
    environment = "dev"
  }
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-webapp"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id
  
  # Site configuration
  site_config {
    # Runtime stack
    application_stack {
      node_version = "16-lts"  # Options: dotnet, java, php, python, node, etc.
    }
    
    # Always On (prevents app from being unloaded)
    always_on = true
    
    # CORS configuration
    cors {
      allowed_origins     = ["https://example.com"]
      support_credentials = true
    }
    
    # HTTP version
    http2_enabled = true
    
    # IP restrictions
    ip_restriction {
      ip_address = "203.0.113.0/24"
      name       = "example-rule"
      priority   = 100
      action     = "Allow"
    }
    
    # Health check
    health_check_path = "/health"
    health_check_eviction_time_in_min = 2
    
    # WebSockets
    websockets_enabled = true
    
    # Minimum TLS version
    minimum_tls_version = "1.2"
    
    # FTP deployment
    ftps_state = "Disabled"  # Options: AllAllowed, FtpsOnly, Disabled
    
    # Container configuration (if using containers)
    # container_registry_use_managed_identity = true
    # container_registry_managed_identity_client_id = azurerm_user_assigned_identity.example.client_id
  }
  
  # App settings (environment variables)
  app_settings = {
    "WEBSITE_NODE_DEFAULT_VERSION" = "~16"
    "APPINSIGHTS_INSTRUMENTATIONKEY" = azurerm_application_insights.example.instrumentation_key
    "DATABASE_URL" = "postgresql://user:password@example-server.postgres.database.azure.com:5432/mydb"
    "ENVIRONMENT" = "development"
  }
  
  # Connection strings
  connection_string {
    name  = "Database"
    type  = "PostgreSQL"
    value = "postgresql://user:password@example-server.postgres.database.azure.com:5432/mydb"
  }
  
  # Identity (system or user assigned)
  identity {
    type = "SystemAssigned"
  }
  
  # Storage account for logs and backups
  storage_account {
    name         = "logs"
    type         = "AzureBlob"
    account_name = azurerm_storage_account.example.name
    access_key   = azurerm_storage_account.example.primary_access_key
    share_name   = "logs"
    mount_path   = "/var/log/webapp"
  }
  
  # Backup configuration
  backup {
    name                = "example-backup"
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
      retention_period_days = 30
      start_time           = "2023-01-01T04:00:00Z"
    }
    storage_account_url = "https://${azurerm_storage_account.example.name}.blob.core.windows.net/${azurerm_storage_container.backups.name}${data.azurerm_storage_account_sas.example.sas}"
  }
  
  # Logs configuration
  logs {
    application_logs {
      file_system_level = "Information"  # Options: Off, Error, Warning, Information, Verbose
      
      azure_blob_storage {
        level             = "Information"
        retention_in_days = 30
        sas_url           = "https://${azurerm_storage_account.example.name}.blob.core.windows.net/${azurerm_storage_container.logs.name}${data.azurerm_storage_account_sas.example.sas}"
      }
    }
    
    http_logs {
      file_system {
        retention_in_days = 7
        retention_in_mb   = 35
      }
    }
  }
  
  # Custom domains and SSL
  custom_domain_verification_id = azurerm_app_service_certificate_order.example.domain_verification_token
  
  # Sticky settings (don't replace during deployment)
  sticky_settings {
    app_setting_names       = ["ENVIRONMENT"]
    connection_string_names = ["Database"]
  }
  
  # Tags
  tags = {
    environment = "dev"
    application = "example"
  }
}

# Storage account for logs and backups
resource "azurerm_storage_account" "example" {
  name                     = "examplewebappstorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "logs" {
  name                  = "logs"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "backups" {
  name                  = "backups"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

# Application Insights for monitoring
resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

# Custom domain and SSL
resource "azurerm_app_service_certificate_order" "example" {
  name                = "example-cert"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  auto_renew          = true
  validity_in_years   = 1
  
  product_type        = "Standard"  # Options: Standard, WildCard
  
  # Domain name to secure
  distinguished_name  = "CN=example.com"
}

# Custom domain binding
resource "azurerm_app_service_custom_hostname_binding" "example" {
  hostname            = "www.example.com"
  app_service_name    = azurerm_linux_web_app.example.name
  resource_group_name = azurerm_resource_group.example.name
  ssl_state           = "SniEnabled"
  thumbprint          = azurerm_app_service_certificate.example.thumbprint
}
`,
	})

	return templates
}
