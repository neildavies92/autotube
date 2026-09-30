# SWA-47: PostgreSQL and Drizzle foundation

This issue starts from SWA-46 on `ff99a9f`. PostgreSQL 17.11 runs locally, with Drizzle ORM 0.45.3 and Drizzle Kit 0.31.11. The `@autotube/database` package owns validated connection configuration, typed database access, migration tooling and migration history. It is shared because migration commands and the upcoming API/worker persistence consumers need the same schema and connection lifecycle.

The initial generated migration creates only the `app` schema. Drizzle owns its `drizzle.__drizzle_migrations` journal. There are no speculative domain tables. SWA-49 will introduce job/run records and Graphile Worker in its own schema after this foundation merges. API health and the current worker do not yet depend on PostgreSQL.

## Local setup

Install the pinned Node/pnpm versions from the root README and Docker with Compose v2. From the repository root:

```sh
pnpm install --frozen-lockfile
cp .env.example .env
pnpm db:up
pnpm db:migrate
pnpm db:migrate
```

Copy the example only for a new local environment; preserve an existing `.env`. Database migration and integration commands load the root `.env` through Node's `--env-file-if-exists`; explicitly exported environment variables take precedence. Do not commit real credentials. `DATABASE_URL` is required for migration/connection use. Its PostgreSQL protocol, hostname, username and database are validated, as are bounded pool size and timeouts. Validation errors identify field names only; the CLI reports safe failure context without connection strings or raw driver errors.

The local database listens only on `127.0.0.1:55432`. `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` and `POSTGRES_PORT` configure Compose; keep `DATABASE_URL` in step with them. These defaults are for local development. PostgreSQL applies initialization settings only when creating a fresh data directory: changing `.env` does not change an existing volume's credentials.

`infra/compose.yaml` uses the distinct `autotube-v2` Compose project, a persistent named `postgres-data` volume, and a `pg_isready` healthcheck. `pnpm db:up` waits for health. `pnpm db:down` stops/removes its containers while retaining the volume. Do not run the retired v1 Compose/SQL setup or remove its data. There is deliberately no general-purpose reset command and no `down --volumes` in the documented development path; resetting non-test data requires explicit user action.

## Schema changes and migrations

Update `packages/database/src/schema.ts` in the issue that needs a table, then run:

```sh
pnpm db:generate
pnpm db:migrate
```

Generation is offline and does not require database credentials. Review and commit generated SQL, the journal and snapshots together. Applied migrations are immutable; add a new forward migration instead. Do not use Drizzle push to bypass reviewed migration files. The migration runner resolves files relative to its module, applies pending migrations and always closes its pool, including failures. Rerunning skips journaled migrations and preserves data. Run a single migration process at a time; concurrent migration deployment orchestration is outside this local foundation.

Consumers import `createDatabase` and `readDatabaseConfig` from `@autotube/database`. Construct once per process, use the typed `db` handle, and await `close()` in shutdown and error cleanup. Multiple calls share the same close promise. Connection, idle timeout and pool bounds prevent unbounded waits or pools. No raw database or credential errors should be forwarded to HTTP responses or logs.

## Isolated integration checks and reset-for-tests

The test server is separate from development, with loopback port 55433 and ephemeral `tmpfs` storage. No persistent volume is mounted. From the repository root:

```sh
pnpm db:test:up
pnpm test:integration
pnpm db:test:down
```

Set `TEST_DATABASE_ADMIN_URL=postgresql://autotube_test:autotube_test@127.0.0.1:55433/postgres` in `.env` (included in the example). It is mandatory and never falls back to `DATABASE_URL`. The test runner permits only loopback hosts, the `/postgres` maintenance database and no query parameters. The admin needs permission to create/drop databases; the dedicated local test server supplies it. CI uses the same tests against its disposable PostgreSQL service on port 5432.

Each test creates a fresh random `autotube_test_<uuid>` database, migrates and inspects only that database, then drops only the database it successfully created. Tests close application and migration connections and assert zero remaining sessions before cleanup. They prove clean migration application, unchanged journal/data on rerun, idempotent client shutdown, and cleanup after migration failure. No paid or external provider calls occur.

To reset **only test infrastructure**, stop and recreate the dedicated ephemeral server:

```sh
pnpm db:test:down
pnpm db:test:up
pnpm test:integration
```

Every test run already creates a fresh database. Reset is useful after an interrupted test process; deleting the dedicated test container also discards its tmpfs contents. Never point tests at development or production, and never substitute the development Compose file in these reset commands.

## Required verification

```sh
pnpm lint
pnpm typecheck
pnpm test
pnpm build
pnpm smoke
pnpm db:test:up
pnpm test:integration
pnpm db:test:down
```

The normal unit suite requires no database. Integration tests fail with a configuration error if the dedicated test URL is missing. `pnpm check` runs the normal foundation checks; integration is an explicit command and a required CI step. See the PR for the commands actually executed and their results.
