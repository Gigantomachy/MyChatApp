terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 6.52"
    }
  }
}

provider "aws" {
  region = var.aws_region

  # applied to every resource this layer creates, to find and cost-attribute everything with a single tag filter
  default_tags {
    tags = {
      Project   = var.project
      ManagedBy = "terraform"
      Layer     = "network"
    }
  }
}
