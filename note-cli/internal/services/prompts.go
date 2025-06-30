package services

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed prompts
var promptFiles embed.FS

// PromptService handles loading and processing of OpenAI prompts
type PromptService struct{}

// NewPromptService creates a new prompt service instance
func NewPromptService() *PromptService {
	return &PromptService{}
}

// PromptData represents data to be injected into prompt templates
type PromptData struct {
	Content string
}

// GetContentAnalysisPrompts returns the system and user prompts for content analysis
func (s *PromptService) GetContentAnalysisPrompts(content string) (systemPrompt, userPrompt string, err error) {
	systemPrompt, err = s.loadPromptFile("prompts/content_analysis_system.md")
	if err != nil {
		return "", "", fmt.Errorf("failed to load content analysis system prompt: %w", err)
	}

	userTemplate, err := s.loadPromptFile("prompts/content_analysis_user.md")
	if err != nil {
		return "", "", fmt.Errorf("failed to load content analysis user prompt: %w", err)
	}

	userPrompt, err = s.processTemplate(userTemplate, PromptData{Content: content})
	if err != nil {
		return "", "", fmt.Errorf("failed to process content analysis user prompt: %w", err)
	}

	return systemPrompt, userPrompt, nil
}

// GetSummaryPrompts returns the system and user prompts for summarization based on content type
func (s *PromptService) GetSummaryPrompts(contentType, content string) (systemPrompt, userPrompt string, err error) {
	var systemFile, userFile string

	switch contentType {
	case "meeting":
		systemFile = "prompts/meeting_system.md"
		userFile = "prompts/meeting_user.md"
	case "interview":
		systemFile = "prompts/interview_system.md"
		userFile = "prompts/interview_user.md"
	case "lecture":
		systemFile = "prompts/lecture_system.md"
		userFile = "prompts/lecture_user.md"
	case "presentation":
		systemFile = "prompts/presentation_system.md"
		userFile = "prompts/presentation_user.md"
	case "conversation":
		systemFile = "prompts/conversation_system.md"
		userFile = "prompts/conversation_user.md"
	default: // "other" or unknown
		systemFile = "prompts/general_system.md"
		userFile = "prompts/general_user.md"
	}

	systemPrompt, err = s.loadPromptFile(systemFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to load %s system prompt: %w", contentType, err)
	}

	userTemplate, err := s.loadPromptFile(userFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to load %s user prompt: %w", contentType, err)
	}

	userPrompt, err = s.processTemplate(userTemplate, PromptData{Content: content})
	if err != nil {
		return "", "", fmt.Errorf("failed to process %s user prompt: %w", contentType, err)
	}

	return systemPrompt, userPrompt, nil
}

// GetLegacySummaryPrompts returns the legacy summary prompts for backward compatibility
func (s *PromptService) GetLegacySummaryPrompts(content string) (systemPrompt, userPrompt string, err error) {
	systemPrompt, err = s.loadPromptFile("prompts/legacy_summary_system.md")
	if err != nil {
		return "", "", fmt.Errorf("failed to load legacy summary system prompt: %w", err)
	}

	userTemplate, err := s.loadPromptFile("prompts/legacy_summary_user.md")
	if err != nil {
		return "", "", fmt.Errorf("failed to load legacy summary user prompt: %w", err)
	}

	userPrompt, err = s.processTemplate(userTemplate, PromptData{Content: content})
	if err != nil {
		return "", "", fmt.Errorf("failed to process legacy summary user prompt: %w", err)
	}

	return systemPrompt, userPrompt, nil
}

// loadPromptFile loads a prompt file from the embedded filesystem
func (s *PromptService) loadPromptFile(filename string) (string, error) {
	data, err := promptFiles.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", filename, err)
	}
	return string(data), nil
}

// processTemplate processes a template string with the provided data
func (s *PromptService) processTemplate(templateStr string, data PromptData) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
