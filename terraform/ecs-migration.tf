# ECS Task Definition for Database Migrations
resource "aws_ecs_task_definition" "migration" {
  family                   = "${var.project_name}-${var.environment}-migration"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn           = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
    name  = "migration"
    image = var.container_image_uri != "" ? var.container_image_uri : "${data.aws_ecr_repository.app.repository_url}:latest"

    # Override command to run migrations
    command = ["alembic", "upgrade", "head"]

    # Environment variables
    environment = [
      {
        name  = "ENVIRONMENT"
        value = var.environment
      }
    ]

    # Inject secrets from AWS Secrets Manager (synced from Doppler)
    secrets = [
      {
        name      = "DOPPLER_SECRETS_JSON"
        valueFrom = data.aws_secretsmanager_secret_version.doppler_sync.arn
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.app.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "migration"
      }
    }

    essential = true
  }])
}

# Output for migration task definition ARN
output "migration_task_definition_arn" {
  description = "ARN of the migration task definition"
  value       = aws_ecs_task_definition.migration.arn
}
