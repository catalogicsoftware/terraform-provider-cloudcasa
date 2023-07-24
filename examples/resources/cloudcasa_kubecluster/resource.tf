# Define a kubecluster resource and install the agent on the active cluster (using KUBECONFIG env var)
resource "cloudcasa_kubecluster" "exmaple_kubecluster" {
  name = "cloudcasa_example_kubecluster"

  auto_install = true
}
