import { Note } from '@/lib/database';

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
    <div className="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 p-6">
      <div className="flex justify-between items-start mb-3">
        <h2 className="text-xl font-semibold text-gray-900 line-clamp-2">
          {note.title}
        </h2>
        <span className="text-sm text-gray-500 whitespace-nowrap ml-4">
          {formatDate(note.created_at)}
        </span>
      </div>
      
      {tags.length > 0 && (
        <div className="mb-3">
          {tags.map((tag, index) => (
            <span
              key={index}
              className="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full mr-2"
            >
              {tag}
            </span>
          ))}
        </div>
      )}
      
      <div className="text-gray-700 prose prose-sm max-w-none">
        <div className="whitespace-pre-wrap">
          {truncateContent(note.content)}
        </div>
      </div>
      
      <div className="mt-4 pt-4 border-t border-gray-100">
        <div className="flex justify-between items-center text-xs text-gray-500">
          <span>ID: {note.id}</span>
          {note.updated_at !== note.created_at && (
            <span>Updated: {formatDate(note.updated_at)}</span>
          )}
        </div>
      </div>
    </div>
  );
}
