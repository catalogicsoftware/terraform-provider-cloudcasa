terraform {
  required_providers {
    cloudcasa = {
      version = "0.0.21"
      source  = "cloudcasa.io/cloudcasa/cloudcasa"
    }
  }
}

// not used
provider "cloudcasa" {
    email = "jgarner@catalogicsoftware.com"
    apikey = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InM3NmtuNThRT2liTXRfZnNpVFlLMCJ9.eyJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9jb3VudHJ5IjoiVW5pdGVkIFN0YXRlcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL3RpbWV6b25lIjoiQW1lcmljYS9OZXdfWW9yayIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvdW50cnlfY29kZSI6IlVTIiwiaHR0cDovL3d3dy5jbG91ZGNhc2EuaW8vY291bnRyeV9jb2RlMyI6IlVTQSIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2ZpcnN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9sYXN0TmFtZSI6Ii0iLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9qb2JUaXRsZSI6IkRldm9wcyIsImh0dHA6Ly93d3cuY2xvdWRjYXNhLmlvL2NvbXBhbnkiOiJDYXRhbG9naWMgU29mdHdhcmUiLCJodHRwOi8vd3d3LmNsb3VkY2FzYS5pby9hd3NfbWFya2V0cGxhY2VfdG9rZW4iOiItIiwibmlja25hbWUiOiJqZ2FybmVyIiwibmFtZSI6IkpvbmF0aGFuIEdhcm5lciIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yOTlhNmJhNjhlNjEwOGFiYjY1MmY4ZTkwZTM0YjVhNj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRmpnLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDIzLTAyLTI0VDE5OjUxOjQ1LjQzNloiLCJlbWFpbCI6ImpnYXJuZXJAY2F0YWxvZ2ljc29mdHdhcmUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vYXV0aC5jbG91ZGNhc2EuaW8vIiwiYXVkIjoiSkVKU3plblBGeE5FUFEwaDY0ZDIzZTZRMEdKNXpRanQiLCJpYXQiOjE2NzcyNjgzMDcsImV4cCI6MTY3NzI3NTUwNywic3ViIjoiYXV0aDB8NWZhYzQ4NDg0MWQ3MDgwMDY4YTA2ZGM5Iiwic2lkIjoiZGZ6RWh5cmFobDg2bGZGc1FPQkc2WTBHb3FqLXVmamsiLCJub25jZSI6IlltUk9Ta0psT1hSSmMyMVNORnBpTW5nNVZqWTVWMnBGUkZSNmNtbzBMVUpMYkRodldHZCtPR0orTVE9PSJ9.yGCZl1jApKHLXEeIzXs-5bgunUdL6GHEns9bY-URMyMCmRKlMc7vAv12U1mwjwYs_3WWyPsZtvvdpS7y9c62hPVkR6XzhLr0JUDtwr8zC6KqvGfYSCXCwQZCroSERCRcC0-0OiaG-HrKyxiArgituMcXztyMHNdLAPEEzI5YvAPF3JV2AnDx57UsnCMJUii2KdCO9HQFtowPwilnFIved4pyisawWHy3lQtcXgfi_afiBGvEPQRYnWkh1kYnDMinbTKwe4eQCpofCmDMvwWB9p68iyHZWoUG6KM2gTFEctkIvONsAFbbeqvTmx3pwulExJXwxUUFqiXEgQBrTV9CpQ"
}

resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true

  # Auto installation requires KUBECONFIG var to be set
  # otherwise cluster creation will fail

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
