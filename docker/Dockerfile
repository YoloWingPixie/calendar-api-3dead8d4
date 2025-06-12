# Multi-stage Dockerfile for Calendar API (Go)
# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for Go modules and HTTPS)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with cache mount for better performance
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY src/ ./src/

# Accept build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build the application with optimizations and version information
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-extldflags \"-static\" -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${DATE}" \
    -o calendar-api ./src

# Stage 2: Runtime
FROM alpine:3.18

# Install ca-certificates for HTTPS and postgresql-client for health checks
RUN apk add --no-cache ca-certificates postgresql-client curl

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/calendar-api .

# Change ownership to non-root user
RUN chown appuser:appuser /app/calendar-api

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# Run the application
CMD ["./calendar-api"]
