import { NextResponse } from 'next/server';
import { configService } from '@/lib/config';
import sqlite3 from 'sqlite3';
import path from 'path';

interface TableInfo {
  name: string;
  columns: Array<{
    name: string;
    type: string;
    nullable: boolean;
    primaryKey: boolean;
  }>;
  rowCount: number;
}

interface DatabaseValidation {
  connected: boolean;
  tables: TableInfo[];
  errors: string[];
  version?: string;
}

interface SqliteVersionRow {
  version: string;
}

interface SqliteTableRow {
  name: string;
}

interface SqliteColumnInfo {
  name: string;
  type: string;
  notnull: number;
  pk: number;
}

interface SqliteCountRow {
  count: number;
}

export async function GET() {
  try {
    const config = await configService.getConfig();
    
    if (!config || !config.database_path) {
      return NextResponse.json({
        success: false,
        error: 'Database path not configured'
      }, { status: 400 });
    }

    const validation: DatabaseValidation = {
      connected: false,
      tables: [],
      errors: []
    };

    // Resolve the database path
    const dbPath = config.database_path.startsWith('~') 
      ? path.join(process.env.HOME || '', config.database_path.slice(1))
      : config.database_path;

    return new Promise<NextResponse>((resolve) => {
      const db = new sqlite3.Database(dbPath, sqlite3.OPEN_READONLY, async (err) => {
        if (err) {
          validation.errors.push(`Failed to connect to database: ${err.message}`);
          resolve(NextResponse.json({
            success: true,
            validation
          }));
          return;
        }

        validation.connected = true;

        try {
          // Get SQLite version
          db.get('SELECT sqlite_version() as version', (err, row: SqliteVersionRow | undefined) => {
            if (!err && row) {
              validation.version = row.version;
            }
          });

          // Get list of tables
          db.all(
            "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'",
            async (err, tables: SqliteTableRow[]) => {
              if (err) {
                validation.errors.push(`Failed to get tables: ${err.message}`);
                db.close();
                resolve(NextResponse.json({
                  success: true,
                  validation
                }));
                return;
              }

              let processedTables = 0;
              const totalTables = tables.length;

              if (totalTables === 0) {
                db.close();
                resolve(NextResponse.json({
                  success: true,
                  validation
                }));
                return;
              }

              for (const table of tables) {
                const tableName = table.name;
                
                // Get table info
                db.all(`PRAGMA table_info(${tableName})`, (err, columns: SqliteColumnInfo[]) => {
                  if (err) {
                    validation.errors.push(`Failed to get info for table ${tableName}: ${err.message}`);
                  } else {
                    // Get row count
                    db.get(`SELECT COUNT(*) as count FROM ${tableName}`, (err, countRow: SqliteCountRow | undefined) => {
                      const rowCount = err ? 0 : (countRow?.count || 0);
                      
                      if (!err) {
                        const tableInfo: TableInfo = {
                          name: tableName,
                          columns: columns.map(col => ({
                            name: col.name,
                            type: col.type,
                            nullable: !col.notnull,
                            primaryKey: !!col.pk
                          })),
                          rowCount
                        };
                        validation.tables.push(tableInfo);
                      } else {
                        validation.errors.push(`Failed to get row count for table ${tableName}: ${err.message}`);
                      }

                      processedTables++;
                      if (processedTables === totalTables) {
                        // Sort tables by name
                        validation.tables.sort((a, b) => a.name.localeCompare(b.name));
                        
                        db.close();
                        resolve(NextResponse.json({
                          success: true,
                          validation
                        }));
                      }
                    });
                  }
                  
                  if (err) {
                    processedTables++;
                    if (processedTables === totalTables) {
                      validation.tables.sort((a, b) => a.name.localeCompare(b.name));
                      db.close();
                      resolve(NextResponse.json({
                        success: true,
                        validation
                      }));
                    }
                  }
                });
              }
            }
          );
        } catch (error) {
          validation.errors.push(`Database validation error: ${error instanceof Error ? error.message : 'Unknown error'}`);
          db.close();
          resolve(NextResponse.json({
            success: true,
            validation
          }));
        }
      });
    });
  } catch (error) {
    console.error('Failed to validate database:', error);
    
    return NextResponse.json({
      success: false,
      error: error instanceof Error ? error.message : 'Failed to validate database'
    }, { status: 500 });
  }
}
