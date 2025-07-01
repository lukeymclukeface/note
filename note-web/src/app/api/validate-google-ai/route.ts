import { NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

interface GoogleAIValidationResult {
  valid: boolean;
  error?: string;
  needsAuth?: boolean;
  details?: {
    projectId?: string;
    account?: string;
    location?: string;
    gcloudVersion?: string;
    services?: string[];
  };
}

export async function POST(request: Request) {
  try {
    const { projectId, location } = await request.json();

    if (!projectId) {
      return NextResponse.json({
        success: false,
        error: 'Project ID is required'
      }, { status: 400 });
    }

    const validation: GoogleAIValidationResult = {
      valid: false
    };

    // Check if gcloud CLI is installed
    try {
      const { stdout: versionOutput } = await execAsync('gcloud version --format="value(version)"');
      validation.details = {
        gcloudVersion: versionOutput.trim(),
        projectId,
        location: location || 'us-central1'
      };
    } catch {
      validation.error = 'Google Cloud CLI (gcloud) is not installed or not in PATH';
      return NextResponse.json({
        success: true,
        validation
      });
    }

    // Check current authentication status
    try {
      const { stdout: accountOutput } = await execAsync('gcloud auth list --filter=status:ACTIVE --format="value(account)"');
      const activeAccount = accountOutput.trim();
      
      if (!activeAccount) {
        validation.needsAuth = true;
        validation.error = 'No active Google Cloud authentication found. Please run "gcloud auth login"';
        return NextResponse.json({
          success: true,
          validation
        });
      }

      validation.details!.account = activeAccount;
    } catch {
      validation.needsAuth = true;
      validation.error = 'Failed to check authentication status. Please run "gcloud auth login"';
      return NextResponse.json({
        success: true,
        validation
      });
    }

    // Check if the project exists and is accessible
    try {
      await execAsync(`gcloud projects describe ${projectId} --format="value(projectId)"`);
    } catch {
      validation.error = `Project "${projectId}" not found or not accessible. Check project ID and permissions.`;
      return NextResponse.json({
        success: true,
        validation
      });
    }

    // Check if required APIs are enabled
    try {
      const { stdout: servicesOutput } = await execAsync(
        `gcloud services list --enabled --project=${projectId} --filter="name:(aiplatform.googleapis.com OR speech.googleapis.com)" --format="value(config.name)"`
      );
      
      const enabledServices = servicesOutput.trim().split('\n').filter(s => s);
      validation.details!.services = enabledServices;

      const requiredServices = ['aiplatform.googleapis.com', 'speech.googleapis.com'];
      const missingServices = requiredServices.filter(service => !enabledServices.includes(service));

      if (missingServices.length > 0) {
        validation.error = `Missing required API services: ${missingServices.join(', ')}. Enable them in the Google Cloud Console.`;
        return NextResponse.json({
          success: true,
          validation
        });
      }
    } catch {
      validation.error = 'Failed to check API services. Ensure you have proper permissions.';
      return NextResponse.json({
        success: true,
        validation
      });
    }

    // Try to make a simple API call to verify everything works
    try {
      // Test with a simple location list call to AI Platform
      await execAsync(
        `gcloud ai-platform locations list --project=${projectId} --format="value(name)" --limit=1`
      );
      
      validation.valid = true;
      validation.error = undefined;
    } catch (error) {
      const err = error as { message?: string };
      validation.error = `API test failed: ${err.message || 'Unknown error'}`;
    }

    return NextResponse.json({
      success: true,
      validation
    });

  } catch (error) {
    console.error('Google AI validation error:', error);
    
    return NextResponse.json({
      success: false,
      error: error instanceof Error ? error.message : 'Failed to validate Google AI configuration'
    }, { status: 500 });
  }
}
