import { NextRequest, NextResponse } from 'next/server';
import { writeFile, mkdir } from 'fs/promises';
import { join } from 'path';
import { homedir } from 'os';

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get('audio') as File;
    
    if (!file) {
      return NextResponse.json({ error: 'No audio file provided' }, { status: 400 });
    }

    // Create .noteai/recordings directory if it doesn't exist
    const noteaiDir = join(homedir(), '.noteai');
    const recordingsDir = join(noteaiDir, 'recordings');
    
    try {
      await mkdir(recordingsDir, { recursive: true });
    } catch {
      // Directory might already exist, ignore error
    }

    // Generate filename with timestamp
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const filename = `recording_${timestamp}.webm`;
    const filepath = join(recordingsDir, filename);

    // Convert file to buffer and save
    const bytes = await file.arrayBuffer();
    const buffer = Buffer.from(bytes);
    
    await writeFile(filepath, buffer);

    return NextResponse.json({ 
      success: true, 
      filename,
      filepath,
      size: buffer.length 
    });

  } catch (error) {
    console.error('Error uploading recording:', error);
    return NextResponse.json(
      { error: 'Failed to upload recording' },
      { status: 500 }
    );
  }
}
