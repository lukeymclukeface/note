# Testing Documentation

This document describes the comprehensive testing strategy implemented for the note-server application, following the requirements to use Go's standard `testing` package, `httptest` for HTTP handlers, mocked interfaces, and `nhooyr.io/websocket` for WebSocket testing.

## Testing Architecture

### 1. HTTP Handler Testing with `httptest`

**Location:** `internal/http/handlers_test.go`

- **Framework:** Go's standard `testing` package + `httptest`
- **Coverage:** All HTTP endpoints including health, transcription, and summarization
- **Features:**
  - Request/response validation
  - Status code verification
  - JSON response structure validation
  - Multipart form handling for file uploads
  - Error handling scenarios
  - Method validation

**Example Test:**
```go
func TestTranscribeHandler(t *testing.T) {
    // Create mock services
    mockTranscriber := &MockTranscriber{
        TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
            return "Hello, world!", nil
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
    
    // Verify response
    assert.Equal(t, http.StatusOK, w.Code)
    // ... additional assertions
}
```

### 2. Mock Interfaces

**Locations:** 
- `internal/http/handlers_test.go`
- `internal/service/transcribe_test.go`
- `internal/service/summarize_test.go`

**Mock Transcriber Interface:**
```go
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
```

**Mock Summarizer Interface:**
```go
type MockSummarizer struct {
    SummarizeTextFunc func(ctx context.Context, text string, maxWords int) (string, error)
}

func (m *MockSummarizer) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
    if m.SummarizeTextFunc != nil {
        return m.SummarizeTextFunc(ctx, text, maxWords)
    }
    return "mock summary", nil
}
```

**Benefits:**
- Isolated testing without external dependencies
- Controlled error simulation
- Deterministic test results
- Fast test execution

### 3. WebSocket Testing with `nhooyr.io/websocket`

**Location:** `internal/ws/hub_test.go`

**Features:**
- Real WebSocket connection testing
- Binary message handling
- JSON message validation
- Connection lifecycle management
- Error handling scenarios
- Multiple connection testing

**Example WebSocket Test:**
```go
func TestTranscribeHub_WebSocketConnection(t *testing.T) {
    // Setup mock transcriber
    mockTranscriber := &MockTranscriber{
        TranscribeAudioFunc: func(ctx context.Context, audioData []byte) (string, error) {
            return "WebSocket test transcription", nil
        },
    }
    
    // Create hub and start it
    transcribeService := service.NewTranscribeServiceWithTranscriber(mockTranscriber)
    hub := NewTranscribeHub(transcribeService)
    go hub.Run()
    defer hub.Shutdown()

    // Create test server
    server := httptest.NewServer(http.HandlerFunc(hub.ServeTranscribeWS))
    defer server.Close()

    // Connect using nhooyr.io/websocket
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
    ctx := context.Background()
    conn, _, err := websocket.Dial(ctx, wsURL, nil)
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close(websocket.StatusNormalClosure, "test completed")

    // Send binary audio data
    audioData := []byte("fake audio data")
    err = conn.Write(ctx, websocket.MessageBinary, audioData)
    if err != nil {
        t.Fatalf("Failed to send data: %v", err)
    }

    // Read and verify response messages
    for {
        msgType, data, err := conn.Read(ctx)
        if err != nil {
            t.Fatalf("Failed to read message: %v", err)
        }
        
        if msgType == websocket.MessageText {
            var msg TranscribeMessage
            json.Unmarshal(data, &msg)
            
            if msg.Type == "final" {
                assert.Equal(t, "WebSocket test transcription", msg.Text)
                break
            }
        }
    }
}
```

### 4. Service Layer Testing

**Locations:**
- `internal/service/transcribe_test.go`
- `internal/service/summarize_test.go`

**Coverage:**
- Interface implementations
- Business logic validation
- Error handling
- Context cancellation
- Stream processing
- Edge cases

### 5. Utility Function Testing

**Location:** `internal/util/util_test.go`

