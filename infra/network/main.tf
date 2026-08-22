# EKS requires subnets in at least two AZs even though we only run one node.
data "aws_availability_zones" "available" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

locals {
  azs = slice(data.aws_availability_zones.available.names, 0, 3)
}

# node sits in a public subnet with a public IP and reaches the internet directly through the IGW
# security group (managed by the EKS module) allows no inbound traffic from the internet
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 6.6"

  name = "${var.project}-vpc"
  cidr = var.vpc_cidr

  azs = local.azs

  # /24 per AZ carved out of the /16: 10.0.0.0/24, 10.0.1.0/24
  public_subnets = [for i, az in local.azs : cidrsubnet(var.vpc_cidr, 8, i)]

  enable_nat_gateway = false

  # Without this the node launches with no public IP, cannot reach the
  # internet, and therefore never joins the cluster.
  map_public_ip_on_launch = true

  # Required for EKS: nodes resolve the cluster endpoint by DNS name.
  enable_dns_hostnames = true
  enable_dns_support   = true

  # the tag the ALB Controller looks for when you later add an Ingress.
  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }
}
