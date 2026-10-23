---
paths:
  - "web/**"
---

# Web conventions

## Rendering and control flow
- No ternaries in JSX. A single optional block uses `{condition && <X />}`; two mutually exclusive blocks use an early return before the main JSX.
- No `else` anywhere — use early returns. Nested ifs move into a helper function with early returns.

```tsx
// wrong
{isNew ? <FieldA /> : <FieldB />}

// right
if (!isNew) return <FieldB />
return <FieldA />
```

## Types and values
- `readonly` on every type and array (`readonly string[]`, `Readonly<Props>`).
- No inline types in JSX or function signatures — name them (`ProductSheetProps`, not `Props`).
- No magic numbers: `const OTP_CODE_LENGTH = 6`.
- Reusable lists use `as const` with a derived type: `const ITEMS = ["a", "b"] as const; type Item = typeof ITEMS[number]`.

## Component placement
- Used by one page only → that page's directory: `app/routes/<page>/_components/`.
- Used by more than one page → `app/components/`.
- Generic UI primitive without domain logic → `app/components/ui/`, built from the existing primitives (`Button`, `Sheet`, ...). Never put domain components there.
- Known debt: `app/components/subscription/` and `app/components/storefront/` are used only by the storefront page and should move to `app/routes/storefront/_components/`. Do not add new files there.

## i18n
- Every user-visible string goes through i18n — no string literals in JSX.
- Format placeholders (`(11) 99999-9999`, `00000-000`) are component constants, not translation keys.
- Server locale comes from the `i18nextMiddleware` registered in `root.tsx`: in loaders/actions use `getLocale(context)` / `getInstance(context)` from `~/i18n.server`. Never create an i18next instance or read the locale from the request.

## Authentication
- Tokens never reach the browser: access and refresh tokens live only in `HttpOnly` cookies, and `authMiddleware` (`~/lib/session.server`) refreshes the access token on the server before any loader/action.
- In loaders/actions use `getAccessToken(context)` — never read the cookie, accept a token in a body, or return a token to the client.
- Start/end a session with `createSessionHeaders(tokens)` / `destroySessionHeaders()`. The client only knows `isLoggedIn: boolean` from the loader.
- No `localStorage`/`sessionStorage` for anything session-related.
