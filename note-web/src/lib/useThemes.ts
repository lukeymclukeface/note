'use client';

import { useTheme } from '@/providers/ThemeProvider';
import { Sun, Moon, Monitor } from 'lucide-react';

export const themes = [
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
  { value: 'system', label: 'System', icon: Monitor },
] as const;

export function useThemes() {
  const { theme, setTheme } = useTheme();
  return { theme, setTheme, themes };
}
