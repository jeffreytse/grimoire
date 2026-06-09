---
name: apply-chain-of-responsibility-pattern
description: Use when more than one object may handle a request and the handler isn't known a priori — letting you pass the request along a chain of handlers until one handles it, without coupling the sender to any specific handler.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 223–232; Java Servlet Filter chain; Express.js middleware; Python logging handler chain
tags: [design-patterns, behavioral, chain-of-responsibility, oop, developer, request-handling, decoupling]
related: [apply-command-pattern, apply-decorator-pattern, apply-solid-principles]
---

# Apply Chain of Responsibility Pattern

Pass a request along a chain of handlers; each handler decides to process it or forward it to the next handler.

## Why This Is Best Practice

**Adopted by:** Java Servlet Filter API (every Java web request passes through a
configurable filter chain — the foundation of Spring Security, CORS handling, and
logging in all Java web apps), Express.js `app.use()` middleware (the core request
pipeline model), Python's `logging` module (each logger's handlers form a chain),
and browser DOM event propagation (bubbling is a chain-of-responsibility).
**Impact:** GoF documents that without Chain of Responsibility, the sender must know
which object handles its request — adding conditional routing logic that must be updated
every time a new handler is added. Java Servlet's filter chain allows adding security,
compression, or logging filters without modifying the servlet or other filters.
**Why best:** Direct `if/elif` dispatch — the alternative — couples the sender to all
possible handlers and requires modification for every new handler. Chain of Responsibility
lets handlers be added, removed, and reordered without touching the sender or other handlers.

Sources: Gamma et al. (1994) pp. 223–232; Java Servlet Filter specification;
Express.js middleware documentation

## Steps

### Step 1: Define the handler interface with a method to set the next handler

```python
from abc import ABC, abstractmethod
from typing import Optional

class SupportHandler(ABC):
    def __init__(self):
        self._next: Optional["SupportHandler"] = None

    def set_next(self, handler: "SupportHandler") -> "SupportHandler":
        self._next = handler
        return handler   # allows chaining: h1.set_next(h2).set_next(h3)

    @abstractmethod
    def handle(self, ticket_level: int) -> str | None: ...
```

### Step 2: Implement concrete handlers that decide to handle or pass

```python
class FrontlineSupport(SupportHandler):
    def handle(self, level: int) -> str | None:
        if level <= 1:
            return "Frontline resolved the ticket"
        if self._next:
            return self._next.handle(level)
        return None

class TechSupport(SupportHandler):
    def handle(self, level: int) -> str | None:
        if level <= 3:
            return "Tech support resolved the ticket"
        if self._next:
            return self._next.handle(level)
        return None

class EngineeringSupport(SupportHandler):
    def handle(self, level: int) -> str | None:
        return "Engineering resolved the ticket"   # always handles
```

### Step 3: Assemble the chain at the composition root

```python
frontline = FrontlineSupport()
tech = TechSupport()
engineering = EngineeringSupport()

frontline.set_next(tech).set_next(engineering)
```

### Step 4: Send requests through the chain without knowing which handler responds

```python
for level in [1, 2, 3, 4, 5]:
    result = frontline.handle(level)
    print(f"Level {level}: {result}")
```

### Step 5: Decide whether unhandled requests are an error or valid

If every request must be handled, add a final handler that always handles (catches all):

```python
class FallbackHandler(SupportHandler):
    def handle(self, level: int) -> str | None:
        return f"No handler found for level {level} — escalated to manager"
```

Alternatively, let `None` return indicate "not handled" and check at the call site.

## When NOT to Use

- **When exactly one handler always applies** — if the handler is deterministic, dispatch directly. Chain adds overhead with no flexibility benefit.
- **When request handling must be auditable** — chains make it hard to trace which handler acted. Log entry/exit in each handler if auditability is required.

## Common Mistakes

**Forgetting to call `self._next.handle()`** — if a handler doesn't forward when it can't handle the request, the chain silently swallows it. Always forward when not handling.

**Circular chains.** If `h1.next = h2` and `h2.next = h1`, the chain loops forever. Validate chain configuration or use a framework (Spring Security, Express) that prevents this.

**Coupling handler order to business logic.** If the chain only works in one specific order, it's not really flexible. Handlers should be order-independent for the cases they handle.
