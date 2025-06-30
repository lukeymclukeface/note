import { loadConfig, getConfigFilePath, configExists, maskApiKey } from '@/lib/config';

export default function SettingsPage() {
  const config = loadConfig();
  const configPath = getConfigFilePath();
  const hasConfig = configExists();

  if (!hasConfig) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="max-w-4xl mx-auto p-6">
          <div className="text-center py-12">
            <div className="text-6xl mb-4">⚙️</div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Not Found</h1>
            <p className="text-gray-600 dark:text-gray-300 mb-8">
              No configuration file found. Please run the CLI setup first:
            </p>
            <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 text-left max-w-md mx-auto">
              <code className="text-sm font-mono text-gray-900 dark:text-gray-100">
                cd note-cli<br />
                ./note setup
              </code>
            </div>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-4">
              Expected config location: <code className="text-xs bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">{configPath}</code>
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (!config) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="max-w-4xl mx-auto p-6">
          <div className="text-center py-12">
            <div className="text-6xl mb-4">❌</div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Configuration Error</h1>
            <p className="text-gray-600 dark:text-gray-300 mb-8">
              Unable to read the configuration file. Please check file permissions.
            </p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Config file: <code className="text-xs bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-2 py-1 rounded">{configPath}</code>
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

      {/* Configuration Display */}
      <div className="space-y-8">
        {/* General Settings */}
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">General Settings</h2>
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Editor</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600">
                  {config.editor || 'nano'}
                </p>
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Date Format</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600">
                  {config.date_format || '2006-01-02'}
                </p>
              </div>
              
              <div className="md:col-span-2">
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Default Tags</label>
                <div className="mt-1">
                  {config.default_tags && config.default_tags.length > 0 ? (
                    <div className="flex flex-wrap gap-2">
                      {config.default_tags.map((tag, index) => (
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
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">OpenAI API Key</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600 font-mono">
                  {maskApiKey(config.openai_key)}
                </p>
                {!config.openai_key && (
                  <p className="text-xs text-red-600 dark:text-red-400 mt-1">⚠️ API key not configured</p>
                )}
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Transcription Model</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600">
                  {config.transcription_model || 'whisper-1'}
                </p>
              </div>
              
              <div className="md:col-span-2">
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Summary Model</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600">
                  {config.summary_model || 'gpt-3.5-turbo'}
                </p>
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
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Notes Directory</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600 font-mono text-xs">
                  {config.notes_dir || '(not set)'}
                </p>
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Database Path</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600 font-mono text-xs">
                  {config.database_path || '(not set)'}
                </p>
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-600 dark:text-gray-300">Configuration File</label>
                <p className="text-sm text-gray-900 dark:text-gray-100 bg-gray-50 dark:bg-gray-700 px-3 py-2 rounded border border-gray-200 dark:border-gray-600 font-mono text-xs">
                  {configPath}
                </p>
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
                <div className={`w-3 h-3 rounded-full mr-3 ${config.openai_key ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  OpenAI API Key {config.openai_key ? 'Configured' : 'Missing'}
                </span>
              </div>
              
              <div className="flex items-center">
                <div className={`w-3 h-3 rounded-full mr-3 ${config.database_path ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  Database {config.database_path ? 'Configured' : 'Missing'}
                </span>
              </div>
              
              <div className="flex items-center">
                <div className={`w-3 h-3 rounded-full mr-3 ${config.notes_dir ? 'bg-green-500' : 'bg-red-500'}`}></div>
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  Notes Directory {config.notes_dir ? 'Set' : 'Missing'}
                </span>
              </div>
            </div>
          </div>
        </section>

        {/* Instructions */}
        <section>
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">Configuration Help</h2>
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
            <div className="space-y-4">
              <div>
                <h3 className="text-sm font-semibold text-blue-900 dark:text-blue-200">To modify settings:</h3>
                <p className="text-sm text-blue-800 dark:text-blue-300 mt-1">
                  Use the CLI configuration commands to update these settings.
                </p>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <h4 className="font-medium text-blue-900 dark:text-blue-200">View Config:</h4>
                  <code className="block bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-2 py-1 rounded mt-1 font-mono text-xs">
                    ./note config
                  </code>
                </div>
                
                <div>
                  <h4 className="font-medium text-blue-900 dark:text-blue-200">Set API Key:</h4>
                  <code className="block bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-2 py-1 rounded mt-1 font-mono text-xs">
                    ./note config set openai_key sk-...
                  </code>
                </div>
                
                <div>
                  <h4 className="font-medium text-blue-900 dark:text-blue-200">Configure Models:</h4>
                  <code className="block bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-2 py-1 rounded mt-1 font-mono text-xs">
                    ./note config model
                  </code>
                </div>
                
                <div>
                  <h4 className="font-medium text-blue-900 dark:text-blue-200">Run Setup:</h4>
                  <code className="block bg-blue-100 dark:bg-blue-800 text-blue-900 dark:text-blue-100 px-2 py-1 rounded mt-1 font-mono text-xs">
                    ./note setup
                  </code>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
      </div>
    </div>
  );
}

// Generate metadata for the page
export function generateMetadata() {
  return {
    title: 'Settings - Note AI',
    description: 'View and manage Note AI application configuration',
  };
}
