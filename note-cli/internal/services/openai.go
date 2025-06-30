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
	// Truncate text if it's too long for the API (rough estimate: 1 token ≈ 4 characters)
	// GPT-3.5-turbo has a 4096 token limit, leaving room for system message and response
	maxInputLength := 100000 // About 3000 tokens
	if len(text) > maxInputLength {
		text = text[:maxInputLength] + "\n\n[Content truncated due to length...]"
	}

	url := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": s.summaryModel,
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": `You are a helpful assistant that summarizes text. 
Your task is to create concise and informative summaries of various types of content, including meeting notes, interviews, and conversations. 
Please ensure your summaries are clear and structured. 
All outputs should be in markdown format only with no wrapping of the response with any explanations of the output.`,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(`Please provide a summary of the following text.

If the text is of a meeting, summarize the key topics covered and any important conclusions or action items. The output should include the following sections:
# Title
### Overview
### Key topics
### Outcome
### Action items (with responsible parties and deadlines if mentioned)

If the text is of an interview, summarize the main questions asked and the responses given.
If the text is of a conversation, please summarize the main points discussed and identify key speakers if possible.

Here is the text:
		
%s`, text),
			},
		},
		"max_completion_tokens": 100000, // Allow for a longer summary
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
