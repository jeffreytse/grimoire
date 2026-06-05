---
name: design-event-driven-architecture
description: Use when designing systems that need loose coupling between services, asynchronous processing of high-volume events, or real-time data propagation across multiple consumers
source: Hohpe & Woolf "Enterprise Integration Patterns" (2003); Richardson "Microservices Patterns" (2018); Stopford "Designing Event-Driven Systems" (O'Reilly 2018)
tags: [architecture, events, microservices, messaging]
verified: true
---

# Design Event-Driven Architecture

Architect systems where components communicate by producing and consuming events, enabling loose coupling, independent scaling, and resilience to partial failures.

## Why This Is Best Practice

**Adopted by:** LinkedIn (Kafka origin), Uber (event sourcing for trips), Netflix (events for recommendations), Amazon (SNS/SQS for order processing)
**Impact:** LinkedIn reduced inter-service coupling from 100% synchronous to < 20% after Kafka adoption; event-driven architectures handle 10-100x more throughput than synchronous RPC for the same infrastructure (Stopford 2018); Martin Fowler: "the most important characteristic of microservices is that they are organized around business capabilities" — events encode business capabilities
**Why best:** Synchronous calls create cascading failures; events decouple producers from consumers, enable fan-out to multiple consumers, and provide a durable audit log of system state changes

Sources: Hohpe & Woolf "Enterprise Integration Patterns" Addison-Wesley (2003); Richardson "Microservices Patterns" Manning (2018); Stopford "Designing Event-Driven Systems" O'Reilly (2018)

## Steps

1. **Model the domain as events** — Identify domain events: "OrderPlaced", "PaymentProcessed", "InventoryReserved". Events are immutable facts in the past tense. They represent what happened, not commands to do something. Name events from the business domain, not technical operations.

2. **Choose an event streaming platform** — Apache Kafka for high-throughput, durable, replayable event streams (log-based). AWS SQS/SNS for simpler queue-based messaging with managed infrastructure. RabbitMQ for complex routing patterns and traditional message queuing. Kafka is the standard for event-driven systems requiring replay and event sourcing.

3. **Design event schemas** — Define event schema: event type, event ID (UUID), timestamp, source service, version, and payload. Use schema registry (Confluent Schema Registry, AWS Glue Schema Registry) to enforce and evolve schemas. Prefer Avro or Protobuf over JSON for production volume; JSON for developer ergonomics in low-volume cases.

4. **Apply the outbox pattern for reliability** — Never publish an event and update a database in separate transactions; one will fail, leaving state inconsistent. Use the transactional outbox: write the event to an outbox table in the same DB transaction as the state change, then a separate process publishes from the outbox to the event stream. This guarantees at-least-once delivery.

5. **Design for idempotent consumers** — Events are delivered at least once; consumers may process the same event multiple times. Every consumer must be idempotent: processing the same event twice must produce the same result as processing it once. Use event ID deduplication: store processed event IDs and skip duplicates.

6. **Define consumer groups and partitioning** — Kafka: partition the event stream by a natural key (order ID, user ID) to ensure related events are ordered. Assign consumer groups so each service receives all events independently. Partitioning determines parallelism and ordering guarantees.

7. **Handle failures with dead letter queues** — Events that cannot be processed (schema violation, downstream failure) must not be silently dropped or block the consumer. Route failed events to a dead letter topic after N retry attempts. Monitor DLQ depth as a service health metric. Implement automated reprocessing after root cause is fixed.

8. **Implement event versioning** — Events are immutable once published but schemas evolve. Strategy: add new fields as optional (backward compatible). Never remove or rename existing fields. Use schema version in the event header. Consumers must handle unknown fields gracefully (ignore, don't fail).

9. **Provide event replay capability** — Kafka retains events for a configurable period (default 7 days; set longer for event sourcing). New consumers can replay from the beginning of the log to rebuild state. Design consumers to handle replay efficiently. This enables: onboarding new services, disaster recovery, and debugging.

10. **Monitor event pipeline health** — Track: consumer lag (events published but not yet processed), DLQ depth, event processing latency (p99), and schema validation error rate. Alert on consumer lag growth (consumer falling behind producer) as a precursor to processing failure. Use Kafka's consumer group offset tracking for lag measurement.

## Rules

- Events are facts, not commands; "OrderPlaced" not "PlaceOrder"; event-driven is not RPC with messaging.
- Producers must not know about consumers; adding a new consumer must require zero changes to the producer.
- Never rely on event ordering across different event types; ordering is only guaranteed within a partition.
- The outbox pattern is mandatory for any state change that must be reflected in an event; dual-write without it is an eventual inconsistency waiting to happen.

## Common Mistakes

- **Synchronous event publishing** — calling an API to publish an event synchronously re-introduces the coupling that events are meant to eliminate; publish asynchronously via the outbox.
- **Fat events with too much data** — events that carry entire entity state cause coupling through data; carry only what changed and event consumers can fetch additional data if needed (or use event sourcing).
- **No schema registry** — consumers break silently when producers change event schemas without coordination; schema registry enforces compatibility.
- **Ignoring consumer lag** — a consumer falling behind by millions of events is an operational emergency, not just a performance issue; alert on lag.

## When NOT to Use

- Simple CRUD applications where synchronous request-response is simpler and sufficient
- Operations requiring immediate consistency guarantees (financial transactions needing synchronous confirmation)
- Teams without Kafka/messaging operational expertise — operational complexity of event streaming is significant; start with synchronous services and evolve
