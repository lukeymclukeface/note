// Get the start of the week (Sunday) for a given date
export function getWeekStart(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay();
  const diff = d.getDate() - day;
  const weekStart = new Date(d.setDate(diff));
  weekStart.setHours(0, 0, 0, 0);
  return weekStart;
}

// Get the end of the week (Saturday) for a given date
export function getWeekEnd(date: Date): Date {
  const weekStart = getWeekStart(date);
  const weekEnd = new Date(weekStart);
  weekEnd.setDate(weekStart.getDate() + 6);
  weekEnd.setHours(23, 59, 59, 999);
  return weekEnd;
}

// Get an array of dates for the week
export function getWeekDates(weekStart: Date): Date[] {
  const dates = [];
  for (let i = 0; i < 7; i++) {
    const date = new Date(weekStart);
    date.setDate(weekStart.getDate() + i);
    dates.push(date);
  }
  return dates;
}

// Format date for display
export function formatDate(date: Date): string {
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  });
}

// Format time for display
export function formatTime(date: Date): string {
  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });
}

// Format duration in minutes to readable string
export function formatDuration(minutes: number): string {
  if (minutes < 60) {
    return `${minutes}m`;
  }
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  if (remainingMinutes === 0) {
    return `${hours}h`;
  }
  return `${hours}h ${remainingMinutes}m`;
}

// Get hours array for calendar grid
export function getHoursArray(): number[] {
  return Array.from({ length: 24 }, (_, i) => i);
}

// Calculate position of event in calendar grid
export function getEventPosition(start: Date, end: Date, dayStart: Date): {
  top: number;
  height: number;
} {
  const dayStartTime = dayStart.getTime();
  const eventStart = Math.max(start.getTime(), dayStartTime);
  const eventEnd = Math.min(end.getTime(), dayStartTime + 24 * 60 * 60 * 1000);
  
  const startMinutes = (eventStart - dayStartTime) / (1000 * 60);
  const durationMinutes = (eventEnd - eventStart) / (1000 * 60);
  
  const hourHeight = 60; // pixels per hour
  const top = (startMinutes / 60) * hourHeight;
  const height = Math.max((durationMinutes / 60) * hourHeight, 20); // Minimum 20px height
  
  return { top, height };
}

// Check if two events overlap
export function eventsOverlap(event1: { start: Date; end: Date }, event2: { start: Date; end: Date }): boolean {
  return event1.start < event2.end && event2.start < event1.end;
}

// Get current week start
export function getCurrentWeekStart(): Date {
  return getWeekStart(new Date());
}

// Navigate to previous week
export function getPreviousWeek(currentWeekStart: Date): Date {
  const prevWeek = new Date(currentWeekStart);
  prevWeek.setDate(currentWeekStart.getDate() - 7);
  return prevWeek;
}

// Navigate to next week
export function getNextWeek(currentWeekStart: Date): Date {
  const nextWeek = new Date(currentWeekStart);
  nextWeek.setDate(currentWeekStart.getDate() + 7);
  return nextWeek;
}

// Check if date is today
export function isToday(date: Date): boolean {
  const today = new Date();
  return date.toDateString() === today.toDateString();
}

// Get week range string for display
export function getWeekRangeString(weekStart: Date): string {
  const weekEnd = getWeekEnd(weekStart);
  
  if (weekStart.getMonth() === weekEnd.getMonth()) {
    return `${weekStart.toLocaleDateString('en-US', { month: 'long', day: 'numeric' })} - ${weekEnd.getDate()}, ${weekStart.getFullYear()}`;
  } else {
    return `${weekStart.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${weekEnd.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}, ${weekStart.getFullYear()}`;
  }
}
