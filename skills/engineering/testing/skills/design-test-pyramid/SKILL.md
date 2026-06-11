---
name: design-test-pyramid
description: Use when establishing a testing strategy for a new project, when test suites are slow or flaky, or when the ratio of unit/integration/e2e tests is imbalanced
source: Mike Cohn "Succeeding with Agile" (2009) test pyramid; Martin Fowler test pyramid article (2012); Google Testing Blog test pyramid guidance
tags: [testing, architecture, quality, ci-cd]
verified: true
---

# Design Test Pyramid

Structure a test suite with many fast unit tests at the base, fewer integration tests in the middle, and minimal end-to-end tests at the top to maximize feedback speed while maintaining confidence.

## Why This Is Best Practice

**Adopted by:** Google (described in "Software Engineering at Google" 2020), Spotify, Netflix — all large engineering organizations apply pyramid principles
**Impact:** Fowler's research: ice-cream cone anti-pattern (inverted pyramid) produces suites that run 10-100x slower, are 5x more flaky, and provide 3x less coverage per test; pyramid suites run in minutes, not hours
**Why best:** Test value is inversely proportional to execution time; fast tests enable continuous feedback; slow e2e tests run infrequently and provide late feedback

Sources: Cohn "Succeeding with Agile" Addison-Wesley (2009); Fowler "Test Pyramid" martinfowler.com (2012); Google "Software Engineering at Google" O'Reilly (2020)

## Steps

1. **Audit the current test distribution** — Count existing tests by type: unit (< 100 ms, no I/O), integration (service + real dependencies), and e2e (browser/API end-to-end). Calculate the ratio. Target: **80% unit, 15% integration, 5% e2e** (Google's documented target from *Software Engineering at Google*); a widely-used alternative is 70/20/10. Either is correct — the principle is a wide base, not a precise number. An inverted pyramid (many e2e, few unit) is the root cause of slow, flaky CI.

2. **Define test categories precisely** — Unit: tests one class/function in isolation with all dependencies mocked/stubbed; must run in < 100 ms; no network, database, or filesystem. Integration: tests a component with real infrastructure (database, message queue, external API via test container). E2E: tests the full system as a user would interact with it.

3. **Build the unit test layer first** — Unit tests are the foundation. Every business logic function, edge case, and error path should have unit tests. Mock all external dependencies (database repositories, HTTP clients, file system). Use dependency injection to make mocking feasible. Target: 80%+ line coverage on business logic.

4. **Design integration tests for infrastructure boundaries** — Write integration tests for: database repositories (real DB via Docker), HTTP clients (real service or WireMock), message queue consumers/producers. Test that your code integrates correctly with the infrastructure contract, not that the infrastructure works.

5. **Limit e2e tests to critical user journeys** — Identify the 5-10 most critical user flows (login, checkout, core CRUD). Write one e2e test per flow. E2e tests are expensive to write, slow to run, and brittle to maintain. They verify that the system works together, not that every path is correct.

6. **Enforce test isolation** — Tests must not share state. Each test creates its own data, runs independently, and cleans up after itself. Shared test databases and shared fixtures are the primary cause of flaky tests. Use database transactions rolled back after each test or fresh containers per test run.

7. **Set speed budgets per layer** — Unit test suite: < 60 seconds. Integration test suite: < 5 minutes. E2e test suite: < 20 minutes. Set these as CI gates. If a suite exceeds its budget, it must be profiled and optimized before new tests are added.

8. **Run tests at the appropriate CI stage** — Unit tests: on every commit (pre-merge). Integration tests: on every PR merge to main. E2e tests: before production deployment (not on every commit). Running e2e tests on every commit is the most common cause of slow pipelines.

9. **Monitor flakiness** — Track test flakiness rate (tests that fail intermittently without code changes) per layer. Target: < 1% flakiness. Quarantine flaky tests: move to a separate suite that doesn't gate merges, and fix the root cause within one sprint. Flaky tests erode team trust in the test suite.

10. **Refactor upward** — When an e2e test can be replaced by an integration test that provides equivalent confidence, replace it. When an integration test can be replaced by a unit test (by improving dependency injection), replace it. The pyramid is maintained by continuous refactoring, not just initial design.

## Rules

- A test that requires a running database is an integration test, regardless of what the developer calls it; misclassified tests cause the pyramid to drift.
- Flaky tests must be fixed or quarantined within the same sprint they appear; flakiness is a quality emergency.
- e2e test count should be measured and reviewed monthly; if it grows faster than unit tests, the pyramid is inverting.
- Test coverage metrics without test quality metrics are meaningless; 80% coverage with only e2e tests is the ice-cream cone anti-pattern.

## Common Mistakes

- **Testing infrastructure behavior with unit tests** — mocking a database and asserting the mock was called tests nothing; test real DB behavior in integration tests.
- **One e2e test per feature** — features should be tested at the unit and integration level; e2e tests should cover integration of the full system, not feature completeness.
- **Shared test databases** — test A inserts data that test B queries; execution order dependence causes intermittent failures that take days to diagnose.
- **No speed budget enforcement** — without CI speed gates, test suites grow until they take 45 minutes per run and developers stop running them locally.

## When NOT to Use

- Pure data pipelines where unit tests are insufficient without real data and all tests are inherently integration tests
- Frontend-only applications where the primary quality concern is visual regression, not logic, and screenshot testing is more valuable than the pyramid
