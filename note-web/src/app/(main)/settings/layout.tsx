'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Settings, Bot, Database, FileText } from 'lucide-react';

const settingsNavItems = [
  { label: 'General', href: '/settings', icon: Settings },
  { label: 'AI Settings', href: '/settings/ai', icon: Bot },
  { label: 'Database', href: '/settings/database', icon: Database },
  { label: 'Raw Config', href: '/settings/raw', icon: FileText },
];

export default function SettingsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-background">
      <div className="flex">
        {/* Left Sidebar */}
        <div className="w-64 bg-card border-r border-border min-h-screen">
          <div className="p-6">
            <h1 className="text-2xl font-bold mb-6">Settings</h1>
            <nav className="space-y-1">
              {settingsNavItems.map((item) => {
                const isActive = pathname === item.href;
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                      isActive
                        ? 'bg-primary/10 text-primary'
                        : 'text-muted-foreground hover:bg-accent hover:text-foreground'
                    }`}
                  >
                    <item.icon className="mr-3 h-5 w-5" />
                    {item.label}
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1">
          <div className="max-w-4xl p-6">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}
