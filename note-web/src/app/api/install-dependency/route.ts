import { NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

const BREW_PACKAGES: Record<string, string> = {
  'FFmpeg': 'ffmpeg',
  'FFprobe': 'ffmpeg', // ffprobe comes with ffmpeg
  'Google Cloud CLI': 'google-cloud-sdk'
};

export async function POST(request: Request) {
  try {
    const { dependency } = await request.json();
    
    if (!dependency || !BREW_PACKAGES[dependency]) {
      return NextResponse.json({
        success: false,
        error: 'Invalid dependency specified'
      }, { status: 400 });
    }

    // First check if brew is available
    try {
      await execAsync('brew --version');
    } catch {
      return NextResponse.json({
        success: false,
        error: 'Homebrew is not installed or not available'
      }, { status: 400 });
    }

    const packageName = BREW_PACKAGES[dependency];
    console.log(`Installing ${dependency} using brew install ${packageName}`);

    try {
      // Install the package
      const { stdout, stderr } = await execAsync(`brew install ${packageName}`, {
        timeout: 300000 // 5 minute timeout
      });

      console.log('Install stdout:', stdout);
      console.log('Install stderr:', stderr);

      return NextResponse.json({
        success: true,
        message: `Successfully installed ${dependency}`,
        output: stdout || stderr
      });
      
    } catch (installError: unknown) {
      console.error('Installation failed:', installError);
      
      const err = installError as { signal?: string; message?: string };
      
      // Check if it's a timeout
      if (err.signal === 'SIGTERM') {
        return NextResponse.json({
          success: false,
          error: 'Installation timed out (took longer than 5 minutes)'
        }, { status: 408 });
      }

      // Check if package is already installed
      if (err.message?.includes('already installed')) {
        return NextResponse.json({
          success: true,
          message: `${dependency} is already installed`,
          output: err.message
        });
      }

      return NextResponse.json({
        success: false,
        error: err.message || 'Installation failed'
      }, { status: 500 });
    }

  } catch (error: unknown) {
    console.error('Request parsing failed:', error);
    return NextResponse.json({
      success: false,
      error: 'Invalid request'
    }, { status: 400 });
  }
}
