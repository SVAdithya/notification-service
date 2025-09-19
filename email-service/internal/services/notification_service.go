package services

import (
	"context"
	"fmt"
	"log"

	"email-service/internal/config"
	"email-service/internal/models"
)

// NotificationService handles email notification processing
type NotificationService struct {
	config         *config.Config
	emailClient    *EmailClient
	messageBuilder *MessageBuilder
	kafkaService   *KafkaService
}

// NewNotificationService creates a new notification service
func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		config:         cfg,
		emailClient:    NewEmailClient(&cfg.Email),
		messageBuilder: NewMessageBuilder(),
		kafkaService:   NewKafkaService(&cfg.Kafka),
	}
}

// ProcessNotification processes a notification payload and sends email
func (ns *NotificationService) ProcessNotification(ctx context.Context, payload *models.NotificationPayload) error {
	// Build email message
	message, err := ns.messageBuilder.BuildEmailMessage(payload)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Validate email message
	if err := ns.messageBuilder.ValidateEmailMessage(message); err != nil {
		return fmt.Errorf("invalid email message: %w", err)
	}

	// Send email
	if err := ns.emailClient.SendEmail(ctx, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully. Notification ID: %s, Recipient: %s, Subject: %s",
		payload.NotificationID, message.To, message.Subject)

	return nil
}

// SendAck sends acknowledgment to Kafka
func (ns *NotificationService) SendAck(notificationID string, status models.Status, details string) error {
	return ns.kafkaService.SendAck(notificationID, status, details)
}

// StartConsumer starts the Kafka consumer
func (ns *NotificationService) StartConsumer(ctx context.Context) error {
	return ns.kafkaService.StartConsumer(ctx, ns.ProcessNotification)
}

// TestEmailConnection tests the SMTP connection
func (ns *NotificationService) TestEmailConnection() error {
	return ns.emailClient.TestConnection()
}

// Close closes all service connections
func (ns *NotificationService) Close() error {
	if ns.kafkaService != nil {
		return ns.kafkaService.Close()
	}
	return nil
}