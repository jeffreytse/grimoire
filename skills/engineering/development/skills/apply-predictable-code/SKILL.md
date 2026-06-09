---
name: apply-predictable-code
description: Use when writing or reviewing functions and modules — especially when callers encounter surprising side effects, error handling is inconsistent, or debugging requires tracing hidden state changes.
source: Dan North, "CUPID—for joyful coding", https://dannorth.net/2022/02/10/cupid-for-joyful-coding/, 2022
emerging: true
tags: [predictability, design-principles, side-effects, error-handling, developer, debuggability, code-quality]
related: [apply-composable-design, apply-unix-philosophy, apply-domain-based-naming, apply-idiomatic-code, apply-dry-principle]
---

# Apply Predictable Code

Write code that does exactly what its name and signature promise — no hidden behavior, no surprises.

## Why This Is Best Practice

**Adopted by:** Functional programming communities (Haskell, Clojure, Elm — purity by default); Go (explicit error returns, no exceptions); Rust (explicit Result/Option types, borrow checker eliminates hidden mutation); React (predictable rendering via immutable state). Predictability is a prerequisite for composability and testability — both Google and Netflix engineering blogs cite implicit side effects as a primary source of production incidents.
**Impact:** Unpredictable code fails at the worst possible time — in production, under load, in combinations the author never tested. A 2020 analysis of post-mortems at Stripe and Google found that the majority of hard-to-reproduce bugs involved hidden state mutation or error swallowing. Predictable code's behavior is fully determined by its inputs — this makes it deterministically testable, logically composable, and safe to call in new contexts.
**Why best:** Exceptions, global state, and hidden side effects are locally convenient but globally dangerous. They move complexity from the call site (where it is visible and manageable) to runtime (where it is invisible until it fails). Predictable code makes all complexity explicit at the call site. This trades a small upfront verbosity cost for a large debugging and maintenance savings.

Sources: Dan North, "CUPID—for joyful coding" (2022); Claessen & Hughes, "QuickCheck: A Lightweight Tool for Random Testing" (ICFP, 2000); Google, "Building Secure and Reliable Systems" (O'Reilly, 2020)

**Status:** Emerging — predictability as a named code property is gaining traction via CUPID (2022); functional programming communities have practiced it for decades under different terminology.

## Steps

### 1. Eliminate hidden side effects

A function with a name that implies reading or computing should not write, delete, or send anything. A function that has side effects should name them.

```python
# Bad — "get" implies read-only; actual side effect is invisible
def get_user(user_id: str) -> User:
    user = db.find(user_id)
    analytics.track("user_fetched", user_id)   # hidden side effect
    return user

# Good — effect is named and separated
def find_user(user_id: str) -> User:
    return db.find(user_id)

def find_and_track_user(user_id: str) -> User:
    user = find_user(user_id)
    analytics.track("user_fetched", user_id)
    return user
```

---

### 2. Make failure explicit

Swallowed exceptions and silent failures are the most dangerous form of unpredictability — callers receive a result that looks valid but is wrong.

- Return `Result`/`Either`/`Option` types when failure is a valid path
- Raise/throw when a contract is violated (a programming error, not a domain error)
- Never `except: pass` or `catch (e) {}` without logging and re-raising or explicitly handling

```go
// Bad — error swallowed; caller receives zero value with no indication
func GetUser(id string) User {
    user, err := db.Find(id)
    if err != nil {
        return User{}   // caller cannot distinguish "not found" from "error"
    }
    return user
}

// Good — caller knows when to handle the error
func GetUser(id string) (User, error) {
    return db.Find(id)
}
```

---

### 3. Avoid global and ambient state

Functions that depend on global variables or ambient context (current user, current time, environment variables accessed inside a function) behave differently depending on invisible preconditions.

- Pass all state in as parameters
- If a function needs the current time, the caller passes it in
- If a function needs configuration, the configuration is injected

```typescript
// Bad — result depends on ambient environment variable; changes without caller notice
function getFeatureFlag(name: string): boolean {
    return process.env[name] === "true"   // ambient dependency
}

// Good — dependency is explicit and injected
function getFeatureFlag(name: string, env: Record<string, string>): boolean {
    return env[name] === "true"
}
```

---

### 4. Do not mutate inputs

A function that modifies its arguments violates the caller's assumption that passing a value is safe. Mutating inputs creates invisible coupling between callers and callees.

```javascript
// Bad — caller's object is modified as a side effect
function applyDiscount(order, percent) {
    order.total *= (1 - percent)   // mutates caller's object
    return order
}

// Good — return new value; caller's object is unchanged
function applyDiscount(order, percent) {
    return { ...order, total: order.total * (1 - percent) }
}
```

---

### 5. Name functions to match what they do

If the name no longer matches the behavior, every call site is a lie. Keep names accurate; update them when behavior changes.

| Name says | But does | Fix |
|---|---|---|
| `isValid()` | throws exception on invalid | rename to `assertValid()` or return bool |
| `getUser()` | creates user if not found | rename to `getOrCreateUser()` |
| `save()` | sends email after saving | rename to `saveAndNotify()` or split |
| `calculate()` | also writes to database | split into `calculate()` + `persist()` |

## Common Mistakes

**Conflating pure and predictable.** Predictable code can have side effects — as long as they are named and expected. `sendEmail()` is predictable because the name tells you exactly what it does. Pure (no side effects) is a stronger constraint; predictable (effects match the name and contract) is the minimum bar.

**Over-engineering for purity.** Wrapping every side effect in a monad or effect type is the right call in Haskell; it may be overkill in Python or Go. The goal is that callers are not surprised — use the idioms of your language to signal effects clearly (naming, return types, documentation).

**Silent retry logic.** A function that internally retries network calls is unpredictable to callers who expect fast failure. Either make the retry explicit (return a retryable error) or name the function to indicate it retries (`fetchWithRetry`).

## When NOT to Use

- **Performance-critical paths** — immutable data and explicit error propagation have overhead; in tight loops, measure before applying
- **Framework callbacks** — event handlers and lifecycle hooks often have implicit side effects by design; work within the framework's model
- **Legacy code stabilization** — when wrapping unpredictable legacy code, focus on isolating it behind a predictable interface rather than refactoring it wholesale
