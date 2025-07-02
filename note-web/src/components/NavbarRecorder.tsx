'use client';

import { useState, useRef, useCallback } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Mic, Loader2 } from 'lucide-react';

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
        <Button
          onClick={startRecording}
          variant="default"
          size="sm"
          className="bg-red-600 hover:bg-red-700 text-white border-red-600 hover:border-red-700 shadow-sm"
          title="Start Recording"
        >
          <Mic className="w-4 h-4" />
          <span className="text-sm ml-0.5 font-medium">Record</span>
        </Button>
      ) : isRecording ? (
        <Button
          onClick={stopRecording}
          variant="default"
          size="sm"
          className="bg-red-800 hover:bg-red-900 text-white border-red-800 hover:border-red-900 shadow-sm"
          title="Stop Recording"
        >
          <div className="w-3 h-3 bg-white rounded-sm animate-pulse"></div>
          <span className="text-sm font-mono ml-2 font-medium">{formatTime(recordingTime)}</span>
        </Button>
      ) : (
        <div className="flex items-center gap-2 px-3 py-2 text-gray-400 dark:text-gray-500">
          <Loader2 className="w-4 h-4 animate-spin" />
          <span className="hidden sm:block text-sm">Saving...</span>
        </div>
      )}
    </div>
  );
}
