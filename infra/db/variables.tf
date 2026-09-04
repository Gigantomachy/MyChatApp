variable "aws_region" {
  description = "AWS region. Must match the network and cluster layers."
  type        = string
  default     = "ca-central-1"
}

variable "project" {
  type    = string
  default = "mychatapp"
}

variable "instance_type" {
  description = "Cassandra node instance type (t3.medium = 4 GiB, ~$30/mo each)."
  type        = string
  default     = "t3.medium"
}

variable "node_ips" {
  description = "Fixed private IPs, one per AZ subnet (10.0.0/24, 10.0.1/24, 10.0.2/24)."
  type        = list(string)
  default     = ["10.0.0.10", "10.0.1.10", "10.0.2.10"]
}

variable "root_disk_size" {
  type    = number
  default = 30
}

variable "data_disk_size" {
  description = "Cassandra data volume per node. Stopped instances still pay for this."
  type        = number
  default     = 50
}

variable "ssh_cidr" {
  description = "CIDR allowed SSH access. Narrow to your home IP (e.g. 1.2.3.4/32) after first run."
  type        = string
  default     = "0.0.0.0/0"
}

variable "network_state_path" {
  type    = string
  default = "../network/terraform.tfstate"
}