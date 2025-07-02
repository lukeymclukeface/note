import { getAllInterviews } from '@/lib/database';
import InterviewCard from '@/components/InterviewCard';
import { Card, CardContent } from '@/components/ui/card';

export default function InterviewsPage() {
  const interviews = getAllInterviews();

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Interviews</h1>
          <p className="text-muted-foreground">
            {interviews.length === 0 
              ? 'No interviews found. Import audio files or text documents to create interview records.' 
              : `${interviews.length} interview${interviews.length !== 1 ? 's' : ''} found`
            }
          </p>
        </header>
        
        {interviews.length === 0 ? (
          <Card className="p-8 text-center">
            <CardContent>
              <div className="text-muted-foreground mb-4">
                <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2-2v2m8 0V4a2 2 0 00-2-2H6a2 2 0 00-2 2v2m0 0v6.586a1 1 0 00.293.707l2.828 2.828a1 1 0 00.707.293h8.344a1 1 0 00.707-.293l2.828-2.828a1 1 0 00.293-.707V6z" />
                </svg>
              </div>
              <h3 className="text-lg font-medium mb-2">No interviews yet</h3>
              <p className="text-muted-foreground mb-4">
                Get started by importing your first interview:
              </p>
              <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
                note import /path/to/interview.mp3
              </code>
            </CardContent>
          </Card>
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
