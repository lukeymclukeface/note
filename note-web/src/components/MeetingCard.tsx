import { Meeting } from '@/lib/database';
import Link from 'next/link';

interface MeetingCardProps {
  meeting: Meeting;
}

export default function MeetingCard({ meeting }: MeetingCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatMeetingDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg transition-shadow duration-200">
      <div className="flex items-start justify-between mb-3">
        <h3 className="text-lg font-semibold line-clamp-2">
          <Link 
            href={`/meetings/${meeting.id}`}
            className="text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
          >
            {meeting.title}
          </Link>
        </h3>
        <div className="flex items-center space-x-2 ml-4">
          {meeting.recording_id && (
            <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              🎵 Recording
            </span>
          )}
          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
            📅 Meeting
          </span>
        </div>
      </div>

      {meeting.meeting_date && (
        <div className="mb-3">
          <span className="text-sm font-medium text-gray-600">Meeting Date: </span>
          <span className="text-sm text-gray-800">{formatMeetingDate(meeting.meeting_date)}</span>
        </div>
      )}

      {meeting.attendees && (
        <div className="mb-3">
          <span className="text-sm font-medium text-gray-600">Attendees: </span>
          <span className="text-sm text-gray-800">{meeting.attendees}</span>
        </div>
      )}

      {meeting.location && (
        <div className="mb-3">
          <span className="text-sm font-medium text-gray-600">Location: </span>
          <span className="text-sm text-gray-800">{meeting.location}</span>
        </div>
      )}

      {meeting.summary && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-600 mb-1">Summary:</h4>
          <p className="text-sm text-gray-700 line-clamp-3">{meeting.summary}</p>
        </div>
      )}

      {meeting.content && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-600 mb-1">Content:</h4>
          <p className="text-sm text-gray-700 line-clamp-3">{meeting.content}</p>
        </div>
      )}

      {meeting.tags && (
        <div className="mb-4">
          <div className="flex flex-wrap gap-1">
            {meeting.tags.split(',').map((tag, index) => (
              <span
                key={index}
                className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
              >
                {tag.trim()}
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="flex justify-between items-center text-xs text-gray-500 border-t border-gray-100 pt-3">
        <span>Created: {formatDate(meeting.created_at)}</span>
        <span>ID: {meeting.id}</span>
      </div>
    </div>
  );
}
