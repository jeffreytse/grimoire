---
name: apply-proxy-pattern
description: Use when you need to control access to an object — adding lazy initialization, access control, logging, caching, or remote access — without changing the object's interface or the client's code.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 207–217; Java RMI (remote proxy); Spring AOP (@Transactional, @Cacheable use proxies); Python unittest.mock; Hibernate lazy loading
tags: [design-patterns, structural, proxy, oop, developer, access-control, lazy-loading, caching]
related: [apply-decorator-pattern, apply-facade-pattern, apply-solid-principles]
---

# Apply Proxy Pattern

Provide a surrogate or placeholder for another object to control access to it.

## Why This Is Best Practice

**Adopted by:** Spring Framework's AOP (every `@Transactional` and `@Cacheable` bean
is a proxy — used in the majority of Spring applications), Java RMI (remote proxy for
distributed objects — Java's original distributed computing model), Hibernate's lazy
loading (collection proxies — default in every Hibernate application), and Python's
`unittest.mock` (a virtual proxy over any object).
**Impact:** Spring's transaction proxy is cited in the Spring documentation as the
reason `@Transactional` requires no manual connection management — 80%+ reduction in
transaction-handling boilerplate (Spring documentation). Hibernate's virtual proxy for
lazy loading avoids loading entire object graphs; without it, a single entity load
triggers cascading database queries.
**Why best:** Adding cross-cutting concerns (logging, auth, caching) directly to the
subject class violates SRP — the class now handles both its own logic and the concern.
A proxy keeps the subject focused while the proxy layer intercepts access.

Sources: Gamma et al. (1994) pp. 207–217; Spring AOP documentation;
Hibernate ORM documentation (proxy-based lazy loading)

## Steps

### Step 1: Define the subject interface shared by real subject and proxy

```python
from abc import ABC, abstractmethod

class ImageLoader(ABC):
    @abstractmethod
    def load(self, path: str) -> bytes: ...

    @abstractmethod
    def display(self) -> None: ...
```

### Step 2: Implement the real subject

```python
class RealImageLoader(ImageLoader):
    def __init__(self, path: str):
        self._path = path
        self._data: bytes | None = None

    def load(self, path: str) -> bytes:
        print(f"[IO] Loading {path} from disk")
        with open(path, "rb") as f:
            return f.read()

    def display(self) -> None:
        if self._data:
            print(f"Displaying {len(self._data)} bytes")
```

### Step 3: Implement the proxy — wraps the real subject, intercepts calls

**Virtual proxy (lazy initialization):**

```python
class LazyImageProxy(ImageLoader):
    def __init__(self, path: str):
        self._path = path
        self._real: RealImageLoader | None = None

    def _ensure_loaded(self):
        if self._real is None:
            self._real = RealImageLoader(self._path)
            self._real._data = self._real.load(self._path)

    def load(self, path: str) -> bytes:
        self._ensure_loaded()
        return self._real._data

    def display(self) -> None:
        self._ensure_loaded()
        self._real.display()
```

**Protection proxy (access control):**

```python
class AuthImageProxy(ImageLoader):
    def __init__(self, real: ImageLoader, user_role: str):
        self._real = real
        self._role = user_role

    def load(self, path: str) -> bytes:
        if self._role != "admin":
            raise PermissionError("only admins can load raw images")
        return self._real.load(path)

    def display(self) -> None:
        self._real.display()
```

### Step 4: Client uses proxy through the interface — unaware of proxy

```python
loader: ImageLoader = LazyImageProxy("/images/banner.png")
# File not loaded yet
loader.display()   # now loaded — only when needed
```

### Step 5: Choose proxy type based on purpose

| Type | Purpose | Example |
|------|---------|---------|
| Virtual | Lazy initialization | Hibernate lazy collections |
| Remote | Local stand-in for remote object | Java RMI stub |
| Protection | Access control | Auth proxy |
| Caching | Memoize expensive calls | HTTP cache proxy |
| Logging | Audit access | `@Transactional` logging |

## When NOT to Use

- **Simple access with no cross-cutting concerns** — if the subject needs no interception, a proxy adds indirection for no benefit.
- **When the performance cost of indirection exceeds the benefit** — a caching proxy for a trivial computation adds overhead. Profile first.

## Common Mistakes

**Proxy that changes the interface.** A proxy must implement exactly the same interface as the subject. If it changes method signatures, it's a Decorator or Adapter, not a Proxy.

**Forgetting thread safety in virtual proxies.** The "check then initialize" in `_ensure_loaded()` is a race condition. Use a lock or `functools.cached_property` (Python) for thread-safe lazy initialization.

**Using Proxy when Decorator is correct.** Decorator adds behavior dynamically at runtime for an object known to the client. Proxy controls access to an object the client typically doesn't hold directly. The distinction matters for testability.
