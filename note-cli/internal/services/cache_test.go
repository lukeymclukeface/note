package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheService_CreateProcessingSession(t *testing.T) {
	cacheService := NewCacheService()
	
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	
	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}
	
	if session.TempDir == "" {
		t.Error("Session temp directory should not be empty")
	}
	
	// Check that directory was created
	if _, err := os.Stat(session.TempDir); os.IsNotExist(err) {
		t.Error("Session temp directory should exist")
	}
	
	// Cleanup
	session.Cleanup()
}

func TestProcessingSession_CacheInputFile(t *testing.T) {
	cacheService := NewCacheService()
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	defer session.Cleanup()
	
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "This is a test file for caching"
	
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Cache the file
	cachedFile, err := session.CacheInputFile(testFile)
	if err != nil {
		t.Fatalf("Failed to cache input file: %v", err)
	}
	
	// Verify cached file exists
	if _, err := os.Stat(cachedFile.CachePath); os.IsNotExist(err) {
		t.Error("Cached file should exist")
	}
	
	// Verify content is the same
	cachedContent, err := os.ReadFile(cachedFile.CachePath)
	if err != nil {
		t.Fatalf("Failed to read cached file: %v", err)
	}
	
	if string(cachedContent) != testContent {
		t.Error("Cached file content should match original")
	}
	
	// Verify metadata
	if cachedFile.OriginalPath != testFile {
		t.Error("Original path should be preserved")
	}
	
	if cachedFile.Hash == "" {
		t.Error("File hash should be generated")
	}
}

func TestProcessingSession_CreateTempFile(t *testing.T) {
	cacheService := NewCacheService()
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	defer session.Cleanup()
	
	tempPath, err := session.CreateTempFile("test", ".txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	// Check that the path is in the session directory
	if !filepath.HasPrefix(tempPath, session.TempDir) {
		t.Error("Temp file should be in session directory")
	}
	
	// Check that file has correct extension
	if filepath.Ext(tempPath) != ".txt" {
		t.Error("Temp file should have correct extension")
	}
	
	// Write to the file to verify it's accessible
	testContent := "temp file content"
	if err := os.WriteFile(tempPath, []byte(testContent), 0644); err != nil {
		t.Errorf("Failed to write to temp file: %v", err)
	}
}

func TestProcessingSession_SaveOutputFile(t *testing.T) {
	cacheService := NewCacheService()
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	defer session.Cleanup()
	
	testContent := "# Test Output\nThis is test output content"
	outputPath, err := session.SaveOutputFile("test", testContent)
	if err != nil {
		t.Fatalf("Failed to save output file: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file should exist")
	}
	
	// Verify content
	savedContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if string(savedContent) != testContent {
		t.Error("Saved content should match original")
	}
	
	// Test retrieval
	retrievedPath, exists := session.GetOutputFile("test")
	if !exists {
		t.Error("Should be able to retrieve saved output file")
	}
	
	if retrievedPath != outputPath {
		t.Error("Retrieved path should match saved path")
	}
}

func TestProcessingSession_MoveToFinalDestination(t *testing.T) {
	cacheService := NewCacheService()
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	defer session.Cleanup()
	
	// Create test output files
	testContent1 := "Test content 1"
	testContent2 := "Test content 2"
	
	_, err = session.SaveOutputFile("file1", testContent1)
	if err != nil {
		t.Fatalf("Failed to save output file 1: %v", err)
	}
	
	_, err = session.SaveOutputFile("file2", testContent2)
	if err != nil {
		t.Fatalf("Failed to save output file 2: %v", err)
	}
	
	// Create destination directory
	destDir := t.TempDir()
	
	// Move files to final destination
	filesToMove := map[string]string{
		"file1": "final1.md",
		"file2": "final2.md",
	}
	
	err = session.MoveToFinalDestination(destDir, filesToMove)
	if err != nil {
		t.Fatalf("Failed to move files to final destination: %v", err)
	}
	
	// Verify files exist in destination
	final1Path := filepath.Join(destDir, "final1.md")
	final2Path := filepath.Join(destDir, "final2.md")
	
	if _, err := os.Stat(final1Path); os.IsNotExist(err) {
		t.Error("Final file 1 should exist in destination")
	}
	
	if _, err := os.Stat(final2Path); os.IsNotExist(err) {
		t.Error("Final file 2 should exist in destination")
	}
	
	// Verify content
	content1, err := os.ReadFile(final1Path)
	if err != nil {
		t.Fatalf("Failed to read final file 1: %v", err)
	}
	
	if string(content1) != testContent1 {
		t.Error("Final file 1 content should match original")
	}
}

func TestProcessingSession_Cleanup(t *testing.T) {
	cacheService := NewCacheService()
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		t.Fatalf("Failed to create processing session: %v", err)
	}
	
	// Create some temp files
	tempPath, err := session.CreateTempFile("test", ".txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	if err := os.WriteFile(tempPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	
	// Save some output files
	_, err = session.SaveOutputFile("output", "test output")
	if err != nil {
		t.Fatalf("Failed to save output file: %v", err)
	}
	
	sessionDir := session.TempDir
	
	// Verify files exist before cleanup
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Error("Temp file should exist before cleanup")
	}
	
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		t.Error("Session directory should exist before cleanup")
	}
	
	// Cleanup
	err = session.Cleanup()
	if err != nil {
		t.Fatalf("Failed to cleanup session: %v", err)
	}
	
	// Verify files are removed after cleanup
	if _, err := os.Stat(sessionDir); !os.IsNotExist(err) {
		t.Error("Session directory should be removed after cleanup")
	}
}

func TestCacheService_InitializeCache(t *testing.T) {
	cacheService := NewCacheService()
	
	err := cacheService.InitializeCache()
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	
	// This test is mainly to ensure no errors occur during initialization
	// The actual directory creation is tested implicitly by other tests
}

func TestGenerateFileHash(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "This is a test file for hashing"
	
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	hash1, err := generateFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to generate file hash: %v", err)
	}
	
	if hash1 == "" {
		t.Error("Hash should not be empty")
	}
	
	// Generate hash again - should be the same
	hash2, err := generateFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to generate file hash second time: %v", err)
	}
	
	if hash1 != hash2 {
		t.Error("Hash should be consistent for same file")
	}
	
	// Modify file and check hash changes
	modifiedContent := testContent + " modified"
	if err := os.WriteFile(testFile, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	
	hash3, err := generateFileHash(testFile)
	if err != nil {
		t.Fatalf("Failed to generate hash for modified file: %v", err)
	}
	
	if hash1 == hash3 {
		t.Error("Hash should change when file content changes")
	}
}

func TestCacheCopyFile(t *testing.T) {
	// Create source file
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")
	testContent := "This is test content for copying"
	
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Copy file
	err := copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
	
	// Verify destination exists
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("Destination file should exist after copy")
	}
	
	// Verify content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	
	if string(dstContent) != testContent {
		t.Error("Destination file content should match source")
	}
}

func TestMoveFile(t *testing.T) {
	// Create source file
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")
	testContent := "This is test content for moving"
	
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Move file
	err := moveFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Failed to move file: %v", err)
	}
	
	// Verify source no longer exists
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("Source file should not exist after move")
	}
	
	// Verify destination exists
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("Destination file should exist after move")
	}
	
	// Verify content
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	
	if string(dstContent) != testContent {
		t.Error("Destination file content should match original source")
	}
}
