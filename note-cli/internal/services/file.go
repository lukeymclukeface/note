package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"note-cli/internal/constants"
)

// FileService handles file system operations
type FileService struct{}

// NewFileService creates a new file service instance
func NewFileService() *FileService {
	return &FileService{}
}

// CopyFile copies a file from src to dst
func (s *FileService) CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return err
	}

	return nil
}

// CreateNoteDirectory creates a directory for a note with a timestamp
func (s *FileService) CreateNoteDirectory(baseName string) (string, error) {
	notesDir, err := constants.GetNotesDir()
	if err != nil {
		return "", fmt.Errorf("failed to get notes directory: %w", err)
	}

	folderName := fmt.Sprintf("%s_%s", baseName, time.Now().Format("20060102_150405"))
	destinationDir := filepath.Join(notesDir, folderName)

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	return destinationDir, nil
}

// SaveMarkdownFiles saves transcript and summary as markdown files
func (s *FileService) SaveMarkdownFiles(transcript, summary, transcriptionPath, summaryPath string) error {
	// Write transcription file
	transcriptionContent := fmt.Sprintf("# Transcription\n\n%s\n", transcript)
	if err := os.WriteFile(transcriptionPath, []byte(transcriptionContent), 0644); err != nil {
		return fmt.Errorf("failed to write transcription file: %w", err)
	}

	// Write summary file
	summaryContent := fmt.Sprintf("# Summary\n\n%s\n", summary)
	if err := os.WriteFile(summaryPath, []byte(summaryContent), 0644); err != nil {
		return fmt.Errorf("failed to write summary file: %w", err)
	}

	return nil
}

// GetCompatibleFiles finds files with supported extensions in a directory
func (s *FileService) GetCompatibleFiles(directory string) ([]string, error) {
	var files []string
	extensions := []string{"*.mp3", "*.wav", "*.m4a", "*.ogg", "*.flac", "*.md", "*.txt"}

	if directory == "" {
		directory = "."
	}

	for _, ext := range extensions {
		pattern := filepath.Join(directory, ext)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		files = append(files, matches...)
	}

	return files, nil
}

// ExtractFolderFromContent extracts the folder name from note content
func (s *FileService) ExtractFolderFromContent(content string) string {
	// Look for "Folder: folderName" in the content
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Folder: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Folder: "))
		}
	}
	return ""
}

// EnsureDirectoryExists creates a directory if it doesn't exist
func (s *FileService) EnsureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}
