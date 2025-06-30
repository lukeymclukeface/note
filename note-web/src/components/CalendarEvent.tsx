import { CalendarEvent } from '@/lib/database';
import { formatTime, formatDuration } from '@/lib/dateUtils';

interface CalendarEventProps {
  event: CalendarEvent;
  style: {
    top: number;
    height: number;
  };
}

export default function CalendarEventComponent({ event, style }: CalendarEventProps) {
  const getEventColor = (type: string) => {
    switch (type) {
      case 'recording':
        return {
          bg: 'bg-blue-500',
          border: 'border-blue-600',
          text: 'text-white'
        };
      case 'note':
        return {
          bg: 'bg-green-500',
          border: 'border-green-600',
          text: 'text-white'
        };
      default:
        return {
          bg: 'bg-gray-500',
          border: 'border-gray-600',
          text: 'text-white'
        };
    }
  };

  const colors = getEventColor(event.type);
  const isShort = style.height < 40;

  return (
    <div
      className={`absolute left-1 right-1 rounded-md shadow-sm border ${colors.bg} ${colors.border} ${colors.text} p-2 cursor-pointer hover:shadow-md transition-shadow`}
      style={{
        top: `${style.top}px`,
        height: `${style.height}px`,
        zIndex: 10
      }}
      title={`${event.title}\n${formatTime(event.start)} - ${formatTime(event.end)}\nDuration: ${formatDuration(event.duration)}`}
    >
      <div className="overflow-hidden">
        <div className={`font-medium truncate ${isShort ? 'text-xs' : 'text-sm'}`}>
          {event.title}
        </div>
        {!isShort && (
          <>
            <div className="text-xs opacity-90 truncate">
              {formatTime(event.start)}
            </div>
            <div className="text-xs opacity-75">
              {formatDuration(event.duration)}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
