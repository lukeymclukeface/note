'use client';

import { useState, useRef } from 'react';
import { CheckCircle, XCircle, Music, FileText } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';

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

    // Validate file size (16GB limit)
    const maxSize = 16 * 1024 * 1024 * 1024; // 16GB
    if (file.size > maxSize) {
      setUploadResult({
        success: false,
        message: 'File size must be less than 16GB'
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
    <div className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Upload Files</h1>
          <p className="text-muted-foreground">
            Upload audio files or text documents to import into Note AI
          </p>
        </header>

        {/* Upload Area */}
        <div className="mb-8">
          <Card
            className={`relative border-2 border-dashed p-8 text-center transition-colors cursor-pointer ${
              isDragging
                ? 'border-primary bg-primary/5'
                : 'border-border hover:border-primary/50'
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
                  <div className="mx-auto h-16 w-16 text-muted-foreground">
                    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-lg">
                      <span className="font-medium text-primary">Click to upload</span> or drag and drop
                    </p>
                    <p className="text-sm text-muted-foreground mt-1">
                      Maximum file size: 16GB
                    </p>
                  </div>
                </>
              )}
            </div>
          </Card>
        </div>

        {/* Upload Result */}
        {uploadResult && (
          <Alert className={`mb-8 ${
            uploadResult.success 
              ? 'border-green-200 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-950 dark:text-green-200'
              : 'border-destructive bg-destructive/10 text-destructive'
          }`}>
            <div className="flex items-center">
              <div className="mr-3">
                {uploadResult.success ? (
                  <CheckCircle className="h-5 w-5" />
                ) : (
                  <XCircle className="h-5 w-5" />
                )}
              </div>
              <div>
                <AlertDescription className="font-medium">
                  {uploadResult.success ? 'Upload Successful!' : 'Upload Failed'}
                </AlertDescription>
                <AlertDescription className="text-sm mt-1">
                  {uploadResult.message}
                </AlertDescription>
                {uploadResult.success && uploadResult.filename && (
                  <AlertDescription className="text-sm mt-2">
                    File saved as: <code className="bg-secondary px-2 py-1 rounded text-xs">{uploadResult.filename}</code>
                  </AlertDescription>
                )}
              </div>
            </div>
          </Alert>
        )}

        {/* Supported File Types */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="text-lg">Supported File Types</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <h3 className="text-sm font-medium mb-3 flex items-center">
                  <Music className="mr-2 h-4 w-4" /> Audio Files
                </h3>
                <div className="space-y-2">
                  {supportedTypes.audio.map(type => (
                    <div key={type} className="flex items-center text-sm text-muted-foreground">
                      <span className="w-12 text-muted-foreground">{type}</span>
                      <span>Audio recording</span>
                    </div>
                  ))}
                </div>
              </div>
              
              <div>
                <h3 className="text-sm font-medium mb-3 flex items-center">
                  <FileText className="mr-2 h-4 w-4" /> Text Documents
                </h3>
                <div className="space-y-2">
                  {supportedTypes.text.map(type => (
                    <div key={type} className="flex items-center text-sm text-muted-foreground">
                      <span className="w-12 text-muted-foreground">{type}</span>
                      <span>Text document</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Next Steps */}
        <Card className="bg-primary/5 border-primary/20">
          <CardHeader>
            <CardTitle className="text-lg text-primary">What happens next?</CardTitle>
          </CardHeader>
          <CardContent className="text-sm space-y-2">
            <p>1. Your file will be saved to <code className="bg-secondary px-2 py-1 rounded">~/.noteai/import/</code></p>
            <p>2. Use the CLI to process the uploaded file:</p>
            <Card className="bg-secondary p-3 mt-2">
              <code className="text-foreground">note import ~/.noteai/import/your-file-name</code>
            </Card>
            <p>3. The processed content will appear in your notes, meetings, or interviews</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
