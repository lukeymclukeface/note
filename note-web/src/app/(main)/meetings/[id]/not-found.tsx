import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🤝</div>
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Meeting Not Found</h1>
        <p className="text-gray-600 mb-8">
          The meeting you&apos;re looking for doesn&apos;t exist or may have been removed.
        </p>
        <div className="space-x-4">
          <Link
            href="/meetings"
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 transition-colors"
          >
            ← Back to Meetings
          </Link>
          <Link
            href="/"
            className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors"
          >
            Go Home
          </Link>
        </div>
      </div>
    </div>
  );
}
