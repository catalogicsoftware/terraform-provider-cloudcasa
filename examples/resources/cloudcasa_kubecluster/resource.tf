# Example kubecluster
resource "cloudcasa_kubecluster" "example" {
  name = "test_terraform_cluster"

  auto_install = true
}
