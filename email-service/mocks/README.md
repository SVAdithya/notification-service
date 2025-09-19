# Email Service Mock Payloads

This folder contains example payloads for testing the email service with different types of email notifications.

## Available Mock Files

### 📧 Order Confirmation

- **`mock-email-payload.json`**: Standard order confirmation email
- Use case: E-commerce order confirmations, purchase receipts

### 👋 Welcome Email

- **`mock-welcome-email.json`**: New user welcome and verification email
- Use case: User onboarding, account verification

### 🔒 Password Reset

- **`mock-reset-password.json`**: Password reset request email
- Use case: Security notifications, password recovery

## Usage Examples

### Direct Service Testing

```bash
# Test order confirmation email
curl -X POST http://localhost:8084/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-email-payload.json

# Test welcome email
curl -X POST http://localhost:8084/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-welcome-email.json

# Test password reset email
curl -X POST http://localhost:8084/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-reset-password.json
```

### Producer Service Testing

```bash
# Send via producer service
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "email",
    "payload": '$(cat mocks/mock-email-payload.json)'
  }'
```

### Bulk Testing

```bash
# Test all email types
for file in mocks/mock-*.json; do
  echo "Testing email: $file"
  curl -X POST http://localhost:8084/send-test \
    -H "Content-Type: application/json" \
    -d @$file
  sleep 2  # Allow time between emails
done
```

## Email Configuration

### Subject Line

Set the email subject using the `channelConfig.subject` field:

```json
"channelConfig": {
  "subject": "Your Custom Subject Here"
}
```

### Template Parameters

Use `{parameterName}` placeholders in both subject and body that will be replaced with values from the `params` object.

## Required Fields

All email payloads must include:

- `notificationId`: Unique identifier for tracking
- `to`: Valid email address
- `templateBody`: Email content with parameter placeholders
- `channelConfig.subject`: Email subject line

## Email Best Practices

1. **Subject Lines**: Keep under 50 characters for mobile compatibility
2. **Content**: Use plain text with clear formatting
3. **Parameters**: Always provide all required template parameters
4. **Testing**: Use valid test email addresses you control
5. **Personalization**: Include recipient name and relevant details

## SMTP Configuration

Before testing, ensure your email service is configured with:

- `EMAIL_SMTP_HOST`: SMTP server (e.g., smtp.gmail.com)
- `EMAIL_SMTP_PORT`: SMTP port (usually 587)
- `EMAIL_SMTP_USER`: Your email address
- `EMAIL_SMTP_PASSWORD`: Your email password or app password
- `EMAIL_SENDER`: From email address

## Security Notes

- Never commit real email addresses to version control
- Use app passwords instead of regular passwords for Gmail
- Test with your own email addresses first
- Be mindful of rate limits when bulk testing