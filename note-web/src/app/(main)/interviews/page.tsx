import { getAllInterviews } from '@/lib/database';
import InterviewCard from '@/components/InterviewCard';
import { Card, CardContent } from '@/components/ui/card';
import { Users } from 'lucide-react';

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
                <Users className="mx-auto h-16 w-16" />
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
