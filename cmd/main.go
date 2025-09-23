package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"paperwork-service/internal/config"
	"paperwork-service/internal/handlers"
	"paperwork-service/internal/middleware"
	"paperwork-service/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := initLogger(cfg)
	defer logger.Sync()

	logger.Info("Starting Art Battle Paperwork Service",
		zap.String("version", "1.0.0"),
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.Port))

	// Initialize services
	eventService := services.NewEventService(logger, cfg.SupabaseURL)
	pdfService := services.NewPaperworkPDFService(logger, cfg.TemplatesPath)

	// Initialize handlers
	paperworkHandler := handlers.NewPaperworkHandler(logger, eventService, pdfService)

	// Setup router
	router := setupRouter(logger, paperworkHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second, // Longer timeout for PDF generation
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}
}

// initLogger initializes the logger based on environment
func initLogger(cfg *config.Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	if cfg.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger
}

// setupRouter configures the HTTP router with all routes and middleware
func setupRouter(logger *zap.Logger, paperworkHandler *handlers.PaperworkHandler) http.Handler {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/api/v1/health", paperworkHandler.HealthCheck).Methods("GET")

	// Public paperwork generation endpoint
	router.HandleFunc("/api/v1/event-pdf/{eid}", paperworkHandler.GenerateEventPaperwork).Methods("GET")

	// Root redirect
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/health", http.StatusTemporaryRedirect)
	}).Methods("GET")

	// Apply middleware
	corsMiddleware := middleware.CORSMiddleware()
	loggingMiddleware := middleware.LoggingMiddleware(logger)

	// Wrap router with middleware
	handler := corsMiddleware.Handler(router)
	handler = loggingMiddleware(handler)

	return handler
}