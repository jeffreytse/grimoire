---
name: apply-memento-pattern
description: Use when you need to capture and restore an object's internal state — for undo/redo, rollback, or snapshot functionality — without exposing or violating the object's encapsulation.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 283–291; text editor undo history; database savepoints; game save states; Redux state snapshots
tags: [design-patterns, behavioral, memento, oop, developer, undo-redo, snapshot, state-management]
related: [apply-command-pattern, apply-state-pattern, apply-solid-principles]
---

# Apply Memento Pattern

Capture an object's internal state in an external object (memento) so the state can be restored later, without exposing internal details.

## Why This Is Best Practice

**Adopted by:** Every text editor's undo system (VS Code, IntelliJ, Vim) — each undo
operation restores a memento of prior editor state. Database savepoints (`SAVEPOINT` in
SQL) are mementos stored by the transaction engine. Git commits are mementos of
filesystem state — each commit captures a snapshot restorable at any time. Redux
time-travel debugging (DevTools' "jump to state") replays memento snapshots.
**Impact:** GoF documents that without Memento, undo requires either breaking
encapsulation (the undo manager reads private fields) or duplicating all object state
externally. Memento lets the object produce its own opaque snapshot, maintaining
encapsulation while enabling undo.
**Why best:** The alternative — Command pattern with undo methods — requires every
command to implement a perfect inverse operation. For complex state, that is
error-prone and sometimes impossible (file content after a regex replace can't always
be inverted). Memento captures state directly, making restoration exact.

Sources: Gamma et al. (1994) pp. 283–291; Git internals documentation; Redux DevTools

## Steps

### Step 1: Define the memento — an opaque snapshot of originator state

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class EditorMemento:
    text: str
    cursor_position: int
    # frozen=True: memento is immutable — protects the snapshot
```

Make the memento's internal structure opaque to the caretaker (undo history).
`frozen=True` (or a private constructor) enforces this.

### Step 2: Add `save()` and `restore()` to the originator

```python
class TextEditor:
    def __init__(self):
        self._text = ""
        self._cursor = 0

    def type(self, text: str) -> None:
        self._text = self._text[:self._cursor] + text + self._text[self._cursor:]
        self._cursor += len(text)

    def save(self) -> EditorMemento:
        return EditorMemento(self._text, self._cursor)

    def restore(self, memento: EditorMemento) -> None:
        self._text = memento.text
        self._cursor = memento.cursor_position

    @property
    def content(self) -> str:
        return self._text
```

The originator creates and restores its own mementos — it's the only one that knows
what state needs capturing.

### Step 3: Build the caretaker (undo history) — stores mementos, never reads them

```python
class UndoHistory:
    def __init__(self, editor: TextEditor):
        self._editor = editor
        self._history: list[EditorMemento] = []

    def save(self) -> None:
        self._history.append(self._editor.save())

    def undo(self) -> None:
        if self._history:
            self._editor.restore(self._history.pop())
```

The caretaker holds mementos but treats them as opaque objects — it never inspects
or modifies their contents.

### Step 4: Use at every state-changing operation

```python
editor = TextEditor()
history = UndoHistory(editor)

history.save()
editor.type("Hello")

history.save()
editor.type(" World")

print(editor.content)   # "Hello World"

history.undo()
print(editor.content)   # "Hello"

history.undo()
print(editor.content)   # ""
```

### Step 5: For large state, use incremental or compressed mementos

If storing full snapshots is expensive (large documents, game world state), store deltas:

```python
@dataclass(frozen=True)
class IncrementalMemento:
    operation: str   # "insert" or "delete"
    position: int
    text: str
    # reconstruct full state by replaying forward or backward
```

## When NOT to Use

- **When state is small and Commands support perfect undo** — if every operation has a cheap inverse, `apply-command-pattern` with `undo()` is simpler and uses less memory.
- **When snapshots are prohibitively expensive** — if a full-state snapshot requires deep-copying gigabytes, use incremental/delta snapshots or event sourcing.

## Common Mistakes

**Caretaker that reads memento internals.** The caretaker should treat mementos as black boxes. If it reads `memento.text` to display a preview, it couples to the internal structure — which must then stay stable.

**Mutable mementos.** If a memento can be modified after capture, the "snapshot" is no longer reliable. Always use immutable (frozen) mementos.

**Not saving before every mutation.** If the caretaker misses a `save()` call, an intermediate state is lost from the undo chain. Consider making `save()` automatic in the originator's mutating methods.
