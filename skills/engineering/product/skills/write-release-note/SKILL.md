---
name: write-release-note
description: Use when writing user-facing release notes for a software product release
source: Stripe Developer Documentation style guide; GitHub Release Notes guidelines; Apple App Store release note best practices
tags: [release-notes, product, documentation, communication, changelog]
verified: true
---

# Write Release Note

Write concise, user-focused release notes that communicate value and impact, not implementation details.

## Why This Is Best Practice

**Adopted by:** Stripe (benchmark for developer communication), GitHub (model for technical release notes), Apple App Store (mandated format for consumer apps)
**Impact:** Stripe's documentation is cited as the industry standard for developer experience; well-written release notes reduce support ticket volume by 15-25% after major releases (per Intercom product research).

**Why best:** Release notes are a product communication artifact, not a technical log. Users care about what changed for them, not how it was implemented. Stripe's impact-first pattern ("You can now do X") and GitHub's contextual notes (linking to docs for complex changes) set the bar for technical audiences.

## Steps

1. **Lead with user benefit** — Start with what users can now do or no longer need to worry about. Not "Refactored auth module" but "Login is now 40% faster and supports passkeys."
2. **Group by impact type** — Use consistent categories: New, Improved, Fixed, Security, Breaking Changes. Order: New → Improved → Fixed → Security → Breaking.
3. **Write one sentence per item** — State the change and its user impact. Link to documentation for anything requiring more than one sentence.
4. **Highlight breaking changes prominently** — Use a visible header or warning; describe what breaks and provide the migration path or link to the migration guide.
5. **Target the right audience** — Developer-facing APIs: include code examples for breaking changes. Consumer apps: plain language, no jargon, benefit-oriented.
6. **Include version and date** — Always: version number, release date, and relevant links (full changelog, migration guide, docs).
7. **Review for jargon** — Read as a non-expert user; replace internal terms with plain language equivalents.

## Rules

- Every breaking change must be explicitly labeled and include a migration path.
- Security fixes must be included; describe the class of vulnerability (not the exploit detail) and advise users to upgrade.
- Omit implementation details (refactors, dependency bumps, internal tooling) unless they affect users.
- Use present or future tense for new features ("You can now..."), past tense for fixes ("Fixed an issue where...").

## Examples

Good (Stripe-style):
> **New:** Webhooks now support retry configuration. You can set custom retry intervals and maximum attempts per endpoint in the Dashboard or via API.
>
> **Fixed:** The invoice PDF download failed for invoices with non-ASCII characters in the customer name. This is now resolved.
>
> **Breaking:** The `charge.create` endpoint no longer accepts `source` as a bare card object. Pass a `PaymentMethod` ID instead. [Migration guide →]

## Common Mistakes

- **"Various bug fixes and performance improvements"** — tells users nothing; enumerate specific fixes.
- **Commit message verbatim** — developer-oriented messages lack context for users.
- **Omitting the "so what"** — "Added retry logic" → who cares; "Failed webhook deliveries are now retried automatically for 72 hours" → clear user value.

## When NOT to Use

- When the release contains only internal refactors, dependency upgrades, or infrastructure changes that have no observable effect on user behavior, publishing a release note creates noise without informing any user decision.
- When the change is a hotfix for an actively exploited security vulnerability, publish a minimal security advisory through your coordinated disclosure process instead of a detailed release note that could aid attackers before users patch.
- When the audience is exclusively internal engineering teams consuming an internal library with no external contract, a commit log or internal Slack post is sufficient; a polished user-facing release note format is mismatched to the audience.
