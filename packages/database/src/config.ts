import { z } from 'zod'

export const postgresUrl = z.string().min(1).refine((value) => {
  try {
    const url = new URL(value)
    return ['postgres:', 'postgresql:'].includes(url.protocol)
      && url.hostname.length > 0
      && url.username.length > 0
      && url.pathname.length > 1
      && url.hash === ''
  } catch {
    return false
  }
}, 'Must be a PostgreSQL connection URL with host, username and database')

const environmentSchema = z.object({
  DATABASE_URL: postgresUrl,
  DATABASE_POOL_MAX: z.coerce.number().int().min(1).max(20).default(5),
  DATABASE_CONNECTION_TIMEOUT_MS: z.coerce.number().int().min(1).max(60_000).default(5_000),
  DATABASE_IDLE_TIMEOUT_MS: z.coerce.number().int().min(1).max(300_000).default(10_000),
})

export function readDatabaseConfig(environment: NodeJS.ProcessEnv = process.env) {
  const result = environmentSchema.safeParse(environment)
  if (!result.success) {
    // Zod input/errors can contain credentials. Report field names only.
    const fields = [...new Set(result.error.issues.map((issue) => issue.path.join('.')))]
    throw new Error(`Invalid database configuration: ${fields.join(', ')}`)
  }
  return {
    connectionString: result.data.DATABASE_URL,
    max: result.data.DATABASE_POOL_MAX,
    connectionTimeoutMillis: result.data.DATABASE_CONNECTION_TIMEOUT_MS,
    idleTimeoutMillis: result.data.DATABASE_IDLE_TIMEOUT_MS,
  }
}

export type DatabaseConfig = ReturnType<typeof readDatabaseConfig>
