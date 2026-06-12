---
name: design-event-driven-system
description: Use when designing systems that communicate through events, including event sourcing and CQRS architectures
source: Martin Fowler "Event-Driven Architecture" (martinfowler.com); Greg Young "CQRS and Event Sourcing" talks; Vernon "Implementing Domain-Driven Design" (Addison-Wesley 2013)
tags: [event-driven, event-sourcing, cqrs, architecture, messaging, kafka]
verified: true
---

# Design Event-Driven System

Design systems where components communicate via immutable events, enabling loose coupling, auditability, and temporal decoupling.

## Why This Is Best Practice

**Adopted by:** LinkedIn (Kafka origin), Uber, Airbnb, AWS (EventBridge), Martin Fowler's canonical patterns
**Impact:** LinkedIn's move to event-driven with Kafka handled 7 trillion messages/day; event sourcing provides a complete audit log by construction, eliminating a common compliance requirement.

**Why best:** Events as first-class citizens decouple producers from consumers in time and space. Event sourcing makes the audit log the source of truth. CQRS separates read and write models, allowing each to scale and evolve independently.

## Steps

1. **Identify domain events** — Use event storming: list all significant things that happen in the business domain (past tense: `OrderPlaced`, `PaymentFailed`).
2. **Classify event types** — Distinguish: domain events (business facts), integration events (cross-service), and commands (intent to act — not events).
3. **Design the event schema** — Include: event type, version, aggregate ID, timestamp, correlation ID, and payload. Use schema registry (Avro/Protobuf + Confluent Schema Registry) to enforce compatibility.
4. **Choose event transport** — Kafka for high-throughput durable streams; RabbitMQ for transient messaging; AWS EventBridge for cloud-native routing; in-process event bus for within-service.
5. **Apply CQRS if read/write loads differ** — Separate command handlers (write side) from query models (read side); project events into optimized read stores.
6. **Handle eventual consistency** — Design UIs and downstream consumers to tolerate lag; use correlation IDs to track saga progress.
7. **Plan event versioning** — Use upcasters or parallel event types for schema evolution; never mutate published event schemas.

## Rules

- Events are immutable facts — never delete or mutate a published event.
- Include correlation and causation IDs in every event for distributed tracing.
- Consumers must be idempotent — the same event may be delivered more than once (at-least-once delivery).
- Avoid event sourcing for simple CRUD domains — the complexity cost exceeds the benefit.

## Examples

Order saga (choreography):
1. `OrderService` publishes `OrderPlaced`.
2. `InventoryService` listens → reserves stock → publishes `StockReserved`.
3. `PaymentService` listens → charges → publishes `PaymentConfirmed`.
4. `OrderService` listens → marks order `Confirmed`.

No orchestrator; each service reacts to events independently.

## Common Mistakes

- **Event as RPC** — publishing `GetUserRequest` as an event defeats the purpose; use synchronous calls for queries.
- **Fat events with full entity state** — creates implicit coupling; prefer thin events with IDs and let consumers fetch if needed, or use event carried state transfer intentionally.
- **No dead-letter queue** — failed events are silently dropped; always configure a DLQ.
