import Database from 'better-sqlite3';
import path from 'path';
import os from 'os';

export interface Note {
  id: number;
  title: string;
  content: string;
  summary: string;
  tags: string;
  recording_id?: number;
  created_at: string;
  updated_at: string;
}

export interface Meeting {
  id: number;
  title: string;
  content: string;
  summary: string;
  attendees: string;
  location: string;
  tags: string;
  recording_id?: number;
  meeting_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Interview {
  id: number;
  title: string;
  content: string;
  summary: string;
  interviewee: string;
  interviewer: string;
  company: string;
  position: string;
  tags: string;
  recording_id?: number;
  interview_date?: string;
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

export interface CalendarEvent {
  id: number;
  title: string;
  start: Date;
  end: Date;
  type: 'recording' | 'note';
  duration: number; // in minutes
  content?: string;
  tags?: string;
  filename?: string;
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
      SELECT id, title, content, summary, tags, recording_id, created_at, updated_at 
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

// Get all meetings from the database
export function getAllMeetings(): Meeting[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, attendees, location, tags, 
             recording_id, meeting_date, created_at, updated_at 
      FROM meetings 
      ORDER BY created_at DESC
    `);
    const meetings = stmt.all() as Meeting[];
    db.close();
    return meetings;
  } catch (error) {
    console.error('Error fetching meetings:', error);
    return [];
  }
}

// Get all interviews from the database
export function getAllInterviews(): Interview[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, interviewee, interviewer, 
             company, position, tags, recording_id, interview_date, 
             created_at, updated_at 
      FROM interviews 
      ORDER BY created_at DESC
    `);
    const interviews = stmt.all() as Interview[];
    db.close();
    return interviews;
  } catch (error) {
    console.error('Error fetching interviews:', error);
    return [];
  }
}

// Get notes by tag
export function getNotesByTag(tag: string): Note[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, tags, recording_id, created_at, updated_at 
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
      SELECT id, title, content, summary, tags, recording_id, created_at, updated_at 
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

// Get a single meeting by ID
export function getMeetingById(id: number): Meeting | null {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, attendees, location, tags, 
             recording_id, meeting_date, created_at, updated_at 
      FROM meetings 
      WHERE id = ?
    `);
    const meeting = stmt.get(id) as Meeting | undefined;
    db.close();
    return meeting || null;
  } catch (error) {
    console.error('Error fetching meeting by ID:', error);
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

// Get recordings for a specific date range
export function getRecordingsInRange(startDate: Date, endDate: Date): Recording[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, filename, file_path, start_time, end_time, duration, 
             file_size, format, sample_rate, channels, created_at 
      FROM recordings 
      WHERE start_time >= ? AND start_time <= ?
      ORDER BY start_time ASC
    `);
    const recordings = stmt.all(startDate.toISOString(), endDate.toISOString()) as Recording[];
    db.close();
    return recordings;
  } catch (error) {
    console.error('Error fetching recordings in range:', error);
    return [];
  }
}

// Convert recordings to calendar events
export function getCalendarEvents(startDate: Date, endDate: Date): CalendarEvent[] {
  try {
    const recordings = getRecordingsInRange(startDate, endDate);
    
    return recordings.map(recording => {
      const start = new Date(recording.start_time);
      const end = new Date(recording.end_time);
      const durationMinutes = Math.round(recording.duration / (1000 * 1000 * 1000 * 60)); // Convert nanoseconds to minutes
      
      return {
        id: recording.id,
        title: recording.filename.replace(/\.[^/.]+$/, ''), // Remove file extension
        start,
        end,
        type: 'recording' as const,
        duration: durationMinutes,
        filename: recording.filename,
      };
    });
  } catch (error) {
    console.error('Error fetching calendar events:', error);
    return [];
  }
}

// Get events for a specific week
export function getWeekEvents(weekStart: Date): CalendarEvent[] {
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekStart.getDate() + 7);
  weekEnd.setHours(23, 59, 59, 999);
  
  return getCalendarEvents(weekStart, weekEnd);
}

// Get recent notes (last N items)
export function getRecentNotes(limit: number = 5): Note[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, tags, recording_id, created_at, updated_at 
      FROM notes 
      ORDER BY created_at DESC
      LIMIT ?
    `);
    const notes = stmt.all(limit) as Note[];
    db.close();
    return notes;
  } catch (error) {
    console.error('Error fetching recent notes:', error);
    return [];
  }
}

// Get recent meetings (last N items)
export function getRecentMeetings(limit: number = 5): Meeting[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, attendees, location, tags, 
             recording_id, meeting_date, created_at, updated_at 
      FROM meetings 
      ORDER BY created_at DESC
      LIMIT ?
    `);
    const meetings = stmt.all(limit) as Meeting[];
    db.close();
    return meetings;
  } catch (error) {
    console.error('Error fetching recent meetings:', error);
    return [];
  }
}

// Get recent interviews (last N items)
export function getRecentInterviews(limit: number = 5): Interview[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, title, content, summary, interviewee, interviewer, 
             company, position, tags, recording_id, interview_date, 
             created_at, updated_at 
      FROM interviews 
      ORDER BY created_at DESC
      LIMIT ?
    `);
    const interviews = stmt.all(limit) as Interview[];
    db.close();
    return interviews;
  } catch (error) {
    console.error('Error fetching recent interviews:', error);
    return [];
  }
}

// Get recent recordings (last N items)
export function getRecentRecordings(limit: number = 5): Recording[] {
  try {
    const db = getDatabase();
    const stmt = db.prepare(`
      SELECT id, filename, file_path, start_time, end_time, duration, 
             file_size, format, sample_rate, channels, created_at 
      FROM recordings 
      ORDER BY created_at DESC
      LIMIT ?
    `);
    const recordings = stmt.all(limit) as Recording[];
    db.close();
    return recordings;
  } catch (error) {
    console.error('Error fetching recent recordings:', error);
    return [];
  }
}

// Get summary statistics
export function getDashboardStats() {
  try {
    const db = getDatabase();
    
    const notesCount = db.prepare('SELECT COUNT(*) as count FROM notes').get() as { count: number };
    const meetingsCount = db.prepare('SELECT COUNT(*) as count FROM meetings').get() as { count: number };
    const interviewsCount = db.prepare('SELECT COUNT(*) as count FROM interviews').get() as { count: number };
    const recordingsCount = db.prepare('SELECT COUNT(*) as count FROM recordings').get() as { count: number };
    
    const totalDuration = db.prepare('SELECT SUM(duration) as total FROM recordings').get() as { total: number | null };
    
    db.close();
    
    return {
      notes: notesCount.count,
      meetings: meetingsCount.count,
      interviews: interviewsCount.count,
      recordings: recordingsCount.count,
      totalDurationMinutes: totalDuration.total ? Math.round(totalDuration.total / (1000 * 1000 * 1000 * 60)) : 0,
    };
  } catch (error) {
    console.error('Error fetching dashboard stats:', error);
    return {
      notes: 0,
      meetings: 0,
      interviews: 0,
      recordings: 0,
      totalDurationMinutes: 0,
    };
  }
}
