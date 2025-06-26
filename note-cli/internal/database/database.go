package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

// Recording represents an audio recording in the database
type Recording struct {
	ID         int           `json:"id"`
	Filename   string        `json:"filename"`
	FilePath   string        `json:"file_path"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	FileSize   int64         `json:"file_size"`
	Format     string        `json:"format"`
	SampleRate int           `json:"sample_rate"`
	Channels   int           `json:"channels"`
	CreatedAt  time.Time     `json:"created_at"`
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
	createNotesTableSQL := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		tags TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createNotesTableSQL); err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	// Create recordings table
	createRecordingsTableSQL := `
	CREATE TABLE IF NOT EXISTS recordings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		duration INTEGER NOT NULL,
		file_size INTEGER NOT NULL,
		format TEXT NOT NULL,
		sample_rate INTEGER NOT NULL,
		channels INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createRecordingsTableSQL); err != nil {
		return fmt.Errorf("failed to create recordings table: %w", err)
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

// ensureRecordingsTable creates the recordings table if it doesn't exist
func ensureRecordingsTable(db *sql.DB) error {
	createRecordingsTableSQL := `
	CREATE TABLE IF NOT EXISTS recordings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		duration INTEGER NOT NULL,
		file_size INTEGER NOT NULL,
		format TEXT NOT NULL,
		sample_rate INTEGER NOT NULL,
		channels INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createRecordingsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create recordings table: %w", err)
	}
	return nil
}

// SaveRecording saves a new recording to the database
func SaveRecording(dbPath string, recording *Recording) error {
	db, err := Connect(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Ensure recordings table exists
	if err := ensureRecordingsTable(db); err != nil {
		return err
	}

	query := `
		INSERT INTO recordings (filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, 
		recording.Filename,
		recording.FilePath,
		recording.StartTime,
		recording.EndTime,
		int64(recording.Duration), // Store duration as nanoseconds
		recording.FileSize,
		recording.Format,
		recording.SampleRate,
		recording.Channels,
		recording.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save recording: %w", err)
	}

	return nil
}

// ListRecordings retrieves all recordings from the database
func ListRecordings(dbPath string) ([]Recording, error) {
	db, err := Connect(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Ensure recordings table exists
	if err := ensureRecordingsTable(db); err != nil {
		return nil, err
	}

	query := `SELECT id, filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels, created_at 
			  FROM recordings ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recordings: %w", err)
	}
	defer rows.Close()

	var recordings []Recording
	for rows.Next() {
		var recording Recording
		var durationNanos int64
		err := rows.Scan(
			&recording.ID,
			&recording.Filename,
			&recording.FilePath,
			&recording.StartTime,
			&recording.EndTime,
			&durationNanos,
			&recording.FileSize,
			&recording.Format,
			&recording.SampleRate,
			&recording.Channels,
			&recording.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recording: %w", err)
		}
		recording.Duration = time.Duration(durationNanos)
		recordings = append(recordings, recording)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recordings: %w", err)
	}

	return recordings, nil
}
