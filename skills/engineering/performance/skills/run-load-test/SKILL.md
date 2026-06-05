---
name: run-load-test
description: Use when validating system performance before a launch, during capacity planning, after significant code or infrastructure changes, or when investigating SLO latency budget breaches
source: Google SRE load testing guidance; k6/JMeter/Locust best practices; Gregg "Systems Performance" (2020) capacity chapter
tags: [performance, testing, reliability, capacity]
verified: true
---

# Run Load Test

Execute structured load tests to measure system throughput, latency, and failure characteristics under simulated production traffic before issues affect real users.

## Why This Is Best Practice

**Adopted by:** Google (load testing as part of production readiness review), Amazon, Netflix (load testing before every major launch), all engineering orgs with SLO commitments
**Impact:** Load testing before launch catches 80% of performance regressions that would otherwise become production incidents (Google SRE data); cost of fixing performance issues post-launch is 100x higher than pre-launch; capacity headroom discovered in load testing prevents 3 AM pages
**Why best:** Synthetic traffic in testing never reproduces real concurrency behavior; load testing is the only way to discover how a system behaves under real-world concurrent load before real users experience it

Sources: Gregg "Systems Performance" 2nd ed. Pearson (2020); Google "Production Readiness Review" SRE practice; k6 documentation; Fowler "Load Testing" martinfowler.com

## Steps

1. **Define load test objectives** — What do you want to learn? Maximum throughput (requests per second), latency at target throughput (SLO validation), breaking point (at what load does the system degrade), or recovery behavior (does the system recover after overload)? Different objectives require different test shapes.

2. **Model production traffic** — Analyze production access logs (last 30 days) to understand: request mix (% GET vs POST, endpoint distribution), concurrency (peak concurrent users), think time (time between user requests), and data distribution (most popular items, p50/p95/p99 request sizes). Model your load test to match this distribution.

3. **Choose a load testing tool** — k6 (JavaScript, open source, excellent for API testing, integrates with CI/CD), Apache JMeter (Java, GUI, high enterprise adoption), Locust (Python, programmatic, good for complex scenarios), Gatling (Scala, high performance, good reports). k6 is the current standard for API load testing; JMeter for complex enterprise scenarios.

4. **Write representative test scripts** — Implement the top 5-10 user journeys covering 80% of production traffic. Parameterize test data (users, IDs, search terms) to avoid cache-hit-only scenarios. Include authentication flows. Add realistic think time (0.5-2 s between requests to simulate human behavior). Hardcoded test data produces artificially optimistic cache performance.

5. **Start with a baseline test** — Run with 10% of expected peak load for 5 minutes. Verify the test itself works correctly (no script errors, responses as expected). Establish baseline metrics: throughput, p50/p95/p99 latency, error rate. This baseline is the reference for all subsequent tests.

6. **Run ramp-up load test** — Gradually increase load from 0 to 150% of expected peak over 30-60 minutes. Observe: at what load does latency begin to increase (the knee of the curve)? At what load does error rate exceed SLO threshold? At what load does the system enter saturation? Record the throughput at each threshold.

7. **Run a sustained load test** — Apply 80% of peak load for 60-90 minutes. Test for: memory leaks (growing heap over time), connection pool exhaustion, database connection limits, and log disk space exhaustion. Short tests miss time-based degradation; sustained tests reveal them.

8. **Observe all system layers during the test** — CPU, memory, disk I/O, and network utilization on application servers; database query times and connection count; cache hit rates; downstream service response times; thread pool and connection pool metrics. Load testing without observing infrastructure metrics is incomplete; the bottleneck may be anywhere in the stack.

9. **Analyze bottlenecks** — Identify the first resource to saturate under load. Common bottlenecks: CPU (insufficient compute), memory (GC pressure, leaks), database connections (pool exhaustion), I/O (disk or network bandwidth), and downstream dependencies (slow third-party APIs). Fix the identified bottleneck and re-run — there will always be another.

10. **Document results and thresholds** — Write a load test report: test scenario, traffic model, peak throughput achieved, latency at peak, breaking point, identified bottlenecks, and recommendations. Define pass/fail criteria for future tests (SLO thresholds). Automate load test execution in CI/CD with pass/fail gates for regression detection.

## Rules

- Never run load tests against production without explicit stakeholder sign-off and a defined abort procedure; a misconfigured load test can cause a real outage.
- Use production-like data volumes; testing against an empty database or a 100-row table misses all query planner and index behavior at scale.
- Load test results are valid only for the infrastructure they were run against; results from a 2-core test server don't predict production behavior on a 32-core fleet.
- Fix identified bottlenecks before accepting the load test results; "we know it breaks at X load" without a remediation plan is a known liability.

## Common Mistakes

- **Testing only the happy path** — production traffic includes authentication, error handling, and retries; happy-path-only tests miss bottlenecks in error paths that are triggered under load.
- **Ignoring warm-up period** — the first 5-10 minutes of a load test include JVM JIT compilation, cache warming, and connection pool filling; exclude this from performance metrics.
- **Caching artifacts from test data** — if all test users hit the same 10 IDs, cache hit rates are 99%; production cache hit rates may be 40%; parameterize test data to match production distribution.
- **No sustained test** — 5-minute load tests miss memory leaks, connection pool exhaustion, and log rotation issues that appear only after 30+ minutes under load.

## When NOT to Use

- Production systems during peak traffic windows (schedule tests during low-traffic hours or in dedicated load test environments)
- Systems with strict rate limits on third-party dependencies that load testing would exhaust for real users
- Prototypes not yet at production scale where performance optimization would be premature
