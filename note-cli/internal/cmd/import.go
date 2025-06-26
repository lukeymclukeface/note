package cmd

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"note-cli/internal/config"
	"note-cli/internal/constants"
	"note-cli/internal/database"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import and transcribe an audio file",
	Long: `Import an audio file, convert it to MP3 if needed, add it to the recordings database, 
transcribe it using OpenAI's Whisper API, create a summary, and save as a note.

If no file is specified, you'll be presented with a list of compatible audio files 
in the current directory to choose from.

Supported audio formats: mp3, wav, m4a, ogg, flac

Requires:
- OpenAI API key configured (run 'note setup')
- ffmpeg installed for format conversion`,
Args: cobra.MaximumNArgs(1),
	RunE: importFile,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importFile(cmd *cobra.Command, args []string) error {
var filePath string

	if len(args) == 1 {
		filePath = args[0]
	} else {
		// No argument provided, find compatible files in the current directory
		var files []string
		extensions := []string{"*.mp3", "*.wav", "*.m4a", "*.ogg", "*.flac"}
		
		for _, ext := range extensions {
			matches, err := filepath.Glob(ext)
			if err != nil {
				continue
			}
			files = append(files, matches...)
		}
		
		if len(files) == 0 {
			return fmt.Errorf("no compatible audio files found in the current directory")
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

	// Validate file
	if !isValidAudioFile(filePath) {
		return fmt.Errorf("invalid audio file format. Supported formats: mp3, wav, m4a, ogg, flac. Got: %s", filepath.Ext(filePath))
	}

	fmt.Println("🎵 Processing audio file...")

	// Convert to MP3 if necessary
	mp3Path, err := convertToMP3WithSpinner(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert file to MP3: %w", err)
	}

	// Load config for database path
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if cfg.DatabasePath == "" {
		return fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	if cfg.OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not configured. Please run 'note setup' first")
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found. Please install with 'brew install ffmpeg' or run 'note setup'")
	}

	// Create a unique folder under notes directory
	notesDir, err := constants.GetNotesDir()
	if err != nil {
		return fmt.Errorf("failed to get notes directory: %w", err)
	}

	// Create a unique folder for this import
	originalFilename := strings.TrimSuffix(filepath.Base(mp3Path), filepath.Ext(mp3Path))
	folderName := fmt.Sprintf("%s_%s", originalFilename, time.Now().Format("20060102_150405"))
	destinationDir := filepath.Join(notesDir, folderName)

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy MP3 to the new directory
	filename := filepath.Base(mp3Path)
	newFilePath := filepath.Join(destinationDir, filename)
	if err := copyFileWithSpinner(mp3Path, newFilePath); err != nil {
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

	if err := saveRecordingWithSpinner(cfg.DatabasePath, &recording); err != nil {
		return fmt.Errorf("failed to save recording: %w", err)
	}

	// Transcribe and summarize
	transcript, err := transcribeFileWithSpinner(newFilePath, cfg.OpenAIKey)
	if err != nil {
		return fmt.Errorf("failed to transcribe file: %w", err)
	}

	summary, err := summarizeTextWithSpinner(transcript, cfg.OpenAIKey)
	if err != nil {
		return fmt.Errorf("failed to summarize transcription: %w", err)
	}

	// Save transcription and summary as separate markdown files
	transcriptionPath := filepath.Join(destinationDir, "transcription.md")
	summaryPath := filepath.Join(destinationDir, "summary.md")

	if err := saveMarkdownFilesWithSpinner(transcript, summary, transcriptionPath, summaryPath); err != nil {
		return fmt.Errorf("failed to save markdown files: %w", err)
	}

	// Connect to database for note creation with metadata only
	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Create a metadata note that references the files
	title := fmt.Sprintf("Audio Import: %s", originalFilename)
	content := fmt.Sprintf("Audio file imported and processed.\n\nFiles:\n- Audio: %s\n- Transcription: %s\n- Summary: %s\n\nFolder: %s",
		filename, "transcription.md", "summary.md", folderName)
	tags := "imported,audio,metadata"

	if err := createNoteWithSpinner(db, title, content, tags); err != nil {
		return fmt.Errorf("failed to save metadata note: %w", err)
	}

	fmt.Printf("✅ File '%s' successfully imported and processed:\n", filename)
	fmt.Printf("   📁 Saved to recordings directory\n")
	fmt.Printf("   🎵 Added to recordings database\n")
	fmt.Printf("   📝 Transcribed and summarized\n")
	fmt.Printf("   📋 Note created: %s\n", title)
	return nil
}

func isValidAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".wav" || ext == ".m4a" || ext == ".ogg" || ext == ".flac"
}

func convertToMP3(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".mp3" {
		return filePath, nil // Already an MP3
	}

	outputPath := strings.TrimSuffix(filePath, ext) + ".mp3"
	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-acodec", "libmp3lame",
		"-ab", "128k",
		"-ar", "44100", // Standardize sample rate
		"-ac", "1", // Convert to mono
		"-y", // Overwrite output file
		outputPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("FFmpeg error output: %s\n", stderr.String())
		return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return outputPath, nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return err
	}

	return nil
}

func transcribeFile(filePath, apiKey string) (string, error) {
	// Load config to get the transcription model
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Use configured transcription model or fallback to default
	transcriptionModel := cfg.TranscriptionModel
	if transcriptionModel == "" {
		transcriptionModel = "whisper-1"
	}

	url := "https://api.openai.com/v1/audio/transcriptions"
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("model", transcriptionModel)
	if err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	transcript, ok := result["text"].(string)
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	return transcript, nil
}

func summarizeText(text, apiKey string) (string, error) {
	// Load config to get the summary model
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// Use configured summary model or fallback to default
	model := cfg.SummaryModel
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	url := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "You are a helpful assistant that creates concise summaries of transcribed audio content.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Please provide a concise summary of the following transcribed audio:\n\n%s", text),
			},
		},
		"max_tokens":  300,
		"temperature": 0.7,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("unexpected response structure")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("unexpected response structure")
	}

	return strings.TrimSpace(content), nil
}

