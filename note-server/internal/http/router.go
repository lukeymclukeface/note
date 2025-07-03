package http

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/your-org/note-server/internal/ws"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(transcribeHub *ws.TranscribeHub) http.Handler {
	return NewRouterWithHandlers(transcribeHub, NewHandlers())
}

// NewRouterWithHandlers creates a new HTTP router with injected handlers for testing
func NewRouterWithHandlers(transcribeHub *ws.TranscribeHub, handlers *Handlers) http.Handler {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	
	// Health check endpoint
	r.Get("/healthz", handlers.HealthHandler)
	
	// API routes
	r.Post("/transcribe", handlers.TranscribeHandler)
	r.Post("/summarize", handlers.SummarizeHandler)
	r.Route("/api", func(r chi.Router) {
		// Notes endpoints
		r.Get("/notes", handlers.GetNotes)
		r.Post("/notes", handlers.CreateNote)
		
		// Meetings endpoints
		r.Get("/meetings", handlers.GetMeetings)
		r.Get("/meetings/{id}", handlers.GetMeeting)
		
		// Interviews endpoints
		r.Get("/interviews", handlers.GetInterviews)
		
		// Recordings endpoints
		r.Get("/recordings", handlers.GetRecordings)
		r.Get("/recordings/{id}", handlers.GetRecording)
		r.Get("/recordings/{id}/audio", handlers.GetRecordingAudio)
		r.Post("/upload-recording", handlers.UploadRecording)
		
		// Configuration endpoints
		r.Get("/config", handlers.GetConfig)
		r.Put("/config", handlers.SetConfig)
		r.Get("/config/raw", handlers.GetConfigRaw)
	})
	
	// WebSocket endpoints
	r.Get("/ws", WebSocketHandler)
	r.Get("/ws/transcribe", transcribeHub.ServeTranscribeWS)
	
	return r
}
