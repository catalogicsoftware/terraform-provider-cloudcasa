terraform {
  required_providers {
    cloudcasa = {
      version = "1.0.0"
      source  = "cloudcasa.io/cloudcasa/cloudcasa"
    }
  }
}

provider "cloudcasa" {
  apikey = "API KEY HERE"
}

resource "cloudcasa_kubecluster" "testcluster" {
  name = "test_terraform_cluster"

  auto_install = true
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

  run_on_apply = true       # If true, the backup will run on each "terraform apply"
}

resource "cloudcasa_kubebackup" "test_offload" {
  name = "test_terraform_offload"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true
  copy_persistent_volumes = true

  run_on_apply = false

  policy_id = resource.cloudcasa_policy.testpolicy.id  
}

resource "cloudcasa_kubebackup" "example_kubebackup" {
  name = "Example Kubebackup"
  kubecluster_id = cloudcasa_kubecluster.example_cluster.id
  
  policy_id = cloudcasa_policy.example_schedule.id
  
  all_namespaces = true
  snapshot_persistent_volumes = true
  
  copy_persistent_volumes = true
  objectstore_id = cloudcasa_objectstore.example_s3.id
  copy_options = jsonencode({
    "pvc_parallelism": 4,
    "total_file_parallelism": 48,
    "upload_speed_limit": 100000000
  })
}

output "testcluster_data" {
  value = resource.cloudcasa_kubecluster.testcluster
}

output "testpolicy_data" {
  value = resource.cloudcasa_policy.testpolicy
}

output "testsnapshot_data" {
  value = resource.cloudcasa_kubebackup.test_snapshot
}

output "testoffload_data" {
  value = resource.cloudcasa_kubebackup.test_offload
}
