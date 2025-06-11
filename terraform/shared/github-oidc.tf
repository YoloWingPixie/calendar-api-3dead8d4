# Data source to get current AWS account ID
data "aws_caller_identity" "current" {}

# GitHub OIDC Provider
resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]

  # GitHub's OIDC thumbprints (these are stable)
  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd"
  ]

  tags = {
    Name        = "github-actions-oidc"
    Description = "OIDC provider for GitHub Actions"
    Project     = var.project_name
  }
}

# IAM Role for GitHub Actions
resource "aws_iam_role" "github_actions" {
  name = "github-actions-${var.project_name}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            # Update these values with your GitHub org/username and repo
            # Format: "repo:OWNER/REPO:*"
            "token.actions.githubusercontent.com:sub" = [
              "repo:YoloWingPixie/calendar-api:*",
              "repo:YoloWingPixie/calendar-api:ref:refs/heads/*",
              "repo:YoloWingPixie/calendar-api:ref:refs/tags/*",
              "repo:YoloWingPixie/calendar-api:pull_request"
            ]
          }
        }
      }
    ]
  })

  tags = {
    Name        = "github-actions-${var.project_name}"
    Description = "Role for GitHub Actions to deploy ${var.project_name}"
    Project     = var.project_name
  }
}

# Policy for GitHub Actions - ECR Access
resource "aws_iam_role_policy" "github_actions_ecr" {
  name = "ecr-access"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:ListImages",
          "ecr:DescribeImages"
        ]
        Resource = "*"
      }
    ]
  })
}

# Policy for GitHub Actions - Terraform State Access
resource "aws_iam_role_policy" "github_actions_terraform_state" {
  name = "terraform-state-access"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "s3:GetBucketVersioning"
        ]
        Resource = "arn:aws:s3:::calendar-api-terraform-state-${data.aws_caller_identity.current.account_id}"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "arn:aws:s3:::calendar-api-terraform-state-${data.aws_caller_identity.current.account_id}/*"
      }
    ]
  })
}

# Policy for GitHub Actions - Terraform Deployment Permissions
resource "aws_iam_role_policy" "github_actions_terraform_deploy" {
  name = "terraform-deploy"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      # VPC and Networking
      {
        Effect = "Allow"
        Action = [
          "ec2:*"
        ]
        Resource = "*"
      },
      # ECS
      {
        Effect = "Allow"
        Action = [
          "ecs:*"
        ]
        Resource = "*"
      },
      # RDS
      {
        Effect = "Allow"
        Action = [
          "rds:*"
        ]
        Resource = "*"
      },
      # IAM (limited to specific roles/policies)
      {
        Effect = "Allow"
        Action = [
          "iam:CreateRole",
          "iam:DeleteRole",
          "iam:AttachRolePolicy",
          "iam:DetachRolePolicy",
          "iam:PutRolePolicy",
          "iam:DeleteRolePolicy",
          "iam:GetRole",
          "iam:GetRolePolicy",
          "iam:ListRolePolicies",
          "iam:ListAttachedRolePolicies",
          "iam:PassRole",
          "iam:TagRole",
          "iam:UntagRole",
          "iam:ListInstanceProfilesForRole"
        ]
        Resource = [
          "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.project_name}-*",
          "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/aws-service-role/*"
        ]
      },
      # ALB
      {
        Effect = "Allow"
        Action = [
          "elasticloadbalancing:*"
        ]
        Resource = "*"
      },
      # CloudWatch Logs
      {
        Effect = "Allow"
        Action = [
          "logs:*"
        ]
        Resource = "*"
      },
      # Secrets Manager (read-only for Doppler secrets)
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret",
          "secretsmanager:ListSecrets"
        ]
        Resource = "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:/calendar-api/*"
      },
      # Route53 (if using custom domains)
      {
        Effect = "Allow"
        Action = [
          "route53:*"
        ]
        Resource = "*"
      },
      # ACM (for SSL certificates)
      {
        Effect = "Allow"
        Action = [
          "acm:*"
        ]
        Resource = "*"
      }
    ]
  })
}

# Output the role ARN for use in GitHub Secrets
output "github_actions_role_arn" {
  value       = aws_iam_role.github_actions.arn
  description = "ARN of the IAM role for GitHub Actions. Add this to GitHub secrets as AWS_ROLE_ARN"
}

output "github_oidc_provider_arn" {
  value       = aws_iam_openid_connect_provider.github.arn
  description = "ARN of the GitHub OIDC provider"
}
