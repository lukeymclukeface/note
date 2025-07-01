import { NextResponse } from 'next/server';

interface OpenAIValidationResult {
  valid: boolean;
  error?: string;
  details?: {
    organizationId?: string;
    models?: string[];
  };
}

export async function POST(request: Request) {
  try {
    const { apiKey } = await request.json();

    if (!apiKey) {
      return NextResponse.json({
        success: false,
        error: 'API key is required'
      }, { status: 400 });
    }

    // Test the API key by making a request to list models
    const response = await fetch('https://api.openai.com/v1/models', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      let errorMessage = 'Invalid API key';
      
      try {
        const errorData = await response.json();
        if (errorData.error?.message) {
          errorMessage = errorData.error.message;
        }
      } catch {
        // If we can't parse the error, use status text
        errorMessage = `HTTP ${response.status}: ${response.statusText}`;
      }

      const result: OpenAIValidationResult = {
        valid: false,
        error: errorMessage
      };

      return NextResponse.json({
        success: true,
        validation: result
      });
    }

    // Parse the successful response
    const data = await response.json();
    
    // Extract useful information
    const models = data.data?.slice(0, 10).map((model: { id: string }) => model.id) || [];
    
    // Try to get organization info (this might fail but that's ok)
    let organizationId: string | undefined;
    try {
      const orgResponse = await fetch('https://api.openai.com/v1/organizations', {
        headers: {
          'Authorization': `Bearer ${apiKey}`,
          'Content-Type': 'application/json',
        },
      });
      
      if (orgResponse.ok) {
        const orgData = await orgResponse.json();
        organizationId = orgData.data?.[0]?.id;
      }
    } catch {
      // Ignore errors getting organization info
    }

    const result: OpenAIValidationResult = {
      valid: true,
      details: {
        organizationId,
        models,
      }
    };

    return NextResponse.json({
      success: true,
      validation: result
    });

  } catch (error) {
    console.error('OpenAI validation error:', error);
    
    return NextResponse.json({
      success: false,
      error: error instanceof Error ? error.message : 'Failed to validate API key'
    }, { status: 500 });
  }
}
