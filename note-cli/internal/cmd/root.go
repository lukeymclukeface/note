package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global verbose flag
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "note",
	Short: "A simple note-taking CLI application",
	Long: `Note CLI is a command-line application for managing notes.
You can create, list, edit, and delete notes from the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Note CLI!")
		fmt.Println("Use 'note help' to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global verbose flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output with detailed information")
}

// IsVerbose returns the current state of the verbose flag
func IsVerbose() bool {
	return verbose
}
