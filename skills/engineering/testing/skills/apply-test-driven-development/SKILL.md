---
name: apply-test-driven-development
description: Use when implementing any business logic, fixing bugs, or when the correctness of a function's behavior under various inputs needs to be proved
source: 'Kent Beck "Test Driven Development: By Example" (2002); Erdogmus et al. "On the Effectiveness of TDD" IEEE Software (2005); Beck "TDD is Dead. Long Live Testing?" debate'
tags: [testing, development, tdd, code-quality]
verified: true
---

# Apply Test-Driven Development

Write a failing test before writing implementation code to drive design, prove correctness, and build a regression safety net incrementally.

## Why This Is Best Practice

**Adopted by:** Pivotal Labs (mandated 100%), Thoughtworks, Extreme Programming practitioners; Kent Beck at Facebook; Google's internal testing culture
**Impact:** Erdogmus et al. IEEE Software: TDD reduces defect density by 40-80% in controlled studies; IBM case study: 40% reduction in bugs with 15-35% development time increase that is recovered in reduced debugging; Microsoft Research: 15-35% time overhead, 60-90% fewer bugs
**Why best:** Tests written after implementation verify what code does, not what it should do; tests written first force API design from the caller's perspective, producing simpler interfaces

Sources: Beck "Test Driven Development: By Example" Addison-Wesley (2002); Erdogmus, Morisio & Torchiano IEEE Transactions on Software Engineering (2005); George & Williams "An Initial Investigation of TDD" ACM (2003)

## Steps

1. **Red: write a failing test first** — Before writing any implementation, write the smallest test that captures the next behavior the system must have. Run it and confirm it fails. A test that passes without implementation is not testing anything. The test must fail for the right reason (not a compilation error, but an assertion failure).

2. **Write the test from the caller's perspective** — The test is the first client of your code. Write the API you wish existed. `result = Calculator.add(2, 3); assert result == 5`. This produces caller-friendly APIs with minimal coupling. Awkward tests reveal awkward APIs before the API is written.

3. **Green: write the minimum implementation to pass** — Write the simplest code that makes the test pass. Resist the urge to generalize prematurely. Hardcoding return values is acceptable temporarily. The goal of the green phase is only to pass the test, not to write the final implementation. Speed matters here.

4. **Refactor: clean the code** — Once the test passes, improve the implementation: rename variables, extract methods, remove duplication, apply patterns. Run tests after every refactoring step. The tests are the safety net that makes refactoring safe. Never refactor when tests are red.

5. **Repeat the cycle for each behavior** — Each RED-GREEN-REFACTOR cycle adds one behavior. Work in the smallest possible steps. A cycle should take 2-10 minutes, not 30 minutes. If a cycle is taking too long, the behavior being tested is too large; decompose it.

6. **Test one behavior per test** — Each test has exactly one reason to fail. One assertion (or one set of closely related assertions) per test. Tests named with "and" are two tests. Multiple behaviors per test make failure diagnosis harder and test maintenance more complex.

7. **Use test doubles strategically** — Mock external dependencies (HTTP clients, database adapters) at the boundary. Don't mock the system under test or its close collaborators. Over-mocking produces tests that pass even when business logic is wrong. Prefer real objects for value objects and business logic.

8. **Triangulate to drive generalization** — If the implementation is hardcoded (return 5), write a second test with different inputs that forces generalization (add(3, 4) == 7). Triangulation is the mechanism that evolves hardcoded implementations into real algorithms through additional tests.

9. **Keep the test suite fast** — TDD only works when tests run in < 1 second per cycle. Slow tests interrupt flow and lead to running them less frequently. Tests that require a database or network are integration tests; they are valuable but are not the TDD inner loop. Mock I/O to keep unit test cycles fast.

10. **Commit after each green cycle** — Small, frequent commits tied to passing test cycles make reverting cheap. "Revert to last passing state" is trivially fast with small commits. Long uncommitted sessions after many red cycles make error recovery painful.

## Rules

- Never write implementation before a failing test; code without a test is unvalidated assumption.
- Never move on from red to green with a passing test that passed for the wrong reason; investigate.
- The refactor step is mandatory, not optional; TDD without refactoring produces tested spaghetti code.
- Test names must describe behavior, not implementation: `should_return_error_when_input_is_negative`, not `test_method_X`.

## Common Mistakes

- **Writing tests after implementation** — tests written after code tend to match what the code does rather than what it should do; edge cases and error paths are under-tested.
- **Testing implementation details** — tests that assert internal method calls or private state break on refactoring; test behavior (inputs and outputs), not implementation.
- **Red-green without refactor** — skipping the refactor phase accumulates technical debt even with test coverage; TDD is a three-phase cycle.
- **Failing for the wrong reason** — a compilation error is not a meaningful red phase; fix compilation errors first, then confirm the test fails on an assertion.

## When NOT to Use

- Exploratory programming and research spikes where the goal is to learn, not to ship; write tests after the spike when the design is understood
- Legacy code with no test infrastructure where integration tests or characterization tests must be established before TDD is feasible
- Highly uncertain requirements that change faster than tests can be written and maintained
