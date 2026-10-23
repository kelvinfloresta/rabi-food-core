---
paths:
  - "api/**"
---

# API architecture

Feature-based modular monolith: `api/features/<feature>/{controller,gateway,routes,usecases}`. Follow the layout and file naming of an existing feature (e.g. `features/category`).

## Boundaries
- Business logic lives in use cases only — not in controllers, not in gateways.
- Imports between features must be one-directional and intentional. If an import direction looks wrong, the modeling is wrong; do not hide it behind an interface. Known violation: `payment` and `subscription` depend on each other — do not add to it.
- A local interface for another feature's use case only when isolated testability is an explicit requirement; otherwise import the concrete type.
- New services are registered in `libs/di/new.go` with `do.Provide`. Adapters copy `config` values when created, so nothing may build the injector before config is final (tests override config in `fixtures/environment.go`).

## Security
- Tenant access model (why: `docs/decisions/001-shared-customer-accounts.md`):
  - **Customers** can use any tenant: customer-facing routes (catalog, order, subscription, payment) take the tenant from the `X-Tenant-ID` header (`middlewares.RequireTenantID`). No membership is required, by design. The header only selects the tenant — it authorizes nothing.
  - **Tenant owners/staff** act through a scoped token: `exchange-token` / `refresh-scoped` check membership (`GetMembership`, 403 otherwise) and put `tenant_id` and `tenant_role` in the claims.
- Acting **as** the tenant — changing its data, or reading data across its customers (e.g. every customer's subscriptions, `GET /tenant/me`) — requires the scoped token: `middlewares.RequireTenant` on the route, tenant from `app_context.GetSession(ctx)`, never from the header.
- On customer routes, queries must be scoped to the session user, or any user can read other customers' data by sending their tenant's ID.
- A route serving both (e.g. `GET /subscription`, `GET /order`) skips `RequireTenant` and branches on `session.IsTenant()`: session tenant when true, header tenant plus session user otherwise.
- `UserID` comes from the session; tenant and user are never read from the request body.
- Token failures (missing, malformed, invalid) are always `401 INVALID_TOKEN`, via the jwtware `ErrorHandler` in `libs/http/fiber_adapter.go`.
- SQL is built with squirrel/placeholders, never by string concatenation.

## Errors
- Domain errors are `*errs.ApiError` declared in `libs/errs/<feature>_error.go` via `newErr(CODE, status)`. Never `errors.New` in use cases or gateways.
- Compare errors with `errors.Is` / `errors.As`.

## Handlers
- Parse, then `validator.V.Struct(data)`, then set `logger.GetWideEvent(uctx).Event = "<verb>-<entity>"`, then call the use case. Some read handlers still miss the WideEvent event — add it when touching them.
- A `BodyParser` error must become a 400. Never `return ctx.JSON(err)`: it answers **200** with the parse error — a known bug in most create/patch handlers; fix it when touching them.

## Health check
- `GET /health` (`libs/http/health_check.go`) is registered before the middlewares and is probed by the ALB; its path must match `infra/terraform/alb.tf`.
- It checks only dependencies the API cannot run without (database, GoTrue). Never add third-party services (e.g. Stripe): their outage would drain every task.
