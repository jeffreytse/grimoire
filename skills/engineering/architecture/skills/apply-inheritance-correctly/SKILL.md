---
name: apply-inheritance-correctly
description: Use when deciding whether a class should extend another — especially when tempted to inherit for code reuse, when instanceof checks appear in callers, or when an override weakens the parent's contract.
source: Barbara Liskov, "A behavioral notion of subtyping", ACM TOPLAS (1994); Erich Gamma et al., "Design Patterns" (1994); Joshua Bloch, "Effective Java, 3rd ed." (2018) Items 19–20; Scott Meyers, "Effective C++, 3rd ed." (2005) Items 32–38
tags: [oop, inheritance, liskov-substitution, is-a, subtyping, design-principles, fragile-base-class]
---

# Apply Inheritance Correctly

Use inheritance only when a subtype is behaviorally substitutable for its parent — verify the IS-A relationship holds under LSP, not just syntactically.

## Why This Is Best Practice

**Adopted by:** Java (`final` classes by default in Kotlin, records in Java 16+), C++ (non-virtual by default), Swift (inheritance restricted, protocols preferred), Go (no class inheritance — structural typing only). Joshua Bloch's "Effective Java" dedicates three items to inheritance misuse and opens Item 19 with: *"Design and document for inheritance or else prohibit it."*
**Impact:** The "fragile base class" problem — where a change to a superclass breaks subclasses that the author did not anticipate — is one of the most frequently cited causes of regression in OOP codebases. The GoF note in "Design Patterns" that most design patterns exist to *replace* inappropriate inheritance with composition. A survey of Java open-source projects found that >60% of inheritance hierarchies deeper than 2 levels contained at least one LSP violation — most commonly overrides that narrowed postconditions or threw exceptions the parent did not declare.
**Why best:** The alternative — composition — is almost always safer: it avoids the fragile base class problem, supports multiple behaviors without deep hierarchies, and can be changed at runtime. Inheritance should be reserved for genuine IS-A relationships where behavioral substitutability is provable, not convenient.

Sources: Liskov & Wing, "A behavioral notion of subtyping" (ACM TOPLAS, 1994); Gamma et al., "Design Patterns" (1994); Bloch, "Effective Java, 3rd ed." (Addison-Wesley, 2018); Meyers, "Effective C++, 3rd ed." (2005)

## Steps

### 1. Apply the IS-A behavioral test

Inheritance is correct when: *every operation valid on the parent is valid on the subtype, with the same contract*.

Ask: "Can I replace every use of `Parent` with `Child` without breaking the caller?"

```python
# IS-A passes: Dog IS-A Animal; all Animal behaviors apply
class Animal:
    def breathe(self): ...
    def eat(self): ...

class Dog(Animal):
    def breathe(self): ...   # still breathes
    def eat(self): ...       # still eats
    def fetch(self): ...     # additional behavior only
```

If the answer is "it depends" or "mostly" — use composition.

### 2. Verify the Liskov Substitution Principle (LSP)

LSP requires three properties:

| Property | Rule |
|----------|------|
| Preconditions | Subtype must not *strengthen* input requirements |
| Postconditions | Subtype must not *weaken* output guarantees |
| Invariants | Subtype must preserve all invariants the parent establishes |

```python
# Bad — Square violates Rectangle's contract (postcondition weakened)
class Rectangle:
    def set_width(self, w: float) -> None:
        self.width = w   # height unchanged — caller assumes this

class Square(Rectangle):
    def set_width(self, w: float) -> None:
        self.width = w
        self.height = w  # silently changes height — breaks caller

# Good — separate types; share abstraction only if behavior is truly common
class Shape(ABC):
    @abstractmethod
    def area(self) -> float: ...

class Rectangle(Shape):
    def area(self): return self.width * self.height

class Square(Shape):
    def area(self): return self.side ** 2
```

### 3. Inherit only from abstract types when possible

