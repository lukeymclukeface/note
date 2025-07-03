# Note

[![Test Web](https://github.com/lukeymclukeface/note/workflows/Test%20Web/badge.svg)](https://github.com/lukeymclukeface/note/actions/workflows/test-web.yml)
[![Test Server](https://github.com/lukeymclukeface/note/workflows/Test%20Server/badge.svg)](https://github.com/lukeymclukeface/note/actions/workflows/test-server.yml)

A comprehensive AI-powered note-taking system with web interface that provides audio recording, transcription, and intelligent summarization capabilities.

## Project Structure

```
note/
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

This repository contains two main components:

- **[note-web](./note-web/)** - Next.js web application for browsing and managing notes
- **[note-server](./note-server/)** - Go-based backend API server for real-time functionality

## Features

### Audio Processing
- Record audio directly through the web interface
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
- Web-based note browsing and management
- Bulk operations and search capabilities

### User Experience
- Modern web interface built with Next.js
- Real-time updates via WebSocket
- Responsive design for desktop and mobile
- Intuitive file upload and management
- Cross-platform web accessibility

## Server Component (note-server)

- Go-based backend API
- WebSocket for real-time updates
- Audio transcription & text summarization (OpenAI, Whisper)
- REST endpoints (`/api/notes`, `/api/transcribe`, etc.)
- Docker-based deployment, health-check on `/healthz`
- Link: "See full docs → note-server/README.md"

## Quick Start

### Prerequisites
- Docker and Docker Compose (for containerized deployment)
- Or for local development:
  - Go 1.23+ (for backend server)
  - Node.js 18+ (for web interface)
  - FFmpeg (for audio processing)
- OpenAI API key

### Web Application Setup
1. Clone the repository
2. Set up environment variables with your OpenAI API key
3. Use Docker Compose for quick deployment (see below)
4. Access the web interface at `http://localhost:3000`
5. Start recording or importing audio files through the web interface

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

#### Manual Development Setup
```bash
# Backend server
cd note-server
go mod download
go run cmd/server/main.go

# Web frontend (in another terminal)
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
