package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/database"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var recordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "Manage audio recordings",
	Long:  `Manage audio recordings with subcommands to list and delete recordings.`,
}

var recordingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all audio recordings",
	Long:  `Display a list of all audio recordings with their details including duration, file size, and creation time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listRecordings(); err != nil {
			fmt.Printf("❌ Error listing recordings: %v\n", err)
			os.Exit(1)
		}
	},
}

var recordingsDeleteCmd = &cobra.Command{
	Use:   "delete [recording-id]",
	Short: "Delete an audio recording",
	Long:  `Delete an audio recording by ID. If no ID is provided, you'll be prompted to select one from a list.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteRecording(args); err != nil {
			fmt.Printf("❌ Error deleting recording: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(recordingsCmd)
	recordingsCmd.AddCommand(recordingsListCmd)
	recordingsCmd.AddCommand(recordingsDeleteCmd)
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

func deleteRecording(args []string) error {
	// Load config to get database path
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.DatabasePath == "" {
		return fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	var recordingID int

	// If ID provided as argument, use it
	if len(args) > 0 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid recording ID: %s", args[0])
		}
		recordingID = id
	} else {
		// Otherwise, prompt user to select from list
		recordings, err := database.ListRecordings(cfg.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to retrieve recordings: %w", err)
		}

		if len(recordings) == 0 {
			fmt.Println("📭 No recordings found to delete.")
			return nil
		}

		// Create options for selection
		var options []huh.Option[int]
		for _, recording := range recordings {
			// Check if file exists
			exists := "✅"
			if _, err := os.Stat(recording.FilePath); os.IsNotExist(err) {
				exists = "❌"
			}
			
			duration := recording.Duration.Round(time.Second)
			createdStr := recording.CreatedAt.Format("2006-01-02 15:04")
			label := fmt.Sprintf("%s %s (ID: %d, %v, %s)", exists, recording.Filename, recording.ID, duration, createdStr)
			
			options = append(options, huh.NewOption(label, recording.ID))
		}

		// Prompt user to select recording
		var selectedID int
		err = huh.NewSelect[int]().
			Title("Select recording to delete:").
			Options(options...).
			Value(&selectedID).
			Run()

		if err != nil {
			return fmt.Errorf("failed to select recording: %w", err)
		}

		recordingID = selectedID
	}

	// Delete the recording from database
	deletedRecording, err := database.DeleteRecording(cfg.DatabasePath, recordingID)
	if err != nil {
		return fmt.Errorf("failed to delete recording: %w", err)
	}

	// Ask if user wants to delete the file as well
	var deleteFile bool
	fileExists := true
	if _, err := os.Stat(deletedRecording.FilePath); os.IsNotExist(err) {
		fileExists = false
	}

	if fileExists {
		err = huh.NewConfirm().
			Title(fmt.Sprintf("Also delete the audio file '%s'?", deletedRecording.Filename)).
			Description("This will permanently remove the file from your system.").
			Affirmative("Yes, delete file").
			Negative("No, keep file").
			Value(&deleteFile).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if deleteFile {
			if err := os.Remove(deletedRecording.FilePath); err != nil {
				fmt.Printf("⚠️  Warning: Failed to delete file %s: %v\n", deletedRecording.FilePath, err)
			} else {
				fmt.Printf("🗑️  Deleted file: %s\n", deletedRecording.FilePath)
			}
		}
	}

	fmt.Printf("✅ Recording '%s' (ID: %d) deleted from database\n", deletedRecording.Filename, deletedRecording.ID)

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
