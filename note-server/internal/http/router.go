package http

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/your-org/note-server/internal/ws"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(transcribeHub *ws.TranscribeHub) http.Handler {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	
	// Initialize handlers with dependencies
	handlers := NewHandlers()
	
	// Health check endpoint
	r.Get("/healthz", handlers.HealthHandler)
	
	// API routes
	r.Post("/transcribe", handlers.TranscribeHandler)
	r.Post("/summarize", handlers.SummarizeHandler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/notes", NotesHandler)
		r.Post("/notes", NotesHandler)
		
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
