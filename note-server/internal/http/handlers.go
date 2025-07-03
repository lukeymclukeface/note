package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/your-org/note-server/internal/config"
	"github.com/your-org/note-server/internal/database"
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

// GetNotes handles GET /api/notes requests
func (h *Handlers) GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// TODO: Implement database query for notes
	response := map[string]any{
		"success": true,
		"notes":   []any{}, // Empty for now
	}

	util.WriteJSONSuccess(w, response)
}

// CreateNote handles POST /api/notes requests
func (h *Handlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// TODO: Implement note creation
	response := map[string]any{
		"success": true,
		"message": "Note creation not yet implemented",
	}

	util.WriteJSONSuccess(w, response)
}

// GetMeetings handles GET /api/meetings requests
func (h *Handlers) GetMeetings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// TODO: Implement database query for meetings
	response := map[string]any{
		"success":  true,
		"meetings": []any{}, // Empty for now
	}

	util.WriteJSONSuccess(w, response)
}

// GetMeeting handles GET /api/meetings/{id} requests
func (h *Handlers) GetMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// TODO: Extract ID from URL and query database
	util.WriteJSONError(w, http.StatusNotFound, "Meeting not found")
}

// GetInterviews handles GET /api/interviews requests
func (h *Handlers) GetInterviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// TODO: Implement database query for interviews
	response := map[string]any{
		"success":    true,
		"interviews": []any{}, // Empty for now
	}

	util.WriteJSONSuccess(w, response)
}

// GetRecordings handles GET /api/recordings requests
func (h *Handlers) GetRecordings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Query recordings from database
	recordings, err := database.GetRecordings()
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get recordings: %v", err))
		return
	}

	response := map[string]any{
		"success":    true,
		"recordings": recordings,
	}

	// Write response directly without extra wrapping
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRecording handles GET /api/recordings/{id} requests
func (h *Handlers) GetRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := ""
	if len(path) > len("/api/recordings/") {
		idStr = path[len("/api/recordings/"):]
		// Remove any trailing path segments (like /audio)
		if slashIndex := strings.Index(idStr, "/"); slashIndex != -1 {
			idStr = idStr[:slashIndex]
		}
	}

	if idStr == "" {
		util.WriteJSONError(w, http.StatusBadRequest, "Recording ID is required")
		return
	}

	// Parse ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Invalid recording ID")
		return
	}

	// Query recording from database
	recording, err := database.GetRecording(id)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get recording: %v", err))
		return
	}

	if recording == nil {
		util.WriteJSONError(w, http.StatusNotFound, "Recording not found")
		return
	}

	response := map[string]any{
		"success":   true,
		"recording": recording,
	}

	// Write response directly without extra wrapping
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRecordingAudio handles GET /api/recordings/{id}/audio requests
func (h *Handlers) GetRecordingAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := ""
	if len(path) > len("/api/recordings/") {
		idStr = path[len("/api/recordings/"):]
		// Remove the /audio suffix
		if strings.HasSuffix(idStr, "/audio") {
			idStr = idStr[:len(idStr)-6]
		}
	}

	if idStr == "" {
		util.WriteJSONError(w, http.StatusBadRequest, "Recording ID is required")
		return
	}

	// Parse ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Invalid recording ID")
		return
	}

	// Query recording from database to get file path
	recording, err := database.GetRecording(id)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get recording: %v", err))
		return
	}

	if recording == nil {
		util.WriteJSONError(w, http.StatusNotFound, "Recording not found")
		return
	}

	// For now, since we don't have actual file storage, return a mock response
	// In a real implementation, you would:
	// 1. Get the file_path from the recording
	// 2. Open and stream the audio file
	// 3. Set appropriate headers (Content-Type, Content-Length, etc.)
	
	filePath, ok := recording["file_path"].(string)
	if !ok {
		util.WriteJSONError(w, http.StatusInternalServerError, "Invalid file path")
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		util.WriteJSONError(w, http.StatusNotFound, "Audio file not found")
		return
	}

	// Set proper content type based on file extension
	filename, ok := recording["filename"].(string)
	if !ok {
		filename = "audio.webm"
	}

	// Determine content type from file extension
	contentType := "audio/webm" // default
	ext := filepath.Ext(filename)
	switch ext {
	case ".webm":
		contentType = "audio/webm"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	case ".ogg":
		contentType = "audio/ogg"
	case ".m4a":
		contentType = "audio/mp4"
	}

	// For development: Check if file contains actual audio data or just text
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, "Could not read file info")
		return
	}

	// If file is very small, it's likely a test file with text content
	if fileInfo.Size() < 1024 { // Less than 1KB
		// Return a minimal valid audio response to prevent player errors
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
		// Don't write any content - the audio player will handle the empty stream gracefully
		return
	}

	// Set headers before serving file
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	// Add CORS headers for audio playback
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range")

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// UploadRecording handles POST /api/upload-recording requests
func (h *Handlers) UploadRecording(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse multipart form (32MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	// Get the audio file from the form
	file, header, err := r.FormFile("audio")
	if err != nil {
		util.WriteJSONError(w, http.StatusBadRequest, "No audio file provided")
		return
	}
	defer file.Close()

	// Get optional timestamps
	startTimeStr := r.FormValue("startTime")
	endTimeStr := r.FormValue("endTime")

	// Use provided times or fallback to current time
	var startTime, endTime time.Time
	if startTimeStr != "" {
		if parsed, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			startTime = parsed
		} else {
			startTime = time.Now()
		}
	} else {
		startTime = time.Now()
	}

	if endTimeStr != "" {
		if parsed, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			endTime = parsed
		} else {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("recording_%s.webm", timestamp)

	// Calculate duration
	durationMs := endTime.Sub(startTime).Milliseconds()
	if durationMs < 0 {
		durationMs = 1000 // Default 1 second if invalid times
	}
	durationSeconds := int(durationMs / 1000)

	// Create recordings directory if it doesn't exist
	recordingsDir := "/tmp/recordings"
	if err := os.MkdirAll(recordingsDir, 0755); err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create recordings directory: %v", err))
		return
	}

	// Create file path
	filePath := filepath.Join(recordingsDir, filename)

	// Save the uploaded file to disk
	outFile, err := os.Create(filePath)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create file: %v", err))
		return
	}
	defer outFile.Close()

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Copy the uploaded file content to disk
	_, err = io.Copy(outFile, file)
	if err != nil {
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save file: %v", err))
		return
	}

	// Save recording metadata to database
	recordingID, err := database.AddRecording(
		filename,
		filePath,
		startTime,
		endTime,
		durationSeconds,
		int(header.Size),
		"webm", // format
		44100,  // sample rate (default)
		2,      // channels (default)
	)
	if err != nil {
		// Clean up the file if database save fails
		os.Remove(filePath)
		util.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save recording metadata: %v", err))
		return
	}

	response := map[string]any{
		"success":     true,
		"filename":    filename,
		"recordingId": recordingID,
		"size":        header.Size,
		"duration":    durationMs,
		"message":     "Recording metadata saved to database",
	}

	util.WriteJSONSuccess(w, response)
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
