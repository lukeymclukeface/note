package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   JSONResponse
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			response: JSONResponse{
				Success: true,
				Data:    "test data",
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"success":true,"data":"test data"}`,
		},
		{
			name:       "error response",
			statusCode: http.StatusBadRequest,
			response: JSONResponse{
				Success: false,
				Error:   "test error",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"success":false,"error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSONResponse(w, tt.statusCode, tt.response)

			if w.Code != tt.wantStatus {
				t.Errorf("WriteJSONResponse() status = %v, want %v", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSONResponse() Content-Type = %v, want %v", contentType, "application/json")
			}

			body := w.Body.String()
			// Parse JSON to ensure it's valid and compare
			var got, want map[string]any
			if err := json.Unmarshal([]byte(body), &got); err != nil {
				t.Fatalf("WriteJSONResponse() returned invalid JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.wantBody), &want); err != nil {
				t.Fatalf("Test case has invalid JSON: %v", err)
			}

			if len(got) != len(want) {
				t.Errorf("WriteJSONResponse() body = %v, want %v", body, tt.wantBody)
				return
			}

			for key, wantVal := range want {
				if gotVal, ok := got[key]; !ok || gotVal != wantVal {
					t.Errorf("WriteJSONResponse() body[%s] = %v, want %v", key, gotVal, wantVal)
				}
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSONError(w, http.StatusNotFound, "not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("WriteJSONError() status = %v, want %v", w.Code, http.StatusNotFound)
	}

	var response JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("WriteJSONError() returned invalid JSON: %v", err)
	}

	if response.Success != false {
		t.Errorf("WriteJSONError() success = %v, want %v", response.Success, false)
	}

	if response.Error != "not found" {
		t.Errorf("WriteJSONError() error = %v, want %v", response.Error, "not found")
	}
}

func TestWriteJSONSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]string{"key": "value"}
	WriteJSONSuccess(w, testData)

	if w.Code != http.StatusOK {
		t.Errorf("WriteJSONSuccess() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("WriteJSONSuccess() returned invalid JSON: %v", err)
	}

	if response.Success != true {
		t.Errorf("WriteJSONSuccess() success = %v, want %v", response.Success, true)
	}

	// Convert data back to map for comparison
	dataMap, ok := response.Data.(map[string]any)
	if !ok {
		t.Errorf("WriteJSONSuccess() data is not a map")
	}

	if dataMap["key"] != "value" {
		t.Errorf("WriteJSONSuccess() data = %v, want %v", dataMap, testData)
	}
}