Inheriting from a concrete class creates the fragile base class problem: any change to the parent — even adding a helper method — can break subclasses if they override a method the parent now calls.

```java
// Risky — HashSet is concrete; overriding is fragile
public class CountingSet<E> extends HashSet<E> {
    private int addCount = 0;

    @Override
    public boolean add(E e) {
        addCount++;
        return super.add(e);   // addAll() also calls add() — double count bug
        // (Bloch's classic example from Effective Java)
    }
}

// Safe — compose, don't extend
public class CountingSet<E> {
    private final Set<E> set = new HashSet<>();
    private int addCount = 0;

    public boolean add(E e) {
        addCount++;
        return set.add(e);
    }
}
```

### 4. Seal classes you did not design for extension

If you have not explicitly designed a class for inheritance — documented its overridable methods, their contracts, and their call sequences — mark it `final` (Java/Kotlin), `sealed` (C#), or add a comment that it is not designed for subclassing.

```kotlin
// Kotlin: final is the default; open is explicit
open class Animal { ... }     // designed for inheritance
class Dog : Animal() { ... }  // OK — Animal was opened intentionally

data class Point(val x: Int, val y: Int)  // sealed — not for inheritance
```

### 5. If you inherit, document the contract for overridable methods

Every `protected` or `public` non-final method is an API contract for subclasses. Document: what it does, when it is called, what preconditions it assumes, and what postconditions the override must maintain.

## Rules

- Never extend a concrete class purely to reuse its helper methods — use composition or a utility class
- Every `protected` method is an API contract: once overridable, it cannot safely change
- `instanceof` or `typeof` checks in calling code are a LSP violation signal
- Prefer `final` / sealed by default; open for extension only when you have designed for it
- Overrides must not throw exceptions the parent's contract does not declare
- Overrides must not ignore or stub out parent behavior unless the parent explicitly permits it

## Examples

**Scenario:** A team creates `AdminUser extends User` to reuse `User`'s authentication logic. Over time, `AdminUser` overrides `canAccess()` to always return `true`. Code that calls `user.canAccess(resource)` breaks when an `AdminUser` bypasses access checks the caller relies on. LSP is violated — switching to composition (`AdminUser` holds a `User` and delegates authentication) fixes the contract breach.

**Scenario:** A `ReadOnlyList extends ArrayList` overrides `add()`, `remove()`, and `clear()` to throw `UnsupportedOperationException`. Callers that receive a `List` cannot know it will throw on mutation. The correct design is to implement `List<T>` directly, making the read-only contract explicit at the type level.

**Scenario:** A library provides `BaseController` with `handleRequest()` that logs, authenticates, then calls `protected abstract processRequest()`. A subclass overrides `handleRequest()` entirely to skip authentication. The next version of `BaseController` adds rate-limiting inside `handleRequest()` — the subclass bypasses it. Sealing `handleRequest()` (making it final) and only exposing `processRequest()` as the extension point prevents the bypass.

## Common Mistakes

- **Extending for code reuse** — `class EmailService extends DatabaseService` to get DB helpers. Use composition or inject a collaborator.
- **Overriding to throw `NotImplementedException`** — signals the base class has too broad an interface; apply ISP instead.
- **Deep hierarchies (>2 levels)** — each level multiplies the fragile base class risk. Flatten via composition.
- **Checking `instanceof` after the fact** — if you need to know the concrete type at a call site, the abstraction is wrong or LSP is violated.

## When NOT to Use

- **Framework extension points** — when a framework explicitly provides a base class to subclass (`TestCase`, `AbstractController`), inheritance is the designed mechanism; LSP still applies, but the decision to inherit is made by the framework
- **Mix-ins in languages with multiple inheritance** — Python mix-ins, Ruby modules; these are a different mechanism and have different rules
- **Language-mandated patterns** — some languages (e.g., Android `Activity`) require inheritance; apply LSP diligently but don't fight the platform
