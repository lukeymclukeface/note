// Package ws provides WebSocket functionality for real-time audio transcription.
//
// WebSocket Transcription Endpoint: /ws/transcribe
//
// Usage:
//   1. Connect to ws://localhost:8080/ws/transcribe
//   2. Send binary audio chunks (opus/pcm format) via WebSocket binary messages
//   3. Receive JSON messages with transcription results:
//      - { "type": "partial", "text": "..." } - Incremental transcription
//      - { "type": "final", "text": "..." } - Final transcription result
//      - { "type": "error", "text": "..." } - Error message
//
// Features:
//   - Concurrent connection limit (100 max)
//   - Context-based cancellation
//   - Goroutine-based message processing
//   - Automatic ping/pong for connection health
//   - Graceful shutdown support
package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/your-org/note-server/internal/service"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 1024 * 1024 // 1MB for audio chunks

	// Maximum number of concurrent connections
	maxConnections = 100
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// TranscribeMessage represents messages sent to/from the transcription WebSocket
type TranscribeMessage struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// TranscribeHub manages WebSocket connections for audio transcription
type TranscribeHub struct {
	// Registered clients
	clients map[*TranscribeClient]bool
	mutex   sync.RWMutex

	// Register requests from clients
	register chan *TranscribeClient

	// Unregister requests from clients
	unregister chan *TranscribeClient

	// Transcription service
	transcribeService *service.TranscribeService

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// TranscribeClient represents a WebSocket client for transcription
type TranscribeClient struct {
	hub    *TranscribeHub
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTranscribeHub creates a new transcription WebSocket hub
func NewTranscribeHub(transcribeService *service.TranscribeService) *TranscribeHub {
	ctx, cancel := context.WithCancel(context.Background())
	return &TranscribeHub{
		clients:           make(map[*TranscribeClient]bool),
		register:          make(chan *TranscribeClient),
		unregister:        make(chan *TranscribeClient),
		transcribeService: transcribeService,
		ctx:               ctx,
		cancel:            cancel,
	}
}

// Run starts the transcription hub
func (h *TranscribeHub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.register:
			h.mutex.Lock()
			if len(h.clients) >= maxConnections {
				h.mutex.Unlock()
				client.conn.Close()
				log.Printf("Connection rejected: maximum connections (%d) reached", maxConnections)
				continue
			}
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Transcription client connected. Active connections: %d", len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				client.cancel()
				log.Printf("Transcription client disconnected. Active connections: %d", len(h.clients))
			}
			h.mutex.Unlock()
		}
	}
}

// Shutdown gracefully shuts down the hub
func (h *TranscribeHub) Shutdown() {
	h.cancel()
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for client := range h.clients {
		client.conn.Close()
	}
}

// ServeTranscribeWS handles WebSocket connection requests for transcription
func (h *TranscribeHub) ServeTranscribeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(h.ctx)
	client := &TranscribeClient{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		ctx:    ctx,
		cancel: cancel,
	}

	h.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// readPump handles incoming audio data from the client
func (c *TranscribeClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		messageType, audioData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Only process binary messages (audio data)
		if messageType != websocket.BinaryMessage {
			continue
		}

		// Process audio chunk in a separate goroutine to avoid blocking
		go c.processAudioChunk(audioData)
	}
}

// processAudioChunk processes an audio chunk and sends transcription results
func (c *TranscribeClient) processAudioChunk(audioData []byte) {
	if len(audioData) == 0 {
		return
	}

	// Create a context with timeout for the transcription
	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()

	// Call the transcription service (currently a stub)
	// In a real implementation, this would stream to a transcription service
	// and receive partial and final results
	go func() {
		// Simulate partial result after a short delay
		time.Sleep(100 * time.Millisecond)
		partialMsg := TranscribeMessage{
			Type: "partial",
			Text: "Processing audio...",
		}
		c.sendTranscribeMessage(partialMsg)

		// Call the transcription service
		text, err := c.hub.transcribeService.TranscribeAudio(ctx, audioData)
		if err != nil {
			log.Printf("Transcription error: %v", err)
			errorMsg := TranscribeMessage{
				Type: "error",
				Text: "Transcription failed",
			}
			c.sendTranscribeMessage(errorMsg)
			return
		}

		// Send final result
		finalMsg := TranscribeMessage{
			Type: "final",
			Text: text,
		}
		c.sendTranscribeMessage(finalMsg)
	}()
}

// sendTranscribeMessage sends a transcription message to the client
func (c *TranscribeClient) sendTranscribeMessage(msg TranscribeMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal transcription message: %v", err)
		return
	}

	select {
	case c.send <- data:
	case <-c.ctx.Done():
	default:
		// Channel is full, close the connection
		close(c.send)
		c.hub.unregister <- c
	}
}

// writePump handles outgoing messages to the client
func (c *TranscribeClient) writePump() {
	defer c.conn.Close()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return

		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WriteMessage error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
