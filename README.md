# AutoTube v2

AutoTube will discover evidence-backed evergreen YouTube opportunities, maintain a ranked topic backlog, produce original videos, publish on a controlled schedule and learn from performance. The first MVP ends at an inspectable topic backlog.

This repository provides the **development and persistence foundations**. Domain workflows are not implemented.

## Architecture

- `apps/api`: Fastify process with a credential-free `GET /health` endpoint and Zod-validated startup settings.
- `apps/worker`: separate, idle Node process; no jobs or polling yet.
- `apps/web`: React/Vite foundation page, with a development proxy to the API.
- `packages/database`: typed Drizzle connection and migration tooling; no application/domain tables yet.
- Strict TypeScript, pnpm, ESLint and Vitest across one workspace. Keep shared modules bounded to actual infrastructure needs.

PostgreSQL/Drizzle is provided by SWA-47. The API/operational UI is tracked in SWA-48, and Graphile Worker in SWA-49. Remotion/FFmpeg belongs to video production. Grafana is deferred to performance reporting; the React app will focus on workflow controls and evidence inspection.

## Local setup

Install Node **24.21.0 LTS** (also pinned in `.nvmrc`) and pnpm **10.34.6**. With nvm available:

```sh
nvm install
nvm use
npm install --global pnpm@10.34.6
pnpm install --frozen-lockfile
pnpm dev
```

`pnpm dev` starts all three apps. Open http://127.0.0.1:5173. The API listens on http://127.0.0.1:8080; the worker prints that it is idle. Stop with Ctrl+C. No `.env` file, database, Docker, Go or provider credentials are required.

To run one app in a separate terminal:

```sh
pnpm dev:api
pnpm dev:worker
pnpm dev:web
```

The web development server forwards `/api/*` to `http://127.0.0.1:8080/*`, removing `/api`. The page's health link therefore requests `/api/health`. Ports are strict: startup fails if occupied. Optional shell variables `AUTOTUBE_HOST` and `AUTOTUBE_PORT` override the API bind address only; `.env.example` documents defaults and files are not loaded automatically. If changing the API port, update the development proxy target in `apps/web/vite.config.ts` too. Only bind beyond loopback deliberately; authentication is not implemented.

## Verification

Run exactly the same commands as CI from the repository root:

```sh
pnpm install --frozen-lockfile
pnpm check
```

`check` runs these commands in sequence, stopping on failure:

```sh
pnpm lint
pnpm typecheck
pnpm test
pnpm build
pnpm smoke
```

Tests use Fastify injection and synthetic environment settings, without live providers. The smoke check starts compiled API/worker processes, checks startup failures and graceful shutdown, starts Vite to check the development proxy, then serves the built web assets. Leave ports **8080, 5173 and 4173** free before running it. It cleans up its own processes. It checks HTTP/assets, not browser rendering.

After `pnpm build`, run built entrypoints separately:

```sh
pnpm --filter @autotube/api start
pnpm --filter @autotube/worker start
pnpm --filter @autotube/web preview
```

Vite preview serves the built web at http://127.0.0.1:4173 and is a local asset preview, not a production deployment. The `/api` proxy is development-only.

## Local PostgreSQL

Docker Compose provides a PostgreSQL 17.11 instance on loopback port 55432, with a persistent v2 volume and healthcheck. The API health endpoint remains independent of the database. To configure and migrate the local database:

```sh
cp .env.example .env # only if .env does not already exist
pnpm db:up
pnpm db:migrate
```

`pnpm db:down` stops the development service without removing its volume. `pnpm db:generate` generates future Drizzle migrations for review; never rewrite an applied migration. No general reset command is provided. The initial migration creates an application schema and Drizzle's migration history, with no speculative product tables.

Use a separate ephemeral PostgreSQL instance for integration tests:

```sh
pnpm db:test:up
TEST_DATABASE_ADMIN_URL=postgresql://autotube_test:autotube_test@127.0.0.1:55433/postgres pnpm test:integration
pnpm db:test:down
```

Every test run creates a new uniquely named test database and drops only the database it created. Rerunning the test command is the reset-for-tests workflow; it never falls back to `DATABASE_URL`. CI runs these same integration tests against its own PostgreSQL service. See [database setup and safety boundaries](docs/swa-47-postgres.md) for exact configuration, migration and isolated-test behaviour.

## Contributing

Read [AGENTS.md](AGENTS.md), the current [Linear issue](https://linear.app/swale-solutions/issue/SWA-46/overhaul-existing-repository-into-a-typescript-workspace) and [AutoTube v2 decisions](https://linear.app/swale-solutions/project/autotube-v2-69e62d254724). Each issue gets one bounded feature branch and PR, Conventional Commits, and exact verification results. Never push directly to `main`, rewrite shared history, commit secrets or reset existing data.

See [migration notes](docs/swa-46-migration.md) for the starting commit, retained behaviour, v1 recovery and the next work. V1 code and its fixtures remain available in Git history.
