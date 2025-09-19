package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	
	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
)

// Producer handles producing messages to Kafka
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.KafkaConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Broker),
			Topic:    cfg.AckTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// SendAck sends acknowledgment to Kafka
func (p *Producer) SendAck(ctx context.Context, notificationID string, status models.AckStatus, details string) error {
	ack := &models.AckPayload{
		NotificationID: notificationID,
		Status:         status,
		Details:        details,
		Timestamp:      time.Now().UTC(),
	}

	msg, err := json.Marshal(ack)
	if err != nil {
		return fmt.Errorf("failed to marshal ack: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(notificationID),
		Value: msg,
	}

	if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("failed to write ack to Kafka: %w", err)
	}

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}