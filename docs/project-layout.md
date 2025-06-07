# Project Layout

This document describes the recommended folder structure for the Calendar API project, following modern FastAPI best practices for 2024.

## Overview

The project follows a domain-driven, three-layer architecture that separates concerns and promotes maintainability. The structure is designed to scale with the application while keeping related functionality grouped together.

## Directory Structure

```
calendar-api/
├── src/                           # Application source code
│   ├── __init__.py
│   ├── main.py                    # FastAPI app entry point
│   ├── config.py                  # Settings management with pydantic-settings
│   ├── database.py                # Async SQLAlchemy engine & session factory
│   ├── dependencies.py            # Common dependencies (get_db, get_current_user)
│   ├── exceptions.py              # Custom exception classes
│   ├── middleware.py              # Authentication & request middleware
│   │
│   ├── api/                       # API layer - HTTP interface
│   │   ├── __init__.py
│   │   ├── v1/                    # API version 1
│   │   │   ├── __init__.py
│   │   │   ├── users.py           # User endpoints (/api/v1/users)
│   │   │   ├── calendars.py       # Calendar endpoints (/api/v1/calendars)
│   │   │   └── events.py          # Event endpoints (/api/v1/events)
│   │   └── deps.py                # API-specific dependencies
│   │
│   ├── core/                      # Core business utilities
│   │   ├── __init__.py
│   │   ├── security.py            # Authentication, API key validation, hashing
│   │   └── constants.py           # Application-wide constants
│   │
│   ├── models/                    # SQLAlchemy ORM models
│   │   ├── __init__.py
│   │   ├── base.py                # Base model class with common fields
│   │   ├── user.py                # User table definition
│   │   ├── calendar.py            # Calendar table definition
│   │   └── event.py               # CalendarEvent table definition
│   │
│   ├── schemas/                   # Pydantic models for validation
│   │   ├── __init__.py
│   │   ├── user.py                # User request/response models
│   │   ├── calendar.py            # Calendar request/response models
│   │   ├── event.py               # Event request/response models
│   │   └── common.py              # Shared schemas (ErrorResponse, Pagination)
│   │
│   └── services/                  # Business logic layer
│       ├── __init__.py
│       ├── user.py                # User business operations
│       ├── calendar.py            # Calendar business operations
│       └── event.py               # Event business operations
│
├── tests/                         # Test suite
│   ├── __init__.py
│   ├── conftest.py                # Pytest fixtures and test configuration
│   ├── test_config.py             # Test-specific settings
│   ├── api/                       # API endpoint tests
│   │   ├── __init__.py
│   │   ├── test_users.py          # User endpoint tests
│   │   ├── test_calendars.py      # Calendar endpoint tests
│   │   └── test_events.py         # Event endpoint tests
│   ├── services/                  # Business logic tests
│   │   ├── __init__.py
│   │   ├── test_user_service.py   # User service tests
│   │   ├── test_calendar_service.py # Calendar service tests
│   │   └── test_event_service.py  # Event service tests
│   └── models/                    # Model and database tests
│       ├── __init__.py
│       └── test_models.py         # ORM model tests
│
├── alembic/                       # Database migrations
│   ├── versions/                  # Migration files (YYYY-MM-DD_descriptive_name.py)
│   ├── alembic.ini               # Alembic configuration
│   ├── env.py                    # Migration environment setup
│   └── script.py.mako            # Migration template
│
├── terraform/                     # Infrastructure as Code (OpenTofu)
│   ├── environments/              # Environment-specific configs
│   │   ├── dev/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   └── terraform.tfvars
│   │   ├── staging/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   └── terraform.tfvars
│   │   └── prod/
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       └── terraform.tfvars
│   ├── modules/                   # Reusable terraform modules
│   │   ├── ecs/                   # ECS Fargate configuration
│   │   ├── rds/                   # RDS PostgreSQL setup
│   │   ├── networking/            # VPC, subnets, security groups
│   │   └── alb/                   # Application Load Balancer
│   └── backend.tf                 # Terraform state backend config
│
├── docker/                        # Docker configuration
│   ├── Dockerfile                 # Multi-stage build for app
│   ├── docker-compose.yml         # Local development environment
│   └── .dockerignore              # Build exclusions
│
├── .github/                       # GitHub configuration
│   └── workflows/                 # GitHub Actions CI/CD
│       ├── ci.yml                 # Continuous Integration (test, lint, type check)
│       ├── build.yml              # Build and push Docker images
│       ├── deploy.yml             # Reusable deployment workflow
│       └── cd.yml                 # Continuous Deployment (orchestrates build + deploy)
│
├── scripts/                       # Utility scripts
│   ├── init_db.py                # Database initialization for development
│   └── seed_data.py              # Development data seeding
│
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore patterns
├── pyproject.toml                 # Project dependencies and metadata
├── uv.lock                        # Locked dependency versions
├── Taskfile.yml                   # Task runner configuration
├── README.md                      # Project overview
├── CLAUDE.md                      # Claude Code guidance
└── docs/                          # Project documentation
    ├── project-layout.md          # This file
    ├── prd.md                     # Product requirements
    ├── domain-model.md            # Business domain model
    ├── data-model.md              # Database schema
    ├── openapi.yaml               # API specification
    └── adr/                       # Architecture decision records
```

## Core Components

