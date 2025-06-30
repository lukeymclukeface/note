'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { clsx } from 'clsx';
import { ThemeSelector } from './ThemeSelector';

const navigation = [
  { name: 'Dashboard', href: '/', icon: '🏠' },
  { name: 'Notes', href: '/notes', icon: '📝' },
  { name: 'Meetings', href: '/meetings', icon: '🤝' },
  { name: 'Interviews', href: '/interviews', icon: '💼' },
  { name: 'Calendar', href: '/calendar', icon: '📅' },
  { name: 'Recordings', href: '/recordings', icon: '🎤' },
];

const secondaryNavigation = [
  { name: 'Settings', href: '/settings', icon: '⚙️' },
];

export default function Navigation() {
  const pathname = usePathname();

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
          </div>
          
          {/* Secondary Navigation */}
          <div className="hidden sm:ml-6 sm:flex sm:space-x-8 sm:items-center">
            {secondaryNavigation.map((item) => (
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
            <ThemeSelector />
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
