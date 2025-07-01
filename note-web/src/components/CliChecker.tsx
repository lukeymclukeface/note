'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter, usePathname } from 'next/navigation';

interface CliStatus {
  checking: boolean;
  available: boolean;
  version?: string;
  error?: string;
}

export default function CliChecker({ children }: { children: React.ReactNode }) {
  const [cliStatus, setCliStatus] = useState<CliStatus>({ checking: true, available: true });
  const router = useRouter();
  const pathname = usePathname();

  const checkCliAvailability = useCallback(async () => {
    try {
      const response = await fetch('/api/cli-check');
      const data = await response.json();
      
      if (data.success) {
        const status = {
          checking: false,
          available: data.available,
          version: data.version,
          error: data.error
        };
        setCliStatus(status);
        
        // Redirect to CLI setup page if CLI is not available
        if (!data.available && pathname !== '/cli-setup') {
          router.push('/cli-setup');
        }
      } else {
        setCliStatus({
          checking: false,
          available: false,
          error: 'Failed to check CLI status'
        });
        
        // Redirect to CLI setup page on error too
        if (pathname !== '/cli-setup') {
          router.push('/cli-setup');
        }
      }
    } catch (error) {
      console.error('Failed to check CLI availability:', error);
      setCliStatus({
        checking: false,
        available: false,
        error: 'Network error checking CLI status'
      });
      
      // Redirect to CLI setup page on network error
      if (pathname !== '/cli-setup') {
        router.push('/cli-setup');
      }
    }
  }, [pathname, router]);

  useEffect(() => {
    // Don't check if we're already on the CLI setup page
    if (pathname === '/cli-setup') {
      setCliStatus({ checking: false, available: true });
      return;
    }
    
    checkCliAvailability();
  }, [pathname, checkCliAvailability]);

  // Show loading state while checking
  if (cliStatus.checking) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <svg className="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p className="text-gray-600 dark:text-gray-300">Checking CLI availability...</p>
        </div>
      </div>
    );
  }

  // Always show children - redirect logic is handled in the effect
  return <>{children}</>;
}
