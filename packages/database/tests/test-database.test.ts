import { describe, expect, it } from 'vitest'
import { readTestAdminUrl } from './test-database.js'

describe('test database safety boundary', () => {
  it('rejects missing, remote, application-database and connection-override URLs', () => {
    for (const value of [
      undefined,
      'postgresql://user:secret@remote.example/postgres',
      'postgresql://user:secret@localhost/autotube',
      'postgresql://user:secret@localhost/postgres?host=remote.example',
    ]) {
      expect(() => readTestAdminUrl({ TEST_DATABASE_ADMIN_URL: value })).toThrow()
    }
    expect(() => readTestAdminUrl({ DATABASE_URL: 'postgresql://user:secret@localhost/postgres' })).toThrow()
  })

  it('accepts the explicit loopback admin URL used by local tests and CI', () => {
    expect(readTestAdminUrl({
      TEST_DATABASE_ADMIN_URL: 'postgresql://autotube_test:secret@127.0.0.1:55433/postgres',
    }).pathname).toBe('/postgres')
  })
})
