import { getMeetingById } from '@/lib/database';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import { MarkdownRenderer } from '@/components/MarkdownRenderer';

interface MeetingPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function MeetingPage({ params }: MeetingPageProps) {
  const { id } = await params;
  const meetingId = parseInt(id);
  
  if (isNaN(meetingId)) {
    notFound();
  }

  const meeting = getMeetingById(meetingId);

  if (!meeting) {
    notFound();
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatMeetingDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-4xl mx-auto p-6">
        {/* Breadcrumb */}
        <nav className="mb-6">
          <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
            <Link href="/meetings" className="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              Meetings
            </Link>
            <span>/</span>
            <span className="text-gray-900 dark:text-white font-medium">Meeting Details</span>
          </div>
        </nav>

        {/* Header */}
        <div className="mb-8">
          <div className="flex items-start justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white leading-tight">
              {meeting.title}
            </h1>
            <div className="flex items-center space-x-2 ml-4">
              {meeting.recording_id && (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
                  🎵 Recording Available
                </span>
              )}
              <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200">
                📅 Meeting
              </span>
            </div>
          </div>

          {/* Meeting Metadata */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <div>
              <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Created:</span>
              <p className="text-sm text-gray-900 dark:text-gray-100">{formatDate(meeting.created_at)}</p>
            </div>
          
            {meeting.meeting_date && (
              <div>
                <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Meeting Date:</span>
                <p className="text-sm text-gray-900 dark:text-gray-100">{formatMeetingDate(meeting.meeting_date)}</p>
              </div>
            )}
          
            {meeting.attendees && (
              <div>
                <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Attendees:</span>
                <p className="text-sm text-gray-900 dark:text-gray-100">{meeting.attendees}</p>
              </div>
            )}
          
            {meeting.location && (
              <div>
                <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Location:</span>
                <p className="text-sm text-gray-900 dark:text-gray-100">{meeting.location}</p>
              </div>
            )}
          
            <div>
              <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Meeting ID:</span>
              <p className="text-sm text-gray-900 dark:text-gray-100">{meeting.id}</p>
            </div>
          
            {meeting.recording_id && (
              <div>
                <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Recording ID:</span>
                <p className="text-sm text-gray-900 dark:text-gray-100">{meeting.recording_id}</p>
              </div>
            )}
        </div>
      </div>

        {/* Main Content */}
        <div className="space-y-8">
          {/* Summary Section */}
          {meeting.summary && (
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">
                📋 Meeting Summary
              </h2>
              <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
                <MarkdownRenderer content={meeting.summary} />
              </div>
            </section>
          )}

          {/* Full Content Section */}
          {meeting.content && (
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">
                📄 Full Content
              </h2>
              <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6 shadow-sm">
                <MarkdownRenderer content={meeting.content} />
              </div>
            </section>
          )}

          {/* Tags Section */}
          {meeting.tags && (
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 dark:text-white mb-4">
                🏷️ Tags
              </h2>
              <div className="flex flex-wrap gap-2">
                {meeting.tags.split(',').map((tag, index) => (
                  <span
                    key={index}
                    className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200"
                  >
                    {tag.trim()}
                  </span>
                ))}
              </div>
            </section>
          )}
      </div>

        {/* Actions */}
        <div className="mt-12 pt-6 border-t border-gray-200 dark:border-gray-700">
          <div className="flex justify-between items-center">
            <Link
              href="/meetings"
              className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              ← Back to Meetings
            </Link>
            
            <div className="flex space-x-3">
              {meeting.recording_id && (
                <Link
                  href={`/recordings/${meeting.recording_id}`}
                  className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors"
                >
                  🎵 View Recording
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Generate metadata for the page
export async function generateMetadata({ params }: MeetingPageProps) {
  const { id } = await params;
  const meetingId = parseInt(id);
  
  if (isNaN(meetingId)) {
    return {
      title: 'Meeting Not Found',
    };
  }

  const meeting = getMeetingById(meetingId);

  if (!meeting) {
    return {
      title: 'Meeting Not Found',
    };
  }

  return {
    title: `${meeting.title} - Meeting Details`,
    description: meeting.summary ? meeting.summary.substring(0, 160) + '...' : 'Meeting details and summary',
  };
}
