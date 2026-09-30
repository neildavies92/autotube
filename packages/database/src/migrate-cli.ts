import { readDatabaseConfig } from './config.js'
import { applyMigrations } from './migrate.js'

try {
  const config = readDatabaseConfig()
  await applyMigrations(config)
  console.log('Database migrations applied')
} catch {
  // Driver errors can embed connection URLs or SQL containing sensitive data.
  console.error('Database migration failed. Check DATABASE_URL, database availability and migration files.')
  process.exitCode = 1
}
