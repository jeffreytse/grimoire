---
name: apply-solid-principles
description: Use when designing or reviewing classes and modules in OOP code — especially when tests are hard to write, changes ripple unexpectedly, or classes accumulate unrelated responsibilities.
source: Robert C. Martin, "Agile Software Development, Principles, Patterns, and Practices" (2003); Martin, "Clean Code" (2008); Martin "Design Principles and Design Patterns" (2000)
tags: [solid, oop, design-principles, coupling, testability, developer, defect-reduction, maintainability]
---

# Apply SOLID Principles

Apply the five SOLID design principles to OOP code to reduce coupling, improve testability, and make change safe.

## Why This Is Best Practice

**Adopted by:** Google (internal code review guidelines), Microsoft (.NET Framework Design Guidelines), Amazon (principal engineer design standards), ThoughtWorks Technology Radar — majority-adoption standard across OOP-heavy engineering organizations.
**Impact:** Dependency Inversion — the D in SOLID — is the primary enabler of unit-testable code; Google's Engineering Practices documentation and ThoughtWorks Technology Radar cite DI as foundational for testable design. Yamashita & Moonen (ICSM 2012) found high coupling and low cohesion — the exact defects SRP and DIP address — are the strongest predictors of maintainability degradation across 74 Java projects. SonarQube ships SOLID violation rules as default quality gates, directly connecting SOLID compliance to automated defect prevention pipelines.
**Why best:** Competing alternatives — anemic domain models (pure data + service layers with no design rules), procedural decomposition, or ad-hoc OOP with no explicit principles — either defer coupling problems to runtime failures or leave them unaddressed entirely. SOLID is the only OOP design framework with a published acronym, universal naming, tooling support (SonarQube, IDE refactoring tools), and a formal specification — making it auditable and enforceable in code review rather than relying on individual judgment.

Sources: Martin, "Agile Software Development, Principles, Patterns, and Practices" (Prentice Hall, 2003); Martin, "Clean Code" (Prentice Hall, 2008); Google Engineering Practices; Microsoft .NET Design Guidelines

## Steps

### 1. S — Single Responsibility Principle

**Rule:** A class should have one reason to change — one actor whose needs it serves.

**Signal it's violated:** Class name contains "And", "Manager", "Handler", or "Util" with >200 lines. Method `saveUserAndSendEmail()` is one method doing two jobs.

**Fix:** Split by actor. If the finance team and the HR team both need `Employee` to change for different reasons, that's two classes.

```python
# Bad — mixes persistence, business logic, and formatting
class Employee:
    def calculate_pay(self): ...
    def save_to_db(self): ...        # persistence actor
    def generate_report(self): ...  # reporting actor

# Good — one reason to change each
class Employee:
    def calculate_pay(self): ...

class EmployeeRepository:
    def save(self, employee): ...

class EmployeeReporter:
    def generate_report(self, employee): ...
```

---

### 2. O — Open/Closed Principle

**Rule:** Open for extension, closed for modification. Add behavior by adding code, not editing existing code.

**Signal it's violated:** Adding a new feature requires editing an `if/elif` chain or `switch` statement inside an existing class. Every new payment type requires editing `PaymentProcessor`.

**Fix:** Extract the varying behavior to an abstraction. New variants implement the abstraction; core logic never changes.

```python
# Bad — every new shape requires editing area()
def area(shape):
    if shape.type == "circle":
        return 3.14 * shape.radius ** 2
    elif shape.type == "square":
        return shape.side ** 2
    # adding triangle means editing this function

# Good — new shapes extend without modifying existing code
class Shape(ABC):
    @abstractmethod
    def area(self) -> float: ...

class Circle(Shape):
    def area(self): return 3.14 * self.radius ** 2

class Square(Shape):
    def area(self): return self.side ** 2

class Triangle(Shape):    # new shape, zero edits to existing code
    def area(self): return 0.5 * self.base * self.height
```

---

### 3. L — Liskov Substitution Principle

**Rule:** Subtypes must be substitutable for their base type without breaking program correctness. Callers must not need to know which subtype they have.

