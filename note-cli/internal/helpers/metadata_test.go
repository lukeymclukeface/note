package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetFileHash(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	content := "Hello, World!"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Calculate hash
	hash1, err := GetFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}
	
	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}
	
	// Calculate hash again - should be the same
	hash2, err := GetFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to calculate hash second time: %v", err)
	}
	
	if hash1 != hash2 {
		t.Errorf("Expected same hash, got %s and %s", hash1, hash2)
	}
	
	// Change file content
	err = os.WriteFile(testFile, []byte("Different content"), 0644)
	if err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}
	
	// Hash should be different
	hash3, err := GetFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to calculate hash after change: %v", err)
	}
	
	if hash1 == hash3 {
		t.Error("Expected different hash after content change")
	}
}

func TestGetFileHashNonExistent(t *testing.T) {
	_, err := GetFileHash("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		content  []byte
		expected string
	}{
		{"test.txt", []byte("Hello"), "text/plain; charset=utf-8"},
		{"test.json", []byte(`{"key": "value"}`), "application/json"},
		{"test.md", []byte("# Markdown"), "text/plain; charset=utf-8"},
	}
	
	tmpDir := t.TempDir()
	
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(testFile, test.content, 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			
			mimeType, err := GetMimeType(testFile)
			if err != nil {
				t.Fatalf("Failed to get MIME type: %v", err)
			}
			
			if mimeType != test.expected {
				t.Errorf("Expected MIME type %s, got %s", test.expected, mimeType)
			}
		})
	}
}

func TestGetFileType(t *testing.T) {
	tests := []struct {
		mimeType string
		expected string
	}{
		{"audio/mpeg", "audio"},
		{"audio/wav", "audio"},
		{"video/mp4", "video"},
		{"video/avi", "video"},
		{"text/plain", "text"},
		{"application/json", "text"},
		{"application/xml", "text"},
		{"application/yaml", "text"},
		{"image/jpeg", "image"},
		{"image/png", "image"},
		{"application/pdf", "document"},
		{"application/octet-stream", "other"},
		{"unknown/type", "other"},
	}
	
	for _, test := range tests {
		t.Run(test.mimeType, func(t *testing.T) {
			result := GetFileType(test.mimeType)
			if result != test.expected {
				t.Errorf("Expected file type %s for MIME type %s, got %s", test.expected, test.mimeType, result)
			}
		})
	}
}

func TestExtractFileMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	content := "Line 1\nLine 2\nLine 3\nHello world test"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Get file info for comparison
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}
	
	metadata, err := ExtractFileMetadata(testFile)
	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}
	
	// Test basic metadata
	if metadata.FileName != "test.txt" {
		t.Errorf("Expected filename test.txt, got %s", metadata.FileName)
	}
	
	if metadata.FileExt != ".txt" {
		t.Errorf("Expected file extension .txt, got %s", metadata.FileExt)
	}
	
	// Check that last modified time is close to file info time
	timeDiff := metadata.LastModified.Sub(fileInfo.ModTime())
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	if timeDiff > time.Second {
		t.Errorf("Last modified time differs too much from file info")
	}
	
	// Test text-specific metadata
	if metadata.LineCount == nil {
		t.Error("Expected line count to be set for text file")
	} else if *metadata.LineCount != 4 {
		t.Errorf("Expected 4 lines, got %d", *metadata.LineCount)
	}
	
	if metadata.WordCount == nil {
		t.Error("Expected word count to be set for text file")
	} else if *metadata.WordCount != 9 { // "Line", "1", "Line", "2", "Line", "3", "Hello", "world", "test"
		t.Errorf("Expected 9 words, got %d", *metadata.WordCount)
	}
	
	if metadata.CharCount == nil {
		t.Error("Expected char count to be set for text file")
	} else if *metadata.CharCount != len(content) {
		t.Errorf("Expected %d characters, got %d", len(content), *metadata.CharCount)
	}
	
	if metadata.Encoding == nil {
		t.Error("Expected encoding to be set for text file")
	} else if *metadata.Encoding != "UTF-8" {
		t.Errorf("Expected UTF-8 encoding, got %s", *metadata.Encoding)
	}
}

