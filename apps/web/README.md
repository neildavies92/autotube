# Operational web shell

Start `pnpm dev:api` and `pnpm dev:web` in separate terminals, or run `pnpm dev` for all apps. Visit http://127.0.0.1:5173. The Overview and API connection navigation provides a local operational shell; workflow controls and evidence inspection follow in later issues. Reporting belongs in future Grafana work.

Vite development and preview proxy `/api/*` to `http://127.0.0.1:8080`, stripping `/api`. After `pnpm build`, `pnpm --filter @autotube/web preview` serves built assets on http://127.0.0.1:4173 with the same proxy. Keep the API running. If changing the API port, update the shared proxy target in `vite.config.ts`. Neither Vite server is a public deployment configuration.

The shell checks `/api/health` on mount, validates the JSON response with Zod and displays loading, available-at-last-check or unavailable. Network, timeout, HTTP and contract failures all show an actionable retry state; responses are not cached. Requests time out after five seconds and are cancelled on unmount. `Check again`/`Retry connection` perform a fresh check; availability is not continuously polled and does not imply database/worker/provider readiness.

`pnpm test` exercises loading, success, navigation, network/HTTP/JSON/schema failure, recovery and cancellation using synthetic HTTP responses. `pnpm smoke` verifies real development/preview proxy connectivity and API unavailability. Run root `pnpm lint`, `pnpm typecheck`, `pnpm test`, `pnpm build`, `pnpm smoke` before delivery.
