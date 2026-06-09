---
name: apply-law-of-demeter
description: Use when a class reaches through intermediate objects to access deep properties, when a change to an internal object forces changes in unrelated callers, or when code chains method calls three or more levels deep.
source: Lieberherr & Holland, "Assuring Good Style for Object-Oriented Programs", IEEE Software, 1989
tags: [tight-coupling, train-wreck, law-of-demeter, developer, encapsulation, loose-coupling, maintainability, defect-reduction]
---

# Apply Law of Demeter

A method should only talk to its immediate neighbors: the object itself, its parameters, objects it created, and its direct fields. Don't reach through one object to access another.

## Why This Is Best Practice

**Adopted by:** Google (Engineering Practices documentation), cited as a core OOP discipline by Martin Fowler ("Refactoring", 2nd ed.) and Robert C. Martin ("Clean Code", ch. 6). Codified as a foundational object-oriented design principle in IEEE Software (Lieberherr & Holland, 1989) and taught in software engineering curricula at MIT, CMU, and Stanford.
**Impact:** Lieberherr & Holland (1989, IEEE Software) demonstrated that LoD violations are the primary source of class coupling that spreads change across a codebase — a single internal structural change forces edits in every class that violated LoD against it. Palomba et al. (2018, IEEE TSE) found that classes violating encapsulation and reaching into neighbor internals have 2–3x higher change-coupling frequency, meaning a change in one class statistically predicts a change in LoD-violating callers.
**Why best:** Deep method chains (`a.getB().getC().getD().doSomething()`) couple the caller to the entire internal structure of A, B, C, and D simultaneously. If any intermediate object changes its type, renames a method, or is removed, every chain that passes through it breaks. The alternative — "tell, don't ask" — delegates responsibility to the object that owns the data, keeping callers ignorant of internal structure. Getter chaining and fluent builder APIs are common sources of false LoD violations; the rule applies to structural navigation, not to well-designed fluent interfaces.

Sources: Lieberherr & Holland, "Assuring Good Style for Object-Oriented Programs" (IEEE Software, 1989); Palomba et al., "Revisiting the Impact of Code Smells on Software Maintainability" (IEEE TSE, 2018); Martin, "Clean Code" (Prentice Hall, 2008), ch. 6

## Steps

### 1. Identify the violation: the train wreck

A train wreck is a chain of calls navigating through intermediate objects to reach deeply nested data.

```python
# Violation — caller knows about Order, Customer, Address, and City simultaneously
city = order.get_customer().get_address().get_city()

# Violation — three levels of structural knowledge baked into one line
discount = cart.get_user().get_membership().get_discount_rate()
```

**Rule of thumb:** More than one `.` accessing a different object type = LoD candidate. (Fluent chains on the same type — `"hello".strip().upper().split()` — are not LoD violations.)

---

### 2. Apply "tell, don't ask"

Instead of asking an object for its internals and acting on them, tell the object to do the work itself.

```python
# Bad — caller asks for internals, then acts on them
if order.get_customer().get_membership().is_premium():
    price = price * 0.9

# Good — tell Order to apply its own discount; caller doesn't know about Membership
price = order.apply_membership_discount(price)
```

The object that owns the data is the right place for the logic. Move the logic to it, not the data to the caller.

---

### 3. Add delegation methods to the immediate owner

When callers repeatedly reach through an object to get the same nested value, add a method on the intermediate object that delegates.

```python
# Bad — every caller navigates Customer → Address → City
def format_shipping_label(order):
    city = order.get_customer().get_address().get_city()
    zip_code = order.get_customer().get_address().get_zip()
    return f"{city}, {zip_code}"

# Good — Order delegates to Customer; caller only knows Order
class Order:
    def shipping_city(self):
        return self.customer.address.city      # internal to Order

    def shipping_zip(self):
        return self.customer.address.zip_code  # internal to Order

def format_shipping_label(order):
    return f"{order.shipping_city()}, {order.shipping_zip()}"
```

---

### 4. Inject the dependency directly when delegation is awkward

If the caller genuinely needs the nested object, inject it directly rather than navigating to it.

```python
# Bad — Service navigates through Repository to get the config it needs
class ReportService:
    def __init__(self, repo):
        self.repo = repo

    def generate(self):
        format = self.repo.get_config().get_report_format()  # LoD violation

# Good — inject the config directly; ReportService doesn't know about Repository internals
class ReportService:
    def __init__(self, config):
        self.config = config

    def generate(self):
        format = self.config.report_format   # direct dependency, no navigation
```

---

### 5. Distinguish structural navigation from fluent APIs

LoD applies to structural navigation across object boundaries. It does not apply to:

| Pattern | LoD applies? |
|---|---|
| `order.customer.address.city` | Yes — navigating through object ownership |
| `"text".strip().upper()` | No — fluent string methods on same type |
| `queryset.filter(active=True).order_by("name")` | No — builder/query DSL on same logical object |
| `response.json()["data"]["user"]["id"]` | Yes — navigating through data structure layers |

For builder/fluent APIs, the chain operates on one logical entity. For structural navigation, each `.` crosses into a different object's private territory.

---

### 6. Audit existing code with a grep pattern

Find LoD candidates in a codebase:

```bash
# Python — chains with 3+ method calls (adjust for your language)
grep -rn "\.\w\+()\.\w\+()\.\w\+()" src/ --include="*.py"

# Java/Kotlin
grep -rn "\.\w\+()\.\w\+()\.\w\+()" src/ --include="*.java" --include="*.kt"

# JavaScript/TypeScript
grep -rn "\.\w\+()\.\w\+()\.\w\+()" src/ --include="*.ts" --include="*.js"
```

Triage each hit: structural navigation (fix) vs fluent API (acceptable).

## When NOT to Use

- **Data Transfer Objects (DTOs)** and plain value structs: `address.city` on a pure data container is not a LoD violation — there is no behavior or encapsulation to protect.
- **Fluent builder and query APIs**: `db.query().where().orderBy().limit()` is intentionally chained on a single logical operation. Don't add delegation wrappers just to break the chain.
- **Well-known value object traversal**: `LocalDate.now().plusDays(7).format(formatter)` traverses a value, not an object graph with encapsulated state.
- **When the delegation method would be meaningless**: if `order.shipping_city()` is only ever called once in one place, the delegation method adds indirection without reducing coupling. Apply LoD where violations recur.

## Common Mistakes

**Adding getters to satisfy LoD superficially.** `order.getCustomerCity()` that internally calls `customer.getAddress().getCity()` just moves the violation inside the class. The delegation method should encapsulate a real responsibility, not paper over the chain.

**Confusing LoD with "no dots."** `self.config.timeout` accessing a plain data field on a directly-held object is fine. LoD targets behavioral navigation through object chains, not property access on owned data.

**Over-applying to data pipelines.** ETL, data transformation, and functional pipelines naturally chain operations. Applying LoD here produces artificial wrapping that obscures the pipeline logic without reducing coupling.
