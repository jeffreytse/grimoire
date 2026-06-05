---
name: design-contract-testing
description: Use when multiple teams maintain independent services that communicate via APIs, or when integration test environments are slow, expensive, or unavailable
source: Pact contract testing framework (DiUS Computing); Richardson "Microservices Patterns" (2018) Ch. 9; Martin Fowler "IntegrationContractTest" article
tags: [testing, microservices, api, integration]
verified: true
---

# Design Contract Testing

Define and verify shared API contracts between consumer and provider services independently, so teams can evolve and deploy services without end-to-end integration test environments.

## Why This Is Best Practice

**Adopted by:** Atlassian (built Pact), Google (uses protobuf contracts), Amazon (contract tests for internal APIs), Seek (pioneered consumer-driven contracts)
**Impact:** Teams using contract testing deploy services independently 4x more often (Atlassian case study); eliminates integration environment dependency; finds breaking changes in minutes instead of hours in e2e test suites
**Why best:** End-to-end integration tests require all services running simultaneously, are slow, and become a deployment bottleneck; contract tests verify interface compatibility in isolation in seconds

Sources: Fowler "IntegrationContractTest" martinfowler.com (2011); Richardson "Microservices Patterns" Manning (2018); Pact Foundation documentation

## Steps

1. **Choose consumer-driven contracts** — The consumer (API caller) defines the contract: the specific request it sends and the minimum response fields it needs. The provider verifies it can satisfy that contract. Consumer-driven contracts prevent providers from breaking consumers with changes consumers don't care about, while allowing providers to extend their API freely.

2. **Select a contract testing tool** — Pact is the de facto standard (supports REST, GraphQL, gRPC; multiple language SDKs). Spring Cloud Contract for JVM/Spring ecosystems. AsyncAPI for event-driven contracts. Choose based on your stack; Pact has the broadest language support.

3. **Write consumer-side contract tests** — In the consumer service: write a test that defines the interaction (HTTP request + expected response fields). Use Pact's mock server as the provider. The test generates a pact file (JSON) that describes the contract. These tests run independently without the real provider.

4. **Publish contracts to a broker** — Use PactFlow or a self-hosted Pact Broker to store and version pact files. The broker is the central registry of contracts between services. Consumers push their pact files after tests pass. Providers pull contracts from the broker to verify against.

5. **Write provider verification tests** — In the provider service: pull the consumer's pact file from the broker. Run the actual provider (or a test double) against each interaction defined in the contract. The provider verifies it can satisfy every consumer's expectations. Provider tests fail if an API change breaks a consumer contract.

6. **Integrate into CI/CD** — Consumer pipeline: run consumer contract tests, publish pact file to broker. Provider pipeline: pull all consumer contracts, run verification, publish results to broker. Block provider deployment if any consumer contract verification fails. This is the safety net that prevents breaking changes.

7. **Use can-i-deploy before deployment** — Pact Broker's `can-i-deploy` CLI command checks: "Has this version of the provider been verified against the version of the consumer currently deployed in production?" Integrate as a CI gate before production deployment. Only deploy when compatibility is confirmed.

8. **Version contracts alongside services** — Tag contracts with the service version (git SHA or semver). This allows rollback: if the new provider breaks the old consumer, you can query whether the previous provider version was compatible. Version tagging enables safe rolling deployments.

9. **Handle contract evolution** — For breaking changes: deprecate the old endpoint and run old and new contracts simultaneously during a migration period. Never delete a contract until all consumers have verified they no longer depend on the old interaction. Use Pact's pending pacts feature for in-progress contract changes.

10. **Establish contract ownership** — Each consumer owns the contracts it generates. Each provider owns the verification. The contract broker owns the compatibility matrix. No team can force a breaking change without failing another team's contract tests; this creates automatic cross-team communication.

## Rules

- Contract tests verify interfaces, not behavior; they don't replace unit tests or integration tests — they replace the need for a shared live integration environment.
- Consumers must only specify the fields they actually use; over-specifying makes contracts brittle and breaks unnecessarily when the provider adds optional fields.
- Providers must never break a published contract without coordinating consumer migration; `can-i-deploy` enforces this automatically.
- Contract files are checked into the consumer's repository as well as the broker; they are code, not generated artifacts to be ignored.

## Common Mistakes

- **Testing the entire response schema** — a consumer that asserts every field of a 50-field response will break whenever the provider adds a field; only assert fields the consumer uses.
- **Running provider tests against mocks** — provider tests must run against the real provider implementation (or a close test double); mocking the provider in provider tests verifies nothing.
- **No broker** — passing contract files via file system or shared repo is unmanageable at scale; a broker is required for `can-i-deploy` and compatibility matrix management.
- **Replacing integration tests entirely** — contract tests verify interface compatibility, not that the provider's business logic is correct; integration tests are still needed for end-to-end behavior verification.

## When NOT to Use

- Monolithic applications where services share a process and in-process calls don't cross network boundaries
- External third-party APIs where you control only the consumer and cannot run provider verification
- Simple CRUD APIs with a single consumer where the overhead of contract infrastructure exceeds the benefit
