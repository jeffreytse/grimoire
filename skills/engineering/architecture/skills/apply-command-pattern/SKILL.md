---
name: apply-command-pattern
description: Use when you need to parameterize objects with operations, queue or log requests, support undoable operations, or implement transactional behavior — by encapsulating each request as an object with a uniform execute interface.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 233–242; Java javax.swing.Action; Python asyncio task queue; Git (every commit is a command); Redux actions (JavaScript state management)
tags: [design-patterns, behavioral, command, oop, developer, undo-redo, queuing, transactional]
related: [apply-chain-of-responsibility-pattern, apply-memento-pattern, apply-solid-principles]
---

# Apply Command Pattern

Encapsulate a request as an object, allowing requests to be parameterized, queued, logged, and undone.

## Why This Is Best Practice

**Adopted by:** Java's `javax.swing.Action` (every button and menu item in Java Swing
wraps a Command — in every Java GUI application), Redux (JavaScript's most-used state
management library — every state change is a dispatched Action, i.e., a Command),
Git's commit model (every commit is an immutable Command object with undo via `revert`),
and database transaction logs (write-ahead logging records Commands for recovery).
**Impact:** GoF documents that Command enables undo/redo, logging, and transactional
rollback without changing invoker or receiver code. Redux's adoption of Command
(85M weekly npm downloads) demonstrates the pattern's value at scale: time-travel
debugging (replay/undo arbitrary state changes) is only possible because every change
is a serializable Command object.
**Why best:** Direct method calls — the alternative — cannot be queued, logged, or
undone without invasive change to every caller. Command objects turn operations into
first-class values that can be stored, transmitted, and replayed.

Sources: Gamma et al. (1994) pp. 233–242; Redux documentation; Git internals documentation

## Steps

### Step 1: Define the command interface

```python
from abc import ABC, abstractmethod

class Command(ABC):
    @abstractmethod
    def execute(self) -> None: ...

    @abstractmethod
    def undo(self) -> None: ...
```

### Step 2: Implement concrete commands that wrap a receiver

```python
class TextEditor:
    def __init__(self):
        self.text = ""

class InsertTextCommand(Command):
    def __init__(self, editor: TextEditor, position: int, text: str):
        self._editor = editor
        self._position = position
        self._text = text

    def execute(self) -> None:
        self._editor.text = (
            self._editor.text[:self._position]
            + self._text
            + self._editor.text[self._position:]
        )

    def undo(self) -> None:
        self._editor.text = (
            self._editor.text[:self._position]
            + self._editor.text[self._position + len(self._text):]
        )
```

### Step 3: Build an invoker that manages command history

```python
class CommandHistory:
    def __init__(self):
        self._history: list[Command] = []

    def execute(self, command: Command) -> None:
        command.execute()
        self._history.append(command)

    def undo(self) -> None:
        if self._history:
            self._history.pop().undo()
```

### Step 4: Use the invoker — sender doesn't know the receiver

```python
editor = TextEditor()
history = CommandHistory()

history.execute(InsertTextCommand(editor, 0, "Hello"))
history.execute(InsertTextCommand(editor, 5, " World"))
print(editor.text)   # "Hello World"

history.undo()
print(editor.text)   # "Hello"
```

### Step 5: For queuing or batch execution, collect commands before executing

```python
batch: list[Command] = [
    InsertTextCommand(editor, 0, "Line 1\n"),
    InsertTextCommand(editor, 7, "Line 2\n"),
]

for cmd in batch:
    history.execute(cmd)
```

Macro recording, test scripts, and transaction replay all use this pattern.

## When NOT to Use

- **Simple, non-undoable operations** — if undo/redo, queuing, and logging are not required, encapsulating the call in a Command object adds overhead without benefit.
- **When operations have no receiver** — if the command is just a function call with no state to capture, a first-class function (lambda, partial) is simpler.

## Common Mistakes

**Command that holds too much state.** If a command stores a deep snapshot of receiver state to support undo, consider `apply-memento-pattern` for the snapshot instead.

**Forgetting that undo must perfectly reverse execute.** `undo()` must leave the system in exactly the state it was in before `execute()`. Test `execute(); undo()` as a no-op pair.

**Mutable command state.** A command stored in history must not change after execution. If parameters can change externally, copy them in the constructor.
