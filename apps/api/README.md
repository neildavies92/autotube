# AutoTube API

Go API service for the local-first AutoTube MVP.

This scaffold intentionally starts with the Go standard library:

- `net/http` and `http.ServeMux` for routing
- `log/slog` for structured logs
- `http.Server` timeouts and graceful shutdown
- `httptest` for handler tests

Run from the repository root:

```bash
npm run dev:api
```

Test from the repository root:

```bash
npm run test:api
```