### Source Code (`src/`)

The main application code follows a layered architecture:

#### API Layer (`api/`)
- **Purpose**: Handle HTTP requests and responses
- **Responsibilities**:
  - Route definition using FastAPI's APIRouter
  - Request validation via Pydantic models
  - Response serialization
  - HTTP status codes and error handling
- **Key Principle**: No business logic - delegates to service layer

#### Service Layer (`services/`)
- **Purpose**: Implement business logic and orchestrate operations
- **Responsibilities**:
  - Business rule enforcement
  - Transaction management
  - Cross-entity operations
  - Complex queries and data aggregation
- **Key Principle**: Framework-agnostic business logic

#### Data Layer (`models/`)
- **Purpose**: Define database schema using SQLAlchemy ORM
- **Responsibilities**:
  - Table definitions
  - Relationships between entities
  - Database constraints
  - Model-level validations
- **Key Principle**: Pure data representation, no business logic

#### Schema Layer (`schemas/`)
- **Purpose**: Define data contracts for API communication
- **Responsibilities**:
  - Request body validation
  - Response serialization models
  - Data transformation between API and internal formats
- **Key Principle**: Separate from ORM models for flexibility

#### Core Utilities (`core/`)
- **Purpose**: Shared utilities and cross-cutting concerns
- **Contents**:
  - Security functions (password hashing, API key validation)
  - Application constants
  - Common utilities

### Tests (`tests/`)

The test structure mirrors the source code for easy navigation:

- **API Tests**: Test HTTP endpoints, status codes, and response formats
- **Service Tests**: Test business logic in isolation
- **Model Tests**: Test database constraints and relationships
- **Integration Tests**: Test full request-response cycles

### Database Migrations (`alembic/`)

- **Naming Convention**: `YYYY-MM-DD_descriptive_slug.py`
- **Version Control**: All schema changes tracked via migrations
- **Async Support**: Configured for async SQLAlchemy

### Infrastructure as Code (`terraform/`)

Manages AWS infrastructure using OpenTofu (Terraform-compatible):

- **Environment Separation**: Separate configurations for dev/staging/prod
- **Modular Design**: Reusable modules for common infrastructure patterns
- **Components**:
  - ECS Fargate for application hosting
  - RDS PostgreSQL for database
  - ALB for load balancing
  - VPC and networking configuration

### Docker Configuration (`docker/`)

- **Dockerfile**: Multi-stage build following ADR-007 guidelines
- **docker-compose.yml**: Local development environment with database
- **.dockerignore**: Optimize build context

### CI/CD Pipeline (`.github/workflows/`)

GitHub Actions workflows implementing trunk-based development:

- **ci.yml**: Continuous Integration - runs on every push
  - Run tests (unit and integration)
  - Run linter (ruff)
  - Type checking
  - Security scanning

- **build.yml**: Docker image building
  - Multi-stage Docker build
  - Push to Amazon ECR
  - Tag with git SHA and environment

- **deploy.yml**: Reusable deployment workflow
  - Takes environment as input
  - Updates ECS service
  - Runs database migrations
  - Health checks

- **cd.yml**: Continuous Deployment orchestrator
  - Triggers on main branch push
  - Calls build.yml
  - Calls deploy.yml for each environment (dev → staging → prod)
  - Implements deployment gates between environments

- **Integration**: Doppler for secrets management

## Development Workflow

1. **New Feature Development**:
   - Define Pydantic schemas in `schemas/`
   - Create/update SQLAlchemy models in `models/`
   - Implement business logic in `services/`
   - Add API endpoints in `api/v1/`
   - Write tests in corresponding `tests/` directories

2. **Database Changes**:
   - Modify models in `models/`
   - Generate migration: `task db:revision -- "description"`
   - Review and run migration: `task db:migrate`

3. **Infrastructure Changes**:
   - Update terraform modules or environment configs
   - Plan changes: `task terraform:plan -- <env>`
   - Apply changes: `task terraform:apply -- <env>`

4. **Testing**:
   - Unit tests for individual components
   - Integration tests for API endpoints
   - Run all tests: `task test`

5. **Local Development**:
   - Start services: `task dev`
   - Database setup: `task db:setup`
   - Seed data: `task db:seed`

## Best Practices

1. **Separation of Concerns**:
   - Keep API layer thin - validation and serialization only
   - Business logic belongs in service layer
   - Database queries through service layer or repositories

2. **Dependency Injection**:
   - Use FastAPI's dependency system
   - Database sessions via `get_db` dependency
   - Current user via `get_current_user` dependency

3. **Async Patterns**:
   - Use async/await throughout the stack
   - Async SQLAlchemy sessions
   - Avoid blocking operations in async routes

4. **Error Handling**:
   - Define custom exceptions in `exceptions.py`
   - Consistent error response format
   - Proper HTTP status codes

5. **Configuration**:
   - All settings in `config.py` using pydantic-settings
   - Environment variables for secrets
   - `.env.example` as template

6. **Code Organization**:
   - One class/function per file when appropriate
   - Group related functionality
   - Clear, descriptive naming

7. **Infrastructure**:
   - Keep terraform DRY with modules
   - Environment parity where possible
   - Document any environment-specific differences

This structure provides a solid foundation for building and maintaining the Calendar API while following FastAPI best practices and supporting cloud-native deployment on AWS.
