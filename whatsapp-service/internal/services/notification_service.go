package services

import (
	"context"
	"fmt"
	"log"

	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
	"whatsapp-service/internal/utils"
)

// NotificationService handles WhatsApp notification processing
type NotificationService struct {
	config         *config.Config
	whatsappClient *WhatsAppClient
	messageBuilder *MessageBuilder
	kafkaService   *KafkaService
}

// NewNotificationService creates a new notification service
func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		config:         cfg,
		whatsappClient: NewWhatsAppClient(&cfg.WhatsApp),
		messageBuilder: NewMessageBuilder(),
		kafkaService:   NewKafkaService(&cfg.Kafka),
	}
}

// ProcessNotification processes a notification payload and sends WhatsApp message
func (ns *NotificationService) ProcessNotification(ctx context.Context, payload *models.NotificationPayload) error {
	// Validate payload
	if err := payload.IsValid(); err != nil {
		return fmt.Errorf("invalid notification payload: %w", err)
	}

	// Clean and validate phone number
	cleanTo, err := utils.CleanPhoneNumber(payload.To)
	if err != nil {
		return fmt.Errorf("invalid phone number %s: %w", payload.To, err)
	}

	// Build WhatsApp message
	message, err := ns.messageBuilder.BuildMessage(payload, cleanTo)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// Send message via WhatsApp API
	response, err := ns.whatsappClient.SendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	log.Printf("WhatsApp message sent successfully. ID: %s, Recipient: %s",
		response.Messages[0].ID, cleanTo)

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

// VerifyWebhook verifies WhatsApp webhook requests
func (ns *NotificationService) VerifyWebhook(verifyToken, challenge string) (string, error) {
	return ns.whatsappClient.VerifyWebhook(verifyToken, challenge)
}

// Close closes all service connections
func (ns *NotificationService) Close() error {
	if ns.kafkaService != nil {
		return ns.kafkaService.Close()
	}
	return nil
}
