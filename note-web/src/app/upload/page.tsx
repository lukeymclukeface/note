'use client';

import { useState, useRef } from 'react';

interface UploadResult {
  success: boolean;
  message: string;
  filename?: string;
}

export default function UploadPage() {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadResult, setUploadResult] = useState<UploadResult | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const supportedTypes = {
    audio: ['.mp3', '.wav', '.m4a', '.aac', '.flac', '.ogg'],
    text: ['.txt', '.md', '.doc', '.docx', '.pdf']
  };

  const allSupportedTypes = [...supportedTypes.audio, ...supportedTypes.text];

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const files = Array.from(e.dataTransfer.files);
    if (files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleFileUpload = async (file: File) => {
    // Validate file type
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase();
    if (!allSupportedTypes.includes(fileExtension)) {
      setUploadResult({
        success: false,
        message: `Unsupported file type. Supported types: ${allSupportedTypes.join(', ')}`
      });
      return;
    }

    // Validate file size (50MB limit)
    const maxSize = 50 * 1024 * 1024; // 50MB
    if (file.size > maxSize) {
      setUploadResult({
        success: false,
        message: 'File size must be less than 50MB'
      });
      return;
    }

    setIsUploading(true);
    setUploadResult(null);

    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });

      const result: UploadResult = await response.json();
      setUploadResult(result);

    } catch {
      setUploadResult({
        success: false,
        message: 'Upload failed. Please try again.'
      });
    } finally {
      setIsUploading(false);
      // Reset file input
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };


  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">Upload Files</h1>
          <p className="text-gray-600 dark:text-gray-300">
            Upload audio files or text documents to import into Note AI
          </p>
        </header>

        {/* Upload Area */}
        <div className="mb-8">
          <div
            className={`relative border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
              isDragging
                ? 'border-blue-400 bg-blue-50 dark:bg-blue-950 dark:border-blue-500'
                : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
            } ${
              isUploading ? 'opacity-50 pointer-events-none' : ''
            }`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <input
              ref={fileInputRef}
              type="file"
              className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
              accept={allSupportedTypes.join(',')}
              onChange={handleFileSelect}
              disabled={isUploading}
            />
            
            <div className="space-y-4">
              {isUploading ? (
                <>
                  <div className="mx-auto h-16 w-16 animate-spin">
                    <svg className="h-16 w-16 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>
                  <p className="text-lg text-gray-600 dark:text-gray-300">Uploading...</p>
                </>
              ) : (
                <>
                  <div className="mx-auto h-16 w-16 text-gray-400 dark:text-gray-500">
                    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-lg text-gray-600 dark:text-gray-300">
                      <span className="font-medium text-blue-600 dark:text-blue-400">Click to upload</span> or drag and drop
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                      Maximum file size: 50MB
                    </p>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Upload Result */}
        {uploadResult && (
          <div className={`mb-8 p-4 rounded-lg ${
            uploadResult.success 
              ? 'bg-green-50 dark:bg-green-900 text-green-800 dark:text-green-200 border border-green-200 dark:border-green-700'
              : 'bg-red-50 dark:bg-red-900 text-red-800 dark:text-red-200 border border-red-200 dark:border-red-700'
          }`}>
            <div className="flex items-center">
              <div className="mr-3">
                {uploadResult.success ? '✅' : '❌'}
              </div>
              <div>
                <p className="font-medium">{uploadResult.success ? 'Upload Successful!' : 'Upload Failed'}</p>
                <p className="text-sm mt-1">{uploadResult.message}</p>
                {uploadResult.success && uploadResult.filename && (
                  <p className="text-sm mt-2">
                    File saved as: <code className="bg-white dark:bg-gray-800 px-2 py-1 rounded text-xs">{uploadResult.filename}</code>
                  </p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Supported File Types */}
        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Supported File Types</h2>
          
          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3 flex items-center">
                🎵 Audio Files
              </h3>
              <div className="space-y-2">
                {supportedTypes.audio.map(type => (
                  <div key={type} className="flex items-center text-sm text-gray-600 dark:text-gray-400">
                    <span className="w-12 text-gray-500 dark:text-gray-500">{type}</span>
                    <span>Audio recording</span>
                  </div>
                ))}
              </div>
            </div>
            
            <div>
              <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3 flex items-center">
                📄 Text Documents
              </h3>
              <div className="space-y-2">
                {supportedTypes.text.map(type => (
                  <div key={type} className="flex items-center text-sm text-gray-600 dark:text-gray-400">
                    <span className="w-12 text-gray-500 dark:text-gray-500">{type}</span>
                    <span>Text document</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Next Steps */}
        <div className="mt-8 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
          <h2 className="text-lg font-semibold text-blue-900 dark:text-blue-100 mb-2">What happens next?</h2>
          <div className="text-sm text-blue-800 dark:text-blue-200 space-y-2">
            <p>1. Your file will be saved to <code className="bg-blue-100 dark:bg-blue-900 px-2 py-1 rounded">~/.noteai/import/</code></p>
            <p>2. Use the CLI to process the uploaded file:</p>
            <div className="bg-blue-100 dark:bg-blue-900 p-3 rounded mt-2">
              <code className="text-blue-900 dark:text-blue-100">note import ~/.noteai/import/your-file-name</code>
            </div>
            <p>3. The processed content will appear in your notes, meetings, or interviews</p>
          </div>
        </div>
      </div>
    </div>
  );
}
