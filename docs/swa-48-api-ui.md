# SWA-48 — API boundary and operational shell

Base: merged SWA-46 on `ff99a9f`. Current SWA-48 acceptance criteria and AutoTube v2 decisions were read in Linear before implementation.

## Delivered boundary

- Fastify retains `GET /health` and its health response fields. Zod validates environment, health query and response contracts.
- A UUID request correlation header is accepted or generated. Invalid supplied IDs return a sanitized 400 using a fresh UUID; response IDs also appear in structured request logs.
- Validation, parser, not-found and unexpected errors share one envelope. Logs allowlist request metadata and omit raw URLs, queries, headers, bodies and exception details; sensitive log keys are additionally redacted.
- Startup validates settings before binding. SIGINT/SIGTERM drain Fastify with an idempotent shutdown and ten-second failure deadline.
- React supplies Overview/API connection navigation and runtime-validated loading, last-check availability, failure and retry states. Five-second health request timeout; cancellation on unmount.
- Development and built-preview proxies share the local API target. The OpenAPI contract and app READMEs describe the boundary and exact commands.

No domain endpoints, database coupling, worker jobs, discovery, reporting, authentication product or deployment are introduced. Health remains process liveness only. The duplicated small health schema is intentional at independently compiled HTTP boundaries; a shared contracts package is unnecessary for this one response.

## Verification

From the repository root with pinned Node/pnpm:

```sh
pnpm install --frozen-lockfile
pnpm lint
pnpm typecheck
pnpm test
pnpm build
pnpm smoke
```

On 30 September 2026, `pnpm check` passed: lint, all app/tool typechecks, 22 tests across two files, all app builds and the full process smoke check. The smoke required loopback socket access outside the execution sandbox.

API tests inject requests without opening ports. They cover retained health behavior, config rejection, invalid query/correlation input, not-found, malformed-body and malformed-URL errors, unexpected exceptions and invalid server response contracts, log correlation and secret omission/redaction. React tests cover loading/navigation/success, network/HTTP/JSON/schema failures, retry recovery and cancellation. These run locally without provider calls or credentials.

The smoke script uses clean synthetic environment settings and real processes. It checks development and preview proxy connectivity, matching correlation IDs, validation error envelopes, API unavailability through the proxy, config/port failure exits and both SIGINT/SIGTERM shutdown. It requires local loopback binding on ports 8080, 5173 and 4173. Exact execution outcomes are recorded in the PR; sandboxed environments that prohibit socket binding must run this command with local network permission.

## SWA-49 handoff

After SWA-47 and SWA-48 are merged (or explicitly available as a verified stacked base), durable execution can consume the database package and the validated operational API boundary. Preserve liveness semantics, propagate correlation IDs to jobs, add graceful closure of worker/database resources, and add only the job/run telemetry schema needed for that issue. Keep usage/cost values currency-aware and distinguish unavailable cost from zero. No SWA-49 implementation is included here.
