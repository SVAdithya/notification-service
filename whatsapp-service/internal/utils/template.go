package utils

import (
	"strings"
)

// RenderTemplate replaces template variables with actual values
func RenderTemplate(template string, params map[string]string) string {
	if template == "" {
		return ""
	}

	result := template
	for key, value := range params {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// GetLanguageCode extracts language code from locale
func GetLanguageCode(locale string) string {
	if locale == "" {
		return "en" // default to English
	}

	// Extract language code from locale (e.g., "en_US" -> "en")
	if strings.Contains(locale, "_") {
		return strings.Split(locale, "_")[0]
	}

	return locale
}

// ValidateTemplate checks if template has required parameters
func ValidateTemplate(template string, params map[string]string) []string {
	var missingParams []string

	// Find all placeholders in template
	start := 0
	for {
		startIdx := strings.Index(template[start:], "{")
		if startIdx == -1 {
			break
		}

		startIdx += start
		endIdx := strings.Index(template[startIdx:], "}")
		if endIdx == -1 {
			break
		}

		endIdx += startIdx
		param := template[startIdx+1 : endIdx]

		// Check if parameter exists
		if _, exists := params[param]; !exists {
			// Check if already in missing list
			found := false
			for _, missing := range missingParams {
				if missing == param {
					found = true
					break
				}
			}
			if !found {
				missingParams = append(missingParams, param)
			}
		}

		start = endIdx + 1
	}

	return missingParams
}