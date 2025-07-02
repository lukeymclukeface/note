import { Note } from '@/lib/database';
import { Music, FileText } from 'lucide-react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

interface NoteCardProps {
  note: Note;
}

export default function NoteCard({ note }: NoteCardProps) {
  // Format the date
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // Parse tags
  const tags = note.tags ? note.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) : [];

  // Truncate content for preview
  const truncateContent = (content: string, maxLength: number = 300) => {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + '...';
  };

  return (
    <Card className="hover:shadow-md transition-shadow duration-200">
      <CardHeader className="pb-3">
        <div className="flex justify-between items-start">
          <h2 className="text-xl font-semibold line-clamp-2">
            {note.title}
          </h2>
          <div className="flex items-center space-x-2 ml-4">
            {note.recording_id && (
              <Badge variant="secondary" className="text-xs">
                <Music className="mr-1 h-3 w-3" /> Recording
              </Badge>
            )}
            <Badge variant="outline" className="text-xs">
              <FileText className="mr-1 h-3 w-3" /> Note
            </Badge>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        
        <div className="text-sm text-muted-foreground mb-3">
          {formatDate(note.created_at)}
        </div>
        
        {tags.length > 0 && (
          <div className="mb-3 flex flex-wrap gap-1">
            {tags.map((tag, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>
        )}
        
        {note.summary && (
          <div className="mb-4">
            <h4 className="text-sm font-medium text-muted-foreground mb-1">Summary:</h4>
            <div className="text-sm">
              {truncateContent(note.summary)}
            </div>
          </div>
        )}
        
        <div className="prose prose-sm max-w-none">
          <h4 className="text-sm font-medium text-muted-foreground mb-1">Content:</h4>
          <div className="whitespace-pre-wrap text-sm">
            {truncateContent(note.content)}
          </div>
        </div>
        
        <div className="mt-4 pt-4 border-t border-border">
          <div className="flex justify-between items-center text-xs text-muted-foreground">
            <span>ID: {note.id}</span>
            {note.updated_at !== note.created_at && (
              <span>Updated: {formatDate(note.updated_at)}</span>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
