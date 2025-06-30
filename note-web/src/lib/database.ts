import Database from 'better-sqlite3';
import path from 'path';
import os from 'os';

export interface Note {
  id: number;
  title: string;
  content: string;
  tags: string;
  created_at: string;
  updated_at: string;
}

export interface Recording {
  id: number;
  filename: string;
  file_path: string;
  start_time: string;
  end_time: string;
  duration: number;
  file_size: number;
  format: string;
  sample_rate: number;
  channels: number;
  created_at: string;
}

// Get the database path (same as CLI: ~/.noteai/notes.db)
function getDatabasePath(): string {
  const homeDir = os.homedir();
  return path.join(homeDir, '.noteai', 'notes.db');
}

// Initialize database connection
function getDatabase(): Database.Database {
  const dbPath = getDatabasePath();
  const db = new Database(dbPath, { readonly: true });
  return db;
}

// Get all notes from the database
export function getAllNotes(): Note[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, tags, created_at, updated_at 
      FROM notes 
      ORDER BY created_at DESC
    `);
    const notes = stmt.all() as Note[];
    db.close();
    return notes;
  } catch (error) {
    console.error('Error fetching notes:', error);
    return [];
  }
}

// Get notes by tag
export function getNotesByTag(tag: string): Note[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, tags, created_at, updated_at 
      FROM notes 
      WHERE tags LIKE ? 
      ORDER BY created_at DESC
    `);
    const notes = stmt.all(`%${tag}%`) as Note[];
    db.close();
    return notes;
  } catch (error) {
    console.error('Error fetching notes by tag:', error);
    return [];
  }
}

// Get a single note by ID
export function getNoteById(id: number): Note | null {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, tags, created_at, updated_at 
      FROM notes 
      WHERE id = ?
    `);
    const note = stmt.get(id) as Note | undefined;
    db.close();
    return note || null;
  } catch (error) {
    console.error('Error fetching note by ID:', error);
    return null;
  }
}

// Get all recordings from the database
export function getAllRecordings(): Recording[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, filename, file_path, start_time, end_time, duration, 
             file_size, format, sample_rate, channels, created_at 
      FROM recordings 
      ORDER BY created_at DESC
    `);
    const recordings = stmt.all() as Recording[];
    db.close();
    return recordings;
  } catch (error) {
    console.error('Error fetching recordings:', error);
    return [];
  }
}
