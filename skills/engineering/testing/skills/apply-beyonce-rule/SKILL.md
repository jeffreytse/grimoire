---
name: apply-beyonce-rule
description: Use when deciding whether to protect a behavior with a test, responding to an unexpected breakage, or reviewing code that depends on undocumented behavior — if a behavior matters enough to depend on, it needs a test; if it breaks without one, you have no standing to complain.
source: 'Winters, Manshreck & Wright "Software Engineering at Google" (O\'Reilly, 2020) ch.11; Titus Winters "C++ as a Live at Head Language" (CppCon 2017); attributed to Beyoncé "Single Ladies" (Columbia Records, 2008)'
tags: [testing, contracts, dependencies, api-design, software-engineering, behavioral-testing]
verified: true
---

# Apply Beyoncé Rule

If you liked it then you shoulda put a ring on it: if a behavior matters enough to depend on, write a test for it — if it breaks and no test caught it, you implicitly said you didn't care.

## Why This Is Best Practice

**Origin:** The Beyoncé Rule emerged at Google as the practical complement to Hyrum's Law. While Hyrum's Law observes that all observable behavior becomes depended upon, the Beyoncé Rule assigns ownership: dependency without a test is an implicit statement that breakage is acceptable. Titus Winters articulates it in *Software Engineering at Google* (2020, ch.11) in the context of large-scale changes (LSCs) — when Google's infrastructure team makes a codebase-wide change that breaks your code, you cannot object if you had no test. The test is the proof that the behavior mattered.

**Adopted by:** Google's LSC (Large-Scale Change) process is built around this rule — automated tooling makes changes across millions of lines; only tests catch breakage. Teams that complain about LSC-caused failures are asked: "Did you have a test?" If not, the breakage is the team's responsibility, not the LSC author's. The rule is also applied in reverse: if a provider changes behavior that consumers depended on without a test, the provider's change is considered legitimate.

**Impact:** The rule clarifies ownership of behavioral contracts across teams. Without it, "who's responsible for this breakage?" becomes a political argument. With it, the answer is always: whoever didn't have a test. This creates a strong incentive to write tests that actually protect the behaviors you care about, rather than tests that only verify the obvious happy path.

**Why best:** Tests are explicit contracts. Undocumented behavioral dependence is implicit and invisible — no one knows you rely on it until it breaks. A test makes the dependency visible, makes the failure immediate and attributable, and gives the provider an opportunity to discuss the change before shipping. Without a test, the dependency is real but unenforceable.

Sources: *Software Engineering at Google* ch.11; Winters CppCon 2017; Hyrum Wright, hyrumslaw.com

## Steps

### Step 1: Identify behavioral dependencies that lack test coverage

Audit your codebase for dependencies on behaviors you care about but haven't tested. These are the highest-risk implicit contracts:

- Third-party API response shapes, field names, and ordering
- Error codes and status codes from upstream services
- Sort order of results from databases or external APIs
- Timing and retry behavior you've calibrated your system against
- Side effects of calling an external service (emails sent, audit logs written)
- Behavior of internal shared libraries your team doesn't own

```bash
# Find callers of external service that have no corresponding test
grep -rn "PaymentService.charge\|UserService.get" src/ \
  | grep -v "_test\.\|_spec\.\|test_"
```

### Step 2: Write tests for behaviors you depend on

Once identified, write tests that fail if the behavior changes. The test is the "ring" — it claims ownership.

```python
# Depending on sort order from an upstream service? Test it.
def test_product_list_returns_newest_first():
    """We depend on newest-first ordering for our homepage display logic."""
    products = catalog_service.list_products(limit=10)
    timestamps = [p.created_at for p in products]
    assert timestamps == sorted(timestamps, reverse=True), (
        "catalog_service must return products newest-first — homepage logic depends on this"
    )

# Depending on a specific error code? Test it.
def test_charge_returns_insufficient_funds_code_not_message():
    """Billing logic branches on INSUFFICIENT_FUNDS code, not message text."""
    with pytest.raises(PaymentError) as exc_info:
        payment_service.charge(card=expired_card, amount=100)
    assert exc_info.value.code == "INSUFFICIENT_FUNDS"
    # deliberately do not assert on exc_info.value.message — not our contract
```

### Step 3: Distinguish contract tests from implementation tests

For API *providers*, the Beyoncé Rule cuts both ways: tests you write on your own implementation inadvertently commit you to that behavior. Audit your test suite to separate:

