---
paths:
  - "api/features/*/gateway/**"
  - "api/libs/database/**"
---

# Database access

## Context
- `context.Context` is the first parameter of every gateway method.
- Reads (SELECT) pass `ctx` as is, so a client cancel stops them.
- Writes (INSERT/UPDATE/DELETE) use `context.WithoutCancel(ctx)`: keep context values, but a client disconnect must not abort a write halfway.

## Filters and patches
- Filter struct fields never use pointers: the zero value means "no filter" (`""`, `time.Time{}`). Optional booleans use `filter.Bool` from `libs/database/filter`.
- `PatchValues` fields are pointers: `nil` means "do not update". This is the only place pointers are used for optionality.

## Time
- Date columns are `TIMESTAMPTZ`; gateways write `time.Now().UTC()`.
- `PgxAdapter` registers a `TimestamptzCodec{ScanLocation: time.UTC}` in `AfterConnect`, so everything read comes back in UTC regardless of the server timezone. Never change `time.Local` — it is process-global state.
- Timezone-dependent business rules (delivery dates, cutoff, cron) use `time.LoadLocation` + `.In(loc)` explicitly.

## Migrations
- Files in `libs/database/migrations/` with goose annotations; the app binary never runs them (`task migrate` does).
- **Pre-production (current):** schema changes go straight into the table's `create` migration — no incremental `add_x_to_y` files. Recreate the local DB afterwards: `task db-down && task db && task migrate && task seed`.
- **After go-live:** never edit an applied migration; every change is a new migration.
- Integration tests apply every migration on a fresh database: a SQL error at test startup means the migration chain is broken for any new environment.
