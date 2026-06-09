---
name: apply-iterator-pattern
description: Use when you need to access elements of a collection sequentially without exposing its underlying representation — decoupling traversal logic from the collection's data structure.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 257–271; Python iteration protocol (__iter__/__next__); Java Iterable/Iterator; C++ STL iterators; Ruby Enumerable
tags: [design-patterns, behavioral, iterator, oop, developer, collection-traversal, encapsulation]
related: [apply-composite-pattern, apply-visitor-pattern, apply-solid-principles]
---

# Apply Iterator Pattern

Provide a way to sequentially access elements of a collection without exposing its underlying structure.

## Why This Is Best Practice

**Adopted by:** Python's iteration protocol (`__iter__`/`__next__` — every `for` loop
uses it, making it the most-invoked pattern in the language), Java's `Iterable`/
`Iterator` (every Java collection implements it — `for-each` loops, streams, and
collectors all depend on it), C++ STL iterators (the foundation of the entire Standard
Template Library algorithm set), and Ruby's `Enumerable` module.
**Impact:** Python's iterator protocol is cited in the language reference as the reason
generators, comprehensions, and the `for` statement all share one protocol. The unifying
effect: any object implementing `__iter__` works with `for`, `zip`, `map`, `list()`,
and every standard library function — without those functions knowing the collection's
type.
**Why best:** The alternative — index-based traversal — exposes the collection type
(`list[i]` doesn't work on a linked list or a file stream). Iterator abstracts over
any sequential structure, enabling algorithms that work on lists, trees, streams,
database cursors, and network responses uniformly.

Sources: Gamma et al. (1994) pp. 257–271; Python iterator protocol documentation;
Java `Iterable` specification

## Steps

### Step 1: Implement Python's iterator protocol on your collection

```python
class NumberRange:
    def __init__(self, start: int, end: int, step: int = 1):
        self._start = start
        self._end = end
        self._step = step

    def __iter__(self):
        current = self._start
        while current < self._end:
            yield current
            current += self._step
```

Using `yield` creates a generator iterator — the simplest correct implementation.

### Step 2: For stateful iterators with external control, use `__next__` explicitly

```python
class NumberRangeIterator:
    def __init__(self, start: int, end: int, step: int):
        self._current = start
        self._end = end
        self._step = step

    def __iter__(self):
        return self

    def __next__(self):
        if self._current >= self._end:
            raise StopIteration
        value = self._current
        self._current += self._step
        return value
```

Use explicit `__next__` when the iterator must be paused and resumed externally
(e.g., paginated API responses loaded one page at a time).

### Step 3: Keep the iterator separate from the collection for multiple simultaneous traversals

```python
class BookShelf:
    def __init__(self):
        self._books: list[str] = []

    def add(self, book: str):
        self._books.append(book)

    def __iter__(self):
        return iter(self._books)   # each call creates a fresh iterator
```

`iter(self._books)` creates a new list iterator each time, so two simultaneous
`for` loops over the same `BookShelf` don't interfere.

### Step 4: Use standard library iteration — don't reinvent `next()`, `zip()`, `enumerate()`

```python
shelf = BookShelf()
shelf.add("Refactoring")
shelf.add("Clean Code")

for i, book in enumerate(shelf):       # enumerate works — shelf is iterable
    print(f"{i+1}. {book}")

paired = list(zip(shelf, shelf))       # multiple iterators — independent
```

### Step 5: For lazy or infinite sequences, use generators

```python
def fibonacci():
    a, b = 0, 1
    while True:
        yield a
        a, b = b, a + b

for n in fibonacci():                  # infinite iterator
    if n > 100:
        break
    print(n)
```

## When NOT to Use

- **Random access collections where index-based access is the primary use** — if callers need `collection[i]` more than sequential traversal, a list interface is clearer.
- **When only one traversal order exists and the collection is simple** — Python lists already have full iterator support; wrapping them adds nothing.

## Common Mistakes

**Modifying the collection during iteration.** Deleting or inserting elements while iterating over them produces undefined behavior in most languages. Collect items to remove, then remove after the loop.

**Making the collection and iterator the same object.** If `BookShelf.__iter__` returns `self` and `BookShelf.__next__` manages position, two `for` loops share state and conflict. Return a new iterator object.

**Forgetting `StopIteration` in `__next__`.** An iterator that never raises `StopIteration` produces an infinite loop in a `for` loop. Always raise it when exhausted.