// Function to run a task with a simple spinner
func runTaskWithSpinner(message string, task func() (interface{}, error)) (interface{}, error) {
	// Channel to signal completion
	done := make(chan struct{})
	var result interface{}
	var taskErr error

	// Start the task in a goroutine
	go func() {
		defer close(done)
		result, taskErr = task()
	}()

	// Spinner characters
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIndex := 0

	// Show spinner while task is running
	for {
		select {
		case <-done:
			// Clear the spinner line
			fmt.Print("\r\033[K")
			if taskErr != nil {
				fmt.Printf("❌ %s failed: %v\n", message, taskErr)
				return nil, taskErr
			}
			fmt.Printf("✅ %s completed!\n", message)
			return result, nil
		default:
			// Update spinner
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
			fmt.Printf("\r%s %s", spinnerChars[spinnerIndex], style.Render(message))
			spinnerIndex = (spinnerIndex + 1) % len(spinnerChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Wrapper functions for transcription and summarization with spinner
func transcribeFileWithSpinner(filePath, apiKey string) (string, error) {
	task := func() (interface{}, error) {
		return transcribeFile(filePath, apiKey)
	}

	result, err := runTaskWithSpinner("🎙️  Transcribing audio using OpenAI", task)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func summarizeTextWithSpinner(text, apiKey string) (string, error) {
	task := func() (interface{}, error) {
		return summarizeText(text, apiKey)
	}

	result, err := runTaskWithSpinner("📝 Creating summary using OpenAI", task)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

// Wrapper function for MP3 conversion with spinner
func convertToMP3WithSpinner(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".mp3" {
		return filePath, nil // Already an MP3, no conversion needed
	}

	task := func() (interface{}, error) {
		return convertToMP3(filePath)
	}

	result, err := runTaskWithSpinner(fmt.Sprintf("🔄 Converting %s to MP3", ext), task)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

// Wrapper function for file copying with spinner
func copyFileWithSpinner(src, dst string) error {
	task := func() (interface{}, error) {
		err := copyFile(src, dst)
		return nil, err
	}

	_, err := runTaskWithSpinner("📁 Copying MP3 file to notes directory", task)
	return err
}

// Wrapper function for saving recording with spinner
func saveRecordingWithSpinner(dbPath string, recording *database.Recording) error {
	task := func() (interface{}, error) {
		err := database.SaveRecording(dbPath, recording)
		return nil, err
	}

	_, err := runTaskWithSpinner("💾 Saving recording metadata to database", task)
	return err
}

// Wrapper function for saving markdown files with spinner
func saveMarkdownFilesWithSpinner(transcript, summary, transcriptionPath, summaryPath string) error {
	task := func() (interface{}, error) {
		// Write transcription file
		transcriptionContent := fmt.Sprintf("# Transcription\n\n%s\n", transcript)
		if err := os.WriteFile(transcriptionPath, []byte(transcriptionContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write transcription file: %w", err)
		}

		// Write summary file
		summaryContent := fmt.Sprintf("# Summary\n\n%s\n", summary)
		if err := os.WriteFile(summaryPath, []byte(summaryContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write summary file: %w", err)
		}

		return nil, nil
	}

	_, err := runTaskWithSpinner("📝 Saving transcription and summary as markdown files", task)
	return err
}

// Wrapper function for creating note with spinner
func createNoteWithSpinner(db *sql.DB, title, content, tags string) error {
	task := func() (interface{}, error) {
		_, err := database.CreateNote(db, title, content, tags)
		return nil, err
	}

	_, err := runTaskWithSpinner("📋 Saving metadata note to database", task)
	return err
}
