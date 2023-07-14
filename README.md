# CloudCasa Terraform Provider

This is the [CloudCasa](https://cloudcasa.io) provider for [Terraform](https://www.terraform.io/).

This provider allows you to install the CloudCasa agent and manage backups in your Kubernetes cluster using Terraform.

## Contents

* [Requirements](#requirements)
* [Getting Started](#getting-started)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) v1.x
- A CloudCasa API key - Visit [CloudCasa](https://home.cloudcasa.io) to sign up and create an API key under Configuration -> API Keys
- [Go](https://golang.org/doc/install) v1.18.x (to build the provider plugin)
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) (for CloudCasa agent installation)

## Getting Started

For now the provider must be built from source and added to Terraform's development overrides. This means `terraform init` will not update the provider package, and instead the provider will be read directly from the binary path each time we run a terraform command. Because of this, just use `terraform plan` or `terraform apply` and skip `terraform init`.

Check the `examples` folder for Terraform manifest templates and for a yaml containing example wordpress PVCs and deployments to test backups.

### Build from Source

Until we add the provider to the Hashicorp registry, we must build the provider from source in the main directory:
```bash
go build -o terraform-provider-cloudcasa
```

After building the binary, set a developer override for the package by editing your `.terraformrc` file. By default this is at `/home/.terraformrc`. Add the path to the binary's folder as such:

```hcl
provider_installation {
  dev_overrides {
    "cloudcasa.io/cloudcasa/cloudcasa" = "/home/jon/work/terraform-provider-cloudcasa"
  }
  direct {}
}
```

### Initialize the Provider

In your terraform manifest, create and configure the provider:

```hcl
terraform {
  required_providers {
    cloudcasa = {
      version = "1.0.0"
      source  = "cloudcasa.io/cloudcasa/cloudcasa"
    }
  }
}

provider "cloudcasa" {
  apikey = "API_KEY_HERE"
}
```

### Create CloudCasa Resources

#### Kubecluster

Set your KUBECONFIG environment variable to the path of the cluster's kubeconfig file, otherwise agent installation will fail. OR Create a cluster without installing the agent, but this will only create a pending cluster resource in CloudCasa.

```hcl
resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true
}
```

#### Kubebackup

The cloudcasa_kubebackup resource refers to both snapshots and copy backups. Cluster ID of a valid CloudCasa kubecluster is required.

If `run_after_create` is True, the backup will be considered Adhoc and does not require a policy ID. With this setting the backup will run any time we run `terraform apply`, even if the backup has already been created. Terraform will wait up to 5 minutes for the job to complete.

You can set most options that are available in the CloudCasa UI. 

For example, here is an Adhoc snapshot job with namespace filters and a pre-hook:

```hcl
resource "cloudcasa_kubebackup" "test_snapshot" {
  name = "test_terraform_kubebackup"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]

  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  pre_hooks = [
    {template = true, namespaces = ["default", "test-csi-snapshot"], hooks = ["61b3bb7b555abc4d71d0a7bf"]}
  ]

  run_after_create = true       # If true, the backup will run on each "terraform apply"
}
```

Another example, here is a Copy job with a policy attached which will not run_after_create:

```hcl
resource "cloudcasa_kubebackup" "test_offload" {
  name = "test_terraform_offload"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true
  copy_persistent_volumes = true
  delete_snapshot_after_copy = false

  run_after_create = false

  policy_id = resource.cloudcasa_policy.testpolicy.id  
}
```

#### Policy

Policies are required for backups that do not have run_after_create set. They are created by defining a Cron schedule for the job:

```hcl
resource "cloudcasa_policy" "testpolicy" {
  name = "test_terraform_policy"
  timezone = "America/New_York"
  schedules = [
    {
      retention = 30,
      cron_spec = "30 0 * * MON,FRI",
      locked = false,
    }
  ]
}
```

### Importing Resources

You can import existing CloudCasa resources to manage them in Terraform using `terraform import`. For example, assume we have created a policy named "test_manual_policy" in CloudCasa UI. First create an empty resource for this policy:

```hcl
resource "cloudcasa_policy" "importtest" {
  name = "test_manual_policy"   # Name of the policy resource in CloudCasa
}
```

Get the ID from the CloudCasa UI (or casa) and use `terraform import <resource_state_path> <CC ID>`:

```bash
terraform import cloudcasa_policy.importtest 64948e5160a55cbabb5625f5
```

After importing, Terraform will try to apply any local changes to the CloudCasa resource the next time you apply. Check the resource in Terraform using `terraform state show <resource_state_path>` and update the configuration values to match, and make any desired changes. Changes in Terraform will always supercede changes in CloudCasa!

For Terraform v1.4+, you can add an `import` block directly to the Terraform config to avoid using `terraform import` manually each time:

```hcl
import {
  provider =
  to = cloudcasa_policy.importtest
  id = "64948e5160a55cbabb5625f5"
}
resource "cloudcasa_policy" "importtest" {
  name = "test_manual_policy"
}
```
