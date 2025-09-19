package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"email-service/internal/config"
	"email-service/internal/handlers"
	"email-service/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Email Service with SMTP Host: %s", cfg.Email.SMTPHost)

	// Initialize services
	notificationService := services.NewNotificationService(cfg)
	defer notificationService.Close()

	// Test SMTP connection on startup
	if err := notificationService.TestEmailConnection(); err != nil {
		log.Printf("Warning: SMTP connection test failed: %v", err)
		log.Println("Service will start anyway, but emails may fail to send")
	} else {
		log.Println("SMTP connection test successful")
	}

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

		log.Println("Email service stopped")
	})
}

// setupHTTPServer configures and returns HTTP server
func setupHTTPServer(cfg *config.Config, handlers *handlers.HTTPHandlers) *http.Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/actuator/health", handlers.HealthCheck)
	mux.HandleFunc("/send-test", handlers.SendTestHandler)
	mux.HandleFunc("/test-connection", handlers.TestConnectionHandler)

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