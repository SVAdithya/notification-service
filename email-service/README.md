# Email Service - Comprehensive Documentation

## Introduction

The email service is designed with clean architecture principles, providing a production-ready SMTP email service with
template-based notifications, proper validation, and error handling.

## Architecture Overview

The service follows **Clean Architecture** principles, ensuring a clear separation of concerns:

```
├── cmd/                    # Application entry point
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── models/           # Domain models and types
│   ├── services/         # Business logic layer
│   ├── handlers/         # HTTP handlers (interface layer)
│   └── utils/            # Utility functions
├── mocks/                # Test data and examples
├── .env                  # Environment configuration
├── Dockerfile           # Container configuration
└── go.mod              # Go module definition
```

## Key Features

### Clean Architecture

- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Inversion**: Dependencies point inward toward business logic
- **Testability**: Easy to unit test each component in isolation
- **Maintainability**: Changes in one layer don't affect others

### Production Features

- **Email Validation**: RFC-compliant email address validation
- **Header Injection Protection**: Sanitized subjects and headers
- **SMTP Connection Testing**: Startup validation and test endpoint
- **Template Rendering**: Safe parameter substitution
- **Priority Headers**: Email priority classification
- **Message Tracking**: Unique Message-ID generation
- **Graceful Shutdown**: Proper resource cleanup

### Security Features

- **Input Sanitization**: Email headers and subject sanitization
- **Format Validation**: Email address format validation
- **Error Handling**: No sensitive data in error messages
- **Non-root Execution**: Secure Docker container setup

## Project Structure Explained

### `/cmd/main.go`

Application entry point with:

- Configuration loading and validation
- SMTP connection testing on startup
- Service initialization and dependency injection
- HTTP server setup with proper timeouts
- Graceful shutdown handling

### `/internal/config/`

Configuration management:

- Environment variable loading with defaults
- Configuration validation
- SMTP and Kafka settings
- Type-safe configuration structs

### `/internal/models/`

Domain models and types:

- **notification.go**: Core notification payload structures
- **errors.go**: Application-specific error definitions
- Strong typing with proper validation

### `/internal/services/`

Business logic layer:

- **notification_service.go**: Main orchestration service
- **email_client.go**: SMTP communication with proper headers
- **kafka_service.go**: Message queue operations
- **message_builder.go**: Email message construction

### `/internal/handlers/`

HTTP interface layer:

- **http_handlers.go**: REST API endpoints
- Health checks and SMTP connection testing
- Request/response handling with proper error codes

### `/internal/utils/`

Utility functions:

- **email.go**: Email validation and sanitization
- **template.go**: Template rendering and validation
- Reusable helper functions

## Getting Started

### 1. **Configuration**

```bash
# Copy environment template
cp .env.example .env

# Edit configuration with your SMTP credentials
nano .env
```

### 2. **Required Environment Variables**

```env
# SMTP Configuration
EMAIL_SENDER=your-email@gmail.com
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-gmail-app-password

# Kafka Configuration
KAFKA_BROKER=localhost:9092
EMAIL_TOPIC=notification_email_topic
ACK_TOPIC=notification_email_ack_topic
```

### 3. **Running the Service**

#### Local Development

```bash
go run cmd/main.go
```

#### Docker

```bash
docker build -t email-service .
docker run -p 8084:8084 --env-file .env email-service
```

#### Docker Compose

```bash
docker-compose up email-service
```

## Configuration Options

| Variable              | Description             | Default                        |
|-----------------------|-------------------------|--------------------------------|
| `SERVICE_PORT`        | HTTP server port        | `8084`                         |
| `KAFKA_BROKER`        | Kafka broker address    | `localhost:9092`               |
| `EMAIL_TOPIC`         | Kafka consumption topic | `notification_email_topic`     |
| `ACK_TOPIC`           | Acknowledgment topic    | `notification_email_ack_topic` |
| `GROUP_ID`            | Kafka consumer group    | `email-service-group`          |
| `EMAIL_SENDER`        | From email address      | **Required**                   |
| `EMAIL_SMTP_HOST`     | SMTP server host        | `smtp.gmail.com`               |
| `EMAIL_SMTP_PORT`     | SMTP server port        | `587`                          |
| `EMAIL_SMTP_USER`     | SMTP username           | **Required**                   |
| `EMAIL_SMTP_PASSWORD` | SMTP password           | **Required**                   |
| `EMAIL_FROM_NAME`     | Display name            | `Notification System`          |

## API Endpoints

### Health Check

```bash
GET /actuator/health
```

Response:

```json
{
  "status": "UP",
  "service": "email-service"
}
```

### Send Test Email

```bash
POST /send-test
Content-Type: application/json

{
  "notificationId": "test-001",
  "to": "user@example.com",
  "templateBody": "Hello {name}!",
  "params": {"name": "World"},
  "channelConfig": {"subject": "Test Email"}
}
```

### Test SMTP Connection

```bash
GET /test-connection
```

Response:

```json
{
  "status": "SUCCESS",
  "message": "SMTP connection successful"
}
```

## Email Features

### Template Rendering

Use `{parameterName}` placeholders in email content:

