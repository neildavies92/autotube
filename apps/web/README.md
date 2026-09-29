# Web

React/Vite foundation page. Run `pnpm dev:web` from the root and visit http://127.0.0.1:5173. `/api/*` proxies to the local Fastify API, stripping `/api`.

After `pnpm build`, `pnpm --filter @autotube/web preview` serves built assets on port 4173 (no API proxy). Navigation and API availability states are SWA-48 work; workflow controls and evidence inspection follow later.
