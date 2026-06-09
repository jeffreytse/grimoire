---
name: apply-observer-pattern
description: Use when a change in one object must automatically notify and update an open-ended set of dependent objects — without the subject knowing which specific objects depend on it.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 293–303; Java EventListener; Python asyncio events; React useState/useEffect; RxJS (Reactive Extensions)
tags: [design-patterns, behavioral, observer, oop, developer, event-driven, publish-subscribe, loose-coupling]
related: [apply-mediator-pattern, apply-event-driven-architecture, apply-solid-principles]
---

# Apply Observer Pattern

Define a one-to-many dependency between objects so that when one changes state, all its dependents are notified and updated automatically.

## Why This Is Best Practice

**Adopted by:** Java's event model (every `EventListener` in AWT, Swing, and JavaFX is
an Observer — the foundation of all Java UI), Python's `asyncio` event system, React's
hooks model (`useState`/`useEffect` — 100M+ npm downloads/week), RxJS (Reactive
Extensions — the core of Angular, used in millions of applications), and all
publish-subscribe message brokers (Redis Pub/Sub, MQTT).
**Impact:** GoF documents that without Observer, the subject must know all its
dependents explicitly — coupling that scales as O(M×N) with M subjects and N observers.
React's adoption of the Observer model for state updates is cited as the primary
reason UI components can re-render independently on state change without centralized
routing logic.
**Why best:** Polling — the alternative — requires every dependent to periodically check
the subject for changes. This wastes CPU and introduces latency. Direct coupling
(subject calls each dependent's specific method) requires the subject to know all
dependents. Observer decouples both: push notification, open-ended subscriber list.

Sources: Gamma et al. (1994) pp. 293–303; React hooks documentation; RxJS documentation

## Steps

### Step 1: Define the observer interface

```python
from abc import ABC, abstractmethod

class Observer(ABC):
    @abstractmethod
    def update(self, subject: "Subject") -> None: ...
```

### Step 2: Define the subject interface with observer management

```python
class Subject(ABC):
    def __init__(self):
        self._observers: list[Observer] = []

    def attach(self, observer: Observer) -> None:
        self._observers.append(observer)

    def detach(self, observer: Observer) -> None:
        self._observers.remove(observer)

    def notify(self) -> None:
        for observer in list(self._observers):   # copy to allow detach during notify
            observer.update(self)
```

### Step 3: Implement the concrete subject — notifies on state change

```python
class StockPrice(Subject):
    def __init__(self, ticker: str, price: float):
        super().__init__()
        self._ticker = ticker
        self._price = price

    @property
    def ticker(self) -> str:
        return self._ticker

    @property
    def price(self) -> float:
        return self._price

    @price.setter
    def price(self, value: float) -> None:
        self._price = value
        self.notify()   # notify all observers on every price change
```

### Step 4: Implement concrete observers

```python
class PriceAlert(Observer):
    def __init__(self, threshold: float):
        self._threshold = threshold

    def update(self, subject: StockPrice) -> None:
        if subject.price > self._threshold:
            print(f"ALERT: {subject.ticker} hit ${subject.price:.2f} (threshold: ${self._threshold})")

class PriceLogger(Observer):
    def update(self, subject: StockPrice) -> None:
        print(f"LOG: {subject.ticker} = ${subject.price:.2f}")
```

### Step 5: Attach observers and trigger updates

```python
aapl = StockPrice("AAPL", 150.0)
aapl.attach(PriceAlert(threshold=200.0))
aapl.attach(PriceLogger())

aapl.price = 180.0   # PriceLogger notified
aapl.price = 205.0   # both PriceAlert and PriceLogger notified
```

## When NOT to Use

- **When the update order matters** — Observer doesn't guarantee notification order. If order is critical, use a queue or coordinator instead.
- **When cascading updates cause loops** — if observer A's `update()` changes subject B which notifies observer C which changes subject A, you have an infinite loop. Map out update graphs before applying.
- **Simple, one-to-one dependencies** — a direct method call is cleaner when exactly one object cares about the change.

## Common Mistakes

**Memory leaks from undetached observers.** If a transient observer (a short-lived widget) attaches and is never detached, the subject holds a reference preventing garbage collection. Always `detach()` observers when they're done.

**Modifying the observer list during `notify()`.** If an observer's `update()` calls `attach()` or `detach()`, iterating over a live list produces errors. Iterate over a copy of the list (Step 2's `list(self._observers)`).

**Sending too much state in `notify()`.** If every notification sends the entire subject state, observers receive data they don't need. Consider a push-pull model: `notify()` sends a hint; observers call specific getters for what they need.
