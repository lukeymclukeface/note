# Note

A comprehensive AI-powered note-taking system that combines audio recording, transcription, and intelligent summarization capabilities.

## Project Structure

This repository contains two main components:

- **[note-cli](./note-cli/)** - Command-line interface for recording, importing, and managing notes
- **[note-web](./note-web/)** - Next.js web application for browsing and managing notes

## Features

### Audio Processing
- Record audio directly from your microphone
- Import existing audio files (mp3, wav, m4a, ogg, flac)
- Automatic chunked transcription for large files
- Speaker diarization through AI post-processing

### AI Integration
- OpenAI Whisper for transcription
- GPT models for intelligent summarization
- Content type detection (meetings, interviews, lectures, etc.)
- Specialized summarization prompts based on content type
- Configurable models for transcription and summarization

### Content Management
- SQLite database for metadata storage
- Markdown and text file import support
- Interactive note browsing and management
- Bulk operations and search capabilities

### User Experience
- Beautiful terminal UI with Bubble Tea
- Interactive prompts and selections
- Progress spinners for long operations
- Comprehensive verbose logging
- Cross-platform compatibility (macOS focus)

## Quick Start

### Prerequisites
- Go 1.23+ (for CLI)
- Node.js 18+ (for web interface)
- FFmpeg (auto-installed via Homebrew on macOS)
- OpenAI API key

### Setup
1. Clone the repository
2. Set up the CLI: `cd note-cli && go build -o note cmd/note/main.go`
3. Run setup: `./note setup`
4. Configure OpenAI models: `./note config model`
5. Start recording or importing: `./note record` or `./note import`

### Web Interface
```bash
cd note-web
npm install
npm run dev
```

## Configuration

The application stores configuration and data in `~/.noteai/`:
- `config.json` - Application settings and API keys
- `notes.db` - SQLite database
- `recordings/` - Audio recordings
- `notes/` - Generated note files

## Contributing

This is a personal project focused on AI-powered note-taking workflows. The codebase follows modern Go and Next.js best practices with a service-oriented architecture.

## License

Private project - All rights reserved.
