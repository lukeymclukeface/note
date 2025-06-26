package cmd

import (
	"database/sql"
	"fmt"
	"note-cli/internal/constants"
	"note-cli/internal/database"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long: `Display a list of all stored notes with their titles, 
creation dates, and preview of content.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get tag filter from flag
		tagFilter, _ := cmd.Flags().GetString("tag")
		
		if tagFilter != "" {
			fmt.Printf("Listing notes with tag '%s':\n", tagFilter)
		} else {
			fmt.Println("Listing all notes:")
		}
		fmt.Println("==================")
		
		db, err := OpenDatabase()
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			return
		}
		defer db.Close()

		notes, err := ListNotes(db, tagFilter)
		if err != nil {
			fmt.Printf("Error retrieving notes: %v\n", err)
			return
		}

		if len(notes) == 0 {
			if tagFilter != "" {
				fmt.Printf("No notes found with tag '%s'. Use 'note create' to add notes!\n", tagFilter)
			} else {
				fmt.Println("No notes found. Use 'note create' to add your first note!")
			}
			return
		}

		for i, note := range notes {
			fmt.Printf("[%d] %s\n", note.ID, note.Title)
			if note.Tags != "" {
				fmt.Printf("    Tags: %s\n", note.Tags)
			}
			fmt.Printf("    Created: %s\n", note.CreatedAt)
			
			// Show content preview (first 100 characters)
			contentPreview := note.Content
			if len(contentPreview) > 100 {
				contentPreview = contentPreview[:100] + "..."
			}
			fmt.Printf("    Preview: %s\n", contentPreview)
			
			if i < len(notes)-1 {
				fmt.Println()
			}
		}
		
		fmt.Printf("\nTotal: %d note(s)\n", len(notes))
	},
}

// OpenDatabase opens a connection to the database
func OpenDatabase() (*sql.DB, error) {
	dbPath, err := constants.GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}
	return database.Connect(dbPath)
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
