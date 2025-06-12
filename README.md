# Calendar API

[![CI](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml/badge.svg)](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml)

A centralized backend REST API service for calendar and event management, designed to replace fragmented legacy calendar tools within the organization.

## Overview

This project provides a RESTful API for managing users, calendars, and events with robust access control and modern cloud-native architecture.

## Technology Stack

- **Language**: Go 1.24+
- **Framework**: Gorilla Mux
- **Database**: PostgreSQL with database/sql
- **Package Manager**: Go modules
- **Infrastructure**: AWS (ECS, RDS, ALB), managed with Terraform
- **CI/CD**: GitHub Actions
- **Task Runner**: Taskfile

## Prerequisites

- Go 1.24+
- [Task](https://taskfile.dev/) (for development commands)
- Docker (for containerized development)
- PostgreSQL (for local development)
- Doppler CLI (for secrets management)

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd calendar-api

# Install dependencies
task mod

# Run default task
task

# Build the application
task build

# Run locally (requires database and environment setup)
task dev:local
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

# Run development server
task dev:local

# Run in Docker
task dev
```

## Environment Setup

### Local Development
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
For Docker Compose development, set these environment variables:

```bash
# Database Configuration (Docker Compose)
POSTGRES_USER=calendar_user
POSTGRES_PASSWORD=calendar_pass
POSTGRES_DB=calendar_db

# Application Configuration
DEBUG=true
BOOTSTRAP_ADMIN_KEY=dev-test-key-123
```

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

**Bootstrap Admin Key**: When deployed with Doppler, the `TF_VAR_bootstrap_admin_key` can be used to access all endpoints:

```bash
curl -H "X-API-Key: KD5nNdhoFuDRmdwZOh4An61QsrUiojYX" http://localhost:8080/api/events
```

## Documentation

- [Project Layout](docs/project-layout.md)
- [OpenAPI Specification](docs/openapi.yaml)
- [Architecture Decision Records](docs/adr/)
- [Domain Model](docs/domain-model.md)
- [Data Model](docs/data-model.md)

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
