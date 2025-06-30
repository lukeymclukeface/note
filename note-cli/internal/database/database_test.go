package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a temporary SQLite database for testing
func setupTestDB(t *testing.T) (*sql.DB, string) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "note_cli_test")
	require.NoError(t, err)

	// Create database path
	dbPath := filepath.Join(tempDir, "test_notes.db")

	// Initialize the database
	err = Initialize(dbPath)
	require.NoError(t, err)

	// Connect to the database
	db, err := Connect(dbPath)
	require.NoError(t, err)

	// Cleanup function to remove the temporary database
	t.Cleanup(func() {
		db.Close()
		os.RemoveAll(tempDir)
	})

	return db, dbPath
}

func TestInitialize(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "note_cli_test_init")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test_init.db")

	err = Initialize(dbPath)
	require.NoError(t, err)

	// Check that the database file was created
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)

	// Connect and verify tables exist
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Check if tables exist
	tables := []string{"notes", "meetings", "interviews", "recordings"}
	for _, table := range tables {
		var name string
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name)
	}
}

func TestConnect(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "note_cli_test_connect")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test_connect.db")

	// Initialize first
	err = Initialize(dbPath)
	require.NoError(t, err)

	// Test successful connection
	db, err := Connect(dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Test that connection is valid
	err = db.Ping()
	assert.NoError(t, err)
}

