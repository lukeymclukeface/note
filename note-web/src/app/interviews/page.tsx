import { getAllInterviews } from '@/lib/database';
import InterviewCard from '@/components/InterviewCard';

export default function InterviewsPage() {
  const interviews = getAllInterviews();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">Interviews</h1>
          <p className="text-gray-600 dark:text-gray-300">
            {interviews.length === 0 
              ? 'No interviews found. Import audio files or text documents to create interview records.' 
              : `${interviews.length} interview${interviews.length !== 1 ? 's' : ''} found`
            }
          </p>
        </div>
        
        {interviews.length === 0 ? (
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-8 text-center">
            <div className="text-gray-400 dark:text-gray-500 mb-4">
              <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2-2v2m8 0V4a2 2 0 00-2-2H6a2 2 0 00-2 2v2m0 0v6.586a1 1 0 00.293.707l2.828 2.828a1 1 0 00.707.293h8.344a1 1 0 00.707-.293l2.828-2.828a1 1 0 00.293-.707V6z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No interviews yet</h3>
            <p className="text-gray-500 dark:text-gray-400 mb-4">
              Get started by importing your first interview:
            </p>
            <code className="bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-3 py-1 rounded font-mono text-sm">
              note import /path/to/interview.mp3
            </code>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {interviews.map((interview) => (
              <InterviewCard key={interview.id} interview={interview} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
