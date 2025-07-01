'use client';

import { useState, useEffect } from 'react';
import { loadRecordings } from '@/lib/actions/recordings';
import { formatTime, formatDuration } from '@/lib/dateUtils';
import type { Recording } from '@/lib/database';

export default function RecordingsPage() {
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
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-6xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">Audio Recordings</h1>
          <p className="text-gray-600 dark:text-gray-300">
            {recordings.length === 0 
              ? 'No recordings found. Use the record button in the navigation bar to start recording.' 
              : `${recordings.length} recording${recordings.length !== 1 ? 's' : ''} found`
            }
          </p>
        </header>

        {isLoading ? (
          <div className="text-center py-4">
            <span className="text-gray-500 dark:text-gray-400">Loading recordings...</span>
          </div>
        ) : recordings.length === 0 ? (
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-8 text-center">
            <div className="text-gray-400 dark:text-gray-500 mb-4">
              <svg className="mx-auto h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No recordings yet</h3>
            <p className="text-gray-500 dark:text-gray-400 mb-4">
              Get started by clicking the record button in the navigation bar.
            </p>
          </div>
        ) : (
          <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead className="bg-gray-50 dark:bg-gray-900">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Recording
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Start Time
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Duration
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Format
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Size
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Quality
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                  {recordings.map((recording) => {
                    // Convert nanoseconds to seconds, then to minutes
                    const durationSeconds = recording.duration / (1000 * 1000 * 1000);
                    const durationMinutes = Math.round(durationSeconds / 60);
                    const startTime = new Date(recording.start_time);
                    const endTime = new Date(recording.end_time);
                    
                    return (
                      <tr key={recording.id} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-8 w-8">
                              <div className="h-8 w-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                                <svg className="h-4 w-4 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                                </svg>
                              </div>
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-gray-900 dark:text-white">
                                {recording.filename}
                              </div>
                              <div className="text-sm text-gray-500 dark:text-gray-400">
                                ID: {recording.id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900 dark:text-white">
                            {formatDate(recording.start_time)}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {formatTime(startTime)} - {formatTime(endTime)}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200">
                            {formatDuration(durationMinutes)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                          {recording.format.toUpperCase()}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                          {formatFileSize(recording.file_size)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                          <div>{recording.sample_rate / 1000}kHz</div>
                          <div>{recording.channels} ch</div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
