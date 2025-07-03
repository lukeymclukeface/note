import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    
    // Forward audio streaming request to the Go server
    const response = await fetch(`${process.env.NOTE_SERVER_URL || 'http://localhost:8080'}/api/recordings/${id}/audio`, {
      method: 'GET',
    });

    if (!response.ok) {
      if (response.status === 404) {
        return NextResponse.json(
          { error: 'Audio file not found' },
          { status: 404 }
        );
      }
      throw new Error(`Server responded with ${response.status}`);
    }

    // Get the audio stream from the server
    const audioBuffer = await response.arrayBuffer();
    
    // Forward the response headers from the server
    const headers = new Headers();
    headers.set('Content-Type', response.headers.get('Content-Type') || 'audio/mpeg');
    headers.set('Content-Length', response.headers.get('Content-Length') || audioBuffer.byteLength.toString());
    headers.set('Accept-Ranges', 'bytes');
    headers.set('Cache-Control', 'public, max-age=3600');
    
    if (response.headers.get('Content-Disposition')) {
      headers.set('Content-Disposition', response.headers.get('Content-Disposition')!);
    }

    return new NextResponse(audioBuffer, {
      status: 200,
      headers,
    });

  } catch (error) {
    console.error('Error proxying audio file:', error);
    return NextResponse.json(
      { error: 'Failed to serve audio file' },
      { status: 500 }
    );
  }
}
