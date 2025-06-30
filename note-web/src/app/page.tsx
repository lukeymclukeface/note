import Link from 'next/link';
import { getDashboardStats, getRecentNotes, getRecentMeetings, getRecentInterviews, getRecentRecordings } from '@/lib/database';
import NoteCard from '@/components/NoteCard';
import MeetingCard from '@/components/MeetingCard';

export default function Home() {
  const stats = getDashboardStats();
  const recentNotes = getRecentNotes(3);
  const recentMeetings = getRecentMeetings(3);
  const recentInterviews = getRecentInterviews(3);
  const recentRecordings = getRecentRecordings(3);

  const StatCard = ({ title, value, description, icon, href }: {
    title: string;
    value: number | string;
    description: string;
    icon: string;
    href: string;
  }) => (
    <Link href={href} className="group">
      <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:shadow-lg transition-shadow">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-600 dark:text-gray-400">{title}</p>
            <p className="text-3xl font-bold text-gray-900 dark:text-white">{value}</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">{description}</p>
          </div>
          <div className="text-3xl opacity-60 group-hover:opacity-100 transition-opacity">
            {icon}
          </div>
        </div>
      </div>
    </Link>
  );

  const SectionHeader = ({ title, href, count }: { title: string; href: string; count: number }) => (
    <div className="flex items-center justify-between mb-4">
      <h2 className="text-xl font-semibold text-gray-900 dark:text-white">{title}</h2>
      <Link href={href} className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 text-sm font-medium">
        View all ({count})
      </Link>
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">Dashboard</h1>
          <p className="text-gray-600 dark:text-gray-300">
            Welcome back! Here&apos;s an overview of your recent activity.
          </p>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Notes"
            value={stats.notes}
            description="Total notes created"
            icon="📝"
            href="/notes"
          />
          <StatCard
            title="Meetings"
            value={stats.meetings}
            description="Meetings recorded"
            icon="🤝"
            href="/meetings"
          />
          <StatCard
            title="Interviews"
            value={stats.interviews}
            description="Interviews conducted"
            icon="💼"
            href="/interviews"
          />
          <StatCard
            title="Audio Time"
            value={`${stats.totalDurationMinutes}m`}
            description="Total recorded"
            icon="🎤"
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
                  <div key={interview.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:shadow-lg transition-shadow">
                    <div className="flex items-start justify-between mb-3">
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                        {interview.title}
                      </h3>
                      <span className="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap ml-2">
                        💼
                      </span>
                    </div>
                    {interview.company && (
                      <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                        {interview.company} - {interview.position}
                      </p>
                    )}
                    <p className="text-sm text-gray-500 dark:text-gray-400 mb-3 line-clamp-2">
                      {interview.interviewee} interviewed by {interview.interviewer}
                    </p>
                    <div className="flex items-center justify-between text-xs text-gray-400 dark:text-gray-500">
                      <span>{new Date(interview.created_at).toLocaleDateString()}</span>
                      {interview.tags && (
                        <span className="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                          {interview.tags.split(',')[0]}
                        </span>
                      )}
                    </div>
                  </div>
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
                    <div key={recording.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:shadow-lg transition-shadow">
                      <div className="flex items-start justify-between mb-3">
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                          {recording.filename.replace(/\.[^/.]+$/, '')}
                        </h3>
                        <span className="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap ml-2">
                          🎤
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">
                        Duration: {durationMinutes} minutes
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mb-3">
                        Format: {recording.format} • {recording.sample_rate}Hz
                      </p>
                      <div className="flex items-center justify-between text-xs text-gray-400 dark:text-gray-500">
                        <span>{new Date(recording.created_at).toLocaleDateString()}</span>
                        <span>{Math.round(recording.file_size / 1024 / 1024)} MB</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </section>
          )}

          {/* Empty State */}
          {stats.notes === 0 && stats.meetings === 0 && stats.interviews === 0 && stats.recordings === 0 && (
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-8 text-center">
              <div className="text-gray-400 dark:text-gray-500 mb-4">
                <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Welcome to Note AI</h3>
              <p className="text-gray-500 dark:text-gray-400 mb-4">
                Get started by creating your first note using the CLI:
              </p>
              <code className="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-3 py-1 rounded font-mono text-sm">
                note import /path/to/your/audio/file.mp3
              </code>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
