'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function CliSetupPage() {
  const [isChecking, setIsChecking] = useState(false);
  const router = useRouter();

  const checkCliAvailability = async () => {
    setIsChecking(true);
    try {
      const response = await fetch('/api/cli-check');
      const data = await response.json();
      
      if (data.success && data.available) {
        // CLI is now available, redirect to home
        router.push('/');
      } else {
        // Still not available
        alert('Note CLI is still not detected. Please ensure it\'s installed and in your PATH.');
      }
    } catch (error) {
      console.error('Failed to check CLI availability:', error);
      alert('Failed to check CLI status. Please try again.');
    } finally {
      setIsChecking(false);
    }
  };

  const continueToApp = () => {
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="max-w-2xl w-full">
        <Card className="p-8 shadow-lg">
          <CardHeader className="text-center">
            <div className="text-6xl mb-4">🛠️</div>
            <CardTitle className="text-3xl mb-2">
              Note CLI Required
            </CardTitle>
            <p className="text-muted-foreground">
              The Note AI CLI is not installed or not available in your PATH. 
              Please install it to use this application.
            </p>
          </CardHeader>
          <CardContent className="space-y-6">

            {/* Installation Instructions */}
            <Card>
              <CardHeader>
                <CardTitle className="text-xl">
                  Installation Instructions
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <Card className="border-l-4 border-primary bg-primary/5">
                  <CardContent className="pt-4">
                    <h3 className="font-medium mb-2">
                      Option 1: Install from Source
                    </h3>
                    <p className="text-sm text-muted-foreground mb-2">
                      Clone the repository and build the CLI:
                    </p>
                    <Card className="bg-secondary p-3 font-mono text-sm overflow-x-auto">
                      <div>git clone https://github.com/your-repo/note-ai.git</div>
                      <div>cd note-ai</div>
                      <div>make install</div>
                    </Card>
                  </CardContent>
                </Card>
                
                <Card className="border-l-4 border-primary bg-primary/5">
                  <CardContent className="pt-4">
                    <h3 className="font-medium mb-2">
                      Option 2: Add to PATH
                    </h3>
                    <p className="text-sm text-muted-foreground mb-2">
                      If you&apos;ve already built the CLI, make sure it&apos;s in your PATH:
                    </p>
                    <Card className="bg-secondary p-3 font-mono text-sm overflow-x-auto">
                      <div>export PATH=$PATH:/path/to/note-cli</div>
                    </Card>
                  </CardContent>
                </Card>
                
                <Card className="border-l-4 border-primary bg-primary/5">
                  <CardContent className="pt-4">
                    <h3 className="font-medium mb-2">
                      Verify Installation
                    </h3>
                    <p className="text-sm text-muted-foreground mb-2">
                      Test that the CLI is working:
                    </p>
                    <Card className="bg-secondary p-3 font-mono text-sm overflow-x-auto">
                      <div>note --version</div>
                    </Card>
                  </CardContent>
                </Card>
              </CardContent>
            </Card>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              <Button
                onClick={checkCliAvailability}
                disabled={isChecking}
                className="flex items-center justify-center space-x-2"
              >
                {isChecking ? (
                  <>
                    <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Checking...</span>
                  </>
                ) : (
                  <>
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    <span>Check Again</span>
                  </>
                )}
              </Button>
              
              <Button
                onClick={continueToApp}
                variant="outline"
              >
                Continue Anyway
              </Button>
            </div>

            {/* Help Text */}
            <div className="text-center text-sm text-muted-foreground">
              <p>
                Need help? Check the{' '}
                <a 
                  href="https://github.com/your-repo/note-ai#installation" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-primary hover:underline"
                >
                  installation documentation
                </a>
                {' '}or{' '}
                <a 
                  href="https://github.com/your-repo/note-ai/issues" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-primary hover:underline"
                >
                  open an issue
                </a>.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
