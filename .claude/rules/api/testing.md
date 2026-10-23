---
paths:
  - "api/**/*_test.go"
  - "api/fixtures/**"
---

# API tests

- Integration tests live in `features/<feature>/usecases/integration_test.go` as a testify suite. Copy the suite skeleton and the request style from an existing one (e.g. `features/order/usecases/integration_test.go`); reuse helpers from `fixtures/` (`fixtures.Tenant`, `fixtures.Auth`, `fixtures.Order`, ...) before writing new ones.
- Each package starts its own Postgres + GoTrue with testcontainers (`fixtures/environment.go`), then builds the DI. The environment is hermetic: config and secrets are set there, never read from `.env` files. Keep the container images in sync with `infra/api/docker-compose.db.yaml`.
- `fixtures.CleanDatabase` runs before each subtest; tests may assume an empty database.
- For an endpoint, cover at least: happy path, missing auth (401), another tenant's resource (403/404), invalid body (400), missing resource (404), and duplicates (409) when applicable.
- Compare instants with `.Equal()`, never `==`. Tests run with `TZ=America/Sao_Paulo` (Taskfile and CI) on purpose, so any dependency on the server timezone fails; the order test "should return timestamps in UTC regardless of the server timezone" guards the pgx codec.
- Run with `task test` (Docker required). `go test ./features/<feature>/...` runs a single package.
