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
	"github.com/fatih/color"
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
	Short: "Configure OpenAI models for transcription and summaries",
	Long:  `Query OpenAI for available models and select models for transcription and summaries.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configureModels(); err != nil {
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

	// Color styling
	// blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	// fmt.Printf("%s\n\n", blue("📋 Current Configuration"))
	fmt.Println("")

	// Define consistent field width for alignment
	fieldWidth := 22

	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Editor:")), green(cfg.Editor))
	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Date Format:")), green(cfg.DateFormat))
	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Default Tags:")), green(strings.Join(cfg.DefaultTags, ", ")))

	// Mask the OpenAI key for security
	var keyDisplay string
	if cfg.OpenAIKey != "" {
		if len(cfg.OpenAIKey) > 8 {
			keyDisplay = cfg.OpenAIKey[:4] + "..." + cfg.OpenAIKey[len(cfg.OpenAIKey)-4:]
		} else {
			keyDisplay = strings.Repeat("*", len(cfg.OpenAIKey))
		}
	} else {
		keyDisplay = gray("(not set)")
	}
	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "OpenAI Key:")), green(keyDisplay))

	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Transcription Model:")), green(cfg.TranscriptionModel))
	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Summary Model:")), green(cfg.SummaryModel))

	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Notes Directory:")), green(cfg.NotesDir))
	fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%-*s", fieldWidth, "Database Path:")), green(cfg.DatabasePath))

	fmt.Printf("\n%s %s\n", gray("Config file:"), gray(config.ConfigPath()))
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

func configureModels() error {
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

	// Filter models for different purposes
	transcriptionModels := filterTranscriptionModels(models)
	chatModels := filterChatModels(models)

	if len(transcriptionModels) == 0 {
		return fmt.Errorf("no suitable transcription models found")
	}
	if len(chatModels) == 0 {
		return fmt.Errorf("no suitable chat models found")
	}

	// Sort models by name for consistent display
	sort.Slice(transcriptionModels, func(i, j int) bool {
		return transcriptionModels[i].ID < transcriptionModels[j].ID
	})
	sort.Slice(chatModels, func(i, j int) bool {
		return chatModels[i].ID < chatModels[j].ID
	})

	// Show current models
	fmt.Printf("Current transcription model: %s\n", cfg.TranscriptionModel)
	fmt.Printf("Current summary model: %s\n\n", cfg.SummaryModel)

	// Configure transcription model
	var transcriptionOptions []huh.Option[string]
	for _, model := range transcriptionModels {
		label := fmt.Sprintf("%s", model.ID)
		if model.OwnedBy != "" {
			label += fmt.Sprintf(" (by %s)", model.OwnedBy)
		}
		transcriptionOptions = append(transcriptionOptions, huh.NewOption(label, model.ID))
	}

	var selectedTranscriptionModel string
	err = huh.NewSelect[string]().
		Title("Select a model for audio transcription:").
		Description("Choose a model for converting audio to text.").
		Options(transcriptionOptions...).
		Value(&selectedTranscriptionModel).
		Height(10).
		Run()

	if err != nil {
		return fmt.Errorf("failed to select transcription model: %w", err)
	}

	// Configure summary model
	var summaryOptions []huh.Option[string]
	for _, model := range chatModels {
		label := fmt.Sprintf("%s", model.ID)
		if model.OwnedBy != "" {
			label += fmt.Sprintf(" (by %s)", model.OwnedBy)
		}
		summaryOptions = append(summaryOptions, huh.NewOption(label, model.ID))
	}

	var selectedSummaryModel string
	err = huh.NewSelect[string]().
		Title("Select a model for text summarization:").
		Description("Choose a model for generating summaries from transcribed text.").
		Options(summaryOptions...).
		Value(&selectedSummaryModel).
		Height(15).
		Run()

	if err != nil {
		return fmt.Errorf("failed to select summary model: %w", err)
	}

	// Update config with selected models
	cfg.TranscriptionModel = selectedTranscriptionModel
	cfg.SummaryModel = selectedSummaryModel

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ Transcription model updated to: %s\n", selectedTranscriptionModel)
	fmt.Printf("✅ Summary model updated to: %s\n", selectedSummaryModel)
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

func filterTranscriptionModels(models []OpenAIModel) []OpenAIModel {
	var transcriptionModels []OpenAIModel

	for _, model := range models {
		// Filter for audio/speech models suitable for transcription
		// Include Whisper models and other audio-related models
		if strings.Contains(model.ID, "whisper") ||
			strings.Contains(model.ID, "audio") ||
			strings.Contains(model.ID, "speech") ||
			strings.Contains(model.ID, "transcribe") {
			transcriptionModels = append(transcriptionModels, model)
		}
	}

	// If no specific audio models found, also include some general models that can handle audio
	if len(transcriptionModels) == 0 {
		for _, model := range models {
			// Add some general models that might work for audio processing
			if strings.Contains(model.ID, "gpt-4") ||
				strings.Contains(model.ID, "gpt-3.5") {
				transcriptionModels = append(transcriptionModels, model)
			}
		}
	}

	// Always ensure whisper-1 is available as a fallback
	hasWhisper := false
	for _, model := range transcriptionModels {
		if model.ID == "whisper-1" {
			hasWhisper = true
			break
		}
	}
	if !hasWhisper {
		transcriptionModels = append([]OpenAIModel{{
			ID:      "whisper-1",
			Object:  "model",
			OwnedBy: "openai",
		}}, transcriptionModels...)
	}

	return transcriptionModels
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
