package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	
	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
)

// Consumer handles consuming messages from Kafka
type Consumer struct {
	reader *kafka.Reader
}

// MessageHandler defines the interface for handling messages
type MessageHandler interface {
	HandleMessage(ctx context.Context, payload *models.NotificationPayload) error
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.KafkaConfig) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{cfg.Broker},
			Topic:    cfg.Topic,
			GroupID:  cfg.GroupID,
			MinBytes: cfg.MinBytes,
			MaxBytes: cfg.MaxBytes,
		}),
	}
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	log.Printf("Starting Kafka consumer for topic: %s", c.reader.Config().Topic)
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled")
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}

			var payload models.NotificationPayload
			if err := json.Unmarshal(msg.Value, &payload); err != nil {
				log.Printf("Invalid payload: %v", err)
				continue
			}

			log.Printf("Processing WhatsApp notification: %s for %s", 
				payload.NotificationID, payload.To)

			if err := handler.HandleMessage(ctx, &payload); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}