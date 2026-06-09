---
name: apply-mediator-pattern
description: Use when many objects communicate with each other in complex, tangled ways — centralizing all interaction in a mediator object so that components only talk to the mediator, not to each other.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 273–282; chat room servers; aviation air traffic control; React event handling architecture; MediatR library (.NET)
tags: [design-patterns, behavioral, mediator, oop, developer, decoupling, event-coordination]
related: [apply-observer-pattern, apply-facade-pattern, apply-low-coupling]
---

# Apply Mediator Pattern

Define an object that encapsulates how a set of objects interact, promoting loose coupling by preventing direct object-to-object references.

## Why This Is Best Practice

**Adopted by:** MediatR (.NET) — used in the majority of .NET CQRS implementations,
downloaded 100M+ times on NuGet — is a direct implementation of the Mediator pattern.
React's unidirectional data flow (state lifted to parent components that coordinate
children) applies Mediator at the UI level. Message brokers (Kafka, RabbitMQ) are
infrastructure-level mediators — adopted by virtually every large-scale distributed system.
**Impact:** GoF documents that M objects with direct references to each other require
M×(M-1) connection points. With a mediator, each of M objects has one reference (to
the mediator): M connections total. Changing one component's interaction requires
modifying one method in the mediator — not M-1 classes.
**Why best:** Direct object-to-object communication — the alternative — creates a web
of dependencies where a change to one component's protocol cascades to all its direct
communicators. Mediator centralizes that protocol and isolates changes.

Sources: Gamma et al. (1994) pp. 273–282; MediatR documentation (.NET);
React "lifting state up" documentation

## Steps

### Step 1: Define the mediator interface

```python
from abc import ABC, abstractmethod

class ChatMediator(ABC):
    @abstractmethod
    def send_message(self, message: str, sender: "ChatUser") -> None: ...

    @abstractmethod
    def add_user(self, user: "ChatUser") -> None: ...
```

### Step 2: Implement the concrete mediator — coordinates component interactions

```python
class ChatRoom(ChatMediator):
    def __init__(self):
        self._users: list["ChatUser"] = []

    def add_user(self, user: "ChatUser") -> None:
        self._users.append(user)

    def send_message(self, message: str, sender: "ChatUser") -> None:
        for user in self._users:
            if user is not sender:
                user.receive(message, sender.name)
```

### Step 3: Define the colleague (component) — knows only the mediator

```python
class ChatUser:
    def __init__(self, name: str, mediator: ChatMediator):
        self.name = name
        self._mediator = mediator

    def send(self, message: str) -> None:
        print(f"{self.name} sends: {message}")
        self._mediator.send_message(message, self)   # talks to mediator, not other users

    def receive(self, message: str, from_name: str) -> None:
        print(f"{self.name} receives from {from_name}: {message}")
```

### Step 4: Wire up via the mediator — components never reference each other

```python
room = ChatRoom()

alice = ChatUser("Alice", room)
bob = ChatUser("Bob", room)
carol = ChatUser("Carol", room)

room.add_user(alice)
room.add_user(bob)
room.add_user(carol)

alice.send("Hello everyone!")   # Bob and Carol receive — Alice never references them
```

### Step 5: Add coordination logic to the mediator — not to components

If sending a message should log, filter profanity, or throttle, that logic belongs in
`ChatRoom.send_message()` — not in each `ChatUser.send()`. Components stay simple;
the mediator handles coordination.

## When NOT to Use

- **When only 2–3 objects interact** — for a small, stable set of direct collaborators, a mediator adds indirection without payoff. Use direct references.
- **When the mediator becomes a God Object** — if the mediator contains complex business logic for many unrelated interactions, it's an anti-pattern. Split into multiple focused mediators.

## Common Mistakes

**Mediator that leaks component knowledge.** If the mediator imports concrete component types and calls their specific methods, it's tightly coupled to every component it coordinates. Use the component interface, not concrete types.

**Components calling each other through the mediator inefficiently.** If A needs a response from B synchronously, a mediator round-trip adds latency. Direct synchronous calls may be appropriate for tightly-coupled pairs.

**Confusing Mediator with Facade.** Facade simplifies a subsystem for external clients. Mediator coordinates internal communication between components in the same subsystem. Components don't know about a Facade; they do know about the Mediator.
