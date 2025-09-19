package models

import (
	"time"
)

// NotificationPayload represents the incoming notification request
type NotificationPayload struct {
	NotificationID   string                 `json:"notificationId"`
	MessageType      string                 `json:"messageType"`
	To               string                 `json:"to"`
	TemplateBody     string                 `json:"templateBody"`
	Params           map[string]string      `json:"params"`
	ChannelConfig    map[string]interface{} `json:"channelConfig"`
	FallbackChannels []map[string]interface{} `json:"fallbackChannels"`
	Priority         Priority               `json:"priority"`
	Locale           string                 `json:"locale"`
}

// AckPayload represents the acknowledgment response
type AckPayload struct {
	NotificationID string    `json:"notificationId"`
	Status         Status    `json:"status"`
	Details        string    `json:"details"`
	Timestamp      time.Time `json:"timestamp"`
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	To      string            `json:"to"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

// Priority represents message priority levels
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Status represents processing status
type Status string

const (
	StatusSuccess Status = "SUCCESS"
	StatusFailure Status = "FAILURE"
	StatusPending Status = "PENDING"
)

// IsValid checks if the notification payload is valid
func (n *NotificationPayload) IsValid() error {
	if n.NotificationID == "" {
		return ErrInvalidNotificationID
	}
	if n.To == "" {
		return ErrInvalidRecipient
	}
	if n.TemplateBody == "" {
		return ErrMissingContent
	}
	return nil
}

// GetSubject extracts subject from channel config or provides default
func (n *NotificationPayload) GetSubject() string {
	if subject, ok := n.ChannelConfig["subject"].(string); ok {
		return subject
	}
	return "Notification" // default subject
}