package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Note represents a note in the database
type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Initialize creates and sets up the database
func Initialize(dbPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create notes table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		tags TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	return nil
}

// Connect establishes a connection to the database
func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

// CreateNote creates a new note in the database
func CreateNote(db *sql.DB, title, content, tags string) (*Note, error) {
	query := `
		INSERT INTO notes (title, content, tags) 
		VALUES (?, ?, ?) 
		RETURNING id, title, content, tags, created_at, updated_at`
	
	var note Note
	err := db.QueryRow(query, title, content, tags).Scan(
		&note.ID, &note.Title, &note.Content, &note.Tags, 
		&note.CreatedAt, &note.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}
	
	return &note, nil
}

// ListNotes retrieves notes from the database with optional tag filtering
func ListNotes(db *sql.DB, tag string) ([]Note, error) {
	var query string
	var args []interface{}
	
	if tag != "" {
		query = `SELECT id, title, content, tags, created_at, updated_at 
				 FROM notes WHERE tags LIKE ? ORDER BY created_at DESC`
		args = []interface{}{"%" + tag + "%"}
	} else {
		query = `SELECT id, title, content, tags, created_at, updated_at 
				 FROM notes ORDER BY created_at DESC`
	}
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()
	
	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Tags, 
						 &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}
	
	return notes, nil
}
