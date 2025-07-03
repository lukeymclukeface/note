package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/your-org/note-server/internal/config"
	"github.com/your-org/note-server/internal/database"
	apphttp "github.com/your-org/note-server/internal/http"
	"github.com/your-org/note-server/internal/service"
	"github.com/your-org/note-server/internal/ws"
)

func main() {
	// Parse env config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// Setup zerolog
	setupLogger(cfg)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize database
	// Use ~/.noteai/notes.db to match frontend expectations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	dbPath := filepath.Join(homeDir, ".noteai", "notes.db")
	
	// Ensure .noteai directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}
	
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info().Str("database_path", dbPath).Msg("Database initialized")

	// Initialize services
	transcribeService := service.NewTranscribeService()
	transcribeHub := ws.NewTranscribeHub(transcribeService)

	// Start the WebSocket hub
	go transcribeHub.Run()

	// Initialize chi router with WebSocket hub
	router := apphttp.NewRouter(transcribeHub)

	// Create HTTP server
	addr := "0.0.0.0:" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	logger.Info().
		Str("address", addr).
		Str("media_dir", cfg.MediaTmpDir).
		Bool("dev_mode", cfg.DevMode).
		Msg("Starting note server")

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	logger.Info().Msgf("Server listening on %s", addr)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown WebSocket hub first
	transcribeHub.Shutdown()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		logger.Info().Msg("Server exited gracefully")
	}
}

func setupLogger(cfg *config.Config) {
	// Set log level based on config
	switch cfg.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Pretty print logs in development mode
	if cfg.DevMode {
		log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
