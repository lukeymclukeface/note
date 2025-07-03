package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/your-org/note-server/internal/config"
	"github.com/your-org/note-server/internal/service"
	"github.com/your-org/note-server/internal/util"
)

// Handlers holds the service dependencies
type Handlers struct {
	transcribeService *service.TranscribeService
	summarizeService  *service.SummarizeService
	configManager     *config.ConfigManager
}

// NewHandlers creates a new handlers instance
func NewHandlers() *Handlers {
	return &Handlers{
		transcribeService: service.NewTranscribeService(),
		summarizeService:  service.NewSummarizeService(),
		configManager:     config.GetManager(),
	}
}

// NewHandlersWithServices creates handlers with injected services for testing
func NewHandlersWithServices(transcribeService *service.TranscribeService, summarizeService *service.SummarizeService) *Handlers {
	return &Handlers{
		transcribeService: transcribeService,
		summarizeService:  summarizeService,
		configManager:     config.GetManager(),
	}
}

// HealthHandler responds with the server's health status
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// TranscribeHandler handles requests for transcription
func (h *Handlers) TranscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Create temporary file
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "audio_*"+filepath.Ext(header.Filename))
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy uploaded file to temporary file
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Read the file data
	tmpFile.Seek(0, 0)
	audioData, err := io.ReadAll(tmpFile)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, "Failed to read file data")
		return
	}

	// Record start time for duration calculation
	startTime := time.Now()

	// Call transcription service
	ctx := context.Background()
	text, err := h.transcribeService.TranscribeAudio(ctx, audioData)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Transcription failed: %v", err))
		return
	}

	// Calculate duration
	durationMs := time.Since(startTime).Milliseconds()

	// Prepare response
	response := map[string]any{
		"text":        text,
		"duration_ms": durationMs,
	}

	util.WriteJSONSuccess(w, response)
}

// SummarizeRequest represents the request body for summarization
type SummarizeRequest struct {
	Text string `json:"text"`
}

// SummarizeHandler handles requests for summarization
func (h *Handlers) SummarizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var req SummarizeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Text == "" {
		util.WriteJSONError(w, http.StatusBadRequest, "Text field is required")
		return
	}

	// Call summarization service
	ctx := context.Background()
	summary, err := h.summarizeService.SummarizeText(ctx, req.Text)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Summarization failed: %v", err))
		return
	}

	// Prepare response
	response := map[string]any{
		"summary": summary,
	}

	util.WriteJSONSuccess(w, response)
}

// NotesHandler handles requests for notes operations
func NotesHandler(w http.ResponseWriter, r *http.Request) {
	// Implement your logic here
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Implement your WebSocket logic here
}

// GetConfig handles GET /api/config requests
func (h *Handlers) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	appConfig := h.configManager.GetConfig()
	
	// Mask the OpenAI key for security
	if appConfig.OpenAIKey != "" {
		if len(appConfig.OpenAIKey) > 8 {
			appConfig.OpenAIKey = appConfig.OpenAIKey[:4] + "..." + appConfig.OpenAIKey[len(appConfig.OpenAIKey)-4:]
		} else {
			appConfig.OpenAIKey = "***"
		}
	}

	response := map[string]any{
		"success": true,
		"config":  appConfig,
	}

	util.WriteJSONSuccess(w, response)
}

// SetConfig handles PUT /api/config requests
func (h *Handlers) SetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newConfig config.AppConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Save the new configuration
	if err := h.configManager.SetConfig(newConfig); err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, "Failed to save configuration: "+err.Error())
		return
	}

	response := map[string]any{
		"success": true,
		"message": "Configuration saved successfully",
	}

	util.WriteJSONSuccess(w, response)
}

// GetConfigRaw handles GET /api/config/raw requests (returns unmasked config for editing)
func (h *Handlers) GetConfigRaw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	appConfig := h.configManager.GetConfig()

	response := map[string]any{
		"success": true,
		"config":  appConfig,
	}

	util.WriteJSONSuccess(w, response)
}
