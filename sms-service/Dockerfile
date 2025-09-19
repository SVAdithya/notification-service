# ---- Build Stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o sms-service .

# ---- Runtime Stage ----
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /app/sms-service .
EXPOSE 8086
ENTRYPOINT ["/app/sms-service"]
