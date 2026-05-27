# Database Migrations

PostgreSQL migrations live here as paired `*.up.sql` and `*.down.sql` files.

Run the local database from the repository root:

```bash
make db-up
make db-migrate
```

Recreate the local database from scratch:

```bash
make db-reset
```

Roll back the current foundation schema:

```bash
make db-rollback
```

The foundation schema starts with:

- `niches` for user-defined channel niches.
- `workflow_runs` for later research, planning, generation, review, upload, and analytics jobs.
- `workflow_events` for durable workflow progress/provenance records.

`sqlc.yaml` points at the initial migration and `db/queries` to generate type-safe Go query code under `apps/api/internal/storage/db`.
