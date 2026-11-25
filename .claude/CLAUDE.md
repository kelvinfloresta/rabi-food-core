# rabi-food-core

Multi-tenant subscription food ordering platform.

- `api/` — Go backend.
- `web/` — React Router SSR frontend.
- `infra/` — local dev environment (`task db`) and Terraform.
- `docs/product-context.md` — customer profiles, positioning and the index of product decisions (`docs/decisions/`). Read it and the relevant decisions before any UX or product rule change, and flag conflicts with them explicitly.

## Language
Code, comments, technical documentation (READMEs, `.claude/`) and commit messages are in English. Product documents in `docs/` (product context, market analysis, decisions) are in Portuguese (pt-BR), so the product owner can read them without translating.

## Comments
Comments state the purpose of the code. Leave out what can change while the purpose stays the same (other components, endpoints, status codes, file names, current limitations) and rules already in `.claude/rules/`: those details go stale and mislead.

## How we work
- **Features:** plan first (plan mode), iterate on one plan, then implement. No spec/plan/tasks document hierarchy.
- **Small changes and bugfixes:** implement directly.
- **Product and architecture decisions** with lasting consequences (product rules and their reasons, feature boundaries, new infrastructure, cross-cutting patterns): record a short decision in `docs/decisions/` — see the template there.
- **Before committing:** run the checks of the affected subproject (listed in its rules) and ask the `api-reviewer` or `web-reviewer` subagent for a review.

## Where instructions live
- `.claude/rules/` — conventions, scoped with `paths:` so they load only for matching files.
- `.claude/agents/` — reviewers with isolated context. `.claude/skills/` — on-demand workflows (`/product`).
- `.claude/settings.json` — hooks for what must be enforced; each hook reports what it blocked and where its rules live.

## Maintaining these files
- This file loads in every session, whatever the subproject: keep here only what applies to the whole repository. What is specific to `api/`, `web/` or `infra/` goes in a rule under `.claude/rules/`, scoped with `paths:` to that subproject.
- Only write what cannot be inferred by reading the code: intent, invariants, past mistakes, non-obvious constraints. Existing code is the reference for structure and naming.
- Add a rule when Claude makes the same mistake twice, not preemptively.
- When a change makes a rule wrong, fix the rule in the same change. A stale rule is worse than no rule.
- If something must always happen, make it a hook, lint rule or test instead of prose.
