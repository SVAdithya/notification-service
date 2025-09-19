# WhatsApp Service - Clean Architecture

A production-ready WhatsApp Business API integration service built with clean architecture principles, supporting text
messages, media files, and template messages.

## 🏗️ Architecture Overview

This service follows **Clean Architecture** principles with clear separation of concerns:

```
├── cmd/                    # Application entry point
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── models/           # Domain models and types
│   ├── services/         # Business logic layer
│   ├── handlers/         # HTTP handlers (interface layer)
│   └── utils/            # Utility functions
├── .env                  # Environment configuration
├── Dockerfile           # Container configuration
└── go.mod              # Go module definition
```

## 🎯 Key Improvements

### ✅ **Clean Architecture**

- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Inversion**: Dependencies point inward toward business logic
- **Testability**: Easy to unit test each component in isolation
- **Maintainability**: Changes in one layer don't affect others

### ✅ **Industry Standards**

- **Configuration Management**: Environment-based configuration with validation
- **Error Handling**: Proper error types and wrapping
- **Logging**: Structured logging with different levels
- **Graceful Shutdown**: Proper resource cleanup on shutdown
- **Health Checks**: Docker health checks and HTTP endpoints

### ✅ **Production Features**

- **Context Support**: Proper context propagation for timeouts
- **Phone Validation**: International phone number format validation
- **Template Processing**: Safe template rendering with parameter validation
- **Message Types**: Support for text, media, and template messages
- **Webhook Handling**: WhatsApp webhook verification and processing

## 📁 Project Structure Explained

### `/cmd/main.go`

Application entry point with:

- Configuration loading and validation
- Service initialization and dependency injection
- HTTP server setup with proper timeouts
- Graceful shutdown handling
- Signal handling for clean termination

### `/internal/config/`

Configuration management:

- Environment variable loading with defaults
- Configuration validation
- Type-safe configuration structs
- Support for different environments

### `/internal/models/`

Domain models and types:

- **notification.go**: Core notification payload structures
- **whatsapp.go**: WhatsApp API request/response models
- **errors.go**: Application-specific error definitions
- Strong typing with proper validation

### `/internal/services/`

Business logic layer:

- **notification_service.go**: Main orchestration service
- **whatsapp_client.go**: WhatsApp API communication
- **kafka_service.go**: Message queue operations
- **message_builder.go**: Message construction logic

### `/internal/handlers/`

HTTP interface layer:

- **http_handlers.go**: REST API endpoints
- Request/response handling
- Input validation and error responses
- Proper HTTP status codes

### `/internal/utils/`

Utility functions:

- **phone.go**: Phone number validation and formatting
- **template.go**: Template rendering and validation
- Reusable helper functions

## 🚀 Getting Started

### 1. **Configuration**
```bash
# Copy environment template
cp .env.example .env

# Edit configuration with your credentials
nano .env
```

### 2. **Required Environment Variables**
```env
# WhatsApp API Configuration
WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxxx
WHATSAPP_PHONE_NUMBER_ID=123456789012345
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your-webhook-token

# Kafka Configuration  
KAFKA_BROKER=localhost:9092
WHATSAPP_TOPIC=notification_whatsapp_topic
ACK_TOPIC=notification_whatsapp_ack_topic
```

### 3. **Running the Service**

#### Local Development

```bash
go run cmd/main.go
```

#### Docker

```bash
docker build -t whatsapp-service .
docker run -p 8085:8085 --env-file .env whatsapp-service
```

#### Docker Compose

```bash
docker-compose up whatsapp-service
```

## 🔧 Configuration Options

| Variable                   | Description             | Default                           |
|----------------------------|-------------------------|-----------------------------------|
| `SERVICE_PORT`             | HTTP server port        | `8085`                            |
| `KAFKA_BROKER`             | Kafka broker address    | `localhost:9092`                  |
| `WHATSAPP_TOPIC`           | Kafka consumption topic | `notification_whatsapp_topic`     |
| `ACK_TOPIC`                | Acknowledgment topic    | `notification_whatsapp_ack_topic` |
| `GROUP_ID`                 | Kafka consumer group    | `whatsapp-service-group`          |
| `WHATSAPP_ACCESS_TOKEN`    | WhatsApp API token      | **Required**                      |
| `WHATSAPP_PHONE_NUMBER_ID` | Phone number ID         | **Required**                      |
| `WHATSAPP_API_VERSION`     | API version             | `v18.0`                           |
| `LOG_LEVEL`                | Logging level           | `info`                            |

