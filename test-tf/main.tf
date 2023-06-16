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
  apikey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTA2LTE2VDE5OjA1OjI4LjMyN1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2ODY5NDIzMzAsImV4cCI6MTY4Njk1NjczMCwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoiMllObkpLVWlWZlhpeVZaSlRQeGFkaE5ZZmotRFhTbmUiLCJub25jZSI6IlVIRkdSRzEwZGsweFRHdFhWbFpUT1RSQlJEaEtXbFJQUWtwU1RtSk5SRWRCTFdkaVYyTlFMa2N4ZHc9PSJ9.g8Dn3WJO3WHhPlbHeOWZF0zpL49AnANG4DZ-3oFqhBd70Z3GrzifZGWMeYPElfZQwRvYP2gJ-wkL05QtvNDa2kjMHsAdv5P0dnamMcblduAdzj82rDJhNNvLCmIowemkU6P8dcxaztTYkolYkh3KL6x7QyOc2hKGGzwTWs3MocKMWjQjG7VAhYhk2H9ha3Zy5q-KehNiFQVjZhkWUyWtv1AXNQTYQ5uEFT8k1QzWIJLe-yZTHDKO-kLg1BhMu-r721-hD39Vt6m2eCQEL2Ux6wLStRYZZ2QlyPYJauojeyVKTVBUnShaNwtdUKsG4j84w3ZPFZcGwx8C5QbIemzyfQ"

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


resource "cloudcasa_kubebackup" "testbackup" {
  name = "test_terraform_kubeoffload"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]
  snapshot_persistent_volumes = true

  copy_persistent_volumes = true
  # TODO: require cloudcasa_kubeoffload resource for copies

  pre_hooks = [
    {template = true, namespaces = ["default", "test-csi-snapshot"], hooks = ["61b3bb7b555abc4d71d0a7bf"]}
  ]

  run_on_apply = true

}

output "testcluster_data" {
  value = resource.cloudcasa_kubecluster.testcluster
}

output "testbackup_data" {
  value = resource.cloudcasa_kubebackup.testbackup
}
