package helpers

import (
	"database/sql"
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/constants"
	"note-cli/internal/database"
	"note-cli/internal/services"
	"os"
	"path/filepath"
)

// LoadConfigWithValidation loads configuration and validates required fields
func LoadConfigWithValidation() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

// LoadConfigAndDatabase loads config and opens database connection
func LoadConfigAndDatabase() (*config.Config, *sql.DB, error) {
	cfg, err := LoadConfigWithValidation()
	if err != nil {
		return nil, nil, err
	}

	if cfg.DatabasePath == "" {
		return nil, nil, fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return cfg, db, nil
}

// GetDatabaseConnection opens a database connection using the default path
func GetDatabaseConnection() (*sql.DB, error) {
	dbPath, err := constants.GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}
	return database.Connect(dbPath)
}

// ValidateOpenAIKey checks if OpenAI API key is configured
func ValidateOpenAIKey(cfg *config.Config) error {
	if cfg.OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not configured. Please run 'note setup' first")
	}
	return nil
}

// TranscriptionResult holds the result of audio transcription
type TranscriptionResult struct {
	FullTranscript string
	ChunkFiles     []string
}

// TranscribeAudioFile transcribes an audio file using the configured provider with caching
func TranscribeAudioFile(filePath string, outputDir string, verboseLogger *services.VerboseLogger) (*TranscriptionResult, error) {
	return TranscribeAudioFileWithCache(filePath, outputDir, verboseLogger, true)
}

// TranscribeAudioFileWithCache transcribes an audio file with optional caching
func TranscribeAudioFileWithCache(filePath string, outputDir string, verboseLogger *services.VerboseLogger, useCache bool) (*TranscriptionResult, error) {
	// Validate input file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Initialize services
	factory := services.NewProviderFactory(verboseLogger)
	audioService := services.NewAudioService()
	cacheService := services.NewCacheService()

	// Validate audio file format
	if !audioService.IsValidAudioFile(filePath) {
		return nil, fmt.Errorf("unsupported audio file format: %s", filePath)
	}

	// Load configuration
	cfg, err := LoadConfigWithValidation()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create transcription provider
	transcriptionProvider, err := factory.CreateTranscriptionProvider(cfg.TranscriptionProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcription provider: %w", err)
	}

	// Check audio dependencies
	if err := audioService.CheckDependencies(); err != nil {
		return nil, err
	}

	// Initialize cache if requested
	if useCache {
		if err := cacheService.InitializeCache(); err != nil {
			if verboseLogger != nil {
				verboseLogger.Error(err, "Failed to initialize cache, proceeding without cache")
			}
			useCache = false
		}
	}

	var result *services.ChunkedTranscriptionResult
	if useCache {
		// Use cache-based transcription
		result, err = transcribeWithCache(filePath, outputDir, audioService, transcriptionProvider, cacheService, verboseLogger)
	} else {
		// Fallback to direct transcription
		result, err = audioService.TranscribeFileChunked(filePath, outputDir, transcriptionProvider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio file: %w", err)
	}

	return &TranscriptionResult{
		FullTranscript: result.FullTranscript,
		ChunkFiles:     result.ChunkFiles,
	}, nil
}

// transcribeWithCache handles the cached transcription workflow
func transcribeWithCache(filePath, outputDir string, audioService *services.AudioService, provider services.AIProvider, cacheService *services.CacheService, verboseLogger *services.VerboseLogger) (*services.ChunkedTranscriptionResult, error) {
	// Create processing session
	session, err := cacheService.CreateProcessingSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create processing session: %w", err)
	}
	defer func() {
		if cleanupErr := session.Cleanup(); cleanupErr != nil {
			if verboseLogger != nil {
				verboseLogger.Error(cleanupErr, "Failed to cleanup cache session")
			}
		}
	}()

	// Cache the input file
	if verboseLogger != nil {
		verboseLogger.Step("Caching input file", filePath)
	}
	cachedFile, err := session.CacheInputFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to cache input file: %w", err)
	}

	// Transcribe using the cached file
	if verboseLogger != nil {
		verboseLogger.Step("Transcribing from cache", cachedFile.CachePath)
	}
	result, err := audioService.TranscribeFileChunkedWithCache(cachedFile.CachePath, "", provider, session)
	if err != nil {
		return nil, err
	}

	// If outputDir is specified, move relevant files to final destination
	if outputDir != "" {
		if verboseLogger != nil {
			verboseLogger.Step("Moving files to final destination", outputDir)
		}

		// Save the full transcription to cache first
		if _, saveErr := session.SaveOutputFile("transcription", result.FullTranscript); saveErr != nil {
			return nil, fmt.Errorf("failed to save transcription to cache: %w", saveErr)
		}

		// Prepare files to move
		filesToMove := map[string]string{
			"transcription": "transcription.md",
		}

		// Add chunk files if they exist
		for i := 1; i <= len(result.ChunkFiles); i++ {
			chunkKey := fmt.Sprintf("chunk_%02d", i)
			chunkFile := fmt.Sprintf("transcription_chunk_%02d.md", i)
			if _, exists := session.GetOutputFile(chunkKey); exists {
				filesToMove[chunkKey] = chunkFile
			}
		}

		// Move files to final destination
		if err := session.MoveToFinalDestination(outputDir, filesToMove); err != nil {
			return nil, fmt.Errorf("failed to move files to destination: %w", err)
		}

		// Update result chunk files to point to final destination
		updatedChunkFiles := make([]string, 0, len(result.ChunkFiles))
		for i := 1; i <= len(result.ChunkFiles); i++ {
			filename := fmt.Sprintf("transcription_chunk_%02d.md", i)
			updatedChunkFiles = append(updatedChunkFiles, filepath.Join(outputDir, filename))
		}
		result.ChunkFiles = updatedChunkFiles
	}

	return result, nil
}

// SafeStringDeref safely dereferences a string pointer, returning empty string if nil
func SafeStringDeref(s *string) string {
	if s == nil {
		return "<none>"
	}
	return *s
}
