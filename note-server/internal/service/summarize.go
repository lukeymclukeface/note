package service

import (
	"context"
	"fmt"
	"strings"
)

// Summarizer interface allows swapping different summarization implementations
type Summarizer interface {
	SummarizeText(ctx context.Context, text string, maxWords int) (string, error)
}

// SummarizeService handles text summarization
type SummarizeService struct {
	summarizer Summarizer
	defaultMaxWords int
}

// FirstNWordsSummarizer is a stub implementation that returns first N words
type FirstNWordsSummarizer struct{}

// NewSummarizeService creates a new summarization service with default stub
func NewSummarizeService() *SummarizeService {
	return &SummarizeService{
		summarizer: &FirstNWordsSummarizer{},
		defaultMaxWords: 50, // Default to first 50 words
	}
}

// NewSummarizeServiceWithSummarizer creates a service with a custom summarizer
func NewSummarizeServiceWithSummarizer(summarizer Summarizer, defaultMaxWords int) *SummarizeService {
	return &SummarizeService{
		summarizer: summarizer,
		defaultMaxWords: defaultMaxWords,
	}
}

// SummarizeText generates a summary of the given text
func (s *SummarizeService) SummarizeText(ctx context.Context, text string) (string, error) {
	return s.SummarizeTextWithOptions(ctx, text, s.defaultMaxWords)
}

// SummarizeTextWithOptions generates a summary with specific options
func (s *SummarizeService) SummarizeTextWithOptions(ctx context.Context, text string, maxWords int) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("text is empty")
	}
	
	return s.summarizer.SummarizeText(ctx, text, maxWords)
}

// SummarizeText implementation for FirstNWordsSummarizer - returns first N words
func (f *FirstNWordsSummarizer) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("text is empty")
	}
	
	if maxWords <= 0 {
		maxWords = 50 // Default fallback
	}
	
	// Split text into words
	words := strings.Fields(strings.TrimSpace(text))
	
	// If text has fewer words than maxWords, return all words
	if len(words) <= maxWords {
		return strings.Join(words, " "), nil
	}
	
	// Return first N words with ellipsis
	firstNWords := strings.Join(words[:maxWords], " ")
	return firstNWords + "...", nil
}

// SummarizeRequest represents a summarization request
type SummarizeRequest struct {
	Text       string `json:"text"`
	MaxLength  int    `json:"max_length,omitempty"`
	Style      string `json:"style,omitempty"` // e.g., "bullet_points", "paragraph"
	Language   string `json:"language,omitempty"`
}

// SummarizeResponse represents a summarization response
type SummarizeResponse struct {
	Summary    string  `json:"summary"`
	WordCount  int     `json:"word_count"`
	Confidence float64 `json:"confidence"`
}
