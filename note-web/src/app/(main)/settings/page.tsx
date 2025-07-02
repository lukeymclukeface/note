'use client';

import { useState, useEffect } from 'react';
import { Config, HealthCheck, COMMON_EDITORS, DATE_FORMATS } from './types';
import { Loader2, Rocket, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';

export default function GeneralSettingsPage() {
  const [config, setConfig] = useState<Config | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [configExists, setConfigExists] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<Config>({});
  const [tagInput, setTagInput] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<{type: 'success' | 'error', text: string} | null>(null);
  const [systemHealth, setSystemHealth] = useState<HealthCheck[]>([]);
  const [isLoadingHealth, setIsLoadingHealth] = useState(false);
  const [installingDeps, setInstallingDeps] = useState<Set<string>>(new Set());

  useEffect(() => {
    loadConfig();
    loadSystemHealth();
  }, []);

  const loadConfig = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('/api/config');
      const data = await response.json();
      
      if (data.success) {
        setConfig(data.config);
        setFormData(data.config);
        setConfigExists(true);
      } else if (response.status === 404 && data.error === 'Configuration file not found') {
        // Config file doesn't exist - this is a setup scenario, not an error
        setConfig(null);
        setConfigExists(false);
      } else {
        // Actual error reading config
        setConfig(null);
        setConfigExists(true); // File exists but has issues
      }
    } catch (error) {
      console.error('Failed to load config:', error);
      setConfig(null);
      setConfigExists(true); // Assume file exists but network/other issues
    } finally {
      setIsLoading(false);
    }
  };

  const loadSystemHealth = async () => {
    setIsLoadingHealth(true);
    try {
      const response = await fetch('/api/system-health');
      const data = await response.json();
      if (data.success) {
        setSystemHealth(data.checks);
      } else {
        setSystemHealth([]);
      }
    } catch (error) {
      console.error('Failed to load system health:', error);
      setSystemHealth([]);
    } finally {
      setIsLoadingHealth(false);
    }
  };

  const installDependency = async (dependencyName: string) => {
    setInstallingDeps(prev => new Set(prev).add(dependencyName));
    setMessage(null);
    
    try {
      const response = await fetch('/api/install-dependency', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ dependency: dependencyName }),
      });
      
      const result = await response.json();
      
      if (result.success) {
        setMessage({type: 'success', text: result.message});
        // Refresh system health after successful installation
        await loadSystemHealth();
      } else {
        setMessage({type: 'error', text: result.error});
      }
    } catch (error) {
      console.error('Failed to install dependency:', error);
      setMessage({type: 'error', text: 'Failed to install dependency'});
    } finally {
      setInstallingDeps(prev => {
        const newSet = new Set(prev);
        newSet.delete(dependencyName);
        return newSet;
      });
    }
  };

  const isBrewInstalled = () => {
    return systemHealth.find(check => check.name === 'Homebrew')?.status === 'ok';
  };

  const canInstallDependency = (dependencyName: string) => {
    return isBrewInstalled() && ['FFmpeg', 'FFprobe', 'Google Cloud CLI'].includes(dependencyName);
  };

  const handleInputChange = (field: keyof Config, value: string | string[]) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleAddTag = () => {
    if (tagInput.trim() && !formData.default_tags?.includes(tagInput.trim())) {
      const newTags = [...(formData.default_tags || []), tagInput.trim()];
      handleInputChange('default_tags', newTags);
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    const newTags = formData.default_tags?.filter(tag => tag !== tagToRemove) || [];
    handleInputChange('default_tags', newTags);
  };

  const handleSave = async () => {
    setIsSaving(true);
    setMessage(null);
    try {
      console.log('Sending config data:', formData);
      const response = await fetch('/api/config', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      
      console.log('Response status:', response.status);
      console.log('Response headers:', response.headers);
      
      if (!response.ok) {
        const text = await response.text();
        console.error('Response text:', text);
        setMessage({type: 'error', text: `Server error: ${response.status} - ${text.substring(0, 100)}`});
        return;
      }
      
      const result = await response.json();
      console.log('Response result:', result);
      
      if (result.success) {
        setConfig(formData);
        setConfigExists(true);
        setIsEditing(false);
        setMessage({type: 'success', text: 'Configuration saved successfully!'});
      } else {
        setMessage({type: 'error', text: result.error || 'Failed to save configuration'});
      }
    } catch (error) {
      console.error('Failed to save config:', error);
      if (error instanceof Error) {
        setMessage({type: 'error', text: `Error: ${error.message}`});
      } else {
        setMessage({type: 'error', text: 'Failed to save configuration'});
      }
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    setFormData(config || {});
    setIsEditing(false);
    setMessage(null);
  };


  // Show loading state
  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="flex justify-center mb-4">
          <Loader2 className="h-16 w-16 text-gray-400 animate-spin" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Loading Configuration</h1>
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          Checking configuration file...
        </p>
      </div>
    );
  }

  // Show setup component if config doesn't exist
  if (!config && !configExists) {
    return (
      <div className="text-center py-12">
        <div className="flex justify-center mb-4">
          <Rocket className="h-16 w-16 text-gray-400" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Welcome to Note AI</h1>
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          Let&apos;s get you set up! The CLI configuration file doesn&apos;t exist yet.
        </p>
        
        <div className="max-w-2xl mx-auto bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-8 text-left">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">Getting Started</h2>
          <div className="space-y-4">
            <div className="border-l-4 border-blue-500 pl-4">
              <h3 className="font-medium text-gray-900 dark:text-white mb-2">Step 1: Install the CLI</h3>
              <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                Make sure you have the Note AI CLI installed and available in your PATH.
              </p>
            </div>
            
            <div className="border-l-4 border-blue-500 pl-4">
              <h3 className="font-medium text-gray-900 dark:text-white mb-2">Step 2: Initialize Configuration</h3>
              <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                Run any CLI command to create the initial configuration file:
              </p>
              <code className="block bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-3 py-2 rounded font-mono text-sm">
                note --help
              </code>
            </div>
            
            <div className="border-l-4 border-blue-500 pl-4">
              <h3 className="font-medium text-gray-900 dark:text-white mb-2">Step 3: Refresh This Page</h3>
              <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                Once the configuration file is created, refresh this page to access your settings.
              </p>
            </div>
          </div>
          
          <div className="mt-6 flex justify-center">
            <Button 
              onClick={() => window.location.reload()}
            >
              Refresh Page
            </Button>
          </div>
        </div>
      </div>
    );
  }

  // Show error state if config exists but couldn't be loaded
  if (!config && configExists) {
    return (
      <div className="text-center py-12">
        <div className="flex justify-center mb-4">
          <X className="h-16 w-16 text-red-400" />
        </div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Error</h1>
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          Unable to read the configuration file. Please check if the CLI is set up correctly.
        </p>
        <Button 
          onClick={loadConfig}
        >
          Try Again
        </Button>
      </div>
    );
  }

  return (
    <div>
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">General Settings</h1>
        <p className="text-gray-600 dark:text-gray-300">
          Application configuration loaded from the CLI config file.
        </p>
      </div>

        {message && (
          <Alert className="mt-4 mb-6">
            <AlertDescription>{message.text}</AlertDescription>
          </Alert>
        )}

        <form className="space-y-8" onSubmit={(e) => { e.preventDefault(); handleSave(); }}>
          {/* General Settings */}
          <section>
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">General Settings</h2>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Editor</label>
                  <select
                    value={formData.editor || ''}
                    onChange={(e) => handleInputChange('editor', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                    disabled={!isEditing}
                  >
                    {COMMON_EDITORS.map((editor) => (
                      <option key={editor} value={editor}>
                        {editor}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Date Format</label>
                  <select
                    value={formData.date_format || ''}
                    onChange={(e) => handleInputChange('date_format', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                    disabled={!isEditing}
                  >
                    {DATE_FORMATS.map((format) => (
                      <option key={format} value={format}>
                        {format}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Default Tags</label>
                  <div className="mt-1">
                    {isEditing ? (
                      <div>
                        <div className="flex flex-wrap gap-2 mb-2">
                          {formData.default_tags?.map((tag, index) => (
                            <span
                              key={index}
                              className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200"
                            >
                              {tag}
                              <button
                                type="button"
                                className="ml-1 text-blue-800 dark:text-blue-200"
                                onClick={() => handleRemoveTag(tag)}
                              >
                                ×
                              </button>
                            </span>
                          ))}
                        </div>
                        <div className="flex gap-2">
                          <input
                            type="text"
                            value={tagInput}
                            onChange={(e) => setTagInput(e.target.value)}
                            className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            placeholder="Add a tag"
                            onKeyPress={(e) => e.key === 'Enter' && handleAddTag()}
                          />
                          <Button type="button" size="sm" variant="secondary" onClick={handleAddTag}>
                            Add
                          </Button>
                        </div>
                      </div>
                    ) : formData.default_tags && formData.default_tags.length > 0 ? (
                      <div className="flex flex-wrap gap-2">
                        {formData.default_tags.map((tag, index) => (
                          <span
                            key={index}
                            className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200"
                          >
                            {tag}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <p className="text-sm text-gray-500 dark:text-gray-400 italic">No default tags configured</p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </section>

          {/* File Paths */}
          <section>
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">File Paths</h2>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Notes Directory</label>
                  <input
                    type="text"
                    value={formData.notes_dir || ''}
                    onChange={(e) => handleInputChange('notes_dir', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                    disabled={!isEditing}
                    placeholder="/path/to/notes"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Database Path</label>
                  <input
                    type="text"
                    value={formData.database_path || ''}
                    onChange={(e) => handleInputChange('database_path', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                    disabled={!isEditing}
                    placeholder="/path/to/database.db"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Configuration File</label>
                  <input
                    type="text"
                    value="~/.noteai/config.json"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 font-mono text-sm cursor-not-allowed"
                    readOnly
                  />
                </div>
              </div>
            </div>
          </section>

          {/* Status */}
          <section>
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">System Status</h2>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Configuration</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.database_path ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    Database {formData.database_path ? 'Configured' : 'Missing'}
                  </span>
                </div>
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.notes_dir ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    Notes Directory {formData.notes_dir ? 'Set' : 'Missing'}
                  </span>
                </div>
              </div>
              
              <div className="border-t border-gray-200 dark:border-gray-600 pt-4">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-lg font-medium text-gray-900 dark:text-white">System Dependencies</h3>
                  <Button
                    type="button"
                    onClick={loadSystemHealth}
                    size="sm"
                    variant="outline"
                    disabled={isLoadingHealth}
                  >
                    {isLoadingHealth ? (
                      <span className="flex items-center">
                        <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-current" fill="none" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Checking...
                      </span>
                    ) : (
                      'Refresh'
                    )}
                  </Button>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {systemHealth.length > 0 ? (
                    systemHealth.map((check, index) => (
                      <div key={index} className="flex items-start space-x-3">
                        <div className={`w-3 h-3 rounded-full mt-1 flex-shrink-0 ${
                          check.status === 'ok' ? 'bg-green-500' : 
                          check.status === 'missing' ? 'bg-red-500' : 
                          'bg-yellow-500'
                        }`}></div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium text-gray-800 dark:text-gray-200">
                              {check.name}
                            </span>
                            <div className="flex items-center space-x-2">
                              <span className={`text-xs px-2 py-1 rounded-full ${
                                check.status === 'ok' ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' :
                                check.status === 'missing' ? 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200' :
                                'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200'
                              }`}>
                                {check.status === 'ok' ? 'Installed' :
                                 check.status === 'missing' ? 'Missing' :
                                 'Error'}
                              </span>
                              {check.status === 'missing' && canInstallDependency(check.name) && (
                                <button
                                  type="button"
                                  onClick={() => installDependency(check.name)}
                                  disabled={installingDeps.has(check.name)}
                                  className="text-xs px-2 py-1 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-md transition-colors duration-200 flex items-center space-x-1"
                                >
                                  {installingDeps.has(check.name) ? (
                                    <>
                                      <svg className="animate-spin h-3 w-3" fill="none" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                      </svg>
                                      <span>Installing...</span>
                                    </>
                                  ) : (
                                    <>
                                      <svg className="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                                      </svg>
                                      <span>Install</span>
                                    </>
                                  )}
                                </button>
                              )}
                            </div>
                          </div>
                          {check.version && (
                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1 font-mono truncate">
                              {check.version}
                            </p>
                          )}
                          {check.error && (
                            <p className="text-xs text-red-600 dark:text-red-400 mt-1">
                              {check.error}
                            </p>
                          )}
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="col-span-full text-center py-4">
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {isLoadingHealth ? 'Checking system dependencies...' : 'Click Refresh to check system dependencies'}
                      </p>
                    </div>
                  )}
                </div>
                
                {systemHealth.some(check => check.status === 'missing') && (
                  <div className="mt-4 p-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md">
                    <p className="text-sm text-yellow-800 dark:text-yellow-200">
                      <strong>Missing dependencies detected:</strong> Some features may not work properly.
                    </p>
                    {systemHealth.find(check => check.name === 'Homebrew' && check.status === 'missing') ? (
                      <p className="text-sm text-yellow-800 dark:text-yellow-200 mt-1">
                        Install Homebrew first, then use it to install missing tools. Visit{' '}
                        <a href="https://brew.sh" target="_blank" rel="noopener noreferrer" className="underline hover:no-underline">
                          brew.sh
                        </a>{' '}
                        for installation instructions.
                      </p>
                    ) : isBrewInstalled() && systemHealth.some(check => 
                      check.status === 'missing' && canInstallDependency(check.name)
                    ) ? (
                      <p className="text-sm text-yellow-800 dark:text-yellow-200 mt-1">
                        Click the &quot;Install&quot; button next to missing dependencies to install them using Homebrew.
                      </p>
                    ) : null}
                  </div>
                )}
              </div>
            </div>
          </section>

          <div className="flex justify-end mt-6 space-x-4">
            {isEditing ? (
              <>
                <Button type="button" variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button type="submit" disabled={isSaving}>
                  {isSaving ? 'Saving...' : 'Save Changes'}
                </Button>
              </>
            ) : (
              <Button type="button" onClick={() => setIsEditing(true)}>
                Edit Settings
              </Button>
            )}
          </div>
        </form>
    </div>
  );
}
