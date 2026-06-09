---
name: apply-state-pattern
description: Use when an object's behavior changes based on its internal state and the logic for each state is complex — replacing large conditional chains with state objects where each state encapsulates its own behavior.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 305–313; TCP connection state machine; finite state machines in game AI; XState (JavaScript statechart library); vending machine controllers
tags: [design-patterns, behavioral, state, oop, developer, state-machine, conditional-elimination]
related: [apply-strategy-pattern, apply-command-pattern, apply-solid-principles]
---

# Apply State Pattern

Allow an object to alter its behavior when its internal state changes — by delegating behavior to state objects rather than using conditional logic.

## Why This Is Best Practice

**Adopted by:** XState (JavaScript statechart library — 1.5M weekly npm downloads)
implements State pattern as its core model. TCP protocol stack implements state as
objects (CLOSED, LISTEN, SYN_SENT, ESTABLISHED, etc. — each with its own behavior for
each operation). Unity's Animator Controller is a visual state machine based on the
same pattern. Every parser generator (ANTLR, PLY) produces state-based recognizers.
**Impact:** GoF documents that without the State pattern, a class with 4 states and
6 operations has 4×6=24 conditional branches scattered across all methods. The State
pattern consolidates each state's behavior into one class — 4 state classes of 6 methods
each. Adding a 5th state requires one new class, not editing 6 existing methods.
**Why best:** `if/elif` state dispatch — the alternative — scatters state logic across
every method, making it hard to add states or understand what happens in a given state.
State objects make each state's behavior cohesive, isolated, and independently testable.

Sources: Gamma et al. (1994) pp. 305–313; XState documentation; RFC 793 (TCP state machine)

## Steps

### Step 1: Define the state interface with all state-dependent methods

```python
from abc import ABC, abstractmethod

class OrderState(ABC):
    @abstractmethod
    def pay(self, order: "Order") -> None: ...

    @abstractmethod
    def ship(self, order: "Order") -> None: ...

    @abstractmethod
    def cancel(self, order: "Order") -> None: ...
```

### Step 2: Implement concrete state classes — each handles its own transitions

```python
class PendingState(OrderState):
    def pay(self, order: "Order") -> None:
        print("Payment accepted")
        order.state = PaidState()

    def ship(self, order: "Order") -> None:
        print("Cannot ship — order not paid")

    def cancel(self, order: "Order") -> None:
        print("Order cancelled")
        order.state = CancelledState()

class PaidState(OrderState):
    def pay(self, order: "Order") -> None:
        print("Already paid")

    def ship(self, order: "Order") -> None:
        print("Order shipped")
        order.state = ShippedState()

    def cancel(self, order: "Order") -> None:
        print("Cancelling paid order — refund issued")
        order.state = CancelledState()

class ShippedState(OrderState):
    def pay(self, order: "Order") -> None:
        print("Already paid and shipped")

    def ship(self, order: "Order") -> None:
        print("Already shipped")

    def cancel(self, order: "Order") -> None:
        print("Cannot cancel — already shipped")

class CancelledState(OrderState):
    def pay(self, order: "Order") -> None:
        print("Order is cancelled")

    def ship(self, order: "Order") -> None:
        print("Order is cancelled")

    def cancel(self, order: "Order") -> None:
        print("Already cancelled")
```

### Step 3: Context holds the current state and delegates all state-dependent operations

```python
class Order:
    def __init__(self):
        self.state: OrderState = PendingState()

    def pay(self) -> None:
        self.state.pay(self)

    def ship(self) -> None:
        self.state.ship(self)

    def cancel(self) -> None:
        self.state.cancel(self)
```

### Step 4: Transition state inside state objects — not in the context

State transitions (`order.state = PaidState()`) happen inside state methods (Step 2),
not inside `Order`. This keeps the context thin and each state's transitions explicit.

### Step 5: Verify transitions by testing each state in isolation

```python
def test_pending_pay():
    order = Order()
    assert isinstance(order.state, PendingState)
    order.pay()
    assert isinstance(order.state, PaidState)

def test_shipped_cannot_cancel():
    order = Order()
    order.pay()
    order.ship()
    order.cancel()
    assert isinstance(order.state, ShippedState)   # state unchanged
```

## When NOT to Use

- **Simple states with few operations** — 2 states × 2 operations = 4 conditions. A simple `if/else` is cleaner than 2 state classes.
- **When transitions are purely data-driven** — if state transitions come from a config file or database, a data-driven state machine (a transition table) is more maintainable than code-based state objects.

## Common Mistakes

**Context that checks `isinstance(state, ...)` to make decisions.** This reintroduces the conditional logic the pattern eliminates. Delegate the decision to the state object.

**State objects that modify context fields directly.** If `PaidState` sets `order.paid_at = datetime.now()`, the state is coupled to the context's internals. Pass only what state objects need, or use events.

**Missing states in the state interface.** If one state handles `ship()` but another doesn't implement the interface method, calls fail unexpectedly. All states must implement all interface methods.
