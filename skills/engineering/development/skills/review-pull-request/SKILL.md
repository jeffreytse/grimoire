---
name: review-pull-request
description: Use when asked to review a pull request, diff, or code change — whether reviewing as an AI assistant or guiding a human reviewer on what to check and how to give feedback.
source: Google Engineering Practices documentation (google.github.io/eng-practices)
tags: [code-review, pull-request, feedback, reviewer, quality, collaboration]
verified: true
---

# Review Pull Request

Systematically evaluate a pull request for correctness, design, readability, and tests — then deliver feedback that is actionable and kind.

## Why This Is Best Practice

**Adopted by:** Google (published as Engineering Practices documentation), Microsoft (Code With Engineering Playbook), Shopify, Airbnb. Google's process is the most detailed publicly available code review standard and is widely referenced across the industry.

**Impact:** Google's internal research (Sadowski et al., "Modern Code Review: A Case Study at Google", ICSE 2018) found that code review catches ~15% of bugs before production and is the primary mechanism for knowledge transfer across teams. Small, well-reviewed PRs merge 40% faster than large ones (SmartBear survey, 2019, n=600 developers).

**Why best:** Checklist-free "eyeball it" reviews miss systematic categories of issues (security, test coverage, API stability). Google's ordered review — design first, then correctness, then style — prevents bikeshedding on formatting before catching a logic bug. This is superior to reviewing in reading order, which buries structural feedback in line-level noise.

Sources: Google Engineering Practices (google.github.io/eng-practices), Sadowski et al. ICSE 2018, SmartBear State of Code Review 2019

## Steps

### 1. Orient before reading

Read the PR description, linked issue/ticket, and any screenshots. Answer: *what problem is this solving and why now?* If the description is absent, note it as a blocking issue before reviewing further.

### 2. Review design (top-down first)

Before line-level comments, assess the overall approach:
- Is this the right place to solve the problem?
- Does it introduce the right abstraction, or does it overfit to one use case?
- Are there simpler alternatives that achieve the same result?
- Does it match the existing architectural patterns of the codebase?

Flag design concerns as blocking (`[BLOCKING]`) because they may make line-level changes irrelevant.

### 3. Check correctness and logic

- Does the code do what the description claims?
- Are edge cases handled: empty input, nulls, integer overflow, concurrent access?
- Are error paths tested and do they fail cleanly?
- Are there race conditions or TOCTOU (time-of-check/time-of-use) issues?

### 4. Evaluate tests

- Do tests exist for the new behavior?
- Do tests cover the failure paths, not just the happy path?
- Are tests readable as documentation — does each test name describe the scenario?
- Would a future developer be able to diagnose a test failure from the test name alone?

### 5. Check for security issues

- Are inputs validated and sanitized before use?
- Are secrets hardcoded or logged?
- Is authorization checked at the right layer (not just UI)?
- Are SQL queries parameterized? Are file paths validated?

### 6. Review readability and maintainability

- Are names self-explanatory without requiring comments?
- Are comments present where the *why* is non-obvious (not restating the *what*)?
- Is the code consistent with the codebase's style and idioms?
- Are any functions or classes doing too much (violating single responsibility)?

### 7. Check for breaking changes

- Do public API signatures, database schemas, or wire formats change?
- Is there a migration path or deprecation notice?
- Is the change backward-compatible if deployed incrementally?

### 8. Write and deliver feedback

Write comments. Then before posting, re-read them and apply these filters:
- Mark each comment as `[BLOCKING]`, `[SUGGESTION]`, or `[NIT]` so the author knows what must change
- Blocking: correctness bugs, security holes, missing tests for new behavior, broken API contracts
- Suggestion: better approaches the author should consider but may reject with justification
- Nit: style preferences — optional, do not block approval
- Phrase blocking comments as specific observations, not personal judgments: "This will panic if `items` is nil" not "you forgot nil handling"
- Acknowledge what is done well — at least one genuine positive comment per review

### 9. Decide: approve, request changes, or comment-only

**Google's standard:** Approve a CL once it *definitely improves the overall code health of the codebase*, even if it isn't perfect. Do not withhold approval waiting for perfection — perfect is the enemy of merged. The goal is incremental improvement, not a single pass to ideal.

- **Approve**: all blocking issues resolved; code improves or at least doesn't degrade codebase health
- **Request changes**: one or more blocking issues remain — correctness bugs, security holes, missing tests for new behavior
- **Comment**: questions or suggestions where you lack enough context to block

**Handling pushback:** When the author disagrees with a comment, consider their argument seriously — they are closer to the code and may be right. If their argument is sound, acknowledge it and drop the comment. If you still believe your suggestion improves code health, explain *why* more fully. Hold firm on code health concerns; defer on style preferences. Never approve a known health problem just to end the argument.

## Rules

- Never block on style issues that a linter should enforce — configure the linter instead
- Never leave a review open for more than one business day; it blocks the author
- If the PR is too large to review meaningfully (>500 lines of logic), say so and ask for it to be split
- Do not approve code you do not understand; ask questions instead
- Do not rewrite the author's code in a comment — describe the problem and let the author solve it

## Examples

**Blocking comment (correct):**
> `[BLOCKING]` `processPayment()` calls the charge API and then writes to the database, but if the DB write fails the charge is not reversed. This will cause double-charges on retry. Consider wrapping in a transaction or implementing idempotency.

**Nit (correct):**
> `[NIT]` Variable name `d` could be `duration` for clarity — take or leave.

**Bad comment (avoid):**
> Why did you write it like this? This is confusing.

## Common Mistakes

- Reviewing in reading order — catches style before design; fix by reviewing design first
- Only commenting on problems — authors disengage; give at least one genuine positive
- Blocking on nits — kills velocity; mark style issues as NIT and move on
- Approving without reading tests — tests are half the code; always read them
- Vague feedback like "this could be better" — always specify *what* and *why*
