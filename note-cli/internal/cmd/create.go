package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [note content]",
	Short: "Create a new note",
	Long: `Create a new note with the provided content.
If no content is provided, you'll be prompted to enter it.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := ""
		for i, arg := range args {
			if i > 0 {
				content += " "
			}
			content += arg
		}
		
		cmd.Printf("Creating note: %s\n", content)
		cmd.Printf("Created at: %s\n", time.Now().Format(time.RFC3339))
		
		// TODO: Implement actual note storage logic
		cmd.Println("Note created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	
	// Add flags for the create command
	createCmd.Flags().StringP("title", "t", "", "Title for the note")
	createCmd.Flags().StringP("tags", "", "", "Comma-separated tags for the note")
}
