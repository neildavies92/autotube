import { randomUUID } from 'node:crypto'
import { Pool } from 'pg'
import { postgresUrl, readDatabaseConfig } from '../src/config.js'

// Never infer this connection from DATABASE_URL. Tests may only administer local postgres.
export function readTestAdminUrl(environment: NodeJS.ProcessEnv = process.env) {
  const parsed = postgresUrl.safeParse(environment.TEST_DATABASE_ADMIN_URL)
  if (!parsed.success) throw new Error('Set TEST_DATABASE_ADMIN_URL to the isolated local test server /postgres database')
  const url = new URL(parsed.data)
  if (!['127.0.0.1', 'localhost', '[::1]'].includes(url.hostname)
      || url.pathname !== '/postgres' || url.search !== '') {
    throw new Error('TEST_DATABASE_ADMIN_URL must target a loopback host and /postgres without query parameters')
  }
  return url
}

export async function withTestDatabase<T>(run: (config: ReturnType<typeof readDatabaseConfig>, admin: Pool, name: string) => Promise<T>) {
  const url = readTestAdminUrl()
  const admin = new Pool({ connectionString: url.toString(), max: 1, connectionTimeoutMillis: 5_000 })
  const name = `autotube_test_${randomUUID().replaceAll('-', '')}`
  let created = false
  try {
    // The identifier is generated here; it never contains caller input.
    await admin.query(`CREATE DATABASE "${name}"`)
    created = true
    url.pathname = `/${name}`
    return await run(readDatabaseConfig({ DATABASE_URL: url.toString() }), admin, name)
  } finally {
    try {
      // Only drop the unique database this invocation successfully created.
      if (created) await admin.query(`DROP DATABASE "${name}" WITH (FORCE)`)
    } finally {
      await admin.end()
    }
  }
}
