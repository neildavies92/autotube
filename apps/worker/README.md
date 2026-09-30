# Worker

Separate TypeScript/Node entrypoint. It stays idle until interrupted and handles SIGINT/SIGTERM. No queue, database, provider calls or jobs are implemented. Graphile Worker execution and telemetry belong to SWA-49.

Run `pnpm dev:worker` from the root. After `pnpm build`, use `pnpm --filter @autotube/worker start`.
