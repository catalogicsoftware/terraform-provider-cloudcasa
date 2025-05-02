# Define a kubecluster resource and install the agent on the active cluster (using KUBECONFIG env var)
resource "cloudcasa_kubecluster" "example_kubecluster" {
  name = "cloudcasa_example_kubecluster"

  auto_install = true
}

# Define a kubecluster resource with a reference to an objectstore for backups
resource "cloudcasa_objectstore" "example_s3" {
  name          = "example-s3-objectstore"
  provider_type = "s3"
  bucket_name   = "my-backup-bucket"
  endpoint_url  = "https://s3.amazonaws.com"
  region        = "us-east-1"
  access_key    = "AKIAIOSFODNN7EXAMPLE"
  secret_key    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}

resource "cloudcasa_kubecluster" "example_with_objectstore" {
  name           = "cluster-with-objectstore"
  auto_install   = true
  objectstore_id = cloudcasa_objectstore.example_s3.id
}
