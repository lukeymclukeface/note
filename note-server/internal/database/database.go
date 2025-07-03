package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS recordings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		file_path TEXT,
		start_time DATETIME,
		end_time DATETIME,
		duration INTEGER,
		file_size INTEGER,
		format TEXT,
		sample_rate INTEGER,
		channels INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

// GetRecordings retrieves all recordings from the database
func GetRecordings() ([]map[string]any, error) {
	rows, err := db.Query("SELECT id, filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels, created_at FROM recordings")
	if err != nil {
		return nil, fmt.Errorf("failed to query recordings: %v", err)
	}
	defer rows.Close()

	var recordings []map[string]any
	for rows.Next() {
		var id int
		var filename, file_path, start_time, end_time, format string
		var duration, file_size, sample_rate, channels int
		var created_at string
		if err := rows.Scan(&id, &filename, &file_path, &start_time, &end_time, &duration, &file_size, &format, &sample_rate, &channels, &created_at); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		recordings = append(recordings, map[string]any{
			"id":          id,
			"filename":    filename,
			"file_path":   file_path,
			"start_time":  start_time,
			"end_time":    end_time,
			"duration":    duration,
			"file_size":   file_size,
			"format":      format,
			"sample_rate": sample_rate,
			"channels":    channels,
			"created_at":  created_at,
		})
	}

	return recordings, nil
}

// AddRecording inserts a new recording into the database
func AddRecording(filename, filePath string, startTime, endTime time.Time, duration, fileSize int, format string, sampleRate, channels int) (int64, error) {
	stmt, err := db.Prepare(`INSERT INTO recordings (filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(filename, filePath, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), duration, fileSize, format, sampleRate, channels)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

// GetRecording retrieves a specific recording by ID from the database
func GetRecording(id int) (map[string]any, error) {
	row := db.QueryRow("SELECT id, filename, file_path, start_time, end_time, duration, file_size, format, sample_rate, channels, created_at FROM recordings WHERE id = ?", id)

	var recordingID int
	var filename, file_path, start_time, end_time, format string
	var duration, file_size, sample_rate, channels int
	var created_at string

	err := row.Scan(&recordingID, &filename, &file_path, &start_time, &end_time, &duration, &file_size, &format, &sample_rate, &channels, &created_at)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // Recording not found
		}
		return nil, fmt.Errorf("failed to scan recording: %v", err)
	}

	recording := map[string]any{
		"id":          recordingID,
		"filename":    filename,
		"file_path":   file_path,
		"start_time":  start_time,
		"end_time":    end_time,
		"duration":    duration,
		"file_size":   file_size,
		"format":      format,
		"sample_rate": sample_rate,
		"channels":    channels,
		"created_at":  created_at,
	}

	return recording, nil
}
