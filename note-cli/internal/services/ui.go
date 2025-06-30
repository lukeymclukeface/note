package services

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// UIService handles common UI operations
type UIService struct{}

// NewUIService creates a new UI service instance
func NewUIService() *UIService {
	return &UIService{}
}

// RunTaskWithSpinner executes a task while showing a spinner
func (s *UIService) RunTaskWithSpinner(message string, task func() (interface{}, error)) (interface{}, error) {
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

// FormatFileSize formats a file size in bytes to a human-readable string
func (s *UIService) FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
