package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// AudioService handles audio file operations
type AudioService struct{}

// NewAudioService creates a new audio service instance
func NewAudioService() *AudioService {
	return &AudioService{}
}

// ConvertToMP3 converts an audio file to MP3 format
func (s *AudioService) ConvertToMP3(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".mp3" {
		return filePath, nil // Already an MP3
	}

	outputPath := strings.TrimSuffix(filePath, ext) + ".mp3"
	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-acodec", "libmp3lame",
		"-ab", "128k",
		"-ar", "44100", // Standardize sample rate
		"-ac", "1",     // Convert to mono
		"-y",           // Overwrite output file
		outputPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("FFmpeg error output: %s\n", stderr.String())
		return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return outputPath, nil
}

// GetAudioDuration returns the duration of an audio file in seconds
func (s *AudioService) GetAudioDuration(filePath string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}
	return int(duration), nil
}

// SplitAudio splits an audio file into a chunk starting at 'start' seconds with 'duration' seconds
func (s *AudioService) SplitAudio(inputPath, outputPath string, start, duration int) error {
	// Use same audio format and quality as the conversion function for consistency
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ss", strconv.Itoa(start),
		"-t", strconv.Itoa(duration),
		"-acodec", "libmp3lame",
		"-ab", "128k",
		"-ar", "44100",
		"-ac", "1",
		"-y", // Overwrite output file if it exists
		outputPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split audio chunk at %ds: %s", start, stderr.String())
	}

	return nil
}

// IsValidAudioFile checks if a file is a supported audio format
func (s *AudioService) IsValidAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".wav" || ext == ".m4a" || ext == ".ogg" || ext == ".flac"
}

// IsValidTextFile checks if a file is a supported text format
func (s *AudioService) IsValidTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".txt"
}

// ChunkedTranscriptionResult represents the result of chunked transcription
type ChunkedTranscriptionResult struct {
	FullTranscript string
	ChunkFiles     []string // Paths to individual chunk transcript files
}

// TranscribeFileChunked transcribes a large audio file in chunks
func (s *AudioService) TranscribeFileChunked(filePath, destinationDir string, provider AIProvider) (*ChunkedTranscriptionResult, error) {
	// Define the chunk duration (e.g., 10 minutes)
	chunkDuration := 10 * 60 // 10 minutes in seconds

	// Get the duration of the file using ffprobe
	duration, err := s.GetAudioDuration(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio duration: %w", err)
	}

	// If file is shorter than chunk duration, transcribe normally
	if duration <= chunkDuration {
		transcript, err := provider.TranscribeAudioFile(filePath)
		if err != nil {
			return nil, err
		}
		return &ChunkedTranscriptionResult{
			FullTranscript: transcript,
			ChunkFiles:     []string{},
		}, nil
	}

	totalChunks := (duration + chunkDuration - 1) / chunkDuration // ceiling division
	var fullTranscript strings.Builder
	var chunkFiles []string

	fmt.Printf("🔄 Audio file is %d minutes long, using chunked transcription...\n", duration/60)

	for chunk := 1; chunk <= totalChunks; chunk++ {
		start := (chunk - 1) * chunkDuration
		end := start + chunkDuration
		if end > duration {
			end = duration
		}

		startMin := float64(start) / 60.0
		endMin := float64(end) / 60.0

		fmt.Printf("📝 Transcribing chunk %d/%d (%.1f-%.1f minutes)\n", chunk, totalChunks, startMin, endMin)

		chunkFilePath := fmt.Sprintf("%s_chunk_%d.wav", filePath, start)

		// Split the audio file into chunks using ffmpeg
		err := s.SplitAudio(filePath, chunkFilePath, start, chunkDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to split audio: %w", err)
		}

		chunkTranscript, err := provider.TranscribeAudioFile(chunkFilePath)
		os.Remove(chunkFilePath) // Clean up the chunk file

		if err != nil {
			return nil, fmt.Errorf("failed to transcribe chunk %d: %w", chunk, err)
		}

		fullTranscript.WriteString(chunkTranscript)
		fullTranscript.WriteString("\n\n")

		// Save individual chunk transcription if destination directory is provided
		if destinationDir != "" {
			chunkFileName := fmt.Sprintf("transcription_chunk_%02d.md", chunk)
			chunkFilePath := filepath.Join(destinationDir, chunkFileName)

			chunkContent := fmt.Sprintf("# Transcription Chunk %d\n\n**Time Range:** %.1f - %.1f minutes\n\n%s\n",
				chunk, startMin, endMin, chunkTranscript)

			if writeErr := os.WriteFile(chunkFilePath, []byte(chunkContent), 0644); writeErr != nil {
				// Don't fail the entire process if we can't save the chunk, just log it
				fmt.Printf("Warning: Failed to save chunk %d transcription: %v\n", chunk, writeErr)
			} else {
				chunkFiles = append(chunkFiles, chunkFilePath)
			}
		}
	}

	return &ChunkedTranscriptionResult{
		FullTranscript: strings.TrimSpace(fullTranscript.String()),
		ChunkFiles:     chunkFiles,
	}, nil
}

// CheckDependencies verifies that required audio tools are available
func (s *AudioService) CheckDependencies() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found. Please install with 'brew install ffmpeg' or run 'note setup'")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return fmt.Errorf("ffprobe not found. This is typically included with ffmpeg. Please install with 'brew install ffmpeg' or run 'note setup'")
	}
	return nil
}
