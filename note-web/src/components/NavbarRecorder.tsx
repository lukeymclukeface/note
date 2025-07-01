'use client';

import { useState, useRef, useCallback } from 'react';
import { usePathname, useRouter } from 'next/navigation';

interface NavbarRecorderProps {
  onRecordingComplete?: () => void;
}

export default function NavbarRecorder({ onRecordingComplete }: NavbarRecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [recordingTime, setRecordingTime] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const recordingStartTimeRef = useRef<Date | null>(null);
  
  const pathname = usePathname();
  const router = useRouter();

  const startRecording = useCallback(async () => {
    try {
      setError(null);
      
      // Request microphone access
      const stream = await navigator.mediaDevices.getUserMedia({ 
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          sampleRate: 44100
        } 
      });
      
      streamRef.current = stream;
      
      // Create MediaRecorder
      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: 'audio/webm;codecs=opus'
      });
      
      mediaRecorderRef.current = mediaRecorder;
      chunksRef.current = [];

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          chunksRef.current.push(event.data);
        }
      };

      mediaRecorder.onstop = async () => {
        setIsProcessing(true);
        
        const audioBlob = new Blob(chunksRef.current, { type: 'audio/webm' });
        const recordingEndTime = new Date();
        
        // Upload the recording and create database entry
        const formData = new FormData();
        formData.append('audio', audioBlob, 'recording.webm');
        if (recordingStartTimeRef.current) {
          formData.append('startTime', recordingStartTimeRef.current.toISOString());
          formData.append('endTime', recordingEndTime.toISOString());
        }

        try {
          const response = await fetch('/api/upload-recording', {
            method: 'POST',
            body: formData
          });

          if (!response.ok) {
            throw new Error('Failed to upload recording');
          }

          const result = await response.json();
          
          if (result.success) {
            onRecordingComplete?.();
            
            // Emit a custom event for real-time updates
            const event = new CustomEvent('recordingCompleted', {
              detail: {
                recordingId: result.recordingId,
                filename: result.filename,
                size: result.size,
                duration: result.duration
              }
            });
            window.dispatchEvent(event);
            
            // If we're on the recordings page, refresh it as fallback
            if (pathname === '/recordings') {
              router.refresh();
            }
          } else {
            throw new Error(result.error || 'Upload failed');
          }
        } catch (uploadError) {
          console.error('Upload error:', uploadError);
          setError('Failed to save recording');
        }
        
        setIsProcessing(false);
        
        // Clean up
        if (streamRef.current) {
          streamRef.current.getTracks().forEach(track => track.stop());
          streamRef.current = null;
        }
      };

      // Start recording
      recordingStartTimeRef.current = new Date();
      mediaRecorder.start(1000);
      setIsRecording(true);
      setRecordingTime(0);

      // Start timer
      timerRef.current = setInterval(() => {
        setRecordingTime(prev => prev + 1);
      }, 1000);

    } catch (err) {
      console.error('Error starting recording:', err);
      if (err instanceof Error) {
        if (err.name === 'NotAllowedError') {
          setError('Microphone access denied');
        } else if (err.name === 'NotFoundError') {
          setError('No microphone found');
        } else {
          setError('Failed to start recording');
        }
      } else {
        setError('Recording error');
      }
    }
  }, [onRecordingComplete, pathname, router]);

  const stopRecording = useCallback(() => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
      
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    }
  }, [isRecording]);

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  // Don't render anything if there's an error - keep navbar clean
  if (error) {
    // Clear error after 3 seconds
    setTimeout(() => setError(null), 3000);
  }

  return (
    <div className="flex items-center">
      {!isRecording && !isProcessing ? (
        <button
          onClick={startRecording}
          className="flex items-center gap-2 px-3 py-2 text-gray-500 dark:text-gray-300 hover:text-red-600 dark:hover:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
          title="Start Recording"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
          </svg>
          <span className="hidden sm:block text-sm">Record</span>
        </button>
      ) : isRecording ? (
        <button
          onClick={stopRecording}
          className="flex items-center gap-2 px-3 py-2 text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
          title="Stop Recording"
        >
          <div className="w-3 h-3 bg-red-600 dark:bg-red-400 rounded-full animate-pulse"></div>
          <span className="text-sm font-mono">{formatTime(recordingTime)}</span>
        </button>
      ) : (
        <div className="flex items-center gap-2 px-3 py-2 text-gray-400 dark:text-gray-500">
          <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span className="hidden sm:block text-sm">Saving...</span>
        </div>
      )}
    </div>
  );
}
