# Calendar API

[![CI](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml/badge.svg)](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml)

A centralized backend REST API service for calendar and event management, designed to replace fragmented legacy calendar tools within the organization.

## Overview

This project provides a RESTful API for managing users, calendars, and events with robust access control and modern cloud-native architecture.

## Documentation

- [Architecture Overview](docs/architecture.md)
- [OpenAPI Specification](docs/openapi.yaml)
- [Swagger Specificaiton](/swagger.yaml)
- [Architecture Decision Records](docs/adr/)
- [Domain Model](docs/domain-model.md)
- [Assumptions](docs/assumptions.md)
- [Questions about the Development of This Project](/docs/FAQ.md)
- [Product Requirements Document derived from original assessment](/docs/prd.md)

## Technology Stack

- **Language**: Go 1.24+
- **Framework**: Gorilla Mux
- **Database**: PostgreSQL with database/sql
- **Package Manager**: Go modules
- **Infrastructure**: AWS (ECS, RDS, ALB), managed with Terraform
- **CI/CD**: GitHub Actions
- **Task Runner**: Taskfile

## Prerequisites

- [Task](https://taskfile.dev/installation/) - Task runner for development commands
- [Docker](https://docs.docker.com/get-docker/) - For containerized development
- [Doppler CLI](https://docs.doppler.com/docs/install-cli) - For secrets management

## Quick Start

1. Setup Doppler per [Doppler Configuration](#doppler-configuration)
2. [Create an AWS Account](https://docs.aws.amazon.com/accounts/latest/reference/manage-acct-creating.html)
3. Setup OIDC authenticaiton for Github to your AWS account per [Setting up OIDC to AWS](#setting-up-oidc-to-aws)
4. Install the Prerequisites

```bash
# Clone the repository
git clone https://github.com/YoloWingPixie/calendar-api-3dead8d4.git
cd calendar-api-3dead8d4

# Install dependencies
task setup

#Setup doppler
doppler login
doppler setup

# Build the application and run locally for development
task dev
```

## Development Commands

```bash
# Show all available tasks
task

# Download dependencies
task mod

# Run tests
task test

# Format code
task format

# Lint code
task lint

# Build application
task build

# Run the application locally
task dev

# Run mod, format, lint, test
task ci
```

## Environment Setup

### Local Development without Doppler
It is possible to run this project locally without Doppler, however task commands will not work properly (as they expect doppler), and this will not be valid for deployment to AWS:
Create a `.env` file with the following variables:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=calendar
DB_SSLMODE=disable

# Server Configuration
PORT=8080
HOST=0.0.0.0

# Application Configuration
DEBUG=true
BOOTSTRAP_ADMIN_KEY=dev-admin-key-123
```

### Docker Development
For Docker Compose development, use Doppler to manage secrets:

```bash
# Run with Doppler
doppler run -- docker compose up
```

The required secrets are managed through Doppler as described in the [Doppler Configuration](#doppler-configuration) section above.

### Setting up OIDC to AWS

1. **Add GitHub as an OIDC identity provider in AWS IAM**  
   IAM console → *Identity providers* → **Add provider** →  
   *Provider type* = **OpenID Connect** | *Provider URL* = `https://token.actions.githubusercontent.com` | *Audience* = `sts.amazonaws.com`  
   [AWS docs](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc.html)

2. **Create an IAM role that trusts that provider**  
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Principal": {
           "Federated": "arn:aws:iam::<ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
         },
         "Action": "sts:AssumeRoleWithWebIdentity",
         "Condition": {
           "StringEquals": {
             "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
             "token.actions.githubusercontent.com:sub": "repo:<OWNER>/<REPO>:ref:refs/heads/*"
           }
         }
       }
     ]
   }
```

Scope the `sub` claim as tightly as possible.
[AWS trust-policy example](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-idp_oidc.html)

3. **Attach the required permissions policy to that role**
   Grant only what the workflow needs (ECR, S3, CloudFormation, …).
   [AWS IAM policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_manage.html)

4. **Request the OIDC token in your workflow**

   ```yaml
   permissions:
     id-token: write   # OIDC token
     contents: read    # allow checkout
   ```

   [GitHub OIDC permissions](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect#adding-permissions-settings)

5. **Assume the role from the job**

   ```yaml
   - uses: actions/checkout@v4

   - uses: aws-actions/configure-aws-credentials@v4
     with:
       role-to-assume: arn:aws:iam::<ACCOUNT_ID>:role/<ROLE_NAME>
       aws-region: us-east-1
   ```

   [configure-aws-credentials action](https://github.com/aws-actions/configure-aws-credentials)

### Doppler Configuration

The application uses Doppler for secrets management. Here are the minimal required secrets that must be configured in Doppler:

```bash
# Required Database Configuration
DATABASE_USERNAME=<your-db-username>

# Required Security Configuration
BOOTSTRAP_ADMIN_KEY=<your-admin-key>
```

Optional secrets (with safe defaults):
```bash
# Server Configuration (defaults: host=0.0.0.0, port=8080)
HOST=<host>
PORT=<port>

# Database Configuration (defaults: host=localhost, port=5432, sslmode=require)
DATABASE_HOST=<db-host>
DATABASE_PORT=<db-port>

# Application Configuration
DEBUG=<true/false>
API_KEY_HEADER=<header-name>  # defaults to X-API-Key
```

Secrets not mentioned here but found in the config are forcibly updated by Doppler after Terraform applies, this includes things like the Database Host, Port, Password, Environment, and Full DB URL.

To set up Doppler:
1. Install Doppler CLI: https://docs.doppler.com/docs/install-cli
2. Login: `doppler login`
3. Setup project: `doppler setup`
4. Configure the required secrets above
5. Run the application with Doppler: `doppler run -- task dev`

For more information:
- [Doppler CLI Documentation](https://docs.doppler.com/docs/cli)
- [Doppler Secrets Management](https://docs.doppler.com/docs/secrets)
- [Doppler Environment Variables](https://docs.doppler.com/docs/environment-variables)
- [Doppler Terraform Integration](https://docs.doppler.com/docs/terraform-provider)

## API Endpoints

### Public Endpoints
- `GET /health` - Health check endpoint (no authentication required)

### Protected Endpoints (require API key)
- `GET /api/events` - List all events
- `POST /api/events` - Create a new event
- `GET /api/events/{id}` - Get event by ID
- `PUT /api/events/{id}` - Update event by ID
- `DELETE /api/events/{id}` - Delete event by ID

### Authentication

All `/api/*` endpoints require authentication via the `X-API-Key` header:

```bash
curl -H "X-API-Key: your-api-key-here" http://localhost:8080/api/events
```

**Bootstrap Admin Key**: When deployed with Doppler, the `BOOTSTRAP_ADMIN_KEY` can be used to access all endpoints:

```bash
curl -H "X-API-Key: IAMSOMERANDOMAPIKEY" http://localhost:8080/api/events
```

## Infrastructure

Terraform manages AWS infrastructure (ECS, RDS, ALB). Key points:

- **Environments**: dev/staging/prod are canonical. Others (PR deployments) use dev Doppler config
- **State**: S3 backend with state locking
- **Deployment**: `terraform apply -var="environment=dev"` from `terraform/` directory
- **Secrets**: Doppler syncs to AWS Secrets Manager. RDS credentials auto-update in Doppler post-deployment
- **Shared resources**: ECR repository created once via `terraform/shared/`

See `terraform/locals.tf` for environment configuration logic.

## License

This is an internal project. All rights reserved.