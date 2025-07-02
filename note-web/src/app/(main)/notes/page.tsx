import { getAllNotes, getAllMeetings, getAllInterviews } from '@/lib/database';
import NoteCard from '@/components/NoteCard';
import MeetingCard from '@/components/MeetingCard';
import InterviewCard from '@/components/InterviewCard';
import { Card, CardContent } from '@/components/ui/card';
import { FileText } from 'lucide-react';
import Link from 'next/link';

export default function NotesPage() {
  const notes = getAllNotes();
  const meetings = getAllMeetings();
  const interviews = getAllInterviews();
  const totalItems = notes.length + meetings.length + interviews.length;

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-semibold mb-2">All Notes</h2>
        <p className="text-muted-foreground">
          {totalItems === 0 
            ? 'No notes found. Import audio files or text documents to create note records.' 
            : `${totalItems} item${totalItems !== 1 ? 's' : ''} found (${notes.length} note${notes.length !== 1 ? 's' : ''}, ${meetings.length} meeting${meetings.length !== 1 ? 's' : ''}, ${interviews.length} interview${interviews.length !== 1 ? 's' : ''})`
          }
        </p>
      </div>

      {totalItems === 0 ? (
        <Card className="p-8 text-center">
          <CardContent>
            <div className="text-muted-foreground mb-4">
              <FileText className="mx-auto h-16 w-16" strokeWidth={0.6} />
            </div>
            <h3 className="text-lg font-medium mb-2">No notes yet</h3>
            <p className="text-muted-foreground mb-4">
              Import audio files or text documents to create note records.
            </p>
            <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
              note import /path/to/your/file.mp3
            </code>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-8">
          {/* Regular Notes Section */}
          {notes.length > 0 && (
            <section>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-xl font-semibold">Notes</h3>
                <p className="text-sm text-muted-foreground">
                  {notes.length} note{notes.length !== 1 ? 's' : ''}
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {notes.slice(0, 6).map((note) => (
                  <NoteCard key={note.id} note={note} />
                ))}
              </div>
            </section>
          )}

          {/* Meetings Section */}
          {meetings.length > 0 && (
            <section>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-xl font-semibold">Recent Meetings</h3>
                <p className="text-sm text-muted-foreground">
                  {meetings.length} meeting{meetings.length !== 1 ? 's' : ''}
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {meetings.slice(0, 6).map((meeting) => (
                  <MeetingCard key={meeting.id} meeting={meeting} />
                ))}
              </div>
              {meetings.length > 6 && (
                <div className="mt-4 text-center">
                  <Link 
                    href="/notes/meetings" 
                    className="text-primary hover:text-primary/80 text-sm font-medium"
                  >
                    View all {meetings.length} meetings →
                  </Link>
                </div>
              )}
            </section>
          )}

          {/* Interviews Section */}
          {interviews.length > 0 && (
            <section>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-xl font-semibold">Recent Interviews</h3>
                <p className="text-sm text-muted-foreground">
                  {interviews.length} interview{interviews.length !== 1 ? 's' : ''}
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {interviews.slice(0, 6).map((interview) => (
                  <InterviewCard key={interview.id} interview={interview} />
                ))}
              </div>
              {interviews.length > 6 && (
                <div className="mt-4 text-center">
                  <Link 
                    href="/notes/interviews" 
                    className="text-primary hover:text-primary/80 text-sm font-medium"
                  >
                    View all {interviews.length} interviews →
                  </Link>
                </div>
              )}
            </section>
          )}
        </div>
      )}
    </div>
  );
}
