# Source File Tracking System

This document describes the source file tracking system implemented in the note CLI application. This system tracks all files that are processed by the summarise command, storing metadata and processing status information.

## Overview

The source file tracking system consists of:

1. **Database Table**: `source_files` table that stores file metadata and processing status
2. **Metadata Extraction**: Comprehensive file metadata extraction for different file types
3. **Processing Status Tracking**: Real-time status updates during file processing
4. **Duplicate Prevention**: Hash-based duplicate detection to prevent reprocessing
5. **Command Interface**: CLI commands to view and manage tracked source files

## Database Schema

### source_files Table

```sql
CREATE TABLE source_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT NOT NULL UNIQUE,
    file_hash TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    file_type TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    metadata TEXT DEFAULT '{}',
    converted_path TEXT,
    processing_status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Fields Description

- **id**: Unique identifier for each source file record
- **file_path**: Absolute path to the original source file (unique constraint)
- **file_hash**: MD5 hash of the file content for duplicate detection
- **file_size**: Size of the file in bytes
- **file_type**: General file type category (audio, video, text, image, document, other)
- **mime_type**: Specific MIME type of the file
- **metadata**: JSON string containing file-specific metadata
- **converted_path**: Path to the converted/processed file (if applicable)
- **processing_status**: Current status (pending, processing, completed, failed)
- **created_at**: Timestamp when the record was created
- **updated_at**: Timestamp when the record was last updated

### Indexes

- `idx_source_files_hash`: Index on `file_hash` for quick duplicate detection
- `idx_source_files_status`: Index on `processing_status` for status filtering

## File Metadata System

### Metadata Structure

The metadata field contains a JSON object with file-specific information:

```json
{
  "file_name": "example.mp3",
  "file_extension": ".mp3",
  "last_modified": "2024-01-01T12:00:00Z",
  "duration": 120.5,
  "sample_rate": 44100,
  "channels": 2,
  "line_count": 100,
  "word_count": 500,
  "char_count": 3000,
  "encoding": "UTF-8",
  "custom": {}
}
```

### Supported File Types

#### Audio Files (.mp3, .wav, .flac, .m4a, .aac)
- Duration (seconds)
- Sample rate (Hz)
- Number of channels
- Bitrate (kbps)
- Codec information

#### Video Files (.mp4, .avi, .mov, .mkv)
- Duration (seconds)
- Resolution (width x height)
- Codec information
- Bitrate (kbps)

#### Text Files (.txt, .md, .json, .xml, .yaml)
- Line count
- Word count
- Character count
- Text encoding
- File structure analysis

#### Image Files (.jpg, .png, .gif, .bmp)
- Dimensions (width x height)
- Color depth
- Compression information

## Processing Status Workflow

### Status Values

1. **pending**: File record created, awaiting processing
2. **processing**: File is currently being processed
3. **completed**: File processing completed successfully
4. **failed**: File processing failed with errors

### Status Transitions

```
pending → processing → completed
              ↓
           failed
```

### Retry Logic

- Files with `failed` status can be retried by running the summarise command again
- The system will update the status to `processing` and attempt to reprocess the file
- Hash-based duplicate detection prevents duplicate processing of successful files

## Command Interface

### List Source Files

View all tracked source files with their status and metadata:

```bash
# List all source files
note list source-files

# Filter by status
note list source-files --status completed
note list source-files --status failed

# Filter by file type
note list source-files --type audio
note list source-files --type text

# Combine filters
note list source-files --status completed --type audio
```

### Summarise Tracking

The summarise command automatically tracks all processed files:

```bash
# Summarise a file (automatically tracked)
note summarise audio.mp3
note summarise document.txt
```

## Duplicate Prevention

### Hash-Based Detection

- MD5 hash is calculated for each file's content
- Files with identical hashes are considered duplicates
- Prevents reprocessing of identical files even with different paths

### Path-Based Detection

- File paths are stored with unique constraints
- Prevents multiple entries for the same file path
- Supports retry logic for failed processing attempts

## Error Handling

### Database Errors

- Connection failures are handled gracefully
- Transaction rollbacks ensure data consistency
- Detailed error messages for troubleshooting

### File Access Errors

- Missing files are detected and reported
- Permission errors are handled appropriately
- Network path issues are caught and logged

## Migration and Compatibility

### Database Migration

The source file tracking system is automatically initialized when:

1. Running `note setup` command
2. First summarise operation (if database exists but table is missing)

### Backward Compatibility

- Existing databases are automatically migrated
- No data loss occurs during migration
- Old functionality remains intact

## Performance Considerations

### Indexing Strategy

- Hash-based lookups are O(1) with proper indexing
- Status filtering is optimized with dedicated indexes
- Large file metadata is stored efficiently as JSON

### Caching

- File metadata extraction is cached during processing
- Temporary files are managed through the cache system
- Memory usage is optimized for large file processing

## Security Considerations

### File Access

- Only processes files in authorized directories
- Validates file paths to prevent directory traversal
- Respects file system permissions

### Data Storage

- Sensitive file paths are stored securely
- No file content is stored in the database
- Metadata is sanitized before storage

## Troubleshooting

### Common Issues

1. **"source_files table not found"**
   - Solution: Run `note setup` to initialize the database

2. **"File already processed"**
   - Solution: Check with `note list source-files` and verify status

3. **"Failed to calculate file hash"**
   - Solution: Ensure file exists and is readable

### Debug Information

Use the verbose flag for detailed processing information:

```bash
note summarise --verbose audio.mp3
note list source-files --verbose
```

## API Reference

### Database Functions

- `CreateSourceFile(db, filePath, fileHash, fileSize, fileType, mimeType, metadata)`
- `GetSourceFileByHash(db, fileHash)`
- `GetSourceFileByPath(db, filePath)`
- `UpdateSourceFileStatus(db, id, status)`
- `UpdateSourceFileConvertedPath(db, id, convertedPath)`
- `ListSourceFiles(db, status)`
- `DeleteSourceFile(db, id)`

### Helper Functions

- `GetFileHash(filePath)`: Calculate MD5 hash of file
- `GetMimeType(filePath)`: Determine MIME type
- `GetFileType(mimeType)`: Categorize file type
- `ExtractFileMetadata(filePath)`: Extract comprehensive metadata
- `ProcessFileForDatabase(filePath)`: Prepare file for database storage

## Future Enhancements

### Planned Features

1. **Advanced Metadata Extraction**
   - Integration with FFmpeg for detailed audio/video metadata
   - EXIF data extraction for images
   - Document structure analysis for PDFs

2. **File Relationships**
   - Track relationships between source and processed files
   - Version history for multiple processing attempts
   - Dependency tracking for multi-file operations

3. **Batch Operations**
   - Bulk import with progress tracking
   - Batch status updates
   - Mass cleanup operations

4. **Web Interface**
   - Web-based file browser
   - Status dashboard
   - Search and filtering capabilities

### Configuration Options

Future versions may include configuration for:

- Metadata extraction depth
- Hash algorithm selection
- Retention policies
- Processing priorities
