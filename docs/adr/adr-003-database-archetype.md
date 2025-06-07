# ADR-003 Database Archetype Choice
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a database solution that can store and manage calendar events, user data, and calendar metadata. The chosen database must support the domain model requirements, provide reliable data persistence, and be cost-effective for development and initial deployment. The solution should align with the chosen AWS platform (from ADR-001) and support modern application patterns.

## Decision Drivers
- Technical fit for the domain model
- Alignment with AWS platform choice
- Support for the domain model's data structure
- Scalability for future growth
- Free tier availability
- Data consistency requirements
- Query performance for calendar operations
- Developer experience and productivity
- Minimize administrative effort

## Consider Options
1. **AWS RDS (PostgreSQL)**
   - Pros:
     - Perfect fit for the relational domain model
     - Native support for all required relationships and constraints
     - Rich query capabilities for complex calendar operations
     - Strong consistency guarantees
     - Better developer experience with familiar SQL patterns
     - Simpler application code due to natural data modeling
     - Free tier available for development
   - Cons:
     - Requires more initial setup compared to serverless options
     - Deployment times compared to serverless options

2. **Amazon DynamoDB**
   - Pros:
     - Serverless operation
     - Low management overhead
     - Built-in scalability
     - AWS native integration
     - Free tier available
     - Deployment times
   - Cons:
     - Requires significant application-level work to model relationships
     - Complex query patterns needed for calendar operations
     - Additional complexity in maintaining data consistency
     - Less intuitive for the domain model
     - More complex application code to handle NoSQL limitations

## Decision Outcome
**AWS RDS (PostgreSQL)** has been selected as the database solution. This decision is driven by the perfect technical fit for our domain model and the availability of a free tier for development. PostgreSQL's relational model naturally supports our calendar domain model, making it the clear choice for both development and production environments.
