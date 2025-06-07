# Calendar API - Product Requirements

## 1. Purpose

This document specifies the requirements for a new internal Calendar API. The objective is to create a backend service for calendar and event data, providing a single source of truth to replace fragmented legacy tools.

## 2. Key Outcomes

*   A RESTful API providing CRUD operations for users, calendars, and events.
*   All underlying infrastructure is defined as code using Terraform.
*   A CI/CD pipeline automates builds, tests, and deployments on merge to `main`.
*   Application logs are centralized in CloudWatch.
*   API contract is defined and discoverable via an OpenAPI specification.
*   Secrets are managed via AWS Secrets Manager.

## 3. Out of Scope (v1)

*   No frontend UI.
*   Authentication is limited to a static API key. No OAuth2/OIDC.
*   No support for recurring events, calendar invitations, or email notifications.
*   Single-region deployment only. No multi-region HA.

## 4. Technical Requirements

### 4.1 Repository & Project Setup

*   **TR-4.1.1:** Initialize Git repository.
*   **TR-4.1.2:** Adopt branching strategy from `adr-009-git-branching-strategy.md`.
*   **TR-4.1.3:** Create directory structure: `src/`, `iac/`, `.github/workflows/`, `docs/`.
*   **TR-4.1.4:** Add standard Python `.gitignore`.
*   **TR-4.1.5:** Create `README.md` with project overview and developer setup instructions.

### 4.2 Application & API

*   **TR-4.2.1: Setup**
    *   **TR-4.2.1.1:** Initialize FastAPI application.
    *   **TR-4.2.1.2:** Set up dependency management with Poetry.
    *   **TR-4.2.1.3:** Implement `GET /health` endpoint returning `{"status": "ok", "version": "...", "timestamp": "..."}`.
*   **TR-4.2.2: Database**
    *   **TR-4.2.2.1:** Create reference `schema.sql` based on `docs/data-model.md`.
    *   **TR-4.2.2.2:** Implement connection logic for PostgreSQL.
    *   **TR-4.2.2.3:** Implement SQLAlchemy ORM models matching the schema.
    *   **TR-4.2.2.4:** Set up Alembic for schema migrations and create initial migration.
*   **TR-4.2.3: Models & Auth**
    *   **TR-4.2.3.1:** Implement Pydantic DTOs from `docs/data-model.md`.
    *   **TR-4.2.3.2:** Implement API key auth dependency (`X-API-Key` header).
*   **TR-4.2.4: Endpoints**
    *   **TR-4.2.4.1:** Implement `POST /users`.
    *   **TR-4.2.4.2:** Implement `POST /calendars`, `GET /calendars`, `GET /calendars/{id}`, `PATCH /calendars/{id}`, `DELETE /calendars/{id}`.
    *   **TR-4.2.4.3:** Implement `POST`, `GET`, `PATCH`, `DELETE` endpoints under `/calendars/{calendar_id}/events`.
    *   **TR-4.2.4.4:** Implement legacy-compatible endpoints: `POST /events`, `GET /events`, `GET /events/{id}`, `PUT /events/{id}`, `DELETE /events/{id}`.

### 4.3 Containerization

*   **TR-4.3.1:** Write a multi-stage `Dockerfile` for production builds.
*   **TR-4.3.2:** Add `.dockerignore` file.
*   **TR-4.3.3:** Create `docker-compose.yml` for local development (API + DB services).

### 4.4 Infrastructure (Terraform)

*   **TR-4.4.1:** Initialize Terraform project with a remote state backend (S3).
*   **TR-4.4.2:** Define network resources: VPC, public/private subnets, NAT Gateway, IGW.
*   **TR-4.4.3:** Define security groups for ALB, ECS service, and RDS instance.
*   **TR-4.4.4:** Provision a multi-AZ PostgreSQL RDS instance in private subnets.
*   **TR-4.4.5:** Provision AWS Secrets Manager secret for DB credentials.
*   **TR-4.4.6:** Provision ECR repository for the application image.
*   **TR-4.4.7:** Provision an Application Load Balancer.
*   **TR-4.4.8:** Provision an ECS cluster and Fargate service definition.
*   **TR-4.4.9:** Provision a CloudWatch Log Group for the ECS service.

### 4.5 CI/CD (GitHub Actions)

*   **TR-4.5.1: Testing**
    *   **TR-4.5.1.1:** Implement unit tests for core logic.
    *   **TR-4.5.1.2:** Implement integration tests for all API endpoints.
*   **TR-4.5.2: CI Workflow (on PR)**
    *   **TR-4.5.2.1:** Add jobs for linting, running tests, and building the Docker image.
*   **TR-4.5.3: CD Workflow (on merge to `main`)**
    *   **TR-4.5.3.1:** Add jobs for pushing image to ECR, running `terraform apply`, and forcing new ECS service deployment.

### 4.6 Documentation

*   **TR-4.6.1:** Maintain `docs/openapi.yaml` as the API contract.
*   **TR-4.6.2:** Create `docs/Architecture.md` with a Mermaid diagram.

## 5. Service Level Objectives (SLOs) & Engineering Standards

*   **SLO-5.1: Performance:** P95 latency for API requests must be < 500ms.
*   **SLO-5.2: Availability:** Target 99.9% uptime. Errors must conform to the `Error` schema in the OpenAPI spec.
*   **STD-5.3: Maintainability:** Code must pass linter checks before merge. All infrastructure changes must be made via Terraform.
*   **STD-5.4: Security:** Require TLS 1.2+ for external traffic. Source secrets from AWS Secrets Manager. IAM roles and Security Groups must be least-privilege.
*   **STD-5.5: Developer Experience:** A new developer must be able to stand up a local environment by following the `README.md`.