func TestConnect_InvalidPath(t *testing.T) {
	// Test connection to non-existent path
	db, err := Connect("/invalid/path/to/database.db")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateNote(t *testing.T) {
	db, _ := setupTestDB(t)

	// Test creating a basic note
	note, err := CreateNote(db, "Test Title", "Test Content", "Test Summary", "test,sample", nil)
	require.NoError(t, err)
	require.NotNil(t, note)

	assert.Greater(t, note.ID, 0)
	assert.Equal(t, "Test Title", note.Title)
	assert.Equal(t, "Test Content", note.Content)
	assert.Equal(t, "Test Summary", note.Summary)
	assert.Equal(t, "test,sample", note.Tags)
	assert.Nil(t, note.RecordingID)
	assert.NotEmpty(t, note.CreatedAt)
	assert.NotEmpty(t, note.UpdatedAt)
}

func TestCreateNote_WithRecordingID(t *testing.T) {
	db, _ := setupTestDB(t)

	recordingID := 123
	note, err := CreateNote(db, "Audio Note", "Audio content", "Audio summary", "audio,test", &recordingID)
	require.NoError(t, err)
	require.NotNil(t, note)

	assert.Equal(t, "Audio Note", note.Title)
	assert.NotNil(t, note.RecordingID)
	assert.Equal(t, 123, *note.RecordingID)
}

func TestCreateNote_EmptyFields(t *testing.T) {
	db, _ := setupTestDB(t)

	// Test with empty optional fields
	note, err := CreateNote(db, "Empty Fields", "Some content", "", "", nil)
	require.NoError(t, err)
	require.NotNil(t, note)

	assert.Equal(t, "Empty Fields", note.Title)
	assert.Equal(t, "Some content", note.Content)
	assert.Equal(t, "", note.Summary)
	assert.Equal(t, "", note.Tags)
}

func TestListNotes_Empty(t *testing.T) {
	db, _ := setupTestDB(t)

	notes, err := ListNotes(db, "")
	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestListNotes_WithData(t *testing.T) {
	db, _ := setupTestDB(t)

	// Create test notes
	_, err := CreateNote(db, "Note 1", "Content 1", "Summary 1", "tag1,common", nil)
	require.NoError(t, err)
	
	_, err = CreateNote(db, "Note 2", "Content 2", "Summary 2", "tag2,common", nil)
	require.NoError(t, err)

	// List all notes
	notes, err := ListNotes(db, "")
	require.NoError(t, err)
	assert.Len(t, notes, 2)

	// Check notes are ordered by creation (newest first typically)
	assert.Equal(t, "Note 1", notes[0].Title)
	assert.Equal(t, "Note 2", notes[1].Title)
}

func TestListNotes_FilterByTag(t *testing.T) {
	db, _ := setupTestDB(t)

	// Create test notes with different tags
	_, err := CreateNote(db, "Work Note", "Work content", "", "work,important", nil)
	require.NoError(t, err)
	
	_, err = CreateNote(db, "Personal Note", "Personal content", "", "personal,life", nil)
	require.NoError(t, err)
	
	_, err = CreateNote(db, "Mixed Note", "Mixed content", "", "work,personal", nil)
	require.NoError(t, err)

	// Test filtering by work tag
	workNotes, err := ListNotes(db, "work")
	require.NoError(t, err)
	assert.Len(t, workNotes, 2) // Work Note and Mixed Note

	// Test filtering by personal tag
	personalNotes, err := ListNotes(db, "personal")
	require.NoError(t, err)
	assert.Len(t, personalNotes, 2) // Personal Note and Mixed Note

	// Test filtering by non-existent tag
	nonExistentNotes, err := ListNotes(db, "nonexistent")
	require.NoError(t, err)
	assert.Empty(t, nonExistentNotes)
}

// Test the Note struct
func TestNoteStruct(t *testing.T) {
	note := Note{
		ID:          1,
		Title:       "Test Note",
		Content:     "Test Content",
		Summary:     "Test Summary",
		Tags:        "test,sample",
		RecordingID: nil,
		CreatedAt:   "2024-01-01T10:00:00Z",
		UpdatedAt:   "2024-01-01T10:00:00Z",
	}

	assert.Equal(t, 1, note.ID)
	assert.Equal(t, "Test Note", note.Title)
	assert.Equal(t, "Test Content", note.Content)
	assert.Equal(t, "Test Summary", note.Summary)
	assert.Equal(t, "test,sample", note.Tags)
	assert.Nil(t, note.RecordingID)
	assert.Equal(t, "2024-01-01T10:00:00Z", note.CreatedAt)
	assert.Equal(t, "2024-01-01T10:00:00Z", note.UpdatedAt)
}

// Test the Meeting struct
func TestMeetingStruct(t *testing.T) {
	meetingDate := "2024-01-01"
	meeting := Meeting{
		ID:          1,
		Title:       "Team Meeting",
		Content:     "Meeting content",
		Summary:     "Meeting summary",
		Attendees:   "John, Jane, Bob",
		Location:    "Conference Room A",
		Tags:        "meeting,team",
		RecordingID: nil,
		MeetingDate: &meetingDate,
		CreatedAt:   "2024-01-01T10:00:00Z",
		UpdatedAt:   "2024-01-01T10:00:00Z",
	}

	assert.Equal(t, 1, meeting.ID)
	assert.Equal(t, "Team Meeting", meeting.Title)
	assert.Equal(t, "Meeting content", meeting.Content)
	assert.Equal(t, "John, Jane, Bob", meeting.Attendees)
	assert.Equal(t, "Conference Room A", meeting.Location)
	assert.NotNil(t, meeting.MeetingDate)
	assert.Equal(t, "2024-01-01", *meeting.MeetingDate)
}

// Test the Interview struct  
func TestInterviewStruct(t *testing.T) {
	interviewDate := "2024-01-01"
	interview := Interview{
		ID:            1,
		Title:         "Software Engineer Interview",
		Content:       "Interview content",
		Summary:       "Interview summary",
		Interviewee:   "John Doe",
		Interviewer:   "Jane Smith",
		Company:       "Tech Corp",
		Position:      "Senior Developer",
		Tags:          "interview,tech",
		RecordingID:   nil,
		InterviewDate: &interviewDate,
		CreatedAt:     "2024-01-01T10:00:00Z",
		UpdatedAt:     "2024-01-01T10:00:00Z",
	}

	assert.Equal(t, 1, interview.ID)
	assert.Equal(t, "Software Engineer Interview", interview.Title)
	assert.Equal(t, "John Doe", interview.Interviewee)
	assert.Equal(t, "Jane Smith", interview.Interviewer)
	assert.Equal(t, "Tech Corp", interview.Company)
	assert.Equal(t, "Senior Developer", interview.Position)
	assert.NotNil(t, interview.InterviewDate)
	assert.Equal(t, "2024-01-01", *interview.InterviewDate)
}

// Test the Recording struct
func TestRecordingStruct(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour)
	duration := endTime.Sub(startTime)

	recording := Recording{
		ID:         1,
		Filename:   "meeting.wav",
		FilePath:   "/path/to/meeting.wav",
		StartTime:  startTime,
		EndTime:    endTime,
		Duration:   duration,
		FileSize:   1024000,
		Format:     "WAV",
		SampleRate: 44100,
		Channels:   2,
		CreatedAt:  startTime,
	}

	assert.Equal(t, 1, recording.ID)
	assert.Equal(t, "meeting.wav", recording.Filename)
	assert.Equal(t, "/path/to/meeting.wav", recording.FilePath)
	assert.Equal(t, startTime, recording.StartTime)
	assert.Equal(t, endTime, recording.EndTime)
	assert.Equal(t, duration, recording.Duration)
	assert.Equal(t, int64(1024000), recording.FileSize)
	assert.Equal(t, "WAV", recording.Format)
	assert.Equal(t, 44100, recording.SampleRate)
	assert.Equal(t, 2, recording.Channels)
}

// Benchmark test for database operations
func BenchmarkCreateNote(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "note_cli_benchmark")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "benchmark.db")
	err = Initialize(dbPath)
	require.NoError(b, err)

	db, err := Connect(dbPath)
	require.NoError(b, err)
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CreateNote(db, "Benchmark Note", "Benchmark Content", "Benchmark Summary", "benchmark", nil)
		require.NoError(b, err)
	}
}

func BenchmarkListNotes(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "note_cli_benchmark")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "benchmark.db")
	err = Initialize(dbPath)
	require.NoError(b, err)

	db, err := Connect(dbPath)
	require.NoError(b, err)
	defer db.Close()

	// Create some test data
	for i := 0; i < 100; i++ {
		_, err := CreateNote(db, "Test Note", "Test Content", "Test Summary", "test", nil)
		require.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ListNotes(db, "")
		require.NoError(b, err)
	}
}
