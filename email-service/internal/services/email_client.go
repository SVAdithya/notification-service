package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"email-service/internal/config"
	"email-service/internal/models"
	"email-service/internal/utils"
)

// EmailClient handles SMTP email communication
type EmailClient struct {
	config *config.EmailConfig
}

// NewEmailClient creates a new SMTP email client
func NewEmailClient(cfg *config.EmailConfig) *EmailClient {
	return &EmailClient{
		config: cfg,
	}
}

// SendEmail sends an email message via SMTP
func (c *EmailClient) SendEmail(ctx context.Context, message *models.EmailMessage) error {
	// Validate email address
	if err := utils.ValidateEmail(message.To); err != nil {
		return fmt.Errorf("invalid recipient email: %w", err)
	}

	// Sanitize subject to prevent header injection
	message.Subject = utils.SanitizeSubject(message.Subject)

	// Create SMTP authentication
	auth := smtp.PlainAuth(
		"",
		c.config.SMTPUser,
		c.config.SMTPPassword,
		c.config.SMTPHost,
	)

	// Build email message
	emailBody, err := c.buildEmailMessage(message)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Send email with timeout context
	smtpAddr := c.config.SMTPHost + ":" + c.config.SMTPPort
	
	// Create a channel to handle the SMTP operation
	errChan := make(chan error, 1)
	go func() {
		err := smtp.SendMail(
			smtpAddr,
			auth,
			c.config.FromEmail,
			[]string{message.To},
			[]byte(emailBody),
		)
		errChan <- err
	}()

	// Wait for completion or timeout
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("%w: %v", models.ErrSMTPError, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: context cancelled", models.ErrTimeout)
	}
}

// buildEmailMessage constructs the email message with proper headers
func (c *EmailClient) buildEmailMessage(message *models.EmailMessage) (string, error) {
	if message.To == "" {
		return "", models.ErrInvalidRecipient
	}
	if message.Subject == "" {
		return "", models.ErrInvalidSubject
	}

	// Build email headers
	var headers []string
	
	// From header with optional display name
	fromAddr := utils.FormatEmailAddress(c.config.FromEmail, c.config.FromName)
	headers = append(headers, "From: "+fromAddr)
	
	// To header
	headers = append(headers, "To: "+message.To)
	
	// Subject header
	headers = append(headers, "Subject: "+message.Subject)
	
	// Content type and encoding
	headers = append(headers, "MIME-Version: 1.0")
	headers = append(headers, "Content-Type: text/plain; charset=UTF-8")
	headers = append(headers, "Content-Transfer-Encoding: 8bit")
	
	// Date header
	headers = append(headers, "Date: "+time.Now().Format(time.RFC1123Z))
	
	// Additional custom headers
	for key, value := range message.Headers {
		headers = append(headers, key+": "+value)
	}
	
	// Message ID for tracking
	messageID := fmt.Sprintf("<%d.%s@%s>", 
		time.Now().Unix(), 
		generateMessageIDLocal(), 
		c.extractDomainFromEmail(c.config.FromEmail))
	headers = append(headers, "Message-ID: "+messageID)

	// Build final message
	headerStr := strings.Join(headers, "\r\n")
	body := strings.ReplaceAll(message.Body, "\n", "\r\n")
	
	return headerStr + "\r\n\r\n" + body + "\r\n", nil
}

// extractDomainFromEmail extracts domain from email address
func (c *EmailClient) extractDomainFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return "localhost"
}

// generateMessageIDLocal generates a unique local part for message ID
func generateMessageIDLocal() string {
	return fmt.Sprintf("notif.%d", time.Now().UnixNano())
}

// TestConnection tests SMTP connection without sending email
func (c *EmailClient) TestConnection() error {
	smtpAddr := c.config.SMTPHost + ":" + c.config.SMTPPort
	
	// Create SMTP authentication
	auth := smtp.PlainAuth(
		"",
		c.config.SMTPUser,
		c.config.SMTPPassword,
		c.config.SMTPHost,
	)

	// Try to establish connection and authenticate
	client, err := smtp.Dial(smtpAddr)
	if err != nil {
		return fmt.Errorf("%w: failed to connect to SMTP server", models.ErrNetworkError)
	}
	defer client.Quit()

	// Test authentication
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%w: SMTP authentication failed", models.ErrAuthenticationFailed)
	}

	return nil
}