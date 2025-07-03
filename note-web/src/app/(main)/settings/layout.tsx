'use client';

import { clsx } from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Settings, Bot, Database, FileText } from 'lucide-react';

const navigation = [
  { name: 'General', href: '/settings', icon: Settings },
  { name: 'AI Settings', href: '/settings/ai', icon: Bot },
  { name: 'Database', href: '/settings/database', icon: Database },
  { name: 'Raw Config', href: '/settings/raw', icon: FileText },
];

const SettingsLayout = ({ children }: { children: React.ReactNode }) => {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Settings</h1>
          <p className="text-muted-foreground">
            Configure your note-taking and AI processing preferences
          </p>
        </header>
        
        <div className="flex gap-8">
          {/* Sidebar */}
          <aside className="w-64 flex-shrink-0">
            <nav className="space-y-1">
              {navigation.map((item) => {
                const isActive = pathname === item.href;
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={clsx(
                      isActive
                        ? 'bg-primary/10 text-primary border-primary'
                        : 'text-muted-foreground hover:text-foreground hover:bg-accent border-transparent',
                      'flex items-center px-3 py-2 text-sm font-medium rounded-md border-l-4 transition-colors'
                    )}
                  >
                    <item.icon className="mr-3 h-5 w-5" />
                    {item.name}
                  </Link>
                );
              })}
            </nav>
          </aside>
          
          {/* Main content */}
          <main className="flex-1 min-w-0">
            {children}
          </main>
        </div>
      </div>
    </div>
  );
};

export default SettingsLayout;
