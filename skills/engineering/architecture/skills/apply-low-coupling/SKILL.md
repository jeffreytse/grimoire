---
name: apply-low-coupling
description: Use when assigning responsibilities to classes — especially when a change to one class ripples into many others, classes are hard to test in isolation, or reusing a class requires dragging in unrelated dependencies.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, low-coupling, design-principles, oop, responsibility-assignment, developer, testability, maintainability]
related: [apply-high-cohesion, apply-information-expert, apply-law-of-demeter, apply-composition-over-inheritance, apply-solid-principles]
---

# Apply Low Coupling

Assign responsibilities so that classes depend on as few other classes as possible, and only on classes unlikely to change.

## Why This Is Best Practice

**Adopted by:** Foundational principle across all major OOP design frameworks — SOLID (Dependency Inversion), Gang of Four ("program to an interface, not an implementation"), GRASP (Larman, 2004). Endorsed by Google (Engineering Practices), Microsoft (.NET Design Guidelines), and ThoughtWorks Technology Radar. Low coupling is the prerequisite for independent deployment in microservices and for test isolation in unit testing.
**Impact:** Yamashita & Moonen (ICSM 2012) measured 74 Java projects and found coupling between objects (CBO) to be the strongest predictor of maintainability degradation — high CBO classes required 2–3x more effort to change than low-CBO counterparts. Google's SRE documentation identifies tight coupling as the primary cause of cascading failures in distributed systems. Highly coupled classes cannot be tested without their dependencies, tripling test setup cost.
**Why best:** Coupling is unavoidable — classes must collaborate. The goal is not zero coupling but *appropriate* coupling: depend on abstractions over concretions, on stable elements over volatile ones, on fewer elements over many. The Law of Demeter is a specific coupling rule; Low Coupling is the general principle that motivates it and all structural decoupling techniques.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Yamashita & Moonen, "Do Code Smells Reflect Important Maintainability Aspects?" (ICSM 2012); Gamma et al., "Design Patterns" (Addison-Wesley, 1994)

## Steps

### 1. Measure coupling: count afferent and efferent dependencies

Before assigning a responsibility, count how many classes the receiving class already depends on.

- **Efferent coupling (Ce):** how many classes this class depends on (outgoing)
- **Afferent coupling (Ca):** how many classes depend on this class (incoming)

A class with Ce > 10 is a coupling hotspot — adding more responsibilities here multiplies the problem. Assign the responsibility elsewhere.

```python
# High coupling — OrderProcessor depends on 6 classes
class OrderProcessor:
    def process(self, order):
        user = UserRepository().find(order.user_id)       # dep 1
        payment = PaymentGateway().charge(order)          # dep 2
        inventory = InventoryService().reserve(order)     # dep 3
        notification = NotificationService().send(user)   # dep 4
        audit = AuditLog().record(order)                  # dep 5
        report = ReportingService().log(order)            # dep 6
```

---

### 2. Depend on abstractions, not concretions

Coupling to a concrete class binds to its implementation details. Coupling to an interface binds only to the contract.

```python
# Bad — coupled to concrete MySQLDatabase; changing DB requires rewriting PaymentService
class PaymentService:
    def __init__(self):
        self.db = MySQLDatabase()    # concrete

# Good — coupled to an interface; any store implementation works
class PaymentService:
    def __init__(self, store: PaymentStore):   # abstract
        self.store = store
```

---

### 3. Assign responsibility to the class that already knows the most

Giving a responsibility to a class that already holds the needed data avoids creating a new dependency. This is Information Expert applied as a coupling reduction tool.

```python
# Bad — OrderController must know both Order and DiscountPolicy to compute total
class OrderController:
    def get_total(self, order, policy):
        return order.subtotal * (1 - policy.rate_for(order))   # knows too much

# Good — Order knows its own total; controller has no coupling to DiscountPolicy
class Order:
    def total(self, policy):
        return self.subtotal * (1 - policy.rate_for(self))
```

---

### 4. Prefer unidirectional dependencies

Bidirectional dependencies (A → B and B → A) are the most fragile coupling. Break cycles by extracting an interface or introducing an intermediary.

```python
# Bad — bidirectional: Order knows Customer, Customer knows Order
class Order:
    def get_customer(self): return self.customer

class Customer:
    def get_orders(self): return self.orders    # both know each other

# Good — unidirectional via a query service; Customer does not know Order
class Order:
    def get_customer(self): return self.customer

class OrderQuery:
    def orders_for(self, customer_id): ...     # third party holds the relationship
```

---

### 5. Prefer stable dependencies over volatile ones

Coupling to a class that changes frequently means changing alongside it. Depend on the things least likely to change.

| More stable (prefer) | Less stable (avoid coupling to) |
|---|---|
| Interfaces and abstractions | Concrete implementations |
| Domain types | Infrastructure types |
| Standard library | Internal utility classes |
| Third-party stable APIs | Internal business rules |

## Common Mistakes

**Coupling to reduce code duplication.** Pulling a dependency into a class just to share a utility method is a coupling mistake. Extract the utility to a shared abstraction instead — both classes depend on the stable abstraction, not on each other.

**Confusing low coupling with no coupling.** A class with zero dependencies is either a value object or a god object hiding its coupling. Appropriate coupling means each dependency is necessary, stable, and abstract.

**Indirection for its own sake.** Adding a layer of indirection without a real coupling reduction does not lower coupling — it just moves it. Each layer of indirection must reduce actual dependencies in the dependency graph.

## When NOT to Use

- **Simple scripts and one-off tools** — coupling overhead is not justified for code with no reuse requirement
- **Framework base classes** — frameworks mandate coupling to their types; work within the model
- **Data classes / value objects** — pure data holders with no behavior have no coupling problem to solve
- **When cohesion suffers** — sometimes reducing coupling requires splitting a class so much that cohesion collapses; balance Low Coupling against High Cohesion (they are complementary GRASP patterns, evaluated together)
