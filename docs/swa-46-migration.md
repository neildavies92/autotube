# SWA-46: TypeScript workspace migration

## Baseline and preservation

Starting `main` commit: **`111581b2bbb16ffc40d74038267b85b6e079a4df`**.

Implementation branch: `neildavies92/swa-46-overhaul-existing-repository-into-a-typescript-workspace`.

A fresh clone was clean before changes. Only `main` was advertised by the remote and the GitHub open-PR search returned no open PRs. No pre-existing local checkout or uncommitted files were modified. There was no applicable `AGENTS.md`; the existing `CLAUDE.md`, Git/PR guidance, application code, tests and build/configuration files were inspected. The current SWA-46 acceptance criteria and AutoTube v2 project decisions (29 September 2026) supersede the old Go/Jira/Notion baseline.

History is preserved. No database was connected to, migrated, reset or dropped; no Docker command was run and no existing volume, environment file or secret was changed. Retiring tracked schema/configuration files does not remove existing databases. SWA-47 must explicitly consider any v1 data before adopting a schema or volume.

Read a historical file without switching the current worktree:

```sh
git show 111581b2bbb16ffc40d74038267b85b6e079a4df:apps/api/internal/youtube/client_test.go
```

For a separate source recovery checkout (choose an unused sibling directory):

```sh
git worktree add --detach ../autotube-v1-reference 111581b2bbb16ffc40d74038267b85b6e079a4df
```

This recovers source only. Do not run historical reset commands or point historical migrations at existing data.

## Reuse and retirement inventory

| Existing component | Decision and future reference |
| --- | --- |
| React/Vite, `src/main.tsx`, HTML entrypoint | Retained; dependencies/configuration moved to the pinned pnpm workspace. |
| `GET /health`, timestamped JSON shape, OpenAPI | Retained in Fastify with matching contract; process liveness only. |
| `/api` development proxy to port 8080 | Retained. Smoke check verifies rewriting and connectivity. |
| Dashboard's speculative queues, publishing cadence and navigation | Replaced by a foundation page; SWA-48 owns the real operational shell and API availability states. |
| Go HTTP/config tests | Behavioural lessons retained: health, unknown routes, local defaults, invalid configuration, graceful shutdown. Vitest and process smoke checks cover the new entrypoints. |
| YouTube client and handler | Retired from active execution, not ported in this issue. At the baseline, `apps/api/internal/youtube/client.go` and `client_test.go` document search/metrics mapping, synthetic video/statistics payloads, empty-ID no-call behaviour and upstream errors. |
| Research handler cache/quota tests | `apps/api/internal/httpserver/server_test.go` at the baseline is useful for later provider work: missing query, provider-disabled route, upstream failure, cache hit/miss and quota attribution. Revalidate API/quota assumptions when implementing discovery. Preserve unknown-vs-zero metrics/cost distinctions rather than copying v1 defaults blindly. |
| Provider/workflow interfaces | Historical reference only. Do not prebuild speculative domain contracts or copy the old queue model over Graphile Worker. |
| Go storage, sqlc queries, SQL migrations and Compose | Retired from active checkout; available at the baseline under `apps/api/internal/storage`, `db`, `infra`. SWA-47 introduces PostgreSQL/Drizzle deliberately. Existing data/volumes remain untouched. |
| Go module/sums, sqlc config, Makefile, npm lockfile | Removed to leave one active pnpm/TypeScript path. No CI workflow existed; SWA-46 adds one. |
| V1 project skills and cleanup helper | Retired: stale Jira/Notion/Go commands and automatic branch cleanup conflicted with v2 delivery. Current conventions consolidated in root `AGENTS.md`; `CLAUDE.md` points there. |

No standalone fixture files were present. The useful provider examples were inline synthetic Go test fixtures; their exact recovery locations are above. V1 tests were inspected but are no longer an active verification target; no live provider calls were made.

## Version choices

Node 24.21.0 is pinned from the supported [Node 24 LTS line](https://nodejs.org/en/about/previous-releases). pnpm 10.34.6 supports this runtime. Exact direct dependencies and the pnpm lockfile were resolved against npm registry metadata during implementation. TypeScript 6.0.3 stays within `typescript-eslint` 8.71.0's supported `<6.1.0` peer range; TypeScript 7 is deliberately not selected. Fastify 5, React 19, Vite 8, Zod 4 and Vitest 5 install with strict peer checking.

## Verification and handoff

The root README is the command reference: `pnpm install --frozen-lockfile`, then `pnpm check` (lint, typecheck, test, build, smoke). Ten Vitest cases cover the retained API boundary and startup validation. Process smoke checks cover the built API/worker, Vite development/proxy and built web assets, invalid startup configuration, occupied ports, and clean API/worker shutdown. No credentials, database or live external services are required. These are HTTP/process checks, not browser E2E tests.

The worker intentionally remains idle. This PR does not implement database migrations, durable execution, domain endpoints, UI workflow controls, discovery, paid generation, rendering, publishing or Grafana.

After SWA-46 is reviewed and merged, **SWA-47** (PostgreSQL/Drizzle) and **SWA-48** (Fastify/operational UI shell) have their only prerequisite satisfied. SWA-49 depends on both of those tickets. Do not treat the minimal health route as completion of SWA-48.
