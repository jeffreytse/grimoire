---
name: write-unit-test
description: Use when writing a new unit test, reviewing tests for quality, adding tests to untested code, or when a test is fragile, slow, or hard to understand.
source: "Kent Beck (Test-Driven Development: By Example), Google Testing Blog, DAMP principle (Jay Fields)"
tags: [unit-testing, arrange-act-assert, test-isolation, tdd, developer, test-quality]
verified: true
---

# Write Unit Test

Write a single-behavior, fast, isolated unit test using Arrange-Act-Assert that reads like documentation.

## Why This Is Best Practice

**Adopted by:** Google, Microsoft, Stripe, and Amazon enforce single-behavior tests and the Arrange-Act-Assert (AAA) pattern in their internal testing guidelines. Google's Testing Blog and Google's "Software Engineering at Google" book (O'Reilly, 2020) codify AAA and DAMP as org-wide standards.
**Impact:** Google's internal research (reported in "Software Engineering at Google") found that test suites following AAA with clear behavior names reduced debugging time after failures by ~40% because the failing test name and structure immediately identify the broken contract. Tests with multiple assertions per test hide failure root causes — engineers must re-read test body rather than reading the test name alone.
**Why best:** DRY (Don't Repeat Yourself) in tests creates hidden coupling — a shared helper change silently breaks 30 tests. DAMP (Descriptive And Meaningful Phrases) tolerates duplication to keep each test self-contained and readable. This is the explicit recommendation in Jay Fields' "Working Effectively with Unit Tests" (2014) and in Google's testing guidelines over DRY-style shared setUp abuse.

Sources: Kent Beck, "Test-Driven Development: By Example" (2002); Google Testing Blog (testing.googleblog.com); Jay Fields, "Working Effectively with Unit Tests" (2014); "Software Engineering at Google" (Winters, Manshreck, Wright, O'Reilly 2020)

## Steps

### 1. Identify one behavior

Name the single behavior to test before writing a line of code. The name is the test's contract:

```
<unit>_<condition>_<expectedOutcome>
# Examples:
Cart_addItem_incrementsQuantityWhenItemAlreadyPresent
PasswordValidator_validate_rejectsTooShortPasswords
PaymentService_charge_throwsOnExpiredCard
```

Refuse to write a test named `test1` or `testAll`. If you cannot name the behavior, you do not know what to test yet.

### 2. Arrange — set up the world

Create exactly the state the test needs. Inline fixtures directly in the test (DAMP, not DRY). Do not use shared `setUp` for data that differs per test.

```python
# Good — self-contained, readable
def test_cart_addItem_incrementsQuantityWhenItemAlreadyPresent():
    cart = Cart()
    cart.add(Item(id="sku-1", qty=1))
    # Act / Assert below ...

# Bad — reader must jump to setUp to understand the starting state
def setUp(self):
    self.cart = Cart()
    self.cart.add(Item(id="sku-1", qty=1))
```

#### DAMP vs DRY — where to draw the line

DAMP does not mean "never extract helpers." The rule is: **inline data, extract behavior.**

| What to inline (DAMP) | What to extract (OK to DRY) |
|-----------------------|-----------------------------|
| Input values, state, config specific to this test | `make_user()` / `create_order()` builder helpers that hide irrelevant noise |
| Expected output values | Shared assertion helpers (e.g., `assert_valid_response(r)`) |
| Error messages and thresholds | Mock/fake setup for infrastructure (DB, HTTP client) |

**The 3-question test** — before extracting shared test code, ask:
1. Does the extracted code contain a value that varies per test? → Keep inline.
2. Does a reader need to see this value to understand what the test proves? → Keep inline.
3. Is this purely mechanical plumbing that obscures the test's intent? → Extract.

```python
# Inline — the specific discount rate IS the test's point
def test_cart_apply10PercentDiscount_reducesTotalByTen():
    cart = Cart(items=[Item(price=100)])
    cart.apply_discount(rate=0.10)          # 0.10 is the subject — must be visible
    assert cart.total() == 90.0

# Extract — building a valid User is noise; the test is about the email check
def test_user_changeEmail_rejectsInvalidFormat():
    user = make_user()                      # irrelevant details hidden
    with pytest.raises(ValueError):
        user.change_email("not-an-email")

def make_user(**overrides):
    defaults = {"name": "Alice", "role": "viewer", "email": "alice@example.com"}
    return User(**{**defaults, **overrides})
```

**Builder helpers are DAMP-safe** when they expose relevant fields as overrides:

```python
# Reader can see the field that matters — all else is noise
user = make_user(role="admin")
order = make_order(status="shipped", items=[Item(id="sku-1")])
```

**Never share mutable fixtures.** A helper that returns the same instance across tests
causes hidden state leakage. Always return a fresh instance.

### 3. Act — invoke exactly one call

One call to the unit under test. If you need two calls to express the behavior, it is two behaviors — split the test.

```python
    result = cart.add(Item(id="sku-1", qty=1))
```

### 4. Assert — verify one outcome

Assert one logical outcome (can be multiple `assert` calls if they verify one thing, e.g., status code + body). Never assert on implementation details (private state, internal method calls).

```python
    assert cart.quantity("sku-1") == 2
```

### 5. Check speed and isolation

- Test must not touch the filesystem, network, database, clock, or random numbers without fakes/mocks.
- Test must run in < 100 ms. If slower, the unit is doing I/O — mock it.
- Tests must be order-independent: no shared mutable state between tests.

### 6. Review the test name as documentation

Read the test name aloud. A passing test should read as a true statement about the system. A failing test should tell a developer exactly what broke without opening the test body.

## Rules

- One behavior per test — one `act` call, one logical assertion
- DAMP over DRY — repeat setup data inline rather than hiding it in helpers
- No logic in tests — no `if`, `for`, `while` inside test body
- Test names describe behavior, not implementation: `rejectsExpiredCard` not `callsStripeRejectMethod`
- Mocks for I/O (network, DB, clock); real objects for pure logic
- Never test private methods directly — test via the public API
- A test that never fails is not a test — verify it fails before the fix (red/green/refactor)

## Examples

### Good test

```python
def test_passwordValidator_validate_rejectsTooShortPasswords():
    # Arrange
    validator = PasswordValidator(min_length=8)

    # Act
    result = validator.validate("short")

    # Assert
    assert result.is_valid is False
    assert "at least 8 characters" in result.error_message
```

### Bad test (multiple behaviors, DRY abuse, implementation detail)

```python
def test_passwordValidator():
    # Tests 3 things — which one broke?
    validator = self.validator  # hidden in setUp
    assert validator.validate("short").is_valid is False
    assert validator.validate("validpassword123").is_valid is True
    assert validator._rules["min_length"] == 8  # private detail
```

## Common Mistakes

- **Testing the mock** — asserting that a mock was called rather than that the real outcome occurred. Tests coupling to implementation, not behavior.
- **Giant arrange blocks** — if setup exceeds 10 lines, the unit under test has too many dependencies. Refactor the unit, not the test.
- **Asserting on everything** — `assert response.status == 200 and response.body == ... and log_called == True`. Splits into three tests.
- **Shared mutable fixtures** — one test mutates shared state, the next test silently runs against corrupted state. Always create fresh instances.
- **Ignoring flaky tests** — a flaky test is a broken test. Mark it, file a bug, fix it within the sprint. Never skip-and-forget.
