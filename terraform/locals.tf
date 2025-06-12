# Central configuration for environment logic
# This file serves as the single source of truth for environment configuration

locals {
  # Canonical environments - these have their own Doppler configs
  canonical_environments = ["dev", "stg", "staging", "prod"]

  # Determine if current environment is canonical
  is_canonical_environment = contains(local.canonical_environments, var.environment)

  # Doppler config to use (non-canonical environments use dev)
  doppler_config = var.environment == "prod" ? "prd" : (local.is_canonical_environment ? var.environment : "dev")

  # AWS Secrets Manager config name to use (matches var.environment)
  aws_sm_config_name = local.is_canonical_environment ? var.environment : "dev"

  # Environment configuration
  environment_config = {
    # Debug mode is enabled for all non-prod environments
    debug = var.environment != "prod"

    # Ephemeral environments have shorter retention periods
    log_retention_days = var.is_ephemeral ? 1 : 7

    # Database backup settings
    backup_retention_days = var.is_ephemeral ? 0 : (var.environment == "prod" ? 30 : 7)

    # Deletion protection only for production
    deletion_protection = var.environment == "prod"
  }
}
