terraform {
  required_providers {
    cloudcasa = {
      version = "0.0.1"
      source  = "cloudcasa.io/cloudcasa/cloudcasa"
    }
  }
}

# provider "kubernetes" {
#   config_path    = "~/work/test-eks-cluster.yaml"
# }

provider "cloudcasa" {
  apikey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTA2LTI3VDE5OjQyOjQ3Ljg1N1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2ODc4OTQ5NjksImV4cCI6MTY4NzkwOTM2OSwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoiYkl5VTZ6ZmMwN3ZpaGN1SEx4Y2tYMVM2U28zb2dDV2MiLCJub25jZSI6IlpDNXhMbUp5WWpFNFNWQjJWME5TYjJsVExVcDZUMjVQV2tJek1XdHhiemRRYkUxVVpFZFNiVk5WVUE9PSJ9.AqdjO5sy6wCT4z20gmsb0DphdnF26JUjxZj4U5G2l4cq2vEQfGbqxAFMRskYF68SBnNkmxtnw_6VIT9cASHq_NiZbJdXDGB51KpuZF58Gwaq3E2WeCSfY4RgQJWjnbE5bmFC5Zxbs5S_N3Hmuwol-jO71EE_clF9pPd6v1hgydlV5rW23szx3cnAGhq63ITJ4nA4xbdXGwWF4mnZK7EZZGI2ROUXe6TCPfBbBfJcXTObggo_eQvurdIcl8Y1nnTuvjBhSVFkZfCmb68NBSBQRUIDDCJW8hAs5CACORXRSbaCdTaUtN7Oer0TQNsfFE4_ulwkMsfXj9f58pLGklWpjg"

  # kubernetes {
  #   config_path = "~/work/test-eks-cluster.yaml"
  # }

}

// TODO: look at all the ways Helm is able to supply cluster credentials
// https://registry.terraform.io/providers/hashicorp/helm/latest/docs
// 3-4 different methods. we should use the same
// https://github.com/hashicorp/terraform-provider-helm/blob/main/helm/provider.go


// TODO: Use datasource instead of resource?
// https://developer.hashicorp.com/terraform/language/data-sources
// Datasources are defined outside of Terraform
// Data block requests that Terraform READS from a given source
// So not for kubeclusters, because those are defined by the user
// So anything is created by CLOUDCASA = Datasource

// managed resources = "resources"
// "while managed resources cause Terraform to create, update, and delete 
// infrastructure objects, data resources cause Terraform only to read objects"

resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true

  # Auto installation requires KUBECONFIG var to be set
  # otherwise cluster creation will fail
  # TODO: verify KUBECONFIG is set when auto_install is true

  # TODO: backup bucket config/advanced options

  # Optional: Install wordpress deployment
  # In a real environment this should be added as a standalone resource
  # provisioner "local-exec" {
  #   command = "kubectl apply -f deployments/minikube-wordpress-application.yaml"
  # }

}

resource "cloudcasa_policy" "testpolicy" {
  name = "test_terraform_policy"
  timezone = "America/New_York"
  schedules = [
    {
      retention = 12,
      cron_spec = "30 0 * * MON,FRI",
      locked = false,
    }
  ]
}

resource "cloudcasa_policy" "importtest" {
  name = "test_manual_policy"

}

resource "cloudcasa_kubebackup" "testbackup" {
  name = "test_terraform_kubeoffload"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]
  snapshot_persistent_volumes = true

  copy_persistent_volumes = false
  # TODO: require cloudcasa_kubeoffload resource for copies

  pre_hooks = [
    {template = true, namespaces = ["default", "test-csi-snapshot"], hooks = ["61b3bb7b555abc4d71d0a7bf"]}
  ]

  run_after_create = true

  policy_id = resource.cloudcasa_policy.testpolicy.id  

}

# resource "cloudcasa_kubebackup" "testbackup" {
#   name = "test_terraform_offload_2"
#   kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

#   all_namespaces = false
#   select_namespaces = [
#     "test-csi-snapshot"
#   ]

#   snapshot_persistent_volumes = true
#   copy_persistent_volumes = true

#   run_after_create = true

#   pre_hooks = [
#     {template = true, namespaces = ["default", "test-csi-snapshot"], hooks = ["61b3bb7b555abc4d71d0a7bf"]}
#   ]
# }

output "testcluster_data" {
  value = resource.cloudcasa_kubecluster.testcluster
}

output "testbackup_data" {
  value = resource.cloudcasa_kubebackup.testbackup
}

output "testpolicy_data" {
  value = resource.cloudcasa_policy.testpolicy
}
