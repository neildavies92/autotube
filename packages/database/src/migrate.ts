import { fileURLToPath } from 'node:url'
import { migrate } from 'drizzle-orm/node-postgres/migrator'
import { createDatabase } from './client.js'
import type { DatabaseConfig } from './config.js'

export async function applyMigrations(config: DatabaseConfig) {
  const connection = createDatabase({ ...config, max: 1 })
  try {
    await migrate(connection.db, {
      migrationsFolder: fileURLToPath(new URL('../migrations', import.meta.url)),
      migrationsSchema: 'drizzle',
      migrationsTable: '__drizzle_migrations',
    })
  } finally {
    await connection.close()
  }
}
