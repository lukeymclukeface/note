'use client';

import { useState, useRef, useEffect } from 'react';
import { Play, Pause, Volume2, AlertCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface AudioPlayerProps {
  src: string;
  filename: string;
  databaseDuration?: number; // Duration in seconds from database
}

export default function AudioPlayer({ src, filename, databaseDuration }: AudioPlayerProps) {
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
    
    const handleError = (e: Event) => {
      console.error('Audio error event:', e);
      const audioElement = e.target as HTMLAudioElement;
      console.error('Audio error details:', {
        error: audioElement.error,
        networkState: audioElement.networkState,
        readyState: audioElement.readyState,
        src: audioElement.src
      });
      
      // Only show error if we don't have database duration as fallback
      if (!databaseDuration || databaseDuration <= 0) {
        setAudioError('Unable to load audio file. The file may be corrupted or in an unsupported format.');
        setIsLoaded(false);
        setIsPlaying(false);
      } else {
        // Use database duration and hide the error
        console.log('Using database duration as fallback despite audio error:', databaseDuration);
        setDuration(databaseDuration);
        setIsLoaded(true);
        setAudioError(null);
      }
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
    <div>
      <audio ref={audioRef} src={src} preload="metadata" />
      
      <div className="space-y-6">
        {/* Play/Pause Button */}
        <div className="flex items-center justify-center">
          <Button
            onClick={togglePlay}
            size="lg"
            className="w-16 h-16 rounded-full bg-primary hover:bg-primary/90"
            aria-label={isPlaying ? 'Pause' : 'Play'}
          >
            {isPlaying ? (
              <Pause className="w-6 h-6" />
            ) : (
              <Play className="w-6 h-6 ml-1" />
            )}
          </Button>
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
            className="w-full h-2 bg-secondary rounded-lg appearance-none cursor-pointer slider disabled:opacity-50 disabled:cursor-not-allowed"
          />
          <div className="flex justify-between text-sm text-muted-foreground">
            <span>{formatTimeDisplay(currentTime)}</span>
            <span>{isLoaded ? formatTimeDisplay(duration) : 'Loading...'}</span>
          </div>
        </div>

        {/* Volume Control */}
        <div className="flex items-center space-x-3">
          <Volume2 className="w-4 h-4 text-muted-foreground" />
          <input
            type="range"
            min="0"
            max="1"
            step="0.1"
            value={volume}
            onChange={handleVolumeChange}
            className="flex-1 h-2 bg-secondary rounded-lg appearance-none cursor-pointer slider"
          />
          <span className="text-sm text-muted-foreground w-8 text-right">
            {Math.round(volume * 100)}%
          </span>
        </div>

        {/* Error Message */}
        {audioError && (
          <div className="text-center text-destructive text-sm bg-destructive/10 p-3 rounded-lg">
            <AlertCircle className="w-5 h-5 mx-auto mb-1" />
            {audioError}
          </div>
        )}

        {/* File Info */}
        {!audioError && (
          <div className="text-sm text-muted-foreground text-center">
            Playing: {filename}
          </div>
        )}
      </div>
      
      <style jsx>{`
        .slider::-webkit-slider-thumb {
          appearance: none;
          height: 16px;
          width: 16px;
          border-radius: 50%;
          background: hsl(var(--primary));
          cursor: pointer;
          border: 2px solid hsl(var(--background));
          box-shadow: 0 0 0 1px hsl(var(--border));
        }
        
        .slider::-moz-range-thumb {
          height: 16px;
          width: 16px;
          border-radius: 50%;
          background: hsl(var(--primary));
          cursor: pointer;
          border: 2px solid hsl(var(--background));
          box-shadow: 0 0 0 1px hsl(var(--border));
        }
      `}</style>
    </div>
  );
}
