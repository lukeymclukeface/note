package cmd

import (
	"fmt"
	"note-cli/internal/helpers"
	"note-cli/internal/services"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var transcribeCmd = &cobra.Command{
	Use:   "transcribe [file]",
	Short: "Transcribe an audio file using AI services",
	Long: `Transcribe a provided audio file into text using AI services configured in the application.
	Supported formats are mp3, wav, m4a, ogg, flac.`,
	Args:  cobra.ExactArgs(1),
	RunE:  transcribeAudio,
}

func init() {
	rootCmd.AddCommand(transcribeCmd)
	
	// Add command flags
	transcribeCmd.Flags().StringP("output", "o", "", "Output file path for the transcription (default: print to stdout)")
	transcribeCmd.Flags().StringP("format", "f", "text", "Output format: text, markdown (default: text)")
	transcribeCmd.Flags().StringP("dir", "d", "", "Output directory for chunk files (optional)")
	transcribeCmd.Flags().BoolP("chunks", "c", false, "Save individual chunk transcriptions to separate files")
}

func transcribeAudio(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	verboseLogger := services.NewVerboseLogger(IsVerbose())
	verboseLogger.StartCommand("transcribe", args)
	start := time.Now()
	defer func() {
		verboseLogger.EndCommand("transcribe", time.Since(start), true)
	}()

	// Get command line flags
	outputFile, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	outputDir, _ := cmd.Flags().GetString("dir")
	saveChunks, _ := cmd.Flags().GetBool("chunks")

	// Validate format
	if format != "text" && format != "markdown" {
		return fmt.Errorf("invalid format: %s. Must be 'text' or 'markdown'", format)
	}

	fmt.Printf("🎵 Transcribing audio file: %s\n", filePath)

	// Determine output directory for chunks
	chunkDir := ""
	if saveChunks {
		if outputDir != "" {
			chunkDir = outputDir
		} else if outputFile != "" {
			// Use the directory of the output file
			chunkDir = filepath.Dir(outputFile)
		} else {
			// Use current directory
			chunkDir = "."
		}
		// Ensure the directory exists
		if err := os.MkdirAll(chunkDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Use the helper function to transcribe the audio file
	result, err := helpers.TranscribeAudioFile(filePath, chunkDir, verboseLogger)
	if err != nil {
		return fmt.Errorf("transcription failed: %w", err)
	}

	// Format the transcript based on the requested format
	var formattedTranscript string
	if format == "markdown" {
		baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		formattedTranscript = fmt.Sprintf("# Transcription: %s\n\n**Source File:** %s\n**Generated:** %s\n\n## Content\n\n%s\n", 
			baseName, filePath, time.Now().Format("2006-01-02 15:04:05"), result.FullTranscript)
	} else {
		formattedTranscript = result.FullTranscript
	}

	// Output the transcription
	if outputFile != "" {
		// Save to file
		if err := os.WriteFile(outputFile, []byte(formattedTranscript), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("✅ Transcription saved to: %s\n", outputFile)
		if saveChunks && len(result.ChunkFiles) > 0 {
			fmt.Printf("📁 Chunk files saved to: %s\n", chunkDir)
			fmt.Printf("📝 Generated %d chunk files\n", len(result.ChunkFiles))
		}
	} else {
		// Print to stdout
		fmt.Println("\n📝 Transcription Result:")
		fmt.Println("========================================")
		fmt.Println(formattedTranscript)
		fmt.Println("========================================")
		if saveChunks && len(result.ChunkFiles) > 0 {
			fmt.Printf("\n📁 Chunk files saved to: %s\n", chunkDir)
			fmt.Printf("📝 Generated %d chunk files\n", len(result.ChunkFiles))
		}
	}

	return nil
}

