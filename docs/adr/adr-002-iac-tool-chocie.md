# ADR-002 IaC Tool Choice
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a tool to manage the infrastructure as code (IaC) for the project. The chosen tool must be a first class citizen in the industry, support modern DevOps practices, and be supported by all observability tools. The tool should enable efficient infrastructure management, support the chosen AWS platform (from ADR-001), and provide a robust ecosystem for infrastructure automation.

## Decision Drivers
- Provider ecosystem breadth and maturity
- State management capabilities
- Infrastructure drift detection and correction
- Multi-cloud support
- Team collaboration features
- Development experience and tooling
- Community size and activity
- Documentation quality and completeness

## Consider Options
1. **Terraform**
   - Pros:
     - Industry standard
     - Extensive provider support
     - Large community and ecosystem
     - Well-documented
     - Remote state management
   - Cons:
     - HashiCorp's new licensing model
     - Limited programming language support
     - HCL syntax limitations

2. **Pulumi**
   - Pros:
     - Full programming language support
     - Strong typing and IDE integration
     - Modern development experience
     - Good documentation
   - Cons:
     - Smaller community
     - Higher learning curve
     - Proprietary state management
     - Cost considerations for teams

3. **CDK**
   - Pros:
     - Native AWS integration
     - TypeScript/JavaScript support
     - Good documentation
     - AWS-backed
   - Cons:
     - AWS-specific
     - Limited multi-cloud support
     - Higher learning curve
     - Less mature than alternatives

4. **OpenTofu**
   - Pros:
     - Open source fork of Terraform
     - Compatible with Terraform providers
     - Growing community
     - No licensing restrictions
     - Same benefits as Terraform
   - Cons:
     - Newer project
     - Smaller community than Terraform
     - May lag behind Terraform updates

## Decision Outcome
**OpenTofu** has been selected as the IaC tool. This decision is driven by its compatibility with Terraform's ecosystem while providing an open-source licensing model. OpenTofu offers the best balance of industry-standard practices, community support, and future cost predictability. While Terraform could be considered for the same reasons, OpenTofu's licensing model provides better long-term cost certainty for the organization.