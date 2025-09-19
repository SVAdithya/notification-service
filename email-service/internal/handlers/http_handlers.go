package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"email-service/internal/models"
	"email-service/internal/services"
)

// HTTPHandlers contains all HTTP handlers
type HTTPHandlers struct {
	notificationService *services.NotificationService
}

// NewHTTPHandlers creates new HTTP handlers
func NewHTTPHandlers(notificationService *services.NotificationService) *HTTPHandlers {
	return &HTTPHandlers{
		notificationService: notificationService,
	}
}

// HealthCheck handles health check requests
func (h *HTTPHandlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":  "UP",
		"service": "email-service",
	}
	
	json.NewEncoder(w).Encode(response)
}

// SendTestHandler handles direct test email requests
func (h *HTTPHandlers) SendTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload models.NotificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid JSON payload: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process notification directly (bypass Kafka)
	if err := h.notificationService.ProcessNotification(r.Context(), &payload); err != nil {
		log.Printf("Failed to send test email: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":         "sent",
		"notificationId": payload.NotificationID,
		"message":        "Test email sent successfully",
		"recipient":      payload.To,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TestConnectionHandler handles SMTP connection testing
func (h *HTTPHandlers) TestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.notificationService.TestEmailConnection(); err != nil {
		log.Printf("SMTP connection test failed: %v", err)
		response := map[string]interface{}{
			"status":  "FAIL",
			"message": "SMTP connection failed",
			"error":   err.Error(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"status":  "SUCCESS",
		"message": "SMTP connection successful",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ErrorHandler handles application errors
func (h *HTTPHandlers) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("HTTP Error: %v", err)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	
	response := map[string]interface{}{
		"error":   "Internal server error",
		"message": err.Error(),
	}
	
	json.NewEncoder(w).Encode(response)
}