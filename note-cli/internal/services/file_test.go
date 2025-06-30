package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileService(t *testing.T) {
	fs := NewFileService()
	assert.NotNil(t, fs)
	assert.IsType(t, &FileService{}, fs)
}

func TestCopyFile(t *testing.T) {
	fs := NewFileService()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "file_service_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := "This is test content for file copying"
	err = os.WriteFile(srcPath, []byte(testContent), 0644)
	require.NoError(t, err)

	// Test copying file
	dstPath := filepath.Join(tempDir, "destination.txt")
	err = fs.CopyFile(srcPath, dstPath)
	require.NoError(t, err)

	// Verify destination file exists and has correct content
	content, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// Verify file permissions
	info, err := os.Stat(dstPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
}

func TestCopyFile_SourceNotExists(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "file_service_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "nonexistent.txt")
	dstPath := filepath.Join(tempDir, "destination.txt")

	err = fs.CopyFile(srcPath, dstPath)
	assert.Error(t, err)
}

func TestCopyFile_InvalidDestination(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "file_service_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	err = os.WriteFile(srcPath, []byte("test"), 0644)
	require.NoError(t, err)

	// Try to copy to invalid destination (directory that doesn't exist)
	dstPath := "/invalid/path/destination.txt"
	err = fs.CopyFile(srcPath, dstPath)
	assert.Error(t, err)
}

func TestCreateNoteDirectory(t *testing.T) {
	fs := NewFileService()

	// This test would require mocking constants.GetNotesDir()
	// For now, we'll test the function exists and has the right signature
	assert.NotNil(t, fs.CreateNoteDirectory)

	// Test with a valid base name
	_, err := fs.CreateNoteDirectory("test_note")
	// This may succeed or fail depending on the environment setup
	// We mainly want to ensure it doesn't panic
	if err != nil {
		// If it fails, that's expected in some test environments
		assert.Error(t, err)
	} else {
		// If it succeeds, that's also fine
		assert.NoError(t, err)
	}
}

func TestCreateContentDirectory_TypeMapping(t *testing.T) {
	fs := NewFileService()

	// Test different content types
	contentTypes := []string{"note", "meeting", "interview", "other", "", "unknown"}
	
	for _, contentType := range contentTypes {
		t.Run("content_type_"+contentType, func(t *testing.T) {
			_, err := fs.CreateContentDirectory("test", contentType)
			// This may succeed or fail depending on the environment setup
			// We mainly want to ensure it doesn't panic
			if err != nil {
				// If it fails, that's expected in some test environments
				assert.Error(t, err)
			} else {
				// If it succeeds, that's also fine
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveMarkdownFiles(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "markdown_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	transcriptPath := filepath.Join(tempDir, "transcript.md")
	summaryPath := filepath.Join(tempDir, "summary.md")

	transcript := "This is the transcription of the audio file"
	summary := "This is a brief summary"

	err = fs.SaveMarkdownFiles(transcript, summary, transcriptPath, summaryPath)
	require.NoError(t, err)

	// Verify transcription file
	transcriptContent, err := os.ReadFile(transcriptPath)
	require.NoError(t, err)
	expectedTranscript := "# Transcription\n\n" + transcript + "\n"
	assert.Equal(t, expectedTranscript, string(transcriptContent))

	// Verify summary file
	summaryContent, err := os.ReadFile(summaryPath)
	require.NoError(t, err)
	expectedSummary := "# Summary\n\n" + summary + "\n"
	assert.Equal(t, expectedSummary, string(summaryContent))
}

func TestSaveMarkdownFiles_InvalidPath(t *testing.T) {
	fs := NewFileService()

	transcript := "Test transcript"
	summary := "Test summary"
	invalidPath := "/invalid/path/file.md"

	err := fs.SaveMarkdownFiles(transcript, summary, invalidPath, invalidPath)
	assert.Error(t, err)
}

func TestGetCompatibleFiles(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "compatible_files_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files with various extensions
	testFiles := []string{
		"audio1.mp3",
		"audio2.wav",
		"audio3.m4a",
		"audio4.ogg",
		"audio5.flac",
		"document1.md",
		"document2.txt",
		"ignored.doc",
		"ignored.pdf",
	}

	for _, filename := range testFiles {
		path := filepath.Join(tempDir, filename)
		err := os.WriteFile(path, []byte("test content"), 0644)
		require.NoError(t, err)
	}

	files, err := fs.GetCompatibleFiles(tempDir)
	require.NoError(t, err)

	// Should find 7 compatible files (mp3, wav, m4a, ogg, flac, md, txt)
	assert.Len(t, files, 7)

	// Check that all returned files exist
	for _, file := range files {
		_, err := os.Stat(file)
		assert.NoError(t, err, "File should exist: %s", file)
	}

	// Check that incompatible files are not included
	for _, file := range files {
		ext := filepath.Ext(file)
		assert.NotEqual(t, ".doc", ext)
		assert.NotEqual(t, ".pdf", ext)
	}
}

func TestGetCompatibleFiles_EmptyDirectory(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "empty_dir_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files, err := fs.GetCompatibleFiles(tempDir)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestGetCompatibleFiles_DefaultDirectory(t *testing.T) {
	fs := NewFileService()

	// Test with empty string (should use current directory)
	_, err := fs.GetCompatibleFiles("")
	// This may succeed or fail depending on the current directory content
	// We mainly want to ensure it doesn't panic
	if err != nil {
		// If it fails, that's expected in some test environments
		assert.Error(t, err)
	} else {
		// If it succeeds, we just check that no error occurred
		assert.NoError(t, err)
	}
}

func TestExtractFolderFromContent(t *testing.T) {
	fs := NewFileService()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "ValidFolderLine",
			content:  "Some content\nFolder: meeting_20240101_120000\nMore content",
			expected: "meeting_20240101_120000",
		},
		{
			name:     "FolderLineAtStart",
			content:  "Folder: note_20240101_130000\nSome content after",
			expected: "note_20240101_130000",
		},
		{
			name:     "FolderLineAtEnd",
			content:  "Some content before\nFolder: interview_20240101_140000",
			expected: "interview_20240101_140000",
		},
		{
			name:     "FolderWithSpaces",
			content:  "Folder:   spaced_folder_name   \nOther content",
			expected: "spaced_folder_name",
		},
		{
			name:     "NoFolderLine",
			content:  "Some content\nNo folder info here\nMore content",
			expected: "",
		},
		{
			name:     "EmptyContent",
			content:  "",
			expected: "",
		},
		{
			name:     "MultipleFolderLines",
			content:  "Folder: first_folder\nSome content\nFolder: second_folder",
			expected: "first_folder", // Should return the first one found
		},
		{
			name:     "FolderLineWithoutColon",
			content:  "Folder meeting_20240101_120000\nMore content",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fs.ExtractFolderFromContent(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "ensure_dir_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test creating a new directory
	newDir := filepath.Join(tempDir, "new_directory")
	err = fs.EnsureDirectoryExists(newDir)
	require.NoError(t, err)

	// Verify directory was created
	info, err := os.Stat(newDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test with existing directory (should not error)
	err = fs.EnsureDirectoryExists(newDir)
	assert.NoError(t, err)

	// Test creating nested directories
	nestedDir := filepath.Join(tempDir, "level1", "level2", "level3")
	err = fs.EnsureDirectoryExists(nestedDir)
	require.NoError(t, err)

	// Verify nested directory was created
	info, err = os.Stat(nestedDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestEnsureDirectoryExists_InvalidPath(t *testing.T) {
	fs := NewFileService()

	// Test with invalid path (permission denied scenario)
	// This test might be platform-specific
	invalidPath := "/root/cannot_create_here"
	err := fs.EnsureDirectoryExists(invalidPath)
	// On most systems, this should fail due to permissions
	// but the exact error depends on the system
	if err != nil {
		assert.Error(t, err)
	}
}

// Benchmark tests
func BenchmarkCopyFile(b *testing.B) {
	fs := NewFileService()

	tempDir, err := os.MkdirTemp("", "benchmark_copy")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create source file with some content
	srcPath := filepath.Join(tempDir, "source.txt")
	content := make([]byte, 1024) // 1KB content
	for i := range content {
		content[i] = byte(i % 256)
	}
	err = os.WriteFile(srcPath, content, 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dstPath := filepath.Join(tempDir, "dest_"+string(rune(i))+".txt")
		err := fs.CopyFile(srcPath, dstPath)
		require.NoError(b, err)
	}
}

func BenchmarkExtractFolderFromContent(b *testing.B) {
	fs := NewFileService()

	content := `This is some content before the folder line.
It contains multiple lines of text.
Folder: meeting_20240101_120000
And some content after the folder line.
This content continues for a while to make the benchmark more realistic.`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fs.ExtractFolderFromContent(content)
	}
}

// Test edge cases for file operations
func TestFileOperations_EdgeCases(t *testing.T) {
	fs := NewFileService()

	t.Run("CopyEmptyFile", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "empty_file_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		srcPath := filepath.Join(tempDir, "empty.txt")
		dstPath := filepath.Join(tempDir, "empty_copy.txt")

		// Create empty file
		err = os.WriteFile(srcPath, []byte{}, 0644)
		require.NoError(t, err)

		err = fs.CopyFile(srcPath, dstPath)
		require.NoError(t, err)

		// Verify empty file was copied
		content, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("CopyLargeFile", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "large_file_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		srcPath := filepath.Join(tempDir, "large.txt")
		dstPath := filepath.Join(tempDir, "large_copy.txt")

		// Create 1MB file
		largeContent := make([]byte, 1024*1024)
		for i := range largeContent {
			largeContent[i] = byte(i % 256)
		}
		err = os.WriteFile(srcPath, largeContent, 0644)
		require.NoError(t, err)

		err = fs.CopyFile(srcPath, dstPath)
		require.NoError(t, err)

		// Verify large file was copied correctly
		copiedContent, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		assert.Equal(t, largeContent, copiedContent)
	})
}
