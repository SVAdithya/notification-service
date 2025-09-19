package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"whatsapp-service/internal/config"
	"whatsapp-service/internal/models"
)

// WhatsAppClient handles communication with WhatsApp Business API
type WhatsAppClient struct {
	config     *config.WhatsAppConfig
	httpClient *http.Client
}

// NewWhatsAppClient creates a new WhatsApp API client
func NewWhatsAppClient(cfg *config.WhatsAppConfig) *WhatsAppClient {
	return &WhatsAppClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a message through WhatsApp Business API
func (c *WhatsAppClient) SendMessage(ctx context.Context, message *models.WhatsAppMessage) (*models.WhatsAppResponse, error) {
	// Marshal message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrMarshalFailed, err)
	}

	// Create request URL
	url := fmt.Sprintf("%s/%s/%s/messages",
		c.config.BaseURL,
		c.config.APIVersion,
		c.config.PhoneNumberID,
	)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrNetworkError, err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	// Parse success response
	var response models.WhatsAppResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrUnmarshalFailed, err)
	}

	// Validate response
	if len(response.Messages) == 0 {
		return nil, fmt.Errorf("no message ID in response")
	}

	return &response, nil
}

// handleAPIError processes API error responses
func (c *WhatsAppClient) handleAPIError(statusCode int, body []byte) error {
	var whatsappErr models.WhatsAppError
	if json.Unmarshal(body, &whatsappErr) == nil {
		return fmt.Errorf("%w: %s (code: %d, type: %s)",
			models.ErrAPICallFailed,
			whatsappErr.Error.Message,
			whatsappErr.Error.Code,
			whatsappErr.Error.Type,
		)
	}

	return fmt.Errorf("%w: HTTP %d - %s",
		models.ErrAPICallFailed,
		statusCode,
		string(body),
	)
}

// VerifyWebhook verifies webhook requests from WhatsApp
func (c *WhatsAppClient) VerifyWebhook(verifyToken, challenge string) (string, error) {
	if verifyToken != c.config.WebhookVerifyToken {
		return "", fmt.Errorf("invalid verify token")
	}
	return challenge, nil
}