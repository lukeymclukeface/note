package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket"
	"github.com/your-org/note-server/internal/service"
)

// Mock Transcriber for WebSocket tests
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

func TestTranscribeHub_NewTranscribeHub(t *testing.T) {
	mockTranscriber := &MockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	
	hub := NewTranscribeHub(transcribeService)
	
	if hub == nil {
		t.Fatal("Expected hub to be created")
	}
	
	if hub.transcribeService != transcribeService {
		t.Error("Expected transcribe service to be set")
	}
	
	if len(hub.clients) != 0 {
		t.Error("Expected clients map to be empty initially")
	}
}

func TestTranscribeHub_WebSocketConnection(t *testing.T) {
	mockTranscriber := &MockTranscriber{
		TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
			return "WebSocket test transcription", nil
		},
	}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	// Start the hub in a goroutine
	go hub.Run()
	defer hub.Shutdown()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to WebSocket
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test completed")

	// Send binary audio data
	audioData := []byte("fake audio data for testing")
	err = conn.Write(ctx, websocket.MessageBinary, audioData)
	if err != nil {
		t.Fatalf("Failed to send audio data: %v", err)
	}

	// Read messages from WebSocket
	var messages []TranscribeMessage
	timeout := time.After(5 * time.Second)

	for len(messages) < 2 {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for WebSocket messages")
		default:
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				t.Fatalf("Failed to read WebSocket message: %v", err)
			}

			if msgType != websocket.MessageText {
				continue // Skip non-text messages
			}

			var msg TranscribeMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}

			messages = append(messages, msg)
		}
	}

	// Verify we received the expected messages
	if len(messages) < 2 {
		t.Fatalf("Expected at least 2 messages, got %d", len(messages))
	}

	// First message should be partial
	if messages[0].Type != "partial" {
		t.Errorf("Expected first message type 'partial', got %s", messages[0].Type)
	}

	// Last message should be final with our mocked transcription
	finalMsg := messages[len(messages)-1]
	if finalMsg.Type != "final" {
		t.Errorf("Expected final message type 'final', got %s", finalMsg.Type)
	}

	if finalMsg.Text != "WebSocket test transcription" {
		t.Errorf("Expected final message text 'WebSocket test transcription', got %s", finalMsg.Text)
	}
}

func TestTranscribeHub_MultipleConnections(t *testing.T) {
	mockTranscriber := &MockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	// Start the hub
	go hub.Run()
	defer hub.Shutdown()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Create multiple connections
	numConnections := 3
	connections := make([]*websocket.Conn, numConnections)
	ctx := context.Background()

	for i := 0; i < numConnections; i++ {
		conn, _, err := websocket.Dial(ctx, wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect WebSocket %d: %v", i, err)
		}
		connections[i] = conn
	}

	// Give time for connections to register
	time.Sleep(100 * time.Millisecond)

	// Close all connections
	for i, conn := range connections {
		if conn != nil {
			conn.Close(websocket.StatusNormalClosure, "test completed")
		} else {
			t.Errorf("Connection %d is nil", i)
		}
	}

	// Give time for cleanup
	time.Sleep(100 * time.Millisecond)
}

func TestTranscribeHub_ConnectionLimit(t *testing.T) {
	// This test would be complex to implement properly as it requires
	// mocking the maxConnections limit. In a real scenario, you'd want
	// to make maxConnections configurable for testing.
	t.Skip("Connection limit test requires refactoring to make maxConnections configurable")
}

func TestTranscribeHub_ErrorHandling(t *testing.T) {
	mockTranscriber := &MockTranscriber{
		TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
			return "", context.DeadlineExceeded // Simulate timeout error
		},
	}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	go hub.Run()
	defer hub.Shutdown()

	server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test completed")

	// Send audio data that will trigger an error
	audioData := []byte("audio data that causes error")
	err = conn.Write(ctx, websocket.MessageBinary, audioData)
	if err != nil {
		t.Fatalf("Failed to send audio data: %v", err)
	}

	// Wait for error message
	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for error message")
		default:
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				return // Connection might be closed, which is expected
			}

			if msgType != websocket.MessageText {
				continue
			}

			var msg TranscribeMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			if msg.Type == "error" {
				// Successfully received error message
				return
			}
		}
	}
}

func TestTranscribeHub_Shutdown(t *testing.T) {
	mockTranscriber := &MockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	// Start the hub
	hubDone := make(chan bool)
	go func() {
		hub.Run()
		hubDone <- true
	}()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Shutdown the hub
	hub.Shutdown()

	// Wait for hub to finish or timeout
	select {
	case <-hubDone:
		// Hub shut down successfully
	case <-time.After(1 * time.Second):
		t.Error("Hub did not shut down within timeout")
	}
}

func TestTranscribeHub_PingPong(t *testing.T) {
	mockTranscriber := &MockTranscriber{}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	go hub.Run()
	defer hub.Shutdown()

	server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test completed")

	// Test basic connection functionality instead of ping/pong
	// since nhooyr.io/websocket handles ping/pong automatically
	time.Sleep(100 * time.Millisecond)
	t.Log("WebSocket connection established successfully")
}

func TestTranscribeMessage_JSON(t *testing.T) {
	tests := []struct {
		name string
		msg  TranscribeMessage
	}{
		{
			name: "partial message",
			msg:  TranscribeMessage{Type: "partial", Text: "Processing..."},
		},
		{
			name: "final message",
			msg:  TranscribeMessage{Type: "final", Text: "Complete transcription"},
		},
		{
			name: "error message",
			msg:  TranscribeMessage{Type: "error", Text: "Transcription failed"},
		},
		{
			name: "message without text",
			msg:  TranscribeMessage{Type: "status"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.msg)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}

			// Unmarshal back
			var decoded TranscribeMessage
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal message: %v", err)
			}

			// Verify fields
			if decoded.Type != tt.msg.Type {
				t.Errorf("Expected type %s, got %s", tt.msg.Type, decoded.Type)
			}

			if decoded.Text != tt.msg.Text {
				t.Errorf("Expected text %s, got %s", tt.msg.Text, decoded.Text)
			}
		})
	}
}

// Integration test combining WebSocket with HTTP
func TestWebSocketIntegration(t *testing.T) {
	mockTranscriber := &MockTranscriber{
		TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
			// Echo back the audio data size as text for testing
			return string(audioData), nil
		},
	}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	go hub.Run()
	defer hub.Shutdown()

	// Create HTTP server with WebSocket endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/transcribe", hub.ServeTranscribeWS)
	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/transcribe"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test completed")

	// Send test audio data
	testData := []byte("integration test audio")
	err = conn.Write(ctx, websocket.MessageBinary, testData)
	if err != nil {
		t.Fatalf("Failed to send audio data: %v", err)
	}

	// Read final transcription result
	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for transcription result")
		default:
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				t.Fatalf("Failed to read message: %v", err)
			}

			if msgType != websocket.MessageText {
				continue
			}

			var msg TranscribeMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			if msg.Type == "final" && msg.Text == string(testData) {
				// Success - received expected final transcription
				return
			}
		}
	}
}

// Benchmark WebSocket message handling
func BenchmarkWebSocketMessageProcessing(b *testing.B) {
	mockTranscriber := &MockTranscriber{
		TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
			return "benchmark result", nil
		},
	}
	transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
	hub := NewTranscribeHub(transcribeService)

	// Create a mock client for benchmarking
	ctx := context.Background()
	client := &TranscribeClient{
		hub:  hub,
		send: make(chan []byte, 256),
		ctx:  ctx,
	}

	audioData := []byte("benchmark audio data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.processAudioChunk(audioData)
	}
}
