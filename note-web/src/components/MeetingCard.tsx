import { Meeting } from '@/lib/database';
import Link from 'next/link';
import { Music, Calendar } from 'lucide-react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

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
    <Card className="hover:shadow-lg transition-shadow duration-200">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <h3 className="text-lg font-semibold line-clamp-2">
            <Link 
              href={`/meetings/${meeting.id}`}
              className="hover:text-primary transition-colors duration-200"
            >
              {meeting.title}
            </Link>
          </h3>
          <div className="flex items-center space-x-2 ml-4">
            {meeting.recording_id && (
              <Badge variant="secondary" className="text-xs">
                <Music className="mr-1 h-3 w-3" /> Recording
              </Badge>
            )}
            <Badge variant="outline" className="text-xs">
              <Calendar className="mr-1 h-3 w-3" /> Meeting
            </Badge>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">

        {meeting.meeting_date && (
          <div className="mb-3">
            <span className="text-sm font-medium text-muted-foreground">Meeting Date: </span>
            <span className="text-sm">{formatMeetingDate(meeting.meeting_date)}</span>
          </div>
        )}

        {meeting.attendees && (
          <div className="mb-3">
            <span className="text-sm font-medium text-muted-foreground">Attendees: </span>
            <span className="text-sm">{meeting.attendees}</span>
          </div>
        )}

        {meeting.location && (
          <div className="mb-3">
            <span className="text-sm font-medium text-muted-foreground">Location: </span>
            <span className="text-sm">{meeting.location}</span>
          </div>
        )}

        {meeting.summary && (
          <div className="mb-4">
            <h4 className="text-sm font-medium text-muted-foreground mb-1">Summary:</h4>
            <p className="text-sm line-clamp-3">{meeting.summary}</p>
          </div>
        )}

        {meeting.content && (
          <div className="mb-4">
            <h4 className="text-sm font-medium text-muted-foreground mb-1">Content:</h4>
            <p className="text-sm line-clamp-3">{meeting.content}</p>
          </div>
        )}

        {meeting.tags && (
          <div className="mb-4">
            <div className="flex flex-wrap gap-1">
              {meeting.tags.split(',').map((tag, index) => (
                <Badge key={index} variant="secondary" className="text-xs">
                  {tag.trim()}
                </Badge>
              ))}
            </div>
          </div>
        )}

        <div className="flex justify-between items-center text-xs text-muted-foreground border-t border-border pt-3">
          <span>Created: {formatDate(meeting.created_at)}</span>
          <span>ID: {meeting.id}</span>
        </div>
      </CardContent>
    </Card>
  );
}
