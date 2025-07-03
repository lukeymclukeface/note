package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// AppConfig represents the application configuration stored in JSON
type AppConfig struct {
	// OpenAI Configuration
	OpenAIKey string `json:"openai_key,omitempty"`
	
	// Other AI providers can be added here
	// GoogleProjectID string `json:"google_project_id,omitempty"`
	// GoogleLocation   string `json:"google_location,omitempty"`
	
	// Model configurations
	TranscriptionProvider string `json:"transcription_provider,omitempty"`
	TranscriptionModel    string `json:"transcription_model,omitempty"`
	SummaryProvider       string `json:"summary_provider,omitempty"`
	SummaryModel          string `json:"summary_model,omitempty"`
}

// ConfigManager handles application configuration persistence
type ConfigManager struct {
	configPath string
	config     *AppConfig
	mutex      sync.RWMutex
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	// Use the same path as the CLI for consistency
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir not available
		homeDir = "."
	}
	
	configDir := filepath.Join(homeDir, ".noteai")
	configPath := filepath.Join(configDir, "config.json")
	
	return &ConfigManager{
		configPath: configPath,
		config:     &AppConfig{},
	}
}

// Load loads the configuration from the JSON file
func (cm *ConfigManager) Load() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// File doesn't exist, start with empty config
		cm.config = &AppConfig{}
		return nil
	}
	
	// Read the config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse JSON
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	cm.config = &config
	return nil
}

// Save saves the current configuration to the JSON file
func (cm *ConfigManager) Save() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	// Ensure config directory exists
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal config to JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigManager) GetConfig() AppConfig {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	config := *cm.config
	return config
}

// SetConfig updates the configuration
func (cm *ConfigManager) SetConfig(config AppConfig) error {
	cm.mutex.Lock()
	cm.config = &config
	cm.mutex.Unlock()
	
	return cm.Save()
}

// GetOpenAIKey returns the OpenAI API key if configured
func (cm *ConfigManager) GetOpenAIKey() string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	return cm.config.OpenAIKey
}

// SetOpenAIKey updates the OpenAI API key
func (cm *ConfigManager) SetOpenAIKey(key string) error {
	cm.mutex.Lock()
	cm.config.OpenAIKey = key
	cm.mutex.Unlock()
	
	return cm.Save()
}

// HasOpenAIKey returns true if OpenAI key is configured
func (cm *ConfigManager) HasOpenAIKey() bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	return cm.config.OpenAIKey != ""
}

// Global instance for easy access
var defaultManager *ConfigManager
var once sync.Once

// GetManager returns the default configuration manager instance
func GetManager() *ConfigManager {
	once.Do(func() {
		defaultManager = NewConfigManager()
		// Load existing config on first access
		_ = defaultManager.Load() // Ignore error, will use empty config
	})
	return defaultManager
}
