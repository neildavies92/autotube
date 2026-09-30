# API

Run `pnpm dev:api` from the root, or `pnpm --filter @autotube/api start` after `pnpm build`. Defaults are `127.0.0.1:8080`; `AUTOTUBE_HOST`, `AUTOTUBE_PORT` and `AUTOTUBE_LOG_LEVEL` override them. Log levels are `fatal`, `error`, `warn`, `info` (default), `debug`, `trace` and `silent`. Export settings in the shell; app startup does not load `.env` automatically. Configuration errors report field names only.

`GET /health` retains the v1 health response fields: `service`, `status`, `message`, `checkedAt`. It measures process liveness, not database, worker or provider readiness. Unknown health query parameters return 400. The response is validated with Zod; the web client independently validates the contract at its HTTP boundary. The [OpenAPI contract](../../contracts/openapi/openapi.yaml) documents responses.

Each response carries `x-request-id`. Supply a UUID in this header to correlate an existing request; omitted IDs are generated with `randomUUID()`. Invalid/oversized/duplicate IDs return 400 with a fresh UUID, never the supplied value. Correlation IDs are tracing metadata, not authentication or authorization.

Errors share `{ "error": { "code": "...", "message": "...", "requestId": "..." } }`. Validation returns `INVALID_REQUEST`; malformed HTTP bodies return `BAD_REQUEST`; unknown routes return `NOT_FOUND`; unexpected exceptions return `INTERNAL_ERROR`. Raw exception and submitted input values are never returned. Fastify parser errors preserve appropriate 4xx status codes.

Structured JSON logs include request ID, method, response status and duration. Request URLs, queries, headers, bodies and exception details are omitted; secret fields have an additional redaction layer. Add safe, explicit metadata when diagnosing new boundaries. Never log full provider/database errors or arbitrary request objects.

`SIGINT` and `SIGTERM` close Fastify and drain accepted work; repeated signals do not trigger concurrent close calls. A ten-second deadline exits with failure if shutdown cannot finish. Startup/port errors exit nonzero without logging environment values.

Verification: `pnpm lint`, `pnpm typecheck`, `pnpm test`, `pnpm build`, `pnpm smoke` from the root. API tests use in-process injection and synthetic failures; smoke uses real local processes and verifies startup errors, proxy connectivity and both shutdown signals. No provider credentials are required.
