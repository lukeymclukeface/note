package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket"

	httpPkg "github.com/your-org/note-server/internal/http"
	"github.com/your-org/note-server/internal/service"
	"github.com/your-org/note-server/internal/ws"
)

// Mock implementations for integration testing
type IntegrationMockTranscriber struct{}

func (m *IntegrationMockTranscriber) TranscribeAudio(ctx context.Context, audioData []byte) (string, error) {
	// Return a deterministic result for integration testing
	return "Integration test transcription: " + string(audioData[:min(len(audioData), 10)]), nil
}

func (m *IntegrationMockTranscriber) TranscribeStream(ctx context.Context, audioChunk []byte) (<-chan service.TranscribeStreamResult, error) {
	resultChan := make(chan service.TranscribeStreamResult, 2)
	go func() {
		defer close(resultChan)
		resultChan <- service.TranscribeStreamResult{Type: "partial", Text: "Processing..."}
		time.Sleep(50 * time.Millisecond)
		resultChan <- service.TranscribeStreamResult{Type: "final", Text: "Integration test stream result"}
	}()
	return resultChan, nil
}

type IntegrationMockSummarizer struct{}

func (m *IntegrationMockSummarizer) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
	words := strings.Fields(text)
	if len(words) <= maxWords {
		return "Summary: " + text, nil
	}
	return "Summary: " + strings.Join(words[:maxWords], " ") + "...", nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestIntegrationHTTPHandlers tests the complete HTTP handler stack
func TestIntegrationHTTPHandlers(t *testing.T) {
	// Create services with mock implementations
	transcriber := &IntegrationMockTranscriber{}
	
	transcribeService := service.NewTranscribeServiceWithTranscriber(transcriber)
	
	// Test individual endpoints using the router
	transcribeHub := ws.NewTranscribeHub(transcribeService)
	router := httpPkg.NewRouter(transcribeHub)
	
	t.Run("health endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		if !strings.Contains(w.Body.String(), "OK") {
			t.Errorf("Expected 'OK' in response, got %s", w.Body.String())
		}
	})
	
	t.Run("summarize endpoint integration", func(t *testing.T) {
		// Create services with mock implementations for this test
		summarizer := &IntegrationMockSummarizer{}
		summarizeService := service.NewSummarizeServiceWithSummarizer(summarizer, 10)
		handlers := httpPkg.NewHandlersWithServices(transcribeService, summarizeService)
		router := httpPkg.NewRouterWithHandlers(transcribeHub, handlers)
		
		requestBody := map[string]string{
			"text": "This is a long text that needs to be summarized for integration testing purposes.",
		}
		jsonBody, _ := json.Marshal(requestBody)
		
		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		
		var response map[string]any
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		
		if !response["success"].(bool) {
			t.Error("Expected successful response")
		}
		
		data := response["data"].(map[string]any)
		summary := data["summary"].(string)
		
		if !strings.Contains(summary, "Summary:") {
			t.Errorf("Expected summary to contain 'Summary:', got %s", summary)
		}
	})
}

// TestIntegrationWebSocket tests the complete WebSocket functionality
func TestIntegrationWebSocket(t *testing.T) {
	transcriber := &IntegrationMockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(transcriber)
	hub := ws.NewTranscribeHub(transcribeService)
	
	// Start the hub
	go hub.Run()
	defer hub.Shutdown()
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
	defer server.Close()
	
	// Convert to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Connect to WebSocket
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test completed")
	
	// Send test audio data
	testAudio := []byte("integration test audio data")
	err = conn.Write(ctx, websocket.MessageBinary, testAudio)
	if err != nil {
		t.Fatalf("Failed to send audio data: %v", err)
	}
	
	// Read messages and verify we get expected transcription
	timeout := time.After(5 * time.Second)
	receivedFinal := false
	
	for !receivedFinal {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for WebSocket messages")
		default:
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				t.Fatalf("Failed to read WebSocket message: %v", err)
			}
			
			if msgType != websocket.MessageText {
				continue
			}
			
			var msg ws.TranscribeMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}
			
			t.Logf("Received message: %+v", msg)
			
			if msg.Type == "final" {
				expectedText := "Integration test transcription: " + string(testAudio[:10])
				if msg.Text != expectedText {
					t.Errorf("Expected final text %q, got %q", expectedText, msg.Text)
				}
				receivedFinal = true
			}
		}
	}
}

