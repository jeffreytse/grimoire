---
name: write-unit-test
description: Use when implementing any function, method, or module that contains logic — before or immediately after writing the code.
source: Kent Beck Test-Driven Development (2002), Google Testing Blog, xUnit pattern (Meszaros 2007)
tags: [testing, tdd, developer, defect-reduction]
stable: true
---

# Write Unit Test

Verify a single unit of logic in isolation with an automated, repeatable test.

## Why This Is Best Practice

**Adopted by:** Google (mandated in Google Engineering Practices, 2014), Microsoft,
Meta, Stripe, Airbnb, and virtually every mature engineering organization. Kent Beck
formalized it with TDD in 2002; it has been mainstream for 20+ years without serious
challenge.
**Impact:** Unit tests reduce defect escape rate by 40–80% depending on coverage
(Capers Jones, Software Engineering Best Practices, 2010). Google reports that code
with unit test coverage has 50% fewer production incidents than untested code (Google
SRE Book, 2016). Teams practicing TDD deliver 15–35% fewer defects with no significant
productivity loss after initial learning (IBM Research, George & Williams 2004).
**Why best:** No other technique catches logic errors earlier or cheaper. Manual testing
is not repeatable; integration tests are too slow for rapid iteration; code review misses
subtle logic errors. Unit tests are the only practice that makes refactoring safe at speed.

Sources: Kent Beck TDD (2002), Google Engineering Practices, IBM Research George & Williams (2004)

## Steps

1. Write the test before writing the implementation (RED): define expected inputs and outputs.
2. Run the test — it must fail. A test that passes before implementation tests nothing.
3. Write the minimum code to make the test pass (GREEN).
4. Refactor the implementation while keeping the test green (REFACTOR).
5. Repeat for each distinct behavior.

## Rules

- Test behavior, not implementation: test what a function does, not how it does it.
- Each test must be independent — no shared mutable state between tests.
- Name tests descriptively: `it('returns null when input is empty')` not `it('test1')`.
- Keep tests fast: unit tests must complete in milliseconds, not seconds.

## Common Mistakes

- Testing the framework instead of application logic.
- Writing tests that only test the happy path — edge cases and error paths matter.
- Sharing state between tests so that test order affects results.
- Coupling tests to implementation details, causing them to break during refactoring.
