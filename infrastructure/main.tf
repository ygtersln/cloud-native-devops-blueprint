terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    # Used for OpenTofu / Kubernetes
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
  
  # Recommended for production: S3 backend with DynamoDB locking
  # backend "s3" {
  #   bucket         = "devops-blueprint-tfstate"
  #   key            = "global/s3/terraform.tfstate"
  #   region         = "eu-central-1"
  #   dynamodb_table = "terraform-locks"
  #   encrypt        = true
  # }
}

provider "aws" {
  region = var.aws_region
}

# Example: Calling the EKS Module
# module "eks" {
#   source = "./modules/aws-eks"
#   cluster_name = var.cluster_name
#   vpc_id = module.vpc.vpc_id
#   subnet_ids = module.vpc.private_subnets
# }
