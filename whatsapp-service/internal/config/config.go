package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Kafka    KafkaConfig    `json:"kafka"`
	WhatsApp WhatsAppConfig `json:"whatsapp"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// KafkaConfig contains Kafka configuration
type KafkaConfig struct {
	Broker         string `json:"broker"`
	ConsumerGroup  string `json:"consumer_group"`
	Topic          string `json:"topic"`
	AckTopic       string `json:"ack_topic"`
	MinBytes       int    `json:"min_bytes"`
	MaxBytes       int    `json:"max_bytes"`
}

// WhatsAppConfig contains WhatsApp Business API configuration
type WhatsAppConfig struct {
	AccessToken      string `json:"access_token"`
	PhoneNumberID    string `json:"phone_number_id"`
	BusinessAccountID string `json:"business_account_id"`
	APIVersion       string `json:"api_version"`
	WebhookVerifyToken string `json:"webhook_verify_token"`
	BaseURL          string `json:"base_url"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVICE_PORT", "8085"),
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Kafka: KafkaConfig{
			Broker:        getEnv("KAFKA_BROKER", "localhost:9092"),
			ConsumerGroup: getEnv("GROUP_ID", "whatsapp-service-group"),
			Topic:         getEnv("WHATSAPP_TOPIC", "notification_whatsapp_topic"),
			AckTopic:      getEnv("ACK_TOPIC", "notification_whatsapp_ack_topic"),
			MinBytes:      10000,  // 10KB
			MaxBytes:      10000000, // 10MB
		},
		WhatsApp: WhatsAppConfig{
			AccessToken:        getEnv("WHATSAPP_ACCESS_TOKEN", ""),
			PhoneNumberID:      getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
			BusinessAccountID:  getEnv("WHATSAPP_BUSINESS_ACCOUNT_ID", ""),
			APIVersion:         getEnv("WHATSAPP_API_VERSION", "v18.0"),
			WebhookVerifyToken: getEnv("WHATSAPP_WEBHOOK_VERIFY_TOKEN", ""),
			BaseURL:            "https://graph.facebook.com",
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate validates required configuration fields
func (c *Config) validate() error {
	if c.WhatsApp.AccessToken == "" {
		return fmt.Errorf("WHATSAPP_ACCESS_TOKEN is required")
	}
	if c.WhatsApp.PhoneNumberID == "" {
		return fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID is required")
	}

	if c.Kafka.Broker == "" {
		return fmt.Errorf("KAFKA_BROKER is required")
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