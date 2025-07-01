import { NextResponse } from 'next/server';
import { validateDatabase } from '@/lib/database';

export async function GET() {
  try {
    const result = validateDatabase();
    
    return NextResponse.json({
      success: true,
      ...result
    });
  } catch (error) {
    console.error('Failed to validate database:', error);
    
    return NextResponse.json({
      success: false,
      error: error instanceof Error ? error.message : 'Failed to validate database'
    }, { status: 500 });
  }
}
