---
name: design-test-strategy
description: Use when starting a new project, onboarding a codebase, auditing test coverage quality, or deciding where to invest testing effort across unit, integration, and end-to-end layers.
source: Mike Cohn (Succeeding with Agile, 2009), Martin Fowler (martinfowler.com/bliki/TestPyramid), Google Testing Blog
tags: [test-pyramid, test-strategy, coverage, integration-testing, e2e, developer, qa]
verified: true
---

# Design Test Strategy

Map the codebase to a test pyramid — assign each behavior to the correct layer so the suite is fast, trustworthy, and cheap to maintain.

## Why This Is Best Practice

**Adopted by:** Google, Spotify, Netflix, and Thoughtworks cite the test pyramid explicitly in their engineering blogs and internal standards. Google's "Software Engineering at Google" (O'Reilly 2020) dedicates three chapters to the pyramid and reports that Google targets 70% unit / 20% integration / 10% E2E across most services. Spotify's "Testing at Spotify" post (2018) describes the same ratio as the outcome of hard-won experience abandoning an inverted pyramid.
**Impact:** An inverted pyramid (many E2E, few unit tests) is the most common cause of slow CI. At Google, test suites that violated the pyramid took 10–30× longer to run and had 3–5× higher flakiness rates compared to pyramid-conformant suites (reported in "Software Engineering at Google", Ch. 11). The cost of a flaky E2E test is 20–40 minutes of engineer time to triage per failure (PagerDuty Engineering, 2020).
**Why best:** Unit tests are ~1000× faster than E2E tests (milliseconds vs. seconds) and have zero flakiness from infrastructure. Tests should fail only because the code is wrong, not because a network partition occurred. The pyramid shapes investment toward tests that give fast, deterministic signal while using E2E only to verify the integration seams that cannot be validated in isolation.

Sources: Mike Cohn, "Succeeding with Agile" (2009); Martin Fowler, "TestPyramid" (martinfowler.com); "Software Engineering at Google" (O'Reilly 2020, Ch. 11–14); Google Testing Blog (testing.googleblog.com)

## Steps

### 1. Audit what exists

Categorize every existing test as unit, integration, or E2E:

```bash
# Count test files by layer (adapt to project conventions)
find . -path '*/unit/*' -name '*test*' | wc -l
find . -path '*/integration/*' -name '*test*' | wc -l
find . -path '*/e2e/*' -name '*test*' | wc -l
```

Plot the current shape. Is it a pyramid, a rectangle, or an inverted pyramid?

### 2. Map risk to layer

For each system concern, assign the correct layer:

| Concern | Correct layer | Reason |
|---|---|---|
| Business logic, pure functions | Unit | No I/O, instant feedback |
| Single service + one real dependency (DB, queue) | Integration | Validates contract with real infra |
| Authentication, authorization rules | Integration | Security requires real tokens/sessions |
| Critical user flows (checkout, login, signup) | E2E | Cross-service contract |
| All UI permutations, edge case inputs | Unit (not E2E) | Too slow and fragile at E2E layer |
| Infrastructure config, networking | Infrastructure test (separate) | Not a code test |

### 3. Set coverage targets per layer

Use these as starting targets, adjusted for risk profile:

- **Unit tests:** 70–80% line coverage on business logic packages. 0% on infrastructure adapters (DB clients, HTTP clients) — test via integration instead.
- **Integration tests:** One happy path + one error path per external dependency (each DB table, each queue, each external API).
- **E2E tests:** Top 5–10 user journeys. If you have more than 20 E2E tests, audit them — most belong lower.

### 4. Define the test contract per layer

Write a one-paragraph definition for your project:

```
Unit: Tests in /unit. No network, no filesystem, no clock. Must run in < 5s total.
      Mock all I/O at the boundary. Use real objects for domain logic.

Integration: Tests in /integration. Run against a local Docker compose stack.
             One real external dep per test. Allowed to be slow (< 5 min total).
             Tagged @integration, excluded from default CI fast-path.

E2E: Tests in /e2e. Run against a deployed staging environment.
     Covers checkout, login, and account flows only.
     Run nightly and on release branches. Tagged @e2e.
```

### 5. Plug gaps — lowest layer first

Add missing tests starting at the lowest applicable layer. Ask for each gap: "Can this be tested with a unit test?" If yes, write a unit test. Only go higher if the behavior requires real infrastructure.

### 6. Gate CI on the pyramid

Configure CI to run layers in order, fail fast:

```yaml
# Example: GitHub Actions
jobs:
  unit:     # runs on every push, < 3 min
  integration:  # runs on PR merge, < 10 min, needs docker
  e2e:      # runs on release branch, < 30 min, needs staging
```

## Rules

- Never write an E2E test for behavior testable at the unit level — this is waste
- Integration tests own the boundary: one real dep, everything else mocked
- E2E tests must not share state between runs — each test sets up and tears down its own data
- Flaky tests are fixed within the same sprint — not skipped, not retried silently
- Coverage percentage alone is a vanity metric; what matters is behavior coverage (each branch of business logic has a test)
- A test that takes > 1 second is not a unit test — find and remove the I/O

## Common Mistakes

- **Ice cream cone** — many E2E, some integration, few unit. Slow, flaky, expensive to fix.
- **Coverage theater** — 90% line coverage on generated/trivial code, 30% on critical paths.
- **Integration tests that call production APIs** — leaks money and data; always use test doubles or local stubs.
- **No layer discipline** — unit tests that open sockets, E2E tests that assert on DB state directly. Blurs the contract and makes failures hard to diagnose.
- **One giant E2E test per feature** — forces sequential execution, balloons runtime as features grow.
