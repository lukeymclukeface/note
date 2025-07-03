export interface AppConfig {
  notes_dir?: string;
  editor?: string;
  date_format?: string;
  default_tags?: string[];
  database_path?: string;
  // OpenAI Configuration
  openai_key?: string;
  // Google AI Configuration
  google_project_id?: string;
  google_location?: string;
  // Model/Provider Configuration
  transcription_provider?: string;
  transcription_model?: string;
  summary_provider?: string;
  summary_model?: string;
}

// Fetch configuration from the API
export async function loadConfig(): Promise<AppConfig | null> {
  try {
    const response = await fetch('/api/config');
    const data = await response.json();
    
    if (data.success) {
      return data.config || data.data?.config || null;
    }
    return null;
  } catch (error) {
    console.error('Error loading config:', error);
    return null;
  }
}

// Save configuration via the API
export async function saveConfig(config: AppConfig): Promise<void> {
  try {
    const response = await fetch('/api/config', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(config),
    });
    
    const result = await response.json();
    
    if (!result.success) {
      throw new Error(result.error || 'Failed to save configuration');
    }
  } catch (error) {
    console.error('Error saving config:', error);
    throw error;
  }
}

// Mask the OpenAI API key for security display
export function maskApiKey(key: string | undefined): string {
  if (!key) return '(not set)';
  
  if (key.length > 8) {
    return key.slice(0, 4) + '...' + key.slice(-4);
  } else {
    return '*'.repeat(key.length);
  }
}

// Service object for easier importing
export const configService = {
  getConfig: loadConfig,
  saveConfig,
  maskApiKey
};
