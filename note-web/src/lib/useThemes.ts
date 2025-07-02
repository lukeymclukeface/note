'use client';

import { useTheme } from '@/providers/ThemeProvider';

export const themes = [
  { value: 'light', label: 'Light', icon: '☀️' },
  { value: 'dark', label: 'Dark', icon: '🌙' },
  { value: 'system', label: 'System', icon: '💻' },
] as const;

export function useThemes() {
  const { theme, setTheme } = useTheme();
  return { theme, setTheme, themes };
}
