package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type AppConfig struct {
	KafkaBroker string
	SMSTopic    string
	AckTopic    string
	GroupID     string
	ServicePort string
}

type NotificationPayload struct {
	NotificationId   string                 `json:"notificationId"`
	MessageType      string                 `json:"messageType"`
	To               string                 `json:"to"`
	TemplateBody     string                 `json:"templateBody"`
	Params           map[string]string      `json:"params"`
	ChannelConfig    map[string]interface{} `json:"channelConfig"`
	FallbackChannels []map[string]interface{} `json:"fallbackChannels"`
	Priority         string                 `json:"priority"`
	Locale           string                 `json:"locale"`
}

type AckPayload struct {
	NotificationId string `json:"notificationId"`
	Status         string `json:"status"`
	Details        string `json:"details"`
	Timestamp      string `json:"timestamp"`
}

func getenv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func getConfigFromEnv() AppConfig {
	return AppConfig{
		KafkaBroker: getenv("KAFKA_BROKER", "localhost:9092"),
		SMSTopic:    getenv("SMS_TOPIC", "notification_sms_topic"),
		AckTopic:    getenv("ACK_TOPIC", "notification_sms_ack_topic"),
		GroupID:     getenv("GROUP_ID", "sms-service-group"),
		ServicePort: getenv("SERVICE_PORT", "8086"),
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	config := getConfigFromEnv()
	
	fmt.Printf("SMS Service starting on port %s\n", config.ServicePort)
	fmt.Printf("Connected to Kafka broker: %s\n", config.KafkaBroker)

	// Start Kafka consumer in a goroutine
	go consumeSMSTopic(config)

	// Health check endpoint
	http.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	// Test endpoint for sending SMS
	http.HandleFunc("/send-test", handleSendTest(config))

	fmt.Printf("SMS Service endpoints:\n")
	fmt.Printf("  Health: http://localhost:%s/actuator/health\n", config.ServicePort)
	fmt.Printf("  Test: POST http://localhost:%s/send-test\n", config.ServicePort)

	if err := http.ListenAndServe(":"+config.ServicePort, nil); err != nil {
		panic(err)
	}
}

func consumeSMSTopic(cfg AppConfig) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    cfg.SMSTopic,
		GroupID:  cfg.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx := context.Background()
	fmt.Printf("Starting to consume SMS topic: %s\n", cfg.SMSTopic)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Error reading Kafka message: %v\n", err)
			continue
		}

		var payload NotificationPayload
		if err := json.Unmarshal(m.Value, &payload); err != nil {
			fmt.Printf("Invalid payload: %v\n", err)
			continue
		}

		fmt.Printf("Processing SMS notification: %s for %s\n", payload.NotificationId, payload.To)

		if err := handleSMSPayload(payload, cfg); err != nil {
			fmt.Printf("Error sending SMS: %v\n", err)
			produceAck(cfg, payload.NotificationId, "FAILURE", fmt.Sprintf("Failed to send SMS: %v", err))
		} else {
			fmt.Println("SMS sent successfully (simulated)")
			produceAck(cfg, payload.NotificationId, "SUCCESS", "SMS sent successfully")
		}
	}
}

func handleSMSPayload(payload NotificationPayload, config AppConfig) error {
	// Clean phone number
	to := cleanPhoneNumber(payload.To)
	
	// Render template with parameters
	message := renderTemplate(payload.TemplateBody, payload.Params)
	
	// TODO: Integrate with actual SMS provider (Twilio, AWS SNS, etc.)
	// For now, we'll simulate sending
	return sendSMS(to, message)
}

func sendSMS(to, message string) error {
	// Simulate SMS sending
	fmt.Printf("📱 SMS TO: %s\n", to)
	fmt.Printf("📱 MESSAGE: %s\n", message)
	
	// TODO: Replace with actual SMS provider integration
	// Examples:
	// - Twilio: https://github.com/twilio/twilio-go
	// - AWS SNS: https://github.com/aws/aws-sdk-go
	// - Other SMS providers
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	// Return success for simulation
	return nil
}

func handleSendTest(config AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload NotificationPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := handleSMSPayload(payload, config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"sent"}`))
	}
}

func renderTemplate(template string, params map[string]string) string {
	for k, v := range params {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template
}

func cleanPhoneNumber(phone string) string {
	// Remove all non-digit characters except +
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	// Ensure it starts with + if it doesn't already
	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}
	
	return cleaned
}

func produceAck(cfg AppConfig, notificationId, status, details string) {
	w := kafka.Writer{
		Addr:  kafka.TCP(cfg.KafkaBroker),
		Topic: cfg.AckTopic,
	}
	defer w.Close()

	ack := AckPayload{
		NotificationId: notificationId,
		Status:         status,
		Details:        details,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
	}

	msg, _ := json.Marshal(ack)
	if err := w.WriteMessages(context.Background(), kafka.Message{
		Value: msg,
	}); err != nil {
		fmt.Printf("Failed to write ack to Kafka: %v\n", err)
	}
}
