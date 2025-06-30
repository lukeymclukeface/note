package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"note-cli/internal/config"
	"os"
	"path/filepath"
	"strings"
)

// OpenAIService handles all OpenAI API interactions
type OpenAIService struct {
	apiKey             string
	transcriptionModel string
	summaryModel       string
	promptService      *PromptService
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService() (*OpenAIService, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.OpenAIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured. Please run 'note setup' first")
	}

	transcriptionModel := cfg.TranscriptionModel
	if transcriptionModel == "" {
		transcriptionModel = "whisper-1"
	}

	summaryModel := cfg.SummaryModel
	if summaryModel == "" {
		summaryModel = "gpt-3.5-turbo"
	}

	return &OpenAIService{
		apiKey:             cfg.OpenAIKey,
		transcriptionModel: transcriptionModel,
		summaryModel:       summaryModel,
		promptService:      NewPromptService(),
	}, nil
}

// TranscribeAudioFile transcribes an audio file using OpenAI Whisper
func (s *OpenAIService) TranscribeAudioFile(filePath string) (string, error) {
	url := "https://api.openai.com/v1/audio/transcriptions"
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("model", s.transcriptionModel)
	if err != nil {
		return "", err
	}

	// Add prompt to encourage speaker identification
	speakerPrompt := "The following audio contains multiple speakers. Please transcribe the entire audio and identify speakers as Speaker 1, Speaker 2, etc. when possible."
	err = writer.WriteField("prompt", speakerPrompt)
	if err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	transcript, ok := result["text"].(string)
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	return transcript, nil
}

// SummarizeText creates a summary of the provided text using OpenAI
func (s *OpenAIService) SummarizeText(text string) (string, error) {
	// Truncate text if it's too long for the API
	maxInputLength := 100000 // About 25k tokens
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "\n\n[Content truncated due to length...]"
	}

	// Get legacy prompts for backward compatibility
	systemPrompt, userPrompt, err := s.promptService.GetLegacySummaryPrompts(text)
	if err != nil {
		return "", fmt.Errorf("failed to load legacy summary prompts: %w", err)
	}

	url := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": s.summaryModel,
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": userPrompt,
			},
		},
		"max_completion_tokens": 100000,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the error response body for more details
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected response status: %s - %s", resp.Status, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("unexpected response structure")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	return strings.TrimSpace(content), nil
}

// OpenAIModel represents an OpenAI model
type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// OpenAIModelsResponse represents the response from OpenAI models endpoint
type OpenAIModelsResponse struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}

// GetAvailableModels fetches available models from OpenAI API
func (s *OpenAIService) GetAvailableModels() ([]OpenAIModel, error) {
	url := "https://api.openai.com/v1/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var modelsResponse OpenAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return modelsResponse.Data, nil
}

// ContentAnalysis represents the result of content analysis
type ContentAnalysis struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
}

// AnalyzeContentAndGenerateTitle analyzes content to determine its type and generates an appropriate title
func (s *OpenAIService) AnalyzeContentAndGenerateTitle(content string) (*ContentAnalysis, error) {
	// Truncate content if it's too long for analysis
	maxInputLength := 50000 // About 12k tokens for analysis
	analysisContent := content
	if len(content) > maxInputLength {
		analysisContent = content[:maxInputLength] + "\n\n[Content truncated for analysis...]"
	}

	// Get content analysis prompts
	systemPrompt, userPrompt, err := s.promptService.GetContentAnalysisPrompts(analysisContent)
	if err != nil {
		return nil, fmt.Errorf("failed to load content analysis prompts: %w", err)
	}

	url := "https://api.openai.com/v1/chat/completions"
	
	// Create request payload
	requestPayload := map[string]interface{}{
		"model": s.summaryModel,
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": userPrompt,
			},
		},
		"max_completion_tokens": 200,
	}
	
	// Only add temperature for models that support it (not o3-mini)
	if !strings.Contains(strings.ToLower(s.summaryModel), "o3") {
		requestPayload["temperature"] = 0.3
	}
	
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response status: %s - %s", resp.Status, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, errors.New("unexpected response structure")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response structure")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response structure")
	}

	content_response, ok := message["content"].(string)
	if !ok {
		return nil, errors.New("unexpected response structure")
	}

	// Parse the JSON response
	var analysis ContentAnalysis
	if err := json.Unmarshal([]byte(strings.TrimSpace(content_response)), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis response: %w", err)
	}

	return &analysis, nil
}

// SummarizeByContentType creates a specialized summary based on the content type
func (s *OpenAIService) SummarizeByContentType(text string, contentType string) (string, error) {
	// Truncate text if it's too long for the API
	maxInputLength := 100000 // About 25k tokens
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "\n\n[Content truncated due to length...]"
	}

	// Get specialized prompts based on content type
	systemPrompt, userPrompt, err := s.promptService.GetSummaryPrompts(contentType, text)
	if err != nil {
		return "", fmt.Errorf("failed to load summary prompts for %s: %w", contentType, err)
	}

	url := "https://api.openai.com/v1/chat/completions"
	
	// Create request payload
	requestPayload := map[string]interface{}{
		"model": s.summaryModel,
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": systemPrompt,
			},
			{
				"role": "user",
				"content": userPrompt,
			},
		},
		"max_completion_tokens": 100000,
	}
	
	// Only add temperature for models that support it (not o3-mini)
	if !strings.Contains(strings.ToLower(s.summaryModel), "o3") {
		requestPayload["temperature"] = 0.7
	}
	
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected response status: %s - %s", resp.Status, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("unexpected response structure")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	return strings.TrimSpace(content), nil
}
