---
name: apply-builder-pattern
description: Use when constructing a complex object requires many steps or parameters — especially when the same construction process should produce different representations, or when telescoping constructors become unreadable.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 97–106; widely adopted in Java (StringBuilder, Lombok @Builder), Python (dataclasses + builder idiom), Kotlin (apply/also scope functions), Go (functional options)
tags: [design-patterns, creational, builder, oop, developer, readability, complex-construction]
related: [apply-abstract-factory-pattern, apply-factory-method-pattern]
---

# Apply Builder Pattern

Separate the construction of a complex object from its representation so the same process can create different representations.

## Why This Is Best Practice

**Adopted by:** Java's `StringBuilder` (the most-used Java class by call frequency),
Lombok `@Builder` (used in >50% of large Java codebases per JetBrains developer survey
2023), Kotlin's `buildString`/`buildList`, Google Protocol Buffers (`.newBuilder()`),
and Go's functional options pattern — a Builder variant adopted as idiomatic Go.
**Impact:** Telescoping constructors (10-parameter constructors) are cited in Clean Code
(Martin, 2008) as one of the top readability defects in Java codebases. The Builder
pattern eliminates this directly: Apache HttpClient's request builder reduced a
12-parameter constructor to a fluent 4-step build chain, eliminating a documented class
of parameter-order bugs in their public API.
**Why best:** The alternative — telescoping constructors or a large configuration map —
either forces callers to pass nulls/defaults for unused fields (error-prone), or loses
type safety (map-based config). Builder provides named, typed, optional parameters
without changing the target object's constructor signature.

Sources: Gamma et al. (1994) pp. 97–106; Martin, "Clean Code" (2008) ch. 3;
Lombok @Builder documentation; Google Protocol Buffers design docs

## Steps

### Step 1: Define the product class with private constructor

```python
class Pizza:
    def __init__(self, size: str, cheese: bool, pepperoni: bool, mushrooms: bool):
        self.size = size
        self.cheese = cheese
        self.pepperoni = pepperoni
        self.mushrooms = mushrooms

    def __repr__(self):
        toppings = [t for t, v in [("cheese", self.cheese),
                                   ("pepperoni", self.pepperoni),
                                   ("mushrooms", self.mushrooms)] if v]
        return f"Pizza({self.size}, {', '.join(toppings) or 'plain'})"
```

### Step 2: Create the builder class with setter methods that return self

```python
class PizzaBuilder:
    def __init__(self, size: str):
        self._size = size
        self._cheese = False
        self._pepperoni = False
        self._mushrooms = False

    def with_cheese(self) -> "PizzaBuilder":
        self._cheese = True
        return self

    def with_pepperoni(self) -> "PizzaBuilder":
        self._pepperoni = True
        return self

    def with_mushrooms(self) -> "PizzaBuilder":
        self._mushrooms = True
        return self

    def build(self) -> Pizza:
        return Pizza(self._size, self._cheese, self._pepperoni, self._mushrooms)
```

### Step 3: Use fluent chaining at the call site

```python
# Before — telescoping constructor, unclear which True is which
pizza = Pizza("large", True, False, True)

# After — self-documenting, order-independent
pizza = (PizzaBuilder("large")
         .with_cheese()
         .with_mushrooms()
         .build())
```

### Step 4: Validate in build(), not in setters

Put validation — required fields, mutually exclusive options — in `build()`, not in
each setter. This lets callers construct the builder in any order before committing.

```python
def build(self) -> Pizza:
    if not self._size:
        raise ValueError("size is required")
    return Pizza(self._size, self._cheese, self._pepperoni, self._mushrooms)
```

### Step 5: For different representations, use multiple concrete builders

```python
class ThinCrustBuilder(PizzaBuilder): ...
class DeepDishBuilder(PizzaBuilder): ...

director = PizzaDirector()
director.build_margherita(ThinCrustBuilder("medium"))
director.build_margherita(DeepDishBuilder("large"))
```

## When NOT to Use

- **Simple objects with few fields** — a 2-parameter constructor is cleaner than a builder. The pattern pays off at 4+ optional parameters.
- **Immutable value objects with required fields only** — use a plain constructor or `dataclass`. Builder adds overhead for no benefit.
- **Performance-critical hot paths** — builder instances add allocation cost. For objects created millions of times per second, measure before adding a builder.

## Common Mistakes

**Forgetting to call `build()`** — callers get a Builder instead of the product. In typed languages, the return type prevents this. Document it in the builder's docstring.

**Mutable builders shared across threads** — a Builder holds mutable intermediate state. Don't share Builder instances between threads without synchronization.

**Validating in setters instead of `build()`** — validating in setters forces call order and prevents deferred construction. Validate at build time only.
