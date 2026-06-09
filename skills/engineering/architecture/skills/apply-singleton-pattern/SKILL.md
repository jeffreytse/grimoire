---
name: apply-singleton-pattern
description: Use when exactly one instance of a class must exist across the entire application — such as a configuration store, connection pool, or logger — and global access to that instance is required.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 127–136; Python logging module (module-level singletons); Java Spring @Bean(singleton scope); Go sync.Once
tags: [design-patterns, creational, singleton, oop, developer, global-state, instance-control]
related: [apply-solid-principles, apply-factory-method-pattern]
---

# Apply Singleton Pattern

Ensure a class has exactly one instance and provide a global access point to it.

## Why This Is Best Practice

**Adopted by:** Spring Framework (default bean scope is singleton — the most-used Java
framework, 65%+ of Java enterprise apps per JetBrains 2023 survey), Python's `logging`
module (each `getLogger(name)` returns the same Logger instance for that name), Go's
`sync.Once` (idiomatic Go singleton construction), and Java `Runtime.getRuntime()`
(JVM runtime singleton, in every JVM since Java 1.0).
**Impact:** GoF documents that uncontrolled global state causes initialization ordering
bugs and duplicate resource allocation. Connection pools (HikariCP, the most-used Java
connection pool) are singletons by necessity — two pools to the same DB double
connections and break max-connection limits.
**Why best:** The alternative — passing the shared instance everywhere via constructor
(dependency injection) — is architecturally cleaner and preferred when testability
matters. Singleton is correct when the instance is a true system-level resource (one
config, one log system, one thread pool) and DI wiring overhead is unjustified.

Sources: Gamma et al. (1994) pp. 127–136; Spring Framework documentation (Bean scopes);
Go `sync.Once` documentation; HikariCP design rationale

## Steps

### Step 1: Make the constructor private (or restrict it)

```python
class Configuration:
    _instance: "Configuration | None" = None

    def __new__(cls) -> "Configuration":
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialized = False
        return cls._instance

    def __init__(self):
        if self._initialized:
            return
        self._data: dict = {}
        self._initialized = True
```

### Step 2: Use `__new__` (Python) or `getInstance()` (Java/C++) as the access point

```python
# All callers get the same instance
config1 = Configuration()
config2 = Configuration()
assert config1 is config2   # True
```

### Step 3: Make thread-safe with a lock for lazy initialization

```python
import threading

class ThreadSafeConfiguration:
    _instance = None
    _lock = threading.Lock()

    def __new__(cls):
        if cls._instance is None:
            with cls._lock:
                if cls._instance is None:   # double-checked locking
                    cls._instance = super().__new__(cls)
        return cls._instance
```

Double-checked locking avoids acquiring the lock on every access after initialization.

### Step 4: Prefer module-level singletons in Python over class-based ones

Python modules are singletons by default — imported once, cached in `sys.modules`:

```python
# config.py — the module IS the singleton
_data: dict = {}

def get(key: str):
    return _data[key]

def set(key: str, value):
    _data[key] = value
```

```python
# caller
import config
config.set("debug", True)
```

### Step 5: In test environments, provide a reset mechanism

```python
@classmethod
def _reset(cls):
    cls._instance = None   # allow re-initialization in tests
```

Mark this method as test-only. Production code never calls it.

## When NOT to Use

- **When testability matters** — Singletons that hold mutable state make tests order-dependent. Prefer dependency injection so tests can pass isolated instances.
- **When multiple instances might ever be needed** — a "singleton database" becomes a problem when you add a read replica. If the constraint might relax, use DI from the start.
- **When it's just convenient global access** — global variables and singletons have the same test and coupling problems. The pattern is justified only when one instance is a hard architectural constraint, not convenience.

## Common Mistakes

**Not guarding against subclassing.** If the singleton class can be subclassed, a subclass can create a second instance. Use `final` (Java) or `__init_subclass__` guards (Python) to prevent this.

**Ignoring thread safety.** A naive `if not _instance: _instance = cls()` has a race condition. Use `sync.Once` (Go), `synchronized` (Java), or double-checked locking (Python).

**Using Singleton where Dependency Injection belongs.** If the "singleton" needs to be swapped in tests, replaced in different environments, or extended, it's not a true singleton — it's a shared service that should be injected.
