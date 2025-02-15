---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudcasa_kubecluster Resource - terraform-provider-cloudcasa"
subcategory: ""
description: |-
  CloudCasa kubecluster configuration
---

# cloudcasa_kubecluster (Resource)

CloudCasa kubecluster configuration

## Example Usage

```terraform
# Define a kubecluster resource and install the agent on the active cluster (using KUBECONFIG env var)
resource "cloudcasa_kubecluster" "exmaple_kubecluster" {
  name = "cloudcasa_example_kubecluster"

  auto_install = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) CloudCasa resource name

### Optional

- `auto_install` (Boolean) Automatically install the CloudCasa agent and register the current kubernetes cluster. Uses KUBECONFIG environment variable for cluster context.Set to false to manage or reference a CloudCasa kubecluster resource without installing the agent on a cluster.

### Read-Only

- `agent_url` (String) CloudCasa Kubeagent installation manifest URL
- `created` (String) Creation time of the CloudCasa resource
- `etag` (String) Etag generated by CloudCasa, used for updating resources in place
- `id` (String) CloudCasa resource ID
- `links` (Map of String) Related resources from CloudCasa
- `status` (Map of String) Cluster status from CloudCasa
- `updated` (String) Last update time of the CloudCasa resource

## Import

Import is supported using the following syntax:

```shell
# Kubecluster can be imported by specifying the CloudCasa ID
terraform import cloudcasa_kubecluster.example 64123abcdef123abcdef123a
```
