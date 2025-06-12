# ADR-007 Container Base Image Strategy
**Status**: Accepted
**Date**: 2025-06-06
**Stakeholders**:
    - Zachary Maynard

## Context and Problem Statement
The project requires a base container image for the application. The chosen image must provide a secure, efficient runtime environment while minimizing attack surface and resource usage. The solution should align with our Fargate deployment strategy and support efficient container builds and deployments.

## Decision Drivers
- Image size
- Security posture
- Build performance
- Runtime performance
- Package availability
- Update frequency
- Vulnerability surface
- Resource efficiency

## Consider Options
1. **golang:1.24-alpine (build) + alpine:3.18 (runtime)**
   - Pros:
     - Minimal runtime image size
     - Excellent for static Go binaries
     - Reduced attack surface
     - Fast deployments
     - No Python dependencies
     - Good security posture
     - Official Go and Alpine images
   - Cons:
     - Requires multi-stage build
     - Debugging tools limited in runtime image
     - Alpine musl libc (irrelevant for static Go binaries)

2. **python:3.13-slim**
   - Pros:
     - Minimal base image size
     - Regular security updates
     - Essential build tools included
     - Good package availability
     - Efficient resource usage
     - Faster deployments
     - Reduced attack surface
     - Official Python image
   - Cons:
     - Limited debugging tools
     - May need additional packages
     - Less comprehensive than full image

3. **python:3.13**
   - Pros:
     - Complete Python environment
     - All development tools included
     - Easier debugging
     - More packages pre-installed
   - Cons:
     - Larger image size
     - More attack surface
     - Slower deployments
     - Higher resource usage

4. **python:3.13-alpine**
   - Pros:
     - Smallest image size
     - Minimal attack surface
     - Fast deployments
   - Cons:
     - musl libc compatibility issues
     - Limited package availability
     - Build performance overhead
     - Potential runtime issues

## Decision Outcome
**UPDATE (2025-06):** The project now uses a multi-stage Docker build: **golang:1.24-alpine** for building the Go application, and **alpine:3.18** as the runtime image. This approach produces a minimal, secure container with only the statically compiled Go binary and essential runtime dependencies. The result is a significantly smaller attack surface, faster deployments, and improved security posture, fully aligned with best practices for Go applications in containerized environments. The original Python-based analysis remains below for historical context.

---

**HISTORICAL:**
python:3.13-slim was originally selected as the base container image. This decision was driven by the need for a balance between image size, security, and functionality. The slim variant provided essential build tools while maintaining a small footprint, making it ideal for Fargate deployment. While the full image offered more tools and the Alpine variant was smaller, the slim image provided the best balance of features, security, and efficiency for the original Python FastAPI application.