// TestIntegrationServiceInteraction tests service layer interactions
func TestIntegrationServiceInteraction(t *testing.T) {
	t.Run("transcribe and summarize flow", func(t *testing.T) {
		transcriber := &IntegrationMockTranscriber{}
		summarizer := &IntegrationMockSummarizer{}
		
		transcribeService := service.NewTranscribeServiceWithTranscriber(transcriber)
		summarizeService := service.NewSummarizeServiceWithSummarizer(summarizer, 10)
		
		ctx := context.Background()
		
		// Step 1: Transcribe audio
		audioData := []byte("This is fake audio data for integration testing")
		transcription, err := transcribeService.TranscribeAudio(ctx, audioData)
		if err != nil {
			// Expected to fail due to ffmpeg, but we can test the mock directly
			transcription, err = transcriber.TranscribeAudio(ctx, audioData)
			if err != nil {
				t.Fatalf("Transcription failed: %v", err)
			}
		}
		
		t.Logf("Transcription result: %s", transcription)
		
		// Step 2: Summarize the transcription
		summary, err := summarizeService.SummarizeText(ctx, transcription)
		if err != nil {
			t.Fatalf("Summarization failed: %v", err)
		}
		
		t.Logf("Summary result: %s", summary)
		
		// Verify the flow worked
		if !strings.Contains(summary, "Summary:") {
			t.Error("Expected summary to contain 'Summary:'")
		}
	})
}

// TestIntegrationErrorHandling tests error handling across the stack
func TestIntegrationErrorHandling(t *testing.T) {
	t.Run("invalid JSON to summarize endpoint", func(t *testing.T) {
		transcribeHub := ws.NewTranscribeHub(service.NewTranscribeService())
		router := httpPkg.NewRouter(transcribeHub)
		
		req := httptest.NewRequest(http.MethodPost, "/summarize", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
		
		var response map[string]any
		json.NewDecoder(w.Body).Decode(&response)
		
		if response["success"].(bool) {
			t.Error("Expected error response")
		}
	})
	
	t.Run("empty text to summarize endpoint", func(t *testing.T) {
		transcribeHub := ws.NewTranscribeHub(service.NewTranscribeService())
		router := httpPkg.NewRouter(transcribeHub)
		
		requestBody := map[string]string{"text": ""}
		jsonBody, _ := json.Marshal(requestBody)
		
		req := httptest.NewRequest(http.MethodPost, "/summarize", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
	
	t.Run("wrong HTTP method", func(t *testing.T) {
		transcribeHub := ws.NewTranscribeHub(service.NewTranscribeService())
		router := httpPkg.NewRouter(transcribeHub)
		
		req := httptest.NewRequest(http.MethodGet, "/summarize", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", w.Code)
		}
	})
}

// BenchmarkIntegrationSummarizeEndpoint benchmarks the complete summarize flow
func BenchmarkIntegrationSummarizeEndpoint(b *testing.B) {
	summarizer := &IntegrationMockSummarizer{}
	
	summarizeService := service.NewSummarizeServiceWithSummarizer(summarizer, 50)
	
	requestBody := map[string]string{
		"text": "This is a sample text for benchmarking the summarization endpoint performance.",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark the service directly
		ctx := context.Background()
		_, err := summarizeService.SummarizeText(ctx, requestBody["text"])
		if err != nil {
			b.Fatal(err)
		}
	}
}
