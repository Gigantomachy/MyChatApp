variable "aws_region" {
  description = "AWS region. Must match the network layer."
  type        = string
  default     = "ca-central-1"
}

variable "project" {
  description = "Short name, also used as the cluster name."
  type        = string
  default     = "mychatapp"
}

variable "kubernetes_version" {
  description = <<-EOT
    EKS Kubernetes minor version.

    1.34 is a deliberate pick: 1.36 is the newest but 1.33 leaves standard
    support on 2026-08-31, which is soon. 1.34 gives plenty of runway without
    being bleeding edge.
  EOT
  type        = string
  default     = "1.34"
}

variable "node_instance_type" {
  description = <<-EOT
    Instance type for the single node.

    t3.large (2 vCPU / 8 GiB) is close to the practical floor: Cassandra,
    the Go backend, the nginx frontend, and the kube-system pods all share it.
    t3.medium (4 GiB) will make Cassandra unhappy even with a capped heap.
  EOT
  type        = string
  default     = "t3.large"
}

variable "node_disk_size" {
  description = "Root EBS volume size (GiB) for the node."
  type        = number
  default     = 30
}

variable "network_state_path" {
  description = "Relative path to the network layer's local state file."
  type        = string
  default     = "../network/terraform.tfstate"
}
