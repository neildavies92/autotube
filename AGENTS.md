# AutoTube contributor and agent instructions

## Source of truth and scope

Read the current Linear issue and AutoTube v2 project decisions before editing. The issue's acceptance criteria are authoritative. Linear owns the v2 backlog and decisions; v1 Jira/Notion plans are historical references. If Linear cannot be accessed, say so and use the supplied issue scope without claiming to have read it.

Build incrementally: the first MVP ends at a ranked, evidence-backed topic backlog. Do not implement downstream discovery, generation, rendering, publishing or reporting in a foundation issue.

## Architecture and conventions

- Strict TypeScript, Node 24.21.0 LTS, pnpm 10.34.6; exact direct dependency pins and a committed pnpm lockfile.
- Fastify in `apps/api`, React/Vite in `apps/web`, a separate Node process in `apps/worker`.
- Keep modules within their app until there is an actual shared consumer. Do not import another app's internals.
- Use Zod at external/environment/HTTP/job boundaries when introduced. Keep startup separate from app construction for tests.
- Preserve the API error envelope and UUID request correlation. Log allowlisted metadata only; never log raw exception messages, request URLs with queries, headers or bodies. Add explicit redaction tests with synthetic sentinel secrets.
- Web availability is a runtime-validated, bounded request with loading/failure/retry states; it describes process liveness only. Keep navigation focused on available operational controls.
- Node apps use ESM and `.js` specifiers for relative imports; TypeScript compiles to `dist`. Web uses Vite/bundler resolution.
- PostgreSQL/Drizzle lives in `packages/database`; Graphile Worker arrives in SWA-49. No Redis or speculative downstream tables.
- React is for controls and evidence inspection; Grafana reporting and Remotion/FFmpeg are later work.
- Never log credentials or make live paid/provider calls in tests or CI. Use synthetic fixtures and local test doubles. Capture timestamps, correlation, durations, usage and currency-aware costs as the relevant operations are introduced; unknown cost is not zero.
- For future jobs, test failures and retries, preserve idempotency, reconcile external side effects and add spending/publishing limits before unattended operations.

## Exact commands

From the root, with the pinned runtime/package manager:

```sh
pnpm install --frozen-lockfile
pnpm dev
```

Individual development commands: `pnpm dev:api`, `pnpm dev:worker`, `pnpm dev:web`.

Required pre-PR checks (also run by `pnpm check` and CI):

```sh
pnpm lint
pnpm typecheck
pnpm test
pnpm build
pnpm smoke
```

Keep ports 8080, 5173 and 4173 free for smoke checks. Add meaningful tests at changed boundaries; do not inflate counts with constants or implementation-mirroring assertions. Unit/process checks need no database or provider secrets. Persistence changes also require `pnpm test:integration` against the dedicated local test instance; see `docs/swa-47-postgres.md`. Never substitute the development `DATABASE_URL` for `TEST_DATABASE_ADMIN_URL`.

## Issue-sized delivery

1. Inspect applicable instructions, Git status, existing code/tests and open PRs. Preserve uncommitted work.
2. Create a short-lived feature branch (e.g. `feat/SWA-47-postgres`, or Linear's suggested branch). Never commit/push directly to `main`.
3. Keep one bounded issue per PR; use Conventional Commits. Update README, contracts and migration notes when behaviour/setup changes.
4. Run the exact checks above and inspect the diff. Use `.github/pull_request_template.md`; include the Linear issue, commands/results and remaining limitations.
5. Move the issue to In Progress when work begins and In Review when the PR is ready if access permits. Report the next unblocked issue.
6. Do not merge, deploy, force-push, delete existing data or rewrite shared history without explicit authorization. Database resets must target isolated test data or be explicitly requested.

The v1 baseline is recorded in `docs/swa-46-migration.md`. Retired SQL and Compose files are historical; do not run them against existing data as part of v2 setup.

## Database changes

- `pnpm db:up` / `pnpm db:down` manage the persistent local v2 instance. Stop without removing volumes.
- `pnpm db:migrate` applies committed migrations; `pnpm db:generate` creates reviewable migrations from the Drizzle schema. Do not mutate applied migration files.
- `pnpm db:test:up`, then `TEST_DATABASE_ADMIN_URL=postgresql://autotube_test:autotube_test@127.0.0.1:55433/postgres pnpm test:integration`, then `pnpm db:test:down` runs isolated persistence checks. Each test run owns a freshly created database; reset means rerun, never resetting the development database.
- Keep API liveness independent of PostgreSQL. Add application tables only with their owning issue; SWA-49 introduces run/attempt records.