func TestExtractFileMetadataAudio(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")
	
	// Create a dummy MP3 file (just for testing file extension logic)
	content := "fake mp3 content"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	metadata, err := ExtractFileMetadata(testFile)
	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}
	
	// Test basic metadata
	if metadata.FileName != "test.mp3" {
		t.Errorf("Expected filename test.mp3, got %s", metadata.FileName)
	}
	
	if metadata.FileExt != ".mp3" {
		t.Errorf("Expected file extension .mp3, got %s", metadata.FileExt)
	}
	
	// Test audio-specific metadata (basic defaults)
	if metadata.SampleRate == nil {
		t.Error("Expected sample rate to be set for audio file")
	} else if *metadata.SampleRate != 44100 {
		t.Errorf("Expected sample rate 44100, got %d", *metadata.SampleRate)
	}
	
	if metadata.Channels == nil {
		t.Error("Expected channels to be set for audio file")
	} else if *metadata.Channels != 2 {
		t.Errorf("Expected 2 channels, got %d", *metadata.Channels)
	}
}

func TestSerializeDeserializeMetadata(t *testing.T) {
	// Create test metadata
	sampleRate := 44100
	channels := 2
	lineCount := 10
	encoding := "UTF-8"
	
	metadata := &FileMetadata{
		FileName:     "test.mp3",
		FileExt:      ".mp3",
		LastModified: time.Now(),
		SampleRate:   &sampleRate,
		Channels:     &channels,
		LineCount:    &lineCount,
		Encoding:     &encoding,
		Custom:       map[string]interface{}{"custom_field": "custom_value"},
	}
	
	// Serialize
	jsonStr, err := SerializeMetadata(metadata)
	if err != nil {
		t.Fatalf("Failed to serialize metadata: %v", err)
	}
	
	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}
	
	// Deserialize
	deserialized, err := DeserializeMetadata(jsonStr)
	if err != nil {
		t.Fatalf("Failed to deserialize metadata: %v", err)
	}
	
	// Verify deserialized data
	if deserialized.FileName != metadata.FileName {
		t.Errorf("Expected filename %s, got %s", metadata.FileName, deserialized.FileName)
	}
	
	if deserialized.FileExt != metadata.FileExt {
		t.Errorf("Expected file extension %s, got %s", metadata.FileExt, deserialized.FileExt)
	}
	
	if deserialized.SampleRate == nil || *deserialized.SampleRate != *metadata.SampleRate {
		t.Errorf("Sample rate mismatch")
	}
	
	if deserialized.Channels == nil || *deserialized.Channels != *metadata.Channels {
		t.Errorf("Channels mismatch")
	}
	
	if deserialized.LineCount == nil || *deserialized.LineCount != *metadata.LineCount {
		t.Errorf("Line count mismatch")
	}
	
	if deserialized.Encoding == nil || *deserialized.Encoding != *metadata.Encoding {
		t.Errorf("Encoding mismatch")
	}
	
	// Check custom field
	if deserialized.Custom == nil {
		t.Error("Expected custom field to be preserved")
	} else if customValue, ok := deserialized.Custom["custom_field"]; !ok || customValue != "custom_value" {
		t.Error("Custom field not preserved correctly")
	}
}

func TestDeserializeMetadataInvalid(t *testing.T) {
	_, err := DeserializeMetadata("invalid json")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestProcessFileForDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	content := "Hello, World!\nThis is a test file."
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Get expected file size
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}
	expectedSize := fileInfo.Size()
	
	// Process file for database
	fileHash, mimeType, fileType, fileSize, metadata, err := ProcessFileForDatabase(testFile)
	if err != nil {
		t.Fatalf("Failed to process file for database: %v", err)
	}
	
	// Verify hash
	if fileHash == "" {
		t.Error("Expected non-empty file hash")
	}
	
	// Verify MIME type
	if !strings.HasPrefix(mimeType, "text/") {
		t.Errorf("Expected text MIME type, got %s", mimeType)
	}
	
	// Verify file type
	if fileType != "text" {
		t.Errorf("Expected file type 'text', got %s", fileType)
	}
	
	// Verify file size
	if fileSize != expectedSize {
		t.Errorf("Expected file size %d, got %d", expectedSize, fileSize)
	}
	
	// Verify metadata is valid JSON
	if metadata == "" {
		t.Error("Expected non-empty metadata")
	}
	
	// Try to deserialize metadata
	_, err = DeserializeMetadata(metadata)
	if err != nil {
		t.Errorf("Failed to deserialize generated metadata: %v", err)
	}
}

