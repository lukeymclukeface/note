package helpers

import (
	"database/sql"
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/constants"
	"note-cli/internal/database"
)

// LoadConfigWithValidation loads configuration and validates required fields
func LoadConfigWithValidation() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

// LoadConfigAndDatabase loads config and opens database connection
func LoadConfigAndDatabase() (*config.Config, *sql.DB, error) {
	cfg, err := LoadConfigWithValidation()
	if err != nil {
		return nil, nil, err
	}

	if cfg.DatabasePath == "" {
		return nil, nil, fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return cfg, db, nil
}

// GetDatabaseConnection opens a database connection using the default path
func GetDatabaseConnection() (*sql.DB, error) {
	dbPath, err := constants.GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}
	return database.Connect(dbPath)
}

// ValidateOpenAIKey checks if OpenAI API key is configured
func ValidateOpenAIKey(cfg *config.Config) error {
	if cfg.OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not configured. Please run 'note setup' first")
	}
	return nil
}
