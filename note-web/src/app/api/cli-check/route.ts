import { NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

export async function GET() {
  try {
    // Try to run 'note --version' to check if CLI is available
    const { stdout } = await execAsync('note --version', { 
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
  } catch (error: unknown) {
    // Check if it's a "command not found" error vs other errors
    const errorObj = error as { code?: number; message?: string; stderr?: string };
    const isNotFound = errorObj.code === 127 || 
                       errorObj.message?.includes('command not found') ||
                       errorObj.message?.includes('not found') ||
                       errorObj.stderr?.includes('command not found');
    
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
      error: errorObj.message,
      message: 'Note CLI appears to be installed but encountered an error'
    });
  }
}
