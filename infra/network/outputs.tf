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

output "ecr_repository_urls" {
  description = "Map of repo name to push URL."
  value       = { for k, v in aws_ecr_repository.app : k => v.repository_url }
}

output "ecr_registry" {
  description = "Registry hostname, for `docker login`."
  value       = split("/", values(aws_ecr_repository.app)[0].repository_url)[0]
}
