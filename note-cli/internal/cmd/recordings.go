package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/database"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var recordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "List all audio recordings",
	Long:  `Display a list of all audio recordings with their details including duration, file size, and creation time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listRecordings(); err != nil {
			fmt.Printf("❌ Error listing recordings: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(recordingsCmd)
}

func listRecordings() error {
	// Load config to get database path
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.DatabasePath == "" {
		return fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	// Get recordings from database
	recordings, err := database.ListRecordings(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to retrieve recordings: %w", err)
	}

	if len(recordings) == 0 {
		fmt.Println("📭 No recordings found.")
		fmt.Println("Use 'note record' to create your first recording.")
		return nil
	}

	fmt.Printf("🎵 Found %d recording(s):\n\n", len(recordings))

	for i, recording := range recordings {
		// Check if file still exists
		exists := "✅"
		if _, err := os.Stat(recording.FilePath); os.IsNotExist(err) {
			exists = "❌"
		}

		// Format duration
		duration := recording.Duration.Round(time.Second)
		
		// Format file size
		sizeStr := formatFileSize(recording.FileSize)
		
		// Format created date
		createdStr := recording.CreatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("%d. %s %s\n", i+1, exists, recording.Filename)
		fmt.Printf("   Duration: %v | Size: %s | Created: %s\n", duration, sizeStr, createdStr)
		fmt.Printf("   Format: %s | %d Hz | %d channel(s)\n", recording.Format, recording.SampleRate, recording.Channels)
		fmt.Printf("   Path: %s\n", recording.FilePath)
		fmt.Println()
	}

	// Show directory info
	homeDir, _ := os.UserHomeDir()
	recordingsDir := filepath.Join(homeDir, ".note-cli", "recordings")
	fmt.Printf("📁 Recordings directory: %s\n", recordingsDir)

	return nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
