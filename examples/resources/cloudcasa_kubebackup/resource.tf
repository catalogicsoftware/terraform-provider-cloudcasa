# Example kubebackup
resource "cloudcasa_kubebackup" "example" {
  name = "test_terraform_kubebackup"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]

  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  pre_hooks = [
    {template = true, namespaces = ["default", "test-csi-snapshot"], hooks = ["61b3bb7b555abc4d71d0a7bf"]}
  ]

  run_after_create = true
}
