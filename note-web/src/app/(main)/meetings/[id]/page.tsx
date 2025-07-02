import { redirect } from 'next/navigation';

interface MeetingPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function MeetingPage({ params }: MeetingPageProps) {
  const { id } = await params;
  redirect(`/notes/meetings/${id}`);
}
