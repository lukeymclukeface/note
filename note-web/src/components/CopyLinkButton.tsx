'use client';

import { Button } from '@/components/ui/button';

interface CopyLinkButtonProps {
  className?: string;
  variant?: 'default' | 'outline' | 'secondary' | 'ghost' | 'link' | 'destructive';
}

export default function CopyLinkButton({ className, variant = 'outline' }: CopyLinkButtonProps) {
  const handleCopyLink = () => {
    navigator.clipboard.writeText(window.location.href);
    // TODO: Add toast notification here
  };

  return (
    <Button 
      variant={variant}
      className={className}
      onClick={handleCopyLink}
    >
      Copy Link
    </Button>
  );
}
