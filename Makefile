SHELL := /bin/sh

.DEFAULT_GOAL := help

.PHONY: help setup dev-api dev-web test test-api typecheck-web build build-web check clean

help: ## Show available local development commands.
	@awk 'BEGIN {FS = ":.*## "; printf "AutoTube local development commands:\n\n"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

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

clean: ## Remove local build and coverage output.
	rm -rf apps/web/dist apps/web/coverage coverage.out
