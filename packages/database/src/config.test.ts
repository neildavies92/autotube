import { describe, expect, it } from 'vitest'
import { readDatabaseConfig } from './config.js'

const connectionString = 'postgresql://autotube:secret@localhost:55432/autotube'

describe('database environment boundary', () => {
  it('requires a PostgreSQL URL and rejects invalid pool limits without exposing input', () => {
    for (const environment of [
      {},
      { DATABASE_URL: 'https://user:top-secret@localhost/db' },
      { DATABASE_URL: 'postgresql://user:top-secret@localhost/' },
      { DATABASE_URL: connectionString, DATABASE_POOL_MAX: '0' },
      { DATABASE_URL: connectionString, DATABASE_CONNECTION_TIMEOUT_MS: 'NaN' },
      { DATABASE_URL: connectionString, DATABASE_IDLE_TIMEOUT_MS: '0' },
    ]) {
      expect(() => readDatabaseConfig(environment)).toThrow(/^Invalid database configuration:/)
      try { readDatabaseConfig(environment) } catch (error) {
        expect(String(error)).not.toContain('secret')
      }
    }
  })

  it('parses explicit pool settings and leaves the caller environment untouched', () => {
    const environment = {
      DATABASE_URL: connectionString,
      DATABASE_POOL_MAX: '2',
      DATABASE_CONNECTION_TIMEOUT_MS: '4000',
      DATABASE_IDLE_TIMEOUT_MS: '2000',
    }
    expect(readDatabaseConfig(environment)).toEqual({
      connectionString, max: 2, connectionTimeoutMillis: 4000, idleTimeoutMillis: 2000,
    })
    expect(environment.DATABASE_POOL_MAX).toBe('2')
  })
})
