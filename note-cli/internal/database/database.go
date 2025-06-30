package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Note represents a general note in the database
type Note struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Summary     string  `json:"summary"`
	Tags        string  `json:"tags"`
	RecordingID *int    `json:"recording_id"`  // Optional reference to recording
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Meeting represents a meeting in the database
type Meeting struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Summary     string  `json:"summary"`
	Attendees   string  `json:"attendees"`
	Location    string  `json:"location"`
	Tags        string  `json:"tags"`
	RecordingID *int    `json:"recording_id"`  // Optional reference to recording
	MeetingDate *string `json:"meeting_date"`  // When the meeting actually occurred
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Interview represents an interview in the database
type Interview struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Summary        string  `json:"summary"`
	Interviewee    string  `json:"interviewee"`
	Interviewer    string  `json:"interviewer"`
	Company        string  `json:"company"`
	Position       string  `json:"position"`
	Tags           string  `json:"tags"`
	RecordingID    *int    `json:"recording_id"`  // Optional reference to recording
	InterviewDate  *string `json:"interview_date"`  // When the interview actually occurred
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// Recording represents an audio recording file in the database
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

	// Create recordings table first (referenced by other tables)
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

	// Create notes table
	createNotesTableSQL := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT DEFAULT '',
		tags TEXT DEFAULT '',
		recording_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recording_id) REFERENCES recordings(id) ON DELETE SET NULL
	);`

	if _, err := db.Exec(createNotesTableSQL); err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	// Create meetings table
	createMeetingsTableSQL := `
	CREATE TABLE IF NOT EXISTS meetings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT DEFAULT '',
		attendees TEXT DEFAULT '',
		location TEXT DEFAULT '',
		tags TEXT DEFAULT '',
		recording_id INTEGER,
		meeting_date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recording_id) REFERENCES recordings(id) ON DELETE SET NULL
	);`

	if _, err := db.Exec(createMeetingsTableSQL); err != nil {
		return fmt.Errorf("failed to create meetings table: %w", err)
	}

	// Create interviews table
	createInterviewsTableSQL := `
	CREATE TABLE IF NOT EXISTS interviews (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT DEFAULT '',
		interviewee TEXT DEFAULT '',
		interviewer TEXT DEFAULT '',
		company TEXT DEFAULT '',
		position TEXT DEFAULT '',
		tags TEXT DEFAULT '',
		recording_id INTEGER,
		interview_date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recording_id) REFERENCES recordings(id) ON DELETE SET NULL
	);`

	if _, err := db.Exec(createInterviewsTableSQL); err != nil {
		return fmt.Errorf("failed to create interviews table: %w", err)
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
func CreateNote(db *sql.DB, title, content, summary, tags string, recordingID *int) (*Note, error) {
	query := `
		INSERT INTO notes (title, content, summary, tags, recording_id) 
		VALUES (?, ?, ?, ?, ?) 
		RETURNING id, title, content, summary, tags, recording_id, created_at, updated_at`
	
	var note Note
	err := db.QueryRow(query, title, content, summary, tags, recordingID).Scan(
		&note.ID, &note.Title, &note.Content, &note.Summary, &note.Tags, &note.RecordingID,
		&note.CreatedAt, &note.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}
	
	return &note, nil
}

// CreateMeeting creates a new meeting in the database
func CreateMeeting(db *sql.DB, title, content, summary, attendees, location, tags string, recordingID *int, meetingDate *string) (*Meeting, error) {
	query := `
		INSERT INTO meetings (title, content, summary, attendees, location, tags, recording_id, meeting_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING id, title, content, summary, attendees, location, tags, recording_id, meeting_date, created_at, updated_at`
	
	var meeting Meeting
	err := db.QueryRow(query, title, content, summary, attendees, location, tags, recordingID, meetingDate).Scan(
		&meeting.ID, &meeting.Title, &meeting.Content, &meeting.Summary, &meeting.Attendees, &meeting.Location,
		&meeting.Tags, &meeting.RecordingID, &meeting.MeetingDate, &meeting.CreatedAt, &meeting.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}
	
	return &meeting, nil
}

// CreateInterview creates a new interview in the database
func CreateInterview(db *sql.DB, title, content, summary, interviewee, interviewer, company, position, tags string, recordingID *int, interviewDate *string) (*Interview, error) {
	query := `
		INSERT INTO interviews (title, content, summary, interviewee, interviewer, company, position, tags, recording_id, interview_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		RETURNING id, title, content, summary, interviewee, interviewer, company, position, tags, recording_id, interview_date, created_at, updated_at`
	
	var interview Interview
	err := db.QueryRow(query, title, content, summary, interviewee, interviewer, company, position, tags, recordingID, interviewDate).Scan(
		&interview.ID, &interview.Title, &interview.Content, &interview.Summary, &interview.Interviewee, &interview.Interviewer,
		&interview.Company, &interview.Position, &interview.Tags, &interview.RecordingID, &interview.InterviewDate,
		&interview.CreatedAt, &interview.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create interview: %w", err)
	}
	
	return &interview, nil
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

// DeleteRecording deletes a recording from the database by ID
func DeleteRecording(dbPath string, recordingID int) (*Recording, error) {
	db, err := Connect(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Ensure recordings table exists
	if err := ensureRecordingsTable(db); err != nil {
		return nil, err
	}

	// First, get the recording details before deletion
	var recording Recording
	var durationNanos int64
	querySelect := `SELECT id, filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels, created_at 
					FROM recordings WHERE id = ?`
	
	err = db.QueryRow(querySelect, recordingID).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recording with ID %d not found", recordingID)
		}
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}
	recording.Duration = time.Duration(durationNanos)

	// Delete the recording from database
	queryDelete := `DELETE FROM recordings WHERE id = ?`
	result, err := db.Exec(queryDelete, recordingID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete recording: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("recording with ID %d not found", recordingID)
	}

	return &recording, nil
}
