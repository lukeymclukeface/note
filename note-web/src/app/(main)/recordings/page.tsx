'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function RecordingsRedirect() {
  const router = useRouter();

  useEffect(() => {
    router.replace('/import');
  }, [router]);

  return null;
}
