'use client';

import { useState, useEffect, useRef } from 'react';
import { useParams, useRouter } from 'next/navigation';
import type { Recording } from '@/lib/database';
import { formatDuration } from '@/lib/dateUtils';

interface AudioPlayerProps {
  src: string;
  filename: string;
  databaseDuration?: number; // Duration in seconds from database
}

function AudioPlayer({ src, filename, databaseDuration }: AudioPlayerProps) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isLoaded, setIsLoaded] = useState(false);
  const [audioError, setAudioError] = useState<string | null>(null);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const updateTime = () => {
      if (!isNaN(audio.currentTime)) {
        setCurrentTime(audio.currentTime);
      }
    };
    
    const updateDuration = () => {
      console.log('updateDuration called, audio.duration:', audio.duration, 'isFinite:', isFinite(audio.duration), 'databaseDuration:', databaseDuration);
      if (audio.duration && isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
        setIsLoaded(true);
        setAudioError(null);
      } else if (databaseDuration && databaseDuration > 0) {
        console.log('Using database duration as fallback:', databaseDuration);
        setDuration(databaseDuration);
        setIsLoaded(true);
        setAudioError(null);
      } else if (audio.duration === Infinity || isNaN(audio.duration)) {
        console.warn('Invalid audio duration detected:', audio.duration);
        if (!databaseDuration || databaseDuration <= 0) {
          setAudioError('Unable to determine audio duration. The audio file may be corrupted or streaming.');
        }
      }
    };
    
    const handleLoadedData = () => {
      console.log('handleLoadedData called, audio.duration:', audio.duration);
      if (audio.duration && isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
        setIsLoaded(true);
        setAudioError(null);
      } else if (databaseDuration && databaseDuration > 0) {
        console.log('Using database duration as fallback in handleLoadedData:', databaseDuration);
        setDuration(databaseDuration);
        setIsLoaded(true);
        setAudioError(null);
      }
    };
    
    const handleCanPlay = () => {
      console.log('handleCanPlay called, audio.duration:', audio.duration);
      if (audio.duration && isFinite(audio.duration) && audio.duration > 0) {
        setDuration(audio.duration);
        setIsLoaded(true);
        setAudioError(null);
      } else if (databaseDuration && databaseDuration > 0) {
        console.log('Using database duration as fallback in handleCanPlay:', databaseDuration);
        setDuration(databaseDuration);
        setIsLoaded(true);
        setAudioError(null);
      }
    };
    
    const handleEnded = () => setIsPlaying(false);
    
    const handlePlay = () => setIsPlaying(true);
    const handlePause = () => setIsPlaying(false);
    
    const handleError = () => {
      setAudioError('Unable to load audio file. The file may be corrupted or in an unsupported format.');
      setIsLoaded(false);
      setIsPlaying(false);
    };

    audio.addEventListener('timeupdate', updateTime);
    audio.addEventListener('loadedmetadata', updateDuration);
    audio.addEventListener('loadeddata', handleLoadedData);
    audio.addEventListener('canplay', handleCanPlay);
    audio.addEventListener('ended', handleEnded);
    audio.addEventListener('play', handlePlay);
    audio.addEventListener('pause', handlePause);
    audio.addEventListener('error', handleError);

    // Try to load duration immediately if already available
    if (audio.readyState >= 1) {
      updateDuration();
    }

    return () => {
      audio.removeEventListener('timeupdate', updateTime);
      audio.removeEventListener('loadedmetadata', updateDuration);
      audio.removeEventListener('loadeddata', handleLoadedData);
      audio.removeEventListener('canplay', handleCanPlay);
      audio.removeEventListener('ended', handleEnded);
      audio.removeEventListener('play', handlePlay);
      audio.removeEventListener('pause', handlePause);
      audio.removeEventListener('error', handleError);
    };
  }, [databaseDuration]);

  // Initialize with database duration if available and audio duration isn't loaded
  useEffect(() => {
    if (databaseDuration && databaseDuration > 0 && (!duration || duration <= 0)) {
      console.log('Initializing with database duration:', databaseDuration);
      setDuration(databaseDuration);
      setIsLoaded(true);
      setAudioError(null);
    }
  }, [databaseDuration, duration]);

  const togglePlay = () => {
    const audio = audioRef.current;
    if (!audio) return;

    if (isPlaying) {
      audio.pause();
    } else {
      audio.play();
    }
    setIsPlaying(!isPlaying);
  };

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const audio = audioRef.current;
    if (!audio || !duration || !isFinite(duration) || duration <= 0) return;

    const newTime = parseFloat(e.target.value);
    if (!isNaN(newTime) && isFinite(newTime) && newTime >= 0 && newTime <= duration) {
      audio.currentTime = newTime;
      setCurrentTime(newTime);
    }
  };

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const audio = audioRef.current;
    if (!audio) return;

    const newVolume = parseFloat(e.target.value);
    audio.volume = newVolume;
    setVolume(newVolume);
  };

  const formatTimeDisplay = (time: number) => {
    if (!time || isNaN(time) || !isFinite(time)) {
      return '0:00';
    }
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Audio Player</h3>
      
      <audio ref={audioRef} src={src} preload="metadata" />
      
      <div className="space-y-4">
        {/* Play/Pause Button */}
        <div className="flex items-center justify-center">
          <button
            onClick={togglePlay}
            className="flex items-center justify-center w-12 h-12 bg-blue-600 hover:bg-blue-700 rounded-full text-white transition-colors"
            aria-label={isPlaying ? 'Pause' : 'Play'}
          >
            {isPlaying ? (
              <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
              </svg>
            ) : (
              <svg className="w-6 h-6 ml-1" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z" />
              </svg>
            )}
          </button>
        </div>

        {/* Progress Bar */}
        <div className="space-y-2">
          <input
            type="range"
            min="0"
            max={isLoaded && duration && isFinite(duration) && duration > 0 ? duration : 100}
            value={isLoaded && isFinite(currentTime) ? currentTime : 0}
            onChange={handleSeek}
            disabled={!isLoaded || !duration || !isFinite(duration) || duration <= 0}
            className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer slider disabled:opacity-50 disabled:cursor-not-allowed"
          />
          <div className="flex justify-between text-sm text-gray-500 dark:text-gray-400">
            <span>{formatTimeDisplay(currentTime)}</span>
            <span>{isLoaded ? formatTimeDisplay(duration) : 'Loading...'}</span>
          </div>
        </div>

        {/* Volume Control */}
        <div className="flex items-center space-x-2">
          <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z" />
          </svg>
          <input
            type="range"
            min="0"
            max="1"
            step="0.1"
            value={volume}
            onChange={handleVolumeChange}
            className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer slider"
          />
          <span className="text-sm text-gray-500 dark:text-gray-400 w-8 text-right">
            {Math.round(volume * 100)}%
          </span>
        </div>

        {/* Error Message */}
        {audioError && (
          <div className="text-center text-red-600 dark:text-red-400 text-sm bg-red-50 dark:bg-red-900/20 p-3 rounded-lg">
            <svg className="w-5 h-5 mx-auto mb-1" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {audioError}
          </div>
        )}

        {/* File Info */}
        {!audioError && (
          <div className="text-sm text-gray-500 dark:text-gray-400 text-center">
            Playing: {filename}
          </div>
        )}
      </div>
    </div>
  );
}

