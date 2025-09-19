package service

import (
	"context"
	"fmt"
	"log"

	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
	"whatsapp-service/internal/whatsapp"
)

// WhatsAppService handles WhatsApp message sending business logic
type WhatsAppService struct {
	client         *whatsapp.Client
	messageBuilder *whatsapp.MessageBuilder
	ackProducer    AckProducer
}

// AckProducer defines interface for sending acknowledgments
type AckProducer interface {
	SendAck(ctx context.Context, notificationID string, status models.AckStatus, details string) error
}

// NewWhatsAppService creates a new WhatsApp service
func NewWhatsAppService(cfg *config.WhatsAppConfig, ackProducer AckProducer) *WhatsAppService {
	return &WhatsAppService{
		client:         whatsapp.NewClient(cfg),
		messageBuilder: whatsapp.NewMessageBuilder(),
		ackProducer:    ackProducer,
	}
}

// HandleMessage processes notification payload and sends WhatsApp message
func (s *WhatsAppService) HandleMessage(ctx context.Context, payload *models.NotificationPayload) error {
	// Build appropriate message type
	message := s.messageBuilder.BuildFromPayload(payload)

	// Send message via WhatsApp API
	response, err := s.client.SendMessage(ctx, message)
	if err != nil {
		// Send failure acknowledgment
		if ackErr := s.ackProducer.SendAck(ctx, payload.NotificationID, models.AckStatusFailure,
			fmt.Sprintf("Failed to send WhatsApp message: %v", err)); ackErr != nil {
			log.Printf("Failed to send failure ack: %v", ackErr)
		}
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	// Log success
	if len(response.Messages) > 0 {
		log.Printf("WhatsApp message sent successfully. ID: %s", response.Messages[0].ID)
	}

	// Send success acknowledgment
	if err := s.ackProducer.SendAck(ctx, payload.NotificationID, models.AckStatusSuccess,
		"WhatsApp message sent successfully"); err != nil {
		log.Printf("Failed to send success ack: %v", err)
		// Don't return error as the main operation (sending message) was successful
	}

	return nil
}

// SendTestMessage sends a test message directly
func (s *WhatsAppService) SendTestMessage(ctx context.Context, payload *models.NotificationPayload) error {
	message := s.messageBuilder.BuildFromPayload(payload)
	
	response, err := s.client.SendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send test message: %w", err)
	}

	if len(response.Messages) > 0 {
		log.Printf("Test message sent successfully. ID: %s", response.Messages[0].ID)
	}

	return nil
}

// VerifyWebhook verifies webhook challenge
func (s *WhatsAppService) VerifyWebhook(verifyToken, challenge string) (string, error) {
	return s.client.VerifyWebhook(verifyToken, challenge)
}

// ProcessWebhook processes incoming webhook notifications
func (s *WhatsAppService) ProcessWebhook(ctx context.Context, payload []byte) error {
	// Log the webhook payload for now
	// In production, you might want to parse and handle different webhook events
	log.Printf("Received webhook payload: %s", string(payload))
	
	// TODO: Implement webhook event processing
	// - Message status updates (sent, delivered, read)
	// - User message responses
	// - Other webhook events
	
	return nil
}