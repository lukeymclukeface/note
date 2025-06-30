# Note CLI

A powerful command-line interface for AI-powered note-taking, audio recording, and content management.

## Installation

### Build from Source
```bash
go build -o note cmd/note/main.go
```

### Dependencies
- Go 1.23+
- FFmpeg (automatically installed via Homebrew during setup)
- OpenAI API key

## Initial Setup

Run the interactive setup wizard:
```bash
./note setup
```

This will:
- Install FFmpeg via Homebrew if needed
- Prompt for your OpenAI API key
- Initialize the SQLite database
- Create necessary directories

## Configuration

Configure AI models for transcription and summarization:
```bash
./note config model
```

View current configuration:
```bash
./note config
```

## Commands

### Recording Audio
```bash
# Start interactive recording
./note record

# Record with specific device
./note record --device "Built-in Microphone"
```

### Importing Content
```bash
# Import audio file
./note import audio.mp3

# Import text/markdown file
./note import document.md

# Interactive file selection (if no file specified)
./note import
```

### Managing Notes
```bash
# Interactive note browser
./note list

# View all recordings
./note recordings

# Delete specific recording
./note recordings --delete recording_id
```

### Content Creation
```bash
# Create new note from scratch
./note create
```

## Features

### Audio Processing
- **Multi-device recording**: Select from available audio input devices
- **Format support**: mp3, wav, m4a, ogg, flac
- **Chunked transcription**: Automatic splitting for large files (>25MB)
- **Speaker diarization**: AI-powered speaker identification

### AI Integration
- **Whisper transcription**: High-quality speech-to-text
- **Intelligent summarization**: Context-aware summaries
- **Content type detection**: Automatically detects meetings, interviews, lectures
- **Custom prompts**: Embedded markdown templates for different content types

### User Interface
- **Interactive prompts**: Beautiful terminal UI with Bubble Tea
- **Progress indicators**: Real-time spinners for long operations
- **Colorized output**: Enhanced readability with colors and formatting
- **Verbose logging**: Detailed execution logs with `--verbose` flag

### Data Management
- **SQLite database**: Efficient metadata storage
- **File organization**: Structured storage in `~/.noteai/`
- **Backup-friendly**: All data in user home directory

## Architecture

The CLI follows a clean service-oriented architecture:

```
cmd/note/main.go           # Entry point
internal/
├── cmd/                   # Command implementations
├── config/                # Configuration management
├── database/              # SQLite operations
├── services/              # Core business logic
│   ├── audio.go          # Audio processing
│   ├── openai.go         # AI integration
│   ├── file.go           # File operations
│   ├── prompts.go        # Template management
│   ├── ui.go             # Terminal UI
│   └── verbose.go        # Logging
├── constants/             # Application constants
└── helpers/               # Utility functions
```

## Configuration Files

All configuration stored in `~/.noteai/`:

- `config.json` - Settings and API keys
- `notes.db` - SQLite database
- `recordings/` - Audio files
- `notes/` - Generated markdown files
- `prompts/` - AI prompt templates (embedded)

## Example Workflow

1. **Setup**: `./note setup`
2. **Configure models**: `./note config model`
3. **Record meeting**: `./note record`
4. **Import existing audio**: `./note import meeting.mp3`
5. **Browse notes**: `./note list`
6. **View recordings**: `./note recordings`

## Verbose Mode

Enable detailed logging for any command:
```bash
./note import audio.mp3 --verbose
```

Shows:
- API request/response details
- File processing steps
- Token usage and costs
- Database operations
- Timing information

## Error Handling

The CLI provides comprehensive error handling:
- Network connectivity issues
- API rate limiting
- File permission problems
- Audio format validation
- Database constraints

## Performance

- **Chunked processing**: Large files split into 10-minute segments
- **Concurrent operations**: Parallel transcription of chunks
- **Progress tracking**: Real-time status updates
- **Memory efficient**: Streaming audio processing
- **Token optimization**: Smart prompt management
