# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Setup
make setup              # npm install (run once after clone)
cp .env.example .env    # create local env file

# Development (run in separate terminals)
make dev-api            # go run ./apps/api/cmd/api
make dev-web            # vite dev server on 127.0.0.1:5173

# Database
make db-up              # start PostgreSQL via Docker Compose
make db-migrate         # apply migrations
make db-reset           # wipe and re-apply migrations from scratch

# Verification (run before opening a PR)
make check              # test-api + typecheck-web + build-web

# Individual checks
make test-api           # go test ./...
make typecheck-web      # tsc -b
make build-web          # tsc -b && vite build

# Code generation
make sqlc-generate      # regenerate apps/api/internal/storage/db from SQL

# Post-merge cleanup
make post-merge-cleanup BRANCH=<merged-branch>
```

Run a single Go test:
```bash
go test ./apps/api/internal/config/... -run TestLoad
```

## Architecture

### Overview

AutoTube is a monorepo with a Go API (`apps/api`) and a React/Vite dashboard (`apps/web`). The Vite dev server proxies `/api/*` → `http://127.0.0.1:8080` (stripping the `/api` prefix), so the frontend calls `/api/health` while the Go handler registers `GET /health`.

### Go API layers (`apps/api/internal/`)

| Package | Role |
|---|---|
| `config` | Loads `Config` from env vars via `LoadFromEnv`. Two opt-in flags: `RequireDatabase` and `RequireProviderSecrets` gate which missing vars are hard errors. |
| `httpserver` | `net/http` + `http.ServeMux`. Middleware stack: `recoverer → requestLogger → mux`. Add routes in `New()`. |
| `storage` | `OpenPostgres` opens a `*sql.DB`. sqlc-generated code lives in `storage/db` — do not hand-edit those files. |
| `providers` | Interface definitions only (no implementations yet): `YouTubeProvider`, `TrendsProvider`, `LLMProvider`, `ImageProvider`, `VoiceProvider`, `RenderProvider`, `UploadProvider`. |
| `workflow` | Domain types for the job queue: `Job`, `JobType`, `JobStatus`, `EnqueueRequest`. Interfaces: `Queue`, `EventSink`, `Handler`. No implementation yet. |

### Database

Schema: `db/migrations/000001_initial_schema.up.sql`  
Queries: `db/queries/niches.sql`, `db/queries/workflows.sql`  
Generated Go: `apps/api/internal/storage/db/` (via `sqlc generate`)

Three tables: `niches`, `workflow_runs`, `workflow_events`. `updated_at` is maintained by a trigger. When adding queries, edit `db/queries/*.sql` then run `make sqlc-generate`.

### Workflow job types

Defined in `workflow/workflow.go`: `research_niche`, `generate_draft`, `render_video`, `upload_video`, `sync_analytics`, `review_readiness`. Valid statuses: `queued → running → succeeded | failed | canceled`.

### Config env vars

| Var | Default | Notes |
|---|---|---|
| `AUTOTUBE_ENV` | `development` | |
| `AUTOTUBE_HTTP_ADDR` | `:8080` | |
| `DATABASE_URL` | — | Required only when `RequireDatabase` is set |
| `YOUTUBE_API_KEY`, `LLM_API_KEY`, etc. | — | Required only when `RequireProviderSecrets` is set |

## Conventions

- Branches: `feat/AT-NN-short-description`
- Commits: Conventional Commits (`feat:`, `fix:`, `chore:`, etc.)
- PRs: squash merge into `main`; use `.github/pull_request_template.md`
- Jira owns task tracking; Notion owns long-form docs (architecture, decisions, compliance)
- `main` is protected — all changes via PR
- The OpenAPI contract lives at `contracts/openapi/openapi.yaml`; keep it in sync with handler changes
