package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestFirstNWordsSummarizer_SummarizeText(t *testing.T) {
	summarizer := &FirstNWordsSummarizer{}

	tests := []struct {
		name     string
		text     string
		maxWords int
		wantText string
		wantErr  bool
	}{
		{
			name:     "text shorter than maxWords",
			text:     "Hello world",
			maxWords: 5,
			wantText: "Hello world",
			wantErr:  false,
		},
		{
			name:     "text longer than maxWords",
			text:     "This is a very long text that should be truncated to the first few words only",
			maxWords: 5,
			wantText: "This is a very long...",
			wantErr:  false,
		},
		{
			name:     "empty text",
			text:     "",
			maxWords: 5,
			wantText: "",
			wantErr:  true,
		},
		{
			name:     "zero maxWords defaults to 50",
			text:     "This is a test text with many words that should be truncated at fifty words maximum by default behavior",
			maxWords: 0,
			wantText: "This is a test text with many words that should be truncated at fifty words maximum by default behavior",
			wantErr:  false,
		},
		{
			name:     "negative maxWords defaults to 50",
			text:     "Short text",
			maxWords: -1,
			wantText: "Short text",
			wantErr:  false,
		},
		{
			name:     "text with extra whitespace",
			text:     "  This   has   extra   spaces  ",
			maxWords: 3,
			wantText: "This has extra...",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := summarizer.SummarizeText(ctx, tt.text, tt.maxWords)

			if (err != nil) != tt.wantErr {
				t.Errorf("SummarizeText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantText {
				t.Errorf("SummarizeText() = %v, want %v", got, tt.wantText)
			}
		})
	}
}

func TestSummarizeService_SummarizeText(t *testing.T) {
	t.Run("with default service", func(t *testing.T) {
		service := NewSummarizeService()
		ctx := context.Background()

		text := "This is a test text with exactly ten words in it"
		result, err := service.SummarizeText(ctx, text)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Should return all words since it's less than default 50
		if result != text {
			t.Errorf("Expected %q, got %q", text, result)
		}
	})

	t.Run("with custom mock summarizer", func(t *testing.T) {
		mockSummarizer := &MockSummarizer{
			SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
				return fmt.Sprintf("Mock summary of: %s (max %d words)", text, maxWords), nil
			},
		}

		service := NewSummarizeServiceWithSummarizer(mockSummarizer, 25)
		ctx := context.Background()

		text := "Original text"
		result, err := service.SummarizeText(ctx, text)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "Mock summary of: Original text (max 25 words)"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("empty text returns error", func(t *testing.T) {
		service := NewSummarizeService()
		ctx := context.Background()

		_, err := service.SummarizeText(ctx, "")
		if err == nil {
			t.Error("Expected error for empty text")
		}
	})
}

func TestSummarizeService_SummarizeTextWithOptions(t *testing.T) {
	mockSummarizer := &MockSummarizer{
		SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
			// Return first N words to verify maxWords is passed correctly
			words := strings.Fields(text)
			if len(words) <= maxWords {
				return strings.Join(words, " "), nil
			}
			return strings.Join(words[:maxWords], " ") + "...", nil
		},
	}

	service := NewSummarizeServiceWithSummarizer(mockSummarizer, 50)
	ctx := context.Background()

	tests := []struct {
		name     string
		text     string
		maxWords int
		expected string
	}{
		{
			name:     "custom maxWords less than text",
			text:     "One two three four five six seven eight",
			maxWords: 3,
			expected: "One two three...",
		},
		{
			name:     "custom maxWords more than text",
			text:     "One two three",
			maxWords: 5,
			expected: "One two three",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.SummarizeTextWithOptions(ctx, tt.text, tt.maxWords)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSummarizeService_ErrorHandling(t *testing.T) {
	mockSummarizer := &MockSummarizer{
		SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
			return "", fmt.Errorf("summarization failed")
		},
	}

	service := NewSummarizeServiceWithSummarizer(mockSummarizer, 50)
	ctx := context.Background()

	_, err := service.SummarizeText(ctx, "test text")
	if err == nil {
		t.Error("Expected error from mock summarizer")
	}

	expectedError := "summarization failed"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

// Test with context cancellation
func TestSummarizeService_ContextCancellation(t *testing.T) {
	mockSummarizer := &MockSummarizer{
		SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				return "summary", nil
			}
		},
	}

	service := NewSummarizeServiceWithSummarizer(mockSummarizer, 50)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := service.SummarizeText(ctx, "test text")
	if err == nil {
		t.Error("Expected error due to cancelled context")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// Mock Summarizer for testing
type MockSummarizer struct {
	SummarizeTextFunc func(ctx context.Context, text string, maxWords int) (string, error)
}

func (m *MockSummarizer) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
	if m.SummarizeTextFunc != nil {
		return m.SummarizeTextFunc(ctx, text, maxWords)
	}
	return "mock summary", nil
}

// Edge case tests
func TestSummarizeService_EdgeCases(t *testing.T) {
	service := NewSummarizeService()
	ctx := context.Background()

	t.Run("text with only whitespace", func(t *testing.T) {
		_, err := service.SummarizeText(ctx, "   \t\n  ")
		if err == nil {
			t.Error("Expected error for whitespace-only text")
		}
	})

	t.Run("very long text", func(t *testing.T) {
		// Create text with 100 words
		words := make([]string, 100)
		for i := range words {
			words[i] = fmt.Sprintf("word%d", i)
		}
		longText := strings.Join(words, " ")

		result, err := service.SummarizeText(ctx, longText)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Should be truncated to first 50 words + "..."
		resultWords := strings.Fields(strings.TrimSuffix(result, "..."))
		if len(resultWords) != 50 {
			t.Errorf("Expected 50 words in summary, got %d", len(resultWords))
		}

		if !strings.HasSuffix(result, "...") {
			t.Error("Expected summary to end with '...'")
		}
	})

	t.Run("text with special characters", func(t *testing.T) {
		text := "Hello, world! This is a test with special characters: @#$%^&*()"
		result, err := service.SummarizeText(ctx, text)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result != text {
			t.Errorf("Expected %q, got %q", text, result)
		}
	})
}

// Benchmarks
func BenchmarkFirstNWordsSummarizer_SummarizeText(b *testing.B) {
	summarizer := &FirstNWordsSummarizer{}
	ctx := context.Background()
	text := "This is a benchmark test with multiple words that will be used to test the performance of the summarization function"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := summarizer.SummarizeText(ctx, text, 10)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSummarizeService_WithMock(b *testing.B) {
	mockSummarizer := &MockSummarizer{
		SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
			return "benchmark summary", nil
		},
	}

	service := NewSummarizeServiceWithSummarizer(mockSummarizer, 50)
	ctx := context.Background()
	text := "This is a benchmark test with multiple words for performance testing"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.SummarizeText(ctx, text)
		if err != nil {
			b.Fatal(err)
		}
	}
}
