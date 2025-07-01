'use client';

import { useState, useEffect } from 'react';
// Client-side function to mask API keys
const maskApiKey = (key: string | undefined): string => {
  if (!key) return '(not set)';
  if (key.length <= 8) return '***';
  return key.substring(0, 4) + '***' + key.substring(key.length - 4);
};

type Config = {
  editor?: string;
  date_format?: string;
  default_tags?: string[];
  openai_key?: string;
  transcription_provider?: string;
  transcription_model?: string;
  summary_provider?: string;
  summary_model?: string;
  google_project_id?: string;
  google_location?: string;
  notes_dir?: string;
  database_path?: string;
};

const AI_PROVIDERS = ['openai', 'google'];
const OPENAI_TRANSCRIPTION_MODELS = ['whisper-1'];
const OPENAI_SUMMARY_MODELS = ['gpt-3.5-turbo', 'gpt-4', 'gpt-4-turbo', 'gpt-4o', 'gpt-4o-mini'];
const GOOGLE_TRANSCRIPTION_MODELS = ['chirp', 'chirp_2', 'gpt-4o-transcribe'];
const GOOGLE_SUMMARY_MODELS = ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-1.0-pro'];
const GOOGLE_LOCATIONS = ['us-central1', 'us-east1', 'us-west1', 'europe-west1', 'asia-southeast1'];
const COMMON_EDITORS = ['nano', 'vim', 'emacs', 'code', 'subl'];
const DATE_FORMATS = ['2006-01-02', '01/02/2006', '02-01-2006', 'Jan 2, 2006'];

export default function SettingsPage() {
  const [config, setConfig] = useState<Config | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<Config>({});
  const [tagInput, setTagInput] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<{type: 'success' | 'error', text: string} | null>(null);

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

  const getTranscriptionModels = () => {
    return formData.transcription_provider === 'google' ? GOOGLE_TRANSCRIPTION_MODELS : OPENAI_TRANSCRIPTION_MODELS;
  };

  const getSummaryModels = () => {
    return formData.summary_provider === 'google' ? GOOGLE_SUMMARY_MODELS : OPENAI_SUMMARY_MODELS;
  };

  if (!config) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="max-w-4xl mx-auto p-6">
          <div className="text-center py-12">
            <div className="text-6xl mb-4">❌</div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Error</h1>
            <p className="text-gray-600 dark:text-gray-300 mb-8">
              Unable to read the configuration file. Please check if the CLI is set up.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-4xl mx-auto p-6">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Settings</h1>
          <p className="text-gray-600 dark:text-gray-300">
            Application configuration loaded from the CLI config file.
          </p>
        </div>

        {message && (
          <div className={`alert mt-4 mb-6 ${message.type === 'success' ? 'alert-success' : 'alert-error'}`}>
            {message.text}
          </div>
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
                          <button type="button" className="btn btn-sm btn-secondary" onClick={handleAddTag}>
                            Add
                          </button>
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

          {/* AI Configuration */}
          <section>
            <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">AI Configuration</h2>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
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
                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">⚠️ API key not configured</p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Transcription Provider</label>
                  <select
                    value={formData.transcription_provider || ''}
                    onChange={(e) => handleInputChange('transcription_provider', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                    disabled={!isEditing}
                  >
                    {AI_PROVIDERS.map((provider) => (
                      <option key={provider} value={provider}>
                        {provider}
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
                    disabled={!isEditing}
                  >
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
                    {AI_PROVIDERS.map((provider) => (
                      <option key={provider} value={provider}>
                        {provider}
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
                    disabled={!isEditing}
                  >
                    {getSummaryModels().map((model) => (
                      <option key={model} value={model}>
                        {model}
                      </option>
                    ))}
                  </select>
                </div>
                {(formData.transcription_provider === 'google' || formData.summary_provider === 'google') && (
                  <>
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
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Google Location</label>
                      <select
                        value={formData.google_location || ''}
                        onChange={(e) => handleInputChange('google_location', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-500 dark:disabled:text-gray-400"
                        disabled={!isEditing}
                      >
                        {GOOGLE_LOCATIONS.map((location) => (
                          <option key={location} value={location}>
                            {location}
                          </option>
                        ))}
                      </select>
                    </div>
                  </>
                )}
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
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="flex items-center">
                  <div className={`w-3 h-3 rounded-full mr-3 ${formData.openai_key ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span className="text-sm text-gray-800 dark:text-gray-200">
                    OpenAI API Key {formData.openai_key ? 'Configured' : 'Missing'}
                  </span>
                </div>
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
                Edit Settings
              </button>
            )}
          </div>
        </form>
      </div>
    </div>
  );
}
