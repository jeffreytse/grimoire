---
name: apply-composite-pattern
description: Use when you need to represent part-whole hierarchies of objects and want clients to treat individual objects and compositions of objects uniformly — without special-casing leaf vs. container nodes.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 163–173; Java Swing (Component/Container hierarchy); Python ast module; HTML DOM; React component tree
tags: [design-patterns, structural, composite, oop, developer, tree-structure, part-whole-hierarchy]
related: [apply-iterator-pattern, apply-visitor-pattern, apply-solid-principles]
---

# Apply Composite Pattern

Compose objects into tree structures to represent part-whole hierarchies, allowing clients to treat individual objects and compositions uniformly.

## Why This Is Best Practice

**Adopted by:** Java Swing's `Component`/`Container` hierarchy (a Container is a
Component that holds Components — the foundation of every Java GUI), the HTML DOM
(every node is an Element that can contain other Elements), Python's `ast` module
(every AST node is either a leaf or a container node with the same visitor interface),
and React's component tree (every component renders leaf elements or child components
uniformly).
**Impact:** GoF documents that without Composite, client code must distinguish between
leaf and container nodes with instanceof checks at every recursive step. In the Swing
hierarchy, removing that distinction meant GUI layout algorithms could recurse through
any depth of nesting with identical code — no special-casing.
**Why best:** The alternative — separate handling for leaves and containers — produces
client code littered with `if isinstance(node, Container)` checks. These checks must
be updated whenever a new node type is added. Composite removes that coupling entirely.

Sources: Gamma et al. (1994) pp. 163–173; Java Swing API documentation; HTML DOM spec

## Steps

### Step 1: Define the component interface for both leaves and composites

```python
from abc import ABC, abstractmethod

class FileSystemItem(ABC):
    def __init__(self, name: str):
        self.name = name

    @abstractmethod
    def size(self) -> int: ...

    @abstractmethod
    def display(self, indent: int = 0) -> None: ...
```

### Step 2: Implement the leaf — no children

```python
class File(FileSystemItem):
    def __init__(self, name: str, size_bytes: int):
        super().__init__(name)
        self._size = size_bytes

    def size(self) -> int:
        return self._size

    def display(self, indent: int = 0) -> None:
        print(" " * indent + f"📄 {self.name} ({self._size}B)")
```

### Step 3: Implement the composite — holds and delegates to children

```python
class Directory(FileSystemItem):
    def __init__(self, name: str):
        super().__init__(name)
        self._children: list[FileSystemItem] = []

    def add(self, item: FileSystemItem) -> None:
        self._children.append(item)

    def remove(self, item: FileSystemItem) -> None:
        self._children.remove(item)

    def size(self) -> int:
        return sum(child.size() for child in self._children)

    def display(self, indent: int = 0) -> None:
        print(" " * indent + f"📁 {self.name}/")
        for child in self._children:
            child.display(indent + 2)
```

### Step 4: Use the component interface uniformly — no type checks needed

```python
root = Directory("project")
src = Directory("src")
src.add(File("main.py", 1024))
src.add(File("utils.py", 512))
root.add(src)
root.add(File("README.md", 256))

root.display()       # works on any depth
print(root.size())   # 1792 — recursively sums without instanceof checks
```

### Step 5: Decide where to put child management — component or composite only

GoF offers two options:
- **Safety**: put `add()`/`remove()` only on `Composite` — type-safe but forces casts
- **Transparency**: put them on `Component` with default no-ops on `Leaf` — uniform interface but semantically odd

Choose transparency when client code must recurse without knowing types; choose safety when the type distinction matters for correctness.

## When NOT to Use

- **Flat collections** — if the structure is never nested (a list of items with no tree), Composite adds complexity without benefit.
- **When leaf/composite behaviors diverge significantly** — if composites need many methods that leaves can't reasonably implement, the shared interface becomes a leaky abstraction.

## Common Mistakes

**Accessing children through the component interface.** If clients call `node.children` directly, they've bypassed the abstraction. Keep traversal recursive via the component's methods.

**Forgetting to delegate to children in composite operations.** A `size()` that returns `self._size` in the composite (not summing children) is a silent bug. Every composite operation must delegate.

**Building unbounded trees without cycle detection.** Composite trees can have cycles if a directory is added as its own child. Validate on `add()` if the domain requires it.
