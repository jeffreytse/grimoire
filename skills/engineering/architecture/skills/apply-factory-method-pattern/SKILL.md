---
name: apply-factory-method-pattern
description: Use when a class must create objects but should not be responsible for deciding which concrete class to instantiate — letting subclasses or configuration determine the type at runtime.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 107–116; adopted in Java (Collection.iterator(), JDBC DriverManager), Python (logging.getLogger()), Django (model managers)
tags: [design-patterns, creational, factory-method, oop, developer, extensibility, decoupling]
related: [apply-abstract-factory-pattern, apply-solid-principles]
---

# Apply Factory Method Pattern

Define an interface for creating an object, but let subclasses decide which class to instantiate.

## Why This Is Best Practice

**Adopted by:** Java's `Collection.iterator()` (the most-called factory method in the
Java standard library), JDBC `DriverManager.getConnection()`, Python's
`logging.getLogger()`, Django ORM model managers, and Spring `BeanFactory` — all of
which allow callers to receive a product without knowing its concrete type.
**Impact:** GoF documents that Factory Method eliminates the coupling between
application-specific classes and the objects they must create. Django's model manager
factory is one of the most-cited examples of the pattern enabling framework
extensibility: users override the factory to return custom QuerySet subclasses without
modifying the ORM.
**Why best:** The alternative — `new ConcreteProduct()` at every call site — couples the
caller to the specific class. When the concrete type must change (new database driver,
new log handler, new collection type), every call site changes. Factory Method
centralizes that coupling in one overridable method.

Sources: Gamma et al. (1994) pp. 107–116; JDBC API specification; Django ORM
documentation (Manager class)

## Steps

### Step 1: Define the product interface

```python
from abc import ABC, abstractmethod

class Notification(ABC):
    @abstractmethod
    def send(self, message: str) -> None: ...
```

### Step 2: Define the creator with a factory method

```python
class NotificationService(ABC):
    @abstractmethod
    def create_notification(self) -> Notification:
        """Factory method — subclasses override this."""
        ...

    def notify(self, message: str) -> None:
        notification = self.create_notification()   # calls factory method
        notification.send(message)
```

### Step 3: Implement concrete creators that override the factory method

```python
class EmailNotificationService(NotificationService):
    def create_notification(self) -> Notification:
        return EmailNotification()

class SMSNotificationService(NotificationService):
    def create_notification(self) -> Notification:
        return SMSNotification()
```

### Step 4: Keep the creator's core logic independent of the concrete product

The `notify()` method in Step 2 works identically regardless of which
`create_notification()` returns. That independence is the pattern's value — the creator
algorithm is reusable across all product variants.

### Step 5: For stateless factories, use a static/class method variant

When subclassing is unnecessary, use a parameterized factory method:

```python
class NotificationFactory:
    @staticmethod
    def create(channel: str) -> Notification:
        match channel:
            case "email": return EmailNotification()
            case "sms":   return SMSNotification()
            case _: raise ValueError(f"unknown channel: {channel}")
```

This variant is simpler than the full GoF structure when the creator hierarchy is not needed.

## When NOT to Use

- **Families of related objects** — if multiple product types must stay consistent (button + checkbox), use `apply-abstract-factory-pattern` instead.
- **Single, fixed product type** — if the type never varies, `new Product()` is cleaner. Factory Method adds a layer only when variation is the point.
- **When a constructor suffices** — modern Python `__init__` with keyword arguments and type hints is often cleaner than a one-off factory for simple objects.

## Common Mistakes

**Naming the factory method `create()`** — generic names lose the domain signal. Name it for the product: `create_notification()`, `create_connection()`, `build_query()`.

**Putting business logic inside the factory method** — the factory method should only instantiate. Initialization or validation belongs in the product's `__init__`, not the factory.

**Forgetting to parameterize when subclassing is unnecessary** — the static variant (Step 5) is valid and often simpler. Don't force a class hierarchy when a single parameterized method suffices.
