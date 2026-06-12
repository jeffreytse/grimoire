---
name: audit-test-coverage
description: Use when evaluating the adequacy of a test suite, identifying coverage gaps, or deciding where to add tests
source: Martin Fowler "TestCoverage" (martinfowler.com); mutation testing (Pitest for Java, mutmut for Python); Google Testing Blog
tags: [test-coverage, mutation-testing, quality, testing, audit]
verified: true
---

# Audit Test Coverage

Evaluate test suite quality by analyzing both line/branch coverage metrics and mutation testing scores to identify meaningful gaps.

## Why This Is Best Practice

**Adopted by:** Google (internal coverage standards), PIT mutation testing (Pitest — used across Java OSS ecosystem)
**Impact:** Mutation testing (Pitest) studies show that 80% line coverage can correspond to as low as 40% mutation score — meaning half the bugs survive despite passing coverage thresholds. Google's internal research found diminishing returns above 85% line coverage for most code.

**Why best:** Line coverage is a floor, not a ceiling. It tells you which lines execute during tests, not whether tests would catch a bug. Mutation testing injects real bugs and checks whether tests detect them — a far stronger signal of test quality.

## Steps

1. **Measure line and branch coverage** — Run coverage tool (Istanbul/nyc for JS, Coverage.py for Python, JaCoCo for Java, go test -cover). Focus on branch coverage, not just line coverage.
2. **Identify uncovered critical paths** — Filter coverage report to business-critical modules (payment, auth, data integrity). Ignore generated code, migrations, and trivial getters.
3. **Run mutation testing** — Run Pitest (Java), mutmut/Cosmic Ray (Python), or Stryker (JS/TS) on the critical modules. Target ≥70% mutation score.
4. **Classify survivors** — For each surviving mutant (bug not caught by tests): determine if it represents a real risk or an untestable equivalent mutation.
5. **Write missing tests** — Prioritize: (a) uncovered branches in critical paths, (b) surviving mutants in high-risk code, (c) error paths and edge cases.
6. **Set and enforce thresholds** — Configure CI to fail below agreed coverage minimums; do not set global thresholds above 80% — it incentivizes trivial tests.
7. **Report trends** — Track coverage over time; a declining trend signals test debt accumulation.

## Rules

- 100% coverage is not the goal — it incentivizes testing implementation details over behavior.
- Never chase coverage by writing assertions-free tests ("assertion-free testing anti-pattern").
- Prioritize coverage of code that changes frequently and code that handles money, auth, or data integrity.
- Exclude generated code, vendor code, and config from coverage metrics.

## Examples

Coverage gap analysis:
- `PaymentService.refund()` — 0% branch coverage, 4 survivors in mutation test
- Risk: high (financial impact)
- Action: write tests for partial refund, over-refund guard, and currency mismatch cases

## Common Mistakes

- **Treating line coverage as quality signal** — a test that calls code without asserting outcomes increases coverage without value.
- **Ignoring branch coverage** — a function with 100% line coverage may have untested `else` paths.
- **Setting coverage targets globally** — boilerplate code and generated code dilute meaningful coverage goals.

## When NOT to Use

- When the codebase is a throwaway prototype or spike with a planned rewrite, investing in coverage audits and mutation testing produces findings that will be discarded along with the code.
- When the module under review is entirely generated code (e.g., ORM migrations, protobuf stubs, GraphQL schema types), coverage metrics on generated output measure the generator, not engineering decisions, and should be excluded rather than audited.
- When the team is in an active incident response, a coverage audit is a non-urgent quality task that should not compete for attention with restoring system availability.
