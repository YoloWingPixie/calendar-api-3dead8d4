# Terraform backend configuration will be added later
# This will use S3 for remote state storage

terraform {
  backend "s3" {
    bucket       = "calendar-api-terraform-state-655593807337"
    key          = "calendar-api/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true  # S3 native locking (2024 best practice)
  }
}
