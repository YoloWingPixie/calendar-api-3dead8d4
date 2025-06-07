# ADR-004 Application Compute Platform
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a compute platform to host the calendar API application. The chosen platform must support the application's scalability requirements, align with AWS platform choice (from ADR-001), and provide efficient resource utilization. The solution should support containerized applications and enable easy scaling based on demand.

## Decision Drivers
- Container orchestration capabilities
- Resource utilization efficiency
- Operational complexity
- Scaling capabilities
- Cost effectiveness
- Deployment flexibility
- Monitoring and observability
- Team expertise requirements

## Consider Options
1. **Amazon ECS with Fargate**
   - Pros:
     - Serverless container execution
     - No EC2 management overhead
     - Automatic scaling
     - Pay-per-task pricing
     - Cost-effective for current scale
     - Native AWS integration
     - Sufficient for current RPS (1.2 peak)
     - Zero server maintenance
   - Cons:
     - Less flexible than Kubernetes
     - Limited multi-cloud support
     - Fewer advanced features

2. **Amazon EKS**
   - Pros:
     - Full Kubernetes capabilities
     - Better for very high scale
     - More flexible orchestration
     - Strong multi-cloud support
     - Advanced features available
   - Cons:
     - Higher operational complexity
     - More resource overhead
     - Overkill for current RPS
     - Higher base cost

3. **EC2 (Not Considered)**
   - Pros:
     - Direct control
     - No container overhead
   - Cons:
     - Manual scaling
     - Higher operational burden
     - Legacy approach
     - Not suitable for modern applications

## Decision Outcome
**Amazon ECS with Fargate** has been selected as the application compute platform. This decision is driven by the current scale requirements (1.2 RPS peak) and the need for operational simplicity. Fargate's serverless nature eliminates the need for EC2 instance management while providing automatic scaling capabilities. While EKS would be more suitable for significantly higher scale (1000x current RPS), ECS with Fargate provides sufficient capabilities for the current workload with minimal operational overhead. The migration path from ECS to EKS is well-documented and straightforward if future scale requirements necessitate it.
