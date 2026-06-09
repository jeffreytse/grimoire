---
name: apply-composable-design
description: Use when designing modules, services, or functions — especially when units are hard to reuse independently, replacements require touching many files, or combining two features requires coupling their internals.
source: Dan North, "CUPID—for joyful coding", https://dannorth.net/2022/02/10/cupid-for-joyful-coding/, 2022
emerging: true
tags: [composability, design-principles, coupling, modularity, developer, replaceability, maintainability]
related: [apply-unix-philosophy, apply-domain-based-naming, apply-predictable-code, apply-idiomatic-code, apply-composition-over-inheritance, apply-solid-principles]
---

# Apply Composable Design

Design software in small, focused units with minimal surface areas that can be combined freely without knowing each other's internals.

## Why This Is Best Practice

**Adopted by:** Unix (foundational design philosophy); Go standard library (small interfaces, single-method where possible); AWS Lambda / serverless ecosystem (each function does one thing); HashiCorp tools (each binary is self-contained and composable via APIs). The composability property is a prerequisite of microservices, serverless, and plugin architectures — all dominant patterns in modern software.
**Impact:** Composable units are independently replaceable. When a unit can be swapped without touching its consumers, teams can evolve components at different speeds, conduct A/B experiments at the module level, and limit blast radius of failures. Google's Site Reliability Engineering documentation identifies tight coupling as the primary cause of cascading failures — composable units contain those failures at boundaries.
**Why best:** Non-composable software grows through accretion — every new feature reaches into existing units, widening surface areas and deepening coupling. The only escape is a rewrite. Composable design prevents accretion: each unit's surface area is a contract; internals change freely as long as the contract holds. This is distinct from the OOP "composition over inheritance" pattern — composability is a property of any unit (functions, services, CLIs, modules) not just classes.

Sources: Dan North, "CUPID—for joyful coding" (2022); Eric Raymond, "The Art of Unix Programming" (2003); Google, "Site Reliability Engineering" (O'Reilly, 2016)

**Status:** Emerging — gaining adoption since Dan North's 2022 publication; not yet majority-adopted across all engineering organizations.

## Steps

### 1. Minimize the public surface area

Every exported function, method, or endpoint is a commitment. Each one must be maintained, versioned, and explained.

- Expose only what callers need to achieve their goal — nothing more
- Default to private/unexported; promote to public when there is a concrete use case
- Prefer a single general function over many specialized variants

```go
// Bad — many entry points, callers must know which to call
func SaveUser(u User) error { ... }
func SaveUserWithAudit(u User, by string) error { ... }
func SaveUserToShardDB(u User, shard int) error { ... }

// Good — one entry point, options injected
type UserStore interface {
    Save(ctx context.Context, u User) error
}
```

---

### 2. Make dependencies explicit

Hidden dependencies make composition impossible — you can't combine units when their needs are invisible.

- Pass all dependencies in via parameters, constructors, or interfaces
- Never reach for global state, singletons, or ambient context
- If a function needs the current time, pass it in — don't call `time.Now()` inside

```python
# Bad — hidden dependency on time; impossible to test or compose predictably
def is_expired(token):
    return token.expiry < datetime.now()   # grabs global time

# Good — time is explicit; callers control it
def is_expired(token, now):
    return token.expiry < now
```

---

### 3. Design units that can be used alone

A composable unit works standalone — it does not require its siblings to be present.

Test: can you import and use this module in a blank project with no other modules from this codebase?

- If no, identify the hidden dependency and extract or inject it
- Circular imports/dependencies are a hard signal that units are not composable

```typescript
// Bad — UserService requires OrderService requires UserService (circular)
class UserService {
    constructor(private orders: OrderService) {}
    getProfile(id: string) {
        return { user: this.findUser(id), orders: this.orders.forUser(id) }
    }
}

// Good — profile assembly is a separate unit; each service works alone
class UserService { findUser(id: string): User { ... } }
class OrderService { forUser(id: string): Order[] { ... } }
class ProfileService {
    constructor(private users: UserService, private orders: OrderService) {}
    getProfile(id: string) { ... }
}
```

---

### 4. Connect via interfaces, not implementations

Composable units know each other's contracts, not internals.

- Define interfaces at the consumer — not the provider (inversion of dependency ownership)
- Accept interfaces in function signatures; return concrete types
- Smaller interfaces = more composable ("accept io.Reader, not *os.File")

```go
// Bad — consumer hardwired to concrete type
func ProcessReport(r *CSVReader) error { ... }

// Good — consumer defines what it needs; any reader works
type LineReader interface {
    ReadLines() ([]string, error)
}
func ProcessReport(r LineReader) error { ... }
```

---

### 5. Apply the one-sentence test

If you cannot describe a unit's purpose in one sentence without using "and", it is not composable — it is two units fused together.

- "Fetches user data **and** sends a welcome email" → split
- "Validates the payment amount" → composable
- "Parses the request **and** writes the response" → split

Split at the "and": each half becomes its own composable unit.

## Common Mistakes

**Confusing composable with composed.** A class that uses composition internally (has-a relationships) is not necessarily composable from the outside. Composability is an external property — can consumers combine this unit freely? Internal structure is separate.

**Micro-splitting to absurdity.** Composability does not mean every function stands alone. Units that always appear together and have no independent use case should stay together. Split only when independent reuse exists or is imminent.

**Wide interfaces.** An interface with 8 methods is not composable — consumers must implement or mock all 8 to use it. Narrow interfaces (1–3 methods) are more composable.

## When NOT to Use

- **Prototypes and exploratory code** — composability adds surface area definition cost; defer until the design stabilizes
- **Framework entry points** — framework-mandated base classes or handlers are not composable by design; use them as-is
- **Performance-critical inner loops** — abstraction layers add overhead; measure before applying to hot paths
- **Very small codebases** — a 300-line script has no composability problem to solve
