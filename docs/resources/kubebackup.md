---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudcasa_kubebackup Resource - terraform-provider-cloudcasa"
subcategory: ""
description: |-
  CloudCasa kubebackup configuration
---

# cloudcasa_kubebackup (Resource)

CloudCasa kubebackup configuration

## Example Usage

```terraform
# Define a basic snapshot and run on apply
resource "cloudcasa_kubebackup" "adhoc_snapshot_example" {
  name = "cloudcasa_adhoc_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  run_after_create = true
}

# Define a basic snapshot to run on a schedule (requires a policy)
resource "cloudcasa_kubebackup" "scheduled_snapshot_example" {
  name = "cloudcasa_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = true
  snapshot_persistent_volumes = true
  copy_persistent_volumes = false

  run_after_create = false
  policy_id = resource.cloudcasa_policy.example_policy.id  
}

# Define a snapshot with namespace selection and app hooks
resource "cloudcasa_kubebackup" "custom_snapshot_example" {
  name = "cloudcasa_custom_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]
  pre_hooks = [
    {
      template = true,
      namespaces = ["default", "test-csi-snapshot"],
      hooks = ["61b3bb7b555abc4d71d0a7bf"]
    }
  ]

  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  run_after_create = true
}

# Define a Copy backup to offload Persistent Volume data
resource "cloudcasa_kubebackup" "copy_example" {
  name = "cloudcasa_copy_example"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = true
  delete_snapshot_after_copy = false

  run_after_create = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `all_namespaces` (Boolean) Set to backup all namespaces, otherwise use the select_namespaces attribute to list namespaces
- `kubecluster_id` (String) ID of the kubecluster to back up
- `name` (String) CloudCasa resource name
- `snapshot_persistent_volumes` (Boolean) Set to snapshot persistent volumes. If false, PVs will be ignored

### Optional

- `copy_persistent_volumes` (Boolean) If true, persistent volume data will be copied and offloaded to S3 storage. This will create and manage an associated kubeoffload resource in CloudCasa
- `delete_snapshot_after_copy` (Boolean) Set to delete resource snapshots after performing data offload
- `kubeoffload_id` (String) ID of the associated kubeoffload resource created for Copy backups
- `policy_id` (String) ID of a policy for scheduling this backup
- `post_hooks` (Attributes List) Post-backup app hooks to execute. See https://docs.cloudcasa.io/help/configuration-apphook.html for details (see [below for nested schema](#nestedatt--post_hooks))
- `pre_hooks` (Attributes List) Pre-backup app hooks to execute. See https://docs.cloudcasa.io/help/configuration-apphook.html for details (see [below for nested schema](#nestedatt--pre_hooks))
- `retention` (Number) Number of days to retain backup data for
- `run_after_create` (Boolean) Set to run the backup immediately after creation or update. If enabled, this will also cause the backup to run on each terraform apply
- `select_namespaces` (List of String) List of namespaces to include in the backup

### Read-Only

- `created` (String) Creation time of the CloudCasa resource
- `etag` (String) Etag generated by CloudCasa, used for updating resources in place
- `id` (String) CloudCasa resource ID
- `offload_etag` (String) Etag of the associated offload resource generated by CloudCasa, used for updating resources in place
- `updated` (String) Last update time of the CloudCasa resource

<a id="nestedatt--post_hooks"></a>
### Nested Schema for `post_hooks`

Required:

- `hooks` (List of String) ID of a hook created in CloudCasa
- `namespaces` (List of String) List of namespaces to run the selected hook in
- `template` (Boolean) Set to use a predefined hook template


<a id="nestedatt--pre_hooks"></a>
### Nested Schema for `pre_hooks`

Required:

- `hooks` (List of String) ID of a hook created in CloudCasa
- `namespaces` (List of String) List of namespaces to run the selected hook in
- `template` (Boolean) Set to use a predefined hook template

## Import

Import is supported using the following syntax:

```shell
# Kubebackup can be imported by specifying the CloudCasa ID
terraform import cloudcasa_kubebackup.example 64123abcdef123abcdef123a
```
