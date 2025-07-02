'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { clsx } from 'clsx';
import { Home, FileText, Users, Briefcase, Calendar, Import, User, Settings, Sun, Moon, Monitor, Menu, Check } from 'lucide-react';
import NavbarRecorder from './NavbarRecorder';
import { UserDropdownMenu } from './UserDropdownMenu';
import { useTheme } from '@/providers/ThemeProvider';

const navigation = [
  { name: 'Dashboard', href: '/', icon: Home },
  { name: 'Notes', href: '/notes', icon: FileText },
  { name: 'Meetings', href: '/meetings', icon: Users },
  { name: 'Interviews', href: '/interviews', icon: Briefcase },
  { name: 'Calendar', href: '/calendar', icon: Calendar },
  { name: 'Import', href: '/import', icon: Import },
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
              {navigation.map((item) => {
                const isActive = item.href === '/import' 
                  ? pathname?.startsWith('/import') || pathname === '/recordings' || pathname === '/upload'
                  : pathname === item.href;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={clsx(
                      isActive
                        ? 'border-primary text-foreground'
                        : 'border-transparent text-muted-foreground hover:border-border hover:text-foreground',
                      'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-colors'
                    )}
                  >
                    <item.icon className="mr-2 h-4 w-4" />
                    {item.name}
                  </Link>
                );
              })}
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
            <div className="sm:hidden flex items-center space-x-2">
              <UserDropdownMenu />
              <button
                type="button"
                className="inline-flex items-center justify-center p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent"
                aria-controls="mobile-menu"
                aria-expanded="false"
              >
                <span className="sr-only">Open main menu</span>
                <Menu className="h-6 w-6" />
              </button>
            </div>
          </div>
        </div>
        
        {/* Mobile menu */}
        <div className="sm:hidden" id="mobile-menu">
          <div className="pt-2 pb-3 space-y-1">
            {navigation.map((item) => {
              const isActive = item.href === '/import' 
                ? pathname?.startsWith('/import') || pathname === '/recordings' || pathname === '/upload'
                : pathname === item.href;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={clsx(
                    isActive
                      ? 'bg-primary/10 border-primary text-primary'
                      : 'border-transparent text-muted-foreground hover:bg-accent hover:border-border hover:text-foreground',
                    'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
                  )}
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.name}
                </Link>
              );
            })}
            
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
                    <Check className="inline ml-2 h-4 w-4 text-primary" />
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