```json
{
  "templateBody": "Hello {name}, your order #{orderNumber} is ready!",
  "params": {
    "name": "John Doe",
    "orderNumber": "12345"
  }
}
```

### Subject Configuration

Set subject via channel configuration:

```json
{
  "channelConfig": {
    "subject": "Order Confirmation - #{orderNumber}"
  }
}
```

### Priority Headers

Email priority is automatically mapped to headers:

- `urgent` → X-Priority: 1 (Highest)
- `high` → X-Priority: 2 (High)
- `medium` → X-Priority: 3 (Normal)
- `low` → X-Priority: 4 (Low)

### Custom Headers

Add custom headers via channel configuration:

```json
{
  "channelConfig": {
    "headers": {
      "X-Campaign-ID": "summer2023",
      "Reply-To": "support@example.com"
    }
  }
}
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/utils
```

### Integration Tests

```bash
# Test with mock email
curl -X POST http://localhost:8084/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-email-payload.json

# Test SMTP connection
curl http://localhost:8084/test-connection
```

### Mock Files

The service includes comprehensive mock files in `/mocks/`:

- `mock-email-payload.json`: Standard order confirmation
- `mock-welcome-email.json`: User onboarding email
- `mock-reset-password.json`: Security notification

## Monitoring & Observability

### Health Checks

- **HTTP Health Check**: `GET /actuator/health`
- **SMTP Connection Test**: `GET /test-connection`
- **Docker Health Check**: Built-in container monitoring
- **Startup Validation**: SMTP connection tested on startup

### Logging

- **Structured Logging**: JSON format ready
- **Request Tracking**: Notification ID correlation
- **Error Details**: Comprehensive error information
- **Debug Mode**: `LOG_LEVEL=debug`

## Security Features

### Input Validation

- **Email Format**: RFC-compliant validation
- **Subject Sanitization**: Header injection prevention
- **Template Safety**: Parameter validation
- **Error Sanitization**: No sensitive data leakage

### SMTP Security

- **TLS/SSL Support**: Encrypted connections
- **Authentication**: SMTP AUTH support
- **Connection Validation**: Startup and runtime testing
- **Timeout Handling**: Prevents hanging connections

## Production Deployment

### Docker Best Practices

- Multi-stage build for smaller images
- Non-root user execution
- Health checks included
- Timezone configuration

### Environment Setup

```bash
# Gmail App Password Setup
1. Enable 2FA on Gmail account
2. Generate app password: Google Account > Security > App passwords
3. Use app password as EMAIL_SMTP_PASSWORD

# Custom SMTP Server
EMAIL_SMTP_HOST=mail.yourcompany.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=notifications@yourcompany.com
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: email-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: email-service
  template:
    metadata:
      labels:
        app: email-service
    spec:
      containers:
      - name: email-service
        image: email-service:latest
        ports:
        - containerPort: 8084
        env:
        - name: EMAIL_SMTP_PASSWORD
          valueFrom:
            secretKeyRef:
              name: email-secrets
              key: smtp-password
        livenessProbe:
          httpGet:
            path: /actuator/health
            port: 8084
          initialDelaySeconds: 30
        readinessProbe:
          httpGet:
            path: /test-connection
            port: 8084
          initialDelaySeconds: 10
```

## Development

### Adding New Features

1. Define models in `/internal/models/`
2. Implement business logic in `/internal/services/`
3. Add HTTP handlers in `/internal/handlers/`
4. Update configuration if needed
5. Add tests and mock data

### Common SMTP Providers

#### Gmail

```env
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
# Use app password, not regular password
```

#### Outlook/Hotmail

```env
EMAIL_SMTP_HOST=smtp.live.com
EMAIL_SMTP_PORT=587
```

#### SendGrid

```env
EMAIL_SMTP_HOST=smtp.sendgrid.net
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=apikey
EMAIL_SMTP_PASSWORD=your-sendgrid-api-key
```

#### AWS SES

```env
EMAIL_SMTP_HOST=email-smtp.us-east-1.amazonaws.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your-ses-smtp-username
EMAIL_SMTP_PASSWORD=your-ses-smtp-password
```

## Metrics (Ready for Implementation)

The service is structured to easily add:

- Prometheus metrics for email delivery rates
- Email open/click tracking
- SMTP performance metrics
- Custom business metrics

## Benefits Achieved

### **Maintainability**

- **Clear Structure**: Easy to navigate and modify
- **Separation of Concerns**: Changes isolated to specific layers
- **Configuration Management**: Centralized settings
- **Error Handling**: Comprehensive and consistent

### **Reliability**

- **Input Validation**: Prevents invalid emails
- **Connection Testing**: Validates SMTP on startup
- **Graceful Errors**: Proper error handling and reporting
- **Resource Management**: Clean shutdown and cleanup

### **Scalability**

- **Stateless Design**: Easy horizontal scaling
- **Kafka Integration**: High-throughput message processing
- **Connection Pooling**: Efficient SMTP connections
- **Load Balancer Ready**: Health checks for LB integration

This refactored email service provides a solid foundation for enterprise email notifications with proper architecture,
comprehensive testing capabilities, and production deployment readiness.
