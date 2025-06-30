package cmd

import (
	"fmt"
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

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import and process an audio or text file",
	Long: `Import and process a file. For audio files, convert to MP3 if needed, add to recordings database, 
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
	RunE: importFile,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importFile(cmd *cobra.Command, args []string) error {
	var filePath string
	
	// Initialize services with verbose logging
	verboseLogger := services.NewVerboseLogger(IsVerbose())
	verboseLogger.StartCommand("import", args)
	start := time.Now()
	successful := true
	
	defer func() {
		verboseLogger.EndCommand("import", time.Since(start), successful)
	}()

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
			Title("Select a file to import:").
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

	// Determine file type and process accordingly
	verboseLogger.Step("Determining file type", fmt.Sprintf("Extension: %s", filepath.Ext(filePath)))
	if audioService.IsValidTextFile(filePath) {
		verboseLogger.Debug("Processing as text file")
		err := processTextFile(filePath, uiService, fileService, verboseLogger)
		if err != nil {
			successful = false
		}
		return err
	} else if audioService.IsValidAudioFile(filePath) {
		verboseLogger.Debug("Processing as audio file")
		err := processAudioFile(filePath, audioService, fileService, uiService, verboseLogger)
		if err != nil {
			successful = false
		}
		return err
	} else {
		successful = false
		err := fmt.Errorf("invalid file format. Supported formats: mp3, wav, m4a, ogg, flac, md, txt. Got: %s", filepath.Ext(filePath))
		verboseLogger.Error(err, "File type validation failed")
		return err
	}
}

func processTextFile(filePath string, uiService *services.UIService, fileService *services.FileService, verboseLogger *services.VerboseLogger) error {
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
	openaiService, err := services.NewOpenAIService(verboseLogger)
	if err != nil {
		verboseLogger.Error(err, "Failed to initialize OpenAI service")
		return err
	}

	// Analyze content type and generate title
	analysis, err := uiService.RunTaskWithSpinner("🔍 Analyzing content type and generating title", func() (interface{}, error) {
		return openaiService.AnalyzeContentAndGenerateTitle(text)
	})
	if err != nil {
		return fmt.Errorf("failed to analyze content: %w", err)
	}

	contentAnalysis := analysis.(*services.ContentAnalysis)
	fmt.Printf("📋 Detected content type: %s\n", contentAnalysis.ContentType)
	fmt.Printf("📝 Generated title: %s\n", contentAnalysis.Title)

	// Create specialized summary based on content type
	summary, err := uiService.RunTaskWithSpinner(fmt.Sprintf("📝 Creating %s summary using OpenAI", contentAnalysis.ContentType), func() (interface{}, error) {
		return openaiService.SummarizeByContentType(text, contentAnalysis.ContentType)
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
	summaryPath := filepath.Join(destinationDir, "summary.md")
	transcriptionPath := filepath.Join(destinationDir, "transcription.md") // Empty for text files
	if err := fileService.SaveMarkdownFiles("", summary.(string), transcriptionPath, summaryPath); err != nil {
		return fmt.Errorf("failed to save summary file: %w", err)
	}

	// Create a metadata note in the database
	folderName := filepath.Base(destinationDir)
	noteContent := fmt.Sprintf("Text file imported and processed.\n\nContent Type: %s\n\nFiles:\n- Summary: %s\n\nFolder: %s",
		contentAnalysis.ContentType, "summary.md", folderName)
	tags := fmt.Sprintf("imported,text,%s,summary", contentAnalysis.ContentType)

	_, err = uiService.RunTaskWithSpinner("📋 Creating metadata note", func() (interface{}, error) {
		_, createErr := database.CreateNote(db, contentAnalysis.Title, noteContent, tags)
		return nil, createErr
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

func processAudioFile(filePath string, audioService *services.AudioService, fileService *services.FileService, uiService *services.UIService, verboseLogger *services.VerboseLogger) error {
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

	// Convert to MP3 if necessary
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

	// Add recording to database
	now := time.Now()
	recording := database.Recording{
		Filename:   filename,
		FilePath:   newFilePath,
		StartTime:  now,
		EndTime:    now,
		Duration:   0, // We don't know the actual recording duration for imported files
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

	// Initialize OpenAI service
	verboseLogger.Step("Initializing OpenAI service", "")
	openaiService, err := services.NewOpenAIService(verboseLogger)
	if err != nil {
		verboseLogger.Error(err, "Failed to initialize OpenAI service")
		return err
	}

	// Transcribe using chunked approach for longer files
	transcriptResult, err := uiService.RunTaskWithSpinner("📝 Transcribing audio using OpenAI", func() (interface{}, error) {
		return audioService.TranscribeFileChunked(newFilePath, destinationDir, openaiService)
	})
	if err != nil {
		return fmt.Errorf("failed to transcribe file: %w", err)
	}

	transcript := transcriptResult.(*services.ChunkedTranscriptionResult).FullTranscript

	// Analyze content type and generate title
	analysis, err := uiService.RunTaskWithSpinner("🔍 Analyzing content type and generating title", func() (interface{}, error) {
		return openaiService.AnalyzeContentAndGenerateTitle(transcript)
	})
	if err != nil {
		return fmt.Errorf("failed to analyze content: %w", err)
	}

	contentAnalysis := analysis.(*services.ContentAnalysis)
	fmt.Printf("📋 Detected content type: %s\n", contentAnalysis.ContentType)
	fmt.Printf("📝 Generated title: %s\n", contentAnalysis.Title)

	// Create specialized summary based on content type
	summary, err := uiService.RunTaskWithSpinner(fmt.Sprintf("📄 Creating %s summary from transcription", contentAnalysis.ContentType), func() (interface{}, error) {
		return openaiService.SummarizeByContentType(transcript, contentAnalysis.ContentType)
	})
	if err != nil {
		return fmt.Errorf("failed to summarize transcription: %w", err)
	}

	// Save transcription and summary as separate markdown files
	transcriptionPath := filepath.Join(destinationDir, "transcription.md")
	summaryPath := filepath.Join(destinationDir, "summary.md")

	_, err = uiService.RunTaskWithSpinner("📝 Saving transcription and summary files", func() (interface{}, error) {
		return nil, fileService.SaveMarkdownFiles(transcript, summary.(string), transcriptionPath, summaryPath)
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
	noteContent := fmt.Sprintf("Audio file imported and processed.\n\nContent Type: %s\n\nFiles:\n- Audio: %s\n- Transcription: %s\n- Summary: %s\n\nFolder: %s",
		contentAnalysis.ContentType, filename, "transcription.md", "summary.md", folderName)
	tags := fmt.Sprintf("imported,audio,%s,transcription,summary", contentAnalysis.ContentType)

	_, err = uiService.RunTaskWithSpinner("📋 Creating metadata note", func() (interface{}, error) {
		_, createErr := database.CreateNote(db, contentAnalysis.Title, noteContent, tags)
		return nil, createErr
	})
	if err != nil {
		return fmt.Errorf("failed to save metadata note: %w", err)
	}

	fmt.Printf("✅ File '%s' successfully imported and processed:\n", filename)
	fmt.Printf("   📁 Saved to notes directory: %s\n", destinationDir)
	fmt.Printf("   🎵 Added to recordings database\n")
	fmt.Printf("   📝 Transcribed and summarized\n")
	fmt.Printf("   📋 Note created: %s\n", contentAnalysis.Title)
	return nil
}
