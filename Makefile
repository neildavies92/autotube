SHELL := /bin/sh

ifneq (,$(wildcard .env))
include .env
export
endif

POSTGRES_DB ?= autotube
POSTGRES_USER ?= autotube

DB_COMPOSE := docker compose --env-file .env -f infra/compose.yaml

.DEFAULT_GOAL := help

.PHONY: help setup dev-api dev-web test test-api typecheck-web build build-web check db-up db-down db-migrate db-rollback db-reset sqlc-generate post-merge-cleanup clean

help: ## Show available local development commands.
	@awk 'BEGIN {FS = ":.*## "; printf "AutoTube local development commands:\n\n"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-24s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Install project dependencies.
	npm install

dev-api: ## Run the Go API locally.
	npm run dev:api

dev-web: ## Run the React/Vite dashboard locally.
	npm run dev:web

test: ## Run all currently available tests.
	npm test

test-api: ## Run Go API tests.
	npm run test:api

typecheck-web: ## Run the frontend TypeScript typecheck.
	npm run typecheck:web

build: ## Build all currently buildable workspaces.
	npm run build

build-web: ## Build the React/Vite dashboard.
	npm run build:web

check: test typecheck-web build ## Run the local pre-PR verification suite.

db-up: ## Start local PostgreSQL.
	$(DB_COMPOSE) up -d postgres

db-down: ## Stop local PostgreSQL.
	$(DB_COMPOSE) down

db-migrate: ## Apply SQL migrations to local PostgreSQL.
	$(DB_COMPOSE) exec -T postgres sh -c 'until pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"; do sleep 1; done'
	$(DB_COMPOSE) exec -T postgres psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" -v ON_ERROR_STOP=1 < db/migrations/000001_initial_schema.up.sql

db-rollback: ## Roll back SQL migrations from local PostgreSQL.
	$(DB_COMPOSE) exec -T postgres psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" -v ON_ERROR_STOP=1 < db/migrations/000001_initial_schema.down.sql

db-reset: ## Recreate local PostgreSQL from scratch and apply migrations.
	$(DB_COMPOSE) down --volumes
	$(DB_COMPOSE) up -d postgres
	$(DB_COMPOSE) exec -T postgres sh -c 'until pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"; do sleep 1; done'
	$(DB_COMPOSE) exec -T postgres psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" -v ON_ERROR_STOP=1 < db/migrations/000001_initial_schema.up.sql

sqlc-generate: ## Regenerate type-safe Go query code from SQL.
	sqlc generate

post-merge-cleanup: ## Clean up locally after a PR merge. Usage: make post-merge-cleanup BRANCH=<merged-branch>
	sh scripts/git/post-merge-cleanup.sh "$(BRANCH)"

clean: ## Remove local build and coverage output.
	rm -rf apps/web/dist apps/web/coverage coverage.out
