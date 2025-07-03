import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const file = formData.get('audio') as File;
    const startTimeStr = formData.get('startTime') as string;
    const endTimeStr = formData.get('endTime') as string;
    
    if (!file) {
      return NextResponse.json({ error: 'No audio file provided' }, { status: 400 });
    }
    
    // Create FormData to send to the Go server
    const serverFormData = new FormData();
    serverFormData.append('audio', file);
    if (startTimeStr) serverFormData.append('startTime', startTimeStr);
    if (endTimeStr) serverFormData.append('endTime', endTimeStr);
    
    // Forward the upload request to the Go server
    const response = await fetch(`${process.env.NOTE_SERVER_URL || 'http://localhost:8080'}/api/upload-recording`, {
      method: 'POST',
      body: serverFormData,
    });

    if (!response.ok) {
      throw new Error(`Server responded with ${response.status}`);
    }

    const result = await response.json();
    return NextResponse.json(result);

  } catch (error) {
    console.error('Error uploading recording:', error);
    return NextResponse.json(
      { success: false, error: 'Failed to upload recording' },
      { status: 500 }
    );
  }
}