| Test type | What it commits you to | Intentional? |
|-----------|------------------------|--------------|
| Tests the documented contract | The public API | Yes — this is the point |
| Tests error message text | Your exact wording | Probably not |
| Tests JSON field ordering | Your serialization order | Probably not |
| Tests internal sort behavior | Your algorithm's incidental output | Probably not |

Remove or relax tests that commit you to behavior you don't intend to guarantee:

```python
# BAD — this test commits you to the exact error message text forever
def test_invalid_input_error():
    with pytest.raises(ValueError, match="User not found in database"):
        service.get_user(-1)

# GOOD — tests the error type and code, not the message
def test_invalid_input_raises_not_found():
    with pytest.raises(NotFoundError) as exc_info:
        service.get_user(-1)
    assert exc_info.value.code == "USER_NOT_FOUND"
```

### Step 4: When accepting or rejecting a breakage, apply the rule

When a behavior breaks — whether from a dependency update, a platform change, or a colleague's refactor — the first question is: "Did we have a test?"

**No test existed:**
- The breakage is your responsibility. Add a test now to prevent silent future breakage.
- If the behavior change is unacceptable, negotiate with the provider — but acknowledge you had no formal claim.

**Test existed and caught it:**
- The test fulfilled its purpose. You have standing to negotiate with the provider.
- The provider should treat this as a breaking change requiring a deprecation window.

```python
# When a dependency update breaks behavior you depend on:
# 1. Write a regression test that captures the expected behavior
def test_upstream_returns_iso8601_timestamps():
    """Regression: upstream changed to epoch seconds in v2.1 — we need ISO 8601."""
    response = upstream_client.get_event(event_id="evt_123")
    # Verify format: should raise if not ISO 8601
    datetime.fromisoformat(response.created_at)

# 2. Pin the dependency until the behavior is restored or you migrate
```

### Step 5: Apply in code review

In code review, flag when a consumer depends on undocumented behavior without a test:

```python
# Code review comment:
# ⚠️ This code depends on the sort order of `inventory.list()`,
# which isn't documented as part of its contract.
# If this ordering matters (it looks like it does for the display logic),
# add a test that asserts the expected order, and/or confirm
# with the inventory team that they consider ordering stable.
for item in inventory.list():
    display_items.append(item)
```

### Step 6: Use the rule to prioritize test coverage

The Beyoncé Rule gives you a decision heuristic for where to spend test-writing effort:

1. **High dependency, no test** → write the test immediately
2. **High dependency, test exists** → good; ensure the test runs in CI
3. **Low dependency, no test** → acceptable; note it as a known risk
4. **Testing behavior you don't guarantee** → remove or relax the test

## Rules

- A behavioral dependency without a test is a claim you're willing to accept silent breakage. Be intentional about that choice.
- Tests you write on your own implementation commit you to that behavior — don't test internal details you're not prepared to support.
- When breakage occurs, the test audit comes first: did a test catch it? If not, the ownership of the fix belongs to the consumer.
- The "ring" (test) must run in CI — a test that only runs locally is not a ring.
- Distinguish the documented contract (test assertively) from incidental behavior (test loosely or not at all).

## Common Mistakes

**Testing the behavior you happen to produce, not the behavior you intend to guarantee:** Every test is a promise. Tests on error messages, internal field names, or implementation details commit you to behavior you probably intend to change freely.

**Treating breakage as the provider's fault when no test existed:** If upstream changes behavior and your system breaks silently, you owned the risk by not having a test. The Beyoncé Rule is symmetric — apply it to yourself first.

**Writing tests after the breakage instead of after the dependency is established:** A test written after a fire catches the specific fire; it often misses adjacent behaviors. Write the test when you first establish the dependency.

**Pinning dependencies instead of testing them:** Pinning a library version avoids breakage but doesn't capture *why* — the next engineer who unpins it won't know what to watch for. Test the behavior; pin as a last resort.

## Examples

**Google LSC workflow:**
Google's infrastructure team makes automated changes to all Go code. If your code breaks: did you have a test? If yes — the LSC was a breaking change; the LSC author is responsible for migrating you. If no — your code broke because you had an untested dependency; fix it yourself. The rule makes large-scale automated refactoring safe at scale.

**Dependency update breaks JSON parsing:**
```python
# Before: upstream returned {"ts": "2024-01-01T00:00:00Z"}
# After v2.0 update: upstream returns {"ts": 1704067200}
# Symptom: silent parsing failure in production for 6 hours

# Beyoncé Rule verdict: no test existed on timestamp format.
# Fix: add the test, then decide whether to pin or migrate.
def test_event_timestamp_is_iso8601():
    event = client.get_event("evt_123")
    datetime.fromisoformat(event.ts)  # raises if not ISO 8601
```
