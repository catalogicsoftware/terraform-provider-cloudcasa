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

Below is a small example of how to initialize the provider, install the CloudCasa agent, and take a snapshot of a cluster. For more details and examples check the `docs` directory.

### Initialize the Provider

In your terraform manifest, create and configure the provider:

```hcl
terraform {
  required_providers {
    cloudcasa = {
      version = "1.2.1"
      source = "catalogicsoftware/cloudcasa"
    }
  }
}

provider "cloudcasa" {
  apikey = "API_KEY_HERE"

  # Optional: For Selfhosted CloudCasa, set server URL
  # cloudcasa_url = "https://cloudcasa.example.com"
  # Optional: allow insecure TLS connections to CloudCasa server
  # insecure_tls = true
}
```

### Create CloudCasa Resources

#### Kubecluster

A cloudcasa_kubecluster resource represents a Kubernetes cluster. You can import an existing CloudCasa
cluster using `terraform import` or define a new cluster.

To automatically install the CloudCasa agent on a cluster set `auto_install` to `true`. The provider
will apply the agent spec using the environment variable `KUBECONFIG` to find the cluster context.

```hcl
resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true
}
```

#### Kubebackup

The cloudcasa_kubebackup resource refers to both snapshots and copy backups. Cluster ID of a valid CloudCasa kubecluster is required.

If `run_on_apply` is True, the backup will be considered Adhoc and does not require a policy ID. With this setting the backup will run any time we run `terraform apply`, even if the backup has already been created. Terraform will wait up to 5 minutes for the job to complete.

You can set most options that are available in the CloudCasa UI. 

For example, here is a simple Adhoc snapshot job:

```hcl
resource "cloudcasa_kubebackup" "adhoc_snapshot_example" {
  name = "cloudcasa_adhoc_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  run_on_apply = true
}
```
For more examples see the `docs` directory.

#### Policy

Policies are required for backups that do not have run_on_apply set. They are created by defining a Cron schedule for the job:

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
