import { NextRequest, NextResponse } from 'next/server';

// Get the server URL from environment or default to localhost
function getServerUrl(): string {
  return process.env.NOTE_SERVER_URL || 'http://localhost:8080';
}

export async function PUT(request: NextRequest) {
  try {
    const serverUrl = getServerUrl();
    const updatedConfig = await request.json();
    
    // Proxy the request to the server
    const response = await fetch(`${serverUrl}/api/config`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updatedConfig),
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      return NextResponse.json(
        { success: false, error: `Server error: ${response.status} - ${errorText}` },
        { status: response.status }
      );
    }
    
    const result = await response.json();
    return NextResponse.json(result);
  } catch (error) {
    console.error('Error proxying config update to server:', error);
    return NextResponse.json(
      { success: false, error: `Failed to connect to server: ${error instanceof Error ? error.message : 'Unknown error'}` },
      { status: 500 }
    );
  }
}

export async function GET() {
  try {
    const serverUrl = getServerUrl();
    
    // Proxy the request to the server
    const response = await fetch(`${serverUrl}/api/config`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      return NextResponse.json(
        { success: false, error: `Server error: ${response.status} - ${errorText}` },
        { status: response.status }
      );
    }
    
    const result = await response.json();
    return NextResponse.json(result);
  } catch (error) {
    console.error('Error proxying config request to server:', error);
    return NextResponse.json(
      { success: false, error: `Failed to connect to server: ${error instanceof Error ? error.message : 'Unknown error'}` },
      { status: 500 }
    );
  }
}
