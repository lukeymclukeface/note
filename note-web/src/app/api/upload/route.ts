import { NextRequest, NextResponse } from 'next/server';
import { writeFile, mkdir } from 'fs/promises';
import { existsSync } from 'fs';
import path from 'path';
import os from 'os';

const UPLOAD_DIR = path.join(os.homedir(), '.noteai', 'import');
const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

const SUPPORTED_TYPES = [
  // Audio files
  '.mp3', '.wav', '.m4a', '.aac', '.flac', '.ogg',
  // Text files
  '.txt', '.md', '.doc', '.docx', '.pdf'
];

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get('file') as File;

    if (!file) {
      return NextResponse.json(
        { success: false, message: 'No file provided' },
        { status: 400 }
      );
    }

    // Validate file size
    if (file.size > MAX_FILE_SIZE) {
      return NextResponse.json(
        { success: false, message: 'File size exceeds 50MB limit' },
        { status: 400 }
      );
    }

    // Validate file type
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase();
    if (!SUPPORTED_TYPES.includes(fileExtension)) {
      return NextResponse.json(
        { 
          success: false, 
          message: `Unsupported file type. Supported types: ${SUPPORTED_TYPES.join(', ')}` 
        },
        { status: 400 }
      );
    }

    // Ensure upload directory exists
    if (!existsSync(UPLOAD_DIR)) {
      await mkdir(UPLOAD_DIR, { recursive: true });
    }

    // Generate unique filename to avoid conflicts
    const timestamp = Date.now();
    const sanitizedName = file.name.replace(/[^a-zA-Z0-9.-]/g, '_');
    const filename = `${timestamp}_${sanitizedName}`;
    const filepath = path.join(UPLOAD_DIR, filename);

    // Convert file to buffer and save
    const bytes = await file.arrayBuffer();
    const buffer = Buffer.from(bytes);
    
    await writeFile(filepath, buffer);

    // Get file stats for response
    const stats = {
      filename,
      originalName: file.name,
      size: file.size,
      type: file.type,
      uploadPath: filepath
    };

    return NextResponse.json({
      success: true,
      message: `File "${file.name}" uploaded successfully`,
      filename,
      stats
    });

  } catch (error) {
    console.error('Upload error:', error);
    
    return NextResponse.json(
      { 
        success: false, 
        message: 'Internal server error during upload' 
      },
      { status: 500 }
    );
  }
}

// Handle other HTTP methods
export async function GET() {
  return NextResponse.json(
    { message: 'Upload endpoint - POST only' },
    { status: 405 }
  );
}

export async function PUT() {
  return NextResponse.json(
    { message: 'Upload endpoint - POST only' },
    { status: 405 }
  );
}

export async function DELETE() {
  return NextResponse.json(
    { message: 'Upload endpoint - POST only' },
    { status: 405 }
  );
}
