package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestCommand creates a test root command with the create command attached
func setupTestCommand() (*cobra.Command, *bytes.Buffer) {
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(createCmd)

	output := new(bytes.Buffer)
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	return rootCmd, output
}

func TestCreateCommand_BasicNote(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", "Test note content"})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: Test note content")
	assert.Contains(t, output.String(), "Note created successfully!")

	// Check that a timestamp is included
	outputStr := output.String()
	assert.Contains(t, outputStr, "Created at:")

	// Verify timestamp format is valid
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Created at:") {
			timestampStr := strings.TrimPrefix(line, "Created at: ")
			_, err := time.Parse(time.RFC3339, timestampStr)
			assert.NoError(t, err, "Timestamp should be in RFC3339 format")
			break
		}
	}
}

func TestCreateCommand_MultipleWords(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", "This", "is", "a", "multi-word", "note"})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: This is a multi-word note")
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_WithTitleFlag(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", "--title", "My Important Note", "This is the content"})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: This is the content")
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_WithTagsFlag(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", "--tags", "work,important,meeting", "Meeting notes content"})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: Meeting notes content")
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_WithBothFlags(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", "-t", "Project Planning", "--tags", "project,planning", "Content for project planning"})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: Content for project planning")
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_NoArguments(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create"})
	err := rootCmd.Execute()

	// Should fail because MinimumNArgs(1) is set
	assert.Error(t, err)
	assert.Contains(t, output.String(), "Error:") // Cobra should output an error message
}

func TestCreateCommand_EmptyContent(t *testing.T) {
	rootCmd, output := setupTestCommand()

	rootCmd.SetArgs([]string{"create", ""})
	err := rootCmd.Execute()

	// Should still work but with empty content
	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: ")
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_SpecialCharacters(t *testing.T) {
	rootCmd, output := setupTestCommand()

	content := "Note with special chars: @#$%^&*()_+-={}[]|\\:;\"'<>?,./"
	rootCmd.SetArgs([]string{"create", content})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: "+content)
	assert.Contains(t, output.String(), "Note created successfully!")
}

func TestCreateCommand_LongContent(t *testing.T) {
	rootCmd, output := setupTestCommand()

	// Create a long content string
	longContent := strings.Repeat("This is a very long note content. ", 100)
	rootCmd.SetArgs([]string{"create", longContent})
	err := rootCmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, output.String(), "Creating note: "+longContent)
	assert.Contains(t, output.String(), "Note created successfully!")
}

