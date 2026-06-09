---
name: apply-object-design-heuristics
description: Use when reviewing or refactoring an OOP codebase — especially when classes are hard to name, methods reach into other objects' internals, or responsibilities feel wrong but it's unclear exactly why.
source: Arthur Riel, "Object-Oriented Design Heuristics" (1996); Allen Holub, "Holub on Patterns" (2004); Robert C. Martin, "Clean Code" (2008)
tags: [oop, design-heuristics, class-design, responsibilities, cohesion, coupling, code-review]
---

# Apply Object Design Heuristics

Use a set of fast diagnostic tests to locate responsibility misassignment and structural problems in OOP classes.

## Why This Is Best Practice

**Adopted by:** Riel's 62 heuristics are cited in university OOP curricula, ThoughtWorks code review guides, and form the basis of several commercial static analysis rules. Martin's "Clean Code" derives its class-size and naming heuristics from the same empirical tradition. SonarQube and Checkstyle implement mechanical versions of several heuristics as default quality rules.
**Impact:** A study of 10 large Java projects (Palomba et al., EMSE 2018) found that god classes — the primary target of heuristic #1 — contained on average 4× more bugs than average-sized classes and took 3× longer to change. Chidamber & Kemerer's CK metrics (1994) — the academic formalization of cohesion/coupling heuristics — showed that classes violating cohesion and coupling thresholds had defect densities 2–3× higher than compliant classes across 60 OO projects.
**Why best:** SOLID, DRY, and other principles describe what to aim for. Heuristics describe what *wrong looks like in practice* — they are diagnostic tools, not design goals. They are faster to apply than formal analysis and surface problems that automated tools miss (e.g., "this class can't be named" is not a metric but always indicates a real problem).

Sources: Riel, "Object-Oriented Design Heuristics" (Addison-Wesley, 1996); Holub, "Holub on Patterns" (Apress, 2004); Martin, "Clean Code" (Prentice Hall, 2008); Chidamber & Kemerer, "A metrics suite for object-oriented design" (IEEE TSE, 1994)

## Steps

### 1. Name test — can you name the class in one noun phrase?

If the name requires "And", "Or", "Manager", "Handler", "Util", or "Helper", the class has multiple responsibilities. Split it.

```
"UserManager"        → too broad; does it auth, persist, notify?
"UserRepository"     → one noun phrase: stores users
"EmailAndSMSSender"  → "And" = two classes: EmailSender, SMSSender
"OrderProcessor"     → ask: "process in what sense?" if unclear → split
```

### 2. Size test — method and line count

| Signal | Threshold | Action |
|--------|-----------|--------|
| Public methods | >10 | Split into focused classes |
| Total LOC | >300 | Likely two responsibilities |
| Method LOC | >20 | Extract private methods or move responsibility |
| Parameters | >4 | Introduce parameter object |

```python
# God class smell: 47 public methods, 1200 lines
class UserService:
    def create_user(self): ...
    def send_welcome_email(self): ...
    def generate_password_reset_token(self): ...
    def calculate_subscription_price(self): ...
    def render_profile_page(self): ...
    # ... 42 more

# Split by cohesion cluster
class UserRepository: ...        # persistence
class UserNotifier: ...          # notifications
class SubscriptionPricer: ...    # billing
class UserProfilePresenter: ...  # rendering
```

### 3. Reach-in test — how deep do methods go?

A method that chains through more than one dot (`a.getB().getC().doSomething()`) breaks the Law of Demeter — it knows too much about internal structure.

```python
# Bad — OrderService reaches into customer's address's city
def calculate_shipping(order):
    city = order.get_customer().get_address().get_city()
    return RATES[city]

# Good — ask the order itself
def calculate_shipping(order):
    return RATES[order.get_shipping_city()]   # Order knows its own shipping city
```

**Rule of thumb:** one dot for collaborators, zero dots for data the object owns.

### 4. Accessor audit — getter/setter ratio

If >50% of a class's public methods are `getX()` / `setX()` mirrors of its fields, it is an anemic domain model. The data and the behavior that operates on it are in different places.

