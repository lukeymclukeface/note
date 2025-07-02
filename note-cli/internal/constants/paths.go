package constants

import (
	"os"
	"path/filepath"
)

// BaseDirectoryName is the name of the main application directory
const BaseDirectoryName = ".noteai"

// GetBaseDir returns the base directory path for the application
func GetBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, BaseDirectoryName), nil
}

// GetNotesDir returns the notes directory path
func GetNotesDir() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "notes"), nil
}

// GetMeetingsDir returns the meetings directory path
func GetMeetingsDir() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "meetings"), nil
}

// GetInterviewsDir returns the interviews directory path
func GetInterviewsDir() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "interviews"), nil
}

// GetRecordingsDir returns the recordings directory path
func GetRecordingsDir() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "recordings"), nil
}

// GetDatabasePath returns the database file path
func GetDatabasePath() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "notes.db"), nil
}

// GetConfigPath returns the config file path
func GetConfigPath() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, "config.json"), nil
}

// GetCacheDir returns the cache directory path
func GetCacheDir() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, ".cache"), nil
}

// GetTempDir returns a temporary directory within the cache for processing
func GetTempDir() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "temp"), nil
}
