package services

import "fmt"

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	// TranscribeAudioFile transcribes an audio file
	TranscribeAudioFile(filePath string) (string, error)
	
	// SummarizeText creates a summary of the provided text
	SummarizeText(text string) (string, error)
	
	// AnalyzeContentAndGenerateTitle analyzes content to determine its type and generates an appropriate title
	AnalyzeContentAndGenerateTitle(content string) (*ContentAnalysis, error)
	
	// SummarizeByContentType creates a specialized summary based on the content type
	SummarizeByContentType(text string, contentType string) (string, error)
	
	// GetProviderName returns the name of the provider
	GetProviderName() string
	
	// GetAvailableModels returns available models for this provider
	GetAvailableModels() ([]Model, error)
}

// Model represents a generic AI model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Type        string `json:"type"` // "transcription", "chat", "both"
}

// ProviderFactory creates AI providers based on configuration
type ProviderFactory struct {
	verboseLogger *VerboseLogger
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(verboseLogger *VerboseLogger) *ProviderFactory {
	return &ProviderFactory{
		verboseLogger: verboseLogger,
	}
}

// CreateTranscriptionProvider creates a provider for transcription
func (f *ProviderFactory) CreateTranscriptionProvider(providerName string) (AIProvider, error) {
	switch providerName {
	case "openai":
		return NewOpenAIService(f.verboseLogger)
	case "google":
		return NewGoogleAIService(f.verboseLogger)
	default:
		return nil, fmt.Errorf("unsupported transcription provider: %s", providerName)
	}
}

// CreateSummaryProvider creates a provider for summarization
func (f *ProviderFactory) CreateSummaryProvider(providerName string) (AIProvider, error) {
	switch providerName {
	case "openai":
		return NewOpenAIService(f.verboseLogger)
	case "google":
		return NewGoogleAIService(f.verboseLogger)
	default:
		return nil, fmt.Errorf("unsupported summary provider: %s", providerName)
	}
}
