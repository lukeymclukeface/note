package services

import (
	"crypto/md5"
	"fmt"
	"io"
	"note-cli/internal/constants"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CacheService handles temporary file operations and cleanup
type CacheService struct {
	tempDirs []string // Track temp directories for cleanup
}

// NewCacheService creates a new cache service instance
func NewCacheService() *CacheService {
	return &CacheService{
		tempDirs: make([]string, 0),
	}
}

// CacheFile represents a cached file with metadata
type CacheFile struct {
	OriginalPath string
	CachePath    string
	Hash         string
	Size         int64
	ModTime      time.Time
}

// ProcessingSession represents a temporary processing session
type ProcessingSession struct {
	ID       string
	TempDir  string
	service  *CacheService
	files    map[string]*CacheFile
	cleanup  []string // Files/dirs to cleanup
}

// CreateProcessingSession creates a new temporary processing session
func (s *CacheService) CreateProcessingSession() (*ProcessingSession, error) {
	// Create unique session ID
	sessionID := fmt.Sprintf("session_%d_%d", time.Now().Unix(), os.Getpid())
	
	// Create temp directory for this session
	tempDir, err := constants.GetTempDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get temp directory: %w", err)
	}
	
	sessionTempDir := filepath.Join(tempDir, sessionID)
	if err := os.MkdirAll(sessionTempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session temp directory: %w", err)
	}
	
	session := &ProcessingSession{
		ID:      sessionID,
		TempDir: sessionTempDir,
		service: s,
		files:   make(map[string]*CacheFile),
		cleanup: make([]string, 0),
	}
	
	// Track for cleanup
	s.tempDirs = append(s.tempDirs, sessionTempDir)
	
	return session, nil
}

// CacheInputFile copies input file to cache and returns cached path
func (session *ProcessingSession) CacheInputFile(inputPath string) (*CacheFile, error) {
	// Get file info
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat input file: %w", err)
	}
	
	// Generate cache filename with hash
	hash, err := generateFileHash(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file hash: %w", err)
	}
	
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	cacheFilename := fmt.Sprintf("%s_%s%s", baseName, hash[:8], ext)
	cachePath := filepath.Join(session.TempDir, "input", cacheFilename)
	
	// Ensure input directory exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache input directory: %w", err)
	}
	
	// Copy file to cache
	if err := copyFile(inputPath, cachePath); err != nil {
		return nil, fmt.Errorf("failed to copy file to cache: %w", err)
	}
	
	cacheFile := &CacheFile{
		OriginalPath: inputPath,
		CachePath:    cachePath,
		Hash:         hash,
		Size:         fileInfo.Size(),
		ModTime:      fileInfo.ModTime(),
	}
	
	session.files["input"] = cacheFile
	session.cleanup = append(session.cleanup, cachePath)
	
	return cacheFile, nil
}

// CreateTempFile creates a temporary file for processing
func (session *ProcessingSession) CreateTempFile(name, extension string) (string, error) {
	filename := fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), extension)
	tempPath := filepath.Join(session.TempDir, "processing", filename)
	
	// Ensure processing directory exists
	if err := os.MkdirAll(filepath.Dir(tempPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create temp processing directory: %w", err)
	}
	
	// Track for cleanup
	session.cleanup = append(session.cleanup, tempPath)
	
	return tempPath, nil
}

// SaveOutputFile saves a processed file to cache with a specific key
func (session *ProcessingSession) SaveOutputFile(key, content string) (string, error) {
	outputDir := filepath.Join(session.TempDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	
	filename := fmt.Sprintf("%s_%d.md", key, time.Now().UnixNano())
	outputPath := filepath.Join(outputDir, filename)
	
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write output file: %w", err)
	}
	
	// Don't add to cleanup - these are the desired outputs
	return outputPath, nil
}

// GetOutputFile retrieves the path of a cached output file
func (session *ProcessingSession) GetOutputFile(key string) (string, bool) {
	outputDir := filepath.Join(session.TempDir, "output")
	pattern := filepath.Join(outputDir, fmt.Sprintf("%s_*.md", key))
	
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", false
	}
	
	// Return the most recent match
	return matches[len(matches)-1], true
}

// MoveToFinalDestination moves selected files from cache to final destination
func (session *ProcessingSession) MoveToFinalDestination(outputDir string, filesToMove map[string]string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	for cacheKey, finalName := range filesToMove {
		cachedPath, exists := session.GetOutputFile(cacheKey)
		if !exists {
			return fmt.Errorf("cached file not found for key: %s", cacheKey)
		}
		
		finalPath := filepath.Join(outputDir, finalName)
		
		// Move file from cache to final destination
		if err := moveFile(cachedPath, finalPath); err != nil {
			return fmt.Errorf("failed to move %s to final destination: %w", cacheKey, err)
		}
	}
	
	return nil
}

// Cleanup removes all temporary files and directories for this session
func (session *ProcessingSession) Cleanup() error {
	var errors []string
	
	// Remove all tracked files
	for _, path := range session.cleanup {
		if err := os.RemoveAll(path); err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove %s: %v", path, err))
		}
	}
	
	// Remove session temp directory
	if err := os.RemoveAll(session.TempDir); err != nil {
		errors = append(errors, fmt.Sprintf("failed to remove session directory %s: %v", session.TempDir, err))
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// CleanupOldCache removes old cache files (older than 24 hours)
func (s *CacheService) CleanupOldCache() error {
	cacheDir, err := constants.GetCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get cache directory: %w", err)
	}
	
	cutoff := time.Now().Add(-24 * time.Hour)
	
	return filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		
		// Skip directories and recent files
		if info.IsDir() || info.ModTime().After(cutoff) {
			return nil
		}
		
		// Remove old files
		if err := os.Remove(path); err != nil {
			// Log but don't fail the cleanup
			fmt.Printf("Warning: failed to remove old cache file %s: %v\n", path, err)
		}
		
		return nil
	})
}

// InitializeCache ensures cache directory structure exists
func (s *CacheService) InitializeCache() error {
	cacheDir, err := constants.GetCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get cache directory: %w", err)
	}
	
	// Create cache directory structure
	dirs := []string{
		cacheDir,
		filepath.Join(cacheDir, "temp"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// Helper functions

func generateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	
	_, err = io.Copy(destination, source)
	return err
}

func moveFile(src, dst string) error {
	// Try to rename first (fastest for same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	
	// Fall back to copy + delete
	if err := copyFile(src, dst); err != nil {
		return err
	}
	
	return os.Remove(src)
}
