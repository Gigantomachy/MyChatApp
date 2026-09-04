data "terraform_remote_state" "network" {
  backend = "local"

  config = {
    path = var.network_state_path
  }
}

locals {
  vpc_id     = data.terraform_remote_state.network.outputs.vpc_id
  subnet_ids = data.terraform_remote_state.network.outputs.public_subnet_ids
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.24"

  name               = var.project
  kubernetes_version = var.kubernetes_version

  vpc_id     = local.vpc_id
  subnet_ids = local.subnet_ids

  # Needed for kubectl from local machine. With no NAT gateway and no VPN, this is also how we reach the api
  endpoint_public_access = true

  # creates an EKS access entry granting cluster-admin to whichever IAM identity runs terraform apply. 
  # Without this we cant talk to our cluster.
  enable_cluster_creator_admin_permissions = true

  addons = {
    vpc-cni = {
      before_compute = true
    }

    eks-pod-identity-agent = {
      before_compute = true
    }

    coredns    = {}
    kube-proxy = {}
  }

  eks_managed_node_groups = {
    default = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = [var.node_instance_type]
      subnet_ids     = local.subnet_ids   # all public subnets = both AZs

      iam_role_additional_policies = {
        ssm = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
      }

      # One node per AZ (3 public subnets across 3 AZs).
      min_size     = 3
      max_size     = 3
      desired_size = 3

      block_device_mappings = {
        xvda = {
          device_name = "/dev/xvda"
          ebs = {
            volume_size           = var.node_disk_size
            volume_type           = "gp3"
            encrypted             = true
            delete_on_termination = true
          }
        }
      }
    }
  }

  node_security_group_additional_rules = {
    node_to_node_all_ports = {
      description = "Allow node-to-node on all TCP ports (kube-proxy NodePort DNAT targets pod ports like 80)"
      protocol    = "tcp"
      from_port   = 1
      to_port     = 65535
      type        = "ingress"
      self        = true
    }
  }
}
