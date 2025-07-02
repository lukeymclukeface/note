'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function UploadRedirect() {
  const router = useRouter();

  useEffect(() => {
    router.replace('/import/upload');
  }, [router]);

  return null;
}
