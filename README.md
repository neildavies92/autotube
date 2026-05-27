# AutoTube

AutoTube is a local-first, compliance-first application for planning and producing original YouTube explainer videos from niche and trend data.

The project is being built as a learning vehicle for full-stack engineering, data analysis, and agentic systems, with a deliberate emphasis on understanding Go and its ecosystem deeply rather than hiding the foundations behind heavy frameworks.

## What It Will Do

- Research user-defined niches and identify promising topic opportunities.
- Produce weekly plans capped at 2-3 publishable videos.
- Generate original draft packages, assets, previews, and metadata.
- Keep provenance, quality checks, and human approval in the workflow.
- Upload approved videos privately before any manual public release.

## Why It Exists

The aim is to build a data-driven content pipeline that can support steady channel growth while remaining original, reviewable, and compliant by design.

## Current Status

AutoTube is in **Milestone 1: Project foundation plus basic niche management**.

The initial Go API and React/Vite workspace scaffold is in progress under AT-13. Current foundation work includes project governance, documentation, Git workflow, shared Codex skills, and the first runnable backend/frontend project structure.

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

Run the React/Vite dashboard:

```bash
make dev-web
```

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
make clean
```

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
