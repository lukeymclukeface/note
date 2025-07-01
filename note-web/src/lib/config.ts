import fs from 'fs';
import path from 'path';
import os from 'os';

export interface AppConfig {
  notes_dir: string;
  editor: string;
  date_format: string;
  default_tags: string[];
  database_path: string;
  // OpenAI Configuration
  openai_key: string;
  // Google AI Configuration
  google_project_id: string;
  google_location: string;
  // Model/Provider Configuration
  transcription_provider: string;
  transcription_model: string;
  summary_provider: string;
  summary_model: string;
}

// Get the config file path (same as CLI: ~/.noteai/config.json)
function getConfigPath(): string {
  const homeDir = os.homedir();
  return path.join(homeDir, '.noteai', 'config.json');
}

// Load configuration from the CLI config file
export function loadConfig(): AppConfig | null {
  try {
    const configPath = getConfigPath();
    
    // Check if config file exists
    if (!fs.existsSync(configPath)) {
      return null;
    }

    // Read and parse config file
    const configData = fs.readFileSync(configPath, 'utf8');
    const config = JSON.parse(configData);

    return {
      notes_dir: config.notes_dir || '',
      editor: config.editor || 'nano',
      date_format: config.date_format || '2006-01-02',
      default_tags: config.default_tags || [],
      openai_key: config.openai_key || '',
      database_path: config.database_path || '',
      transcription_model: config.transcription_model || 'whisper-1',
      summary_model: config.summary_model || 'gpt-3.5-turbo',
      transcription_provider: config.transcription_provider || 'openai',
      summary_provider: config.summary_provider || 'openai',
      google_project_id: config.google_project_id || '',
      google_location: config.google_location || 'us-central1',
    };
  } catch (error) {
    console.error('Error loading config:', error);
    return null;
  }
}

// Get config file path for display
export function getConfigFilePath(): string {
  return getConfigPath();
}

// Check if config file exists
export function configExists(): boolean {
  try {
    return fs.existsSync(getConfigPath());
  } catch {
    return false;
  }
}

// Save configuration to file
export function saveConfig(config: AppConfig): void {
  try {
    const configPath = getConfigPath();
    const configDir = path.dirname(configPath);
    
    // Create config directory if it doesn't exist
    if (!fs.existsSync(configDir)) {
      fs.mkdirSync(configDir, { recursive: true });
    }
    
    // Write config to file
    fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
  } catch (error) {
    console.error('Error saving config:', error);
    throw error;
  }
}

// Mask the OpenAI API key for security display
export function maskApiKey(key: string): string {
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
  getConfigPath,
  configExists,
  maskApiKey
};
