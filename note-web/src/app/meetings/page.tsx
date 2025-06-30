import { getAllMeetings } from '@/lib/database';
import MeetingCard from '@/components/MeetingCard';

export default function MeetingsPage() {
  const meetings = getAllMeetings();

  return (
    <div className="max-w-6xl mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Meetings</h1>
        <p className="text-gray-600 mt-2">
          {meetings.length} meeting{meetings.length !== 1 ? 's' : ''} found
        </p>
      </div>
      
      {meetings.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🤝</div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No meetings yet</h3>
          <p className="text-gray-600">
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
  );
}
