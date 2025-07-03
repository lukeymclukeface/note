# Note

[![Test CLI](https://github.com/lukeymclukeface/note/workflows/Test%20CLI/badge.svg)](https://github.com/lukeymclukeface/note/actions/workflows/test-cli.yml)
[![Test Web](https://github.com/lukeymclukeface/note/workflows/Test%20Web/badge.svg)](https://github.com/lukeymclukeface/note/actions/workflows/test-web.yml)
[![Test Server](https://github.com/lukeymclukeface/note/workflows/Test%20Server/badge.svg)](https://github.com/lukeymclukeface/note/actions/workflows/test-server.yml)

A comprehensive AI-powered note-taking system that combines audio recording, transcription, and intelligent summarization capabilities.

## Project Structure

```
note/
├── note-cli/                 # Command-line interface
│   ├── cmd/
│   ├── internal/
│   └── ...
├── note-server/              # Go backend API server
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── ...
├── note-web/                 # Next.js web application
│   ├── src/
│   ├── public/
│   └── ...
├── Makefile                  # Build and development tasks
├── docker-compose.yml        # Docker orchestration
├── docker-compose.dev.yml    # Development environment
└── README.md                 # This file
```

This repository contains three main components:

- **[note-cli](./note-cli/)** - Command-line interface for recording, importing, and managing notes
- **[note-web](./note-web/)** - Next.js web application for browsing and managing notes
- **[note-server](./note-server/)** - Go-based backend API server for real-time functionality

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

## Server Component (note-server)

- Go-based backend API
- WebSocket for real-time updates
- Audio transcription & text summarization (OpenAI, Whisper)
- REST endpoints (`/api/notes`, `/api/transcribe`, etc.)
- Docker-based deployment, health-check on `/healthz`
- Link: "See full docs → note-server/README.md"

## Quick Start

### Prerequisites
- Go 1.23+ (for CLI)
- Node.js 18+ (for web interface)
- FFmpeg (auto-installed via Homebrew on macOS)
- OpenAI API key

### CLI Setup
1. Clone the repository
2. Set up the CLI: `cd note-cli && go build -o note cmd/note/main.go`
3. Run setup: `./note setup`
4. Configure OpenAI models: `./note config model`
5. Start recording or importing: `./note record` or `./note import`

## Running the Application

### Docker Compose (production)
For a complete production deployment with all services:
```bash
make quick-start
# or
make up
```

### Docker Compose (development with hot reload)
For development with hot reload capabilities:
```bash
make quick-dev
# or
make dev
```

## Development Commands

### Individual Service Development
To run individual services for development:

#### Backend Server Only
```bash
make dev-server
```

#### Web Frontend Only
```bash
make dev-web
```

#### Manual CLI Build
```bash
cd note-cli
go build -o note cmd/note/main.go
./note setup
./note config model
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
