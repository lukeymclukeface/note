'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { clsx } from 'clsx';
import { useState, useRef, useEffect } from 'react';
import { ThemeSelector } from './ThemeSelector';

const navigation = [
  { name: 'Dashboard', href: '/', icon: '🏠' },
  { name: 'Notes', href: '/notes', icon: '📝' },
  { name: 'Meetings', href: '/meetings', icon: '🤝' },
  { name: 'Interviews', href: '/interviews', icon: '💼' },
  { name: 'Calendar', href: '/calendar', icon: '📅' },
];

const importNavigation = [
  { name: 'Recordings', href: '/recordings', icon: '🎤' },
  { name: 'Upload', href: '/upload', icon: '📤' },
];

const secondaryNavigation = [
  { name: 'Settings', href: '/settings', icon: '⚙️' },
];

export default function Navigation() {
  const pathname = usePathname();
  const [isImportDropdownOpen, setIsImportDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

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

  return (
    <nav className="bg-white dark:bg-gray-900 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-6xl mx-auto px-4">
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
                  <span className="mr-2">{item.icon}</span>
                  {item.name}
                </Link>
              ))}
              
              {/* Import Dropdown */}
              <div className="relative" ref={dropdownRef}>
                <button
                  type="button"
                  onClick={() => setIsImportDropdownOpen(!isImportDropdownOpen)}
                  className={clsx(
                    importNavigation.some(item => pathname === item.href)
                      ? 'border-blue-500 text-gray-900 dark:text-white'
                      : 'border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                    'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium h-16'
                  )}
                >
                  <span className="mr-1">📥</span>
                  Import
                  <svg className="ml-1 h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
                  </svg>
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
                          <span className="mr-2">{item.icon}</span>
                          {item.name}
                        </Link>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          
          {/* Secondary Navigation */}
          <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
            {secondaryNavigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={clsx(
                  pathname === item.href
                    ? 'border-blue-500 text-gray-900 dark:text-white'
                    : 'border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                  'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium h-16'
                )}
              >
                <span className="mr-2">{item.icon}</span>
                {item.name}
              </Link>
            ))}
            <div className="flex items-center">
              <ThemeSelector />
            </div>
          </div>
        </div>
        
        {/* Mobile menu button and theme selector */}
          <div className="sm:hidden flex items-center space-x-3">
            <ThemeSelector />
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
                <span className="mr-2">{item.icon}</span>
                {item.name}
              </Link>
            ))}
            
            {/* Import section in mobile */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3 mt-3">
              <div className="pl-3 pr-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                📥 Import
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
                  <span className="mr-2">{item.icon}</span>
                  {item.name}
                </Link>
              ))}
            </div>
            
            {/* Secondary Navigation in Mobile */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3">
              {secondaryNavigation.map((item) => (
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
                  <span className="mr-2">{item.icon}</span>
                  {item.name}
                </Link>
              ))}
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