func TestProcessFileForDatabaseNonExistent(t *testing.T) {
	_, _, _, _, _, err := ProcessFileForDatabase("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestExtractTextMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	
	// Test different text patterns
	tests := []struct {
		name         string
		content      string
		expectedLines int
		expectedWords int
		expectedChars int
	}{
		{
			name:         "simple",
			content:      "Hello world",
			expectedLines: 1,
			expectedWords: 2,
			expectedChars: 11,
		},
		{
			name:         "multiline",
			content:      "Line 1\nLine 2\nLine 3",
			expectedLines: 3,
			expectedWords: 6,
			expectedChars: 20,
		},
		{
			name:         "empty_lines",
			content:      "Line 1\n\nLine 3",
			expectedLines: 3,
			expectedWords: 4,
			expectedChars: 14,
		},
		{
			name:         "trailing_newline",
			content:      "Line 1\nLine 2\n",
			expectedLines: 3, // Split creates empty string after final newline
			expectedWords: 4,
			expectedChars: 14, // "Line 1\nLine 2\n" is 14 characters
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := os.WriteFile(testFile, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			
			metadata := &FileMetadata{}
			err = extractTextMetadata(testFile, metadata)
			if err != nil {
				t.Fatalf("Failed to extract text metadata: %v", err)
			}
			
			if metadata.LineCount == nil || *metadata.LineCount != test.expectedLines {
				t.Errorf("Expected %d lines, got %v", test.expectedLines, metadata.LineCount)
			}
			
			if metadata.WordCount == nil || *metadata.WordCount != test.expectedWords {
				t.Errorf("Expected %d words, got %v", test.expectedWords, metadata.WordCount)
			}
			
			if metadata.CharCount == nil {
				t.Error("Expected char count to be set")
			} else if *metadata.CharCount != test.expectedChars {
				t.Errorf("Expected %d characters, got %d", test.expectedChars, *metadata.CharCount)
			}
			
			if metadata.Encoding == nil || *metadata.Encoding != "UTF-8" {
				t.Errorf("Expected UTF-8 encoding, got %v", metadata.Encoding)
			}
		})
	}
}

func TestExtractAudioVideoMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		filename string
		expectSampleRate bool
		expectChannels bool
	}{
		{"test.mp3", true, true},
		{"test.wav", true, true},
		{"test.flac", true, true},
		{"test.m4a", true, true},
		{"test.aac", true, true},
		{"test.mp4", false, false},
		{"test.avi", false, false},
		{"test.mov", false, false},
		{"test.mkv", false, false},
		{"test.unknown", false, false},
	}
	
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, test.filename)
			err := os.WriteFile(testFile, []byte("fake content"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			
			metadata := &FileMetadata{}
			err = extractAudioVideoMetadata(testFile, metadata)
			if err != nil {
				t.Fatalf("Failed to extract audio/video metadata: %v", err)
			}
			
			if test.expectSampleRate {
				if metadata.SampleRate == nil {
					t.Error("Expected sample rate to be set")
				} else if *metadata.SampleRate != 44100 {
					t.Errorf("Expected sample rate 44100, got %d", *metadata.SampleRate)
				}
			} else {
				if metadata.SampleRate != nil {
					t.Error("Expected sample rate to be nil")
				}
			}
			
			if test.expectChannels {
				if metadata.Channels == nil {
					t.Error("Expected channels to be set")
				} else if *metadata.Channels != 2 {
					t.Errorf("Expected 2 channels, got %d", *metadata.Channels)
				}
			} else {
				if metadata.Channels != nil {
					t.Error("Expected channels to be nil")
				}
			}
		})
	}
}
