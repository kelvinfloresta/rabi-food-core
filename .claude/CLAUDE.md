# rabi-food-core

Multi-tenant subscription food ordering platform.

- `api/` — Go modular monolith (Fiber, pgx, samber/do). Tests: `task test` (Docker only).
- `web/` — React Router 8 SSR app (Tailwind v4, shadcn/ui, i18next). Check: `npm run typecheck` in `web/`.
- `infra/` — local dev compose (`task db`) and Terraform (AWS ECS + ALB).
- `docs/product-context.md` — customer profiles and product decisions. Read it before any UX or product rule change, and flag conflicts with it explicitly.

## Language
All code, comments, documentation and commit messages are in English.

## How we work
- **Features:** plan first (plan mode), iterate on one plan, then implement. No spec/plan/tasks document hierarchy.
- **Small changes and bugfixes:** implement directly.
- **Business rules (pricing, cutoff, billing, tenant isolation):** the integration tests are the spec. A rule change comes with a test that would fail without it.
- **Architecture decisions** with lasting consequences (feature boundaries, new infrastructure, cross-cutting patterns): record a short ADR in `docs/decisions/` — see the template there.
- **Before committing:** run the checks of the affected subproject and ask the `api-reviewer` or `web-reviewer` subagent for a review.

## Where instructions live
- `.claude/rules/` — conventions, scoped with `paths:` so they load only for matching files.
- `.claude/agents/` — reviewers with isolated context. `.claude/skills/` — on-demand workflows (`/product`).
- `.claude/settings.json` — hooks for what must be enforced (Go files are gofmt-ed after every edit).

## Maintaining these files
- Only write what cannot be inferred by reading the code: intent, invariants, past mistakes, non-obvious constraints. Existing code is the reference for structure and naming.
- Add a rule when Claude makes the same mistake twice, not preemptively.
- When a change makes a rule wrong, fix the rule in the same change. A stale rule is worse than no rule.
- If something must always happen, make it a hook, lint rule or test instead of prose.
