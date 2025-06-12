# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete Go migration system with version tracking and detailed logging
- Doppler secrets integration via DOPPLER_SECRETS_JSON parsing
- API key authentication middleware with bootstrap admin key support
- User management system with database storage

### Changed
- Migrated from Python FastAPI to Go with Gorilla Mux
- Replaced Alembic with custom Go migration system
- Updated all GitHub workflows for Go development
- Simplified deployment: migrations now run automatically on application startup
- Removed separate ECS migration task definition (no longer needed)

### Removed
- Alembic migration system and configuration files
- Python-specific dependencies and configuration
- RUN_MIGRATIONS_ONLY environment variable and migration-only mode
- Separate ECS migration task definition in Terraform

### Added
- SQLAlchemy ORM models for User, Calendar, and CalendarEvent entities
- Pydantic schemas for all API request/response validation
- Database connection and session management setup
- Alembic database migration configuration
- Initial database migration script (001_initial_schema.py)
- Automated database migrations in GitHub Actions deployment workflow
- Database-related Taskfile commands (db:migrate, db:revision, db:downgrade, db:history, db:current)

### Changed
- Updated Calendar model to include missing `public_write` field
- Updated CalendarEvent model to use correct field names (`title` instead of `event_name`)
- Added `creator_user_id` and `is_all_day` fields to CalendarEvent model
- Enhanced deploy-common.yml workflow to automatically run Alembic migrations after Terraform deployment

### Deprecated
- None

### Removed
- Redundant /src/database/ directory and schema.sql file (using Alembic for all migrations)

### Fixed
- None

### Security
- None

## [0.2.2] - 2025-06-11
### Added
- Initial project setup
- Basic infrastructure provisioning
- CI/CD pipeline foundation
- Initial project setup with FastAPI
- Basic project structure following clean architecture principles
- Health check endpoint (`GET /api/v1/health`)
- CI/CD pipeline with GitHub Actions
- Infrastructure as Code using Terraform
  - VPC and networking setup
  - RDS PostgreSQL instance
  - ECS Fargate service
  - Application Load Balancer
  - Security groups and IAM roles
  - ECR repository
  - CloudWatch log group
- Project documentation
  - Product Requirements Document (PRD)
  - Data model specification
  - Domain model documentation
  - OpenAPI specification
  - Project layout documentation
  - Architecture Decision Records (ADRs)
