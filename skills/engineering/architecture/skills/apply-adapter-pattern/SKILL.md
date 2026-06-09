---
name: apply-adapter-pattern
description: Use when you need to use an existing class but its interface is incompatible with what your code expects — converting one interface into another without modifying either.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 139–150; Java InputStreamReader (adapts InputStream to Reader); Python io.TextIOWrapper; Django REST Framework serializers
tags: [design-patterns, structural, adapter, oop, developer, interface-compatibility, legacy-integration]
related: [apply-facade-pattern, apply-solid-principles]
---

# Apply Adapter Pattern

Convert the interface of a class into the interface clients expect, enabling classes with incompatible interfaces to work together.

## Why This Is Best Practice

**Adopted by:** Java's `InputStreamReader` (adapts byte-stream `InputStream` to
character-stream `Reader` — in every JVM), Python's `io.TextIOWrapper` (same
adaptation in the standard library), Django REST Framework's Serializer (adapts ORM
models to JSON-serializable dicts), and every database driver layer (JDBC, ODBC) which
adapts vendor-specific protocols to a standard interface.
**Impact:** GoF documents that the Adapter is the primary mechanism for integrating
legacy or third-party code into a new architecture without modifying either side.
Every major language's standard library uses it to bridge abstraction layers — evidence
of universal adoption across all major ecosystems.
**Why best:** The alternative — modifying the existing class — creates coupling to the
consumer and may be impossible (third-party code, compiled libraries). The alternative
of modifying the client breaks encapsulation. Adapter wraps without touching either side.

Sources: Gamma et al. (1994) pp. 139–150; Java I/O class hierarchy documentation;
Python `io` module documentation

## Steps

### Step 1: Identify the target interface the client expects

```python
from abc import ABC, abstractmethod

class JSONLogger(ABC):
    @abstractmethod
    def log(self, event: dict) -> None: ...
```

### Step 2: Identify the adaptee (existing class with incompatible interface)

```python
class LegacyLogger:
    def write_line(self, timestamp: str, level: str, message: str) -> None:
        print(f"[{timestamp}] {level}: {message}")
```

### Step 3: Write the adapter — wraps adaptee, implements target interface

```python
from datetime import datetime

class LegacyLoggerAdapter(JSONLogger):
    def __init__(self, legacy: LegacyLogger):
        self._legacy = legacy

    def log(self, event: dict) -> None:
        self._legacy.write_line(
            timestamp=event.get("ts", datetime.now().isoformat()),
            level=event.get("level", "INFO"),
            message=event.get("message", ""),
        )
```

### Step 4: Wire the adapter at the composition root

```python
legacy = LegacyLogger()
logger: JSONLogger = LegacyLoggerAdapter(legacy)

# Client uses only the target interface — unaware of LegacyLogger
logger.log({"level": "ERROR", "message": "payment failed", "ts": "2026-06-10T12:00:00"})
```

### Step 5: For two-way adaptation, create bidirectional adapters

If both sides need to call each other, write two adapters — one for each direction.
A single class implementing both interfaces creates hidden coupling and is hard to test.

## When NOT to Use

- **When you control both sides** — if you own both the client and the class being used, change one interface to match the other directly. Adapter is for cases where modification is impossible or undesirable.
- **When the interfaces differ too fundamentally** — if the adaptee's semantics are incompatible (not just syntactically different), an adapter produces a leaky abstraction. Consider a `apply-facade-pattern` instead.

## Common Mistakes

**Adapting too much.** An adapter that translates complex business logic is a Facade or a service, not an adapter. Keep adapters thin — only interface translation, no logic.

**Modifying the adaptee to make the adapter simpler.** If you can modify the adaptee, there may be no need for an adapter at all. If you can't, don't.

**Forgetting to handle adapter errors.** When the adaptee throws an exception using its own error types, the adapter should either translate them to the target interface's error types or let them propagate with clear context.
