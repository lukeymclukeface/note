import { NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

interface HealthCheck {
  name: string;
  status: 'ok' | 'missing' | 'error';
  version?: string;
  error?: string;
}

async function checkCommand(commandName: string, versionCommand: string, name: string): Promise<HealthCheck> {
  try {
    // First check if the command exists using 'which'
    await execAsync(`which ${commandName}`);
    
    // If it exists, get the version
    const { stdout, stderr } = await execAsync(versionCommand);
    const output = stdout.trim() || stderr.trim();
    
    return {
      name,
      status: 'ok',
      version: output
    };
  } catch (error: any) {
    // Check if it's a "command not found" error from 'which' or the version command
    if (error.code === 127 || error.code === 1 || 
        error.message.includes('command not found') || 
        error.message.includes('not found')) {
      return {
        name,
        status: 'missing'
      };
    }
    
    return {
      name,
      status: 'error',
      error: error.message
    };
  }
}

export async function GET() {
  try {
    const checks = await Promise.all([
      checkCommand('brew', 'brew --version', 'Homebrew'),
      checkCommand('ffmpeg', 'ffmpeg -version 2>/dev/null | head -1', 'FFmpeg'),
      checkCommand('ffprobe', 'ffprobe -version 2>/dev/null | head -1', 'FFprobe'),
      checkCommand('gcloud', 'gcloud version 2>/dev/null | head -1', 'Google Cloud CLI')
    ]);

    return NextResponse.json({
      success: true,
      checks
    });
  } catch (error) {
    console.error('System health check failed:', error);
    return NextResponse.json({
      success: false,
      error: 'Failed to perform system health check'
    }, { status: 500 });
  }
}
