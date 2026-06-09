---
name: apply-template-method-pattern
description: Use when multiple classes share the same algorithm skeleton but differ in specific steps — defining the invariant structure once in a base class and letting subclasses override only the varying parts.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 325–330; Java AbstractList (size() and get() are abstract); Django class-based views; JUnit lifecycle hooks (setUp/tearDown)
tags: [design-patterns, behavioral, template-method, oop, developer, code-reuse, algorithm-skeleton]
related: [apply-strategy-pattern, apply-inheritance-correctly, apply-solid-principles]
---

# Apply Template Method Pattern

Define the skeleton of an algorithm in a base class, deferring specific steps to subclasses — so subclasses can redefine parts of the algorithm without changing its overall structure.

## Why This Is Best Practice

**Adopted by:** Java's `AbstractList` (implements the full `List` contract based on
just `size()` and `get()` — subclasses provide those two methods, inheriting 30+ free
implementations), Django's class-based views (`View.dispatch()` is the template method
— subclasses override `get()`, `post()`, etc.), JUnit's test lifecycle
(`setUp()`/`tearDown()` are hook methods in the template), and Java's `HttpServlet`
(`service()` dispatches to overridable `doGet()`/`doPost()`).
**Impact:** GoF documents that Template Method is one of the most fundamental techniques
for code reuse via inheritance. Django's class-based views reduce view code by 60–80%
for standard CRUD operations — because the template (`dispatch`, `render`) is in the
base class and only the business logic step changes.
**Why best:** Code duplication — the alternative — copies the algorithm skeleton into
every subclass. A bug or enhancement to the skeleton requires editing every copy.
Template Method edits once; all subclasses inherit the fix.

Sources: Gamma et al. (1994) pp. 325–330; Django class-based views documentation;
Java `AbstractList` source

## Steps

### Step 1: Identify the algorithm skeleton and the varying steps

```
Skeleton (invariant): read data → process → format output → write result
Varying steps:        how to read, how to process, how to format
```

### Step 2: Implement the template method in the base class — call abstract steps

```python
from abc import ABC, abstractmethod

class ReportGenerator(ABC):
    def generate(self, data_source: str) -> str:
        """Template method — defines the algorithm skeleton."""
        raw_data = self._read_data(data_source)
        processed = self._process(raw_data)
        formatted = self._format(processed)
        self._write(formatted)
        return formatted

    @abstractmethod
    def _read_data(self, source: str) -> list: ...

    @abstractmethod
    def _process(self, data: list) -> dict: ...

    @abstractmethod
    def _format(self, data: dict) -> str: ...

    def _write(self, output: str) -> None:
        """Hook — optional override. Default: print."""
        print(output)
```

### Step 3: Implement concrete subclasses — override only the varying steps

```python
class CSVReportGenerator(ReportGenerator):
    def _read_data(self, source: str) -> list:
        import csv
        with open(source) as f:
            return list(csv.reader(f))

    def _process(self, data: list) -> dict:
        return {"rows": len(data), "cols": len(data[0]) if data else 0}

    def _format(self, data: dict) -> str:
        return f"CSV Report: {data['rows']} rows × {data['cols']} columns"

class JSONReportGenerator(ReportGenerator):
    def _read_data(self, source: str) -> list:
        import json
        with open(source) as f:
            return json.load(f)

    def _process(self, data: list) -> dict:
        return {"count": len(data)}

    def _format(self, data: dict) -> str:
        return f"JSON Report: {data['count']} records"
```

### Step 4: Distinguish abstract steps from hooks

- **Abstract steps** (marked `@abstractmethod`): must be overridden — they have no default
- **Hooks** (concrete but empty or no-op): optionally overridden — they extend behavior without replacing it

```python
def _write(self, output: str) -> None:
    print(output)   # hook — override to write to file, send email, etc.
```

### Step 5: Mark the template method as final (or document it as non-overridable)

The template method defines the algorithm contract. Subclasses should override only the steps, not the skeleton.

```python
def generate(self, data_source: str) -> str:
    # Python has no `final` keyword — document explicitly
    # Subclasses: override _read_data, _process, _format, _write only
    ...
```

## When NOT to Use

- **When subclasses override the template method itself** — this defeats the pattern. If subclasses need full algorithm control, `apply-strategy-pattern` is correct (composition over inheritance).
- **Deep inheritance hierarchies** — Template Method encourages subclassing. If the hierarchy grows beyond 2–3 levels, composition becomes harder to navigate. Prefer Strategy at that point.
- **When "varying steps" are independent of each other** — if steps A and B vary independently, Strategy with two injected behaviors is more flexible than a subclass hierarchy.

## Common Mistakes

**Subclass overriding the template method.** If `CSVReportGenerator.generate()` overrides `ReportGenerator.generate()`, the skeleton is broken. The template method is the invariant; override only the abstract steps.

**Too many abstract steps.** If every step is abstract, there's no shared skeleton — subclasses duplicate the orchestration. Keep the skeleton meaningful (3–6 steps); if one step always needs overriding alone, consider a Strategy for that step.

**Confusing Template Method with Strategy.** Template Method uses inheritance — subclass IS-A base. Strategy uses composition — context HAS-A strategy. If the "algorithm" varies per instance at runtime, use Strategy. If it varies per class, use Template Method.
