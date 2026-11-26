---
paths:
  - "api/**"
---

# Errors and logs

Each request (or cron run) writes one log line: the wide event (`logger.WideEvent`). A few indexed fields with few distinct values (`event`, `status`, `error_code`) find and count failures. The free-text `error` field is not indexed: it is read only after finding the line, for the external cause (driver message, response body). Logs are paid by volume, so the goal is short lines and few indexed fields, not more lines.

## Error codes locate the failure
- Domain errors are `*errs.ApiError` declared in `libs/errs/<feature>_error.go` via `newErr(CODE, status)`. `errors.New` in `features/` fails the lint (forbidigo).
- The code is assigned once, where the failure is first detected: a gateway wrapping a database or HTTP error, or a use case rejecting a business rule or an unusable response. Wrap the cause as `fmt.Errorf("%w: %w", errs.ErrX, err)`, ApiError first: the first ApiError in the chain is the response code and the `error_code` log field. Never return `errs.ErrX` in place of `err` (the log loses the cause).
- Above that point the error goes up with its code untouched. A caller, including a use case calling another, adds context only as text: `fmt.Errorf("create order: %w", err)`. Wrapping it with the caller's own code hides where it failed and can change the response (a 404 becomes a 500). The exception is a caller whose code tells the client to act differently (`SIGNUP_SIGNIN_FAILED`: the account exists, sign in instead): it goes first with `%w: %w`, so the inner code stays in the chain and in `error`.
- Every failure point gets its own code, internal ones included (500, e.g. `SUBSCRIPTION_STATE_VERIFICATION_FAILED`), because with `%w: %w` the code is the only pointer to the failing line. When several points of one operation share a code (e.g. begin, insert and commit of one write) and the cause's own text does not tell them apart, add a few words of context: `fmt.Errorf("%w: insert tenant: %w", errs.ErrX, err)`.
- `INTERNAL_ERROR` is only the fallback for errors nobody anticipated — never wrap `errs.ErrInternal` on purpose. In the log it means a failure point without a code.
- Always `%w`, never `%v`: the chain is how the code and the original error reach the error handler and the log. Our codes are read from the chain by type (`errors.Is`, `errors.As`, `errors.AsType`), never parsed from the error text.
- Known violation, to fix when touching it: `STRIPE_SERVICE_FAILURE` is shared by every Stripe operation.
- Exception: the health check runs before the request middlewares, so it has no wide event and logs its failures on lines of its own.

## Wide event
- Do not log a separate line where an error happens and then return it: the middleware logs every returned error once, on the request line.
- Each use case writes the fields of its own domain (`SetProductID`, `TrackOrderPayment`), whoever called it. The caller records nothing about the call: `event` is the flow, `error_code` the failure point.
- Set a field as soon as its value is known, before anything that can fail (input IDs at the start of the use case): a failed request's line must still say what it was working on. Values that only exist after a step (a created ID) are set right after it.
- New `WideEvent` fields are cleared in `reset()` (events are pooled: a missed field leaks into the next request's line) and covered by a test that reads the emitted line (`fixtures.Logger`).
- Entity fields hold one value, and the last write wins. Never set them inside a loop over items: record a count, and put the failing item's ID in the error (`FullCode`, e.g. `errs.ProductNotFound(ids...)`).
- Record in the wide event only failures that do not reach the response (handled, retried, degraded): they leave no trace in the error chain.
- A field meant to be filtered or counted has few distinct values (codes, statuses, kinds). IDs only for the request's own entities.
