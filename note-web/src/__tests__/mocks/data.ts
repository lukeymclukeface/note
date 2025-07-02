import { Note, Meeting, Interview, Recording } from '@/lib/database'
import { AppConfig } from '@/lib/config'

export const mockNote: Note = {
  id: 1,
  title: 'Test Note Title',
  content: 'This is a test note content with some detailed information about the topic.',
  summary: 'This is a test summary of the note.',
  tags: 'test,mock,unit-test',
  recording_id: 1,
  created_at: '2023-12-01T10:00:00.000Z',
  updated_at: '2023-12-01T10:00:00.000Z',
}

export const mockNoteWithoutRecording: Note = {
  id: 2,
  title: 'Note Without Recording',
  content: 'This note does not have an associated recording.',
  summary: 'Summary without recording.',
  tags: 'text,import',
  created_at: '2023-12-02T14:30:00.000Z',
  updated_at: '2023-12-02T14:30:00.000Z',
}

export const mockMeeting: Meeting = {
  id: 1,
  title: 'Weekly Team Standup',
  content: 'Discussion about project progress, blockers, and next steps.',
  summary: 'Team discussed current sprint progress. No major blockers identified.',
  attendees: 'John Doe, Jane Smith, Bob Johnson',
  location: 'Conference Room A',
  tags: 'meeting,standup,team',
  recording_id: 2,
  meeting_date: '2023-12-01',
  created_at: '2023-12-01T09:00:00.000Z',
  updated_at: '2023-12-01T09:00:00.000Z',
}

export const mockMeetingMinimal: Meeting = {
  id: 2,
  title: 'Quick Sync',
  content: 'Brief discussion about urgent matters.',
  summary: 'Resolved immediate issues.',
  attendees: '',
  location: '',
  tags: 'meeting,sync',
  created_at: '2023-12-02T16:00:00.000Z',
  updated_at: '2023-12-02T16:00:00.000Z',
}

export const mockInterview: Interview = {
  id: 1,
  title: 'Frontend Developer Interview',
  content: 'Technical interview focusing on React and TypeScript skills.',
  summary: 'Candidate showed strong technical skills and good problem-solving abilities.',
  interviewee: 'Alex Smith',
  interviewer: 'Sarah Wilson',
  company: 'Tech Corp',
  position: 'Senior Frontend Developer',
  tags: 'interview,frontend,technical',
  recording_id: 3,
  interview_date: '2023-12-01',
  created_at: '2023-12-01T14:00:00.000Z',
  updated_at: '2023-12-01T14:00:00.000Z',
}

export const mockRecording: Recording = {
  id: 1,
  filename: 'meeting-recording.mp3',
  file_path: '/path/to/recording.mp3',
  start_time: '2023-12-01T09:00:00.000Z',
  end_time: '2023-12-01T10:00:00.000Z',
  duration: 3600000, // 1 hour in milliseconds
  file_size: 5242880, // 5MB
  format: 'mp3',
  sample_rate: 44100,
  channels: 2,
  created_at: '2023-12-01T09:00:00.000Z',
}

export const mockConfig: AppConfig = {
  notes_dir: '/Users/test/.noteai/notes',
  editor: 'code',
  date_format: '2006-01-02',
  default_tags: ['imported', 'ai-processed'],
  openai_key: 'sk-test1234567890abcdef',
  database_path: '/Users/test/.noteai/notes.db',
  transcription_model: 'whisper-1',
  summary_model: 'gpt-4',
  google_project_id: 'test-project-id',
  google_location: 'us-central1',
  transcription_provider: 'openai',
  summary_provider: 'openai',
}

export const mockConfigMissingKey: AppConfig = {
  notes_dir: '/Users/test/.noteai/notes',
  editor: 'nano',
  date_format: '2006-01-02',
  default_tags: [],
  openai_key: '',
  database_path: '/Users/test/.noteai/notes.db',
  transcription_model: 'whisper-1',
  summary_model: 'gpt-3.5-turbo',
  google_project_id: '',
  google_location: 'us-central1',
  transcription_provider: 'openai',
  summary_provider: 'openai',
}

// Arrays for list components
export const mockNotes: Note[] = [mockNote, mockNoteWithoutRecording]
export const mockMeetings: Meeting[] = [mockMeeting, mockMeetingMinimal]
export const mockInterviews: Interview[] = [mockInterview]
export const mockRecordings: Recording[] = [mockRecording]
