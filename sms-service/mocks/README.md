# SMS Service Mock Payloads

This folder contains example payloads for testing the SMS service with different types of text messages.

## Available Mock Files

### 📱 Verification SMS

- **`mock-sms-payload.json`**: Standard verification code SMS
- Use case: 2FA codes, account verification, security alerts

### 🚨 Alert SMS

- **`mock-alert-sms.json`**: Critical system alert SMS
- Use case: Emergency notifications, system alerts, security breaches

## Usage Examples

### Direct Service Testing

```bash
# Test verification SMS
curl -X POST http://localhost:8086/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-sms-payload.json

# Test alert SMS
curl -X POST http://localhost:8086/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-alert-sms.json
```

### Producer Service Testing

```bash
# Send via producer service
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "sms",
    "payload": '$(cat mocks/mock-sms-payload.json)'
  }'
```

### Bulk Testing

```bash
# Test multiple SMS types
for file in mocks/mock-*.json; do
  echo "Testing SMS: $file"
  curl -X POST http://localhost:8086/send-test \
    -H "Content-Type: application/json" \
    -d @$file
  sleep 1  # Brief pause between messages
done
```

## SMS Best Practices

### Message Length

- **Standard SMS**: 160 characters max
- **Long SMS**: Automatically split into multiple messages
- **Unicode**: 70 characters max (for emojis, special characters)

### Content Guidelines

1. **Be Concise**: SMS has character limits
2. **Clear Action**: Include clear next steps if needed
3. **Time Sensitive**: SMS is for urgent/immediate notifications
4. **Professional**: Keep tone appropriate for your brand

### Template Variables

Use `{parameterName}` placeholders that will be replaced with values from the `params` object.

## Required Fields

All SMS payloads must include:

- `notificationId`: Unique identifier for tracking
- `to`: Valid international phone number (e.g., "+1234567890")
- `templateBody`: SMS content with parameter placeholders

## Phone Number Format

- **International Format**: Always use + followed by country code
- **Valid Examples**: "+1234567890", "+44123456789", "+91123456789"
- **Invalid Examples**: "123-456-7890", "(123) 456-7890", "123456789"

## Priority Levels

Set appropriate priority for SMS delivery:

- `"urgent"`: Emergency alerts, security notifications
- `"high"`: Verification codes, important updates
- `"medium"`: General notifications, confirmations
- `"low"`: Marketing messages, non-critical updates

## SMS Provider Integration

The current SMS service is a simulation. To integrate with real SMS providers:

### Popular SMS Providers

1. **Twilio**: Most popular, easy integration
2. **AWS SNS**: Part of AWS ecosystem
3. **Vonage (Nexmo)**: Good international coverage
4. **MessageBird**: European focus
5. **Plivo**: Cost-effective option

### Integration Example (Twilio)

```go
// In internal/services/sms_client.go
import "github.com/twilio/twilio-go"

client := twilio.NewRestClient()
params := &api.CreateMessageParams{}
params.SetTo(to)
params.SetFrom(twilioNumber)
params.SetBody(message)
resp, err := client.Api.CreateMessage(params)
```

## Testing Notes

### Development Testing

- Current implementation logs SMS to console
- No actual SMS delivery in development mode
- Check service logs for sent message simulation

### Production Setup

- Configure SMS provider credentials
- Set up proper error handling
- Implement delivery status tracking
- Add rate limiting for compliance

## Configuration

Before testing, ensure SMS service environment variables are set:

```env
# SMS Provider Configuration (example for future implementation)
SMS_PROVIDER=twilio
SMS_ACCOUNT_SID=your-account-sid
SMS_AUTH_TOKEN=your-auth-token
SMS_FROM_NUMBER=+1234567890

# Kafka Configuration
KAFKA_BROKER=localhost:9092
SMS_TOPIC=notification_sms_topic
ACK_TOPIC=notification_sms_ack_topic
```

## Compliance and Regulations

### Important Considerations

1. **Opt-in Required**: Users must consent to receive SMS
2. **Opt-out Support**: Include STOP/UNSUBSCRIBE options
3. **Time Restrictions**: Respect sending time zones and hours
4. **Rate Limits**: Comply with carrier and provider limits
5. **Content Restrictions**: Avoid spam-like content

### Legal Requirements

- **TCPA Compliance** (US): Obtain proper consent
- **GDPR** (EU): Handle personal data appropriately
- **Local Laws**: Check regulations in target countries

## Monitoring and Debugging

### Health Check

```bash
curl http://localhost:8086/actuator/health
```

### Service Logs

```bash
# Docker Compose
docker-compose logs -f sms-service

# Local Development
# Check console output for simulated SMS messages
```

### Kafka Monitoring

- Visit http://localhost:8080 for Kafka UI
- Monitor `notification_sms_topic` for incoming messages
- Check `notification_sms_ack_topic` for delivery confirmations