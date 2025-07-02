package helpers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

// FileMetadata represents metadata for different file types
type FileMetadata struct {
	// Common metadata
	FileName     string    `json:"file_name"`
	FileExt      string    `json:"file_extension"`
	LastModified time.Time `json:"last_modified"`
	
	// Audio/Video metadata
	Duration     *float64 `json:"duration,omitempty"`     // Duration in seconds
	SampleRate   *int     `json:"sample_rate,omitempty"`  // Sample rate for audio
	Channels     *int     `json:"channels,omitempty"`     // Number of audio channels
	Bitrate      *int     `json:"bitrate,omitempty"`      // Bitrate in kbps
	Codec        *string  `json:"codec,omitempty"`        // Audio/video codec
	
	// Text metadata
	LineCount    *int     `json:"line_count,omitempty"`   // Number of lines in text files
	WordCount    *int     `json:"word_count,omitempty"`   // Number of words in text files
	CharCount    *int     `json:"char_count,omitempty"`   // Number of characters in text files
	Encoding     *string  `json:"encoding,omitempty"`     // Text encoding
	
	// Image metadata (if needed in future)
	Width        *int     `json:"width,omitempty"`        // Image width
	Height       *int     `json:"height,omitempty"`       // Image height
	
	// Custom metadata
	Custom       map[string]interface{} `json:"custom,omitempty"`
}

// GetFileHash calculates MD5 hash of a file
func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetMimeType determines the MIME type of a file
func GetMimeType(filePath string) (string, error) {
	// First try using file extension
	ext := filepath.Ext(filePath)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType, nil
	}

	// Fallback to content detection
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to detect mime type: %w", err)
	}

	return mtype.String(), nil
}

// GetFileType determines the general file type category
func GetFileType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "text/") || 
		 mimeType == "application/json" ||
		 mimeType == "application/xml" ||
		 mimeType == "application/yaml":
		return "text"
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case mimeType == "application/pdf":
		return "document"
	default:
		return "other"
	}
}

// ExtractFileMetadata extracts metadata from a file based on its type
func ExtractFileMetadata(filePath string) (*FileMetadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	metadata := &FileMetadata{
		FileName:     filepath.Base(filePath),
		FileExt:      filepath.Ext(filePath),
		LastModified: fileInfo.ModTime(),
	}

	mimeType, err := GetMimeType(filePath)
	if err != nil {
		return metadata, nil // Return basic metadata even if mime detection fails
	}

	fileType := GetFileType(mimeType)
	
	switch fileType {
	case "audio", "video":
		if err := extractAudioVideoMetadata(filePath, metadata); err != nil {
			// Log error but don't fail completely
			fmt.Printf("Warning: Failed to extract audio/video metadata: %v\n", err)
		}
	case "text":
		if err := extractTextMetadata(filePath, metadata); err != nil {
			// Log error but don't fail completely
			fmt.Printf("Warning: Failed to extract text metadata: %v\n", err)
		}
	}

	return metadata, nil
}

// extractAudioVideoMetadata extracts metadata for audio and video files
func extractAudioVideoMetadata(filePath string, metadata *FileMetadata) error {
	// For now, we'll use basic file info
	// In a production environment, you might want to use libraries like:
	// - github.com/dhowden/tag for audio metadata
	// - FFmpeg bindings for comprehensive audio/video metadata
	
	// Basic implementation - you can enhance this with proper audio/video libraries
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// Set some defaults based on file extension
	switch ext {
	case ".mp3", ".wav", ".flac", ".m4a", ".aac":
		// Default audio settings - in practice, these should be extracted from the file
		sampleRate := 44100
		channels := 2
		metadata.SampleRate = &sampleRate
		metadata.Channels = &channels
	case ".mp4", ".avi", ".mov", ".mkv":
		// Default video settings
		// These should be extracted using proper video libraries
	}
	
	return nil
}

// extractTextMetadata extracts metadata for text files
func extractTextMetadata(filePath string, metadata *FileMetadata) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open text file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read text file: %w", err)
	}

	text := string(content)
	
	// Count lines
	lines := strings.Split(text, "\n")
	lineCount := len(lines)
	metadata.LineCount = &lineCount
	
	// Count words (simple implementation)
	words := strings.Fields(text)
	wordCount := len(words)
	metadata.WordCount = &wordCount
	
	// Count characters
	charCount := len(text)
	metadata.CharCount = &charCount
	
	// Basic encoding detection (simplified)
	encoding := "UTF-8" // Default assumption
	metadata.Encoding = &encoding
	
	return nil
}

// SerializeMetadata converts FileMetadata to JSON string
func SerializeMetadata(metadata *FileMetadata) (string, error) {
	data, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to serialize metadata: %w", err)
	}
	return string(data), nil
}

// DeserializeMetadata converts JSON string back to FileMetadata
func DeserializeMetadata(jsonData string) (*FileMetadata, error) {
	var metadata FileMetadata
	if err := json.Unmarshal([]byte(jsonData), &metadata); err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
	}
	return &metadata, nil
}

// ProcessFileForDatabase prepares file information for database storage
func ProcessFileForDatabase(filePath string) (fileHash, mimeType, fileType string, fileSize int64, metadata string, err error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", "", 0, "", fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize = fileInfo.Size()

	// Calculate hash
	fileHash, err = GetFileHash(filePath)
	if err != nil {
		return "", "", "", 0, "", fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Get MIME type
	mimeType, err = GetMimeType(filePath)
	if err != nil {
		return "", "", "", 0, "", fmt.Errorf("failed to get MIME type: %w", err)
	}

	// Determine file type
	fileType = GetFileType(mimeType)

	// Extract metadata
	fileMetadata, err := ExtractFileMetadata(filePath)
	if err != nil {
		return "", "", "", 0, "", fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Serialize metadata
	metadata, err = SerializeMetadata(fileMetadata)
	if err != nil {
		return "", "", "", 0, "", fmt.Errorf("failed to serialize metadata: %w", err)
	}

	return fileHash, mimeType, fileType, fileSize, metadata, nil
}
