---
name: apply-flyweight-pattern
description: Use when an application creates a very large number of fine-grained objects whose combined memory cost is prohibitive — by sharing the intrinsic (immutable) state across instances and storing only extrinsic (context-specific) state outside.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 195–206; Java String pool (String.intern()); Python integer and string interning (sys.intern); game engine particle systems; text rendering glyph caches
tags: [design-patterns, structural, flyweight, oop, developer, memory-optimization, large-scale-objects]
related: [apply-prototype-pattern, apply-singleton-pattern]
---

# Apply Flyweight Pattern

Share immutable intrinsic state across many fine-grained objects to reduce memory usage when the number of object instances is very large.

## Why This Is Best Practice

**Adopted by:** Java's String pool (`String.intern()` — every JVM uses it for string
literals, reducing string object count by 60–80% in typical applications), Python's
integer interning (integers -5 to 256 are singletons in CPython), Java's `Character`
cache (characters 0–127), and every major game engine's particle system (10,000+
particles sharing a single mesh/texture object).
**Impact:** GoF documents a text editor example where storing one glyph object per
character in a 100,000-character document would require 100,000 objects. With Flyweight,
one object per distinct glyph character (typically ~100 unique glyphs) serves all
100,000 positions — a 1000× reduction in object count.
**Why best:** The alternative — one full object per instance — is correct and simple
but prohibitive at scale. Flyweight applies only when profiling confirms object count is
the memory bottleneck; it trades simplicity for a controlled reduction in allocation cost.

Sources: Gamma et al. (1994) pp. 195–206; Java String pool documentation;
CPython source (`Objects/longobject.c` integer cache)

## Steps

### Step 1: Separate intrinsic (shared) from extrinsic (context-specific) state

```
Intrinsic state (shared): character glyph image, font, size
Extrinsic state (per-instance): position on page (x, y), color

Particle system:
Intrinsic: mesh, texture, animation frames
Extrinsic: position, velocity, lifetime
```

Intrinsic state must be immutable — it's shared across all instances that use it.

### Step 2: Make the flyweight class hold only intrinsic state

```python
class GlyphFlyweight:
    def __init__(self, character: str, font: str, size: int):
        self.character = character   # intrinsic — shared
        self.font = font
        self.size = size

    def render(self, x: int, y: int, color: str) -> None:
        # x, y, color are extrinsic — passed in, not stored
        print(f"Render '{self.character}' at ({x},{y}) in {color} [{self.font} {self.size}pt]")
```

### Step 3: Build a factory that returns shared flyweight instances

```python
class GlyphFactory:
    _pool: dict[tuple, GlyphFlyweight] = {}

    @classmethod
    def get(cls, character: str, font: str, size: int) -> GlyphFlyweight:
        key = (character, font, size)
        if key not in cls._pool:
            cls._pool[key] = GlyphFlyweight(character, font, size)
        return cls._pool[key]
```

### Step 4: Store extrinsic state outside the flyweight — in the client context

```python
class TextDocument:
    def __init__(self):
        self._chars: list[tuple[GlyphFlyweight, int, int, str]] = []

    def add_character(self, ch: str, font: str, size: int, x: int, y: int, color: str):
        glyph = GlyphFactory.get(ch, font, size)   # shared flyweight
        self._chars.append((glyph, x, y, color))   # extrinsic state stored here

    def render(self):
        for glyph, x, y, color in self._chars:
            glyph.render(x, y, color)
```

### Step 5: Measure first — only apply where profiling shows object count is the bottleneck

```python
import sys
# Before: 100,000 GlyphFlyweight objects
# After:  ~95 unique glyphs in pool
print(f"Pool size: {len(GlyphFactory._pool)}")
print(f"Pool memory: {sum(sys.getsizeof(g) for g in GlyphFactory._pool.values())} bytes")
```

## When NOT to Use

- **When object count is small** — if an application creates hundreds, not tens of thousands of objects, the pattern adds complexity with no measurable benefit.
- **When extrinsic state is large** — if the per-instance state dominates memory usage, sharing the intrinsic portion saves little. Profile before applying.
- **When objects must maintain identity** — if callers compare objects by identity (`is`) rather than value, sharing breaks expected behavior.

## Common Mistakes

**Storing mutable state in the flyweight.** If the "shared" state changes per use, it is extrinsic, not intrinsic. Mutable flyweight state causes one caller's changes to silently affect all others sharing the same object.

**Not using a factory.** Flyweights only save memory if they're actually shared. Without a factory/pool, each caller creates a new instance and the pattern achieves nothing.

**Applying prematurely.** Flyweight is a memory optimization, not an architectural pattern. Apply it after profiling confirms object count is the problem, not before.
