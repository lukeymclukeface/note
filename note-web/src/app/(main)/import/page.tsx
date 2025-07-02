'use client';

import { useState, useEffect } from 'react';
import { loadRecordings } from '@/lib/actions/recordings';
import { formatTime, formatDuration } from '@/lib/dateUtils';
import type { Recording } from '@/lib/database';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import { Mic } from 'lucide-react';

export default function ImportPage() {
  const [recordings, setRecordings] = useState<Recording[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchRecordings = async () => {
    try {
      const data = await loadRecordings();
      setRecordings(data);
    } catch (error) {
      console.error('Error loading recordings:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchSingleRecording = async (recordingId: number) => {
    try {
      const response = await fetch(`/api/recordings/${recordingId}`);
      if (response.ok) {
        const result = await response.json();
        if (result.success && result.recording) {
          // Add the new recording to the top of the list
          setRecordings(prev => [result.recording, ...prev]);
        }
      }
    } catch (error) {
      console.error('Error fetching new recording:', error);
    }
  };

  useEffect(() => {
    fetchRecordings();
    
    // Listen for recording completion events
    const handleRecordingCompleted = (event: CustomEvent) => {
      const { recordingId } = event.detail;
      if (recordingId) {
        // Fetch the complete recording data and add it to the list
        fetchSingleRecording(recordingId);
      }
    };
    
    // Add event listener for real-time updates
    window.addEventListener('recordingCompleted', handleRecordingCompleted as EventListener);
    
    // Cleanup
    return () => {
      window.removeEventListener('recordingCompleted', handleRecordingCompleted as EventListener);
    };
  }, []);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

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

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-semibold mb-2">Audio Recordings</h2>
        <p className="text-muted-foreground">
          {recordings.length === 0 
            ? 'No recordings found. Use the record button in the navigation bar to start recording.' 
            : `${recordings.length} recording${recordings.length !== 1 ? 's' : ''} found`
          }
        </p>
      </div>

      {isLoading ? (
        <Card>
          <CardContent className="p-6">
            <div className="space-y-3">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-3/4" />
            </div>
          </CardContent>
        </Card>
      ) : recordings.length === 0 ? (
        <Card className="p-8 text-center">
          <CardContent>
            <div className="text-muted-foreground mb-4">
              <Mic className="mx-auto h-16 w-16" />
            </div>
            <h3 className="text-lg font-medium mb-2">No recordings yet</h3>
            <p className="text-muted-foreground mb-4">
              Get started by clicking the record button in the navigation bar.
            </p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Recording</TableHead>
                  <TableHead>Start Time</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead>Format</TableHead>
                  <TableHead>Size</TableHead>
                  <TableHead>Quality</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recordings.map((recording) => {
                  // Convert nanoseconds to seconds, then to minutes
                  const durationSeconds = recording.duration / (1000 * 1000 * 1000);
                  const durationMinutes = Math.round(durationSeconds / 60);
                  const startTime = new Date(recording.start_time);
                  const endTime = new Date(recording.end_time);
                  
                  return (
                    <TableRow key={recording.id}>
                      <TableCell>
                        <div className="flex items-center">
                          <div className="flex-shrink-0 h-8 w-8">
                            <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
                              <Mic className="h-4 w-4 text-primary" />
                            </div>
                          </div>
                          <div className="ml-4">
                            <div className="text-sm font-medium">
                              <a href={`/recordings/${recording.id}`} className="text-primary hover:text-primary/80 hover:underline">
                                {recording.filename.replace(/\.[^/.]+$/, '')}
                              </a>
                            </div>
                            <div className="text-sm text-muted-foreground">
                              ID: {recording.id}
                            </div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="text-sm">
                          {formatDate(recording.start_time)}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {formatTime(startTime)} - {formatTime(endTime)}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant="secondary">
                          {formatDuration(durationMinutes)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-sm">
                        {recording.format.toUpperCase()}
                      </TableCell>
                      <TableCell className="text-sm">
                        {formatFileSize(recording.file_size)}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        <div>{recording.sample_rate / 1000}kHz</div>
                        <div>{recording.channels} ch</div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
