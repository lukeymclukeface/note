package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port string `envconfig:"PORT" default:"8080"`
	Host string `envconfig:"HOST" default:"localhost"`

	// OpenAI configuration
	OpenAIKey string `envconfig:"OPENAI_KEY"`

	// Media and file handling
	MediaTmpDir string `envconfig:"MEDIA_TMP_DIR" default:"/tmp/note-media"`

	// Logging configuration
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// WebSocket configuration
	WSReadBufferSize  int `envconfig:"WS_READ_BUFFER_SIZE" default:"1024"`
	WSWriteBufferSize int `envconfig:"WS_WRITE_BUFFER_SIZE" default:"1024"`

	// Audio processing configuration
	MaxAudioDuration int    `envconfig:"MAX_AUDIO_DURATION" default:"300"` // seconds
	AudioFormat      string `envconfig:"AUDIO_FORMAT" default:"wav"`

	// Development mode
	DevMode bool `envconfig:"DEV_MODE" default:"false"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// OpenAI key is now optional and managed through JSON config
	
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if c.MediaTmpDir == "" {
		return fmt.Errorf("MEDIA_TMP_DIR cannot be empty")
	}

	return nil
}

// Address returns the full server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
