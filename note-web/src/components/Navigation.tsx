'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { clsx } from 'clsx';
import { Home, FileText, Users, Briefcase, Calendar, Import, Mic, Upload, User, Settings, Sun, Moon, Monitor, ChevronDown } from 'lucide-react';
import NavbarRecorder from './NavbarRecorder';
import { UserDropdownMenu } from './UserDropdownMenu';
import { useTheme } from '@/providers/ThemeProvider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

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
  const { theme, setTheme } = useTheme();

  return (
    <nav className="bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b border-border">
      <div className="px-4 w-full">
        <div className="flex justify-between h-16">
          <div className="flex">
            {/* Logo */}
            <div className="flex-shrink-0 flex items-center">
              <Link href="/" className="text-xl font-bold text-foreground">
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
                      ? 'border-primary text-foreground'
                      : 'border-transparent text-muted-foreground hover:border-border hover:text-foreground',
                    'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors'
                  )}
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              ))}
              
              {/* Import Dropdown */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    className={clsx(
                      importNavigation.some(item => pathname === item.href)
                        ? 'border-primary text-foreground'
                        : 'border-transparent text-muted-foreground hover:border-border hover:text-foreground',
                      'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium h-16 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 transition-colors'
                    )}
                  >
                    <Import className="mr-1 h-4 w-4" />
                    Import
                    <ChevronDown className="ml-1 h-4 w-4" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start">
                  {importNavigation.map((item) => (
                    <DropdownMenuItem key={item.name} asChild>
                      <Link
                        href={item.href}
                        className={clsx(
                          pathname === item.href
                            ? 'bg-accent text-accent-foreground'
                            : '',
                          'flex items-center w-full'
                        )}
                      >
                        <item.icon className="mr-2 h-4 w-4" />
                        {item.name}
                      </Link>
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
          
          {/* Right side: Secondary Navigation + Mobile Menu */}
          <div className="flex items-center">
{/* Desktop User Menu */}
            <div className="hidden sm:flex sm:items-center sm:space-x-3">
              <NavbarRecorder />
              <UserDropdownMenu />
            </div>
            
            {/* Mobile menu button */}
            <div className="sm:hidden flex items-center space-x-2">
              <UserDropdownMenu />
              <button
                type="button"
                className="inline-flex items-center justify-center p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent"
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
                    ? 'bg-primary/10 border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:bg-accent hover:border-border hover:text-foreground',
                  'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
                )}
              >
                <item.icon className="mr-2 h-4 w-4" />
                {item.name}
              </Link>
            ))}
            
            {/* Import section in mobile */}
            <div className="border-t border-border pt-3 mt-3">
              <div className="pl-3 pr-4 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider flex items-center">
                <Import className="mr-2 h-4 w-4" />
                Import
              </div>
              {importNavigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={clsx(
                    pathname === item.href
                      ? 'bg-primary/10 border-primary text-primary'
                      : 'border-transparent text-muted-foreground hover:bg-accent hover:border-border hover:text-foreground',
                    'block pl-6 pr-4 py-2 border-l-4 text-base font-medium'
                  )}
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              ))}
            </div>
            
            {/* User section in mobile */}
            <div className="border-t border-border pt-3">
              <div className="pl-3 pr-4 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider flex items-center">
                <User className="mr-2 h-4 w-4" />
                User
              </div>
              
              {/* Settings Link */}
              <Link
                href="/settings"
                className={clsx(
                pathname === '/settings'
                    ? 'bg-primary/10 border-primary text-primary'
                    : 'border-transparent text-muted-foreground hover:bg-accent hover:border-border hover:text-foreground',
                  'block pl-6 pr-4 py-2 border-l-4 text-base font-medium'
                )}
              >
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </Link>
              
              {/* Theme Options */}
              <div className="pl-6 pr-4 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                Theme
              </div>
              {themes.map((themeOption) => (
                <button
                  key={themeOption.value}
                  onClick={() => setTheme(themeOption.value)}
                  className={clsx(
                    theme === themeOption.value
                      ? 'bg-primary/10 border-primary text-primary'
                      : 'border-transparent text-muted-foreground hover:bg-accent hover:border-border hover:text-foreground',
                    'w-full text-left block pl-9 pr-4 py-2 border-l-4 text-base font-medium'
                  )}
                >
                  <themeOption.icon className="mr-2 h-4 w-4" />
                  {themeOption.label}
                  {theme === themeOption.value && (
                    <svg
                      className="inline ml-2 h-4 w-4 text-primary"
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
