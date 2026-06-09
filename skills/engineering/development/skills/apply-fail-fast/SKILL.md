---
name: apply-fail-fast
description: Use when writing functions, modules, or services to assert internal assumptions — preconditions, postconditions, invariants — so that violated assumptions produce an immediate, loud error at the point of violation rather than propagating corrupt state downstream.
source: Hunt & Thomas, "The Pragmatic Programmer" (2nd ed. 2019) §23; Hoare, "An Axiomatic Basis for Computer Programming", CACM 1969; Martin Fowler, "Fail Fast" (martinfowler.com, 2004); widely adopted at Google, Amazon (SRE), Erlang/OTP
tags: [fail-fast, assertions, invariants, error-handling, developer, debugging, reliability, defect-reduction]
related: [validate-external-input, diagnose-bug]
---

# Apply Fail-Fast

Assert internal assumptions; crash loudly at the violation site, not downstream.

## Why This Is Best Practice

**Adopted by:** Google Engineering Practices (assert liberally, crash on invariant
violation), Amazon SRE (loud failures over silent degradation), Hunt & Thomas
"Pragmatic Programmer" ("dead programs tell no lies"), Erlang OTP ("let it crash"
philosophy), and all major languages' standard library design (Java `Objects.requireNonNull`,
Python `assert`, Go `panic`, Rust `unwrap`/`expect`).
**Impact:** Microsoft Research found that 70% of production bugs in distributed systems
originate from error propagation, not the original fault (Nanz et al., 2015). Terminating
at the violation site reduces mean diagnosis time from hours to minutes because the stack
trace points directly at the violated assumption.
**Why best:** The alternative — defensive recovery ("if null, use default; if wrong type,
coerce") — hides programmer errors and allows corrupt state to propagate until it
manifests as a mysterious failure far from the origin. Fail-fast makes the defect visible
immediately.

Sources: Hunt & Thomas, "Pragmatic Programmer" 2nd ed. §23; Hoare (1969); Fowler,
martinfowler.com/ieeeSoftware/failFast.pdf; Erlang OTP Design Principles

## Steps

### Step 1: Distinguish programmer errors from expected failures

| Error type | Source | Correct response |
|------------|--------|-----------------|
| Programmer error | Violated assumption (null arg, empty list, illegal state) | Fail-fast: assert/panic/throw unchecked exception |
| Expected failure | External input, network, user action | Graceful error handling (return error, throw checked exception) |

Apply fail-fast **only** to programmer errors. For external failures, use
`validate-external-input`.

### Step 2: Assert preconditions at function entry

Check every assumption the function makes about its arguments before using them.

```python
def process_order(order_id: int, items: list[Item]) -> Receipt:
    assert order_id > 0, f"order_id must be positive, got {order_id}"
    assert len(items) > 0, "items must not be empty"
    assert all(i.quantity > 0 for i in items), "all item quantities must be positive"
    # ... implementation
```

### Step 3: Assert invariants at state transitions

After every operation that changes state, assert the invariant still holds.

```python
def transfer(self, amount: Decimal, target: Account) -> None:
    self.balance -= amount
    target.balance += amount
    assert self.balance >= 0, f"transfer created negative balance: {self.balance}"
    assert amount > 0, f"transfer amount must be positive: {amount}"
```

### Step 4: Log the violated assumption before crashing

Include both the expected condition and the violating value in the assertion message.

```python
# Poor — no context
assert user is not None

# Good — states the assumption and the violating value
assert user is not None, f"expected authenticated user in session {session_id}, got None"
```

### Step 5: Never catch-and-continue a programmer error

If code catches an assertion/panic at a high level to "recover," it defeats fail-fast.
Allow programmer errors to propagate to the top and crash (or restart the unit of work).

```python
# Wrong — hides the bug
try:
    process_order(order_id, items)
except AssertionError:
    return default_receipt  # bug is now silent

# Right — let it surface
process_order(order_id, items)  # if it throws, let it propagate
```

### Step 6: In languages where `assert` can be disabled, use explicit guards

Python's `assert` is disabled with `-O`. For critical invariants, use an explicit raise:

```python
if order_id <= 0:
    raise ValueError(f"order_id must be positive, got {order_id}")
```

## When NOT to Use

- **Data from external boundaries** — use `validate-external-input` instead. Failed
  external validation is NOT a programmer error; it is an expected condition.
- **Recoverable user-triggered errors** — invalid form input, rate limiting, not-found
  responses. These are expected; handle them gracefully.
- **High-frequency inner loops where assertion cost is profiler-confirmed significant** —
  guard with a compile-time or configuration flag, document the trade-off.

## Common Mistakes

**Using fail-fast for expected external failures.** If the network times out or a user
enters a bad value, that is not a programmer error. Panicking on expected conditions
makes the system fragile, not robust.

**Catching and swallowing assertion failures.** A catch-all `except Exception: pass`
at the top level silently buries programmer errors. Let them crash with a full stack trace.

**Asserting without a message.** A bare `assert x is not None` produces a useless
`AssertionError` in production logs. Always include the assumption and the violating value.
