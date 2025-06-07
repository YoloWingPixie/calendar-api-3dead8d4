# ADR-005 API Language and Framework
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a language and framework for implementing the calendar API. The chosen stack must support rapid development while maintaining performance and reliability. The solution should align with the chosen AWS platform and container strategy, while enabling efficient implementation of the domain model and business logic.

## Decision Drivers
- Development velocity
- Runtime performance
- Type safety
- Database integration
- API documentation
- Testing capabilities
- Container compatibility
- Community support

## Consider Options
1. **Python with FastAPI**
   - Pros:
     - Rapid development with Pydantic
     - Excellent SQLAlchemy integration
     - Built-in OpenAPI documentation
     - Async support
     - Strong type hints
     - Fast development cycle
     - Rich ecosystem (Alembic, uv)
     - Easy to understand and maintain
   - Cons:
     - Slower runtime than Go
     - Higher memory usage
     - GIL limitations
     - Less suitable for CPU-bound tasks

2. **Go**
   - Pros:
     - Excellent performance
     - Low memory footprint
     - Strong concurrency
     - Static typing
     - Fast compilation
     - Good container support
     - Built-in testing
   - Cons:
     - Steeper learning curve
     - More boilerplate code
     - Less mature ORM options
     - Slower development cycle
     - Less expressive type system

3. **TypeScript with Express/NestJS**
   - Pros:
     - JavaScript ecosystem
     - Good type system
     - Familiar to web developers
     - Rich middleware ecosystem
   - Cons:
     - Runtime type checking
     - Less mature database tooling
     - Higher memory usage
     - Less suitable for backend services

## Decision Outcome
**Python with FastAPI** has been selected as the API language and framework. While Go would provide better runtime performance and resource utilization, Python's development speed and rich ecosystem (Pydantic, SQLAlchemy, Alembic, uv) make it the better choice given the time constraints. The performance difference is acceptable for the current scale (1.2 RPS peak), and the development velocity benefits outweigh the runtime performance considerations. FastAPI was chosen over Flask for its modern async support, built-in OpenAPI documentation, and better type safety through Pydantic. 