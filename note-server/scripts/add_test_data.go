package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/your-org/note-server/internal/database"
)

func main() {
	// Initialize database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	dbPath := filepath.Join(homeDir, ".noteai", "notes.db")
	
	// Ensure .noteai directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}
	
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Printf("Database initialized at: %s", dbPath)
// Add sample recordings
	if _, err := database.AddRecording("test1.wav", "/path/to/test1.wav", time.Now(), time.Now().Add(1*time.Hour), 3600, 1024, "wav", 44100, 2); err != nil {
		log.Fatalf("Failed to add recording: %v", err)
	}
	if _, err := database.AddRecording("test2.wav", "/path/to/test2.wav", time.Now(), time.Now().Add(1*time.Hour), 3600, 2048, "wav", 44100, 2); err != nil {
		log.Fatalf("Failed to add recording: %v", err)
	}
	log.Println("Sample recordings added!")
}
