---
name: apply-high-cohesion
description: Use when assigning responsibilities to classes — especially when a class is hard to name, does work unrelated to its core concept, or a single change requires understanding an entire class to avoid breaking something.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, high-cohesion, design-principles, oop, responsibility-assignment, developer, readability, maintainability]
related: [apply-low-coupling, apply-information-expert, apply-solid-principles, apply-unix-philosophy]
---

# Apply High Cohesion

Assign responsibilities so that each class remains focused on one strongly related set of concerns that can be understood and explained as a unit.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004); SOLID SRP (Martin, 2003); Unix philosophy; Go standard library design (small, focused interfaces). High cohesion is a foundational metric in software quality measurement — Chidamber & Kemerer (1994, IEEE TSE) established the Lack of Cohesion of Methods (LCOM) metric, now a standard quality gate in SonarQube, CodeClimate, and NDepend.
**Impact:** Chidamber & Kemerer (1994) showed that LCOM is inversely correlated with fault-proneness — low-cohesion classes have statistically higher defect density. Palomba et al. (2018, IEEE TSE) confirmed that God Classes (the archetypal low-cohesion pattern) have 3–4x higher change frequency than average classes. High-cohesion classes are self-documenting: the class name describes everything inside it.
**Why best:** Low cohesion grows through feature accretion — each new feature lands in the nearest existing class rather than the right class. The result is a God Class that knows everything and changes for every reason. High cohesion is not achieved by refusing to add features; it is achieved by asking "does this feature belong here, or in a focused class I haven't created yet?" The distinction between High Cohesion (GRASP) and SRP (SOLID) is axis: SRP asks "how many actors would cause this to change?"; High Cohesion asks "how related are the operations this class performs?"

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Chidamber & Kemerer, "A Metrics Suite for Object-Oriented Design" (IEEE TSE, 1994); Palomba et al., "Revisiting the Impact of Code Smells" (IEEE TSE, 2018)

## Steps

### 1. Apply the name test

A highly cohesive class can be named with one noun and a short description. If the name requires "And" or a vague term like "Manager", "Handler", or "Util", it is likely doing too much.

| Name | Signal |
|---|---|
| `OrderCalculator` | Cohesive — one focused operation |
| `UserManager` | Suspicious — "manage" hides multiple concerns |
| `OrderAndPaymentService` | Low cohesion — explicitly two concerns |
| `Utils` | No cohesion — a wastebasket class |

---

### 2. Measure LCOM: count unrelated method groups

List all instance methods. Draw a line between two methods if they share an instance variable. Disconnected groups signal low cohesion — the class contains multiple independent clusters that belong in separate classes.

```python
# Low cohesion — two disconnected groups: reporting and user management
class UserService:
    def find_user(self, id): return self.user_db.find(id)       # group A: uses user_db
    def update_email(self, id, email): ...                      # group A: uses user_db
    def generate_monthly_report(self): ...                      # group B: uses report_db
    def export_report_csv(self): ...                            # group B: uses report_db
    # A and B never share an instance variable — this is two classes

# High cohesion — split into two focused classes
class UserRepository:
    def find(self, id): ...
    def update_email(self, id, email): ...

class UserReportGenerator:
    def generate_monthly(self): ...
    def export_csv(self): ...
```

---

### 3. Assign new responsibilities to the most cohesion-preserving class

When adding a feature, evaluate each candidate class by asking: "does this responsibility relate to the existing operations in that class?"

- If yes and LCOM would not increase: assign it there
- If no or LCOM would increase: create a new focused class
- If yes but the class is already large: extract the related group first, then assign

---

### 4. Split God Classes by responsibility cluster

A God Class is a class that has grown to encompass unrelated responsibilities. Split it by identifying distinct clusters of methods and data, each of which belongs together.

```python
# God Class — OrderSystem does everything
class OrderSystem:
    def place_order(self): ...        # ordering
    def cancel_order(self): ...       # ordering
    def charge_payment(self): ...     # payment
    def refund_payment(self): ...     # payment
    def send_confirmation(self): ...  # notification
    def send_cancellation(self): ...  # notification
    def generate_invoice(self): ...   # invoicing

# Split into 4 cohesive classes
class OrderManager:     # ordering
class PaymentProcessor: # payment
class Notifier:         # notification
class InvoiceGenerator: # invoicing
```

---

### 5. Balance against low coupling

Splitting a class too aggressively raises coupling — each new class needs to know about others. Find the natural split: separate what has different rates of change or different collaborators, not what merely appears in different lines.

A class with 8 cohesive methods serving one concept is better than 4 classes of 2 methods each that all call each other.

## Common Mistakes

**Splitting by operation type instead of concept.** Moving all "create" methods to one class and all "update" methods to another produces structural cohesion but not conceptual cohesion. Group by domain concept, not by CRUD verb.

**Equating cohesion with size.** A small class can be incoherent (two unrelated methods). A large class can be highly cohesive (50 methods all serving one complex domain concept, like a geometry engine). Measure by relatedness, not line count.

**Confusing High Cohesion with SRP.** SRP asks: how many actors need this to change? High Cohesion asks: how related are the responsibilities? Both point at the same smell (bloated classes) but from different angles. Applying both gives the clearest guidance.

## When NOT to Use

- **Facade classes** — deliberately aggregate multiple subsystems behind a simplified interface; their broad responsibility is intentional
- **Entry-point / bootstrap classes** — application startup code necessarily wires many components together; this is coordination, not low cohesion
- **Value objects and data classes** — pure data structures with no behavior have no cohesion to measure
