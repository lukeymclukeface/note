package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long: `Manage the note-cli application configuration stored in ~/.noteai/config.json.
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

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Add subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)
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
