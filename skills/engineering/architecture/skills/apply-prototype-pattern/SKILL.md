---
name: apply-prototype-pattern
description: Use when creating a new object by copying an existing instance is cheaper or simpler than constructing from scratch — especially when object initialization is expensive or when the exact type of object to create is determined at runtime.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 117–126; Python built-in copy module (copy/deepcopy); JavaScript Object.create() and spread syntax; Java Cloneable interface
tags: [design-patterns, creational, prototype, oop, developer, cloning, performance]
related: [apply-abstract-factory-pattern, apply-builder-pattern]
---

# Apply Prototype Pattern

Create new objects by copying an existing prototype rather than constructing from scratch.

## Why This Is Best Practice

**Adopted by:** JavaScript (the prototype-based object model is the language's core
inheritance mechanism; `Object.create()` is a direct implementation), Python (`copy`
module's `copy()` and `deepcopy()` are standard library prototype operations), Java
(`Cloneable` interface and `Object.clone()`), and Unity game engine (prefab
instantiation is a prototype clone operation at the engine level).
**Impact:** GoF documents the prototype's primary benefit: when the cost of creating an
object is significantly greater than the cost of copying one, prototypes reduce
construction cost to memcopy speed. In game engines (Unity), instantiating thousands of
enemies from a single prefab prototype is an order-of-magnitude faster than constructing
each from scratch with physics body initialization.
**Why best:** The alternative — reconstructing the full object state every time — is
wasteful when the state was expensive to compute. The prototype also eliminates
coupling to constructor signatures: the client copies the object without knowing its
concrete class.

Sources: Gamma et al. (1994) pp. 117–126; Python `copy` module documentation;
Unity Manual: Instantiating Prefabs

## Steps

### Step 1: Add a `clone()` method to the prototype interface

```python
from abc import ABC, abstractmethod
import copy

class Prototype(ABC):
    @abstractmethod
    def clone(self) -> "Prototype": ...
```

### Step 2: Implement `clone()` using shallow or deep copy

```python
class GameEnemy(Prototype):
    def __init__(self, name: str, health: int, position: tuple, inventory: list):
        self.name = name
        self.health = health
        self.position = position
        self.inventory = inventory

    def clone(self) -> "GameEnemy":
        return copy.deepcopy(self)   # deep copy: inventory list is independent
```

Use `deepcopy` when the prototype contains mutable nested objects that each clone
should own independently. Use `copy` (shallow) only when nested objects are shared
by design or are immutable.

### Step 3: Create a prototype registry for named prototypes

```python
class EnemyRegistry:
    def __init__(self):
        self._prototypes: dict[str, GameEnemy] = {}

    def register(self, name: str, prototype: GameEnemy) -> None:
        self._prototypes[name] = prototype

    def create(self, name: str) -> GameEnemy:
        prototype = self._prototypes.get(name)
        if prototype is None:
            raise KeyError(f"no prototype registered for '{name}'")
        return prototype.clone()

registry = EnemyRegistry()
registry.register("orc", GameEnemy("Orc", 100, (0, 0), ["sword"]))

# Later — create many orcs cheaply
orc1 = registry.create("orc")
orc2 = registry.create("orc")
```

### Step 4: Mutate the clone, not the prototype

After cloning, modify only the clone's instance-specific state:

```python
orc1.position = (10, 5)
orc1.health = 80   # damaged variant
# registry's prototype is unchanged
```

## When NOT to Use

- **Simple objects with cheap construction** — if `new Object()` is nearly free, cloning adds complexity with no benefit.
- **Objects with non-copyable resources** — file handles, network connections, and database connections cannot be meaningfully cloned. Use a factory instead.
- **Deep object graphs with circular references** — deepcopy handles these but can be slow and surprising. Test explicitly.

## Common Mistakes

**Shallow copy when deep copy is needed.** If the prototype holds a list, a shallow copy shares the list between prototype and clone. Mutating the clone's list mutates the prototype. Use `deepcopy` unless sharing is intentional.

**Modifying the prototype instead of the clone.** The registry prototype is the template. Any mutation to it affects all future clones. Treat registered prototypes as read-only after registration.

**Implementing `clone()` by calling the constructor.** Calling `ConcreteProduct(self.field1, self.field2, ...)` in `clone()` defeats the pattern — it re-couples `clone()` to the constructor signature. Use `copy.deepcopy(self)` or equivalent.
