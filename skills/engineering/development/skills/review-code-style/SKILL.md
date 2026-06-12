---
name: review-code-style
description: Use when reviewing code for style consistency, readability, or adherence to a team or language style guide
source: Google Engineering Practices / Google Style Guides (google.github.io/styleguide)
tags: [code-review, style, linting, readability, quality]
verified: true
---

# Review Code Style

Systematically review code against a defined style guide using automated linters and human judgment.

## Why This Is Best Practice

**Adopted by:** Google, Airbnb, Meta, Microsoft (all publish and enforce internal style guides)
**Impact:** Consistent style reduces code review time by up to 20% and onboarding time by 30% (Google Engineering Productivity research)

**Why best:** Style consistency lowers cognitive load when reading code written by others. Automated linting catches mechanical issues before human review, freeing reviewers to focus on logic and design. Google's style guides are widely adopted as authoritative baselines across the industry.

## Steps

1. **Configure a linter** — Install and configure the canonical linter for the language (ESLint for JS/TS, Pylint/Ruff for Python, golangci-lint for Go, Checkstyle for Java). Point it at the team's chosen style guide profile.
2. **Run linter in CI** — Fail the pipeline on any lint error so no violation reaches review.
3. **Adopt a formatter** — Use an opinionated auto-formatter (Prettier, Black, gofmt) to eliminate whitespace and formatting debates entirely.
4. **Identify guide violations manually** — Read the diff for naming, structure, and comment quality that linters miss (e.g., vague variable names, missing doc-comments on public APIs).
5. **Leave actionable comments** — Reference the specific style rule and show the corrected form, not just the problem.
6. **Track recurring issues** — Log repeated violations; feed them back into linter config or team docs.

## Rules

- Never block a review solely on style when a linter or formatter could auto-fix it — configure the tool instead.
- Distinguish must-fix (violates agreed guide) from suggestion (personal preference).
- Apply the same guide uniformly; do not enforce stricter rules on certain contributors.
- Do not mix style-only commits with feature/bug commits — keep them separate.

## Examples

Bad: `var x = getUserData()` — vague name, `var` in modern JS.
Good: `const userData = fetchUserProfile(userId)` — descriptive name, correct keyword, explicit argument.

Review comment: "Per the Airbnb JS guide §13.1, prefer `const` over `let`/`var` when the binding is not reassigned. Please update."

## Common Mistakes

- **Running the linter locally only** — violations slip through; enforce in CI.
- **Reviewing style before logic** — wastes time if the logic will be refactored.
- **No agreed baseline** — without a shared guide, style debates are subjective and unresolvable.

## When NOT to Use

- When reviewing a proof-of-concept or throwaway script where style consistency has no long-term maintenance value, a style review wastes review bandwidth better spent on correctness.
- When the codebase has no linter configured and no style guide agreed upon, a style review before establishing that baseline will produce unresolvable subjective disagreements rather than actionable fixes.
- When reviewing a logic-heavy diff where the author is mid-refactor, block on logic correctness first; style feedback during structural churn will be invalidated by the next revision.
