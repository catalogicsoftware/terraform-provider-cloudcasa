---
page_title: "cloudcasa_objectstore Resource - terraform-provider-cloudcasa"
subcategory: ""
description: |-
  Manages a CloudCasa objectstore resource for backup storage.
---

# cloudcasa_objectstore (Resource)

The objectstore resource allows you to configure and manage storage locations for CloudCasa backups. CloudCasa supports both S3 and Azure Blob Storage as backend storage providers.

## Example Usage

### AWS S3 Example

```terraform
resource "cloudcasa_objectstore" "s3_example" {
  name          = "my-s3-storage"
  provider_type = "s3"
  bucket_name   = "my-backup-bucket"
  endpoint_url  = "https://s3.amazonaws.com"
  region        = "us-east-1"
  access_key    = "AKIAIOSFODNN7EXAMPLE"
  secret_key    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

### Azure Storage Example

```terraform
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
```

### Private Objectstore With Proxy Cluster

```terraform
resource "cloudcasa_objectstore" "private_s3" {
  name          = "private-s3-storage"
  provider_type = "s3"
  bucket_name   = "private-backup-bucket"
  endpoint_url  = "https://s3.private.example.com"
  private       = true
  proxy_cluster = cloudcasa_kubecluster.proxy.id
}
```

## Argument Reference

The following arguments are supported:

### Common Fields
* `name` - (Required) The name of the objectstore resource.
* `provider_type` - (Required) The provider type for the objectstore. Valid values are "s3" or "azure".
* `private` - (Optional) Whether the storage location is isolated from CloudCasa servers. If true, `proxy_cluster` must be specified. Defaults to false.
* `proxy_cluster` - (Optional) The ID of the proxy cluster used for connecting to the object store. Required if `private` is true.
* `skip_tls_validation` - (Optional) Whether to skip TLS certificate validation when connecting to the storage provider. Defaults to false.

### AWS S3-specific Fields
* `bucket_name` - (Required for S3) The name of the S3 bucket.
* `endpoint_url` - (Required for S3) The endpoint URL for the S3 provider.
* `region` - (Optional for S3) The AWS region where the bucket is located.
* `access_key` - (Optional) The access key for authenticating with the S3 provider.
* `secret_key` - (Optional) The secret key for authenticating with the S3 provider.

### Azure-specific Fields
* `subscription_id` - (Required for Azure) The Azure subscription ID.
* `tenant_id` - (Required for Azure) The Azure tenant ID.
* `client_id` - (Required for Azure) The Azure client ID.
* `client_secret` - (Required for Azure) The Azure client secret.
* `resource_group_name` - (Required for Azure) The Azure resource group name.
* `storage_account_name` - (Required for Azure) The Azure storage account name.
* `region` - (Required for Azure) The Azure region.
* `government_cloud` - (Optional) Whether the Azure storage is in the Azure Government Cloud. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the objectstore resource.
* `created` - The timestamp when the objectstore was created.
* `updated` - The timestamp when the objectstore was last updated.
* `etag` - The ETag of the objectstore resource.

## Import

Objectstores can be imported using the resource `id`, e.g.,

```
$ terraform import cloudcasa_objectstore.example 1234567890abcdef12345678
``` 