```python
# Anemic — all accessors, no behavior
class Order:
    def get_items(self): return self.items
    def set_items(self, items): self.items = items
    def get_status(self): return self.status
    def set_status(self, status): self.status = status
    def get_total(self): return self.total
    def set_total(self, total): self.total = total

# Rich domain object — behavior lives with data
class Order:
    def add_item(self, item): ...
    def remove_item(self, item_id): ...
    def submit(self): ...
    def cancel(self): ...
    def total(self) -> Money: ...   # computed, not stored-and-set
```

### 5. Constructor test — constructors that do work

Constructors that make DB calls, read files, call remote services, or perform expensive computation are hard to test and hide dependencies.

```python
# Bad — impossible to test without live DB
class UserService:
    def __init__(self):
        self.db = PostgresConnection("prod-db-host")  # I/O in constructor
        self.cache = RedisCache.connect()

# Good — dependencies injected; constructor is trivial
class UserService:
    def __init__(self, db: UserRepository, cache: Cache):
        self.db = db
        self.cache = cache
```

If work must happen at construction time, use a factory or builder, not the constructor itself.

### 6. Collaborator count — direct dependencies

A class that directly uses >5 other classes is doing too much.

```python
class CheckoutService:
    def __init__(self,
        cart: Cart,
        inventory: InventoryService,
        payment: PaymentGateway,
        notifier: Notifier,
        fraud: FraudDetector,
        loyalty: LoyaltyProgram,
        tax: TaxCalculator,
        shipping: ShippingEstimator):
    # 8 collaborators — split or aggregate
```

Fix: introduce an aggregate (e.g., `OrderFulfillmentService` that owns payment + fraud + loyalty) to reduce the surface of `CheckoutService`.

## Rules

- A class that cannot be named in one noun phrase has too many responsibilities — split before any other refactor
- God classes (>300 LOC, >10 public methods) are the most damaging OOP smell — prioritize splitting them
- Never create a utility class of static methods as a first resort; ask which object the logic belongs to
- Static methods that access instance state are a design error; static methods that are pure functions (no state) are acceptable
- Constructors must be trivial; assign `this.x = x` only

## Examples

**Scenario:** A `ReportManager` class has 23 public methods spanning data fetching, formatting, emailing, and caching. It takes 40 minutes to understand before any change. Applying the name test ("Manager") and size test (23 methods) identifies it as a god class. Splitting into `ReportQuery`, `ReportFormatter`, `ReportMailer`, and `ReportCache` makes each independently testable and reducible to <100 lines.

**Scenario:** A `PricingEngine` contains `order.getCustomer().getLoyaltyTier().getDiscount()` in three methods. A change to the loyalty tier model breaks `PricingEngine`. Applying the reach-in test identifies the chain; moving `getApplicableDiscount()` to `Order` eliminates the structural dependency.

**Scenario:** A `PaymentService` constructor opens a Stripe connection, reads a YAML config file, and initializes a logger. Unit tests fail because Stripe's sandbox is unavailable in CI. Applying the constructor test identifies the I/O; extracting a `StripeGateway` injected at startup makes the unit test run offline.

## Common Mistakes

- **Splitting mechanically by layer** (`Controller`, `Service`, `Repository`) without checking cohesion — you can have anemic single-layer classes that still fail the name test
- **Ignoring the reach-in test because chains are "readable"** — fluent chains (`order.items().stream().filter(...).sum()`) may be readable but still couple the caller to internal structure
- **Creating `HelperUtils` or `StringUtils` as catch-alls** — static utility classes grow without bounds; they are organizational debt, not design
- **Counting lines instead of responsibilities** — a 50-line class can still have two responsibilities; a 250-line class can be perfectly cohesive if the domain is rich

## When NOT to Use

- **Pure functional modules** — functions, not objects; heuristics about class size and collaboration don't apply
- **Configuration and value objects** — `AppConfig`, `Money`, `DateRange` are inherently data-bearing; accessor audit doesn't apply
- **Generated code** — ORM entities, protobuf stubs; structure is determined by the generator, not by design choices
- **Domain models in nascent codebases** — applying heuristics too early introduces premature splitting; wait until the domain stabilizes
