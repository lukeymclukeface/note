'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { formatDuration } from '@/lib/dateUtils';
import { Music, Loader2, AlertCircle } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import AudioPlayer from '@/components/AudioPlayer';
import CopyLinkButton from '@/components/CopyLinkButton';
import type { Recording } from '@/lib/database';

export default function RecordingDetailsPage() {
  const params = useParams();
  const [recording, setRecording] = useState<Recording | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRecording = async () => {
      if (!params?.id) {
        setError('Invalid recording ID');
        setIsLoading(false);
        return;
      }

      try {
        const response = await fetch(`/api/recordings/${params.id}`);
        if (!response.ok) {
          throw new Error('Recording not found');
        }
        const result = await response.json();
        if (result.success && result.recording) {
          setRecording(result.recording);
        } else {
          throw new Error('Invalid response format');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load recording');
      } finally {
        setIsLoading(false);
      }
    };

    fetchRecording();
  }, [params?.id]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <div className="max-w-7xl mx-auto px-4 py-8">
          {/* Breadcrumb Skeleton */}
          <nav className="mb-6">
            <div className="flex items-center space-x-2 text-sm">
              <Skeleton className="h-4 w-12" />
              <span>/</span>
              <Skeleton className="h-4 w-20" />
              <span>/</span>
              <Skeleton className="h-4 w-32" />
            </div>
          </nav>

          {/* Header Skeleton */}
          <header className="mb-8">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <Skeleton className="h-10 w-80 mb-2" />
              </div>
              <div className="flex items-center space-x-2 ml-4">
                <Skeleton className="h-6 w-20 rounded-full" />
              </div>
            </div>
            <Skeleton className="h-4 w-96" />
          </header>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Audio Player Skeleton */}
            <div className="lg:col-span-2">
              <Card>
                <CardContent className="p-6">
                  <Skeleton className="h-6 w-32 mb-4" />
                  <div className="space-y-6">
                    {/* Play button skeleton */}
                    <div className="flex items-center justify-center">
                      <Skeleton className="w-16 h-16 rounded-full" />
                    </div>
                    {/* Progress bar skeleton */}
                    <div className="space-y-2">
                      <Skeleton className="h-2 w-full rounded-lg" />
                      <div className="flex justify-between">
                        <Skeleton className="h-4 w-12" />
                        <Skeleton className="h-4 w-16" />
                      </div>
                    </div>
                    {/* Volume control skeleton */}
                    <div className="flex items-center space-x-3">
                      <Skeleton className="w-4 h-4" />
                      <Skeleton className="h-2 flex-1 rounded-lg" />
                      <Skeleton className="h-4 w-8" />
                    </div>
                    {/* File info skeleton */}
                    <div className="text-center">
                      <Skeleton className="h-4 w-48 mx-auto" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Recording Details Skeleton */}
            <div className="space-y-6">
              {/* Recording Information */}
              <Card>
                <CardContent className="p-6">
                  <Skeleton className="h-6 w-40 mb-4" />
                  <div className="grid grid-cols-1 gap-4">
                    {Array.from({ length: 6 }).map((_, i) => (
                      <div key={i}>
                        <Skeleton className="h-4 w-24 mb-1" />
                        <Skeleton className="h-4 w-32" />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              {/* Recording Timeline */}
              <Card>
                <CardContent className="p-6">
                  <Skeleton className="h-6 w-36 mb-4" />
                  <div className="space-y-3">
                    {Array.from({ length: 3 }).map((_, i) => (
                      <div key={i}>
                        <Skeleton className="h-4 w-20 mb-1" />
                        <Skeleton className="h-4 w-40" />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              {/* Actions */}
              <Card>
                <CardContent className="p-6">
                  <Skeleton className="h-6 w-16 mb-4" />
                  <div className="space-y-3">
                    <Skeleton className="h-10 w-full rounded-md" />
                    <Skeleton className="h-10 w-full rounded-md" />
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>

          {/* Bottom Actions Skeleton */}
          <div className="mt-12 pt-6 border-t border-border">
            <div className="flex justify-between items-center">
              <Skeleton className="h-10 w-40 rounded-md" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="p-8 text-center">
          <CardContent>
            <div className="text-destructive mb-4">
              <AlertCircle className="mx-auto h-16 w-16" />
            </div>
            <h3 className="text-lg font-medium mb-2">Error Loading Recording</h3>
            <p className="text-muted-foreground mb-4">{error}</p>
            <Button asChild>
              <Link href="/import">
                Back to Recordings
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!recording) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card className="p-8 text-center">
          <CardContent>
            <h3 className="text-lg font-medium mb-2">Recording Not Found</h3>
            <Button asChild>
              <Link href="/import">
                Back to Recordings
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

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
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };

  // Convert nanoseconds to seconds, then to minutes
  const durationSeconds = recording.duration / (1000 * 1000 * 1000);
  const durationMinutes = Math.round(durationSeconds / 60);
  
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Breadcrumb */}
        <nav className="mb-6">
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <Link href="/import" className="hover:text-primary transition-colors">
              Import
            </Link>
            <span>/</span>
            <Link href="/import" className="hover:text-primary transition-colors">
              Recordings
            </Link>
            <span>/</span>
            <span className="text-foreground font-medium">Recording Details</span>
          </div>
        </nav>

        {/* Header */}
        <header className="mb-8">
          <div className="flex items-start justify-between mb-4">
            <h1 className="text-4xl font-bold mb-2">
              {recording.filename.replace(/\.[^/.]+$/, '')}
            </h1>
            <div className="flex items-center space-x-2 ml-4">
              <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary/10 text-primary">
                <Music className="mr-1 h-4 w-4" /> Recording
              </span>
            </div>
          </div>
          <p className="text-muted-foreground">
            Audio recording details and playback controls
          </p>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Audio Player */}
          <div className="lg:col-span-2">
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold mb-4">Audio Player</h3>
                <AudioPlayer 
                  src={`/api/recordings/${recording.id}/audio`}
                  filename={recording.filename}
                  databaseDuration={durationSeconds}
                />
              </CardContent>
            </Card>
          </div>

          {/* Recording Details */}
          <div className="space-y-6">
            {/* Recording Information */}
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold mb-4">Recording Information</h3>
                
                <div className="grid grid-cols-1 gap-4">
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">File Name:</span>
                    <p className="text-sm">{recording.filename}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Recording ID:</span>
                    <p className="text-sm">#{recording.id}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Duration:</span>
                    <p className="text-sm">{formatDuration(durationMinutes)}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">File Size:</span>
                    <p className="text-sm">{formatFileSize(recording.file_size)}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Format:</span>
                    <p className="text-sm">{recording.format.toUpperCase()}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Audio Quality:</span>
                    <div className="text-sm">
                      <p>{recording.sample_rate / 1000}kHz</p>
                      <p>{recording.channels} {recording.channels === 1 ? 'channel' : 'channels'}</p>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Recording Timeline */}
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold mb-4">Recording Timeline</h3>
                
                <div className="space-y-3">
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Started:</span>
                    <p className="text-sm">{formatDate(recording.start_time)}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Ended:</span>
                    <p className="text-sm">{formatDate(recording.end_time)}</p>
                  </div>
                  
                  <div>
                    <span className="text-sm font-medium text-muted-foreground">Created:</span>
                    <p className="text-sm">{formatDate(recording.created_at)}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Actions */}
            <Card>
              <CardContent className="p-6">
                <h3 className="text-lg font-semibold mb-4">Actions</h3>
                
                <div className="space-y-3">
                  <Button asChild className="w-full">
                    <a
                      href={`/api/recordings/${recording.id}/audio`}
                      download={recording.filename}
                    >
                      Download Recording
                    </a>
                  </Button>
                  
                  <CopyLinkButton 
                    variant="outline" 
                    className="w-full"
                  />
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Actions */}
        <div className="mt-12 pt-6 border-t border-border">
          <div className="flex justify-between items-center">
            <Button variant="outline" asChild>
              <Link href="/import">
                ← Back to Recordings
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
