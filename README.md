# AutoTube

AutoTube is a local-first, compliance-first application for researching niches and producing original evergreen YouTube explainer videos for a single MVP channel.

The project is being built as a learning vehicle for full-stack engineering, data analysis, and agentic systems, with a deliberate emphasis on understanding Go and its ecosystem deeply rather than hiding the foundations behind heavy frameworks.

## What It Will Do

- Research user-defined niches and identify durable evergreen topic opportunities.
- Produce weekly plans capped at 2-3 publishable videos for the MVP channel.
- Generate original draft packages, assets, previews, and metadata.
- Keep provenance, quality checks, and human approval in the workflow.
- Upload approved videos privately before any manual public release.

## Why It Exists

The aim is to build a data-driven evergreen content pipeline that can support steady channel growth while remaining original, reviewable, and compliant by design.

## Current Status

AutoTube is in **Milestone 1: Project foundation plus basic niche management**.

Current foundation work includes project governance, documentation, Git workflow, shared Codex skills, the first runnable backend/frontend project structure, a local PostgreSQL persistence foundation, and the base dashboard health/status integration.

V1 is intentionally single-channel. Multi-channel portfolio management is a future capability, but implementation should avoid hard-coding choices that would block it later.

## Planned Stack

- Backend: Go 1.24, standard-library-first
- HTTP: `net/http` with `http.ServeMux`
- Logging: `log/slog`
- Database: PostgreSQL via `database/sql`
- Query generation: `sqlc`
- Frontend: React + Vite + TypeScript
- API contract: REST + OpenAPI
- Durable background jobs: deferred until the first real async workflows are needed

## Repository Layout

```text
apps/api             Go API and backend foundations
apps/web             React/Vite dashboard
contracts/openapi    REST API contract
db/migrations        PostgreSQL migrations
infra                Local services and Docker Compose
skills               Project-owned Codex skills
```

## Getting Started

Prerequisites:

- Go 1.24
- Node 22, as declared in `.nvmrc`
- npm
- Make
- Docker with Docker Compose
- sqlc, for regenerating Go query code
- GitHub CLI, authenticated for PR-aware cleanup commands

Install frontend dependencies:

```bash
make setup
```

Create a local environment file:

```bash
cp .env.example .env
```

Run the Go API:

```bash
make dev-api
```

Run the React/Vite dashboard in another terminal:

```bash
make dev-web
```

The dashboard proxies `/api/*` requests to the local API at `127.0.0.1:8080`, so keep the API running to see the live health/status state.

Run the local verification suite:

```bash
make check
```

Start PostgreSQL and apply the initial schema:

```bash
make db-up
make db-migrate
```

Recreate the database from scratch:

```bash
make db-reset
```

The Makefile is a thin wrapper around the root npm scripts and Go commands, so the underlying commands remain visible in `package.json` while local development gets one consistent entry point.

Useful individual commands:

```bash
make test-api
make typecheck-web
make build-web
make sqlc-generate
make post-merge-cleanup BRANCH=<merged-branch>
make clean
```

After merging a pull request, run `make post-merge-cleanup BRANCH=<merged-branch>` to fetch and prune remote refs, switch back to `main`, update it from GitHub, and delete the merged local branch. The cleanup helper uses a normal safe branch delete first, then only force-deletes a squash-merged branch when GitHub confirms the PR was merged into `main`.

## Working Conventions

- `main` is protected; all changes should land through pull requests.
- Work on short-lived branches named like `feat/AT-13-short-description`.
- Use Conventional Commits.
- Rebase feature branches during development and squash merge PRs into `main`.
- Keep Notion updated when architecture, requirements, workflow, or durable decisions change.
- Use the pull request template in `.github/` for every change.

## Documentation

Long-form product, architecture, compliance, operating, and decision documentation is maintained in Notion. Jira owns task tracking; the repository README is intentionally kept as the concise front door for contributors and future readers.

## Maintainer

Maintained by Neil Davies.
