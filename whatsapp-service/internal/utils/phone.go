package utils

import (
	"regexp"
	"strings"

	"whatsapp-service/internal/models"
)

// phoneNumberRegex validates phone number format (international format)
var phoneNumberRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// CleanPhoneNumber cleans and validates phone number
func CleanPhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", models.ErrInvalidPhoneNumber
	}

	// Remove all non-digit characters except +
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Ensure it starts with + if it doesn't already
	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}

	// Validate format
	if !phoneNumberRegex.MatchString(cleaned) {
		return "", models.ErrInvalidPhoneNumber
	}

	return cleaned, nil
}

// IsValidPhoneNumber checks if phone number is in valid international format
func IsValidPhoneNumber(phone string) bool {
	cleaned, err := CleanPhoneNumber(phone)
	if err != nil {
		return false
	}
	return phoneNumberRegex.MatchString(cleaned)
}