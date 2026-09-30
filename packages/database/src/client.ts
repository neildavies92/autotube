import { drizzle } from 'drizzle-orm/node-postgres'
import { Pool } from 'pg'
import type { DatabaseConfig } from './config.js'
import * as schema from './schema.js'

export function createDatabase(config: DatabaseConfig) {
  const pool = new Pool(config)
  // pg emits errors for idle connections; never print the connection or driver error.
  pool.on('error', () => console.error('Database idle connection failed'))
  const db = drizzle(pool, { schema })
  let closing: Promise<void> | undefined
  return {
    db,
    close: () => (closing ??= pool.end()),
  }
}

export type Database = ReturnType<typeof createDatabase>['db']