**Coverage:**
- JSON response formatting
- Error handling utilities
- Timestamp generation
- HTTP response writing

### 6. Integration Testing

**Location:** `internal/integration_test.go`

**Features:**
- End-to-end workflow testing
- Multi-component interaction
- Real HTTP request/response cycles
- WebSocket integration
- Error propagation testing

## Test Execution

### Running Individual Test Suites

```bash
# HTTP handler tests
go test -v ./internal/http/...

# WebSocket tests
go test -v ./internal/ws/...

# Service layer tests
go test -v ./internal/service/...

# Utility tests
go test -v ./internal/util/...

# Integration tests
go test -v ./internal/integration_test.go
```

### Running All Tests

```bash
# All tests
go test ./internal/...

# With verbose output
go test -v ./internal/...

# With coverage
go test -cover ./internal/...
```

## Test Coverage Areas

### ✅ HTTP Handlers
- **Health endpoint:** GET requests, method validation
- **Transcription endpoint:** File upload, multipart forms, error scenarios
- **Summarization endpoint:** JSON request/response, validation, error handling
- **Router integration:** Middleware, routing, CORS

### ✅ WebSocket Functionality
- **Connection management:** Connect, disconnect, multiple clients
- **Message handling:** Binary audio data, JSON responses, streaming
- **Hub operations:** Client registration, message broadcasting
- **Error scenarios:** Invalid data, connection failures, timeouts

### ✅ Service Layer
- **Transcription service:** Audio processing, streaming, mocked implementations
- **Summarization service:** Text processing, word limits, edge cases
- **Interface compliance:** Mock implementations, dependency injection

### ✅ Utilities
- **JSON handling:** Response formatting, error responses, success responses
- **HTTP utilities:** Status codes, headers, content types
- **Time utilities:** Timestamp generation, formatting

### ✅ Integration
- **Full stack testing:** HTTP → Service → Mock implementations
- **Error propagation:** End-to-end error handling
- **Real protocols:** Actual HTTP and WebSocket connections

## Performance Testing

### Benchmarks Included

```bash
# Run benchmarks
go test -bench=. -benchmem ./internal/...
```

**Benchmark Coverage:**
- HTTP handler performance
- JSON serialization/deserialization
- Service layer operations
- WebSocket message processing

## Mocking Strategy

### Interface-Based Mocking
All external dependencies are abstracted behind interfaces:
- `Transcriber` interface for audio transcription
- `Summarizer` interface for text summarization

### Mock Implementation Benefits
1. **No External Dependencies:** Tests run without requiring actual transcription/summarization services
2. **Controlled Behavior:** Precise control over mock responses and errors
3. **Fast Execution:** No network calls or heavy processing
4. **Deterministic Results:** Consistent test outcomes
5. **Error Simulation:** Easy testing of error scenarios

### Mock Injection
Services accept mock implementations through constructor functions:
```go
// Production
service := NewTranscribeService()

// Testing with mock
mockTranscriber := &MockTranscriber{...}
service := NewTranscribeServiceWithTranscriber(mockTranscriber)
```

## Best Practices Implemented

1. **Test Organization:** Clear separation of unit, integration, and benchmark tests
2. **Mock Management:** Reusable mock implementations across test files
3. **Error Testing:** Comprehensive error scenario coverage
4. **Resource Cleanup:** Proper cleanup of resources (connections, servers, etc.)
5. **Timeout Handling:** Appropriate timeouts for async operations
6. **Context Usage:** Proper context propagation and cancellation testing
7. **Test Data:** Deterministic test data for consistent results

## Dependencies

```go
// Testing dependencies
require (
    "nhooyr.io/websocket" // WebSocket client for testing
    // Standard library: testing, httptest, context, json, etc.
)
```

## Continuous Integration

The test suite is designed to run in CI environments:
- No external service dependencies
- Deterministic execution
- Proper timeout handling
- Clear pass/fail indicators
- Coverage reporting

This comprehensive testing approach ensures reliability, maintainability, and confidence in the note-server application's functionality.
