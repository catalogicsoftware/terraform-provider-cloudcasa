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
    apikey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAzLTEzVDIxOjQ2OjQ5Ljg1NFoiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2Nzg3NDQwMTEsImV4cCI6MTY3ODc1MTIxMSwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoid2pPT1Q3Uzd1MkJvTEhwYjlaYnYwNk5Tcmk2UHZqT08iLCJub25jZSI6IlduYzVXbjVOY25wek16RlpibWQyTldzeVZWWjVNMlIzVlZGRGNVMUtaVlUwV25ob1ZXRkliMWRUZHc9PSJ9.dFhbnKlfKwEevnbbamWADV3vjZyZPjSP5ztwGNCQgQaSXUBzDbNPQaC6ZL-tamdPwCGyfj3JfTMvgrrAYirA6kE5YowdHmFD2E7YgAyprGeWhjRYBpG_imEUWERkJxhq-Ae7RpEqFxYc8eqPNBYlgvwSY4u4HB9LwnWupAZtXLgDeAFNa6TFybEP61DJDBD13_RS5rCPPrsZHkq7aQuxEr0XDvwRtROLaBpfOyFgJs3XbXA8QkwxnaJH6YxDugUmsQF_Ip4pVzp6eOtBLJ12YU2_HICpChBFO7LiT1AxNm9OJ5MzOd6JxHtHUljAjGp27ADJRxVA0BAGXPP3oG2pZA"

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

// managed resources = "resources"
// "while managed resources cause Terraform to create, update, and delete 
// infrastructure objects, data resources cause Terraform only to read objects"


resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true

  # Auto installation requires KUBECONFIG var to be set
  # otherwise cluster creation will fail

  # This is an Alternade method we can document?
  # provisioner "local-exec" {
  #   command = "kubectl apply -f ${self.agent_url}"
  # }

  # TODO:
  # - add auto_install option to tf resource
  # - verify KUBECONFIG is set when auto_install is true
  # - perform agent apply & wait for ACTIVE state when auto_install true
}

resource "cloudcasa_kubebackup" "testbackup" {
  name = "test_terraform_kubebackup"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  run_on_apply = true

}

output "testcluster_data" {
  value = resource.cloudcasa_kubecluster.testcluster
}

output "testbackup_data" {
  value = resource.cloudcasa_kubebackup.testbackup
}
