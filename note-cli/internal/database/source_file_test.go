package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestSourceFileDB(t *testing.T) (*sql.DB, func()) {
	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	
	// Initialize the database
	if err := Initialize(dbPath); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	
	// Connect to the database
	db, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}
	
	return db, cleanup
}

func TestCreateSourceFile(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Test data
	filePath := "/test/path/file.mp3"
	fileHash := "abc123def456"
	fileSize := int64(1024000)
	fileType := "audio"
	mimeType := "audio/mpeg"
	metadata := `{"file_name": "file.mp3", "duration": 120.5}`
	
	// Create source file
	sourceFile, err := CreateSourceFile(db, filePath, fileHash, fileSize, fileType, mimeType, metadata)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Verify the result
	if sourceFile.ID == 0 {
		t.Error("Expected non-zero ID")
	}
	if sourceFile.FilePath != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, sourceFile.FilePath)
	}
	if sourceFile.FileHash != fileHash {
		t.Errorf("Expected file hash %s, got %s", fileHash, sourceFile.FileHash)
	}
	if sourceFile.FileSize != fileSize {
		t.Errorf("Expected file size %d, got %d", fileSize, sourceFile.FileSize)
	}
	if sourceFile.FileType != fileType {
		t.Errorf("Expected file type %s, got %s", fileType, sourceFile.FileType)
	}
	if sourceFile.MimeType != mimeType {
		t.Errorf("Expected mime type %s, got %s", mimeType, sourceFile.MimeType)
	}
	if sourceFile.Metadata != metadata {
		t.Errorf("Expected metadata %s, got %s", metadata, sourceFile.Metadata)
	}
	if sourceFile.ProcessingStatus != "pending" {
		t.Errorf("Expected processing status 'pending', got %s", sourceFile.ProcessingStatus)
	}
	if sourceFile.ConvertedPath != nil {
		t.Error("Expected converted path to be nil")
	}
	if sourceFile.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
	if sourceFile.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}
}

func TestCreateSourceFileDuplicate(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Test data
	filePath := "/test/path/file.mp3"
	fileHash := "abc123def456"
	fileSize := int64(1024000)
	fileType := "audio"
	mimeType := "audio/mpeg"
	metadata := `{"file_name": "file.mp3"}`
	
	// Create first source file
	_, err := CreateSourceFile(db, filePath, fileHash, fileSize, fileType, mimeType, metadata)
	if err != nil {
		t.Fatalf("Failed to create first source file: %v", err)
	}
	
	// Try to create duplicate (should fail due to unique constraint on file_path)
	_, err = CreateSourceFile(db, filePath, fileHash+"different", fileSize, fileType, mimeType, metadata)
	if err == nil {
		t.Error("Expected error when creating duplicate file path, got nil")
	}
}

func TestGetSourceFileByHash(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Test data
	filePath := "/test/path/file.mp3"
	fileHash := "abc123def456"
	fileSize := int64(1024000)
	fileType := "audio"
	mimeType := "audio/mpeg"
	metadata := `{"file_name": "file.mp3"}`
	
	// Create source file
	original, err := CreateSourceFile(db, filePath, fileHash, fileSize, fileType, mimeType, metadata)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Retrieve by hash
	retrieved, err := GetSourceFileByHash(db, fileHash)
	if err != nil {
		t.Fatalf("Failed to get source file by hash: %v", err)
	}
	
	if retrieved == nil {
		t.Fatal("Expected source file, got nil")
	}
	
	// Verify the result
	if retrieved.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, retrieved.ID)
	}
	if retrieved.FileHash != fileHash {
		t.Errorf("Expected file hash %s, got %s", fileHash, retrieved.FileHash)
	}
}

func TestGetSourceFileByHashNotFound(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Try to get non-existent file
	retrieved, err := GetSourceFileByHash(db, "nonexistent")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if retrieved != nil {
		t.Error("Expected nil for non-existent file, got source file")
	}
}

