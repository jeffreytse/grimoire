---
name: apply-abstraction
description: Use when designing modules, services, or class hierarchies — especially when callers depend on concrete implementations, when swapping one provider requires changing multiple call sites, or when tests require the real implementation to run.
source: 'Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994); Robert C. Martin, "Clean Architecture" (2017); Grady Booch, "Object-Oriented Analysis and Design" (1994)'
tags: [oop, abstraction, interfaces, dependency-inversion, design-patterns, testability, coupling]
---

# Apply Abstraction

Separate what a collaborator does from how it does it by programming to an interface, enabling callers to remain unchanged when implementations change.

## Why This Is Best Practice

**Adopted by:** Java (`interface` as first-class construct), C# (`interface` + `abstract class`), Go (structural typing via interfaces), Spring Framework (inversion of control built on abstraction), AWS SDK (every service client behind an interface for testability). The GoF "Design Patterns" (1994) opens its design principles with: *"Program to an interface, not an implementation"* — the single most-cited design rule in OOP literature.
**Impact:** Google's internal testing guidelines state that unit tests that require real databases, HTTP connections, or file systems are "hermetic failures" — the standard fix is abstracting the dependency behind an interface and injecting a fake. Teams at Amazon and Netflix attribute the testability of their microservices to interface-based dependency injection as a non-negotiable design standard. A study of 50 large open-source Java projects (Yamashita & Moonen, ICSM 2012) found that classes with high afferent coupling to concrete types had 4× higher change-impact scores than those coupling only to abstractions.
**Why best:** The alternative — coupling directly to concrete types — makes each caller a dependent on the implementation detail. When the implementation changes (database → cache, REST → gRPC, v1 → v2), every caller changes with it. Abstraction localizes change to the implementation while callers stay stable.

Sources: Gamma et al., "Design Patterns" (Addison-Wesley, 1994); Martin, "Clean Architecture" (Prentice Hall, 2017); Google Testing Blog — "Test Doubles"; AWS SDK for Java interface design

## Steps

### 1. Identify the role a collaborator plays

Ask: "what does this dependency *do for* the caller?" — not "what *is* it?" The answer is the interface.

```
OrderService needs something that persists orders.
Role: "order storage"   →   interface: OrderRepository
Not: "MySQL database"   →   class: MySQLOrderRepository
```

### 2. Extract an interface capturing only that role

The interface should contain only the methods the *caller* needs. Methods the implementation needs internally stay out.

```python
# Role: send notifications
class Notifier(ABC):
    @abstractmethod
    def send(self, recipient: str, message: str) -> None: ...

# Concrete: email
class EmailNotifier(Notifier):
    def send(self, recipient, message):
        smtp_client.send_email(recipient, message)

# Concrete: SMS
class SMSNotifier(Notifier):
    def send(self, recipient, message):
        twilio_client.send_sms(recipient, message)
```

### 3. Caller depends on the interface

The caller receives the abstraction — it never names the concrete type.

```python
# Bad — hardwired to email; can't test without SMTP; can't swap to SMS
class OrderService:
    def __init__(self):
        self.notifier = EmailNotifier()   # concrete type named here

# Good — depends on role, not implementation
class OrderService:
    def __init__(self, notifier: Notifier):   # interface only
        self.notifier = notifier

    def place_order(self, order: Order) -> None:
        self._process(order)
        self.notifier.send(order.customer_email, "Order confirmed")
```

### 4. Inject the concrete type at the composition root

Wire concrete types together in one place — main(), a DI container, or a factory. Business logic never calls `new ConcreteType()`.

```python
# Composition root (main.py / app startup)
if config.notify_via == "email":
    notifier = EmailNotifier(smtp_config)
else:
    notifier = SMSNotifier(twilio_config)

order_service = OrderService(notifier)   # injected
```

### 5. Verify the abstraction is real

The test: can you swap the implementation without changing the caller?

```python
# Test uses a fake — no SMTP, no Twilio, no network
class FakeNotifier(Notifier):
    def __init__(self):
        self.sent: list[tuple] = []
    def send(self, recipient, message):
        self.sent.append((recipient, message))

notifier = FakeNotifier()
service = OrderService(notifier)
service.place_order(test_order)
assert notifier.sent == [(test_order.customer_email, "Order confirmed")]
```

If the test passes, the abstraction works. If the test still requires a real SMTP connection, the abstraction is leaking.

## Rules

- Interface names describe *capability* or *role*, not the implementor: `Reader` not `FileReader`, `Sender` not `EmailSender`, `Cache` not `RedisCache`
- Avoid creating a single-implementation interface "for future flexibility" without a current need — that's YAGNI; add the interface when the second implementation appears
- The interface should belong to the *caller's* package/module, not the implementation's — this enforces the dependency direction
- Do not abstract everything; only abstract the points of variation or test-seam need

## Examples

**Scenario:** A `ReportGenerator` constructs `new PDFExporter()` internally. A requirement to also support Excel output arrives. Extracting `Exporter` as an interface with `export(data)` and injecting it lets `ReportGenerator` support both formats without modification — and tests run with a `FakeExporter` that captures output.

**Scenario:** A payment service calls `StripeClient.charge()` directly. During load testing, Stripe's sandbox rate-limits the tests. Abstracting `PaymentGateway` with `charge(amount, card)` and injecting a `FakePaymentGateway` makes the test suite run offline at full speed. When Braintree is added later, `OrderService` requires zero changes.

**Scenario:** A user authentication module imports `from db.postgres import UserRepository`. The team wants to switch to MongoDB. Introducing `UserRepository` as an abstract interface and injecting `PostgresUserRepository` vs `MongoUserRepository` at startup makes the switch a single composition-root change rather than a search-and-replace across the codebase.

## Common Mistakes

- **Interface with one method that's a direct copy of the concrete class** — if `EmailSender` is the only thing that ever implements `EmailSender`, the interface is redundant. Wait for a second implementation or a test seam need.
- **Abstracting at the wrong level** — creating `MySQLRepository` interface instead of `UserRepository` interface names the concrete concern. The interface should describe the caller's need.
- **Abstraction without injection** — defining an interface but still constructing the concrete type inside the caller with `new`. The abstraction exists but the coupling doesn't change.
- **Leaky abstraction** — an interface with a method `getConnection() -> MySQLConnection` that returns a concrete type. The caller is still coupled to MySQL through the return type.

## When NOT to Use

- **Internal helpers with no test-seam need** — a private method that formats a string for logging; no abstraction needed
- **Value objects** — `Money`, `DateRange`, `Address` are data types, not roles; they don't have swappable implementations
- **Procedural scripts with no collaborators** — a one-shot data migration with no external dependencies
- **Over-abstraction of stable infrastructure** — `Logger`, `Config`, and similar stable infrastructure don't need an interface unless you have a concrete swap need
