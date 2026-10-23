# 002. Subscription is a commitment; cancellation and refund

- Status: Accepted
- Date: 2026-05-01

## Context
The subscription model is what makes planning and fair pricing possible for the establishment. Trials and opt-out toggles weaken it, but customers fear being locked in, and Brazilian consumer law (CDC) grants a right of withdrawal.

## Decision
- **Always a commitment, never a trial.** The customer signs up and pays for the chosen plan; deliveries happen as scheduled.
- **Auto-renewal is not a toggle.** It is a fact of the subscription. The UI states it plainly: *"your subscription renews automatically at the end of the cycle. You can cancel renewal at any time."*
- **Cancellation without refund:** at any time. Already-paid deliveries continue until exhausted; no penalty, no refund — the subscription stops renewing.
- **Cancellation with refund (CDC):** only within 7 days of the **first delivery** (not the signup date). The refund is automatic and returns the value of undelivered orders.
- **Substitution is not a refund.** "I'd prefer another product" is a change request with no cost or friction. "I want my money back" is a refund, valid only in the 7-day window. The UI communicates them separately.

## Consequences
- The subscription stays a real commitment while the customer always has a clear exit.
- Refund rules depend on the first delivery date, so it must be known per subscription.
- Every flow that mentions cancelling must keep substitution and refund apart.
