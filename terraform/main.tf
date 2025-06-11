# Main configuration file for Calendar API infrastructure

# Local values for resource naming and tagging
locals {
  name_prefix = "${var.project_name}-${var.environment}"

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
    Ephemeral   = var.is_ephemeral
  }
}
