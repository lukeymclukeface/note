package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"os/exec"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the note application",
	Long:  `Check and configure the dependencies required for the note application, including ffmpeg and OpenAI API key.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up note...")
		fmt.Println("==================")
		
		ffmpegOK := checkFFmpeg()
		openaiOK := checkOpenAIKey()
		
		fmt.Println()
		if ffmpegOK && openaiOK {
			fmt.Println("✅ Setup complete! All dependencies are configured.")
		} else {
			fmt.Println("⚠️  Setup incomplete. Please resolve the issues above.")
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func checkFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("❌ ffmpeg is not installed.")
		fmt.Println("   Install with: brew install ffmpeg")
		return false
	} else {
		fmt.Println("✅ ffmpeg is installed.")
		return true
	}
}

func checkOpenAIKey() bool {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return false
	}

	if cfg.OpenAIKey == "" {
		fmt.Println("❌ OpenAI API key is missing.")
		fmt.Println("   Set with: note config set openai_key <your-key>")
		return false
	} else {
		fmt.Println("✅ OpenAI API key is configured.")
		return true
	}
}
