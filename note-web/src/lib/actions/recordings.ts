'use server';

import { getAllRecordings } from '@/lib/database';
import type { Recording } from '@/lib/database';

export async function loadRecordings(): Promise<Recording[]> {
  try {
    return getAllRecordings();
  } catch (error) {
    console.error('Error loading recordings:', error);
    return [];
  }
}
