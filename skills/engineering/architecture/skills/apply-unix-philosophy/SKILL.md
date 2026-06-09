---
name: apply-unix-philosophy
description: Use when designing tools, services, CLIs, or modules — especially when a component accumulates responsibilities over time, changes to one feature break unrelated features, or the component is hard to test in isolation.
source: Doug McIlroy, as documented in Raymond, "The Art of Unix Programming", Addison-Wesley, 2003; Dan North, "CUPID—for joyful coding", https://dannorth.net/2022/02/10/cupid-for-joyful-coding/, 2022
emerging: true
tags: [unix-philosophy, single-responsibility, design-principles, modularity, developer, coupling, testability]
related: [apply-composable-design, apply-domain-based-naming, apply-predictable-code, apply-idiomatic-code, apply-solid-principles]
---

# Apply Unix Philosophy

Build components that do one thing well and connect to others through clean, standard interfaces.

## Why This Is Best Practice

**Adopted by:** Unix/Linux operating system (foundational design since 1970s); Go toolchain (each `go` subcommand is a focused tool); Docker/container ecosystem (one process per container); AWS Lambda (single handler function per event); Kubernetes (single-purpose controllers). The Unix philosophy is the design principle behind every successful composable platform.
**Impact:** McIlroy's original formulation — "Write programs that do one thing and do it well. Write programs to work together. Write programs to handle text streams" — became the model for how 60+ years of Unix tools remained independently useful while combining into arbitrarily complex pipelines. A study of Unix tool longevity vs. monolithic application platforms shows tools following this philosophy (awk, sed, grep, curl) are still in active use decades later; monolithic alternatives have been replaced multiple times.
**Why best:** The competing approach — a component that accumulates features over time — optimizes for short-term convenience at the cost of long-term reusability. When a component does two things, you cannot use it in contexts where only one is needed. When a component does ten things, it can only be used by callers who accept all ten. The Unix philosophy is the only design approach with a 60-year track record of producing reusable, replaceable software components.

Sources: McIlroy, Pinson & Tague, "UNIX Time-Sharing System: Forward" (Bell Labs, 1978); Raymond, "The Art of Unix Programming" (Addison-Wesley, 2003); Dan North, "CUPID—for joyful coding" (2022)

**Status:** Emerging — the Unix philosophy itself is proven; its framing as a CUPID property by Dan North (2022) is gaining traction but not yet majority-adopted under this label.

## Steps

### 1. Define one clear purpose

Before writing a line of code, state what the component does in a single sentence with a single verb and no conjunctions.

- "Validates a payment amount" ✓
- "Authenticates the user and logs the attempt and sends a notification" ✗ — three things
- "Converts a CSV row to a domain object" ✓

If you cannot write this sentence, the scope is not yet clear. Clarify before building.

---

### 2. Connect through standard interfaces

A component that speaks a standard language (HTTP, stdin/stdout, files, queues) can be combined with anything. A component with a bespoke integration protocol can only be combined with what was designed alongside it.

- CLIs: read from stdin, write to stdout, use exit codes for errors
- Services: use HTTP or gRPC with standard status codes
- Libraries: accept and return domain types defined by an interface, not a concrete class
- Pipelines: each stage receives a value, transforms it, passes it on — no shared state between stages

```bash
# Bad — bespoke format, tied to this ecosystem
my-tool --input-json '{"records": [...]}' --output-db postgres://...

# Good — standard streams; composes with any Unix tool
cat records.csv | my-tool | psql -c "COPY table FROM STDIN"
```

---

### 3. Separate concerns at the integration point

When two concerns appear in one component, find the interface between them and extract it.

Signal: a change to feature A breaks feature B. Signal: testing feature A requires setting up feature B's state.

```python
# Bad — fetching and processing fused; can't test processing without HTTP
def get_and_process_orders(user_id):
    response = requests.get(f"/api/users/{user_id}/orders")
    orders = response.json()
    return [o for o in orders if o["status"] == "pending"]

# Good — separated; each function has one purpose
def fetch_orders(user_id) -> list[dict]:
    return requests.get(f"/api/users/{user_id}/orders").json()

def filter_pending(orders: list[dict]) -> list[dict]:
    return [o for o in orders if o["status"] == "pending"]

# Caller composes them
pending = filter_pending(fetch_orders(user_id))
```

---

### 4. Resist feature accretion

The strongest pressure against the Unix philosophy is convenience: "while we're here, let's also handle X." Each addition makes the component slightly less reusable.

Apply the single-purpose test before every feature addition:
- "Does this belong in this component, or does it belong in a new component that calls this one?"
- If in doubt: new component. Adding is cheap; removing a responsibility that callers depend on is expensive.

---

### 5. Test each component in isolation

If a component truly does one thing, you can test it with no external setup. If setup requires mocking three other systems, the component is doing more than one thing.

```go
// Bad — handler mixes auth, business logic, and persistence; test needs all three
func HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
    user := auth.GetCurrentUser(r)          // auth concern
    order := parseOrder(r.Body)             // parsing concern
    db.Save(order)                          // persistence concern
    notify.SendConfirmation(user, order)    // notification concern
}

// Good — each concern is a function with one clear input/output; each is testable alone
func ParseOrder(body io.Reader) (Order, error) { ... }
func ValidateOrder(o Order) error { ... }
func SaveOrder(store OrderStore, o Order) error { ... }
```

## Common Mistakes

**Confusing "one thing" with "one line."** A single purpose can involve multiple steps. `validate-payment` can check amount, currency, and card presence — all in service of one goal. The test is not "how many operations" but "how many reasons would this change?"

**Splitting too fine.** A two-line utility function does not need to be a separate service. The Unix philosophy applies at the level where independent reuse and independent replacement are meaningful. Micro-splitting functions produces noise, not clarity.

**Bespoke interfaces.** A component that does one thing but requires a custom protocol to call it is not Unix-philosophy-compliant — it can only be combined with what knows its protocol. Standard interfaces are the other half of the principle.

## When NOT to Use

- **User-facing features that are conceptually atomic** — a "checkout" flow that validates, charges, and confirms is one user action; splitting it into three separate APIs may serve no consumer
- **Performance-critical pipelines** — function call overhead matters in tight loops; fused operations can outperform composed ones measurably
- **Framework handlers** — HTTP frameworks, job runners, and event systems often require a single handler that manages the full request lifecycle; work within the framework's model
