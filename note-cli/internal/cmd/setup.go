package cmd

import (
	"fmt"
	"note-cli/internal/config"
	"note-cli/internal/database"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the note application",
	Long:  `Interactive setup for the note application. Checks for required dependencies and configuration, offering to install missing packages and configure settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔧 Setting up note...")
		fmt.Println("====================")
		fmt.Println()

		// Run interactive setup
		runInteractiveSetup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}


func runInteractiveSetup() {
	// Check brew installation first
	brewInstalled := checkBrew()

	// Check dependencies
	ffmpegInstalled := checkFFmpegInstalled()

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return
	}

	// Handle ffmpeg installation
	if !ffmpegInstalled {
		if brewInstalled {
			var installFFmpeg bool
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("FFmpeg is required but not installed.").
						Description("Would you like me to install it using Homebrew?").
						Value(&installFFmpeg),
				),
			)

			if err := form.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if installFFmpeg {
				fmt.Println("📦 Installing ffmpeg...")
				if err := installPackageWithBrew("ffmpeg"); err != nil {
					fmt.Printf("❌ Failed to install ffmpeg: %v\n", err)
					fmt.Println("Please install manually with: brew install ffmpeg")
				} else {
					fmt.Println("✅ FFmpeg installed successfully!")
				}
			} else {
				fmt.Println("⚠️  FFmpeg not installed. You can install it later with: brew install ffmpeg")
			}
		} else {
			fmt.Println("❌ FFmpeg is required but not installed.")
			fmt.Println("   Please install Homebrew first, then run: brew install ffmpeg")
		}
	} else {
		fmt.Println("✅ FFmpeg is already installed.")
	}

	// Handle missing configuration
	if cfg.OpenAIKey == "" {
		fmt.Println()

		var provideKey bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("OpenAI API key is required for AI features.").
					Description("Would you like to provide your OpenAI API key now?").
					Value(&provideKey),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if provideKey {
			var apiKey string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter your OpenAI API key:").
						Description("You can get this from https://platform.openai.com/api-keys").
						Value(&apiKey).
						EchoMode(huh.EchoModePassword).
						Validate(func(str string) error {
							if strings.TrimSpace(str) == "" {
								return fmt.Errorf("API key cannot be empty")
							}
							if !strings.HasPrefix(str, "sk-") {
								return fmt.Errorf("OpenAI API keys should start with 'sk-'")
							}
							return nil
						}),
				),
			)

			if err := form.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Save the API key
			cfg.OpenAIKey = strings.TrimSpace(apiKey)
			if err := config.Save(cfg); err != nil {
				fmt.Printf("❌ Error saving configuration: %v\n", err)
				return
			}

			fmt.Println("✅ OpenAI API key saved successfully!")
		} else {
			fmt.Println("⚠️  OpenAI API key not configured. You can set it later with: note config set openai_key <your-key>")
		}
	} else {
		fmt.Println("✅ OpenAI API key is already configured.")
	}

	// Handle database setup
	if cfg.DatabasePath == "" {
		fmt.Println()

		var setupDb bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Database is not configured.").
					Description("Would you like to set up the database now?").
					Value(&setupDb),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if setupDb {
			// Use default database path
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("❌ Error getting home directory: %v\n", err)
				return
			}

			dbPath := homeDir + "/.note-cli/notes.db"

			// Run database setup
			if err := setupDatabase(dbPath); err != nil {
				fmt.Printf("❌ Error setting up database: %v\n", err)
				return
			}

			// Save database path to config
			cfg.DatabasePath = dbPath
			if err := config.Save(cfg); err != nil {
				fmt.Printf("❌ Error saving configuration: %v\n", err)
				return
			}

			fmt.Println("✅ Database setup completed!")
		} else {
			fmt.Println("⚠️  Database not configured. You can set it up later with: note setup")
		}
	} else {
		fmt.Println("✅ Database is already configured.")
	}

	// Final status
	fmt.Println()

	// Re-check everything after setup attempts
	ffmpegInstalled = checkFFmpegInstalled()
	cfg, err = config.Load()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return
	}

	if ffmpegInstalled && cfg.OpenAIKey != "" && cfg.DatabasePath != "" {
		fmt.Println("🎉 Setup complete! All dependencies and configuration are ready.")
	} else {
		fmt.Println("⚠️  Setup completed with some items remaining. Run 'note setup' again if needed.")
	}
}

func checkBrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func checkFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func installPackageWithBrew(packageName string) error {
	cmd := exec.Command("brew", "install", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func setupDatabase(dbPath string) error {
	fmt.Printf("🗄️  Setting up database at %s...\n", dbPath)

	// Initialize the database
	if err := database.Initialize(dbPath); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	fmt.Println("✅ Database initialized successfully!")
	return nil
}
