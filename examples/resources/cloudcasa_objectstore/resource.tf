# AWS S3 Example
resource "cloudcasa_objectstore" "s3_example" {
  name          = "my-s3-storage"
  provider_type = "s3"
  bucket_name   = "my-backup-bucket"
  endpoint_url  = "https://s3.amazonaws.com"
  region        = "us-east-1"
  access_key    = "AKIAIOSFODNN7EXAMPLE"
  secret_key    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

# Azure Storage Example
resource "cloudcasa_objectstore" "azure_example" {
  name                = "my-azure-storage"
  provider_type       = "azure"
  subscription_id     = "00000000-0000-0000-0000-000000000000"
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "client_secret_value"
  region              = "eastus"
  resource_group_name = "my-resource-group"
  storage_account_name = "mystorageaccount"
}

# Private Objectstore With Proxy Cluster
resource "cloudcasa_kubecluster" "proxy" {
  name        = "proxy-cluster"
  auto_install = true
}

resource "cloudcasa_objectstore" "private_s3" {
  name          = "private-s3-storage"
  provider_type = "s3"
  bucket_name   = "private-backup-bucket"
  endpoint_url  = "https://s3.private.example.com"
  private       = true
  proxy_cluster = cloudcasa_kubecluster.proxy.id
} 