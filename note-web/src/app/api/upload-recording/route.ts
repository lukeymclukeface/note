import { NextRequest, NextResponse } from 'next/server';
import { writeFile, mkdir } from 'fs/promises';
import { join } from 'path';
import { homedir } from 'os';
import { insertRecording } from '@/lib/database';

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get('audio') as File;
    const startTimeStr = formData.get('startTime') as string;
    const endTimeStr = formData.get('endTime') as string;
    
    if (!file) {
      return NextResponse.json({ error: 'No audio file provided' }, { status: 400 });
    }
    
    // Use provided times or fallback to current time
    const startTime = startTimeStr ? new Date(startTimeStr) : new Date();
    const endTime = endTimeStr ? new Date(endTimeStr) : new Date();

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
    
    // Calculate duration in nanoseconds (compatible with CLI format)
    const durationMs = endTime.getTime() - startTime.getTime();
    const durationNs = durationMs * 1000000; // Convert to nanoseconds
    
    // Insert recording into database
    const recordingId = insertRecording({
      filename,
      file_path: filepath,
      start_time: startTime.toISOString(),
      end_time: endTime.toISOString(),
      duration: durationNs,
      file_size: buffer.length,
      format: 'webm',
      sample_rate: 44100, // From MediaRecorder config
      channels: 1, // Typically mono for voice recordings
    });
    
    if (!recordingId) {
      console.warn('Failed to create database entry for recording:', filename);
    }

    return NextResponse.json({ 
      success: true, 
      filename,
      filepath,
      size: buffer.length,
      recordingId,
      duration: durationMs 
    });

  } catch (error) {
    console.error('Error uploading recording:', error);
    return NextResponse.json(
      { error: 'Failed to upload recording' },
      { status: 500 }
    );
  }
}
