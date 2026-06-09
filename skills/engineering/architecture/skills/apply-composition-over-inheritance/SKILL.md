---
name: apply-composition-over-inheritance
description: Use when designing relationships between classes or modules — especially when tempted to extend a class to reuse its behavior, or when an inheritance hierarchy is becoming deep or brittle.
source: 'Gamma, Helm, Johnson & Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software", Addison-Wesley, 1994'
tags: [tight-coupling, composition, design-patterns, oop, developer, fragile-base-class, flexibility, maintainability]
---

# Apply Composition Over Inheritance

Assemble behavior by combining objects rather than extending class hierarchies. Inherit only when a true "is-a" relationship exists and the substitution guarantee holds.

## Why This Is Best Practice

**Adopted by:** Gang of Four (the foundational OOP design patterns book explicitly recommends this in its introduction); Google (Engineering Practices); Go language (no class inheritance by design); Rust (no class inheritance by design); Kotlin (classes are `final` by default, requiring explicit opt-in to inheritance) — the trend across modern language design is to make composition the default path.
**Impact:** Taivalsaari (1996, IEEE TSE) found inheritance is the single greatest source of tight coupling in OOP systems and the primary barrier to component reuse. Inheritance creates compile-time binding between parent and child that cannot be altered at runtime; composition creates runtime flexibility where behavior can be swapped, extended, or mocked independently. The Gang of Four (1994) cite this principle as the root motivation behind 12 of their 23 design patterns — Strategy, Decorator, Bridge, Composite, Proxy, and others all exist specifically to replace inheritance with composition.
**Why best:** Inheritance solves one problem (reuse) while creating three: it exposes internal implementation to subclasses (fragile base class problem), it locks the class hierarchy at design time, and it forces every subclass to accept the full contract of the parent even when only part is needed. Composition solves the same reuse problem without these costs. Mixins and traits are a middle path — they avoid deep hierarchies but can create implicit dependencies; composition with explicit interfaces is more auditable.

Sources: Gamma, Helm, Johnson & Vlissides, "Design Patterns" (Addison-Wesley, 1994); Taivalsaari, "On the Notion of Inheritance" (IEEE TSE, 1996); Seemann, "Dependency Injection Principles, Practices, and Patterns" (Manning, 2019)

## Steps

### 1. Apply the "is-a" vs "has-a" test

Before inheriting, answer: "Is every instance of the subclass genuinely a kind of the parent, in every context, forever?"

| Relationship | Correct model |
|---|---|
| `Dog` is-a `Animal` — always, substitutable | Inheritance (if LSP holds) |
| `Car` has-a `Engine` — uses it, isn't one | Composition |
| `Stack` is-a `Vector` — false (Java's mistake) | Composition |
| `Button` has-a `ClickHandler` — behavior varies | Composition |
| `Square` is-a `Rectangle` — fails LSP | Composition |

If the answer is "mostly" or "in most cases", use composition. Inheritance requires "always, without exception."

---

### 2. Replace behavior inheritance with strategy injection

When a subclass only exists to override one method, extract that method into an interface and inject it.

```python
# Bad — inheritance just to vary one behavior
class EmailNotifier:
    def send(self, message):
        smtp.send(message)

class SlackNotifier(EmailNotifier):
    def send(self, message):       # overrides entire parent contract
        slack.post(message)

class SMSNotifier(EmailNotifier):
    def send(self, message):
        sms.send(message)

# Good — compose with a strategy; no inheritance needed
class Notifier:
    def __init__(self, channel):   # inject the behavior
        self.channel = channel

    def send(self, message):
        self.channel.send(message)

# Each channel is independent; Notifier is closed to modification
email_notifier = Notifier(EmailChannel())
slack_notifier = Notifier(SlackChannel())
```

---

### 3. Replace structural inheritance with delegation

When a class inherits to gain access to parent methods, delegate instead.

```python
# Bad — Stack inherits Vector to reuse add/remove; exposes all Vector methods
class Stack(list):
    def push(self, item): self.append(item)
    def pop(self): return super().pop()
    # but now callers can also insert(0, x) or sort() — breaking the stack contract

# Good — compose with a list; expose only what Stack should expose
class Stack:
    def __init__(self):
        self._items = []           # private; callers can't reach it

    def push(self, item):
        self._items.append(item)

    def pop(self):
        return self._items.pop()

    def peek(self):
        return self._items[-1]
```

---

### 4. Replace mixin hierarchies with explicit composition

Deep mixin chains hide dependencies and make behavior unpredictable. Name and inject each capability explicitly.

```python
# Bad — mixin chain; hard to know what LoggingMixin needs to function
class LoggingMixin:
    def log(self, msg): print(f"[{self.__class__.__name__}] {msg}")

class RetryMixin:
    def with_retry(self, fn, retries=3): ...

class PaymentService(LoggingMixin, RetryMixin, BaseService):
    ...   # which mixin does what? what does BaseService expose?

# Good — compose named dependencies explicitly
class PaymentService:
    def __init__(self, logger, retry_policy):
        self.logger = logger
        self.retry_policy = retry_policy

    def charge(self, amount):
        self.logger.info(f"charging {amount}")
        return self.retry_policy.execute(lambda: self._do_charge(amount))
```

---

### 5. Use composition to enable testing

Inherited behavior is hard to isolate in tests. Composed dependencies can be replaced with fakes.

```python
# Bad — FileReport inherits from BaseReport; test requires filesystem
class BaseReport:
    def save(self, path): open(path, 'w').write(self.render())

class SalesReport(BaseReport):
    def render(self): ...

# in test — must mock filesystem or use temp files
report = SalesReport()
report.save("/tmp/test.csv")   # touches disk

# Good — inject the storage; test injects a fake
class SalesReport:
    def __init__(self, storage):
        self.storage = storage

    def export(self, path):
        self.storage.write(path, self.render())

# in test — no disk, no mocking framework needed
fake_storage = FakeStorage()
SalesReport(fake_storage).export("test.csv")
assert fake_storage.written["test.csv"] == expected
```

## When NOT to Use

Use inheritance instead when:

- **True polymorphism is required**: callers must treat `Dog` and `Cat` identically as `Animal` — Liskov Substitution holds completely.
- **Framework extension points**: frameworks (Django views, React components) define base classes as intentional extension points. Use them.
- **Template Method pattern**: a fixed algorithm with one or two variable steps that vary per subclass — inheritance is cleaner than injecting a single hook function.
- **Shallow hierarchy (≤ 2 levels)**: deep hierarchies are the problem. A single level of inheritance for genuine "is-a" is fine.

Composition also has limits:

- **Performance-critical hot paths**: virtual dispatch carries a small overhead vs direct calls. In tight loops, measure before abstracting.
- **Value types and data classes**: pure data structures with no behavior don't need either model. Use plain structs or dataclasses.
- **Framework-mandated inheritance**: when a framework requires it (e.g., `unittest.TestCase`, Django `View`), use it. Don't fight the framework.

## Common Mistakes

**Composition for its own sake.** Wrapping every single method of an inner object in a pass-through method is worse than inheritance. If you're writing `def foo(self): return self._inner.foo()` for every method, reconsider — either inheritance is genuinely correct here, or you need a different abstraction.

**Deep composition trees.** `A` composes `B` which composes `C` which composes `D` — if the call chain is opaque, you've traded a deep inheritance tree for a deep dependency tree. Keep composition graphs shallow and dependencies named.

**Confusing composition with aggregation.** Composition: the inner object's lifetime is owned by the outer object. Aggregation: the inner object is shared and outlives the outer. Both are valid; know which you're using.
