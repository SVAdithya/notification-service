package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whatsapp-service/internal/config"
	"whatsapp-service/internal/handlers"
	"whatsapp-service/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting WhatsApp Service with Phone Number ID: %s", cfg.WhatsApp.PhoneNumberID)

	// Initialize services
	notificationService := services.NewNotificationService(cfg)
	defer notificationService.Close()

	// Initialize HTTP handlers
	httpHandlers := handlers.NewHTTPHandlers(notificationService)

	// Start Kafka consumer in background
	consumerCtx, cancelConsumer := context.WithCancel(context.Background())
	go func() {
		if err := notificationService.StartConsumer(consumerCtx); err != nil && err != context.Canceled {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Setup HTTP server
	server := setupHTTPServer(cfg, httpHandlers)

	// Start server in background
	go func() {
		log.Printf("HTTP server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	waitForShutdown(func() {
		log.Println("Shutting down gracefully...")

		// Cancel Kafka consumer
		cancelConsumer()

		// Shutdown HTTP server
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		log.Println("Service stopped")
	})
}

// setupHTTPServer configures and returns HTTP server
func setupHTTPServer(cfg *config.Config, handlers *handlers.HTTPHandlers) *http.Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/actuator/health", handlers.HealthCheck)
	mux.HandleFunc("/webhook", handlers.WebhookHandler)
	mux.HandleFunc("/send-test", handlers.SendTestHandler)

	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
}

// waitForShutdown waits for shutdown signal and executes cleanup function
func waitForShutdown(cleanup func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %s", sig)

	cleanup()
}

func init() {
	// Set log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Set log level based on environment
	if level := os.Getenv("LOG_LEVEL"); level == "debug" {
		log.SetOutput(os.Stdout)
	}
}