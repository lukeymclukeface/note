'use client';

import { useState } from 'react';
import { CalendarEvent } from '@/lib/database';
import CalendarEventComponent from './CalendarEvent';
import {
  getWeekDates,
  formatDate,
  getHoursArray,
  getEventPosition,
  isToday,
  getPreviousWeek,
  getNextWeek,
  getCurrentWeekStart,
  getWeekRangeString
} from '@/lib/dateUtils';

interface WeekCalendarProps {
  events: CalendarEvent[];
}

export default function WeekCalendar({ events }: WeekCalendarProps) {
  const [currentWeekStart, setCurrentWeekStart] = useState(getCurrentWeekStart());
  
  const weekDates = getWeekDates(currentWeekStart);
  const hours = getHoursArray();
  
  // Filter events for current week
  const weekEvents = events.filter(event => {
    const eventDate = new Date(event.start);
    return eventDate >= currentWeekStart && eventDate < getNextWeek(currentWeekStart);
  });
  
  // Group events by day
  const eventsByDay = weekDates.map(date => {
    const dayStart = new Date(date);
    dayStart.setHours(0, 0, 0, 0);
    const dayEnd = new Date(date);
    dayEnd.setHours(23, 59, 59, 999);
    
    return weekEvents.filter(event => {
      const eventStart = new Date(event.start);
      return eventStart >= dayStart && eventStart <= dayEnd;
    });
  });
  
  const navigateToPrevWeek = () => {
    setCurrentWeekStart(getPreviousWeek(currentWeekStart));
  };
  
  const navigateToNextWeek = () => {
    setCurrentWeekStart(getNextWeek(currentWeekStart));
  };
  
  const goToToday = () => {
    setCurrentWeekStart(getCurrentWeekStart());
  };
  
  const formatHour = (hour: number) => {
    if (hour === 0) return '12 AM';
    if (hour === 12) return '12 PM';
    if (hour < 12) return `${hour} AM`;
    return `${hour - 12} PM`;
  };

  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-sm">
      {/* Calendar Header */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              {getWeekRangeString(currentWeekStart)}
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              {weekEvents.length} recording{weekEvents.length !== 1 ? 's' : ''} this week
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={navigateToPrevWeek}
              className="p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
              title="Previous week"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <button
              onClick={goToToday}
              className="px-3 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
            >
              Today
            </button>
            <button
              onClick={navigateToNextWeek}
              className="p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
              title="Next week"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Calendar Grid */}
      <div className="overflow-hidden">
        {/* Days Header */}
        <div className="grid grid-cols-8 border-b border-gray-200 dark:border-gray-700">
          <div className="p-2 text-xs font-medium text-gray-500 dark:text-gray-400 text-center border-r border-gray-200 dark:border-gray-700">
            Time
          </div>
          {weekDates.map((date, index) => (
            <div
              key={index}
              className={`p-2 text-sm font-medium text-center border-r border-gray-200 dark:border-gray-700 last:border-r-0 ${
                isToday(date) ? 'bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300' : 'text-gray-700 dark:text-gray-300'
              }`}
            >
              <div className={isToday(date) ? 'font-semibold' : ''}>{formatDate(date)}</div>
            </div>
          ))}
        </div>
        
        {/* Calendar Body */}
        <div className="relative">
          <div className="grid grid-cols-8">
            {/* Time Column */}
            <div className="border-r border-gray-200 dark:border-gray-700">
              {hours.map(hour => (
                <div
                  key={hour}
                  className="h-[60px] border-b border-gray-100 dark:border-gray-700 p-1 text-xs text-gray-500 dark:text-gray-400 text-right"
                  style={{ lineHeight: '1' }}
                >
                  {hour === 0 ? '' : formatHour(hour)}
                </div>
              ))}
            </div>
            
            {/* Day Columns */}
            {weekDates.map((date, dayIndex) => {
              const dayStart = new Date(date);
              dayStart.setHours(0, 0, 0, 0);
              
              return (
                <div key={dayIndex} className="relative border-r border-gray-200 dark:border-gray-700 last:border-r-0">
                  {/* Hour Grid */}
                  {hours.map(hour => (
                    <div
                      key={hour}
                      className={`h-[60px] border-b border-gray-100 dark:border-gray-700 ${
                        isToday(date) ? 'bg-blue-25 dark:bg-blue-950' : ''
                      }`}
                    />
                  ))}
                  
                  {/* Events */}
                  {eventsByDay[dayIndex].map(event => {
                    const position = getEventPosition(event.start, event.end, dayStart);
                    return (
                      <CalendarEventComponent
                        key={event.id}
                        event={event}
                        style={position}
                      />
                    );
                  })}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
