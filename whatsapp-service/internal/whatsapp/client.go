package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
)

// Client represents WhatsApp API client
type Client struct {
	config     *config.WhatsAppConfig
	httpClient *http.Client
}

// NewClient creates a new WhatsApp client
func NewClient(cfg *config.WhatsAppConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.RequestTimeout) * time.Second,
		},
	}
}

// SendMessage sends a message via WhatsApp API
func (c *Client) SendMessage(ctx context.Context, message *models.WhatsAppMessage) (*models.WhatsAppResponse, error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s/messages",
		c.config.APIBaseURL, c.config.APIVersion, c.config.PhoneNumberID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var whatsappErr models.WhatsAppError
		if json.Unmarshal(body, &whatsappErr) == nil {
			return nil, fmt.Errorf("WhatsApp API error: %s (code: %d)",
				whatsappErr.Error.Message, whatsappErr.Error.Code)
		}
		return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(body))
	}

	var response models.WhatsAppResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Messages) == 0 {
		return nil, fmt.Errorf("no message ID in response")
	}

	return &response, nil
}

// VerifyWebhook verifies webhook challenge
func (c *Client) VerifyWebhook(verifyToken, challenge string) (string, error) {
	if verifyToken == c.config.WebhookVerifyToken {
		return challenge, nil
	}
	return "", fmt.Errorf("invalid verify token")
}

// CleanPhoneNumber formats phone number for WhatsApp API
func CleanPhoneNumber(phone string) string {
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

// GetLanguageCode extracts language code from locale
func GetLanguageCode(locale string) string {
	if locale == "" {
		return "en"
	}

	// Extract language code from locale (e.g., "en_US" -> "en")
	if strings.Contains(locale, "_") {
		return strings.Split(locale, "_")[0]
	}

	return locale
}