# Calendar API

[![CI](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml/badge.svg)](https://github.com/YoloWingPixie/calendar-api/actions/workflows/ci.yml)

A centralized backend REST API service for calendar and event management, designed to replace fragmented legacy calendar tools within the organization.

## Overview

This project provides a RESTful API for managing users, calendars, and events with robust access control and modern cloud-native architecture.

## Technology Stack

- **Language**: Python 3.13+
- **Framework**: FastAPI
- **Database**: PostgreSQL with SQLAlchemy ORM and Alembic migrations
- **Package Manager**: uv
- **Infrastructure**: AWS (ECS, RDS, ALB), managed with Terraform
- **CI/CD**: GitHub Actions

## Prerequisites

- Python 3.13+
- [uv](https://docs.astral.sh/uv/) package manager
- [Task](https://taskfile.dev/) (for development commands)
- Docker (for containerized development)
- PostgreSQL (for local development)

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd calendar-api

# Install dependencies
uv sync

# Run default task
task
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
