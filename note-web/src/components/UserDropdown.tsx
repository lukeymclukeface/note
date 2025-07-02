'use client';

import Link from 'next/link';
import { useState, useRef, useEffect } from 'react';
import { useTheme } from '@/providers/ThemeProvider';
import { Settings, Sun, Moon, Monitor, ChevronDown, User } from 'lucide-react';

const themes = [
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
  { value: 'system', label: 'System', icon: Monitor },
] as const;

export function UserDropdown() {
  const { theme, setTheme } = useTheme();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="inline-flex items-center px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
        aria-expanded={isOpen}
        aria-haspopup="true"
      >
        <User className="text-lg" />
        <ChevronDown className="ml-2 -mr-1 h-4 w-4" />
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-50">
          <div className="py-1">
            {/* Settings Link */}
            <Link
              href="/settings"
              onClick={() => setIsOpen(false)}
              className="text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 block px-4 py-2 text-sm"
            >
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </Link>
            
            {/* Divider */}
            <div className="border-t border-gray-200 dark:border-gray-600 my-1"></div>
            
            {/* Theme Options */}
            <div className="px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Theme
            </div>
            {themes.map((themeOption) => (
              <button
                key={themeOption.value}
                onClick={() => {
                  setTheme(themeOption.value);
                  setIsOpen(false);
                }}
                className={`${
                  theme === themeOption.value
                    ? 'bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-gray-100'
                    : 'text-gray-700 dark:text-gray-200'
                } group flex items-center px-4 py-2 text-sm w-full text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors`}
              >
                <themeOption.icon className="mr-3 h-4 w-4" />
                <span>{themeOption.label}</span>
                {theme === themeOption.value && (
                  <svg
                    className="ml-auto h-4 w-4 text-blue-600 dark:text-blue-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
