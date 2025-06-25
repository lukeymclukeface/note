package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "note-cli",
	Short: "A simple note-taking CLI application",
	Long: `Note CLI is a command-line application for managing notes.
You can create, list, edit, and delete notes from the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Note CLI!")
		fmt.Println("Use 'note-cli help' to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