## 📡 API Endpoints

### Health Check
```bash
GET /actuator/health
```

Response:

```json
{
  "status": "UP",
  "service": "whatsapp-service"
}
```

### Send Test Message

```bash
POST /send-test
Content-Type: application/json

{
  "notificationId": "test-001",
  "to": "+1234567890",
  "templateBody": "Hello {name}!",
  "params": {"name": "World"}
}
```

### Webhook Endpoint

```bash
GET /webhook?hub.verify_token=TOKEN&hub.challenge=CHALLENGE
POST /webhook
```

## 📝 Message Types

### 1. Text Messages
```json
{
  "notificationId": "text-001",
  "to": "+1234567890",
  "templateBody": "Hello {name}, your order #{order} is ready!",
  "params": {
    "name": "John",
    "order": "12345"
  }
}
```

### 2. Media Messages
```json
{
  "notificationId": "media-001",
  "to": "+1234567890",
  "templateBody": "Here's your receipt!",
  "mediaUrl": "https://example.com/receipt.jpg",
  "mediaType": "image"
}
```

### 3. Template Messages
```json
{
  "notificationId": "template-001",
  "to": "+1234567890",
  "templateName": "order_confirmation",
  "params": {
    "customer_name": "John",
    "order_number": "12345"
  }
}
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/utils
```

### Integration Tests
```bash
# Test with real WhatsApp API (requires valid credentials)
curl -X POST http://localhost:8085/send-test \
  -H "Content-Type: application/json" \
  -d @test-payload.json
```

## 📊 Monitoring & Observability

### Health Checks

- **HTTP Health Check**: `GET /actuator/health`
- **Docker Health Check**: Built-in container health monitoring
- **Kubernetes Health Check**: Ready for k8s liveness/readiness probes

### Logging

```bash
# Debug mode
LOG_LEVEL=debug go run cmd/main.go

# Production logging
LOG_LEVEL=info go run cmd/main.go
```

### Metrics (Ready for Implementation)

The service is structured to easily add:

- Prometheus metrics
- Distributed tracing
- Custom business metrics
- Performance monitoring

## 🔒 Security Features

### Input Validation

- Phone number format validation
- Template parameter sanitization
- JSON payload validation
- Error message sanitization

### Authentication

- WhatsApp webhook token verification
- API token validation
- Secure credential handling

### Best Practices

- No sensitive data in logs
- Proper error handling without data leakage
- Non-root Docker user
- Health check endpoints

## 🚀 Production Deployment

### Docker Best Practices

- Multi-stage build for smaller images
- Non-root user execution
- Health checks included
- Timezone configuration

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whatsapp-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: whatsapp-service
  template:
    metadata:
      labels:
        app: whatsapp-service
    spec:
      containers:
      - name: whatsapp-service
        image: whatsapp-service:latest
        ports:
        - containerPort: 8085
        env:
        - name: WHATSAPP_ACCESS_TOKEN
          valueFrom:
            secretKeyRef:
              name: whatsapp-secrets
              key: access-token
        livenessProbe:
          httpGet:
            path: /actuator/health
            port: 8085
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /actuator/health
            port: 8085
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 🔄 Scaling & Performance

### Horizontal Scaling

- Stateless service design
- Kafka consumer group balancing
- Multiple service instances supported
- Load balancer ready

### Performance Optimizations

- HTTP connection pooling
- Kafka batch processing
- Efficient phone number validation
- Template caching ready

## 🛠️ Development

### Adding New Features

1. Define models in `/internal/models/`
2. Implement business logic in `/internal/services/`
3. Add HTTP handlers in `/internal/handlers/`
4. Update configuration if needed
5. Add tests for all components

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Security scan
gosec ./...

# Dependency check
go mod tidy
```

## 📖 Additional Resources

- [WhatsApp Business API Documentation](https://developers.facebook.com/docs/whatsapp)
- [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

This refactored service provides a solid foundation for enterprise-grade WhatsApp integration with proper architecture,
testing, and deployment capabilities.