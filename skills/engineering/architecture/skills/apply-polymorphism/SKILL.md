---
name: apply-polymorphism
description: Use when behavior varies based on type — especially when if/switch statements check the type of an object to decide what to do, or when adding a new variant requires editing an existing class.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, polymorphism, design-principles, oop, responsibility-assignment, developer, open-closed, extensibility]
related: [apply-protected-variations, apply-solid-principles, apply-composition-over-inheritance]
---

# Apply Polymorphism

When behavior varies by type, assign the responsibility for that variation to the types themselves using polymorphic operations — not to a caller that checks the type.

## Why This Is Best Practice

**Adopted by:** Core OOP principle since Simula 67 (1967) and Smalltalk (1972); formalized in Gang of Four "Design Patterns" (1994) — 12 of the 23 patterns are applications of this principle (Strategy, Command, State, Visitor, Observer, etc.). Directly enables SOLID's Open/Closed Principle. Every modern OOP language enforces this via virtual methods, interfaces, and abstract classes.
**Impact:** Type-checking code (`if type == "A"`, `switch(type)`) is a closed set: adding a new type always requires editing the conditional. Polymorphic code is an open set: adding a new type means adding a new class that implements the interface — existing code requires no edits. Google's internal refactoring guidelines cite conditional type checks as one of the top patterns targeted for polymorphic replacement.
**Why best:** Type-checking callers violate the Open/Closed Principle (SOLID O) and create coupling between the caller and every concrete type it checks. Polymorphism moves the type-specific behavior to the type itself — the only class that needs to change when a type's behavior changes is that type's class. This is the mechanism behind the Strategy, State, Command, and Template Method patterns.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Gamma et al., "Design Patterns" (Addison-Wesley, 1994); Martin, "Agile Software Development" (Prentice Hall, 2003)

## Steps

### 1. Identify the type-checking smell

Any code that branches on an object's type is a candidate for polymorphic replacement.

```python
# Type-checking smell — every new payment type requires editing this function
def process_payment(payment):
    if payment.type == "credit_card":
        stripe.charge(payment.amount, payment.card_number)
    elif payment.type == "bank_transfer":
        ach.transfer(payment.amount, payment.account_number)
    elif payment.type == "crypto":
        blockchain.send(payment.amount, payment.wallet)
    # Adding "gift_card" requires editing this function
```

---

### 2. Define a polymorphic interface

Create an interface (or abstract class) that declares the operation. Each type gets its own class that implements the operation for that type's behavior.

```python
from abc import ABC, abstractmethod

class PaymentMethod(ABC):
    @abstractmethod
    def process(self, amount: float) -> Receipt: ...

class CreditCardPayment(PaymentMethod):
    def process(self, amount):
        return stripe.charge(amount, self.card_number)

class BankTransferPayment(PaymentMethod):
    def process(self, amount):
        return ach.transfer(amount, self.account_number)

class CryptoPayment(PaymentMethod):
    def process(self, amount):
        return blockchain.send(amount, self.wallet)

# Caller is closed to modification — no type checks
def process_payment(payment: PaymentMethod, amount: float):
    return payment.process(amount)
```

---

### 3. Migrate type-checking conditionals incrementally

For existing type-checking code, replace one branch at a time:

1. Introduce the interface
2. Move the first branch's behavior to a concrete class
3. Replace that branch with a polymorphic call
4. Repeat for remaining branches
5. Remove the conditional

```python
# Step 1: introduce interface, keep conditional for now
class PaymentMethod(ABC):
    @abstractmethod
    def process(self, amount): ...

# Step 2-3: migrate credit card, use polymorphic call for it
class CreditCardPayment(PaymentMethod):
    def process(self, amount):
        return stripe.charge(amount, self.card_number)

def process_payment(payment):
    if isinstance(payment, CreditCardPayment):   # now polymorphic for this type
        return payment.process(payment.amount)
    elif payment.type == "bank_transfer":        # still transitional
        ach.transfer(...)
    # continue migrating...
```

---

### 4. Use polymorphism for state variation

When an object's behavior changes based on its internal state, use the State pattern — each state is a class implementing the state interface.

```python
# Bad — Order behavior varies by status string; adding a status requires editing everywhere
class Order:
    def can_cancel(self):
        if self.status == "pending": return True
        if self.status == "shipped": return False
        if self.status == "delivered": return False
        if self.status == "refunded": return False

# Good — each state encapsulates its own behavior
class OrderState(ABC):
    @abstractmethod
    def can_cancel(self) -> bool: ...

class PendingState(OrderState):
    def can_cancel(self): return True

class ShippedState(OrderState):
    def can_cancel(self): return False

class Order:
    def can_cancel(self):
        return self.state.can_cancel()   # delegates to state; no type checking
```

---

### 5. Validate at system boundaries, not with type checks

The one acceptable place for type-based conditionals is at system entry points where raw data (JSON, form input) must be converted to typed domain objects. The conditional lives at the boundary; once inside the system, everything is polymorphic.

```python
# Acceptable — type conditional at deserialization boundary only
def deserialize_payment(data: dict) -> PaymentMethod:
    match data["type"]:
        case "credit_card":   return CreditCardPayment(**data)
        case "bank_transfer": return BankTransferPayment(**data)
        case "crypto":        return CryptoPayment(**data)
        case _:               raise ValueError(f"Unknown type: {data['type']}")

# After this point: all code is polymorphic, no type checks
```

## Common Mistakes

**Using `isinstance()` as polymorphism.** Adding `isinstance(obj, ConcreteClass)` guards before calling type-specific methods is still type-checking, not polymorphism. Every `isinstance` check is a candidate for a polymorphic interface.

**Polymorphism for trivial variation.** Not every boolean variation needs a class hierarchy. `if is_debug: log()` does not need a `DebugMode` class. Apply polymorphism when behavior is substantial and types are extensible.

**Deep inheritance hierarchies.** A polymorphic type hierarchy 4+ levels deep is harder to follow than a type check. Keep hierarchies shallow (1–2 levels). Prefer interfaces over inheritance for behavior variation.

## When NOT to Use

- **Finite, never-extending type sets** — if there will always and only be exactly two cases (`True`/`False`, `enabled`/`disabled`) and no third case is possible, a simple conditional is clearer
- **Data dispatch at serialization boundaries** — type-based routing from raw input to typed objects is acceptable at system boundaries (see Step 5)
- **Performance-critical paths** — virtual dispatch has a small overhead compared to direct calls; in tight loops, measure before introducing polymorphism
