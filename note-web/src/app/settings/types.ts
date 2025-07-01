export type Config = {
  editor?: string;
  date_format?: string;
  default_tags?: string[];
  notes_dir?: string;
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
};

export interface HealthCheck {
  name: string;
  status: 'ok' | 'missing' | 'error';
  version?: string;
  error?: string;
}

// Client-side function to mask API keys
export const maskApiKey = (key: string | undefined): string => {
  if (!key) return '(not set)';
  if (key.length <= 8) return '***';
  return key.substring(0, 4) + '***' + key.substring(key.length - 4);
};

export const AI_PROVIDERS = ['openai', 'google'];
export const OPENAI_TRANSCRIPTION_MODELS = ['whisper-1'];
export const OPENAI_SUMMARY_MODELS = ['gpt-3.5-turbo', 'gpt-4', 'gpt-4-turbo', 'gpt-4o', 'gpt-4o-mini'];
export const GOOGLE_TRANSCRIPTION_MODELS = ['chirp', 'chirp_2', 'gpt-4o-transcribe'];
export const GOOGLE_SUMMARY_MODELS = ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-1.0-pro'];
export const GOOGLE_LOCATIONS = ['us-central1', 'us-east1', 'us-west1', 'europe-west1', 'asia-southeast1'];
export const COMMON_EDITORS = ['nano', 'vim', 'emacs', 'code', 'subl'];
export const DATE_FORMATS = ['2006-01-02', '01/02/2006', '02-01-2006', 'Jan 2, 2006'];
