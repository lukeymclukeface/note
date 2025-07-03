package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TranscribeStreamResult represents a streaming transcription result
type TranscribeStreamResult struct {
	Type string `json:"type"` // "partial" or "final"
	Text string `json:"text"`
}

// Transcriber interface allows swapping different transcription implementations
type Transcriber interface {
	TranscribeAudio(ctx context.Context, audioData []byte) (string, error)
	TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error)
}

// TranscribeService handles audio transcription
type TranscribeService struct {
	transcriber Transcriber
}

// PlaceholderTranscriber is a dummy implementation for development
type PlaceholderTranscriber struct{}

// NewTranscribeService creates a new transcription service
func NewTranscribeService() *TranscribeService {
	return &TranscribeService{
		transcriber: &PlaceholderTranscriber{},
	}
}

// NewTranscribeServiceWithTranscriber creates a service with a custom transcriber
func NewTranscribeServiceWithTranscriber(transcriber Transcriber) *TranscribeService {
	return &TranscribeService{
		transcriber: transcriber,
	}
}

// convertToWav converts audio data to WAV format using ffmpeg
func (s *TranscribeService) convertToWav(ctx context.Context, audioData []byte) ([]byte, error) {
	// Create temporary files for input and output
	tempDir, err := os.MkdirTemp("", "audio_convert_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	inputFile := filepath.Join(tempDir, "input_audio")
	outputFile := filepath.Join(tempDir, "output.wav")

	// Write input audio data to temp file
	if err := os.WriteFile(inputFile, audioData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	// Run ffmpeg to convert to WAV
	cmd := exec.CommandContext(ctx, "ffmpeg", 
		"-i", inputFile,
		"-ar", "16000", // 16kHz sample rate
		"-ac", "1",    // mono
		"-f", "wav",
		"-y", // overwrite output file
		outputFile,
	)

	// Capture stderr for error reporting
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg conversion failed: %w, stderr: %s", err, stderr.String())
	}

	// Read the converted WAV file
	wavData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read converted file: %w", err)
	}

	return wavData, nil
}

// TranscribeAudio transcribes audio data to text
func (s *TranscribeService) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("audio data is empty")
	}

	// Convert to WAV format using ffmpeg
	wavData, err := s.convertToWav(ctx, audioData)
	if err != nil {
		return "", fmt.Errorf("audio conversion failed: %w", err)
	}

	// Use the configured transcriber to process the WAV data
	return s.transcriber.TranscribeAudio(ctx, wavData)
}

// TranscribeAudio implementation for PlaceholderTranscriber
func (p *PlaceholderTranscriber) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("audio data is empty")
	}
	
	// Placeholder implementation - returns dummy string
	return "This is a placeholder transcription result from dummy model", nil
}

// TranscribeStream implementation for PlaceholderTranscriber
func (p *PlaceholderTranscriber) TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error) {
	if len(audioChunk) == 0 {
		return nil, fmt.Errorf("audio chunk is empty")
	}
	
	// Create a channel for streaming results
	resultChan := make(chan TranscribeStreamResult, 10)
	
	// Simulate streaming transcription in a goroutine
	go func() {
		defer close(resultChan)
		
		// Simulate partial result
		select {
		case resultChan <- TranscribeStreamResult{
			Type: "partial",
			Text: "Processing audio chunk (placeholder)...",
		}:
		case <-ctx.Done():
			return
		}
		
		// Simulate processing time
		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return
		}
		
		// Send final result
		select {
		case resultChan <- TranscribeStreamResult{
			Type: "final",
			Text: "This is a placeholder streaming transcription result from dummy model",
		}:
		case <-ctx.Done():
		}
	}()
	
	return resultChan, nil
}

// TranscribeRequest represents a transcription request
type TranscribeRequest struct {
	AudioData []byte `json:"audio_data"`
	Format    string `json:"format"`
	Language  string `json:"language,omitempty"`
}

// TranscribeResponse represents a transcription response
type TranscribeResponse struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Duration   float64 `json:"duration"`
}

// TranscribeStream processes audio chunks and returns transcription results via a channel
func (s *TranscribeService) TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error) {
	return s.transcriber.TranscribeStream(ctx, audioChunk)
}
