# Architecture decision records

One file per decision with lasting consequences: feature boundaries, new infrastructure, cross-cutting patterns. Not for features or bugfixes — those are planned in the conversation and recorded by the code and its tests.

Name files `NNN-short-title.md` (e.g. `001-payment-subscription-boundary.md`). Never rewrite an accepted record: when a decision changes, add a new one and mark the old one `Superseded by NNN`.

## Template

```markdown
# NNN. Title

- Status: Proposed | Accepted | Superseded by NNN
- Date: YYYY-MM-DD

## Context
What forces the decision, in a few sentences.

## Decision
What we do.

## Consequences
What becomes easier, what becomes harder, and what we accept.
```
