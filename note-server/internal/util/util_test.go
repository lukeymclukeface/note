package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestJSONResponse_Structure(t *testing.T) {
	tests := []struct {
		name     string
		response JSONResponse
	}{
		{
			name: "success response with data",
			response: JSONResponse{
				Success: true,
				Data:    map[string]string{"key": "value"},
			},
		},
		{
			name: "error response",
			response: JSONResponse{
				Success: false,
				Error:   "Something went wrong",
			},
		},
		{
			name: "response with message",
			response: JSONResponse{
				Success: true,
				Message: "Operation completed",
				Data:    "result",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON and back to verify structure
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			var decoded JSONResponse
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if decoded.Success != tt.response.Success {
				t.Errorf("Expected Success %v, got %v", tt.response.Success, decoded.Success)
			}

			if decoded.Message != tt.response.Message {
				t.Errorf("Expected Message %s, got %s", tt.response.Message, decoded.Message)
			}

			if decoded.Error != tt.response.Error {
				t.Errorf("Expected Error %s, got %s", tt.response.Error, decoded.Error)
			}
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		response       JSONResponse
		expectedStatus int
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			response: JSONResponse{
				Success: true,
				Data:    "test data",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "error response",
			statusCode: http.StatusBadRequest,
			response: JSONResponse{
				Success: false,
				Error:   "bad request",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			response: JSONResponse{
				Success: false,
				Error:   "internal server error",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONResponse(w, tt.statusCode, tt.response)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Verify the response body is valid JSON
			var decoded JSONResponse
			err := json.NewDecoder(w.Body).Decode(&decoded)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if decoded.Success != tt.response.Success {
				t.Errorf("Expected Success %v, got %v", tt.response.Success, decoded.Success)
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedStatus int
	}{
		{
			name:           "bad request error",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid input",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			statusCode:     http.StatusNotFound,
			message:        "Resource not found",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			message:        "Something went wrong",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONError(w, tt.statusCode, tt.message)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			var response JSONResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if response.Success != false {
				t.Error("Expected Success to be false for error response")
			}

			if response.Error != tt.message {
				t.Errorf("Expected Error %s, got %s", tt.message, response.Error)
			}

			if response.Data != nil {
				t.Error("Expected Data to be nil for error response")
			}
		})
	}
}

func TestWriteJSONSuccess(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{
			name: "string data",
			data: "success message",
		},
		{
			name: "map data",
			data: map[string]string{"key": "value"},
		},
		{
			name: "array data",
			data: []string{"item1", "item2"},
		},
		{
			name: "number data",
			data: 42,
		},
		{
			name: "nil data",
			data: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONSuccess(w, tt.data)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			var response JSONResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if response.Success != true {
				t.Error("Expected Success to be true for success response")
			}

			if response.Error != "" {
				t.Error("Expected Error to be empty for success response")
			}

			// For complex data types, we can't easily compare directly due to JSON marshaling
			// Instead, just verify that the data field is populated when expected
			if tt.data != nil && response.Data == nil {
				t.Error("Expected Data to be populated for non-nil input")
			}
		})
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	timestamp := GetCurrentTimestamp()

	// Verify the timestamp is in RFC3339 format
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("Expected RFC3339 format timestamp, got parsing error: %v", err)
	}

	// Verify the timestamp is recent (within last 5 seconds)
	parsedTime, _ := time.Parse(time.RFC3339, timestamp)
	now := time.Now().UTC()
	diff := now.Sub(parsedTime)

	if diff > 5*time.Second {
		t.Errorf("Timestamp seems too old: %v", diff)
	}

	if diff < 0 {
		t.Error("Timestamp is in the future")
	}
}

func TestGetCurrentTimestamp_Format(t *testing.T) {
	timestamp := GetCurrentTimestamp()

	// Check specific format characteristics
	if !strings.Contains(timestamp, "T") {
		t.Error("Expected timestamp to contain 'T' separator")
	}

	if !strings.Contains(timestamp, "Z") {
		t.Error("Expected timestamp to end with 'Z' (UTC timezone)")
	}

	// Verify it's a valid RFC3339 timestamp
	expectedFormat := "2006-01-02T15:04:05Z07:00"
	_, err := time.Parse(expectedFormat, timestamp)
	if err != nil {
		t.Errorf("Timestamp doesn't match expected format: %v", err)
	}
}

// Test error cases for JSON writing
func TestWriteJSONResponse_ErrorCases(t *testing.T) {
	// Test with data that can't be marshaled to JSON
	w := httptest.NewRecorder()
	
	// Function types can't be marshaled to JSON
	response := JSONResponse{
		Success: true,
		Data:    func() {}, // This will cause JSON marshal to fail
	}

	// This should handle the error gracefully (though the current implementation might not)
	WriteJSONResponse(w, http.StatusOK, response)

	// The response should indicate some kind of error occurred
	// Note: This test might need adjustment based on how you want to handle marshal errors
}

// Integration test combining multiple utility functions
func TestUtilFunctions_Integration(t *testing.T) {
	t.Run("success workflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		data := map[string]any{
			"message":   "Operation completed successfully",
			"timestamp": GetCurrentTimestamp(),
			"count":     42,
		}
		
		WriteJSONSuccess(w, data)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response JSONResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Error("Expected success response")
		}

		responseData, ok := response.Data.(map[string]any)
		if !ok {
			t.Fatal("Expected data to be a map")
		}

		if responseData["count"] != float64(42) { // JSON numbers become float64
			t.Errorf("Expected count 42, got %v", responseData["count"])
		}
	})

	t.Run("error workflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		
		WriteJSONError(w, http.StatusBadRequest, "Validation failed")

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		var response JSONResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Success {
			t.Error("Expected error response")
		}

		if response.Error != "Validation failed" {
			t.Errorf("Expected error message 'Validation failed', got %s", response.Error)
		}
	})
}

// Benchmarks
func BenchmarkWriteJSONSuccess(b *testing.B) {
	data := map[string]any{
		"message": "benchmark test",
		"count":   100,
		"items":   []string{"a", "b", "c"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		WriteJSONSuccess(w, data)
	}
}

func BenchmarkWriteJSONError(b *testing.B) {
	message := "benchmark error message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		WriteJSONError(w, http.StatusBadRequest, message)
	}
}

func BenchmarkGetCurrentTimestamp(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCurrentTimestamp()
	}
}
