# WhatsApp Service Mock Payloads

This folder contains example payloads for testing the WhatsApp service with different message types.

## Available Mock Files

### 📝 Text Messages

- **`mock-text-payload.json`**: Simple text message with parameter substitution
- Use case: Basic notifications, alerts, confirmations

### 🖼️ Image Messages

- **`mock-image-payload.json`**: Image message with caption
- Use case: Receipts, QR codes, product images

### 📄 Document Messages

- **`mock-document-payload.json`**: Document/file sharing with filename
- Use case: Invoices, PDFs, reports, contracts

### 📋 Template Messages

- **`mock-template-payload.json`**: Pre-approved WhatsApp Business template
- Use case: Marketing messages, order confirmations, appointment reminders

## Usage Examples

### Direct Service Testing

```bash
# Test text message
curl -X POST http://localhost:8085/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-text-payload.json

# Test image message
curl -X POST http://localhost:8085/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-image-payload.json

# Test document message
curl -X POST http://localhost:8085/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-document-payload.json

# Test template message
curl -X POST http://localhost:8085/send-test \
  -H "Content-Type: application/json" \
  -d @mocks/mock-template-payload.json
```

### Producer Service Testing

```bash
# Send via producer service
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "whatsapp",
    "payload": '$(cat mocks/mock-text-payload.json)'
  }'
```

### Bulk Testing

```bash
# Test all message types in sequence
for file in mocks/mock-*.json; do
  echo "Testing: $file"
  curl -X POST http://localhost:8085/send-test \
    -H "Content-Type: application/json" \
    -d @$file
  sleep 1
done
```

## Customization

To create your own test payloads:

1. Copy an existing mock file
2. Update the `notificationId` to be unique
3. Modify the `to` field with a valid phone number
4. Adjust message content and parameters
5. Set appropriate `mediaUrl` for media messages
6. Ensure `templateName` exists in WhatsApp Business Manager for template messages

## Required Fields

All payloads must include:

- `notificationId`: Unique identifier for tracking
- `to`: Valid international phone number (e.g., "+1234567890")
- Either `templateBody` or `templateName` for message content

## Media Requirements

For media messages:

- `mediaUrl`: Must be publicly accessible HTTPS URL
- `mediaType`: One of "image", "document", "audio", "video"
- File size limits: Images (5MB), Documents (100MB), Audio/Video (16MB)

## Template Requirements

For template messages:

- Template must be created and approved in WhatsApp Business Manager
- Parameter order in `params` must match template definition
- Language code in `locale` must match template language