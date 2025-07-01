import { NextRequest, NextResponse } from 'next/server';
import { getDatabase } from '@/lib/database';

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
      SELECT id, filename, file_path, start_time, end_time, duration, 
             file_size, format, sample_rate, channels, created_at 
      FROM recordings 
      WHERE id = ?
    `);
    
    const recording = stmt.get(recordingId);
    db.close();

    if (!recording) {
      return NextResponse.json(
        { error: 'Recording not found' },
        { status: 404 }
      );
    }

    return NextResponse.json({
      success: true,
      recording
    });

  } catch (error) {
    console.error('Error fetching recording:', error);
    return NextResponse.json(
      { error: 'Failed to fetch recording' },
      { status: 500 }
    );
  }
}
