package config

import (
	"encoding/json"
	"fmt"
	"note-cli/internal/constants"
	"os"
)

// Config represents the application configuration
type Config struct {
	NotesDir               string   `json:"notes_dir"`
	Editor                 string   `json:"editor"`
	DateFormat             string   `json:"date_format"`
	DefaultTags            []string `json:"default_tags"`
	OpenAIKey              string   `json:"openai_key"`
	DatabasePath           string   `json:"database_path"`
	TranscriptionModel     string   `json:"transcription_model"`
	SummaryModel           string   `json:"summary_model"`
	// Provider configuration
	TranscriptionProvider  string   `json:"transcription_provider"`
	SummaryProvider        string   `json:"summary_provider"`
	// Google AI configuration
	GoogleProjectID        string   `json:"google_project_id"`
	GoogleLocation         string   `json:"google_location"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	notesDir, _ := constants.GetNotesDir()
	return &Config{
		NotesDir:              notesDir,
		Editor:                "nano",
		DateFormat:            "2006-01-02",
		DefaultTags:           []string{},
		OpenAIKey:             "",
		DatabasePath:          "",
		TranscriptionModel:    "whisper-1",
		SummaryModel:          "gpt-3.5-turbo",
		TranscriptionProvider: "openai",
		SummaryProvider:       "openai",
		GoogleProjectID:       "",
		GoogleLocation:        "us-central1",
	}
}

// ConfigDir returns the configuration directory path
func ConfigDir() string {
	baseDir, _ := constants.GetBaseDir()
	return baseDir
}

// ConfigPath returns the configuration file path
func ConfigPath() string {
	configPath, _ := constants.GetConfigPath()
	return configPath
}

// Load reads the configuration from the config file
func Load() (*Config, error) {
	configPath := ConfigPath()
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		config := DefaultConfig()
		if err := Save(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}
	
	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Migration: Set default values for new fields if they're empty
	if config.TranscriptionModel == "" {
		config.TranscriptionModel = "whisper-1"
	}
	if config.SummaryModel == "" {
		config.SummaryModel = "gpt-3.5-turbo"
	}
	if config.TranscriptionProvider == "" {
		config.TranscriptionProvider = "openai"
	}
	if config.SummaryProvider == "" {
		config.SummaryProvider = "openai"
	}
	if config.GoogleLocation == "" {
		config.GoogleLocation = "us-central1"
	}
	
	return &config, nil
}

// Save writes the configuration to the config file
func Save(config *Config) error {
	configDir := ConfigDir()
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to config file
	configPath := ConfigPath()
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// EnsureConfigExists ensures the config directory and file exist
func EnsureConfigExists() error {
	_, err := Load()
	return err
}
