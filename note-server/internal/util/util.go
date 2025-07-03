package util

import (
	"github.com/your-org/note-server/pkg/response"
	"github.com/your-org/note-server/pkg/timeutil"
	"net/http"
)

// Re-export shared utilities for backward compatibility
// These will be removed in a future version, use pkg/ directly

// JSONResponse is deprecated, use pkg/response.JSONResponse
type JSONResponse = response.JSONResponse

// WriteJSONResponse is deprecated, use pkg/response.WriteJSONResponse
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp JSONResponse) {
	response.WriteJSONResponse(w, statusCode, resp)
}

// WriteJSONError is deprecated, use pkg/response.WriteJSONError
func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	response.WriteJSONError(w, statusCode, message)
}

// WriteJSONSuccess is deprecated, use pkg/response.WriteJSONSuccess
func WriteJSONSuccess(w http.ResponseWriter, data any) {
	response.WriteJSONSuccess(w, data)
}

// GetCurrentTimestamp is deprecated, use pkg/timeutil.GetCurrentTimestamp
func GetCurrentTimestamp() string {
	return timeutil.GetCurrentTimestamp()
}
