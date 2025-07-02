import { getAllMeetings } from '@/lib/database';
import MeetingCard from '@/components/MeetingCard';
import { Users } from 'lucide-react';

export default function MeetingsPage() {
  const meetings = getAllMeetings();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-6xl mx-auto p-6">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Meetings</h1>
          <p className="text-gray-600 dark:text-gray-300 mt-2">
            {meetings.length} meeting{meetings.length !== 1 ? 's' : ''} found
          </p>
        </div>
      
        {meetings.length === 0 ? (
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-12 text-center">
            <div className="flex justify-center mb-4">
              <Users className="h-16 w-16 text-gray-400 dark:text-gray-500" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No meetings yet</h3>
            <p className="text-gray-600 dark:text-gray-300">
              Import audio files or text documents to create meeting records.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {meetings.map((meeting) => (
              <MeetingCard key={meeting.id} meeting={meeting} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
