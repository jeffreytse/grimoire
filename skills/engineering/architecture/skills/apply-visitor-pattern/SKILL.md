---
name: apply-visitor-pattern
description: Use when you need to perform many distinct operations on a stable object structure — adding new operations without modifying the element classes by separating the operation from the data structure.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 331–344; Java compiler AST visitors; Python ast.NodeVisitor; LLVM IR pass infrastructure; checkstyle/PMD code analysis
tags: [design-patterns, behavioral, visitor, oop, developer, open-closed, operation-extension, ast]
related: [apply-composite-pattern, apply-iterator-pattern, apply-solid-principles]
---

# Apply Visitor Pattern

Represent an operation to be performed on elements of an object structure — letting you add new operations without changing the element classes.

## Why This Is Best Practice

**Adopted by:** Java compiler's AST visitors (the `javax.lang.model` API uses Visitor
for all annotation processing — in every Java build), Python's `ast.NodeVisitor` (the
foundation of Flake8, mypy, Bandit, and all Python static analysis tools), LLVM's IR
pass infrastructure (every optimization is a Visitor over the IR — the core of Clang,
Rust's compiler backend, and Swift), and PMD/Checkstyle (code analysis tools visited
80M+ downloads).
**Impact:** GoF documents that without Visitor, adding a new operation to an object
structure requires modifying every element class. Python's `ast.NodeVisitor` allows any
tool to add new analysis operations (type checking, linting, security scanning) without
touching the AST node classes. The entire Python static analysis ecosystem builds on
this single design decision.
**Why best:** Adding operations to each class directly — the alternative — violates OCP
and bloats element classes with unrelated operations. Visitor keeps element classes
focused on data representation; operations are concentrated in visitor classes.

Sources: Gamma et al. (1994) pp. 331–344; Python `ast.NodeVisitor` documentation;
LLVM Pass documentation

## Steps

### Step 1: Define the visitor interface with a method for each element type

```python
from abc import ABC, abstractmethod

class ShapeVisitor(ABC):
    @abstractmethod
    def visit_circle(self, circle: "Circle") -> None: ...

    @abstractmethod
    def visit_rectangle(self, rectangle: "Rectangle") -> None: ...

    @abstractmethod
    def visit_triangle(self, triangle: "Triangle") -> None: ...
```

### Step 2: Add `accept()` to element classes — double dispatch

```python
class Shape(ABC):
    @abstractmethod
    def accept(self, visitor: ShapeVisitor) -> None: ...

class Circle(Shape):
    def __init__(self, radius: float):
        self.radius = radius

    def accept(self, visitor: ShapeVisitor) -> None:
        visitor.visit_circle(self)   # double dispatch: runtime type of both self and visitor

class Rectangle(Shape):
    def __init__(self, width: float, height: float):
        self.width = width
        self.height = height

    def accept(self, visitor: ShapeVisitor) -> None:
        visitor.visit_rectangle(self)
```

### Step 3: Implement concrete visitors — each adds one operation across all elements

```python
class AreaCalculator(ShapeVisitor):
    def __init__(self):
        self.total_area = 0.0

    def visit_circle(self, circle: Circle) -> None:
        import math
        self.total_area += math.pi * circle.radius ** 2

    def visit_rectangle(self, rectangle: Rectangle) -> None:
        self.total_area += rectangle.width * rectangle.height

    def visit_triangle(self, triangle: "Triangle") -> None:
        self.total_area += 0.5 * triangle.base * triangle.height

class SVGExporter(ShapeVisitor):
    def __init__(self):
        self.output = []

    def visit_circle(self, circle: Circle) -> None:
        self.output.append(f'<circle r="{circle.radius}"/>')

    def visit_rectangle(self, rectangle: Rectangle) -> None:
        self.output.append(f'<rect width="{rectangle.width}" height="{rectangle.height}"/>')

    def visit_triangle(self, triangle: "Triangle") -> None:
        self.output.append(f'<polygon points="..."/>')
```

### Step 4: Traverse the structure — apply any visitor to any element

```python
shapes: list[Shape] = [Circle(5), Rectangle(3, 4), Circle(2)]

area_calc = AreaCalculator()
for shape in shapes:
    shape.accept(area_calc)
print(f"Total area: {area_calc.total_area:.2f}")

exporter = SVGExporter()
for shape in shapes:
    shape.accept(exporter)
print("\n".join(exporter.output))
```

Adding a new operation (e.g., `JSONExporter`) requires one new Visitor class — no
changes to `Circle`, `Rectangle`, or `Triangle`.

### Step 5: For recursive structures (trees), traverse in the visitor

```python
class PrintVisitor(ast.NodeVisitor):
    def visit_FunctionDef(self, node):
        print(f"Function: {node.name}")
        self.generic_visit(node)   # recurse into children
```

`generic_visit` (Python's built-in) handles recursion, so the visitor focuses on
the operation, not the traversal.

## When NOT to Use

- **When the element class structure changes frequently** — adding a new element type requires updating every existing visitor. Visitor is best when element classes are stable and operations vary.
- **When elements have very different interfaces** — if the visitor's `visit_X` methods can't access what they need through the element's public interface, the visitor forces breaking encapsulation.

## Common Mistakes

**Forgetting `accept()` on a new element class.** If a new element doesn't implement `accept()`, visitors skip it silently. All element classes must implement `accept()`.

**Visitor that builds state from one traversal and uses it in another.** If `AreaCalculator.total_area` is read mid-traversal by another visitor, state is inconsistent. Each traversal should be independent.

**Using Visitor when Iterator suffices.** If all you need is sequential access to elements without type-specific handling, `apply-iterator-pattern` is simpler. Visitor pays off when behavior differs meaningfully per element type.
