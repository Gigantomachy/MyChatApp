output "cluster_name" {
  description = "EKS cluster name."
  value       = module.eks.cluster_name
}

output "cluster_endpoint" {
  description = "Kubernetes API server endpoint."
  value       = module.eks.cluster_endpoint
}

output "node_security_group_id" {
  description = "Security group attached to the node. Useful when debugging connectivity."
  value       = module.eks.node_security_group_id
}

output "configure_kubectl" {
  description = "Run this to point kubectl at the new cluster."
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}

output "ecr_repository_urls" {
  description = "Passed through from the network layer for convenience when tagging images."
  value       = data.terraform_remote_state.network.outputs.ecr_repository_urls
}
