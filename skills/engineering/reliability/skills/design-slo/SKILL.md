---
name: design-slo
description: Use when defining or reviewing service reliability targets, error budgets, or SLI/SLO/SLA structures for a service
source: "Site Reliability Engineering (Google, 2016) Ch. 4; The Site Reliability Workbook (Google, 2018) Ch. 2; Alex Hidalgo 'Implementing Service Level Objectives' (2020)"
tags: [sre, slo, reliability, error-budget, sli, availability, monitoring]
verified: true
---

# Design SLO

Define measurable reliability commitments that align engineering effort with what users actually care about.

## Why This Is Best Practice

**Adopted by:** Google, Spotify, Dropbox, Netflix, and the majority of mature SRE organizations; mandated practice in the Google SRE model adopted by thousands of engineering organizations globally

**Impact:** Google's SRE book reports error budget policies cut feature/reliability conflict resolution time significantly; Dropbox publicly documented SLO adoption reducing production incidents by 30% within one year

**Why best:** Without SLOs, reliability work is driven by gut feel and loudest complaint. SLOs make reliability a quantified engineering decision: when error budget is full, ship features; when it burns fast, halt releases and fix reliability. This makes the reliability/feature trade-off objective and removes it from politics.

## Steps

1. **Identify the SLI (Service Level Indicator)** — choose a metric that represents user happiness: availability (successful requests / total requests), latency (p99 response time), or error rate; prefer request-based over time-based SLIs
2. **Set the SLO target** — pick a number the service can actually achieve today, then tighten it over time; 99.9% ("three nines") is a reasonable starting target for most web APIs
3. **Calculate the error budget** — error budget = (1 - SLO) × time period; 99.9% monthly SLO = 43.8 minutes of allowable downtime per month
4. **Define the measurement window** — rolling 28-day windows are preferred over calendar months; they avoid cliff effects at month boundaries
5. **Establish the error budget policy** — document in writing: what happens when budget is 50% consumed? 100% consumed? Who is notified? What releases are paused?
6. **Set burn rate alerts** — alert when the error budget is burning faster than baseline; a 1x burn rate depletes budget exactly on time; a 6x burn rate is an alert-worthy signal
7. **Review SLOs quarterly** — tighten targets as reliability improves; loosen if consistently in breach with no user impact signal

## Rules

- SLOs must measure what users experience, not what servers report
- Do not set an SLO you cannot measure with existing instrumentation — build instrumentation first
- 100% is never the right SLO: it eliminates the error budget and makes every incident a policy violation
- SLAs (contractual commitments to customers) must always be weaker than your internal SLOs

## Examples

**SLI:** Percentage of HTTP requests to `/api/checkout` completing in under 800ms, measured at the load balancer.
**SLO:** 99.5% of checkout requests complete under 800ms over a rolling 28-day window.
**Error budget:** 0.5% of requests × 28 days = roughly 3 hours of budget before the policy triggers a feature freeze.

## Common Mistakes

- Setting SLOs without an error budget policy: an SLO without a policy is just a dashboard metric
- Measuring availability from the server side: a 200 OK from a server that returned empty data is not a successful user request
- Setting aspirational SLOs that the service has never achieved: unreachable targets destroy team credibility

## When NOT to Use

- The service has no instrumentation and no historical request-success or latency data — SLOs require a measurable SLI; invest in observability instrumentation before defining targets.
- The service is an internal batch processing job with no real-time user impact — availability-based SLOs are designed for request-serving systems; batch jobs need throughput and lateness SLIs, which require a different design approach.
- The team has not agreed on an error budget policy — an SLO without a documented policy for what happens when budget depletes is a vanity metric; design the policy first or the SLO will be ignored in practice.
