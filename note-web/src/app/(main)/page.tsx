import Link from 'next/link';
import { getDashboardStats, getRecentNotes, getRecentMeetings, getRecentInterviews, getRecentRecordings } from '@/lib/database';
import NoteCard from '@/components/NoteCard';
import MeetingCard from '@/components/MeetingCard';
import { Users, Mic, Briefcase, FileText } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

export default function Home() {
  const stats = getDashboardStats();
  const recentNotes = getRecentNotes(3);
  const recentMeetings = getRecentMeetings(3);
  const recentInterviews = getRecentInterviews(3);
  const recentRecordings = getRecentRecordings(3);

  const StatCard = ({ title, value, description, icon: Icon, href }: {
    title: string;
    value: number | string;
    description: string;
    icon: React.ComponentType<{ className?: string }>;
    href: string;
  }) => (
    <Link href={href} className="group">
      <Card className="hover:shadow-lg transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-muted-foreground">{title}</p>
              <p className="text-3xl font-bold">{value}</p>
              <p className="text-sm text-muted-foreground">{description}</p>
            </div>
            <div className="opacity-60 group-hover:opacity-100 transition-opacity">
              <Icon className="h-8 w-8 text-muted-foreground" />
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );

  const SectionHeader = ({ title, href, count }: { title: string; href: string; count: number }) => (
    <div className="flex items-center justify-between mb-4">
      <h2 className="text-xl font-semibold">{title}</h2>
      <Link href={href} className="text-primary hover:text-primary/80 text-sm font-medium transition-colors">
        View all ({count})
      </Link>
    </div>
  );

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back! Here&apos;s an overview of your recent activity.
          </p>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Notes"
            value={stats.notes}
            description="Total notes created"
            icon={FileText}
            href="/notes"
          />
          <StatCard
            title="Meetings"
            value={stats.meetings}
            description="Meetings recorded"
            icon={Users}
            href="/meetings"
          />
          <StatCard
            title="Interviews"
            value={stats.interviews}
            description="Interviews conducted"
            icon={Briefcase}
            href="/interviews"
          />
          <StatCard
            title="Audio Time"
            value={`${stats.totalDurationMinutes}m`}
            description="Total recorded"
            icon={Mic}
            href="/recordings"
          />
        </div>

        {/* Recent Activity Sections */}
        <div className="space-y-8">
          {/* Recent Notes */}
          {recentNotes.length > 0 && (
            <section>
              <SectionHeader title="Recent Notes" href="/notes" count={stats.notes} />
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {recentNotes.map((note) => (
                  <NoteCard key={note.id} note={note} />
                ))}
              </div>
            </section>
          )}

          {/* Recent Meetings */}
          {recentMeetings.length > 0 && (
            <section>
              <SectionHeader title="Recent Meetings" href="/meetings" count={stats.meetings} />
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {recentMeetings.map((meeting) => (
                  <MeetingCard key={meeting.id} meeting={meeting} />
                ))}
              </div>
            </section>
          )}

          {/* Recent Interviews */}
          {recentInterviews.length > 0 && (
            <section>
              <SectionHeader title="Recent Interviews" href="/interviews" count={stats.interviews} />
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {recentInterviews.map((interview) => (
                    <Card key={interview.id} className="hover:shadow-lg transition-shadow">
                      <CardHeader className="pb-3">
                        <div className="flex items-start justify-between">
                          <CardTitle className="text-lg truncate">
                            {interview.title}
                          </CardTitle>
                          <Badge variant="outline" className="text-xs ml-2">
                            <Briefcase className="h-3 w-3 mr-1" />
                            Interview
                          </Badge>
                        </div>
                      </CardHeader>
                      <CardContent className="pt-0">
                        {interview.company && (
                          <p className="text-sm mb-2">
                            {interview.company} - {interview.position}
                          </p>
                        )}
                        <p className="text-sm text-muted-foreground mb-3 line-clamp-2">
                          {interview.interviewee} interviewed by {interview.interviewer}
                        </p>
                        <div className="flex items-center justify-between text-xs text-muted-foreground">
                          <span>{new Date(interview.created_at).toLocaleDateString()}</span>
                          {interview.tags && (
                            <Badge variant="secondary" className="text-xs">
                              {interview.tags.split(',')[0]}
                            </Badge>
                          )}
                        </div>
                      </CardContent>
                    </Card>
                ))}
              </div>
            </section>
          )}

          {/* Recent Recordings */}
          {recentRecordings.length > 0 && (
            <section>
              <SectionHeader title="Recent Recordings" href="/recordings" count={stats.recordings} />
              <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {recentRecordings.map((recording) => {
                  const durationMinutes = Math.round(recording.duration / (1000 * 1000 * 1000 * 60));
                  return (
                    <Card key={recording.id} className="hover:shadow-lg transition-shadow">
                      <CardHeader className="pb-3">
                        <div className="flex items-start justify-between">
                          <CardTitle className="text-lg truncate">
                            {recording.filename.replace(/\.[^/.]+$/, '')}
                          </CardTitle>
                          <Badge variant="outline" className="text-xs ml-2">
                            <Mic className="h-3 w-3 mr-1" />
                            Recording
                          </Badge>
                        </div>
                      </CardHeader>
                      <CardContent className="pt-0">
                        <p className="text-sm mb-2">
                          Duration: {durationMinutes} minutes
                        </p>
                        <p className="text-sm text-muted-foreground mb-3">
                          Format: {recording.format} • {recording.sample_rate}Hz
                        </p>
                        <div className="flex items-center justify-between text-xs text-muted-foreground">
                          <span>{new Date(recording.created_at).toLocaleDateString()}</span>
                          <Badge variant="secondary" className="text-xs">
                            {Math.round(recording.file_size / 1024 / 1024)} MB
                          </Badge>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </section>
          )}

          {/* Empty State */}
          {stats.notes === 0 && stats.meetings === 0 && stats.interviews === 0 && stats.recordings === 0 && (
            <Card className="p-8 text-center">
              <CardContent>
                <div className="text-muted-foreground mb-4">
                  <FileText className="mx-auto h-16 w-16" />
                </div>
                <h3 className="text-lg font-medium mb-2">Welcome to Note AI</h3>
                <p className="text-muted-foreground mb-4">
                  Get started by creating your first note using the CLI:
                </p>
                <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
                  note import /path/to/your/audio/file.mp3
                </code>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