func TestGetSourceFileByPath(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Test data
	filePath := "/test/path/file.mp3"
	fileHash := "abc123def456"
	fileSize := int64(1024000)
	fileType := "audio"
	mimeType := "audio/mpeg"
	metadata := `{"file_name": "file.mp3"}`
	
	// Create source file
	original, err := CreateSourceFile(db, filePath, fileHash, fileSize, fileType, mimeType, metadata)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Retrieve by path
	retrieved, err := GetSourceFileByPath(db, filePath)
	if err != nil {
		t.Fatalf("Failed to get source file by path: %v", err)
	}
	
	if retrieved == nil {
		t.Fatal("Expected source file, got nil")
	}
	
	// Verify the result
	if retrieved.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, retrieved.ID)
	}
	if retrieved.FilePath != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, retrieved.FilePath)
	}
}

func TestUpdateSourceFileStatus(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Create source file
	sourceFile, err := CreateSourceFile(db, "/test/file.mp3", "hash123", 1024, "audio", "audio/mpeg", "{}")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Update status
	newStatus := "completed"
	err = UpdateSourceFileStatus(db, sourceFile.ID, newStatus)
	if err != nil {
		t.Fatalf("Failed to update source file status: %v", err)
	}
	
	// Verify update
	updated, err := GetSourceFileByHash(db, "hash123")
	if err != nil {
		t.Fatalf("Failed to get updated source file: %v", err)
	}
	
	if updated.ProcessingStatus != newStatus {
		t.Errorf("Expected status %s, got %s", newStatus, updated.ProcessingStatus)
	}
	
	// Check that updated_at changed (allow for same timestamp due to SQLite precision)
	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Error("Expected updated_at to be at least equal to created_at")
	}
}

func TestUpdateSourceFileStatusNotFound(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Try to update non-existent file
	err := UpdateSourceFileStatus(db, 999, "completed")
	if err == nil {
		t.Error("Expected error when updating non-existent file, got nil")
	}
}

func TestUpdateSourceFileConvertedPath(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Create source file
	sourceFile, err := CreateSourceFile(db, "/test/file.mp3", "hash123", 1024, "audio", "audio/mpeg", "{}")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Update converted path
	convertedPath := "/converted/file.wav"
	err = UpdateSourceFileConvertedPath(db, sourceFile.ID, convertedPath)
	if err != nil {
		t.Fatalf("Failed to update source file converted path: %v", err)
	}
	
	// Verify update
	updated, err := GetSourceFileByHash(db, "hash123")
	if err != nil {
		t.Fatalf("Failed to get updated source file: %v", err)
	}
	
	if updated.ConvertedPath == nil {
		t.Error("Expected converted path to be set")
	} else if *updated.ConvertedPath != convertedPath {
		t.Errorf("Expected converted path %s, got %s", convertedPath, *updated.ConvertedPath)
	}
}

func TestListSourceFiles(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Create multiple source files with different statuses
	files := []struct {
		path   string
		hash   string
		status string
	}{
		{"/test/file1.mp3", "hash1", "pending"},
		{"/test/file2.mp3", "hash2", "processing"},
		{"/test/file3.mp3", "hash3", "completed"},
		{"/test/file4.mp3", "hash4", "pending"},
	}
	
	for _, file := range files {
		sourceFile, err := CreateSourceFile(db, file.path, file.hash, 1024, "audio", "audio/mpeg", "{}")
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}
		
		// Update status if not pending (default)
		if file.status != "pending" {
			err = UpdateSourceFileStatus(db, sourceFile.ID, file.status)
			if err != nil {
				t.Fatalf("Failed to update source file status: %v", err)
			}
		}
	}
	
	// List all files
	allFiles, err := ListSourceFiles(db, "")
	if err != nil {
		t.Fatalf("Failed to list all source files: %v", err)
	}
	
	if len(allFiles) != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), len(allFiles))
	}
	
	// List files by status
	pendingFiles, err := ListSourceFiles(db, "pending")
	if err != nil {
		t.Fatalf("Failed to list pending source files: %v", err)
	}
	
	expectedPending := 2 // file1 and file4
	if len(pendingFiles) != expectedPending {
		t.Errorf("Expected %d pending files, got %d", expectedPending, len(pendingFiles))
	}
	
	// Verify all returned files have pending status
	for _, file := range pendingFiles {
		if file.ProcessingStatus != "pending" {
			t.Errorf("Expected pending status, got %s", file.ProcessingStatus)
		}
	}
}