export default function RecordingDetailsPage() {
  const params = useParams();
  const router = useRouter();
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

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <span className="text-gray-500 dark:text-gray-400">Loading recording...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 mb-4">
            <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Error Loading Recording</h3>
          <p className="text-gray-500 dark:text-gray-400 mb-4">{error}</p>
          <button
            onClick={() => router.push('/recordings')}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
          >
            Back to Recordings
          </button>
        </div>
      </div>
    );
  }

  if (!recording) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">Recording Not Found</h3>
          <button
            onClick={() => router.push('/recordings')}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
          >
            Back to Recordings
          </button>
        </div>
      </div>
    );
  }

  // Convert nanoseconds to seconds, then to minutes
  const durationSeconds = recording.duration / (1000 * 1000 * 1000);
  const durationMinutes = Math.round(durationSeconds / 60);
  
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => router.push('/recordings')}
            className="flex items-center text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 mb-4"
          >
            <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Back to Recordings
          </button>
          
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">
            {recording.filename.replace(/\.[^/.]+$/, '')}
          </h1>
          <p className="text-gray-600 dark:text-gray-300">
            Recording Details
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Audio Player */}
          <div className="lg:col-span-2">
            <AudioPlayer 
              src={`/api/recordings/${recording.id}/audio`}
              filename={recording.filename}
              databaseDuration={durationSeconds}
            />
          </div>

          {/* Recording Details */}
          <div className="space-y-6">
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Recording Information</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    File Name
                  </label>
                  <p className="text-gray-900 dark:text-white">{recording.filename}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Recording ID
                  </label>
                  <p className="text-gray-900 dark:text-white">#{recording.id}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Start Time
                  </label>
                  <p className="text-gray-900 dark:text-white">{formatDate(recording.start_time)}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    End Time
                  </label>
                  <p className="text-gray-900 dark:text-white">{formatDate(recording.end_time)}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Duration
                  </label>
                  <p className="text-gray-900 dark:text-white">{formatDuration(durationMinutes)}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    File Size
                  </label>
                  <p className="text-gray-900 dark:text-white">{formatFileSize(recording.file_size)}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Format
                  </label>
                  <p className="text-gray-900 dark:text-white">{recording.format.toUpperCase()}</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Audio Quality
                  </label>
                  <div className="text-gray-900 dark:text-white">
                    <p>{recording.sample_rate / 1000}kHz</p>
                    <p>{recording.channels} {recording.channels === 1 ? 'channel' : 'channels'}</p>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1">
                    Created
                  </label>
                  <p className="text-gray-900 dark:text-white">{formatDate(recording.created_at)}</p>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Actions</h3>
              
              <div className="space-y-3">
                <a
                  href={`/api/recordings/${recording.id}/audio`}
                  download={recording.filename}
                  className="block w-full bg-blue-600 hover:bg-blue-700 text-white text-center px-4 py-2 rounded-lg transition-colors"
                >
                  Download Recording
                </a>
                
                <button
                  onClick={() => {
                    navigator.clipboard.writeText(window.location.href);
                    // You could add a toast notification here
                  }}
                  className="block w-full bg-gray-600 hover:bg-gray-700 text-white text-center px-4 py-2 rounded-lg transition-colors"
                >
                  Copy Link
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <style jsx>{`
        .slider::-webkit-slider-thumb {
          appearance: none;
          height: 16px;
          width: 16px;
          border-radius: 50%;
          background: #3b82f6;
          cursor: pointer;
          border: 2px solid #ffffff;
          box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
        }
        
        .slider::-moz-range-thumb {
          height: 16px;
          width: 16px;
          border-radius: 50%;
          background: #3b82f6;
          cursor: pointer;
          border: 2px solid #ffffff;
          box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1);
        }
      `}</style>
    </div>
  );
}
