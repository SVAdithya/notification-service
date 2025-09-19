package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Configurable via env or flags
type AppConfig struct {
	KafkaBroker      string
	EmailTopic       string
	AckTopic         string
	GroupID          string
	EmailSender      string
	EmailSMTPHost    string
	EmailSMTPPort    string
	EmailSMTPUser    string
	EmailSMTPPassword string
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
		KafkaBroker:      "localhost:9092",
		EmailTopic:       "notification_email_topic",
		AckTopic:         "notification_email_ack_topic",
		GroupID:          getenv("EMAIL_CONSUMER_GROUP", "email-service-group"),
		EmailSender:      getenv("EMAIL_SENDER", ""),
		EmailSMTPHost:    getenv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
		EmailSMTPPort:    getenv("EMAIL_SMTP_PORT", "587"),
		EmailSMTPUser:    getenv("EMAIL_SMTP_USER", ""),         // Your email ID/login
		EmailSMTPPassword:getenv("EMAIL_SMTP_PASSWORD", ""),     // Your password/app password
	}
}

// NotificationPayload models the payload structure (partial, can be expanded)
type NotificationPayload struct {
	NotificationId string                 `json:"notificationId"`
	MessageType    string                 `json:"messageType"`
	To             string                 `json:"to"`
	TemplateBody   string                 `json:"templateBody"`
	Params         map[string]string      `json:"params"`
	ChannelConfig  map[string]interface{} `json:"channelConfig"`
	FallbackChannels []map[string]interface{} `json:"fallbackChannels"`
	Priority       string                 `json:"priority"`
	Locale         string                 `json:"locale"`
}

type AckPayload struct {
	NotificationId string `json:"notificationId"`
	Status         string `json:"status"`
	Details        string `json:"details"`
	Timestamp      string `json:"timestamp"`
}

func main() {
	config := getConfigFromEnv()
	fmt.Println("BROKER ADDRESS IN USE:", config.KafkaBroker)

	payload, err := loadPayload("mock-payload.json")
	if err != nil {
		fmt.Println("Error loading payload:", err)
	} else if err := handlePayload(payload, config); err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("SUCCESS: Email sent")
	}

	http.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"UP"}`))
	})

	fmt.Println("Service running. Health endpoint on :8084/actuator/health")
	if err := http.ListenAndServe(":8084", nil); err != nil {
		panic(err)
	}
}

func loadPayload(path string) (NotificationPayload, error) {
	var payload NotificationPayload
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return payload, err
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func handlePayload(payload NotificationPayload, config AppConfig) error {
	body := renderTemplate(payload.TemplateBody, payload.Params)
	return sendEmail(config, payload.To, "Notification", body)
}

func consumeEmailTopic(cfg AppConfig) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		Topic: cfg.EmailTopic,
		GroupID: cfg.GroupID,
		MinBytes: 10e3, MaxBytes: 10e6,
	})
	defer r.Close()
	ctx := context.Background()
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error reading Kafka message:", err)
			continue
		}
		var payload NotificationPayload
		if err := json.Unmarshal(m.Value, &payload); err != nil {
			fmt.Println("Invalid payload:", err)
			continue
		}
		if err := handlePayload(payload, cfg); err != nil {
			fmt.Println("Error sending email:", err)
			produceAck(cfg, payload.NotificationId, "FAILURE", "Failed to send email")
		} else {
			produceAck(cfg, payload.NotificationId, "SUCCESS", "Email sent")
		}
	}
}

func renderTemplate(template string, params map[string]string) string {
	for k, v := range params {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template
}

func sendEmail(cfg AppConfig, to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		cfg.EmailSMTPUser,
		cfg.EmailSMTPPassword,
		cfg.EmailSMTPHost,
	)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(
		cfg.EmailSMTPHost+":"+cfg.EmailSMTPPort,
		auth,
		cfg.EmailSender,
		[]string{to},
		msg,
	)
}

func produceAck(cfg AppConfig, notificationId, status, details string) {
	w := kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBroker),
		Topic:    cfg.AckTopic,
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
		fmt.Println("Failed to write ack to Kafka:", err)
	}
}
