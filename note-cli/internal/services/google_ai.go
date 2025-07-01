package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"note-cli/internal/config"
	"os/exec"
	"strings"
	"time"
)

// GoogleAIService handles interactions with Google AI
type GoogleAIService struct {
	projectID        string
	location         string
	transcriptionModel string
	summaryModel       string
	verboseLogger      *VerboseLogger
}

// NewGoogleAIService creates a new Google AI service instance
func NewGoogleAIService(verboseLogger *VerboseLogger) (*GoogleAIService, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Set default models for Google AI
	transcriptionModel := "gemini-1.5-flash"
	summaryModel := "gemini-1.5-flash"

	if verboseLogger != nil {
		verboseLogger.Config("Google Project ID", cfg.GoogleProjectID)
		verboseLogger.Config("Google Location", cfg.GoogleLocation)
		verboseLogger.Config("Transcription Model", transcriptionModel)
		verboseLogger.Config("Summary Model", summaryModel)
	}

	return &GoogleAIService{
		projectID:          cfg.GoogleProjectID,
		location:           cfg.GoogleLocation,
		transcriptionModel: transcriptionModel,
		summaryModel:       summaryModel,
		verboseLogger:      verboseLogger,
	}, nil
}

// TranscribeAudioFile transcribes an audio file
// Note: Google Vertex AI Gemini models don't directly support audio transcription.
// This would typically use Google's Speech-to-Text API instead.
func (s *GoogleAIService) TranscribeAudioFile(filePath string) (string, error) {
	return "", fmt.Errorf("audio transcription is not supported by Google Vertex AI Gemini models. Please use OpenAI for transcription or configure Google Speech-to-Text API separately")
}

// SummarizeText summarizes text
func (s *GoogleAIService) SummarizeText(text string) (string, error) {
	// Truncate text if it's too long
	maxInputLength := 100000 // About 25k tokens
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "\n\n[Content truncated due to length...]"
	}

	prompt := fmt.Sprintf("Please create a comprehensive summary of the following text:\n\n%s", text)
	return s.generateContent(prompt, s.summaryModel)
}

// AnalyzeContentAndGenerateTitle analyzes content and generates a title
func (s *GoogleAIService) AnalyzeContentAndGenerateTitle(content string) (*ContentAnalysis, error) {
	// Truncate content if it's too long for analysis
	maxInputLength := 50000 // About 12k tokens for analysis
	analysisContent := content
	if len(content) > maxInputLength {
		analysisContent = content[:maxInputLength] + "\n\n[Content truncated for analysis...]"
	}

	prompt := fmt.Sprintf(`Analyze the following content and determine its type and generate an appropriate title.

Content types to choose from:
- meeting: For meeting notes, discussions, team calls
- interview: For job interviews, candidate evaluations, Q&A sessions
- lecture: For educational content, presentations, talks
- note: For general notes, personal thoughts, documentation
- general: For any other type of content

Please respond with a JSON object containing the content type and title:
{"content_type": "meeting", "title": "Weekly Team Standup - Project Updates"}

Content to analyze:
%s`, analysisContent)

	response, err := s.generateContent(prompt, s.summaryModel)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var analysis ContentAnalysis
	if err := json.Unmarshal([]byte(strings.TrimSpace(response)), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}

	return &analysis, nil
}

// SummarizeByContentType summarizes text by content type
func (s *GoogleAIService) SummarizeByContentType(text string, contentType string) (string, error) {
	// Truncate text if it's too long
	maxInputLength := 100000 // About 25k tokens
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "\n\n[Content truncated due to length...]"
	}

	var prompt string
	switch contentType {
	case "meeting":
		prompt = fmt.Sprintf(`Create a detailed meeting summary of the following content. Include:
- Key discussion points
- Decisions made
- Action items
- Next steps

Content:
%s`, text)
	case "interview":
		prompt = fmt.Sprintf(`Create a detailed interview summary of the following content. Include:
- Candidate responses
- Key qualifications discussed
- Interview highlights
- Assessment notes

Content:
%s`, text)
	case "lecture":
		prompt = fmt.Sprintf(`Create a detailed lecture summary of the following content. Include:
- Main topics covered
- Key concepts
- Important details
- Learning objectives

Content:
%s`, text)
	default:
		prompt = fmt.Sprintf("Please create a comprehensive summary of the following %s content:\n\n%s", contentType, text)
	}

	return s.generateContent(prompt, s.summaryModel)
}

// GetProviderName returns the provider name
func (s *GoogleAIService) GetProviderName() string {
	return "google"
}

// GetAvailableModels returns available models
func (s *GoogleAIService) GetAvailableModels() ([]Model, error) {
	// Return common Gemini models available in Vertex AI
	models := []Model{
		{
			ID:          "gemini-1.5-flash",
			Name:        "Gemini 1.5 Flash",
			Description: "Fast and efficient model for text generation and analysis",
			Provider:    "google",
			Type:        "chat",
		},
		{
			ID:          "gemini-1.5-pro",
			Name:        "Gemini 1.5 Pro",
			Description: "Advanced model for complex text generation and analysis",
			Provider:    "google",
			Type:        "chat",
		},
		{
			ID:          "gemini-1.0-pro",
			Name:        "Gemini 1.0 Pro",
			Description: "Standard model for text generation and analysis",
			Provider:    "google",
			Type:        "chat",
		},
	}

	return models, nil
}

// fetchAccessToken fetches the access token from gcloud
func (s *GoogleAIService) fetchAccessToken() (string, error) {
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// generateContent makes API calls to Google's Vertex AI
func (s *GoogleAIService) generateContent(prompt, model string) (string, error) {
	token, err := s.fetchAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to fetch access token: %w", err)
	}

	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:generateContent",
		s.location, s.projectID, s.location, model)

	// Prepare request body for Vertex AI Gemini API
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 8192,
			"temperature":     0.7,
			"topP":            0.8,
			"topK":            40,
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	if s.verboseLogger != nil {
		s.verboseLogger.Step("Sending request to Google AI", fmt.Sprintf("Model: %s, Prompt length: %d chars", model, len(prompt)))
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		if s.verboseLogger != nil {
			s.verboseLogger.Error(err, "Google AI request failed")
		}
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if s.verboseLogger != nil {
		s.verboseLogger.API("POST", url, requestBody, string(respBody), resp.StatusCode, duration)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s - %s", resp.Status, string(respBody))
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the generated text from the response
	candidates, ok := response["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid candidate structure")
	}

	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content structure")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("no parts in content")
	}

	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid part structure")
	}

	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text in part")
	}

	return strings.TrimSpace(text), nil
}
