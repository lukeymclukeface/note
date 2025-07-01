import { NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

export async function GET() {
  try {
    // Try to run 'note --version' to check if CLI is available
    const { stdout, stderr } = await execAsync('note --version', { 
      timeout: 5000,
      env: { ...process.env, PATH: process.env.PATH }
    });
    
    // If we get here, the CLI is available
    return NextResponse.json({ 
      success: true, 
      available: true,
      version: stdout.trim(),
      message: 'Note CLI is available'
    });
  } catch (error: any) {
    // Check if it's a "command not found" error vs other errors
    const isNotFound = error.code === 127 || 
                       error.message?.includes('command not found') ||
                       error.message?.includes('not found') ||
                       error.stderr?.includes('command not found');
    
    if (isNotFound) {
      return NextResponse.json({ 
        success: true, 
        available: false,
        message: 'Note CLI is not installed or not in PATH'
      });
    }
    
    // Some other error occurred (CLI might be installed but having issues)
    return NextResponse.json({ 
      success: true, 
      available: true, // Assume it's installed but having issues
      error: error.message,
      message: 'Note CLI appears to be installed but encountered an error'
    });
  }
}
