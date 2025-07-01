import { NextRequest, NextResponse } from 'next/server';
import fs from 'fs';
import path from 'path';
import os from 'os';

// Get the config file path (same as CLI: ~/.noteai/config.json)
function getConfigPath(): string {
  const homeDir = os.homedir();
  return path.join(homeDir, '.noteai', 'config.json');
}

export async function PUT(request: NextRequest) {
  try {
    console.log('PUT request received for config update');
    const configPath = getConfigPath();
    console.log('Config path:', configPath);
    
    const updatedConfig = await request.json();
    console.log('Received config data:', updatedConfig);

    // Validate that the config directory exists
    const configDir = path.dirname(configPath);
    if (!fs.existsSync(configDir)) {
      console.log('Creating config directory:', configDir);
      fs.mkdirSync(configDir, { recursive: true });
    }

    // Write the updated configuration
    console.log('Writing config to file');
    fs.writeFileSync(configPath, JSON.stringify(updatedConfig, null, 2));
    console.log('Config written successfully');

    return NextResponse.json({ success: true, message: 'Configuration updated successfully' });
  } catch (error) {
    console.error('Error updating config:', error);
    console.error('Error stack:', error instanceof Error ? error.stack : 'Unknown error');
    return NextResponse.json(
      { success: false, error: `Failed to update configuration: ${error instanceof Error ? error.message : 'Unknown error'}` },
      { status: 500 }
    );
  }
}

export async function GET() {
  try {
    const configPath = getConfigPath();
    
    if (!fs.existsSync(configPath)) {
      return NextResponse.json(
        { success: false, error: 'Configuration file not found' },
        { status: 404 }
      );
    }

    const configData = fs.readFileSync(configPath, 'utf8');
    const config = JSON.parse(configData);

    return NextResponse.json({ success: true, config });
  } catch (error) {
    console.error('Error reading config:', error);
    return NextResponse.json(
      { success: false, error: 'Failed to read configuration' },
      { status: 500 }
    );
  }
}
