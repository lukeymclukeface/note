'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

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
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
      <div className="max-w-2xl w-full">
        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-8 shadow-lg">
          <div className="text-center mb-6">
            <div className="text-6xl mb-4">🛠️</div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
              Note CLI Required
            </h1>
            <p className="text-gray-600 dark:text-gray-300">
              The Note AI CLI is not installed or not available in your PATH. 
              Please install it to use this application.
            </p>
          </div>

          <div className="space-y-6">
            {/* Installation Instructions */}
            <div className="bg-gray-50 dark:bg-gray-900 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
                Installation Instructions
              </h2>
              
              <div className="space-y-4">
                <div className="border-l-4 border-blue-500 pl-4">
                  <h3 className="font-medium text-gray-900 dark:text-white mb-2">
                    Option 1: Install from Source
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                    Clone the repository and build the CLI:
                  </p>
                  <div className="bg-gray-800 dark:bg-gray-700 rounded-md p-3 font-mono text-sm text-green-400 overflow-x-auto">
                    <div>git clone https://github.com/your-repo/note-ai.git</div>
                    <div>cd note-ai</div>
                    <div>make install</div>
                  </div>
                </div>
                
                <div className="border-l-4 border-blue-500 pl-4">
                  <h3 className="font-medium text-gray-900 dark:text-white mb-2">
                    Option 2: Add to PATH
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                    If you&apos;ve already built the CLI, make sure it&apos;s in your PATH:
                  </p>
                  <div className="bg-gray-800 dark:bg-gray-700 rounded-md p-3 font-mono text-sm text-green-400 overflow-x-auto">
                    <div>export PATH=$PATH:/path/to/note-cli</div>
                  </div>
                </div>
                
                <div className="border-l-4 border-blue-500 pl-4">
                  <h3 className="font-medium text-gray-900 dark:text-white mb-2">
                    Verify Installation
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                    Test that the CLI is working:
                  </p>
                  <div className="bg-gray-800 dark:bg-gray-700 rounded-md p-3 font-mono text-sm text-green-400 overflow-x-auto">
                    <div>note --version</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-3 justify-center">
              <button
                onClick={checkCliAvailability}
                disabled={isChecking}
                className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-md font-medium transition-colors duration-200 flex items-center justify-center space-x-2"
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
              </button>
              
              <button
                onClick={continueToApp}
                className="px-6 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-md font-medium transition-colors duration-200"
              >
                Continue Anyway
              </button>
            </div>

            {/* Help Text */}
            <div className="text-center text-sm text-gray-500 dark:text-gray-400">
              <p>
                Need help? Check the{' '}
                <a 
                  href="https://github.com/your-repo/note-ai#installation" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-blue-600 dark:text-blue-400 hover:underline"
                >
                  installation documentation
                </a>
                {' '}or{' '}
                <a 
                  href="https://github.com/your-repo/note-ai/issues" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-blue-600 dark:text-blue-400 hover:underline"
                >
                  open an issue
                </a>.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
