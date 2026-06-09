---
name: apply-indirection
description: Use when two components are directly coupled and that coupling causes problems — especially when changing one component forces changes in the other, or when a component cannot be tested without its direct dependency.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, indirection, design-principles, oop, decoupling, developer, testability, maintainability]
related: [apply-low-coupling, apply-protected-variations, apply-pure-fabrication, apply-solid-principles]
---

# Apply Indirection

Assign responsibility to an intermediate object to decouple two components that would otherwise depend directly on each other.

## Why This Is Best Practice

**Adopted by:** The mechanism behind the majority of Gang of Four design patterns (Adapter, Facade, Proxy, Mediator, Broker, Bridge — all are applications of Indirection). David Wheeler's maxim — "Any problem in computer science can be solved by adding another layer of indirection" — is one of the most cited principles in software engineering. Dependency injection containers, service meshes, and API gateways are all industrial-scale applications of Indirection.
**Impact:** Direct coupling between two volatile components means every change in one propagates to the other. An intermediate object isolates each side: component A only knows the intermediary's interface, and component B only knows the intermediary's interface. Changes on either side are absorbed by the intermediary. This is the mechanism that makes microservices independently deployable — the service mesh is the intermediary between every pair of services.
**Why best:** Indirection is the general principle from which specific patterns emerge. When two classes are directly coupled and that coupling is causing problems, the solution is almost always to introduce an intermediary that owns the relationship. The choice of which intermediary pattern to use (Adapter, Facade, Proxy, Mediator) is determined by what kind of coupling is being broken.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Gamma et al., "Design Patterns" (Addison-Wesley, 1994); Wheeler, "Computer Programming and Computer Systems" (Academic Press, 1963)

## Steps

### 1. Identify the direct coupling that is causing a problem

Not all direct coupling needs indirection. Add an intermediary only when the direct coupling creates a concrete problem.

**Signals that indirection is warranted:**
- Changing component B requires editing component A (change coupling)
- Testing component A requires setting up component B (test coupling)
- Component A cannot be reused in a context where component B is unavailable
- Component B is third-party and its API is unstable

```python
# Direct coupling — OrderService imports and instantiates Stripe directly
from stripe import StripeClient

class OrderService:
    def charge(self, order):
        client = StripeClient(api_key=STRIPE_KEY)
        return client.charges.create(amount=order.total, currency="usd")
        # Problem: can't test without network; can't swap payment processor
```

---

### 2. Introduce an intermediary interface

Define an interface that captures what component A needs from component B. Component A depends on the interface; component B implements it. The interface is the intermediary.

```python
from abc import ABC, abstractmethod

# Intermediary: PaymentGateway interface
class PaymentGateway(ABC):
    @abstractmethod
    def charge(self, amount: float, currency: str) -> ChargeResult: ...

# B implements the interface
class StripeGateway(PaymentGateway):
    def charge(self, amount, currency):
        return StripeClient(STRIPE_KEY).charges.create(amount=amount, currency=currency)

# A depends on the interface only; knows nothing about Stripe
class OrderService:
    def __init__(self, gateway: PaymentGateway):
        self.gateway = gateway

    def charge(self, order):
        return self.gateway.charge(order.total, order.currency)

# Test injects a fake — no network needed
class FakeGateway(PaymentGateway):
    def charge(self, amount, currency):
        return ChargeResult(status="ok", charge_id="test-123")
```

---

### 3. Choose the intermediary pattern that fits the coupling problem

| Problem | Intermediary pattern |
|---|---|
| Interface mismatch (API shapes differ) | Adapter |
| Complex subsystem with many classes | Facade |
| Adding behavior to object without modifying it | Proxy / Decorator |
| Many-to-many communication between components | Mediator |
| Abstracting remote/local call difference | Proxy |
| Shielding against interface volatility | Protected Variations |

---

### 4. Apply Indirection to reverse a dependency direction

When a high-level module depends on a low-level module (violating SOLID DIP), introduce an interface owned by the high-level module. The low-level module implements it. Dependency direction reverses.

```python
# Bad — business logic (high-level) depends on MySQL (low-level)
class OrderService:
    def find(self, id):
        return MySQLConnection().query("SELECT * FROM orders WHERE id = ?", id)

# Good — interface owned by business layer; MySQLRepository is a low-level detail
class OrderRepository(ABC):       # lives in business layer
    @abstractmethod
    def find(self, id) -> Order: ...

class MySQLOrderRepository(OrderRepository):   # lives in infrastructure layer
    def find(self, id):
        return MySQLConnection().query("SELECT * FROM orders WHERE id = ?", id)

class OrderService:
    def __init__(self, repo: OrderRepository): ...  # high-level depends on abstraction
```

---

### 5. Limit indirection depth

Every additional layer of indirection adds cognitive load. Two or three layers between a caller and the actual behavior is the practical limit before the code becomes harder to trace than the original coupling.

Rule: each layer of indirection must solve a specific, named coupling problem. If you cannot articulate what problem a layer solves, remove it.

## Common Mistakes

**Indirection for its own sake.** Wrapping a stable, never-changing class in an interface just because "we might swap it later" adds complexity without reducing coupling. Apply indirection when coupling is causing a current problem or a highly likely future problem.

**The wrong intermediary.** Introducing a Mediator when an Adapter is needed, or a Facade when a Proxy is needed, creates the right structure but solves the wrong problem. Identify the coupling type first (Step 3), then choose the pattern.

**Interface explosion.** Creating an interface for every class produces a 1:1 mapping that provides no actual decoupling — callers still know exactly which implementation they will get via the interface name. Interfaces should represent roles, not types.

## When NOT to Use

- **Stable third-party libraries** — wrapping `json.loads()` in an interface adds no value; the library's API is stable and the implementation will never be swapped
- **Value objects and data structures** — no behavior to decouple; pass them directly
- **Performance-critical paths** — every layer of indirection has dispatch overhead; measure before applying in tight loops
- **Small codebases** — the overhead of interfaces and intermediaries pays off in large codebases; in a 500-line service, direct calls are clearer
