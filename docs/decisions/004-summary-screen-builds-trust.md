# 004. The summary screen is where trust is built

- Status: Accepted
- Date: 2026-05-01

## Context
The final step of the subscription wizard is where the customer decides whether to commit. A commitment without a trial (see 002) only works if the customer trusts what they are agreeing to.

## Decision
The summary is a trust-building screen, not a bureaucratic checklist. It answers the questions the customer has not asked out loud:

- When will I receive? Concrete dates, not just weekdays.
- How much will I pay, and when?
- What happens if I want to cancel? Substitution and refund shown separately (see 002).
- Can I change later?

## Consequences
- The summary needs concrete delivery dates and charge dates before the customer confirms, so they must be computable before the subscription exists.
- It is the most important step of the flow; changes to it deserve the most care.
