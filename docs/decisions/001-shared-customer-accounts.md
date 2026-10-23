# 001. Shared customer accounts across tenants

- Status: Accepted
- Date: 2026-10-23

## Context
Customers buy from several local establishments. Asking them to sign up again at each one adds friction exactly where the Decider profile arrives (QR code, shared link) and wants to subscribe right away.

## Decision
A customer account belongs to the platform, not to a tenant.

- One user (GoTrue) is reused across all tenants. A customer can browse, order and subscribe at any tenant without membership; customer-facing routes take the tenant from the `X-Tenant-ID` header, which only selects where the customer acts.
- Acting **as** a tenant — changing its data or reading data across its customers — requires a token scoped to that tenant, issued by `exchange-token` / `refresh-scoped` only to members (`iam.tenant_members`).
- `iam.tenant_customers` records the business relationship (who is a customer of a tenant, via `enroll-customers`). It does not gate access.

## Consequences
Easier:
- One signup, profile and payment method for every establishment.
- A customer of one establishment is one step away from another, so tenants benefit from each other's customers.
- A single user pool: no per-tenant authentication or account duplication.

Harder, and accepted:
- **Security relies on every query.** The header authorizes nothing, so every customer-route query must be scoped to the session user; forgetting it exposes other customers' data (it happened once, in the subscription list). Enforced by the rule in `.claude/rules/api/architecture.md` and by integration tests.
- **Data ownership and LGPD.** The platform holds the customer's account; which customer data each tenant may see, who is the data controller, and what happens to it when a tenant leaves must be defined before launch.
- **Profile data is global.** Phone, address and name changes apply to every tenant; there is no per-tenant version of a customer.
- **No per-tenant blocking yet.** A tenant cannot ban a customer without new state (likely on `iam.tenant_customers`).
- **Platform branding.** Login and OTP emails are the platform's, not the establishment's.

Open question: the product positions itself against marketplaces, and shared accounts with cross-tenant discovery resemble one. Where the line is (e.g. no commission, no tenants competing for the same customer, no cross-promotion) is not decided yet.
