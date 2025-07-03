import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params;
    
    // Forward request to the Go server
    const response = await fetch(`${process.env.NOTE_SERVER_URL || 'http://localhost:8080'}/api/meetings/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Server responded with ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error fetching meeting:', error);
    return NextResponse.json(
      { success: false, error: 'Failed to fetch meeting' },
      { status: 500 }
    );
  }
}
