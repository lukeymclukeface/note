import { getAllNotes } from '@/lib/database';
import NoteCard from '@/components/NoteCard';
import { Card, CardContent } from '@/components/ui/card';
import { FileText } from 'lucide-react';

export default function NotesPage() {
  const notes = getAllNotes();

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">My Notes</h1>
          <p className="text-muted-foreground">
            {notes.length === 0 
              ? 'No notes found. Create your first note using the CLI.' 
              : `${notes.length} note${notes.length !== 1 ? 's' : ''} found`
            }
          </p>
        </header>

        {notes.length === 0 ? (
          <Card className="p-8 text-center">
            <CardContent>
              <div className="text-muted-foreground mb-4">
                <FileText className="mx-auto h-16 w-16" strokeWidth={0.6} />
              </div>
              <h3 className="text-lg font-medium mb-2">No notes yet</h3>
              <p className="text-muted-foreground mb-4">
                Get started by creating your first note using the CLI:
              </p>
              <code className="bg-secondary text-secondary-foreground px-3 py-1 rounded font-mono text-sm">
                note import /path/to/your/audio/file.mp3
              </code>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {notes.map((note) => (
              <NoteCard key={note.id} note={note} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
