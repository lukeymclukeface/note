'use client';

import { useState, useRef, useEffect } from 'react';
import { themes, useThemes } from '@/lib/useThemes';
import { ChevronDown, Check } from 'lucide-react';

export function ThemeSelector() {
  const { theme, setTheme } = useThemes();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const currentTheme = themes.find(t => t.value === theme) || themes[0];

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
    <div className="relative" data-testid="theme-selector" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="inline-flex items-center px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
        aria-expanded={isOpen}
        aria-haspopup="true"
      >
        <currentTheme.icon className="mr-2 h-4 w-4" />
        <span className="hidden sm:inline">{currentTheme.label}</span>
        <ChevronDown className="ml-2 -mr-1 h-4 w-4" />
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-50">
          <div className="py-1" role="menu" aria-orientation="vertical">
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
                role="menuitem"
              >
              <themeOption.icon className="mr-3 h-4 w-4" />
              <span>{themeOption.label}</span>
                {theme === themeOption.value && (
                  <Check className="ml-auto h-4 w-4 text-blue-600 dark:text-blue-400" />
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
