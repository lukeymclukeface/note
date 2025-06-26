package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/constants"
	"note-cli/internal/database"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record audio from microphone",
	Long:  `Record audio from the default microphone using ffmpeg. Press Enter to stop recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := startRecording(); err != nil {
			fmt.Printf("❌ Error during recording: %v\n", err)
			os.Exit(1)
		}
	},
}

var selectDevice bool

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().BoolVar(&selectDevice, "select-device", false, "Force device selection instead of auto-selecting MacBook microphone")
}

func startRecording() error {
	// Load config to get database path
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.DatabasePath == "" {
		return fmt.Errorf("database not configured. Please run 'note setup' first")
	}

	// Check if ffmpeg is available
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found. Please install with 'brew install ffmpeg' or run 'note setup'")
	}

	// Get available audio devices
	microphoneIndex, err := selectMicrophone()
	if err != nil {
		return fmt.Errorf("failed to select microphone: %w", err)
	}

	// Create recordings directory
	recordingsDir, err := constants.GetRecordingsDir()
	if err != nil {
		return fmt.Errorf("failed to get recordings directory: %w", err)
	}
	if err := os.MkdirAll(recordingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create recordings directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("recording_%s.mp3", timestamp)
	filepath := filepath.Join(recordingsDir, filename)

	fmt.Println("🎤 Starting audio recording...")
	fmt.Println("📁 Saving to:", filepath)
	fmt.Println("⏹️  Press Enter to stop recording")
	fmt.Println()

	// Start recording
	startTime := time.Now()
	inputDevice := fmt.Sprintf(":%d", microphoneIndex)
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation", // Use AVFoundation on macOS
		"-i", inputDevice, // Use selected microphone
		"-ac", "1", // Mono audio
		"-ar", "44100", // Sample rate
		"-acodec", "libmp3lame", // MP3 audio codec (LAME encoder)
		"-ab", "128k", // Audio bitrate for compression
		"-y", // Overwrite output file if exists
		filepath,
	)

	// Capture ffmpeg output for debugging and provide stdin
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Create stdin pipe to send 'q' command to ffmpeg
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	fmt.Printf("🔊 Recording started (PID: %d)\n", cmd.Process.Pid)

	// Wait for user input in a goroutine
	done := make(chan bool)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadLine()
		done <- true
	}()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-done:
		fmt.Println("\n⏹️  Stopping recording...")
	case <-sigChan:
		fmt.Println("\n⏹️  Recording interrupted...")
	}

	// Stop ffmpeg gracefully by sending 'q' command
	if _, err := stdin.Write([]byte("q\n")); err != nil {
		fmt.Printf("Warning: Failed to send quit command to ffmpeg: %v\n", err)
		// Fallback to kill
		if killErr := cmd.Process.Kill(); killErr != nil {
			fmt.Printf("Warning: Failed to kill ffmpeg process: %v\n", killErr)
		}
	}
	stdin.Close()

	// Wait for process to finish and check for errors
	if err := cmd.Wait(); err != nil {
		// Check if there's useful error information in stderr
		if stderr.Len() > 0 {
			fmt.Printf("FFmpeg error output: %s\n", stderr.String())
		}
		fmt.Printf("Warning: ffmpeg process ended with error: %v\n", err)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Give ffmpeg a moment to finalize the file
	time.Sleep(100 * time.Millisecond)

	// Check if file was created and get its size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		// Print ffmpeg errors if available
		if stderr.Len() > 0 {
			fmt.Printf("FFmpeg error output: %s\n", stderr.String())
		}
		return fmt.Errorf("recording file not found: %w", err)
	}

	// Save recording to database
	recording := database.Recording{
		Filename:   filename,
		FilePath:   filepath,
		StartTime:  startTime,
		EndTime:    endTime,
		Duration:   duration,
		FileSize:   fileInfo.Size(),
		Format:     "mp3",
		SampleRate: 44100, // Sample rate
		Channels:   1,     // Mono recording
		CreatedAt:  time.Now(),
	}

	if err := database.SaveRecording(cfg.DatabasePath, &recording); err != nil {
		fmt.Printf("⚠️  Recording saved but failed to update database: %v\n", err)
	} else {
		fmt.Println("✅ Recording saved to database")
	}

	fmt.Printf("🎵 Recording completed!\n")
	fmt.Printf("   Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("   Size: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("   File: %s\n", filepath)

	return nil
}

type AudioDevice struct {
	Index int
	Name  string
}

func selectMicrophone() (int, error) {
	// Get list of audio devices from ffmpeg
	deviceCmd := exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")
	var deviceOutput bytes.Buffer
	deviceCmd.Stderr = &deviceOutput
	deviceCmd.Run() // This always "fails" but gives us device list

	if deviceOutput.Len() == 0 {
		return 0, fmt.Errorf("no audio devices found")
	}

	// Parse audio devices from ffmpeg output
	audioDevices, err := parseAudioDevices(deviceOutput.String())
	if err != nil {
		return 0, fmt.Errorf("failed to parse audio devices: %w", err)
	}

	if len(audioDevices) == 0 {
		return 0, fmt.Errorf("no audio input devices found")
	}

	// Look for MacBook microphone automatically (unless --select-device flag is used)
	if !selectDevice {
		for _, device := range audioDevices {
			if strings.Contains(strings.ToLower(device.Name), "macbook") {
				fmt.Printf("🎤 Auto-selected: %s\n", device.Name)
				return device.Index, nil
			}
		}
	}

	// If no MacBook microphone found, let user choose
	fmt.Println("🎤 Available microphones:")
	for _, device := range audioDevices {
		fmt.Printf("  [%d] %s\n", device.Index, device.Name)
	}
	fmt.Println()

	// Create options for the form
	options := make([]huh.Option[int], len(audioDevices))
	for i, device := range audioDevices {
		options[i] = huh.NewOption(device.Name, device.Index)
	}

	var selectedIndex int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a microphone to use for recording:").
				Options(options...).
				Value(&selectedIndex),
		),
	)

	if err := form.Run(); err != nil {
		return 0, fmt.Errorf("failed to get user selection: %w", err)
	}

	return selectedIndex, nil
}

func parseAudioDevices(output string) ([]AudioDevice, error) {
	var devices []AudioDevice

	// Look for the audio devices section
	lines := strings.Split(output, "\n")
	inAudioSection := false

	for _, line := range lines {
		// Start of audio devices section
		if strings.Contains(line, "AVFoundation audio devices:") {
			inAudioSection = true
			continue
		}

		// End of section or start of next section
		if inAudioSection && (strings.Contains(line, "Error opening input") ||
			strings.TrimSpace(line) == "" && !strings.Contains(line, "[AVFoundation")) {
			break
		}

		// Parse device lines like: "[AVFoundation indev @ 0x...] [0] MacBook Pro Microphone"
		if inAudioSection && strings.Contains(line, "[AVFoundation indev") {
			// Use regex to extract device index and name
			re := regexp.MustCompile(`\[(\d+)\]\s+(.+)$`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				index, err := strconv.Atoi(matches[1])
				if err != nil {
					continue
				}
				name := strings.TrimSpace(matches[2])
				devices = append(devices, AudioDevice{
					Index: index,
					Name:  name,
				})
			}
		}
	}

	return devices, nil
}
