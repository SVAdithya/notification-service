# ---- Build Stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o email-service .

# ---- Runtime Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/email-service .
EXPOSE 8084
ENTRYPOINT ["/app/email-service"]
