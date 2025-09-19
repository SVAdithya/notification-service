# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o email-service ./cmd

# ---- Runtime Stage ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl tzdata
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/email-service .

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8084

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8084/actuator/health || exit 1

ENTRYPOINT ["./email-service"]
