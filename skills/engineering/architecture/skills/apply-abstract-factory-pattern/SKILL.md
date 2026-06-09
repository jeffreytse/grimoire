---
name: apply-abstract-factory-pattern
description: Use when a system must create families of related objects and must remain independent of how those objects are created — so that swapping one product family for another requires no changes to client code.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 87–95; adopted in Java AWT/Swing (LookAndFeel), .NET WPF (theme factories), Apple AppKit (UI factory hierarchy)
tags: [design-patterns, creational, abstract-factory, oop, developer, extensibility, product-family]
related: [apply-factory-method-pattern, apply-solid-principles]
---

# Apply Abstract Factory Pattern

Provide an interface for creating families of related objects without specifying their concrete classes.

## Why This Is Best Practice

**Adopted by:** Java AWT/Swing (the `LookAndFeel` system is an abstract factory for UI
components), .NET WPF theme factories, Apple AppKit (NSWindow/NSView factory hierarchy),
and Google Guice/Spring (module-level provider factories) — among the most downloaded
frameworks in each ecosystem.
**Impact:** GoF documents that Abstract Factory eliminates conditional logic that otherwise
appears at every object-creation site when switching product families. In practice, adding
a new UI theme in Java Swing requires implementing one factory class — not editing every
component constructor. The pattern makes "swap the whole family" a one-class change.
**Why best:** The alternative — instantiating concrete classes directly — scatters
product-family selection across the codebase. Every `new WindowsButton()` is a site that
must be updated when switching to `MacButton`. Abstract Factory centralizes that decision
in one place and enforces consistency: you cannot accidentally mix products from different
families (a Windows button with a Mac checkbox).

Sources: Gamma et al. (1994) pp. 87–95; Freeman & Freeman, "Head First Design Patterns"
(2004) ch. 4; Java AWT documentation

## Steps

### Step 1: Define the abstract factory interface

Declare a creation method for each product type in the family:

```python
from abc import ABC, abstractmethod

class UIFactory(ABC):
    @abstractmethod
    def create_button(self) -> Button: ...

    @abstractmethod
    def create_checkbox(self) -> Checkbox: ...
```

### Step 2: Define abstract product interfaces

```python
class Button(ABC):
    @abstractmethod
    def render(self) -> str: ...

class Checkbox(ABC):
    @abstractmethod
    def render(self) -> str: ...
```

### Step 3: Implement concrete factories and products per family

```python
class WindowsFactory(UIFactory):
    def create_button(self) -> Button:
        return WindowsButton()
    def create_checkbox(self) -> Checkbox:
        return WindowsCheckbox()

class MacFactory(UIFactory):
    def create_button(self) -> Button:
        return MacButton()
    def create_checkbox(self) -> Checkbox:
        return MacCheckbox()
```

### Step 4: Write client code against the abstract factory only

```python
class Application:
    def __init__(self, factory: UIFactory):
        self.button = factory.create_button()
        self.checkbox = factory.create_checkbox()

    def render(self):
        print(self.button.render())
        print(self.checkbox.render())

# Swap families without touching Application
app = Application(WindowsFactory())
app = Application(MacFactory())
```

### Step 5: Select the concrete factory at the entry point (not inside components)

Decide which family to use once — at startup, from config, or from environment — then
pass the factory down. Components never call `WindowsFactory()` directly.

## When NOT to Use

- **Single-product systems** — if only one type of object is being created, `apply-factory-method-pattern` is simpler.
- **Rarely-changing product families** — if there's only ever one family and no plans to add another, the interface adds indirection with no payoff.
- **When products in a family don't need to be consistent** — if mixing products across families is acceptable, the constraint enforced by Abstract Factory is unnecessary.

## Common Mistakes

**Adding new product types to an existing factory.** Adding `create_tooltip()` to `UIFactory` requires changing every concrete factory. This violates OCP. Prefer extending via composition rather than widening the factory interface.

**Using Abstract Factory for a single product.** If you only have one kind of product, Factory Method is the right pattern — Abstract Factory is for families.

**Putting factory selection logic inside a component.** `if os == "windows": factory = WindowsFactory()` inside a `render()` method defeats the point. Selection belongs at the composition root.
