package response

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard JSON response
type JSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// WriteJSONResponse writes a JSON response to the HTTP response writer
func WriteJSONResponse(w http.ResponseWriter, statusCode int, response JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSONResponse(w, statusCode, JSONResponse{
		Success: false,
		Error:   message,
	})
}

// WriteJSONSuccess writes a JSON success response
func WriteJSONSuccess(w http.ResponseWriter, data any) {
	WriteJSONResponse(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    data,
	})
}
