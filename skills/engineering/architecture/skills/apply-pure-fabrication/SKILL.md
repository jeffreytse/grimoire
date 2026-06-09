---
name: apply-pure-fabrication
description: Use when no domain class is a suitable home for a responsibility — especially when assigning a responsibility to a domain object would violate high cohesion or low coupling by introducing infrastructure concerns (database, email, logging) into domain objects.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, pure-fabrication, design-principles, oop, responsibility-assignment, developer, separation-of-concerns, low-coupling]
related: [apply-information-expert, apply-low-coupling, apply-high-cohesion, apply-indirection, apply-domain-based-naming]
---

# Apply Pure Fabrication

When assigning a responsibility to a domain class would violate High Cohesion or Low Coupling, create an artificial class — one with no domain counterpart — solely to achieve a good design.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004); the design justification behind Service classes in Domain-Driven Design (Evans, 2003), Repository classes, and all infrastructure-layer objects. The Repository pattern (Fowler, "Patterns of Enterprise Application Architecture", 2002) is a canonical Pure Fabrication: no business domain has a "repository" — it is an invented construct to separate persistence from domain logic.
**Impact:** Without Pure Fabrication, the Information Expert principle leads to domain objects absorbing infrastructure responsibilities: `Order` would handle its own database persistence, `Customer` would send its own emails, `Invoice` would generate its own PDFs. Domain objects become coupled to databases, email servers, and file systems — impossible to test without infrastructure setup. Pure Fabrication is the escape valve that keeps domain objects clean.
**Why best:** The tension in GRASP is between Information Expert (put behavior with data) and Low Coupling/High Cohesion (keep classes focused and dependency-light). Pure Fabrication resolves this tension by creating a class specifically designed to hold the infrastructure responsibility. The domain class stays clean; the Pure Fabrication absorbs the infrastructure coupling. This is intentional design, not a workaround.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Evans, "Domain-Driven Design" (Addison-Wesley, 2003); Fowler, "Patterns of Enterprise Application Architecture" (Addison-Wesley, 2002)

## Steps

### 1. Recognize when Information Expert would violate cohesion or coupling

The trigger for Pure Fabrication is a conflict between Information Expert (put the responsibility with the data) and High Cohesion / Low Coupling.

```python
# Information Expert says: Order has the data for persistence → Order should persist itself
class Order:
    def save(self):
        db.execute("INSERT INTO orders ...", self.id, self.total, ...)   # Infrastructure in domain
        # Problem: Order now depends on DB driver; can't test without database
        # Problem: Order is responsible for both business rules AND persistence

# Pure Fabrication resolves the conflict
class OrderRepository:   # Invented class — no "repository" in business domain
    def save(self, order: Order) -> None:
        db.execute("INSERT INTO orders ...", order.id, order.total, ...)
```

---

### 2. Name Pure Fabrications clearly

A Pure Fabrication has no domain counterpart, so its name describes its technical role, not a business concept. This is one of the exceptions to `apply-domain-based-naming`.

Common Pure Fabrication names and their responsibilities:

| Name | Responsibility |
|---|---|
| `XRepository` | Persistence operations for domain type X |
| `XService` | Domain operations that span multiple types |
| `XFactory` | Complex creation of domain objects |
| `XMapper` / `XAssembler` | Translation between layers (domain ↔ DTO) |
| `XValidator` | Validation logic separated from domain object |
| `XNotifier` / `XEmailer` | Notification / messaging concerns |
| `XAdapter` | Wrapping an external API in a domain interface |

---

### 3. Keep domain objects free of infrastructure

The rule of thumb: domain objects should not import or instantiate infrastructure classes (databases, email clients, HTTP clients, file systems).

```python
# Bad — Customer knows about email; couples domain to email infrastructure
class Customer:
    def send_welcome_email(self):
        smtp = smtplib.SMTP("mail.example.com")   # infrastructure in domain
        smtp.sendmail(...)

# Good — WelcomeEmailer is a Pure Fabrication; Customer knows nothing about email
class WelcomeEmailer:
    def __init__(self, smtp_client):
        self.smtp = smtp_client

    def send(self, customer: Customer) -> None:
        self.smtp.sendmail(
            to=customer.email,
            subject="Welcome!",
            body=f"Hello {customer.name}..."
        )
```

---

### 4. Use Pure Fabrication for cross-cutting algorithms

When an algorithm spans multiple domain objects but doesn't belong on any one of them, a Pure Fabrication holds it.

```python
# Algorithm: calculate shipping cost — depends on Order, Customer, and external rate API
# Not cohesive on Order (knows too much about Customer and rates)
# Not cohesive on Customer (has nothing to do with Customer's core concerns)

class ShippingCalculator:   # Pure Fabrication — no domain concept, but cohesion justified
    def __init__(self, rate_api: ShippingRateAPI):
        self.rate_api = rate_api

    def calculate(self, order: Order, customer: Customer) -> Money:
        zone = self.rate_api.get_zone(customer.shipping_address)
        return self.rate_api.get_rate(zone, order.weight)
```

---

### 5. Keep Pure Fabrications cohesive

A Pure Fabrication that absorbs too many responsibilities becomes a service layer God Object — the same problem it was meant to prevent. Apply High Cohesion to Pure Fabrications as strictly as to domain objects.

```python
# Bad — OrderService is a dumping ground
class OrderService:
    def place_order(self): ...
    def cancel_order(self): ...
    def generate_invoice(self): ...      # different concern
    def send_confirmation_email(self): ... # different concern
    def calculate_shipping(self): ...    # different concern

# Good — each Pure Fabrication has one clear responsibility
class OrderPlacementService: ...
class InvoiceGenerator: ...
class OrderNotifier: ...
class ShippingCalculator: ...
```

## Common Mistakes

**Pure Fabrication as a wastebasket.** A class named `OrderUtils`, `OrderHelper`, or `OrderManager` with no clear purpose is not a Pure Fabrication — it is a junk drawer. Pure Fabrications have a single, named technical responsibility.

**Abusing Pure Fabrication to justify anemic domain models.** If ALL behavior is in Pure Fabrications and domain objects have only getters/setters, the balance is wrong. Information Expert should be the first choice; Pure Fabrication is the override when cohesion or coupling is genuinely violated.

**Confusing Pure Fabrication with Domain Service.** A Domain Service (DDD terminology) handles domain logic that doesn't naturally fit on one entity. A Pure Fabrication handles infrastructure concerns. Both are invented classes, but they serve different purposes and belong in different layers.

## When NOT to Use

- **When a domain class can absorb the responsibility without violating cohesion or low coupling** — Information Expert is still the right choice; Pure Fabrication is only for overrides
- **Simple CRUD applications** — in applications where all operations are data transformations with no domain logic, the distinction between domain and Pure Fabrication collapses; just write focused service/handler classes
