'use client';

import { useState, useRef, useCallback } from 'react';
import { Mic, Square, Loader2 } from 'lucide-react';

interface AudioRecorderProps {
  onRecordingComplete?: () => void;
}

export default function AudioRecorder({ onRecordingComplete }: AudioRecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [recordingTime, setRecordingTime] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

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
        
        // Upload the recording
        const formData = new FormData();
        formData.append('audio', audioBlob, 'recording.webm');

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
          } else {
            throw new Error(result.error || 'Upload failed');
          }
        } catch (uploadError) {
          console.error('Upload error:', uploadError);
          setError('Failed to save recording. Please try again.');
        }
        
        setIsProcessing(false);
        
        // Clean up
        if (streamRef.current) {
          streamRef.current.getTracks().forEach(track => track.stop());
          streamRef.current = null;
        }
      };

      // Start recording
      mediaRecorder.start(1000); // Collect data every second
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
          setError('Microphone access denied. Please allow microphone access and try again.');
        } else if (err.name === 'NotFoundError') {
          setError('No microphone found. Please check your audio devices.');
        } else {
          setError('Failed to start recording. Please check your microphone.');
        }
      } else {
        setError('An unknown error occurred while starting recording.');
      }
    }
  }, [onRecordingComplete]);

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

  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-medium text-gray-900 dark:text-white">
          Record Audio
        </h3>
        {isRecording && (
          <div className="flex items-center text-red-600 dark:text-red-400">
            <div className="w-3 h-3 bg-red-600 dark:bg-red-400 rounded-full animate-pulse mr-2"></div>
            <span className="text-sm font-mono">{formatTime(recordingTime)}</span>
          </div>
        )}
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md">
          <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
        </div>
      )}

      <div className="flex items-center gap-4">
        {!isRecording && !isProcessing ? (
          <button
            onClick={startRecording}
            className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
          >
            <Mic className="w-5 h-5" />
            Start Recording
          </button>
        ) : isRecording ? (
          <button
            onClick={stopRecording}
            className="flex items-center gap-2 px-4 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded-lg transition-colors"
          >
            <Square className="w-5 h-5" />
            Stop Recording
          </button>
        ) : (
          <button
            disabled
            className="flex items-center gap-2 px-4 py-2 bg-gray-400 text-white rounded-lg cursor-not-allowed"
          >
            <Loader2 className="w-5 h-5 animate-spin" />
            Processing...
          </button>
        )}

        {isRecording && (
          <div className="text-sm text-gray-600 dark:text-gray-400">
            Recording audio from microphone...
          </div>
        )}
      </div>

      <div className="mt-4 text-xs text-gray-500 dark:text-gray-400">
        <p>• Recordings are saved to ~/.noteai/recordings/</p>
        <p>• Audio format: WebM with Opus codec</p>
        <p>• Make sure to allow microphone access when prompted</p>
      </div>
    </div>
  );
}
