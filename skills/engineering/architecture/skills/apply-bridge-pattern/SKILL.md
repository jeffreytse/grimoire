---
name: apply-bridge-pattern
description: Use when you need to vary both an abstraction and its implementation independently — avoiding an exponential class hierarchy that results from combining them through inheritance.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 151–161; Java AWT (abstraction: Shape hierarchy; implementation: rendering pipeline); Android View system; JDBC (abstraction: Connection; implementation: vendor driver)
tags: [design-patterns, structural, bridge, oop, developer, extensibility, abstraction-implementation]
related: [apply-adapter-pattern, apply-composition-over-inheritance, apply-solid-principles]
---

# Apply Bridge Pattern

Decouple an abstraction from its implementation so that both can vary independently.

## Why This Is Best Practice

**Adopted by:** Java AWT rendering pipeline (Shape abstraction + platform renderer
implementation — ships in every JVM), JDBC (Connection abstraction + vendor-specific
driver implementation — every Java database program uses this), Android's View system
(View abstraction + Canvas implementation), and .NET's Graphics class hierarchy.
**Impact:** GoF shows that without Bridge, combining M abstractions × N implementations
requires M×N classes. With Bridge, it requires M+N classes. A real example: a graphics
library with 3 shapes (Circle, Square, Triangle) × 3 renderers (SVG, Canvas, PDF)
needs 9 classes without Bridge, 6 with it — and adding a 4th renderer requires 3 new
classes vs. 1.
**Why best:** The alternative — inheritance — permanently binds abstraction to
implementation. `SVGCircle extends Circle` cannot use a PDF renderer without a new
subclass. Bridge separates the dimensions, making each independently extensible.

Sources: Gamma et al. (1994) pp. 151–161; Java AWT documentation; JDBC specification

## Steps

### Step 1: Identify the two independent variation dimensions

```
Abstraction dimension: Shape — Circle, Square, Triangle
Implementation dimension: Renderer — SVGRenderer, CanvasRenderer, PDFRenderer
```

Any combination should work without a dedicated class for each.

### Step 2: Define the implementation interface

```python
from abc import ABC, abstractmethod

class Renderer(ABC):
    @abstractmethod
    def render_circle(self, x: float, y: float, radius: float) -> None: ...

    @abstractmethod
    def render_rect(self, x: float, y: float, w: float, h: float) -> None: ...
```

### Step 3: Define the abstraction holding a reference to the implementation

```python
class Shape(ABC):
    def __init__(self, renderer: Renderer):
        self._renderer = renderer   # the bridge

    @abstractmethod
    def draw(self) -> None: ...
```

### Step 4: Extend both hierarchies independently

```python
class Circle(Shape):
    def __init__(self, renderer: Renderer, x: float, y: float, radius: float):
        super().__init__(renderer)
        self._x, self._y, self._radius = x, y, radius

    def draw(self) -> None:
        self._renderer.render_circle(self._x, self._y, self._radius)

class SVGRenderer(Renderer):
    def render_circle(self, x, y, radius):
        print(f'<circle cx="{x}" cy="{y}" r="{radius}"/>')

    def render_rect(self, x, y, w, h):
        print(f'<rect x="{x}" y="{y}" width="{w}" height="{h}"/>')
```

### Step 5: Compose at runtime — any abstraction with any implementation

```python
svg = SVGRenderer()
pdf = PDFRenderer()

Circle(svg, 10, 10, 5).draw()   # SVG circle
Circle(pdf, 10, 10, 5).draw()   # PDF circle — same Circle class, different renderer
```

Adding a new renderer requires one new `Renderer` subclass. Adding a new shape
requires one new `Shape` subclass. Neither side affects the other.

## When NOT to Use

- **When only one dimension varies** — if the implementation never changes, composition or a simple strategy suffices. Bridge adds abstraction overhead without payoff.
- **When both hierarchies are shallow (1–2 subclasses each)** — the simplification only pays off at 3+ variants per dimension.

## Common Mistakes

**Confusing Bridge with Adapter.** Adapter makes incompatible interfaces work together after the fact. Bridge separates abstraction from implementation from the start, by design. If you're retrofitting, it's Adapter.

**Putting abstraction logic in the implementation.** The implementation interface should only expose primitives (render_circle, render_rect). If it contains domain logic, the bridge leaks.

**Over-engineering for a single renderer.** If there's only one implementation and no plan to add more, Bridge is premature. Add it when the second implementation dimension appears.
