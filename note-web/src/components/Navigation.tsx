'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { clsx } from 'clsx';
import { useState, useRef, useEffect } from 'react';
import { Home, FileText, Users, Briefcase, Calendar, Import, Mic, Upload, User, Settings, Sun, Moon, Monitor, ChevronDown } from 'lucide-react';
import NavbarRecorder from './NavbarRecorder';
import { UserDropdownMenu } from './UserDropdownMenu';
import { useTheme } from '@/providers/ThemeProvider';

const navigation = [
  { name: 'Dashboard', href: '/', icon: Home },
  { name: 'Notes', href: '/notes', icon: FileText },
  { name: 'Meetings', href: '/meetings', icon: Users },
  { name: 'Interviews', href: '/interviews', icon: Briefcase },
  { name: 'Calendar', href: '/calendar', icon: Calendar },
];

const importNavigation = [
  { name: 'Recordings', href: '/recordings', icon: Mic },
  { name: 'Upload', href: '/upload', icon: Upload },
];

const themes = [
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
  { value: 'system', label: 'System', icon: Monitor },
] as const;

export default function Navigation() {
  const pathname = usePathname();
  const [isImportDropdownOpen, setIsImportDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const { theme, setTheme } = useTheme();

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsImportDropdownOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Handle keyboard navigation
  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Escape') {
      setIsImportDropdownOpen(false);
    }
  };

  return (
    <nav className="bg-white dark:bg-gray-900 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="px-4 w-full">
        <div className="flex justify-between h-16">
          <div className="flex">
            {/* Logo */}
            <div className="flex-shrink-0 flex items-center">
              <Link href="/" className="text-xl font-bold text-gray-900 dark:text-white">
                Note AI
              </Link>
            </div>
            
            {/* Navigation Links */}
            <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={clsx(
                    pathname === item.href
                      ? 'border-blue-500 text-gray-900 dark:text-white'
                      : 'border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                    'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium'
                  )}
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              ))}
              
              {/* Import Dropdown */}
              <div className="relative" ref={dropdownRef}>
                <button
                  type="button"
                  onClick={() => setIsImportDropdownOpen(!isImportDropdownOpen)}
                  onKeyDown={handleKeyDown}
                  className={clsx(
                    importNavigation.some(item => pathname === item.href)
                      ? 'border-blue-500 text-gray-900 dark:text-white'
                      : 'border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                    'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium h-16 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors'
                  )}
                  aria-expanded={isImportDropdownOpen}
                  aria-haspopup="true"
                  tabIndex={0}
                >
                  <Import className="mr-1 h-4 w-4" />
                  Import
                  <ChevronDown className="ml-1 h-4 w-4" />
                </button>
                
                {isImportDropdownOpen && (
                  <div className="absolute left-0 mt-2 w-48 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-50">
                    <div className="py-1">
                      {importNavigation.map((item) => (
                        <Link
                          key={item.name}
                          href={item.href}
                          onClick={() => setIsImportDropdownOpen(false)}
                          className={clsx(
                            pathname === item.href
                              ? 'bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300'
                              : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700',
                            'block px-4 py-2 text-sm'
                          )}
                        >
                          <item.icon className="mr-2 h-4 w-4" />
                          {item.name}
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
          
          {/* Right side: Secondary Navigation + Mobile Menu */}
          <div className="flex items-center">
{/* Desktop User Menu */}
            <div className="hidden sm:flex sm:items-center sm:space-x-4">
              <NavbarRecorder />
              <UserDropdownMenu />
            </div>
            
            {/* Mobile menu button */}
            <div className="sm:hidden flex items-center space-x-3">
              <UserDropdownMenu />
              <button
                type="button"
                className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 dark:text-gray-300 hover:text-gray-500 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800"
                aria-controls="mobile-menu"
                aria-expanded="false"
              >
                <span className="sr-only">Open main menu</span>
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
        
        {/* Mobile menu */}
        <div className="sm:hidden" id="mobile-menu">
          <div className="pt-2 pb-3 space-y-1">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={clsx(
                  pathname === item.href
                    ? 'bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-700 dark:text-blue-300'
                    : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                  'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
                )}
              >
                <item.icon className="mr-2 h-4 w-4" />
                {item.name}
              </Link>
            ))}
            
            {/* Import section in mobile */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3 mt-3">
              <div className="pl-3 pr-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider flex items-center">
                <Import className="mr-2 h-4 w-4" />
                Import
              </div>
              {importNavigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={clsx(
                    pathname === item.href
                      ? 'bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-700 dark:text-blue-300'
                      : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                    'block pl-6 pr-4 py-2 border-l-4 text-base font-medium'
                  )}
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              ))}
            </div>
            
            {/* User section in mobile */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3">
              <div className="pl-3 pr-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider flex items-center">
                <User className="mr-2 h-4 w-4" />
                User
              </div>
              
              {/* Settings Link */}
              <Link
                href="/settings"
                className={clsx(
                  pathname === '/settings'
                    ? 'bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-700 dark:text-blue-300'
                    : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                  'block pl-6 pr-4 py-2 border-l-4 text-base font-medium'
                )}
              >
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </Link>
              
              {/* Theme Options */}
              <div className="pl-6 pr-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Theme
              </div>
              {themes.map((themeOption) => (
                <button
                  key={themeOption.value}
                  onClick={() => setTheme(themeOption.value)}
                  className={clsx(
                    theme === themeOption.value
                      ? 'bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-700 dark:text-blue-300'
                      : 'border-transparent text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                    'w-full text-left block pl-9 pr-4 py-2 border-l-4 text-base font-medium'
                  )}
                >
                  <themeOption.icon className="mr-2 h-4 w-4" />
                  {themeOption.label}
                  {theme === themeOption.value && (
                    <svg
                      className="inline ml-2 h-4 w-4 text-blue-600 dark:text-blue-400"
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
        </div>
      </div>
    </nav>
  );
}
