package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"note-cli/internal/config"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long: `Manage the note application configuration stored in ~/.noteai/config.json.
Use subcommands to view, edit, or reset configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show current config when no subcommand is provided
		showConfig()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
  - notes_dir: Directory where notes are stored
  - editor: Default editor for editing notes
  - date_format: Date format for timestamps (Go time format)
  - default_tags: Comma-separated list of default tags
  - openai_key: OpenAI API key for transcription and summarization`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		setConfig(key, value)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Reset the configuration file to default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		resetConfig()
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Display the path to the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.ConfigPath())
	},
}

var configModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Configure OpenAI model for summaries",
	Long:  `Query OpenAI for available models and select one for generating summaries.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configureModel(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Add subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configModelCmd)
}

func showConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")
	fmt.Printf("Notes Directory: %s\n", cfg.NotesDir)
	fmt.Printf("Editor: %s\n", cfg.Editor)
	fmt.Printf("Date Format: %s\n", cfg.DateFormat)
	fmt.Printf("Default Tags: %s\n", strings.Join(cfg.DefaultTags, ", "))
	
	// Mask the OpenAI key for security
	if cfg.OpenAIKey != "" {
		if len(cfg.OpenAIKey) > 8 {
			fmt.Printf("OpenAI Key: %s...%s\n", cfg.OpenAIKey[:4], cfg.OpenAIKey[len(cfg.OpenAIKey)-4:])
		} else {
			fmt.Printf("OpenAI Key: %s\n", strings.Repeat("*", len(cfg.OpenAIKey)))
		}
	} else {
		fmt.Printf("OpenAI Key: (not set)\n")
	}
	
	fmt.Printf("AI Model: %s\n", cfg.AIModel)
	fmt.Printf("Database Path: %s\n", cfg.DatabasePath)
	
	fmt.Printf("\nConfig file: %s\n", config.ConfigPath())
}

func setConfig(key, value string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	
	switch key {
	case "notes_dir":
		cfg.NotesDir = value
	case "editor":
		cfg.Editor = value
	case "date_format":
		cfg.DateFormat = value
	case "default_tags":
		// Split comma-separated tags and trim whitespace
		tags := strings.Split(value, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		cfg.DefaultTags = tags
	case "openai_key":
		cfg.OpenAIKey = value
	default:
		fmt.Printf("Unknown configuration key: %s\n", key)
		fmt.Println("Available keys: notes_dir, editor, date_format, default_tags, openai_key")
		return
	}
	
	if err := config.Save(cfg); err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		return
	}
	
	fmt.Printf("Configuration updated: %s = %s\n", key, value)
}

func resetConfig() {
	cfg := config.DefaultConfig()
	if err := config.Save(cfg); err != nil {
		fmt.Printf("Error resetting configuration: %v\n", err)
		return
	}
	
	fmt.Println("Configuration reset to defaults.")
	showConfig()
}

// OpenAI API structures
type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type OpenAIModelsResponse struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}

func configureModel() error {
	// Load current config to get API key
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not configured. Please run 'note setup' or 'note config set openai_key <your-key>' first")
	}

	fmt.Println("🔍 Fetching available OpenAI models...")

	// Query OpenAI for available models
	models, err := fetchOpenAIModels(cfg.OpenAIKey)
	if err != nil {
		return fmt.Errorf("failed to fetch models: %w", err)
	}

	// Filter for chat models suitable for summaries
	chatModels := filterChatModels(models)
	if len(chatModels) == 0 {
		return fmt.Errorf("no suitable chat models found")
	}

	// Sort models by name for consistent display
	sort.Slice(chatModels, func(i, j int) bool {
		return chatModels[i].ID < chatModels[j].ID
	})

	// Create options for selection
	var options []huh.Option[string]
	for _, model := range chatModels {
		// Create a descriptive label
		label := fmt.Sprintf("%s", model.ID)
		if model.OwnedBy != "" {
			label += fmt.Sprintf(" (by %s)", model.OwnedBy)
		}
		options = append(options, huh.NewOption(label, model.ID))
	}

	// Show current model
	fmt.Printf("Current model: %s\n\n", cfg.AIModel)

	// Prompt user to select model
	var selectedModel string
	err = huh.NewSelect[string]().
		Title("Select an OpenAI model for generating summaries:").
		Description("Choose a model that will be used for summarizing transcribed audio.").
		Options(options...).
		Value(&selectedModel).
		Height(15).
		Run()

	if err != nil {
		return fmt.Errorf("failed to select model: %w", err)
	}

	// Update config with selected model
	cfg.AIModel = selectedModel
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ Model updated to: %s\n", selectedModel)
	return nil
}

func fetchOpenAIModels(apiKey string) ([]OpenAIModel, error) {
	url := "https://api.openai.com/v1/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var modelsResponse OpenAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return modelsResponse.Data, nil
}

func filterChatModels(models []OpenAIModel) []OpenAIModel {
	var chatModels []OpenAIModel

	for _, model := range models {
		// Filter for chat completion models
		// Include GPT models that are suitable for chat completions
		if strings.Contains(model.ID, "gpt-") && 
		   !strings.Contains(model.ID, "instruct") &&
		   !strings.Contains(model.ID, "embedding") &&
		   !strings.Contains(model.ID, "whisper") &&
		   !strings.Contains(model.ID, "tts") &&
		   !strings.Contains(model.ID, "dall-e") {
			chatModels = append(chatModels, model)
		}
	}

	return chatModels
}
