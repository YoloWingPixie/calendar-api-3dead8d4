variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (e.g., prod, staging, pr-123)"
  type        = string
  validation {
    condition     = can(regex("^[a-z0-9-]+$", var.environment))
    error_message = "Environment must contain only lowercase letters, numbers, and hyphens"
  }
}

variable "is_ephemeral" {
  description = "Whether this is an ephemeral environment (e.g., PR environments)"
  type        = bool
  default     = false
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "calendar-api"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC (/28 is smallest allowed, providing 16 IPs)"
  type        = string
  default     = "10.0.0.0/28"
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
  default     = ["us-east-1a"]
}

variable "database_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "database_multi_az" {
  description = "Enable Multi-AZ for RDS (set to false for ephemeral environments)"
  type        = bool
  default     = false
}


variable "database_allocated_storage" {
  description = "Allocated storage for RDS in GB"
  type        = number
  default     = 20
}

variable "database_username" {
  description = "Master username for RDS"
  type        = string
  default     = "calendar_admin"
  sensitive   = true
}

variable "ecs_task_cpu" {
  description = "CPU units for ECS task (256 = 0.25 vCPU - minimum for Fargate)"
  type        = string
  default     = "256"
}

variable "ecs_task_memory" {
  description = "Memory for ECS task in MB (512 = 0.5 GB - minimum for Fargate)"
  type        = string
  default     = "512"
}

variable "app_port" {
  description = "Port the application listens on"
  type        = number
  default     = 8000
}

variable "min_tasks" {
  description = "Minimum number of ECS tasks"
  type        = number
  default     = 1
}

variable "max_tasks" {
  description = "Maximum number of ECS tasks"
  type        = number
  default     = 1
}

variable "doppler_token" {
  description = "Doppler service token for Terraform provider (read-only)"
  type        = string
  sensitive   = true
}

variable "docker_image_tag" {
  description = "Docker image tag to deploy"
  type        = string
  default     = "latest"
}

variable "pr_number" {
  description = "PR number for ephemeral environments"
  type        = number
  default     = null
}
