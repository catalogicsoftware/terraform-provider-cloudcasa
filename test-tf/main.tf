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
    apikey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAzLTA5VDE5OjU1OjEzLjQ3N1oiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NzgzOTE3MTQsImV4cCI6MTY3ODM5ODkxNCwic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoiTEEwTjBncm50dFlTWFNpVDhjWGNWTHZ2ZXlFNWpVc2oiLCJub25jZSI6ImRpNHhZM2xQYnpSVmIwdGlXV2xHVmt4WlZVWlZSa2hsWVhWb2NEZEdjR3gyWldOek5YZHNUSFYwWHc9PSJ9.RLWfs_ttDvZKZWviJTUHnW9htXWkFwcxWgt2xj7-TZGBVi0cTwYodjnUUTUsYw9ZVaLvt34KgDQqva-ioRC-QSPA70s5YtFHppFdO2TPzd1j3_S_ldVnqsAa6U9D7kLBAeV8ADcbNOjeajtrfiCCzT4qjtsVBX7ukL8nCxEi8lfdZcYqLxUSv0pF2F2Gx_Al5yPCG4T7ycjlauM_Bn0elI09U-kwVRdg9HIqPetVXjq2U2PRZnLacGI5gmdrnvj9_KQhzOAMHYZjcOzFhwAuxJFsOAR63vAS6omqoXAfnf8b_sfx_M5IYe8LP0pB6kewACP2pfDsrvwmpIO49G0ltw"

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

output "testcluster_data" {
  value = resource.cloudcasa_kubecluster.testcluster
}
