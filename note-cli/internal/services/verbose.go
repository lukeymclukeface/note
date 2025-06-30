package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

// VerboseLogger handles verbose output throughout the application
type VerboseLogger struct {
	enabled bool
}

// NewVerboseLogger creates a new verbose logger instance
func NewVerboseLogger(enabled bool) *VerboseLogger {
	return &VerboseLogger{enabled: enabled}
}

// IsEnabled returns whether verbose logging is enabled
func (v *VerboseLogger) IsEnabled() bool {
	return v.enabled
}

// Info logs general information
func (v *VerboseLogger) Info(format string, args ...interface{}) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	cyan := color.New(color.FgCyan)
	cyan.Printf("[%s] INFO: %s\n", timestamp, fmt.Sprintf(format, args...))
}

// Debug logs debug information  
func (v *VerboseLogger) Debug(format string, args ...interface{}) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	blue := color.New(color.FgBlue)
	blue.Printf("[%s] DEBUG: %s\n", timestamp, fmt.Sprintf(format, args...))
}

// API logs API request/response information
func (v *VerboseLogger) API(method, url string, requestBody interface{}, responseBody interface{}, statusCode int, duration time.Duration) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	magenta := color.New(color.FgMagenta)
	
	magenta.Printf("[%s] API: %s %s\n", timestamp, method, url)
	
	if requestBody != nil {
		v.logJSON("Request Body", requestBody)
	}
	
	yellow := color.New(color.FgYellow)
	yellow.Printf("         Status: %d, Duration: %v\n", statusCode, duration)
	
	if responseBody != nil {
		v.logJSON("Response Body", responseBody)
	}
	fmt.Println()
}

// File logs file operation information
func (v *VerboseLogger) File(operation, filePath string, size int64) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	green := color.New(color.FgGreen)
	
	absPath, _ := filepath.Abs(filePath)
	sizeStr := ""
	if size > 0 {
		sizeStr = fmt.Sprintf(" (%s)", formatFileSize(size))
	}
	
	green.Printf("[%s] FILE: %s - %s%s\n", timestamp, operation, absPath, sizeStr)
}

// Database logs database operation information
func (v *VerboseLogger) Database(operation, details string) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	red := color.New(color.FgRed)
	red.Printf("[%s] DB: %s - %s\n", timestamp, operation, details)
}

// Config logs configuration information
func (v *VerboseLogger) Config(key, value string) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	white := color.New(color.FgWhite)
	
	// Mask sensitive values
	displayValue := value
	if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "secret") {
		if len(value) > 8 {
			displayValue = value[:4] + "..." + value[len(value)-4:]
		} else {
			displayValue = "***"
		}
	}
	
	white.Printf("[%s] CONFIG: %s = %s\n", timestamp, key, displayValue)
}

// Process logs process/command execution information
func (v *VerboseLogger) Process(command string, args []string, duration time.Duration) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	cyan := color.New(color.FgCyan)
	
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}
	
	cyan.Printf("[%s] PROCESS: %s (took %v)\n", timestamp, fullCommand, duration)
}

// Error logs error information
func (v *VerboseLogger) Error(err error, context string) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	red := color.New(color.FgRed, color.Bold)
	red.Printf("[%s] ERROR: %s - %v\n", timestamp, context, err)
}

// Step logs step information for command execution
func (v *VerboseLogger) Step(stepName, details string) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("[%s] STEP: %s\n", timestamp, stepName)
	if details != "" {
		fmt.Printf("         %s\n", details)
	}
}

// StartCommand logs the beginning of a command execution
func (v *VerboseLogger) StartCommand(commandName string, args []string) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	bold := color.New(color.Bold, color.FgWhite)
	bold.Printf("\n[%s] ========== COMMAND START: %s ==========\n", timestamp, commandName)
	
	if len(args) > 0 {
		fmt.Printf("         Arguments: %s\n", strings.Join(args, " "))
	}
	
	// Show working directory
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("         Working Directory: %s\n", wd)
	}
	fmt.Println()
}

// EndCommand logs the end of a command execution
func (v *VerboseLogger) EndCommand(commandName string, duration time.Duration, success bool) {
	if !v.enabled {
		return
	}
	timestamp := time.Now().Format("15:04:05.000")
	bold := color.New(color.Bold, color.FgWhite)
	
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	
	bold.Printf("\n[%s] ========== COMMAND END: %s (%s) - %v ==========\n\n", timestamp, commandName, status, duration)
}

// logJSON logs JSON data with proper formatting
func (v *VerboseLogger) logJSON(label string, data interface{}) {
	if !v.enabled {
		return
	}
	
	var jsonBytes []byte
	var err error
	
	switch d := data.(type) {
	case string:
		jsonBytes = []byte(d)
	case []byte:
		jsonBytes = d
	default:
		jsonBytes, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("         %s: %v (failed to marshal JSON: %v)\n", label, data, err)
			return
		}
	}
	
	// Try to format as JSON if it's valid JSON
	var formatted interface{}
	if err := json.Unmarshal(jsonBytes, &formatted); err == nil {
		if prettyBytes, err := json.MarshalIndent(formatted, "         ", "  "); err == nil {
			gray := color.New(color.FgHiBlack)
			gray.Printf("         %s:\n%s\n", label, string(prettyBytes))
			return
		}
	}
	
	// Fallback to raw output
	gray := color.New(color.FgHiBlack)
	gray.Printf("         %s: %s\n", label, string(jsonBytes))
}

// formatFileSize formats file size in human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
