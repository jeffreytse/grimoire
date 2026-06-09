---
name: apply-decorator-pattern
description: Use when you need to add responsibilities to individual objects dynamically — without modifying the class and without the combinatorial explosion that results from subclassing every combination of features.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 175–184; Java I/O streams (BufferedInputStream, GZIPOutputStream); Python @functools.wraps; Django middleware stack
tags: [design-patterns, structural, decorator, oop, developer, extensibility, dynamic-behavior]
related: [apply-composition-over-inheritance, apply-proxy-pattern, apply-solid-principles]
---

# Apply Decorator Pattern

Attach additional responsibilities to an object dynamically by wrapping it in decorator objects that share the same interface.

## Why This Is Best Practice

**Adopted by:** Java I/O (`BufferedInputStream(new FileInputStream(path))` is a
decorator chain — used in virtually every Java file operation), Python's `@functools.wraps`
and the entire decorator syntax (`@`, `@login_required`, `@cache` — one of Python's most
distinctive features), Django's middleware stack (each middleware wraps the next,
adding auth, CSRF, session handling), and JavaScript's `Array.prototype.reduce`
-based pipeline patterns.
**Impact:** Java's I/O library design is cited in GoF as the canonical Decorator
example. Without Decorator, adding buffering + gzip + encryption to file I/O would
require 3×3=9 dedicated subclasses for each combination. With Decorator, any combination
is achieved by wrapping: 3 decorators compose into any of 9 behaviors with 3 classes.
**Why best:** Subclassing for each combination is the alternative — it produces
exponential class growth (M features × N base types). Decorator provides the same
flexibility from M+N classes with runtime composition.

Sources: Gamma et al. (1994) pp. 175–184; Python `functools` documentation;
Django middleware documentation

## Steps

### Step 1: Define the component interface

```python
from abc import ABC, abstractmethod

class TextProcessor(ABC):
    @abstractmethod
    def process(self, text: str) -> str: ...
```

### Step 2: Implement the concrete component (base behavior)

```python
class PlainTextProcessor(TextProcessor):
    def process(self, text: str) -> str:
        return text
```

### Step 3: Write the base decorator — wraps a component, delegates by default

```python
class TextProcessorDecorator(TextProcessor):
    def __init__(self, wrapped: TextProcessor):
        self._wrapped = wrapped

    def process(self, text: str) -> str:
        return self._wrapped.process(text)
```

### Step 4: Extend the base decorator — add behavior before/after delegation

```python
class UpperCaseDecorator(TextProcessorDecorator):
    def process(self, text: str) -> str:
        return super().process(text).upper()

class TrimDecorator(TextProcessorDecorator):
    def process(self, text: str) -> str:
        return super().process(text).strip()

class PrefixDecorator(TextProcessorDecorator):
    def __init__(self, wrapped: TextProcessor, prefix: str):
        super().__init__(wrapped)
        self._prefix = prefix

    def process(self, text: str) -> str:
        return f"{self._prefix}{super().process(text)}"
```

### Step 5: Compose decorators at the call site — any combination, any order

```python
processor = PrefixDecorator(
    UpperCaseDecorator(
        TrimDecorator(PlainTextProcessor())
    ),
    prefix=">> "
)

print(processor.process("  hello world  "))
# >> HELLO WORLD
```

Order matters: decorators apply from innermost to outermost. Trim before UpperCase
is different from UpperCase before Trim only if case affects whitespace — but always
be explicit about order.

## When NOT to Use

- **When behavior is fixed at design time** — if the combination of features never changes at runtime, a single concrete class is simpler.
- **When decorator chains become deep** — chains of 5+ decorators are hard to debug. If the wrapping count grows unbounded, consider a pipeline or middleware pattern instead.
- **In languages with native function decoration** — Python's `@` syntax and `functools` handle function-level decoration natively. Use them directly rather than building class-based decorators for functions.

## Common Mistakes

**Decorator that doesn't call `super().process()`** — a decorator that forgets to delegate completely replaces the wrapped behavior instead of extending it. Every decorator must call its wrapped component unless replacement is intentional.

**Breaking the interface contract.** A decorator may not add, remove, or change the signature of methods on the interface. It wraps, not transforms. Interface changes belong in Adapter.

**Using Decorator when Strategy is the right tool.** If the variation is selecting one algorithm from several alternatives (not stacking multiple behaviors), `apply-strategy-pattern` is cleaner.
