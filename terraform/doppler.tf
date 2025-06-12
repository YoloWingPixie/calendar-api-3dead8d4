# Doppler to AWS Secrets Manager Integration

# Data source to reference the synced secret from Doppler
# Doppler should be configured to sync to AWS Secrets Manager using the "Single-Secret JSON" strategy
# The secret name pattern in AWS SM should be: /doppler/3dead8d4/<environment>
# Non-canonical environments use the dev Doppler config (defined in locals.tf)
data "aws_secretsmanager_secret" "doppler_sync" {
  name = "/calendar-api/${local.aws_sm_config_name}/"
}

# Get the current version of the secret
data "aws_secretsmanager_secret_version" "doppler_sync" {
  secret_id = data.aws_secretsmanager_secret.doppler_sync.id
}

# IAM policy for ECS task execution role to access the Doppler-synced secret
resource "aws_iam_role_policy" "ecs_secrets_access" {
  name = "${local.name_prefix}-ecs-secrets-access"
  role = aws_iam_role.ecs_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = data.aws_secretsmanager_secret.doppler_sync.arn
      }
    ]
  })
}

# Automatic Doppler secret updates after RDS creation
# These resources automatically update Doppler with RDS connection details

# Update DATABASE_URL in Doppler after RDS is created
resource "doppler_secret" "database_url" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_URL"
  value   = "postgresql://${var.database_username}:${random_password.db_password.result}@${aws_db_instance.main.endpoint}/${aws_db_instance.main.db_name}"

  # Ensure this only runs after RDS is fully created
  depends_on = [aws_db_instance.main]
}

# Also store individual components for flexibility
resource "doppler_secret" "database_host" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_HOST"
  value   = aws_db_instance.main.address
}

resource "doppler_secret" "database_port" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_PORT"
  value   = tostring(aws_db_instance.main.port)
}

resource "doppler_secret" "database_password" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_PASSWORD"
  value   = random_password.db_password.result
}

resource "doppler_secret" "database_name" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_NAME"
  value   = aws_db_instance.main.db_name
}

resource "doppler_secret" "database_username" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DATABASE_USERNAME"
  value   = var.database_username
}

# Environment-specific settings
resource "doppler_secret" "environment" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "ENVIRONMENT"
  value   = var.environment
}

resource "doppler_secret" "debug" {
  project = "3dead8d4"
  config  = local.doppler_config
  name    = "DEBUG"
  value   = tostring(local.environment_config.debug)
}
