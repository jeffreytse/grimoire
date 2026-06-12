---
name: design-load-test
description: Use when designing or running load, stress, or performance tests to validate system behavior under traffic
source: Steve Smith "The Art of Application Performance Testing" (O'Reilly 2015); k6 documentation (k6.io); Gatling documentation (gatling.io)
tags: [load-testing, performance, k6, gatling, stress-testing, scalability]
verified: true
---

# Design Load Test

Design and execute load tests that reveal system performance characteristics, bottlenecks, and breaking points before production traffic does.

## Why This Is Best Practice

**Adopted by:** Grafana Labs (k6 as OSS standard), Netflix (Chaos + load testing), Amazon (load tests every release)
**Impact:** Amazon's 100ms latency improvement drove 1% revenue increase; load testing before peak events (Black Friday) prevents outages that cost retailers $220k/minute average downtime cost (Gartner).

**Why best:** Load testing answers questions that unit and integration tests cannot: Where does the system break? What is the maximum sustainable throughput? Where is the bottleneck? Without this data, capacity planning is guesswork.

## Steps

1. **Define objectives** — Specify target RPS (requests per second), P95/P99 latency SLOs, and error rate thresholds before writing a single test.
2. **Profile realistic traffic** — Analyze production access logs to model realistic request distribution (not just the happy path). Include think time and session patterns.
3. **Choose the tool** — k6 (JavaScript, modern, CI-friendly), Gatling (Scala/Java, excellent reporting), Locust (Python, flexible), JMeter (enterprise, complex). Prefer k6 for greenfield.
4. **Design test scenarios** — Baseline (expected load), stress (2-3× baseline), spike (sudden burst), soak (sustained load over hours to detect memory leaks).
5. **Isolate the test environment** — Use a production-like environment; never load-test production without traffic shaping. Ensure DBs have realistic data volumes.
6. **Run and instrument** — Collect server metrics (CPU, memory, DB connections, GC pauses) alongside client-side latency. Stream to Grafana or Datadog.
7. **Analyze and document** — Identify the bottleneck (DB, CPU, network, external API); report P50/P95/P99 latency, error rate, and throughput at each load level.

## Rules

- Always define success criteria before running — a test without a pass/fail threshold is just an observation.
- Load test with realistic data volumes — testing against an empty DB misses index degradation at scale.
- Include think time between requests to avoid unrealistic burst patterns.
- Automate load tests in CI for smoke-level validation; run full suites before releases.

## Examples

k6 baseline scenario:
```js
export const options = {
  stages: [
    { duration: '2m', target: 100 },  // ramp up
    { duration: '5m', target: 100 },  // sustain
    { duration: '1m', target: 0 },    // ramp down
  ],
  thresholds: { http_req_duration: ['p(95)<500'], http_req_failed: ['rate<0.01'] },
};
```

## Common Mistakes

- **Testing a single endpoint** — real traffic hits many endpoints; model realistic user journeys.
- **No baseline before optimization** — you can't prove an optimization worked without a pre-change measurement.
- **Ignoring downstream dependencies** — a load test that mocks the database hides the most common bottleneck.
