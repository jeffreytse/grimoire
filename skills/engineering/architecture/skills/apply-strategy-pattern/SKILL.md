---
name: apply-strategy-pattern
description: Use when you have multiple algorithms or behaviors that can be swapped for a given task — encapsulating each algorithm in its own class and making them interchangeable at runtime.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 315–323; Java Comparator; Python sorted(key=); scikit-learn estimator API; AWS SDK retry strategies
tags: [design-patterns, behavioral, strategy, oop, developer, algorithm-selection, open-closed]
related: [apply-state-pattern, apply-template-method-pattern, apply-solid-principles]
---

# Apply Strategy Pattern

Define a family of algorithms, encapsulate each one, and make them interchangeable — letting the algorithm vary independently from the clients that use it.

## Why This Is Best Practice

**Adopted by:** Java's `Comparator` interface (Strategy for sorting — used in every
`Collections.sort()` call), Python's `sorted(key=...)` (a Strategy as a first-class
function), scikit-learn's estimator API (every ML algorithm is a Strategy implementing
`fit()`/`predict()` — the library's entire design), and AWS SDK's configurable
retry strategies (`RetryPolicy` interface).
**Impact:** GoF documents that Strategy eliminates conditional logic that selects
algorithms at runtime. scikit-learn's adoption of Strategy is directly responsible for
its plug-and-play model comparison: swapping `RandomForestClassifier` for `SVC` requires
one line — because both implement the same Strategy interface.
**Why best:** `if/elif` dispatch — the alternative — couples the client to every
algorithm variant and requires modification for every new algorithm. Strategy makes each
algorithm independently changeable, testable, and composable without touching the client.

Sources: Gamma et al. (1994) pp. 315–323; Java `Comparator` specification;
scikit-learn API documentation

## Steps

### Step 1: Define the strategy interface

```python
from abc import ABC, abstractmethod

class SortStrategy(ABC):
    @abstractmethod
    def sort(self, data: list) -> list: ...
```

### Step 2: Implement concrete strategies — one class per algorithm

```python
class BubbleSortStrategy(SortStrategy):
    def sort(self, data: list) -> list:
        data = list(data)
        n = len(data)
        for i in range(n):
            for j in range(0, n - i - 1):
                if data[j] > data[j + 1]:
                    data[j], data[j + 1] = data[j + 1], data[j]
        return data

class QuickSortStrategy(SortStrategy):
    def sort(self, data: list) -> list:
        if len(data) <= 1:
            return data
        pivot = data[len(data) // 2]
        left = [x for x in data if x < pivot]
        mid = [x for x in data if x == pivot]
        right = [x for x in data if x > pivot]
        return self.sort(left) + mid + self.sort(right)
```

### Step 3: Context holds a strategy reference — delegates to it

```python
class Sorter:
    def __init__(self, strategy: SortStrategy):
        self._strategy = strategy

    def set_strategy(self, strategy: SortStrategy) -> None:
        self._strategy = strategy

    def sort(self, data: list) -> list:
        return self._strategy.sort(data)
```

### Step 4: Swap strategies at runtime without changing context

```python
sorter = Sorter(BubbleSortStrategy())
print(sorter.sort([3, 1, 4, 1, 5]))   # [1, 1, 3, 4, 5]

sorter.set_strategy(QuickSortStrategy())
print(sorter.sort([3, 1, 4, 1, 5]))   # [1, 1, 3, 4, 5] — same result, different algorithm
```

### Step 5: In Python, use first-class functions as lightweight strategies

```python
from typing import Callable

class Sorter:
    def __init__(self, strategy: Callable[[list], list] = sorted):
        self._strategy = strategy

    def sort(self, data: list) -> list:
        return self._strategy(data)

# Use built-ins or lambdas as strategies — no class needed
sorter = Sorter(strategy=lambda d: sorted(d, reverse=True))
```

Use class-based strategies when state is needed; function strategies for stateless algorithms.

## When NOT to Use

- **Only one algorithm ever** — if the algorithm never varies, a direct method call is simpler. Strategy pays off only when multiple alternatives exist.
- **Trivial variation** — if the difference between strategies is a single parameter (ascending vs descending sort), a `reverse=True` flag is cleaner than two strategy classes.

## Common Mistakes

**Passing strategy selection logic into the strategy.** A Strategy that asks "which variant should I use?" is no longer a Strategy — it's a dispatcher. Each Strategy implements exactly one algorithm.

**State leakage between calls.** Strategies that store results from previous `sort()` calls produce bugs when reused. Strategies should be stateless or clearly document what state they maintain.

**Using Strategy when Template Method is the right pattern.** If all algorithms share the same skeleton but differ only in specific steps, `apply-template-method-pattern` is more appropriate than several full Strategy implementations.
