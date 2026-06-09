---
name: apply-interpreter-pattern
description: Use when you need to interpret sentences in a simple language — defining a grammar as a class hierarchy where each rule is a class and interpreting an expression is traversing that hierarchy.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 243–255; Python ast module; SQL query parsers; regular expression engines; Excel formula evaluation
tags: [design-patterns, behavioral, interpreter, oop, developer, language-processing, expression-evaluation]
related: [apply-composite-pattern, apply-visitor-pattern, apply-solid-principles]
---

# Apply Interpreter Pattern

Define a grammar for a language and represent each grammar rule as a class, using composition to interpret expressions.

## Why This Is Best Practice

**Adopted by:** Python's `ast` module (every Python expression is an AST node that
interprets itself — the foundation of all Python compilers, linters, and tools), Java's
`javax.script` (script interpreter architecture), SQL database engines (query expressions
are parsed into an AST where each node knows how to evaluate itself), and Excel (cell
formula evaluation traverses an expression tree).
**Impact:** GoF documents that Interpreter makes grammar rules independently changeable
and combinable. Python's `ast` module design enables the entire ecosystem of static
analysis tools (Flake8, mypy, Bandit) — each tool visits the same expression tree
without modifying the parser.
**Why best:** The alternative — a large `if/elif` parser — puts all rules in one method.
Adding a new grammar rule requires modifying that method. Interpreter makes each rule
an independent class — adding a rule means adding a class, not editing existing code.

Sources: Gamma et al. (1994) pp. 243–255; Python `ast` module documentation;
ANTLR4 documentation (widely-used parser generator that embodies the Interpreter pattern)

## Steps

### Step 1: Define the grammar (small and simple — this pattern suits narrow DSLs)

```
<expression> ::= <number>
               | <expression> '+' <expression>
               | <expression> '*' <expression>
```

### Step 2: Define the abstract expression interface

```python
from abc import ABC, abstractmethod

class Expression(ABC):
    @abstractmethod
    def interpret(self, context: dict) -> int: ...
```

### Step 3: Implement terminal expressions (leaves — no sub-expressions)

```python
class NumberExpression(Expression):
    def __init__(self, value: int):
        self._value = value

    def interpret(self, context: dict) -> int:
        return self._value

class VariableExpression(Expression):
    def __init__(self, name: str):
        self._name = name

    def interpret(self, context: dict) -> int:
        return context[self._name]
```

### Step 4: Implement non-terminal expressions (composites — have sub-expressions)

```python
class AddExpression(Expression):
    def __init__(self, left: Expression, right: Expression):
        self._left = left
        self._right = right

    def interpret(self, context: dict) -> int:
        return self._left.interpret(context) + self._right.interpret(context)

class MultiplyExpression(Expression):
    def __init__(self, left: Expression, right: Expression):
        self._left = left
        self._right = right

    def interpret(self, context: dict) -> int:
        return self._left.interpret(context) * self._right.interpret(context)
```

### Step 5: Build and interpret an expression tree

```python
# Represents: (x + 3) * 2
context = {"x": 5}

expr = MultiplyExpression(
    AddExpression(VariableExpression("x"), NumberExpression(3)),
    NumberExpression(2)
)

print(expr.interpret(context))   # (5 + 3) * 2 = 16
```

Note: expression trees are typically built by a parser, not manually. Interpreter
defines the evaluation step; parsing is a separate concern.

## When NOT to Use

- **Complex grammars** — GoF explicitly warns against using Interpreter for complex languages. Use a real parser generator (ANTLR, PLY, lark-parser) instead.
- **Performance-critical expression evaluation** — recursive tree traversal is slow for large expressions. Compiled languages parse to bytecode for a reason.
- **When an existing library covers the domain** — regular expressions (`re`), SQL (`sqlparse`), and math expressions (`sympy`) have mature parsers. Don't reimplement them.

## Common Mistakes

**Using Interpreter for a complex language.** If the grammar has more than 10 rules or requires backtracking, use a parser generator. Interpreter doesn't scale to real languages.

**Mixing parsing and interpretation.** The Interpreter pattern covers evaluation of an already-parsed expression tree. Building the tree (parsing) is a separate step and should be separate code.

**Not handling context errors.** If `context["x"]` is missing, `VariableExpression.interpret()` raises `KeyError`. Validate context before interpretation or catch and wrap errors.
