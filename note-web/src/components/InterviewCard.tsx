import { Interview } from '@/lib/database';

interface InterviewCardProps {
  interview: Interview;
}

export default function InterviewCard({ interview }: InterviewCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatInterviewDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <div className="bg-white rounded-lg shadow-md border border-gray-200 p-6 hover:shadow-lg transition-shadow duration-200">
      <div className="flex items-start justify-between mb-3">
        <h3 className="text-lg font-semibold text-gray-900 line-clamp-2">
          {interview.title}
        </h3>
        <div className="flex items-center space-x-2 ml-4">
          {interview.recording_id && (
            <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              🎵 Recording
            </span>
          )}
          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
            💼 Interview
          </span>
        </div>
      </div>

      {interview.interview_date && (
        <div className="mb-3">
          <span className="text-sm font-medium text-gray-600">Interview Date: </span>
          <span className="text-sm text-gray-800">{formatInterviewDate(interview.interview_date)}</span>
        </div>
      )}

      <div className="grid grid-cols-2 gap-3 mb-3">
        {interview.interviewee && (
          <div>
            <span className="text-sm font-medium text-gray-600">Interviewee: </span>
            <span className="text-sm text-gray-800">{interview.interviewee}</span>
          </div>
        )}
        {interview.interviewer && (
          <div>
            <span className="text-sm font-medium text-gray-600">Interviewer: </span>
            <span className="text-sm text-gray-800">{interview.interviewer}</span>
          </div>
        )}
      </div>

      <div className="grid grid-cols-2 gap-3 mb-3">
        {interview.company && (
          <div>
            <span className="text-sm font-medium text-gray-600">Company: </span>
            <span className="text-sm text-gray-800">{interview.company}</span>
          </div>
        )}
        {interview.position && (
          <div>
            <span className="text-sm font-medium text-gray-600">Position: </span>
            <span className="text-sm text-gray-800">{interview.position}</span>
          </div>
        )}
      </div>

      {interview.summary && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-600 mb-1">Summary:</h4>
          <p className="text-sm text-gray-700 line-clamp-3">{interview.summary}</p>
        </div>
      )}

      {interview.content && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-600 mb-1">Content:</h4>
          <p className="text-sm text-gray-700 line-clamp-3">{interview.content}</p>
        </div>
      )}

      {interview.tags && (
        <div className="mb-4">
          <div className="flex flex-wrap gap-1">
            {interview.tags.split(',').map((tag, index) => (
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
        <span>Created: {formatDate(interview.created_at)}</span>
        <span>ID: {interview.id}</span>
      </div>
    </div>
  );
}
