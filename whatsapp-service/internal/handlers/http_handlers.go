package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"whatsapp-service/internal/models"
	"whatsapp-service/internal/services"
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
		"service": "whatsapp-service",
	}
	
	json.NewEncoder(w).Encode(response)
}

// WebhookHandler handles WhatsApp webhook requests
func (h *HTTPHandlers) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleWebhookVerification(w, r)
	case http.MethodPost:
		h.handleWebhookNotification(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWebhookVerification handles webhook verification
func (h *HTTPHandlers) handleWebhookVerification(w http.ResponseWriter, r *http.Request) {
	verifyToken := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if verifyToken == "" || challenge == "" {
		http.Error(w, "Missing verification parameters", http.StatusBadRequest)
		return
	}

	responseChallenge, err := h.notificationService.VerifyWebhook(verifyToken, challenge)
	if err != nil {
		log.Printf("Webhook verification failed: %v", err)
		http.Error(w, "Verification failed", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseChallenge))
}

// handleWebhookNotification handles webhook notifications from WhatsApp
func (h *HTTPHandlers) handleWebhookNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Log webhook data for debugging
	log.Printf("Received webhook notification: %s", string(body))

	// TODO: Process webhook notifications (message status updates, delivery receipts, etc.)
	// This could include updating message status in database, handling user responses, etc.

	w.WriteHeader(http.StatusOK)
}

// SendTestHandler handles direct test message requests
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
		log.Printf("Failed to send test message: %v", err)
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":         "sent",
		"notificationId": payload.NotificationID,
		"message":        "Test message sent successfully",
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