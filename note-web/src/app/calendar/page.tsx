import { getAllRecordings } from '@/lib/database';
import WeekCalendar from '@/components/WeekCalendar';

export default function CalendarPage() {
  // Get all recordings and convert to calendar events
  const recordings = getAllRecordings();
  
  const events = recordings.map(recording => {
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

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Recording Calendar</h1>
          <p className="text-gray-600">
            View your recorded sessions in a weekly calendar format
          </p>
        </header>

        {recordings.length === 0 ? (
          <div className="bg-white border border-gray-200 rounded-lg p-8 text-center">
            <div className="text-gray-400 mb-4">
              <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No recordings yet</h3>
            <p className="text-gray-500 mb-4">
              Start recording audio sessions using the CLI to see them on the calendar:
            </p>
            <code className="bg-gray-100 text-gray-800 px-3 py-1 rounded font-mono text-sm">
              note record
            </code>
          </div>
        ) : (
          <WeekCalendar events={events} />
        )}
      </div>
    </div>
  );
}
