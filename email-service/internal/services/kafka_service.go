package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"email-service/internal/config"
	"email-service/internal/models"
)

// KafkaService handles Kafka operations
type KafkaService struct {
	config *config.KafkaConfig
	reader *kafka.Reader
	writer *kafka.Writer
}

// NotificationProcessor is a function type for processing notifications
type NotificationProcessor func(context.Context, *models.NotificationPayload) error

// NewKafkaService creates a new Kafka service
func NewKafkaService(cfg *config.KafkaConfig) *KafkaService {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Broker},
		Topic:    cfg.Topic,
		GroupID:  cfg.ConsumerGroup,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
	})

	writer := kafka.Writer{
		Addr:     kafka.TCP(cfg.Broker),
		Topic:    cfg.AckTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaService{
		config: cfg,
		reader: reader,
		writer: &writer,
	}
}

// StartConsumer starts consuming messages from Kafka
func (ks *KafkaService) StartConsumer(ctx context.Context, processor NotificationProcessor) error {
	log.Printf("Starting Kafka consumer for topic: %s", ks.config.Topic)

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer context cancelled")
			return ctx.Err()
		default:
			// Read message with timeout
			m, err := ks.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}

			// Process message
			if err := ks.processMessage(ctx, m, processor); err != nil {
				log.Printf("Error processing message: %v", err)
			}

			// Commit message
			if err := ks.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

// processMessage processes a single Kafka message
func (ks *KafkaService) processMessage(ctx context.Context, message kafka.Message, processor NotificationProcessor) error {
	// Parse notification payload
	var payload models.NotificationPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return ks.SendAck(payload.NotificationID, models.StatusFailure, "Invalid JSON payload")
	}

	log.Printf("Processing email notification: %s for recipient: %s", payload.NotificationID, payload.To)

	// Process notification
	if err := processor(ctx, &payload); err != nil {
		log.Printf("Failed to process notification %s: %v", payload.NotificationID, err)
		return ks.SendAck(payload.NotificationID, models.StatusFailure, err.Error())
	}

	// Send success acknowledgment
	return ks.SendAck(payload.NotificationID, models.StatusSuccess, "Email sent successfully")
}

// SendAck sends acknowledgment to Kafka
func (ks *KafkaService) SendAck(notificationID string, status models.Status, details string) error {
	ack := models.AckPayload{
		NotificationID: notificationID,
		Status:         status,
		Details:        details,
		Timestamp:      time.Now(),
	}

	msg, err := json.Marshal(ack)
	if err != nil {
		return fmt.Errorf("failed to marshal ack: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return ks.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(notificationID),
		Value: msg,
	})
}

// Close closes Kafka connections
func (ks *KafkaService) Close() error {
	var readerErr, writerErr error

	if ks.reader != nil {
		readerErr = ks.reader.Close()
	}

	if ks.writer != nil {
		writerErr = ks.writer.Close()
	}

	if readerErr != nil {
		return readerErr
	}
	return writerErr
}