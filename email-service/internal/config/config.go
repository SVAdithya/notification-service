package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig `json:"server"`
	Kafka  KafkaConfig  `json:"kafka"`
	Email  EmailConfig  `json:"email"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	Broker        string `json:"broker"`
	ConsumerGroup string `json:"consumer_group"`
	Topic         string `json:"topic"`
	AckTopic      string `json:"ack_topic"`
	MinBytes      int    `json:"min_bytes"`
	MaxBytes      int    `json:"max_bytes"`
}

// EmailConfig contains SMTP email configuration
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     string `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	FromEmail    string `json:"from_email"`
	FromName     string `json:"from_name"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVICE_PORT", "8084"),
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Kafka: KafkaConfig{
			Broker:        getEnv("KAFKA_BROKER", "localhost:9092"),
			ConsumerGroup: getEnv("GROUP_ID", "email-service-group"),
			Topic:         getEnv("EMAIL_TOPIC", "notification_email_topic"),
			AckTopic:      getEnv("ACK_TOPIC", "notification_email_ack_topic"),
			MinBytes:      10000,    // 10KB
			MaxBytes:      10000000, // 10MB
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("EMAIL_SMTP_PORT", "587"),
			SMTPUser:     getEnv("EMAIL_SMTP_USER", ""),
			SMTPPassword: getEnv("EMAIL_SMTP_PASSWORD", ""),
			FromEmail:    getEnv("EMAIL_SENDER", ""),
			FromName:     getEnv("EMAIL_FROM_NAME", "Notification System"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate validates required configuration fields
func (c *Config) validate() error {
	if c.Email.FromEmail == "" {
		return fmt.Errorf("EMAIL_SENDER is required")
	}
	if c.Email.SMTPUser == "" {
		return fmt.Errorf("EMAIL_SMTP_USER is required")
	}
	if c.Email.SMTPPassword == "" {
		return fmt.Errorf("EMAIL_SMTP_PASSWORD is required")
	}
	return nil
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}