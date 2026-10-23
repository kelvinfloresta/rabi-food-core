# 003. Two entry points: landing page and vitrine

- Status: Accepted
- Date: 2026-05-01

## Context
The two customer profiles (see `docs/product-context.md`) arrive with different needs: the Discoverer must understand the subscription model before committing; the Decider already knows what they want and needs minimum friction.

## Decision
- **Landing page** for the Discoverer (search, social media, word of mouth). It communicates value before the subscription flow: convenience and organization, competitive pricing, fresh food, and support for local businesses.
- **Vitrine** for the Decider (QR code at the establishment, link shared by another customer). It *is* the product listing where the customer picks what to subscribe to: no "subscribe" button leading to another listing.
- The vitrine does not explain the product; it links to the landing page.

## Consequences
- Each page has one job, and neither duplicates the other's.
- Links that are shared or printed (QR codes) point to the vitrine, not the landing page.
