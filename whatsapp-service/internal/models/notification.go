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
	MediaURL         string                 `json:"mediaUrl,omitempty"`
	MediaType        MediaType              `json:"mediaType,omitempty"`
	TemplateName     string                 `json:"templateName,omitempty"`
}

// AckPayload represents the acknowledgment response
type AckPayload struct {
	NotificationID string    `json:"notificationId"`
	Status         Status    `json:"status"`
	Details        string    `json:"details"`
	Timestamp      time.Time `json:"timestamp"`
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

// MediaType represents supported media types
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeVideo    MediaType = "video"
)

// IsValid checks if the notification payload is valid
func (n *NotificationPayload) IsValid() error {
	if n.NotificationID == "" {
		return ErrInvalidNotificationID
	}
	if n.To == "" {
		return ErrInvalidRecipient
	}
	if n.TemplateBody == "" && n.TemplateName == "" {
		return ErrMissingContent
	}
	return nil
}