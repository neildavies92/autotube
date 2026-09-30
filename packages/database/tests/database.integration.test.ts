import { sql } from 'drizzle-orm'
import { describe, expect, it } from 'vitest'
import { createDatabase } from '../src/client.js'
import { applyMigrations } from '../src/migrate.js'
import { withTestDatabase } from './test-database.js'

describe('PostgreSQL migrations and connection lifecycle', () => {
  it('migrates an empty database, preserves data on rerun and closes all owned connections', async () => {
    await withTestDatabase(async (config, admin, name) => {
      const connection = createDatabase(config)
      try {
        const before = await connection.db.execute<{ schema_name: string }>(sql`select schema_name from information_schema.schemata where schema_name = 'app'`)
        expect(before.rows).toEqual([])
        await applyMigrations(config)
        const schemas = await connection.db.execute(sql`select schema_name from information_schema.schemata where schema_name = 'app'`)
        expect(schemas.rowCount).toBe(1)
        const first = await connection.db.execute(sql`select hash, created_at from drizzle.__drizzle_migrations order by id`)
        expect(first.rowCount).toBe(1)
        // A synthetic fixture proves rerunning migrations does not reset existing data.
        await connection.db.execute(sql`create table app.integration_fixture (value text not null)`)
        await connection.db.execute(sql`insert into app.integration_fixture values ('preserve me')`)
        await applyMigrations(config)
        const second = await connection.db.execute(sql`select hash, created_at from drizzle.__drizzle_migrations order by id`)
        expect(second.rows).toEqual(first.rows)
        const fixture = await connection.db.execute(sql`select value from app.integration_fixture`)
        expect(fixture.rows).toEqual([{ value: 'preserve me' }])
      } finally {
        // Closing can be safely shared by signal handlers and error cleanup.
        await Promise.all([connection.close(), connection.close()])
      }
      const connections = await admin.query<{ count: string }>('select count(*) from pg_stat_activity where datname = $1', [name])
      expect(connections.rows[0]?.count).toBe('0')
      await expect(connection.db.execute(sql`select 1`)).rejects.toThrow()
    })
  })

  it('closes the migration pool after a failed migration', async () => {
    await withTestDatabase(async (config, admin, name) => {
      const fixture = createDatabase(config)
      try {
        // A conflicting relation makes journal creation fail before application migration.
        await fixture.db.execute(sql`create schema drizzle`)
        await fixture.db.execute(sql`create view drizzle.__drizzle_migrations as select 1 as id`)
      } finally {
        await fixture.close()
      }
      await expect(applyMigrations(config)).rejects.toThrow()
      const connections = await admin.query<{ count: string }>('select count(*) from pg_stat_activity where datname = $1', [name])
      expect(connections.rows[0]?.count).toBe('0')
    })
  })
})
