import { NextRequest, NextResponse } from 'next/server';
import { getDatabase } from '@/lib/database';
import fs from 'fs';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const resolvedParams = await params;
    const recordingId = parseInt(resolvedParams.id);
    
    if (isNaN(recordingId)) {
      return NextResponse.json(
        { error: 'Invalid recording ID' },
        { status: 400 }
      );
    }

    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT file_path, format, filename
      FROM recordings 
      WHERE id = ?
    `);
    
    const recording = stmt.get(recordingId) as { file_path: string; format: string; filename: string } | undefined;
    db.close();

    if (!recording) {
      return NextResponse.json(
        { error: 'Recording not found' },
        { status: 404 }
      );
    }

    // Check if file exists
    if (!fs.existsSync(recording.file_path)) {
      return NextResponse.json(
        { error: 'Audio file not found' },
        { status: 404 }
      );
    }

    // Read the file
    const fileBuffer = fs.readFileSync(recording.file_path);
    
    // Determine content type based on format
    let contentType = 'audio/mpeg'; // default
    switch (recording.format.toLowerCase()) {
      case 'wav':
        contentType = 'audio/wav';
        break;
      case 'mp3':
        contentType = 'audio/mpeg';
        break;
      case 'ogg':
        contentType = 'audio/ogg';
        break;
      case 'webm':
        contentType = 'audio/webm';
        break;
      default:
        contentType = 'audio/mpeg';
    }

    // Return the audio file with appropriate headers
    return new NextResponse(fileBuffer, {
      status: 200,
      headers: {
        'Content-Type': contentType,
        'Content-Length': fileBuffer.length.toString(),
        'Accept-Ranges': 'bytes',
        'Cache-Control': 'public, max-age=3600',
        'Content-Disposition': `inline; filename="${recording.filename}"`,
      },
    });

  } catch (error) {
    console.error('Error serving audio file:', error);
    return NextResponse.json(
      { error: 'Failed to serve audio file' },
      { status: 500 }
    );
  }
}
