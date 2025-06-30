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

	// Initialize services
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
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Determine file type and process accordingly
	if audioService.IsValidTextFile(filePath) {
		return processTextFile(filePath, uiService, fileService)
	} else if audioService.IsValidAudioFile(filePath) {
		return processAudioFile(filePath, audioService, fileService, uiService)
	} else {
		return fmt.Errorf("invalid file format. Supported formats: mp3, wav, m4a, ogg, flac, md, txt. Got: %s", filepath.Ext(filePath))
	}
}

func processTextFile(filePath string, uiService *services.UIService, fileService *services.FileService) error {
	// Load config and validate
	cfg, err := helpers.LoadConfigWithValidation()
	if err != nil {
		return err
	}

	if err := helpers.ValidateOpenAIKey(cfg); err != nil {
		return err
	}

	// Read text content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read text file: %w", err)
	}

	text := string(content)

	// Initialize OpenAI service and summarize
	openaiService, err := services.NewOpenAIService()
	if err != nil {
		return err
	}

	summary, err := uiService.RunTaskWithSpinner("📝 Creating summary using OpenAI", func() (interface{}, error) {
		return openaiService.SummarizeText(text)
	})
	if err != nil {
		return fmt.Errorf("failed to summarize text: %w", err)
	}

	// Create notes directory
	originalFilename := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	destinationDir, err := fileService.CreateNoteDirectory(originalFilename + "_summary")
	if err != nil {
		return err
	}

	// Save the summary
	summaryPath := filepath.Join(destinationDir, "summary.md")
	transcriptionPath := filepath.Join(destinationDir, "transcription.md") // Empty for text files
	if err := fileService.SaveMarkdownFiles("", summary.(string), transcriptionPath, summaryPath); err != nil {
		return fmt.Errorf("failed to save summary file: %w", err)
	}

	fmt.Printf("✅ Text file '%s' successfully processed and summarized:\n", filePath)
	fmt.Printf("   📝 Summary saved at: %s\n", summaryPath)

	return nil
}

func processAudioFile(filePath string, audioService *services.AudioService, fileService *services.FileService, uiService *services.UIService) error {
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
	openaiService, err := services.NewOpenAIService()
	if err != nil {
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

	// Summarize the transcript
	summary, err := uiService.RunTaskWithSpinner("📄 Creating summary from transcription", func() (interface{}, error) {
		return openaiService.SummarizeText(transcript)
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

	// Create a metadata note that references the files
	title := fmt.Sprintf("Audio Import: %s", originalFilename)
	folderName := filepath.Base(destinationDir)
	content := fmt.Sprintf("Audio file imported and processed.\n\nFiles:\n- Audio: %s\n- Transcription: %s\n- Summary: %s\n\nFolder: %s",
		filename, "transcription.md", "summary.md", folderName)
	tags := "imported,audio,metadata"

	_, err = uiService.RunTaskWithSpinner("📋 Creating metadata note", func() (interface{}, error) {
		_, createErr := database.CreateNote(db, title, content, tags)
		return nil, createErr
	})
	if err != nil {
		return fmt.Errorf("failed to save metadata note: %w", err)
	}

	fmt.Printf("✅ File '%s' successfully imported and processed:\n", filename)
	fmt.Printf("   📁 Saved to notes directory: %s\n", destinationDir)
	fmt.Printf("   🎵 Added to recordings database\n")
	fmt.Printf("   📝 Transcribed and summarized\n")
	fmt.Printf("   📋 Note created: %s\n", title)
	return nil
}
