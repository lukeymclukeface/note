package cmd

import (
	"bytes"
	"database/sql"
	"note-cli/internal/database"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockListNotes simulates the database.ListNotes function for testing
func mockListNotes(db *sql.DB, tag string) ([]database.Note, error) {
	if tag == "nonexistent" {
		return []database.Note{}, nil
	}
	
	notes := []database.Note{
		{
			ID:        1,
			Title:     "Test Note 1",
			Content:   "Content of test note 1",
			Tags:      "test,sample",
			CreatedAt: "2024-01-01T10:00:00Z",
		},
		{
			ID:        2,
			Title:     "Test Note 2",
			Content:   "Content of test note 2",
			Tags:      "test,work",
			CreatedAt: "2024-01-02T11:00:00Z",
		},
	}
	
	// Filter by tag if specified
	if tag != "" {
		var filteredNotes []database.Note
		for _, note := range notes {
			if containsTag(note.Tags, tag) {
				filteredNotes = append(filteredNotes, note)
			}
		}
		return filteredNotes, nil
	}
	
	return notes, nil
}

// containsTag checks if a tag exists in the comma-separated tags string
func containsTag(tags, tag string) bool {
	return tags == tag || 
		   len(tags) > len(tag) && tags[:len(tag)+1] == tag+"," ||
		   len(tags) > len(tag) && tags[len(tags)-len(tag)-1:] == ","+tag ||
		   len(tags) > len(tag)*2+1 && findInTags(tags, ","+tag+",")
}

func findInTags(tags, pattern string) bool {
	for i := 0; i <= len(tags)-len(pattern); i++ {
		if tags[i:i+len(pattern)] == pattern {
			return true
		}
	}
	return false
}

// setupTestListCommand creates a test root command with the list command attached
func setupTestListCommand() (*cobra.Command, *bytes.Buffer) {
	// Create root command and add list command
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(listCmd)

	output := new(bytes.Buffer)
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	return rootCmd, output
}

func TestListCommand_Basic(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(listCmd)

	output := new(bytes.Buffer)
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	// This test checks that the command configuration is correct
	// We can't easily test the Run function due to interactive prompts
	// and database dependencies, so we test the command setup instead
	assert.NotNil(t, listCmd.Run)
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List all notes", listCmd.Short)
}

func TestListCommand_WithTagFlag(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(listCmd)

	// Test that the tag flag is properly registered
	tagFlag := listCmd.Flags().Lookup("tag")
	require.NotNil(t, tagFlag)
	assert.Equal(t, "", tagFlag.DefValue)
	
	// Test shorthand exists by checking if ShorthandLookup returns a flag
	shorthandFlag := listCmd.Flags().ShorthandLookup("t")
	assert.NotNil(t, shorthandFlag)
	assert.Equal(t, tagFlag, shorthandFlag)
}

func TestListCommand_WithRecentFlag(t *testing.T) {
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.AddCommand(listCmd)

	// Test that the recent flag is properly registered
	recentFlag := listCmd.Flags().Lookup("recent")
	require.NotNil(t, recentFlag)
	assert.Equal(t, "false", recentFlag.DefValue)
	
	// Test shorthand exists by checking if ShorthandLookup returns a flag
	shorthandFlag := listCmd.Flags().ShorthandLookup("r")
	assert.NotNil(t, shorthandFlag)
	assert.Equal(t, recentFlag, shorthandFlag)
}

// Test the mock functions directly
func TestMockListNotes_AllNotes(t *testing.T) {
	notes, err := mockListNotes(nil, "")
	require.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, "Test Note 1", notes[0].Title)
	assert.Equal(t, "Test Note 2", notes[1].Title)
}

func TestMockListNotes_FilterByTag(t *testing.T) {
	notes, err := mockListNotes(nil, "test")
	require.NoError(t, err)
	assert.Len(t, notes, 2) // Both notes have "test" tag
	
	notes, err = mockListNotes(nil, "work")
	require.NoError(t, err)
	assert.Len(t, notes, 1) // Only second note has "work" tag
	assert.Equal(t, "Test Note 2", notes[0].Title)
}

func TestMockListNotes_NoMatchingTag(t *testing.T) {
	notes, err := mockListNotes(nil, "nonexistent")
	require.NoError(t, err)
	assert.Len(t, notes, 0)
}

func TestContainsTag(t *testing.T) {
	tests := []struct {
		tags     string
		tag      string
		expected bool
	}{
		{"test", "test", true},
		{"test,work", "test", true},
		{"test,work", "work", true},
		{"important,test,work", "test", true},
		{"test,work,urgent", "work", true},
		{"testing", "test", false},
		{"", "test", false},
		{"work", "test", false},
		{"important,urgent", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.tags+"_contains_"+tt.tag, func(t *testing.T) {
			result := containsTag(tt.tags, tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test command configuration
func TestListCommand_Configuration(t *testing.T) {
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List all notes", listCmd.Short)
	assert.Contains(t, listCmd.Long, "Display a list of all stored notes")
	assert.NotNil(t, listCmd.Run)
}

// Test that flags are properly defined
func TestListCommand_FlagDefinitions(t *testing.T) {
	// Check tag flag
	tagFlag := listCmd.Flags().Lookup("tag")
	require.NotNil(t, tagFlag, "tag flag should be defined")
	assert.Equal(t, "string", tagFlag.Value.Type())
	assert.Equal(t, "", tagFlag.DefValue)

	// Check recent flag  
	recentFlag := listCmd.Flags().Lookup("recent")
	require.NotNil(t, recentFlag, "recent flag should be defined")
	assert.Equal(t, "bool", recentFlag.Value.Type())
	assert.Equal(t, "false", recentFlag.DefValue)
}
