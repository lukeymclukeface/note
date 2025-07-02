import { getAllMeetings } from '@/lib/database';
import MeetingCard from '@/components/MeetingCard';
import { Users } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

export default function MeetingsPage() {
  const meetings = getAllMeetings();

  return (
    <div>
        <div className="mb-6">
          <h2 className="text-2xl font-semibold mb-2">Meetings</h2>
          <p className="text-muted-foreground">
            {meetings.length === 0 
              ? 'No meetings found. Import audio files or text documents to create meeting records.' 
              : `${meetings.length} meeting${meetings.length !== 1 ? 's' : ''} found`
            }
          </p>
        </div>
      
        {meetings.length === 0 ? (
          <Card className="p-8 text-center">
            <CardContent>
              <div className="text-muted-foreground mb-4">
                <Users className="mx-auto h-16 w-16" />
              </div>
              <h3 className="text-lg font-medium mb-2">No meetings yet</h3>
              <p className="text-muted-foreground mb-4">
                Import audio files or text documents to create meeting records.
              </p>
              <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
                note import /path/to/your/meeting.mp3
              </code>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {meetings.map((meeting) => (
              <MeetingCard key={meeting.id} meeting={meeting} />
            ))}
          </div>
        )}
    </div>
  );
}
