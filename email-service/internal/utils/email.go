package utils

import (
	"regexp"
	"strings"

	"email-service/internal/models"
)

// emailRegex validates email address format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email address format
func ValidateEmail(email string) error {
	if email == "" {
		return models.ErrInvalidRecipient
	}

	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return models.ErrInvalidEmailFormat
	}

	return nil
}

// FormatEmailAddress formats email address with optional display name
func FormatEmailAddress(email, name string) string {
	if name == "" {
		return email
	}
	return name + " <" + email + ">"
}

// SanitizeSubject removes potentially dangerous characters from subject
func SanitizeSubject(subject string) string {
	// Remove line breaks and other control characters that could be used for header injection
	subject = strings.ReplaceAll(subject, "\n", " ")
	subject = strings.ReplaceAll(subject, "\r", " ")
	subject = strings.ReplaceAll(subject, "\t", " ")
	
	// Trim whitespace
	subject = strings.TrimSpace(subject)
	
	// Limit length to prevent extremely long subjects
	if len(subject) > 200 {
		subject = subject[:197] + "..."
	}
	
	return subject
}

// IsValidEmailAddress checks if email address is valid
func IsValidEmailAddress(email string) bool {
	return ValidateEmail(email) == nil
}