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

	url := "https://api.openai.com/v1/chat/completions"
	
	// Create request payload
	requestPayload := map[string]interface{}{
		"model": s.summaryModel,
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": `You are a content analysis assistant. Your task is to analyze text content and determine its type, then generate an appropriate title.

Analyze the provided content and respond with a JSON object containing:
1. "content_type": one of "meeting", "interview", "lecture", "conversation", "presentation", "other"
2. "title": a concise, descriptive title (max 60 characters) that captures the essence of the content

For meetings: focus on main topics discussed
For interviews: focus on the interviewee and main subject
For lectures: focus on the topic being taught
For conversations: focus on the main discussion points
For presentations: focus on the subject being presented
For other: create a general descriptive title

Respond ONLY with valid JSON, no additional text.`,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(`Analyze this content and provide the JSON response:\n\n%s`, analysisContent),
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

	// Create specialized prompts based on content type
	var systemPrompt, userPrompt string

	switch contentType {
	case "meeting":
		systemPrompt = `You are a meeting minutes assistant. Create professional meeting summaries that capture key decisions, action items, and important discussions. Use clear structure and markdown formatting.`
		userPrompt = `Please create a comprehensive meeting summary with the following structure:

# Meeting Summary

## Overview
Brief description of the meeting purpose and participants (if identifiable)

## Key Topics Discussed
Main agenda items and discussion points

## Decisions Made
Important decisions and conclusions reached

## Action Items
Tasks assigned with responsible parties and deadlines (if mentioned)

## Next Steps
Planned follow-up actions or future meetings

Here is the meeting content:

%s`

	case "interview":
		systemPrompt = `You are an interview summary assistant. Create structured summaries of interviews that capture key questions, responses, and insights. Use clear formatting and highlight important points.`
		userPrompt = `Please create a comprehensive interview summary with the following structure:

# Interview Summary

## Overview
Brief description of the interview context and participants

## Key Questions & Responses
Main questions asked and significant responses given

## Key Insights
Important insights, opinions, or information shared

## Notable Quotes
Significant quotes or statements (if any)

## Conclusion
Main takeaways from the interview

Here is the interview content:

%s`

	case "lecture":
		systemPrompt = `You are an educational content assistant. Create structured summaries of lectures that capture key concepts, learning objectives, and important information. Use clear academic formatting.`
		userPrompt = `Please create a comprehensive lecture summary with the following structure:

# Lecture Summary

## Topic Overview
Main subject and learning objectives

## Key Concepts
Important concepts and theories covered

## Main Points
Detailed breakdown of the lecture content

## Examples & Illustrations
Key examples or case studies mentioned

## Conclusion
Summary of main takeaways and learning outcomes

Here is the lecture content:

%s`

	case "presentation":
		systemPrompt = `You are a presentation summary assistant. Create structured summaries of presentations that capture key points, data, and conclusions. Use clear formatting appropriate for business contexts.`
		userPrompt = `Please create a comprehensive presentation summary with the following structure:

# Presentation Summary

## Overview
Main topic and presentation objectives

## Key Points
Main arguments and supporting information

## Data & Evidence
Important statistics, facts, or evidence presented

## Conclusions
Main conclusions and recommendations

## Call to Action
Next steps or actions recommended (if any)

Here is the presentation content:

%s`

	case "conversation":
		systemPrompt = `You are a conversation summary assistant. Create structured summaries of conversations that capture main topics, participant perspectives, and key outcomes. Use clear formatting.`
		userPrompt = `Please create a comprehensive conversation summary with the following structure:

# Conversation Summary

## Participants
Key participants in the conversation (if identifiable)

## Main Topics
Primary subjects discussed

## Key Points
Important points raised by different participants

## Agreements & Disagreements
Areas of consensus and differing viewpoints

## Outcomes
Any decisions, agreements, or next steps discussed

Here is the conversation content:

%s`

	default: // "other" or unknown
		systemPrompt = `You are a general content summary assistant. Create clear, structured summaries that capture the main points and important information from various types of content. Use appropriate markdown formatting.`
		userPrompt = `Please create a comprehensive summary with the following structure:

# Content Summary

## Overview
Brief description of the content type and main subject

## Key Points
Main topics and important information covered

## Details
Significant details and supporting information

## Conclusion
Main takeaways and summary of the content

Here is the content:

%s`
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
				"content": fmt.Sprintf(userPrompt, text),
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
