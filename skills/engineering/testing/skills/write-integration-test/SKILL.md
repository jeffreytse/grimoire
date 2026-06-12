---
name: write-integration-test
description: Use when writing tests that verify the interaction between two or more components, services, or external systems
source: Google Testing Blog (testing.googleblog.com); Martin Fowler "IntegrationTest" (martinfowler.com); Pact contract testing (pact.io)
tags: [integration-testing, testing, contract-tests, test-pyramid, quality]
verified: true
---

# Write Integration Test

Write tests that verify component interactions — database, API, message bus, or service-to-service — with real or realistic dependencies.

## Why This Is Best Practice

**Adopted by:** Google (testing pyramid doctrine), Pact.io (used by ITV, Dius, DiUS, many enterprises for contract tests)
**Impact:** Google's testing pyramid allocates 15% of test effort to integration tests; Pact contract tests reduced ITV's integration regression cycle from 3 days to 20 minutes.

**Why best:** Unit tests verify logic in isolation; integration tests verify that the parts work together. Without them, systems pass unit tests but fail at the boundary — database constraints, serialization mismatches, API version drift. Contract tests (Pact) are a lighter alternative for service-to-service boundaries.

## Steps

1. **Identify the integration boundary** — What two components are you testing together? (App ↔ DB, Service A ↔ Service B, App ↔ Message Queue)
2. **Choose the right integration style** — Broad integration test (real dependencies, slower) vs. narrow/contract test (mocked transport, consumer-driven). Use contract tests for external services you don't control.
3. **Set up realistic test data** — Use fixtures or factories; never depend on production data. Use database transactions or container resets to isolate tests.
4. **Spin up dependencies** — Use Docker Compose or Testcontainers to run real DB/broker instances in CI; avoid in-memory fakes for integration tests (they hide real behavior).
5. **Write the test** — Arrange real inputs → Act through the public interface → Assert on real observable output (DB state, HTTP response, emitted event). Do not assert on internal implementation.
6. **Assert side effects** — Verify state changes in all touched systems (e.g., row inserted in DB AND event published to queue).
7. **Clean up** — Tear down or roll back test data after each test to prevent interference.

## Rules

- Integration tests must pass in CI against the same dependency versions used in production.
- Never use `time.Sleep` to wait for async operations — use polling with timeout or event-driven assertions.
- One integration test per boundary scenario — do not retest unit-level logic here.
- Scope each test to a single integration point; multi-hop tests are end-to-end tests, not integration tests.

## Examples

```python
def test_create_order_persists_to_db(db_session, api_client):
    response = api_client.post("/orders", json={"item": "widget", "qty": 2})
    assert response.status_code == 201
    order = db_session.query(Order).filter_by(id=response.json()["id"]).one()
    assert order.qty == 2
    assert order.status == "pending"
```

## Common Mistakes

- **In-memory database substitutes** — SQLite behaves differently from Postgres; use real engines.
- **Shared test data** — tests interfere with each other; always isolate state.
- **Testing too much in one test** — a test that spans 5 services is an E2E test with integration-test maintenance cost.
