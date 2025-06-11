# Data source to reference the shared ECR repository
# This should be created once outside of environment-specific Terraform
data "aws_ecr_repository" "app" {
  name = var.project_name
}

# Note: The actual ECR repository should be created in a separate
# "shared infrastructure" Terraform configuration that runs once
# for the entire project, not per environment
