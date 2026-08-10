output "vpc_id" {
  description = "ID of the VPC."
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs. Used for both the control plane ENIs and the node group."
  value       = module.vpc.public_subnets
}

output "azs" {
  description = "Availability Zones in use."
  value       = local.azs
}

output "aws_region" {
  description = "Region, so the cluster layer can inherit it."
  value       = var.aws_region
}
