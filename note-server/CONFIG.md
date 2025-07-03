# Configuration Guide

## Overview

The note server now supports flexible configuration management. OpenAI API keys and other AI settings are no longer required at startup and can be configured through the web interface.

## Configuration Methods

### 1. Environment Variables (Server Settings)
Basic server configuration is still handled through environment variables:

```bash
PORT=8080                    # Server port
HOST=localhost              # Server host
MEDIA_TMP_DIR=/tmp/note-media # Temporary media storage
LOG_LEVEL=info              # Logging level
DEV_MODE=false              # Development mode
```

### 2. JSON Configuration File (AI Settings)
AI-related settings are stored in `~/.noteai/config.json` and can be managed through:
- Web interface at `/settings/ai` (recommended)
- Direct API calls to `/api/config`
- Manual editing of the JSON file

## Configuration Endpoints

### GET /api/config
Returns the current configuration with masked API keys.

### PUT /api/config
Updates the configuration with new values.

### GET /api/config/raw
Returns the unmasked configuration for editing (use with caution).

## Configuration Structure

The JSON configuration file supports the following fields:

```json
{
  "openai_key": "sk-...",
  "transcription_provider": "openai",
  "transcription_model": "whisper-1",
  "summary_provider": "openai", 
  "summary_model": "gpt-4"
}
```

## Migration from Environment Variables

If you were previously using the `OPENAI_KEY` environment variable:

1. Start the server (it no longer requires the OpenAI key)
2. Go to the web interface at `/settings/ai`
3. Enter your OpenAI API key and save
4. Remove the `OPENAI_KEY` environment variable

## Docker Usage

The server now mounts the configuration directory to persist settings:

```yaml
volumes:
  - ~/.noteai:/home/app/.noteai
```

This ensures your AI configuration persists across container restarts.

## Security Notes

- API keys are stored in `~/.noteai/config.json` with 600 permissions
- The `/api/config` endpoint masks API keys in responses
- Use `/api/config/raw` only when necessary for editing