func TestDeleteSourceFile(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Create source file
	sourceFile, err := CreateSourceFile(db, "/test/file.mp3", "hash123", 1024, "audio", "audio/mpeg", "{}")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Delete the file
	deleted, err := DeleteSourceFile(db, sourceFile.ID)
	if err != nil {
		t.Fatalf("Failed to delete source file: %v", err)
	}
	
	// Verify deleted file information
	if deleted.ID != sourceFile.ID {
		t.Errorf("Expected deleted ID %d, got %d", sourceFile.ID, deleted.ID)
	}
	if deleted.FilePath != sourceFile.FilePath {
		t.Errorf("Expected deleted path %s, got %s", sourceFile.FilePath, deleted.FilePath)
	}
	
	// Verify file is actually deleted
	retrieved, err := GetSourceFileByHash(db, "hash123")
	if err != nil {
		t.Fatalf("Unexpected error when trying to retrieve deleted file: %v", err)
	}
	
	if retrieved != nil {
		t.Error("Expected nil for deleted file, got source file")
	}
}

func TestDeleteSourceFileNotFound(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Try to delete non-existent file
	_, err := DeleteSourceFile(db, 999)
	if err == nil {
		t.Error("Expected error when deleting non-existent file, got nil")
	}
}

func TestSourceFileTimestamps(t *testing.T) {
	db, cleanup := setupTestSourceFileDB(t)
	defer cleanup()
	
	// Record time before creation
	beforeCreate := time.Now()
	time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference
	
	// Create source file
	sourceFile, err := CreateSourceFile(db, "/test/file.mp3", "hash123", 1024, "audio", "audio/mpeg", "{}")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference
	afterCreate := time.Now()
	
	// Verify creation timestamps (be more lenient with time ranges)
	timeTolerance := 5 * time.Second
	if sourceFile.CreatedAt.Before(beforeCreate.Add(-timeTolerance)) || sourceFile.CreatedAt.After(afterCreate.Add(timeTolerance)) {
		t.Errorf("Created timestamp %v is not within expected range %v to %v", sourceFile.CreatedAt, beforeCreate.Add(-timeTolerance), afterCreate.Add(timeTolerance))
	}
	
	if sourceFile.UpdatedAt.Before(beforeCreate.Add(-timeTolerance)) || sourceFile.UpdatedAt.After(afterCreate.Add(timeTolerance)) {
		t.Errorf("Updated timestamp %v is not within expected range %v to %v", sourceFile.UpdatedAt, beforeCreate.Add(-timeTolerance), afterCreate.Add(timeTolerance))
	}
	
	// Record time before update
	time.Sleep(10 * time.Millisecond)
	beforeUpdate := time.Now()
	time.Sleep(10 * time.Millisecond)
	
	// Update the file
	err = UpdateSourceFileStatus(db, sourceFile.ID, "completed")
	if err != nil {
		t.Fatalf("Failed to update source file: %v", err)
	}
	
	time.Sleep(10 * time.Millisecond)
	afterUpdate := time.Now()
	
	// Verify update timestamps
	updated, err := GetSourceFileByHash(db, "hash123")
	if err != nil {
		t.Fatalf("Failed to get updated source file: %v", err)
	}
	
	// Created timestamp should not change
	if !updated.CreatedAt.Equal(sourceFile.CreatedAt) {
		t.Error("Created timestamp should not change on update")
	}
	
	// Updated timestamp should change (be more lenient)
	if updated.UpdatedAt.Before(beforeUpdate.Add(-timeTolerance)) || updated.UpdatedAt.After(afterUpdate.Add(timeTolerance)) {
		t.Errorf("Updated timestamp %v is not within expected range %v to %v after update", updated.UpdatedAt, beforeUpdate.Add(-timeTolerance), afterUpdate.Add(timeTolerance))
	}
	
	// Allow for timestamps to be equal (SQLite precision)
	if updated.UpdatedAt.Before(sourceFile.UpdatedAt) {
		t.Error("Updated timestamp should be at least equal to original timestamp")
	}
}
