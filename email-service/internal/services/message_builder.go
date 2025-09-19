package services

import (
	"email-service/internal/models"
	"email-service/internal/utils"
)

// MessageBuilder handles creation of email messages
type MessageBuilder struct{}

// NewMessageBuilder creates a new message builder
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

// BuildEmailMessage creates an email message from notification payload
func (mb *MessageBuilder) BuildEmailMessage(payload *models.NotificationPayload) (*models.EmailMessage, error) {
	// Validate payload
	if err := payload.IsValid(); err != nil {
		return nil, err
	}

	// Get subject from channel config or use default
	subject := payload.GetSubject()
	
	// Render templates with parameters
	renderedSubject, renderedBody := utils.RenderEmailTemplate(subject, payload.TemplateBody, payload.Params)

	// Create email message
	message := &models.EmailMessage{
		To:      payload.To,
		Subject: renderedSubject,
		Body:    renderedBody,
		Headers: make(map[string]string),
	}

	// Add priority header if specified
	if payload.Priority != "" {
		message.Headers["X-Priority"] = mb.mapPriorityToHeader(payload.Priority)
	}

	// Add notification ID for tracking
	if payload.NotificationID != "" {
		message.Headers["X-Notification-ID"] = payload.NotificationID
	}

	// Add custom headers from channel config
	if headers, ok := payload.ChannelConfig["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				message.Headers[key] = strValue
			}
		}
	}

	return message, nil
}

// mapPriorityToHeader maps priority enum to email header value
func (mb *MessageBuilder) mapPriorityToHeader(priority models.Priority) string {
	switch priority {
	case models.PriorityUrgent:
		return "1 (Highest)"
	case models.PriorityHigh:
		return "2 (High)"
	case models.PriorityMedium:
		return "3 (Normal)"
	case models.PriorityLow:
		return "4 (Low)"
	default:
		return "3 (Normal)"
	}
}

// ValidateEmailMessage validates the constructed email message
func (mb *MessageBuilder) ValidateEmailMessage(message *models.EmailMessage) error {
	if err := utils.ValidateEmail(message.To); err != nil {
		return err
	}

	if message.Subject == "" {
		return models.ErrInvalidSubject
	}

	if message.Body == "" {
		return models.ErrMissingContent
	}

	return nil
}