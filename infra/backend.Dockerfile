# Backend Dockerfile
# Multi-stage build for optimized production image

# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files first for better layer caching
COPY app/backend/go.mod app/backend/go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download && go mod verify

# Copy source code
COPY app/backend/ ./

# Build the application
# CGO_ENABLED=0 for static binary (no external dependencies)
# -ldflags="-w -s" to strip debug info and reduce binary size
# -trimpath to remove absolute paths for reproducible builds
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -trimpath \
    -o server .

# Runtime stage - use Alpine for small size with wget for health checks
FROM alpine:latest

WORKDIR /app

# Install only ca-certificates and wget (minimal dependencies)
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Copy binary from builder
COPY --from=builder /build/server .

# Change ownership
RUN chown appuser:appuser /app/server

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Run the application
CMD ["./server"]
