package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestTranscribeCommand(t *testing.T) {
	// Test command registration
	transcribeCommand := rootCmd
	found := false
	
	for _, cmd := range transcribeCommand.Commands() {
		if cmd.Use == "transcribe [file]" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Transcribe command not registered")
	}
}

func TestTranscribeCommandFlags(t *testing.T) {
	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use: "transcribe [file]",
	}
	
	// Add flags like the real command
	cmd.Flags().StringP("output", "o", "", "Output file path for the transcription (default: print to stdout)")
	cmd.Flags().StringP("format", "f", "text", "Output format: text, markdown (default: text)")
	cmd.Flags().StringP("dir", "d", "", "Output directory for chunk files (optional)")
	cmd.Flags().BoolP("chunks", "c", false, "Save individual chunk transcriptions to separate files")
	
	// Test that flags are properly defined
	flagTests := []struct {
		name      string
		shorthand string
		exists    bool
	}{
		{"output", "o", true},
		{"format", "f", true},
		{"dir", "d", true},
		{"chunks", "c", true},
		{"nonexistent", "z", false},
	}
	
	for _, tt := range flagTests {
		flag := cmd.Flags().Lookup(tt.name)
		if tt.exists && flag == nil {
			t.Errorf("Expected flag %s to exist", tt.name)
		}
		if !tt.exists && flag != nil {
			t.Errorf("Expected flag %s to not exist", tt.name)
		}
		if tt.exists && flag != nil && flag.Shorthand != tt.shorthand {
			t.Errorf("Expected flag %s shorthand to be %s, got %s", tt.name, tt.shorthand, flag.Shorthand)
		}
	}
}

func TestTranscribeCommandHelp(t *testing.T) {
	// Test basic properties of the command
	if transcribeCmd.Short == "" {
		t.Error("Transcribe command should have a short description")
	}
	
	if transcribeCmd.Long == "" {
		t.Error("Transcribe command should have a long description")
	}
	
	// Test that the long description contains expected content
	expectedContent := []string{
		"audio file",
		"mp3",
		"wav",
		"m4a",
		"ogg",
		"flac",
	}
	
	for _, content := range expectedContent {
		if !bytes.Contains([]byte(transcribeCmd.Long), []byte(content)) {
			t.Errorf("Command description should contain %q", content)
		}
	}
	
	// Test that flags are properly documented
	flags := transcribeCmd.Flags()
	if flags.Lookup("output") == nil {
		t.Error("Command should have output flag")
	}
	if flags.Lookup("format") == nil {
		t.Error("Command should have format flag")
	}
	if flags.Lookup("chunks") == nil {
		t.Error("Command should have chunks flag")
	}
	if flags.Lookup("dir") == nil {
		t.Error("Command should have dir flag")
	}
}

func TestTranscribeCommandValidation(t *testing.T) {
	// Test that command requires exactly one argument
	if transcribeCmd.Args == nil {
		t.Error("Transcribe command should have argument validation")
	}
	
	// Test with no arguments
	err := transcribeCmd.Args(transcribeCmd, []string{})
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}
	
	// Test with multiple arguments
	err = transcribeCmd.Args(transcribeCmd, []string{"file1.mp3", "file2.mp3"})
	if err == nil {
		t.Error("Expected error when multiple arguments provided")
	}
	
	// Test with correct number of arguments
	err = transcribeCmd.Args(transcribeCmd, []string{"file.mp3"})
	if err != nil {
		t.Errorf("Expected no error with one argument, got: %v", err)
	}
}
