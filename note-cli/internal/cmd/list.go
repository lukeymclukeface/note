package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long: `Display a list of all stored notes with their titles, 
creation dates, and preview of content.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all notes:")
		fmt.Println("==================")
		
		// TODO: Implement actual note retrieval logic
		fmt.Println("No notes found. Use 'note create' to add your first note!")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	
	// Add flags for the list command
	listCmd.Flags().StringP("tag", "t", "", "Filter notes by tag")
	listCmd.Flags().BoolP("recent", "r", false, "Show only recent notes")
}
