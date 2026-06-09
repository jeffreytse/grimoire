---
name: apply-information-expert
description: Use when deciding which class should be responsible for a computation, validation, or operation — especially when the same data is accessed from multiple places to perform the same calculation, or when a class delegates all its work to another class that holds the actual data.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, information-expert, design-principles, oop, responsibility-assignment, developer, cohesion, maintainability]
related: [apply-creator-pattern, apply-high-cohesion, apply-low-coupling, apply-controller-pattern, apply-solid-principles]
---

# Apply Information Expert

Assign a responsibility to the class that has the information needed to fulfill it.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004) — the most frequently applied of the nine GRASP patterns, and the first principle to consider when assigning any responsibility. Consistent with the OOP encapsulation principle (Bjarne Stroustrup, "The C++ Programming Language") and the "Tell, Don't Ask" design heuristic (Martin Fowler, "Refactoring"). Applied in Google's internal style guides ("avoid anemic domain models").
**Impact:** When behavior lives where the data is, there is no need to expose data to external callers to perform operations on it. This reduces getter/setter proliferation, lowers coupling (callers don't need to know the internal structure), and raises cohesion (the class has a reason to exist beyond holding data). Anemic domain models — where all behavior is in service classes and domain objects are pure data — are the primary cause of "behavior scattered across the codebase," a pattern Martin Fowler identifies as a maintainability antipattern.
**Why best:** The alternative — placing behavior in a class that lacks the needed data — requires the class to retrieve data from multiple sources, creating unnecessary coupling. Information Expert is the default starting point: try assigning the responsibility to the class with the most relevant data first. Override with Low Coupling or High Cohesion only when the default assignment produces an obvious violation of those principles.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Fowler, "Refactoring" 2nd ed. (Addison-Wesley, 2018); Stroustrup, "The C++ Programming Language" (Addison-Wesley, 2013)

## Steps

### 1. Identify what information the responsibility requires

Before assigning a responsibility, enumerate the data it needs. The class that owns the most of that data is the Information Expert.

Example: computing an order total requires line items, quantities, and prices. `Order` holds all of these → `Order` is the expert → `Order.total()` is the correct assignment.

```python
# Bad — OrderController reaches into Order's data to compute the total
class OrderController:
    def get_total(self, order):
        return sum(item.price * item.quantity for item in order.line_items)
        # ^ accessing order's data from outside to do order's job

# Good — Order has the data; Order computes the total
class Order:
    def total(self):
        return sum(item.price * item.quantity for item in self.line_items)
```

---

### 2. Move scattered calculations back to their data owner

When the same data is accessed from multiple callers to perform the same computation, that computation belongs on the data owner.

```python
# Bad — discount calculation duplicated across multiple callers
# In OrderController:
discount = order.membership_tier == "gold" and order.total > 100
# In InvoiceService:
discount = order.membership_tier == "gold" and order.total > 100
# In ReportGenerator:
discount = order.membership_tier == "gold" and order.total > 100

# Good — Order knows its own discount eligibility
class Order:
    def eligible_for_discount(self):
        return self.membership_tier == "gold" and self.total() > 100
```

---

### 3. Apply "tell, don't ask"

Asking an object for its data so you can operate on it violates Information Expert. Telling an object to perform an operation respects it.

| Ask (violates) | Tell (respects) |
|---|---|
| `if account.get_balance() < 0: account.freeze()` | `account.freeze_if_overdrawn()` |
| `tax = order.get_subtotal() * tax_rate` | `tax = order.calculate_tax(tax_rate)` |
| `if user.get_role() == "admin": ...` | `if user.is_admin(): ...` |

---

### 4. Recognize when to override: cohesion and coupling conflicts

Information Expert is the first principle to apply, not the only one. Override it when:

- **High Cohesion conflict:** The expert class already has too many responsibilities. Assigning this one would bloat it beyond its natural purpose.
- **Low Coupling conflict:** The expert class would need a new, unrelated dependency to fulfill the responsibility. A Pure Fabrication is better.

```python
# Override scenario: Order is the expert for payment processing (has amount, currency)
# BUT adding payment logic to Order couples it to PaymentGateway (external service)
# LOW COUPLING overrides INFORMATION EXPERT here → create PaymentService

class PaymentService:   # Pure Fabrication — no domain concept, but cohesion/coupling justified
    def charge(self, order, gateway):
        return gateway.charge(order.total(), order.currency)
```

---

### 5. Avoid the anemic domain model

An anemic domain model has domain objects (Order, Customer, Product) with only getters/setters and no behavior. All behavior sits in service classes that reach into these objects.

Signs of an anemic model:
- Domain classes have only `get_*` and `set_*` methods
- Service classes are 500+ lines of procedural code operating on domain objects
- Changing a business rule requires editing 3–5 service classes

Fix: push behavior back to the domain class that owns the data. Start with the most duplicated calculations.

## Common Mistakes

**Applying Information Expert to infrastructure concerns.** Database access, HTTP calls, and file I/O should not live on domain classes even if those classes have the relevant data. Use the Pure Fabrication pattern for infrastructure — keep domain objects free of external dependencies.

**Confusing "has the data" with "knows about the data."** If Order knows about Customer and Customer knows the shipping address, Order is not the expert for address formatting — Customer is. Follow the chain to find the true owner.

**Over-applying against cohesion.** Putting every calculation related to `User` on the `User` class produces a 2000-line God Class. Use High Cohesion to constrain Information Expert: the expert assignment is correct until the class becomes incoherent.

## When NOT to Use

- **Cross-cutting concerns** (logging, metrics, caching): no single domain class owns these; they belong in infrastructure layers
- **Reporting and analytics**: queries that aggregate data from many sources belong in dedicated read models or query services, not on domain objects
- **When the expert would gain an inappropriate dependency**: a domain class should not import an email library just because it has the data needed to compose a notification
