'use client';

import { useState, useEffect } from 'react';
import { Config } from '../types';

interface TableInfo {
  name: string;
  columns: Array<{
    name: string;
    type: string;
    nullable: boolean;
    primaryKey: boolean;
  }>;
  rowCount: number;
}

interface DatabaseValidation {
  connected: boolean;
  tables: TableInfo[];
  errors: string[];
  version?: string;
}

export default function DatabaseSettingsPage() {
  const [config, setConfig] = useState<Config | null>(null);
  const [validation, setValidation] = useState<DatabaseValidation | null>(null);
  const [isValidating, setIsValidating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await fetch('/api/config');
      const data = await response.json();
      if (data.success) {
        setConfig(data.config);
      } else {
        setConfig(null);
      }
    } catch (error) {
      console.error('Failed to load config:', error);
      setConfig(null);
    }
  };

  const validateDatabase = async () => {
    setIsValidating(true);
    setError(null);
    try {
      const response = await fetch('/api/database/validate');
      const data = await response.json();
      if (data.success) {
        setValidation(data.validation);
      } else {
        setError(data.error || 'Failed to validate database');
        setValidation(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to validate database');
      setValidation(null);
    } finally {
      setIsValidating(false);
    }
  };

  if (!config) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">❌</div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Error</h1>
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          Unable to read the configuration file. Please check if the CLI is set up.
        </p>
      </div>
    );
  }

  return (
    <div>
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Database Settings</h1>
        <p className="text-gray-600 dark:text-gray-300">
          Manage database configuration and validate table structure.
        </p>
      </div>

      {/* Database Configuration */}
      <section className="mb-8">
        <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">Database Configuration</h2>
        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Database Path</label>
              <input
                type="text"
                value={config.database_path || ''}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 font-mono text-sm cursor-not-allowed"
                readOnly
              />
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Database path is configured in the general settings
              </p>
            </div>
            <div className="flex items-center justify-between pt-4">
              <div className="flex items-center">
                <div className={`w-3 h-3 rounded-full mr-3 ${config.database_path ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-sm text-gray-800 dark:text-gray-200">
                  Database {config.database_path ? 'Configured' : 'Not Configured'}
                </span>
              </div>
              <button
                onClick={validateDatabase}
                disabled={isValidating || !config.database_path}
                className="btn btn-primary btn-sm"
              >
                {isValidating ? (
                  <span className="flex items-center">
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-current" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Validating...
                  </span>
                ) : (
                  'Validate Database'
                )}
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Validation Results */}
      {error && (
        <section className="mb-8">
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
            <div className="flex items-center mb-2">
              <div className="text-red-500 mr-2">❌</div>
              <h3 className="text-lg font-medium text-red-900 dark:text-red-200">Validation Error</h3>
            </div>
            <p className="text-red-800 dark:text-red-300">{error}</p>
          </div>
        </section>
      )}

      {validation && (
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">Database Validation Results</h2>
          
          {/* Connection Status */}
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm mb-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className={`w-4 h-4 rounded-full mr-3 ${validation.connected ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  Database Connection: {validation.connected ? 'Connected' : 'Failed'}
                </h3>
              </div>
              {validation.version && (
                <span className="text-sm text-gray-500 dark:text-gray-400 font-mono">
                  SQLite {validation.version}
                </span>
              )}
            </div>
          </div>

          {/* Validation Errors */}
          {validation.errors.length > 0 && (
            <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6 mb-6">
              <div className="flex items-center mb-3">
                <div className="text-yellow-500 mr-2">⚠️</div>
                <h3 className="text-lg font-medium text-yellow-900 dark:text-yellow-200">Validation Issues</h3>
              </div>
              <ul className="list-disc list-inside space-y-1">
                {validation.errors.map((error, index) => (
                  <li key={index} className="text-yellow-800 dark:text-yellow-300 text-sm">{error}</li>
                ))}
              </ul>
            </div>
          )}

          {/* Tables */}
          <div className="space-y-6">
            <h3 className="text-xl font-medium text-gray-900 dark:text-white">Database Tables</h3>
            {validation.tables.length > 0 ? (
              <div className="grid gap-6">
                {validation.tables.map((table) => (
                  <div key={table.name} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="text-lg font-medium text-gray-900 dark:text-white">{table.name}</h4>
                      <span className="text-sm text-gray-500 dark:text-gray-400">
                        {table.rowCount} rows
                      </span>
                    </div>
                    <div className="overflow-x-auto">
                      <table className="w-full text-sm">
                        <thead>
                          <tr className="border-b border-gray-200 dark:border-gray-600">
                            <th className="text-left py-2 text-gray-700 dark:text-gray-200">Column</th>
                            <th className="text-left py-2 text-gray-700 dark:text-gray-200">Type</th>
                            <th className="text-left py-2 text-gray-700 dark:text-gray-200">Nullable</th>
                            <th className="text-left py-2 text-gray-700 dark:text-gray-200">Primary Key</th>
                          </tr>
                        </thead>
                        <tbody>
                          {table.columns.map((column, index) => (
                            <tr key={index} className="border-b border-gray-100 dark:border-gray-700">
                              <td className="py-2 font-mono text-gray-900 dark:text-gray-100">{column.name}</td>
                              <td className="py-2 font-mono text-gray-600 dark:text-gray-300">{column.type}</td>
                              <td className="py-2">
                                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                                  column.nullable 
                                    ? 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200'
                                    : 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200'
                                }`}>
                                  {column.nullable ? 'Yes' : 'No'}
                                </span>
                              </td>
                              <td className="py-2">
                                {column.primaryKey && (
                                  <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
                                    PK
                                  </span>
                                )}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <div className="text-gray-400 dark:text-gray-500 text-4xl mb-2">📊</div>
                <p className="text-gray-500 dark:text-gray-400">No tables found in the database</p>
              </div>
            )}
          </div>
        </section>
      )}

      {!validation && !isValidating && !error && (
        <div className="text-center py-12">
          <div className="text-gray-400 dark:text-gray-500 text-6xl mb-4">🗄️</div>
          <h3 className="text-xl font-medium text-gray-900 dark:text-white mb-2">Database Validation</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-6">
            Click &quot;Validate Database&quot; to check your database configuration and table structure.
          </p>
        </div>
      )}
    </div>
  );
}
