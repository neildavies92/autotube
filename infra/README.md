# Local Infrastructure

Local development services live here.

## PostgreSQL

`compose.yaml` runs PostgreSQL 16 for local development. Copy the example environment file before starting services:

```bash
cp .env.example .env
```

The committed example uses these defaults:

- database: `autotube`
- user: `autotube`
- password: `autotube`
- port: `5432`

Start it from the repository root:

```bash
make db-up
```

Recreate it from scratch and apply migrations:

```bash
make db-reset
```

Use this connection string for local API development:

```text
postgres://autotube:autotube@localhost:5432/autotube?sslmode=disable
```

Keep real local credentials in `.env`; it is ignored by Git.
