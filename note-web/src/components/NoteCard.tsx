import { Note } from '@/lib/database';
import { Music, FileText } from 'lucide-react';

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
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 p-6">
      <div className="flex justify-between items-start mb-3">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white line-clamp-2">
          {note.title}
        </h2>
        <div className="flex items-center space-x-2 ml-4">
          {note.recording_id && (
            <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
              <Music className="mr-1 h-3 w-3" /> Recording
            </span>
          )}
          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
            <FileText className="mr-1 h-3 w-3" /> Note
          </span>
        </div>
      </div>
      
      <div className="text-sm text-gray-500 dark:text-gray-400 mb-3">
        {formatDate(note.created_at)}
      </div>
      
      {tags.length > 0 && (
        <div className="mb-3">
          {tags.map((tag, index) => (
            <span
              key={index}
              className="inline-block bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs px-2 py-1 rounded-full mr-2"
            >
              {tag}
            </span>
          ))}
        </div>
      )}
      
      {note.summary && (
        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Summary:</h4>
          <div className="text-gray-700 dark:text-gray-300 text-sm">
            {truncateContent(note.summary)}
          </div>
        </div>
      )}
      
      <div className="text-gray-700 dark:text-gray-300 prose prose-sm max-w-none">
        <h4 className="text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">Content:</h4>
        <div className="whitespace-pre-wrap text-sm">
          {truncateContent(note.content)}
        </div>
      </div>
      
      <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
        <div className="flex justify-between items-center text-xs text-gray-500 dark:text-gray-400">
          <span>ID: {note.id}</span>
          {note.updated_at !== note.created_at && (
            <span>Updated: {formatDate(note.updated_at)}</span>
          )}
        </div>
      </div>
    </div>
  );
}
