---
paths:
  - "api/**/*_test.go"
  - "api/fixtures/**"
---

# API tests

## Structure
- Integration tests live in `features/<feature>/usecases/integration_test.go` as a testify suite. Copy the suite skeleton and request style from an existing one (e.g. `features/order/usecases/integration_test.go`); reuse `fixtures/` helpers (`fixtures.Tenant`, `fixtures.Auth`, `fixtures.Order`, ...) before writing new ones.
- A fixture exposes every input its request accepts (body, query, filters) as a parameter, before `auth`, even when no current test uses it — it is the extension point for future tests. Never remove such a parameter because every caller passes `nil`.
- Body inputs are pointers; for fixtures that create or configure something, a `nil` input means "valid defaults".
- Write every case inside `t.Run`: `fixtures.CleanDatabase` runs in `SetupSubTest`, so only subtests start from an empty database.
- Each package starts its own Postgres + GoTrue with testcontainers (`fixtures/environment.go`) and only then builds the DI. Database, GoTrue, `AUTH_SECRET` and the API port are set there, never read from `.env`. Keep the container images in sync with `infra/api/docker-compose.db.yaml`.

## Cases to cover for an endpoint
- Happy path, asserting the persisted result (read it back), not only the status.
- Management routes (owner/staff): another tenant's resource → **404**, never 403, so existence is not revealed; list endpoints must simply not return it. Customer routes are tenant-agnostic by design (tenant from `X-Tenant-ID`), so test that data is scoped to the header's tenant, not that access is denied.
- Customer routes: also test that a user sending another tenant's `X-Tenant-ID` sees only their own data (`fixtures.Auth.UserAuth(t, tenantID, "")` signs up a user with no membership).
- Owner/staff tokens in fixtures are signed directly (`fixtures.Auth.TenantOwnerAuth`), so the membership check in `exchange-token` / `refresh-scoped` needs its own auth tests — not covered yet.
- Insufficient role (e.g. customer on a management route) → 403 `FORBIDDEN`.
- Token failures → 401 `INVALID_TOKEN`. Covered once for all protected routes in `Test_AuthIntegration_ProtectedRouteToken`; do not repeat per endpoint.
- Missing required fields → 400 with code `VALIDATION_ERROR`.
- Missing resource → 404.
- Domain conflicts (duplicates, state transitions) → assert the domain error code (`errs.ErrX.Code`), not only the HTTP status.
- Known bug: a malformed JSON body returns **200** in handlers that do `return ctx.JSON(err)` after `BodyParser` (most create/patch handlers). The correct behavior is 400; do not write tests that assert the 200.

## Time
- Compare instants with `.Equal()`, never `==`.
- `task test` and CI run with `TZ=America/Sao_Paulo` on purpose, so any dependency on the server timezone fails. A plain `go test` uses your machine's timezone — run `TZ=America/Sao_Paulo go test ./features/<feature>/...` to reproduce CI for one package.
- The order test "should return timestamps in UTC regardless of the server timezone" guards the pgx codec in `PgxAdapter`.
