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

// RenderEmailTemplate renders both subject and body templates
func RenderEmailTemplate(subject, body string, params map[string]string) (string, string) {
	renderedSubject := RenderTemplate(subject, params)
	renderedBody := RenderTemplate(body, params)
	return renderedSubject, renderedBody
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

// EscapeHTML escapes HTML characters in email content
func EscapeHTML(text string) string {
	replacements := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&#39;",
	}

	result := text
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}