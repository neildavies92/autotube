# AutoTube Web

React + Vite + TypeScript dashboard shell for the local-first AutoTube MVP.

The dashboard calls the Go API health endpoint through the Vite development proxy:

```bash
/api/health -> http://127.0.0.1:8080/health
```

Start the API before the dashboard to see the live `API online` state. If the API is not running, the dashboard shows its API offline state.

Run from the repository root:

```bash
npm run dev:web
```

Build from the repository root:

```bash
npm run build:web
```
