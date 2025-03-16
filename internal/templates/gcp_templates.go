// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package templates

// generateGCPTemplates generates templates for Google Cloud resources
func generateGCPTemplates() []Template {
	var templates []Template

	// GCP Compute Instance
	templates = append(templates, Template{
		Provider:    "gcp",
		Resource:    "compute_instance",
		DisplayName: "Compute Instance",
		Description: "Google Cloud Compute Instance with all common configuration options",
		Category:    "Compute",
		Tags:        "compute,vm,instance",
		Content: `# Google Cloud Compute Instance Template
# This template includes all common properties for a GCP Compute Instance

resource "google_compute_instance" "example" {
  name         = "example-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  
  # Allow stopping for update
  allow_stopping_for_update = true
  
  # Delete protection
  deletion_protection = false
  
  # Boot disk
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      size  = 50  # GB
      type  = "pd-ssd"  # pd-standard, pd-balanced, pd-ssd, pd-extreme
      
      # Disk labels
      labels = {
        environment = "dev"
      }
    }
    
    # Auto-delete disk when instance is deleted
    auto_delete = true
    
    # Disk encryption
    kms_key_self_link = google_kms_crypto_key.example.id
  }
  
  # Additional disk
  attached_disk {
    source      = google_compute_disk.example.id
    device_name = "data-disk"
    mode        = "READ_WRITE"  # READ_ONLY, READ_WRITE
  }

  # Network interface
  network_interface {
    network = "default"
    
    # Uncomment to assign a public IP
    access_config {
      nat_ip = google_compute_address.example.address
      network_tier = "PREMIUM"  # PREMIUM, STANDARD
    }
    
    # Alias IP ranges
    alias_ip_range {
      ip_cidr_range = "/24"
      subnetwork_range_name = "secondary-range"
    }
  }
  
  # Additional network interface
  network_interface {
    network = google_compute_network.example.id
    subnetwork = google_compute_subnetwork.example.id
  }

  # Metadata
  metadata = {
    ssh-keys = "user:${file("~/.ssh/id_rsa.pub")}"
    startup-script = <<-EOF
      #!/bin/bash
      apt-get update
      apt-get install -y nginx
      echo "Hello from OpenTofu on GCP" > /var/www/html/index.html
    EOF
  }
  
  # Metadata startup script from file
  metadata_startup_script = file("startup.sh")

  # Service account
  service_account {
    email  = google_service_account.example.email
    scopes = ["cloud-platform"]
  }

  # Labels
  labels = {
    environment = "dev"
    application = "example"
  }
  
  # Tags for firewall rules
  tags = ["web", "app"]
  
  # Scheduling options
  scheduling {
    preemptible         = false
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"  # MIGRATE, TERMINATE
    provisioning_model  = "STANDARD" # STANDARD, SPOT
    
    # Node affinity for sole-tenant nodes
    node_affinities {
      key      = "compute.googleapis.com/node-group-name"
      operator = "IN"
      values   = [google_compute_node_group.example.name]
    }
  }
  
  # Shielded VM options
  shielded_instance_config {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }
  
  # Confidential computing
  confidential_instance_config {
    enable_confidential_compute = true
  }
  
  # Advanced machine features
  advanced_machine_features {
    enable_nested_virtualization = true
    threads_per_core             = 2
  }
  
  # Resource policies (e.g., for scheduled snapshots)
  resource_policies = [google_compute_resource_policy.example.id]
}

# Static IP address
resource "google_compute_address" "example" {
  name         = "example-address"
  address_type = "EXTERNAL"  # INTERNAL, EXTERNAL
  region       = "us-central1"
  
  # For internal addresses
  # subnetwork = google_compute_subnetwork.example.id
  # address    = "10.0.0.10"
  
  # Network tier for external addresses
  network_tier = "PREMIUM"  # PREMIUM, STANDARD
  
  # Labels
  labels = {
    environment = "dev"
  }
}

# Additional disk
resource "google_compute_disk" "example" {
  name  = "data-disk"
  type  = "pd-ssd"
  zone  = "us-central1-a"
  size  = 100  # GB
  
  # Source image or snapshot
  # image = "debian-cloud/debian-11"
  # snapshot = google_compute_snapshot.example.id
  
  # Disk encryption
  disk_encryption_key {
    kms_key_self_link = google_kms_crypto_key.example.id
  }
  
  # Physical block size
  physical_block_size_bytes = 4096
  
  # Labels
  labels = {
    environment = "dev"
  }
}

# KMS key for disk encryption
resource "google_kms_key_ring" "example" {
  name     = "example-keyring"
  location = "global"
}

resource "google_kms_crypto_key" "example" {
  name     = "example-key"
  key_ring = google_kms_key_ring.example.id
  
  # Rotation period
  rotation_period = "7776000s"  # 90 days
  
  # Purpose
  purpose = "ENCRYPT_DECRYPT"
}
`,
	})

	// GCP Cloud Storage Bucket
	templates = append(templates, Template{
		Provider:    "gcp",
		Resource:    "storage_bucket",
		DisplayName: "Cloud Storage Bucket",
		Description: "Google Cloud Storage Bucket with all common configuration options",
		Category:    "Storage",
		Tags:        "storage,bucket,gcs",
		Content: `# Google Cloud Storage Bucket Template
# This template includes all common properties for a GCP Storage Bucket

resource "google_storage_bucket" "example" {
  name          = "example-bucket-name"
  location      = "US"  # Options: US, EU, ASIA, or specific regions like us-central1
  
  # Storage class
  storage_class = "STANDARD"  # STANDARD, NEARLINE, COLDLINE, ARCHIVE
  
  # Uniform bucket-level access
  uniform_bucket_level_access = true
  
  # Public access prevention
  public_access_prevention = "enforced"  # enforced, inherited
  
  # Versioning
  versioning {
    enabled = true
  }
  
  # Website configuration
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
  
  # CORS configuration
  cors {
    origin          = ["https://example.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
  
  # Lifecycle rules
  lifecycle_rule {
    condition {
      age                   = 30
      created_before        = "2023-01-01"
      with_state            = "ARCHIVED"  # LIVE, ARCHIVED, ANY
      matches_storage_class = ["STANDARD"]
      num_newer_versions    = 3
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }
  
  lifecycle_rule {
    condition {
      age                = 90
      matches_prefix     = ["logs/"]
      matches_suffix     = [".log"]
      num_newer_versions = 3
    }
    action {
      type = "Delete"
    }
  }
  
  # Retention policy
  retention_policy {
    retention_period = 2592000  # 30 days in seconds
    is_locked        = false
  }
  
  # Logging
  logging {
    log_bucket        = google_storage_bucket.logs.name
    log_object_prefix = "bucket-logs"
  }
  
  # Object encryption
  encryption {
    default_kms_key_name = google_kms_crypto_key.example.id
  }
  
  # Custom placement config for dual-region buckets
  # custom_placement_config {
  #   data_locations = ["US-EAST1", "US-WEST1"]
  # }
  
  # Autoclass for automatic transition between storage classes
  autoclass {
    enabled = true
    terminal_storage_class = "ARCHIVE"
  }
  
  # Labels
  labels = {
    environment = "dev"
    application = "example"
  }
  
  # Force destroy (allows Terraform to destroy non-empty buckets)
  force_destroy = false
  
  # Requester pays
  requester_pays = false
}

# IAM policy for bucket
resource "google_storage_bucket_iam_binding" "example" {
  bucket = google_storage_bucket.example.name
  role   = "roles/storage.objectViewer"
  
  members = [
    "user:jane@example.com",
    "group:admins@example.com",
    "serviceAccount:my-service-account@project-id.iam.gserviceaccount.com",
  ]
}

# Default object ACL
resource "google_storage_bucket_acl" "example" {
  bucket = google_storage_bucket.example.name
  
  # Predefined ACL
  predefined_acl = "private"  # private, publicRead, publicReadWrite, etc.
  
  # Or custom role entities
  role_entity = [
    "OWNER:user-jane@example.com",
    "READER:group-admins@example.com",
  ]
}

# Signed URL for temporary access
data "google_iam_policy" "example" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "serviceAccount:${google_service_account.example.email}",
    ]
    condition {
      title       = "Temporary access"
      description = "Temporary access for 1 hour"
      expression  = "request.time < timestamp('2023-12-31T23:59:59Z')"
    }
  }
}

# Service account for signed URLs
resource "google_service_account" "example" {
  account_id   = "storage-account"
  display_name = "Storage Service Account"
}

# Logs bucket
resource "google_storage_bucket" "logs" {
  name          = "example-logs-bucket"
  location      = "US"
  storage_class = "STANDARD"
  
  # Logs retention
  lifecycle_rule {
    condition {
      age = 30
    }
    action {
      type = "Delete"
    }
  }
}
`,
	})

	// GCP Cloud SQL Instance
	templates = append(templates, Template{
		Provider:    "gcp",
		Resource:    "cloud_sql",
		DisplayName: "Cloud SQL Instance",
		Description: "Google Cloud SQL Instance with all common configuration options",
		Category:    "Database",
		Tags:        "database,sql,mysql,postgres",
		Content: `# Google Cloud SQL Instance Template
# This template includes all common properties for a GCP Cloud SQL Instance

resource "google_sql_database_instance" "example" {
  name             = "example-instance"
  database_version = "POSTGRES_14"  # MYSQL_8_0, POSTGRES_14, etc.
  region           = "us-central1"
  
  # Delete protection
  deletion_protection = false
  
  # Settings
  settings {
    # Instance type and resources
    tier              = "db-custom-2-8192"  # db-f1-micro, db-g1-small, custom tiers
    edition           = "ENTERPRISE"        # ENTERPRISE, ENTERPRISE_PLUS
    availability_type = "REGIONAL"          # ZONAL, REGIONAL
    disk_type         = "PD_SSD"            # PD_SSD, PD_HDD
    disk_size         = 100                 # GB
    disk_autoresize   = true
    disk_autoresize_limit = 500             # GB
    
    # Connectivity
    ip_configuration {
      ipv4_enabled        = true
      private_network     = google_compute_network.example.id
      require_ssl         = true
      ssl_mode            = "TRUSTED_CLIENT_CERTIFICATE"
      allocated_ip_range  = google_compute_global_address.example.name
      
      # Authorized networks
      authorized_networks {
        name  = "office"
        value = "203.0.113.0/24"
      }
      
      # PSC (Private Service Connect) settings
      psc_config {
        psc_enabled               = true
        allowed_consumer_projects = ["my-project-id"]
      }
    }
    
    # Backup configuration
    backup_configuration {
      enabled                        = true
      start_time                     = "02:00"  # UTC time (HH:MM format)
      location                       = "us"     # us, eu, asia
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 30
        retention_unit   = "COUNT"  # COUNT, TIME
      }
    }
    
    # Maintenance window
    maintenance_window {
      day          = 7   # 1-7 (Monday-Sunday)
      hour         = 3   # 0-23
      update_track = "stable"  # stable, preview, canary
    }
    
    # Insights configuration
    insights_config {
      query_insights_enabled  = true
      query_string_length     = 4096
      record_application_tags = true
      record_client_address   = true
      query_plans_per_minute  = 10
    }
    
    # Database flags
    database_flags {
      name  = "max_connections"
      value = "100"
    }
    
    database_flags {
      name  = "log_min_duration_statement"
      value = "1000"  # ms
    }
    
    # User labels
    user_labels = {
      environment = "dev"
      application = "example"
    }
    
    # Password policy
    password_validation_policy {
      min_length                  = 8
      complexity                  = "COMPLEXITY_DEFAULT"
      reuse_interval              = 5
      disallow_username_substring = true
      password_change_interval    = "30d"
      enable_password_policy      = true
    }
    
    # Data cache configuration
    data_cache_config {
      data_cache_enabled = true
    }
    
    # SQL Server specific settings (for SQL Server only)
    # sql_server_audit_config {
    #   bucket             = google_storage_bucket.example.name
    #   upload_interval    = "300s"
    #   retention_interval = "86400s"
    # }
    
    # Active Directory configuration (for SQL Server only)
    # active_directory_config {
    #   domain = "example.com"
    # }
    
    # Deny maintenance period
    deny_maintenance_period {
      start_date = "2023-11-01"
      end_date   = "2023-12-31"
      time       = "00:00:00"
    }
    
    # Connector enforcement
    connector_enforcement = "REQUIRED"  # NOT_REQUIRED, REQUIRED
  }
  
  # Read replicas
  replica_configuration {
    failover_target = false
    
    # MySQL specific settings
    connect_retry_interval    = 60  # seconds
    dump_file_path            = "gs://example-bucket/example-dump.sql"
    master_heartbeat_period   = 100000  # microseconds
    password                  = "replicate"
    username                  = "replicate"
    verify_server_certificate = true
    client_certificate        = google_sql_ssl_cert.client.cert
    client_key                = google_sql_ssl_cert.client.private_key
    ca_certificate            = google_sql_ssl_cert.server.server_ca_cert
  }
  
  # Restore from backup
  # restore_backup_context {
  #   backup_run_id = google_sql_backup_run.example.id
  #   project       = "my-project-id"
  #   instance_id   = "source-instance"
  # }
  
  # Clone from source
  # clone {
  #   source_instance_name = google_sql_database_instance.source.id
  #   point_in_time        = "2023-01-01T00:00:00Z"
  # }
  
  # Depends on VPC
  depends_on = [
    google_service_networking_connection.private_vpc_connection
  ]
}

# Database
resource "google_sql_database" "example" {
  name     = "example-database"
  instance = google_sql_database_instance.example.name
  charset  = "UTF8"
  collation = "en_US.UTF8"
}

# Database user
resource "google_sql_user" "example" {
  name     = "example-user"
  instance = google_sql_database_instance.example.name
  password = "example-password"  # Use Secret Manager in production
  
  # User type
  type = "BUILT_IN"  # BUILT_IN, CLOUD_IAM_USER, CLOUD_IAM_SERVICE_ACCOUNT
  
  # Password policy
  password_policy {
    allowed_failed_attempts      = 3
    password_expiration_duration = "30d"
    enable_failed_attempts_check = true
    enable_password_verification = true
  }
}

# SSL certificate
resource "google_sql_ssl_cert" "client" {
  common_name = "example-cert"
  instance    = google_sql_database_instance.example.name
}

# VPC network for private IP
resource "google_compute_network" "example" {
  name                    = "example-network"
  auto_create_subnetworks = false
}

# Global address for private services access
resource "google_compute_global_address" "example" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.example.id
}

# Private services connection
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.example.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.example.name]
}

# IAM binding for Cloud SQL
resource "google_project_iam_binding" "cloud_sql_client" {
  project = "my-project-id"
  role    = "roles/cloudsql.client"
  
  members = [
    "serviceAccount:${google_service_account.example.email}",
  ]
}

# Service account for Cloud SQL
resource "google_service_account" "example" {
  account_id   = "cloud-sql-client"
  display_name = "Cloud SQL Client Service Account"
}
`,
	})

	return templates
}
