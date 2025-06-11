# Outputs for the Calendar API infrastructure

output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "ecr_repository_url" {
  description = "URL of the ECR repository"
  value       = data.aws_ecr_repository.app.repository_url
}

output "database_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "database_port" {
  description = "RDS instance port"
  value       = aws_db_instance.main.port
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.app.name
}

output "cloudwatch_log_group" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.app.name
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "environment" {
  description = "Environment name"
  value       = var.environment
}

output "doppler_secret_name" {
  description = "AWS Secrets Manager secret name for Doppler sync"
  value       = "/doppler/3dead8d4/${var.environment}"
}

output "doppler_secrets_updated" {
  description = "Doppler secrets have been automatically updated with RDS details"
  value       = "DATABASE_URL and related secrets have been updated in Doppler project '3dead8d4' config '${local.doppler_config}'"
}

output "environment_info" {
  description = "Environment configuration information"
  value = {
    environment              = var.environment
    is_canonical            = local.is_canonical_environment
    doppler_config_used     = local.doppler_config
    canonical_environments  = local.canonical_environments
  }
}
