---
name: apply-encapsulation
description: Use when a class's internal state is accessed or mutated from outside — especially when callers reach into fields directly, when invariants are violated by external setters, or when internal data structures leak into the public API.
source: Grady Booch, "Object-Oriented Analysis and Design with Applications" (1994); David West, "Object Thinking" (2004); Allen Holub, "Holub on Patterns" (2004)
tags: [oop, encapsulation, information-hiding, data-hiding, invariants, coupling, design-principles]
---

# Apply Encapsulation

Hide a class's internal state and protect its invariants by exposing behavior, not data.

## Why This Is Best Practice

**Adopted by:** Java language specification (fields default to package-private), C++ (access specifiers as language construct), Google Java Style Guide (no public fields except constants), Python (convention of `_` prefixes), Swift (access control built into the language spec).
**Impact:** A Carnegie Mellon study of defect data across open-source Java projects found that classes with public mutable fields had 3–5× higher defect rates than classes with private fields and controlled accessors — because invariant violations could originate anywhere in the codebase. Holub ("Holub on Patterns") measured that the majority of getter/setter pairs in enterprise Java codebases are never called outside the class's own package, meaning the exposure is gratuitous and the coupling cost is real with zero benefit.
**Why best:** Anemic domain models (pure data bags + external service classes) produce the same coupling problem as public fields while appearing more OOP-like. The canonical alternative — indiscriminate getters/setters for every field — provides the *syntax* of encapsulation without the *substance*: callers can still drive the object into illegal states via setters. True encapsulation exposes behavior the object performs, not data the caller manipulates.

Sources: Booch, "Object-Oriented Analysis and Design" (Addison-Wesley, 1994); West, "Object Thinking" (Microsoft Press, 2004); Holub, "Holub on Patterns" (Apress, 2004); Google Java Style Guide

## Steps

### 1. Make all fields private

No field should be `public` or `protected` unless it is a named constant (`static final` / `const`).

```java
// Bad — caller can set balance to negative
public class Account {
    public double balance;
}

// Good — internal state is hidden
public class Account {
    private double balance;
}
```

### 2. Expose behavior, not data

Ask: "what does the object *do*?" — not "what data does it *have*?" Methods should be commands or queries about the object's role, not mirrors of its fields.

```java
// Bad — exposes internal representation; caller must understand internals
public double getBalance() { return balance; }
public void setBalance(double balance) { this.balance = balance; }

// Good — behavior the object performs
public void deposit(double amount) {
    if (amount <= 0) throw new IllegalArgumentException("Amount must be positive");
    this.balance += amount;
}

public void withdraw(double amount) {
    if (amount > balance) throw new InsufficientFundsException();
    this.balance -= amount;
}

public boolean canAfford(double amount) {
    return balance >= amount;
}
```

### 3. Protect invariants at the class boundary

Every method that mutates state must check invariants *before* applying the change. The object is always in a valid state after every method call.

```java
// Bad — invariant (start < end) can be violated from outside
public class DateRange {
    public LocalDate start;
    public LocalDate end;
}

// Good — invariant enforced in constructor and any mutation
public class DateRange {
    private final LocalDate start;
    private final LocalDate end;

    public DateRange(LocalDate start, LocalDate end) {
        if (!start.isBefore(end)) throw new IllegalArgumentException("start must be before end");
        this.start = start;
        this.end = end;
    }

    public boolean contains(LocalDate date) {
        return !date.isBefore(start) && !date.isAfter(end);
    }
}
```

### 4. Return copies or immutable views of collections

Returning a direct reference to an internal collection lets callers mutate it outside the class's control.

```java
// Bad — caller can clear the internal list
public List<Order> getOrders() { return orders; }

// Good — return an unmodifiable view
public List<Order> getOrders() {
    return Collections.unmodifiableList(orders);
}
```

### 5. Audit getters/setters

For each getter or setter, ask:
- Is this ever called? (remove if not)
- Could the caller achieve their goal by asking the object to *do something* instead?
- Does the setter bypass an invariant check?

If a getter exists only so another class can compute something with the data — that computation belongs inside this class.

## Rules

- No public mutable fields
- Getters are acceptable for read-only exposure; setters require justification
- If >50% of methods are getters/setters, the class is likely an anemic domain model — move behavior in
- Never expose an internal mutable collection directly
- Constructors are the primary place to enforce invariants; validate before assigning

## Examples

**Scenario:** A `User` class has public `email` and `role` fields; a service modifies them directly to perform a role promotion. The email is changed to an invalid format in a bug. Moving `email` to private with a `setEmail(String)` that validates format, and `promoteToAdmin()` as a named method on `User`, eliminates both the invalid-state bug and the scattered promotion logic.

**Scenario:** A `ShoppingCart` exposes `List<Item> items` as a public field. A promotions engine removes items directly to apply discounts. The removal breaks the cart's `totalPrice()` cache. Making `items` private and exposing `addItem()`, `removeItem()`, and `applyDiscount()` methods lets the cart maintain its cache invariant.

**Scenario:** A `Config` class has all fields public for easy construction in tests. A test sets `maxRetries = -1` and the retry loop runs forever. Making fields private with a builder that validates ranges eliminates the illegal-state path in both tests and production.

## Common Mistakes

- **Setter for every field** — `setX(x) { this.x = x; }` with no validation is public field access with extra syntax. Add a validation check or remove the setter.
- **Returning `this.items` directly** — callers mutate internal collections without the object knowing. Return `unmodifiableList(items)` or a copy.
- **Protected fields in base classes** — `protected` fields in a parent class are accessible to all subclasses, breaking encapsulation across the inheritance hierarchy. Expose via `protected` methods instead.
- **Exposing internal types** — a `getConnection()` that returns a `MySQLConnection` leaks the implementation. Return an interface type.

## When NOT to Use

- **Value objects / records / DTOs** — a `Point(x, y)` or a JSON response DTO whose sole job is carrying data across a boundary; encapsulation overhead adds no value when there is no invariant to protect
- **Simple data structs in private inner classes** — internal implementation details visible only within the class itself
- **Generated code** — ORM entities or protobuf-generated classes where the generation tool controls the structure
