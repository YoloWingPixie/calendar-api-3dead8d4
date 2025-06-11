# Random password for RDS (only used for initial creation)
# Actual password management will be through Doppler
resource "random_password" "db_password" {
  length  = 32
  special = true
  # RDS doesn't allow these characters: / @ " space
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}"
  subnet_ids = aws_subnet.public[*].id

  tags = {
    Name = "${var.project_name}-${var.environment}-db-subnet-group"
  }
}

# RDS PostgreSQL Instance
resource "aws_db_instance" "main" {
  identifier = "${var.project_name}-${var.environment}"

  # Engine
  engine         = "postgres"
  engine_version = "16"

  # Instance
  instance_class    = var.database_instance_class
  allocated_storage = var.database_allocated_storage
  storage_type      = "gp3"
  storage_encrypted = true

  # Database
  db_name  = "calendar_db"
  username = var.database_username
  password = random_password.db_password.result

  # Network
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = true # Using security groups for access control

  # Availability
  multi_az = var.database_multi_az

  # Backup
  backup_retention_period = local.environment_config.backup_retention_days
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"

  # Performance Insights (disabled for cost savings)
  performance_insights_enabled = false

  # Other settings
  auto_minor_version_upgrade = true
  skip_final_snapshot       = var.is_ephemeral
  final_snapshot_identifier = var.is_ephemeral ? null : "${var.project_name}-${var.environment}-final-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  deletion_protection       = local.environment_config.deletion_protection

  tags = {
    Name        = "${var.project_name}-${var.environment}-rds"
    Environment = var.environment
    Ephemeral   = var.is_ephemeral
  }
}

# Store connection info in Parameter Store for Doppler to read
resource "aws_ssm_parameter" "db_host" {
  name  = "/${var.project_name}/${var.environment}/db/host"
  type  = "String"
  value = aws_db_instance.main.address

  tags = {
    Name        = "${var.project_name}-${var.environment}-db-host"
    Environment = var.environment
  }
}

resource "aws_ssm_parameter" "db_port" {
  name  = "/${var.project_name}/${var.environment}/db/port"
  type  = "String"
  value = aws_db_instance.main.port

  tags = {
    Name        = "${var.project_name}-${var.environment}-db-port"
    Environment = var.environment
  }
}

# Note: The actual database password should be managed through Doppler
# This password is only for initial RDS creation
