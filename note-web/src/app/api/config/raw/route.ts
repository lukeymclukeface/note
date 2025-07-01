import { NextResponse } from 'next/server';
import { configService } from '@/lib/config';

export async function GET() {
  try {
    const config = await configService.getConfig();
    
    return NextResponse.json({
      success: true,
      config
    });
  } catch (error) {
    console.error('Failed to load raw config:', error);
    
    return NextResponse.json({
      success: false,
      error: error instanceof Error ? error.message : 'Failed to load configuration'
    }, { status: 500 });
  }
}
