# ---- Build Stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o whatsapp-service .

# ---- Runtime Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/whatsapp-service .
EXPOSE 8085
ENTRYPOINT ["/app/whatsapp-service"]
