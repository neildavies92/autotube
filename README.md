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

The repository is being prepared before the application scaffold is added. Current foundation work includes project governance, documentation, Git workflow, and shared Codex skills that keep delivery consistent across machines.

## Planned Stack

- Backend: Go 1.24, standard-library-first
- HTTP: `net/http` with `http.ServeMux`
- Logging: `log/slog`
- Database: PostgreSQL via `database/sql`
- Query generation: `sqlc`
- Frontend: React + Vite + TypeScript
- API contract: REST + OpenAPI
- Durable background jobs: deferred until the first real async workflows are needed

## Planned Repository Layout

```text
apps/api             Go API and backend foundations
apps/web             React/Vite dashboard
contracts/openapi    REST API contract
db/migrations        PostgreSQL migrations
infra                Local services and Docker Compose
skills               Project-owned Codex skills
```

## Getting Started

The runnable application scaffold has not been added yet. Until then:

1. Review the project documentation in Notion for the current product and architecture decisions.
2. Use Jira as the execution backlog.
3. Follow the repository Git workflow before making changes.
4. Use the project-owned skills in [`skills/`](skills/) when working on AutoTube from another machine.

## Working Conventions

- `main` is protected; all changes should land through pull requests.
- Work on short-lived branches named like `feat/AT-13-short-description`.
- Use Conventional Commits.
- Rebase feature branches during development and squash merge PRs into `main`.
- Keep Notion updated when architecture, requirements, workflow, or durable decisions change.

## Documentation

Long-form product, architecture, compliance, operating, and decision documentation is maintained in Notion. Jira owns task tracking; the repository README is intentionally kept as the concise front door for contributors and future readers.

## Maintainer

Maintained by Neil Davies.
