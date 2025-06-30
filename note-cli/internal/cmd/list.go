package cmd

import (
	"database/sql"
	"fmt"
	"note-cli/internal/constants"
	"note-cli/internal/database"
	"note-cli/internal/helpers"
	"note-cli/internal/services"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long: `Display a list of all stored notes with their titles, 
creation dates, and preview of content.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInteractiveList(cmd); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// OpenDatabase opens a connection to the database
func OpenDatabase() (*sql.DB, error) {
	return helpers.GetDatabaseConnection()
}

// ListNotes retrieves notes from the database
func ListNotes(db *sql.DB, tag string) ([]database.Note, error) {
	return database.ListNotes(db, tag)
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags for the list command
	listCmd.Flags().StringP("tag", "t", "", "Filter notes by tag")
	listCmd.Flags().BoolP("recent", "r", false, "Show only recent notes")
}

func runInteractiveList(cmd *cobra.Command) error {
	// Get tag filter from flag
	tagFilter, _ := cmd.Flags().GetString("tag")

	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	notes, err := ListNotes(db, tagFilter)
	if err != nil {
		return fmt.Errorf("failed to retrieve notes: %w", err)
	}

	if len(notes) == 0 {
		if tagFilter != "" {
			fmt.Printf("📭 No notes found with tag '%s'\n", tagFilter)
		} else {
			fmt.Println("📭 No notes found")
		}
		fmt.Println("Use 'note create' to add your first note!")
		return nil
	}

	// Create options for note selection
	var options []huh.Option[int]
	for _, note := range notes {
		label := fmt.Sprintf("%s\n    %s", note.Title, note.CreatedAt)
		options = append(options, huh.NewOption(label, note.ID))
	}

	// Add "Exit" option
	options = append(options, huh.NewOption("← Exit", -1))

	// Prompt user to select a note
	var selectedID int
	title := "Select a note:"
	if tagFilter != "" {
		title = fmt.Sprintf("Select a note (filtered by tag '%s'):", tagFilter)
	}

	err = huh.NewSelect[int]().
		Title(title).
		Options(options...).
		Value(&selectedID).
		Height(15).
		Run()

	if err != nil {
		return fmt.Errorf("failed to select note: %w", err)
	}

	// Handle exit selection
	if selectedID == -1 {
		return nil
	}

	// Find the selected note
	var selectedNote *database.Note
	for _, note := range notes {
		if note.ID == selectedID {
			selectedNote = &note
			break
		}
	}

	if selectedNote == nil {
		return fmt.Errorf("selected note not found")
	}

	// Show note actions
	return showNoteActions(selectedNote)
}

func showNoteActions(note *database.Note) error {
	fmt.Printf("\n📝 Selected: %s\n", note.Title)
	fmt.Printf("🏷️ Tags: %s\n", note.Tags)
	fmt.Printf("📅 Created: %s\n", note.CreatedAt)
	fmt.Println()

	// Check if this is an imported audio note
	isAudioNote := strings.Contains(note.Tags, "imported") && strings.Contains(note.Tags, "audio")

	var actionOptions []huh.Option[string]

	if isAudioNote {
		// For audio notes, offer different options
		actionOptions = []huh.Option[string]{
			huh.NewOption("📂 Open folder in Finder", "open_folder"),
			huh.NewOption("📄 View summary", "view_summary"),
			huh.NewOption("📜 View transcription", "view_transcript"),
			huh.NewOption("🗑️ Delete note", "delete"),
			huh.NewOption("← Back to list", "back"),
		}
	} else {
		// For regular notes
		actionOptions = []huh.Option[string]{
			huh.NewOption("📄 View content", "view_content"),
			huh.NewOption("🗑️ Delete note", "delete"),
			huh.NewOption("← Back to list", "back"),
		}
	}

	var selectedAction string
	err := huh.NewSelect[string]().
		Title("What would you like to do?").
		Options(actionOptions...).
		Value(&selectedAction).
		Run()

	if err != nil {
		return fmt.Errorf("failed to select action: %w", err)
	}

	switch selectedAction {
	case "back":
		return nil
	case "view_content":
		return viewNoteContent(note)
	case "view_summary":
		return viewAudioNoteSummary(note)
	case "view_transcript":
		return viewAudioNoteTranscript(note)
	case "open_folder":
		return openNoteFolder(note)
	case "delete":
		return deleteNote(note)
	default:
		return fmt.Errorf("unknown action: %s", selectedAction)
	}
}

func viewNoteContent(note *database.Note) error {
	fmt.Printf("\n📝 %s\n", note.Title)
	fmt.Printf("🏷️ Tags: %s\n", note.Tags)
	fmt.Printf("📅 Created: %s\n", note.CreatedAt)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(note.Content)
	fmt.Println(strings.Repeat("=", 50))

	// Wait for user to press enter
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
	return nil
}

func viewAudioNoteSummary(note *database.Note) error {
	// Extract folder name from content
	fileService := services.NewFileService()
	folderName := fileService.ExtractFolderFromContent(note.Content)
	if folderName == "" {
		return fmt.Errorf("could not determine note folder")
	}

	notesDir, err := constants.GetNotesDir()
	if err != nil {
		return fmt.Errorf("failed to get notes directory: %w", err)
	}

	summaryPath := filepath.Join(notesDir, folderName, "summary.md")

	// Check if summary file exists
	if _, err := os.Stat(summaryPath); os.IsNotExist(err) {
		return fmt.Errorf("summary file not found at: %s", summaryPath)
	}

	// Read and display summary
	content, err := os.ReadFile(summaryPath)
	if err != nil {
		return fmt.Errorf("failed to read summary: %w", err)
	}

	fmt.Printf("\n📄 Summary: %s\n", note.Title)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(string(content))
	fmt.Println(strings.Repeat("=", 50))

	// Wait for user to press enter
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
	return nil
}

func viewAudioNoteTranscript(note *database.Note) error {
	// Extract folder name from content
	fileService := services.NewFileService()
	folderName := fileService.ExtractFolderFromContent(note.Content)
	if folderName == "" {
		return fmt.Errorf("could not determine note folder")
	}

	notesDir, err := constants.GetNotesDir()
	if err != nil {
		return fmt.Errorf("failed to get notes directory: %w", err)
	}

	transcriptPath := filepath.Join(notesDir, folderName, "transcription.md")

	// Check if transcript file exists
	if _, err := os.Stat(transcriptPath); os.IsNotExist(err) {
		return fmt.Errorf("transcription file not found at: %s", transcriptPath)
	}

	// Read and display transcript
	content, err := os.ReadFile(transcriptPath)
	if err != nil {
		return fmt.Errorf("failed to read transcription: %w", err)
	}

	fmt.Printf("\n📜 Transcription: %s\n", note.Title)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println(string(content))
	fmt.Println(strings.Repeat("=", 50))

	// Wait for user to press enter
	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
	return nil
}

func openNoteFolder(note *database.Note) error {
	// Extract folder name from content
	fileService := services.NewFileService()
	folderName := fileService.ExtractFolderFromContent(note.Content)
	if folderName == "" {
		return fmt.Errorf("could not determine note folder")
	}

	notesDir, err := constants.GetNotesDir()
	if err != nil {
		return fmt.Errorf("failed to get notes directory: %w", err)
	}

	folderPath := filepath.Join(notesDir, folderName)

	// Check if folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return fmt.Errorf("folder not found at: %s", folderPath)
	}

	// Open folder in Finder (macOS)
	cmd := exec.Command("open", folderPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open folder: %w", err)
	}

	fmt.Printf("📂 Opened folder: %s\n", folderPath)
	return nil
}

func deleteNote(note *database.Note) error {
	// Confirm deletion
	var confirm bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete '%s'?", note.Title)).
		Description("This action cannot be undone.").
		Affirmative("Yes, delete").
		Negative("Cancel").
		Value(&confirm).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirm {
		fmt.Println("❌ Deletion cancelled")
		return nil
	}

	// Delete from database
	db, err := OpenDatabase()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	query := "DELETE FROM notes WHERE id = ?"
	result, err := db.Exec(query, note.ID)
	if err != nil {
		return fmt.Errorf("failed to delete note from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found in database")
	}

	// For audio notes, ask if they want to delete the folder too
	isAudioNote := strings.Contains(note.Tags, "imported") && strings.Contains(note.Tags, "audio")
	if isAudioNote {
		fileService := services.NewFileService()
		folderName := fileService.ExtractFolderFromContent(note.Content)
		if folderName != "" {
			notesDir, err := constants.GetNotesDir()
			if err == nil {
				folderPath := filepath.Join(notesDir, folderName)
				if _, err := os.Stat(folderPath); err == nil {
					var deleteFolder bool
					err = huh.NewConfirm().
						Title("Also delete the audio files and folder?").
						Description(fmt.Sprintf("This will permanently delete: %s", folderPath)).
						Affirmative("Yes, delete folder").
						Negative("No, keep files").
						Value(&deleteFolder).
						Run()

					if err == nil && deleteFolder {
						if err := os.RemoveAll(folderPath); err != nil {
							fmt.Printf("⚠️  Warning: Failed to delete folder %s: %v\n", folderPath, err)
						} else {
							fmt.Printf("🗑️  Deleted folder: %s\n", folderPath)
						}
					}
				}
			}
		}
	}

	fmt.Printf("✅ Note '%s' deleted successfully\n", note.Title)
	return nil
}
