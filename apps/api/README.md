# API

Fastify/TypeScript foundation. `GET /health` is process liveness only, not database or provider readiness. Run `pnpm dev:api` from the root. After `pnpm build`, use `pnpm --filter @autotube/api start`.

See the root [README](../../README.md) for settings and verification. Consistent HTTP errors, correlation and the operational API boundary are SWA-48 work.
