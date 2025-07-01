'use client';

import { useState, useEffect } from 'react';

export default function RawConfigPage() {
  const [config, setConfig] = useState<string>('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadRawConfig();
  }, []);

  const loadRawConfig = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('/api/config/raw');
      if (response.ok) {
        const data = await response.json();
        if (data.success) {
          setConfig(JSON.stringify(data.config, null, 2));
          setError(null);
        } else {
          setError(data.error || 'Failed to load configuration');
        }
      } else {
        setError(`HTTP ${response.status}: ${response.statusText}`);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load configuration');
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600 dark:text-gray-300">Loading configuration...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">❌</div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Error</h1>
        <p className="text-red-600 dark:text-red-400 mb-8">{error}</p>
        <button
          onClick={loadRawConfig}
          className="btn btn-primary"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Raw Configuration</h1>
        <p className="text-gray-600 dark:text-gray-300">
          View the raw JSON configuration file from ~/.noteai/config.json
        </p>
      </div>
      
      <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-medium text-gray-900 dark:text-white">Configuration File Contents</h2>
          <button
            onClick={loadRawConfig}
            className="btn btn-sm btn-outline"
          >
            Reload
          </button>
        </div>
        
        <pre className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-600 rounded-md p-4 overflow-auto text-sm font-mono whitespace-pre-wrap">
          <code className="text-gray-800 dark:text-gray-200">{config}</code>
        </pre>
      </div>
    </div>
  );
}
