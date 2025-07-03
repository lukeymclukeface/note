package service

import (
	"context"
	"testing"
	"time"
)

// Test the placeholder transcriber implementation
func TestPlaceholderTranscriber_TranscribeAudio(t *testing.T) {
	transcriber := &PlaceholderTranscriber{}

	tests := []struct {
		name      string
		audioData []byte
		wantErr   bool
		wantText  string
	}{
		{
			name:      "valid audio data",
			audioData: []byte("fake audio data"),
			wantErr:   false,
			wantText:  "This is a placeholder transcription result from dummy model",
		},
		{
			name:      "empty audio data",
			audioData: []byte{},
			wantErr:   true,
		},
		{
			name:      "nil audio data",
			audioData: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := transcriber.TranscribeAudio(ctx, tt.audioData)

			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantText {
				t.Errorf("TranscribeAudio() = %v, want %v", got, tt.wantText)
			}
		})
	}
}

func TestPlaceholderTranscriber_TranscribeStream(t *testing.T) {
	transcriber := &PlaceholderTranscriber{}

	tests := []struct {
		name       string
		audioChunk []byte
		wantErr    bool
	}{
		{
			name:       "valid audio chunk",
			audioChunk: []byte("fake audio chunk"),
			wantErr:    false,
		},
		{
			name:       "empty audio chunk",
			audioChunk: []byte{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resultChan, err := transcriber.TranscribeStream(ctx, tt.audioChunk)

			if (err != nil) != tt.wantErr {
				t.Errorf("TranscribeStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Collect results from channel
			var results []TranscribeStreamResult
			for result := range resultChan {
				results = append(results, result)
			}

			if len(results) != 2 {
				t.Errorf("Expected 2 results, got %d", len(results))
				return
			}

			if results[0].Type != "partial" {
				t.Errorf("Expected first result type 'partial', got %s", results[0].Type)
			}

			if results[1].Type != "final" {
				t.Errorf("Expected second result type 'final', got %s", results[1].Type)
			}
		})
	}
}

func TestTranscribeService_TranscribeAudio(t *testing.T) {
	t.Run("with mock transcriber", func(t *testing.T) {
		mockTranscriber := &MockTranscriber{
			TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
				return "mocked transcription", nil
			},
		}

		service := NewTranscribeServiceWithTranscriber(mockTranscriber)
		ctx := context.Background()

		// Mock transcribers should skip ffmpeg conversion and work correctly
		result, err := service.TranscribeAudio(ctx, []byte("fake audio"))
		
		// We expect this to succeed with the mock
		if err != nil {
			t.Errorf("Expected success with mock transcriber, but got error: %v", err)
		}
		
		if result != "mocked transcription" {
			t.Errorf("Expected 'mocked transcription', got '%s'", result)
		}
	})

	t.Run("empty audio data", func(t *testing.T) {
		service := NewTranscribeService()
		ctx := context.Background()

		_, err := service.TranscribeAudio(ctx, []byte{})
		if err == nil {
			t.Error("Expected error for empty audio data")
		}
	})
}

func TestTranscribeService_TranscribeStream(t *testing.T) {
	mockTranscriber := &MockTranscriber{
		TranscribeStreamFunc: func(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error) {
			resultChan := make(chan TranscribeStreamResult, 1)
			go func() {
				defer close(resultChan)
				resultChan <- TranscribeStreamResult{Type: "final", Text: "mocked stream result"}
			}()
			return resultChan, nil
		},
	}

	service := NewTranscribeServiceWithTranscriber(mockTranscriber)
	ctx := context.Background()

	resultChan, err := service.TranscribeStream(ctx, []byte("fake audio chunk"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var results []TranscribeStreamResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
		return
	}

	if results[0].Text != "mocked stream result" {
		t.Errorf("Expected 'mocked stream result', got %s", results[0].Text)
	}
}

// Test context cancellation for streaming
func TestTranscribeService_TranscribeStream_Cancellation(t *testing.T) {
	transcriber := &PlaceholderTranscriber{}
	ctx, cancel := context.WithCancel(context.Background())

	resultChan, err := transcriber.TranscribeStream(ctx, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	// Cancel context immediately
	cancel()

	// Should not hang and should close channel
	timeout := time.After(1 * time.Second)
	select {
	case <-resultChan:
		// Channel should close
	case <-timeout:
		t.Error("Channel did not close after context cancellation")
	}
}

// Mock Transcriber for testing (same as in handlers_test.go)
type MockTranscriber struct {
	TranscribeAudioFunc  func(ctx context.Context, audioData []byte) (string, error)
	TranscribeStreamFunc func(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error)
}

func (m *MockTranscriber) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	if m.TranscribeAudioFunc != nil {
		return m.TranscribeAudioFunc(ctx, audioData)
	}
	return "mock transcription result", nil
}

func (m *MockTranscriber) TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan TranscribeStreamResult, error) {
	if m.TranscribeStreamFunc != nil {
		return m.TranscribeStreamFunc(ctx, audioChunk)
	}
	// Default mock implementation
	resultChan := make(chan TranscribeStreamResult, 2)
	go func() {
		defer close(resultChan)
		resultChan <- TranscribeStreamResult{Type: "partial", Text: "mock partial"}
		resultChan <- TranscribeStreamResult{Type: "final", Text: "mock final result"}
	}()
	return resultChan, nil
}

// Benchmarks
func BenchmarkPlaceholderTranscriber_TranscribeAudio(b *testing.B) {
	transcriber := &PlaceholderTranscriber{}
	ctx := context.Background()
	audioData := []byte("benchmark audio data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transcriber.TranscribeAudio(ctx, audioData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTranscribeService_WithMock(b *testing.B) {
	mockTranscriber := &MockTranscriber{
		TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
			return "benchmark result", nil
		},
	}

	service := NewTranscribeServiceWithTranscriber(mockTranscriber)
	ctx := context.Background()
	audioData := []byte("benchmark audio data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will still try to call convertToWav, but the mock will be called if conversion succeeds
		service.TranscribeAudio(ctx, audioData)
	}
}