**Signal it's violated:** `isinstance()` / `typeof` checks in calling code. Overridden methods throw `NotImplementedException`. Square extends Rectangle but breaks `set_width()` behavior.

**Fix:** If a subtype must restrict or ignore inherited behavior, it's not a true subtype — redesign the hierarchy or use composition.

```python
# Bad — Square breaks Rectangle's contract
class Rectangle:
    def set_width(self, w): self.width = w
    def set_height(self, h): self.height = h

class Square(Rectangle):
    def set_width(self, w):
        self.width = w
        self.height = w   # silently breaks caller assumptions

# Good — model separately; use a common abstraction only if behavior is truly shared
class Shape(ABC):
    @abstractmethod
    def area(self) -> float: ...

class Rectangle(Shape):
    def area(self): return self.width * self.height

class Square(Shape):
    def area(self): return self.side ** 2
```

---

### 4. I — Interface Segregation Principle

**Rule:** Clients should not depend on methods they don't use. Split fat interfaces into focused ones.

**Signal it's violated:** Implementing class has empty, stub, or `raise NotImplementedException` methods. Interface has >7 methods serving callers with different needs.

**Fix:** Split by client. `Printer` only needs `print()`; `Scanner` only needs `scan()`. Don't force `Printer` to implement `scan()`.

```python
# Bad — all implementors forced to implement all methods
class Machine(ABC):
    def print(self): ...
    def scan(self): ...
    def fax(self): ...

class SimplePrinter(Machine):
    def print(self): ...
    def scan(self): raise NotImplementedError    # forced stub
    def fax(self): raise NotImplementedError     # forced stub

# Good — implement only what you use
class Printer(ABC):
    def print(self): ...

class Scanner(ABC):
    def scan(self): ...

class SimplePrinter(Printer):
    def print(self): ...   # only what it needs

class AllInOne(Printer, Scanner):
    def print(self): ...
    def scan(self): ...
```

---

### 5. D — Dependency Inversion Principle

**Rule:** High-level modules should not depend on low-level modules. Both should depend on abstractions. Pass dependencies in, don't construct them inside.

**Signal it's violated:** `new ConcreteService()` inside a class. `import MySQLDatabase` in business logic. Tests require a live database or HTTP connection.

**Fix:** Define an interface in the high-level module; pass the implementation via constructor injection. Tests inject fakes; production injects real implementations.

```python
# Bad — OrderService is hardwired to MySQL; impossible to unit test
class OrderService:
    def __init__(self):
        self.db = MySQLDatabase()    # concrete dependency baked in

    def place_order(self, order):
        self.db.save(order)

# Good — depend on abstraction; inject the implementation
class OrderRepository(ABC):
    @abstractmethod
    def save(self, order): ...

class OrderService:
    def __init__(self, repo: OrderRepository):   # injected
        self.repo = repo

    def place_order(self, order):
        self.repo.save(order)

# Test injects fake; production injects real
class FakeRepository(OrderRepository):
    def save(self, order): self.saved = order

service = OrderService(FakeRepository())   # no DB needed
```

---

### 6. Apply in code review

Check each principle as a dedicated pass:

| Principle | Code review question |
|-----------|---------------------|
| S | Does this class have more than one reason to change? |
| O | Does adding behavior require editing existing logic? |
| L | Can every subtype be used everywhere the base type is used? |
| I | Does any implementor have empty or stub methods? |
| D | Does any class construct its own dependencies? |

## When NOT to Use

SOLID has a cost: more files, more abstractions, more indirection. Skip it for:

- **Scripts and one-off tools** — 50-line script that runs once; no extension needed
- **Data classes / value objects** — pure data holders with no behavior; SRP doesn't apply
- **Prototypes** — defer until the design stabilizes; premature abstraction is waste
- **Small functions** — SOLID is OOP class design guidance; pure functions have different rules

The test: "Will this code be changed by multiple people for different reasons over 6+ months?" If yes, apply SOLID. If no, skip.
