# ADR-006 Secrets Management Strategy
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a solution for managing sensitive configuration data and secrets across different environments. The chosen solution must support secure storage, access control, and integration with our CI/CD pipeline and container platform. The solution should enable easy secret rotation and environment-specific configuration management.

## Decision Drivers
- Integration capabilities
- Access control granularity
- Secret rotation support
- Environment management
- Cost effectiveness
- Developer experience
- Security compliance
- Deployment simplicity

## Consider Options
1. **Doppler**
   - Pros:
     - Generous free tier
     - Native Fargate integration
     - GitHub Actions integration
     - Multi-environment support
     - Simple CLI and SDK
     - Secret versioning
     - Access control per environment
     - Easy local development
   - Cons:
     - Less mature than AWS Secrets Manager
     - Smaller community
     - Third-party dependency

2. **AWS Secrets Manager**
   - Pros:
     - Native AWS integration
     - Strong security features
     - Automatic rotation
     - IAM integration
     - Well-documented
   - Cons:
     - Higher cost
     - AWS-specific
     - More complex setup
     - Less developer-friendly

3. **GitHub Actions Secrets**
   - Pros:
     - Native GitHub integration
     - Simple to use
     - Free for public repos
   - Cons:
     - Limited to GitHub
     - No secret rotation
     - Basic access control
     - Not suitable for runtime

4. **HashiCorp Vault**
   - Pros:
     - Enterprise-grade security
     - Dynamic secrets
     - Extensive features
     - Multi-cloud support
   - Cons:
     - Complex to operate
     - Overkill for current needs
     - Higher operational overhead
     - Steeper learning curve

## Decision Outcome
**Doppler** has been selected as the secrets management solution. This decision is driven by its excellent integration with our chosen stack (Fargate, GitHub Actions) and its generous free tier. While AWS Secrets Manager would provide deeper AWS integration and HashiCorp Vault would offer more advanced features, Doppler's simplicity, multi-environment support, and developer-friendly approach make it the most practical choice for our current needs. The solution provides a good balance of security, ease of use, and cost-effectiveness while supporting our container and CI/CD workflows. 