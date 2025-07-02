package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/database"
	"note-cli/internal/helpers"
	"note-cli/internal/services"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var summariseCmd = &cobra.Command{
	Use:   "summarise [file]",
	Short: "Summarise and process an audio or text file",
	Long: `Summarise and process a file. For audio files, convert to MP3 if needed, add to recordings database, 
transcribe, create a summary, and save as a note.
For text files, read the content, create a summary, and save as a note.

If no file is specified, you'll be presented with a list of compatible files 
in the current directory to choose from.

Supported formats:
- Audio: mp3, wav, m4a, ogg, flac
- Text: md, txt

Requires:
- OpenAI API key configured (run 'note setup')
- ffmpeg installed for audio format conversion`,
	Args: cobra.MaximumNArgs(1),
	RunE: summariseFile,
}

func init() {
	rootCmd.AddCommand(summariseCmd)
}

func summariseFile(cmd *cobra.Command, args []string) error {
	var filePath string
	
	// Initialize services with verbose logging
	verboseLogger := services.NewVerboseLogger(IsVerbose())
	verboseLogger.StartCommand("summarise", args)
	start := time.Now()
	successful := true
	
	defer func() {
		verboseLogger.EndCommand("summarise", time.Since(start), successful)
	}()


	// Initialize provider factory
	factory := services.NewProviderFactory(verboseLogger)

	// Load config to get providers
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Section for selecting transcription provider
	transcriptionProvider, err := factory.CreateTranscriptionProvider(cfg.TranscriptionProvider)
	if err != nil {
		return fmt.Errorf("failed to create transcription provider: %w", err)
	}

	// Section for selecting summary provider
	summaryProvider, err := factory.CreateSummaryProvider(cfg.SummaryProvider)
	if err != nil {
		return fmt.Errorf("failed to create summary provider: %w", err)
	}

	audioService := services.NewAudioService()
	fileService := services.NewFileService()
	uiService := services.NewUIService()

	if len(args) == 1 {
		filePath = args[0]
	} else {
		// No argument provided, find compatible files in the current directory
		files, err := fileService.GetCompatibleFiles(".")
		if err != nil {
			return fmt.Errorf("failed to scan directory: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("no compatible files found in the current directory. Supported formats: mp3, wav, m4a, ogg, flac, md, txt")
		}

		// Create options for selection
		var options []huh.Option[string]
		for _, file := range files {
			options = append(options, huh.NewOption(file, file))
		}

		selectField := huh.NewSelect[string]().
			Title("Select a file to summarise:").
			Options(options...).
			Value(&filePath)

		form := huh.NewForm(huh.NewGroup(selectField))

		if err := form.Run(); err != nil {
			return fmt.Errorf("file selection cancelled")
		}
	}

	// Check if file exists
	verboseLogger.Step("Validating input file", fmt.Sprintf("File path: %s", filePath))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		successful = false
		verboseLogger.Error(err, "File validation failed")
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Initialize database connection for source file tracking
	cfg, db, err := helpers.LoadConfigAndDatabase()
	if err != nil {
		successful = false
		verboseLogger.Error(err, "Failed to initialize database")
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Check if file is already being processed or has been processed
	verboseLogger.Step("Checking file processing status", "")
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		successful = false
		verboseLogger.Error(err, "Failed to get absolute path")
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if file already exists in source_files table
	existingFile, err := database.GetSourceFileByPath(db, absFilePath)
	if err != nil {
		successful = false
		verboseLogger.Error(err, "Failed to check existing file")
		return fmt.Errorf("failed to check existing file: %w", err)
	}

	if existingFile != nil {
		if existingFile.ProcessingStatus == "completed" {
			return fmt.Errorf("file has already been processed successfully (ID: %d, converted to: %s)", existingFile.ID, helpers.SafeStringDeref(existingFile.ConvertedPath))
		} else if existingFile.ProcessingStatus == "processing" {
			return fmt.Errorf("file is currently being processed (ID: %d)", existingFile.ID)
		} else if existingFile.ProcessingStatus == "failed" {
			verboseLogger.Debug(fmt.Sprintf("File previously failed processing (ID: %d), retrying", existingFile.ID))
			// Update status to processing to retry
			if err := database.UpdateSourceFileStatus(db, existingFile.ID, "processing"); err != nil {
				verboseLogger.Error(err, "Failed to update file status for retry")
				return fmt.Errorf("failed to update file status for retry: %w", err)
			}
		}
	} else {
		// Create new source file record
		verboseLogger.Step("Creating source file record", "")
		fileHash, mimeType, fileType, fileSize, metadata, err := helpers.ProcessFileForDatabase(absFilePath)
		if err != nil {
			successful = false
			verboseLogger.Error(err, "Failed to process file metadata")
			return fmt.Errorf("failed to process file metadata: %w", err)
		}

		// Check if file with same hash already exists
		existingByHash, err := database.GetSourceFileByHash(db, fileHash)
		if err != nil {
			successful = false
			verboseLogger.Error(err, "Failed to check file by hash")
			return fmt.Errorf("failed to check file by hash: %w", err)
		}

		if existingByHash != nil {
			return fmt.Errorf("a file with identical content has already been processed (original: %s, ID: %d)", existingByHash.FilePath, existingByHash.ID)
		}

		// Create source file record
		existingFile, err = database.CreateSourceFile(db, absFilePath, fileHash, fileSize, fileType, mimeType, metadata)
		if err != nil {
			successful = false
			verboseLogger.Error(err, "Failed to create source file record")
			return fmt.Errorf("failed to create source file record: %w", err)
		}

		// Update status to processing
		if err := database.UpdateSourceFileStatus(db, existingFile.ID, "processing"); err != nil {
			verboseLogger.Error(err, "Failed to update file status to processing")
			return fmt.Errorf("failed to update file status to processing: %w", err)
		}
	}

	// Determine file type and process accordingly
	verboseLogger.Step("Determining file type", fmt.Sprintf("Extension: %s", filepath.Ext(filePath)))
	var processingErr error
	var outputPath string
	
	if audioService.IsValidTextFile(filePath) {
		verboseLogger.Debug("Processing as text file")
		outputPath, processingErr = processTextFileWithTracking(absFilePath, uiService, fileService, verboseLogger, summaryProvider)
	} else if audioService.IsValidAudioFile(filePath) {
		verboseLogger.Debug("Processing as audio file")
		outputPath, processingErr = processAudioFileWithTracking(absFilePath, audioService, fileService, uiService, verboseLogger, transcriptionProvider, summaryProvider)
	} else {
		processingErr = fmt.Errorf("invalid file format. Supported formats: mp3, wav, m4a, ogg, flac, md, txt. Got: %s", filepath.Ext(filePath))
		verboseLogger.Error(processingErr, "File type validation failed")
	}
	
	// Update source file status based on processing result
	if processingErr != nil {
		successful = false
		// Mark as failed
		if updateErr := database.UpdateSourceFileStatus(db, existingFile.ID, "failed"); updateErr != nil {
			verboseLogger.Error(updateErr, "Failed to update file status to failed")
		}
		return processingErr
	} else {
		// Mark as completed and set converted path if available
		if outputPath != "" {
			if updateErr := database.UpdateSourceFileConvertedPath(db, existingFile.ID, outputPath); updateErr != nil {
				verboseLogger.Error(updateErr, "Failed to update converted path")
			}
		}
		if updateErr := database.UpdateSourceFileStatus(db, existingFile.ID, "completed"); updateErr != nil {
			verboseLogger.Error(updateErr, "Failed to update file status to completed")
			return fmt.Errorf("processing succeeded but failed to update status: %w", updateErr)
		}
	}
	
	return nil
}

func processTextFile(filePath string, uiService *services.UIService, fileService *services.FileService, verboseLogger *services.VerboseLogger, summaryProvider services.AIProvider) error {
	// Load config and validate
	cfg, db, err := helpers.LoadConfigAndDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := helpers.ValidateOpenAIKey(cfg); err != nil {
		return err
	}

	// Read text content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read text file: %w", err)
	}

	text := string(content)

	// Initialize OpenAI service
	verboseLogger.Step("Initializing OpenAI service", "")
// Analyze content type and generate title
	analysis, err := uiService.RunTaskWithSpinner("🔍 Analyzing content type and generating title", func() (interface{}, error) {
		return summaryProvider.AnalyzeContentAndGenerateTitle(text)
	})
	if err != nil {
		return fmt.Errorf("failed to analyze content: %w", err)
	}

	contentAnalysis := analysis.(*services.ContentAnalysis)
	fmt.Printf("📋 Detected content type: %s\n", contentAnalysis.ContentType)
	fmt.Printf("📝 Generated title: %s\n", contentAnalysis.Title)

	// Create specialized summary based on content type
	summaryResult, err := uiService.RunTaskWithSpinner(fmt.Sprintf("📝 Creating %s summary using OpenAI", contentAnalysis.ContentType), func() (interface{}, error) {
return summaryProvider.SummarizeByContentType(text, contentAnalysis.ContentType)
	})
	if err != nil {
		return fmt.Errorf("failed to summarize text: %w", err)
	}

	// Create content directory using the generated title and detected content type
	safeTitle := strings.ReplaceAll(contentAnalysis.Title, "/", "-")
	safeTitle = strings.ReplaceAll(safeTitle, ":", "-")
	destinationDir, err := fileService.CreateContentDirectory(safeTitle, contentAnalysis.ContentType)
	if err != nil {
		return err
	}

// Save the summary
summaryStr := summaryResult.(string)
summaryPath := filepath.Join(destinationDir, "summary.md")
transcriptionPath := filepath.Join(destinationDir, "transcription.md") // Empty for text files
if err := fileService.SaveMarkdownFiles("", summaryStr, transcriptionPath, summaryPath); err != nil {
	return fmt.Errorf("failed to save summary file: %w", err)
}

	// Create a metadata note in the database
	folderName := filepath.Base(destinationDir)
	noteContent := fmt.Sprintf("Text file summarised and processed.\n\nContent Type: %s\n\nFiles:\n- Summary: %s\n\nFolder: %s",
		contentAnalysis.ContentType, "summary.md", folderName)
	tags := fmt.Sprintf("summarised,text,%s,summary", contentAnalysis.ContentType)

	_, err = uiService.RunTaskWithSpinner("📋 Creating metadata note", func() (interface{}, error) {
		switch contentAnalysis.ContentType {
		case "meeting":
			_, createErr := database.CreateMeeting(db, contentAnalysis.Title, noteContent, summaryStr, "", "", tags, nil, nil)
			return nil, createErr
		case "interview":
			_, createErr := database.CreateInterview(db, contentAnalysis.Title, noteContent, summaryStr, "", "", "", "", tags, nil, nil)
			return nil, createErr
		default:
			_, createErr := database.CreateNote(db, contentAnalysis.Title, noteContent, summaryStr, tags, nil)
			return nil, createErr
		}
	})
	if err != nil {
		return fmt.Errorf("failed to save metadata note: %w", err)
	}

	fmt.Printf("✅ Text file '%s' successfully processed and summarized:\n", filePath)
	fmt.Printf("   📁 Saved to notes directory: %s\n", destinationDir)
	fmt.Printf("   📝 Summary saved at: %s\n", summaryPath)
	fmt.Printf("   📋 Note created: %s\n", contentAnalysis.Title)

	return nil
}

func processAudioFile(filePath string, audioService *services.AudioService, fileService *services.FileService, uiService *services.UIService, verboseLogger *services.VerboseLogger, transcriptionProvider services.AIProvider, summaryProvider services.AIProvider) error {
	fmt.Println("🎵 Processing audio file...")

	// Check dependencies
	if err := audioService.CheckDependencies(); err != nil {
		return err
	}

	// Load config and validate
	cfg, db, err := helpers.LoadConfigAndDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := helpers.ValidateOpenAIKey(cfg); err != nil {
		return err
	}

	// Initialize cache service
	cacheService := services.NewCacheService()
	if err := cacheService.InitializeCache(); err != nil {
		verboseLogger.Error(err, "Failed to initialize cache, proceeding without cache")
	}

	// Convert to MP3 if necessary (using cache)
	mp3Path, err := uiService.RunTaskWithSpinner("🔄 Converting to MP3 format", func() (interface{}, error) {
		return audioService.ConvertToMP3(filePath)
	})
	if err != nil {
		return fmt.Errorf("failed to convert file to MP3: %w", err)
	}

	// Create a unique folder under notes directory
	originalFilename := strings.TrimSuffix(filepath.Base(mp3Path.(string)), filepath.Ext(mp3Path.(string)))
	destinationDir, err := fileService.CreateNoteDirectory(originalFilename)
	if err != nil {
		return err
	}

	// Copy MP3 to the new directory
	filename := filepath.Base(mp3Path.(string))
	newFilePath := filepath.Join(destinationDir, filename)
	
	_, err = uiService.RunTaskWithSpinner("📁 Copying MP3 file to notes directory", func() (interface{}, error) {
		return nil, fileService.CopyFile(mp3Path.(string), newFilePath)
	})
	if err != nil {
		return fmt.Errorf("failed to copy MP3 file: %w", err)
	}

	// Get file info for complete recording metadata
	fileInfo, err := os.Stat(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Get actual audio duration using ffprobe
	durationSeconds, err := audioService.GetAudioDuration(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to get audio duration: %w", err)
	}
	duration := time.Duration(durationSeconds) * time.Second

	// Add recording to database
	now := time.Now()
	recording := database.Recording{
		Filename:   filename,
		FilePath:   newFilePath,
		StartTime:  now,
		EndTime:    now.Add(duration), // Set end time based on actual duration
		Duration:   duration,           // Use actual duration from ffprobe
		FileSize:   fileInfo.Size(),
		Format:     "mp3",
		SampleRate: 44100, // Default sample rate
		Channels:   1,     // Assume mono
		CreatedAt:  now,
	}

	_, err = uiService.RunTaskWithSpinner("💾 Saving recording metadata to database", func() (interface{}, error) {
		return nil, database.SaveRecording(cfg.DatabasePath, &recording)
	})
	if err != nil {
		return fmt.Errorf("failed to save recording: %w", err)
	}

	// Transcribe using helper function
	transcriptResult, err := uiService.RunTaskWithSpinner("📝 Transcribing audio", func() (interface{}, error) {
		return helpers.TranscribeAudioFile(newFilePath, destinationDir, verboseLogger)
	})
	if err != nil {
		return fmt.Errorf("failed to transcribe file: %w", err)
	}

	transcript := transcriptResult.(*helpers.TranscriptionResult).FullTranscript

	// Analyze content type and generate title
	analysis, err := uiService.RunTaskWithSpinner("🔍 Analyzing content type and generating title", func() (interface{}, error) {
		return summaryProvider.AnalyzeContentAndGenerateTitle(transcript)
	})
	if err != nil {
		return fmt.Errorf("failed to analyze content: %w", err)
	}

	contentAnalysis := analysis.(*services.ContentAnalysis)
	fmt.Printf("📋 Detected content type: %s\n", contentAnalysis.ContentType)
	fmt.Printf("📝 Generated title: %s\n", contentAnalysis.Title)

	// Create specialized summary based on content type
	summaryResult, err := uiService.RunTaskWithSpinner(fmt.Sprintf("📄 Creating %s summary from transcription", contentAnalysis.ContentType), func() (interface{}, error) {
return summaryProvider.SummarizeByContentType(transcript, contentAnalysis.ContentType)
	})
	if err != nil {
		return fmt.Errorf("failed to summarize transcription: %w", err)
	}

// Save transcription and summary as separate markdown files
summaryStr := summaryResult.(string)
transcriptionPath := filepath.Join(destinationDir, "transcription.md")
summaryPath := filepath.Join(destinationDir, "summary.md")

_, err = uiService.RunTaskWithSpinner("📝 Saving transcription and summary files", func() (interface{}, error) {
	return nil, fileService.SaveMarkdownFiles(transcript, summaryStr, transcriptionPath, summaryPath)
})
	if err != nil {
		return fmt.Errorf("failed to save markdown files: %w", err)
	}

	// Update destination directory to use the generated title and content type
	safeTitle := strings.ReplaceAll(contentAnalysis.Title, "/", "-")
	safeTitle = strings.ReplaceAll(safeTitle, ":", "-")
	newDestinationDir, err := fileService.CreateContentDirectory(safeTitle, contentAnalysis.ContentType)
	if err != nil {
		return err
	}

	// Move files to the new directory with proper title
	if newDestinationDir != destinationDir {
		// Copy files to new location
		newAudioPath := filepath.Join(newDestinationDir, filename)
		if err := fileService.CopyFile(filepath.Join(destinationDir, filename), newAudioPath); err != nil {
			return fmt.Errorf("failed to move MP3 file: %w", err)
		}

		// Copy markdown files to new location
		newTranscriptionPath := filepath.Join(newDestinationDir, "transcription.md")
		newSummaryPath := filepath.Join(newDestinationDir, "summary.md")
		
		if err := fileService.CopyFile(transcriptionPath, newTranscriptionPath); err != nil {
			return fmt.Errorf("failed to move transcription file: %w", err)
		}
		
		if err := fileService.CopyFile(summaryPath, newSummaryPath); err != nil {
			return fmt.Errorf("failed to move summary file: %w", err)
		}

		// Remove old directory
		if err := os.RemoveAll(destinationDir); err != nil {
			fmt.Printf("Warning: failed to remove old directory %s: %v\n", destinationDir, err)
		}

		// Update paths
		transcriptionPath = newTranscriptionPath
		summaryPath = newSummaryPath
		destinationDir = newDestinationDir
	}

	// Create a metadata note that references the files
	folderName := filepath.Base(destinationDir)
	noteContent := fmt.Sprintf("Audio file summarised and processed.\n\nContent Type: %s\n\nFiles:\n- Audio: %s\n- Transcription: %s\n- Summary: %s\n\nFolder: %s",
		contentAnalysis.ContentType, filename, "transcription.md", "summary.md", folderName)
	tags := fmt.Sprintf("summarised,audio,%s,transcription,summary", contentAnalysis.ContentType)

	_, err = uiService.RunTaskWithSpinner("📋 Creating metadata note", func() (interface{}, error) {
		switch contentAnalysis.ContentType {
		case "meeting":
			_, createErr := database.CreateMeeting(db, contentAnalysis.Title, noteContent, summaryStr, "", "", tags, nil, nil)
			return nil, createErr
		case "interview":
			_, createErr := database.CreateInterview(db, contentAnalysis.Title, noteContent, summaryStr, "", "", "", "", tags, nil, nil)
			return nil, createErr
		default:
			_, createErr := database.CreateNote(db, contentAnalysis.Title, noteContent, summaryStr, tags, nil)
			return nil, createErr
		}
	})
	if err != nil {
		return fmt.Errorf("failed to save metadata note: %w", err)
	}

	fmt.Printf("✅ File '%s' successfully summarised and processed:\n", filename)
	fmt.Printf("   📁 Saved to notes directory: %s\n", destinationDir)
	fmt.Printf("   🎵 Added to recordings database\n")
	fmt.Printf("   📝 Transcribed and summarized\n")
	fmt.Printf("   📋 Note created: %s\n", contentAnalysis.Title)
	return nil
}

// processTextFileWithTracking processes a text file and returns the output path
func processTextFileWithTracking(filePath string, uiService *services.UIService, fileService *services.FileService, verboseLogger *services.VerboseLogger, summaryProvider services.AIProvider) (string, error) {
	err := processTextFile(filePath, uiService, fileService, verboseLogger, summaryProvider)
	if err != nil {
		return "", err
	}
	
	// Return the base directory where files are stored
	// For text files, we can't predict the exact path without reproducing the logic
	// but we can return a general indication
	return "processed-to-notes-directory", nil
}

// processAudioFileWithTracking processes an audio file and returns the output path
func processAudioFileWithTracking(filePath string, audioService *services.AudioService, fileService *services.FileService, uiService *services.UIService, verboseLogger *services.VerboseLogger, transcriptionProvider services.AIProvider, summaryProvider services.AIProvider) (string, error) {
	err := processAudioFile(filePath, audioService, fileService, uiService, verboseLogger, transcriptionProvider, summaryProvider)
	if err != nil {
		return "", err
	}
	
	// Return the base directory where files are stored
	// For audio files, we can't predict the exact path without reproducing the logic
	// but we can return a general indication
	return "processed-to-notes-directory", nil
}
