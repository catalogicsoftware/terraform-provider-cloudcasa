# Define a basic snapshot and run on apply
resource "cloudcasa_kubebackup" "adhoc_snapshot_example" {
  name = "cloudcasa_adhoc_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  run_on_apply = true
}

# Define a basic snapshot to run on a schedule (requires a policy)
resource "cloudcasa_kubebackup" "scheduled_snapshot_example" {
  name = "cloudcasa_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = true
  snapshot_persistent_volumes = true
  copy_persistent_volumes = false

  run_on_apply = false
  policy_id = resource.cloudcasa_policy.example_policy.id  
}

# Define a snapshot with namespace selection and app hooks
resource "cloudcasa_kubebackup" "custom_snapshot_example" {
  name = "cloudcasa_custom_snapshot_example"
  kubecluster_id = resource.cloudcasa_kubecluster.example.id

  all_namespaces = false
  select_namespaces = [
    "test-csi-snapshot"
  ]
  pre_hooks = [
    {
      template = true,
      namespaces = ["default", "test-csi-snapshot"],
      hooks = ["61b3bb7b555abc4d71d0a7bf"]
    }
  ]

  snapshot_persistent_volumes = true

  copy_persistent_volumes = false

  run_on_apply = true
}

# Define a Copy backup to offload Persistent Volume data
resource "cloudcasa_kubebackup" "copy_example" {
  name = "cloudcasa_copy_example"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = true
  delete_snapshot_after_copy = false

  run_on_apply = true
}

# Define a Copy backup using a specific objectstore
resource "cloudcasa_objectstore" "example_s3" {
  name          = "example-s3-objectstore"
  provider_type = "s3"
  bucket_name   = "my-backup-bucket"
  endpoint_url  = "https://s3.amazonaws.com"
  region        = "us-east-1" 
  access_key    = "AKIAIOSFODNN7EXAMPLE"
  secret_key    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

resource "cloudcasa_kubebackup" "copy_with_objectstore" {
  name = "cloudcasa_copy_with_objectstore"
  kubecluster_id = resource.cloudcasa_kubecluster.testcluster.id

  all_namespaces = true
  snapshot_persistent_volumes = true

  copy_persistent_volumes = true
  delete_snapshot_after_copy = false
  objectstore_id = cloudcasa_objectstore.example_s3.id

  run_on_apply = true
}
