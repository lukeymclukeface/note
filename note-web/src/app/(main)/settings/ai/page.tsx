'use client';

import { useState, useEffect } from 'react';
import { X, AlertTriangle, Search, Loader2 } from 'lucide-react';
import { Config, maskApiKey, AI_PROVIDERS, OPENAI_TRANSCRIPTION_MODELS, OPENAI_SUMMARY_MODELS, GOOGLE_TRANSCRIPTION_MODELS, GOOGLE_SUMMARY_MODELS, GOOGLE_LOCATIONS } from '../types';

interface OpenAIValidationResult {
  valid: boolean;
  error?: string;
  details?: {
    organizationId?: string;
    models?: string[];
  };
}

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

export default function AISettingsPage() {
  const [config, setConfig] = useState<Config | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<Config>({});
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<{type: 'success' | 'error', text: string} | null>(null);
  const [isValidatingOpenAI, setIsValidatingOpenAI] = useState(false);
  const [openAIValidation, setOpenAIValidation] = useState<OpenAIValidationResult | null>(null);
  const [isValidatingGoogleAI, setIsValidatingGoogleAI] = useState(false);
  const [googleAIValidation, setGoogleAIValidation] = useState<GoogleAIValidationResult | null>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await fetch('/api/config');
      const data = await response.json();
      if (data.success) {
        setConfig(data.config);
        setFormData(data.config);
      } else {
        setConfig(null);
      }
    } catch (error) {
      console.error('Failed to load config:', error);
      setConfig(null);
    }
  };

  const getTranscriptionModels = () => {
    return formData.transcription_provider === 'google' ? GOOGLE_TRANSCRIPTION_MODELS : OPENAI_TRANSCRIPTION_MODELS;
  };

  const getSummaryModels = () => {
    return formData.summary_provider === 'google' ? GOOGLE_SUMMARY_MODELS : OPENAI_SUMMARY_MODELS;
  };

  const handleInputChange = (field: keyof Config, value: string | string[]) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSave = async () => {
    setIsSaving(true);
    setMessage(null);
    try {
      const response = await fetch('/api/config', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      
      if (!response.ok) {
        const text = await response.text();
        setMessage({type: 'error', text: `Server error: ${response.status} - ${text.substring(0, 100)}`});
        return;
      }
      
      const result = await response.json();
      
      if (result.success) {
        setConfig(formData);
        setIsEditing(false);
        setMessage({type: 'success', text: 'AI configuration saved successfully!'});
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

  const validateOpenAIKey = async () => {
    if (!formData.openai_key) {
      setMessage({type: 'error', text: 'No OpenAI API key to validate'});
      return;
    }

    setIsValidatingOpenAI(true);
    setMessage(null);
    setOpenAIValidation(null);
    
    try {
      const response = await fetch('/api/validate-openai', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ apiKey: formData.openai_key }),
      });

      if (!response.ok) {
        const text = await response.text();
        setMessage({type: 'error', text: `Validation failed: ${response.status} - ${text.substring(0, 100)}`});
        return;
      }

      const result = await response.json();

      if (result.success) {
        setOpenAIValidation(result.validation);
        if (result.validation.valid) {
          setMessage({type: 'success', text: 'OpenAI API key is valid and working!'});
        } else {
          setMessage({type: 'error', text: `OpenAI API key validation failed: ${result.validation.error}`});
        }
      } else {
        setMessage({type: 'error', text: result.error || 'Failed to validate API key'});
      }
    } catch (error) {
      console.error('Failed to validate OpenAI key:', error);
      setMessage({type: 'error', text: 'Failed to validate API key'});
    } finally {
      setIsValidatingOpenAI(false);
    }
  };

  const validateGoogleAI = async () => {
    if (!formData.google_project_id) {
      setMessage({type: 'error', text: 'No Google Project ID to validate'});
      return;
    }

    setIsValidatingGoogleAI(true);
    setMessage(null);
    setGoogleAIValidation(null);
    
    try {
      const response = await fetch('/api/validate-google-ai', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          projectId: formData.google_project_id,
          location: formData.google_location 
        }),
      });

      if (!response.ok) {
        const text = await response.text();
        setMessage({type: 'error', text: `Validation failed: ${response.status} - ${text.substring(0, 100)}`});
        return;
      }

      const result = await response.json();

      if (result.success) {
        setGoogleAIValidation(result.validation);
        if (result.validation.valid) {
          setMessage({type: 'success', text: 'Google AI configuration is valid and working!'});
        } else if (result.validation.needsAuth) {
          setMessage({type: 'error', text: `Google AI validation failed: ${result.validation.error}. Please run "gcloud auth login" in your terminal.`});
        } else {
          setMessage({type: 'error', text: `Google AI validation failed: ${result.validation.error}`});
        }
      } else {
        setMessage({type: 'error', text: result.error || 'Failed to validate Google AI configuration'});
      }
    } catch (error) {
      console.error('Failed to validate Google AI:', error);
      setMessage({type: 'error', text: 'Failed to validate Google AI configuration'});
    } finally {
      setIsValidatingGoogleAI(false);
    }
  };

  if (!config) {
    return (
      <div className="text-center py-12">
        <div className="flex justify-center mb-4">
          <X size={64} className="text-red-500" />
        </div>
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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">AI Settings</h1>
        <p className="text-gray-600 dark:text-gray-300">
          Configure AI providers, models, and API credentials for transcription and summarization.
        </p>
      </div>

      {message && (
        <div className={`alert mb-6 ${message.type === 'success' ? 'alert-success' : 'alert-error'}`}>
          {message.text}
        </div>
      )}

      <form className="space-y-8" onSubmit={(e) => { e.preventDefault(); handleSave(); }}>
        {/* OpenAI Configuration */}
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">OpenAI Configuration</h2>
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">OpenAI API Key</label>
                <input
                  type={isEditing ? "text" : "password"}
                  value={isEditing ? (formData.openai_key || '') : maskApiKey(formData.openai_key)}
                  onChange={(e) => handleInputChange('openai_key', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 font-mono focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing}
                  placeholder={isEditing ? "Enter your OpenAI API key" : ""}
                />
                {!formData.openai_key && (
                  <p className="text-xs text-red-600 dark:text-red-400 mt-1 flex items-center">
                    <AlertTriangle size={12} className="mr-1" />
                    API key not configured
                  </p>
                )}
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Get your API key from <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="underline">OpenAI Platform</a>
                </p>
                
                {/* Validation Button and Results */}
                {formData.openai_key && (
                  <div className="mt-4">
                    <button
                      type="button"
                      onClick={validateOpenAIKey}
                      disabled={isValidatingOpenAI}
                      className="btn btn-sm btn-outline"
                    >
                      {isValidatingOpenAI ? (
                        <span className="flex items-center">
                          <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                          Validating...
                        </span>
                      ) : (
                        <span className="flex items-center">
                          <Search size={16} className="mr-2" />
                          Validate API Key
                        </span>
                      )}
                    </button>
                    
                    {/* Validation Results */}
                    {openAIValidation && (
                      <div className={`mt-3 p-3 rounded-md border ${
                        openAIValidation.valid 
                          ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800'
                          : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800'
                      }`}>
                        <div className="flex items-center mb-2">
                          <div className={`w-3 h-3 rounded-full mr-2 ${
                            openAIValidation.valid ? 'bg-green-500' : 'bg-red-500'
                          }`}></div>
                          <span className={`text-sm font-medium ${
                            openAIValidation.valid 
                              ? 'text-green-800 dark:text-green-200'
                              : 'text-red-800 dark:text-red-200'
                          }`}>
                            {openAIValidation.valid ? 'API Key Valid' : 'API Key Invalid'}
                          </span>
                        </div>
                        
                        {openAIValidation.valid && openAIValidation.details && (
                          <div className="text-xs space-y-1">
                            {openAIValidation.details.organizationId && (
                              <p className="text-green-700 dark:text-green-300">
                                <span className="font-medium">Organization:</span> {openAIValidation.details.organizationId}
                              </p>
                            )}
                            {openAIValidation.details.models && openAIValidation.details.models.length > 0 && (
                              <p className="text-green-700 dark:text-green-300">
                                <span className="font-medium">Available models:</span> {openAIValidation.details.models.slice(0, 5).join(', ')}
                                {openAIValidation.details.models.length > 5 && ` (+${openAIValidation.details.models.length - 5} more)`}
                              </p>
                            )}
                          </div>
                        )}
                        
                        {!openAIValidation.valid && openAIValidation.error && (
                          <p className="text-xs text-red-700 dark:text-red-300">
                            <span className="font-medium">Error:</span> {openAIValidation.error}
                          </p>
                        )}
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </section>

        {/* Google AI Configuration */}
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">Google AI Configuration</h2>
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Google Project ID</label>
                <input
                  type="text"
                  value={formData.google_project_id || ''}
                  onChange={(e) => handleInputChange('google_project_id', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing}
                  placeholder="your-google-project-id"
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Find your project ID in the <a href="https://console.cloud.google.com/" target="_blank" rel="noopener noreferrer" className="underline">Google Cloud Console</a>
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Google Location</label>
                <select
                  value={formData.google_location || ''}
                  onChange={(e) => handleInputChange('google_location', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing}
                >
                  <option value="">Select a location...</option>
                  {GOOGLE_LOCATIONS.map((location) => (
                    <option key={location} value={location}>
                      {location}
                    </option>
                  ))}
                </select>
              </div>
              
              {/* Google AI Validation */}
              {formData.google_project_id && (
                <div className="md:col-span-2 mt-4">
                  <button
                    type="button"
                    onClick={validateGoogleAI}
                    disabled={isValidatingGoogleAI}
                    className="btn btn-sm btn-outline"
                  >
                    {isValidatingGoogleAI ? (
                      <span className="flex items-center">
                        <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                        Validating...
                      </span>
                    ) : (
                      <span className="flex items-center">
                        <Search size={16} className="mr-2" />
                        Validate Google AI
                      </span>
                    )}
                  </button>
                  
                  {/* Google AI Validation Results */}
                  {googleAIValidation && (
                    <div className={`mt-3 p-3 rounded-md border ${
                      googleAIValidation.valid 
                        ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800'
                        : googleAIValidation.needsAuth
                        ? 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800'
                        : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800'
                    }`}>
                      <div className="flex items-center mb-2">
                        <div className={`w-3 h-3 rounded-full mr-2 ${
                          googleAIValidation.valid 
                            ? 'bg-green-500' 
                            : googleAIValidation.needsAuth 
                            ? 'bg-yellow-500' 
                            : 'bg-red-500'
                        }`}></div>
                        <span className={`text-sm font-medium ${
                          googleAIValidation.valid 
                            ? 'text-green-800 dark:text-green-200'
                            : googleAIValidation.needsAuth
                            ? 'text-yellow-800 dark:text-yellow-200'
                            : 'text-red-800 dark:text-red-200'
                        }`}>
                          {googleAIValidation.valid 
                            ? 'Google AI Valid' 
                            : googleAIValidation.needsAuth 
                            ? 'Authentication Required'
                            : 'Google AI Invalid'}
                        </span>
                      </div>
                      
                      {googleAIValidation.valid && googleAIValidation.details && (
                        <div className="text-xs space-y-1">
                          {googleAIValidation.details.account && (
                            <p className="text-green-700 dark:text-green-300">
                              <span className="font-medium">Account:</span> {googleAIValidation.details.account}
                            </p>
                          )}
                          {googleAIValidation.details.gcloudVersion && (
                            <p className="text-green-700 dark:text-green-300">
                              <span className="font-medium">gcloud CLI:</span> {googleAIValidation.details.gcloudVersion}
                            </p>
                          )}
                          {googleAIValidation.details.services && googleAIValidation.details.services.length > 0 && (
                            <p className="text-green-700 dark:text-green-300">
                              <span className="font-medium">Enabled APIs:</span> {googleAIValidation.details.services.join(', ')}
                            </p>
                          )}
                        </div>
                      )}
                      
                      {googleAIValidation.needsAuth && (
                        <div className="text-xs space-y-2">
                          <p className="text-yellow-700 dark:text-yellow-300">
                            <span className="font-medium">Action Required:</span> Run the following command in your terminal:
                          </p>
                          <code className="block bg-yellow-100 dark:bg-yellow-800/30 text-yellow-800 dark:text-yellow-200 px-2 py-1 rounded text-xs font-mono">
                            gcloud auth login
                          </code>
                        </div>
                      )}
                      
                      {!googleAIValidation.valid && !googleAIValidation.needsAuth && googleAIValidation.error && (
                        <p className="text-xs text-red-700 dark:text-red-300">
                          <span className="font-medium">Error:</span> {googleAIValidation.error}
                        </p>
                      )}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </section>

        {/* Model/Provider Configuration */}
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">Model & Provider Configuration</h2>
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Transcription Provider</label>
                <select
                  value={formData.transcription_provider || ''}
                  onChange={(e) => handleInputChange('transcription_provider', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing}
                >
                  <option value="">Select a provider...</option>
                  {AI_PROVIDERS.map((provider) => (
                    <option key={provider} value={provider}>
                      {provider.charAt(0).toUpperCase() + provider.slice(1)}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Transcription Model</label>
                <select
                  value={formData.transcription_model || ''}
                  onChange={(e) => handleInputChange('transcription_model', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing || !formData.transcription_provider}
                >
                  <option value="">Select a model...</option>
                  {getTranscriptionModels().map((model) => (
                    <option key={model} value={model}>
                      {model}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Summary Provider</label>
                <select
                  value={formData.summary_provider || ''}
                  onChange={(e) => handleInputChange('summary_provider', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing}
                >
                  <option value="">Select a provider...</option>
                  {AI_PROVIDERS.map((provider) => (
                    <option key={provider} value={provider}>
                      {provider.charAt(0).toUpperCase() + provider.slice(1)}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Summary Model</label>
                <select
                  value={formData.summary_model || ''}
                  onChange={(e) => handleInputChange('summary_model', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                  disabled={!isEditing || !formData.summary_provider}
                >
                  <option value="">Select a model...</option>
                  {getSummaryModels().map((model) => (
                    <option key={model} value={model}>
                      {model}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            
            {/* Configuration Status */}
            <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-600">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Configuration Status</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.openai_key ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    OpenAI API Key {formData.openai_key ? 'Configured' : 'Missing'}
                  </span>
                </div>
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.google_project_id ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    Google Project {formData.google_project_id ? 'Configured' : 'Missing'}
                  </span>
                </div>
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.transcription_provider && formData.transcription_model ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    Transcription {formData.transcription_provider && formData.transcription_model ? 'Configured' : 'Incomplete'}
                  </span>
                </div>
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.summary_provider && formData.summary_model ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    Summary {formData.summary_provider && formData.summary_model ? 'Configured' : 'Incomplete'}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div className="flex justify-end mt-6 space-x-4">
          {isEditing ? (
            <>
              <button type="button" className="btn btn-outline" onClick={handleCancel}>
                Cancel
              </button>
              <button type="submit" className={`btn btn-success ${isSaving ? 'loading' : ''}`} disabled={isSaving}>
                Save Changes
              </button>
            </>
          ) : (
            <button type="button" className="btn btn-primary" onClick={() => setIsEditing(true)}>
              Edit AI Settings
            </button>
          )}
        </div>
      </form>
    </div>
  );
}
