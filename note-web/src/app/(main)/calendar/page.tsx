import { getAllRecordings } from '@/lib/database';
import WeekCalendar from '@/components/WeekCalendar';
import { Card, CardContent } from '@/components/ui/card';
import { Calendar } from 'lucide-react';

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
    <div className="min-h-screen bg-background">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Recording Calendar</h1>
          <p className="text-muted-foreground">
            View your recorded sessions in a weekly calendar format
          </p>
        </header>

        {recordings.length === 0 ? (
          <Card className="p-8 text-center">
            <CardContent>
              <div className="text-muted-foreground mb-4">
                <Calendar className="mx-auto h-16 w-16" />
              </div>
              <h3 className="text-lg font-medium mb-2">No recordings yet</h3>
              <p className="text-muted-foreground mb-4">
                Start recording audio sessions using the CLI to see them on the calendar:
              </p>
              <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
                note record
              </code>
            </CardContent>
          </Card>
        ) : (
          <WeekCalendar events={events} />
        )}
      </div>
    </div>
  );
}
