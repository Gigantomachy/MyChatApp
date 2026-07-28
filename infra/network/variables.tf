variable "aws_region" {
  description = "AWS region to deploy into. Must match the cluster layer."
  type        = string
  default     = "ca-central-1"
}

variable "project" {
  description = "Short name prefixed onto resource names."
  type        = string
  default     = "mychatapp"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}
