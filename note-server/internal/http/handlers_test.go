package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/note-server/internal/service"
	"github.com/your-org/note-server/internal/ws"
)

// Mock Transcriber implementation
type MockTranscriber struct {
	TranscribeAudioFunc  func(ctx context.Context, audioData []byte) (string, error)
	TranscribeStreamFunc func(ctx context.Context, audioChunk []byte) (<-chan service.TranscribeStreamResult, error)
}

func (m *MockTranscriber) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	if m.TranscribeAudioFunc != nil {
		return m.TranscribeAudioFunc(ctx, audioData)
	}
	return "mock transcription result", nil
}

func (m *MockTranscriber) TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan service.TranscribeStreamResult, error) {
	if m.TranscribeStreamFunc != nil {
		return m.TranscribeStreamFunc(ctx, audioChunk)
	}
	// Default mock implementation
	resultChan := make(chan service.TranscribeStreamResult, 2)
	go func() {
		defer close(resultChan)
		resultChan <- service.TranscribeStreamResult{Type: "partial", Text: "mock partial"}
		resultChan <- service.TranscribeStreamResult{Type: "final", Text: "mock final result"}
	}()
	return resultChan, nil
}

// Mock Summarizer implementation
type MockSummarizer struct {
	SummarizeTextFunc func(ctx context.Context, text string, maxWords int) (string, error)
}

func (m *MockSummarizer) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
	if m.SummarizeTextFunc != nil {
		return m.SummarizeTextFunc(ctx, text, maxWords)
	}
	return "mock summary", nil
}

// Helper function to create handlers with mocked services
func createHandlersWithMocks(transcriber service.Transcriber, summarizer service.Summarizer) *Handlers {
	transcribeService := service.NewTranscribeServiceWithTranscriber(transcriber)
	summarizeService := service.NewSummarizeServiceWithSummarizer(summarizer, 50)
	
	return &Handlers{
		transcribeService: transcribeService,
		summarizeService:  summarizeService,
	}
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET health check should return OK",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "POST to health check should return method not allowed",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewHandlers()
			req := httptest.NewRequest(tt.method, "/healthz", nil)
			w := httptest.NewRecorder()

			handlers.HealthHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestTranscribeHandler(t *testing.T) {
	t.Run("successful transcription", func(t *testing.T) {
		mockTranscriber := &MockTranscriber{
			TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
				return "Hello, world!", nil
			},
		}
		
		handlers := createHandlersWithMocks(mockTranscriber, &MockSummarizer{})

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		// Add file field
		fileWriter, err := writer.CreateFormFile("file", "test.wav")
		if err != nil {
			t.Fatal(err)
		}
		
		// Write fake audio data
		audioData := []byte("fake audio data")
		_, err = fileWriter.Write(audioData)
		if err != nil {
			t.Fatal(err)
		}
		
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/transcribe", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handlers.TranscribeHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]any
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		data, ok := response["data"].(map[string]any)
		if !ok {
			t.Fatal("expected data field in response")
		}

		if text, ok := data["text"].(string); !ok || text != "Hello, world!" {
			t.Errorf("expected text 'Hello, world!', got %v", data["text"])
		}

		if _, ok := data["duration_ms"]; !ok {
			t.Error("expected duration_ms field in response")
		}
	})

	t.Run("transcription service error", func(t *testing.T) {
		mockTranscriber := &MockTranscriber{
			TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
				return "", fmt.Errorf("transcription failed")
			},
		}
		
		handlers := createHandlersWithMocks(mockTranscriber, &MockSummarizer{})

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, _ := writer.CreateFormFile("file", "test.wav")
		fileWriter.Write([]byte("fake audio data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/transcribe", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handlers.TranscribeHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		handlers := NewHandlers()
		req := httptest.NewRequest(http.MethodGet, "/transcribe", nil)
		w := httptest.NewRecorder()

		handlers.TranscribeHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("no file provided", func(t *testing.T) {
		handlers := NewHandlers()
		
		// Create empty multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/transcribe", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		handlers.TranscribeHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid multipart form", func(t *testing.T) {
		handlers := NewHandlers()
		req := httptest.NewRequest(http.MethodPost, "/transcribe", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "multipart/form-data")
		w := httptest.NewRecorder()

		handlers.TranscribeHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestSummarizeHandler(t *testing.T) {
	t.Run("successful summarization", func(t *testing.T) {
		mockSummarizer := &MockSummarizer{
			SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
				return "This is a summary of: " + text, nil
			},
		}
		
		handlers := createHandlersWithMocks(&MockTranscriber{}, mockSummarizer)

		requestBody := SummarizeRequest{
			Text: "This is a long text that needs to be summarized.",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.SummarizeHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]any
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		data, ok := response["data"].(map[string]any)
		if !ok {
			t.Fatal("expected data field in response")
		}

		expectedSummary := "This is a summary of: This is a long text that needs to be summarized."
		if summary, ok := data["summary"].(string); !ok || summary != expectedSummary {
			t.Errorf("expected summary %q, got %v", expectedSummary, data["summary"])
		}
	})

	t.Run("summarization service error", func(t *testing.T) {
		mockSummarizer := &MockSummarizer{
			SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
				return "", fmt.Errorf("summarization failed")
			},
		}
		
		handlers := createHandlersWithMocks(&MockTranscriber{}, mockSummarizer)

		requestBody := SummarizeRequest{Text: "Test text"}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.SummarizeHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		handlers := NewHandlers()
		req := httptest.NewRequest(http.MethodGet, "/summarize", nil)
		w := httptest.NewRecorder()

		handlers.SummarizeHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		handlers := NewHandlers()
		req := httptest.NewRequest(http.MethodPost, "/summarize", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.SummarizeHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("empty text field", func(t *testing.T) {
		handlers := NewHandlers()
		requestBody := SummarizeRequest{Text: ""}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handlers.SummarizeHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Integration test with the router
func TestHandlersIntegration(t *testing.T) {
	t.Run("health endpoint integration", func(t *testing.T) {
		transcribeHub := createMockTranscribeHub()
		router := NewRouter(transcribeHub)

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

// Helper function to create a mock transcribe hub
func createMockTranscribeHub() *ws.TranscribeHub {
	mockTranscriber := &MockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	return ws.NewTranscribeHub(transcribeService)
}

// Benchmarks
func BenchmarkHealthHandler(b *testing.B) {
	handlers := NewHandlers()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handlers.HealthHandler(w, req)
	}
}

func BenchmarkSummarizeHandler(b *testing.B) {
	mockSummarizer := &MockSummarizer{
		SummarizeTextFunc: func(ctx context.Context, text string, maxWords int) (string, error) {
			return "benchmark summary", nil
		},
	}
	
	handlers := createHandlersWithMocks(&MockTranscriber{}, mockSummarizer)
	
	requestBody := SummarizeRequest{Text: "This is test text for benchmarking."}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handlers.SummarizeHandler(w, req)
	}
}
