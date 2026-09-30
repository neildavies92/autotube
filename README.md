# AutoTube v2

AutoTube will discover evidence-backed evergreen YouTube opportunities, maintain a ranked topic backlog, produce original videos, publish on a controlled schedule and learn from performance. The first MVP ends at an inspectable topic backlog.

This repository currently provides the **development foundation only** (SWA-46).

## Architecture

- `apps/api`: Fastify health boundary with Zod validation, consistent errors, request correlation IDs and sanitized structured logs.
- `apps/worker`: separate, idle Node process; no jobs or polling yet.
- `apps/web`: React/Vite operational navigation and API availability checks, with loading, failure and retry states.
- Strict TypeScript, pnpm, ESLint and Vitest across one workspace. Shared packages are added only when multiple consumers need them; none are needed yet.

PostgreSQL/Drizzle comes in SWA-47, the API/operational UI in SWA-48, and Graphile Worker in SWA-49. Remotion/FFmpeg belongs to video production. Grafana is deferred to performance reporting; the React app will focus on workflow controls and evidence inspection.

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

The web development server forwards `/api/*` to `http://127.0.0.1:8080/*`, removing `/api`. The operational shell requests `/api/health` and validates the response before reporting availability. Ports are strict: startup fails if occupied. Optional shell variables `AUTOTUBE_HOST` and `AUTOTUBE_PORT` override the API bind address only; `.env.example` documents defaults and files are not loaded automatically. If changing the API port, update the development proxy target in `apps/web/vite.config.ts` too. Only bind beyond loopback deliberately; authentication is not implemented.

## API and operational shell

The shell provides Overview and API connection views. It checks API availability on mount and on demand, with a five-second timeout, loading state and retry after failure. Availability describes the last successful check of the API process; it does not imply database, worker or provider readiness. Workflow controls are introduced with their owning issues.

`GET /health` returns the retained health response and `x-request-id`. Callers may supply a UUID `x-request-id`; invalid IDs and unexpected health query parameters receive a consistent error envelope. Unknown routes and internal errors use the same envelope without exposing raw exceptions or request values. Logs record correlation, status and duration with sensitive data omitted/redacted. Optional `AUTOTUBE_LOG_LEVEL` configures logging; startup configuration errors expose field names only.

See [API/UI behaviour and verification](docs/swa-48-api-ui.md) for the full contract, shutdown and tests. API settings are shell environment variables, not automatically loaded from `.env`.

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

Vite preview serves the built web at http://127.0.0.1:4173 and is a local asset preview, not a production deployment. Development and preview share the `/api` proxy to the local API.

## Contributing

Read [AGENTS.md](AGENTS.md), the current [Linear issue](https://linear.app/swale-solutions/issue/SWA-46/overhaul-existing-repository-into-a-typescript-workspace) and [AutoTube v2 decisions](https://linear.app/swale-solutions/project/autotube-v2-69e62d254724). Each issue gets one bounded feature branch and PR, Conventional Commits, and exact verification results. Never push directly to `main`, rewrite shared history, commit secrets or reset existing data.

See [migration notes](docs/swa-46-migration.md) for the starting commit, retained behaviour, v1 recovery and the next work. V1 code and its fixtures remain available in Git history.
