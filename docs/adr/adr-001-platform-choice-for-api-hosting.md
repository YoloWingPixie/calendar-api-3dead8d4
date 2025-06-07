# ADR-001 Platform Choice for API Hosting
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a reliable, secure, scalable, and cost-effective platform to host the API, its database, and supporting infrastructure. The chosen platform must support modern DevOps practices, including CI/CD pipelines, IaC, accessibility for obvservability tooling, modern AAA (Authentication, Authorization, and Accounting) practices, and support for modern programming languages and frameworks.

## Decision Drivers
- Industry leading cloud platform
- First class IaC support for most industry standard IaC tools
- Deployment is likely to be completable within free tier of platform
- Platform provides first class support for security and compliance tooling
- Platform is well supported by all observability tools

## Consider Options
- DigitalOcean
- AWS
- GCP
- Azure

## Decision Outome
**AWS** due to its current use by the organization while meeting all decision drivers. 